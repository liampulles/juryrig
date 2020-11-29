package command

import (
	"fmt"

	"github.com/liampulles/juryrig/internal/wire"
)

// Gen is a program directive used to generate Go code.
func Gen(args []string, wiring *wire.Wiring) error {
	fmt.Println("Hello World!")
	return nil
}
