---
purpose: "Registry architecture, protocol, rollout, and governance guidance for sloth."
status: "draft"
owner: "platform"
last_updated: "2026-05-10"
related_docs:
  - "docs/IDEAS.md"
  - "docs/MILESTONES.md"
  - "docs/COMPONENT-CONTRACTS.md"
---

# sloth Registry

Date: 2026-05-10
Status: Draft

## 1) Scope and Goal

The sloth registry distributes reusable theme sets and component variants for sloth-based projects.

It should support:

- pulling theme sets or individual components
- plugin-version-scoped component contract retrieval and distribution
- compatibility checks against sloth plugin and frontend renderer versions
- delivery of shadcn/base UI compatible files into project folders
- optional future paid and private package workflows

## 1.1) Near-Term Delivery Approach (OCI Registry First)

Before a standalone registry metadata API is introduced, sloth uses OCI artifacts in GitHub Container Registry (GHCR) as the contract distribution backend.

Near-term approach:

- source contracts from `packages/contracts/src/contracts/components/`
- package the contract folder as versioned OCI artifacts in GHCR
- source component-contract schema from `packages/contracts` and publish versioned schema artifacts to GHCR
- use CLI-level abstraction over OCI internals via contract commands
- use ORAS Go SDK in CLI for pull/list behavior

Schema URL policy:

- keep `$schema` in contract files pointed to a stable HTTPS document URL
- treat GHCR as immutable artifact storage and provenance source, not the canonical `$schema` URL endpoint
- mirror promoted schema versions from contracts release artifacts into docs-hosted static paths for validator/editor compatibility

This keeps distribution immutable and CDN-backed while avoiding custom registry API complexity.

Detailed integration and execution planning is tracked in a dedicated OCI migration kanban.

## 2) Difficulty Assessment

Short answer:

- MVP free registry is medium difficulty.
- Private/auth registry is medium to high.
- Paid marketplace is high.

Main complexity is not file hosting. It is trust, entitlement, and lifecycle management:

- compatibility resolution across plugin/renderer/dependencies
- artifact integrity and package signing
- namespace ownership and moderation
- billing and entitlement enforcement

## 3) Recommended Architecture

Use a split architecture:

- Registry API service:
  - package metadata
  - version resolution
  - search
  - auth and entitlements
- Artifact storage:
  - immutable tar.gz/zip artifacts in S3-compatible storage
- CDN:
  - fast, cacheable global downloads
- Database:
  - PostgreSQL for packages, versions, owners, entitlements

Suggested stack:

- API: Node.js (Fastify/NestJS) or Go (Fiber/Chi)
- DB: PostgreSQL
- Storage: S3, R2, or MinIO
- CDN: Cloudflare or CloudFront
- Auth: OAuth + PAT
- Payments (later): Stripe

Pragmatic recommendation for sloth now:

- keep CLI in Go (Cobra)
- start registry API in Node.js for faster product iteration

## 4) Protocol and Contracts

Protocol (near-term):

- artifact protocol: OCI distribution API (GHCR)
- client library: ORAS Go SDK
- integrity: OCI digest verification plus optional SHA256 metadata in artifact annotations
- optional later: signature verification (cosign/sigstore)

Core package manifest example:

```json
{
  "name": "marketing-hero-pack",
  "version": "1.0.0",
  "kind": "theme",
  "compatibility": {
    "sloth": ">=0.1.0",
    "renderer": ">=0.1.0"
  },
  "artifacts": {
    "components": "./components",
    "templates": "./templates",
    "ui": "./components/ui"
  }
}
```

## 5) CLI Pull Flow

Proposed user-facing flow keeps `contracts` abstraction and hides OCI internals:

1. `sloth contracts ls --version <x.y.z|latest>` resolves available contracts for a release.
2. `sloth contracts pull --name <contract> --version <x.y.z>` pulls a single contract from OCI-backed release payload.
3. CLI verifies compatibility constraints and writes local contract files.
4. CLI writes lock metadata for deterministic updates.

## 6) CLI Commands (Contract-Oriented)

Suggested commands:

- `sloth contracts ls --version <x.y.z|latest>`
- `sloth contracts pull --name <contract> --version <x.y.z|latest>`

Optional future expansion:

- `sloth contracts pull --all --version <x.y.z|latest>`

Useful flags:

- `--host`
- `--authorization-token`
- `--yes-to-all`
- `--dry-run`
- `--format json|table`

## 7) Security Baseline

Before public launch, enforce:

- immutable published versions
- checksum verification on every install
- namespace ownership checks
- moderation workflow for content
- audit logs for publish/unpublish and access changes

## 8) Rollout Plan

Phase 1: Free public registry MVP

- contracts package sourced immutable artifacts published to GHCR as OCI artifacts
- component-contract schema versions published to GHCR as OCI artifacts from `packages/contracts` source
- docs-hosted schema URL retained as canonical `$schema` endpoint and synchronized from published schema artifacts
- CLI ORAS-based list and pull behavior under `contracts` command group
- integration testing against local Zot OCI registry in docker compose on-demand profile

Phase 2: Private packages

- PAT auth, access control, team namespaces

Phase 3: Paid marketplace

- Stripe integration, entitlement checks, billing and licensing workflows

## 9) Open Decisions

- Should package install allow file overwrite by default or prompt first?
- Should compatibility checks be strict-fail or warning-first?
- Should registry packages include executable frontend runtime code, or config-only artifacts?
- What namespace policy should be used for organizations and publishers?
- Should schema publication and contract publication share a single OCI artifact layout or separate repositories in GHCR?
