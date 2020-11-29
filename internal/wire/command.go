package wire

import (
	"strings"

	"github.com/liampulles/juryrig/internal/gen"

	"github.com/liampulles/juryrig/internal/config"
)

var commandRegistry map[string]Command = map[string]Command{
	"gen": gen.Command,
}

// Command defines a method which takes arguments and environment config,
// does something, and returns nil if successful, otherwise an error.
type Command func(args []string, cfg *config.Config) error

func determineCommand(arg string) *Command {
	cmd, ok := commandRegistry[arg]
	if !ok {
		return nil
	}
	return &cmd
}

func availableCommands() string {
	var cmdNames []string
	for k := range commandRegistry {
		cmdNames = append(cmdNames, k)
	}
	return "[" + strings.Join(cmdNames, ", ") + "]"
}
