package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestReadDir(t *testing.T) {
	dir := t.TempDir()

	// Создаем файлы с разными содержимым
	files := map[string]string{
		"VAR1": "value1",
		"VAR2": "value2\nwith newline",
		"VAR3": "", // пустой файл
	}
	for name, content := range files {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("Error creating file %s: %v", path, err)
		}
	}

	env, err := ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir failed: %v", err)
	}

	// Проверка переменных
	for key, expected := range files {
		value := env[key]
		if value != expected {
			t.Errorf("For %s: expected %q, got %q", key, expected, value)
		}
	}
}

func TestReadDir_ZeroBytesFile(t *testing.T) {
	dir := t.TempDir()

	// Создаем файл с нулевым содержимым
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
	// Проверяем, что содержимое заменено на новую строку с байтами \n
	expected := string(bytes.ReplaceAll([]byte{0, 0, 0}, []byte("\x00"), []byte("\n")))
	if val != expected {
		t.Errorf("Expected %q, got %q", expected, val)
	}
}

func TestReadDirFailOpenDir(t *testing.T) {
	// Попытка прочитать несуществующую директорию
	_, err := ReadDir("/path/to/nonexistent/dir")
	if err == nil {
		t.Fatal("Expected error for non-existent directory, got nil")
	}
}

func TestReadDirFileReadError(t *testing.T) {
	// Создаем временную папку
	dir := t.TempDir()

	// Создаем файл, затем изменяем его разрешения, чтобы вызвать ошибку чтения
	filename := "VAR"
	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, []byte("val"), 0o644); err != nil {
		t.Fatalf("Error writing file: %v", err)
	}
	// Убираем все разрешения чтобы чтение вызвало ошибку
	os.Chmod(path, 0o000)

	// Восстановим разрешения после теста
	defer os.Chmod(path, 0o644)

	_, err := ReadDir(dir)
	if err == nil {
		t.Fatal("Expected error on file read, got nil")
	}
}
