---
id: "8"
title: "L3 Lessons Audit Batch 2 (gotcha-breaking-change to gotcha-fix-task-claim)"
priority: "P1"
estimated_time: "2h"
dependencies: [2, 3, 4, 5, 6]
type: "doc"
complexity: "high"
mainSession: false
---

# 8: L3 Lessons Audit Batch 2 (gotcha-breaking-change to gotcha-fix-task-claim)

## Description
Audit 20 lesson files from docs/lessons/: gotcha-breaking-change-quality-gate-deadlock.md through gotcha-fix-task-broad-scope.md. Classify each item (code-reference, process-standard, experience-summary), assess validity using L3 structured rules, detect duplicates via topic clustering, and mark every item as valid/outdated/duplicate/needs-update with justification.

## Target Files
1. gotcha-breaking-change-quality-gate-deadlock.md
2. gotcha-breaking-task-quality-gate-test-scope.md
3. gotcha-characterization-test-vs-refactoring.md
4. gotcha-dispatcher-ignores-compilation-diagnostics.md
5. gotcha-docs-only-needs-code-audit.md
6. gotcha-drift-detection-task-runtime.md
7. gotcha-duplicate-test-runs.md
8. gotcha-e2e-env-override-isolation.md
9. gotcha-e2e-script-generation.md
10. gotcha-e2e-skill-monorepo-path-mismatch.md
11. gotcha-e2e-test-binary-isolation.md
12. gotcha-e2e-test-quality-antipatterns.md
13. gotcha-embedded-template-name-mismatch.md
14. gotcha-eval-loop-decision-gate.md
15. gotcha-eval-prd-use-zcode-agents.md
16. gotcha-eval-reviser-too-many-attacks.md
17. gotcha-eval-rollback-destroys-improvements.md
18. gotcha-eval-rubric-misses-disguised-patches.md
19. gotcha-eval-subagent-type.md
20. gotcha-fix-task-broad-scope.md

## Reference Files
- `docs/proposals/global-doc-code-audit/proposal.md` — L3 Knowledge Base Audit Flow, Success Criteria, L3 Validity Rules

## Affected Files

### Create
| File | Description |
|------|-------------|
| docs/features/global-doc-code-audit/audit/l3-lessons-batch2-report.md | L3 lessons batch 2 audit report |

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
