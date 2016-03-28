package main

import (
	"bytes"
	"os"
	"testing"
)

func TestRunHelp(t *testing.T) {
	args := []string{
		"--help",
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
		"--list",
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
		"--deauth",
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

	nova := NewNova()
	_, err := os.Stat(nova.credentialCachePath())
	if err == nil {
		t.Errorf("Credential cache file sill exists")
	}
}
