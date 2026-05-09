package source

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type indexFile struct {
	Items []struct {
		Version string `json:"version"`
	} `json:"items"`
}

type manifestFile struct {
	Version       string `json:"version"`
	SchemaVersion string `json:"schemaVersion"`
	Components    map[string]struct {
		ContractPath string `json:"contractPath"`
		ContentHash  string `json:"contentHash"`
	} `json:"components"`
}

type LocalRegistryResolver struct {
	RootPath string
}

func NewLocalRegistryResolver(rootPath string) *LocalRegistryResolver {
	return &LocalRegistryResolver{RootPath: rootPath}
}

func (r *LocalRegistryResolver) ListContracts(pluginVersion string) ([]Contract, error) {
	resolvedVersion := pluginVersion
	if resolvedVersion == "" {
		latest, err := r.latestVersion()
		if err != nil {
			return nil, err
		}
		resolvedVersion = latest
	}

	manifest, err := r.readManifest(resolvedVersion)
	if err != nil {
		return nil, err
	}

	items := make([]Contract, 0, len(manifest.Components))
	for name, comp := range manifest.Components {
		payload, err := r.readContractPayload(resolvedVersion, comp.ContractPath)
		if err != nil {
			return nil, err
		}
		label, _ := payload["label"].(string)
		version, _ := payload["version"].(string)
		if version == "" {
			version = manifest.Version
		}
		items = append(items, Contract{
			Name:          name,
			Label:         label,
			Version:       version,
			SchemaVersion: manifest.SchemaVersion,
			ContentHash:   comp.ContentHash,
			Payload:       payload,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})

	return items, nil
}

func (r *LocalRegistryResolver) GetContract(pluginVersion string, name string) (*Contract, error) {
	contracts, err := r.ListContracts(pluginVersion)
	if err != nil {
		return nil, err
	}
	return FindByName(contracts, name)
}

func (r *LocalRegistryResolver) latestVersion() (string, error) {
	indexPath := filepath.Join(r.RootPath, "index.json")
	raw, err := os.ReadFile(indexPath)
	if err != nil {
		return "", fmt.Errorf("read contract index: %w", err)
	}

	idx := &indexFile{}
	if err := json.Unmarshal(raw, idx); err != nil {
		return "", fmt.Errorf("parse contract index: %w", err)
	}
	if len(idx.Items) == 0 {
		return "", fmt.Errorf("contract index has no items")
	}

	return idx.Items[len(idx.Items)-1].Version, nil
}

func (r *LocalRegistryResolver) readManifest(version string) (*manifestFile, error) {
	manifestPath := filepath.Join(r.RootPath, version, "manifest.json")
	raw, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("read manifest for %s: %w", version, err)
	}

	manifest := &manifestFile{}
	if err := json.Unmarshal(raw, manifest); err != nil {
		return nil, fmt.Errorf("parse manifest for %s: %w", version, err)
	}
	return manifest, nil
}

func (r *LocalRegistryResolver) readContractPayload(version string, contractPath string) (map[string]any, error) {
	resolvedPath := filepath.Join(r.RootPath, version, contractPath)
	raw, err := os.ReadFile(resolvedPath)
	if err != nil {
		return nil, fmt.Errorf("read contract payload %s: %w", resolvedPath, err)
	}
	payload := map[string]any{}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, fmt.Errorf("parse contract payload %s: %w", resolvedPath, err)
	}
	return payload, nil
}
