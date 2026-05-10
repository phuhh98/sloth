---
purpose: "Milestone 1 board for contract schema realignment and docs synchronization."
status: "active"
owner: "platform-and-cli"
last_updated: "2026-05-10"
related_docs:
  - "docs/SCHEMA-DRAFTS.md"
  - "docs/COMPONENT-CONTRACTS.md"
  - "docs/IDEAS.md"
  - "docs/MILESTONES.md"
---

# Milestone Kanban: Contract Schema Realignment

## Scope

- Milestone: Milestone 1 contract schema realignment
- Goal: revise the schema draft to match the intended contract model and align docs with the CLI pull -> verify -> push -> plugin lifecycle flow
- Constraints: docs/spec work only; no plugin implementation in this cycle
- milestone_updated_at: 2026-05-10

## Task Decomposition Rules

- Split schema semantics, docs alignment, and approval into small executable tasks.
- Keep any future implementation work blocked until the revised schema draft is approved.
- Prefer doc-local edits before cross-doc consistency passes.

## Kanban

### To Do

- [ ] DOC-001 Revise component contract schema draft
  - what: rewrite the schema draft to use kind-specific rules
  - do: add the fixed layout preset set, block grid-span semantics, and kind-driven shape notes
  - next: review the revised JSON examples for clarity
  - deps: none
  - requires-confirmation: false
  - status: todo

- [ ] DOC-002 Add reusable SEO schema definition
  - what: define SEO as a shared schema object
  - do: update the draft to use one reusable SEO object for page and contract schemas
  - next: cross-link the shared SEO shape in docs
  - deps: DOC-001
  - requires-confirmation: false
  - status: todo

- [ ] DOC-003 Clarify layout and section semantics
  - what: document how layout, section, and block differ
  - do: write the rules for layout presets, section full-width behavior, and block span ownership
  - next: validate the builder-facing wording
  - deps: DOC-001
  - requires-confirmation: false
  - status: todo

- [ ] DOC-004 Align CLI and host boundary language
  - what: document pull -> verify -> push and ingest ownership
  - do: update the contracts policy doc to match the simulation OpenAPI flow
  - next: ensure the ingest endpoint and lifecycle wording are consistent
  - deps: none
  - requires-confirmation: false
  - status: todo

- [ ] DOC-005 Update terminology in IDEAS
  - what: align the product model with the revised schema
  - do: replace ambiguous kind/type language and add layoutPreset, gridSpan, and SEO reuse notes
  - next: verify the TypeScript examples still read cleanly
  - deps: DOC-001
  - requires-confirmation: false
  - status: todo

- [ ] DOC-006 Update milestone status notes
  - what: reflect the spec-first correction track in milestones
  - do: replace outdated plugin-implementation notes with the revised scope boundary
  - next: keep milestone notes consistent with the new board
  - deps: DOC-001
  - requires-confirmation: false
  - status: todo

- [ ] DOC-007 Reconcile component-hub base contract notes
  - what: keep the shared base contract aligned with the revised schema
  - do: add the shared SEO and layout semantics to the shortlist notes
  - next: confirm the hub doc does not drift from schema draft language
  - deps: DOC-001
  - requires-confirmation: false
  - status: todo

- [ ] DOC-008 Final cross-doc consistency pass
  - what: check all updated docs for terminology and scope consistency
  - do: review links, names, and the no-implementation boundary across docs
  - next: approve the spec set for the next implementation milestone
  - deps: DOC-002, DOC-004, DOC-005, DOC-006, DOC-007
  - requires-confirmation: true
  - status: todo

### In Progress

- [ ] None

### Blocked

- [ ] IMPL-001 Plugin lifecycle implementation
  - what: create or update plugin ingest behavior based on the approved schema
  - do: wait for schema and docs approval before touching plugin code
  - next: implement lifecycle materialization after approval
  - deps: DOC-008
  - requires-confirmation: false
  - status: blocked

### Done

- [x] Initial architecture review and scope reset
  - what: reviewed current contracts state against the revised target model
  - do: captured the corrected architecture direction and scope limits
  - next: none
  - deps: none
  - status: done

## Dependency Plan

- DOC-001 -> DOC-002
- DOC-001 -> DOC-003
- DOC-001 -> DOC-005
- DOC-001 -> DOC-006
- DOC-001 -> DOC-007
- DOC-002, DOC-004, DOC-005, DOC-006, DOC-007 -> DOC-008
- DOC-008 -> IMPL-001

## Notes

- Risks:
  - contract schema semantics can drift from builder expectations if the fixed layout presets are not called out explicitly.
  - reusable SEO fields can fragment if page and component docs do not point at the same shared object shape.
  - plugin implementation may be started too early without the spec approval gate.
- Decisions:
  - keep this cycle spec-first and docs-only.
  - keep the fixed layout preset set minimal and sloth-owned.
  - use a reusable SEO object definition rather than repeating SEO fields in each schema section.
- Next:
  - finish the docs updates, then approve the board for the next implementation cycle.

## Archival

When this milestone is complete, move this file to `docs/archive/` and create a fresh board in `docs/` for the next milestone.
