---
purpose: "Initial draft JSON artifacts for sloth component/page contract schema and reusable schema definitions."
status: "draft"
owner: "platform-and-cli"
last_updated: "2026-05-10"
related_docs:
  - "docs/IDEAS.md"
  - "docs/COMPONENT-CONTRACTS.md"
  - "docs/MILESTONES.md"
---

# sloth Initial Schema Drafts

Date: 2026-05-10
Status: Draft

## Architecture Review Notes (2026-05-10)

- CLI is the enablement layer: users pull remote contracts, then push verified payloads to host ingest endpoints.
- The host/plugin API surface may be derived from the OpenAPI simulation example, but its behavior boundary stays stable.
- Plugin lifecycle materialization is an implementation concern after ingest; it is not the source of truth for contract shape.
- The component contract schema is the primary schema artifact. Keep Strapi content-type implementation out of this cycle.
- `Page` does not require persisted `contractRefs` for this cycle.
- Contract linkage can be derived at runtime from `puckConfig` and component records.
- Keep the page model minimal now: page identity, route, dataset, and editor/runtime config.
- Revisit explicit `contractRefs` only if runtime profiling shows a concrete query/performance bottleneck.

## 1) Component Contract Schema Draft

The component contract schema should remain kind-driven and kind-specific.

- `layout` contracts use a fixed preset set owned by sloth.
- `section` contracts expand to full width inside their parent layout.
- `block` contracts declare how much of the layout grid they occupy.
- SEO data should be defined once as a reusable object definition and referenced wherever needed.

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://schemas.sloth.dev/component-contract/1.0.0",
  "title": "Sloth Component Contract",
  "type": "object",
  "required": ["name", "label", "kind", "version", "schemaVersion"],
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
    },
    "layoutPreset": {
      "type": "string",
      "enum": [
        "header",
        "footer",
        "content-12cols",
        "aside-left",
        "aside-right"
      ]
    },
    "gridSpan": {
      "type": "object",
      "required": ["cols"],
      "additionalProperties": false,
      "properties": {
        "cols": { "type": "integer", "minimum": 1, "maximum": 12 },
        "rows": { "type": "integer", "minimum": 1 }
      }
    },
    "breakpointGridSpec": {
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "xs": { "$ref": "#/$defs/gridSpan" },
        "sm": { "$ref": "#/$defs/gridSpan" },
        "md": { "$ref": "#/$defs/gridSpan" },
        "lg": { "$ref": "#/$defs/gridSpan" },
        "xl": { "$ref": "#/$defs/gridSpan" }
      }
    }
  },
  "properties": {
    "name": {
      "type": "string",
      "pattern": "^[a-z0-9]+(?:-[a-z0-9]+)*$"
    },
    "label": {
      "type": "string",
      "minLength": 1
    },
    "kind": {
      "type": "string",
      "enum": ["layout", "section", "block"]
    },
    "componentKind": {
      "type": "string",
      "description": "Business-purpose classifier such as CTA, Carousel, HeroSection"
    },
    "layoutPreset": {
      "$ref": "#/$defs/layoutPreset"
    },
    "version": {
      "type": "string",
      "pattern": "^(0|[1-9]\\d*)\\.(0|[1-9]\\d*)\\.(0|[1-9]\\d*)(?:-[0-9A-Za-z.-]+)?(?:\\+[0-9A-Za-z.-]+)?$"
    },
    "schemaVersion": {
      "type": "string"
    },
    "category": {
      "type": "string"
    },
    "gridSpan": {
      "$ref": "#/$defs/gridSpan"
    },
    "responsiveGrid": {
      "$ref": "#/$defs/breakpointGridSpec"
    },
    "dataset": {
      "type": "array",
      "minItems": 1,
      "items": {
        "type": "object",
        "required": ["key", "label", "type"],
        "additionalProperties": false,
        "properties": {
          "key": {
            "type": "string",
            "pattern": "^[a-zA-Z][a-zA-Z0-9_]*$"
          },
          "label": {
            "type": "string"
          },
          "type": {
            "type": "string",
            "enum": ["string", "number", "option", "dynamic"]
          },
          "required": {
            "type": "boolean"
          },
          "options": {
            "type": "array",
            "items": {
              "type": "object",
              "required": ["label", "value"],
              "additionalProperties": false,
              "properties": {
                "label": { "type": "string" },
                "value": {}
              }
            }
          },
          "value": {},
          "valueDropdown": {
            "type": "object",
            "required": ["contentType", "path"],
            "additionalProperties": false,
            "properties": {
              "contentType": { "type": "string" },
              "path": { "type": "string" },
              "multiple": { "type": "boolean" }
            }
          }
        }
      }
    },
    "renderMeta": {
      "type": "object",
      "required": ["rendererKey"],
      "additionalProperties": false,
      "properties": {
        "rendererKey": { "type": "string" },
        "supportsChildren": { "type": "boolean" }
      }
    },
    "seo": {
      "$ref": "#/$defs/seo"
    }
  }
}
```

Kind-specific notes:

- `layout` contracts must use a fixed preset from `layoutPreset` and may also carry `responsiveGrid` guidance for the builder.
- `section` contracts should be authored as full-width content regions; the layout determines the visible columns.
- `block` contracts should use `gridSpan` to communicate the user-visible size in the builder, with responsive behavior handled by the layout grid and component implementation.
- `seo` should be reused as a shared object definition instead of repeating ad hoc SEO fields in each schema.

## 2) Page Composition Draft

Current Milestone 1 design profile:

- active fields in the page model: `name`, `label`, `pageType`, `route`, `dataset`, `puckConfig`
- deferred fields: `compiledConfig`, `seo`
- removed from the page model for now: persisted `contractRefs`
- page SEO should reuse the shared `seo` object definition from the component schema draft

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
