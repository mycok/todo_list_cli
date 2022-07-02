package cmds

import (
	"io"
)

// Command must be implemented by terminal command types.
type Command interface {
	GetName() string
	GetUsage() string
	Run(w io.Writer, args ...string) error
}

// Cmd type represents a terminal command.
type Cmd struct {
	Name   string
	Usage  string
	Action func(w io.Writer, args ...string) error
}

// String returns the name of the Cmd instance.
func (c Cmd) String() string {
	return c.Name
}

// GetName returns the name of the Cmd instance.
func (c Cmd) GetName() string {
	return c.Name
}

// GetUsage returns the usage string of the Cmd instance.
func (c Cmd) GetUsage() string {
	return c.Usage
}

// Run executes the action function of the Cmd instance.
func (c Cmd) Run(w io.Writer, args ...string) error {
	return c.Action(w, args...)
}
