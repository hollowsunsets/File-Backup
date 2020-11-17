package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hollowsunsets/fcopy/cmd"
	"os"
)

func main() {
	sess, err := session.NewSession()
	_, err = sess.Config.Credentials.Get()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	cmd.Execute()
}



