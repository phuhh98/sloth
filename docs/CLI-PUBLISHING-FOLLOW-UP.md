---
purpose: "Track post-Milestone-2 follow-up tasks for CLI npm publishing automation and CI/CD"
status: "todo"
owner: "platform"
created_date: "2026-05-10"
related_docs:
  - "docs/MILESTONES.md"
  - "docs/archive/KANBAN-MILESTONE-2.md"
---

# CLI Publishing Automation - Follow-up Tasks

This document tracks follow-up work required after Milestone 2 CLI completion to enable automated publishing to npm.

## Context

Milestone 2 CLI implementation (packages/cli) is now complete with:

- Local `.sloth/` workspace management
- Component contract commands (init, list, inspect, add, verify, push)
- YAML/ENV/default configuration precedence
- Cross-platform Go binary build pipeline
- Config-driven npm publish package generation
- Taskfile commands for build/test/release-prep

**Gap:** No CI/CD workflow exists to automate npm publishing on version tags or releases.

## Required Tasks

### 1. GitHub Actions CI/CD Workflow for CLI Release

**Purpose:** Automate npm publishing when a release tag is pushed or release is created.

**Deliverables:**

- `.github/workflows/cli-publish.yml` workflow that:
  - Triggers on version tags (e.g., `cli@v*` or `@sloth/cli@v*`)
  - Runs `task cli-release-prep` to generate cross-platform binaries and checksums
  - Generates/updates publish-packages for each supported platform
  - Publishes to npm using authenticated `npm publish` or pnpm equivalent
  - Creates GitHub release with checksums and binary artifacts
  - Supports dry-run mode for validation before actual publish

**Acceptance criteria:**

- Workflow tests successfully on a dry-run tag
- Produces correct platform-specific package directories
- Generates release notes with checksums
- Does not publish on non-release commits

### 2. npm Publishing Credentials and Access Control

**Purpose:** Secure npm publishing authentication and scoped package management.

**Deliverables:**

- Configure npm token/authentication in GitHub secrets (e.g., `NPM_TOKEN`)
- Verify `@sloth/cli` scope is registered and ownership is set
- Document npm publishing policies (e.g., SemVer versioning, tags, dist-tags)
- Set up npm two-factor authentication (2FA) policy if applicable
- Create access control rules for who can publish (optional for private org)

**Acceptance criteria:**

- GitHub Actions workflow can authenticate to npm via token
- CI can successfully publish test version to npm registry (dry-run or separate registry)
- npm scope and package ownership are correct

### 3. Release Tagging and Version Management Strategy

**Purpose:** Establish clear conventions for tagging releases and updating package versions.

**Deliverables:**

- Document version numbering convention (e.g., semantic versioning, aligned with plugin versions)
- Define tag naming convention (e.g., `cli/v0.0.2`, `@sloth/cli@0.0.2`)
- Create script or documentation for incrementing versions in:
  - `packages/cli/package.json`
  - Potentially `packages/cli/go.mod` and version constants in Go code
- Add pre-release and post-release checklist (commit message format, CHANGELOG updates, etc.)

**Acceptance criteria:**

- Tag format is consistent and machine-parseable
- Version update process is repeatable and documented
- Release checklist prevents accidental publishes with incorrect versions

### 4. Documentation: CLI Release Runbook

**Purpose:** Provide step-by-step guide for maintainers to publish a new CLI version.

**Deliverables:**

- Document in `packages/cli/RELEASE.md` or `docs/CLI-RELEASE-RUNBOOK.md` with:
  - Pre-release checklist (tests, docs, version updates)
  - How to create a release tag locally and push it
  - Expected CI/CD behavior (what workflow runs, what checks)
  - How to verify the npm package after publish
  - Rollback procedure if publish fails or is incorrect

**Acceptance criteria:**

- New maintainer can follow runbook without extra help
- Runbook is tested (dry-run or actual publish)
- Rollback steps are clear and include npm unpublish guidance if needed

## Task Dependencies

```
[1] Workflow setup -> [2] Auth/credentials -> [3] Versioning -> [4] Documentation
```

## Implementation Order

1. **Task 1:** Set up GitHub Actions workflow for publish trigger and dry-run testing.
2. **Task 2:** Configure npm token and validate publish permissions.
3. **Task 3:** Define versioning and tagging strategy.
4. **Task 4:** Create release runbook for team reference.

## Success Metrics

- [ ] CI workflow triggers on version tag push
- [ ] CLI package published to npm successfully on first release
- [ ] Checksums and binaries are available in release assets
- [ ] Team can run full release cycle with single tag push + CI automation
- [ ] Rollback procedure tested and documented

## Notes

- Consider using `changesets` or similar tool for automated changelog and version management if scaling to multiple packages.
- Keep npm token rotation schedule and security review as ongoing operational tasks.
- Monitor npm registry for any package transfer or access issues post-publish.
