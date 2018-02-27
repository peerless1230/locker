package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/peerless1230/locker/cgroups"
	"github.com/peerless1230/locker/cgroups/subsystems"
	"github.com/peerless1230/locker/common"

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
Params: tty bool, cmdArray []string, res *subsystems.ResourceLimitConfig, volumeSlice []string, containerName string
Return:
*/
func Run(tty bool, cmdArray []string, res *subsystems.ResourceLimitConfig, volumeSlice []string, containerName string) {
	containerID := common.RandStringBytes(64)
	parent, writePipe := container.NewParentProcess(tty, volumeSlice, containerID)
	// check parent process inited successfully.
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	//record container info
	containerName, err := recordContainerInfo(parent.Process.Pid, cmdArray, containerName, containerID)
	if err != nil {
		log.Errorf("Record container info error %v", err)
		return
	}

	cgroupManager := cgroups.NewCgroupManager(containerName)
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)
	sendInitCommand(cmdArray, writePipe)
	err = parent.Wait()

	layerPath := filepath.Join(rootLAYER, containerID)
	container.CleanUpVolumes(volumeSlice, layerPath)
	container.CleanUpOverlayFS(layerPath)
	deleteContainerInfo(containerID)
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

func recordContainerInfo(containerPID int, commandArray []string, containerName string, containerID string) (string, error) {
	id := containerID
	createTime := time.Now().Format("2006-01-02 15:04:05")
	command := strings.Join(commandArray, "")
	containerInfo := &container.Info{
		ID:          id,
		Pid:         strconv.Itoa(containerPID),
		Command:     command,
		CreatedTime: createTime,
		Status:      container.RUNNING,
		Name:        containerName,
	}

	jsonBytes, err := json.Marshal(containerInfo)
	if err != nil {
		log.Errorf("Record container info error %v", err)
		return "", err
	}
	jsonStr := string(jsonBytes)

	dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerID)
	if err := os.MkdirAll(dirURL, 0622); err != nil {
		log.Errorf("Mkdir error %s error %v", dirURL, err)
		return "", err
	}
	fileName := filepath.Join(dirURL, container.ConfigName)
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		log.Errorf("Create file %s error %v", fileName, err)
		return "", err
	}
	if _, err := file.WriteString(jsonStr); err != nil {
		log.Errorf("File write string error %v", err)
		return "", err
	}

	return containerName, nil
}

func deleteContainerInfo(containerID string) {
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerID)
	common.RmdirAll(dirURL)
}
