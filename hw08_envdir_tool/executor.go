package main

import (
	"io"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) < 1 {
		return 1
	}

	startCmd := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec

	startCmd.Stdout = os.Stdout
	startCmd.Stderr = os.Stderr

	stdin, err := startCmd.StdinPipe()
	if err != nil {
		return 1
	}
	defer stdin.Close()

	// start coping input data from parent to child process
	go func() {
		io.Copy(stdin, os.Stdin)
	}()

	// calculate invironments
	for envName, envValue := range env {
		if envValue.NeedRemove {
			os.Unsetenv(envName)
			continue
		}
		os.Setenv(envName, envValue.Value)
	}

	// update envs
	startCmd.Env = os.Environ()

	// start process
	err = startCmd.Start()
	if err != nil {
		return 1
	}

	// wait process finish
	err = startCmd.Wait()
	if err != nil {
		return 1
	}

	return 0
}
