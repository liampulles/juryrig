package command

import (
	"fmt"
	"strings"

	"github.com/liampulles/juryrig/internal/config"
	"github.com/liampulles/juryrig/internal/parse"
)

// Command runs commands
type Command interface {
	Run(args []string) error
}

// Manager implements Command by delegating to other commands
type Manager struct {
	commands map[string]Command
}

var _ Command = &Manager{}

// NewManager is a constructor
func NewManager(commands map[string]Command) *Manager {
	return &Manager{
		commands: commands,
	}
}

// Run runs the manager
func (m *Manager) Run(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("need at least one arg for the command to run. commands: %s",
			m.listCommands())
	}
	name := args[0]
	cmd := m.findCommand(name)
	if cmd == nil {
		return fmt.Errorf("no known command with name \"%s\". commands: %s",
			name, m.listCommands())
	}
	if err := cmd.Run(args[1:]); err != nil {
		return fmt.Errorf("\"%s\" failed: %w", name, err)
	}
	return nil
}

func (m *Manager) listCommands() string {
	var names []string
	for name := range m.commands {
		names = append(names, name)
	}
	return fmt.Sprintf("[%s]", strings.Join(names, ", "))
}

// Note: result is nillable
func (m *Manager) findCommand(name string) Command {
	for cmdName, cmd := range m.commands {
		if cmdName == name {
			return cmd
		}
	}
	return nil
}

// Gen implements command to generate go files
type Gen struct {
	cfgService config.Service
}

var _ Command = &Gen{}

// NewGen is a constructor
func NewGen(cfgService config.Service) *Gen {
	return &Gen{
		cfgService: cfgService,
	}
}

// Run runs gen
func (g *Gen) Run(args []string) error {
	cfg, err := g.cfgService.Read()
	if err != nil {
		return fmt.Errorf("could not fetch config: %w", err)
	}

	_, err = parse.ParseFileWithMapperDefs(cfg.BaseFilename)
	if err != nil {
		return fmt.Errorf("could not parse file %s: %w", cfg.BaseFilename, err)
	}
	return nil
}
