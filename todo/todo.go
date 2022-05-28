package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// Perform a fmt.Stringer interface satisfaction check.
var _ fmt.Stringer = (*List)(nil)

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
	js, err := json.Marshal(l)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, js, 0644)
}

// Load opens the provided file name, decodes the JSON data and parses it into a todo list type.
func (l *List) Load(filename string) error {
	fileData, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return err
	}

	if len(fileData) == 0 {
		return nil
	}

	return json.Unmarshal(fileData, l)
}

// ListItems lists all todo list items.
func (l *List) ListItems(details, completed *bool) {
	if *details {
		fmt.Print(l.listItemDetails(*completed))
	} else {
		fmt.Print(l.list(*completed))
	}
}

func (l *List) listItemDetails(completed bool) string {
	var formattedStr string

	formatted := ""
	allFormatted := ""
	dateFormat := "02-01-2006 15:04"

	for i, t := range *l {
		prefix := "   "

		if t.Done {
			prefix = "✅ "
			formattedStr = fmt.Sprintf(
				"%s%d: %s - created %s - completed  %s\n",
				prefix,
				i+1,
				t.Task,
				t.CreatedAt.Format(dateFormat),
				t.CompletedAt.Format(dateFormat),
			)

			allFormatted += formattedStr
		} else {
			formattedStr = fmt.Sprintf(
				"%s%d: %s - created %s\n",
				prefix,
				i+1,
				t.Task,
				t.CreatedAt.Format(dateFormat),
			)

			allFormatted += formattedStr
			formatted += formattedStr
		}
	}

	if completed {
		return allFormatted
	}

	return formatted
}

func (l *List) list(completed bool) string {
	if completed {
		return l.String()
	}

	formatted := ""

	for i, t := range *l {
		prefix := "   "

		if !t.Done {
			formatted += fmt.Sprintf("%s%d: %s\n", prefix, i+1, t.Task)
		}
	}

	return formatted
}

// String returns a formatted todo.Task string.
func (l *List) String() string {
	formatted := ""

	for i, t := range *l {
		prefix := "   "

		if t.Done {
			prefix = "✅ "
		}

		formatted += fmt.Sprintf("%s%d: %s\n", prefix, i+1, t.Task)
	}

	return formatted
}
