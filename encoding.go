package starenv

import (
	"encoding/base64"
	"encoding/hex"
)

// Base64 decodes base64 encoded ref and returns it. Default tag for this
// derefer is "b64".
func Base64(ref string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(ref)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Hex decodes hex encoded ref and returns it. Default tag for this derefer is
// "hex".
func Hex(ref string) (string, error) {
	b, err := hex.DecodeString(ref)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
