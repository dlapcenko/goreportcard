package download

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/vcs"
	"net/url"
)

// Download takes a user-provided string that represents a remote
// Go repository, and attempts to download it in a way similar to go get.
// It is forgiving in terms of the exact path given: the path may have
// a scheme or username, which will be trimmed.
func Download(path, dest string) (root *vcs.RepoRoot, err error) {
	return download(path, dest, true)
}

func download(path, dest string, firstAttempt bool) (root *vcs.RepoRoot, err error) {
	vcs.ShowCmd = false

	path, err = Clean(path)
	if err != nil {
		return root, err
	}

	root, err = vcs.RepoRootForImportPath(path, true)
	if err != nil {
		return root, err
	}

	localDirPath := filepath.Join(dest, root.Root, "..")

	if err := os.MkdirAll(localDirPath, 0777); err != nil {
		return root, err
	}

	fullLocalPath := filepath.Join(dest, root.Root)

	if root.VCS.Name == "Git" {
		root.VCS.CreateCmd = "clone --depth 1 {repo} {dir}"
		root.VCS.TagSyncDefault = ""
	}
	var rootRepo = root.Repo

	if _, err := url.Parse(root.Repo); err != nil {
		log.Printf("WARN: could not parse root.Repo: %v", err)
	}

	if err := root.VCS.Create(fullLocalPath, rootRepo); err != nil {
		return root, err
	}

	if err := root.VCS.TagSync(fullLocalPath, ""); err != nil && firstAttempt {
		// may have been rebased; we delete the directory, then try one more time:
		log.Printf("Failed to update %q (%v), trying again...", root.Repo, err.Error())
		err = os.RemoveAll(fullLocalPath)
		if err != nil {
			log.Printf("Failed to delete directory %s", fullLocalPath)
		}
		return download(path, dest, false)
	}

	return root, err
}

// Clean trims any URL parts, like the scheme or username, that might be present
// in a user-submitted URL
func Clean(path string) (string, error) {
	importPath := trimUsername(trimScheme(path))
	root, err := vcs.RepoRootForImportPath(importPath, true)
	if err != nil {
		return "", err
	}
	if root != nil && (root.Root == "" || root.Repo == "") {
		return root.Root, errors.New("empty repo root")
	}
	return root.Root, err
}

// trimScheme removes a scheme (e.g. https://) from the URL for more
// convenient pasting from browsers.
func trimScheme(repo string) string {
	schemeSep := "://"
	schemeSepIdx := strings.Index(repo, schemeSep)
	if schemeSepIdx > -1 {
		return repo[schemeSepIdx+len(schemeSep):]
	}

	return repo
}

// trimUsername removes the username for a URL, if it is present
func trimUsername(repo string) string {
	usernameSep := "@"
	usernameSepIdx := strings.Index(repo, usernameSep)
	if usernameSepIdx > -1 {
		return repo[usernameSepIdx+len(usernameSep):]
	}

	return repo
}

// exists returns whether the given file or directory exists or not
// from http://stackoverflow.com/a/10510783
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
