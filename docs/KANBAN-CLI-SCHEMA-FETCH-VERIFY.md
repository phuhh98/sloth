---
purpose: "CLI follow-up: fetch JSON schema from GHCR OCI and validate contract payload structure in verify and pull commands."
status: "active"
owner: "cli"
last_updated: "2026-05-12"
related_docs:
  - "docs/MILESTONES.md"
  - "docs/COMPONENT-CONTRACTS.md"
  - ".github/workflows/ghcr-contract-artifacts.yml"
---

# Kanban: CLI Schema Fetch & Verify Hardening

## Scope

- Milestone: Milestone 2 follow-up
- Goal: replace the concrete (hardcoded) schema version semver-range check in `contracts verify` with live JSON schema fetch from the GHCR-published OCI schema artifact; validate the contract payload structure against the fetched schema; make the schema registry path configurable
- Schema OCI ref pattern: `ghcr.io/<owner>/sloth/schemas/component-contract:<schemaVersion>` (media type `application/schema+json`)
- Contracts OCI ref pattern: `ghcr.io/<owner>/sloth/contracts:<version>` (unchanged)
- Constraints: reuse the existing `registry.OCIClient` and ORAS infrastructure; keep schema fetch opt-out possible via flag for offline use
- milestone_updated_at: 2026-05-12

## Task Decomposition Rules

- Each task touches one package boundary at a time.
- Schema config, schema client, schema source interface, and command wiring are separate tasks.
- Tests accompany each implementation task.

## Kanban

### To Do

- [ ] SFV-001 Add schema registry config field
  - what: extend `pkg/config` to carry a schema registry repository path
  - do: add `SchemaRepository` field to `Registry` struct (yaml: `schemaRepository`); add `EnvRegistrySchemaRepository = "SLOTH_REGISTRY_SCHEMA_REPOSITORY"` constant; set default to `phuhh98/sloth/schemas/component-contract`; wire the env override into `ResolveConfig`
  - next: SFV-002 can read the resolved value from the profile
  - deps: none
  - requires-confirmation: false
  - status: todo

- [ ] SFV-002 Add schema OCI pull method to registry client
  - what: teach `registry.OCIClient` to pull the JSON schema document for a given schema version
  - do: add `SchemaMediaType = "application/schema+json"` constant; add `PullSchema(ctx, schemaVersion string) ([]byte, error)` method that pulls the single-layer OCI artifact from `<schemaRegistryHost>/<schemaRepository>:<schemaVersion>` and returns the raw JSON bytes; add unit test with a mock layer descriptor
  - next: SFV-003 wraps this behind the source interface
  - deps: SFV-001
  - requires-confirmation: false
  - status: todo

- [ ] SFV-003 Add SchemaFetcher interface and OCI implementation in pkg/source
  - what: define a `SchemaFetcher` interface so commands stay decoupled from the OCI client
  - do: add `type SchemaFetcher interface { FetchSchema(schemaVersion string) ([]byte, error) }` in `pkg/source`; add `OCISchemaFetcher` backed by the new `OCIClient.PullSchema`; add a `NopSchemaFetcher` that returns nil for offline/test use
  - next: SFV-004 can inject the fetcher into the verify flow
  - deps: SFV-002
  - requires-confirmation: false
  - status: todo

- [ ] SFV-004 Wire SchemaFetcher into contracts verify command
  - what: validate contract payload structure against the fetched JSON schema in addition to the existing semver range check
  - do: build `SchemaFetcher` in `BuildRuntime` using the resolved schema repository; in `buildVerifyInput` / `verify.Run`, after the semver check fetch the schema document and validate the payload with `github.com/santhosh-tekuri/jsonschema/v6` (already indirectly available, or add the dep); add `--skip-schema-validation` flag to `contracts verify` for offline scenarios; add test cases for valid payload, structural violation, and skip flag behaviour
  - next: SFV-005 applies the same to contracts pull
  - deps: SFV-003
  - requires-confirmation: false
  - status: todo

- [ ] SFV-005 Wire SchemaFetcher into contracts pull command
  - what: validate the pulled contract payload against the fetched JSON schema immediately after download
  - do: after `resolvePullContract` succeeds, call the schema fetcher for the contract's `schemaVersion` and validate; emit a warning (not a hard error) when `--skip-schema-validation` is set or the fetcher returns no schema; add test for pull-with-schema-validation path
  - next: SFV-006 re-baselines integration tests
  - deps: SFV-003, SFV-004
  - requires-confirmation: false
  - status: todo

- [ ] SFV-006 Re-baseline integration tests and update docs
  - what: ensure existing CLI integration tests still pass and add new schema-validation test coverage
  - do: update `root_integration_test.go` and `root_zot_integration_test.go` mock fixtures to include a `schemaVersion` field matching a seeded schema artifact; add a schema-validation integration path to the Zot integration test (seeding the schema OCI artifact alongside the contracts OCI artifact); update CLI public docs page to mention `--skip-schema-validation` and the schema registry config key
  - next: milestone can be marked complete
  - deps: SFV-004, SFV-005
  - requires-confirmation: true
  - status: todo

### In Progress

- [ ] None

### Blocked

- [ ] None

### Done

- [ ] None
