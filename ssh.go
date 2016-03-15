package main

import (
	"os"
	"os/exec"

	log "github.com/Sirupsen/logrus"
)

type Ssh struct {
	config Config
}

func (s *Ssh) Run() error {
	var server string
	if s.config.SshUser != "" {
		server += s.config.SshUser + "@"
	}
	server += s.config.SshHost

	cmd := exec.Command(s.config.SshCommand, append(s.config.SshOptions, server, s.config.SshRemoteCommand)...)

	log.Debugf("ssh command:%v", cmd.Args)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
