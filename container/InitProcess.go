package container

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unicode"

	log "github.com/Sirupsen/logrus"
	"github.com/peerless1230/locker/common"
)

/*
NewPipe is used to create a pipe, if the pipe create failed,
it return nil, nil, err to avoid more risks on pipe operations.
Params:
Return: *os.File, *os.File, error
*/
func readInitCommands() []string {
	// use extra fd to get the init pipe
	read := os.NewFile(uintptr(3), "pipe")
	defer read.Close()
	cmds, err := ioutil.ReadAll(read)
	if err != nil {
		log.Errorf("commands: %s", string(cmds))

		log.Errorf("Get init commands error: %v", err)
		return nil
	}
	strCmds := string(cmds)
	return strings.FieldsFunc(strCmds, unicode.IsSpace)
}

/*
NewContainerInitProcess is used to init the container's environment
and mount /proc
Params:
Return: error
*/
func NewContainerInitProcess() error {
	/* find the extra_file's fd
	cmd := exec.Command("ls", "-l", "/proc/self/fd")
	outPipe, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("Get outputPipe error: %v", err)
	}
	log.Debugf("outPipe is %s", string(outPipe))
	cmd.Start()

	cmd.Wait()*/

	commands := readInitCommands()
	if len(commands) == 0 || commands == nil {
		return fmt.Errorf("Get init commmand error, there is not any command")
	}

	log.Debugf("SetHostname")
	common.CheckError(syscall.Sethostname([]byte("locker-container")))
	log.Debugf("Chroot")
	common.CheckError(syscall.Chroot("/home/encore/alpine_stress"))
	common.CheckError(os.Chdir("/"))
	mountFlags := syscall.MS_NOEXEC | syscall.MS_NODEV | syscall.MS_NOSUID
	log.Infof("Mount /proc")
	syscall.Mount("proc", "proc", "proc", uintptr(mountFlags), "")
	// find the command's path from system PATH
	cmdPath, err := exec.LookPath(commands[0])
	if err != nil {
		log.Errorf("Find command %s from PATH error: %v", commands[0], err)
		return err
	}
	log.Debugf("Found command in %s", cmdPath)
	syscall.Exec(cmdPath, commands, os.Environ())
	log.Infof("Umount /proc")
	syscall.Unmount("proc", 0)
	return nil
}
