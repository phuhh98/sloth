---
purpose: "Execution board for pivoting contract distribution from release-manifest migration to OCI artifacts in GHCR with CLI abstraction."
status: "active"
owner: "platform"
last_updated: "2026-05-10T22:05:00Z"
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
- milestone_updated_at: 2026-05-10T22:05:00Z

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

- [ ] None
  - what: no remaining open OCI pivot tasks
  - do: prepare merge summary and follow-up backlog items outside this board
  - next: archive board when merge is complete
  - deps: none
  - status: todo

## In Progress

- [ ] None
  - what: no active in-progress tasks
  - do: maintain done state and prepare archival/rollout notes
  - next: archive board after merge
  - deps: none
  - status: todo

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

- [x] OCI-001: Add OCI registry config model to CLI
  - what: Add config fields for registry host/repository and auth usage in CLI.
  - do: Added registry host/repository/useAuthorizationToken fields with env support and defaults, plus precedence tests in CLI config resolution.
  - next: Implement ORAS-backed registry client and hook resolver selection.
  - deps: none
  - status: done

- [x] OCI-002: Implement ORAS-backed registry client in CLI
  - what: Add GHCR OCI list/pull primitives in Go using ORAS SDK.
  - do: Added `pkg/registry` ORAS client for version listing and digest-safe release payload pull, then wired `--source oci` to an OCI resolver in contracts command path.
  - next: Complete resolver migration coverage with OCI-focused command integration tests.
  - deps: OCI-001
  - status: done

- [x] OCI-003: Replace manifest-based contracts resolver with OCI resolver
  - what: Stop depending on docs static manifest/index for contract listing.
  - do: Added OCI source integration tests for list/add command paths and hardened localhost OCI behavior for stable resolver execution in tests and local registry flows.
  - next: Add command UX polish (`ls`/`--version`) and single-contract pull command.
  - deps: OCI-002
  - status: done

- [x] OCI-004: Add `contracts ls` command alias and version selector
  - what: Normalize list command naming and include `--version <x.y.z|latest>` semantics.
  - do: Added `contracts ls` alias, introduced `--version` with backward-compatible `--plugin-version`, and updated CLI docs/examples.
  - next: Implement `contracts pull` command.
  - deps: OCI-003
  - status: done

- [x] OCI-005: Add single-contract pull command
  - what: Support pulling one contract by name from selected release version.
  - do: Added `sloth contracts pull --name --version [--out]` with workspace write + lock update behavior and integration coverage for success, missing contract, and missing version paths.
  - next: Validate end-to-end flow against real Zot registry in OCI-007.
  - deps: OCI-003
  - status: done

- [x] OCI-006: Add local Zot compose service for integration tests
  - what: Provide on-demand local OCI registry for CLI integration tests.
  - do: Added dedicated `docker-compose.oci-registry.yaml`, Zot config at `packages/cli/testdata/zot/config.json`, and Taskfile commands `oci-registry-up`/`oci-registry-down`.
  - next: Seed test artifacts and add real-registry integration tests in OCI-007.
  - deps: none
  - status: done

- [x] OCI-008: Add GitHub Actions workflow for GHCR contract artifact publish
  - what: Publish contract release and schema artifacts to GHCR on release.
  - do: Added workflow `.github/workflows/ghcr-contract-artifacts.yml` with `packages: write`, release/manual triggers, tag-derived version resolution, dry-run support, ORAS GHCR publish for contracts+schema, and digest outputs in workflow summary.
  - next: Switch schema artifact source from docs static path to `packages/contracts` source in OCI-013.
  - deps: OCI-002
  - status: done

- [x] OCI-007: Add CLI integration tests against Zot
  - what: Validate list/pull behavior against real OCI registry implementation.
  - do: Added env-gated real Zot integration test (`TestContractsListAndPullWithZot`) with compose setup/teardown and seeded OCI artifact push via ORAS Go SDK; added `task cli-test-zot-integration`; wired CI workflow step in docs-pages pipeline.
  - next: Expand CI matrix if multi-version registry fixtures are needed.
  - deps: OCI-005, OCI-006
  - status: done

