package main

import (
	"os"

	"github.com/liampulles/juryrig/internal/wire"
)

func main() {
	// Delegate functionality to Run, do the difficult-to-test
	// bits here.
	os.Exit(wire.Run(os.Args, os.Environ()))
}
