package urlutil

import (
	"regexp"
	"strings"
)

var (
	validPrefixes = map[string][]string{
		"url": {"http://", "https://"},
		"git": {"git://", "github.com/", "git@"},
	}
	urlPathWithSuffix = regexp.MustCompile(".git(?:#.+)?$")
)

// IsURL returns true if the provided str is an HTTP(s) URL
func IsURL(str string) bool {
	return checkURL(str, "url")
}

// IsGitURL returns true if the provided str is a git repository URL
func IsGitURL(str string) bool {
	if IsURL(str) && urlPathWithSuffix.MatchString(str) {
		return true
	}
	return checkURL(str, "git")
}

func checkURL(str, kind string) bool {
	for _, prefix := range validPrefixes[kind] {
		if strings.HasPrefix(str, prefix) {
			return true
		}
	}
	return false
}
