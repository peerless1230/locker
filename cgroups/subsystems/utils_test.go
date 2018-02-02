package subsystems

import (
	"testing"
)

func TestFindCgroupMountpoint(t *testing.T) {
	str := FindCgroupMountpoint("memory")
	if str == "" {
		t.Fatalf("Find memory Cgroup Mountpoint error.")
	}
	t.Logf("Found Memory Cgroup Mountpoint: %s", str)
	str = FindCgroupMountpoint("cpu")
	if str == "" {
		t.Fatalf("Find cpu Cgroup Mountpoint error.")
	}
	t.Logf("Found cpu Cgroup Mountpoint: %s", str)
	str = FindCgroupMountpoint("cpuset")
	if str == "" {
		t.Fatalf("Find cpuset Cgroup Mountpoint error.")
	}
	t.Logf("Found cpuset Cgroup Mountpoint: %s", str)
}
