package main

import (
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
		expected Environment
	}{}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

		})
	}
}
