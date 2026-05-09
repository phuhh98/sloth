---
applyTo:
  - "packages/component-hub/src/contracts/**/*.json"
  - "contracts/schema/component-contract.schema.json"
  - "apps/docs/static/schemas/component-contract/**/*.json"
---

# Component Contract Authoring

When creating or editing sloth component contracts:

- Treat `packages/component-hub/src/contracts/**` as the source of truth for contract instances.
- Organize contract source by release version first: `src/contracts/<release-version>/components/<component-name>/contract.json`.
- Each release folder must have a `manifest.json` that describes the full contract set for that release.
- Every contract instance file must include `$schema` pointing at the hosted contract schema URL for its `schemaVersion`.
- Keep the release folder name aligned with each contract file's `version` field.
- Keep release manifest `components.<name>.contentHash` aligned with the exact contents of each `contract.json` file in that release.
- Keep `name` as a lowercase slug and keep `renderMeta.rendererKey` aligned with the frontend implementation key.
- Contract releases are immutable after introduction. If any contract changes or a new contract is added, create a new release folder rather than editing an old release in place.
- When a newer release is added, keep older release folders in place and set `deprecatedAt` on each non-latest release manifest to at least 6 months in the future.
- When the contract schema changes, update both:
  - `contracts/schema/component-contract.schema.json`
  - `apps/docs/static/schemas/component-contract/**/schema.json`
- Contract instance files are schema instances, not schemas themselves: use `$schema`, not `$id`.
