---
purpose: "Initial draft JSON artifacts for sloth component/page content-types and component contract schema."
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

- `Page` does not require persisted `contractRefs` for Milestone 1.
- Contract linkage can be derived at runtime from `puckConfig` and component records.
- Keep page model minimal now: page identity, route, dataset, and editor/runtime config.
- Revisit explicit `contractRefs` only if runtime profiling shows a concrete query/performance bottleneck.

## 1) Strapi Content-Type Draft: Component

Suggested target file in plugin:

- `server/src/content-types/component/schema.json`

```json
{
  "kind": "collectionType",
  "collectionName": "sloth_components",
  "info": {
    "singularName": "component",
    "pluralName": "components",
    "displayName": "Component",
    "description": "Business-purpose component entity materialized from contract ingest"
  },
  "options": {
    "draftAndPublish": true
  },
  "pluginOptions": {},
  "attributes": {
    "name": {
      "type": "string",
      "required": true,
      "unique": true
    },
    "label": {
      "type": "string",
      "required": true
    },
    "componentKind": {
      "type": "enumeration",
      "enum": [
        "CTA",
        "Carousel",
        "HeroSection",
        "AsideLayout",
        "Header",
        "Footer",
        "Custom"
      ],
      "required": true
    },
    "contractName": {
      "type": "string",
      "required": true
    },
    "contractVersion": {
      "type": "string",
      "required": true
    },
    "schemaVersion": {
      "type": "string",
      "required": true
    },
    "contractHash": {
      "type": "string"
    },
    "description": {
      "type": "text"
    },
    "metadata": {
      "type": "json"
    },
    "contractPayload": {
      "type": "json",
      "required": true
    }
  }
}
```

## 2) Strapi Content-Type Draft: Page

Current Milestone 1 implementation profile:

- active fields in plugin: `name`, `label`, `pageType`, `route`, `dataset`, `puckConfig`
- deferred fields: `compiledConfig`, `seo`
- removed from current draft model: persisted `contractRefs`

Suggested target file in plugin:

- `server/src/content-types/page/schema.json`

```json
{
  "kind": "collectionType",
  "collectionName": "sloth_pages",
  "info": {
    "singularName": "page",
    "pluralName": "pages",
    "displayName": "Page",
    "description": "Page configuration and first-level linked data mapping"
  },
  "options": {
    "draftAndPublish": true
  },
  "pluginOptions": {},
  "attributes": {
    "name": {
      "type": "string",
      "required": true,
      "unique": true
    },
    "label": {
      "type": "string",
      "required": true
    },
    "pageType": {
      "type": "enumeration",
      "enum": ["static", "dynamic"],
      "required": true
    },
    "route": {
      "type": "string",
      "required": true,
      "unique": true
    },
    "dataset": {
      "type": "json",
      "required": true
    },
    "puckConfig": {
      "type": "json",
      "required": true
    },
    "compiledConfig": {
      "type": "json"
    },
    "seo": {
      "type": "json"
    }
  }
}
```

## 3) Initial `$schema` Draft: Component Contract

Suggested target file for published schema bundle:

- `contracts/schema/component-contract.schema.json`

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://schemas.sloth.dev/component-contract/1.0.0",
  "title": "Sloth Component Contract",
  "type": "object",
  "required": ["name", "label", "kind", "version", "schemaVersion", "dataset"],
  "additionalProperties": false,
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
    }
  }
}
```
