package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/phuhh98/sloth/packages/cli/pkg/config"
	"github.com/phuhh98/sloth/packages/cli/pkg/registry"
	"github.com/phuhh98/sloth/packages/cli/pkg/source"
	"github.com/spf13/cobra"
)

type Options struct {
	WorkingDir    string
	ConfigPath    string
	Profile       string
	HostOverride  string
	TokenOverride string
	Format        string
	Source        string
	PluginVersion string
}

func NewRootCommand() *cobra.Command {
	wd, _ := os.Getwd()
	opts := &Options{
		WorkingDir: wd,
		Format:     "table",
		Source:     "local",
	}

	rootCmd := &cobra.Command{
		Use:   "sloth",
		Short: "sloth CLI for component contract management",
	}

	rootCmd.PersistentFlags().StringVar(&opts.ConfigPath, "config", "", "Path to .sloth/config.yaml")
	rootCmd.PersistentFlags().StringVar(&opts.Profile, "profile", "", "Profile name from config.yaml")
	rootCmd.PersistentFlags().StringVarP(&opts.HostOverride, "host", "H", "", "Host URL override")
	rootCmd.PersistentFlags().StringVarP(&opts.TokenOverride, "authorization-token", "T", "", "Authorization token override")
	rootCmd.PersistentFlags().StringVar(&opts.Format, "format", "table", "Output format: table|json")

	rootCmd.AddCommand(newInitCommand(opts))
	rootCmd.AddCommand(newContractsCommand(opts))

	return rootCmd
}

func (o *Options) ResolveConfig() (*config.File, string, config.Profile, error) {
	env := config.EnvSettingsFromOS()
	path := config.ResolvePath(o.WorkingDir, o.ConfigPath, env.ConfigPath)
	cfg, err := config.Load(path)
	if err != nil {
		if !errors.Is(err, config.ErrConfigNotFound) {
			return nil, "", config.Profile{}, err
		}
		cfg = nil
	}

	selectedProfile := strings.TrimSpace(o.Profile)
	if selectedProfile == "" {
		selectedProfile = env.Profile
	}

	name := selectedProfile
	profile := config.Profile{}
	if cfg != nil {
		name, profile, err = cfg.ResolveProfile(selectedProfile)
		if err != nil {
			return nil, "", config.Profile{}, err
		}
	}

	if strings.TrimSpace(profile.Host) == "" {
		profile.Host = env.Host
	}
	if strings.TrimSpace(profile.AuthorizationToken) == "" {
		profile.AuthorizationToken = env.AuthorizationToken
	}
	if strings.TrimSpace(profile.Registry.Host) == "" {
		profile.Registry.Host = env.RegistryHost
	}
	if strings.TrimSpace(profile.Registry.Repository) == "" {
		profile.Registry.Repository = env.RegistryRepository
	}
	if profile.Registry.UseAuthorizationToken == nil {
		profile.Registry.UseAuthorizationToken = env.RegistryUseAuth
	}
	if strings.TrimSpace(profile.Host) == "" {
		profile.Host = config.DefaultHost
	}
	if strings.TrimSpace(profile.Registry.Host) == "" {
		profile.Registry.Host = config.DefaultRegistryHost
	}
	if strings.TrimSpace(profile.Registry.Repository) == "" {
		profile.Registry.Repository = config.DefaultRegistryRepository
	}
	if profile.Registry.UseAuthorizationToken == nil {
		v := config.DefaultRegistryUseAuthorizationToken
		profile.Registry.UseAuthorizationToken = &v
	}
	if strings.TrimSpace(name) == "" {
		name = config.DefaultProfileName
	}

	if o.HostOverride != "" {
		profile.Host = o.HostOverride
	}
	if o.TokenOverride != "" {
		profile.AuthorizationToken = o.TokenOverride
	}
	if profile.Host == "" {
		return nil, "", config.Profile{}, fmt.Errorf("resolved profile %q has empty host", name)
	}
	return cfg, name, profile, nil
}

func (o *Options) Resolver() source.Resolver {
	if strings.TrimSpace(o.Source) == "oci" {
		_, _, profile, err := o.ResolveConfig()
		if err != nil {
			return source.NewErrorResolver(err)
		}

		client, err := registry.NewOCIClient(registry.OCIClientOptions{
			Host:                  profile.Registry.Host,
			Repository:            profile.Registry.Repository,
			AuthorizationToken:    config.EffectiveToken(profile, o.TokenOverride),
			UseAuthorizationToken: profile.Registry.UseAuthorizationToken != nil && *profile.Registry.UseAuthorizationToken,
		})
		if err != nil {
			return source.NewErrorResolver(err)
		}

		return source.NewOCIRegistryResolver(client)
	}

	base := filepath.Join(o.WorkingDir, "apps", "docs", "static", "registry", "contracts")
	return source.NewLocalRegistryResolver(base)
}
