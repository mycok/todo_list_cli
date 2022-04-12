package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/myok/todo_list_cli/todo"
)

var todoFileName = ".todo.json"

func main() {
	task := flag.String("task", "", "Todo item to be added to the todo list")
	list := flag.Bool("list", false, "List all available todo items")
	done := flag.Int("done", 0, "Mark the todo list item as complete")

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
	case *task != "":
		l.Add(*task)
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)

			os.Exit(1)
		}

	default:
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}
}
