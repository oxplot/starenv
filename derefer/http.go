package derefer

import (
	"errors"
	"io"
	"net/http"
	"net/url"
)

// HTTP derefs a URL to its content. Default tag for this derefer is "https".
// Tag "http" exists for the insecure config.
type HTTP struct {
	// By default, if no scheme is provided, https is assumed. Setting this to
	// true will instead assume http.
	DefaultInsecure bool
}

// Deref treats ref as a URL and returns the response body of a GET request.
// You can leave off the scheme and take advantange of the default tag for this derefer:
//
//	*https://example.com
//
// instead of:
//
//	*https:https://example.com
func (h *HTTP) Deref(ref string) (string, error) {
	u, err := url.Parse(ref)
	if err != nil {
		return "", errors.New("cannot parse " + ref + " as url")
	}
	if u.Scheme == "" {
		if h.DefaultInsecure {
			u.Scheme = "http"
		} else {
			u.Scheme = "https"
		}
	}
	resp, err := http.Get(u.String())
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close }()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
