---
id: "2"
title: "Extract phase-detection, db-schema, and existing-code-split rule files"
priority: "P1"
estimated_time: "1h"
dependencies: []
type: "doc"
mainSession: false
---

# 2: Extract phase-detection, db-schema, and existing-code-split rule files

## Description

Extract three smaller conditional rule files from the current `SKILL.md`. Each file handles a specific conditional concern that currently uses inline tag blocks:
- `rules/phase-detection.md` (~2KB) — phase detection from PRD/design, currently in Step 2 (lines 133-161)
- `rules/db-schema.md` (~1KB) — DB schema task creation, currently in `<HAS_DB>` block (lines 232-241)
- `rules/existing-code-split.md` (~1.5KB) — shared code modification split, currently in Step 4a (lines 222-231)

## Reference Files
- `docs/proposals/simplify-breakdown-tasks-prompt/proposal.md` — Source proposal (Rule File Extraction Plan)
- `plugins/forge/skills/breakdown-tasks/SKILL.md` — Current skill file (source of extracted content)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/breakdown-tasks/rules/phase-detection.md` | Phase detection rule file (~2KB) |
| `plugins/forge/skills/breakdown-tasks/rules/db-schema.md` | DB schema rule file (~1KB) |
| `plugins/forge/skills/breakdown-tasks/rules/existing-code-split.md` | Existing code modification split rule file (~1.5KB) |

### Modify
| File | Changes |
|------|---------|
| (none) | SKILL.md is modified in task 3 |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria

**rules/phase-detection.md:**
- [ ] Contains three-tier detection: explicit (highest priority) → heuristic → fallback
- [ ] Contains explicit detection patterns: flow diagram diamond nodes, PRD sections named "Round/Phase/Stage 1/2"
- [ ] Contains heuristic detection patterns: sequential markers, conditional transitions, go/no-go checkpoints, gated prose (both English and Chinese patterns)
- [ ] Contains fallback: artifact-driven decomposition when no phases detected
- [ ] Contains `phase-inventory.json` format specification
- [ ] Load condition at top: "Load IF PRD has phase/gate structure (detected by explicit or heuristic patterns)"
- [ ] Maintenance note listing skeleton dependencies (Step 2, Step 3)

**rules/db-schema.md:**
- [ ] Contains schema task creation rules (one task per entity in er-diagram.md)
- [ ] Contains acceptance criteria: "DDL executes without error", "all FK references resolve", "indexes created"
- [ ] Contains breaking classification: ALTER existing table → `breaking: true`, all CREATE TABLE new → `breaking: false`
- [ ] Contains scope assignment: `scope: "backend"`
- [ ] Contains dependency rule: depends on interface tasks (migration may need type information)
- [ ] Load condition at top: "Load IF `design/er-diagram.md` exists"
- [ ] Maintenance note listing skeleton dependencies (Step 2 element mapping, Step 4a)

**rules/existing-code-split.md:**
- [ ] Contains artifact-update + feature sub-task split procedure
- [ ] Contains sub-ID convention: `<seq>.<sub>a` for shared artifact update, `<seq>.<sub>b` for feature implementation
- [ ] Contains when-to-apply threshold: >5 downstream files OR spans multiple architectural layers
- [ ] Contains `breaking: true` requirement for shared artifact update sub-task
- [ ] Contains exclusion: purely additive new code does not need splitting
- [ ] Load condition at top: "Load IF tech-design references modifications to existing shared code"
- [ ] Maintenance note listing skeleton dependencies (Step 4a, Step 5)

**All files:**
- [ ] Each file is independently understandable
- [ ] Each file uses skill-relative paths (compatible with forge distribution model)
- [ ] Each file includes guard clause for malformed/empty artifacts

## Hard Rules

- All file paths in rule files must use skill-relative references (compatible with forge distribution model under `${CLAUDE_SKILL_DIR}`)
- Do NOT modify `SKILL.md` — that is task 3

## Implementation Notes

These three files are simpler than `ui-placement.md` because each has a single conditional tag (or no tag at all, in the case of existing-code-split which is always present but only relevant for certain features). The `existing-code-split` content is currently always in SKILL.md (not gated by a tag), but it only applies when the tech-design modifies existing shared code — making it a good candidate for conditional loading.
