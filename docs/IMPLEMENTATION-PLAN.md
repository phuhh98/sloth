---
purpose: "High-level implementation plan and timeline estimate based on product docs."
status: "draft"
owner: "product-and-architecture"
last_updated: "2026-05-10"
related_docs:
  - "docs/IDEAS.md"
  - "docs/MILESTONES.md"
  - "docs/COMPONENT-CONTRACTS.md"
  - "docs/REGISTRY.md"
  - "docs/SCHEMA-DRAFTS.md"
---

# sloth Implementation Plan

Date: 2026-05-10

## Current Baseline

Milestone 1 is partially complete.

What is already in place:

- Plugin content types for Component and Page.
- Content API routes for inspection, contract schema, ingest, and page delivery.
- Service and controller scaffolding with Document Service API usage.

Main gap:

- Admin routes and Puck-based builder UI are still scaffold-level and need full MVP implementation.

## Assumptions

- Team size is 1 to 3 engineers.
- Delivery follows roadmap order from docs IDEAS and docs MILESTONES.
- CLI and component hub begin from near-empty package scaffolds.

## High-Level Phases

## Phase 1: Complete Milestone 1

Scope:

- Implement admin APIs for components and pages, including compile action.
- Harden ingest and delivery contract behavior.
- Add auth and policies on write endpoints.
- Build admin builder MVP: component palette, editor, dataset mapping, preview.

Deliverables:

- Stable admin and content API contracts.
- End-to-end page authoring flow in admin.
- Runtime delivery contract ready for frontend consumption examples.

Exit criteria:

- CRUD and compile work end-to-end.
- Ingest and page delivery validated with real data.
- Milestone status updated in docs MILESTONES.

## Phase 2: Milestone 2 CLI Contract Workflow

Scope:

- Build CLI commands: list, inspect, add, verify, push.
- Implement local .sloth folder conventions and lock metadata.
- Add compatibility gate, drift checks, dry-run behavior, push pipeline.
- Package Go CLI binaries for macOS, Linux, and Windows.
- Distribute CLI through npm with package metadata and platform-aware binary resolution.

Deliverables:

- Reliable contract sync from local to host.
- Blocking compatibility checks with clear errors.
- npm-installable CLI package that executes the correct binary by platform.

Exit criteria:

- CLI completes verify -> inspect -> compare -> ingest pipeline safely.
- Incompatible versions are blocked before write.
- npm install and first-run command succeed on macOS, Linux, and Windows using published prebuilt binaries.

## Phase 3: Milestone 3 Component Hub Starter + Renderer Examples

Scope:

- Publish starter packs and package descriptors.
- Provide reference runtime integration examples for first-level page payload strategy.

Deliverables:

- At least one high-quality starter pack.
- Working sample showing page consumption and component mapping.

Exit criteria:

- New project can install pack and render a page with minimal custom wiring.

## Phase 4: Milestone 4 Public Registry MVP

Scope:

- Build metadata API, immutable artifacts, and checksum verification.
- Add CLI registry workflows for search, info, add, update.

Deliverables:

- Public free registry MVP.
- Compatibility-aware package install flow.

Exit criteria:

- Users can discover and install package artifacts reliably.

## Phase 5: Milestone 5 Private and Paid Extensions

Scope:

- Add authentication, private namespaces, entitlement checks, and billing integration.

Deliverables:

- Private package support.
- Paid package entitlement lifecycle.

Exit criteria:

- Access control and billing are production-safe with audit trails.

## ADR Checkpoints

Create ADRs early for these decisions:

- CLI verification ownership versus host ingest ownership.
- compiledConfig generation timing: on save, on publish, or both.
- Runtime payload boundary and first-level data policy.
- Registry API language and stack.
- Strict-fail versus warning-first compatibility policy.
- CLI distribution strategy: npm-delivered prebuilt binaries versus direct go install.

## Timeline Estimate

## Solo full-time (strong Strapi plus Go plus TypeScript)

- Milestone 1: 3 to 5 weeks
- Milestone 2: 4 to 6 weeks
- Milestone 3: 2 to 4 weeks
- Milestone 4: 4 to 7 weeks
- Milestone 5: 6 to 10 weeks

Total:

- Milestones 1 to 3: about 9 to 15 weeks
- Milestones 1 to 5: about 19 to 32 weeks

## Small team (2 to 3 engineers with parallel tracks)

- Milestone 1: 2 to 4 weeks
- Milestone 2: 3 to 5 weeks
- Milestone 3: 2 to 3 weeks
- Milestone 4: 3 to 5 weeks
- Milestone 5: 4 to 8 weeks

Total:

- Milestones 1 to 3: about 7 to 12 weeks
- Milestones 1 to 5: about 14 to 25 weeks

## Practical Comment

If the priority is fastest usable value, complete Milestones 1 and 2 first. This yields a useful foundation in about 2 to 3 months for solo execution, or about 1.5 to 2 months for a small focused team.

## Risks That Can Extend Timeline

- Admin builder complexity and UX iteration loop.
- Contract schema drift between plugin and CLI.
- Runtime delivery enrichment and content query edge cases.
- Registry security and governance requirements arriving earlier than planned.

## Update Rule

When implementation status changes:

1. Update this plan if scope or sequencing changes.
2. Update docs MILESTONES status and notes.
3. Reflect strategic direction changes in docs IDEAS.
