# Website

This website is built using [Docusaurus](https://docusaurus.io/), a modern static website generator.

## Documentation Groups

- Consumer-facing pages live in `apps/docs/docs/consumers/`.
- Repository-maintainer pages live in `apps/docs/docs/repo-developers/`.
- Keep cross-links explicit when consumer pages reference maintainer workflows.

## Schema Sync Note

Before building docs for release, sync promoted schema versions from GHCR artifacts into:

`apps/docs/static/schemas/<artifact>/<version>/schema.json`

Then build Docusaurus output so canonical `$schema` URLs resolve from docs static hosting.

## Installation

```bash
yarn
```

## Local Development

```bash
yarn start
```

This command starts a local development server and opens up a browser window. Most changes are reflected live without having to restart the server.

## Build

```bash
yarn build
```

This command generates static content into the `build` directory and can be served using any static contents hosting service.

## Deployment

Using SSH:

```bash
USE_SSH=true yarn deploy
```

Not using SSH:

```bash
GIT_USER=<Your GitHub username> yarn deploy
```

If you are using GitHub pages for hosting, this command is a convenient way to build the website and push to the `gh-pages` branch.
