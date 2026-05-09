package app

import "github.com/spf13/cobra"

func newContractsCommand(opts *Options) *cobra.Command {
	contractsCmd := &cobra.Command{
		Use:   "contracts",
		Short: "Manage local and remote component contracts",
	}

	contractsCmd.PersistentFlags().StringVar(&opts.PluginVersion, "plugin-version", "", "Target plugin version")
	contractsCmd.PersistentFlags().StringVar(&opts.Source, "source", "local", "Contract source: local")

	contractsCmd.AddCommand(newContractsListCommand(opts))
	contractsCmd.AddCommand(newContractsInspectCommand(opts))
	contractsCmd.AddCommand(newContractsAddCommand(opts))
	contractsCmd.AddCommand(newContractsVerifyCommand(opts))
	contractsCmd.AddCommand(newContractsPushCommand(opts))

	return contractsCmd
}
