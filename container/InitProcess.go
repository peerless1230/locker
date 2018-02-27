package container

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"unicode"

	log "github.com/Sirupsen/logrus"
	"github.com/peerless1230/locker/common"
)

const tempRootfs = ".put_old" // temp path for put_out rootfs

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

	setUpSystemMount()

	log.Debugf("SetHostname to locker-container")
	common.CheckError(syscall.Sethostname([]byte("locker-container")))

	// find the command's path from system PATH
	cmdPath, err := exec.LookPath(commands[0])
	if err != nil {
		log.Errorf("Find command %s from PATH error: %v", commands[0], err)
		return err
	}
	log.Debugf("Found command in %s", cmdPath)
	syscall.Exec(cmdPath, commands, os.Environ())
	log.Infof("Umount /proc")
	syscall.Unmount("proc", syscall.MNT_DETACH)
	return nil
}

/*
pivotRoot use syscall.pivotRoot to change the rootfs
then umount put_old.
Params: rootPath string
Return: error
*/
func pivotRoot(rootPath string) error {
	// remount new-rootfs and change it's fstype
	log.Debugf("Remouning new-rootfs.")

	if err := syscall.Mount(rootPath, rootPath, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("Remount rootfs to itself error: %v", err)
	}
	pivotPath := filepath.Join(rootPath, tempRootfs)

	if err := common.MkdirAll(pivotPath, 0755); err != nil {
		return err
	}

	log.Debugf("Starting PivotRoot.")
	if err := syscall.PivotRoot(rootPath, pivotPath); err != nil {
		return fmt.Errorf("PivotRoot %s to %s error: %v", rootPath, pivotPath, err)
	}
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("Chdir to / error: %v", err)
	}

	pivotPath = filepath.Join("/", tempRootfs)
	// umount old_rootfs with waiting it's busy state flag
	if err := syscall.Unmount(pivotPath, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("Unmount old_rootfs error: %v", err)
	}
	err := syscall.Rmdir(pivotPath)
	log.Debugf("Removing %s", pivotPath)
	return err
}

/*
setUpSystemMount mount /proc and /dev
Params:
Return: error
*/
func setUpSystemMount() {

	pwd, err := os.Getwd()
	if err != nil {
		log.Errorf("Get current dir error: %v", err)
	}
	log.Debugf("Local dir is %s", pwd)

	mountFlags := syscall.MS_NOEXEC | syscall.MS_NODEV | syscall.MS_NOSUID
	log.Debugf("Mount /proc")
	if err := syscall.Mount("proc", filepath.Join(pwd, "proc"), "proc", uintptr(mountFlags), ""); err != nil {
		log.Debugf("Mount /proc error: %v", err)

	}

	log.Debugf("Mount /dev")
	mountFlags = syscall.MS_STRICTATIME | syscall.MS_NOSUID
	if syscall.Mount("tmpfs", filepath.Join(pwd, "dev"), "tmpfs", uintptr(mountFlags), "mode=0755"); err != nil {
		log.Debugf("Mount /dev error: %v", err)

	}

	mountFlags = syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV | syscall.MS_RDONLY
	log.Debugf("Mount /sys")
	if err := syscall.Mount("sys", filepath.Join(pwd, "sys"), "sysfs", uintptr(mountFlags), ""); err != nil {
		log.Debugf("Mount /sys error: %v", err)

	}

	pivotRoot(pwd)

}
