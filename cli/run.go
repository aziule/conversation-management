package cli

import (
	"flag"

	"github.com/aziule/conversation-management/app"
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
	app.Run(c.configFilePath)

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
