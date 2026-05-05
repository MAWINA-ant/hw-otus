package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("directrory is not exist", func(t *testing.T) {
		_, err := ReadDir("./testdata/envi")
		require.Equal(t, "stat ./testdata/envi: no such file or directory", err.Error())
	})

	t.Run("path is not directory", func(t *testing.T) {
		_, err := ReadDir("./testdata/echo.sh")
		require.Equal(t, err, ErrNotDirectory)
	})

	tests := []struct {
		name     string
		path     string
		err      error
		expected Environment
	}{
		{"empty direcrtory", "./emptyDir", nil, make(Environment)},
		{"env directory", "./testdata/env", nil, Environment{
			"BAR":   {"bar", false},
			"EMPTY": {"", false},
			"FOO":   {"   foo\nwith new line", false},
			"HELLO": {"\"hello\"", false},
			"UNSET": {"", true},
		}},
	}
	os.Mkdir("emptyDir", os.ModeDir)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			environment, err := ReadDir(tc.path)
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.expected, environment)
		})
	}
	os.Remove("emptyDir")
}
