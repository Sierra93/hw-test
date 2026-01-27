package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
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
		if value.Value != expected {
			t.Errorf("For %s: expected %q, got %q", key, expected, value.Value)
		}
	}
}

func TestReadDir_ZeroBytesFile(t *testing.T) {
	dir := t.TempDir()

	// Создаем файл с тремя нулевыми байтами
	filename := "VAR_ZERO"
	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, []byte{0, 0, 0}, 0o644); err != nil {
		t.Fatalf("Error creating file: %v", err)
	}

	env, err := ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir failed: %v", err)
	}

	val, ok := env[filename]
	if !ok {
		t.Fatalf("Variable %s not found in env", filename)
	}

	// 1. Заменяем 0 на \n -> получаем "\n\n\n"
	expected := string(bytes.ReplaceAll([]byte{0, 0, 0}, []byte("\x00"), []byte("\n")))
	// 2. Обрезаем последний \n (как это делает ваша функция) -> получаем "\n\n"
	expected = strings.TrimSuffix(expected, "\n")

	if val.Value != expected {
		t.Errorf("Expected %q (len %d), got %q (len %d)", expected, len(expected), val.Value, len(val.Value))
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
