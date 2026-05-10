---
purpose: "Execution board for migrating component-hub contracts from version-first source folders to component-first source with deterministic release generation."
status: "superseded"
owner: "platform"
last_updated: "2026-05-10T11:10:00Z"
related_docs:
  - "docs/COMPONENT-HUB-DOCS-INTEGRATION.md"
  - "docs/MILESTONES.md"
  - "docs/REGISTRY.md"
  - "docs/IDEAS.md"
  - "docs/KANBAN-OCI-REGISTRY-PIVOT.md"
---

# Milestone Kanban: Contract Release Model Migration

Superseded: This board is abandoned in favor of OCI registry strategy tracked in `docs/KANBAN-OCI-REGISTRY-PIVOT.md`.

Use this board to control migration from source layout `src/contracts/<version>/components/...` to `src/contracts/components/...` while preserving immutable versioned published artifacts.

## Scope

- Milestone: Contract Release Model Migration (post Milestone 3)
- Goal: Use component-first source contracts and deterministic release snapshots for versioned artifacts.
- Constraints: Keep published `registry/contracts/<version>/...` layout compatible in this migration; max 20 tasks.
- milestone_updated_at: 2026-05-10T11:10:00Z

## Task Decomposition Rules

- Split into small executable tasks.
- Prefer package-local tasks before cross-package integration tasks.
- Minimize tasks that require concurrent edits in different packages.
- Define clear dependency order.
- For tasks requiring user confirmation (frontend review, code style review, docker/browser verification), mark with `requires-confirmation: true`.

## Kanban

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

- [x] REL-005: Update registry build to generate versioned artifacts from release ledger
  - what: Preserve dist compatibility while moving source model.
  - do: Updated build-registry.mjs to materialize releases from ledger entries (preferred path) while maintaining fallback to legacy folder-based copying. Verified deterministic rebuilds produce identical artifacts. Ledger validation integrated into build pipeline.
  - next: Add one-time migration script for existing releases.
  - deps: REL-003, REL-004
  - status: done

- [ ] REL-009: Update release workflow tests for new behavior
  - what: Cover create/verify/materialize and immutability enforcement.
  - do: Rewrite workflow test fixtures for component-first source and release ledger.
  - next: Add regression case for hash drift rejection.
  - deps: REL-004, REL-006
  - status: todo

- [ ] REL-010: Update contract policy tests for new invariants
  - what: Replace folder path assumptions with release ledger and component-first checks.
  - do: Rewrite validation tests and error expectations.
  - next: Align all messages with new policy vocabulary.
  - deps: REL-003, REL-006
  - status: todo

- [ ] REL-011: Update docs and repo instructions for source model change
  - what: Align docs and instruction rules with component-first source and versioned release generation flow.
  - do: Update integration docs, schema docs if needed, and contract authoring instruction.
  - next: Update milestone notes for migration status.
  - deps: REL-005, REL-007
  - status: todo

- [ ] REL-012: Run full verification gate before merge
  - what: Validate migration end-to-end.
  - do: Run component-hub tests, docs registry tests, docs build, and deterministic rebuild checks.
  - next: Prepare rollout summary.
  - deps: REL-009, REL-010, REL-011
  - status: todo

## In Progress

- [ ] REL-009: Update release workflow tests for new behavior
  - what: Cover create/verify/materialize and immutability enforcement.
  - do: Rewrite workflow test fixtures for component-first source and release ledger.
  - next: Add regression case for hash drift rejection.
  - deps: REL-004, REL-006
  - status: in-progress

## Done

- [x] REL-008: Track docs registry artifacts in git
  - what: Commit versioned generated artifacts under docs static registry.
  - do: Updated docs gitignore rules to stop ignoring `apps/docs/static/registry/**` and clarified docs sync workflow assumptions to commit generated artifacts. Verified with `registry:prepare` and git status.
  - next: Update release workflow tests for new behavior.
  - deps: REL-007
  - status: done

- [x] REL-007: Update docs registry sync to preserve immutable historical versions
  - what: Remove destructive sync behavior.
  - do: Updated sync-component-hub-registry script to preserve existing versioned artifacts and only refresh mutable index/state files. Added automated tests proving historical folder preservation and stable revision behavior.
  - next: Track docs registry artifacts in git.
  - deps: REL-005
  - status: done

