package main

import (
	"flag"
	"fmt"
	"os"
	"bufio"
	"io"
	"strings"

	"github.com/myok/todo_list_cli/todo"
)

var todoFileName = ".todo.json"

func main() {
	add := flag.Bool("add", false, "Add new todo item to the todo list")
	list := flag.Bool("list", false, "List all available todo items")
	done := flag.Int("done", 0, "Mark todo list item as complete")
	del := flag.Int("del", 0, "Delete todo list item from the list")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s tool. Developed by mycok\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright @2022\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage information:\n")

		flag.PrintDefaults()
	}

	flag.Parse()

	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	l := &todo.List{}

	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}

	switch {
	case *list:
		fmt.Print(l)
	case *done > 0:
		if err := l.Complete(*done); err != nil {
			fmt.Fprintln(os.Stderr, err)

			os.Exit(1)
		}

		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)

			os.Exit(1)
		}
	case *add:
		// If any args (excluding flags) are provided, they will be used as the name
		// of the new todo item. else we will read from user input.
		t, err := readTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)

			os.Exit(1)
		}

		l.Add(t)

		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)

			os.Exit(1)
		}
	case *del > 0:
		if err := l.Delete(*del); err != nil {
			fmt.Fprintln(os.Stderr, err)

			os.Exit(1)
		}

		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)

			os.Exit(1)
		}

	default:
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}
}

func readTask(r io.Reader, args... string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	// Only open an interactive shell session if no args have been provided.
	s := bufio.NewScanner(r)
	// Blocking call.
	s.Scan()
	if err := s.Err(); err != nil {
		return "", err
	}

	if len(s.Text()) == 0 {
		return "", fmt.Errorf("Task cannot be blank")
	}

	return s.Text(), nil
}
