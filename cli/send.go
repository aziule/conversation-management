package cli

import "fmt"

// ReceiveCommand is the command responsible for sending a message to the user.
// Its purpose is mainly to test the bot.
type SendCommand struct{}

// Execute runs the command
func (c *SendCommand) Execute(args []string) error {
	fmt.Println("hey")
	return nil
}

// GetName returns the command's name, to be used when invoking it from the cli
func (c *SendCommand) GetName() string {
	return "Send"
}
