package cgroups

import (
	"github.com/peerless1230/locker/cgroups/subsystems"

	log "github.com/Sirupsen/logrus"
)

/*
CgroupManager implemented Subsystem interface
*/
type CgroupManager struct {
	// relative path to root hierarchy
	Path string
	// config of ResourceLimit in a struct
	Resource *subsystems.ResourceLimitConfig
}

/*
NewCgroupManager return a new CgroupManager by path
Params: path string
Return: *CgroupManager
*/
func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{
		Path: path,
	}
}

/*
Set all kinds of cgroups
Params: res *subsystems.ResourceLimitConfig
Return: error
*/
func (c *CgroupManager) Set(res *subsystems.ResourceLimitConfig) error {
	for _, subSysIns := range subsystems.SubSystemIns {
		subSysIns.Set(c.Path, res)
	}
	return nil
}

/*
Apply the process into cgroup by pid
Params: pid int
Return: error
*/
func (c *CgroupManager) Apply(pid int) error {
	for _, subSysIns := range subsystems.SubSystemIns {
		subSysIns.Apply(c.Path, pid)
	}
	return nil
}

/*
Destroy all kinds of cgroups
Params:
Return: error
*/
func (c *CgroupManager) Destroy() error {
	for _, subSysIns := range subsystems.SubSystemIns {
		if err := subSysIns.Remove(c.Path); err != nil {
			log.Warnf("remove cgroup fail %v", err)
		}
	}
	return nil
}
