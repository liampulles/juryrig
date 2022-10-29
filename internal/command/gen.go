package command

import (
	"errors"
	"flag"
	"fmt"

	"github.com/liampulles/juryrig/internal/config"
	"github.com/liampulles/juryrig/internal/parse"
)

// Gen implements Command to generate go files.
type Gen struct {
	cfgService config.Service
}

var _ Command = &Gen{} //nolint:exhaustruct

// NewGen is a constructor.
func NewGen(cfgService config.Service) *Gen {
	return &Gen{
		cfgService: cfgService,
	}
}

// Run generates mappers, as per a the spec in the comments of the file.
func (g *Gen) Run(args []string) error {
	// Read args
	_, err := g.parseArgs(args)
	if err != nil {
		return err
	}

	// Read config
	cfg, err := g.cfgService.Read()
	if err != nil {
		return fmt.Errorf("could not fetch config: %w", err)
	}

	// Parse mappers
	_, err = parse.Read(cfg.BaseFilename)
	if err != nil {
		return fmt.Errorf("could not parse file %s: %w", cfg.BaseFilename, err)
	}

	return nil
}

type arguments struct {
	OutputFile string
}

var ErrInvalidArgs = errors.New("invalid args")

func (g *Gen) parseArgs(args []string) (arguments, error) {
	fs := flag.NewFlagSet("gen", flag.ContinueOnError)
	outputFile := fs.String("o", "", "output file")

	if err := fs.Parse(args); err != nil {
		fs.Usage()
		return arguments{}, fmt.Errorf("could not parse args for gen: %w", err)
	}

	if *outputFile == "" {
		return arguments{}, ErrInvalidArgs
	}

	return arguments{
		OutputFile: *outputFile,
	}, nil
}
