package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	"github.com/mycok/todo_list_cli/colors"
)

// Perform a fmt.Stringer interface satisfaction check.
var _ fmt.Stringer = (*List)(nil)

var (
	fmtWithDetail    = "%s%s%d: \t%s\t%s\t%s\t%s\n"
	fmtWithoutDetail = "%s%s%d: \t%s\t%s\n"
)

const prefix string = " " // prefix should be a singe space string.

type todo struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

// List represents a slice of todo type items.
type List []todo

// Add creates a new todo item and appends it to the list.
func (l *List) Add(task string) {
	t := todo{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}

	*l = append(*l, t)
}

// Complete marks a todo item of a specific index as complete by setting
// item.Done = true and item.CompletedAt = time.Now().
func (l *List) Complete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("Item %d does not exist", i)
	}

	// Adjusting user provided index to match zero based slice index.
	idx := i - 1
	ls[idx].Done = true
	ls[idx].CompletedAt = time.Now()

	return nil
}

// Delete a todo item of a specific index from the list.
func (l *List) Delete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("Item %d does not exist", i)
	}

	*l = append(ls[:i-1], ls[i:]...)

	return nil
}

// Save encodes the list as JSON and persists it to file using the provided file name.
func (l *List) Save(filename string) error {
	data, err := json.Marshal(l)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// Load opens the provided file name, decodes the JSON data and parses
// it into a todo list type.
func (l *List) Load(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return err
	}

	if len(data) == 0 {
		return nil
	}

	return json.Unmarshal(data, l)
}

// ListItems lists all todo list items.
func (l *List) ListItems(w io.Writer, details, completed bool) {
	msg := fmt.Sprintf(
		"%sno tasks available!. use %s-add%s %sto add tasks%s",
		colors.Yellow, colors.Cyan, colors.Reset, colors.Yellow, colors.Reset)

	tw := tabwriter.NewWriter(w, 0, 0, 1, ' ', tabwriter.Debug)

	switch {
	case details:
		if l.listItemDetails(completed) == "" {
			fmt.Fprintln(w, msg)
		} else {
			// Display the task headers using the same format as fmtWithDetail
			// variable replacing only the integer ID with a string ID title.
			fmt.Fprintf(
				tw, "%s%s%s \t%s\t%s\t%s\t%s\n",
				prefix, colors.Yellow, "ID", "TASK", "CREATED", "COMPLETED", colors.Reset,
			)

			// Display the list of items separated by new lines
			fmt.Fprintln(tw, l.listItemDetails(completed))

			tw.Flush()
		}

	default:
		if l.list(completed) == "" {
			fmt.Fprintln(w, msg)
		} else {
			// Display the task headers using the same format as fmtWithoutDetail
			// variable replacing only the integer ID with a string ID title.
			fmt.Fprintf(
				tw, "%s%s%s \t%s\t%s\n",
				prefix, colors.Yellow, "ID", "TASK", colors.Reset,
			)

			fmt.Fprintln(tw, l.list(completed))

			tw.Flush()
		}
	}
}

func (l *List) listItemDetails(completed bool) string {
	withoutCompleted := ""
	withCompleted := ""
	dateFormat := "02-01-2006 15:04"

	// Return an empty string in case there are no task items recorded.
	if len(*l) == 0 {
		return ""
	}

	for i, t := range *l {
		if t.Done {
			withCompleted += fmt.Sprintf(
				fmtWithDetail,
				prefix,
				colors.Green,
				i+1,
				t.Task,
				t.CreatedAt.Format(dateFormat),
				t.CompletedAt.Format(dateFormat),
				colors.Reset,
			)

		} else {
			formattedStr := fmt.Sprintf(
				fmtWithDetail,
				prefix,
				colors.White,
				i+1,
				t.Task,
				t.CreatedAt.Format(dateFormat),
				"",
				colors.Reset,
			)

			withCompleted += formattedStr
			withoutCompleted += formattedStr
		}
	}

	if completed {
		return withCompleted
	}

	return withoutCompleted
}

func (l *List) list(completed bool) string {
	// Return both completed and pending tasks.
	if completed {
		return l.String()
	}

	formatted := ""

	// Return an empty string in case there are no task items recorded.
	if len(*l) == 0 {
		return formatted
	}

	// return only pending tasks.
	for i, t := range *l {
		if !t.Done {
			formatted += fmt.Sprintf(
				fmtWithoutDetail,
				prefix,
				colors.White,
				i+1,
				t.Task,
				colors.Reset,
			)
		}
	}

	return formatted
}

// String returns a formatted todo.Task string. it returns both pending and
// completed tasks.
func (l *List) String() string {
	var taskColor colors.Color

	formatted := ""

	// Return an empty string in case there are no task items recorded.
	if len(*l) == 0 {
		return formatted
	}

	for i, t := range *l {
		if t.Done {
			taskColor = colors.Green
		} else {
			taskColor = colors.White
		}

		formatted += fmt.Sprintf(
			fmtWithoutDetail,
			prefix,
			taskColor,
			i+1,
			t.Task,
			colors.Reset,
		)
	}

	return formatted
}
