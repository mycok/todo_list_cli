package cmds

import (
	"errors"
	"flag"
	"fmt"
)

// ErrDuplicateCmd represents a duplicate command error.
var ErrDuplicateCmd = errors.New("duplicate command")

// ErrDuplicateFlag represents a duplicate command flag error.
var ErrDuplicateFlag = errors.New("duplicate flag")

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

var cmdFlagArgs map[string]flag.Value

// AddFlag add a command flag to the map.
func AddFlag(name string, val flag.Value) error {
	flag, ok := cmdFlagArgs[name]
	if ok {
		return fmt.Errorf("%w: flag %q already exists", ErrDuplicateFlag, flag)
	}

	cmdFlagArgs[name] = val

	return nil
}

// GetFlag retrieves the command flag that matches the provided name.
func GetFlag(name string) flag.Value {
	flag, ok := cmdFlagArgs[name]
	if !ok {
		return nil
	}

	return flag
}

// CmdFlags returns the command flags map.
func CmdFlags() map[string]flag.Value {
	return cmdFlagArgs
}

func init() {
	commands = make(map[string]Command)
	cmdFlagArgs = make(map[string]flag.Value, 3)
}
