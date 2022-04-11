package todo_test

import (
	"os"
	"testing"

	"github.com/myok/todo_list_cli/todo"
)

func TestTodoListFunctionality(t *testing.T) {
	l := todo.List{}
	todoName := "test mark as complete functionality"

	l.Add("test mark as complete functionality")

	if len(l) != 1 {
		t.Errorf("Expected: 1 item in the list but Got: %d", len(l))
	}

	if l[0].Task != todoName {
		t.Errorf("Expected: %q, Got: %q", todoName, l[0].Task)
	}

	l.Complete(1)
	if !l[0].Done {
		t.Errorf("Expected %q to be true, Got: %t", l[0].Task, l[0].Done)
	}

	l.Delete(1)
	if len(l) > 0 {
		t.Errorf("Expected list to be empty, Got: %+#v", l[0])
	}
}

func TestSaveGetTodoList(t *testing.T) {
	l1 := todo.List{}
	l2 := todo.List{}

	todoName := "new task"
	l1.Add(todoName)

	if l1[0].Task != todoName {
		t.Errorf("Expected %q, Got: %q", todoName, l1[0].Task)
	}

	tf, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("Error creating temp file: %s", err)
	}

	defer os.Remove(tf.Name())

	if err := l1.Save(tf.Name()); err != nil {
		t.Fatalf("Error saving todo list to file: %s", err)
	}

	if err := l2.Get(tf.Name()); err != nil {
		t.Fatalf("Error retrieving todo list from file: %s", err)
	}

	if len(l2) != 1 {
		t.Errorf("Expected todo list to contain at least one todo item: Got: %d", len(l2))
	}

	if l1[0].Task != l2[0].Task {
		t.Errorf("Task %q should match %q task", l1[0].Task, l2[0].Task)
	}
}
