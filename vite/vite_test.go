package vite

import (
	"os"
	"path/filepath"
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
			if err != tt.expectErr {
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
			if err != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}

			if tt.manifest && hash == "" {
				t.Error("expected non-empty hash")
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
	"_shared-B7PI925R.js": {
		"file": "assets/shared-B7PI925R.js",
		"name": "shared",
		"css": ["assets/shared-ChJ_j-JJ.css"]
	},
	"views/foo.js": {
		"file": "assets/foo-BRBmoGS9.js",
		"name": "foo",
		"src": "views/foo.js",
		"isEntry": true,
		"imports": ["_shared-B7PI925R.js"],
		"css": ["assets/foo-5UjPuW-k.css"]
	}
}`

func createViteWithManifest(t *testing.T, manifest string) *Vite {
	t.Helper()

	return New(t.TempDir(), "build", fstest.MapFS{
		"build/manifest.json": &fstest.MapFile{Data: []byte(manifest)},
	})
}

func TestConcurrentManifestAccess(t *testing.T) {
	v := createViteWithManifest(t, testManifestWithImports)

	var wg sync.WaitGroup

	for range 50 {

		wg.Go(func() {

			asset, err := v.Asset("views/foo.js")
			if err != nil {
				t.Error(err)

				return
			}

			if asset != "/build/assets/foo-BRBmoGS9.js" {
				t.Errorf("unexpected asset: %v", asset)
			}

			css, err := v.CSS("views/foo.js")
			if err != nil {
				t.Error(err)

				return
			}

			if len(css) != 1 {
				t.Errorf("expected 1 css file, got: %v", css)
			}
		})
	}

	wg.Wait()
}
