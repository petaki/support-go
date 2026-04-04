package forms

import "testing"

func TestUsernameRegexp(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid Alphanumeric", "john123", true},
		{"Valid With Dots", "john.doe", true},
		{"Valid With Hyphens", "john-doe", false},
		{"Valid With Underscores", "john_doe", true},
		{"Valid Mixed", "john.doe_123_test", true},
		{"Empty String", "", false},
		{"With Spaces", "john doe", false},
		{"With Special Chars", "john@doe", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UsernameRegexp.MatchString(tt.input)
			if tt.expected != got {
				t.Errorf("expected: %v, got: %v for input: %v", tt.expected, got, tt.input)
			}
		})
	}
}

func TestEmailRegexp(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid Simple", "test@example.com", true},
		{"Valid With Dots", "test.name@example.com", true},
		{"Valid With Plus", "test+tag@example.com", true},
		{"Valid Subdomain", "test@sub.example.com", true},
		{"Empty String", "", false},
		{"No At Sign", "testexample.com", false},
		{"No Domain", "test@", false},
		{"No Local Part", "@example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EmailRegexp.MatchString(tt.input)
			if tt.expected != got {
				t.Errorf("expected: %v, got: %v for input: %v", tt.expected, got, tt.input)
			}
		})
	}
}
