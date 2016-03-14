package main

import (
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
)

const (
	DEFAULT_SSH_COMMAND = "ssh"
	APPNAME             = "novassh"
	VERSION             = "0.1"
)

type Config struct {
	// Executable name of SSH
	SshCommand string

	// Option flags for SSH command
	SshOptions []string

	// Hostname to connect to the instance
	SshHost string

	// Username of SSH
	SshUser string

	// Command-name to be run on the instance
	SshRemoteCommand string
}

func (c *Config) ParseArgs(args []string) (exitWithHelp bool, err error) {
	// Environments
	if os.Getenv("NOVASSH_COMMAND") != "" {
		c.SshCommand = os.Getenv("NOVASSH_COMMAND")
	}

	// Aeguments
	i := 0
	sshargs := []string{}
	for i < len(args) {
		arg := args[i]
		if arg == "--novassh-debug" {
			// Enable debug
			log.SetLevel(log.DebugLevel)
			enableDebugTransport()

		} else if arg == "--novassh-command" {
			// Detects SSH command
			i++
			c.SshCommand = args[i]

		} else if arg == "--help" {
			sshargs = []string{}
			break

		} else {
			sshargs = append(sshargs, arg)
		}
		i++
	}

	// Set default SSH command if not set
	if c.SshCommand == "" {
		c.SshCommand = DEFAULT_SSH_COMMAND
	}

	// Any SSH args was not provided.
	if len(sshargs) == 0 {
		help()
		return true, nil
	}

	return false, c.parseSshArgs(sshargs)
}

func (c *Config) parseSshArgs(args []string) (err error) {
	nova := NewNova()
	if err := nova.Init(); err != nil {
		return err
	}

	found := false
	pos := len(args) - 1 // position of machine name in arguments
	for pos >= 0 {
		arg := args[pos]
		found, err = c.resolveMachineName(nova, arg)
		if err != nil {
			return err

		} else if found {
			break
		}
		pos--
	}

	if found {
		if pos > 0 {
			c.SshOptions = args[:pos]
		}
		if len(args) > 1 {
			c.SshRemoteCommand = strings.Join(args[pos+1:], " ")
		}
		log.Debugf("The machine is found: ipaddr=%s, args=%v command=%s", c.SshHost, c.SshOptions, c.SshRemoteCommand)
		return nil

	} else {
		return fmt.Errorf("Could not found the machine.")
	}
}

func (c *Config) resolveMachineName(nova *nova, arg string) (found bool, err error) {
	var user, instancename string

	if strings.Index(arg, "@") > 0 {
		ss := strings.Split(arg, "@")
		user = ss[0]
		instancename = ss[1]
	} else {
		instancename = arg
	}

	log.Debugf("Try to find the machine: instance-name=%s", instancename)

	machine, err := nova.Find(instancename)
	if err != nil {
		// error
		return false, err

	} else if machine != nil {
		// Found
		c.SshHost = machine.Ipaddr
		c.SshUser = user
		return true, nil

	} else {
		// Not found
		log.Debugf("No match: name=%s", instancename)
		return false, nil
	}
}

func help() {
	fmt.Fprintf(os.Stdout, `NAME:
	%s - The ssh wrapper program to connect OpenStack instance(nova) with the instance name.

USAGE:
	%s [ssh-options] user@hostname [comamnd]

VERSION:
	%s

OPTIONS:
	--novassh-command: Specify SSH command (default: "ssh").
	--novassh-debug:   output some debug messages.
	--help:            print this message.

ENVIRONMENTS:
	NOVASSH_COMMAND: Specify SSH command (default: "ssh").

`, APPNAME, APPNAME, VERSION)
}
