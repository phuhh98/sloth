package app

import (
	"fmt"

	"github.com/phuhh98/sloth/packages/cli/pkg/config"
	"github.com/phuhh98/sloth/packages/cli/pkg/host"
	"github.com/phuhh98/sloth/packages/cli/pkg/source"
)

type Runtime struct {
	WorkingDir    string
	Format        string
	Source        string
	PluginVersion string
	ConfigPath    string
	Config        *config.File
	ProfileName   string
	Profile       config.Profile
	Resolver      source.Resolver
}

func (o *Options) BuildRuntime() (*Runtime, error) {
	cfg, profileName, profile, err := o.ResolveConfig()
	if err != nil {
		return nil, err
	}

	resolver, err := o.buildResolver(profile)
	if err != nil {
		return nil, fmt.Errorf("build contract resolver: %w", err)
	}

	return &Runtime{
		WorkingDir:    o.WorkingDir,
		Format:        o.Format,
		Source:        o.Source,
		PluginVersion: o.PluginVersion,
		ConfigPath:    o.ConfigPath,
		Config:        cfg,
		ProfileName:   profileName,
		Profile:       profile,
		Resolver:      resolver,
	}, nil
}

func (r *Runtime) HostClient(tokenOverride string) *host.Client {
	token := config.EffectiveToken(r.Profile, tokenOverride)
	return host.NewClient(r.Profile.Host, token)
}