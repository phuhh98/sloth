package app

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/phuhh98/sloth/packages/cli/pkg/lock"
	"github.com/phuhh98/sloth/packages/cli/pkg/output"
	"github.com/phuhh98/sloth/packages/cli/pkg/source"
	"github.com/phuhh98/sloth/packages/cli/pkg/workspace"
	"github.com/spf13/cobra"
)

type pulledContractMeta struct {
	Name          string
	Version       string
	SchemaVersion string
	Payload       map[string]any
}

type pullWriteResult struct {
	Path        string
	ContentHash string
}

func newContractsPullCommand(opts *Options) *cobra.Command {
	var contractName string
	var outPath string

	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull one contract by name from the selected release version",
		RunE: func(cmd *cobra.Command, args []string) error {
			runtime, err := opts.BuildRuntime()
			if err != nil {
				return err
			}

			meta, err := resolvePullContract(runtime, contractName)
			if err != nil {
				return err
			}

			writeResult, err := writePulledContract(runtime, meta, outPath)
			if err != nil {
				return err
			}

			return printPullOutput(cmd, runtime, meta, writeResult)
		},
	}

	cmd.Flags().StringVar(&contractName, "name", "", "Contract name to pull")
	cmd.Flags().StringVar(&outPath, "out", "", "Output path for the pulled contract file")
	return cmd
}

func resolvePullContract(runtime *Runtime, contractName string) (*pulledContractMeta, error) {
	if strings.TrimSpace(contractName) == "" {
		return nil, fmt.Errorf("--name is required")
	}

	contract, err := runtime.Resolver.GetContract(runtime.PluginVersion, contractName)
	if err != nil {
		return nil, err
	}
	return newPulledContractMeta(contract, runtime.PluginVersion, contractName)
}

func newPulledContractMeta(contract *source.Contract, requestedVersion string, requestedName string) (*pulledContractMeta, error) {
	if contract == nil || len(contract.Payload) == 0 {
		return nil, fmt.Errorf("contract %q has empty payload", requestedName)
	}

	name := strings.TrimSpace(payloadString(contract.Payload, "name"))
	if name == "" {
		name = strings.TrimSpace(contract.Name)
	}

	version := strings.TrimSpace(payloadString(contract.Payload, "version"))
	if version == "" {
		version = strings.TrimSpace(contract.Version)
	}
	if version == "" {
		version = strings.TrimSpace(requestedVersion)
	}

	schemaVersion := strings.TrimSpace(payloadString(contract.Payload, "schemaVersion"))

	return &pulledContractMeta{
		Name:          name,
		Version:       version,
		SchemaVersion: schemaVersion,
		Payload:       contract.Payload,
	}, nil
}

func writePulledContract(runtime *Runtime, meta *pulledContractMeta, outPath string) (*pullWriteResult, error) {
	if strings.TrimSpace(outPath) == "" {
		return writePulledContractToWorkspace(runtime, meta)
	}

	resolvedOutPath, err := resolveOutputPath(outPath, meta.Name, meta.Version)
	if err != nil {
		return nil, err
	}
	writtenPath, hash, err := writeContractPayload(resolvedOutPath, meta.Payload)
	if err != nil {
		return nil, err
	}

	return &pullWriteResult{Path: writtenPath, ContentHash: hash}, nil
}

func writePulledContractToWorkspace(runtime *Runtime, meta *pulledContractMeta) (*pullWriteResult, error) {
	if err := workspace.Init(runtime.WorkingDir); err != nil {
		return nil, err
	}

	writtenPath, hash, schemaVersion, err := workspace.SaveContract(runtime.WorkingDir, meta.Payload)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(meta.SchemaVersion) == "" {
		meta.SchemaVersion = schemaVersion
	}

	if err := upsertLockForPull(runtime, meta, hash); err != nil {
		return nil, err
	}

	return &pullWriteResult{Path: writtenPath, ContentHash: hash}, nil
}

func upsertLockForPull(runtime *Runtime, meta *pulledContractMeta, contentHash string) error {
	lp := workspace.LockPath(runtime.WorkingDir)
	lf, err := lock.Read(lp)
	if err != nil {
		return err
	}

	lf.Upsert(lock.Entry{
		Name:          meta.Name,
		Version:       meta.Version,
		SchemaVersion: meta.SchemaVersion,
		Source:        runtime.Source,
		ContentHash:   contentHash,
		LastSyncedAt:  time.Now().UTC().Format(time.RFC3339),
	})
	return lock.Write(lp, lf)
}

func printPullOutput(cmd *cobra.Command, runtime *Runtime, meta *pulledContractMeta, writeResult *pullWriteResult) error {
	format, err := output.ParseFormat(runtime.Format)
	if err != nil {
		return err
	}

	if format == output.FormatJSON {
		return output.PrintJSON(cmd.OutOrStdout(), map[string]any{
			"name":          meta.Name,
			"version":       meta.Version,
			"schemaVersion": meta.SchemaVersion,
			"source":        runtime.Source,
			"path":          writeResult.Path,
			"contentHash":   writeResult.ContentHash,
		})
	}

	rows := [][]string{{meta.Name, meta.Version, meta.SchemaVersion, writeResult.Path}}
	return output.PrintTable(cmd.OutOrStdout(), []string{"NAME", "VERSION", "SCHEMA", "PATH"}, rows)
}

func payloadString(payload map[string]any, key string) string {
	v, _ := payload[key].(string)
	return v
}

func resolveOutputPath(raw string, contractName string, contractVersion string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", fmt.Errorf("output path is required")
	}

	info, err := os.Stat(trimmed)
	if err == nil {
		if info.IsDir() {
			return filepath.Join(trimmed, fmt.Sprintf("%s@%s.json", contractName, contractVersion)), nil
		}
		return trimmed, nil
	}
	if !os.IsNotExist(err) {
		return "", fmt.Errorf("stat output path %q: %w", trimmed, err)
	}

	if strings.HasSuffix(trimmed, string(os.PathSeparator)) {
		return filepath.Join(trimmed, fmt.Sprintf("%s@%s.json", contractName, contractVersion)), nil
	}

	return trimmed, nil
}

func writeContractPayload(path string, payload map[string]any) (string, string, error) {
	raw, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return "", "", fmt.Errorf("marshal contract payload: %w", err)
	}
	raw = append(raw, '\n')

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", "", fmt.Errorf("create output directory: %w", err)
	}
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		return "", "", fmt.Errorf("write contract file: %w", err)
	}

	h := sha256.Sum256(raw)
	return path, hex.EncodeToString(h[:]), nil
}
