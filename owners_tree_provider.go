package main

import (
	"io/ioutil"
	"path"
	"regexp"
)

type DirEntry struct {
	Path       string
	DirOwners  []string
	FileOwners map[string][]string

	Parent  *DirEntry
	SubDirs []*DirEntry
}

type OwnersTreeProvider struct {
	RootPath            string
	Excludes            []*regexp.Regexp
	ownersFileProcessor *OwnersFileProcessor
}

func (o *OwnersTreeProvider) GetFileTree() (entries *DirEntry, err error) {
	return o.walkDir("", nil)
}

func (o *OwnersTreeProvider) isExcluded(filePath string) bool {
	for _, re := range o.Excludes {
		if re.MatchString(filePath) {
			return true
		}
	}

	return false
}

func (o *OwnersTreeProvider) walkDir(filePath string, parent *DirEntry) (entry *DirEntry, err error) {
	dirEntries, err := ioutil.ReadDir(path.Join(o.RootPath, filePath))
	if err != nil {
		return
	}

	entry = &DirEntry{
		Path:   filePath,
		Parent: parent,
	}

	var subDirs []string // Queue for subDirs to walk.
	hasOwners := false

	for _, ent := range dirEntries {

		if o.isExcluded(filePath) {
			continue
		}

		if ent.IsDir() {
			subDirs = append(subDirs, ent.Name())
		}

		if ent.Name() == "OWNERS" {
			hasOwners = true
		}

	}

	if !hasOwners {
		// If there is no OWNERS file for this dir, we should inherit the parent dir
		// owners.
		if parent != nil {
			entry.DirOwners = parent.DirOwners
		}
	} else {
		ownersContent, err := o.ownersFileProcessor.getOwnersForFile(path.Join(filePath, "OWNERS"))
		if err != nil {
			return nil, err
		}

		entry.DirOwners = ownersContent.dirOwners
		entry.FileOwners = ownersContent.fileOwners
	}

	// Parse subdirectories
	for _, subdir := range subDirs {
		subdirEntry, err := o.walkDir(path.Join(filePath, subdir), entry)
		if err != nil {
			return nil, err
		}
		entry.SubDirs = append(entry.SubDirs, subdirEntry)
	}

	return
}
