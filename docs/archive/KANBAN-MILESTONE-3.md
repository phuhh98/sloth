---
purpose: "Milestone-level implementation Kanban tracking for component hub starter packs and renderer examples"
status: "archived"
owner: "platform"
last_updated: "2026-05-10T11:45:00Z"
related_docs:
  - "docs/IMPLEMENTATION-PLAN.md"
  - "docs/MILESTONES.md"
  - "docs/IDEAS.md"
  - "docs/REGISTRY.md"
  - "docs/KANBAN-OCI-REGISTRY-PIVOT.md"
---

# Milestone Kanban: Milestone 3 - Component Hub Starter + Renderer Examples

Archived: Milestone 3 core scope is complete. OCI pivot and remaining follow-up tasks are tracked in `docs/KANBAN-OCI-REGISTRY-PIVOT.md`.

Use this board for detailed execution tracking of one milestone.

### Scope

- Milestone: Milestone 3 - Component Hub Starter + Renderer Examples
- Goal: Publish starter packs and runtime integration examples for rendering page payloads with component mappings
- Constraints: Keep core milestone under 20 tasks; HUB-000..HUB-012 (13 tasks) plus PUB-001..PUB-004 (4 tasks, Milestone 2 follow-up)
- milestone_updated_at: 2026-05-10T00:40:00Z

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

- [x] HUB-000: Finalize dynamic content component shortlist and SEO mapping
  - what: Lock first-class dynamic content contracts for article/post ecosystems before schema implementation.
  - do: Validate starter set for feed, related carousel, article teaser, author, breadcrumb, and seo-head components.
  - next: Feed approved list into static/dynamic block planning tasks.
  - deps: none
  - requires-confirmation: true
  - status: done

- [x] HUB-008: Finalize static block shortlist and priorities
  - what: Lock static/reusable blocks and sections used across non-dynamic pages.
  - do: Validate Text, Card, Stat, CTA, Features, Testimonials, FAQ, Pricing, and base layout wrappers.
  - next: Feed approved static set into shared block contract draft.
  - deps: none
  - requires-confirmation: true
  - status: done

- [x] HUB-009: Define shared layout/section/block base contract and SEO slots
  - what: Establish common attributes for all component contracts plus SEO behavior boundaries.
  - do: Specify base fields (id, rendererKey, visibility, style tokens, dataSource), section metadata, and block-level SEO contribution model.
  - next: Use shared base contract in HUB-001 starter pack manifest and per-component schemas.
  - deps: HUB-000, HUB-008
  - requires-confirmation: false
  - status: done

- [x] HUB-001: Define starter pack manifest schema
  - what: Create starter pack manifest contract for component hub artifacts.
  - do: Add JSON schema and TS types for pack metadata and component references.
  - next: Implement validator and build script.
  - deps: HUB-009
  - requires-confirmation: false
  - status: done

- [x] HUB-002: Implement starter pack validator and builder
  - what: Build and validate starter pack artifacts from source descriptors.
  - do: Add scripts and tests in packages/component-hub.
  - next: Generate first pack artifact.
  - deps: HUB-001
  - requires-confirmation: false
  - status: done

- [x] HUB-003: Publish first starter pack artifact
  - what: Create one high-quality starter pack with docs metadata.
  - do: Add pack descriptor and generated artifact under docs static registry path.
  - next: Add renderer mapping example.
  - deps: HUB-002
  - requires-confirmation: false
  - status: done

- [x] HUB-004: Add runtime renderer mapping utility
  - what: Implement utility that maps component contract rendererKey to React implementation.
  - do: Add package-level API and tests for mapping behavior.
  - next: Integrate with sample app page payload.
  - deps: HUB-003
  - requires-confirmation: false
  - status: done

- [x] HUB-005: Build frontend runtime example page
  - what: Demonstrate rendering a page payload using starter pack component mappings.
  - do: Add example in apps/frontend with first-level payload handling.
  - next: Document usage in public docs.
  - deps: HUB-004
  - requires-confirmation: true
  - status: done

- [x] HUB-006: Add docs sync workflow for contracts and packs
  - what: Ensure docs static registry receives generated contract/pack artifacts per release.
  - do: Add script and tests for sync behavior.
  - next: Wire task commands and CI hooks.
  - deps: HUB-003
  - requires-confirmation: false
  - status: done

- [x] HUB-007: Document component hub and renderer integration
  - what: Publish public docs for starter packs and renderer integration patterns.
  - do: Add docs pages in apps/docs/docs with examples and versioning guidance.
  - next: Mark milestone complete after validation.
  - deps: HUB-005, HUB-006
  - requires-confirmation: false
  - status: done

- [x] HUB-010: Define OpenAPI contract for sloth CMS-agnostic API surface
  - what: Define `openapi.yaml` for inspection, contract discovery, and ingest APIs using a general CMS integration model (not Strapi-specific).
  - do: Draft OpenAPI paths/schemas based on docs/IDEAS.md intent, including auth, errors, pagination, and versioning policy.
  - next: Use OpenAPI as source for mock-server behavior and integration contract tests.
  - deps: HUB-009
  - requires-confirmation: true
  - status: done

