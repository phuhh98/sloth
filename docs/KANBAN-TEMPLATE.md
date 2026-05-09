---
purpose: "Template for milestone-level implementation Kanban tracking"
status: "active-template"
owner: "platform"
last_updated: "2026-05-10"
related_docs:
  - "docs/IMPLEMENTATION-PLAN.md"
  - "docs/MILESTONES.md"
---

# Milestone Kanban: <Milestone Name>

Use this board for detailed execution tracking of one milestone.

## Scope

- Milestone: <Milestone ID and Name>
- Goal: <What is delivered when milestone is done>
- Constraints: <Key boundaries and non-goals>
- milestone_updated_at: <YYYY-MM-DD>

## Task Decomposition Rules

- Split into small executable tasks.
- Prefer package-local tasks before cross-package integration tasks.
- Minimize tasks that require concurrent edits in different packages.
- Define clear dependency order.
- For tasks requiring user confirmation (frontend review, code style review, docker/browser verification), mark with `requires-confirmation: true`.
  - Agent will use multi-choice selection to get user input before proceeding.
  - Examples: UI changes, design decisions, environment setup verification.

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

Do not add per-task timestamps.

### To Do

- [ ] <Task 1>
  - what: <what this task is>
  - do: <what to do now>
  - next: <what to do next>
  - deps: <dependency or "none">
  - requires-confirmation: false
  - status: todo

- [ ] <Task 2>
  - what: <what this task is>
  - do: <what to do now>
  - next: <what to do next>
  - deps: <dependency or "none">
  - status: todo

### In Progress

- [ ] <Task currently being worked>
      - what: <what this task is>
      - do: <what to do now>
      - next: <what to do next>
      - deps: <dependency or "none">
      - status: in-progress

### Blocked

- [ ] <Task blocked and reason>
      - what: <what this task is>
      - do: <what to do now>
      - next: <what to do next once unblocked>
      - deps: <blocking dependency>
      - status: blocked

### Done

- [x] <Completed task>
      - what: <what this task is>
      - do: <what was completed>
      - next: <next follow-up or "none">
      - deps: <dependency or "none">
      - status: done

## Dependency Plan

- Task <A> -> Task <B>
- Task <B> -> Task <C>

## Notes

- Risks:
  - <Risk>
- Decisions:
  - <Decision>
- Next:
  - <Next task>

## Archival

When this milestone is complete, move this file to `docs/archive/` and create a fresh board in `docs/` for the next milestone.
