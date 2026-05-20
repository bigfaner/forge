---
id: "4"
title: "Inject principles into coding-cleanup template"
priority: "P1"
estimated_time: "20m"
dependencies: []
type: "doc"
mainSession: false
---

# 4: Inject principles into coding-cleanup template

## Description

Add Simplicity First and Surgical Changes principles to `coding-cleanup.md`. These two principles address the most common cleanup anti-patterns: over-cleaning (removing things that look dead but aren't) and scope creep (touching adjacent code while cleaning).

## Reference Files
- `docs/proposals/karpathy-coding-guidelines/proposal.md` — Source proposal

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/prompt/data/coding-cleanup.md` | Add `<CODING_PRINCIPLES>` block with Simplicity First + Surgical Changes after role description, before `## Workflow` |

## Acceptance Criteria
- Contains `<CODING_PRINCIPLES>` block with Simplicity First + Surgical Changes principles
- Positioned after role description, before `## Workflow`
- No semantic overlap with existing Step 2 "Make Improvements" instructions (principles complement, not duplicate)
- Step numbering (Step 1/3, 2/3, 3/3) unchanged
- Template placeholders undisturbed

## Hard Rules
- MUST include only Simplicity First + Surgical Changes — no additional principles

## Implementation Notes
- Cleanup tasks are prone to "while I'm here" syndrome. Surgical Changes should explicitly call out: don't touch code outside the cleanup target.
