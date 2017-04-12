package session

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// A CredentialGetter is used to get valid (temporary) AWS credentials.
type CredentialGetter struct {
	Region        string // Region to use; if unset, defaults to region in profile.
	AssumeRoleTTL time.Duration
}

var DefaultCredentialGetter = &CredentialGetter{}

// Get returns a valid aws session.Session for the given realm. Get reads MFA
// tokens from in; if in is nil, os.Stdin will be used.
//
// The region for the request should be specified in the profile config file.
func Get(realm string, in io.Reader) (*session.Session, error) {
	return DefaultCredentialGetter.Get(realm, in)
}

// Get returns a valid aws session.Session for the given realm. Get reads MFA
// tokens from in; if in is nil, os.Stdin will be used.
//
// Get is designed for local access - getting credentials on your local machine.
func (c *CredentialGetter) Get(realm string, in io.Reader) (*session.Session, error) {
	// aws-vault does not have good API's for programmatic access; a lot of the
	// code we need is buried inside 'package main'. instead, hack
	//
	// in addition I couldn't figure out a way to print exactly 4 environment
	// variables, so just print all of them and filter the ones we need.
	//
	// also TODO: figure out the right realm to use on everyone's machines.
	assumeRoleTTL := 15 * time.Minute
	if c.AssumeRoleTTL != 0 {
		assumeRoleTTL = c.AssumeRoleTTL
	}
	cmd := exec.Command("aws-vault", "exec", "--assume-role-ttl", assumeRoleTTL.String(), realm, "env")
	if in == nil {
		in = os.Stdin
	}
	cmd.Stdin = in
	// MFA prompt gets printed here
	cmd.Stderr = os.Stderr
	var buf bytes.Buffer
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		io.Copy(os.Stdout, &buf)
		return nil, err
	}
	scanner := bufio.NewScanner(&buf)
	var key, secret, token, region string
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) < 2 {
			continue
		}
		switch parts[0] {
		case "AWS_ACCESS_KEY_ID":
			key = parts[1]
		case "AWS_SECRET_ACCESS_KEY":
			secret = parts[1]
		case "AWS_SESSION_TOKEN":
			token = parts[1]
		case "AWS_REGION":
			region = parts[1]
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if c.Region != "" {
		region = c.Region
	}
	creds := credentials.NewStaticCredentials(key, secret, token)
	return session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: creds,
	})
}
