package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"io"
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

	result := m.Run()

	fmt.Println("....Cleaning up....")

	os.Remove(binName)
	os.Remove(fileName)

	os.Exit(result)
}

func TestTodoCLI(t *testing.T) {
	todo := "todo test number 1"
	todo_1 := "test todo from user input"

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cmdPath := filepath.Join(dir, binName)

	t.Run("Add new todo from args", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add", todo)
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Add new todo from user input", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add")
		// Access to the standard input of the current interactive / shell session through a pipe.
		cmdStdin, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}

		// Write / pipe the provided string to the standard input of the current interactive / shell session.
		io.WriteString(cmdStdin, todo_1)
		cmdStdin.Close()


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

		expected := fmt.Sprintf("   1: %s\n   2: %s\n", todo, todo_1)
		if expected != string(out) {
			t.Errorf("Expected: %q, Got: %q instead \n", expected, string(out))
		}
	})
}
