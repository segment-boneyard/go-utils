package utils

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/segmentio/go-utils/session"
)

// An S3Copier copies files from one S3 profile to another.
type S3Copier struct {
	FromProfile string
	ToProfile   string
	Timeout     time.Duration
}

// S3CopySettings are used to copy a single object from one bucket to another.
type S3CopySettings struct {
	// Key to retrieve the object.
	FromKey    string
	FromBucket string
	// Key to put the object to. If empty, uses settings.FromKey.
	ToKey    string
	ToBucket string
	// Region to get the object from. If empty, uses the region from the
	// FromProfile.
	FromRegion string
	// Region to put the object to. If empty, uses the region from the
	// ToProfile.
	ToRegion string
	// ACL to use for putting the object.
	ACL string
}

// DefaultS3Timeout is the default timeout for an S3 copy to complete.
var DefaultS3Timeout = 10 * time.Minute

// NewS3Copier creates an S3Copier for copying objects between profiles.
//
// More configuration options can be set by directly initializing an S3Copier.
func NewS3Copier(fromProfile, toProfile string) *S3Copier {
	return &S3Copier{
		FromProfile: fromProfile,
		ToProfile:   toProfile,
		Timeout:     DefaultS3Timeout,
	}
}

func wrap(err error, reason string) error {
	return fmt.Errorf("Error %s: %v", reason, err)
}

// Copy copies the object from settings.FromBucket to settings.ToBucket using
// the given credentials, and returns an error if any exist.
//
// Copy will ask aws-vault to return valid credentials for s.FromProfile and
// s.ToProfile.
func (s *S3Copier) Copy(settings *S3CopySettings) error {
	c := &session.CredentialGetter{
		Region: settings.FromRegion,
	}
	sess, err := c.Get(s.FromProfile, nil)
	if err != nil {
		return wrap(err, "getting credentials")
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.Timeout)
	defer cancel()
	svc := s3.New(sess)
	inp := new(s3.GetObjectInput)
	inp.SetBucket(settings.FromBucket)
	inp.SetKey(settings.FromKey)
	req, out := svc.GetObjectRequest(inp)
	req.HTTPRequest = req.HTTPRequest.WithContext(ctx)
	if err := req.Send(); err != nil {
		return wrap(err, "getting object")
	}
	fmt.Fprintf(os.Stderr, "downloading s3://%s/%s\n", settings.FromBucket, settings.FromKey)
	// Not happy about this, but the PutObjectInput.Body needs to be
	// a ReadSeeker, so we need to buffer everything.
	body, err := ioutil.ReadAll(out.Body)
	if err != nil {
		return wrap(err, "downloading object")
	}
	if err := out.Body.Close(); err != nil {
		return wrap(err, "closing get request body")
	}

	putc := &session.CredentialGetter{
		Region: settings.ToRegion,
	}
	sess2, err := putc.Get(s.ToProfile, nil)
	if err != nil {
		return wrap(err, "getting credentials for put request")
	}
	svc2 := s3.New(sess2)
	inp2 := new(s3.PutObjectInput)
	inp2.SetBucket(settings.ToBucket)
	if settings.ToKey == "" {
		inp2.SetKey(settings.FromKey)
	} else {
		inp2.SetKey(settings.ToKey)
	}
	inp2.SetBody(bytes.NewReader(body))
	inp2.SetACL(settings.ACL)
	req2, _ := svc2.PutObjectRequest(inp2)
	req2.HTTPRequest = req2.HTTPRequest.WithContext(ctx)
	fmt.Fprintf(os.Stderr, "uploading s3://%s/%s\n", *inp2.Bucket, *inp2.Key)
	start := time.Now()
	if err := req2.Send(); err != nil {
		return wrap(err, "putting object")
	}
	fmt.Fprintf(os.Stderr, "upload completed in %v\n", time.Since(start))
	return nil
}
