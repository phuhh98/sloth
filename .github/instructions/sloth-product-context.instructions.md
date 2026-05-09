---
description: "Use when planning or implementing sloth plugin, admin builder, CLI sync, component hub, registry, content-types, routes, and related architecture decisions. Consult relevant docs files (not only docs/IDEAS.md), maintain docs/MILESTONES.md, and capture new ideas into docs with markdown metadata."
name: "Sloth Product Context"
---

# Sloth Product Context Instruction

- Treat docs/IDEAS.md as the primary product-direction source for sloth.
- Before proposing architecture or implementation details, consult docs/IDEAS.md and any other relevant files under docs/ (for example docs/REGISTRY.md and docs/MILESTONES.md).
- Use each docs file markdown metadata to decide relevance quickly.
- Require markdown metadata frontmatter for docs files that guide product/architecture decisions. Recommended fields:
  - purpose
  - status
  - owner
  - last_updated
  - related_docs
- If a request conflicts with docs/IDEAS.md, call out the conflict clearly and suggest options:
  - follow current IDEAS direction
  - update IDEAS first, then implement
- Keep naming consistent with IDEAS domain language: component config, page template, puckConfig, compiledConfig, component hub, registry.
- Prefer incremental delivery aligned with roadmap phases (plugin/editor and CLI first; registry paid features later).
- For Strapi work, favor Strapi v5 patterns and Document Service API usage.
- Maintain docs/MILESTONES.md as a status tracker of completed/in-progress/not-started work mapped from docs/IDEAS.md roadmap items.
- When new ideas appear during implementation or discussion, proactively suggest saving them:
  - append broad product ideas to docs/IDEAS.md
  - split specialized ideas into focused docs files under docs/
  - every new docs file must include markdown metadata frontmatter describing its general purpose and scope
- When introducing new behavior, suggest which docs files should be updated to keep decisions synchronized.
