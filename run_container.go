package main

import (
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/peerless1230/locker/container"
)

/*
Run is used to Run the command given to container
Params: tty bool, command string, args []string
Return: error
*/
func sendInitCommands(cmdArray []string, writePipe *os.File) {
	cmdStr := strings.Join(cmdArray, " ")
	log.Infof("Container's commands: %s", cmdStr)
	writePipe.WriteString(cmdStr)
	log.Debugf("Write %s to pipe", cmdStr)
	writePipe.Close()
}

/*
Run is used to Run the command given to container
Params: tty bool, command string, args []string
Return: error
*/
func Run(tty bool, cmdArray []string) {
	parent, writePipe := container.NewParentProcess(tty)
	// check parent process inited successfully.
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	sendInitCommands(cmdArray, writePipe)
	parent.Wait()
	os.Exit(-1)
}
