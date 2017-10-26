package main

import (
	"github.com/aziule/conversation-management/cli"
	log "github.com/sirupsen/logrus"
)

func main() {
	cliHandler := cli.NewHandler()
	cliHandler.RegisterCommand(cli.NewRunCommand())
	cliHandler.RegisterCommand(cli.NewReceiveCommand())
	err := cliHandler.Handle()

	if err != nil {
		log.Fatalf("An error occurred when handling the command: %s", err)
	}
}
