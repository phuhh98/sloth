package app

import (
	"fmt"

	"github.com/phuhh98/sloth/packages/cli/pkg/output"
	"github.com/spf13/cobra"
)

func newContractsListCommand(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Aliases: []string{"ls"},
		Short: "List available contracts for a plugin version",
		RunE: func(cmd *cobra.Command, args []string) error {
			runtime, err := opts.BuildRuntime()
			if err != nil {
				return err
			}

			contracts, err := runtime.Resolver.ListContracts(runtime.PluginVersion)
			if err != nil {
				return err
			}

			format, err := output.ParseFormat(runtime.Format)
			if err != nil {
				return err
			}

			if format == output.FormatJSON {
				payload := map[string]any{
					"pluginVersion": runtime.PluginVersion,
					"source":        runtime.Source,
					"contracts":     contracts,
				}
				return output.PrintJSON(cmd.OutOrStdout(), payload)
			}

			rows := make([][]string, 0, len(contracts))
			for _, contract := range contracts {
				rows = append(rows, []string{contract.Name, contract.Label, contract.Version, contract.SchemaVersion})
			}
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "pluginVersion=%s source=%s\n", runtime.PluginVersion, runtime.Source); err != nil {
				return err
			}
			return output.PrintTable(cmd.OutOrStdout(), []string{"NAME", "LABEL", "VERSION", "SCHEMA"}, rows)
		},
	}

	return cmd
}
