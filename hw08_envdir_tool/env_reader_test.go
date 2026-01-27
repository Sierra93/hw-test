package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadValueFromFile(t *testing.T) {
	tests := []struct {
		description string
		file        string
		want        string
	}{
		{
			"ignore second line",
			"testdata/env/BAR",
			"bar",
		},
		{
			"empty value in file",
			"testdata/env/EMPTY",
			"",
		},
		{
			"replace 0x00 with new line",
			"testdata/env/FOO",
			"   foo\nwith new line",
		},
		{
			"set value",
			"testdata/env/HELLO",
			"\"hello\"",
		},
		{
			"unset value",
			"testdata/env/UNSET",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			got, err := readValueFromFile(tt.file)
			require.NoError(t, err)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestReadDir(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		got, err := ReadDir("testdata/env")
		require.NoError(t, err)

		want := Environment{
			"BAR":   {Value: "bar", NeedRemove: false},
			"EMPTY": {Value: "", NeedRemove: false},
			"FOO":   {Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": {Value: "\"hello\"", NeedRemove: false},
			"UNSET": {Value: "", NeedRemove: true},
		}
		assert.Equal(t, want, got)
	})

	t.Run("invalid directory", func(t *testing.T) {
		_, err := ReadDir("testdata/env/BAR")
		assert.Error(t, err)
	})

	t.Run("unsupported file name", func(t *testing.T) {
		_, err := ReadDir("testdata/err")
		assert.ErrorIs(t, err, ErrUnsupportedFileName)
	})
}
