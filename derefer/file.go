package derefer

import (
	"errors"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"strings"
)

// File derefs a file path to the content of the file. Default tag for this
// derefer is "file".
type File struct {
	// ExpandHome when set to true will replace "~/" with the current user's home
	// directory just like a standard POSIX shell would do.
	ExpandHome bool
}

// Deref treats ref as path to a file and returns the content of the file.
func (f *File) Deref(ref string) (string, error) {
	if f.ExpandHome && strings.HasPrefix(ref, "~/") {
		u, err := user.Current()
		if err != nil {
			return "", errors.New("failed to get current user: " + err.Error())
		}
		ref = filepath.Join(u.HomeDir, ref[2:])
	}
	b, err := ioutil.ReadFile(ref)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// TempFile creates a temporary file and stores the content of ref in it and
// returns its path. This is useful for storing secrets in files without
// having their content in the env var. Under most OSes, temp directory content
// is held in memory and isn't written to disk, this ensures a further layer of
// security for secrets.
// Note that the file is NOT deleted automatically.
func TempFile(v string) (string, error) {
	f, err := ioutil.TempFile("", "starenv-*")
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := f.WriteString(v); err != nil {
		return "", err
	}
	return f.Name(), nil
}
