package derefer

import "github.com/oxplot/starenv"

var NoConfigInit = map[string]func() (starenv.Derefer, error){
	"b64": func() (starenv.Derefer, error) {
		return starenv.DereferFunc(Base64), nil
	},
	"hex": func() (starenv.Derefer, error) {
		return starenv.DereferFunc(Hex), nil
	},
	"ssm": func() (starenv.Derefer, error) {
		return NewAWSParameterStore()
	},
	"essm": func() (starenv.Derefer, error) {
		d, err := NewAWSParameterStore()
		if err != nil {
			return nil, err
		}
		d.Decrypt = true
		return d, nil
	},
	"file": func() (starenv.Derefer, error) {
		return &File{ExpandHome: true}, nil
	},
	"gpg": func() (starenv.Derefer, error) {
		return starenv.DereferFunc(GPG), nil
	},
}
