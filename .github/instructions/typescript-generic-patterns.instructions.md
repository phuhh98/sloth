---
description: "Use when writing or reviewing TypeScript in this project. Prefer generics for reusable abstractions and apply suitable design patterns when they improve maintainability, type safety, and extensibility."
name: "TypeScript Generic and Pattern Guidance"
applyTo:
  - "**/*.ts"
  - "**/*.tsx"
---

# TypeScript Generic and Pattern Guidance

- Prefer generics over `any` when behavior is reusable across types.
- Model domain contracts with precise interfaces and type aliases; avoid weakly typed object maps unless justified.
- Use discriminated unions for variant models (for example, static vs dynamic page types) instead of ad hoc runtime branching.
- Introduce design patterns only where they reduce coupling or complexity:
  - strategy for interchangeable behaviors
  - factory for controlled object/config creation
  - adapter for external API or format boundaries
- Keep function signatures explicit and narrow; favor small composable functions.
- Preserve runtime validation at boundaries even with strong static types.
- Avoid over-engineering: apply patterns when there is clear variability, extension pressure, or repeated branching.
- When adding new abstractions, include a short rationale in code comments or PR notes explaining why generics or a pattern is warranted.
