# @sloth/cli Developer README

This document is for contributors working on the CLI implementation.

## Purpose

The CLI is implemented in Go (Cobra) and distributed via npm-compatible packaging.

Primary scope:

- component contract workflows (list, inspect, pull, add, verify, push)
- local `.sloth` workspace management
- cross-platform binary build and package generation

Out of scope:

- page content authoring operations

## Package Layout

- `cmd/sloth/main.go`: binary entrypoint
- `internal/app/`: command wiring and command handlers
- `pkg/config/`: `.sloth/config.yaml` parsing and profile resolution
- `pkg/host/`: HTTP client for host endpoints
- `pkg/source/`: contract source resolution
- `pkg/workspace/`: local `.sloth` folder operations
- `pkg/lock/`: lock file read/write/update
- `pkg/compat/`: schema and version compatibility checks
- `pkg/verify/`: verification and collision rules
- `pkg/output/`: table and json output helpers
- `bin/sloth.js`: npm runtime wrapper to execute platform binary
- `scripts/build-cross.mjs`: cross-platform Go build and checksums
- `scripts/generate-publish-packages.mjs`: generated publish folders
- `scripts/smoke-install.mjs`: wrapper smoke test

## Local Development

From repository root:

```bash
task cli-test
task cli-build-cross
task cli-generate-publish-packages
task cli-smoke-install
```

Direct package commands:

```bash
pnpm --filter @sloth/cli test
pnpm --filter @sloth/cli run build:go
pnpm --filter @sloth/cli run build:publish-packages
pnpm --filter @sloth/cli run smoke:install
```

## Configuration Precedence

Runtime configuration resolves in this order:

1. YAML profile values from `.sloth/config.yaml`
2. environment variables
3. built-in defaults

Environment variables:

- `SLOTH_CONFIG`: config path fallback when `--config` is not provided
- `SLOTH_PROFILE`: profile fallback when `--profile` is not provided
- `SLOTH_HOST`: host fallback when YAML host is unavailable
- `SLOTH_AUTHORIZATION_TOKEN` (or legacy `SLOTH_TOKEN`): token fallback when YAML token is unavailable
- `SLOTH_REGISTRY_HOST`: OCI registry host fallback when YAML registry host is unavailable
- `SLOTH_REGISTRY_REPOSITORY`: OCI repository fallback when YAML registry repository is unavailable
- `SLOTH_REGISTRY_USE_AUTHORIZATION_TOKEN`: OCI token usage flag fallback (`true|false`)

OCI registry config lives under each profile:

```yaml
currentProfile: default
profiles:
  default:
    host: http://localhost:1337
    authorizationToken: ""
    registry:
      host: ghcr.io
      repository: phuhh98/sloth/contracts
      useAuthorizationToken: true
```

Use OCI-backed contract source:

```bash
sloth contracts ls --source oci --version latest
```

Built-in defaults:

- profile name: `default`
- host: `http://localhost:1337`

Note:

- CLI flags still act as explicit overrides when provided (for example `--host`, `--authorization-token`).

## Testing Strategy

Current layers:

- unit tests for config, lock, workspace, compat, verify
- command integration tests in `internal/app/root_integration_test.go`
- mocked host endpoints via Go `httptest`

Why mocks are used:

- keep CI and local runs deterministic
- avoid requiring always-on Strapi host for most command behavior checks

Known gap:

- mock tests do not fully replace live host validation for auth/policy and real ingest behavior

## Host API Assumptions

Current CLI host client targets:

- `GET /sloth/inspection/plugin-status`
- `GET /sloth/inspection/contract-schema`
- `POST /sloth/contracts/ingest`

When plugin routes change, update both:

1. `pkg/host/client.go`
2. integration tests in `internal/app/root_integration_test.go`

## Release Workflow (Developer)

Standard prep flow:

```bash
task cli-release-prep
```

Expected generated artifacts:

- `dist/bin/<os>-<arch>/sloth[.exe]`
- `dist/checksums.txt`
- `dist/publish-packages/*`

`dist/` is generated output and is ignored by git.

## Contributor Notes

- Keep command output stable for both `--format table` and `--format json`.
- Avoid introducing network calls outside host client paths.
- Add or update tests for all behavior changes.
- Keep docs in sync with user-facing command behavior under `apps/docs/docs/`.
