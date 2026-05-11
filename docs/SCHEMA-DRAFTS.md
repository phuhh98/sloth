---
purpose: "Component contract schema design decisions and finalized v0.0.1 spec. Source of truth for schema structure choices."
status: "active"
owner: "platform-and-cli"
last_updated: "2026-05-12"
related_docs:
  - "docs/IDEAS.md"
  - "docs/COMPONENT-CONTRACTS.md"
  - "docs/MILESTONES.md"
  - "apps/docs/docs/consumers/component-contract-schema.md"
  - "apps/docs/docs/repo-developers/contract-schema-v001-implementation.md"
---

# sloth Schema Drafts

Date: 2026-05-10
Status: Active (v0.0.1 finalized 2026-05-12)

## v0.0.1 — Finalized (2026-05-12)

The component contract schema v0.0.1 is finalized. Key decisions recorded here for reference.

### Schema file locations

- Source: `packages/contracts/src/schemas/component-contract/0.0.1/schema.json`
- Docs-hosted (canonical `$schema` URL): `apps/docs/static/schemas/component-contract/0.0.1/schema.json`
- Working draft: `notebooks/component-contract.json`

### Key design decisions

**`schemaVersion` uses `const` not free string.**
Value is locked to `"0.0.1"` in the schema file itself with `"const": "0.0.1"` + `"default": "0.0.1"`. Validators reject any other value automatically. When a new version ships, a new schema file with a different `const` is published — no additional plugin-side version string matching at the validation layer.

**`kind` drives required config via `allOf if/then` at top level.**

- `kind: "layout"` → `layoutConfig` required
- `kind: "block"` → `blockConfig` required
- `kind: "section"` → `sectionConfig` optional (omit for standalone, provide for zone-placed)

**`dataset` conditional requirements are inline, not in `$defs`.**
Previous iteration defined `datasetItemOptionRule` and `datasetItemRelationRule` in `$defs` but never referenced them — they were silently orphaned. The final design inlines the `allOf if/then` rules directly in `dataset.items`. Do not move these to `$defs` without a matching `$ref`.

**`dataset.type` enum**: `"string"`, `"number"`, `"option"`, `"relation"`, `"dynamic"`. The old `"zone"` and `"component"` types were removed. Zones are now free-identifier keys in `layoutConfig.zones`, not dataset fields.

**`relationConfig.resolve`** has two modes:

- `"scalar"` + required `path` field → plugin extracts one scalar value from the related entry
- `"documentId"` → plugin returns the raw Strapi `documentId`; renderer fetches the full entry

**`renderMeta.supportsChildren` removed.** Puck's unnamed default DropZone fallback has no place in a design that uses explicit named zones in `layoutConfig.zones`. Removed to prevent ambiguity.

**Responsive uses an array of breakpoint objects**, not a single enum string. Three semantic breakpoints: `mobile` (no prefix, < 768 px), `tablet` (`md:`, ≥ 768 px), `desktop` (`xl:`, ≥ 1280 px).

**`layoutConfig.zones[].key` is a free identifier** chosen by the contract author. It does not have to match any `dataset` key — the `"zone"` type was removed from `dataset.type`.

**Gap enum** maps to Tailwind gap utilities: `none→gap-0`, `xs→gap-1`, `sm→gap-2`, `md→gap-4`, `lg→gap-6`, `xl→gap-8`.

### Public docs

- User guide: `apps/docs/docs/consumers/component-contract-schema.md`
- Implementation notes: `apps/docs/docs/repo-developers/contract-schema-v001-implementation.md`

## Architecture Review Notes (2026-05-10)

- CLI is the enablement layer: users pull remote contracts, then push verified payloads to host ingest endpoints.
- The host/plugin API surface may be derived from the OpenAPI simulation example, but its behavior boundary stays stable.
- Plugin lifecycle materialization is an implementation concern after ingest; it is not the source of truth for contract shape.
- The component contract schema is the primary schema artifact. Keep Strapi content-type implementation out of this cycle.
- `Page` does not require persisted `contractRefs` for this cycle.
- Contract linkage can be derived at runtime from `puckConfig` and component records.
- Keep the page model minimal now: page identity, route, dataset, and editor/runtime config.
- Revisit explicit `contractRefs` only if runtime profiling shows a concrete query/performance bottleneck.

## 1) Component Contract Schema — v0.0.1 (finalized)

> The design described in this section supersedes all earlier draft ideas.
> The canonical source of truth is `notebooks/component-contract.json` (working draft) and
> `packages/contracts/src/schemas/component-contract/0.0.1/schema.json` (canonical).

The schema is kind-driven. `kind` determines which config block is required:

- `layout` — author defines a free-identifier zone grid via `layoutConfig`. No fixed presets.
- `section` — standalone (omit `sectionConfig`) or zone-placed (provide `sectionConfig` with `colSpan`).
- `block` — always zone-placed; declares its column/row span via `blockConfig`.

