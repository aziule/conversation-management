package cli

import (
	"flag"
	"fmt"
)

// ReceiveCommand is the command responsible for simulating messages received
// from Facebook (or any other platform).
// Its purpose is mainly to test the bot.
type ReceiveCommand struct{}

// Execute runs the command
func (c *ReceiveCommand) Execute() error {
	var dtype string
	dataType := flag.String("data", "abc", "The kind of data to send")
	flag.StringVar(&dtype, "dtype", "abc", "a var")
	flag.Parse()

	fmt.Println("dataType:", *dataType)
	fmt.Println("dtype:", dtype)
	//
	//config, err := app.LoadConfig(*configFlagPath)
	//
	//if err != nil {
	//	log.Fatalf("An error occurred when loading the config: %s", err)
	//}

	//url := "http://localhost:" + config.ListeningPort
	return nil
}

// GetName returns the command's name, to be used when invoking it from the cli
func (c *ReceiveCommand) GetName() string {
	return "receive"
}