- [x] HUB-011: Build component-hub mock server with seeded contract set
  - what: Provide a deterministic mock server inside component-hub that serves a defined starter contract dataset via the OpenAPI spec.
  - do: Implement mock endpoints, seed fixtures (static + dynamic contracts), and runner scripts for local/CI execution.
  - next: Run CLI integration tests against mock server.
  - deps: HUB-003, HUB-010
  - requires-confirmation: false
  - status: done

- [x] HUB-012: Add CLI integration tests against OpenAPI mock server and track gaps
  - what: Validate `sloth` CLI behavior end-to-end using OpenAPI-defined mock server + seeded contracts.
  - do: Add integration test suite for list/inspect/add/verify/push flows, assert outputs/errors, and record known flaws with improvement notes.
  - next: Prioritize fixes and feed defects into next milestone backlog.
  - deps: HUB-011
  - requires-confirmation: false
  - status: done

**Post-Milestone-2 Follow-up: CLI Publishing Automation (requires user intervention)**

- [ ] PUB-001: GitHub Actions CI/CD workflow for CLI release
  - what: Set up GitHub Actions to automate npm publishing on version tags.
  - do: Create `.github/workflows/cli-publish.yml` with cross-platform build, checksum generation, and npm publish steps.
  - next: Test workflow on dry-run tag.
  - deps: none (Milestone 2 complete)
  - requires-confirmation: true
  - status: todo

- [ ] PUB-002: Configure npm publishing credentials and access control
  - what: Set up npm token authentication, verify package scope ownership, and define publishing policy.
  - do: Add npm token to GitHub secrets, verify @sloth/cli scope, document SemVer policy.
  - next: Validate workflow can authenticate to npm.
  - deps: PUB-001
  - requires-confirmation: true
  - status: todo

- [ ] PUB-003: Define release tagging and version management strategy
  - what: Establish version numbering convention and tag format for releases.
  - do: Document tag naming (e.g., cli/v0.0.2), version update process, and pre-release checklist.
  - next: Create release runbook.
  - deps: none
  - requires-confirmation: true
  - status: todo

- [ ] PUB-004: Document CLI release runbook for maintainers
  - what: Provide step-by-step guide for publishing new CLI versions.
  - do: Create packages/cli/RELEASE.md or docs/CLI-RELEASE-RUNBOOK.md with pre-release checklist, tag procedure, CI behavior, and rollback steps.
  - next: Test runbook with first release.
  - deps: PUB-003
  - requires-confirmation: true
  - status: todo

## In Progress

- [ ] None
  - what: no active task
  - do: keep PUB tasks pending until user-driven credential/release-policy confirmation
  - next: run first CLI publish dry-run once secrets and tags are configured
  - deps: PUB-001, PUB-002, PUB-003
  - status: in-progress

## Blocked

- [ ] None
  - what: no blocked task
  - do: continue planned sequence
  - next: raise blockers here when encountered
  - deps: none
  - status: blocked

## Done

- [x] HUB-000..HUB-012
  - what: Milestone 3 core scope implemented (shortlist/base contract decisions, starter packs, runtime mapping example, docs sync, OpenAPI mock server, CLI integration tests)
  - do: validated with component-hub tests, docs registry tests, and CLI Go tests
  - next: close milestone after optional visual review of frontend example and docs site render
  - deps: none
  - status: done

### Dependency Plan

- HUB-000 -> HUB-009
- HUB-008 -> HUB-009
- HUB-009 -> HUB-001
- HUB-001 -> HUB-002
- HUB-002 -> HUB-003
- HUB-003 -> HUB-004
- HUB-003 -> HUB-006
- HUB-004 -> HUB-005
- HUB-005 -> HUB-007
- HUB-006 -> HUB-007
- HUB-009 -> HUB-010
- HUB-003 -> HUB-011
- HUB-010 -> HUB-011
- HUB-011 -> HUB-012
- PUB-001 -> PUB-002
- PUB-003 -> PUB-004

### Notes

- Risks:
  - Runtime payload shape may evolve while plugin delivery contracts are still maturing.
- Decisions:
  - Start with one starter pack before broadening pack catalog.
  - Add HUB-000 and HUB-008 as gating refinement tasks for dynamic and static component contract sets.
  - Add HUB-009 to standardize shared layout/section/block fields and SEO mapping before schema implementation.
  - Add HUB-010 to standardize a CMS-agnostic API contract via OpenAPI before mock/integration testing.
  - Add HUB-011/HUB-012 to validate CLI against deterministic API behavior and capture quality gaps.
  - HUB-001..HUB-007 are Milestone 3 core scope; PUB-001..PUB-004 are Milestone 2 follow-up (CLI publishing automation).
  - User confirmation required for all PUB-\* tasks (credential setup, release policy, runbook validation).
  - HUB-000 and HUB-008 were completed via documented shortlist decisions in `docs/COMPONENT-HUB-BASE-CONTRACT-AND-SHORTLIST.md`.
  - HUB-005 is implemented as a runtime script example and still benefits from optional UI/UX review if promoted to a production frontend route.
- Next:
  - Execute PUB-001..PUB-004 when release credentials and maintainers are available for confirmation.

### Archival

When this milestone is complete, move this file to `docs/archive/` and create a fresh board in `docs/` for the next milestone.

Additional archival rule:

- When a task is completed and its related reference docs are not expected to be reused, move those docs into `docs/archive/`.
