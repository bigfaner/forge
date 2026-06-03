---
id: "7"
title: "L3 Lessons Audit Batch 1 (arch/gotcha-a to gotcha-breaking-change)"
priority: "P1"
estimated_time: "2h"
dependencies: [2, 3, 4, 5, 6]
type: "doc"
complexity: "high"
mainSession: false
---

# 7: L3 Lessons Audit Batch 1 (arch/gotcha-a to gotcha-breaking-change)

## Description
Audit the first 20 lesson files (alphabetically) from docs/lessons/: arch-constant-rename-whack-a-mole.md through gotcha-breaking-change-integration-test-blast-radius.md. Classify each item (code-reference, process-standard, experience-summary), assess validity using L3 structured rules, detect duplicates via topic clustering, and mark every item as valid/outdated/duplicate/needs-update with justification.

## Target Files
1. arch-constant-rename-whack-a-mole.md
2. arch-dispatcher-post-loop-message-misleading.md
3. arch-forge-skill-gap-analysis.md
4. arch-freeform-findings-indirect-influence.md
5. arch-post-loop-artifact-commit-gap.md
6. arch-prototype-navigation-contract.md
7. arch-task-failure-recovery-loop.md
8. arch-task-record-immutability.md
9. arch-test-type-index-chicken-egg.md
10. fix-zsh-compinit-docker.md
11. gotcha-ac-self-report-without-verification.md
12. gotcha-adjacent-section-over-removal.md
13. gotcha-agent-example-over-schema.md
14. gotcha-api-no-api-prefix.md
15. gotcha-auto-gen-tasks-reappear-and-preempt-fix.md
16. gotcha-auto-push-no-upstream.md
17. gotcha-auto-unblock-loop.md
18. gotcha-blocked-task-never-auto-unblocks.md
19. gotcha-brainstorm-challenge-failure.md
20. gotcha-breaking-change-integration-test-blast-radius.md

## Reference Files
- `docs/proposals/global-doc-code-audit/proposal.md` — L3 Knowledge Base Audit Flow, Success Criteria, L3 Validity Rules

## Affected Files

### Create
| File | Description |
|------|-------------|
| docs/features/global-doc-code-audit/audit/l3-lessons-batch1-report.md | L3 lessons batch 1 audit report |

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
