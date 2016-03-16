package main

import (
	"fmt"
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
)

func main() {
	var err error

	c := Config{}
	cmd, err := c.ParseArgs(os.Args[1:])
	if err != nil {
		goto ERROR
	}

	switch cmd {
	case CMD_HELP:
		help(c)

	case CMD_LIST:
		if err = list(c); err != nil {
			goto ERROR
		}

	case CMD_SSH:
		if err = ssh(c); err != nil {
			goto ERROR
		}

	case CMD_DEAUTH:
		if err = deauth(c); err != nil {
			goto ERROR
		}

	default:
		log.Errorf("Undefined command: %s", cmd)
		goto ERROR
	}
	os.Exit(0)

ERROR:
	log.Errorf("%v", err)
	os.Exit(1)
}

func ssh(c Config) error {
	ssh := &Ssh{config: c}
	return ssh.Run()
}

func list(c Config) error {
	nova := NewNova()
	if err := nova.Init(); err != nil {
		return err
	}

	machines, err := nova.List()
	if err != nil {
		return err
	}

	if len(machines) == 0 {
		fmt.Fprintf(os.Stdout, "No server found.\n")
		return nil
	}

	width := 0
	for _, m := range machines {
		if len(m.Name) > width {
			width = len(m.Name)
		}
	}

	format := "%" + strconv.Itoa(-width) + "s %s\n"
	fmt.Fprintf(os.Stdout, format, "[Name]", "[IP Address]")
	for _, m := range machines {
		fmt.Fprintf(os.Stdout, format, m.Name, m.Ipaddr)
	}
	return nil
}

func help(c Config) {
	fmt.Fprintf(os.Stdout, `NAME:
	%s - The ssh wrapper program to connect OpenStack instance(nova) with the instance name.

USAGE:
	%s [ssh-options] user@hostname [comamnd]

VERSION:
	%s

OPTIONS:
	--novassh-list:    Display instances.
	--novassh-command: Specify SSH command (default: "ssh").
	--novassh-debug:   Output some debug messages.
	--help:            Print this message.

ENVIRONMENTS:
	NOVASSH_COMMAND: Specify SSH command (default: "ssh").

`, APPNAME, APPNAME, VERSION)
}

func deauth(c Config) error {
	nova := NewNova()
	return nova.RemoveCredentialCache()
}
