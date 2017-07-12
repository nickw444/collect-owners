package main

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type ownersFileContent struct {
	path       string
	dirOwners  []string
	fileOwners map[string][]string // Email addresses of owners.
}

type OwnersFileProcessor struct {
	RepoRoot string
}

var (
	perFileRegex  = regexp.MustCompile("per-file (.*)=(.*)")
	filePathRegex = regexp.MustCompile("file:/(.*)")
)

// Reads a line and determines if it's a file:// or an owner email and
// handles accordingly.
func (o *OwnersFileProcessor) getOwners(fileOrOwner string) ([]string, error) {
	filePathMatch := filePathRegex.FindAllStringSubmatch(fileOrOwner, -1)
	if filePathMatch != nil {
		next_file := filePathMatch[0][1]
		if strings.HasPrefix(next_file, "/") {
			// File is relative to the repo root. Strip the leading slash
			next_file = next_file[1:]
		} else {
			// Need to resolve the path
			dir, _ := filepath.Split(fileOrOwner)
			next_file = path.Join(dir, next_file)
		}

		content, err := o.getOwnersForFile(next_file)
		return content.dirOwners, err

	}
	return []string{fileOrOwner}, nil
}

func (o *OwnersFileProcessor) getOwnersForFile(filename string) (content *ownersFileContent, err error) {
	bytes, err := ioutil.ReadFile(path.Join(o.RepoRoot, filename))
	if err != nil {
		return
	}
	fileContent := strings.Split(string(bytes), "\n")

	dir, _ := filepath.Split(filename)
	content = &ownersFileContent{}
	content.path = dir
	content.fileOwners = make(map[string][]string)

	for _, line := range fileContent {
		line := strings.Trim(line, " \n")
		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, "set") {
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		perFileMatches := perFileRegex.FindAllStringSubmatch(line, -1)
		if perFileMatches != nil {
			fileGlob := perFileMatches[0][1]
			owners, err := o.getOwners(perFileMatches[0][2])
			if err != nil {
				return nil, err
			}

			for _, owner := range owners {
				content.fileOwners[fileGlob] = append(content.fileOwners[fileGlob], owner)
			}

			continue
		}

		owners, err := o.getOwners(line)
		if err != nil {
			return nil, err
		}

		for _, owner := range owners {
			content.dirOwners = append(content.dirOwners, owner)
		}
	}
	return
}
