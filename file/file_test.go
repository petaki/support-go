package file

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestHash(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		expectErr bool
	}{
		{"Valid File", "hello world", false},
		{"Empty File", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "testfile")

			err := os.WriteFile(path, []byte(tt.content), 0644)
			if err != nil {
				t.Fatal(err)
			}

			hash, err := Hash(path)
			if tt.expectErr && err == nil {
				t.Error("expected error, got nil")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.expectErr && hash == "" {
				t.Error("expected non-empty hash")
			}
		})
	}
}

func TestHashNonExistentFile(t *testing.T) {
	_, err := Hash("/nonexistent/path/file.txt")
	if err == nil {
		t.Error("expected error for non-existent file, got nil")
	}
}

func TestHashConsistency(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "testfile")

	err := os.WriteFile(path, []byte("consistent content"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	hash1, err := Hash(path)
	if err != nil {
		t.Fatal(err)
	}

	hash2, err := Hash(path)
	if err != nil {
		t.Fatal(err)
	}

	if hash1 != hash2 {
		t.Errorf("expected consistent hashes, got: %v and %v", hash1, hash2)
	}
}

func TestHashFromFS(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		expectErr bool
	}{
		{"Valid File", "hello world", false},
		{"Empty File", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapFS := fstest.MapFS{
				"testfile": &fstest.MapFile{
					Data: []byte(tt.content),
				},
			}

			hash, err := HashFromFS("testfile", mapFS)
			if tt.expectErr && err == nil {
				t.Error("expected error, got nil")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.expectErr && hash == "" {
				t.Error("expected non-empty hash")
			}
		})
	}
}

func TestHashFromFSNonExistentFile(t *testing.T) {
	mapFS := fstest.MapFS{}

	_, err := HashFromFS("nonexistent", mapFS)
	if err == nil {
		t.Error("expected error for non-existent file, got nil")
	}
}

func TestHashAndHashFromFSMatch(t *testing.T) {
	content := "matching content"

	dir := t.TempDir()
	path := filepath.Join(dir, "testfile")

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	hash1, err := Hash(path)
	if err != nil {
		t.Fatal(err)
	}

	mapFS := fstest.MapFS{
		"testfile": &fstest.MapFile{
			Data: []byte(content),
		},
	}

	hash2, err := HashFromFS("testfile", mapFS)
	if err != nil {
		t.Fatal(err)
	}

	if hash1 != hash2 {
		t.Errorf("Hash and HashFromFS produced different results: %v vs %v", hash1, hash2)
	}
}
