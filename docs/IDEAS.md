---
purpose: "Primary brainstorming and evolving product direction for sloth."
status: "draft"
owner: "product-and-architecture"
last_updated: "2026-05-10"
related_docs:
  - "docs/MILESTONES.md"
  - "docs/REGISTRY.md"
  - "docs/COMPONENT-CONTRACTS.md"
  - "docs/SCHEMA-DRAFTS.md"
  - "docs/ARCHITECTURE-DIAGRAM.md"
---

# sloth - Brainstorming and Product Direction

Date: 2026-05-10
Status: Draft brainstorming

## 1) Product Vision

sloth aims to provide a complete page-template platform on top of Strapi:

- A Strapi plugin to define and manage page templates (static and dynamic).
- A drag-and-drop Admin UI powered by Puck for composing templates.
- A plugin-owned content model for components and pages, plus a contract interface layer for frontend implementers.
- A developer CLI (Go + Cobra) focused on component contract operations, not page content operations.
- A component hub for reusable theme/component variants (shadcn/ui + Base UI foundations).
- A future registry for public/private (free and paid) templates/themes/components.

## 2) Terminology and Type Model

Use explicit naming to avoid confusion between CMS content-type kind and TypeScript types.

### 2.1 Component Contract Definition

```ts
export type DatasetValueType = "string" | "number" | "option" | "dynamic";

export interface ContractDatasetField {
  key: string; // unique within component
  label: string;
  type: DatasetValueType;
  options?: Array<{ label: string; value: string | number | boolean }>; // for option
  value?: string | number | boolean | null; // for fixed values
  valueDropdown?: {
    contentType: string; // e.g. api::article.article
    path: string; // path in the entry shape, e.g. "attributes.title"
    multiple?: boolean;
  }; // for dynamic values
  required?: boolean;
}

export interface ComponentContract {
  name: string; // globally unique slug-like key
  label: string;
  contentTypeKind: "layout" | "section" | "block";
  category?: string; // optional grouping in builder UI
  version: string; // semver-ish
  schemaVersion: string; // compatible schema track
  dataset: ContractDatasetField[];
  renderMeta?: {
    // maps to frontend renderer contract
    rendererKey: string;
    supportsChildren?: boolean;
  };
}
```

Notes:

- Add `key` for each dataset item, not only label.
- `contentTypeKind` replaces ambiguous `type` naming.
- `schemaVersion` is used for compatibility checks.
- Keep `rendererKey` separate from display `name` to allow renames.

### 2.1.1 Domain Clarification

- `Component` and `Page` are Strapi content-types.
- `Component Contract` is an interface exposed to frontend developers.
- `Component Contract` is related to `Component` but not one-to-one:
  - a contract can back multiple component records
  - component records can evolve while honoring a stable contract family
- component records represent business purpose and kind (CTA, Carousel, HeroSection, AsideLayout, Header, Footer), not concrete frontend implementation.

### 2.2 TypeScript Surfaces

TypeScript interfaces are delivery contracts for integrators, not the CMS content-type definitions themselves.

- Component federated data types:
  - derived from component contract JSON Schema or generated TypeScript types
  - used by developers in their own component implementations
- Page config types:
  - generic page interface for page rendering config and linked content metadata
  - populated payload is first-level only for linked content; deeper content fetch remains consumer responsibility

### 2.3 Page Template Definition

```ts
export interface PageDatasetBinding {
  key: string; // unique on page
  contentType: string; // e.g. api::post.post
  filters?: Record<string, unknown>;
  locale?: string;
  status?: "draft" | "published";
}

export interface PageTemplate {
  name: string; // unique key
  label: string;
  type: "static" | "dynamic";
  route: string; // e.g. /about or /blog/:slug
  dataset: PageDatasetBinding[];

  // Source of truth for editor
  puckConfig: Record<string, unknown>;

  // Optional compiled shape for frontend runtime (faster and safer)
  compiledConfig?: Record<string, unknown>;

  seo?: {
    title?: string;
    description?: string;
    noIndex?: boolean;
  };
}
```

Notes:

- Keeping both `puckConfig` and `compiledConfig` is a good design.
- Treat `puckConfig` as editable source, `compiledConfig` as derived artifact.

## 3) Strapi Plugin Scope

Recommended plugin modules:

- Content-types:
  - `plugin::sloth.component`
  - `plugin::sloth.page`
  - optional `plugin::sloth.theme` (later)
- Services:
  - ingest service (contract payload -> component content-type records)
  - inspection service (plugin status, component inventory, schema compatibility metadata)
  - compiler service (`puckConfig` -> `compiledConfig`)
  - sync service (import/export payloads)
- Controllers + Routes:
  - admin routes for builder and management
  - content-api routes for runtime frontend consumption
- Admin UI:
  - tab/page for drag-and-drop builder (Puck)
  - component palette, dataset mapping panel, preview mode

Important implementation direction for Strapi v5:

- Use Document Service API (`strapi.documents(...)`) as the data access layer.

