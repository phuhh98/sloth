---
purpose: "Track current implementation status mapped to roadmap and priorities defined in docs/IDEAS.md."
status: "active"
owner: "project-core"
last_updated: "2026-05-12"
related_docs:
  - "docs/IDEAS.md"
  - "docs/REGISTRY.md"
  - "docs/COMPONENT-CONTRACTS.md"
---

# sloth Milestones

Date: 2026-05-10
Status source: docs/IDEAS.md (Section 14)

## Status Legend

- Not Started
- In Progress
- Completed
- Blocked

## Milestone 1

Goal:

- finalize plugin content-types and API contracts
- ship admin Puck editor MVP

Status: In Progress
Notes:

- spec-first correction track is now active for contracts, schema semantics, and docs alignment
- CLI/host boundary remains pull -> verify -> push, with host ingest initiating plugin lifecycle materialization
- schema revision now needs fixed layout presets, section full-width semantics, block grid-span semantics, and reusable SEO object definitions
- page model remains minimal for this cycle; persisted `contractRefs` stay deferred unless profiling proves a bottleneck
- plugin implementation work is intentionally deferred until the revised schema draft and milestone board are approved
- remaining work in this cycle is documentation and planning only

## Milestone 2

Goal:

- build CLI contract list/add/verify/push with robust compatibility checks
- local contract/set folder conventions + lock file
- distribute Go CLI via npm package with prebuilt binaries for macOS, Linux, and Windows

Status: Completed
Notes:

- command shape and sync strategy documented in docs/IDEAS.md
- component contract version-gated sync and compatibility-abort behavior documented in docs/COMPONENT-CONTRACTS.md
- page content operations are intentionally excluded from CLI scope
- host API is inspection-first; CLI owns verify/compare/push workflow
- CLI implementation stays in Go + Cobra; npm is the delivery channel for platform binaries and package metadata
- implemented `packages/cli` Go + Cobra scaffold with commands: init, contracts ls/inspect/pull/add/verify/push
- implemented local `.sloth/` workspace init, config profile parsing, lock file management, compatibility checks, and collision detection
- implemented YAML/ENV/default config resolution with explicit precedence for host URL, token, and profile selection
- implemented cross-platform Go binary build pipeline with checksum generation and npm wrapper resolver
- implemented config-driven publish package generation and Taskfile release-prep commands
- added CLI tests and public docs pages at `apps/docs/docs/cli-*.md`
- **follow-up required:** set up GitHub Actions CI workflow for npm publishing, configure npm publishing credentials, and automate version tagging per release
- **follow-up required:** correct CLI schema fetch and verify behaviour — `contracts verify` only does semver range checks; with schema artifacts now published to GHCR (`ghcr.io/<owner>/sloth/schemas/component-contract:<version>`), the CLI must fetch the JSON schema document from OCI and validate the contract payload structure against it; schema registry path must also be configurable (separate from the contracts artifact path); tracked in `docs/KANBAN-CLI-SCHEMA-FETCH-VERIFY.md`

## Milestone 3

Goal:

- component hub starter packs
- runtime renderer integration examples

Status: Completed
Notes:

