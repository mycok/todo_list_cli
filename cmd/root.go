package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/mycok/todo_list_cli/api"
	"github.com/mycok/todo_list_cli/colors"
)

func Execute() error {
	fs, err := handleCmdlineFlags()
	if err != nil {
		return err
	}

	switch os.Args[1] {
	case "add":
		c := api.Get("add")
		// TODO: assert to api.Cmd.
		if c == nil {
			return nil
		}

		if err = c.Run(os.Stdout, fs.Args()...); err != nil {
			return err
		}
	default:
		fmt.Fprintf(
			os.Stdout,
			"%sInvalid command!%s\n",
			colors.Red,
			colors.Reset,
		)
		fmt.Fprintln(os.Stdout)
		// Show usage information
		fs.Usage()
	}

	return nil
}

func handleCmdlineFlags() (*flag.FlagSet, error) {
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
		fmt.Fprintf(
			fs.Output(),
			"%sUsage information:%s\n",
			colors.Magenta,
			colors.Reset,
		)

		fs.PrintDefaults()

		fmt.Println()
		fmt.Fprintf(
			fs.Output(),
			"%sExamples:%s\n",
			colors.Magenta,
			colors.Reset,
		)

		fmt.Println("-add go shopping today [Adds a single item]")
		fmt.Println("-add [Adds multiple items using the input shell]")

	}

	todoFileName := ".todo.json"

	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	fs.String("file", todoFileName, "File name to store todo list")
	fs.Bool("details", false, "List all todo tasks with details")
	fs.Bool("completed", false, "List all todo tasks including the completed")

	// Commands
	// add := fs.Bool("add", false, "Add new todo item to the todo list")
	// list := fs.Bool("list", false, "List all available todo items")
	// complete := fs.Int("complete", 0, "Mark todo list item as complete")
	// del := fs.Int("del", 0, "Delete todo list item from the list")
	var err error

	if err = fs.Parse(os.Args[2:]); err != nil {
		return nil, err
	}

	// Walk through all provided flags [both set and unset] and add each
	// to the cmdFlagArgs map.
	fs.VisitAll(func(f *flag.Flag) {
		err = api.AddFlag(f.Name, f.Value)
		if err != nil {
			if !errors.Is(err, api.ErrDuplicateFlag) {
				return
			}
		}
	})

	if err != nil {
		return nil, err
	}

	return fs, nil
}
