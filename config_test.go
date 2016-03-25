package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
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

	args := []string{
		configTestInstance.Name,
	}

	c := &Config{Args: args}
	cmd, err := c.ParseArgs()
	if cmd != CMD_CONNECT {
		t.Errorf("Command should be CMD_CONNECT: command=%d", cmd)
	} else if err != nil {
		t.Errorf("%v", err)
	}

	if c.ConnType != CON_SSH {
		t.Errorf("ConnType should be CON_SSH: type=%d", c.ConnType)
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

	args := []string{
		"root@" + configTestInstance.Name,
	}

	c := &Config{Args: args}
	cmd, err := c.ParseArgs()
	if cmd != CMD_CONNECT {
		t.Errorf("Command should be CMD_CONNECT: command=%d", cmd)
	} else if err != nil {
		t.Errorf("%v", err)
	}

	if c.ConnType != CON_SSH {
		t.Errorf("ConnType should be CON_SSH: type=%d", c.ConnType)
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

	args := []string{
		"root@" + configTestInstance.Name,
		"test-command",
	}

	c := &Config{Args: args}
	cmd, err := c.ParseArgs()
	if cmd != CMD_CONNECT {
		t.Errorf("Command should be CMD_CONNECT: command=%d", cmd)
	} else if err != nil {
		t.Errorf("%v", err)
	}

	if c.ConnType != CON_SSH {
		t.Errorf("ConnType should be CON_SSH: type=%d", c.ConnType)
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

	args := []string{
		// Port fowarding option for ssh
		"-L",
		"54321:localhost:54321",
		configTestInstance.Name,
	}

	c := &Config{Args: args}
	cmd, err := c.ParseArgs()
	if cmd != CMD_CONNECT {
		t.Errorf("Command should be CMD_CONNECT: command=%d", cmd)
	} else if err != nil {
		t.Errorf("%v", err)
	}

	if c.ConnType != CON_SSH {
		t.Errorf("ConnType should be CON_SSH: type=%d", c.ConnType)
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

func TestHelp(t *testing.T) {
	args := []string{
		"--novassh-help",
	}

	c := &Config{Args: args}
	cmd, err := c.ParseArgs()
	if cmd != CMD_HELP {
		t.Errorf("Command should be CMD_CONNECT: command=%d", cmd)
	} else if err != nil {
		t.Errorf("%v", err)
	}
}

func TestList(t *testing.T) {
	args := []string{
		"--novassh-list",
	}

	c := &Config{
		Stdout: new(bytes.Buffer),
		Stdin:  nil,
		Stderr: nil,
		Args:   args,
	}
	cmd, err := c.ParseArgs()
	if cmd != CMD_LIST {
		t.Errorf("Command should be CMD_CONNECT: command=%d", cmd)
	} else if err != nil {
		t.Errorf("%v", err)
	}
}

func TestDeauth(t *testing.T) {
	args := []string{
		"--novassh-deauth",
	}

	c := &Config{Args: args}
	cmd, err := c.ParseArgs()
	if cmd != CMD_DEAUTH {
		t.Errorf("Command should be CMD_CONNECT: command=%d", cmd)
	} else if err != nil {
		t.Errorf("%v", err)
	}

	nova := NewNova()
	_, err = os.Stat(nova.credentialCachePath())
	if err != nil {
		t.Errorf("Credential cache file sill exists")
	}
}

func TestDebug(t *testing.T) {
	args := []string{
		"--novassh-debug",
	}

	c := &Config{
		Stdout: new(bytes.Buffer),
		Stdin:  nil,
		Stderr: nil,
		Args:   args,
	}
	_, err := c.ParseArgs()
	if err != nil {
		t.Errorf("%v", err)
	}

	// disable debug
	disableDebugTransport()
	logrus.SetLevel(logrus.InfoLevel)
}

func TestConsole(t *testing.T) {
	if configTestInstance == nil {
		t.Skipf("No servers found. Skip this test.")
	}

	args := []string{
		"--novassh-console",
		configTestInstance.Name,
	}

	c := &Config{Args: args}
	_, err := c.ParseArgs()
	if c.ConnType != CON_CONSOLE {
		t.Errorf("ConnType should be CON_CONSOLE: type=%d", c.ConnType)
	} else if err != nil {
		t.Errorf("%v", err)
	}
}
