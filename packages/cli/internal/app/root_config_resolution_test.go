package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/phuhh98/sloth/packages/cli/pkg/config"
)

func TestResolveConfigPrefersYAMLOverEnvAndDefault(t *testing.T) {
	tmp := t.TempDir()
	configPath := filepath.Join(tmp, ".sloth", "config.yaml")
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	content := "currentProfile: dev\nprofiles:\n  dev:\n    host: http://yaml-host:1337\n    authorizationToken: yaml-token\n    registry:\n      host: ghcr.io\n      repository: yaml-org/yaml-repo\n      useAuthorizationToken: false\n"
	if err := os.WriteFile(configPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	t.Setenv(config.EnvHost, "http://env-host:1337")
	t.Setenv(config.EnvAuthorizationToken, "env-token")
	t.Setenv(config.EnvRegistryHost, "env-registry.invalid")
	t.Setenv(config.EnvRegistryRepository, "env-org/env-repo")
	t.Setenv(config.EnvRegistryUseAuth, "true")

	opts := &Options{WorkingDir: tmp}
	_, profileName, profile, err := opts.ResolveConfig()
	if err != nil {
		t.Fatalf("resolve config: %v", err)
	}

	if profileName != "dev" {
		t.Fatalf("expected profile dev, got %q", profileName)
	}
	if profile.Host != "http://yaml-host:1337" {
		t.Fatalf("expected yaml host, got %q", profile.Host)
	}
	if token := config.EffectiveToken(profile, ""); token != "yaml-token" {
		t.Fatalf("expected yaml token, got %q", token)
	}
	if profile.Registry.Host != "ghcr.io" {
		t.Fatalf("expected yaml registry host, got %q", profile.Registry.Host)
	}
	if profile.Registry.Repository != "yaml-org/yaml-repo" {
		t.Fatalf("expected yaml registry repository, got %q", profile.Registry.Repository)
	}
	if profile.Registry.UseAuthorizationToken == nil || *profile.Registry.UseAuthorizationToken {
		t.Fatalf("expected yaml registry useAuthorizationToken false, got %v", profile.Registry.UseAuthorizationToken)
	}
}

func TestResolveConfigUsesEnvWhenYAMLMissing(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv(config.EnvProfile, "ci")
	t.Setenv(config.EnvHost, "http://env-host:1337")
	t.Setenv(config.EnvAuthorizationToken, "env-token")
	t.Setenv(config.EnvRegistryHost, "ghcr.io")
	t.Setenv(config.EnvRegistryRepository, "env-org/env-repo")
	t.Setenv(config.EnvRegistryUseAuth, "0")

	opts := &Options{WorkingDir: tmp}
	_, profileName, profile, err := opts.ResolveConfig()
	if err != nil {
		t.Fatalf("resolve config: %v", err)
	}

	if profileName != "ci" {
		t.Fatalf("expected profile ci from env, got %q", profileName)
	}
	if profile.Host != "http://env-host:1337" {
		t.Fatalf("expected env host, got %q", profile.Host)
	}
	if token := config.EffectiveToken(profile, ""); token != "env-token" {
		t.Fatalf("expected env token, got %q", token)
	}
	if profile.Registry.Host != "ghcr.io" {
		t.Fatalf("expected env registry host, got %q", profile.Registry.Host)
	}
	if profile.Registry.Repository != "env-org/env-repo" {
		t.Fatalf("expected env registry repository, got %q", profile.Registry.Repository)
	}
	if profile.Registry.UseAuthorizationToken == nil || *profile.Registry.UseAuthorizationToken {
		t.Fatalf("expected env registry useAuthorizationToken false, got %v", profile.Registry.UseAuthorizationToken)
	}
}

func TestResolveConfigFallsBackToDefaultHost(t *testing.T) {
	tmp := t.TempDir()
	opts := &Options{WorkingDir: tmp}
	_, profileName, profile, err := opts.ResolveConfig()
	if err != nil {
		t.Fatalf("resolve config: %v", err)
	}

	if profileName != config.DefaultProfileName {
		t.Fatalf("expected profile %q, got %q", config.DefaultProfileName, profileName)
	}
	if profile.Host != config.DefaultHost {
		t.Fatalf("expected default host %q, got %q", config.DefaultHost, profile.Host)
	}
	if profile.Registry.Host != config.DefaultRegistryHost {
		t.Fatalf("expected default registry host %q, got %q", config.DefaultRegistryHost, profile.Registry.Host)
	}
	if profile.Registry.Repository != config.DefaultRegistryRepository {
		t.Fatalf("expected default registry repository %q, got %q", config.DefaultRegistryRepository, profile.Registry.Repository)
	}
	if profile.Registry.UseAuthorizationToken == nil || !*profile.Registry.UseAuthorizationToken {
		t.Fatalf("expected default registry useAuthorizationToken true, got %v", profile.Registry.UseAuthorizationToken)
	}
}
