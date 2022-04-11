package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/myok/todo_list_cli/todo"
)

const todoFileName = ".todo.json"

func main() {
	l := &todo.List{}
	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}

	switch {
	case len(os.Args) == 1:
		for _, item := range *l {
			fmt.Println(item.Task)
		}
	default:
		todoName := strings.Join(os.Args[1:], " ")
		l.Add(todoName)
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)

			os.Exit(1)
		}
	}
}

