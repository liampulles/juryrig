package wire

import (
	"fmt"
	"os"

	goConfig "github.com/liampulles/go-config"
	"github.com/liampulles/juryrig/internal/command"
	"github.com/liampulles/juryrig/internal/config"
)

// Run takes in program arguments and config source, does something,
// and returns an exit code.
func Run(args []string, cfgSource goConfig.Source) int {
	cmdManager := wire(cfgSource)
	if err := cmdManager.Run(args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
		return 1
	}
	return 0
}

func wire(cfgSource goConfig.Source) *command.Manager {
	cfgService := config.NewServiceImpl(cfgSource)

	genCmd := command.NewGen(cfgService)

	return command.NewManager(map[string]command.Command{
		"gen": genCmd,
	})
}
