package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestContractsListJSON(t *testing.T) {
	t.Parallel()
	tmp := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmp, "apps", "docs", "static", "registry", "contracts", "0.0.1", "components", "hero-banner"), 0o755); err != nil {
		t.Fatalf("mkdir registry: %v", err)
	}

	writeJSON(t, filepath.Join(tmp, "apps", "docs", "static", "registry", "contracts", "index.json"), map[string]any{
		"registryFormatVersion": "1",
		"items": []map[string]any{{"version": "0.0.1"}},
	})
	writeJSON(t, filepath.Join(tmp, "apps", "docs", "static", "registry", "contracts", "0.0.1", "manifest.json"), map[string]any{
		"version":       "0.0.1",
		"schemaVersion": "0.0.1",
		"components": map[string]any{
			"hero-banner": map[string]any{"contractPath": "./components/hero-banner/contract.json", "contentHash": "abc"},
		},
	})
	writeJSON(t, filepath.Join(tmp, "apps", "docs", "static", "registry", "contracts", "0.0.1", "components", "hero-banner", "contract.json"), map[string]any{
		"name":          "hero-banner",
		"label":         "Hero Banner",
		"version":       "0.0.1",
		"schemaVersion": "0.0.1",
	})

	oldWD, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldWD) }()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	root := NewRootCommand()
	buf := bytes.NewBuffer(nil)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs([]string{"contracts", "list", "--plugin-version", "0.0.1", "--format", "json"})

	if err := root.Execute(); err != nil {
		t.Fatalf("execute list: %v", err)
	}

	if !bytes.Contains(buf.Bytes(), []byte("hero-banner")) {
		t.Fatalf("expected hero-banner in output: %s", buf.String())
	}
}

func TestContractsInspectJSON(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/sloth/inspection/plugin-status":
			_, _ = w.Write([]byte(`{"pluginName":"sloth-strapi-plugin","pluginVersion":"0.0.1","compatibleSchemaVersions":["0.0.1"],"components":[],"totalComponents":0}`))
		case "/sloth/inspection/contract-schema":
			_, _ = w.Write([]byte(`{"schemaVersion":"0.0.1","schemaUrl":"/sloth/inspection/contract-schema?schemaVersion=0.0.1&inline=true"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	tmp := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmp, ".sloth"), 0o755); err != nil {
		t.Fatalf("mkdir .sloth: %v", err)
	}
	config := []byte("currentProfile: default\nprofiles:\n  default:\n    host: " + server.URL + "\n")
	if err := os.WriteFile(filepath.Join(tmp, ".sloth", "config.yaml"), config, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	oldWD, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldWD) }()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	root := NewRootCommand()
	buf := bytes.NewBuffer(nil)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs([]string{"contracts", "inspect", "--format", "json"})

	if err := root.Execute(); err != nil {
		t.Fatalf("execute inspect: %v", err)
	}

	out := map[string]any{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("parse output json: %v", err)
	}
	if _, ok := out["status"]; !ok {
		t.Fatalf("expected status in output: %s", buf.String())
	}
}

func writeJSON(t *testing.T, path string, v any) {
	t.Helper()
	raw, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal json: %v", err)
	}
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		t.Fatalf("write file %s: %v", path, err)
	}
}
