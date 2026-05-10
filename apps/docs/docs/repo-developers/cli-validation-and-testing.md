---
title: CLI Validation and Testing
---

# CLI Validation and Testing

This page explains how CLI correctness is verified without requiring a permanently running Strapi host.

## Current Automated Coverage

The CLI includes unit and integration tests in `packages/cli`:

- unit tests for compatibility and verification logic
- unit tests for config and lock/workspace file behavior
- integration tests for command execution with mocked HTTP host responses

### Mock Host Integration

Command-level integration tests use Go `httptest` to emulate host endpoints:

- `GET /sloth/inspection/plugin-status`
- `GET /sloth/inspection/contract-schema`

This validates CLI behavior for inspect/list flows without external network dependencies.

## What Is Not Fully Proven by Mocks

Mock tests cannot fully guarantee production behavior for:

- auth and token policy differences
- payload shape drift between plugin and CLI
- runtime latency/timeouts in real environments
- ingest behavior under host-specific validation rules

## Recommended Local Verification

Run the CLI test suite:

```bash
task cli-test
```

Run distribution smoke checks:

```bash
task cli-release-prep
```

Run host integration checks against a local Strapi host:

1. Start host services.
2. Run `sloth contracts inspect`.
3. Run `sloth contracts push --dry-run`.
4. Run `sloth contracts push` with test payloads.

## Recommended Contract Hardening

To reduce risk before release:

- add contract tests for host client request and response shapes
- test non-200 status handling and retry behavior for push
- maintain a versioned API contract document alongside plugin route changes

## Related

- [CLI Command Reference](../consumers/cli-contract)
- [CLI Getting Started](../consumers/cli-getting-started)