## 4) API Surface (Proposed)

CLI-host integration should target a dedicated OpenAPI contract that can be implemented by Strapi plugin APIs or another CMS-compatible host.

### Host Inspection APIs (for CLI)

- `GET /sloth/inspection/plugin-status`
  - returns general host/plugin info for CLI comparison:
    - plugin version
    - compatible `$schema` version(s)
    - current component inventory summary
    - optional contract metadata summary
- `GET /sloth/inspection/contract-schema?schemaVersion=<version>&inline=<boolean>`
  - returns schema link or inline schema payload for requested compatibility version

### Contract Ingest API (Push Target)

- `POST /sloth/contracts/ingest`
  - receives already-verified contracts from CLI
  - host parses payload and creates/updates corresponding `Component` content-type records
  - host returns ingest summary (created/updated/skipped/failed)

### Verification Ownership Rule

- Host API does not provide contract verification workflow endpoints.
- CLI must perform verification and drift checks before calling ingest:
  - schema validation
  - compatibility checks
  - name-collision checks
  - remote comparison checks (missing/extra/pending)

### Page Runtime API

- `GET /sloth/pages/:id/delivery`
  - current behavior (Milestone 1): returns page record payload for delivery use
  - target behavior: returns page rendering config + first-level populated linked content
  - component contract information is resolved dynamically from page config and component records when needed

For deeper linked data, consumers can call default Strapi generated content-type endpoints.

This keeps runtime API generic and avoids parsing every leaf node in plugin runtime.

Current architecture note:

- do not persist `contractRefs` on page content-type for Milestone 1
- derive contract linkage at delivery/compile layer
- add persisted refs later only if measurable performance pressure appears

### Admin API

- `GET /sloth/components`
- `POST /sloth/components`
- `PUT /sloth/components/:documentId`
- `GET /sloth/pages`
- `POST /sloth/pages`
- `PUT /sloth/pages/:documentId`
- `POST /sloth/pages/:documentId/compile`

## 5) Local Filesystem Convention for CLI Sync

A stable local structure helps predictable diff/pull/push:

```txt
.sloth/
  config.yaml
  contracts/
    hero-banner@1.0.0.json
    two-column@1.1.0.json
  sets/
    marketing-core.json
    blog-core.json
  manifests/
    lock.json
```

Where:

- `config.yaml` stores host and auth profile references (do not commit tokens).
- `lock.json` stores remote documentId, version hash, lastSyncedAt.

## 6) CLI (Go + Cobra) Command Design

CLI scope is component contracts only. Page concrete content is not managed by CLI commands.

Suggested command shape:

- `sloth contracts ls --version <version|latest>`
  - show compatible contract catalog with release version header
- `sloth contracts inspect --host <url>`
  - inspect host plugin status and current component inventory summary
- `sloth contracts pull --name <name> --version <version|latest>`
  - pull an individual contract from registry for selected release
- `sloth contracts add component --name <name> --plugin-version <version>`
- `sloth contracts add set --name <set-name> --plugin-version <version>`
- `sloth contracts add --all --plugin-version <version>`
- `sloth contracts verify --file <contract.json> [--plugin-version <version>]`
  - validate schema compatibility and name-collision rules before push
- `sloth contracts push --plugin-version <version>`
  - verify -> inspect remote -> compare drift -> ingest to host

Flags:

- `-f, --force`
- `-o, --out <dir>`
- `-Y, --yes-to-all`
- `-a, --all`
- `-H, --host <url>`
- `-T, --authorization-token <token>`

Additional recommended flags:

- `--profile <name>`: select host/token profile from config file
- `--dry-run`: show what would change without applying
- `--format json|table`: output mode for CI and humans

Registry/source backend stays internal to CLI:

- contracts are published as OCI artifacts to GHCR
- CLI uses ORAS Go SDK for list and pull operations
- users do not need source flags for npm/git/raw registry URLs

Detailed behavior and error contract are defined in docs/COMPONENT-CONTRACTS.md.

## 7) Sync Strategy and Conflict Handling

Recommended behavior:

- Diff basis:
  - structural hash per document (canonicalized JSON)
  - metadata hash (name, version, updatedAt)
- Push policy:
  - default: reject if remote changed since last sync
  - force: overwrite remote after explicit confirmation
- Pull policy:
  - default with `--all`: replace local target directory
  - optional future mode: merge strategy for non-overlapping keys

Compatibility gate:

- before add/push, compare host schema compatibility version with requested contract schema version
- if incompatible, abort task with a blocking error and do not write partial results

Verification gate:

- verify command checks custom contract names do not collide with latest compatible official contracts and host-existing contracts
- push command must re-check remote state before ingest to prevent stale comparisons
- this enables derived projects to create their own contracts safely

This avoids silent drift and makes CI automation reliable.

## 8) Component Hub Direction

Component hub scope:

- Theme packs (multiple coordinated components)
- Individual component variants
- Compatibility metadata:
  - required sloth plugin version
  - required frontend renderer version
  - required base dependencies (shadcn/ui blocks)
