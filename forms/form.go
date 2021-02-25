package forms

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// UsernameRegexp regexp.
var UsernameRegexp = regexp.MustCompile("^[a-zA-Z0-9\\.\\-_]+$")

// EmailRegexp regexp.
var EmailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Form type.
type Form struct {
	Data   map[string]interface{}
	Errors Bag
}

// New function.
func New(data map[string]interface{}) *Form {
	return &Form{
		data,
		map[string][]string{},
	}
}

// NewFromRequest function.
func NewFromRequest(w http.ResponseWriter, r *http.Request) (*Form, error) {
	var data map[string]interface{}

	err := DecodeBody(w, r, &data)
	if err != nil {
		return nil, err
	}

	return New(data), nil
}

// Required function.
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		if f.Data[field] == nil || strings.TrimSpace(fmt.Sprintf("%v", f.Data[field])) == "" {
			f.Errors.Add(field, "The field is required.")
		}
	}
}

// MatchesPattern function.
func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	if f.Data[field] == nil {
		return
	}

	switch value := f.Data[field].(type) {
	case string:
		if value == "" {
			return
		}

		if !pattern.MatchString(value) {
			f.Errors.Add(field, "The field format is invalid.")
		}
	}
}

// Min function.
func (f *Form) Min(field string, min float64) {
	if f.Data[field] == nil {
		return
	}

	switch value := f.Data[field].(type) {
	case float64:
		if value < min {
			f.Errors.Add(field, fmt.Sprintf("The field must be at least %v.", min))
		}
	case string:
		if value == "" {
			return
		}

		if float64(len(value)) < min {
			f.Errors.Add(field, fmt.Sprintf("The field must be at least %v characters.", min))
		}
	}
}

// Max function.
func (f *Form) Max(field string, max float64) {
	if f.Data[field] == nil {
		return
	}

	switch value := f.Data[field].(type) {
	case float64:
		if value > max {
			f.Errors.Add(field, fmt.Sprintf("The field may not be greater than %v.", max))
		}
	case string:
		if value == "" {
			return
		}

		if float64(len(value)) > max {
			f.Errors.Add(field, fmt.Sprintf("The field may not be greater than %v characters.", max))
		}
	}
}

// IsValid function.
func (f *Form) IsValid() bool {
	return len(f.Errors) == 0
}
