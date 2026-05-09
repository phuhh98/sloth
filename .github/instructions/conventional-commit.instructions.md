---
description: "Use when creating git commits in this repository. Enforce conventional commit format, derive commit intent from current Kanban task, and avoid MCP git tools."
name: "Conventional Commit Workflow"
---

# Conventional Commit Workflow

When asked to commit changes in this repository:

- Use CLI git commands in terminal for commit actions.
- Do not use Git MCP tools for add, commit, or other git operations.
- Build commit message using the conventional commit format:
  - type(scope): title
- Add extra details as bullet lines in the commit body when needed:
  - Use lines that start with "- "
  - Keep each bullet concise and change-focused

## Commit Intent Source

- First consult the current milestone Kanban task in docs/ to determine commit intent and scope.
- If the Kanban task context is not sufficient, inspect staged/unstaged diff and derive a precise conventional commit message from actual code changes.
- Prefer the narrowest valid scope based on changed package or module.

## Quality Rules

- Keep title imperative and specific.
- Match type to actual change (feat, fix, refactor, docs, test, chore, perf, build, ci).
- Do not combine unrelated changes into one commit message.
- If changes span distinct concerns, split into multiple conventional commits.
