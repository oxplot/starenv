package derefer

import (
	"encoding/base64"
	"encoding/hex"
)

// Base64 decodes base64 encoded ref and returns it.
func Base64(ref string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(ref)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Hex decodes hex encoded ref and returns it.
func Hex(ref string) (string, error) {
	b, err := hex.DecodeString(ref)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
