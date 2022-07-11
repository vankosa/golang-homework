package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	str := "test_string"
	file, err := os.OpenFile("./test.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o777)
	require.Nil(t, err, err)
	defer file.Close()
	defer os.Remove(file.Name())
	require.Nil(t, err, err)

	_, err = file.WriteString(str)
	require.Nil(t, err, err)

	err = Copy(file.Name(), "./test2.txt", 5, 6)
	require.Nil(t, err, err)

	file2, err := os.OpenFile("./test2.txt", os.O_RDONLY, 0o777)
	require.Nil(t, err, err)
	defer file2.Close()
	defer os.Remove("./test2.txt")

	source, err := ioutil.ReadAll(file2)
	require.Nil(t, err, err)
	require.Equal(t, []byte("string"), source)
}
