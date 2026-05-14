---
id: "4"
title: "Migration: ResolveScope → config.yaml and doc updates"
priority: "P1"
estimated_time: "1h"
dependencies: ["1"]
scope: "all"
breaking: true
type: "implementation"
mainSession: false
---

# 4: Migration: ResolveScope → config.yaml and doc updates

## Description

Migrate scope resolution from `just project-type` subprocess calls to reading `.forge/config.yaml` directly via `forge config get project-type`. Update `ResolveScope()` in Go code, all skill/hook docs that reference `just project-type`, and remove the `project-type` recipe from justfile and templates.

## Reference Files
- `docs/proposals/forge-info-commands/proposal.md` — Source proposal

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/just/just.go` | Rewrite `ResolveScope()` to read config.yaml directly instead of calling `just project-type` subprocess |
| `forge-cli/pkg/just/just_test.go` | Update tests to reflect new ResolveScope behavior |
| `plugins/forge/hooks/guide.md` | Replace `just project-type` with `forge config get project-type` |
| `plugins/forge/commands/fix-bug.md` | Replace `just project-type` with `forge config get project-type` |
| `plugins/forge/commands/execute-task.md` | Replace `just project-type` with `forge config get project-type` |
| `plugins/forge/commands/run-tasks.md` | Replace `just project-type` with `forge config get project-type` |
| `justfile` | Remove `project-type:` recipe |
| `forge-cli/internal/cmd/root.go` | No change needed (already registered) |

### Delete
| File | Reason |
|------|--------|
| `project-type` recipe in justfile | Replaced by config.yaml |
| `project-type` recipe in skill templates | Replaced by config.yaml |

## Acceptance Criteria

- [ ] `ResolveScope()` reads `project-type` from `.forge/config.yaml` via profile package, no subprocess call
- [ ] All skill/hook docs use `forge config get project-type` instead of `just project-type`
- [ ] `justfile` has no `project-type:` recipe
- [ ] Skill justfile templates have no `project-type:` recipe
- [ ] Existing `ResolveScope()` tests updated and passing
- [ ] Scope resolution behavior is identical: mixed → pass scope, frontend/backend → skip scope, missing → skip scope
- [ ] Test coverage ≥ 80% for modified code

## Hard Rules

- `ResolveScope()` must NOT spawn any subprocess — direct config file read only
- Must maintain backward-compatible behavior: if config.yaml doesn't exist or `project-type` is missing, return empty string (skip scope)

## Implementation Notes

- The `ResolveScope()` rewrite should import and use `profile.ReadConfig(projectRoot)` or similar to get `ProjectType` field
- Remove `HasRecipe` call for `project-type` — no longer relevant
- For skill docs: search all `.md` files under `plugins/forge/` for `just project-type` and replace with `forge config get project-type`
- Also check for patterns like `Run ... "just", "project-type"` or `just project-type` in command strings
- For justfile: locate and remove the `project-type:` recipe (likely `@echo "backend"` or similar)
- For template justfiles in `plugins/forge/skills/init-justfile/templates/`: remove `project-type:` recipe from each template
