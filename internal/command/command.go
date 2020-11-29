package command

import (
	"strings"

	"github.com/liampulles/juryrig/internal/wire"
)

var commandRegistry map[string]Command = map[string]Command{
	"gen": Gen,
}

// Command defines a method which takes arguments and environment config,
// does something, and returns nil if successful, otherwise an error.
type Command func(args []string, wiring *wire.Wiring) error

// Determine resolves an arg to a Command, if present
func Determine(arg string) *Command {
	cmd, ok := commandRegistry[arg]
	if !ok {
		return nil
	}
	return &cmd
}

// Available returns a readable string of registered commands.
func Available() string {
	var cmdNames []string
	for k := range commandRegistry {
		cmdNames = append(cmdNames, k)
	}
	return "[" + strings.Join(cmdNames, ", ") + "]"
}
