package derefer

import (
	"encoding/base64"
	"encoding/hex"
)

func Base64(ref string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(ref)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func Hex(ref string) (string, error) {
	b, err := hex.DecodeString(ref)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
