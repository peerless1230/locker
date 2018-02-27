package container

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/peerless1230/locker/common"
)

const dockerLayer = "/var/lib/docker/overlay2/"
const rootLAYER = "/var/lib/locker/overlay2/"
const diffLAYER = "diff"
const workLAYER = "work"
const mergedLAYER = "merged"

// ubuntuImageLayer should be replace as the Docker Overlay2 Image you want to use.
const ubuntuImageLayer = "/var/lib/docker/overlay2/4ee6de34917e8c8afa0ce2b09ccf5bf453fba42af9829b9bad7222f3e9c0ec9d"
const linkFile = "link"
const lowerFile = "lower"

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
Params: tty bool, volumeSlice []string
Return: *exec.Cmd, *os.File
*/
func NewParentProcess(tty bool, volumeSlice []string) (*exec.Cmd, *os.File) {
	// here changed to use pipe pass params to InitProcess.
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		log.Errorf("New pipe error: %v", err)
		return nil, nil
	}

	cmd := exec.Command("/proc/self/exe", "init")

	// add SysProcAttrs set root user for container Init process
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS | syscall.CLONE_NEWUTS | syscall.CLONE_NEWUSER |
			syscall.CLONE_NEWPID | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      0,
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      0,
				Size:        1,
			},
		},
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		log.Debugf("tty is enabled")
	}
	containerID := "92745277a8b052e2c50cf757da7140afabd9f6abbae7b6d6516f944a55658dfc"
	layerPath := filepath.Join(rootLAYER, containerID)
	createOverlayLayers(layerPath)
	mountVolumes(volumeSlice, layerPath)
	mergedDir := filepath.Join(layerPath, mergedLAYER)
	cmd.Dir = mergedDir
	cmd.ExtraFiles = []*os.File{readPipe}
	return cmd, writePipe
}

// create OverlayLayers for container's OverlayFS
func createOverlayLayers(layerPath string) {
	if err := createDiffLayer(layerPath); err != nil {
		log.Fatalf("Create diff dir for container's OverlayFS failed: %v", err)
		os.Exit(1)
	}
	if err := createWorkLayer(layerPath); err != nil {
		log.Fatalf("Create work dir for container's OverlayFS failed: %v", err)
		os.Exit(2)
	}
	if err := createMergedLayer(layerPath); err != nil {
		log.Fatalf("Create merged dir for container's OverlayFS failed: %v", err)
		os.Exit(3)
	}
	mountOverlayFS(layerPath)
}

// create upper dir for container's OverlayFS
func createDiffLayer(layerPath string) error {
	diffLayer := filepath.Join(layerPath, diffLAYER)
	if isExist, _ := common.IsPathOrFileExists(diffLayer); isExist == false {
		if err := common.MkdirAll(diffLayer, 0777); err != nil {
			return err
		}
	}
	return nil
}

// create work dir for container's OverlayFS
func createWorkLayer(layerPath string) error {
	workLayer := filepath.Join(layerPath, workLAYER)
	if isExist, _ := common.IsPathOrFileExists(workLayer); isExist == false {
		if err := common.MkdirAll(workLayer, 0777); err != nil {
			return err
		}
	}
	return nil
}

// create mount dir for container's OverlayFS
func createMergedLayer(layerPath string) error {
	mergedLayer := filepath.Join(layerPath, mergedLAYER)
	if isExist, _ := common.IsPathOrFileExists(mergedLayer); isExist == false {
		if err := common.MkdirAll(mergedLayer, 0777); err != nil {
			return err
		}
	}
	return nil
}

