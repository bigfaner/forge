---
id: "12"
title: "L3 Lessons Audit Batch 6 (gotcha-task-executor-invisible to lesson-vibe)"
priority: "P1"
estimated_time: "2h"
dependencies: [2, 3, 4, 5, 6]
type: "doc"
complexity: "high"
mainSession: false
---

# 12: L3 Lessons Audit Batch 6 (gotcha-task-executor-invisible to lesson-vibe)

## Description
Audit 20 lesson files from docs/lessons/: gotcha-task-executor-invisible-thinking-time.md through lesson-tui-visual-verify.md. Classify each item (code-reference, process-standard, experience-summary), assess validity using L3 structured rules, detect duplicates via topic clustering, and mark every item as valid/outdated/duplicate/needs-update with justification.

## Target Files
1. gotcha-task-executor-invisible-thinking-time.md
2. gotcha-task-executor-never-returns.md
3. gotcha-task-executor-proposal-inaccessible.md
4. gotcha-task-executor-redundant-search.md
5. gotcha-task-executor-stops-at-step1.md
6. gotcha-task-executor-thinking-overhead.md
7. gotcha-task-index-preserve-deps.md
8. gotcha-task-reference-files-scope-creep.md
9. gotcha-task-reference-source-anchor-misread.md
10. gotcha-task-type-documentation-vs-doc.md
11. gotcha-task-type-for-md-files.md
12. gotcha-test-chain-not-linked-to-last-business-gate.md
13. gotcha-test-pipeline-no-languages.md
14. gotcha-test-script-staging-vs-graduation.md
15. hook-stop-e2e-blocking.md
16. lesson-forge-tui-pipeline-gap.md
17. lesson-gate-force-over-fix.md
18. lesson-guide-missing-proposals-dir.md
19. lesson-tui-tech-design-mockup.md
20. lesson-tui-visual-verify.md

## Reference Files
- `docs/proposals/global-doc-code-audit/proposal.md` — L3 Knowledge Base Audit Flow, Success Criteria, L3 Validity Rules

## Affected Files

### Create
| File | Description |
|------|-------------|
| docs/features/global-doc-code-audit/audit/l3-lessons-batch6-report.md | L3 lessons batch 6 audit report |

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
