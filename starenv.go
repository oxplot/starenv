/*
Package starenv implements populating environmental variables from variety of
sources.

# Usage

Simplest way to use starenv is to import the autoload package:

	import _ "github.com/oxplot/starenv/autoload"

The above will iterate through all environmental variables looking for
specially formatted values which tell it where to load the final values from.
After the above is imported, you can use the usual os.Getenv() to get the
value of environmental variables.

# Ref Pipline

The source of a value is defined as a pipeline of "derefer" tags followed by
the reference for the last derefer. Here is an example of a environmental
variable specifying to load its value from a base64 encoded file and decrypt
it using GPG:

	GITHUB_TOKEN=*gpg*b64*file:~/.github_token

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
Load() to populate the env vars. When using the autoload package, .env file is
also read and loaded. To disable this, set DOTENV_ENABLED=0.

Any type that implements Derefer methods can be registered and used in the
pipeline.
*/
package starenv

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Loader holds a registry of derefers which are looked up and applied to
// values of all environmental variables when Load() is called.
type Loader struct {
	// Star is the prefix and separator of derefer tags which defaults to "*".
	Star     string
	derefers map[string]Derefer
}

// NewLoader returns a new loader with empty "" tag mapped to to a passthrough
// derefer. This is needed to allow for environmental variable values which
// start with "*":
//
//	GLOB_PAT=*:*.terraform
//
// The above will resolve to "*.terraform".
func NewLoader() *Loader {
	return &Loader{
		Star:     "*",
		derefers: map[string]Derefer{"": DereferFunc(literalDerefer)},
	}
}

// Register maps a tag with a derefer. Tags that include ":" are unusable.
func (l *Loader) Register(tag string, d Derefer) {
	l.derefers[tag] = d
}

// Load iterates through all environmental variables and recursively derefs
// their values appropriately. If STARENV_ENABLED env var is set to 0, it does
// nothing. Error in loading one variable does not stop the process. All errors
// are collected and returned as a slice. The errored out variables are set to
// empty string.
func (l *Loader) Load() []error {
	if os.Getenv("STARENV_ENABLED") == "0" {
		return nil
	}
	errs := []error{}
	for _, e := range os.Environ() {
		s := strings.SplitN(e, "=", 2)
		k, v := s[0], s[1]
		v, err := l.load(v)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to load env var %s: %w", k, err))
			os.Setenv(k, "")
			continue
		}
		os.Setenv(k, v)
	}
	return errs
}

func (l *Loader) load(ref string) (string, error) {
	if !strings.HasPrefix(ref, l.Star) {
		return ref, nil
	}
	ref = ref[len(l.Star):]
	colIdx := strings.Index(ref, ":")
	if colIdx == -1 {
		return "", errors.New("no colon found")
	}
	tags := strings.Split(ref[:colIdx], l.Star)
	ref = ref[colIdx+1:]
	for i := len(tags) - 1; i >= 0; i-- {
		d, ok := l.derefers[tags[i]]
		if !ok {
			return "", errors.New("no registered derefer with tag " + tags[i])
		}
		var err error
		ref, err = d.Deref(ref)
		if err != nil {
			return "", err
		}
	}
	return ref, nil
}

var (
	// DefaultLoader is the default loader which has all the DefaultDerefers
	// registered with it.
	DefaultLoader = NewLoader()
)

func init() {
	for t, n := range DefaultDerefers {
		DefaultLoader.Register(t, &LazyDerefer{New: n})
	}
}
