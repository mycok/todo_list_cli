package api

import (
	"errors"
	"flag"
	"fmt"
)

// ErrDuplicateFlag represents a duplicate command flag error.
var ErrDuplicateFlag = errors.New("duplicate flag")

var cmdFlags map[string]flag.Value

// AddFlag add a command flag to the map.
func AddFlag(name string, val flag.Value) error {
	flag, ok := cmdFlags[name]
	if ok {
		return fmt.Errorf("%w: flag %q already exists", ErrDuplicateFlag, flag)
	}

	cmdFlags[name] = val

	return nil
}

// GetFlag retrieves the command flag that matches the provided name.
func GetFlag(name string) flag.Value {
	val, ok := cmdFlags[name]
	if !ok {
		return nil
	}

	return val
}

// CmdFlags returns the command flags map.
func CmdFlags() map[string]flag.Value {
	return cmdFlags
}

func init() {
	cmdFlags = make(map[string]flag.Value, 3)
}
