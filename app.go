package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	enableDebugTransport()

	ssh := &SshWrapper{
		SshCommand: "ssh",
	}
	if err := ssh.ParseArgs(os.Args[1:]); err != nil {
		log.Errorf("%v", err)
	}
	ssh.Run()
}
