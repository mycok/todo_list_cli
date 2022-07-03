package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mycok/todo_list_cli/api"
	"github.com/mycok/todo_list_cli/todo"
)

var list = api.Cmd{
	Name:  "list",
	Usage: "List all available todo tasks",
	Action: func(w io.Writer, args ...string) error {
		f := api.GetFlag("file")
		if f == nil {
			return fmt.Errorf("%w: %q", api.ErrMissingFlag, "file")
		}

		d := api.GetFlag("details")
		if d == nil {
			return fmt.Errorf("%w: %q", api.ErrMissingFlag, "details")
		}

		c := api.GetFlag("completed")
		if c == nil {
			return fmt.Errorf("%w: %q", api.ErrMissingFlag, "completed")
		}

		return listAction(
			w, f.Value.String(), d.Value.(flag.Getter).Get().(bool),
			c.Value.(flag.Getter).Get().(bool), args...,
		)
	},
}

func listAction(w io.Writer, fName string, details, completed bool, args ...string) error {
	l := &todo.List{}

	if err := l.Load(fName); err != nil {
		return err
	}

	l.ListItems(w, details, completed)

	return nil
}

func init() {
	if err := api.Register(list); err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}
