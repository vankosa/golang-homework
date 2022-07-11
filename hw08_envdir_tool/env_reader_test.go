package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadDir(t *testing.T) {
	t.Run("dir not found", func(t *testing.T) {
		env, err := ReadDir("not found")
		assert.Error(t, err)
		assert.Len(t, env, 0)
	})

	t.Run("available tests", func(t *testing.T) {
		env, err := ReadDir("./testdata/env")
		assert.NoError(t, err)
		assert.Len(t, env, 5)

		envVal, ok := env["UNSET"]
		assert.True(t, ok)
		expected := EnvValue{
			Value:      "",
			NeedRemove: true,
		}
		assert.Equal(t, envVal, expected)
	})
}
