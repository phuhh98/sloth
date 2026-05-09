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
	content := "currentProfile: dev\nprofiles:\n  dev:\n    host: http://yaml-host:1337\n    authorizationToken: yaml-token\n"
	if err := os.WriteFile(configPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	t.Setenv(config.EnvHost, "http://env-host:1337")
	t.Setenv(config.EnvAuthorizationToken, "env-token")

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
}

func TestResolveConfigUsesEnvWhenYAMLMissing(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv(config.EnvProfile, "ci")
	t.Setenv(config.EnvHost, "http://env-host:1337")
	t.Setenv(config.EnvAuthorizationToken, "env-token")

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
}
