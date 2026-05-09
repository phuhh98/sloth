package compat

import (
	"fmt"
	"strings"

	semver "github.com/Masterminds/semver/v3"
)

func SchemaCompatible(hostSchemaVersion string, contractSchemaVersion string) (bool, error) {
	host, err := semver.NewVersion(normalize(hostSchemaVersion))
	if err != nil {
		return false, fmt.Errorf("invalid host schema version: %w", err)
	}
	contract, err := semver.NewVersion(normalize(contractSchemaVersion))
	if err != nil {
		return false, fmt.Errorf("invalid contract schema version: %w", err)
	}

	if host.Major() != contract.Major() {
		return false, nil
	}
	if host.LessThan(contract) {
		return false, nil
	}

	return true, nil
}

func PluginVersionInRange(pluginVersion string, supportedRange string) (bool, error) {
	if strings.TrimSpace(supportedRange) == "" {
		return true, nil
	}

	constraint, err := semver.NewConstraint(supportedRange)
	if err != nil {
		return false, fmt.Errorf("invalid supported range: %w", err)
	}
	pv, err := semver.NewVersion(normalize(pluginVersion))
	if err != nil {
		return false, fmt.Errorf("invalid plugin version: %w", err)
	}

	return constraint.Check(pv), nil
}

func normalize(version string) string {
	trimmed := strings.TrimSpace(version)
	if strings.HasPrefix(trimmed, "v") {
		return trimmed[1:]
	}
	return trimmed
}
