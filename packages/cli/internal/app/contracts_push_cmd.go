package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/phuhh98/sloth/packages/cli/pkg/config"
	"github.com/phuhh98/sloth/packages/cli/pkg/host"
	"github.com/phuhh98/sloth/packages/cli/pkg/lock"
	"github.com/phuhh98/sloth/packages/cli/pkg/output"
	"github.com/phuhh98/sloth/packages/cli/pkg/workspace"
	"github.com/spf13/cobra"
)

func newContractsPushCommand(opts *Options) *cobra.Command {
	var dryRun bool
	var retries int
	var yesToAll bool

	cmd := &cobra.Command{
		Use:   "push",
		Short: "Push local contracts to host ingest endpoint",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := workspace.Init(opts.WorkingDir); err != nil {
				return err
			}

			_, profileName, profile, err := opts.ResolveConfig()
			if err != nil {
				return err
			}
			token := config.EffectiveToken(profile, opts.TokenOverride)
			client := host.NewClient(profile.Host, token)
			status, err := client.PluginStatus()
			if err != nil {
				return err
			}

			localContracts, names, err := loadLocalContracts(opts.WorkingDir)
			if err != nil {
				return err
			}
			if len(localContracts) == 0 {
				return fmt.Errorf("no local contracts found in %s", workspace.ContractsDir(opts.WorkingDir))
			}

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

			if dryRun {
				return output.PrintJSON(cmd.OutOrStdout(), map[string]any{
					"profile":             profileName,
					"host":                profile.Host,
					"contracts":           names,
					"hostTotalComponents": status.TotalComponents,
					"missingOnHost":       missing,
					"dryRun":              true,
				})
			}

			if !yesToAll {
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Pushing %d contracts to %s\n", len(localContracts), profile.Host); err != nil {
					return err
				}
			}

			var ingest *host.IngestResponse
			var pushErr error
			for attempt := 0; attempt <= retries; attempt++ {
				ingest, pushErr = client.IngestContracts(localContracts)
				if pushErr == nil {
					break
				}
				time.Sleep(200 * time.Millisecond)
			}
			if pushErr != nil {
				return pushErr
			}

			lp := workspace.LockPath(opts.WorkingDir)
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
			if err := lock.Write(lp, lf); err != nil {
				return err
			}

			return output.PrintJSON(cmd.OutOrStdout(), map[string]any{
				"profile": profileName,
				"host":    profile.Host,
				"result":  ingest,
			})
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show push plan without writing to host")
	cmd.Flags().IntVar(&retries, "retries", 1, "Retry count for ingest call")
	cmd.Flags().BoolVarP(&yesToAll, "yes-to-all", "Y", false, "Skip confirmation prompts")
	return cmd
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
