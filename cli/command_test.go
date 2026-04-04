package cli

import (
	"errors"
	"os"
	"strings"
	"testing"
)

func TestCommand_FlagSet(t *testing.T) {
	command := &Command{
		Name: "start",
	}

	fs := command.FlagSet()
	if fs == nil {
		t.Error("expected non-nil FlagSet")
	}

	fs2 := command.FlagSet()
	if fs != fs2 {
		t.Error("expected same FlagSet instance on repeated calls")
	}
}

func TestCommand_Parse(t *testing.T) {
	t.Run("Successful Parse", func(t *testing.T) {
		command := &Command{
			Name:      "start",
			Arguments: []string{"port"},
		}

		args, err := command.Parse([]string{"8080"})
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		if len(args) != 1 || args[0] != "8080" {
			t.Errorf("expected: [8080], got: %v", args)
		}
	})

	t.Run("Missing Arguments", func(t *testing.T) {
		command := &Command{
			Name:      "start",
			Arguments: []string{"port"},
		}

		_, err := command.Parse([]string{})
		if !errors.Is(err, ErrMissingArguments) {
			t.Errorf("expected: %v, got: %v", ErrMissingArguments, err)
		}
	})

	t.Run("With Flags", func(t *testing.T) {
		command := &Command{
			Name:      "start",
			Arguments: []string{"port"},
		}

		var verbose bool
		command.FlagSet().BoolVar(&verbose, "verbose", false, "verbose output")

		args, err := command.Parse([]string{"-verbose", "8080"})
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		if !verbose {
			t.Error("expected verbose flag to be true")
		}

		if len(args) != 1 || args[0] != "8080" {
			t.Errorf("expected: [8080], got: %v", args)
		}
	})
}

func TestCommand_PrintHelp(t *testing.T) {
	command := &Command{
		Name:      "start",
		Arguments: []string{"port", "host"},
	}

	group := &Group{
		Name: "serve",
	}

	output := captureStdout(t, func() {
		result := command.PrintHelp(group)
		if result != Success {
			t.Errorf("expected: %v, got: %v", Success, result)
		}
	})

	for _, expected := range []string{"serve", "start", "<port>", "<host>"} {
		if !strings.Contains(output, expected) {
			t.Errorf("expected output to contain: %v, got: %v", expected, output)
		}
	}
}

func TestCommand_PrintError(t *testing.T) {
	command := &Command{
		Name: "start",
	}

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	oldWriter := ErrorLog.Writer()
	ErrorLog.SetOutput(w)

	result := command.PrintError(errors.New("connection refused"))

	w.Close()
	ErrorLog.SetOutput(oldWriter)

	buf := make([]byte, 4096)
	n, _ := r.Read(buf)
	r.Close()

	if result != Failure {
		t.Errorf("expected: %v, got: %v", Failure, result)
	}

	output := string(buf[:n])
	if !strings.Contains(output, "connection refused") {
		t.Errorf("expected output to contain: connection refused, got: %v", output)
	}
}
