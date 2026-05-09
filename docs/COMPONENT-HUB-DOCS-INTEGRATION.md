---
purpose: "Define how component-hub becomes the source of truth for registry artifacts and React component implementations used by the docs site."
status: "proposed"
owner: "platform"
last_updated: "2026-05-10"
related_docs:
  - "docs/REGISTRY.md"
  - "docs/MILESTONES.md"
  - "docs/COMPONENT-CONTRACTS.md"
---

# Component Hub to Docs Registry Integration

Date: 2026-05-10
Status: Proposed

## Goal

Use `packages/component-hub` as the single source for:

- component contracts
- component manifests
- basic React frontend component implementations

Then, during docs build, publish contract artifacts into:

- `apps/docs/static/registry/contracts/`

This keeps GitHub Pages as a simple static host while moving source ownership to `component-hub`.

## Proposed Source Structure (`packages/component-hub`)

```text
packages/component-hub/
  package.json
  src/
    contracts/
      <release-version>/
        manifest.json
        components/
          <component-name>/
            contract.json
    react/
      <component-name>/
        <version>/
          index.tsx
  scripts/
    build-registry.mjs
  dist/
    registry/
      contracts/
        <release-version>/
          manifest.json
          components/
            <component-name>/
              contract.json
    react/
      <component-name>/
        <version>/
          index.js
```

## Build-Time Integration Flow

1. Install workspace dependencies with `pnpm install`.
2. Build component-hub artifacts:
   - `pnpm --filter @sloth/component-hub run build:registry`
3. Sync built contracts/manifests into docs static hosting path:

- copy from `packages/component-hub/dist/registry/contracts/**`
- to `apps/docs/static/registry/contracts/**`

4. Regenerate index files:

- `apps/docs/static/registry/contracts/index.json`
- `apps/docs/static/registry/index.json`

5. Build docs:
   - `pnpm --filter apps-docs build`

## Proposed Scripts and Tasks

At package level (`packages/component-hub/package.json`):

- `build:registry` — validates contracts and emits immutable artifacts under `dist/registry`
- `build:react` — compiles basic React implementations under `dist/react`
- `build` — runs both targets

At docs level (`apps/docs/package.json`):

- `registry:prepare` — build component-hub contract release artifacts and sync into docs `static/registry/contracts`
- `prebuild` — run `registry:prepare` before docs build
- `prestart` — run `registry:prepare` before local docs dev

At root task level (`Taskfile.yml`):

- `build-component-hub` — run component-hub build
- `sync-docs-registry` — run docs registry sync and index generation
- `build-docs` — depends on component-hub build and docs registry sync before Docusaurus build

## Versioning Rules

- All published artifact versions are immutable.
- Each contract release is path-scoped:
  - `registry/contracts/<release-version>/...`
- Index files can evolve, but versioned payload files cannot be mutated in place.
- Contract changes require a new release version folder that contains the full contract set for that release.

## Validation Rules

Before copying to docs static registry:

- validate each `contract.json` against:
  - `https://phuhh98.github.io/sloth/schemas/component-contract/0.0.1/schema.json`
- ensure every manifest references an existing contract file
- ensure component/version tuple uniqueness

## Why This Direction

- avoids duplicate manual editing in `apps/docs/static`
- keeps ownership in a dedicated package (`component-hub`)
- allows future npm/git distribution using the same artifacts
- keeps GitHub Pages deployment simple and static-host-only

## Open Questions

- Should component-hub publish artifacts as npm package files in addition to workspace sync?
- Should React implementation output include TypeScript declarations for downstream consumers?
- Should docs build fail hard if any contract validation fails, or allow warning mode in early phase?
