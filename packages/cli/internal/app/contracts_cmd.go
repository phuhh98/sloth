package app

import "github.com/spf13/cobra"

func newContractsCommand(opts *Options) *cobra.Command {
	contractsCmd := &cobra.Command{
		Use:   "contracts",
		Short: "Manage local and remote component contracts",
	}

	contractsCmd.PersistentFlags().StringVar(&opts.PluginVersion, "version", "latest", "Target contract release version (x.y.z|latest)")
	contractsCmd.PersistentFlags().StringVar(&opts.PluginVersion, "plugin-version", "latest", "Deprecated alias for --version")
	contractsCmd.PersistentFlags().StringVar(&opts.Source, "source", "local", "Contract source: local|oci")

	contractsCmd.AddCommand(newContractsListCommand(opts))
	contractsCmd.AddCommand(newContractsPullCommand(opts))
	contractsCmd.AddCommand(newContractsInspectCommand(opts))
	contractsCmd.AddCommand(newContractsAddCommand(opts))
	contractsCmd.AddCommand(newContractsVerifyCommand(opts))
	contractsCmd.AddCommand(newContractsPushCommand(opts))

	return contractsCmd
}
