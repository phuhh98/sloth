---
purpose: "Internal architecture summary and pointer to canonical public architecture diagram."
status: "active"
owner: "product-and-architecture"
last_updated: "2026-05-10"
related_docs:
  - "docs/IDEAS.md"
  - "docs/MILESTONES.md"
  - "docs/COMPONENT-CONTRACTS.md"
  - "docs/REGISTRY.md"
  - "apps/docs/docs/repo-developers/architecture.mdx"
---

# sloth Architecture Design Diagram

Date: 2026-05-10
Source of truth: docs/IDEAS.md

## Canonical Diagram Source

The canonical architecture diagram lives in:

- `apps/docs/docs/repo-developers/architecture.mdx`

Use that file as the single source of truth for the Mermaid diagram and public architecture narrative.

## Responsibility Boundaries

- CLI owns verification workflow before push.
- Host plugin owns ingest and materialization into component records.
- Runtime delivery endpoint serves page delivery payload and first-level linked content strategy.
- Registry and component hub are later roadmap phases and remain decoupled from core plugin and CLI MVP.
- Current contract source is `packages/contracts/src/contracts`.
- Current schema source is `packages/contracts/src/schemas`.
- Contract and schema release artifacts are distributed via GHCR OCI.
- Docs-hosted schema URLs in `apps/docs/static/schemas` remain canonical `$schema` endpoints.

## Architecture Notes

- Keep architecture as a modular monolith around Strapi plugin and CLI during Milestones 1 and 2.
- Add registry complexity incrementally after stable plugin and CLI contracts are proven.
- Keep runtime API generic and avoid deep linked-data parsing in plugin runtime.

## Update Rule

- When architecture changes, update `apps/docs/docs/repo-developers/architecture.mdx` first.
- Update this file only for internal planning context or repo-specific notes not suitable for public docs.
