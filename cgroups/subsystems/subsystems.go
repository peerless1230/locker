package subsystems

/*
ResourceLimitConfig indicate the set of limited Resource
*/
type ResourceLimitConfig struct {
	CpuSet       string
	CpuShare     string
	MemeryLimits string
}

var tasksFileName = "tasks"
var memoryLimitsFileName = "memory.limit_in_bytes"
var cpuSharesFileName = "cpu.shares"
var cpusetLimitsFileName = "cpuset.cpus"

/*
SubSystem Interface describe the functions of subsystem
*/
type SubSystem interface {
	// get the name of subsystem
	GetName() string
	// set a cgroup into the subsystem
	Set(cgroupPath string, res *ResourceLimitConfig) error
	// remove a cgroup from the subsystem
	Remove(cgroupPath string) error
	// add a process into the cgroup
	Apply(cgroupPath string, pid int) error
}

/*
SubSystemIns array for all kinds of Subsystems
*/
var (
	SubSystemIns = []SubSystem{
		&MemorySubSystem{},
		&CpuSubSystem{},
		&CpusetSubSystem{},
	}
)