- Preview metadata (thumbnail, description, tags)

Extensibility policy:

- sloth does not force developers to use only sloth official contracts or hub components
- derived projects can define and use custom contracts
- verification rules provide safety without restricting experimentation

A package descriptor per item can look like:

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

## 9) Registry Feasibility: How Hard Is It?

Short answer: medium difficulty for MVP, hard for production-grade paid marketplace.

### 9.1 Difficulty by Phase

- Phase A (Open free registry, no auth): Low to Medium
  - host metadata index + downloadable tarballs
  - CLI install/pull by package name
- Phase B (Auth + private packages): Medium
  - token auth, access control, rate limiting
- Phase C (Paid packages + billing + licensing): High
  - payments, invoices, entitlement checks, fraud/abuse controls

### 9.2 Main Complexity Drivers

- Version/compatibility resolution across plugin, renderer, and UI dependencies.
- Security and trust:
  - package signing
  - integrity hashes
  - provenance and malware scanning pipeline
- Paid access model:
  - entitlement service + offline/CI tokens
- Long-term maintenance:
  - immutable versions, deprecations, migration docs

## 10) Registry Architecture Recommendation

For a pragmatic start, use a split architecture:

- Registry API service:
  - metadata search, package versions, auth, entitlements
- Object storage:
  - package artifacts (tar.gz/zip)
- CDN:
  - fast global downloads
- Database:
  - package metadata, owners, versions, entitlements

Suggested stack:

- API: Node.js (Fastify or NestJS) or Go (Fiber/Chi)
- DB: PostgreSQL
- Storage: S3-compatible (S3, R2, MinIO)
- CDN: Cloudflare or CloudFront
- Auth: OAuth + PAT tokens
- Payments (if paid): Stripe
- Search: Postgres full-text first, later Meilisearch/OpenSearch

If your team is small and already TypeScript-heavy, Node.js for registry is usually faster to ship.
If you want one language for CLI + backend and prioritize throughput, Go backend is excellent.

## 11) Protocol and Package Format

Recommended transport and package contract:

- Metadata protocol: REST JSON
- Artifact protocol: HTTPS download URL + SHA256 checksum
- Optional later: signed manifest (cosign/sigstore style)

CLI flow:

1. Resolve package metadata from registry API.
2. Select version based on compatibility and constraints.
3. Download artifact from CDN URL.
4. Verify SHA256.
5. Extract into target folders:
   - `components/ui`
   - atomic component folders
   - sloth config files
6. Run optional post-install checks.

## 12) Usage Template for Developers

### Publish template

1. Create package with descriptor + artifacts.
2. Run `sloth package build` (future command).
3. Run `sloth registry publish`.

### Install template

1. Run `sloth registry search <keyword>`.
2. Run `sloth registry add <package>@<version>`.
3. Confirm files to be written and dependency additions.

### Sync with Strapi

1. Run `sloth contracts inspect --host <url>`.
2. Run `sloth contracts ls --version <version|latest>`.
3. Run `sloth contracts pull --name <contract> --version <version|latest>` (or `contracts add ...`).
4. Run `sloth contracts verify --file <contract.json>` for custom contracts.
5. Run `sloth contracts push --plugin-version <version>`.

## 13) Security and Governance Baseline

Minimum baseline before public registry launch:

- immutable versions once published
- checksum verification on every install
- namespace ownership validation
- package moderation workflow
- audit logs for publish/unpublish and entitlement changes

## 14) Suggested Execution Roadmap

- Milestone 1:
  - finalize plugin content-types and API contracts
  - ship admin Puck editor MVP
- Milestone 2:
  - build CLI contract list/add/verify/push with robust compatibility checks
  - local contract and set folder conventions + lock file
- Milestone 3:
  - component hub starter packs
  - runtime renderer integration examples and first-level page payload strategy
- Milestone 4:
  - free public registry MVP (metadata + artifacts)
- Milestone 5:
  - private/paid registry extensions

## 15) Open Questions

- Should `compiledConfig` be generated on save, on publish, or both?
- Should dynamic page routes support catch-all segments at MVP stage?
- Is local-first workflow preferred over Strapi-admin-first workflow?
- Should component hub packages include frontend runtime code, or only config/schema?
- What compatibility policy will be enforced (strict vs best-effort)?
- Should host ingest reject only malformed payloads, or also enforce minimal schema-version guards?

## 16) Frontend Rendering Strategy

Two-track strategy is recommended:

- sample frontend folder strategy:
  - provide a reference renderer app for fast adoption
  - include contract-to-component mapping and generic page runtime consumption examples
- separate SDK strategy:
  - provide a focused runtime SDK for rendering and data adapters
  - keep plugin and CLI concerns decoupled from frontend framework details

Recommended rollout:

1. Start with sample frontend for validation speed.
2. Extract stable rendering primitives into an SDK once mapping and runtime contracts settle.

---

This document is intended as a living spec. Keep it updated as decisions become concrete and move finalized architecture choices into ADRs.
