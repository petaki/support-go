package vite

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/petaki/support-go/file"
)

// Vite type.
type Vite struct {
	mu              sync.Mutex
	publicDirectory string
	buildDirectory  string
	assetFS         fs.FS
	manifest        Manifest
}

// New function.
func New(publicDirectory, buildDirectory string, assetFS ...fs.FS) *Vite {
	v := new(Vite)
	v.publicDirectory = publicDirectory
	v.buildDirectory = buildDirectory

	if len(assetFS) > 0 && assetFS[0] != nil {
		v.assetFS = assetFS[0]
	}

	return v
}

// IsRunningHot function.
func (v *Vite) IsRunningHot() bool {
	info, err := os.Stat(v.hotFile())

	return err == nil && info.Mode().IsRegular()
}

// ManifestHash function.
func (v *Vite) ManifestHash() (string, error) {
	if v.IsRunningHot() {
		return "", nil
	}

	var hash string
	var err error

	if v.assetFS != nil {
		hash, err = file.HashFromFS(v.manifestPath(), v.assetFS)
	} else {
		hash, err = file.Hash(v.manifestPath())
	}

	if errors.Is(err, fs.ErrNotExist) {
		return "", fmt.Errorf("%w: %q", ErrManifestNotExist, v.manifestPath())
	}

	if err != nil {
		return "", err
	}

	return hash, nil
}

// Asset function.
func (v *Vite) Asset(asset string) (string, error) {
	if v.IsRunningHot() {
		return v.hotAsset(asset)
	}

	_, chunk, err := v.chunk(asset)
	if err != nil {
		return "", err
	}

	return path.Join("/", v.buildDirectory, chunk.File), nil
}

// InertiaSSRURL function.
func (v *Vite) InertiaSSRURL(defaultURL string) (string, error) {
	if v.IsRunningHot() {
		return v.hotAsset("__inertia_ssr")
	}

	return defaultURL, nil
}

// CSS function.
func (v *Vite) CSS(asset string) ([]string, error) {
	if v.IsRunningHot() {
		return nil, nil
	}

	manifest, chunk, err := v.chunk(asset)
	if err != nil {
		return nil, err
	}

	var css []string

	for _, current := range chunk.Imports {
		css = v.appendAssets(css, manifest[current].CSS)
	}

	return v.appendAssets(css, chunk.CSS), nil
}

func (v *Vite) chunk(asset string) (Manifest, *ManifestChunk, error) {
	manifest, err := v.loadManifest()
	if err != nil {
		return nil, nil, err
	}

	chunk, ok := manifest[asset]
	if !ok {
		return nil, nil, fmt.Errorf("%w: %q", ErrAssetNotExist, asset)
	}

	return manifest, &chunk, nil
}

func (v *Vite) appendAssets(assets []string, files []string) []string {
	for _, current := range files {
		assets = append(assets, path.Join("/", v.buildDirectory, current))
	}

	return assets
}

func (v *Vite) hotAsset(asset string) (string, error) {
	hotFileContent, err := os.ReadFile(v.hotFile())
	if err != nil {
		return "", err
	}

	devServerURL := strings.TrimRight(strings.TrimSpace(string(hotFileContent)), "/")

	return fmt.Sprintf("%s/%s", devServerURL, asset), nil
}

func (v *Vite) loadManifest() (Manifest, error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.manifest != nil {
		return v.manifest, nil
	}

	var manifestContent []byte
	var err error

	if v.assetFS != nil {
		manifestContent, err = fs.ReadFile(v.assetFS, v.manifestPath())
	} else {
		manifestContent, err = os.ReadFile(v.manifestPath())
	}

	if errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("%w: %q", ErrManifestNotExist, v.manifestPath())
	}

	if err != nil {
		return nil, err
	}

	var manifest Manifest

	err = json.Unmarshal(manifestContent, &manifest)
	if err != nil {
		return nil, err
	}

	v.manifest = manifest

	return manifest, nil
}

func (v *Vite) hotFile() string {
	return filepath.Join(v.publicDirectory, "hot")
}

func (v *Vite) manifestPath() string {
	if v.assetFS != nil {
		return path.Join(v.buildDirectory, "manifest.json")
	}

	return filepath.Join(v.publicDirectory, v.buildDirectory, "manifest.json")
}
