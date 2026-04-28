package main

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

var ErrNotDirectory = errors.New("dir path is not a directory")

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
		return nil, ErrNotDirectory
	}
	environment := make(Environment)
	entryes, err := os.ReadDir(dir)
	for _, entry := range entryes {
		if entry.IsDir() {
			continue
		}
		fileName := entry.Name()
		if strings.Contains(fileName, "=") {
			continue
		}
		fileInfo, _ := os.Stat(entry.Name())
		envValue := EnvValue{}
		if fileInfo.Size() == 0 {
			envValue.Value = ""
			envValue.NeedRemove = true
		} else {
			value, err := valueFromFile(entry.Name())
			if err != nil {
				continue
			}
			envValue.Value = value
			envValue.NeedRemove = false
		}
		environment[entry.Name()] = envValue
	}
	return environment, nil
}

func valueFromFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	reader := bufio.NewReader(file)
	firstLine := []byte{}
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return "", err
		}
		if b == 0x00 || b == '\n' {
			break
		}
		firstLine = append(firstLine, b)
	}
	return string(firstLine), nil
}
