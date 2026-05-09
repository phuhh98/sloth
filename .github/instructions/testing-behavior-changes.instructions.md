---
description: "Use when changing testable modules or behavior in TypeScript/JavaScript/Go code. Require adding or updating appropriate tests to verify behavior and prevent regressions."
name: "Testing Behavior Changes"
applyTo:
  - "**/*.ts"
  - "**/*.tsx"
  - "**/*.js"
  - "**/*.jsx"
  - "**/*.mjs"
  - "**/*.cjs"
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
- When creating a custom script (for example under a `scripts/` folder), add a corresponding automated test file in the same change.
- Custom script tests should verify the script's branching or argument-handling behavior and be included in the package test runner.
- If a module is currently hard to test, add a brief note in the change explaining constraints and propose a near-term testability improvement.
- Avoid shipping behavior changes with no test evidence unless explicitly requested by the user.
