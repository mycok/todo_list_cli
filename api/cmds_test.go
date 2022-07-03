package api

import (
	"errors"
	"io"
	"testing"
)

type testCmd struct{}

func (c testCmd) GetName() string  { return "test command" }
func (c testCmd) GetUsage() string { return "command help" }
func (c testCmd) Run(w io.Writer, args ...string) error {
	w.Write([]byte("this command has run without any errors"))

	return nil
}

func TestRegister(t *testing.T) {
	testCases := []struct {
		name        string
		cmd         Command
		expectedErr error
	}{
		{
			name:        "NewCommand",
			cmd:         testCmd{},
			expectedErr: nil,
		},
		{
			name:        "DuplicateCommand",
			cmd:         testCmd{},
			expectedErr: ErrDuplicateCmd,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := Register(tc.cmd)

			if tc.expectedErr != nil {
				if err == nil {
					t.Fatalf("Expected error: %s, but got nil", tc.expectedErr)
				}

				if !errors.Is(err, tc.expectedErr) {
					t.Errorf("Expected err: %s, but got: %s instead", tc.expectedErr, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("Unexpected err: %s", err)
			}

			// Get the added command from the commands map.
			cmd := Get(tc.cmd.GetName())
			if cmd.GetName() != tc.cmd.GetName() {
				t.Errorf("Expected command: %s, but got: %s instead", tc.cmd.GetName(), cmd.GetName())
			}
		})
	}

	t.Cleanup(func() {
		cmds := Commands()
		delete(cmds, testCases[0].cmd.GetName())
	})
}

func TestGet(t *testing.T) {
	c := testCmd{}

	testCases := []struct {
		name           string
		cmdName        string
		expectedResult Command
	}{
		{
			name:           "CmdFound",
			cmdName:        c.GetName(),
			expectedResult: c,
		},
		{
			name:           "CmdNotFound",
			cmdName:        "notFound",
			expectedResult: nil,
		},
	}

	if err := Register(c); err != nil {
		t.Fatal(err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := Get(tc.cmdName)
			if tc.expectedResult != cmd {
				t.Errorf("Expected result: %v, but got: %v instead", tc.expectedResult, cmd)
			}
		})
	}

	t.Cleanup(func() {
		cmds := Commands()
		delete(cmds, c.GetName())
	})
}
