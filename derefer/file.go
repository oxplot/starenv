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
