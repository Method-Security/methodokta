package main

import (
	"flag"
	"os"

	"github.com/method-security/methodokta/cmd"
)

var version = "none"

func main() {
	flag.Parse()

	methodokta := cmd.NewMethodOkta(version)
	methodokta.InitRootCommand()

	if err := methodokta.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
