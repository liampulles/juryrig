package command

import (
	"errors"
	"fmt"
	"strings"
)

// Command runs commands.
type Command interface {
	Run(args []string) error
}

// Manager implements Command by delegating to other commands.
type Manager struct {
	commands map[string]Command
}

var _ Command = &Manager{} //nolint:exhaustruct

// NewManager is a constructor.
func NewManager(commands map[string]Command) *Manager {
	return &Manager{
		commands: commands,
	}
}

var (
	ErrNoCommand      = errors.New("no command given")
	ErrInvalidCommand = errors.New("no such command")
)

// Run runs the manager.
func (m *Manager) Run(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("valid commands: %s: %w",
			m.listCommands(), ErrNoCommand)
	}

	name := args[0]
	cmd := m.findCommand(name)

	if cmd == nil {
		return fmt.Errorf("valid commands: %s: %w",
			m.listCommands(), ErrInvalidCommand)
	}

	if err := cmd.Run(args[1:]); err != nil {
		return fmt.Errorf("\"%s\" failed: %w", name, err)
	}

	return nil
}

func (m *Manager) listCommands() string {
	names := make([]string, len(m.commands))
	i := 0

	for name := range m.commands {
		names[i] = name
		i++
	}

	return fmt.Sprintf("[%s]", strings.Join(names, ", "))
}

// Note: result is nillable.
func (m *Manager) findCommand(name string) Command {
	for cmdName, cmd := range m.commands {
		if cmdName == name {
			return cmd
		}
	}

	return nil
}
