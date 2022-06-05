package derefer

import (
	"errors"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// AWSParameterStore derefs a Parameter Store ARN to its value.
// Default tag for this derefer is "ssm". "pssm" tag is set to Plaintext mode
// of this derefer.
type AWSParameterStore struct {
	sess *session.Session
	// By default, the value is decrypted. Set this to true to retrieve
	// non-encrypted values.
	Plaintext bool
}

func NewAWSParameterStoreWithSession(sess *session.Session) *AWSParameterStore {
	return &AWSParameterStore{sess: sess}
}

func NewAWSParameterStore() (*AWSParameterStore, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, err
	}
	return NewAWSParameterStoreWithSession(sess), nil
}

var awsParamStoreARN = regexp.MustCompile("^(?:arn:aws:ssm:([^:]+):[^:]*:parameter)?(/.+)$")

// Deref returns the value of the parameter ref which is the ARN of the
// resource in the form (where [..] is optional):
//   [arn:aws:ssm:<region>:<account-number>:parameter]/...
func (d *AWSParameterStore) Deref(ref string) (string, error) {
	m := awsParamStoreARN.FindStringSubmatch(ref)
	if m == nil {
		return "", errors.New("invalid AWS param store key " + ref)
	}
	region, name := m[1], m[2]
	if region == "" {
		region = *d.sess.Config.Region
	}
	svc := ssm.New(d.sess, aws.NewConfig().WithRegion(region))
	ssmParam, err := svc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(name),
		WithDecryption: aws.Bool(!d.Plaintext),
	})
	if err != nil {
		return "", errors.New("cannot load AWS param '" + name + "' in region '" + region + "': " + err.Error())
	}
	return *ssmParam.Parameter.Value, nil
}
