package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
		return fmt.Errorf("failed to seek in source file: %w", err)
	}

	const bufferSize = 32 * 1024 // 32 KB буфер
	buf := make([]byte, bufferSize)
	var copied int64
	for copied < limit {
		bytesLeft := limit - copied
		readSize := bufferSize
		if bytesLeft < int64(readSize) {
			readSize = int(bytesLeft)
		}

		n, err := source.Read(buf[:readSize])
		if n > 0 {
			nw, errw := destination.Write(buf[:n])
			if errw != nil {
				return fmt.Errorf("ошибка при записи: %w", errw)
			}
			copied += int64(nw)

			// каждые 1 Мб или по завершении
			if copied%(1024*1024) == 0 || copied == limit {
				fmt.Printf("\rСкопировано %d из %d байт (%.2f%%)", copied, limit, float64(copied)*100/float64(limit))
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("ошибка при чтении: %w", err)
		}
	}
	fmt.Println("\nКопирование завершено.")
	return nil
}
