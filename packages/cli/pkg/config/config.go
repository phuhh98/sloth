package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

const DefaultConfigRelativePath = ".sloth/config.yaml"
const DefaultHost = "http://localhost:1337"
const DefaultProfileName = "default"
const DefaultRegistryHost = "ghcr.io"
const DefaultRegistryRepository = "phuhh98/sloth/contracts"
const DefaultRegistryUseAuthorizationToken = true

const (
	EnvConfigPath          = "SLOTH_CONFIG"
	EnvProfile             = "SLOTH_PROFILE"
	EnvHost                = "SLOTH_HOST"
	EnvAuthorizationToken  = "SLOTH_AUTHORIZATION_TOKEN"
	EnvAuthorizationToken2 = "SLOTH_TOKEN"
	EnvRegistryHost        = "SLOTH_REGISTRY_HOST"
	EnvRegistryRepository  = "SLOTH_REGISTRY_REPOSITORY"
	EnvRegistryUseAuth     = "SLOTH_REGISTRY_USE_AUTHORIZATION_TOKEN"
)

var ErrConfigNotFound = errors.New("sloth config file not found")

type Registry struct {
	Host                  string `yaml:"host,omitempty" json:"host,omitempty"`
	Repository            string `yaml:"repository,omitempty" json:"repository,omitempty"`
	UseAuthorizationToken *bool  `yaml:"useAuthorizationToken,omitempty" json:"useAuthorizationToken,omitempty"`
}

type Profile struct {
	Host               string   `yaml:"host" json:"host"`
	Token              string   `yaml:"token,omitempty" json:"token,omitempty"`
	AuthorizationToken string   `yaml:"authorizationToken,omitempty" json:"authorizationToken,omitempty"`
	Registry           Registry `yaml:"registry,omitempty" json:"registry,omitempty"`
}

type File struct {
	CurrentProfile string             `yaml:"currentProfile" json:"currentProfile"`
	Profiles       map[string]Profile `yaml:"profiles" json:"profiles"`
}

type EnvSettings struct {
	ConfigPath         string
	Profile            string
	Host               string
	AuthorizationToken string
	RegistryHost       string
	RegistryRepository string
	RegistryUseAuth    *bool
}

func EnvSettingsFromOS() EnvSettings {
	settings := EnvSettings{
		ConfigPath:         strings.TrimSpace(os.Getenv(EnvConfigPath)),
		Profile:            strings.TrimSpace(os.Getenv(EnvProfile)),
		Host:               strings.TrimSpace(os.Getenv(EnvHost)),
		RegistryHost:       strings.TrimSpace(os.Getenv(EnvRegistryHost)),
		RegistryRepository: strings.TrimSpace(os.Getenv(EnvRegistryRepository)),
	}
	settings.AuthorizationToken = strings.TrimSpace(os.Getenv(EnvAuthorizationToken))
	if settings.AuthorizationToken == "" {
		settings.AuthorizationToken = strings.TrimSpace(os.Getenv(EnvAuthorizationToken2))
	}
	settings.RegistryUseAuth = parseOptionalBool(os.Getenv(EnvRegistryUseAuth))
	return settings
}

func parseOptionalBool(raw string) *bool {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "1", "true", "yes", "on":
		v := true
		return &v
	case "0", "false", "no", "off":
		v := false
		return &v
	default:
		return nil
	}
}

func ResolvePath(workingDir string, configuredPath string, envPath string) string {
	if configuredPath != "" {
		if filepath.IsAbs(configuredPath) {
			return configuredPath
		}
		return filepath.Join(workingDir, configuredPath)
	}
	if envPath != "" {
		if filepath.IsAbs(envPath) {
			return envPath
		}
		return filepath.Join(workingDir, envPath)
	}
	return filepath.Join(workingDir, DefaultConfigRelativePath)
}

func Load(path string) (*File, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrConfigNotFound
		}
		return nil, fmt.Errorf("read config: %w", err)
	}

	cfg := &File{}
	if err := yaml.Unmarshal(raw, cfg); err != nil {
		return nil, fmt.Errorf("parse config yaml: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (f *File) Validate() error {
	if len(f.Profiles) == 0 {
		return errors.New("config has no profiles")
	}

	if f.CurrentProfile != "" {
		if _, ok := f.Profiles[f.CurrentProfile]; !ok {
			return fmt.Errorf("currentProfile %q does not exist", f.CurrentProfile)
		}
	}

	for name, profile := range f.Profiles {
		if strings.TrimSpace(profile.Host) == "" {
			return fmt.Errorf("profile %q: host is required", name)
		}
	}

	return nil
}

func (f *File) ResolveProfile(selected string) (string, Profile, error) {
	if selected != "" {
		p, ok := f.Profiles[selected]
		if !ok {
			return "", Profile{}, fmt.Errorf("profile %q not found", selected)
		}
		return selected, p, nil
	}

	if f.CurrentProfile != "" {
		return f.CurrentProfile, f.Profiles[f.CurrentProfile], nil
	}

	if p, ok := f.Profiles[DefaultProfileName]; ok {
		return DefaultProfileName, p, nil
	}

	names := make([]string, 0, len(f.Profiles))
	for name := range f.Profiles {
		names = append(names, name)
	}
	sort.Strings(names)
	if len(names) > 0 {
		selectedName := names[0]
		return selectedName, f.Profiles[selectedName], nil
	}

	return "", Profile{}, errors.New("no usable profile found")
}

func EffectiveToken(profile Profile, override string) string {
	if strings.TrimSpace(override) != "" {
		return strings.TrimSpace(override)
	}
	if strings.TrimSpace(profile.AuthorizationToken) != "" {
		return strings.TrimSpace(profile.AuthorizationToken)
	}
	return strings.TrimSpace(profile.Token)
}
