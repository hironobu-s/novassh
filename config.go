package main

import (
	"fmt"
	"os"
	"strings"

	"io"

	log "github.com/Sirupsen/logrus"
)

const (
	DEFAULT_SSH_COMMAND = "ssh"
	APPNAME             = "novassh"
	VERSION             = "0.1"
)

// Commands
const (
	CMD_HELP = iota + 1
	CMD_LIST
	CMD_SSH
	CMD_DEAUTH
)

type Config struct {
	// Outputs
	Stdout io.Writer
	Stdin  io.Reader
	Stderr io.Writer

	// Arguments
	Args []string

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

func (c *Config) ParseArgs() (command int, err error) {
	// Environments
	if os.Getenv("NOVASSH_COMMAND") != "" {
		c.SshCommand = os.Getenv("NOVASSH_COMMAND")
	}

	// Aeguments
	i := 0
	sshargs := []string{}
	for i < len(c.Args) {
		arg := c.Args[i]
		if arg == "--novassh-debug" {
			// Enable debug
			log.SetLevel(log.DebugLevel)
			enableDebugTransport()

		} else if arg == "--novassh-command" {
			// Detects SSH command
			i++
			c.SshCommand = c.Args[i]

		} else if arg == "--novassh-list" {
			// List instances
			command = CMD_LIST

		} else if arg == "--novassh-deauth" {
			// Remove credential cache
			command = CMD_DEAUTH

		} else if arg == "--novassh-help" {
			command = CMD_HELP
			break

		} else {
			command = CMD_SSH
			sshargs = append(sshargs, arg)
		}
		i++
	}

	// Set default SSH command if not set
	if c.SshCommand == "" {
		c.SshCommand = DEFAULT_SSH_COMMAND
	}

	// Display help if no arguments are given
	if command == 0 && len(sshargs) == 0 {
		command = CMD_HELP
	}

	if command == CMD_SSH {
		return CMD_SSH, c.parseSshArgs(sshargs)
	} else {
		return command, nil
	}
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
		log.Debugf("The server is found: ipaddr=%s, args=%v command=%s", c.SshHost, c.SshOptions, c.SshRemoteCommand)
		return nil

	} else {
		return fmt.Errorf("Could not found the server.")
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

	log.Debugf("Try to find the server: instance-name=%s", instancename)

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
