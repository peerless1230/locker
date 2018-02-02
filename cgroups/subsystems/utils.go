package subsystems

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

	"locker/common"

	log "github.com/Sirupsen/logrus"
)

/*
FindCgroupMountpoint used to find the mount from /proc/self/mountinfo
Params: subsysName string
Return: "memory"
*/
func FindCgroupMountpoint(subsysName string) string {
	// find the mountinfo of the process
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		common.CheckError(err)
	}

	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		info := scanner.Text()
		index := strings.Index(info, subsysName)
		if index != -1 {
			fields := strings.Split(info, " ")
			return fields[4]
		}
	}

	if err = scanner.Err(); err != nil {
		common.CheckError(err)
	}

	return ""
}

/*
GetCgroupPath used to get the absolute path of cgroup, if there isn't cgroupPath,
we will create it
Params: subsysName string, cgroupPath string, autoCreate bool
Return: string, error
*/
func GetCgroupPath(subsysName string, cgroupPath string, autoCreate bool) (string, error) {
	cgroupRoot := FindCgroupMountpoint(subsysName)
	absolutePath := path.Join(cgroupRoot, cgroupPath)
	isExsits, err := common.IsPathOrFileExists(absolutePath)
	if err == nil {
		// check file not exsit and autoCreate is enabled
		// avoid to create cgroup when Remove()
		if (isExsits == false) && autoCreate {
			if err := os.Mkdir(absolutePath, 0755); err == nil {
				log.Debugf("Created Cgrouppath %s", absolutePath)
			} else {
				return "", fmt.Errorf("Error in Create cgroup: %v", err)
			}
		}
		return absolutePath, nil

	}
	return "", fmt.Errorf("Cgroup path error %v", err)
}
