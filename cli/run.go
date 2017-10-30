package cli

import (
	"flag"

	"github.com/aziule/conversation-management/app"
	log "github.com/sirupsen/logrus"
)

// RunCommand is the command responsible for running our bot using the given configuration.
// This is the main command.
type RunCommand struct {
	configFilePath string
}

// NewRunCommand returns a new RunCommand
func NewRunCommand() *RunCommand {
	return &RunCommand{}
}

// Usage returns the usage text for the command
func (c *RunCommand) Usage() string {
	return `run [-config=./config.json]:
	Runs the server and listens to incoming messages`
}

// Execute runs the command
func (c *RunCommand) Execute(f *flag.FlagSet) error {
	config, err := app.LoadConfig(c.configFilePath)

	if err != nil {
		log.Fatalf("An error occurred when loading the config: %s", err)
	}

	if config.Debug {
		log.SetLevel(log.DebugLevel)
	}

	app.Run(config)

	return nil
}

// FlagSet returns the command's flag set
func (c *RunCommand) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.configFilePath, "config", "config.json", "Config file path")
}

// Name returns the command's name, to be used when invoking it from the cli
func (c *RunCommand) Name() string {
	return "run"
}
