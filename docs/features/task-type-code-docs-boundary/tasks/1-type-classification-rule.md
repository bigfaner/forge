---
id: "1"
title: "Add code/docs classification rule to type-assignment.md"
priority: "P1"
estimated_time: "15min"
dependencies: []
type: "documentation"
mainSession: false
---

# 1: Add code/docs classification rule to type-assignment.md

## Description
The `type-assignment.md` reference document currently lists type definitions but lacks an explicit classification rule that maps types to quality-gate behavior. Agent assigns type by intent rather than by output artifact, leading to docs-only tasks being marked as `enhancement` or `implementation`.

Add a clear "classify by output artifact" rule with a code/docs/meta classification table, so downstream consumers (guide.md, submit-task, task generators) can reference a single source of truth.

## Reference Files
- `docs/proposals/task-type-code-docs-boundary/proposal.md` — Source proposal

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/references/shared/type-assignment.md` | Add code/docs/meta classification table and "按产出物分类" rule |

## Acceptance Criteria
- [ ] Document contains a classification table: Code types (feature, enhancement, cleanup, refactor, fix → quality-gate), Doc type (documentation → skip), Meta type (gate → special)
- [ ] Rule states "classify by output artifact, not by intent" with examples (e.g., "improving agent prompts in .md files = documentation, not enhancement")
- [ ] Existing type definitions table preserved unchanged

## Implementation Notes
- Keep the existing table as-is; add the new classification section after it
- This file is a shared reference read by quick-tasks and breakdown-tasks skills via `${CLAUDE_SKILL_DIR}/../../references/shared/type-assignment.md`
