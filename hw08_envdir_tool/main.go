package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args

	envs, err := ReadDir(args[1])
	if err != nil {
		log.Fatal(err)
	}

	RunCmd(args[2:], envs)
}
