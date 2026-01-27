package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 0
	}

	for name, val := range env {
		if _, ok := os.LookupEnv(name); !ok {
			if err := os.Setenv(name, val.Value); err != nil {
				return 1
			}
		}

		if err := os.Unsetenv(name); err != nil {
			return 1
		}

		if !val.NeedRemove {
			if err := os.Setenv(name, val.Value); err != nil {
				return 1
			}
		}
	}

	var args []string
	if len(cmd) > 1 {
		args = cmd[1:]
	}

	command := exec.Command(cmd[0], args...) //nolint:gosec
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr

	err := command.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		return 1
	}

	return 0
}
