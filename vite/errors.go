package vite

import "errors"

var (
	// ErrManifestNotExist error.
	ErrManifestNotExist = errors.New("vite: the manifest does not exist")

	// ErrAssetNotExist error.
	ErrAssetNotExist = errors.New("vite: the asset does not exist")
)
