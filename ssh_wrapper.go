package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type SshWrapper struct {
	SshCommand string
	Args       []string
	User       string
	Host       string
	Port       string
	Command    string
}

func (s *SshWrapper) ParseArgs(args []string) (err error) {
	nova := NewNova()
	if err := nova.Init(); err != nil {
		return err
	}

	found := false
	pos := len(args) - 1 // position of machine name in arguments
	for pos >= 0 {
		arg := args[pos]
		found, err = s.resolveMachineName(nova, arg)
		if err != nil {
			return err

		} else if found {
			break
		}
		pos--
	}

	if found {
		if pos > 0 {
			s.Args = args[:pos]
		}
		if len(args) > 1 {
			s.Command = strings.Join(args[pos+1:], " ")
		}
		log.Debugf("The machine is found: ipaddr=%s, args=%v command=%s", s.Host, s.Args, s.Command)
		return nil

	} else {
		return fmt.Errorf("Could not found the machine.")
	}
}

func (s *SshWrapper) resolveMachineName(nova *nova, arg string) (found bool, err error) {
	var user, name, port string

	if strings.Index(arg, "@") > 0 {
		ss := strings.Split(arg, "@")
		user = ss[0]
		name = ss[1]
	} else {
		name = arg
	}

	if strings.Index(name, ":") > 0 {
		ss := strings.Split(name, ":")
		name = ss[0]
		port = ss[1]
	}

	log.Debugf("Try to find the machine: name=%s", name)

	machine, err := nova.Find(name)
	if err != nil {
		// error
		return false, err

	} else if machine != nil {
		// Found
		s.Port = port
		s.Host = machine.Ipaddr
		s.User = user
		return true, nil

	} else {
		// Not found
		log.Debugf("No match: name=%s", name)
		return false, nil
	}
}

func (s *SshWrapper) Run() {
	var server string
	if s.User != "" {
		server += s.User + "@"
	}
	server += s.Host
	if s.Port != "" {
		server += ":" + s.Port
	}

	cmd := exec.Command(s.SshCommand, append(s.Args, server, s.Command)...)

	log.Debugf("ssh command:%v", cmd.Args)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
