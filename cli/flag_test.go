package cli

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestCommand_String(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		t.Setenv("APP_ADDR", "")

		command := &Command{Name: "serve"}
		addr := command.String("addr", "APP_ADDR", ":4000", "Application Address")

		_, err := command.Parse([]string{})
		if err != nil {
			t.Fatal(err)
		}

		if *addr != ":4000" {
			t.Errorf("expected: :4000, got: %v", *addr)
		}
	})

	t.Run("Env", func(t *testing.T) {
		t.Setenv("APP_ADDR", ":5000")

		command := &Command{Name: "serve"}
		addr := command.String("addr", "APP_ADDR", ":4000", "Application Address")

		_, err := command.Parse([]string{})
		if err != nil {
			t.Fatal(err)
		}

		if *addr != ":5000" {
			t.Errorf("expected: :5000, got: %v", *addr)
		}
	})

	t.Run("Flag Overrides Env", func(t *testing.T) {
		t.Setenv("APP_ADDR", ":5000")

		command := &Command{Name: "serve"}
		addr := command.String("addr", "APP_ADDR", ":4000", "Application Address")

		_, err := command.Parse([]string{"-addr", ":6000"})
		if err != nil {
			t.Fatal(err)
		}

		if *addr != ":6000" {
			t.Errorf("expected: :6000, got: %v", *addr)
		}
	})

	t.Run("Usage", func(t *testing.T) {
		command := &Command{Name: "serve"}
		command.String("addr", "APP_ADDR", ":4000", "Application Address")

		var buf bytes.Buffer
		command.FlagSet().SetOutput(&buf)
		command.FlagSet().PrintDefaults()

		if !strings.Contains(buf.String(), "(env APP_ADDR)") {
			t.Errorf("expected usage to contain: (env APP_ADDR), got: %v", buf.String())
		}
	})
}

func TestCommand_Int(t *testing.T) {
	t.Run("Env", func(t *testing.T) {
		t.Setenv("COMMAND_HISTORY_LIMIT", "42")

		command := &Command{Name: "serve"}
		limit := command.Int("limit", "COMMAND_HISTORY_LIMIT", 100, "Command History Limit")

		_, err := command.Parse([]string{})
		if err != nil {
			t.Fatal(err)
		}

		if *limit != 42 {
			t.Errorf("expected: 42, got: %v", *limit)
		}
	})

	for _, arguments := range [][]string{{}, {"-limit", "5"}} {
		t.Run("Invalid Env", func(t *testing.T) {
			t.Setenv("COMMAND_HISTORY_LIMIT", "abc")

			command := &Command{Name: "serve"}
			command.Int("limit", "COMMAND_HISTORY_LIMIT", 100, "Command History Limit")

			_, err := command.Parse(arguments)
			if !errors.Is(err, ErrInvalidEnvValue) {
				t.Fatalf("expected: %v, got: %v", ErrInvalidEnvValue, err)
			}

			for _, expected := range []string{"COMMAND_HISTORY_LIMIT", "abc"} {
				if !strings.Contains(err.Error(), expected) {
					t.Errorf("expected error to contain: %v, got: %v", expected, err)
				}
			}
		})
	}
}

func TestCommand_Bool(t *testing.T) {
	t.Run("Env", func(t *testing.T) {
		t.Setenv("APP_DEBUG", "1")

		command := &Command{Name: "serve"}
		debug := command.Bool("debug", "APP_DEBUG", false, "Application Debug Mode")

		_, err := command.Parse([]string{})
		if err != nil {
			t.Fatal(err)
		}

		if !*debug {
			t.Error("expected debug to be true")
		}
	})

	t.Run("Invalid Env", func(t *testing.T) {
		t.Setenv("APP_DEBUG", "nope")

		command := &Command{Name: "serve"}
		command.Bool("debug", "APP_DEBUG", false, "Application Debug Mode")

		_, err := command.Parse([]string{})
		if !errors.Is(err, ErrInvalidEnvValue) {
			t.Errorf("expected: %v, got: %v", ErrInvalidEnvValue, err)
		}
	})
}

func TestCommand_Secret(t *testing.T) {
	t.Run("Env", func(t *testing.T) {
		t.Setenv("APP_KEY", "s3cret")

		command := &Command{Name: "serve"}
		key := command.Secret("key", "APP_KEY", "Application Key")

		_, err := command.Parse([]string{})
		if err != nil {
			t.Fatal(err)
		}

		if *key != "s3cret" {
			t.Errorf("expected: s3cret, got: %v", *key)
		}
	})

	t.Run("Flag Overrides Env", func(t *testing.T) {
		t.Setenv("APP_KEY", "s3cret")

		command := &Command{Name: "serve"}
		key := command.Secret("key", "APP_KEY", "Application Key")

		_, err := command.Parse([]string{"-key", "other"})
		if err != nil {
			t.Fatal(err)
		}

		if *key != "other" {
			t.Errorf("expected: other, got: %v", *key)
		}
	})

	t.Run("Usage", func(t *testing.T) {
		t.Setenv("APP_KEY", "s3cret")

		command := &Command{Name: "serve"}
		command.Secret("key", "APP_KEY", "Application Key")

		var buf bytes.Buffer
		command.FlagSet().SetOutput(&buf)
		command.FlagSet().PrintDefaults()

		if !strings.Contains(buf.String(), "$APP_KEY") {
			t.Errorf("expected usage to contain: $APP_KEY, got: %v", buf.String())
		}

		if strings.Contains(buf.String(), "s3cret") {
			t.Errorf("expected usage to hide the secret, got: %v", buf.String())
		}
	})
}
