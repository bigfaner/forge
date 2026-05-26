---
title: "Surface Rule File Conventions"
domains: [surface, rules, recipe, mixed-project, orchestration, data-propagation, naming]
---

# Surface Rule File Conventions

_Source: feature/surface-aware-justfile_

## Rule File Format

### TECH-surface-rules-001: Surface Rule File Format Convention

**Requirement**: Surface rule files follow the path pattern `rules/surfaces/<type>.md` and are consumed by both init-justfile (recipe generation) and run-tests (orchestration). Each file contains: (1) orchestration sequence table with exit code handling per step, (2) recipe invocation contract table defining just signature, exit code semantics, and (3) journey filter strategy mapping @tag to surface type. Files are loaded dynamically based on surface-type detected at runtime.
**Scope**: [CROSS]
**Source**: feature/surface-aware-justfile TECH-003

## Recipe Naming

### TECH-surface-rules-002: Recipe Naming Convention for Mixed Projects

**Requirement**: In mixed projects (multiple surfaces), recipes use `<action>-<surface-key>` naming pattern (e.g., `dev-admin-panel`, `probe-payment-service`). Aggregate recipes without prefix serve as default entry points (e.g., `dev` starts all dev servers in dependency order). Surface-key in recipe names must use `[a-zA-Z0-9_-]` characters only. Orchestration order is recorded in justfile header comments for run-tests parsing.
**Scope**: [CROSS]
**Source**: feature/surface-aware-justfile TECH-004

## Data Propagation

### TECH-surface-rules-003: Cross-Layer Surface Data Propagation

**Requirement**: Surface information flows through a 7-step chain:
1. `.forge/config.yaml` surfaces field -> `forge surfaces` CLI (file read + CLI query)
2. `forge surfaces` CLI -> breakdown-tasks/quick-tasks skill (JSON stdout)
3. Skill -> task frontmatter (YAML fields: surface-key, surface-type)
4. Frontmatter -> index.json Task Go struct (JSON serialization)
5. index.json -> run-tests skill (Go function call)
6. run-tests -> execution strategy rule file (file path: `rules/surfaces/{type}.md`)
7. `forge surfaces` CLI -> init-justfile skill (CLI stdout)

Fallback chain: task frontmatter -> `forge surfaces <path>` longest-prefix-match -> error exit.
**Scope**: [CROSS]
**Source**: feature/surface-aware-justfile TECH-006
