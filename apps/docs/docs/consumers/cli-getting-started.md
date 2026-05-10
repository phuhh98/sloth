---
title: CLI Getting Started
---

# sloth CLI Getting Started

The sloth CLI manages component contract workflows.

## Scope

The CLI is focused on component contracts only.

- in scope: list, inspect, add, verify, push component contracts
- out of scope: page content authoring and runtime page editing

## Initialize Workspace

Run init from your project root:

```bash
sloth init
```

This creates a local workspace:

```text
.sloth/
  config.yaml
  contracts/
  sets/
  manifests/
    lock.json
```

## Configure Host Profiles

The CLI resolves configuration with explicit precedence:

1. **YAML config** from `.sloth/config.yaml` (highest priority)
2. **Environment variables** (when YAML values missing)
3. **Built-in defaults** (fallback)
4. **Runtime flags** override all (explicit override only)

### YAML Configuration

Set host URL and token in `.sloth/config.yaml`:

```yaml
# yaml-language-server: $schema=https://phuhh98.github.io/sloth/schemas/cli-config/0.0.1/schema.json
currentProfile: default
profiles:
  default:
    host: http://localhost:1337
    authorizationToken: ""
  production:
    host: https://api.production.example.com
    authorizationToken: ""
```

Schema reference URL:

- [https://phuhh98.github.io/sloth/schemas/cli-config/0.0.1/schema.json](https://phuhh98.github.io/sloth/schemas/cli-config/0.0.1/schema.json)

Use profile selection:

```bash
sloth contracts inspect --profile default
sloth contracts push --profile production
```

### Environment Variable Fallback

If YAML config is missing or incomplete, the CLI checks these environment variables:

- `SLOTH_CONFIG`: path to config file (default: `.sloth/config.yaml`)
- `SLOTH_PROFILE`: profile name (default: `default`)
- `SLOTH_HOST`: host URL (default: `http://localhost:1337`)
- `SLOTH_AUTHORIZATION_TOKEN` or `SLOTH_TOKEN`: authorization token

Example in CI/CD:

```bash
export SLOTH_HOST=https://api.production.example.com
export SLOTH_AUTHORIZATION_TOKEN=secret-token-xyz
sloth contracts push --version 0.0.1
```

### Runtime Flag Override

Use flags to override YAML and ENV values explicitly:

```bash
sloth contracts inspect --host https://custom-host:1337 --authorization-token custom-token
```

Supported flags:

- `--config`: config file path
- `--profile`: profile name
- `--host -H`: host URL override
- `--authorization-token -T`: token override

## Typical Workflow

1. Inspect host status.

```bash
sloth contracts inspect --format table
```

2. List available contracts for a plugin version.

```bash
sloth contracts ls --version 0.0.1 --format table
```

3. Add contracts locally.

```bash
sloth contracts add --all --version 0.0.1
```

4. Verify local contract files.

```bash
sloth contracts verify --file .sloth/contracts/hero-banner@0.0.1.json --version 0.0.1
```

5. Push verified contracts.

```bash
sloth contracts push --version 0.0.1 --dry-run
sloth contracts push --version 0.0.1
```

## Next

- See [CLI Command Reference](./cli-contract) for all commands and flags.
- See [CLI Validation and Testing](../repo-developers/cli-validation-and-testing) for mock coverage and host verification strategy.
- See [CLI Distribution and Release](../repo-developers/cli-distribution-and-release) for cross-platform build and npm packaging.
