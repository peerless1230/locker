package subsystems

import (
	"bufio"
	"os"
	"strings"

	"../../common"
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
		common.Check(err)
	}

	defer f.Close()
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		info := scanner.Text()
		index := strings.Index(info, subsysName)
		if index != -1 {
			fields := strings.Split(subsysName, " ")
			return fields[4]
		}
	}

	if err = scanner.Err(); err != nil {
		common.Check(err)
	}

	return ""
}

/*
GetCgroupPath used to get the absolute path of cgroup
Params: subsysName string, cgroupPath string
Return: "memory"
*/
func GetCgroupPath(subsysName string, cgroupPath string) (string, error) {
	cgroupRoot := FindCgroupMountpoint(subsysName)
	return "", nil
}
