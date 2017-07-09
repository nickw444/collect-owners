package main

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/nickw444/collect-owners/usernamedb"
)

type ownersFileContent struct {
	path       string
	dirOwners  []string
	fileOwners map[string][]string // Email addresses of owners.
}

// OwnersFileEntry is a mapping of a glob'd path to a list of owner GH usernames
type OwnersFileEntry struct {
	Glob   string
	Owners []string // Usernames that have been processed by UsernameDB.
}

// OwnersWalker walks given owners files to produce an array of OwnersWalkerEntries.
type OwnersWalker struct {
	UsernameDB    usernamedb.UsernameDB
	RepoRoot      string
	FileProvider  *FileProvider
	FileProcessor *OwnersFileProcessor

	fileContent []*ownersFileContent
}

func (o *OwnersWalker) Walk() (err error) {
	o.fileContent = []*ownersFileContent{}
	ownerFiles, err := o.FileProvider.GetFiles()
	if err != nil {
		return
	}

	for _, file := range ownerFiles {
		content, err := o.FileProcessor.getOwnersForFile(file)
		if err != nil {
			return err
		}
		o.fileContent = append(o.fileContent, content)
	}

	return
}

func (o *OwnersWalker) CollectEntries() (entries []*OwnersFileEntry) {

	for _, content := range o.fileContent {

		// Add an entry for the Directory:
		dirOwnersUsernames := o.UsernameDB.ToUsernames(content.dirOwners)
		if len(dirOwnersUsernames) > 0 {
			dirGlob := content.path + "*"
			entry := &OwnersFileEntry{
				Glob:   dirGlob,
				Owners: dirOwnersUsernames,
			}
			entries = append(entries, entry)
		}

		// Add an entry for each per-file entry:
		for fileGlob, fileOwners := range content.fileOwners {
			fileOwnersUsernames := o.UsernameDB.ToUsernames(fileOwners)
			if len(fileOwnersUsernames) > 0 {
				glob := path.Join(content.path, fileGlob)
				entry := &OwnersFileEntry{
					Glob:   glob,
					Owners: fileOwnersUsernames,
				}
				entries = append(entries, entry)
			}
		}

	}

	return
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
