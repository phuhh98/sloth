---
sidebar_position: 7
---

# Renderer Integration Example

Milestone 3 introduces a runtime mapping utility that resolves `rendererKey` values from page payload nodes to renderer implementations.

## Runtime Utility

Use `createRendererRegistry` from component-hub runtime:

```js
import { createRendererRegistry } from "@sloth/component-hub/runtime/renderer-mapping";
```

## First-Level Payload Strategy

- Page payload includes ordered component nodes.
- Each node uses `rendererKey` plus lightweight props.
- First-level linked content is provided as a lookup map.
- Renderer functions resolve linked entries by id when needed.

Example script: `apps/frontend/runtime-example.mjs`.

Run it locally:

```bash
node apps/frontend/runtime-example.mjs
```

This prints rendered HTML fragments for the seeded page payload.
