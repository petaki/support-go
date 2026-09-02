package vite

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"testing/fstest"
)

const testManifest = `{"resources/js/app.js":{"file":"assets/app-abc123.js","src":"resources/js/app.js","isEntry":true,"css":["assets/app-def456.css"]}}`

const testHotContent = "http://localhost:5173"

func createTempPublicDir(t *testing.T, hot bool, manifest bool) string {
	t.Helper()

	dir := t.TempDir()

	if hot {
		err := os.WriteFile(filepath.Join(dir, "hot"), []byte(testHotContent), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	if manifest {
		buildDir := filepath.Join(dir, "build")

		err := os.MkdirAll(buildDir, 0755)
		if err != nil {
			t.Fatal(err)
		}

		err = os.WriteFile(filepath.Join(buildDir, "manifest.json"), []byte(testManifest), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	return dir
}

func TestIsRunningHot(t *testing.T) {
	tests := []struct {
		name     string
		hot      bool
		expected bool
	}{
		{"Hot", true, true},
		{"Not Hot", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := createTempPublicDir(t, tt.hot, false)
			v := New(dir, "build")

			got := v.IsRunningHot()
			if tt.expected != got {
				t.Errorf("expected: %v, got: %v", tt.expected, got)
			}
		})
	}
}

func TestIsRunningHotWithUnreachableHotFile(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "public")

	err := os.WriteFile(dir, []byte("not a directory"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	v := New(dir, "build")

	if v.IsRunningHot() {
		t.Error("expected not hot when the hot file path is unreachable")
	}
}

func TestIsRunningHotWithHotDirectory(t *testing.T) {
	dir := t.TempDir()

	err := os.Mkdir(filepath.Join(dir, "hot"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	v := New(dir, "build")

	if v.IsRunningHot() {
		t.Error("expected not hot when the hot file is a directory")
	}
}

func TestManifestHash(t *testing.T) {
	tests := []struct {
		name      string
		hot       bool
		manifest  bool
		expectErr error
	}{
		{"With Manifest", false, true, nil},
		{"Hot Returns Empty", true, false, nil},
		{"No Manifest", false, false, ErrManifestNotExist},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := createTempPublicDir(t, tt.hot, tt.manifest)
			v := New(dir, "build")

			hash, err := v.ManifestHash()
			if !errors.Is(err, tt.expectErr) {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}

			if tt.hot && hash != "" {
				t.Errorf("expected empty hash, got: %v", hash)
			}

			if tt.manifest && !tt.hot && hash == "" {
				t.Error("expected non-empty hash")
			}
		})
	}
}

func TestManifestHashWithUnreachableManifest(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "public")

	err := os.WriteFile(dir, []byte("not a directory"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	v := New(dir, "build")

	_, err = v.ManifestHash()
	if err == nil {
		t.Fatal("expected an error")
	}

	if errors.Is(err, ErrManifestNotExist) {
		t.Errorf("expected the underlying error, got: %v", err)
	}
}

func TestManifestHashWithFS(t *testing.T) {
	tests := []struct {
		name      string
		manifest  bool
		expectErr error
	}{
		{"With Manifest", true, nil},
		{"No Manifest", false, ErrManifestNotExist},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			mapFS := fstest.MapFS{}

			if tt.manifest {
				mapFS["build/manifest.json"] = &fstest.MapFile{
					Data: []byte(testManifest),
				}
			}

			v := New(dir, "build", mapFS)

			hash, err := v.ManifestHash()
			if !errors.Is(err, tt.expectErr) {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}

			if tt.manifest && hash == "" {
				t.Error("expected non-empty hash")
			}
		})
	}
}

func TestAssetWithoutManifest(t *testing.T) {
	tests := []struct {
		name    string
		assetFS []fs.FS
	}{
		{"From Disk", nil},
		{"From FS", []fs.FS{fstest.MapFS{}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New(t.TempDir(), "build", tt.assetFS...)

			_, err := v.Asset("resources/js/app.js")
			if !errors.Is(err, ErrManifestNotExist) {
				t.Errorf("expected: %v, got: %v", ErrManifestNotExist, err)
			}
		})
	}
}

func TestAsset(t *testing.T) {
	tests := []struct {
		name     string
		hot      bool
		manifest bool
		expected string
	}{
		{"Hot", true, false, "http://localhost:5173/resources/js/app.js"},
		{"Production", false, true, "/build/assets/app-abc123.js"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := createTempPublicDir(t, tt.hot, tt.manifest)
			v := New(dir, "build")

			got, err := v.Asset("resources/js/app.js")
			if err != nil {
				t.Fatal(err)
			}

			if tt.expected != got {
				t.Errorf("expected: %v, got: %v", tt.expected, got)
			}
		})
	}
}

func TestAssetWithFS(t *testing.T) {
	dir := t.TempDir()
	mapFS := fstest.MapFS{
		"build/manifest.json": &fstest.MapFile{
			Data: []byte(testManifest),
		},
	}

	v := New(dir, "build", mapFS)

	got, err := v.Asset("resources/js/app.js")
	if err != nil {
		t.Fatal(err)
	}

	expected := "/build/assets/app-abc123.js"
	if expected != got {
		t.Errorf("expected: %v, got: %v", expected, got)
	}
}

func TestAssetWithTrailingSlashInHotFile(t *testing.T) {
	dir := t.TempDir()

	err := os.WriteFile(filepath.Join(dir, "hot"), []byte(testHotContent+"/\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	v := New(dir, "build")

	got, err := v.Asset("resources/js/app.js")
	if err != nil {
		t.Fatal(err)
	}

	expected := "http://localhost:5173/resources/js/app.js"
	if expected != got {
		t.Errorf("expected: %v, got: %v", expected, got)
	}
}

func TestAssetNotExist(t *testing.T) {
	v := createViteWithManifest(t, testManifestWithImports)

	_, err := v.Asset("resources/js/missing.js")
	if !errors.Is(err, ErrAssetNotExist) {
		t.Fatalf("expected: %v, got: %v", ErrAssetNotExist, err)
	}

	if !strings.Contains(err.Error(), "resources/js/missing.js") {
		t.Errorf("expected the asset name in: %v", err)
	}

	_, err = v.CSS("resources/js/missing.js")
	if !errors.Is(err, ErrAssetNotExist) {
		t.Errorf("expected: %v, got: %v", ErrAssetNotExist, err)
	}
}

func TestInertiaSSRURL(t *testing.T) {
	tests := []struct {
		name     string
		hot      bool
		expected string
	}{
		{"Hot", true, "http://localhost:5173/__inertia_ssr"},
		{"Not Hot", false, "http://127.0.0.1:13714/render"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := createTempPublicDir(t, tt.hot, false)
			v := New(dir, "build")

			got, err := v.InertiaSSRURL("http://127.0.0.1:13714/render")
			if err != nil {
				t.Fatal(err)
			}

			if tt.expected != got {
				t.Errorf("expected: %v, got: %v", tt.expected, got)
			}
		})
	}
}

func TestCSS(t *testing.T) {
	tests := []struct {
		name     string
		hot      bool
		manifest bool
		expected []string
	}{
		{"Hot", true, false, nil},
		{"Production", false, true, []string{"/build/assets/app-def456.css"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := createTempPublicDir(t, tt.hot, tt.manifest)
			v := New(dir, "build")

			got, err := v.CSS("resources/js/app.js")
			if err != nil {
				t.Fatal(err)
			}

			if tt.expected == nil && got != nil {
				t.Errorf("expected nil, got: %v", got)
			}

			if tt.expected != nil {
				if len(tt.expected) != len(got) {
					t.Errorf("expected: %v, got: %v", tt.expected, got)
				} else {
					for i, css := range tt.expected {
						if css != got[i] {
							t.Errorf("expected: %v, got: %v", css, got[i])
						}
					}
				}
			}
		})
	}
}

func TestCSSWithFS(t *testing.T) {
	dir := t.TempDir()
	mapFS := fstest.MapFS{
		"build/manifest.json": &fstest.MapFile{
			Data: []byte(testManifest),
		},
	}

	v := New(dir, "build", mapFS)

	got, err := v.CSS("resources/js/app.js")
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{"/build/assets/app-def456.css"}
	if len(expected) != len(got) {
		t.Errorf("expected: %v, got: %v", expected, got)
	} else {
		for i, css := range expected {
			if css != got[i] {
				t.Errorf("expected: %v, got: %v", css, got[i])
			}
		}
	}
}

const testManifestWithImports = `{
		"_shared-4Bu55YQl.js": {
			"file": "assets/shared-4Bu55YQl.js",
			"name": "shared",
			"css": [
				"assets/shared-Dn0TWFOf.css"
			]
		},
		"_shared-Dn0TWFOf.css": {
			"file": "assets/shared-Dn0TWFOf.css",
			"src": "_shared-Dn0TWFOf.css"
		},
		"resources/css/app.css": {
			"file": "assets/app-W1erjkBN.css",
			"name": "app",
			"names": [
				"app.css"
			],
			"src": "resources/css/app.css",
			"isEntry": true
		},
		"resources/js/app.js": {
			"file": "assets/app-_mEpyO-B.js",
			"name": "app",
			"src": "resources/js/app.js",
			"isEntry": true,
			"imports": [
				"_shared-4Bu55YQl.js"
			],
			"dynamicImports": [
				"resources/js/lazy.js"
			],
			"css": [
				"assets/app-W1erjkBN.css"
			]
		},
		"resources/js/lazy.js": {
			"file": "assets/lazy-CxGDFGtv.js",
			"name": "lazy",
			"src": "resources/js/lazy.js",
			"isDynamicEntry": true,
			"imports": [
				"_shared-4Bu55YQl.js"
			]
		},
		"resources/js/second.js": {
			"file": "assets/second-zqNUbXCm.js",
			"name": "second",
			"src": "resources/js/second.js",
			"isEntry": true,
			"imports": [
				"_shared-4Bu55YQl.js"
			]
		}
	}`

func createViteWithManifest(t *testing.T, manifest string) *Vite {
	t.Helper()

	return New(t.TempDir(), "build", fstest.MapFS{
		"build/manifest.json": &fstest.MapFile{Data: []byte(manifest)},
	})
}

func TestCSSWithImports(t *testing.T) {
	tests := []struct {
		name     string
		asset    string
		expected []string
	}{
		{
			"Imported Chunk CSS First",
			"resources/js/app.js",
			[]string{"/build/assets/shared-Dn0TWFOf.css", "/build/assets/app-W1erjkBN.css"},
		},
		{
			"Only Imported Chunk CSS",
			"resources/js/second.js",
			[]string{"/build/assets/shared-Dn0TWFOf.css"},
		},
		{
			"CSS Entrypoint",
			"resources/css/app.css",
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := createViteWithManifest(t, testManifestWithImports)

			got, err := v.CSS(tt.asset)
			if err != nil {
				t.Fatal(err)
			}

			if len(tt.expected) != len(got) {
				t.Fatalf("expected: %v, got: %v", tt.expected, got)
			}

			for i, css := range tt.expected {
				if css != got[i] {
					t.Errorf("expected: %v, got: %v", css, got[i])
				}
			}
		})
	}
}

func TestConcurrentManifestAccess(t *testing.T) {
	v := createViteWithManifest(t, testManifestWithImports)

	var wg sync.WaitGroup

	for range 50 {

		wg.Go(func() {

			asset, err := v.Asset("resources/js/app.js")
			if err != nil {
				t.Error(err)

				return
			}

			if asset != "/build/assets/app-_mEpyO-B.js" {
				t.Errorf("unexpected asset: %v", asset)
			}

			css, err := v.CSS("resources/js/app.js")
			if err != nil {
				t.Error(err)

				return
			}

			if len(css) != 2 {
				t.Errorf("expected 2 css files, got: %v", css)
			}
		})
	}

	wg.Wait()
}
