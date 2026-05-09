---
description: "Use when planning or executing implementation milestones. Manage work with Markdown Kanban boards in docs/, archive completed milestone boards, and decompose milestones into executable tasks with minimal cross-package coupling."
name: "Implementation Kanban Management"
---

# Implementation Kanban Management

When working on implementation plans and milestones in this repository:

- Use Markdown Kanban boards in `docs/` as the detailed task tracker for execution.
- Prefer the VS Code extension `holooooo.markdown-kanban` format so board state is easy to monitor in-editor.
- Keep the Kanban board more detailed than the high-level implementation plan.
- Create one Kanban board per stage or major step from the high-level implementation plan.
- For each milestone, explicitly break down the milestone into smaller executable tasks.
- **Limit each Kanban board to a maximum of 20 tasks.** If a milestone exceeds 20 tasks, divide it into separate Kanban boards (e.g., KANBAN-MILESTONE-2-PHASE-A.md, KANBAN-MILESTONE-2-PHASE-B.md).
- Decompose tasks to reduce inter-package coupling and avoid requiring simultaneous changes across multiple packages when possible.
- Sequence tasks so each ticket can be completed independently with clear inputs and outputs.
- Keep a clear `To Do`, `In Progress`, `Blocked`, and `Done` flow in each board.
- Use `##` headings only for Kanban columns (`To Do`, `In Progress`, `Blocked`, `Done`) because markdown-kanban treats H2 as columns.
- Use `###` or lower heading levels for non-column sections (for example scope, notes, dependency plan, archival) to avoid creating unintended columns.
- Keep task cards short and concise. Each task card should include only:
  - what this task is
  - what to do
  - what to do next
  - dependencies
  - requires-confirmation flag (for tasks needing user review)
  - current status
- Mark tasks with `requires-confirmation: true` when they need user intervention for:
  - Frontend or UI changes requiring visual review
  - Code style, logic, or design decisions requiring approval
  - Environment setup or docker/browser verification
  - When marked as requiring confirmation, agent will use `vscode_askQuestions` to get multi-choice user input before proceeding.
- Avoid chronological logs or per-task timestamp history.
- Record timestamp only at milestone level (for example `milestone_updated_at`).
- When a milestone is completed, move its Kanban file from `docs/` to `docs/archive/`.
- When starting work on a new milestone from the implementation plan, initialize a new Kanban board in `docs/`.
- Keep board titles and ticket labels aligned with milestone names used in `docs/MILESTONES.md`.
- Update the board continuously during execution so users can monitor completed work and next tasks without reading commit history.
