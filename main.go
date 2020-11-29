package main

import (
	"os"

	goConfig "github.com/liampulles/go-config"

	"github.com/liampulles/juryrig/internal/app"
)

func main() {
	// Delegate functionality to Run, do the difficult-to-test
	// bits here.
	os.Exit(app.Run(os.Args, goConfig.NewEnvSource()))
}
