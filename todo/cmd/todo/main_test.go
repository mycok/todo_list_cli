package main

import (
	"testing"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	binName = "todo"
	fileName = ".todo.json"
)

func TestMain(m *testing.M) {
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

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cmdPath := filepath.Join(dir, binName)

	t.Run("Add new todo", func(t *testing.T) {
		cmd := exec.Command(cmdPath, strings.Split(todo, " ")...)
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("List todos", func(t *testing.T) {
		cmd := exec.Command(cmdPath)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := todo + "\n"
		if expected != string(out) {
			t.Errorf("Expected: %q, Got: %q instead \n", expected, string(out))
		}
	})
}
