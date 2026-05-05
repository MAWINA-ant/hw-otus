package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
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
	entryes, _ := os.ReadDir(dir)
	for _, entry := range entryes {
		if entry.IsDir() {
			continue
		}
		fileName := entry.Name()
		if strings.Contains(fileName, "=") {
			continue
		}
		fileInfo, err := entry.Info()
		if err != nil {
			fmt.Println("stat???? ", entry.Name(), err)
			continue
		}
		envValue := EnvValue{}
		if fileInfo.Size() == 0 {
			envValue.Value = ""
			envValue.NeedRemove = true
		} else {
			value, err := valueFromFile(filepath.Join(dir, entry.Name()))
			if err != nil {
				fmt.Println("Value from ", entry.Name())
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
		fmt.Println("not open")
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var firstLine string
	if scanner.Scan() {
		firstLine = scanner.Text()
	}
	firstLineBytes := bytes.ReplaceAll([]byte(firstLine), []byte("\x00"), []byte("\n"))
	firstLine = strings.TrimRight(string(firstLineBytes), " \t")
	return firstLine, nil
}
