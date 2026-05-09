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

Status: Not Started
Notes:

- command shape and sync strategy documented in docs/IDEAS.md
- component contract version-gated sync and compatibility-abort behavior documented in docs/COMPONENT-CONTRACTS.md
- page content operations are intentionally excluded from CLI scope
- host API is inspection-first; CLI owns verify/compare/push workflow

## Milestone 3

Goal:

- component hub starter packs
- runtime renderer integration examples

Status: Not Started
Notes:

- depends on stable plugin and CLI contract
- include rendering strategy decision: sample frontend-only or extracted SDK

## Milestone 4

Goal:

- free public registry MVP (metadata + artifacts)

Status: Not Started
Notes:

- detailed architecture and rollout in docs/REGISTRY.md

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
