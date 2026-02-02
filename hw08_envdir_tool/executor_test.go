package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCmd(t *testing.T) {
	t.Run("empty command", func(t *testing.T) {
		code := RunCmd(nil, nil)
		assert.Equal(t, 0, code)
	})
	t.Run("equal exit code", func(t *testing.T) {
		code := RunCmd([]string{"/bin/bash", "testdata/echo_fail.sh"}, nil)
		assert.Equal(t, 2, code)
	})

	t.Run("success", func(t *testing.T) {
		env := Environment{
			"BAR":   {Value: "bar", NeedRemove: false},
			"EMPTY": {Value: "", NeedRemove: false},
			"FOO":   {Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": {Value: "\"hello\"", NeedRemove: false},
			"UNSET": {Value: "", NeedRemove: true},
		}

		code := RunCmd([]string{"/bin/bash", "testdata/echo.sh", "arg1=1", "arg2=2"}, env)
		assert.Equal(t, 0, code)
	})
}
