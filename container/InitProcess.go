package container

import (
	"os"
	"syscall"

	"../common"
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
	common.Check(syscall.Sethostname([]byte("test")))
	log.Infof("Chroot")
	common.Check(syscall.Chroot("/home/encore/busybox"))
	common.Check(os.Chdir("/"))
	mountFlags := syscall.MS_NOEXEC | syscall.MS_NODEV | syscall.MS_NOSUID
	log.Infof("Mount /proc")
	syscall.Mount("proc", "proc", "proc", uintptr(mountFlags), "")
	//argv := splitArgs(args)
	//cmd := exec.Command(command)
	cmd := "/bin/sh"
	log.Infof("Exec cmd: %s", cmd)
	log.Infof("Exec args: %s", args)
	syscall.Exec(cmd, args, os.Environ())
	//common.Check(cmd.Run())
	log.Infof("Umount /proc")
	syscall.Unmount("proc", 0)
	return nil
}
