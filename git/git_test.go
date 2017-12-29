package git

import "testing"

func TestClone(t *testing.T) {
	gitUrl := "https://github.com/gotoolkit/pkg.git"

	_, err := Clone(gitUrl)
	if err != nil {
		t.Fatalf("%q should be clone", gitUrl)
	}

}
