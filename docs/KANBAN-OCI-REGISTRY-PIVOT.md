---
purpose: "Execution board for pivoting contract distribution from release-manifest migration to OCI artifacts in GHCR with CLI abstraction."
status: "active"
owner: "platform"
last_updated: "2026-05-10T11:35:00Z"
related_docs:
  - "docs/REGISTRY.md"
  - "docs/IDEAS.md"
  - "docs/COMPONENT-CONTRACTS.md"
  - "docs/MILESTONES.md"
  - "docs/archive/KANBAN-CONTRACT-RELEASE-MIGRATION.md"
---

# Milestone Kanban: OCI Registry Pivot

Use this board to execute the registry strategy pivot to GHCR OCI artifacts with ORAS in CLI while keeping user-facing contract commands simple (`contracts ls`, `contracts pull`).

## Scope

- Milestone: OCI Registry Pivot (post Milestone 3 adjustment)
- Goal: Replace manifest-ledger-based distribution assumptions with OCI-based contract distribution in GHCR.
- Constraints: Keep user command surface contract-oriented; avoid requiring users to invoke raw registry commands.
- milestone_updated_at: 2026-05-10T11:35:00Z

## Carryover Mapping (From Superseded Milestone 3 Work)

- REL-009 (release workflow tests): kept as OCI-007 and OCI-011 with OCI-focused integration and validation coverage.
- REL-010 (contract policy tests): kept as OCI-011 by re-baselining test/policy expectations around OCI resolver behavior.
- REL-011 (docs/instructions updates): kept as OCI-012 for final user-facing and instruction alignment.
- REL-012 (full verification gate): kept as OCI-010.
- Release-ledger-specific migration assertions: dropped as no longer applicable to OCI strategy.

## Task Decomposition Rules

- Split into small executable tasks.
- Prefer package-local tasks before cross-package integration tasks.
- Minimize tasks requiring simultaneous changes across packages.
- Keep task count <= 20.

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

- [ ] OCI-002: Implement ORAS-backed registry client in CLI
  - what: Add GHCR OCI list/pull primitives in Go using ORAS SDK.
  - do: Add pkg/registry client abstraction and GHCR implementation with digest-safe pull.
  - next: Integrate with contracts resolver.
  - deps: OCI-001
  - status: todo

- [ ] OCI-003: Replace manifest-based contracts resolver with OCI resolver
  - what: Stop depending on docs static manifest/index for contract listing.
  - do: Refactor `pkg/source` + app wiring to consume OCI metadata and payload files.
  - next: Update command integration tests.
  - deps: OCI-002
  - status: todo

- [ ] OCI-004: Add `contracts ls` command alias and version selector
  - what: Normalize list command naming and include `--version <x.y.z|latest>` semantics.
  - do: Update command docs/output and keep backward-compatible alias behavior.
  - next: Update docs and examples.
  - deps: OCI-003
  - status: todo

- [ ] OCI-005: Add single-contract pull command
  - what: Support pulling one contract by name from selected release version.
  - do: Implement `sloth contracts pull --name --version [--out]` and file write behavior.
  - next: Add negative tests for missing contract/version.
  - deps: OCI-003
  - status: todo

- [ ] OCI-006: Add local Zot compose service for integration tests
  - what: Provide on-demand local OCI registry for CLI integration tests.
  - do: Add dedicated compose file/service with profile so default compose remains unchanged.
  - next: Add task commands to start/stop registry test service.
  - deps: none
  - status: todo

- [ ] OCI-007: Add CLI integration tests against Zot
  - what: Validate list/pull behavior against real OCI registry implementation.
  - do: Add test setup/teardown that uses Zot endpoint and seeded contract artifact.
  - next: Gate in CI-friendly mode.
  - deps: OCI-005, OCI-006
  - status: todo

- [ ] OCI-008: Add GitHub Actions workflow for GHCR contract artifact publish
  - what: Publish raw contract folder artifact to GHCR on CLI release.
  - do: Add workflow with `packages: write` permission and release-tag-derived version.
  - next: Add dry-run/manual dispatch path.
  - deps: OCI-002
  - status: todo

- [ ] OCI-009: Deprecate obsolete manifest-ledger migration checks
  - what: Remove or downgrade checks that assume release-manifest migration is active.
  - do: Update lint-staged/hooks/CI steps and docs references; archive or remove obsolete migration scripts/tests that are no longer part of release criteria.
  - next: Re-baseline tests and validation to OCI path.
  - deps: OCI-003, OCI-008
  - status: todo

- [ ] OCI-011: Re-baseline contract validation and test suites for OCI path
  - what: Convert remaining migration-era test expectations into OCI-era behavior checks.
  - do: Update component-hub and CLI tests to remove release-ledger coupling where obsolete; add coverage for version selection, missing contract handling, and resolver fallback behavior.
  - next: Ensure pre-commit and CI checks align with updated test contracts.
  - deps: OCI-004, OCI-005, OCI-009
  - status: todo

- [ ] OCI-012: Finalize docs and repo instructions for OCI contract flow
  - what: Align user docs and repository instructions with `contracts ls/pull` + OCI backend model.
  - do: Update CLI docs pages, docs README references, and instruction files that still describe manifest-ledger migration behavior.
  - next: Include in final verification gate and release notes.
  - deps: OCI-004, OCI-005
  - status: todo

- [ ] OCI-010: Run full verification gate for pivot
  - what: Validate pivot end-to-end across CLI tests, docs build, and workflow checks.
  - do: Execute test matrix and summarize rollout constraints.
  - next: Prepare merge summary.
  - deps: OCI-007, OCI-008, OCI-009, OCI-011, OCI-012
  - status: todo

## In Progress

- [ ] OCI-001: Add OCI registry config model to CLI
  - what: Add config fields for registry host/repository and auth usage in CLI.
  - do: Extend config structs/env resolution and update tests for precedence.
  - next: Wire resolver selection into contracts command path.
  - deps: none
  - status: in-progress

## Blocked

- [ ] None
  - what: no active blockers
  - do: continue with dependency order
  - next: log blockers here when encountered
  - deps: none
  - status: blocked

## Done

- [x] DOC-OCI-000: Update strategy docs before execution planning
  - what: Align architecture docs with OCI pivot and contracts abstraction.
  - do: Updated registry/contracts/milestones docs and marked old migration board superseded before creating this board.
  - next: implement OCI-001.
  - deps: none
  - status: done
