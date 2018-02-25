package container

import (
	"path/filepath"
	"testing"

	"github.com/peerless1230/locker/common"
)

func TestCreateDiffLayer(t *testing.T) {
	containerID := "95741597a8b052e2c50cf757da7140afabd9f6abbae7b6d6516f944a55658dfc"
	layerPath := filepath.Join(rootLAYER, containerID)
	err := createDiffLayer(layerPath)
	if err != nil {
		t.Fatalf("createDiffLayer failed: %v", err)
	}
	diffLayer := filepath.Join(layerPath, diffLAYER)
	common.RmdirAll(diffLayer)
}

func TestCreateMergedLayer(t *testing.T) {
	containerID := "95741597a8b052e2c50cf757da7140afabd9f6abbae7b6d6516f944a55658dfc"
	layerPath := filepath.Join(rootLAYER, containerID)
	err := createMergedLayer(layerPath)
	if err != nil {
		t.Fatalf("createMergedLayer failed: %v", err)
	}
	mergedLayer := filepath.Join(layerPath, mergedLAYER)
	common.RmdirAll(mergedLayer)
}
