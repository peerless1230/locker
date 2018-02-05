package main

import (
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/peerless1230/locker/cgroups"
	"github.com/peerless1230/locker/cgroups/subsystems"
	"github.com/peerless1230/locker/common"

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
func Run(tty bool, cmdArray []string, res *subsystems.ResourceLimitConfig) {
	parent, writePipe := container.NewParentProcess(tty)
	// check parent process inited successfully.
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	cgroupManager := cgroups.NewCgroupManager("locker")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)
	sendInitCommand(cmdArray, writePipe)
	err := parent.Wait()
	common.CheckError(err)
	log.Debugf("Parent process exited.")
	os.Exit(0)
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command: %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
