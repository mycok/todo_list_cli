package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/myok/todo_list_cli/todo"
)

var todoFileName = ".todo.json"

func main() {
	add := flag.Bool("add", false, "Add new todo item to the todo list")
	list := flag.Bool("list", false, "List all available todo items")
	details := flag.Bool("details", false, "List all available todo items showing more details like date & time")
	completed := flag.Bool("completed", false, "List all available todo items including the completed")
	complete := flag.Int("complete", 0, "Mark todo list item as complete")
	del := flag.Int("del", 0, "Delete todo list item from the list")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s tool. Developed by mycok\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright @2022\n")
		fmt.Println()
		fmt.Fprintf(flag.CommandLine.Output(), "Usage information:\n")

		flag.PrintDefaults()

		fmt.Println()
		fmt.Println("Examples:")

		fmt.Println("-add go shopping today [Add a single item]")
		fmt.Println("-add [Add multiple items using the input shell]")

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
		l.ListItems(details, completed)
	case *complete > 0:
		if err := l.Complete(*complete); err != nil {
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
		tasks, err := readTasks(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)

			os.Exit(1)
		}

		for _, t := range tasks {
			l.Add(t)
		}

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

func readTasks(r io.Reader, args ...string) ([]string, error) {
	todos := []string{}

	if len(args) > 0 {
		todos = append(todos, strings.Join(args, " "))
		return todos, nil
	}

	// Only open an interactive shell session if no args have been provided.
	s := bufio.NewScanner(r)
	// Blocking call.
	for {
		s.Scan()

		if len(s.Text()) > 0 {
			todos = append(todos, s.Text())
		} else {
			break
		}
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}