//  mount container's OverlayFS on merged dir
func mountOverlayFS(layerPath string) {
	linkFilePath := path.Join(ubuntuImageLayer, linkFile)
	linkByte, err := ioutil.ReadFile(linkFilePath)
	if err != nil {
		log.Fatalf("Read Image layer's link Failed: %v", err)
		os.Exit(4)
	}

	link := string(linkByte)
	link = dockerLayer + "l/" + link
	lowerFilePath := path.Join(ubuntuImageLayer, lowerFile)
	lowerByte, err := ioutil.ReadFile(lowerFilePath)
	if err != nil {
		log.Fatalf("Read Image layer's lower Failed: %v", err)
		os.Exit(5)
	}
	lower := string(lowerByte)
	lowerArray := strings.Split(lower, ":")
	log.Debugf("LowerArray: %v", lowerArray)
	log.Debugf("LowerArray Length: %v", len(lowerArray))

	for i, ele := range lowerArray {
		tmpStr := dockerLayer + ele
		lowerArray[i] = tmpStr
	}
	linkSlice := []string{link}
	allLower := append(linkSlice, lowerArray...)
	strLower := strings.Join(allLower, ":")
	log.Debugf("Lower: %s", strLower)

	workLayer := filepath.Join(layerPath, workLAYER)
	mergedLayer := filepath.Join(layerPath, mergedLAYER)
	diffLayer := filepath.Join(layerPath, diffLAYER)

	mountStr := "lowerdir=%s,upperdir=%s,workdir=%s"
	mountCmd := fmt.Sprintf(mountStr, strLower, diffLayer, workLayer)

	log.Debugf("Mount OverlayFS")
	mountFlags := syscall.MS_RELATIME
	if syscall.Mount("overlayfs", mergedLayer, "overlay", uintptr(mountFlags), mountCmd); err != nil {
		log.Errorf("Mount overlayfs error: %v", err)
		os.Exit(7)
	}
}

/*
CleanUpOverlayFS is used to umount, then remove OverlayFS layers.
Params: layerPath string
Return:
*/
func CleanUpOverlayFS(layerPath string) {
	umountOverlayFS(layerPath)
	removeOverlayFS(layerPath)
}

// umount OverlayFS layers
func umountOverlayFS(layerPath string) {
	mergedLayer := filepath.Join(layerPath, mergedLAYER)

	if err := syscall.Unmount(mergedLayer, syscall.MNT_DETACH); err != nil {
		log.Errorf("Unmount container's OverlayFS error: %v", err)
	}
}

// remove OverlayFS layers
func removeOverlayFS(layerPath string) {
	mergedLayer := filepath.Join(layerPath, mergedLAYER)
	common.RmdirAll(mergedLayer)
	common.RmdirAll(layerPath)
}

//  mount container's OverlayFS on merged dir
func mountVolumes(volumeSlice []string, layerPath string) {
	volumeMap := parseVolumeSlice(volumeSlice, layerPath)
	createVolumeDirs(volumeMap)
}

// parse the slice of volumes from cli flags
func parseVolumeSlice(volumeSlice []string, layerPath string) map[string]string {
	mergedLayer := filepath.Join(layerPath, mergedLAYER)

	var volumes = make(map[string]string)
	for _, ele := range volumeSlice {
		volTemp := strings.Split(ele, ":")
		volumes[volTemp[0]] = filepath.Join(mergedLayer, volTemp[1])
	}
	return volumes
}

// create Volumes for container
func createVolumeDirs(volumeMap map[string]string) {
	for k, v := range volumeMap {
		parentDir := k
		if isExist, _ := common.IsPathOrFileExists(parentDir); isExist == false {
			common.MkdirAll(parentDir, 0755)
		}
		containerDir := v
		if isExist, _ := common.IsPathOrFileExists(containerDir); isExist == false {
			common.MkdirAll(containerDir, 0755)
		}
		mountFlags := syscall.MS_RELATIME | syscall.MS_BIND | syscall.MS_REC

		syscall.Mount(parentDir, containerDir, "bind", uintptr(mountFlags), "")
	}

}

/*
CleanUpVolumes is used to umount, then remove
Params: volumes map[string]string
Return:
*/
func CleanUpVolumes(volumes []string, layerPath string) {
	volumeMap := parseVolumeSlice(volumes, layerPath)
	umountVolumes(volumeMap)
	removeVolumes(volumeMap)
}

// umount Volumes
func umountVolumes(volumeMap map[string]string) {
	for _, v := range volumeMap {
		log.Debugf("Unmount the container's volume on: %v.", v)
		if err := syscall.Unmount(v, syscall.MNT_DETACH); err != nil {
			log.Errorf("Unmount container's volume error: %v", err)
		}
	}
}

// remove Volumes
func removeVolumes(volumeMap map[string]string) {
	for _, v := range volumeMap {
		// remove dirs of volumes
		common.RmdirAll(v)
	}
}
