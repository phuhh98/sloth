---
purpose: "Capture architecture decisions and open points for distributing the Go CLI through npm with cross-platform binaries."
status: "active"
owner: "platform-and-cli"
last_updated: "2026-05-10"
related_docs:
  - "docs/MILESTONES.md"
  - "docs/IMPLEMENTATION-PLAN.md"
  - "docs/KANBAN-MILESTONE-2.md"
  - "docs/COMPONENT-CONTRACTS.md"
  - "docs/REGISTRY.md"
---

# CLI Distribution Notes: Go Binary via npm

Date: 2026-05-10

## Confirmed Direction

- CLI implementation remains Go + Cobra.
- Distribution channel is npm.
- Delivery target is prebuilt binaries for:
  - macOS
  - Linux
  - Windows
- npm package includes package metadata through package.json.

## Clarification: npm as Source vs npm as Distribution

- npm can be used as a contract artifact source for CLI contract resolution.
- npm can also be used as the CLI runtime delivery channel.
- These are different concerns and should be kept explicit in command and release docs.

## Milestone and Plan Alignment Completed

- Milestone 2 now explicitly includes npm-based binary distribution.
- Phase 2 implementation plan now includes binary packaging and npm distribution as scope.
- Phase 2 exit criteria now requires successful npm install and first-run on macOS/Linux/Windows.
- Milestone 2 kanban now includes dedicated tasks:
  - cross-platform build and checksum pipeline
  - npm wrapper package/resolver
  - install-flow and release contract validation

## npm Registry Binary Hosting Answer

Yes. npm registry can host binaries as package files inside the published package tarball.

Practical interpretation:

- package.json remains required package metadata
- binary files are included as package artifacts
- install/runtime logic resolves OS and architecture to execute the right binary

## Recommended Packaging Patterns

1. Single package with all platform binaries

- Simpler release flow
- Larger install size

2. Meta package plus per-platform packages using optional dependencies

- Smaller install per platform
- More release coordination and version management

## Can We Avoid Many Static Package Folders?

Yes, with one constraint from npm.

- npm publishes one package per package.json tarball.
- That means each published package still needs its own package.json at publish time.
- But those package directories do not need to be manually maintained in source control.

Recommended approach:

- Keep one source package at packages/cli.
- Add one distribution config manifest (for example platforms, package names, binary paths, and version policy).
- Use a generator script to emit temporary publish directories under dist/publish-packages/.
- Publish generated packages from dist instead of hand-maintained per-platform folders.

Benefits:

- one canonical config source
- less duplicated boilerplate across package folders
- easier platform matrix updates

Tradeoff:

- release tooling is slightly more complex and must be tested in CI

## Open Decisions to Resolve

- Choose packaging pattern: single package or per-platform optional packages.
- Define supported architecture matrix explicitly, including amd64 and arm64.
- Decide checksum verification responsibility:
  - release-time only
  - release-time plus runtime verification
- Decide signing/provenance target for milestone scope.
- Confirm fallback behavior when platform binary is missing.

## Suggested Definition of Done for Distribution Track

- Build pipeline produces deterministic binaries for all target OS/arch.
- Checksums generated and published with release artifacts.
- npm package install works in clean environments across target platforms.
- First CLI command execution succeeds after install.
- Error message is clear when platform is unsupported.
- Version mapping between Go binary and npm package is documented.
