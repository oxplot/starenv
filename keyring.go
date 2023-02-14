package starenv

import (
	"errors"
	"strings"

	"github.com/zalando/go-keyring"
)

// Keyring uses system provided secret storage to retrive the secret stored for
// ref and returns it. ref must be formatted as "service/secret" where
// "service" is a grouping under which the secret "secret" is stored. The exact
// definition of "service" depends on the operating system. Refer to
// github.com/zalando/go-keyring module for more info.
//
// On linux, you can store a secret via GNOME's SecretService using the command line:
//
//	secret-tool store --label my_app_secret service my_app username user123
//
// You can then pass "my_app/user123" as ref to Keyring() method to retrieve it.
//
// Default tag for this derefer is "keyring".
func Keyring(ref string) (string, error) {
	s := strings.SplitN(ref, "/", 2)
	if len(s) != 2 {
		return "", errors.New("keyring path must be of form service/key")
	}
	return keyring.Get(s[0], s[1])
}
