package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
)

func main() {
	var err error

	c := Config{}
	exitWithHelp, err := c.ParseArgs(os.Args[1:])

	if exitWithHelp {
		os.Exit(0)

	} else if err != nil {
		log.Errorf("%v", err)
	}

	ssh := &Ssh{config: c}
	ssh.Run()
}
