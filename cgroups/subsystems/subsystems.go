package subsystems

/*
ResourceLimitConfig indicate the set of limited Resource
*/
type ResourceLimitConfig struct {
	CpuSet       string
	CpuShare     string
	MemeryLimits string
}

/*
Subsystem Interface describe the functions of subsystem
*/
type Subsystem interface {
	// get the name of subsystem
	GetName() string
	// set a cgroup into the subsystem
	Set(path string, res *ResourceLimitConfig) error
	// remove a cgroup from the subsystem
	Remove(path string) error
	// add a process into the cgroup
	Apply(path string, pid int) error
}

var (
	SubsystemIns = []Subsystem{
		&MemorySubSystem{},
	}
)
