---
title: "Surface Rule File Conventions"
domains: [surface, rules, recipe, scalar, named, orchestration, data-propagation, naming, text-mode, parsing, scaffold]
---

# Surface Rule File Conventions

_Source: feature/surface-aware-justfile_

## Scaffold Command

### TECH-surface-rules-001: Surface Recipe Generation via Scaffold CLI

**Requirement**: Surface recipe generation is handled by the `forge justfile scaffold` CLI command, not by surface rule files. The command generates placeholder-templated justfile recipes for a given surface type, outputting to stdout. Each surface type has a fixed recipe set: cli/tui produce test + teardown + quality recipes (compile/fmt/lint/unit-test); api/web produce dev + probe + test + teardown + orchestration + quality recipes; mobile produces test-setup + dev + probe + test + teardown + orchestration + quality recipes. The CLI command is the single source of truth for recipe generation — surface rule files (`rules/surfaces/<type>.md`) have been removed.
**Scope**: [CROSS]
**Source**: feature/init-justfile-slim TECH-003 (replaces feature/surface-aware-justfile TECH-003)

## Recipe Naming

### TECH-surface-rules-002: Recipe Naming Convention (Scalar vs Named Surfaces)

**Requirement**: Recipe naming is determined by surface form (scalar vs named), not by project size (single vs multi):
- **Scalar surface** (no key, only type): recipes use the verb directly — `test`, `dev`, `teardown`. No prefix.
- **Named surface** (has key): recipes use `<key>-<verb>` — e.g., `admin-panel-dev`, `payment-api-test`. Applies to both single and multi-surface projects.

Surface-key in recipe names must use `[a-zA-Z0-9_-]` characters only. Aggregate recipes without prefix serve as default entry points (e.g., `dev` starts all dev servers in dependency order). Orchestration order is recorded in justfile header comments for run-tests parsing.
**Scope**: [CROSS]
**Source**: feature/surface-aware-justfile TECH-004

## Data Propagation

### TECH-surface-rules-003: Cross-Layer Surface Data Propagation

**Requirement**: Surface information flows through a 7-step chain:
1. `.forge/config.yaml` surfaces field -> `forge surfaces` CLI (file read + CLI query)
2. `forge surfaces` CLI (text mode) -> breakdown-tasks/quick-tasks skill (text stdout parsed per-line)
3. Skill -> task frontmatter (YAML fields: surface-key, surface-type)
4. Frontmatter -> index.json Task Go struct (JSON serialization)
5. index.json -> run-tests skill (Go function call)
6. run-tests -> execution via `forge quality-gate` (surface-aware lifecycle orchestration)
7. `forge surfaces` CLI (text mode) -> init-justfile/test-guide skill (text stdout)

Text mode parsing rule (unified across all skills): per line, if line contains `=`, split into key (before `=`) and type (after `=`) — named surface; otherwise key is empty and type is the line — scalar surface.

Fallback chain: task frontmatter -> `forge surfaces <path>` longest-prefix-match -> error exit.
**Scope**: [CROSS]
**Source**: feature/surface-aware-justfile TECH-006
