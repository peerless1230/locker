package container

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/peerless1230/locker/common"

	log "github.com/Sirupsen/logrus"
)

/*
NewParentProcess is used to create the parent process of container
Params: tty bool, command string, args []string
Return: error
*/
func NewParentProcess(tty bool, command string, args []string) *exec.Cmd {
	str := common.CombineArgsWithBlank(append([]string{command}, args...))
	log.Infof("After combine: %s", str)
	args = []string{"init", str}
	cmd := exec.Command("/proc/self/exe", args...)
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

	}
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}

	return cmd
}
