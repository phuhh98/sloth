package app

import (
	"fmt"
	"strings"

	"github.com/phuhh98/sloth/packages/cli/pkg/config"
	"github.com/phuhh98/sloth/packages/cli/pkg/host"
	"github.com/phuhh98/sloth/packages/cli/pkg/output"
	"github.com/phuhh98/sloth/packages/cli/pkg/verify"
	"github.com/phuhh98/sloth/packages/cli/pkg/workspace"
	"github.com/spf13/cobra"
)

func newContractsVerifyCommand(opts *Options) *cobra.Command {
	var filePath string
	var supportedRange string

	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify a contract file against compatibility and collision rules",
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(filePath) == "" {
				return fmt.Errorf("--file is required")
			}
			payload, err := workspace.LoadContractFile(filePath)
			if err != nil {
				return err
			}
			name, _ := payload["name"].(string)
			schemaVersion, _ := payload["schemaVersion"].(string)

			hostNames := []string{}
			hostSchema := ""
			if cfg, profileName, profile, err := opts.ResolveConfig(); err == nil {
				_ = cfg
				_ = profileName
				token := config.EffectiveToken(profile, opts.TokenOverride)
				client := host.NewClient(profile.Host, token)
				if status, err := client.PluginStatus(); err == nil {
					for _, component := range status.Components {
						if n, ok := component["name"].(string); ok {
							hostNames = append(hostNames, n)
						}
					}
					if len(status.CompatibleSchemaVersions) > 0 {
						hostSchema = status.CompatibleSchemaVersions[0]
					}
				}
			}

			resolver := opts.Resolver()
			officialCatalog, err := resolver.ListContracts(opts.PluginVersion)
			if err != nil {
				return err
			}
			officialNames := make([]string, 0, len(officialCatalog))
			for _, item := range officialCatalog {
				officialNames = append(officialNames, item.Name)
			}

			localNames, err := workspace.LocalContractNames(opts.WorkingDir)
			if err != nil {
				return err
			}

			result := verify.Run(verify.Input{
				ContractName:          name,
				ContractSchemaVersion: schemaVersion,
				HostSchemaVersion:     hostSchema,
				PluginVersion:         opts.PluginVersion,
				SupportedRange:        supportedRange,
				OfficialCatalogNames:  officialNames,
				HostNames:             hostNames,
				LocalNames:            localNames,
			})

			format, err := output.ParseFormat(opts.Format)
			if err != nil {
				return err
			}
			if format == output.FormatJSON {
				if err := output.PrintJSON(cmd.OutOrStdout(), result); err != nil {
					return err
				}
			} else {
				state := "valid"
				if !result.Valid {
					state = "invalid"
				}
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "verify=%s errors=%d warnings=%d\n", state, len(result.Errors), len(result.Warnings)); err != nil {
					return err
				}
				for _, item := range result.Errors {
					if _, err := fmt.Fprintf(cmd.OutOrStdout(), "ERROR: %s\n", item); err != nil {
						return err
					}
				}
				for _, item := range result.Warnings {
					if _, err := fmt.Fprintf(cmd.OutOrStdout(), "WARN: %s\n", item); err != nil {
						return err
					}
				}
			}

			if !result.Valid {
				return fmt.Errorf("verification failed")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Contract file path")
	cmd.Flags().StringVar(&supportedRange, "supported-range", ">=0.0.1", "Supported plugin semver range")
	return cmd
}
