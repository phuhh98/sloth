package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	digest "github.com/opencontainers/go-digest"
	"github.com/opencontainers/image-spec/specs-go"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/phuhh98/sloth/packages/cli/pkg/registry"
)

func TestContractsListJSON(t *testing.T) {
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
	root.SetArgs([]string{"contracts", "ls", "--version", "0.0.1", "--format", "json"})

	if err := root.Execute(); err != nil {
		t.Fatalf("execute list: %v", err)
	}

	if !bytes.Contains(buf.Bytes(), []byte("hero-banner")) {
		t.Fatalf("expected hero-banner in output: %s", buf.String())
	}
}

func TestContractsInspectJSON(t *testing.T) {
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

func TestContractsListJSONFromOCI(t *testing.T) {
	registry := startMockOCIRegistry(t)
	defer registry.Close()

	tmp := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmp, ".sloth"), 0o755); err != nil {
		t.Fatalf("mkdir .sloth: %v", err)
	}

	host := strings.TrimPrefix(registry.URL, "http://")
	config := []byte("currentProfile: default\nprofiles:\n  default:\n    host: http://localhost:1337\n    registry:\n      host: " + host + "\n      repository: org/contracts\n      useAuthorizationToken: false\n")
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
	root.SetArgs([]string{"contracts", "--source", "oci", "--version", "latest", "ls", "--format", "json"})

	if err := root.Execute(); err != nil {
		t.Fatalf("execute OCI list: %v", err)
	}

	out := map[string]any{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("parse output json: %v", err)
	}
	if out["source"] != "oci" {
		t.Fatalf("expected source=oci, got: %v", out["source"])
	}
	rawContracts, ok := out["contracts"].([]any)
	if !ok || len(rawContracts) != 1 {
		t.Fatalf("expected one contract in OCI list output: %s", buf.String())
	}
}

func TestContractsAddComponentFromOCI(t *testing.T) {
	registry := startMockOCIRegistry(t)
	defer registry.Close()

	tmp := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmp, ".sloth"), 0o755); err != nil {
		t.Fatalf("mkdir .sloth: %v", err)
	}

	host := strings.TrimPrefix(registry.URL, "http://")
	config := []byte("currentProfile: default\nprofiles:\n  default:\n    host: http://localhost:1337\n    registry:\n      host: " + host + "\n      repository: org/contracts\n      useAuthorizationToken: false\n")
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
	root.SetArgs([]string{"contracts", "--source", "oci", "--version", "latest", "add", "component", "--name", "hero-banner", "--format", "json"})

	if err := root.Execute(); err != nil {
		t.Fatalf("execute OCI add component: %v", err)
	}

	if !bytes.Contains(buf.Bytes(), []byte("hero-banner")) {
		t.Fatalf("expected hero-banner in add output: %s", buf.String())
	}

	contractFile := filepath.Join(tmp, ".sloth", "contracts", "hero-banner@0.0.2.json")
	if _, err := os.Stat(contractFile); err != nil {
		t.Fatalf("expected saved contract %s: %v", contractFile, err)
	}
}

func TestContractsPullFromOCIToWorkspace(t *testing.T) {
	registry := startMockOCIRegistry(t)
	defer registry.Close()

	tmp := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmp, ".sloth"), 0o755); err != nil {
		t.Fatalf("mkdir .sloth: %v", err)
	}

	host := strings.TrimPrefix(registry.URL, "http://")
	config := []byte("currentProfile: default\nprofiles:\n  default:\n    host: http://localhost:1337\n    registry:\n      host: " + host + "\n      repository: org/contracts\n      useAuthorizationToken: false\n")
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
	root.SetArgs([]string{"contracts", "--source", "oci", "--version", "latest", "pull", "--name", "hero-banner", "--format", "json"})

	if err := root.Execute(); err != nil {
		t.Fatalf("execute OCI pull: %v", err)
	}

	contractFile := filepath.Join(tmp, ".sloth", "contracts", "hero-banner@0.0.2.json")
	if _, err := os.Stat(contractFile); err != nil {
		t.Fatalf("expected pulled contract %s: %v", contractFile, err)
	}

	lockFile := filepath.Join(tmp, ".sloth", "manifests", "lock.json")
	rawLock, err := os.ReadFile(lockFile)
	if err != nil {
		t.Fatalf("read lock file: %v", err)
	}
	if !bytes.Contains(rawLock, []byte("hero-banner")) {
		t.Fatalf("expected pulled contract in lock file: %s", string(rawLock))
	}
}

func TestContractsPullFromOCIToOutPath(t *testing.T) {
	registry := startMockOCIRegistry(t)
	defer registry.Close()

	tmp := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmp, ".sloth"), 0o755); err != nil {
		t.Fatalf("mkdir .sloth: %v", err)
	}

	host := strings.TrimPrefix(registry.URL, "http://")
	config := []byte("currentProfile: default\nprofiles:\n  default:\n    host: http://localhost:1337\n    registry:\n      host: " + host + "\n      repository: org/contracts\n      useAuthorizationToken: false\n")
	if err := os.WriteFile(filepath.Join(tmp, ".sloth", "config.yaml"), config, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	oldWD, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldWD) }()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	outPath := filepath.Join(tmp, "tmp-out", "hero-banner.json")
	root := NewRootCommand()
	buf := bytes.NewBuffer(nil)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs([]string{"contracts", "--source", "oci", "--version", "latest", "pull", "--name", "hero-banner", "--out", outPath, "--format", "json"})

	if err := root.Execute(); err != nil {
		t.Fatalf("execute OCI pull with out: %v", err)
	}

	if _, err := os.Stat(outPath); err != nil {
		t.Fatalf("expected output file %s: %v", outPath, err)
	}
}

