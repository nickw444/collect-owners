package main

import (
	"fmt"
	"strings"
)

type OwnersFileFormatter struct {
}

func (o *OwnersFileFormatter) Format(entries []*OwnersFileEntry) {
	// Get longest glob:
	longestGlob := longestEntryPathLength(entries)

	for _, entry := range entries {
		fmt.Println(rightPad(entry.Glob, longestGlob), strings.Join(entry.Owners, " "))
	}
}

func rightPad(s string, n int) string {
	if len(s) >= n {
		return s
	}
	return s + strings.Repeat(" ", n-len(s))
}

func longestEntryPathLength(entries []*OwnersFileEntry) int {
	longest := 0
	for _, entry := range entries {
		if len(entry.Glob) > longest {
			longest = len(entry.Glob)
		}
	}
	return longest
}
