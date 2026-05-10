package source

import (
	"context"
	"strings"
	"testing"

	"github.com/phuhh98/sloth/packages/cli/pkg/registry"
)

type stubOCIClient struct {
	release  *registry.Release
	err      error
	seenVers []string
}

func (s *stubOCIClient) PullRelease(_ context.Context, requestedVersion string) (*registry.Release, string, error) {
	s.seenVers = append(s.seenVers, requestedVersion)
	if s.err != nil {
		return nil, "", s.err
	}
	return s.release, requestedVersion, nil
}

func TestOCIRegistryResolverListContractsFallbacks(t *testing.T) {
	t.Parallel()

	client := &stubOCIClient{
		release: &registry.Release{
			Version:       "0.2.0",
			SchemaVersion: "0.0.1",
			Contracts: []registry.Contract{
				{
					Name:        "hero-banner",
					Label:       "Hero Banner",
					ContentHash: "abc",
					Payload: map[string]any{
						"name":          "hero-banner",
						"schemaVersion": "0.0.1",
					},
				},
			},
		},
	}

	resolver := NewOCIRegistryResolver(client)
	contracts, err := resolver.ListContracts("latest")
	if err != nil {
		t.Fatalf("list contracts: %v", err)
	}
	if len(contracts) != 1 {
		t.Fatalf("expected one contract, got %d", len(contracts))
	}
	if contracts[0].Version != "0.2.0" {
		t.Fatalf("expected fallback to release version, got %q", contracts[0].Version)
	}
	if contracts[0].SchemaVersion != "0.0.1" {
		t.Fatalf("expected fallback to release schemaVersion, got %q", contracts[0].SchemaVersion)
	}
	if len(client.seenVers) != 1 || client.seenVers[0] != "latest" {
		t.Fatalf("expected resolver to pass through requested version, got %v", client.seenVers)
	}
}

func TestOCIRegistryResolverGetContractMissing(t *testing.T) {
	t.Parallel()

	resolver := NewOCIRegistryResolver(&stubOCIClient{
		release: &registry.Release{
			Version: "0.2.0",
			Contracts: []registry.Contract{
				{Name: "hero-banner"},
			},
		},
	})

	_, err := resolver.GetContract("0.2.0", "faq-list")
	if err == nil {
		t.Fatalf("expected missing contract error")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestOCIRegistryResolverPropagatesClientError(t *testing.T) {
	t.Parallel()

	resolver := NewOCIRegistryResolver(&stubOCIClient{err: context.DeadlineExceeded})
	_, err := resolver.ListContracts("latest")
	if err == nil {
		t.Fatalf("expected resolver error")
	}
}
