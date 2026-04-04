package forms

import (
	"regexp"
	"testing"
)

func TestNew(t *testing.T) {
	data := map[string]any{
		"name": "John",
	}

	f := New(data)
	if f.Data["name"] != "John" {
		t.Errorf("expected: John, got: %v", f.Data["name"])
	}

	if !f.IsValid() {
		t.Error("expected form to be valid")
	}
}

func TestIsValid(t *testing.T) {
	f := New(map[string]any{})

	if !f.IsValid() {
		t.Error("expected form to be valid")
	}

	f.Errors.Add("field", "error")

	if f.IsValid() {
		t.Error("expected form to be invalid")
	}
}

func TestRequired(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]any
		field   string
		isValid bool
	}{
		{"Valid String", map[string]any{"name": "John"}, "name", true},
		{"Missing Field", map[string]any{}, "name", false},
		{"Nil Field", map[string]any{"name": nil}, "name", false},
		{"Empty String", map[string]any{"name": ""}, "name", false},
		{"Whitespace Only", map[string]any{"name": "   "}, "name", false},
		{"Number Value", map[string]any{"age": float64(25)}, "age", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := New(tt.data)
			f.Required(tt.field)

			if tt.isValid != f.IsValid() {
				t.Errorf("expected valid: %v, got: %v", tt.isValid, f.IsValid())
			}
		})
	}
}

func TestRequiredMultipleFields(t *testing.T) {
	f := New(map[string]any{
		"name": "John",
	})

	f.Required("name", "email")

	if _, ok := f.Errors["name"]; ok {
		t.Error("expected no error for name")
	}

	if _, ok := f.Errors["email"]; !ok {
		t.Error("expected error for email")
	}
}

func TestMatchesPattern(t *testing.T) {
	pattern := regexp.MustCompile(`^[a-z]+$`)

	tests := []struct {
		name    string
		data    map[string]any
		isValid bool
	}{
		{"Matching Pattern", map[string]any{"field": "abc"}, true},
		{"Non-Matching Pattern", map[string]any{"field": "ABC123"}, false},
		{"Nil Field", map[string]any{}, true},
		{"Empty String", map[string]any{"field": ""}, true},
		{"Non-String Type", map[string]any{"field": float64(123)}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := New(tt.data)
			f.MatchesPattern("field", pattern)

			if tt.isValid != f.IsValid() {
				t.Errorf("expected valid: %v, got: %v", tt.isValid, f.IsValid())
			}
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]any
		min     float64
		isValid bool
	}{
		{"Float64 Below Min", map[string]any{"field": float64(3)}, 5, false},
		{"Float64 At Min", map[string]any{"field": float64(5)}, 5, true},
		{"Float64 Above Min", map[string]any{"field": float64(10)}, 5, true},
		{"String Shorter Than Min", map[string]any{"field": "ab"}, 3, false},
		{"String At Min Length", map[string]any{"field": "abc"}, 3, true},
		{"String Above Min Length", map[string]any{"field": "abcdef"}, 3, true},
		{"Nil Field", map[string]any{}, 5, true},
		{"Empty String", map[string]any{"field": ""}, 5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := New(tt.data)
			f.Min("field", tt.min)

			if tt.isValid != f.IsValid() {
				t.Errorf("expected valid: %v, got: %v", tt.isValid, f.IsValid())
			}
		})
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]any
		max     float64
		isValid bool
	}{
		{"Float64 Above Max", map[string]any{"field": float64(10)}, 5, false},
		{"Float64 At Max", map[string]any{"field": float64(5)}, 5, true},
		{"Float64 Below Max", map[string]any{"field": float64(3)}, 5, true},
		{"String Longer Than Max", map[string]any{"field": "abcdef"}, 3, false},
		{"String At Max Length", map[string]any{"field": "abc"}, 3, true},
		{"String Below Max Length", map[string]any{"field": "ab"}, 3, true},
		{"Nil Field", map[string]any{}, 5, true},
		{"Empty String", map[string]any{"field": ""}, 5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := New(tt.data)
			f.Max("field", tt.max)

			if tt.isValid != f.IsValid() {
				t.Errorf("expected valid: %v, got: %v", tt.isValid, f.IsValid())
			}
		})
	}
}