- finalized dynamic/static shortlist and shared base contract + SEO slot model in `docs/COMPONENT-HUB-BASE-CONTRACT-AND-SHORTLIST.md`
- implemented starter pack schema/types, validator/builder scripts, and first artifact `marketing-core@0.0.1`
- moved runtime renderer example out of active code; `packages/component-hub` is now a placeholder while contract tooling is consolidated in `packages/contracts`
- extended docs sync to publish pack artifacts under `apps/docs/static/registry/packs/`
- documented starter packs, runtime integration, and OpenAPI host contract in `apps/docs/docs/`
- added CMS-agnostic OpenAPI spec and contracts package mock server with seeded fixtures
- added CLI integration test flow against the OpenAPI mock server (`list/inspect/add/verify/push`)
- abandoned release-ledger/manifest migration path in favor of OCI registry strategy
- adopted GHCR OCI artifact distribution with ORAS-based CLI abstraction for `contracts ls` and `contracts pull`
- milestone execution board archived at `docs/archive/KANBAN-MILESTONE-3.md`
- OCI pivot execution board archived at `docs/archive/KANBAN-OCI-REGISTRY-PIVOT.md`
- unfinished migration-era tasks (test re-baseline, docs alignment, verification gate) were carried over and refined in `docs/archive/KANBAN-OCI-REGISTRY-PIVOT.md`
- OCI pivot kickoff completed task OCI-001 (CLI registry config model + precedence tests) and moved OCI-002 (ORAS client implementation) to in-progress
- OCI pivot implementation completed OCI-002 by adding ORAS-backed OCI list/pull primitives and `--source oci` resolver wiring in CLI; OCI-003 is now in-progress for migration test coverage and path hardening
- OCI pivot completed OCI-003 and OCI-004 by adding OCI list/add integration coverage, stabilizing localhost OCI behavior, introducing `contracts ls`, and adding `--version` with backward-compatible `--plugin-version`
- OCI pivot completed OCI-005 and OCI-006 by implementing `contracts pull --name --version [--out]` with lock-aware workspace writes, adding pull success/error integration coverage, and adding a dedicated Zot compose profile + Taskfile lifecycle commands for real-registry test setup
- OCI pivot completed OCI-008 by adding a GHCR publish workflow with release/manual triggers, ORAS-based contracts+schema publication, and digest outputs for downstream docs sync automation
- OCI pivot completed OCI-007 by adding an env-gated real Zot integration test (`list` + `pull` against seeded OCI artifact), plus CI/task hooks to run it in Docker-capable environments
- OCI pivot completed OCI-013 by moving component-contract schema source-of-truth to `packages/contracts/src/schemas/component-contract/<version>/schema.json`, wiring docs prebuild schema sync from contracts source, and updating GHCR publish workflow to release schema artifacts from contracts source with docs parity checks
- OCI pivot completed OCI-009 by removing migration-era compare-ref immutability validation from docs-pages CI, replacing it with OCI-era contracts registry build checks, and updating lint-staged contracts gates to run build+tests without release-migration assumptions
- OCI pivot completed OCI-011 by re-baselining CLI test expectations to `contracts ls` + `--version`, and adding OCI resolver fallback/missing-contract error coverage to prevent migration-era regressions
- OCI pivot completed OCI-012 by finishing active docs/instructions alignment for OCI flow (`contracts ls/pull`, `--version`, GHCR schema publication model, and contracts package ownership wording)
- OCI pivot completed OCI-010 full verification gate across CLI/package/docs pipelines, including real Zot OCI integration (`task cli-test-zot-integration`) and docs production build validation
- OCI pivot completed OCI-014 by dropping unused script aliases and removing obsolete Taskfile frontend task call-sites to keep execution surface minimal
- OCI pivot completed OCI-015 by revising public/internal architecture docs, improving docs navigation naming, and cleaning orphaned docs page placement
- OCI pivot completed OCI-016 by adding a versioned CLI config YAML schema in contracts, syncing it to docs static schema hosting, and documenting the canonical schema URL for user config reference

## Milestone 4

Goal:

- free public registry MVP (metadata + artifacts)

Status: Not Started
Notes:

- detailed architecture and rollout in docs/REGISTRY.md
- near-term bootstrap path uses GHCR OCI artifacts and CLI ORAS integration before introducing full metadata API

## Milestone 5

Goal:

- private/paid registry extensions

Status: Not Started
Notes:

- to start after milestone 4 baseline security is in place

## Update Rule

When work is completed or scope changes:

1. Update milestone status and notes here.
2. Reflect strategic changes in docs/IDEAS.md.
3. Update specialized design details in corresponding docs files (for example docs/REGISTRY.md).
