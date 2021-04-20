package cli

import "testing"

func TestColor(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		got      string
	}{
		{"Red", "\u001B[31mtext\u001B[0m", Red("text")},
		{"Green", "\u001B[32mtext\u001B[0m", Green("text")},
		{"Yellow", "\u001B[33mtext\u001B[0m", Yellow("text")},
		{"Blue", "\u001B[34mtext\u001B[0m", Blue("text")},
		{"Purple", "\u001B[35mtext\u001B[0m", Purple("text")},
		{"Cyan", "\u001B[36mtext\u001B[0m", Cyan("text")},
		{"Gray", "\u001B[37mtext\u001B[0m", Gray("text")},
		{"White", "\u001B[97mtext\u001B[0m", White("text")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expected != tt.got {
				t.Errorf("expected: %v, got: %v", tt.expected, tt.got)
			}
		})
	}
}
