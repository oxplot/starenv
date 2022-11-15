package derefer

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// AWSParameterStore derefs a Parameter Store ARN to its value.
// Default tag for this derefer is "ssm". "pssm" tag is set to Plaintext mode
// of this derefer.
type AWSParameterStore struct {
	cfg aws.Config
	// By default, the value is decrypted. Set this to true to retrieve
	// non-encrypted values.
	Plaintext bool
}

// NewAWSParameterStore creates a new AWSParameterStore derefer with default
// configuration.
func NewAWSParameterStore() (*AWSParameterStore, error) {
	c, err := awsconfig.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	return NewAWSParameterStoreWithConfig(c), nil
}

// NewAWSParameterStoreWithConfig creates a new AWSParameterStore derefer with
// the given configuration.
func NewAWSParameterStoreWithConfig(c aws.Config) *AWSParameterStore {
	return &AWSParameterStore{cfg: c}
}

var awsParamStoreARN = regexp.MustCompile("^(?:arn:aws:ssm:([^:]+):[^:]*:parameter)?(/.+)$")

// Deref returns the value of the parameter ref which is the ARN of the
// resource in the form (where [..] is optional):
//
//	[arn:aws:ssm:<region>:<account-number>:parameter]/...
func (d *AWSParameterStore) Deref(ref string) (string, error) {
	m := awsParamStoreARN.FindStringSubmatch(ref)
	if m == nil {
		return "", errors.New("invalid AWS param store key " + ref)
	}
	region, name := m[1], m[2]
	cfg := d.cfg
	if region == "" {
		cfg = cfg.Copy()
		cfg.Region = region
	}
	cl := ssm.NewFromConfig(cfg)
	ssmParam, err := cl.GetParameter(context.Background(), &ssm.GetParameterInput{
		Name:           aws.String(name),
		WithDecryption: aws.Bool(!d.Plaintext),
	})
	if err != nil {
		return "", errors.New("cannot load AWS param '" + name + "' in region '" + region + "': " + err.Error())
	}
	return *ssmParam.Parameter.Value, nil
}

// S3 derefs an S3 bucket and object path to its content.
// Default tag for this derefer is "s3".
type S3 struct {
	cfg aws.Config
}

// NewS3 creates a new S3 derefer with default configuration.
func NewS3() (*S3, error) {
	c, err := awsconfig.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	return NewS3WithConfig(c), nil
}

// NewS3WithConfig creates a new S3 derefer with the given configuration.
func NewS3WithConfig(c aws.Config) *S3 {
	return &S3{cfg: c}
}

var s3PathPat = regexp.MustCompile("^([^/]+)/([^@]+)(?:@(.+))?$")

// Deref returns the content of the S3 object given the bucket and path, and
// optionally the object version. If version is omitted, latest version of the
// object is used.
//
//	bucket/path/to/object[@version]
func (d *S3) Deref(ref string) (string, error) {
	m := s3PathPat.FindStringSubmatch(ref)
	if m == nil {
		return "", errors.New("invalid S3 object path: " + ref)
	}
	var ver *string
	if m[3] != "" {
		ver = aws.String(m[3])
	}
	cl := s3.NewFromConfig(d.cfg)
	o, err := cl.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket:    aws.String(m[1]),
		Key:       aws.String(m[2]),
		VersionId: ver,
	})
	if err != nil {
		return "", fmt.Errorf("cannot load S3 object '%s': %s", ref, err)
	}
	defer o.Body.Close()
	buf, err := ioutil.ReadAll(o.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read S3 object '%s' content: %s", ref, err)
	}
	return string(buf), nil
}
