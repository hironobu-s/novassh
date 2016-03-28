package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
)

const (
	DEFAULT_SSH_COMMAND = "ssh"
	APPNAME             = "novassh"
	VERSION             = "0.2.2"
)

// Commands
const (
	CMD_HELP = iota + 1
	CMD_LIST
	CMD_CONNECT
	CMD_DEAUTH
)

type Command int

func (c *Command) String() string {
	switch *c {
	case CMD_HELP:
		return "HELP"
	case CMD_LIST:
		return "LIST"
	case CMD_CONNECT:
		return "SSH"
	case CMD_DEAUTH:
		return "DEAUTH"
	default:
		return "UNKNOWN"
	}
}

// Connection types
const (
	// Use SSH(default)
	CON_SSH = iota + 1

	// Use serial console via websocket
	CON_CONSOLE
)

type ConnType int

func (c *ConnType) String() string {
	switch *c {
	case CON_SSH:
		return "SSH"
	case CON_CONSOLE:
		return "CONSOLE"
	default:
		return "UNKNOWN"
	}
}

// Config
type Config struct {
	// Outputs
	Stdout io.Writer
	Stdin  io.Reader
	Stderr io.Writer

	// Arguments
	Args []string

	// Connection type
	ConnType ConnType

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

	// Websocket URL (For TYPE_CONSOLE only)
	ConsoleUrl string

	// Authentication cache (Default: false)
	AuthCache bool
}

func (c *Config) ParseArgs() (command Command, err error) {
	// Environments
	if os.Getenv("NOVASSH_COMMAND") != "" {
		c.SshCommand = os.Getenv("NOVASSH_COMMAND")
	}

	// Defaults
	c.SshCommand = DEFAULT_SSH_COMMAND
	c.ConnType = CON_SSH
	c.AuthCache = false

	// Aeguments
	i := 0
	sshargs := []string{}
	for i < len(c.Args) {
		arg := c.Args[i]
		if arg == "--debug" {
			// Enable debug
			log.SetLevel(log.DebugLevel)
			enableDebugTransport()

		} else if arg == "--command" {
			// Detects SSH command
			i++
			c.SshCommand = c.Args[i]

		} else if arg == "--list" {
			// List instances
			command = CMD_LIST

		} else if arg == "--authcache" {
			// Authentication cache
			c.AuthCache = true

		} else if arg == "--deauth" {
			// Remove credential cache
			command = CMD_DEAUTH

		} else if arg == "--help" {
			command = CMD_HELP
			break

		} else if arg == "--console" {
			// Use serial console
			c.ConnType = CON_CONSOLE

		} else {
			command = CMD_CONNECT
			sshargs = append(sshargs, arg)
		}
		i++
	}

	// Display help if no arguments are given
	if command == 0 && len(sshargs) == 0 {
		command = CMD_HELP
	}

	log.Debugf("Command: %s", command.String())

	if command == CMD_CONNECT {
		return CMD_CONNECT, c.parseSshArgs(sshargs)
	} else {
		return command, nil
	}
}

func (c *Config) parseSshArgs(args []string) (err error) {
	nova := NewNova()
	if err := nova.Init(c.AuthCache); err != nil {
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

		// Set console url if type is CON_CONSOLE
		if c.ConnType == CON_CONSOLE {
			url, err := nova.GetConsoleUrl(machine)
			if err != nil {
				return true, err
			}
			c.ConsoleUrl = url
		}

		return true, nil

	} else {
		// Not found
		log.Debugf("No match: name=%s", instancename)
		return false, nil
	}
}
