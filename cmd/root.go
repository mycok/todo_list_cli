package cmd

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mycok/todo_list_cli/api"
	"github.com/mycok/todo_list_cli/colors"
)

func Run(w io.Writer) error {
	fs, err := handleCmdlineFlags(w)
	if err != nil {
		return err
	}

	// Execute commands as provided by the user.
	switch os.Args[1] {
	case "add":
		return executeCmd(w, os.Args[1], fs.Args()...)

	case "list":
		return executeCmd(w, os.Args[1], fs.Args()...)

	case "complete":
		return executeCmd(w, os.Args[1], fs.Args()...)

	default:
		fmt.Fprintf(
			w,
			"%sInvalid command!%s\n",
			colors.Red,
			colors.Reset,
		)
		fmt.Fprintln(w)
		// Show usage information
		fs.Usage()
	}

	return nil
}

func handleCmdlineFlags(w io.Writer) (*flag.FlagSet, error) {
	fs := flag.NewFlagSet("todo_list_cli", flag.ExitOnError)
	fs.SetOutput(w)

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

		fmt.Fprintln(fs.Output())
		fmt.Fprintf(
			fs.Output(),
			"%sUsage information:%s\n",
			colors.Magenta,
			colors.Reset,
		)

		fs.PrintDefaults()

		fmt.Fprintln(fs.Output())
		fmt.Fprintf(
			fs.Output(),
			"%sExamples:%s\n",
			colors.Magenta,
			colors.Reset,
		)

		fmt.Fprintln(fs.Output(), "-add go shopping today [Adds a single item]")
		fmt.Fprintln(fs.Output(), "-add [Adds multiple items using the input shell]")
	}

	todoFileName := ".todo.json"

	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	fs.String("file", todoFileName, "File name to store todo list")
	fs.Bool("details", false, "List all todo tasks with details")
	fs.Bool("completed", false, "List all todo tasks including the completed")

	var err error

	if err = fs.Parse(os.Args[2:]); err != nil {
		return nil, err
	}

	// Walk through all provided flags [both set and unset] and add each
	// to the cmdFlagArgs map.
	fs.VisitAll(func(f *flag.Flag) {
		err = api.AddFlag(f)
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

func executeCmd(w io.Writer, cmd string, args ...string) error {
	c := api.Get(cmd)
	if c == nil {
		return nil
	}

	if err := c.Run(w, args...); err != nil {
		return err
	}

	return nil
}
