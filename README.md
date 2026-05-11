# sloth

A page-template platform built on top of Strapi. sloth provides a Strapi plugin, a developer CLI, and a component hub for defining, distributing, and rendering page templates in frontend applications.

> **Status:** Early development. Core plugin and CLI scaffolding are in progress.

## Repository Structure

```
apps/
  docs/         # Docusaurus documentation site + schema hosting (GitHub Pages)
  frontend/     # Frontend application (planned)
  strapi/       # Strapi CMS host application

packages/
  cli/          # sloth CLI (Go + Cobra) — contract operations
  contracts/    # Contract source, release tooling, and OCI artifacts
  component-hub/# Placeholder for future frontend component/runtime package
  strapi-plugin/# Core sloth Strapi plugin

docs/           # Internal planning docs — architecture, ADRs, ideas
packages/contracts/src/schemas/ # Canonical JSON Schema source files
```

## Prerequisites

| Tool    | Version | Install                             |
| ------- | ------- | ----------------------------------- |
| Node.js | ≥ 20    | [nodejs.org](https://nodejs.org)    |
| pnpm    | 9.x     | `npm i -g pnpm`                     |
| go-task | latest  | `brew install go-task`              |
| Docker  | latest  | [docker.com](https://docker.com)    |
| Go      | ≥ 1.22  | [go.dev](https://go.dev) (CLI only) |

## Getting Started

Install all dependencies from the repo root:

```bash
pnpm install
```

## Available Tasks

View all available task commands:

```bash
task --list
```

Key tasks:

| Task                    | Description                                |
| ----------------------- | ------------------------------------------ |
| `task start-strapi-dev` | Start the Strapi development server        |
| `task build-plugin`     | Build the strapi-plugin package            |
| `task watch-plugin`     | Watch and rebuild the plugin on changes    |
| `task start-docs`       | Start the docs site dev server (port 3000) |
| `task build-docs`       | Build the docs site for production         |
| `task start-dev-all`    | Start all services (Docker + Strapi)       |
| `task stop-dev-all`     | Stop all running services                  |

## Plugin Development

The Strapi plugin lives in `packages/strapi-plugin/`. It is linked into `apps/strapi/` via pnpm workspace.

To work on the plugin while running Strapi:

```bash
# Terminal 1 — watch plugin for changes
task watch-plugin

# Terminal 2 — run Strapi
task start-strapi-dev
```

Type-check without emitting JS:

```bash
cd packages/strapi-plugin
pnpm run test:ts:back   # server TypeScript check
pnpm run test:ts:front  # admin TypeScript check
```

## Documentation

Public documentation is served from `apps/docs/` via GitHub Pages:

- Live site: [https://phuhh98.github.io/sloth/](https://phuhh98.github.io/sloth/)
- Local: `task start-docs`

Internal planning documents are in `docs/` (not public-facing).

## Resources

- [go-task](https://taskfile.dev/)
- [Strapi plugin development](https://docs.strapi.io/cms/plugins-development/create-a-plugin)
- [Strapi Document Service API](https://docs.strapi.io/cms/api/document-service)
