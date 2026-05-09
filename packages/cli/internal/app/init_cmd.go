package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/phuhh98/sloth/packages/cli/pkg/config"
	"github.com/phuhh98/sloth/packages/cli/pkg/workspace"
	"github.com/spf13/cobra"
)

func newInitCommand(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize local .sloth workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := workspace.Init(opts.WorkingDir); err != nil {
				return err
			}
			env := config.EnvSettingsFromOS()
			cfgPath := config.ResolvePath(opts.WorkingDir, opts.ConfigPath, env.ConfigPath)
			if _, err := os.Stat(cfgPath); err != nil {
				if !errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf("check config file: %w", err)
				}
				if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
					return fmt.Errorf("create config directory: %w", err)
				}
				template := "currentProfile: default\nprofiles:\n  default:\n    host: http://localhost:1337\n    authorizationToken: \"\"\n"
				if err := os.WriteFile(cfgPath, []byte(template), 0o644); err != nil {
					return fmt.Errorf("write default config file: %w", err)
				}
			}

			_, err := fmt.Fprintf(cmd.OutOrStdout(), "Initialized workspace at %s\n", workspace.RootPath(opts.WorkingDir))
			return err
		},
	}
}
