# v2-to-v3 Config Migration Guide

How to migrate `.forge/config.yaml` from Forge v2 (test-profiles + capabilities) to Forge v3 (project-type + interfaces).

## Overview

Forge v3 removes the Profile concept entirely. Language is auto-detected from project files. User config shrinks from 3 fields to 2.

**Migration effort**: 2-line manual edit (delete `test-profiles`, rename `capabilities` to `interfaces`). No automated migration tooling.

## Field Mapping

| v2 Field | v3 Replacement | Notes |
|----------|---------------|-------|
| `project-type` | `project-type` (unchanged) | Scope resolution, values: backend, frontend, mixed |
| `test-profiles` | *(removed)* | Language auto-detected from project files. Override via `languages` if needed. |
| `capabilities` | `interfaces` | Same closed enum: web-ui, tui, mobile-ui, api, cli. Renamed for clarity. |
| *(new)* | `languages` | Optional override when auto-detection fails or returns wrong result |

## Profile-to-Language Mapping

v2 profile names are removed. The following table maps each v2 profile to the v3 language that replaces it (auto-detected, no config needed):

| v2 Profile | v3 Language Key | Detection Signal |
|------------|----------------|-----------------|
| `go-test` | `go` | `go.mod` exists at project root |
| `web-playwright` | `javascript` | `package.json` with `@playwright/test` in devDependencies |
| `rust-test` | `rust` | `Cargo.toml` exists at project root |
| `pytest` | `python` | `pyproject.toml` or `requirements.txt` containing "pytest" |
| `java-junit` | `java` | `pom.xml` or `build.gradle` exists at project root |
| `maestro` | `mobile` | `android/` or `ios/` directory exists at project root |

## Before/After Examples

### Single-language Go backend

```yaml
# v2
project-type: backend
test-profiles: [go-test]
capabilities: [api, cli]
```

```yaml
# v3
project-type: backend
interfaces: [api, cli]
# languages: auto-detected from go.mod
```

### Frontend with Playwright

```yaml
# v2
project-type: frontend
test-profiles: [web-playwright]
capabilities: [web-ui, api]
```

```yaml
# v3
project-type: frontend
interfaces: [web-ui, api]
# languages: auto-detected from package.json + @playwright/test
```

### Multi-language project (Go backend + JS frontend)

```yaml
# v2
project-type: mixed
test-profiles: [go-test, web-playwright]
capabilities: [api, cli, web-ui]
```

```yaml
# v3
project-type: mixed
interfaces: [api, cli, web-ui]
# languages: both auto-detected (go.mod + package.json)
```

### Python CLI project

```yaml
# v2
project-type: backend
test-profiles: [pytest]
capabilities: [api, cli]
```

```yaml
# v3
project-type: backend
interfaces: [api, cli]
# languages: auto-detected from pyproject.toml or requirements.txt
```

### Mobile project (Maestro)

```yaml
# v2
project-type: frontend
test-profiles: [maestro]
capabilities: [mobile-ui]
```

```yaml
# v3
project-type: frontend
interfaces: [mobile-ui]
# languages: auto-detected from android/ or ios/ directory
```

### Minimal config (language auto-detected, all interfaces)

```yaml
# v2
project-type: backend
test-profiles: [go-test]
```

```yaml
# v3
project-type: backend
# Both languages and interfaces omitted: language detected from project files,
# interfaces default to all supported types for that language.
# For Go: api, cli
```

## Common Override Patterns

### Multi-language false positive

A Go project with a `package.json` for lint tooling (ESLint, Prettier) may detect both `go` and `javascript`. Suppress the false positive:

```yaml
project-type: backend
interfaces: [api, cli]
languages: [go]              # override: suppress javascript detection
```

Detection is not suppressed globally -- `forge testing detect` still reports all signals. The `languages` field only controls which strategies are activated for test generation and execution.

### Monorepo with language-specific subdirectories

Detection scans the project root only. If a Go service lives in `services/api/` and the root has no `go.mod`, the language is not detected. Specify manually:

```yaml
project-type: backend
interfaces: [api]
languages: [go]              # override: no go.mod at project root
```

This is a known limitation. Subdirectory detection is not supported in v3.0.

### Narrowing interfaces

A Go project supports `api` and `cli` by default. To test only API:

```yaml
project-type: backend
interfaces: [api]            # only API tests, skip CLI
# languages: auto-detected
```

When `interfaces` is omitted, all supported interface types for the detected language are used. When specified, only the listed types are tested.

## Troubleshooting

### No language detected

**Symptom**: `forge testing detect` returns empty output. `forge testing get generate` exits with non-zero code and stderr contains "languages".

**Cause**: No detection signal found at the project root (no `go.mod`, `package.json`, `Cargo.toml`, etc.).

**Fix**: Add `languages` to config.yaml:

```yaml
project-type: backend
interfaces: [api, cli]
languages: [go]              # explicit language selection
```

### Wrong language detected

**Symptom**: `forge testing detect` reports a language you do not expect (e.g., `javascript` in a Go project with a `package.json` for tooling).

**Fix**: Override with `languages` field (see "Multi-language false positive" above).

### Old v2 config still present

**Symptom**: After upgrading to Forge v3, the CLI ignores `test-profiles` and `capabilities` fields.

**Cause**: v3 does not read v2 fields. They are silently ignored.

**Fix**: Remove `test-profiles` and `capabilities` from config.yaml. Add `interfaces` (renamed from `capabilities`). Optionally add `languages` if auto-detection does not work for your project.

### `interfaces` values not recognized

**Symptom**: Error message about invalid interface type.

**Cause**: Valid values are: `web-ui`, `tui`, `mobile-ui`, `api`, `cli` (same closed enum as v2 capabilities).

**Fix**: Check spelling. Values are lowercase with hyphens. No new values were added in v3.

## CLI Command Changes

| v2 Command | v3 Command |
|-----------|-----------|
| `forge profile detect` | `forge testing detect` |
| `forge profile get <name> --generate` | `forge testing get generate` |
| `forge profile get <name> --run` | `forge testing get run` |
| `forge profile get <name> --graduate` | `forge testing get graduate` |
| `forge profile get <name> --justfile` | `forge testing get justfile` |
| `forge profile get <name> --template <file>` | `forge testing get template <file>` |
| *(none)* | `forge testing interfaces` |

Key difference: v3 commands do not take a profile name argument. Language is auto-detected or specified via `--language` flag for multi-language projects.
