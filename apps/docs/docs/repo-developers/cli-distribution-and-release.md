---
title: CLI Distribution and Release
---

# CLI Distribution and Release

The sloth CLI is implemented in Go and distributed through npm-compatible packaging workflows.

## Build Targets

Cross-platform binaries are built for:

- darwin amd64
- darwin arm64
- linux amd64
- linux arm64
- windows amd64
- windows arm64

## Taskfile Commands

Use these commands from repository root:

```bash
task cli-test
task cli-build-cross
task cli-generate-publish-packages
task cli-smoke-install
task cli-release-prep
```

## Output Artifacts

Generated outputs:

- `packages/cli/dist/bin/<os>-<arch>/sloth[.exe]`
- `packages/cli/dist/checksums.txt`
- `packages/cli/dist/publish-packages/*`

## Distribution Config

Platform packaging is driven by:

- `packages/cli/distribution.config.json`

This avoids maintaining many static per-platform package folders in source.

## npm Wrapper Behavior

The wrapper executable script resolves platform and arch at runtime, then runs the matching binary.

- wrapper path: `packages/cli/bin/sloth.js`
- build script: `packages/cli/scripts/build-cross.mjs`
- publish package generator: `packages/cli/scripts/generate-publish-packages.mjs`

## Release Checklist

1. Run `task cli-test`.
2. Run `task cli-release-prep`.
3. Verify checksums and generated package manifests.
4. Validate first-run command on each supported platform in CI or release pipeline.

## Related

- [CLI Getting Started](../consumers/cli-getting-started)
- [CLI Command Reference](../consumers/cli-contract)
