---
created: "2026-05-26"
author: "faner"
status: Draft
---

# Proposal: Unify `forge task add` — Replace `--template` with `--type`

## Problem

`forge task add` has two overlapping flags: `--template` (selects task template file + defaults) and `--type` (sets task classification). The template's defaults include a `Type` field, so `--template fix-task` implicitly sets `--type coding.fix`. Users must understand both concepts to use the command effectively.

### Evidence

- CLI help lists both `--template` and `--type` as independent flags, but they are coupled: template determines type.
- Only 2 templates exist (`fix-task`, `cleanup-task`), each maps 1:1 to a type value (`coding.fix`, `coding.cleanup`).
- Template filenames (`fix-task.md`) differ from type values (`coding.fix`), adding unnecessary indirection.

### Urgency

Forge v3.0.0 is in active development (current branch). This is the right time to clean up the CLI surface before release. Post-release breaking changes are more costly.

## Proposed Solution

Remove `--template` flag. Rename template files to match their type values (`coding.fix.md`, `coding.cleanup.md`). When `--type` is specified, the system checks if a matching template file exists; if so, it loads the template and its defaults. This eliminates the template/type duality — users only need `--type`.

### Innovation Highlights

Straightforward CLI simplification. The key insight is that template filenames can serve double duty as type identifiers, removing the need for a separate mapping layer.

## Requirements Analysis

### Key Scenarios

- **User specifies `--type coding.fix`**: System finds `coding.fix.md` template, loads content + defaults (priority=P0, breaking=true, etc.), creates task with template.
- **User specifies `--type coding.feature`**: No matching template file. System creates a regular task with the type field set, no template content.
- **User specifies `--type` without value**: Error, same as current behavior.
- **Quality gate creates fix task**: Internal caller uses type value `coding.fix` instead of template name `fix-task`.

### Non-Functional Requirements

- **Backward compatibility**: Breaking change, acceptable in v3.0.0 pre-release.

### Constraints & Dependencies

- Templates are embedded at compile time via `//go:embed`. File renames require rebuild.
- `quality_gate.go` calls `addFixTask()` programmatically with template names — must be updated to use type values.

## Alternatives & Industry Benchmarking

### Industry Solutions

Most CLI tools use a single `--type` or `--kind` flag to classify and template items (e.g., `kubectl create --type`, `npm init --type`). Separate template/type flags are uncommon.

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | Zero cost | UX confusion persists, two flags for one concept | Rejected: v3 is the window for breaking changes |
| Keep `--template` hidden, add `--type` lookup | Backward compat | Old scripts still work | Two code paths to maintain | Rejected: adds complexity without benefit |
| **Rename files, remove `--template`** | Industry standard (single type flag) | Clean mental model, one flag does both | Breaking change | **Selected: v3.0.0 allows breaking changes** |

## Feasibility Assessment

### Technical Feasibility

All changes are within the Go CLI codebase. No external dependencies. Template files are embedded, so renaming is a compile-time concern only.

### Resource & Timeline

Small scope: ~10 files, mostly renaming and updating string references. Estimated 1 coding task.

### Dependency Readiness

No external dependencies. Self-contained change.

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| Users need both `--template` and `--type` | Occam's Razor | Overturned: they serve the same purpose for 100% of current templates |
| Template names should be human-friendly short names | XY Detection | Refined: type values (`coding.fix`) are equally readable and more descriptive |

## Scope

### In Scope

- Rename template files: `fix-task.md` → `coding.fix.md`, `cleanup-task.md` → `coding.cleanup.md`
- Remove `--template` flag from `forge task add`
- Update `template.go` defaults map keys to use type values
- Update `--type` to auto-discover matching template file
- Update `quality_gate.go` internal caller to use type values
- Update all related tests

### Out of Scope

- Adding new template types
- Changing template file content
- Changing `--var` behavior
- Adding template discovery from user project directories

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Other code references template names as strings | M | M | Grep all occurrences of `fix-task` and `cleanup-task` across codebase |
| Tests hardcode template names | H | L | Update test fixtures alongside code changes |

## Success Criteria

- [ ] `forge task add --type coding.fix --title "Fix X"` loads `coding.fix.md` template and applies its defaults
- [ ] `forge task add --type coding.feature --title "Build Y"` works without template (no matching file)
- [ ] `forge task add --template fix-task` returns an error (flag removed)
- [ ] `forge task add -h` shows no `--template` flag
- [ ] Quality gate auto-created fix tasks use type value `coding.fix` instead of template name
- [ ] All existing tests pass after rename

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
