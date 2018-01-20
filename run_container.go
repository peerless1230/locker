package main

import (
	"os"

	"./container/"
	log "github.com/Sirupsen/logrus"
)

/*
Run is used to Run the command given to container
Params: tty bool, command string, args []string
Return: error
*/
func Run(tty bool, command string, args []string) {
	parent := container.NewParentProcess(tty, command, args)
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	parent.Wait()
	os.Exit(-1)
}
