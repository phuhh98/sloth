package source

import (
	"context"

	"github.com/phuhh98/sloth/packages/cli/pkg/registry"
)

type OCIReleaseClient interface {
	PullRelease(ctx context.Context, requestedVersion string) (*registry.Release, string, error)
}

type OCIRegistryResolver struct {
	client OCIReleaseClient
}

func NewOCIRegistryResolver(client OCIReleaseClient) *OCIRegistryResolver {
	return &OCIRegistryResolver{client: client}
}

func (r *OCIRegistryResolver) ListContracts(pluginVersion string) ([]Contract, error) {
	release, _, err := r.client.PullRelease(context.Background(), pluginVersion)
	if err != nil {
		return nil, err
	}

	contracts := make([]Contract, 0, len(release.Contracts))
	for _, contract := range release.Contracts {
		version := contract.Version
		if version == "" {
			version = release.Version
		}
		schemaVersion := contract.SchemaVersion
		if schemaVersion == "" {
			schemaVersion = release.SchemaVersion
		}
		contracts = append(contracts, Contract{
			Name:          contract.Name,
			Label:         contract.Label,
			Version:       version,
			SchemaVersion: schemaVersion,
			ContentHash:   contract.ContentHash,
			Payload:       contract.Payload,
		})
	}

	return contracts, nil
}

func (r *OCIRegistryResolver) GetContract(pluginVersion string, name string) (*Contract, error) {
	contracts, err := r.ListContracts(pluginVersion)
	if err != nil {
		return nil, err
	}
	return FindByName(contracts, name)
}
