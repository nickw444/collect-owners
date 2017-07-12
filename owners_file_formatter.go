package main

import (
	"fmt"
	"path"
	"strings"

	"github.com/nickw444/collect-owners/usernamedb"
)

type OwnersFileFormatter struct {
	UsernameDB *usernamedb.UsernameDB
}

type OwnershipRule struct {
	Glob   string
	Owners []string
}

func (o *OwnersFileFormatter) Format(rootEntry *DirEntry) {
	// Get longest glob:
	rules := o.buildRules(rootEntry)
	longestGlob := longestGlobLength(rules)

	for _, rule := range rules {
		resolvedOwners := o.UsernameDB.ToUsernames(rule.Owners)
		if len(resolvedOwners) > 0 {
			fmt.Println(rightPad(rule.Glob, longestGlob), strings.Join(resolvedOwners, " "))
		}
	}
}

func (o *OwnersFileFormatter) buildRules(entry *DirEntry) (rules []*OwnershipRule) {

	rootGlob := entry.Path
	if len(rootGlob) > 0 {
		rootGlob += "/*"
	} else {
		rootGlob += "*"
	}

	rootRule := &OwnershipRule{
		Glob:   rootGlob,
		Owners: entry.DirOwners,
	}

	rules = append(rules, rootRule)

	for fileGlob, fileOwners := range entry.FileOwners {
		rule := &OwnershipRule{
			Glob:   path.Join(entry.Path, fileGlob),
			Owners: fileOwners,
		}
		rules = append(rules, rule)
	}

	for _, subDirEntry := range entry.SubDirs {
		subDirRules := o.buildRules(subDirEntry)
		for _, rule := range subDirRules {
			rules = append(rules, rule)
		}
	}

	return
}

func rightPad(s string, n int) string {
	if len(s) >= n {
		return s
	}
	return s + strings.Repeat(" ", n-len(s))
}

func longestGlobLength(rules []*OwnershipRule) int {
	longest := 0
	for _, rule := range rules {
		if len(rule.Glob) > longest {
			longest = len(rule.Glob)
		}
	}
	return longest
}
