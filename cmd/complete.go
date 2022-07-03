package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/mycok/todo_list_cli/api"
	"github.com/mycok/todo_list_cli/todo"
)

var complete = api.Cmd{
	Name:  "complete",
	Usage: "Mark a task as complete",
	Action: func(w io.Writer, args ...string) error {
		f := api.GetFlag("file")
		if f == nil {
			return fmt.Errorf("%w: %q", api.ErrMissingFlag, "file")
		}

		return completeAction(w, f.Value.String(), args...)
	},
}

func completeAction(w io.Writer, fName string, args ...string) error {
	l := &todo.List{}

	if err := l.Load(fName); err != nil {
		return err
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	if err := l.Complete(id); err != nil {
		return err
	}

	if err := l.Save(fName); err != nil {
		return err
	}

	return nil
}

func init() {
	if err := api.Register(complete); err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}
