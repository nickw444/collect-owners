package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var fileProcessor *OwnersFileProcessor

func TestMain(m *testing.M) {
	fileProcessor = &OwnersFileProcessor{
		RepoRoot: "fixtures/owners_file_processor_test",
	}

	os.Exit(m.Run())
}

func TestGetOwnersForBasicFile(t *testing.T) {
	content, err := fileProcessor.getOwnersForFile("OWNERS")
	assert.NoError(t, err)
	assert.Equal(t, []string{"bar@baz.com", "foo@baz.com"}, content.dirOwners)
	assert.Equal(t, map[string][]string{}, content.fileOwners)
}

func TestGetOwnersForFileWithPerFile(t *testing.T) {
	content, err := fileProcessor.getOwnersForFile("OWNERS-perfile")
	assert.NoError(t, err)
	assert.Equal(t, []string{"bar@baz.com"}, content.dirOwners)
	assert.Equal(t, map[string][]string{
		"*.js":     []string{"js@example.com"},
		"file.es6": []string{"file@example.com"},
	}, content.fileOwners)
}

func TestGetOwnersForFileWithIgnoredSyntax(t *testing.T) {
	content, err := fileProcessor.getOwnersForFile("OWNERS-syntax")
	assert.NoError(t, err)
	assert.Equal(t, []string{}, content.dirOwners)
	assert.Equal(t, map[string][]string{}, content.fileOwners)
}

func TestGetOwnersForFileWithPerFileRelImport(t *testing.T) {
	content, err := fileProcessor.getOwnersForFile("OWNERS-withrelimport")
	assert.NoError(t, err)
	assert.Equal(t, []string{"rel@owners.com"}, content.dirOwners)
	assert.Equal(t, map[string][]string{}, content.fileOwners)
}

func TestGetOwnersForFileWithPerFileAbsImport(t *testing.T) {
	content, err := fileProcessor.getOwnersForFile("OWNERS-withabsimport")
	assert.NoError(t, err)
	assert.Equal(t, []string{"abs@owners.com"}, content.dirOwners)
	assert.Equal(t, map[string][]string{}, content.fileOwners)
}

func TestGetOwnersOwnerName(t *testing.T) {
	owner, err := fileProcessor.getOwners("", "foo@bar.com")
	assert.NoError(t, err)
	assert.Equal(t, []string{"foo@bar.com"}, owner)
}

func TestGetOwnersOwnerFile(t *testing.T) {
	owner, err := fileProcessor.getOwners("", "file://OWNERS")
	assert.NoError(t, err)
	assert.Equal(t, []string{"bar@baz.com", "foo@baz.com"}, owner)
}

func TestGetFileImportPathRel(t *testing.T) {
	newPath := fileProcessor.getFileImportPath("bar", "foo/OWNERS")
	assert.Equal(t, "bar/foo/OWNERS", newPath)
}

func TestGetFileImportPathAbs(t *testing.T) {
	newPath := fileProcessor.getFileImportPath("bar", "/OWNERS")
	assert.Equal(t, "OWNERS", newPath)
}
