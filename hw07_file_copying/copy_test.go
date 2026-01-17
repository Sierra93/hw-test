package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopy(t *testing.T) {
	// Создаем временную директорию для теста
	tempDir := t.TempDir()

	inputFile := filepath.Join(tempDir, "input.txt")
	outputFile := filepath.Join(tempDir, "output.txt")

	content := []byte("Hello, world! This is a test file.")
	if err := os.WriteFile(inputFile, content, 0o644); err != nil {
		t.Fatalf("Failed to create input file: %v", err)
	}

	// Тест успешного копирования
	if err := Copy(inputFile, outputFile, 0, int64(len(content)), false); err != nil {
		t.Fatalf("Copy failed: %v", err)
	}

	result, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if string(result) != string(content) {
		t.Errorf("Expected %s, got %s", string(content), string(result))
	}

	// Тест с offset и limit
	offset := int64(7)
	limit := int64(5) // "world"
	if err := Copy(inputFile, outputFile, offset, limit, false); err != nil {
		t.Fatalf("Copy with offset and limit failed: %v", err)
	}

	result, err = os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	expected := content[offset : offset+limit]
	if string(result) != string(expected) {
		t.Errorf("Expected %s, got %s", string(expected), string(result))
	}

	// Тест с offset больше размера файла
	if err := Copy(inputFile, outputFile, int64(len(content)+10), 10, false); err == nil {
		t.Error("Expected error for offset exceeding file size, got nil")
	}

	// Тест без ограничения (limit=0)
	if err := Copy(inputFile, outputFile, 0, 0, false); err != nil {
		t.Fatalf("Copy with limit=0 failed: %v", err)
	}
	result, err = os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	if string(result) != string(content) {
		t.Errorf("Expected full content, got %s", string(result))
	}
}
