package container

import (
	"os"
	"syscall"

	"github.com/peerless1230/locker/common"

	log "github.com/Sirupsen/logrus"
)

/*
NewContainerInitProcess is used to init the container's environment
and mount /proc
Params: command string, args []string
Return: error
*/
func NewContainerInitProcess(command string) error {
	command = command[1 : len(command)-1]

	log.Infof("Command: %s", command)
	args := []string{"/bin/sh", "-c", command}
	log.Infof("SetHostname")
	common.CheckError(syscall.Sethostname([]byte("test")))
	log.Infof("Chroot")
	common.CheckError(syscall.Chroot("/home/encore/busybox"))
	common.CheckError(os.Chdir("/"))
	mountFlags := syscall.MS_NOEXEC | syscall.MS_NODEV | syscall.MS_NOSUID
	log.Infof("Mount /proc")
	syscall.Mount("proc", "proc", "proc", uintptr(mountFlags), "")
	//argv := splitArgs(args)
	//cmd := exec.Command(command)
	cmd := "/bin/sh"
	log.Infof("Exec cmd: %s", cmd)
	log.Infof("Exec args: %s", args)
	syscall.Exec(cmd, args, os.Environ())
	//common.CheckError(cmd.Run())
	log.Infof("Umount /proc")
	syscall.Unmount("proc", 0)
	return nil
}
