package derefer

import (
	"errors"
	"strings"

	"github.com/zalando/go-keyring"
)

func Keyring(ref string) (string, error) {
	s := strings.SplitN(ref, "/", 2)
	if len(s) != 2 {
		return "", errors.New("keyring path must be of form service/key")
	}
	return keyring.Get(s[0], s[1])
}
