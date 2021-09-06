package starenv

import (
	"errors"
	"os"
	"strings"
)

type Derefer interface {
	Deref(ref string) (value string, err error)
}

type DereferFunc func(ref string) (string, error)

func (d DereferFunc) Deref(ref string) (string, error) {
	return d(ref)
}

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

func (l *Loader) Register(tag string, d Derefer) {
	l.derefers[tag] = d
}

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

func Register(tag string, d Derefer) {
	defaultLoader.Register(tag, d)
}

func Load() error {
	return defaultLoader.Load()
}
