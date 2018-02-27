package main

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/peerless1230/locker/cgroups"
	"github.com/peerless1230/locker/cgroups/subsystems"
	"github.com/peerless1230/locker/common"
	"github.com/syndtr/gocapability/capability"

	"github.com/peerless1230/locker/container"
)

const rootLAYER = "/var/lib/locker/overlay2/"

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
Params: tty bool, cmdArray []string, res *subsystems.ResourceLimitConfig, volumeSlice []string
Return:
*/
func Run(tty bool, cmdArray []string, res *subsystems.ResourceLimitConfig, volumeSlice []string) {
	parent, writePipe := container.NewParentProcess(tty, volumeSlice)
	// check parent process inited successfully.
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	pid, err := capability.NewPid(os.Getpid())
	if err != nil {
		log.Debugf("Set %d CAP_SETGID to setgroup error: %v", pid, err)
	}

	cgroupManager := cgroups.NewCgroupManager("locker")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)
	sendInitCommand(cmdArray, writePipe)
	err = parent.Wait()

	containerID := "92745277a8b052e2c50cf757da7140afabd9f6abbae7b6d6516f944a55658dfc"
	layerPath := filepath.Join(rootLAYER, containerID)
	container.CleanUpVolumes(volumeSlice, layerPath)
	container.CleanUpOverlayFS(layerPath)
	common.CheckError(err)
	log.Debugf("Parent process exited .")
	os.Exit(0)
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command: %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