### What was removed from the early draft and why

| Removed field / concept                                       | Reason                                                                                                |
| ------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- |
| `version` (semver)                                            | Redundant with `schemaVersion`. Removed to reduce surface area.                                       |
| `layoutPreset` enum (`header`, `footer`, `content-12cols`, …) | Replaced by free-form `layoutConfig` — the contract author defines columns and zones directly.        |
| `gridSpan` / `responsiveGrid` ($defs)                         | Replaced by `layoutConfig.zones[].span`, `blockConfig.colSpan`, and per-config `responsive` arrays.   |
| `componentKind` (business classifier)                         | Not implemented in v0.0.1. May revisit if needed for admin filtering.                                 |
| `seo` top-level field                                         | SEO is a page-level concern, not a component contract concern. Not included.                          |
| `$defs` (seo, layoutPreset, gridSpan, breakpointGridSpec)     | All definitions were either removed or inlined. No `$defs` in v0.0.1.                                 |
| `dataset.type: "zone"` and `"component"`                      | Zones are now free-identifier keys in `layoutConfig.zones`. Zone/component types in dataset are gone. |
| `valueDropdown`                                               | Redesigned and renamed to `relationConfig` with `resolve: "scalar" \| "documentId"` modes.            |
| `renderMeta.supportsChildren`                                 | Puck unnamed default DropZone fallback conflicts with explicit `layoutConfig.zones`. Removed.         |

### v0.0.1 top-level required fields

`name`, `label`, `kind`, `schemaVersion`, `dataset`

### `schemaVersion`

Uses `"const": "0.0.1"` + `"default": "0.0.1"`. Not a free string. Any value other than `"0.0.1"` fails validation.

### `dataset.type` enum (v0.0.1)

`"string"`, `"number"`, `"option"`, `"relation"`, `"dynamic"`

`"option"` requires `options[]`. `"relation"` requires `relationConfig` with `resolve: "scalar" | "documentId"`.
When `resolve` is `"scalar"`, `path` is also required. All enforced via inline `allOf if/then` — no `$defs`.

### Responsive breakpoints

Three semantic values: `"mobile"` (no Tailwind prefix, < 768 px), `"tablet"` (`md:`, ≥ 768 px), `"desktop"` (`xl:`, ≥ 1280 px).
Breakpoints are entries in a `responsive` array on `layoutConfig`, `blockConfig`, or `sectionConfig` — not a separate top-level field.

## 2) Page Composition Draft

Current Milestone 1 design profile:

- active fields in the page model: `name`, `label`, `pageType`, `route`, `dataset`, `puckConfig`
- deferred fields: `compiledConfig`, `seo`
- removed from the page model for now: persisted `contractRefs`
- page SEO is a page-level concern; the component contract schema does not define a `seo` object — page SEO must be designed independently

The page model remains minimal and runtime-oriented.

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://schemas.sloth.dev/page-template/1.0.0",
  "title": "Sloth Page Template",
  "type": "object",
  "required": ["name", "label", "pageType", "route", "dataset", "puckConfig"],
  "additionalProperties": false,
  "$defs": {
    "seo": {
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "title": { "type": "string" },
        "description": { "type": "string" },
        "canonicalUrl": { "type": "string" },
        "noIndex": { "type": "boolean" }
      }
    }
  },
  "properties": {
    "name": {
      "type": "string",
      "minLength": 1
    },
    "label": {
      "type": "string",
      "minLength": 1
    },
    "pageType": {
      "type": "string",
      "enum": ["static", "dynamic"]
    },
    "route": {
      "type": "string",
      "minLength": 1
    },
    "dataset": {
      "type": "json",
      "description": "First-level linked data bindings for the page"
    },
    "puckConfig": {
      "type": "json",
      "description": "Source of truth for the editor"
    },
    "compiledConfig": {
      "type": "json"
    },
    "seo": {
      "$ref": "#/$defs/seo"
    }
  }
}
```

Notes:

- Contract linkage should be derived at runtime from `puckConfig` and component records.
- Keep page composition generic so the builder can evolve without forcing contract schema churn.
- Revisit persisted contract references only if profiling shows a concrete query or rendering bottleneck.

## 3) Lifecycle and Host Boundary Notes

- CLI `pull` resolves remote contract releases and materializes local verified payloads.
- CLI `push` sends verified contract payloads to the host ingest endpoint.
- The plugin lifecycle hook then analyzes the pushed contract JSON and creates or updates Strapi components.
- The OpenAPI file in `packages/contracts/openapi/sloth-api.openapi.yaml` is the simulation contract shape for that interaction; the implementation can derive equivalent routes if needed.

These notes are design constraints for the next implementation phase, not direct plugin instructions.
