---
id: "6"
title: "Update gen-test-scripts input discovery + convention loading"
priority: "P1"
estimated_time: "1.5h"
dependencies: ["1"]
type: "documentation"
mainSession: false
noTest: true
---

# 6: Update gen-test-scripts input discovery + convention loading

## Description

Update `gen-test-scripts` SKILL.md to accept per-type test case files as input (in addition to the legacy single `test-cases.md`). Add convention loading that reads per-type instruction frontmatter `conventions` field and loads matching files from `docs/conventions/`.

## Reference Files
- `docs/proposals/gen-test-cases-per-type-dispatch/proposal.md` — Source proposal
- `plugins/forge/skills/gen-test-scripts/SKILL.md` — Current skill to update

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none) |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-test-scripts/SKILL.md` | Update prerequisite check to accept per-type files OR legacy format; update Step 1 input discovery with glob fallback; update Step Actionability gate for per-type file paths; add convention loading after profile resolution |

### Delete
| File | Reason |
|------|--------|
| (none) |

## Acceptance Criteria
- [ ] Prerequisite check: glob `testing/*-test-cases.md` first; if found, accept per-type mode; if empty, fall back to `testing/test-cases.md` for legacy mode
- [ ] Step 1 Read Test Cases: when reading per-type files, skip the type grouping step (file is already single-type)
- [ ] `--type` filter: when using per-type files, `--type` selects which `{type}-test-cases.md` to read; when using legacy, behavior unchanged
- [ ] Step Actionability gate: check eval reports for the specific per-type file being processed (e.g., `testing/eval/` directory for `ui-test-cases.md`)
- [ ] Convention loading: after profile resolution, read the active type's instruction file frontmatter (from `types/{type}.md` in gen-test-cases), extract `conventions` field, load existing files from `docs/conventions/`, skip missing silently
- [ ] gen-test-scripts frontmatter includes `conventions: [testing-isolation.md]` for project-wide conventions

## Hard Rules
- Legacy `testing/test-cases.md` must continue to work without changes — no breaking behavior
- When both per-type files AND legacy `test-cases.md` exist, prefer per-type files (new takes precedence)
- Convention loading must be non-blocking — missing convention files are silently skipped

## Implementation Notes
- The discovery logic is simple: `ls testing/*-test-cases.md` → if files found, use per-type mode; else fall back to `testing/test-cases.md`
- For per-type mode with `--type` filter, only load the matching `{type}-test-cases.md` file
- For per-type mode without `--type`, load all `{type}-test-cases.md` files sequentially (same as current behavior but from separate files)
- Convention loading reuses the same frontmatter-reading pattern from the gen-test-cases dispatcher
