package gen

import (
	"fmt"

	"github.com/liampulles/juryrig/internal/config"
)

// Command is a program directive used to generate Go code.
func Command(args []string, cfg *config.Config) error {
	fmt.Println("Hello World!")
	return nil
}
