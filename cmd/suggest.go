package cmd

import (
	"fmt"
	"io"

	"github.com/agnivade/levenshtein"

	"github.com/mycok/todo_list_cli/api"
	"github.com/mycok/todo_list_cli/colors"
)

var suggest = api.Cmd{
	Action: func(w io.Writer, args ...string) error {
		var list []string

		for _, c := range api.Commands() {
			name := c.GetName()
			d := levenshtein.ComputeDistance(name, args[0])

			if d <= 3 {
				list = append(list, name)
			}
		}

		fmt.Fprintf(w, "Invalid command %s%q%s", colors.Red, args[0], colors.Reset)

		if len(list) == 0 {
			fmt.Fprintln(w)

			return nil
		}

		fmt.Fprint(w, " Did you mean? ")

		for i, v := range list {
			if i > 0 {
				fmt.Fprint(w, ", ")
			}

			fmt.Fprintf(w, "%s%s%s", colors.Green, v, colors.Reset)
		}

		fmt.Fprintln(w)

		return nil
	},
}
