---
id: "11"
title: "L3 Lessons Audit Batch 5 (gotcha-review-task to gotcha-task-executor-invisible)"
priority: "P1"
estimated_time: "2h"
dependencies: [2, 3, 4, 5, 6]
type: "doc"
complexity: "high"
mainSession: false
---

# 11: L3 Lessons Audit Batch 5 (gotcha-review-task to gotcha-task-executor-invisible)

## Description
Audit 20 lesson files from docs/lessons/: gotcha-review-task-incomplete-dependencies.md through gotcha-task-executor-ignores-implementation-notes.md. Classify each item (code-reference, process-standard, experience-summary), assess validity using L3 structured rules, detect duplicates via topic clustering, and mark every item as valid/outdated/duplicate/needs-update with justification.

## Target Files
1. gotcha-review-task-incomplete-dependencies.md
2. gotcha-reviser-agent-long-running.md
3. gotcha-run-tasks-no-auto-test.md
4. gotcha-shared-interface-mock-cascade.md
5. gotcha-skill-step-analysis-paralysis.md
6. gotcha-skip-plan-for-pipeline-change.md
7. gotcha-spec-authority-drift.md
8. gotcha-split-rules-operational-blindness.md
9. gotcha-split-task-missing-shared-setup.md
10. gotcha-stale-skill-cli-flags.md
11. gotcha-stale-state-json-feature-mismatch.md
12. gotcha-stale-test-results-cascade.md
13. gotcha-standard-task-id-collision.md
14. gotcha-stop-hook-non-blocking-error.md
15. gotcha-strategy-bypass-justfile.md
16. gotcha-surface-fields-single-surface-empty.md
17. gotcha-task-cli-path-duplication.md
18. gotcha-task-derivation-over-research.md
19. gotcha-task-executor-auto-claim.md
20. gotcha-task-executor-ignores-implementation-notes.md

## Reference Files
- `docs/proposals/global-doc-code-audit/proposal.md` — L3 Knowledge Base Audit Flow, Success Criteria, L3 Validity Rules

## Affected Files

### Create
| File | Description |
|------|-------------|
| docs/features/global-doc-code-audit/audit/l3-lessons-batch5-report.md | L3 lessons batch 5 audit report |

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
