package app

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestCLIFlowWithOpenAPIMockServer(t *testing.T) {
	tmp := t.TempDir()
	mockPort, cleanup := startOpenAPIMockServer(t)
	defer cleanup()

	writeContractRegistryFixture(t, tmp)

	if err := os.MkdirAll(filepath.Join(tmp, ".sloth"), 0o755); err != nil {
		t.Fatalf("mkdir .sloth: %v", err)
	}
	config := []byte(fmt.Sprintf("currentProfile: default\nprofiles:\n  default:\n    host: http://127.0.0.1:%d\n", mockPort))
	if err := os.WriteFile(filepath.Join(tmp, ".sloth", "config.yaml"), config, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	if _, err := runCommandInDir(t, tmp, "init"); err != nil {
		t.Fatalf("init command failed: %v", err)
	}

	inspectOut, err := runCommandInDir(t, tmp, "contracts", "inspect", "--format", "json")
	if err != nil {
		t.Fatalf("inspect command failed: %v", err)
	}
	if !bytes.Contains(inspectOut, []byte("sloth-host-mock")) {
		t.Fatalf("inspect output missing mock host payload: %s", inspectOut)
	}

	listOut, err := runCommandInDir(t, tmp, "contracts", "list", "--plugin-version", "0.0.1", "--format", "json")
	if err != nil {
		t.Fatalf("list command failed: %v", err)
	}
	if !bytes.Contains(listOut, []byte("hero-banner")) {
		t.Fatalf("list output missing hero-banner: %s", listOut)
	}

	if _, err := runCommandInDir(t, tmp, "contracts", "add", "--all", "--plugin-version", "0.0.1", "--format", "json"); err != nil {
		t.Fatalf("add --all command failed: %v", err)
	}

	customContractPath := filepath.Join(tmp, ".sloth", "contracts", "custom-card@0.0.1.json")
	writeJSONFile(t, customContractPath, map[string]any{
		"name":          "custom-card",
		"label":         "Custom Card",
		"kind":          "block",
		"version":       "0.0.1",
		"schemaVersion": "0.0.1",
		"dataset": []map[string]any{
			{"key": "title", "label": "Title", "type": "string", "required": true},
		},
		"renderMeta": map[string]any{"rendererKey": "custom-card"},
	})

	if _, err := runCommandInDir(
		t,
		tmp,
		"contracts",
		"verify",
		"--file",
		customContractPath,
		"--plugin-version",
		"0.0.1",
		"--supported-range",
		">=0.0.1",
		"--format",
		"json",
	); err != nil {
		t.Fatalf("verify command failed: %v", err)
	}

	pushOut, err := runCommandInDir(
		t,
		tmp,
		"contracts",
		"push",
		"--yes-to-all",
		"--format",
		"json",
	)
	if err != nil {
		t.Fatalf("push command failed: %v", err)
	}

	payload := map[string]any{}
	if err := json.Unmarshal(pushOut, &payload); err != nil {
		t.Fatalf("parse push output: %v", err)
	}

	result, ok := payload["result"].(map[string]any)
	if !ok {
		t.Fatalf("push output missing result field: %s", pushOut)
	}
	totalReceived, ok := result["totalReceived"].(float64)
	if !ok || totalReceived < 2 {
		t.Fatalf("push output unexpected totalReceived: %s", pushOut)
	}
}

func startOpenAPIMockServer(t *testing.T) (int, func()) {
	t.Helper()
	_, currentFile, _, _ := runtime.Caller(0)
	mockScriptPath := filepath.Clean(filepath.Join(currentFile, "..", "..", "..", "..", "component-hub", "scripts", "openapi-mock-server.mjs"))

	cmd := exec.Command("node", mockScriptPath, "--port", "0")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("stdout pipe: %v", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		t.Fatalf("stderr pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatalf("start mock server: %v", err)
	}

	var stderrBuf bytes.Buffer
	go func() {
		_, _ = stderrBuf.ReadFrom(stderr)
	}()

	scanner := bufio.NewScanner(stdout)
	deadline := time.After(10 * time.Second)
	port := 0

	for port == 0 {
		select {
		case <-deadline:
			_ = cmd.Process.Kill()
			t.Fatalf("mock server did not become ready. stderr=%s", stderrBuf.String())
		default:
			if !scanner.Scan() {
				continue
			}
			line := scanner.Text()
			if strings.HasPrefix(line, "SLOTH_MOCK_SERVER_READY:") {
				_, _ = fmt.Sscanf(line, "SLOTH_MOCK_SERVER_READY:%d", &port)
			}
		}
	}

	client := &http.Client{Timeout: 1 * time.Second}
	for attempts := 0; attempts < 30; attempts++ {
		resp, err := client.Get(fmt.Sprintf("http://127.0.0.1:%d/healthz", port))
		if err == nil && resp.StatusCode == 200 {
			_ = resp.Body.Close()
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	cleanup := func() {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		_, _ = cmd.Process.Wait()
	}

	return port, cleanup
}

func writeContractRegistryFixture(t *testing.T, root string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(root, "apps", "docs", "static", "registry", "contracts", "0.0.1", "components", "hero-banner"), 0o755); err != nil {
		t.Fatalf("mkdir registry fixture: %v", err)
	}

	writeJSONFile(t, filepath.Join(root, "apps", "docs", "static", "registry", "contracts", "index.json"), map[string]any{
		"registryFormatVersion": "1",
		"items": []map[string]any{{"version": "0.0.1"}},
	})
	writeJSONFile(t, filepath.Join(root, "apps", "docs", "static", "registry", "contracts", "0.0.1", "manifest.json"), map[string]any{
		"version":       "0.0.1",
		"schemaVersion": "0.0.1",
		"components": map[string]any{
			"hero-banner": map[string]any{"contractPath": "./components/hero-banner/contract.json", "contentHash": "abc"},
		},
	})
	writeJSONFile(t, filepath.Join(root, "apps", "docs", "static", "registry", "contracts", "0.0.1", "components", "hero-banner", "contract.json"), map[string]any{
		"name":          "hero-banner",
		"label":         "Hero Banner",
		"kind":          "section",
		"version":       "0.0.1",
		"schemaVersion": "0.0.1",
		"dataset": []map[string]any{
			{"key": "headline", "label": "Headline", "type": "string", "required": true},
		},
		"renderMeta": map[string]any{"rendererKey": "hero-banner"},
	})
}

func runCommandInDir(t *testing.T, workingDir string, args ...string) ([]byte, error) {
	t.Helper()
	oldWD, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldWD) }()
	if err := os.Chdir(workingDir); err != nil {
		t.Fatalf("chdir %s: %v", workingDir, err)
	}

	root := NewRootCommand()
	buf := bytes.NewBuffer(nil)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	err := root.Execute()
	return bytes.Clone(buf.Bytes()), err
}

func writeJSONFile(t *testing.T, filePath string, payload any) {
	t.Helper()
	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal json for %s: %v", filePath, err)
	}
	if err := os.WriteFile(filePath, raw, 0o644); err != nil {
		t.Fatalf("write json file %s: %v", filePath, err)
	}
}
