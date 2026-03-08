package main

import "strings"

// shortHash returns a safe short version of a commit hash
func shortHash(hash string) string {
	hash = strings.TrimSpace(hash)
	if len(hash) > 7 {
		return hash[:7]
	}
	return hash
}
