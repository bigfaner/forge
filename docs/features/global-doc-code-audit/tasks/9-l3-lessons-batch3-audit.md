---
id: "9"
title: "L3 Lessons Audit Batch 3 (gotcha-fix-task-claim to gotcha-macos-sleep)"
priority: "P1"
estimated_time: "2h"
dependencies: [2, 3, 4, 5, 6]
type: "doc"
complexity: "high"
mainSession: false
---

# 9: L3 Lessons Audit Batch 3 (gotcha-fix-task-claim to gotcha-macos-sleep)

## Description
Audit 20 lesson files from docs/lessons/: gotcha-fix-task-claim-priority.md through gotcha-macos-sleep-kills-subagent-connection.md. Classify each item (code-reference, process-standard, experience-summary), assess validity using L3 structured rules, detect duplicates via topic clustering, and mark every item as valid/outdated/duplicate/needs-update with justification.

## Target Files
1. gotcha-fix-task-claim-priority.md
2. gotcha-fix-task-dependency-chain.md
3. gotcha-fix-task-empty-type.md
4. gotcha-fix-task-index-test-isolation.md
5. gotcha-fix-task-scope-too-broad.md
6. gotcha-fix-task-type-hardcoded.md
7. gotcha-forge-cli-invocation.md
8. gotcha-forge-feature-no-get-subcommand.md
9. gotcha-forge-task-index-always-required.md
10. gotcha-forge-task-index-per-type-duplicate.md
11. gotcha-gen-test-scripts-ts-residue.md
12. gotcha-go-test-staging-graduation-friction.md
13. gotcha-graduation-dual-module-drift.md
14. gotcha-hook-idempotency-feature-complete.md
15. gotcha-hook-startup-errors.md
16. gotcha-hook-unbounded-test-timeout.md
17. gotcha-implementation-type-for-skill-files.md
18. gotcha-journey-hallucination-revision-death-spiral.md
19. gotcha-large-output-stall-subagent.md
20. gotcha-macos-sleep-kills-subagent-connection.md

## Reference Files
- `docs/proposals/global-doc-code-audit/proposal.md` — L3 Knowledge Base Audit Flow, Success Criteria, L3 Validity Rules

## Affected Files

### Create
| File | Description |
|------|-------------|
| docs/features/global-doc-code-audit/audit/l3-lessons-batch3-report.md | L3 lessons batch 3 audit report |

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
