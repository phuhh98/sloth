package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/phuhh98/sloth/packages/cli/pkg/lock"
	"github.com/phuhh98/sloth/packages/cli/pkg/output"
	"github.com/phuhh98/sloth/packages/cli/pkg/workspace"
	"github.com/spf13/cobra"
)

func newContractsAddCommand(opts *Options) *cobra.Command {
	var addAll bool
	var componentName string
	var setName string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add contracts to local workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !addAll {
				return fmt.Errorf("use one of: add --all | add component --name <name> | add set --name <set>")
			}
			return runAddAll(cmd, opts)
		},
	}

	componentCmd := &cobra.Command{
		Use:   "component",
		Short: "Add one component contract",
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(componentName) == "" {
				return fmt.Errorf("--name is required")
			}
			return runAddComponent(cmd, opts, componentName)
		},
	}
	componentCmd.Flags().StringVar(&componentName, "name", "", "Component contract name")

	setCmd := &cobra.Command{
		Use:   "set",
		Short: "Add one set of contracts",
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(setName) == "" {
				return fmt.Errorf("--name is required")
			}
			return runAddSet(cmd, opts, setName)
		},
	}
	setCmd.Flags().StringVar(&setName, "name", "", "Contract set name")

	cmd.Flags().BoolVar(&addAll, "all", false, "Add all compatible contracts")
	cmd.AddCommand(componentCmd)
	cmd.AddCommand(setCmd)

	return cmd
}

func runAddAll(cmd *cobra.Command, opts *Options) error {
	runtime, err := opts.BuildRuntime()
	if err != nil {
		return err
	}
	if err := workspace.Init(runtime.WorkingDir); err != nil {
		return err
	}
	contracts, err := runtime.Resolver.ListContracts(runtime.PluginVersion)
	if err != nil {
		return err
	}
	payloads := make([]map[string]any, 0, len(contracts))
	names := make([]string, 0, len(contracts))
	for _, c := range contracts {
		payloads = append(payloads, c.Payload)
		names = append(names, c.Name)
	}
	return addContracts(cmd, runtime, payloads, names, "all")
}

func runAddComponent(cmd *cobra.Command, opts *Options, componentName string) error {
	runtime, err := opts.BuildRuntime()
	if err != nil {
		return err
	}
	if err := workspace.Init(runtime.WorkingDir); err != nil {
		return err
	}
	match, err := runtime.Resolver.GetContract(runtime.PluginVersion, componentName)
	if err != nil {
		return err
	}
	return addContracts(cmd, runtime, []map[string]any{match.Payload}, []string{match.Name}, "component")
}

func runAddSet(cmd *cobra.Command, opts *Options, setName string) error {
	runtime, err := opts.BuildRuntime()
	if err != nil {
		return err
	}
	if err := workspace.Init(runtime.WorkingDir); err != nil {
		return err
	}
	names, err := loadSetContractNames(runtime.WorkingDir, setName)
	if err != nil {
		return err
	}
	if len(names) == 0 {
		return fmt.Errorf("set %q does not include any contracts", setName)
	}
	payloads := make([]map[string]any, 0, len(names))
	for _, name := range names {
		match, err := runtime.Resolver.GetContract(runtime.PluginVersion, name)
		if err != nil {
			return err
		}
		payloads = append(payloads, match.Payload)
	}
	return addContracts(cmd, runtime, payloads, names, "set")
}

func addContracts(cmd *cobra.Command, runtime *Runtime, payloads []map[string]any, names []string, mode string) error {
	lp := workspace.LockPath(runtime.WorkingDir)
	lf, err := lock.Read(lp)
	if err != nil {
		return err
	}

	written := make([]map[string]string, 0, len(payloads))
	for _, payload := range payloads {
		path, hash, schemaVersion, err := workspace.SaveContract(runtime.WorkingDir, payload)
		if err != nil {
			return err
		}
		name, _ := payload["name"].(string)
		version, _ := payload["version"].(string)
		lf.Upsert(lock.Entry{
			Name:          name,
			Version:       version,
			SchemaVersion: schemaVersion,
			Source:        runtime.Source,
			ContentHash:   hash,
			LastSyncedAt:  time.Now().UTC().Format(time.RFC3339),
		})
		written = append(written, map[string]string{"name": name, "version": version, "path": path})
	}

	if err := lock.Write(lp, lf); err != nil {
		return err
	}

	format, err := output.ParseFormat(runtime.Format)
	if err != nil {
		return err
	}
	if format == output.FormatJSON {
		return output.PrintJSON(cmd.OutOrStdout(), map[string]any{
			"mode":   mode,
			"count":  len(payloads),
			"names":  names,
			"saved":  written,
			"source": runtime.Source,
		})
	}

	rows := make([][]string, 0, len(written))
	for _, item := range written {
		rows = append(rows, []string{item["name"], item["version"], item["path"]})
	}
	return output.PrintTable(cmd.OutOrStdout(), []string{"NAME", "VERSION", "PATH"}, rows)
}

func loadSetContractNames(workingDir string, setName string) ([]string, error) {
	type setFile struct {
		Components []string `json:"components"`
	}

	setPath := filepath.Join(workspace.SetsDir(workingDir), setName+".json")
	raw, err := os.ReadFile(setPath)
	if err != nil {
		return nil, fmt.Errorf("read set file %q: %w", setPath, err)
	}

	parsed := &setFile{}
	if err := json.Unmarshal(raw, parsed); err != nil {
		return nil, fmt.Errorf("parse set file %q: %w", setPath, err)
	}

	return parsed.Components, nil
}
