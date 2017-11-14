package cli

import (
	"bytes"
	"errors"
	"flag"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/aziule/conversation-management/app"
	log "github.com/sirupsen/logrus"
)

// ReceiveCommand is the command responsible for running our bot using the given configuration.
// This is the main command.
type ReceiveCommand struct {
	configFilePath string
	dataFilePath   string
}

// NewReceiveCommand returns a new ReceiveCommand
func NewReceiveCommand() *ReceiveCommand {
	return &ReceiveCommand{}
}

// Usage returns the usage text for the command
func (c *ReceiveCommand) Usage() string {
	return `receive [-config=./config.json] -data=file.json:
	Sends a message to the bot, in order to fake a message sent by a user.`
}

// Execute runs the command
func (c *ReceiveCommand) Execute(f *flag.FlagSet) error {
	// Shared flags between the commands
	config, err := app.LoadConfig(c.configFilePath)

	if err != nil {
		// @todo: move this to the handler
		log.Fatalf("An error occurred when loading the config: %s", err)
	}

	// @todo: move it also somewhere else
	if config.Debug {
		log.SetLevel(log.DebugLevel)
	}

	data, err := ioutil.ReadFile(c.dataFilePath)

	if err != nil {
		return errors.New("Err")
	}

	// For now, only ping localhost
	url := "http://localhost:" + strconv.Itoa(config.ListeningPort) + "/fb"
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	client := http.DefaultClient
	_, err = client.Do(request)

	if err != nil {
		return err
	}

	return nil
}

// FlagSet returns the command's flag set
func (c *ReceiveCommand) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.configFilePath, "config", "config.json", "Config file path")
	f.StringVar(&c.dataFilePath, "data", "", "The data to receive")
}

// Name returns the command's name, to be used when invoking it from the cli
func (c *ReceiveCommand) Name() string {
	return "receive"
}
