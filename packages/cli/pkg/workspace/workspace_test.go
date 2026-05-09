package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitAndSaveContract(t *testing.T) {
	t.Parallel()
	tmp := t.TempDir()
	if err := Init(tmp); err != nil {
		t.Fatalf("init workspace: %v", err)
	}

	if _, err := os.Stat(LockPath(tmp)); err != nil {
		t.Fatalf("lock file missing: %v", err)
	}

	path, hash, schemaVersion, err := SaveContract(tmp, map[string]any{
		"name":          "hero-banner",
		"version":       "1.0.0",
		"schemaVersion": "0.0.1",
		"label":         "Hero Banner",
	})
	if err != nil {
		t.Fatalf("save contract: %v", err)
	}
	if hash == "" || schemaVersion != "0.0.1" {
		t.Fatalf("unexpected save metadata hash=%s schema=%s", hash, schemaVersion)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("saved contract missing: %v", err)
	}

	names, err := LocalContractNames(tmp)
	if err != nil {
		t.Fatalf("local names: %v", err)
	}
	if len(names) != 1 || names[0] != "hero-banner" {
		t.Fatalf("unexpected names: %v", names)
	}

	if _, err := LoadContractFile(filepath.Join(ContractsDir(tmp), "hero-banner@1.0.0.json")); err != nil {
		t.Fatalf("load saved contract: %v", err)
	}
}
