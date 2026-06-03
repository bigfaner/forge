---
title: "Surface CLI Conventions"
domains: [surface, cli, json, stderr, forge-surfaces, structured-output]
---

# Surface CLI Conventions

_Source: feature/surface-aware-justfile_

## JSON Output

### TECH-surface-cli-001: forge surfaces --json Flag Specification

**Requirement**: The `forge surfaces` CLI command supports a `--json` bool flag that outputs structured JSON for machine consumption by skills. Three sub-modes:
- **List mode** (no path, no --types): `{"surfaces": [{"key": "admin-panel", "type": "web"}, ...]}`
- **Query mode** (with path): `[{"key": "admin-panel", "type": "web"}]`; no match returns `[]` with exit 0
- **Types mode** (--types flag): `{"types": ["api", "web"]}`

When surfaces config is missing (all modes): stderr `{"error": "no surface configured; run \`forge init\` to configure surfaces"}`, exit 1.
**Scope**: [CROSS]
**Source**: feature/surface-aware-justfile TECH-001

### TECH-surface-cli-002: --json Mode stderr Format Override

**Requirement**: When `--json` flag is active, ALL output (stdout AND stderr) MUST use structured JSON format. This is an explicit exception to the plain-text stderr format convention (`<context>: <specific-detail>`). Rationale: `--json` consumers are machines (skills via Bash tool), requiring stderr to be parseable by JSON parsers without ambiguity. All error paths under `--json` must use `json.NewEncoder(cmd.ErrOrStderr()).Encode()` with `{"error": "..."}` format; `fmt.Fprintf(os.Stderr, ...)` is prohibited.
**Scope**: [CROSS]
**Source**: feature/surface-aware-justfile TECH-002
