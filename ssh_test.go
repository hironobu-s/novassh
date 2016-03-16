package main

import (
	"bytes"
	"testing"
)

func TestSshRun(t *testing.T) {
	args := []string{
		"dummy-instance",
		"w",
	}

	c := Config{
		Stdout: new(bytes.Buffer),
		Stdin:  nil,
		Stderr: nil,
		Args:   args,
	}

	ssh := Ssh{
		config: c,
	}
	ssh.Run()
}
