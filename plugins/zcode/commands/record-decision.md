---
name: record-decision
description: Record an architecture/technical decision to docs/decisions/ at any stage of development.
allowed_tools: ["Read", "Write", "Edit", "Bash", "AskUserQuestion"]
---

# /record-decision

Archives a single technical or architectural decision into `docs/decisions/` through a 4-round interactive flow. Can be invoked at any point in the development lifecycle.

## When to Use

- Mid-implementation decisions not captured in the design phase
- Historical supplements: backfilling informally-made decisions
- Brainstorm or PRD phase: important technical constraints before a formal tech-design exists

## Process

Follow the record-decision 4-round interaction flow defined in `plugins/zcode/references/shared/decision-logging.md` (Section 3).

### Round summary

| Round | Prompt | Input |
|-------|--------|-------|
| 1 | Select decision type (1-8) | Number selection |
| 2 | Decision description (one sentence) | Text |
| 3 | Decision rationale (one sentence) | Text |
| 4 | Associated feature slug (or skip) | Text or Enter |

### Auto-filled fields

- `Date`: today's date (YYYY-MM-DD)
- `Source`: `<feature-slug>/tech-design.md` if a feature was provided; `manual` if Round 4 was skipped

### After the 4 rounds

1. Append the decision row to `docs/decisions/<type>.md`
2. Update `docs/decisions/manifest.md` (Categories count + Recent Decisions table)

Refer to Sections 6 and 7 of `plugins/zcode/references/shared/decision-logging.md` for the exact row format and manifest update protocol.

### Error handling

If `docs/decisions/` does not exist, auto-create the directory, all 8 type files, and `manifest.md` before writing. See Section 8 of `plugins/zcode/references/shared/decision-logging.md` for all error scenarios.
