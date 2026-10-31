package cli

import "errors"

var (
	// ErrMissingArguments error.
	ErrMissingArguments = errors.New("cli: missing arguments")

	// ErrInvalidEnvValue error.
	ErrInvalidEnvValue = errors.New("cli: invalid env value")
)
