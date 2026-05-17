---
id: "1"
title: "Create /learn skill (SKILL.md + templates)"
priority: "P0"
estimated_time: "2h"
dependencies: []
type: "feature"
mainSession: false
---

# 1: Create /learn skill (SKILL.md + templates)

## Description

Create the unified `/learn` skill that absorbs `/record-decision` and `/learn-lesson` functionality into a single entry point. The skill supports two input modes (interactive + direct args), writes entries first then reports for review, and classifies knowledge into the 4 knowledge directories.

The skill must reuse existing format specifications:
- Decision format: from `plugins/forge/references/shared/decision-logging.md` (Sections 6-7: row format, manifest update)
- Lesson format: from `plugins/forge/skills/learn-lesson/templates/template.md`
- Convention/business-rule format: append entries matching existing `docs/conventions/` and `docs/business-rules/` file patterns

## Reference Files
- `docs/proposals/knowledge-accumulation-loop/proposal.md` — Source proposal
- `plugins/forge/references/shared/decision-logging.md` — Decision archiving format
- `plugins/forge/skills/learn-lesson/SKILL.md` — Existing lesson skill workflow
- `plugins/forge/skills/learn-lesson/templates/template.md` — Lesson file template
- `plugins/forge/skills/learn-lesson/examples/debug-race-condition.md` — Lesson example
- `plugins/forge/skills/consolidate-specs/SKILL.md` — Vocabulary reference for classification

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/learn/SKILL.md` | Unified knowledge accumulation skill |
| `plugins/forge/skills/learn/templates/decision-entry.md` | Decision entry format template (extracted from decision-logging.md) |
| `plugins/forge/skills/learn/templates/lesson-entry.md` | Lesson entry format template (adapted from learn-lesson) |
| `plugins/forge/skills/learn/templates/convention-entry.md` | Convention entry format template |

### Modify
| File | Changes |
|------|---------|
| (none) | |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria
- [ ] SKILL.md defines a skill with frontmatter `name: learn` and clear description
- [ ] Skill supports two input modes: `/learn` (interactive) and `/learn "text"` (direct)
- [ ] Skill workflow: identify knowledge type(s) → classify → write → report
- [ ] Knowledge type identification covers: decision, lesson, convention, business-rule
- [ ] Multi-type capture: single input can produce entries in multiple directories
- [ ] Write-first-then-report: entries are written immediately, all shown in final report for user review
- [ ] Writes to `docs/decisions/` using decision-logging.md row format (Section 6) + manifest update (Section 7)
- [ ] Writes to `docs/lessons/` using learn-lesson template format
- [ ] Writes to `docs/conventions/` by appending entries to existing domain files (or creating new)
- [ ] Writes to `docs/business-rules/` by appending entries to existing domain files (or creating new)
- [ ] Accepts custom vocabulary values (domains, types) without error
- [ ] Detects bulk extraction needs and delegates to `/consolidate-specs`
- [ ] Uses auto-generated vocabulary when available, falls back gracefully when not
- [ ] Category classification reuses 8-category vocabulary from learn-lesson tags + decision-logging type mapping

## Hard Rules
- Must read `plugins/forge/references/shared/decision-logging.md` for decision format specifications — do not reinvent
- Must read `plugins/forge/skills/learn-lesson/templates/template.md` for lesson format — do not reinvent
- All knowledge directory formats must remain compatible with `/consolidate-specs` overlap detection
- No code changes — SKILL.md and templates are prompt-level only

## Implementation Notes
- The classification step should present the 8-category vocabulary as suggestions, not enforced values
- For convention/business-rule entries, reuse the project-global ID encoding from consolidate-specs (BIZ- prefix, TECH- prefix)
- The "write-first, review-after" pattern means the skill does NOT ask for confirmation before writing — it writes and shows what was written in the final report
- The skill should handle the case where `docs/decisions/` doesn't exist yet (auto-create following decision-logging.md Section 8)