- [x] REL-005: Update registry build to generate versioned artifacts from release ledger
  - what: Preserve dist compatibility while moving source model.
  - do: Updated build-registry.mjs to materialize releases from ledger entries (preferred path) while maintaining fallback to legacy folder-based copying. Verified deterministic rebuilds produce identical artifacts. Ledger validation integrated into build pipeline.
  - next: Add one-time migration script for existing releases.
  - deps: REL-003, REL-004
  - status: done

- [x] REL-004: Redesign release workflow command surface
  - what: Move workflow commands to `release create|verify|materialize` semantics.
  - do: Added three new functions (createReleaseLedgerEntry, verifyReleaseIntegrity, materializeRelease) supporting deterministic release generation and verification from ledger. Updated command surface to `release create|verify|materialize` while maintaining backward compatibility with old sync/create commands. Added 3 new tests covering all new workflow paths.
  - next: Update registry build to generate versioned artifacts from release ledger.
  - deps: REL-001, REL-002, REL-003
  - status: done

- [x] REL-002: Add component-first source discovery utility
  - what: Centralize scanning of source contracts in `src/contracts/components/*/contract.json`.
  - do: Build reusable helper for release workflow, policy validation, and registry build scripts.
  - next: Integrate helper into policy and workflow scripts.
  - deps: REL-001
  - status: done

- [x] REL-003: Refactor contract version policy to release-ledger model
  - what: Validate releases from release ledger plus source contracts instead of version folders.
  - do: Added validateLedgerRelease() function that checks ledger integrity and component hash alignment; supports backward-compatible dual-location lookup for contracts during migration; extended validateContracts() to validate ledger entries; updated with 7 new tests covering all validation paths.
  - next: Redesign release workflow command surface.
  - deps: REL-001, REL-002
  - status: done

- [x] REL-001: Define release ledger schema and invariants
  - what: Specify immutable release metadata contract.
  - do: Defined fields (`version`, `components`, `contentHash`, `createdAt`, optional `sourceGitRef`, `deprecatedAt`) and implemented parser/validator helper plus baseline ledger file.
  - next: Consume ledger in policy/build/release workflows.
  - deps: none
  - status: done

- [x] REL-006: Skip one-time migration script (no public releases yet)
  - what: Skip legacy migration since 0.0.1 was internal only.
  - do: Removed legacy release-source files and seeded `src/contracts/components/*/contract.json` from existing generated 0.0.1 artifacts so ledger validation remains green without a migration script.
  - next: Update docs registry sync for immutable versions.
  - deps: none
  - status: done (skipped)

## Blocked

- [ ] None
  - what: no active blockers
  - do: continue with dependency order
  - next: log blockers here when encountered
  - deps: none
  - status: blocked

- [x] REL-000: Migration scope and compatibility boundaries approved
  - what: Confirm target state and compatibility constraints.
  - do: Locked decisions for component-first source and version-first published layout.
  - next: Execute REL-001 onward.
  - deps: none
  - status: done

## Dependency Plan

- REL-001 -> REL-002
- REL-001 -> REL-003
- REL-002 -> REL-003
- REL-001 -> REL-004
- REL-002 -> REL-004
- REL-003 -> REL-004
- REL-003 -> REL-005
- REL-004 -> REL-005
- REL-001 -> REL-006
- REL-002 -> REL-006
- REL-005 -> REL-007
- REL-007 -> REL-008
- REL-004 -> REL-009
- REL-006 -> REL-009
- REL-003 -> REL-010
- REL-006 -> REL-010
- REL-005 -> REL-011
- REL-007 -> REL-011
- REL-009 -> REL-012
- REL-010 -> REL-012
- REL-011 -> REL-012

## Notes

- Risks:
  - Transition period may require compatibility shims while existing release folders still exist.
  - Deterministic materialization for older releases depends on release metadata completeness.
- Decisions:
  - Published artifact layout remains version-first for this migration.
  - Release generation is explicit and command-driven.
  - Release payloads remain immutable once published.
- Next:
  - Implement REL-001 and REL-002 foundation in component-hub scripts and tests.
