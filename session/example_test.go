package session_test

import (
	"fmt"

	"github.com/segmentio/go-utils/session"
)

func ExampleGet() {
	sess, err := session.Get("stage", nil)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Println(sess.Config.Region)
}
