package cli

import (
	"strings"
	"testing"
)

func TestGroup_PrintHelp(t *testing.T) {
	tests := []struct {
		name     string
		group    Group
		expected []string
	}{
		{
			"Single Command",
			Group{
				Name:  "serve",
				Usage: "Server commands",
				Commands: []*Command{
					{
						Name:  "start",
						Usage: "Start the server",
					},
				},
			},
			[]string{"start", "Start the server"},
		},
		{
			"Multiple Commands",
			Group{
				Name:  "serve",
				Usage: "Server commands",
				Commands: []*Command{
					{
						Name:  "start",
						Usage: "Start the server",
					},
					{
						Name:  "stop",
						Usage: "Stop the server",
					},
				},
			},
			[]string{"start", "Start the server", "stop", "Stop the server"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureStdout(t, func() {
				result := tt.group.PrintHelp()
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
