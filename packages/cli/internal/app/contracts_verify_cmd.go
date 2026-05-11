package app

import (
	"fmt"
	"io"
	"strings"

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
			runtime, err := opts.BuildRuntime()
			if err != nil {
				return err
			}
			input, err := buildVerifyInput(runtime, filePath, supportedRange)
			if err != nil {
				return err
			}
			result := verify.Run(input)
			format, err := output.ParseFormat(runtime.Format)
			if err != nil {
				return err
			}
			if err := printVerifyResult(cmd.OutOrStdout(), result, format); err != nil {
				return err
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

func buildVerifyInput(runtime *Runtime, filePath, supportedRange string) (verify.Input, error) {
	payload, err := workspace.LoadContractFile(filePath)
	if err != nil {
		return verify.Input{}, err
	}
	name, _ := payload["name"].(string)
	schemaVersion, _ := payload["schemaVersion"].(string)

	hostNames := []string{}
	hostSchema := ""
	if status, err := runtime.HostClient("").PluginStatus(); err == nil {
		for _, component := range status.Components {
			if n, ok := component["name"].(string); ok {
				hostNames = append(hostNames, n)
			}
		}
		if len(status.CompatibleSchemaVersions) > 0 {
			hostSchema = status.CompatibleSchemaVersions[0]
		}
	}

	officialCatalog, err := runtime.Resolver.ListContracts(runtime.PluginVersion)
	if err != nil {
		return verify.Input{}, err
	}
	officialNames := make([]string, 0, len(officialCatalog))
	for _, item := range officialCatalog {
		officialNames = append(officialNames, item.Name)
	}

	localNames, err := workspace.LocalContractNames(runtime.WorkingDir)
	if err != nil {
		return verify.Input{}, err
	}

	return verify.Input{
		ContractName:          name,
		ContractSchemaVersion: schemaVersion,
		HostSchemaVersion:     hostSchema,
		PluginVersion:         runtime.PluginVersion,
		SupportedRange:        supportedRange,
		OfficialCatalogNames:  officialNames,
		HostNames:             hostNames,
		LocalNames:            localNames,
	}, nil
}

func printVerifyResult(w io.Writer, result verify.Result, format output.Format) error {
	if format == output.FormatJSON {
		return output.PrintJSON(w, result)
	}
	state := "valid"
	if !result.Valid {
		state = "invalid"
	}
	if _, err := fmt.Fprintf(w, "verify=%s errors=%d warnings=%d\n", state, len(result.Errors), len(result.Warnings)); err != nil {
		return err
	}
	for _, item := range result.Errors {
		if _, err := fmt.Fprintf(w, "ERROR: %s\n", item); err != nil {
			return err
		}
	}
	for _, item := range result.Warnings {
		if _, err := fmt.Fprintf(w, "WARN: %s\n", item); err != nil {
			return err
		}
	}
	return nil
}

