---
applyTo: "docs/**,**/*.md,**/*.mdx,apps/docs/**"
---

# Documentation Location Convention

Documentation that needs to be **publicly available** (user guides, API references, schema specs, tutorials, changelogs) must live inside `apps/docs/` and be authored as Docusaurus pages.

## Rules

- **Public docs → `apps/docs/docs/`**: Any `.md` or `.mdx` file intended for external audiences (developers integrating sloth, end-users, contributors) goes here. Docusaurus renders these as versioned, searchable documentation pages at `https://phuhh98.github.io/sloth/`.
- **Schema artifacts → `apps/docs/static/schemas/`**: JSON Schema files that must be hosted at a canonical public URL are placed here under the path structure `static/schemas/<namespace>/<artifact>/<version>/schema.json`.
- **Internal/planning docs → `docs/`** (repo root): Architecture diagrams, ADRs, implementation plans, idea drafts, and milestone tracking that are NOT intended for public users stay in the root `docs/` folder.

## Directory Structure

```
apps/docs/
  docs/                   # Public Docusaurus pages (user-facing)
    intro.mdx
    schemas.md
    <category>/
      _category_.json
      *.md | *.mdx
  static/
    schemas/              # Publicly hosted JSON Schema artifacts
      sloth/
        <artifact-name>/
          <version>/
            schema.json
  src/
    pages/                # Custom React pages (landing page, etc.)
    components/           # Shared React components for the docs site

docs/                     # Internal planning & architecture (NOT public)
  IDEAS.md
  MILESTONES.md
  ARCHITECTURE-DIAGRAM.md
  IMPLEMENTATION-PLAN.md
  adr/
```

## When Adding New Documentation

1. Decide: is this for **external users** or **internal planning**?
2. External → create `.md`/`.mdx` in `apps/docs/docs/` (add to an appropriate category subfolder).
3. Internal → create in root `docs/`.
4. If a new schema version is released, copy the updated schema to `apps/docs/static/schemas/sloth/<artifact>/<version>/schema.json` and update `apps/docs/docs/schemas.md` with the new URL.
5. Run `task build-docs` (or `pnpm --filter apps-docs build`) to verify no broken links before committing.

## Do Not

- Do not add tutorial, blog, or placeholder content from Docusaurus templates.
- Do not place public documentation as raw markdown in the repo root or `docs/`.
- Do not reference `docs/` internal files from `apps/docs/` pages (they are not served by Docusaurus).
