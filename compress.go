package starenv

import (
	"bytes"
	"compress/bzip2"
	"compress/flate"
	"compress/gzip"
	"io/ioutil"
)

// Gzip decompresses gzipped compressed ref and returns it. Default tag for this
// derefer is "gz".
func Gzip(ref string) (string, error) {
	b := bytes.NewBufferString(ref)
	r, err := gzip.NewReader(b)
	if err != nil {
		return "", err
	}
	out, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// Bzip2 decompresses bzipped2 compressed ref and returns it. Default tag for
// this derefer is "bz2".
func Bzip2(ref string) (string, error) {
	b := bytes.NewBufferString(ref)
	out, err := ioutil.ReadAll(bzip2.NewReader(b))
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// Flate decompresses flate compressed ref and returns it. Default tag for this
// derefer is "flate".
func Flate(ref string) (string, error) {
	b := bytes.NewBufferString(ref)
	out, err := ioutil.ReadAll(flate.NewReader(b))
	if err != nil {
		return "", err
	}
	return string(out), nil
}
