---
id: "3"
title: "Inject Surgical Changes principle into coding-refactor template"
priority: "P1"
estimated_time: "20m"
dependencies: []
type: "doc"
mainSession: false
---

# 3: Inject Surgical Changes principle into coding-refactor template

## Description

Add the Surgical Changes principle only to `coding-refactor.md`. The "Think Before Coding" concern is already addressed by the existing Impact Mapping step (Step 2), so no duplication needed.

## Reference Files
- `docs/proposals/karpathy-coding-guidelines/proposal.md` — Source proposal

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/prompt/data/coding-refactor.md` | Add `<CODING_PRINCIPLES>` block with Surgical Changes principle after role + External behavior definition, before `## Pre-check` |

## Acceptance Criteria
- Contains `<CODING_PRINCIPLES>` block with only the Surgical Changes principle
- Positioned after the "External behavior" definition block, before `## Pre-check`
- No overlap with Impact Mapping step (Step 2) content — Impact Mapping maps change scope; Surgical Changes constrains what gets touched during execution
- Pre-check section preserved
- Step numbering (Step 1/4, 2/4, 3/4, 4/4) unchanged
- Template placeholders undisturbed

## Hard Rules
- MUST include ONLY Surgical Changes principle — not Think Before Coding (already covered by Impact Mapping)

## Implementation Notes
- Refactor template is the most complex (170 lines, 4 steps with structural/behavioral sub-workflows). The principle should be brief to avoid further inflating context.
- Surgical Changes principle for refactors should emphasize: don't "improve" code outside the explicit refactoring scope, even if you notice opportunities.
