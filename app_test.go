package main

import (
	"bytes"
	"testing"
)

func TestRunHelp(t *testing.T) {
	args := []string{
		"--novassh-help",
	}

	c := Config{
		Stdout: new(bytes.Buffer),
		Stdin:  nil,
		Stderr: nil,
		Args:   args,
	}

	code := run(c)
	if code != 0 {
		t.Errorf("failure exit code: %d", code)
	}
}

func TestRunList(t *testing.T) {
	args := []string{
		"--novassh-list",
	}

	c := Config{
		Stdout: new(bytes.Buffer),
		Stdin:  nil,
		Stderr: nil,
		Args:   args,
	}
	code := run(c)
	if code != 0 {
		t.Errorf("failure exit code: %d", code)
	}
}

func TestRunDeauth(t *testing.T) {
	args := []string{
		"--novassh-deauth",
	}
	c := Config{
		Stdout: new(bytes.Buffer),
		Stdin:  nil,
		Stderr: nil,
		Args:   args,
	}

	code := run(c)
	if code != 0 {
		t.Errorf("failure exit code: %d", code)
	}
}
