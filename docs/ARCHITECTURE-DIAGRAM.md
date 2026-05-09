---
purpose: "High-level architecture design diagram for sloth based on product direction in docs/IDEAS.md."
status: "draft"
owner: "product-and-architecture"
last_updated: "2026-05-10"
related_docs:
  - "docs/IDEAS.md"
  - "docs/MILESTONES.md"
  - "docs/COMPONENT-CONTRACTS.md"
  - "docs/REGISTRY.md"
---

# sloth Architecture Design Diagram

Date: 2026-05-10
Source of truth: docs/IDEAS.md

## High-Level System Architecture

```mermaid
flowchart LR
  %% Actors
  A[Admin User]
  D[Frontend Runtime Consumer]
  E[Developer]

  %% Local developer workspace
  subgraph L[Local Project]
    C[sloth CLI\nGo + Cobra]
    F[.sloth config and contracts\nconfig.yaml, contracts, sets, lock.json]
  end

  %% Host system
  subgraph H[Strapi Host with sloth plugin]
    B[Admin UI Builder\nPuck + palette + dataset mapping]
    R1[Admin API\ncomponents/pages/compile]
    R2[Content API\ninspection + ingest + page delivery]
    S[Plugin Services\ningest, inspection, compiler, sync]
    T1[plugin::sloth.component]
    T2[plugin::sloth.page]
    G[(Strapi Document Service API)]
  end

  %% Distribution ecosystem
  subgraph X[Distribution Ecosystem]
    K[Contract Source\nnpm or git]
    U[Component Hub\nthemes and variants]
    V[Registry API and Artifacts\nfuture]
  end

  %% Admin flow
  A --> B
  B --> R1
  R1 --> S
  S --> G
  G --> T1
  G --> T2

  %% Runtime flow
  D --> R2
  R2 --> S
  S --> G
  R2 --> D

  %% CLI flow
  E --> C
  C --> F
  C --> K
  C --> R2

  %% Verification ownership rule from IDEAS
  C -. verifies schema, compatibility, and drift before push .-> R2
  R2 -. ingests verified payloads only .-> S

  %% Future ecosystem links
  C -. list/add/update/search .-> V
  U -. publishes packs to .-> V
```

## Responsibility Boundaries

- CLI owns verification workflow before push.
- Host plugin owns ingest and materialization into component records.
- Runtime delivery endpoint serves page delivery payload and first-level linked content strategy.
- Registry and component hub are later roadmap phases and remain decoupled from core plugin and CLI MVP.

## Architecture Notes

- Keep architecture as a modular monolith around Strapi plugin and CLI during Milestones 1 and 2.
- Add registry complexity incrementally after stable plugin and CLI contracts are proven.
- Keep runtime API generic and avoid deep linked-data parsing in plugin runtime.
