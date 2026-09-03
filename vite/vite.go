package vite

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

// Vite type.
type Vite struct {
	mu              sync.Mutex
	publicDirectory string
	buildDirectory  string
	assetFS         fs.FS
	manifest        Manifest
	manifestHash    string
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

	_, err := v.loadManifest()
	if err != nil {
		return "", err
	}

	v.mu.Lock()
	defer v.mu.Unlock()

	return v.manifestHash, nil
}

// Asset function.
func (v *Vite) Asset(asset string) (string, error) {
	if v.IsRunningHot() {
		return v.hotAsset(asset)
	}

	_, chunk, err := v.manifestChunk(asset)
	if err != nil {
		return "", err
	}

	return v.assetPath(chunk.File), nil
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

	manifest, chunk, err := v.manifestChunk(asset)
	if err != nil {
		return nil, err
	}

	var css []string

	for _, current := range chunk.Imports {
		for _, file := range manifest[current].CSS {
			css = append(css, v.assetPath(file))
		}
	}

	for _, file := range chunk.CSS {
		css = append(css, v.assetPath(file))
	}

	return css, nil
}

// Preload function.
func (v *Vite) Preload(asset string) ([]string, error) {
	if v.IsRunningHot() {
		return nil, nil
	}

	manifest, chunk, err := v.manifestChunk(asset)
	if err != nil {
		return nil, err
	}

	var preload []string

	for _, current := range chunk.Imports {
		preload = append(preload, v.assetPath(manifest[current].File))
	}

	return preload, nil
}

func (v *Vite) assetPath(file string) string {
	return path.Join("/", v.buildDirectory, file)
}

func (v *Vite) hotAsset(asset string) (string, error) {
	hotFileContent, err := os.ReadFile(v.hotFile())
	if err != nil {
		return "", err
	}

	devServerURL := strings.TrimRight(strings.TrimSpace(string(hotFileContent)), "/")

	return fmt.Sprintf("%s/%s", devServerURL, asset), nil
}

func (v *Vite) hotFile() string {
	return filepath.Join(v.publicDirectory, "hot")
}

func (v *Vite) manifestChunk(asset string) (Manifest, *ManifestChunk, error) {
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

func (v *Vite) loadManifest() (Manifest, error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.manifest != nil {
		return v.manifest, nil
	}

	var manifestPath string
	var manifestContent []byte
	var err error

	if v.assetFS != nil {
		manifestPath = path.Join(v.buildDirectory, "manifest.json")
		manifestContent, err = fs.ReadFile(v.assetFS, manifestPath)
	} else {
		manifestPath = filepath.Join(v.publicDirectory, v.buildDirectory, "manifest.json")
		manifestContent, err = os.ReadFile(manifestPath)
	}

	if errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("%w: %q", ErrManifestNotExist, manifestPath)
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
	v.manifestHash = fmt.Sprintf("%x", md5.Sum(manifestContent))

	return manifest, nil
}
