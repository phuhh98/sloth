---
purpose: "Milestone-level implementation Kanban tracking for CLI contract management"
status: "active"
owner: "platform"
last_updated: "2026-05-10"
related_docs:
  - "docs/IMPLEMENTATION-PLAN.md"
  - "docs/MILESTONES.md"
  - "docs/IDEAS.md"
  - "docs/COMPONENT-CONTRACTS.md"
---

# Milestone Kanban: Milestone 2 - CLI Contract Management

Use this board for detailed execution tracking of one milestone.

### Scope

- Milestone: Milestone 2 - CLI Contract Management
- Goal: Build Go + Cobra CLI contract management commands (list/add/verify/push) with local .sloth workspace and robust compatibility checks
- Constraints: Scope is component contracts only; exclude page content operations and keep milestone under 20 tasks
- milestone_updated_at: 2026-05-10

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

- [ ] CLI-001: Set up CLI package structure (Go + Cobra scaffold)
  - what: Initialize packages/cli with Go module and Cobra command framework.
  - do: Create cmd/, pkg/, internal/ and add root command entrypoint.
  - next: Start profile config parsing.
  - deps: none
  - requires-confirmation: false
  - status: todo

- [ ] CLI-002: Implement config.yaml parser for host profiles
  - what: Parse .sloth/config.yaml for host URL, token, and profile selection.
  - do: Implement YAML unmarshaling, profile lookup, and validation.
  - next: Use in workspace init and inspect command.
  - deps: CLI-001
  - requires-confirmation: false
  - status: todo

- [ ] CLI-003: Create lock.json format and file I/O operations
  - what: Define lock schema for contract sync metadata and support file operations.
  - do: Implement read/write/merge/update logic with validation.
  - next: Wire into init and push flows.
  - deps: none
  - requires-confirmation: false
  - status: todo

- [ ] CLI-004: Implement local .sloth/ workspace initialization
  - what: Initialize local .sloth structure with contracts, sets, manifests, and defaults.
  - do: Implement init behavior for new and existing workspace states.
  - next: Enable local list/add flows.
  - deps: CLI-002, CLI-003
  - requires-confirmation: false
  - status: todo

- [ ] CLI-005: Implement sloth contracts list with plugin version flag
  - what: List compatible contracts for a plugin version from configured source.
  - do: Add command with --plugin-version and --source flags.
  - next: Add output formatter support.
  - deps: CLI-001, CLI-004
  - requires-confirmation: false
  - status: todo

- [ ] CLI-006: Implement sloth contracts inspect host command
  - what: Inspect host plugin status and current contract schema inventory.
  - do: Add inspect command and host API calls to inspection endpoints.
  - next: Add output formatter support.
  - deps: CLI-001, CLI-002
  - requires-confirmation: false
  - status: todo

- [ ] CLI-007: Add output formatting for list and inspect
  - what: Support json and table output modes for query commands.
  - do: Add shared output package and wire format flag handling.
  - next: Confirm format structure and apply to add/verify flows as needed.
  - deps: CLI-005, CLI-006
  - requires-confirmation: true
  - status: todo

- [ ] CLI-008: Implement sloth contracts add component
  - what: Add one component contract locally by name/version/source.
  - do: Implement command flags, fetch, validate, and save behavior.
  - next: Reuse logic in parent add command.
  - deps: CLI-001, CLI-004, CLI-007
  - requires-confirmation: false
  - status: todo

- [ ] CLI-009: Implement sloth contracts add set
  - what: Add one contract set locally by name/version/source.
  - do: Implement command flags, fetch, validate, and save behavior for sets.
  - next: Reuse logic in parent add command.
  - deps: CLI-001, CLI-004, CLI-007
  - requires-confirmation: false
  - status: todo

- [ ] CLI-010: Implement sloth contracts add --all
  - what: Bulk add all compatible contracts for a plugin version.
  - do: Add parent add command, convert component/set to subcommands, implement --all.
  - next: Feed results into verify workflow.
  - deps: CLI-008, CLI-009
  - requires-confirmation: false
  - status: todo

- [ ] CLI-011: Build schema version compatibility checking logic
  - what: Validate plugin and contract schema compatibility.
  - do: Implement semver strategy and compatibility matrix checks.
  - next: Reuse in verify command.
  - deps: none
  - requires-confirmation: true
  - status: todo

- [ ] CLI-012: Implement name collision detection in verification
  - what: Detect local or remote contract naming collisions.
  - do: Implement collision checks against local workspace and host inventory.
  - next: Integrate with verify command output.
  - deps: CLI-004, CLI-011
  - requires-confirmation: false
  - status: todo

- [ ] CLI-013: Implement sloth contracts verify --file
  - what: Verify contract file against schema and compatibility before push.
  - do: Wire schema compatibility and collision checks into command report.
  - next: Gate push on verify outcome.
  - deps: CLI-011, CLI-012
  - requires-confirmation: true
  - status: todo

- [ ] CLI-014: Implement sloth contracts push with drift and retry
  - what: Push verified contracts and handle drift/retry behavior.
  - do: Implement ingest call, lock comparison, and retry strategy.
  - next: Validate end-to-end with integration tests.
  - deps: CLI-013, CLI-003, CLI-006
  - requires-confirmation: true
  - status: todo

- [ ] CLI-015: Write integration tests and document CLI contract
  - what: Provide end-to-end test coverage and user-facing CLI contract docs.
  - do: Add integration tests and create docs/CLI-CONTRACT.md specification.
  - next: Prepare for merge.
  - deps: CLI-005, CLI-006, CLI-008, CLI-009, CLI-013, CLI-014
  - requires-confirmation: true
  - status: todo

## In Progress

- [ ] None
  - what: no active task
  - do: pick next task from dependency order
  - next: start CLI-001
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

- CLI-001 -> CLI-002
- CLI-002 -> CLI-004
- CLI-003 -> CLI-004
- CLI-004 -> CLI-005
- CLI-005 -> CLI-007
- CLI-006 -> CLI-007
- CLI-007 -> CLI-008
- CLI-007 -> CLI-009
- CLI-008 -> CLI-010
- CLI-009 -> CLI-010
- CLI-011 -> CLI-012
- CLI-011 -> CLI-013
- CLI-012 -> CLI-013
- CLI-013 -> CLI-014
- CLI-014 -> CLI-015

### Notes

- Risks:
  - Host API behavior for inspect/push may differ from assumptions and require adapter changes.
  - Compatibility matrix rules may need iterative refinement once real contracts are validated.
- Decisions:
  - Keep milestone as one board because total tasks (15) stay within the 20-task limit.
  - Keep verification gates before push to reduce invalid host mutations.
- Next:
  - Start CLI-001.

### Archival

When this milestone is complete, move this file to `docs/archive/` and create a fresh board in `docs/` for the next milestone.
