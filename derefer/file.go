package derefer

import (
	"errors"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"strings"
)

type File struct {
	ExpandHome bool
}

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
