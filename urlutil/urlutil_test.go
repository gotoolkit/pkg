package urlutil

import "testing"

var (
	gitUrls = []string{
		"git://github.com/gotoolkit/pkg",
		"https://github.com/gotoolkit/pkg.git",
		"http://github.com/gotoolkit/pkg.git",
	}
	incompleteGitUrls = []string{
		"github.com/gotoolkit/pkg",
	}
	invalidGitUrls = []string{
		"http://github.com/gotoolkit/pkg.git:#branch",
	}
)

func TestIsGitURL(t *testing.T) {
	for _, url := range gitUrls {
		if !IsGitURL(url) {
			t.Fatalf("%q should be detected as valid Git url", url)
		}
	}
	for _, url := range incompleteGitUrls {
		if !IsGitURL(url) {
			t.Fatalf("%q should be detected as valid Git url", url)
		}
	}
	for _, url := range invalidGitUrls {
		if IsGitURL(url) {
			t.Fatalf("%q should be detected as valid Git prefix", url)
		}
	}
}
