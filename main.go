package main

import (
	"fmt"
	"os"

	"github.com/mycok/todo_list_cli/cmd"
)

func main() {
	if err := cmd.Run(os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)

		return
	}
}
