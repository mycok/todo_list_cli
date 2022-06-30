package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"text/tabwriter"
	"time"

	"github.com/mycok/todo_list_cli/colors"
)

var (
	binName  = "todo_cli"
	fileName = ".todo.json"
)

func TestMain(m *testing.M) {
	if os.Getenv("TODO_FILENAME") != "" {
		fileName = os.Getenv("TODO_FILENAME")
	}

	fmt.Println("....building tool....")

	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	cmd := exec.Command("go", "build", "-o", binName)
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot build tool %s: %s", binName, err)

		os.Exit(1)
	}

	fmt.Println("....Running tests....")

	code := m.Run()

	fmt.Println("....Cleaning up....")

	os.Remove(binName)
	os.Remove(fileName)

	os.Exit(code)
}

func TestTodoCLI(t *testing.T) {
	var createdAt string

	todo := "todo test number 1"
	todo1 := "test todo from user input"

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Get the path for the compiled app binary.
	cmdPath := filepath.Join(dir, binName)

	t.Run("Add new todo from args", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add", todo)
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Add new todo from user input", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add")
		// Access to the standard input of the current interactive / shell session
		// through a pipe.
		cmdStdin, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}

		// Write / pipe the provided string to the standard input of the current
		// interactive / shell session.
		io.WriteString(cmdStdin, todo1)
		cmdStdin.Close()

		createdAt = time.Now().Format("02-01-2006 15:04")

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("List todos", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list")
		// Access output from both stdOut and stdErr of the current
		// interactive / shell session.
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		var outBuf bytes.Buffer

		tw := tabwriter.NewWriter(&outBuf, 0, 0, 1, ' ', tabwriter.Debug)

		prefix := " "

		fmt.Fprintf(tw, "%s%s%s \t%s\t%s\n", prefix, colors.Yellow, "ID", "TASK", colors.Reset)
		fmt.Fprintf(tw, "%s%s%d: \t%s\t%s\n", prefix, colors.White, 1, todo, colors.Reset)
		fmt.Fprintf(tw, "%s%s%d: \t%s\t%s\n", prefix, colors.White, 2, todo1, colors.Reset)
		fmt.Fprintln(tw)

		tw.Flush()

		expectedTodoList := outBuf.String()

		if expectedTodoList != string(out) {
			t.Errorf("Expected: %q, Got: %q instead \n", expectedTodoList, string(out))
		}
	})

	t.Run("delete todo", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-del", "1")
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("List todos", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list")
		// Access output from both stdOut and stdErr of the current interactive / shell session.
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		var outBuf bytes.Buffer

		tw := tabwriter.NewWriter(&outBuf, 0, 0, 1, ' ', tabwriter.Debug)

		prefix := " "

		fmt.Fprintf(tw, "%s%s%s \t%s\t%s\n", prefix, colors.Yellow, "ID", "TASK", colors.Reset)
		fmt.Fprintf(tw, "%s%s%d: \t%s\t%s\n", prefix, colors.White, 1, todo1, colors.Reset)
		fmt.Fprintln(tw)

		tw.Flush()

		expectedTodoList := outBuf.String()

		if expectedTodoList != string(out) {
			t.Errorf("Expected: %q, Got: %q instead \n", expectedTodoList, string(out))
		}
	})

	t.Run("List todos with details", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list", "-details")
		// Access output from both stdOut and stdErr of the current interactive / shell session.
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		var outBuf bytes.Buffer

		tw := tabwriter.NewWriter(&outBuf, 0, 0, 1, ' ', tabwriter.Debug)

		prefix := " "

		fmt.Fprintf(
			tw, "%s%s%s \t%s\t%s\t%s\t%s\n",
			prefix, colors.Yellow, "ID", "TASK", "CREATED", "COMPLETED", colors.Reset,
		)

		fmt.Fprintf(
			tw, "%s%s%d: \t%s\t%s\t%s\t%s\n",
			prefix, colors.White, 1, todo1, createdAt, "", colors.Reset,
		)

		fmt.Fprintln(tw)

		tw.Flush()

		expectedTodoList := outBuf.String()

		if expectedTodoList != string(out) {
			t.Errorf("Expected: %q, Got: %q instead \n", expectedTodoList, string(out))
		}
	})

}
