/*
  Package starenv implements populating environmental variables from variety of
  sources.

  Usage

  Simplest way to use starenv is to import the autoload package:
    import _ "github.com/oxplot/starenv/autoload"
  The above will iterate through all environmental variables looking for
  specially formatted values which tell it where to load the final values from.
  After the above is imported, you can use the usual os.Getenv() to get the
  value of environmental variables.

  Ref Pipline

  The source of a value is defined as a pipeline of "derefer" tags followed by
  the reference for the last derefer. Here is an example of a environmental
  variable specifying to load its value from a base64 encoded file and decrypt
  it using GPG:
    GITHUB_TOKEN=*gpg:*b64:*file:~/.github_token
  Each derefer is applied in reverse, starting with "file" which loads the
  content of "~/.github_token". "b64" then decodes it and finally "gpg"
  decrypts it.

  If the value of an environmental variable starts with Loader.Star (which
  defaults to "*"), it is treated as a pipeline. Otherwise, it's treated as a
  literal value and left unchagned. In the unlikely case where a literal value
  starting with Loader.Star is needed, the following can be used:
    GLOB_PAT=*:*.terraform
  Here the blank derefer treats everything after "*:" as literal and returns
  it, thus leading to GLOB_PAT being set to "*.terraform".

  Package autoload registers a set of derefers that are included in derefer
  package with appropriate tags. To have more control over tags and the timing
  of when the loading happens, you can register each derefer manually and call
  Load() to populate the env vars.

  Any type that implements Derefer methods can be registered and used in the
  pipeline.
*/
package starenv

import (
	"errors"
	"os"
	"strings"
)

// Derefer is an interface that wraps Deref method.
//
// Deref method is called with the recursively derefed value of all subsequent
// derefers in the pipeline. ref is therefore a literal by the time it's passed
// to this method.
type Derefer interface {
	Deref(ref string) (value string, err error)
}

// DereferFunc type is an adapter to allow use of ordinary functions as derefers.
type DereferFunc func(ref string) (string, error)

func (d DereferFunc) Deref(ref string) (string, error) {
	return d(ref)
}

// Loader holds a registry of derefers which are looked up and applied to
// values of all environmental variables when Load() is called.
type Loader struct {
	Star     string
	derefers map[string]Derefer
}

func NewLoader() *Loader {
	return &Loader{
		Star:     "*",
		derefers: map[string]Derefer{},
	}
}

// Register maps a tag with a derefer. Empty tag and tags with ":" are
// unusable.
func (l *Loader) Register(tag string, d Derefer) {
	l.derefers[tag] = d
}

// Load iterates through all environmental variables and recursively derefs
// their values appropriately.
func (l *Loader) Load() error {
	for _, e := range os.Environ() {
		s := strings.SplitN(e, "=", 2)
		k, v := s[0], s[1]
		v, err := l.recurseDeref(v)
		if err != nil {
			return errors.New("failed to load env var " + k + ": " + err.Error())
		}
		os.Setenv(k, v)
	}
	return nil
}

func (l *Loader) recurseDeref(ref string) (string, error) {
	if !strings.HasPrefix(ref, l.Star) {
		return ref, nil
	}
	ref = ref[len(l.Star):]
	colIdx := strings.Index(ref, ":")
	if colIdx == -1 {
		return "", errors.New("no colon found")
	}
	if colIdx == 0 { // literal value
		return ref[1:], nil
	}
	tag := ref[:colIdx]
	d, ok := l.derefers[tag]
	if !ok {
		return "", errors.New("no registered derefer with tag " + tag)
	}
	ref, err := l.recurseDeref(ref[colIdx+1:])
	if err != nil {
		return "", err
	}
	return d.Deref(ref)
}

var (
	defaultLoader = NewLoader()
)

// Register maps a tag with a derefer on the default loader.
func Register(tag string, d Derefer) {
	defaultLoader.Register(tag, d)
}

// Load iterates through all environmental variables and recursively derefs
// their values appropriately. Derefers must fist be registered with Register()
// function.
func Load() error {
	return defaultLoader.Load()
}
