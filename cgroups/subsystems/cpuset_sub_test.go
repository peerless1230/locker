package subsystems

import (
	"os"
	"testing"
)

var testCPUSetSub = &CpusetSubSystem{}

/*
TestGetName used to test the return of CPUSetSubSystem.GetName
*/
func TestCPUSetCgroup(t *testing.T) {
	if testCPUSetSub.GetName() != "cpuset" {
		t.FailNow()
	}
	err := testCPUSetSub.Set(testCgroup, &testResConfig)
	if err != nil {
		t.Fatalf("Set cgroup error: %v", err)
	}
	if err := testCPUSetSub.Apply(testCgroup, os.Getpid()); err != nil {
		t.Fatalf("Apply cgroup error: %v", err)
	}
	if err := testCPUSetSub.Set(testSecondCgroup, &testResConfig); err != nil {
		t.Fatalf("Set cgroup error: %v", err)
	}

	if err := testCPUSetSub.Apply(testSecondCgroup, os.Getpid()); err != nil {
		t.Fatalf("Apply cgroup error: %v", err)
	}
	// move the process to parent cgroup, otherwise we can't remove the childs.
	if err := testCPUSetSub.Apply("", os.Getpid()); err != nil {
		t.Fatalf("Apply cgroup error: %v", err)
	}
	if err := testCPUSetSub.Remove(testCgroup); err != nil {
		t.Fatalf("Remove cgroup error: %v", err)
	}
	if err := testCPUSetSub.Remove(testSecondCgroup); err != nil {
		t.Fatalf("Remove cgroup error: %v", err)
	}

}
