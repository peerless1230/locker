package container

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"unicode"

	log "github.com/Sirupsen/logrus"
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
	cmds, err := ioutil.ReadAll(read)
	if err != nil {
		log.Errorf("commands: %s", string(cmds))

		log.Errorf("Get init commands error: %v", err)
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

	cmd := exec.Command("ls", "-l", "/proc/self/fd")
	outPipe, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("Get outputPipe error: %v", err)
	}
	log.Debugf("outPipe is %s", string(outPipe))
	cmd.Start()
	read := os.NewFile(uintptr(3), "pipe")
	cmds, err := ioutil.ReadAll(read)
	if err != nil {
		log.Errorf("commands: %s", string(cmds))

		log.Errorf("Get init commands error: %v", err)
	}
	strCmds := string(cmds)
	//cmds, err := ioutil.ReadAll(outPipe)

	/*if err != nil {
		log.Errorf("commands: %s", string(cmds))

		log.Errorf("Get init commands error: %v", err)
	}
	strCmds := string(cmds)*/
	log.Debugf("strCmds is %s", string(strCmds))

	cmd.Wait()

	/*commands := readInitCommands()
	if len(commands) == 0 {
		return fmt.Errorf("Get init commmand error, there is not any command")
	}

	log.Debugf("SetHostname")
	common.CheckError(syscall.Sethostname([]byte("test")))
	log.Debugf("Chroot")
	common.CheckError(syscall.Chroot("/home/encore/busybox"))
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
	syscall.Unmount("proc", 0)*/
	return nil
}
