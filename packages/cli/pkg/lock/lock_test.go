package lock

import (
	"path/filepath"
	"testing"
)

func TestReadWriteAndMerge(t *testing.T) {
	t.Parallel()
	tmp := t.TempDir()
	path := filepath.Join(tmp, "lock.json")

	f, err := Read(path)
	if err != nil {
		t.Fatalf("read missing lock: %v", err)
	}
	f.Upsert(Entry{Name: "hero-banner", Version: "1.0.0", SchemaVersion: "0.0.1"})
	if err := Write(path, f); err != nil {
		t.Fatalf("write lock: %v", err)
	}

	reloaded, err := Read(path)
	if err != nil {
		t.Fatalf("reload lock: %v", err)
	}
	if _, ok := reloaded.Entries["hero-banner"]; !ok {
		t.Fatalf("expected hero-banner entry")
	}

	reloaded.Merge(&File{Entries: map[string]Entry{"feature-grid": {Name: "feature-grid", Version: "1.0.0", SchemaVersion: "0.0.1"}}})
	if len(reloaded.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(reloaded.Entries))
	}
}
