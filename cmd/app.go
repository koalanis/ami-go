// The app defines the main go process / control flow for the program
// it is responsible for running the program in CLI vs interactive mode
// it also reads the configuration of the bot server before performing assistant actions
package app

import (
	"os"

	"github.com/koalanis/ami-go/cli"
	"github.com/koalanis/ami-go/server"
)

func Run() {
	amigo := cli.ParseArgs()

	if !cli.ValidateExecutionContext(amigo) {
		os.Exit(1)
		return
	}

	if amigo.ServerMode {
		server.ServerInit(amigo)
	} else {
		cli.CliInit(amigo)
	}

	os.Exit(0)
}
