---
purpose: "Track current implementation status mapped to roadmap and priorities defined in docs/IDEAS.md."
status: "active"
owner: "project-core"
last_updated: "2026-05-10"
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

- implemented backend content-types draft in plugin for component and page
- implemented content-api slice:
  - `GET /inspection/plugin-status`
  - `GET /inspection/contract-schema`
  - `POST /contracts/ingest`
  - `GET /pages/:id/delivery`
- current page delivery behavior returns page record payload; first-level linked content population is planned
- contract ingest currently materializes component records via Document Service API
- page model intentionally excludes persisted `contractRefs` in current phase
- remaining work includes admin Puck editor MVP and stronger runtime/delivery enrichment

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
- implemented `packages/cli` Go + Cobra scaffold with commands: init, contracts list/inspect/add/verify/push
- implemented local `.sloth/` workspace init, config profile parsing, lock file management, compatibility checks, and collision detection
- implemented YAML/ENV/default config resolution with explicit precedence for host URL, token, and profile selection
- implemented cross-platform Go binary build pipeline with checksum generation and npm wrapper resolver
- implemented config-driven publish package generation and Taskfile release-prep commands
- added CLI tests and public docs pages at `apps/docs/docs/cli-*.md`
- **follow-up required:** set up GitHub Actions CI workflow for npm publishing, configure npm publishing credentials, and automate version tagging per release

## Milestone 3

Goal:

- component hub starter packs
- runtime renderer integration examples

Status: Not Started
Notes:

- depends on stable plugin and CLI contract
- include rendering strategy decision: sample frontend-only or extracted SDK
- component-hub should own contract + basic React implementation sources
- docs build should sync generated release-versioned contract sets into `apps/docs/static/registry/contracts/`

## Milestone 4

Goal:

- free public registry MVP (metadata + artifacts)

Status: Not Started
Notes:

- detailed architecture and rollout in docs/REGISTRY.md
- near-term bootstrap path uses static docs hosting with component-hub generated artifacts before introducing full metadata API

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
