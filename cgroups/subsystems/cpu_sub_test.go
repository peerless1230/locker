package subsystems

import (
	"os"
	"testing"
)

var testCPUSub = &CpuSubSystem{}

/*
TestGetName used to test the return of CpuSubSystem.GetName
*/
func TestCpuCgroup(t *testing.T) {
	if testCPUSub.GetName() != "cpu" {
		t.FailNow()
	}
	err := testCPUSub.Set(testCgroup, &testResConfig)
	if err != nil {
		t.Fatalf("Set cgroup error: %v", err)
	}
	if err := testCPUSub.Apply(testCgroup, os.Getpid()); err != nil {
		t.Fatalf("Apply cgroup error: %v", err)
	}
	if err := testCPUSub.Set(testSecondCgroup, &testResConfig); err != nil {
		t.Fatalf("Set cgroup error: %v", err)
	}

	if err := testCPUSub.Apply(testSecondCgroup, os.Getpid()); err != nil {
		t.Fatalf("Apply cgroup error: %v", err)
	}
	// move the process to parent cgroup, otherwise we can't remove the childs.
	if err := testCPUSub.Apply("", os.Getpid()); err != nil {
		t.Fatalf("Apply cgroup error: %v", err)
	}
	if err := testCPUSub.Remove(testCgroup); err != nil {
		t.Fatalf("Remove cgroup error: %v", err)
	}
	if err := testCPUSub.Remove(testSecondCgroup); err != nil {
		t.Fatalf("Remove cgroup error: %v", err)
	}

}
