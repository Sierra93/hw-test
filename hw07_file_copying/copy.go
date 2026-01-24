package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == "" || toPath == "" {
		return ErrUnsupportedFile
	}

	fromPathAbs, err := filepath.Abs(fromPath)
	if err != nil {
		return err
	}

	toPathAbs, err := filepath.Abs(toPath)
	if err != nil {
		return err
	}

	if fromPathAbs == toPathAbs {
		return ErrUnsupportedFile
	}

	sourceFileStat, err := os.Stat(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}

	if !sourceFileStat.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	sourceFileSize := sourceFileStat.Size()
	if sourceFileSize < offset {
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 || limit > sourceFileSize-offset {
		limit = sourceFileSize - offset
	}

	source, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer source.Close()

	destination, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destination.Close()

	_, err = source.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	bar := progressbar.DefaultBytes(limit)

	_, err = io.CopyN(io.MultiWriter(destination, bar), source, limit)
	if err != nil && errors.Is(err, io.EOF) {
		return fmt.Errorf("failed to copy data from source file to destination file: %w", err)
	}

	return nil
}
