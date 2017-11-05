// Package cli provides CLI methods to run or interact with bots.
package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

// Command is the main inteface for creating new commands to be used from the CLI
type Command interface {
	Name() string
	Execute(*flag.FlagSet) error
	SetFlags(*flag.FlagSet)
	Usage() string
}

// CliHandler is the handler responsible for executing commands
type CliHandler struct {
	topLevelFlags *flag.FlagSet
	commands      []Command
}

// NewHandler creates a new CliHandler
func NewHandler() *CliHandler {
	handler := &CliHandler{
		topLevelFlags: flag.CommandLine,
	}

	handler.topLevelFlags.Usage = func() { handler.explain() }

	return handler
}

// explain explains how to use the commands and what commands are available
func (h *CliHandler) explain() {
	fmt.Fprintln(os.Stderr, "COMMANDS:")

	for _, c := range h.commands {
		fmt.Fprintln(os.Stderr, c.Usage())
	}
}

// RegisterCommand adds a command to the list of available commands
func (h *CliHandler) RegisterCommand(command Command) {
	h.commands = append(h.commands, command)
}

// HandleCommand handles a command identified by its name
func (h *CliHandler) Handle() error {
	if !flag.Parsed() {
		flag.Parse()
	}

	if h.topLevelFlags.NArg() < 1 {
		h.topLevelFlags.Usage()
		// @todo: handle
		return errors.New("Invalid arguments passed")
	}

	name := h.topLevelFlags.Arg(0)

	for _, command := range h.commands {
		if command.Name() != name {
			continue
		}

		f := flag.NewFlagSet(name, flag.ContinueOnError)
		f.Usage = func() { fmt.Println(command.Usage()) }
		command.SetFlags(f)

		if err := f.Parse(h.topLevelFlags.Args()[1:]); err != nil {
			// @todo: handle
			return errors.New("Error 1")
		}

		return command.Execute(f)
	}

	h.topLevelFlags.Usage()
	// @todo: handle
	return errors.New("Error 2")
}
