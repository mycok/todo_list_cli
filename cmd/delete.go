package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/mycok/todo_list_cli/api"
	"github.com/mycok/todo_list_cli/todo"
)

var del = api.Cmd{
	Name:  "del",
	Usage: "Delete task by providing a task [ID]",
	Action: func(w io.Writer, args ...string) error {
		f := api.GetFlag("file")
		if f == nil {
			return fmt.Errorf("%w: %q", api.ErrFlagNotFound, "file")
		}

		return deleteAction(w, f.Value.String(), args...)
	},
}

func deleteAction(w io.Writer, fName string, args ...string) error {
	var err error

	l := &todo.List{}

	if err = l.Load(fName); err != nil {
		return err
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	if err = l.Delete(id); err != nil {
		return err
	}

	if err = l.Save(fName); err != nil {
		return err
	}

	return nil
}

func init() {
	if err := api.Register(del); err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}
