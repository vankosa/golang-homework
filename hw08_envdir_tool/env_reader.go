package main

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
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
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	result := Environment{}
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		file, err := os.Open(filepath.Join(dir, f.Name()))
		if err != nil {
			return nil, err
		}
		defer file.Close()

		reader := bufio.NewReader(file)
		content, err := reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}

		content = strings.ReplaceAll(strings.TrimRight(content, "\t\n "), "\x00", "\n")

		result[f.Name()] = EnvValue{content, len(content) == 0}
	}

	return result, nil
}