- [x] OCI-013: Move schema source-of-truth to contracts package with GHCR publication
  - what: Keep schema authoring beside contracts in `packages/contracts` and publish immutable schema artifacts to GHCR.
  - do: Added schema source under `packages/contracts/src/schemas/component-contract/<version>/schema.json`; added contracts script to sync/check docs-hosted schema parity; updated docs prebuild sync and GHCR publish workflow to read schema artifact from contracts source path.
  - next: Continue with OCI-009 migration-era check cleanup.
  - deps: OCI-008
  - status: done

- [x] OCI-009: Deprecate obsolete manifest-ledger migration checks
  - what: Remove or downgrade checks that assume release-manifest migration is active.
  - do: Replaced migration-era `validate:contracts --compare-ref` CI gate in docs-pages workflow with OCI-era contracts registry build gate; replaced pre-commit lint-staged contracts gate to run registry build + tests instead of migration immutability compare checks.
  - next: Complete OCI-011 re-baseline for any remaining migration-era test expectations.
  - deps: OCI-003, OCI-008
  - status: done

- [x] OCI-011: Re-baseline contract validation and test suites for OCI path
  - what: Convert remaining migration-era test expectations into OCI-era behavior checks.
  - do: Updated CLI integration tests to prefer `contracts ls` and `--version`; added OCI resolver unit coverage for version pass-through, fallback behavior, missing contract handling, and upstream client error propagation.
  - next: Finish OCI-012 docs/instructions sweep and run final verification gate in OCI-010.
  - deps: OCI-004, OCI-005, OCI-009
  - status: done

- [x] OCI-012: Finalize docs and repo instructions for OCI contract flow
  - what: Align user docs and repository instructions with `contracts ls/pull` + OCI backend model.
  - do: Completed active docs/instructions sweep across root `docs/`, public docs pages, and schema URL docs; updated command examples to `--version` semantics and corrected ownership wording to contracts package where applicable.
  - next: Execute OCI-010 full verification gate and prepare final merge summary.
  - deps: OCI-004, OCI-005, OCI-013
  - status: done

- [x] OCI-010: Run full verification gate for pivot
  - what: Validate pivot end-to-end across CLI tests, docs build, and workflow checks.
  - do: Executed full matrix: CLI Go tests, contracts tests, docs registry tests, docs production build, schema parity check, and real Zot integration test (`task cli-test-zot-integration`) with successful list/pull behavior.
  - next: prepare merge summary and archive OCI pivot board post-merge.
  - deps: OCI-007, OCI-008, OCI-009, OCI-011, OCI-012, OCI-013
  - status: done

- [x] OCI-014: Drop unused scripts and update call sites
  - what: Remove obsolete scripts and update package.json/Taskfile call sites to only use active script paths.
  - do: Removed unused contracts package script aliases with no active call-sites and dropped obsolete `start-frontend-dev` Taskfile task pointing to a non-existent app path.
  - next: keep script/task inventory minimal as new commands are added.
  - deps: none
  - status: done

- [x] OCI-015: Revise docs IA and architecture diagram
  - what: Refresh Docusaurus docs structure/naming and architecture docs to reflect latest OCI/contracts layout.
  - do: Updated public/internal architecture docs for contracts package + GHCR model, improved docs navigation naming (Maintainers), refreshed intro status, and removed orphaned top-level `apps/docs/docs/schemas.md` duplicate to keep audience folders clean.
  - next: continue docs polish with future milestone changes.
  - deps: OCI-014
  - status: done

- [x] OCI-016: Add versioned CLI config schema and publish docs URL
  - what: Provide a versioned schema for `.sloth/config.yaml` in contracts and host it via docs static schemas like component-contract schema.
  - do: Added `cli-config/0.0.1/schema.json` source in contracts, expanded schema sync script/tests for multi-artifact sync, synced docs static schema output, and documented canonical URL usage in consumer docs.
  - next: decide whether to publish CLI config schema as a GHCR OCI schema artifact in a future task.
  - deps: OCI-015
  - status: done
