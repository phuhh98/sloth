---
description: "Use when changing testable modules or behavior in TypeScript/JavaScript/Go code. Require adding or updating appropriate tests to verify behavior and prevent regressions."
name: "Testing Behavior Changes"
applyTo:
  - "**/*.ts"
  - "**/*.tsx"
  - "**/*.js"
  - "**/*.jsx"
  - "**/*.go"
---

# Testing Behavior Changes

- For any behavior change in a testable module, add or update appropriate automated tests in the same change.
- Treat tests as part of the definition of done for behavior changes, not optional follow-up work.
- Choose the narrowest test level that proves behavior clearly:
  - unit tests for pure logic and local branching
  - integration tests for module boundaries and framework wiring
- Cover both expected behavior and key edge cases that could regress.
- If behavior changes intentionally, update existing tests to match the new contract and remove outdated assertions.
- If a module is currently hard to test, add a brief note in the change explaining constraints and propose a near-term testability improvement.
- Avoid shipping behavior changes with no test evidence unless explicitly requested by the user.
