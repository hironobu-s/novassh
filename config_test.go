package main

import (
	"testing"
)

var configTestInstance *machine

func TestInstanceName(t *testing.T) {
	n := NewNova()
	n.Init()

	machines, _ := n.List()
	if len(machines) == 0 {
		t.Skipf("No machines found. Skip this test.")
	}

	configTestInstance = machines[0]
}

// ---------------------------------------

// Only instance name in the arguments.
func TestParseArgs1(t *testing.T) {
	if configTestInstance == nil {
		t.Skipf("No servers found. Skip this test.")
	}

	c := &Config{}

	args := []string{
		configTestInstance.Name,
	}

	cmd, err := c.ParseArgs(args)
	if cmd != CMD_SSH {
		t.Errorf("Command should be CMD_SSH: command=%d", cmd)

	} else if err != nil {
		t.Errorf("%v", err)
	}

	if c.SshHost != configTestInstance.Ipaddr {
		t.Errorf("hostname is not match: %v", c)
	}
	if c.SshUser != "" {
		t.Errorf("username is not match: %v", c)
	}
	if c.SshRemoteCommand != "" {
		t.Errorf("remote-command is not match: %v", c)
	}
}

// Instance name with user in the arguments.
func TestParseArgs2(t *testing.T) {
	if configTestInstance == nil {
		t.Skipf("No servers found. Skip this test.")
	}

	c := &Config{}

	args := []string{
		"root@" + configTestInstance.Name,
	}

	cmd, err := c.ParseArgs(args)
	if cmd != CMD_SSH {
		t.Errorf("Command should be CMD_SSH: command=%d", cmd)

	} else if err != nil {
		t.Errorf("%v", err)
	}

	if c.SshHost != configTestInstance.Ipaddr {
		t.Errorf("hostname is not match: %v", c)
	}
	if c.SshUser != "root" {
		t.Errorf("username is not match: %v", c)
	}
	if c.SshRemoteCommand != "" {
		t.Errorf("remote-command is not match: %v", c)
	}
}

// Instance name with user and remote commands in the arguments
func TestParseArgs3(t *testing.T) {
	if configTestInstance == nil {
		t.Skipf("No servers found. Skip this test.")
	}

	c := &Config{}

	args := []string{
		"root@" + configTestInstance.Name,
		"test-command",
	}

	cmd, err := c.ParseArgs(args)
	if cmd != CMD_SSH {
		t.Errorf("Command should be CMD_SSH: command=%d", cmd)

	} else if err != nil {
		t.Errorf("%v", err)
	}

	if c.SshHost != configTestInstance.Ipaddr {
		t.Errorf("hostname is not match: %v", c)
	}
	if c.SshUser != "root" {
		t.Errorf("username is not match: %v", c)
	}
	if c.SshRemoteCommand != "test-command" {
		t.Errorf("remote-command is not match: %v", c)
	}
}

// With SSH options
func TestParseArgs4(t *testing.T) {
	if configTestInstance == nil {
		t.Skipf("No servers found. Skip this test.")
	}

	c := &Config{}

	args := []string{
		// Port fowarding option for ssh
		"-L",
		"54321:localhost:54321",
		configTestInstance.Name,
	}

	cmd, err := c.ParseArgs(args)
	if cmd != CMD_SSH {
		t.Errorf("Command should be CMD_SSH: command=%d", cmd)

	} else if err != nil {
		t.Errorf("%v", err)
	}

	if c.SshOptions[0] != "-L" || c.SshOptions[1] != "54321:localhost:54321" {
		t.Errorf("ssh options are not match: %v", c)
	}
	if c.SshHost != configTestInstance.Ipaddr {
		t.Errorf("hostname is not match: %v", c)
	}
	if c.SshUser != "" {
		t.Errorf("username is not match: %v", c)
	}
	if c.SshRemoteCommand != "" {
		t.Errorf("remote-command is not match: %v", c)
	}
}
