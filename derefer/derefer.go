package derefer

import "github.com/oxplot/starenv"

var NewDefault = map[string]func() (starenv.Derefer, error){
	"b64": func() (starenv.Derefer, error) {
		return starenv.DereferFunc(Base64), nil
	},
	"hex": func() (starenv.Derefer, error) {
		return starenv.DereferFunc(Hex), nil
	},
	"ssm": func() (starenv.Derefer, error) {
		return NewAWSParameterStore()
	},
	"pssm": func() (starenv.Derefer, error) {
		d, err := NewAWSParameterStore()
		if err != nil {
			return nil, err
		}
		d.Plaintext = true
		return d, nil
	},
	"file": func() (starenv.Derefer, error) {
		return &File{ExpandHome: true}, nil
	},
	"gpg": func() (starenv.Derefer, error) {
		return starenv.DereferFunc(GPG), nil
	},
}

type Lazy struct {
	New func() (starenv.Derefer, error)
	d   starenv.Derefer
}

func (l *Lazy) Deref(ref string) (string, error) {
	if l.d == nil {
		var err error
		l.d, err = l.New()
		if err != nil {
			return "", err
		}
	}
	return l.d.Deref(ref)
}
