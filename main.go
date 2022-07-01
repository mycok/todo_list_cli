package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mycok/todo_list_cli/colors"
	"github.com/mycok/todo_list_cli/todo"
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
		fmt.Fprintf(
			flag.CommandLine.Output(),
			"%stodoCLI tool: Developed by mycok%s\n",
			colors.Cyan, colors.Reset,
		)

		fmt.Fprintf(
			flag.CommandLine.Output(),
			"%s<github.com/mycok>: Copyright @2022%s\n",
			colors.Cyan, colors.Reset,
		)

		fmt.Println()
		fmt.Fprintf(flag.CommandLine.Output(), "%sUsage information:%s\n", colors.Magenta, colors.Reset)

		flag.PrintDefaults()

		fmt.Println()
		fmt.Fprintf(flag.CommandLine.Output(), "%sExamples:%s\n", colors.Magenta, colors.Reset)

		fmt.Println("-add go shopping today [Adds a single item]")
		fmt.Println("-add [Adds multiple items using the input shell]")

	}

	flag.Parse()

	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	if err := run(todoFileName, *del, *complete, *add, *list, *details, *completed); err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}

func run(fName string, del, complete int, add, list, details, completed bool) error {
	var err error

	l := &todo.List{}

	if err = l.Load(fName); err != nil {
		return err
	}

	switch {
	case list:
		l.ListItems(os.Stdout, details, completed)

	case complete > 0:
		if err = l.Complete(complete); err != nil {
			return err
		}

		if err = l.Save(fName); err != nil {
			return err
		}

	case add:
		// If any args (excluding flags) are provided, they will be used as the name
		// of the new todo item. else we will read from standard input.
		tasks, err := readTasksInput(os.Stdin, flag.Args()...)
		if err != nil {
			return err
		}

		for _, t := range tasks {
			l.Add(t)
		}

		if err = l.Save(fName); err != nil {
			return err
		}

	case del > 0:
		if err = l.Delete(del); err != nil {
			return err
		}

		if err = l.Save(fName); err != nil {
			return err
		}

	default:
		fmt.Fprintln(os.Stdout, "Invalid option")
		fmt.Fprintln(os.Stdout)
		// Show usage information
		flag.Usage()
	}

	return nil
}

func readTasksInput(r io.Reader, args ...string) ([]string, error) {
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
			// Case of an empty string.
			break
		}
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}
