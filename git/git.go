package git

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/gotoolkit/pkg/urlutil"
)

type gitRepo struct {
	remote string
	ref    string
	subdir string
}

// Clone clones a repository into a newly created directory which will be under "git-repos"
func Clone(remoteURL string) (string, error) {

	repo, err := parseRemoteURL(remoteURL)
	if err != nil {
		return "", err
	}
	fetch := fetchArgs(repo.remote, repo.ref)

	root, err := ioutil.TempDir("", "git-repos")
	if err != nil {
		return "", nil
	}

	if out, err := gitWithinDir(root, "init"); err != nil {
		return "", errors.Wrapf(err, "failed to init repo at %s: %s", root, out)
	}

	if out, err := gitWithinDir(root, "remote", "add", "origin", repo.remote); err != nil {
		return "", errors.Wrapf(err, "failed add origin repo at %s: %s", repo.remote, out)
	}

	if output, err := gitWithinDir(root, fetch...); err != nil {
		return "", errors.Wrapf(err, "error fetching: %s", output)
	}

	return checkoutGit(root, repo.ref, repo.subdir)
}

func parseRemoteURL(remoteURL string) (gitRepo, error) {
	repo := gitRepo{}
	if !isGitTransport(remoteURL) {
		remoteURL = "https://" + remoteURL
	}
	var fragment string
	if strings.HasPrefix(remoteURL, "git@") {
		parts := strings.SplitN(remoteURL, "#", 2)
		repo.remote = parts[0]
		if len(parts) == 2 {
			fragment = parts[1]
		}
		repo.ref, repo.subdir = getRefAndSubdir(fragment)
	} else {
		u, err := url.Parse(remoteURL)
		if err != nil {
			return repo, err
		}
		repo.ref, repo.subdir = getRefAndSubdir(u.Fragment)
		u.Fragment = ""
		repo.remote = u.String()
	}
	return repo, nil
}

func getRefAndSubdir(fragment string) (ref string, subdir string) {
	refAndDir := strings.SplitN(fragment, ":", 2)
	ref = "master"
	if len(refAndDir[0]) != 0 {
		ref = refAndDir[0]
	}
	if len(refAndDir) > 1 && len(refAndDir[1]) != 0 {
		subdir = refAndDir[1]
	}
	return
}

func fetchArgs(remoteURL string, ref string) []string {
	args := []string{"fetch", "--recurse-submodules=yes"}
	// TODO
	if supportsShallowClone(remoteURL) {
		args = append(args, "--depth", "1")
	}
	return append(args, "origin", ref)
}

// Check if a given git URL supports a shallow git clone,
// i.e. it is a non-HTTP server or a smart HTTP server.
func supportsShallowClone(remoteURL string) bool {
	if urlutil.IsURL(remoteURL) {
		// Check if the HTTP server is smart

		// Smart servers must correctly respond to a query for the git-upload-pack service
		serviceURL := remoteURL + "/info/refs?service=git-upload-pack"

		// Try a HEAD request and fallback to a Get request on error
		res, err := http.Head(serviceURL)
		if err != nil || res.StatusCode != http.StatusOK {
			res, err = http.Get(serviceURL)
			if err == nil {
				res.Body.Close()
			}
			if err != nil || res.StatusCode != http.StatusOK {
				// request failed
				return false
			}
		}

		if res.Header.Get("Content-Type") != "application/x-git-upload-pack-advertisement" {
			// Fallback, not a smart server
			return false
		}
		return true
	}
	// Non-HTTP protocols always support shallow clones
	return true
}

func checkoutGit(root, ref, subdir string) (string, error) {
	if output, err := gitWithinDir(root, "checkout", ref); err != nil {
		if _, err2 := gitWithinDir(root, "checkout", "FETCH_HEAD"); err2 != nil {
			return "", errors.Wrapf(err, "error checking out %s: %s", ref, output)
		}
	}
	return root, nil
}
func gitWithinDir(dir string, args ...string) ([]byte, error) {
	a := []string{"--work-tree", dir, "--git-dir", filepath.Join(dir, ".git")}
	return git(append(a, args...)...)
}

func git(args ...string) ([]byte, error) {
	return exec.Command("git", args...).CombinedOutput()
}

func isGitTransport(str string) bool {
	return urlutil.IsURL(str) || urlutil.IsGitURL(str)
}
