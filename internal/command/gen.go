package command

import (
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

// Run runs gen.
func (g *Gen) Run(args []string) error {
	cfg, err := g.cfgService.Read()
	if err != nil {
		return fmt.Errorf("could not fetch config: %w", err)
	}

	_, err = parse.Read(cfg.BaseFilename)
	if err != nil {
		return fmt.Errorf("could not parse file %s: %w", cfg.BaseFilename, err)
	}

	return nil
}
