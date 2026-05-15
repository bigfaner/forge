---
id: "1"
title: "Add drift detection + auto-fix steps to consolidate-specs SKILL.md"
priority: "P0"
estimated_time: "1-2h"
dependencies: []
type: "documentation"
mainSession: false
---

# 1: Add drift detection + auto-fix steps to consolidate-specs SKILL.md

## Description

Extend `plugins/forge/skills/consolidate-specs/SKILL.md` to add spec drift detection and auto-fix capabilities. After the existing "extract → review → integrate" flow (Steps 1-8), add new Steps 9-11 that verify existing project-level specs still match the current codebase and fix any drift found.

## Reference Files
- `docs/proposals/spec-drift-detection/proposal.md` — Source proposal
- `plugins/forge/skills/consolidate-specs/SKILL.md` — Primary file to modify

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/consolidate-specs/SKILL.md` | Add Steps 9-11 (drift detect, auto-fix, commit), adjust HARD-GATE, update Workflow diagram |

## Acceptance Criteria

- [ ] SKILL.md contains new Steps 9-11 after existing Step 8
- [ ] Step 9: Detect Drift — reads all `docs/business-rules/*.md` and `docs/conventions/*.md`, validates each rule against current code, classifies as `current` / `drifted` / `orphaned`
- [ ] Step 10: Auto-fix Drift — updates drifted rules in-place (preserving project-global IDs), removes orphaned rules (commit message records rule ID + reason), detects implicit new rules from code changes
- [ ] Step 11: Commit Changes — commits modified spec files with descriptive message listing changed rule IDs
- [ ] HARD-GATE updated: second bullet changes from "Do NOT overwrite existing project-level spec files — append only" to "Do NOT overwrite existing project-level spec files — append only, unless drift is detected in Step 9"
- [ ] Workflow diagram updated to include Steps 9-11
- [ ] Skill supports drift-only mode: when no PRD/design exists (quick mode), skip Steps 1-8 and run only Steps 9-11

## Hard Rules

- Only modify `plugins/forge/skills/consolidate-specs/SKILL.md`
- Project-global IDs must be preserved during auto-fix (only update description/behavior text)
- Deleted rules must be recorded in commit message with ID and deletion reason

## Implementation Notes

- Drift detection should compare rule keywords against actual code — not simple text matching (mitigates false positives per proposal risk analysis)
- The drift-only mode flag can be implicit: if `prd/prd-spec.md` and `design/tech-design.md` don't exist, run drift-only
- New implicit rules from code should be extracted with `[CROSS]` classification and presented to user before appending
