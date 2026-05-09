package source

import "fmt"

type Contract struct {
	Name          string         `json:"name"`
	Label         string         `json:"label,omitempty"`
	Version       string         `json:"version"`
	SchemaVersion string         `json:"schemaVersion"`
	ContentHash   string         `json:"contentHash,omitempty"`
	Payload       map[string]any `json:"payload,omitempty"`
}

type Resolver interface {
	ListContracts(pluginVersion string) ([]Contract, error)
	GetContract(pluginVersion string, name string) (*Contract, error)
}

func FindByName(contracts []Contract, name string) (*Contract, error) {
	for _, contract := range contracts {
		if contract.Name == name {
			copyContract := contract
			return &copyContract, nil
		}
	}
	return nil, fmt.Errorf("contract %q not found", name)
}
