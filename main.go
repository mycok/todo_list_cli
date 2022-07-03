package main

import (
	"fmt"
	"os"

	"github.com/mycok/todo_list_cli/cmd"
)

func main() {
	if err := cmd.Execute(os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}

// func run(fs *flag.FlagSet, fName string, del, complete int, add, list, details, completed bool) error {
// 	var err error

// 	l := &todo.List{}

// 	if err = l.Load(fName); err != nil {
// 		return err
// 	}

// 	switch {
// 	case list:
// 		l.ListItems(os.Stdout, details, completed)

// 	case complete > 0:
// 		if err = l.Complete(complete); err != nil {
// 			return err
// 		}

// 		if err = l.Save(fName); err != nil {
// 			return err
// 		}

// 	case add:
// 		If any args (excluding flags) are provided, they will be used as the name
// 		of the new todo item. else we will read from standard input.
// 		tasks, err := readTasksInput(os.Stdin, fs.Args()...)
// 		if err != nil {
// 			return err
// 		}

// 		for _, t := range tasks {
// 			l.Add(t)
// 		}

// 		if err = l.Save(fName); err != nil {
// 			return err
// 		}

// 	case del > 0:
// 		if err = l.Delete(del); err != nil {
// 			return err
// 		}

// 		if err = l.Save(fName); err != nil {
// 			return err
// 		}

// 	default:
// 		fmt.Fprintf(os.Stdout, "%sInvalid option!%s\n", colors.Red, colors.Reset)
// 		fmt.Fprintln(os.Stdout)
// 		// Show usage information
// 		fs.Usage()
// 	}

// 	return nil
// }
