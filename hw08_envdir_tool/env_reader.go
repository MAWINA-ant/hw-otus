package main

import (
	"errors"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	info, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, errors.New("dir path is not a directory")
	}
	var environment Environment
	entryes, err := os.ReadDir(dir)
	for _, entry := range entryes {
		if entry.IsDir() {
			continue
		}
		fileName := entry.Name()
		if strings.Contains(fileName, "=") {
			continue
		}

	}
	return environment, nil
}
