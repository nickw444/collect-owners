package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var formatter = &OwnersFileFormatter{}

func TestBuildRulesNoOwners(t *testing.T) {
	entry := buildEntry("", []string{}, map[string][]string{})
	rules := formatter.buildRules(entry)
	assert.Len(t, rules, 0)
}

func TestBuildRulesDirOwnersRoot(t *testing.T) {
	owners := []string{"foo@example.com", "bar@example.com"}
	entry := buildEntry("", owners, map[string][]string{})
	rules := formatter.buildRules(entry)
	assert.Len(t, rules, 1)
	assert.Equal(t, "*", rules[0].Glob)
	assert.Equal(t, owners, rules[0].Owners)
}

func TestBuildRulesDirOwnersSubDir(t *testing.T) {
	owners := []string{"foo@example.com", "bar@example.com"}
	entry := buildEntry("foo/bar/baz", owners, map[string][]string{})
	rules := formatter.buildRules(entry)
	assert.Len(t, rules, 1)
	assert.Equal(t, "foo/bar/baz/*", rules[0].Glob)
	assert.Equal(t, owners, rules[0].Owners)
}

func TestBuildRulesFileOwnersRoot(t *testing.T) {
	owners := []string{"foo@example.com", "bar@example.com"}
	entry := buildEntry("", []string{}, map[string][]string{
		"*.es6": owners,
	})
	rules := formatter.buildRules(entry)
	assert.Len(t, rules, 1)
	assert.Equal(t, "*.es6", rules[0].Glob)
	assert.Equal(t, owners, rules[0].Owners)
}

func TestBuildRulesFileOwnersSubDir(t *testing.T) {
	owners := []string{"foo@example.com", "bar@example.com"}
	entry := buildEntry("foo/bar/baz", []string{}, map[string][]string{
		"*.es6": owners,
	})
	rules := formatter.buildRules(entry)
	assert.Len(t, rules, 1)
	assert.Equal(t, "foo/bar/baz/*.es6", rules[0].Glob)
	assert.Equal(t, owners, rules[0].Owners)
}

func TestBuildRulesRecursion(t *testing.T) {
	topEntry := buildEntry("foo/bar", []string{"dir@example.com"}, map[string][]string{})
	topEntry.SubDirs = []*DirEntry{
		buildEntry("foo/bar/baz", []string{"subdir@example.com"}, map[string][]string{}),
	}

	rules := formatter.buildRules(topEntry)
	assert.Len(t, rules, 2)
	assert.Equal(t, "foo/bar/*", rules[0].Glob)
	assert.Equal(t, []string{"dir@example.com"}, rules[0].Owners)
	assert.Equal(t, "foo/bar/baz/*", rules[1].Glob)
	assert.Equal(t, []string{"subdir@example.com"}, rules[1].Owners)
}

func buildEntry(path string, dirOwners []string, fileOwners map[string][]string) *DirEntry {
	return &DirEntry{
		Path:       path,
		DirOwners:  dirOwners,
		FileOwners: fileOwners,

		Parent:  nil,
		SubDirs: []*DirEntry{},
	}
}
