package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/opencontainers/image-spec/specs-go"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/phuhh98/sloth/packages/cli/pkg/registry"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/registry/remote"
)

func TestContractsListAndPullWithZot(t *testing.T) {
	if os.Getenv("SLOTH_RUN_ZOT_INTEGRATION") != "1" {
		t.Skip("set SLOTH_RUN_ZOT_INTEGRATION=1 to run Zot integration test")
	}
	if _, err := exec.LookPath("docker"); err != nil {
		t.Skipf("docker not available: %v", err)
	}

	repoRoot, err := findRepoRootFromCurrentFile()
	if err != nil {
		t.Fatalf("resolve repo root: %v", err)
	}

	if out, err := runInDir(repoRoot, "docker", "compose", "-f", "docker-compose.oci-registry.yaml", "--profile", "oci-test", "up", "-d"); err != nil {
		t.Fatalf("start zot compose service: %v\n%s", err, out)
	}
	t.Cleanup(func() {
		_, _ = runInDir(repoRoot, "docker", "compose", "-f", "docker-compose.oci-registry.yaml", "--profile", "oci-test", "down")
	})

	if err := waitForRegistry("http://127.0.0.1:5001/v2/"); err != nil {
		t.Fatalf("wait for zot registry: %v", err)
	}
	releaseVersion := fmt.Sprintf("0.9.9-zot-%d", time.Now().UnixNano())

	if err := seedZotRelease("localhost:5001", "org/contracts", releaseVersion); err != nil {
		t.Fatalf("seed zot release: %v", err)
	}

	tmp := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmp, ".sloth"), 0o755); err != nil {
		t.Fatalf("mkdir .sloth: %v", err)
	}
	config := []byte("currentProfile: default\nprofiles:\n  default:\n    host: http://localhost:1337\n    registry:\n      host: localhost:5001\n      repository: org/contracts\n      useAuthorizationToken: false\n")
	if err := os.WriteFile(filepath.Join(tmp, ".sloth", "config.yaml"), config, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	oldWD, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldWD) }()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	root := NewRootCommand()
	listOut := bytes.NewBuffer(nil)
	root.SetOut(listOut)
	root.SetErr(listOut)
	root.SetArgs([]string{"contracts", "--source", "oci", "--version", releaseVersion, "ls", "--format", "json"})
	if err := root.Execute(); err != nil {
		t.Fatalf("execute list against zot: %v", err)
	}

	parsed := map[string]any{}
	if err := json.Unmarshal(listOut.Bytes(), &parsed); err != nil {
		t.Fatalf("parse list output: %v", err)
	}
	rawContracts, ok := parsed["contracts"].([]any)
	if !ok || len(rawContracts) != 1 {
		t.Fatalf("expected one contract from zot list: %s", listOut.String())
	}

	root = NewRootCommand()
	pullOut := bytes.NewBuffer(nil)
	root.SetOut(pullOut)
	root.SetErr(pullOut)
	root.SetArgs([]string{"contracts", "--source", "oci", "--version", releaseVersion, "pull", "--name", "hero-banner", "--format", "json"})
	if err := root.Execute(); err != nil {
		t.Fatalf("execute pull against zot: %v", err)
	}

	if _, err := os.Stat(filepath.Join(tmp, ".sloth", "contracts", fmt.Sprintf("hero-banner@%s.json", releaseVersion))); err != nil {
		t.Fatalf("expected pulled contract file: %v", err)
	}
}

func findRepoRootFromCurrentFile() (string, error) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("runtime caller unavailable")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "..", "..")), nil
}

func runInDir(dir string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func waitForRegistry(url string) error {
	deadline := time.Now().Add(20 * time.Second)
	client := &http.Client{Timeout: 1 * time.Second}
	for time.Now().Before(deadline) {
		resp, err := client.Get(url)
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("registry not ready at %s", url)
}

func seedZotRelease(host string, repository string, version string) error {
	ctx := context.Background()
	repo, err := remote.NewRepository(fmt.Sprintf("%s/%s", host, repository))
	if err != nil {
		return fmt.Errorf("new remote repository: %w", err)
	}
	repo.PlainHTTP = true

	releaseJSON, err := json.Marshal(map[string]any{
		"version":       version,
		"schemaVersion": "0.0.1",
		"contracts": []map[string]any{
			{
				"name":          "hero-banner",
				"label":         "Hero Banner",
				"version":       version,
				"schemaVersion": "0.0.1",
				"contentHash":   "abc",
				"payload": map[string]any{
					"name":          "hero-banner",
					"label":         "Hero Banner",
					"kind":          "section",
					"version":       version,
					"schemaVersion": "0.0.1",
					"dataset": []map[string]any{
						{"key": "headline", "label": "Headline", "type": "string", "required": true},
					},
					"renderMeta": map[string]any{"rendererKey": "hero-banner"},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("marshal release payload: %w", err)
	}

	layerDesc, err := oras.PushBytes(ctx, repo, registry.ReleaseLayerMediaType, releaseJSON)
	if err != nil {
		return fmt.Errorf("push release payload blob: %w", err)
	}

	configDesc, err := oras.PushBytes(ctx, repo, ocispec.MediaTypeImageConfig, []byte(`{}`))
	if err != nil {
		return fmt.Errorf("push config blob: %w", err)
	}

	manifest := ocispec.Manifest{
		Versioned: specs.Versioned{SchemaVersion: 2},
		MediaType: ocispec.MediaTypeImageManifest,
		Config:    configDesc,
		Layers:    []ocispec.Descriptor{layerDesc},
	}
	manifestBytes, err := json.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("marshal manifest: %w", err)
	}

	if _, err := oras.TagBytes(ctx, repo, ocispec.MediaTypeImageManifest, manifestBytes, version); err != nil {
		return fmt.Errorf("tag manifest for version %s: %w", version, err)
	}

	_, remoteManifestBytes, err := oras.FetchBytes(ctx, repo, version, oras.DefaultFetchBytesOptions)
	if err != nil {
		return fmt.Errorf("verify fetch tagged manifest %s: %w", version, err)
	}

	remoteManifest := &ocispec.Manifest{}
	if err := json.Unmarshal(remoteManifestBytes, remoteManifest); err != nil {
		return fmt.Errorf("verify parse tagged manifest %s: %w", version, err)
	}
	if len(remoteManifest.Layers) == 0 {
		return fmt.Errorf("verify tagged manifest %s has no layers", version)
	}

	return nil
}
