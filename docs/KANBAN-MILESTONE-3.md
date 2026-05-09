---
purpose: "Milestone-level implementation Kanban tracking for component hub starter packs and renderer examples"
status: "active"
owner: "platform"
last_updated: "2026-05-10"
related_docs:
  - "docs/IMPLEMENTATION-PLAN.md"
  - "docs/MILESTONES.md"
  - "docs/IDEAS.md"
  - "docs/REGISTRY.md"
---

# Milestone Kanban: Milestone 3 - Component Hub Starter + Renderer Examples

Use this board for detailed execution tracking of one milestone.

### Scope

- Milestone: Milestone 3 - Component Hub Starter + Renderer Examples
- Goal: Publish starter packs and runtime integration examples for rendering page payloads with component mappings
- Constraints: Keep milestone under 20 tasks, prioritize package-local work in component-hub before cross-package wiring
- milestone_updated_at: 2026-05-10T23:20:00Z

### Task Decomposition Rules

- Split into small executable tasks.
- Prefer package-local tasks before cross-package integration tasks.
- Minimize tasks that require concurrent edits in different packages.
- Define clear dependency order.
- For tasks requiring user confirmation (frontend review, code style review, docker/browser verification), mark with `requires-confirmation: true`.
  - Agent will use multi-choice selection to get user input before proceeding.

### Kanban

Task card format (keep concise):

```text
- [ ] <Task title>
  - what: <what this task is>
  - do: <what to do now>
  - next: <what to do next>
  - deps: <dependency or "none">
  - requires-confirmation: <true|false> (optional, default false)
  - status: <todo|in-progress|blocked|done>
```

## To Do

- [ ] HUB-001: Define starter pack manifest schema
  - what: Create starter pack manifest contract for component hub artifacts.
  - do: Add JSON schema and TS types for pack metadata and component references.
  - next: Implement validator and build script.
  - deps: none
  - requires-confirmation: false
  - status: todo

- [ ] HUB-002: Implement starter pack validator and builder
  - what: Build and validate starter pack artifacts from source descriptors.
  - do: Add scripts and tests in packages/component-hub.
  - next: Generate first pack artifact.
  - deps: HUB-001
  - requires-confirmation: false
  - status: todo

- [ ] HUB-003: Publish first starter pack artifact
  - what: Create one high-quality starter pack with docs metadata.
  - do: Add pack descriptor and generated artifact under docs static registry path.
  - next: Add renderer mapping example.
  - deps: HUB-002
  - requires-confirmation: false
  - status: todo

- [ ] HUB-004: Add runtime renderer mapping utility
  - what: Implement utility that maps component contract rendererKey to React implementation.
  - do: Add package-level API and tests for mapping behavior.
  - next: Integrate with sample app page payload.
  - deps: HUB-003
  - requires-confirmation: false
  - status: todo

- [ ] HUB-005: Build frontend runtime example page
  - what: Demonstrate rendering a page payload using starter pack component mappings.
  - do: Add example in apps/frontend with first-level payload handling.
  - next: Document usage in public docs.
  - deps: HUB-004
  - requires-confirmation: true
  - status: todo

- [ ] HUB-006: Add docs sync workflow for contracts and packs
  - what: Ensure docs static registry receives generated contract/pack artifacts per release.
  - do: Add script and tests for sync behavior.
  - next: Wire task commands and CI hooks.
  - deps: HUB-003
  - requires-confirmation: false
  - status: todo

- [ ] HUB-007: Document component hub and renderer integration
  - what: Publish public docs for starter packs and renderer integration patterns.
  - do: Add docs pages in apps/docs/docs with examples and versioning guidance.
  - next: Mark milestone complete after validation.
  - deps: HUB-005, HUB-006
  - requires-confirmation: false
  - status: todo

## In Progress

- [ ] None
  - what: no active task
  - do: start HUB-001 when milestone 3 execution begins
  - next: scaffold schema and type definitions
  - deps: none
  - status: in-progress

## Blocked

- [ ] None
  - what: no blocked task
  - do: continue planned sequence
  - next: raise blockers here when encountered
  - deps: none
  - status: blocked

## Done

- [x] None
  - what: no completed task yet
  - do: none
  - next: none
  - deps: none
  - status: done

### Dependency Plan

- HUB-001 -> HUB-002
- HUB-002 -> HUB-003
- HUB-003 -> HUB-004
- HUB-003 -> HUB-006
- HUB-004 -> HUB-005
- HUB-005 -> HUB-007
- HUB-006 -> HUB-007

### Notes

- Risks:
  - Runtime payload shape may evolve while plugin delivery contracts are still maturing.
- Decisions:
  - Start with one starter pack before broadening pack catalog.
- Next:
  - Start HUB-001.

### Archival

When this milestone is complete, move this file to `docs/archive/` and create a fresh board in `docs/` for the next milestone.
