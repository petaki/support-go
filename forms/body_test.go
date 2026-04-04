package forms

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecodeBody(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		body        string
		expectErr   bool
		errStatus   int
		errMsg      string
	}{
		{
			"Valid JSON",
			"application/json",
			`{"name":"John"}`,
			false, 0, "",
		},
		{
			"No Content-Type",
			"",
			`{"name":"John"}`,
			false, 0, "",
		},
		{
			"Invalid Content-Type",
			"text/plain",
			`{"name":"John"}`,
			true,
			http.StatusUnsupportedMediaType,
			"Content-Type header contains an invalid value",
		},
		{
			"Empty Body",
			"application/json",
			"",
			true,
			http.StatusBadRequest,
			"Request body must not be empty",
		},
		{
			"Malformed JSON",
			"application/json",
			`{"name":}`,
			true,
			http.StatusBadRequest,
			"Request body contains badly-formed JSON",
		},
		{
			"Multiple JSON Objects",
			"application/json",
			`{"name":"John"}{"name":"Jane"}`,
			true,
			http.StatusBadRequest,
			"Request body must only contain a single JSON object",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			if tt.contentType != "" {
				r.Header.Set("Content-Type", tt.contentType)
			}

			w := httptest.NewRecorder()

			var dst map[string]any
			err := DecodeBody(w, r, &dst)

			if tt.expectErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				formErr, ok := err.(*Error)
				if !ok {
					t.Fatalf("expected *Error, got: %T", err)
				}

				if tt.errStatus != formErr.Status {
					t.Errorf("expected status: %v, got: %v", tt.errStatus, formErr.Status)
				}

				if !strings.Contains(formErr.Msg, tt.errMsg) {
					t.Errorf("expected message containing: %v, got: %v", tt.errMsg, formErr.Msg)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestDecodeBodyTypeMismatch(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"age":"not a number"}`))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	var dst struct {
		Age int `json:"age"`
	}

	err := DecodeBody(w, r, &dst)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	formErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("expected *Error, got: %T", err)
	}

	if formErr.Status != http.StatusBadRequest {
		t.Errorf("expected status: %v, got: %v", http.StatusBadRequest, formErr.Status)
	}
}

func TestBagAdd(t *testing.T) {
	b := Bag{}

	b.Add("name", "is required")
	if len(b["name"]) != 1 {
		t.Errorf("expected 1 error, got: %v", len(b["name"]))
	}

	b.Add("name", "is too short")
	if len(b["name"]) != 2 {
		t.Errorf("expected 2 errors, got: %v", len(b["name"]))
	}

	b.Add("email", "is required")
	if len(b["email"]) != 1 {
		t.Errorf("expected 1 error for email, got: %v", len(b["email"]))
	}
}

func TestErrorError(t *testing.T) {
	e := &Error{
		Status: http.StatusBadRequest,
		Msg:    "test error message",
	}

	if e.Error() != "test error message" {
		t.Errorf("expected: test error message, got: %v", e.Error())
	}
}
