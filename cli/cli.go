package cli

import (
	"errors"
	"flag"
)

// Command is the main inteface for creating new commands to be used from the CLI
type Command interface {
	Name() string
	Execute(*flag.FlagSet) error
	SetFlags(*flag.FlagSet)
}

type CliHandler struct {
	topLevelFlags *flag.FlagSet
	commands      []Command
}

func NewHandler() *CliHandler {
	return &CliHandler{
		topLevelFlags: flag.CommandLine,
	}
}

func (h *CliHandler) RegisterCommand(command Command) {
	h.commands = append(h.commands, command)
}

// HandleCommand handles a command identified by its name
func (h *CliHandler) Handle() error {
	if !flag.Parsed() {
		flag.Parse()
	}

	if h.topLevelFlags.NArg() < 1 {
		// @todo: handle
		return errors.New("Invalid arguments passed")
	}

	name := h.topLevelFlags.Arg(0)

	for _, command := range h.commands {
		if command.Name() != name {
			continue
		}

		f := flag.NewFlagSet(name, flag.ContinueOnError)
		command.SetFlags(f)

		if err := f.Parse(h.topLevelFlags.Args()[1:]); err != nil {
			// @todo: handle
			return errors.New("Error")
		}

		return command.Execute(f)
	}

	// @todo: handle
	return errors.New("Error")
}
