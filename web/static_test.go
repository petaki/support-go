package web

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestStatic(t *testing.T) {
	t.Run("Embedded", func(t *testing.T) {
		files := fstest.MapFS{
			"favicon.ico": &fstest.MapFile{Data: []byte("icon")},
		}

		rr := httptest.NewRecorder()
		Static(false, "", files).ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/favicon.ico", nil))

		if rr.Code != http.StatusOK || rr.Body.String() != "icon" {
			t.Errorf("expected: 200 icon, got: %v %v", rr.Code, rr.Body.String())
		}
	})

	t.Run("Debug", func(t *testing.T) {
		directory := t.TempDir()

		err := os.WriteFile(filepath.Join(directory, "favicon.ico"), []byte("disk"), 0o644)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		Static(true, directory, nil).ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/favicon.ico", nil))

		if rr.Code != http.StatusOK || rr.Body.String() != "disk" {
			t.Errorf("expected: 200 disk, got: %v %v", rr.Code, rr.Body.String())
		}
	})
}
