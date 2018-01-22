package subsystems

/*
MemorySubSystem implemented Subsystem interface
*/
type MemorySubSystem struct {
}

/*
GetName used to return the name of subsystem
Params:
Return: "memory"
*/
func (subsys *MemorySubSystem) GetName() string {
	return "memory"
}

/*
Set used to return the name of subsystem
Params: path string, res *ResourceLimitConfig
Return: error
*/
func (subsys *MemorySubSystem) Set(path string, res *ResourceLimitConfig) error {
	return nil
}

/*
Remove used to remove subsystem by the cgroup path
Params: path string
Return: error
*/
func (subsys *MemorySubSystem) Remove(path string) error {
	return nil
}

/*
Apply used to add a process into the cgroup
Params: path string
Return: error
*/
func (subsys *MemorySubSystem) Apply(path string, pid int) error {
	return nil
}
