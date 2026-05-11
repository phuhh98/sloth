---
purpose: "Track CLI readability and refactoring work for the sloth package."
status: "done"
owner: "platform-and-cli"
last_updated: "2026-05-11"
related_docs:
  - "packages/cli/README.md"
  - "packages/cli/internal/app/root.go"
  - "packages/cli/internal/app/contracts_*.go"
  - "packages/cli/pkg/config/config.go"
---

# Milestone Kanban: CLI Readability Refactor

## Scope

- Goal: reduce command-level orchestration noise and make the CLI easier for humans to read and maintain
- Constraints: behavior must remain unchanged; focus on structure, naming, and orchestration boundaries
- milestone_updated_at: 2026-05-11

## Kanban

### To Do

- [ ] None

### In Progress

- [ ] None

### Blocked

- [ ] None

### Done

- [x] CLI architecture review
  - what: identified readability and refactor opportunities in the CLI package
  - deps: none
  - status: done

- [x] CLI-001 Introduce shared runtime context
  - what: resolve config and resolver once per command run
  - do: extracted Runtime struct and BuildRuntime() in runtime.go; commands now call opts.BuildRuntime() instead of repeating setup
  - deps: none
  - status: done

- [x] CLI-002 Thin out contract command handlers
  - what: move orchestration out of Cobra RunE functions
  - do:
    - Extracted applyEnvToProfile / applyProfileDefaults from ResolveConfig (root.go)
    - Extracted buildVerifyInput / printVerifyResult from contracts_verify_cmd.go
    - Extracted executePush / ingestWithRetry / syncLockAfterPush / printPushDryRun from contracts_push_cmd.go
    - Extracted runAddAll / runAddComponent / runAddSet from contracts_add_cmd.go
    - Removed unused config/host imports from inspect and push cmd files
    - All cognitive complexity violations resolved; all tests passing
  - deps: CLI-001
  - status: done

## Dependency Plan

- CLI-001 -> CLI-002

## Notes

- Risks:
  - refactor can accidentally change behavior if config resolution or source selection is duplicated inconsistently.
  - too much abstraction would make the CLI harder to read instead of easier.
- Decisions:
  - kept the first slice small and behavior-preserving.
  - preferred a shared runtime context over deeper architectural changes.

## Archival

Move this file to `docs/archive/` when no longer needed as a reference.
