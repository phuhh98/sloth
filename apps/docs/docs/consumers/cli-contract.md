---
title: CLI Command Reference
---

# CLI Command Reference

The sloth CLI manages component contracts only.

For onboarding and workflow guidance, start with [CLI Getting Started](./cli-getting-started).

## Commands

- `sloth init`
- `sloth contracts ls --version <version|latest> [--source local|oci] [--format table|json]`
- `sloth contracts inspect --profile <name> [--format table|json]`
- `sloth contracts pull --name <component> --version <version|latest> [--out <path>]`
- `sloth contracts add component --name <component> --version <version|latest>`
- `sloth contracts add set --name <set-name> --version <version|latest>`
- `sloth contracts add --all --version <version|latest>`
- `sloth contracts verify --file <contract.json> --version <version|latest>`
- `sloth contracts push --version <version|latest> [--dry-run] [--retries <n>]`

## Local Workspace

```text
.sloth/
  config.yaml
  contracts/
  sets/
  manifests/
    lock.json
```

### Configuration

`config.yaml` supports profile-based host and token setup:

```yaml
currentProfile: default
profiles:
  default:
    host: http://localhost:1337
    authorizationToken: ""
```

The CLI resolves config using explicit precedence:

1. YAML values from `.sloth/config.yaml`
2. Environment variables (`SLOTH_HOST`, `SLOTH_AUTHORIZATION_TOKEN`, etc.)
3. Built-in defaults
4. Runtime flags (explicit override)

For detailed config examples and environment variable setup, see [CLI Getting Started — Configure Host Profiles](./cli-getting-started#configure-host-profiles).

## Verification Rules

`sloth contracts verify` checks:

- schema compatibility between host schema and contract schema
- plugin version compatibility range
- name collisions with official catalog
- name collisions with host inventory

When schema compatibility fails, the command emits `ERR_SCHEMA_VERSION_INCOMPATIBLE` and exits with a non-zero status.

## Push Pipeline

`sloth contracts push` follows this sequence:

1. Read local contracts from `.sloth/contracts`.
2. Inspect host state from `/sloth/inspection/plugin-status`.
3. Compute drift summary (local missing on host).
4. Call `/sloth/contracts/ingest`.
5. Update `.sloth/manifests/lock.json` with sync metadata.

## Distribution and Release

Use Taskfile commands for build and packaging:

- `task cli-build-cross`
- `task cli-generate-publish-packages`
- `task cli-smoke-install`
- `task cli-release-prep`

Use OCI registry lifecycle commands for local real-registry testing:

- `task oci-registry-up`
- `task oci-registry-down`

Generated artifacts:

- `packages/cli/dist/bin/<os>-<arch>/sloth[.exe]`
- `packages/cli/dist/checksums.txt`
- `packages/cli/dist/publish-packages/*`

## Related

- [CLI Getting Started](./cli-getting-started)
- [CLI Validation and Testing](../repo-developers/cli-validation-and-testing)
- [CLI Distribution and Release](../repo-developers/cli-distribution-and-release)
