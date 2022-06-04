package derefer

import (
	"os"
)

// Env returns value of environmental variable with given name or empty string
// if not set.
func Env(ref string) (string, error) {
	return os.Getenv(ref), nil
}
