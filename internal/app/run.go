package app

import (
	"fmt"
	"os"

	"github.com/liampulles/juryrig/internal/command"
	"github.com/liampulles/juryrig/internal/wire"

	goConfig "github.com/liampulles/go-config"
)

// Run takes in program arguments and config source, does something,
// and returns an exit code.
func Run(args []string, cfgSource goConfig.Source) int {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, args[0]+" requires at least one argument.")
		return 1
	}

	cmdPtr := command.Determine(args[1])
	if cmdPtr == nil {
		fmt.Fprintf(os.Stderr, "%s has no command registered for %s. available commands: %s\n", args[0], args[1], command.Available())
		return 2
	}
	cmd := *cmdPtr

	wiring := wire.Connect(cfgSource)

	if err := cmd(args[1:], wiring); err != nil {
		fmt.Fprintf(os.Stderr, "error running %s: %s\n", args[1], err)
		return 3
	}

	return 0
}
