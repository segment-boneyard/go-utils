package utils_test

import (
	"fmt"

	utils "github.com/segmentio/go-utils"
)

func ExampleS3Copier() {
	copier := utils.NewS3Copier("stageprofile", "prodprofile")
	err := copier.Copy(&utils.S3CopySettings{
		FromKey:    "foo-bar.csv",
		FromBucket: "stagebucket",
		ToBucket:   "prodbucket",
	})
	fmt.Println(err)
}
