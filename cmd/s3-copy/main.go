// s3-copy copies an S3 file from one bucket to another.
//
// Example usage:
//
//     s3-copy --from=akey.csv --from-bucket=stagebucket --from-profile=stage \
//          --from-region=us-west-1 --to-bucket=prodbucket --to-profile=prod
//
// It achieves this by using one set of credentials to download the file to your
// machine, then (optionally) a different set of credentials or regions to
// upload the file.
//
// Due to the limitations of the aws-sdk-go library, the entire file is
// currently downloaded and read into memory before being uploaded again.
//
// Credentials are read out of your environment using `aws-vault`. This is
// designed to be run on a local machine.
//
// To call this function from another Go program, use the S3Copier interface in
// github.com/segmentio/go-utils.
package main

import (
	"flag"
	"fmt"
	"os"

	utils "github.com/segmentio/go-utils"
)

func checkError(err error, reason string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error %s: %v\n", reason, err)
		os.Exit(2)
	}
}

func main() {
	profile := flag.String("profile", "", "Profile to use for get and put. Cannot use with -to-profile or -from-profile")
	from := flag.String("from", "", "S3 file to get")
	fromBucket := flag.String("from-bucket", "", "S3 bucket to get object from")
	fromProfile := flag.String("from-profile", "", "Profile to use to fetch")
	fromRegion := flag.String("from-region", "", "Region to use for fetch (defaults to value in profile)")

	to := flag.String("to", "", "S3 file to put (defaults to the Key for from)")
	toBucket := flag.String("to-bucket", "", "S3 bucket to put object")
	toProfile := flag.String("to-profile", "", "Profile to use for put object")
	toRegion := flag.String("to-region", "", "Region to use for put (defaults to value in profile)")

	acl := flag.String("acl", "", "ACL to use for putting the object")
	timeout := flag.Duration("timeout", utils.DefaultS3Timeout, "Amount of time to allow for the operation")
	flag.Parse()
	if *profile != "" {
		if *toProfile != "" || *fromProfile != "" {
			os.Stderr.WriteString("Cannot set -profile and (-to-profile or -from-profile)\n")
			os.Exit(2)
		}
		*toProfile = *profile
		*fromProfile = *profile
	}
	// TODO: check whether from contains a s3:// URL, and get the from bucket
	// from that if so.
	if *fromBucket == "" {
		os.Stderr.WriteString("Please provide a from bucket\n")
		os.Exit(2)
	}
	if *from == "" {
		os.Stderr.WriteString("Please provide a from argument\n")
		os.Exit(2)
	}
	copier := &utils.S3Copier{
		FromProfile: *fromProfile,
		ToProfile:   *toProfile,
		Timeout:     *timeout,
	}
	if *to == "" {
		*to = *from
	}
	settings := &utils.S3CopySettings{
		FromKey:    *from,
		FromBucket: *fromBucket,
		ToKey:      *to,
		ToBucket:   *toBucket,
		FromRegion: *fromRegion,
		ToRegion:   *toRegion,
		ACL:        *acl,
	}
	if err := copier.Copy(settings); err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		os.Exit(2)
	}
}
