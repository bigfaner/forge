---
id: "2"
title: "Inject Karpathy principles into coding-fix template, replacing existing rules"
priority: "P1"
estimated_time: "30m"
dependencies: []
type: "doc"
mainSession: false
---

# 2: Inject Karpathy principles into coding-fix template, replacing existing rules

## Description

Add Think Before Coding, Simplicity First, and Surgical Changes principles to `coding-fix.md`. Replace the existing `<IMPORTANT>` block containing "MINIMAL CHANGES" and "NO REFACTORING" with the new `<CODING_PRINCIPLES>` block, ensuring semantic equivalence or stronger coverage.

## Reference Files
- `docs/proposals/karpathy-coding-guidelines/proposal.md` — Source proposal

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/prompt/data/coding-fix.md` | Replace first `<IMPORTANT>` block (MINIMAL CHANGES + NO REFACTORING) with `<CODING_PRINCIPLES>` block containing 3 principles |

## Acceptance Criteria
- First `<IMPORTANT>` block (lines 9-12: "MINIMAL CHANGES" + "NO REFACTORING") replaced by `<CODING_PRINCIPLES>` block with Think Before Coding + Simplicity First + Surgical Changes
- Semantic coverage: new principles enforce at least the same constraints as old rules (minimal fix scope, no refactoring)
- Second `<IMPORTANT>` block (lines 33-37: Hard Rules about file scope restrictions) preserved unchanged
- `<CODING_PRINCIPLES>` positioned after role description, before `## Workflow`
- Step numbering (Step 1/4, 2/4, 3/4, 4/4) unchanged
- Template placeholders undisturbed

## Hard Rules
- MUST NOT keep both old `<IMPORTANT>` rules AND new principles — overlap must be resolved by replacement
- MUST preserve the second `<IMPORTANT>` block (Hard Rules) exactly as-is

## Implementation Notes
- "Simplicity First" directly replaces "MINIMAL CHANGES" with broader coverage (no unrequested changes)
- "Surgical Changes" directly replaces "NO REFACTORING" with more actionable guidance (scope boundary, no adjacent code changes)
- "Think Before Coding" is new for fix template — agents often misdiagnose errors by jumping to the first plausible fix without verifying the root cause
