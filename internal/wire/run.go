package wire

import (
	"fmt"
)

// Run takes in program arguments and environment, does something
// and returns an exit code.
func Run(args []string, env []string) int {
	fmt.Println("Hello World!")
	return 0
}
