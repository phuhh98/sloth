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

## 1.1) Near-Term Delivery Approach (Component Hub First)

Before a standalone registry API is introduced, sloth should use `packages/component-hub` as the source of truth for component registry artifacts.

Near-term approach:

- source release-versioned contracts and release manifests from `component-hub`
- copy generated immutable artifacts into docs static hosting path:
  - `apps/docs/static/registry/contracts/`
- publish through GitHub Pages as static JSON endpoints

This enables an immediate usable registry surface without introducing API, auth, or billing complexity.

Detailed implementation steps are captured in `docs/COMPONENT-HUB-DOCS-INTEGRATION.md`.

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

Protocol:

- metadata API: REST JSON
- artifact download: HTTPS URL
- integrity: SHA256 checksum for each artifact
- optional later: signed manifest (sigstore/cosign style)

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

## 5) CLI Install Flow

Proposed flow for `sloth registry add`:

1. Resolve package metadata from registry API.
2. Choose version by compatibility constraints.
3. Download artifact from CDN.
4. Verify SHA256 checksum.
5. Extract into project folders:
   - `components/ui`
   - atomic component folders
   - sloth config/template folders
6. Write lock metadata for deterministic updates.

## 6) CLI Commands (Registry-Oriented)

Suggested commands:

- `sloth registry search <keyword>`
- `sloth registry info <package>`
- `sloth registry add <package>@<version>`
- `sloth registry update <package>|--all`
- `sloth registry remove <package>`
- `sloth registry publish` (future)

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

- component-hub sourced immutable artifacts hosted from docs static path
- static index endpoints for namespace/components/themes/packs
- checksum verification and richer metadata API can be added in next increment

Phase 2: Private packages

- PAT auth, access control, team namespaces

Phase 3: Paid marketplace

- Stripe integration, entitlement checks, billing and licensing workflows

## 9) Open Decisions

- Should package install allow file overwrite by default or prompt first?
- Should compatibility checks be strict-fail or warning-first?
- Should registry packages include executable frontend runtime code, or config-only artifacts?
- What namespace policy should be used for organizations and publishers?
