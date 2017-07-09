package main

import (
	"os"
	"path/filepath"
	"strings"
)

type FileProvider struct {
	RepoRoot string
	Excludes []string
}

// Returns an array of files relative to the repository root dir.
func (f *FileProvider) GetFiles() (files []string, err error) {
	files = []string{}

	err = filepath.Walk(f.RepoRoot, func(fPath string, info os.FileInfo, err error) (walkErr error) {
		fPath, walkErr = filepath.Rel(f.RepoRoot, fPath)
		if walkErr != nil {
			return
		}

		_, filename := filepath.Split(fPath)
		isMatch, walkErr := filepath.Match("OWNERS", filename)
		if walkErr != nil {
			return
		}

		if isMatch {
			for _, exclude := range f.Excludes {
				if strings.Contains(fPath, exclude) {
					return
				}
			}

			files = append(files, fPath)
		}

		return
	})
	return
}
