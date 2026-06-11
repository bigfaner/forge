---
id: "1"
title: "Inject Karpathy principles into coding-feature and coding-enhancement templates"
priority: "P1"
estimated_time: "30m"
dependencies: []
type: "doc"
mainSession: false
---

# 1: Inject Karpathy principles into coding-feature and coding-enhancement templates

## Description

Add Karpathy's four coding principles (Think Before Coding, Simplicity First, Surgical Changes, Goal-Driven Execution) to both `coding-feature.md` and `coding-enhancement.md` templates, wrapped in `<CODING_PRINCIPLES>` XML tags.

Both templates share identical structure (3-step TDD workflow) and receive the same full principle set. The principles act as behavioral guidelines governing the agent's approach throughout execution — not additional workflow steps.

## Reference Files
- `docs/proposals/karpathy-coding-guidelines/proposal.md` — Source proposal
- `docs/conventions/forge-distribution.md` — Forge architecture context

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/prompt/data/coding-feature.md` | Add `<CODING_PRINCIPLES>` block with 4 principles after role description, before `## Workflow` |
| `forge-cli/pkg/prompt/data/coding-enhancement.md` | Add `<CODING_PRINCIPLES>` block with 4 principles after role description, before `## Workflow` |

## Acceptance Criteria
- Both files contain `<CODING_PRINCIPLES>` block (uppercase XML tags) positioned after the role description line, before `## Workflow`
- Block includes all 4 principles: Think Before Coding, Simplicity First, Surgical Changes, Goal-Driven Execution
- No timing conflict with Step 1 (Read Task Definition) — "Think Before Coding" guides Step 1 behavior, does not insert a new step before it
- No semantic overlap with existing `<IMPORTANT>` block — principles complement, not duplicate
- Step numbering (Step 1/3, 2/3, 3/3) unchanged
- Template placeholders (`{{TASK_ID}}`, `{{TASK_FILE}}`, `{{SCOPE}}`, `{{PHASE_SUMMARY}}`) undisturbed

## Hard Rules
- MUST read `docs/proposals/karpathy-coding-guidelines/proposal.md` section "结构设计原则" before writing any principle text
- MUST position `<CODING_PRINCIPLES>` after role line, before `## Workflow` — not inside any step

## Implementation Notes
- Principle text should be concise (~50 words per principle). Verbose principles reduce agent compliance.
- Use "trivial tasks use judgment" caveat in Simplicity First to prevent over-conservatism.
- The `<CODING_PRINCIPLES>` tag creates a clear instruction hierarchy: CODING_PRINCIPLES (behavioral guide) < IMPORTANT (task-level hard constraint) < HARD-GATE (process checkpoint).
