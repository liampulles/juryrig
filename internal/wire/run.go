package wire

import (
	"fmt"
	"os"

	goConfig "github.com/liampulles/go-config"

	"github.com/liampulles/juryrig/internal/config"
)

// Run takes in program arguments and config source, does something,
// and returns an exit code.
func Run(args []string, cfgSource goConfig.Source) int {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, args[0]+" requires at least one argument.")
		return 1
	}

	cmdPtr := determineCommand(args[1])
	if cmdPtr == nil {
		fmt.Fprintf(os.Stderr, "%s has no command registered for %s. available commands: %s\n", args[0], args[1], availableCommands())
		return 2
	}
	cmd := *cmdPtr

	cfg, err := config.Read(cfgSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read config: %s\n", err)
		return 3
	}

	if err := cmd(args[1:], cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error running %s: %s\n", args[1], err)
		return 4
	}

	return 0
}
