---
purpose: "Define component contract distribution, compatibility checks, and CLI behavior for sloth plugin versions."
status: "draft"
owner: "platform-and-cli"
last_updated: "2026-05-10"
related_docs:
  - "docs/IDEAS.md"
  - "docs/REGISTRY.md"
  - "docs/MILESTONES.md"
  - "docs/SCHEMA-DRAFTS.md"
---

# sloth Component Contracts

Date: 2026-05-10
Status: Draft

## 1) Goal

Enable sloth CLI to push component contracts to a Strapi host based on a target sloth Strapi plugin version.

This feature supports:

- clean project bootstrap from a small starter component set
- incremental adoption of newer component contracts over time
- strong compatibility guardrails to avoid schema mismatches
- custom contract authoring in derived projects without forcing sloth official contracts

Boundary rule:

- host API exposes inspection and ingest endpoints
- CLI owns verification and comparison workflow

## 2) Contract Source Resolution

For a requested plugin version, CLI resolves contract schemas from one of two sources:

- npm source:
  - fetch schema bundle from published sloth Strapi plugin package at specified version
- git source:
  - fetch schema bundle from raw git-hosted files at a specific tag, branch, or commit

Source selection should be explicit via command flags.

## 3) Compatibility Rules

Before writing any contract to Strapi host, CLI must validate compatibility:

- read current schema compatibility version from target Strapi host
- compare with requested plugin version compatibility range
- if incompatible, abort the operation and return a clear error

Abort behavior is mandatory. No partial writes when compatibility fails.

Verification ownership:

- CLI performs business verification before push
- host does not expose verify endpoints for CLI verification workflow

## 4) Schema Evolution Policy

Non-breaking policy for general component schema:

- preserve initial general $schema contract shape as much as possible
- prefer adding new components over changing existing contract semantics
- when contract changes are required, gate them behind explicit version compatibility checks

## 5) CLI Behaviors

### 5.1 List Contracts

CLI can list latest available component contracts with plugin version header and component summary list.

CLI should also inspect host status before push planning.

Output should include:

- sloth plugin version and source
- contract schema version
- host component inventory summary (when `inspect` mode is used)
- list of components with general metadata:
  - name
  - label
  - type
  - contract version

### 5.2 Add Contracts

CLI supports add modes:

- add single component contract
- add a named set of component contracts
- add all available contracts

All add modes must run compatibility validation before write.

### 5.3 Verify Contracts

CLI must provide a separate verify command for custom contract files.

Verification checks include:

- schema compatibility version alignment
- contract name collision detection against latest compatible official contract catalog
- contract name collision detection against host-existing contract/component names

If any check fails, verification returns blocking errors and push/add operations must not proceed.

## 6) Proposed Command Shape

- `sloth contracts list --plugin-version <version> [--source npm|git]`
- `sloth contracts inspect --host <url>`
- `sloth contracts add component --name <component> --plugin-version <version> [--source npm|git]`
- `sloth contracts add set --name <set-name> --plugin-version <version> [--source npm|git]`
- `sloth contracts add --all --plugin-version <version> [--source npm|git]`
- `sloth contracts verify --file <contract.json> [--plugin-version <version>]`
- `sloth contracts push --plugin-version <version> [--source npm|git]`

Push pipeline:

1. verify local contracts
2. inspect remote host state
3. compare drift (missing/extra/update)
4. push to host ingest endpoint
5. host materializes component content-type records

Host APIs for CLI:

- `GET /sloth/inspection/plugin-status`
- `GET /sloth/inspection/contract-schema?schemaVersion=<version>&inline=<boolean>`
- `POST /sloth/contracts/ingest`

Source flags:

- npm mode:
  - `--npm-package <name>` optional override
- git mode:
  - `--git-url <raw-base-url>`
  - `--git-ref <tag|branch|commit>`

Behavior flags:

- `--dry-run`
- `--yes-to-all`
- `--host`
- `--authorization-token`

## 7) Error Contract

When incompatible version is detected, return a blocking error:

- code: `ERR_SCHEMA_VERSION_INCOMPATIBLE`
- include:
  - host schema version
  - requested plugin version
  - supported version range
  - remediation hint (use compatible version or update plugin first)

## 8) Milestone Mapping

This feature maps primarily to Milestone 2 in docs/MILESTONES.md:

- CLI sync and robust diffing
- local and remote compatibility-safe updates

## 9) Initial Draft Schemas and Content-Types

Initial JSON drafts for implementation are documented in:

- docs/SCHEMA-DRAFTS.md
