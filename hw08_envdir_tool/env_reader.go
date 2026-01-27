package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

const unsupportedSymbol = "="

var ErrUnsupportedFileName = errors.New("file name contains unsupported symbol '='")

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	env := make(Environment, len(files))

	for _, file := range files {
		finfo, err := file.Info()
		if err != nil {
			return nil, fmt.Errorf("failed to get file info: %w", err)
		}

		if finfo.IsDir() {
			continue
		}

		if strings.Contains(finfo.Name(), unsupportedSymbol) {
			return nil, fmt.Errorf("%q: %w", finfo.Name(), ErrUnsupportedFileName)
		}

		value, err := readValueFromFile(path.Join(dir, finfo.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to read env value: %w", err)
		}

		env[file.Name()] = EnvValue{
			Value:      value,
			NeedRemove: finfo.Size() == 0,
		}
	}

	return env, nil
}

func readValueFromFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	s := strings.ReplaceAll(string(content), "\x00", "\n")
	lines := strings.Split(s, "\n")
	value := lines[0]
	value = strings.TrimRight(value, " \t")

	return value, nil
}
