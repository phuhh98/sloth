package lock

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const CurrentFormatVersion = "1"

type Entry struct {
	Name          string `json:"name"`
	Version       string `json:"version"`
	SchemaVersion string `json:"schemaVersion"`
	Source        string `json:"source"`
	ContentHash   string `json:"contentHash"`
	LastSyncedAt  string `json:"lastSyncedAt,omitempty"`
	RemoteID      string `json:"remoteId,omitempty"`
}

type File struct {
	FormatVersion string           `json:"formatVersion"`
	UpdatedAt     string           `json:"updatedAt"`
	Entries       map[string]Entry `json:"entries"`
}

func New() *File {
	return &File{
		FormatVersion: CurrentFormatVersion,
		UpdatedAt:     time.Now().UTC().Format(time.RFC3339),
		Entries:       map[string]Entry{},
	}
}

func Read(path string) (*File, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return New(), nil
		}
		return nil, fmt.Errorf("read lock file: %w", err)
	}

	f := &File{}
	if err := json.Unmarshal(raw, f); err != nil {
		return nil, fmt.Errorf("parse lock file: %w", err)
	}

	if f.FormatVersion == "" {
		f.FormatVersion = CurrentFormatVersion
	}
	if f.Entries == nil {
		f.Entries = map[string]Entry{}
	}
	if f.UpdatedAt == "" {
		f.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	}

	return f, nil
}

func (f *File) Upsert(entry Entry) {
	if f.Entries == nil {
		f.Entries = map[string]Entry{}
	}
	f.Entries[entry.Name] = entry
	f.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
}

func (f *File) Merge(other *File) {
	if other == nil {
		return
	}
	for _, key := range other.SortedKeys() {
		f.Upsert(other.Entries[key])
	}
}

func (f *File) SortedKeys() []string {
	keys := make([]string, 0, len(f.Entries))
	for key := range f.Entries {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func Write(path string, f *File) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create lock directory: %w", err)
	}

	raw, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal lock file: %w", err)
	}
	raw = append(raw, '\n')

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, raw, 0o644); err != nil {
		return fmt.Errorf("write lock temp file: %w", err)
	}

	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("replace lock file: %w", err)
	}

	return nil
}
