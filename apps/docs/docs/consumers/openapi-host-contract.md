---
sidebar_position: 8
---

# OpenAPI Host Contract

Milestone 3 adds a CMS-agnostic OpenAPI contract for the host inspection and ingest surface used by CLI workflows.

## Spec Location

- Source file: `packages/contracts/openapi/sloth-api.openapi.yaml`

## Covered Endpoints

- `GET /healthz`
- `GET /sloth/inspection/plugin-status`
- `GET /sloth/inspection/contract-schema`
- `GET /sloth/contracts`
- `GET /sloth/contracts/{name}`
- `POST /sloth/contracts/ingest`

## Mock Server

A deterministic mock server implementation is available at:

- `packages/contracts/scripts/openapi-mock-server.mjs`

Run it:

```bash
node packages/contracts/scripts/openapi-mock-server.mjs --port 4010
```

The server prints a ready marker with the bound port and serves seeded contracts from `packages/contracts/src/mock/seed-contracts.json`.