func TestContractsPullFromOCIMissingContract(t *testing.T) {
	registry := startMockOCIRegistry(t)
	defer registry.Close()

	tmp := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmp, ".sloth"), 0o755); err != nil {
		t.Fatalf("mkdir .sloth: %v", err)
	}

	host := strings.TrimPrefix(registry.URL, "http://")
	config := []byte("currentProfile: default\nprofiles:\n  default:\n    host: http://localhost:1337\n    registry:\n      host: " + host + "\n      repository: org/contracts\n      useAuthorizationToken: false\n")
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
	root.SetArgs([]string{"contracts", "--source", "oci", "--version", "latest", "pull", "--name", "missing-contract"})

	err := root.Execute()
	if err == nil {
		t.Fatalf("expected missing contract error")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Fatalf("expected not found error, got: %v", err)
	}
}

func TestContractsPullFromOCIMissingVersion(t *testing.T) {
	registry := startMockOCIRegistry(t)
	defer registry.Close()

	tmp := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmp, ".sloth"), 0o755); err != nil {
		t.Fatalf("mkdir .sloth: %v", err)
	}

	host := strings.TrimPrefix(registry.URL, "http://")
	config := []byte("currentProfile: default\nprofiles:\n  default:\n    host: http://localhost:1337\n    registry:\n      host: " + host + "\n      repository: org/contracts\n      useAuthorizationToken: false\n")
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
	root.SetArgs([]string{"contracts", "--source", "oci", "--version", "9.9.9", "pull", "--name", "hero-banner"})

	err := root.Execute()
	if err == nil {
		t.Fatalf("expected missing version error")
	}
	if !strings.Contains(err.Error(), "fetch OCI manifest") {
		t.Fatalf("expected OCI manifest fetch error, got: %v", err)
	}
}

func startMockOCIRegistry(t *testing.T) *httptest.Server {
	t.Helper()

	releaseJSON := []byte(`{"version":"0.0.2","schemaVersion":"0.0.2","contracts":[{"name":"hero-banner","label":"Hero Banner","version":"0.0.2","schemaVersion":"0.0.2","contentHash":"abc","payload":{"name":"hero-banner","label":"Hero Banner","kind":"section","version":"0.0.2","schemaVersion":"0.0.2","dataset":[{"key":"headline","label":"Headline","type":"string","required":true}],"renderMeta":{"rendererKey":"hero-banner"}}}]}`)
	layerDesc := descriptorForTest(registry.ReleaseLayerMediaType, releaseJSON)
	manifestBytes := manifestJSONForTest(layerDesc)

	payloadByKey := map[string][]byte{
		"org/contracts:0.0.2":               manifestBytes,
		"org/contracts@" + layerDesc.Digest.String(): releaseJSON,
	}

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/v2/" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if !strings.HasPrefix(path, "/v2/") {
			http.NotFound(w, r)
			return
		}

		rest := strings.TrimPrefix(path, "/v2/")
		if strings.HasSuffix(rest, "/tags/list") {
			_ = json.NewEncoder(w).Encode(map[string]any{"name": "org/contracts", "tags": []string{"0.0.2"}})
			return
		}

		if strings.Contains(rest, "/manifests/") {
			parts := strings.SplitN(rest, "/manifests/", 2)
			repo := parts[0]
			ref := parts[1]
			key := repo + ":" + ref
			if strings.HasPrefix(ref, "sha256:") {
				key = repo + "@" + ref
			}
			payload, ok := payloadByKey[key]
			if !ok {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", ocispec.MediaTypeImageManifest)
			_, _ = w.Write(payload)
			return
		}

		if strings.Contains(rest, "/blobs/") {
			parts := strings.SplitN(rest, "/blobs/", 2)
			repo := parts[0]
			d := parts[1]
			payload, ok := payloadByKey[repo+"@"+d]
			if !ok {
				http.NotFound(w, r)
				return
			}
			_, _ = w.Write(payload)
			return
		}

		http.NotFound(w, r)
	}))
}

func descriptorForTest(mediaType string, payload []byte) ocispec.Descriptor {
	d := digest.FromBytes(payload)
	return ocispec.Descriptor{
		MediaType: mediaType,
		Digest:    d,
		Size:      int64(len(payload)),
	}
}

func manifestJSONForTest(layer ocispec.Descriptor) []byte {
	manifest := ocispec.Manifest{
		Versioned: specs.Versioned{SchemaVersion: 2},
		MediaType: ocispec.MediaTypeImageManifest,
		Config:    descriptorForTest(ocispec.MediaTypeImageConfig, []byte(`{}`)),
		Layers:    []ocispec.Descriptor{layer},
	}
	raw, _ := json.Marshal(manifest)
	return raw
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
