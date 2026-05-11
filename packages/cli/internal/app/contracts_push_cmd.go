package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/phuhh98/sloth/packages/cli/pkg/host"
	"github.com/phuhh98/sloth/packages/cli/pkg/lock"
	"github.com/phuhh98/sloth/packages/cli/pkg/output"
	"github.com/phuhh98/sloth/packages/cli/pkg/workspace"
	"github.com/spf13/cobra"
)

type pushFlags struct {
	dryRun   bool
	retries  int
	yesToAll bool
}

func newContractsPushCommand(opts *Options) *cobra.Command {
	var flags pushFlags

	cmd := &cobra.Command{
		Use:   "push",
		Short: "Push local contracts to host ingest endpoint",
		RunE: func(cmd *cobra.Command, args []string) error {
			return executePush(cmd, opts, flags)
		},
	}

	cmd.Flags().BoolVar(&flags.dryRun, "dry-run", false, "Show push plan without writing to host")
	cmd.Flags().IntVar(&flags.retries, "retries", 1, "Retry count for ingest call")
	cmd.Flags().BoolVarP(&flags.yesToAll, "yes-to-all", "Y", false, "Skip confirmation prompts")
	return cmd
}

func executePush(cmd *cobra.Command, opts *Options, flags pushFlags) error {
	runtime, err := opts.BuildRuntime()
	if err != nil {
		return err
	}
	if err := workspace.Init(runtime.WorkingDir); err != nil {
		return err
	}

	client := runtime.HostClient("")
	status, err := client.PluginStatus()
	if err != nil {
		return err
	}

	localContracts, names, err := loadLocalContracts(runtime.WorkingDir)
	if err != nil {
		return err
	}
	if len(localContracts) == 0 {
		return fmt.Errorf("no local contracts found in %s", workspace.ContractsDir(runtime.WorkingDir))
	}

	if flags.dryRun {
		return printPushDryRun(cmd, runtime, names, status)
	}

	if !flags.yesToAll {
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Pushing %d contracts to %s\n", len(localContracts), runtime.Profile.Host); err != nil {
			return err
		}
	}

	ingest, err := ingestWithRetry(client, localContracts, flags.retries)
	if err != nil {
		return err
	}
	if err := syncLockAfterPush(runtime.WorkingDir); err != nil {
		return err
	}
	return output.PrintJSON(cmd.OutOrStdout(), map[string]any{
		"profile": runtime.ProfileName,
		"host":    runtime.Profile.Host,
		"result":  ingest,
	})
}

func printPushDryRun(cmd *cobra.Command, runtime *Runtime, names []string, status *host.PluginStatus) error {
	hostNames := map[string]bool{}
	for _, comp := range status.Components {
		if n, ok := comp["name"].(string); ok {
			hostNames[strings.ToLower(n)] = true
		}
	}
	missing := 0
	for _, name := range names {
		if !hostNames[strings.ToLower(name)] {
			missing++
		}
	}
	return output.PrintJSON(cmd.OutOrStdout(), map[string]any{
		"profile":             runtime.ProfileName,
		"host":                runtime.Profile.Host,
		"contracts":           names,
		"hostTotalComponents": status.TotalComponents,
		"missingOnHost":       missing,
		"dryRun":              true,
	})
}

func ingestWithRetry(client *host.Client, contracts []map[string]any, retries int) (*host.IngestResponse, error) {
	var result *host.IngestResponse
	var err error
	for attempt := 0; attempt <= retries; attempt++ {
		result, err = client.IngestContracts(contracts)
		if err == nil {
			return result, nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return nil, err
}

func syncLockAfterPush(workingDir string) error {
	lp := workspace.LockPath(workingDir)
	lf, err := lock.Read(lp)
	if err != nil {
		return err
	}
	now := time.Now().UTC().Format(time.RFC3339)
	for _, key := range lf.SortedKeys() {
		entry := lf.Entries[key]
		entry.LastSyncedAt = now
		lf.Upsert(entry)
	}
	return lock.Write(lp, lf)
}

func loadLocalContracts(workingDir string) ([]map[string]any, []string, error) {
	dir := workspace.ContractsDir(workingDir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, err
	}
	contracts := []map[string]any{}
	names := []string{}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		payload, err := workspace.LoadContractFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, nil, err
		}
		contracts = append(contracts, payload)
		if name, ok := payload["name"].(string); ok {
			names = append(names, name)
		}
	}
	return contracts, names, nil
}
