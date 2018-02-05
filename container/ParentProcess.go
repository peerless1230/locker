package container

import (
	"os"
	"os/exec"
	"syscall"

	log "github.com/Sirupsen/logrus"
)

/*
NewPipe is used to create a pipe, if the pipe create failed,
it return nil, nil, err to avoid more risks on pipe operations.
Params:
Return: *os.File, *os.File, error
*/
func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err == nil {
		return read, write, err
	}
	return nil, nil, err
}

/*
NewParentProcess is used to create the parent process of container
Params: tty bool
Return: *exec.Cmd, *os.File
*/
func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {
	// here changed to use pipe pass params to InitProcess.
	writePipe, readPipe, err := NewPipe()
	if err != nil {
		log.Errorf("New pipe error: %v", err)
		return nil, nil
	}

	cmd := exec.Command("/proc/self/exe", "init")

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		log.Debugf("tty is enabled")
	}
	cmd.ExtraFiles = []*os.File{readPipe}
	return cmd, writePipe
}
