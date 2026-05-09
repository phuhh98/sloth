package app

import (
	"fmt"

	"github.com/phuhh98/sloth/packages/cli/pkg/config"
	"github.com/phuhh98/sloth/packages/cli/pkg/host"
	"github.com/phuhh98/sloth/packages/cli/pkg/output"
	"github.com/spf13/cobra"
)

func newContractsInspectCommand(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "inspect",
		Short: "Inspect host plugin and schema status",
		RunE: func(cmd *cobra.Command, args []string) error {
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
			schema, err := client.ContractSchema("", false)
			if err != nil {
				return err
			}

			format, err := output.ParseFormat(opts.Format)
			if err != nil {
				return err
			}
			if format == output.FormatJSON {
				return output.PrintJSON(cmd.OutOrStdout(), map[string]any{
					"profile": profileName,
					"status":  status,
					"schema":  schema,
				})
			}

			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "profile=%s host=%s\n", profileName, profile.Host); err != nil {
				return err
			}
			rows := [][]string{{status.PluginName, status.PluginVersion, fmt.Sprintf("%d", status.TotalComponents), schema.SchemaVersion, schema.SchemaURL}}
			return output.PrintTable(cmd.OutOrStdout(), []string{"PLUGIN", "VERSION", "COMPONENTS", "SCHEMA", "SCHEMA_URL"}, rows)
		},
	}
}
