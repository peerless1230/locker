package subsystems

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"github.com/peerless1230/locker/common"

	log "github.com/Sirupsen/logrus"
)

/*
MemorySubSystem implemented Subsystem interface
*/
type MemorySubSystem struct {
}

/*
GetName used to return the name of subsystem
Params:
Return: "memory"
*/
func (subsys *MemorySubSystem) GetName() string {
	return "memory"
}

/*
Set the cgroup's limit config
Params: cgroupPath string, res *ResourceLimitConfig
Return: error
*/
func (subsys *MemorySubSystem) Set(cgroupPath string, res *ResourceLimitConfig) error {
	subsysCgroupPath, err := GetCgroupPath(subsys.GetName(), cgroupPath, true)
	if err == nil {
		// Write the limits to cgroup's config file
		if res.MemoryLimits != "" {
			limitsFilePath := path.Join(subsysCgroupPath, memoryLimitsFileName)
			if err := ioutil.WriteFile(limitsFilePath, []byte(res.MemoryLimits), 0644); err == nil {
				log.Debugf("Write memory Limits: %s to %s", res.MemoryLimits, limitsFilePath)
			} else {
				return fmt.Errorf("Set memory limits failed: %v", err)
			}
			return nil
		}
	}

	return err
}

/*
Remove used to remove subsystem by the cgroup path
Params: cgroupPath string
Return: error
*/
func (subsys *MemorySubSystem) Remove(cgroupPath string) error {
	subsysCgroupPath, err := GetCgroupPath(subsys.GetName(), cgroupPath, false)
	if err == nil {
		return os.RemoveAll(subsysCgroupPath)
	}
	return err
}

/*
Apply used to add a process into the cgroup
Params: cgroupPath string, pid int
Return: error
*/
func (subsys *MemorySubSystem) Apply(cgroupPath string, pid int) error {
	if subsysCgroupPath, err := GetCgroupPath(subsys.GetName(), cgroupPath, false); err == nil {
		tasksFilePath := path.Join(subsysCgroupPath, tasksFileName)
		var taskFile *os.File

		var lineFeed bool
		if isExist, _ := common.IsPathOrFileExists(tasksFilePath); isExist == true {
			taskFile, err = os.OpenFile(tasksFilePath, os.O_RDWR|os.O_APPEND, 0644)
			log.Debugf("Open file : (%s) for add pid", tasksFilePath)
		} else {
			taskFile, err = os.Create(tasksFilePath)
			log.Debugf("Created file : (%s) for add pid", tasksFilePath)
		}
		defer taskFile.Close()
		if err != nil {
			return fmt.Errorf("Open tasks file failed: %v", err)
		}
		tmp := make([]byte, 1)
		fstate, _ := taskFile.Stat()
		if n, _ := taskFile.ReadAt(tmp, fstate.Size()-1); n != 0 {
			if tmp[0] != '\n' {
				lineFeed = false
			}
		}
		_, err := taskFile.Seek(0, os.SEEK_END)
		if err != nil {

		}
		if lineFeed == false {
			taskFile.WriteString(string('\n'))
		}
		taskFile.WriteString(strconv.Itoa(pid))
		log.Debugf("Write pid(%s)into tasks:", pid)
		taskFile.Sync()
	}
	return nil

}
