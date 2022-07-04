package api

import (
	"errors"
	"fmt"
)

var (
	// ErrDuplicateCmd represents a duplicate command error.
	ErrDuplicateCmd = errors.New("duplicate command")
	// ErrInvalidCmd represents a duplicate command error.
	ErrCmdNotFound = errors.New("command not found")
)

var commands map[string]Command

// Register adds a command to the commands map.
func Register(cmd Command) error {
	name := cmd.GetName()

	for k := range commands {
		if k == name {
			return fmt.Errorf("%w: command %q already exists", ErrDuplicateCmd, name)
		}
	}

	commands[name] = cmd

	return nil
}

// Get retrieves a command that matches the provided name.
func Get(name string) Command {
	c, ok := commands[name]
	if !ok {
		return nil
	}

	return c
}

// Commands returns the commands map.
func Commands() map[string]Command {
	return commands
}

func init() {
	commands = make(map[string]Command)
}
