package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadDir(t *testing.T) {
	dir := t.TempDir()
	filesToCreate := map[string]string{
		"VAR1": "value1",
		"VAR2": "value2\nwith newline",
		"VAR3": "",
	}

	expectedValues := map[string]string{
		"VAR1": "value1",
		"VAR2": "value2",
		"VAR3": "",
	}

	for name, content := range filesToCreate {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("Error creating file %s: %v", path, err)
		}
	}

	env, err := ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir failed: %v", err)
	}

	for key, expected := range expectedValues {
		value, ok := env[key]
		if !ok {
			t.Errorf("Variable %s missing in result", key)
			continue
		}
		if value.Value != expected {
			t.Errorf("For %s: expected %q, got %q", key, expected, value.Value)
		}
	}
}

func TestReadDir_ZeroBytesFile(t *testing.T) {
	dir := t.TempDir()

	filename := "VAR_ZERO"
	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, []byte{0, 0, 0}, 0o644); err != nil {
		t.Fatalf("Error creating file: %v", err)
	}

	env, err := ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir failed: %v", err)
	}

	val := env[filename]
	expected := ""
	if val.Value != expected {
		t.Errorf("Expected %q, got %q", expected, val.Value)
	}
}

func TestReadDirFailOpenDir(t *testing.T) {
	_, err := ReadDir("/path/to/nonexistent/dir")
	if err == nil {
		t.Fatal("Expected error for non-existent directory, got nil")
	}
}

func TestReadDirFileReadError(t *testing.T) {
	dir := t.TempDir()

	filename := "VAR"
	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, []byte("val"), 0o644); err != nil {
		t.Fatalf("Error writing file: %v", err)
	}
	os.Chmod(path, 0o000)

	defer os.Chmod(path, 0o644)

	_, err := ReadDir(dir)
	if err == nil {
		t.Fatal("Expected error on file read, got nil")
	}
}
