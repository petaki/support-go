package mix

import "errors"

var (
	// ErrManifestNotExist error.
	ErrManifestNotExist = errors.New("mix: the mix manifest does not exist")
)
