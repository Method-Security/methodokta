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

	methodokta.InitApplicationCommand()
	methodokta.InitDeviceCommand()
	methodokta.InitGroupCommand()
	methodokta.InitOrgCommand()
	methodokta.InitUserCommand()

	if err := methodokta.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
