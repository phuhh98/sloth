package workspace

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/phuhh98/sloth/packages/cli/pkg/lock"
)

const (
	RootDirName      = ".sloth"
	ContractsDirName = "contracts"
	SetsDirName      = "sets"
	ManifestsDirName = "manifests"
	LockFileName     = "lock.json"
)

func RootPath(workingDir string) string {
	return filepath.Join(workingDir, RootDirName)
}

func ContractsDir(workingDir string) string {
	return filepath.Join(RootPath(workingDir), ContractsDirName)
}

func SetsDir(workingDir string) string {
	return filepath.Join(RootPath(workingDir), SetsDirName)
}

func LockPath(workingDir string) string {
	return filepath.Join(RootPath(workingDir), ManifestsDirName, LockFileName)
}

func Init(workingDir string) error {
	dirs := []string{
		RootPath(workingDir),
		ContractsDir(workingDir),
		SetsDir(workingDir),
		filepath.Join(RootPath(workingDir), ManifestsDirName),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create workspace directory %q: %w", dir, err)
		}
	}

	lp := LockPath(workingDir)
	if _, err := os.Stat(lp); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("check lock file: %w", err)
		}
		if err := lock.Write(lp, lock.New()); err != nil {
			return err
		}
	}

	return nil
}

func SaveContract(workingDir string, payload map[string]any) (string, string, string, error) {
	name, _ := payload["name"].(string)
	version, _ := payload["version"].(string)
	if strings.TrimSpace(name) == "" || strings.TrimSpace(version) == "" {
		return "", "", "", fmt.Errorf("contract requires name and version")
	}

	outPath := filepath.Join(ContractsDir(workingDir), fmt.Sprintf("%s@%s.json", name, version))
	raw, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return "", "", "", fmt.Errorf("marshal contract payload: %w", err)
	}
	raw = append(raw, '\n')

	h := sha256.Sum256(raw)
	hash := hex.EncodeToString(h[:])
	if err := os.WriteFile(outPath, raw, 0o644); err != nil {
		return "", "", "", fmt.Errorf("write contract file: %w", err)
	}

	schemaVersion, _ := payload["schemaVersion"].(string)
	return outPath, hash, schemaVersion, nil
}

func LoadContractFile(path string) (map[string]any, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read contract file: %w", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, fmt.Errorf("parse contract file: %w", err)
	}
	return payload, nil
}

func LocalContractNames(workingDir string) ([]string, error) {
	entries, err := os.ReadDir(ContractsDir(workingDir))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read contracts directory: %w", err)
	}

	names := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		name := strings.Split(entry.Name(), "@")[0]
		names = append(names, name)
	}
	sort.Strings(names)
	return names, nil
}
