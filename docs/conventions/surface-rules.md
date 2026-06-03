---
title: "Surface Rule File Conventions"
domains: [surface, rules, recipe, scalar, named, orchestration, data-propagation, naming, text-mode, parsing]
---

# Surface Rule File Conventions

_Source: feature/surface-aware-justfile_

## Rule File Format

### TECH-surface-rules-001: Surface Rule File Format Convention

**Requirement**: Surface rule files follow the path pattern `rules/surfaces/<type>.md` and are consumed by both init-justfile (recipe generation) and run-tests (orchestration). Each file contains: (1) orchestration sequence table with exit code handling per step, (2) recipe invocation contract table defining just signature, exit code semantics, and (3) journey filter strategy mapping @tag to surface type. Files are loaded dynamically based on surface-type detected at runtime.
**Scope**: [CROSS]
**Source**: feature/surface-aware-justfile TECH-003

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
6. run-tests -> execution strategy rule file (file path: `rules/surfaces/{type}.md`)
7. `forge surfaces` CLI (text mode) -> init-justfile/test-guide skill (text stdout)

Text mode parsing rule (unified across all skills): per line, if line contains `=`, split into key (before `=`) and type (after `=`) — named surface; otherwise key is empty and type is the line — scalar surface.

Fallback chain: task frontmatter -> `forge surfaces <path>` longest-prefix-match -> error exit.
**Scope**: [CROSS]
**Source**: feature/surface-aware-justfile TECH-006
