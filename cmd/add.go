package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mycok/todo_list_cli/api"
	"github.com/mycok/todo_list_cli/todo"
)

var add = api.Cmd{
	Name:  "add",
	Usage: "Add new task(s)",
	Action: func(w io.Writer, args ...string) error {
		f := api.GetFlag("file")
		if f == nil {
			return fmt.Errorf("%w: %q", api.ErrMissingFlag, "file")
		}

		return addAction(os.Stdin, os.Stdout, f.Value.String(), args...)
	},
}

func addAction(r io.Reader, w io.Writer, fName string, args ...string) error {
	var err error
	// Initialize an empty List.
	l := &todo.List{}

	if err = l.Load(fName); err != nil {
		return err
	}

	tasks, err := readTaskInput(r, args...)
	if err != nil {
		return err
	}

	for _, t := range tasks {
		l.Add(t)
	}

	if err = l.Save(fName); err != nil {
		return err
	}

	return nil
}

func readTaskInput(r io.Reader, args ...string) ([]string, error) {
	todos := []string{}

	if len(args) > 0 {
		todos = append(todos, strings.Join(args, " "))

		return todos, nil
	}

	// Only open an interactive shell session if no args have been provided.
	s := bufio.NewScanner(r)

	// Blocking call.
	// Capture multi-line input by calling s.Scan in an infinite loop that
	// only terminates if it scans an empty string.
	for {
		s.Scan()

		// Called each time the [Enter] key is pressed to terminate the current
		// string from the standard input / terminal.
		if len(s.Text()) > 0 {
			todos = append(todos, s.Text())
		} else {
			// Case of an empty string.
			break
		}
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func init() {
	if err := api.Register(add); err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}
