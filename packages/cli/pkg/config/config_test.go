package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAndResolveProfile(t *testing.T) {
	t.Parallel()
	tmp := t.TempDir()
	configPath := filepath.Join(tmp, ".sloth", "config.yaml")
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	content := "currentProfile: dev\nprofiles:\n  dev:\n    host: http://localhost:1337\n    authorizationToken: abc\n"
	if err := os.WriteFile(configPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	name, profile, err := cfg.ResolveProfile("")
	if err != nil {
		t.Fatalf("resolve profile: %v", err)
	}
	if name != "dev" {
		t.Fatalf("expected profile dev, got %s", name)
	}
	if profile.Host != "http://localhost:1337" {
		t.Fatalf("unexpected host %s", profile.Host)
	}
	if token := EffectiveToken(profile, ""); token != "abc" {
		t.Fatalf("unexpected token %s", token)
	}
}
