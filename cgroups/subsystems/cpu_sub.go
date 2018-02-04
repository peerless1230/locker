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
CpuSubSystem implemented Subsystem interface
*/
type CpuSubSystem struct {
}

/*
setCPULimit config CPU limit by the params
Params:
Return: "cpu"
*/
func setCPULimit(limitConfig string, limitsFilePath string) error {
	if err := ioutil.WriteFile(limitsFilePath, []byte(limitConfig), 0644); err == nil {
		log.Debugf("Write cpu limit: %s to %s", limitConfig, limitsFilePath)
	} else {
		return fmt.Errorf("Set cpu limit failed: %v", err)
	}
	return nil
}

/*
GetName used to return the name of subsystem
Params:
Return: "cpu"
*/
func (subsys *CpuSubSystem) GetName() string {
	return "cpu"
}

/*
Set the cgroup's limit config
Params: cgroupPath string, res *ResourceLimitConfig
Return: error
*/
func (subsys *CpuSubSystem) Set(cgroupPath string, res *ResourceLimitConfig) error {
	subsysCgroupPath, err := GetCgroupPath(subsys.GetName(), cgroupPath, true)
	if err == nil {
		if res.CPUShare != "" {
			limitsFilePath := path.Join(subsysCgroupPath, cpuSharesFileName)
			err = setCPULimit(res.CPUShare, limitsFilePath)
		}
	}
	if err == nil {
		if res.CPUS != "" {
			limitsFilePath := path.Join(subsysCgroupPath, cpuPeriodFileName)
			err = setCPULimit("100000", limitsFilePath)
			common.CheckError(err)
			//var n float32
			var strCPUQuota string
			limitsFilePath = path.Join(subsysCgroupPath, cpuQuotaFileName)
			n, err := strconv.ParseFloat(res.CPUS, 32)
			common.CheckError(err)
			CPUQuota := 100000 * n
			strCPUQuota = strconv.Itoa(int(CPUQuota))
			err = setCPULimit(strCPUQuota, limitsFilePath)
			common.CheckError(err)
			return nil

		}
	}
	if err == nil {
		if res.CPUPeriod != "" {
			limitsFilePath := path.Join(subsysCgroupPath, cpuPeriodFileName)
			err = setCPULimit(res.CPUPeriod, limitsFilePath)
		}
	}
	if err == nil {
		if res.CPUQuota != "" {
			limitsFilePath := path.Join(subsysCgroupPath, cpuQuotaFileName)
			err = setCPULimit(res.CPUQuota, limitsFilePath)
		}
	}

	return err
}

/*
Remove used to remove subsystem by the cgroup path
Params: cgroupPath string
Return: error
*/
func (subsys *CpuSubSystem) Remove(cgroupPath string) error {
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
func (subsys *CpuSubSystem) Apply(cgroupPath string, pid int) error {
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
