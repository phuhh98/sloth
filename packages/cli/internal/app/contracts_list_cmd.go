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
			resolver := opts.Resolver()
			contracts, err := resolver.ListContracts(opts.PluginVersion)
			if err != nil {
				return err
			}

			format, err := output.ParseFormat(opts.Format)
			if err != nil {
				return err
			}

			if format == output.FormatJSON {
				payload := map[string]any{
					"pluginVersion": opts.PluginVersion,
					"source":        opts.Source,
					"contracts":     contracts,
				}
				return output.PrintJSON(cmd.OutOrStdout(), payload)
			}

			rows := make([][]string, 0, len(contracts))
			for _, contract := range contracts {
				rows = append(rows, []string{contract.Name, contract.Label, contract.Version, contract.SchemaVersion})
			}
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "pluginVersion=%s source=%s\n", opts.PluginVersion, opts.Source); err != nil {
				return err
			}
			return output.PrintTable(cmd.OutOrStdout(), []string{"NAME", "LABEL", "VERSION", "SCHEMA"}, rows)
		},
	}

	return cmd
}
