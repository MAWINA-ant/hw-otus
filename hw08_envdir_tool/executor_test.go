package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	tests := []struct {
		name        string
		cmd         []string
		environment Environment
		expected    int
	}{
		{
			"run echo from testdata empty env",
			[]string{
				"/bin/bash",
				"./testdata/echo.sh",
				"arg1=1",
				"arg2=2"},
			Environment{},
			0,
		},
		{
			"run echo from testdata with env",
			[]string{
				"/bin/bash",
				"./testdata/echo.sh",
				"arg1=1",
				"arg2=2"},
			Environment{
				"BAR":   {"bar", false},
				"EMPTY": {"", false},
				"FOO":   {"   foo\nwith new line", false},
				"HELLO": {"\"hello\"", false},
				"UNSET": {"", true},
			},
			0,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			exitCode := RunCmd(tc.cmd, tc.environment)
			require.Equal(t, tc.expected, exitCode)
		})
	}
}
