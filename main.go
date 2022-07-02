package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mycok/todo_list_cli/cmds"
	"github.com/mycok/todo_list_cli/colors"
	"github.com/mycok/todo_list_cli/todo"
)

func main() {
	fs := flag.NewFlagSet("todo_list_cli", flag.ExitOnError)
	fs.SetOutput(os.Stdout)

	fs.Usage = func() {
		fmt.Fprintf(
			fs.Output(),
			"%stodoCLI tool: Developed by mycok%s\n",
			colors.Cyan, colors.Reset,
		)

		fmt.Fprintf(
			fs.Output(),
			"%s<github.com/mycok>: Copyright @2022%s\n",
			colors.Cyan, colors.Reset,
		)

		fmt.Println()
		fmt.Fprintf(fs.Output(), "%sUsage information:%s\n", colors.Magenta, colors.Reset)

		fs.PrintDefaults()

		fmt.Println()
		fmt.Fprintf(fs.Output(), "%sExamples:%s\n", colors.Magenta, colors.Reset)

		fmt.Println("-add go shopping today [Adds a single item]")
		fmt.Println("-add [Adds multiple items using the input shell]")

	}

	file := fs.String("file", ".todo.json", "File name to store todo list")
	details := fs.Bool("details", false, "List all available todo items showing more details like date & time")
	completed := fs.Bool("completed", false, "List all available todo items including the completed")

	// Commands
	add := fs.Bool("add", false, "Add new todo item to the todo list")
	list := fs.Bool("list", false, "List all available todo items")
	complete := fs.Int("complete", 0, "Mark todo list item as complete")
	del := fs.Int("del", 0, "Delete todo list item from the list")

	if err := fs.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}

	// Add all provided flags to the cmdFlagArgs map.
	for i := 0; i < fs.NFlag(); i++ {
		fs.Visit(func(f *flag.Flag) {
			err := cmds.AddFlag(f.Name, f.Value)
			if err != nil {
				if !errors.Is(err, cmds.ErrDuplicateFlag) {
					fmt.Fprintln(os.Stderr, err)

					os.Exit(1)
				}
			}
		})
	}

	todoFileName := *file

	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	if err := run(fs, todoFileName, *del, *complete, *add, *list, *details, *completed); err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}

func run(fs *flag.FlagSet, fName string, del, complete int, add, list, details, completed bool) error {
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
		tasks, err := readTasksInput(os.Stdin, fs.Args()...)
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
		fmt.Fprintf(os.Stdout, "%sInvalid option!%s\n", colors.Red, colors.Reset)
		fmt.Fprintln(os.Stdout)
		// Show usage information
		fs.Usage()
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
