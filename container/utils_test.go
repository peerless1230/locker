package container

import (
	"path/filepath"
	"testing"
)

func TestCreateDiffLayer(t *testing.T) {
	containerID := "95741597a8b052e2c50cf757da7140afabd9f6abbae7b6d6516f944a55658dfc"
	layerPath := filepath.Join(rootLAYER, containerID)
	err := createDiffLayer(layerPath)
	if err != nil {
		t.Fatalf("createDiffLayer failed: %v", err)
	}
}

func TestCreateMergedLayer(t *testing.T) {
	containerID := "95741597a8b052e2c50cf757da7140afabd9f6abbae7b6d6516f944a55658dfc"
	layerPath := filepath.Join(rootLAYER, containerID)
	err := createMergedLayer(layerPath)
	if err != nil {
		t.Fatalf("createMergedLayer failed: %v", err)
	}
}
