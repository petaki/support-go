package cli

import (
	"os"
	"strings"
	"testing"
)

func TestApp_PrintHelp(t *testing.T) {
	tests := []struct {
		name     string
		app      App
		expected []string
	}{
		{
			"Single Group",
			App{
				Name:    "myapp",
				Version: "1.0.0",
				Groups: []*Group{
					{
						Name:  "serve",
						Usage: "Server commands",
					},
				},
			},
			[]string{"myapp", "1.0.0", "serve", "Server commands"},
		},
		{
			"Multiple Groups",
			App{
				Name:    "myapp",
				Version: "2.0.0",
				Groups: []*Group{
					{
						Name:  "serve",
						Usage: "Server commands",
					},
					{
						Name:  "migrate",
						Usage: "Migration commands",
					},
				},
			},
			[]string{"myapp", "2.0.0", "serve", "Server commands", "migrate", "Migration commands"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureStdout(t, func() {
				result := tt.app.PrintHelp()
				if result != Success {
					t.Errorf("expected: %v, got: %v", Success, result)
				}
			})

			for _, expected := range tt.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("expected output to contain: %v, got: %v", expected, output)
				}
			}
		})
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	oldStdout := os.Stdout
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = oldStdout

	buf := make([]byte, 4096)
	n, _ := r.Read(buf)
	r.Close()

	return string(buf[:n])
}
