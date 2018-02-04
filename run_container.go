package main

import (
	"os"

	log "github.com/Sirupsen/logrus"

	"github.com/peerless1230/locker/container"
)

/*
Run is used to Run the command given to container
Params: tty bool, command string, args []string
Return: error
*/
func Run(tty bool, command string, args []string) {
	parent := container.NewParentProcess(tty, command, args)
	// check parent process inited successfully.
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	parent.Wait()
	os.Exit(-1)
}
