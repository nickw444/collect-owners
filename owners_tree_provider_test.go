package main

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsExcluded(t *testing.T) {
	excludes := []*regexp.Regexp{
		regexp.MustCompile("^.git"),
		regexp.MustCompile("node_modules"),
		regexp.MustCompile("/.*\\.min\\.js"),
	}

	provider := &OwnersTreeProvider{
		RootPath: "",
		Excludes: excludes,
	}

	assert.True(t, provider.isExcluded(".git/objects/AAAAA"))
	assert.True(t, provider.isExcluded(".git"))
	assert.True(t, provider.isExcluded(".git/"))
	assert.True(t, provider.isExcluded("package/x/node_modules/bin"))
	assert.True(t, provider.isExcluded("package/x/build/bundle.min.js"))
}
