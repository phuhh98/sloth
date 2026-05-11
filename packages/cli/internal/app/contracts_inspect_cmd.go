package app

import (
	"fmt"

	"github.com/phuhh98/sloth/packages/cli/pkg/output"
	"github.com/spf13/cobra"
)

func newContractsInspectCommand(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "inspect",
		Short: "Inspect host plugin and schema status",
		RunE: func(cmd *cobra.Command, args []string) error {
			runtime, err := opts.BuildRuntime()
			if err != nil {
				return err
			}
			client := runtime.HostClient("")

			status, err := client.PluginStatus()
			if err != nil {
				return err
			}
			schema, err := client.ContractSchema("", false)
			if err != nil {
				return err
			}

			format, err := output.ParseFormat(runtime.Format)
			if err != nil {
				return err
			}
			if format == output.FormatJSON {
				return output.PrintJSON(cmd.OutOrStdout(), map[string]any{
					"profile": runtime.ProfileName,
					"status":  status,
					"schema":  schema,
				})
			}

			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "profile=%s host=%s\n", runtime.ProfileName, runtime.Profile.Host); err != nil {
				return err
			}
			rows := [][]string{{status.PluginName, status.PluginVersion, fmt.Sprintf("%d", status.TotalComponents), schema.SchemaVersion, schema.SchemaURL}}
			return output.PrintTable(cmd.OutOrStdout(), []string{"PLUGIN", "VERSION", "COMPONENTS", "SCHEMA", "SCHEMA_URL"}, rows)
		},
	}
}
