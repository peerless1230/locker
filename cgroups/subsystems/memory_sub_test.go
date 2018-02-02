package subsystems

import (
	"os"
	"testing"
)

var testMemSub = &MemorySubSystem{}
var testResConfig = ResourceLimitConfig{
	MemeryLimits: "200m",
}
var testCgroup = "testmemorycgroup"
var testSecondCgroup = "testsecondmemorycgroup"

/*
TestGetName used to test the return of MemorySubSystem.GetName
*/
func TestMemoryCgroup(t *testing.T) {
	if testMemSub.GetName() != "memory" {
		t.FailNow()
	}
	err := testMemSub.Set(testCgroup, &testResConfig)
	if err != nil {
		t.Fatalf("Set cgroup error: %v", err)
	}
	if err := testMemSub.Apply(testCgroup, os.Getpid()); err != nil {
		t.Fatalf("Apply cgroup error: %v", err)
	}
	if err := testMemSub.Set(testSecondCgroup, &testResConfig); err != nil {
		t.Fatalf("Set cgroup error: %v", err)
	}

	if err := testMemSub.Apply(testSecondCgroup, os.Getpid()); err != nil {
		t.Fatalf("Apply cgroup error: %v", err)
	}
	// move the process to parent cgroup, otherwise we can't remove the childs.
	if err := testMemSub.Apply("", os.Getpid()); err != nil {
		t.Fatalf("Apply cgroup error: %v", err)
	}
	if err := testMemSub.Remove(testCgroup); err != nil {
		t.Fatalf("Remove cgroup error: %v", err)
	}
	if err := testMemSub.Remove(testSecondCgroup); err != nil {
		t.Fatalf("Remove cgroup error: %v", err)
	}

}
