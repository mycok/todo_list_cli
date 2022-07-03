package api

import (
	"errors"
	"flag"
	"fmt"
)

var (
	// ErrInvalidFlag represents an invalid command flag error.
	ErrInvalidFlag   = errors.New("invalid flag")
	// ErrDuplicateFlag represents a duplicate command flag error.
	ErrDuplicateFlag = errors.New("duplicate flag")
)

var cmdFlags map[string]*flag.Flag

// AddFlag add a command flag to the map.
func AddFlag(f *flag.Flag) error {
	flag, ok := cmdFlags[f.Name]
	if ok {
		return fmt.Errorf("%w: flag %q already exists", ErrDuplicateFlag, flag.Name)
	}

	cmdFlags[f.Name] = f

	return nil
}

// GetFlag retrieves the command flag that matches the provided name.
func GetFlag(name string) *flag.Flag {
	f, ok := cmdFlags[name]
	if !ok {
		return nil
	}

	return f
}

// CmdFlags returns the command flags map.
func CmdFlags() map[string]*flag.Flag {
	return cmdFlags
}

func init() {
	cmdFlags = make(map[string]*flag.Flag, 3)
}
