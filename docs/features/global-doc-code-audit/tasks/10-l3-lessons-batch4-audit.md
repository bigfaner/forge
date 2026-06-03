---
id: "10"
title: "L3 Lessons Audit Batch 4 (gotcha-main-session to gotcha-revert-mid)"
priority: "P1"
estimated_time: "2h"
dependencies: [2, 3, 4, 5, 6]
type: "doc"
complexity: "high"
mainSession: false
---

# 10: L3 Lessons Audit Batch 4 (gotcha-main-session to gotcha-revert-mid)

## Description
Audit 20 lesson files from docs/lessons/: gotcha-main-session-flag.md through gotcha-revert-mid-dispatch.md. Classify each item (code-reference, process-standard, experience-summary), assess validity using L3 structured rules, detect duplicates via topic clustering, and mark every item as valid/outdated/duplicate/needs-update with justification.

## Target Files
1. gotcha-main-session-flag.md
2. gotcha-merge-ghost-revival.md
3. gotcha-pipeline-skill-bypass.md
4. gotcha-post-completion-commit.md
5. gotcha-pre-existing-syntax-errors-block.executor.md
6. gotcha-present-analysis-before-edit.md
7. gotcha-prompt-template-complexity-agnostic.md
8. gotcha-proposal-baseline-drift.md
9. gotcha-proposal-success-criteria-contradiction.md
10. gotcha-quality-gate-buffered-output-appears-dead.md
11. gotcha-quality-gate-cross-feature-pollution.md
12. gotcha-quality-gate-doc-type-ignore.md
13. gotcha-quality-gate-fix-task-loop.md
14. gotcha-quick-tasks-merge-threshold.md
15. gotcha-quick-tasks-no-autochain.md
16. gotcha-quick-tasks-no-commit.md
17. gotcha-quick-tasks-stale-detect-command.md
18. gotcha-recursive-go-test-process-explosion.md
19. gotcha-redundant-manual-e2e-verification.md
20. gotcha-revert-mid-dispatch.md

## Reference Files
- `docs/proposals/global-doc-code-audit/proposal.md` — L3 Knowledge Base Audit Flow, Success Criteria, L3 Validity Rules

## Affected Files

### Create
| File | Description |
|------|-------------|
| docs/features/global-doc-code-audit/audit/l3-lessons-batch4-report.md | L3 lessons batch 4 audit report |

### Modify
| File | Changes |
|------|---------|
| (none) | Audit reads only |

### Delete
| File | Reason |
|------|--------|
| (none) | No deletions |

## Acceptance Criteria
- [ ] All 20 target items classified (code-reference, process-standard, experience-summary)
- [ ] Each item's validity assessed using structured rules (tool change -> outdated, process contradiction -> needs-update, path invalid -> outdated, generalized conclusion -> valid, partial aging -> needs-update)
- [ ] Duplicate detection performed via topic clustering
- [ ] Every item marked as valid/outdated/duplicate/needs-update with justification
- [ ] Cross-layer influence items from L1/L2 reports checked against relevant items
- [ ] Audit report follows unified template

## Hard Rules
- Do NOT modify any code or documentation — audit only
- All audit output written in English

## Implementation Notes
- Apply L3 validity rules from proposal: valid, outdated, duplicate, needs-update
- For code-reference items: verify paths via find/grep
- For process/experience items: assess against current project state (directory structure, toolchain, team conventions)
- Human confirmation required for deletion/merge recommendations — mark these clearly in the report
