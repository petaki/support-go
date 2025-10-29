package vite

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/petaki/support-go/file"
)

// Vite type.
type Vite struct {
	publicDirectory string
	buildDirectory  string
	assetFS         fs.FS
	manifest        Manifest
}

// New function.
func New(publicDirectory, buildDirectory string) *Vite {
	v := new(Vite)
	v.publicDirectory = publicDirectory
	v.buildDirectory = buildDirectory

	return v
}

// NewWithFS function.
func NewWithFS(publicDirectory, buildDirectory string, assetFS fs.FS) *Vite {
	v := New(publicDirectory, buildDirectory)
	v.assetFS = assetFS

	return v
}

// IsRunningHot function.
func (v *Vite) IsRunningHot() bool {
	_, err := os.Stat(v.hotFile())

	return !os.IsNotExist(err)
}

// ManifestHash function.
func (v *Vite) ManifestHash() (string, error) {
	if v.IsRunningHot() {
		return "", nil
	}

	if v.assetFS != nil {
		hash, err := file.HashFromFS(v.manifestPath(), v.assetFS)
		if err != nil {
			return "", ErrManifestNotExist
		}

		return hash, nil
	}

	hash, err := file.Hash(v.manifestPath())
	if err != nil {
		return "", ErrManifestNotExist
	}

	return hash, nil
}

// Asset function.
func (v *Vite) Asset(asset string) (string, error) {
	if v.IsRunningHot() {
		return v.hotAsset(asset)
	}

	chunk, err := v.chunk(asset)
	if err != nil {
		return "", err
	}

	return filepath.Join(v.buildDirectory, chunk.File), nil
}

// CSS function.
func (v *Vite) CSS(asset string) ([]string, error) {
	if v.IsRunningHot() {
		return nil, nil
	}

	chunk, err := v.chunk(asset)
	if err != nil {
		return nil, err
	}

	var css []string

	for _, current := range chunk.CSS {
		css = append(css, filepath.Join(v.buildDirectory, current))
	}

	return css, nil
}

func (v *Vite) chunk(asset string) (*ManifestChunk, error) {
	err := v.ensureManifest()
	if err != nil {
		return nil, err
	}

	chunk, ok := v.manifest[asset]
	if !ok {
		return nil, fmt.Errorf("vite: unable to locate file: %v", asset)
	}

	return &chunk, nil
}

func (v *Vite) hotAsset(asset string) (string, error) {
	hotFileContent, err := os.ReadFile(v.hotFile())
	if err != nil {
		return "", err
	}

	devServerUrl := strings.TrimSpace(string(hotFileContent))

	return fmt.Sprintf("%s/%s", devServerUrl, asset), nil
}

func (v *Vite) ensureManifest() error {
	if v.manifest != nil {
		return nil
	}

	var manifestContent []byte

	if v.assetFS != nil {
		_, err := fs.Stat(v.assetFS, v.manifestPath())
		if errors.Is(err, fs.ErrNotExist) {
			return ErrManifestNotExist
		}

		manifestContent, err = fs.ReadFile(v.assetFS, v.manifestPath())
		if err != nil {
			return err
		}
	} else {
		_, err := os.Stat(v.manifestPath())
		if os.IsNotExist(err) {
			return ErrManifestNotExist
		}

		manifestContent, err = os.ReadFile(v.manifestPath())
		if err != nil {
			return err
		}
	}

	err := json.Unmarshal(manifestContent, &v.manifest)
	if err != nil {
		return err
	}

	return nil
}

func (v *Vite) hotFile() string {
	return filepath.Join(v.publicDirectory, "hot")
}

func (v *Vite) manifestPath() string {
	if v.assetFS != nil {
		return filepath.Join(v.buildDirectory, "manifest.json")
	}

	return filepath.Join(v.publicDirectory, v.buildDirectory, "manifest.json")
}
