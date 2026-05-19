---
created: 2026-05-18
author: "faner"
status: Approved
---

# Proposal: Remove `plugins/forge/references/` Directory

## Problem

`plugins/forge/references/` is a non-standard Claude Code plugin directory that introduces fragile cross-component path traversal (`${CLAUDE_SKILL_DIR}/../../references/shared/`) and hidden dependency coupling across 15 skill/command reference sites.

### Evidence

- Claude Code official plugin spec defines 6 standard directories (commands, agents, skills, hooks, scripts, .mcp.json/.lsp.json). `references/` is not among them.
- The official plugin reference document explicitly warns against `../../` relative paths in skill/command files.
- `profile-detection.md` duplicates detection logic that also exists in Go CLI code (`forge-cli/pkg/profile/`), creating drift risk.
- `forge-config.schema.json` and `forge-config.example.yaml` serve double duty: referenced at plugin runtime AND read by CLI tests via hardcoded source-tree relative path.

### Urgency

The v3.0.0 release is the right time to clean this up. The directory adds unnecessary complexity and creates a maintenance burden that grows with each new shared reference file. Removing it now prevents path traversal issues from accumulating.

## Proposed Solution

Remove `plugins/forge/references/` entirely. Inline all shared protocol content directly into the 15 referencing skill/command files. Move CLI-specific files (schema, example YAML) to the CLI repository. Update `forge-distribution.md` to reflect the removal.

### Innovation Highlights

Straightforward simplification. The key insight is that Claude Code plugins should be self-contained — cross-skill shared state via file references is an anti-pattern that the standard plugin model deliberately avoids. Each skill owning its own protocol copy is the correct plugin-native approach.

## Requirements Analysis

### Key Scenarios

- Skills that reference `decision-logging.md` (consolidate-specs, tech-design, learn) must contain the full decision logging protocol inline
- Skills that reference `knowledge-extraction.md` (write-prd, tech-design) and commands (run-tasks, fix-bug) must contain the full knowledge extraction routine inline
- Skills that reference `type-assignment.md` and `intent-propagation.md` (breakdown-tasks, quick-tasks) must contain the full type mapping and intent propagation logic inline
- Skills that reference `step0-profile-resolution.md` (breakdown-tasks, quick-tasks) must contain the profile resolution rules inline
- The `gen-sitemap` command must have config.yaml and sitemap.json examples inline
- CLI test must read `forge-config.schema.json` and `forge-config.example.yaml` from the CLI repo, not from the plugin source tree

### Non-Functional Requirements

- Zero runtime behavior change — inline is a structural refactor, not a functional change
- All existing path references (`${CLAUDE_SKILL_DIR}/../../references/shared/...`) must be completely removed

### Constraints & Dependencies

- Must update `docs/conventions/forge-distribution.md` to remove references/ documentation
- CLI test files (`config_schema_test.go`) must be updated to read schema from new location

## Alternatives & Industry Benchmarking

### Industry Solutions

Claude Code plugins are designed as self-contained skill packages. The standard model discourages shared state between skills — each skill is an independent unit with its own SKILL.md, optional scripts/, and optional templates/.

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | Zero cost | Path fragility persists, non-standard pattern accumulates | Rejected: cost of delay exceeds cost of fix |
| Gradual deprecation | — | Lower risk per change | Extended transition period with dual patterns | Rejected: only 15 reference sites, one-shot is cleaner |
| Merge to skill templates/ | Standard plugin pattern | Files remain accessible | Still cross-directory dependency, doesn't solve root problem | Rejected: same fragility, different location |
| **Full inline into each skill** | Plugin-native self-containment | Zero cross-skill deps, standard-compliant, future-proof | Content duplication across skills | **Selected: duplication is acceptable for independence** |

## Feasibility Assessment

### Technical Feasibility

Fully feasible. All reference files are plain markdown/JSON/YAML read at runtime via the Read tool. Inlining is a copy-paste operation with no code logic changes.

### Resource & Timeline

Small scope: 10 files to inline across 15 reference sites, plus documentation and CLI test updates. Single-pass refactor.

### Dependency Readiness

No external dependencies. All changes are internal to the forge plugin.

## Scope

### In Scope

- Remove `plugins/forge/references/` directory and all contents
- Inline `decision-logging.md` protocol into consolidate-specs, tech-design, learn skills
- Inline `knowledge-extraction.md` routine into write-prd, tech-design skills and run-tasks, fix-bug commands
- Inline `type-assignment.md` and `intent-propagation.md` into breakdown-tasks, quick-tasks skills
- Inline `step0-profile-resolution.md` into breakdown-tasks, quick-tasks skills
- Inline `config.yaml` and `sitemap.json` examples into gen-sitemap command
- Move `forge-config.schema.json` and `forge-config.example.yaml` to CLI repo (or appropriate location)
- Update `docs/conventions/forge-distribution.md` to remove references/ section
- Update `profile-detection.md` consumers (if any) — note: this file duplicates Go CLI logic, evaluate if removal is sufficient

### Out of Scope

- Other non-standard directories or plugin structure changes
- Changes to the Go CLI profile detection logic
- Refactoring skill internal structure beyond inline changes

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Content drift across inlined copies | M | M | Each skill owns its own copy — when protocol changes, update the skill that needs it. Acceptable tradeoff for independence. |
| Missing a reference site during migration | L | H | Comprehensive grep for `references/shared/` across entire plugin. Verify with `forge eval-forge` after completion. |
| CLI test breakage from schema file move | L | M | Update test file path and verify test passes before merge. |
| Breaking existing user projects | L | H | This is a plugin-internal refactor. No public API changes. Users won't be affected. |

## Success Criteria

- [ ] `plugins/forge/references/` directory no longer exists
- [ ] Zero occurrences of `references/shared/` path in any plugin file
- [ ] All 15 former reference sites contain the inlined content
- [ ] CLI schema test passes with new file location
- [ ] `forge-distribution.md` no longer documents references/ directory
- [ ] `forge eval-forge` audit passes after changes

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
