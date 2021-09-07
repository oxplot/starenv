package derefer

import (
	"errors"
	"os/exec"
	"strings"
)

// GPG takes encrypted content ref, decrypts it and returns it. It calls on
// external gpg command to do this.
// Default tag for this derefer is "gpg".
func GPG(ref string) (string, error) {
	c := exec.Command("gpg", "--decrypt")
	c.Stdin = strings.NewReader(ref)
	w := &strings.Builder{}
	c.Stdout = w
	if err := c.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return "", errors.New("gpg call failed: " + string(exitErr.Stderr))
		}
		return "", errors.New("gpg call failed: " + err.Error())
	}
	return w.String(), nil
}
