package verify

import (
	"fmt"
	"sort"
	"strings"

	"github.com/phuhh98/sloth/packages/cli/pkg/compat"
)

type Input struct {
	ContractName          string
	ContractSchemaVersion string
	HostSchemaVersion     string
	PluginVersion         string
	SupportedRange        string
	OfficialCatalogNames  []string
	HostNames             []string
	LocalNames            []string
}

type Result struct {
	Valid    bool     `json:"valid"`
	Errors   []string `json:"errors"`
	Warnings []string `json:"warnings"`
}

func Run(input Input) Result {
	res := Result{Valid: true}

	if strings.TrimSpace(input.ContractName) == "" {
		res.addError("contract name is required")
	}

	if strings.TrimSpace(input.HostSchemaVersion) != "" && strings.TrimSpace(input.ContractSchemaVersion) != "" {
		ok, err := compat.SchemaCompatible(input.HostSchemaVersion, input.ContractSchemaVersion)
		if err != nil {
			res.addError(fmt.Sprintf("schema compatibility check failed: %v", err))
		} else if !ok {
			res.addError(fmt.Sprintf("ERR_SCHEMA_VERSION_INCOMPATIBLE host=%s contract=%s", input.HostSchemaVersion, input.ContractSchemaVersion))
		}
	}

	if strings.TrimSpace(input.PluginVersion) != "" {
		ok, err := compat.PluginVersionInRange(input.PluginVersion, input.SupportedRange)
		if err != nil {
			res.addError(fmt.Sprintf("plugin range check failed: %v", err))
		} else if !ok {
			res.addError(fmt.Sprintf("plugin version %s not supported by range %s", input.PluginVersion, input.SupportedRange))
		}
	}

	official := normalizeSet(input.OfficialCatalogNames)
	hostSet := normalizeSet(input.HostNames)
	localSet := normalizeSet(input.LocalNames)
	name := strings.ToLower(strings.TrimSpace(input.ContractName))
	if name != "" {
		if official[name] {
			res.addError(fmt.Sprintf("contract name collision with official catalog: %s", input.ContractName))
		}
		if hostSet[name] {
			res.addError(fmt.Sprintf("contract name collision with host inventory: %s", input.ContractName))
		}
		if localSet[name] {
			res.addWarning(fmt.Sprintf("contract name already exists in local workspace: %s", input.ContractName))
		}
	}

	sort.Strings(res.Errors)
	sort.Strings(res.Warnings)
	res.Valid = len(res.Errors) == 0
	return res
}

func normalizeSet(values []string) map[string]bool {
	out := map[string]bool{}
	for _, value := range values {
		trimmed := strings.ToLower(strings.TrimSpace(value))
		if trimmed == "" {
			continue
		}
		out[trimmed] = true
	}
	return out
}

func (r *Result) addError(message string) {
	r.Valid = false
	r.Errors = append(r.Errors, message)
}

func (r *Result) addWarning(message string) {
	r.Warnings = append(r.Warnings, message)
}
