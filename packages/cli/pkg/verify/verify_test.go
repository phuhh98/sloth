package verify

import "testing"

func TestRunBlocksNameCollisionsAndCompatibility(t *testing.T) {
	t.Parallel()
	result := Run(Input{
		ContractName:          "hero-banner",
		ContractSchemaVersion: "0.0.2",
		HostSchemaVersion:     "0.0.1",
		PluginVersion:         "1.0.0",
		SupportedRange:        ">=0.1.0, <1.0.0",
		OfficialCatalogNames:  []string{"hero-banner"},
		HostNames:             []string{"feature-grid"},
	})

	if result.Valid {
		t.Fatalf("expected invalid result")
	}
	if len(result.Errors) < 2 {
		t.Fatalf("expected multiple errors, got %v", result.Errors)
	}
}

func TestRunValid(t *testing.T) {
	t.Parallel()
	result := Run(Input{
		ContractName:          "new-card",
		ContractSchemaVersion: "0.0.1",
		HostSchemaVersion:     "0.0.1",
		PluginVersion:         "0.2.0",
		SupportedRange:        ">=0.1.0, <1.0.0",
		OfficialCatalogNames:  []string{"hero-banner"},
		HostNames:             []string{"feature-grid"},
		LocalNames:            []string{"local-only"},
	})

	if !result.Valid {
		t.Fatalf("expected valid result, got errors %v", result.Errors)
	}
}
