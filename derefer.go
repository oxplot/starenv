package starenv

// NewDefault is a mapping of default tags to derefer creator functions that
// use sensible default config.
// When using autoload package, these derefers are automatically loaded with
// default tags.
var NewDefault = map[string]func() (Derefer, error){
	"env": func() (Derefer, error) {
		return DereferFunc(Env), nil
	},
	"b64": func() (Derefer, error) {
		return DereferFunc(Base64), nil
	},
	"hex": func() (Derefer, error) {
		return DereferFunc(Hex), nil
	},
	"ssm": func() (Derefer, error) {
		return NewAWSParameterStore()
	},
	"pssm": func() (Derefer, error) {
		d, err := NewAWSParameterStore()
		if err != nil {
			return nil, err
		}
		d.Plaintext = true
		return d, nil
	},
	"s3": func() (Derefer, error) {
		return NewS3()
	},
	"file": func() (Derefer, error) {
		return &File{ExpandHome: true}, nil
	},
	"gpg": func() (Derefer, error) {
		return DereferFunc(GPG), nil
	},
	"keyring": func() (Derefer, error) {
		return DereferFunc(Keyring), nil
	},
	"https": func() (Derefer, error) {
		return &HTTP{}, nil
	},
	"http": func() (Derefer, error) {
		return &HTTP{DefaultInsecure: true}, nil
	},
	"tmpfile": func() (Derefer, error) {
		return DereferFunc(TempFile), nil
	},
	"gz": func() (Derefer, error) {
		return DereferFunc(Gzip), nil
	},
	"bz2": func() (Derefer, error) {
		return DereferFunc(Bzip2), nil
	},
	"flate": func() (Derefer, error) {
		return DereferFunc(Flate), nil
	},
}

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

// Deref calls f(ref).
func (d DereferFunc) Deref(ref string) (string, error) {
	return d(ref)
}

func literalDerefer(ref string) (string, error) {
	return ref, nil
}

// LazyDerefer is a derefer that encapsulates a derefer creator function and delays
// its call until the first deref call.
type LazyDerefer struct {
	New func() (Derefer, error)
	d   Derefer
}

// Deref calls the underlying derefer's Deref method.
func (l *LazyDerefer) Deref(ref string) (string, error) {
	if l.d == nil {
		var err error
		l.d, err = l.New()
		if err != nil {
			return "", err
		}
	}
	return l.d.Deref(ref)
}
