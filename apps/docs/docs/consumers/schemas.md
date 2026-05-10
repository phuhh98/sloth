---
sidebar_position: 2
---

# Schemas

This page lists stable, versioned schema and registry URLs served by the sloth docs site.

Contract and schema release artifacts are published to GHCR (OCI) from contracts package release automation.
For JSON Schema tooling compatibility, the canonical `$schema` URL remains the docs-hosted HTTPS URL below.

## Component Contract

- Version: `0.0.1`
- Hosted URL: [https://phuhh98.github.io/sloth/schemas/component-contract/0.0.1/schema.json](https://phuhh98.github.io/sloth/schemas/component-contract/0.0.1/schema.json)

## CLI Config YAML

- Version: `0.0.1`
- Hosted URL: [https://phuhh98.github.io/sloth/schemas/cli-config/0.0.1/schema.json](https://phuhh98.github.io/sloth/schemas/cli-config/0.0.1/schema.json)

You can reference this schema in `.sloth/config.yaml` for editor validation:

```yaml
# yaml-language-server: $schema=https://phuhh98.github.io/sloth/schemas/cli-config/0.0.1/schema.json
currentProfile: default
profiles:
  default:
    host: http://localhost:1337
    authorizationToken: ""
```

## Component Registry

- Registry index: [https://phuhh98.github.io/sloth/registry/index.json](https://phuhh98.github.io/sloth/registry/index.json)
- Contract releases index: [https://phuhh98.github.io/sloth/registry/contracts/index.json](https://phuhh98.github.io/sloth/registry/contracts/index.json)
- Themes index: [https://phuhh98.github.io/sloth/registry/themes/index.json](https://phuhh98.github.io/sloth/registry/themes/index.json)

### Versioned Folder Convention

Use immutable version folders so old references never break:

```text
apps/docs/static/registry/
  index.json
  contracts/
    index.json
    <release-version>/
      manifest.json
      components/
        <component-name>/
          contract.json
  themes/
    index.json
```

Example versioned contract release artifacts:

- [https://phuhh98.github.io/sloth/registry/contracts/0.0.1/manifest.json](https://phuhh98.github.io/sloth/registry/contracts/0.0.1/manifest.json)
- [https://phuhh98.github.io/sloth/registry/contracts/0.0.1/components/hero-banner/contract.json](https://phuhh98.github.io/sloth/registry/contracts/0.0.1/components/hero-banner/contract.json)
- [https://phuhh98.github.io/sloth/registry/contracts/0.0.1/components/article-teaser/contract.json](https://phuhh98.github.io/sloth/registry/contracts/0.0.1/components/article-teaser/contract.json)

## Notes

- Schema URLs should be immutable per version.
- Registry artifacts should be immutable per version as well.
- Publish new versions at new paths instead of replacing existing files.
- GHCR is the immutable artifact/provenance backend; docs URL is the canonical `$schema` endpoint for validators/editors.
- During docs build/release, sync promoted schema versions from GHCR artifacts into `apps/docs/static/schemas/...` before compiling Docusaurus output.
