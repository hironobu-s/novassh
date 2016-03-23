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
		Stdout:   new(bytes.Buffer),
		Stdin:    nil,
		Stderr:   nil,
		ConnType: CON_SSH,
		Args:     args,
	}

	con := Connect{
		config: c,
	}
	con.Run()
}

func TestConsoleRun(t *testing.T) {
	args := []string{}

	c := Config{
		Stdout:   new(bytes.Buffer),
		Stdin:    nil,
		Stderr:   nil,
		ConnType: CON_CONSOLE,
		Args:     args,
	}

	con := Connect{
		config: c,
	}
	con.Run()
}
