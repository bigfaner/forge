---
id: "13"
title: "L3 Final Batch: Remaining Lessons + All Decisions"
priority: "P1"
estimated_time: "2h"
dependencies: [2, 3, 4, 5, 6]
type: "doc"
complexity: "high"
mainSession: false
---

# 13: L3 Final Batch: Remaining Lessons + All Decisions

## Description
Audit the final 13 lesson files (lesson-vibe-coding-scope-control.md through worktree-stale-refs.md) plus all 10 decision files from docs/decisions/ (architecture.md through testing.md). Total: 23 items. Classify each item (code-reference, process-standard, experience-summary), assess validity using L3 structured rules, detect duplicates via topic clustering, and mark every item as valid/outdated/duplicate/needs-update with justification. Decisions are architectural records — assess whether each decision's rationale still holds and whether the described implementation still matches current code.

## Target Files — Lessons (13)
1. lesson-vibe-coding-scope-control.md
2. pattern-compile-check-before-submit.md
3. pattern-dispatcher-auto-verify.md
4. pattern-large-scale-rename.md
5. pattern-sitemap-shared-layout.md
6. pattern-surface-resolution-shortcut.md
7. pattern-task-vs-output-naming.md
8. tool-cli-e2e-lifecycle.md
9. tool-fix-e2e-unknown-placeholder.md
10. tool-justfile-arg-attribute.md
11. tool-record-coverage-capture.md
12. tool-submit-background-timeout.md
13. worktree-stale-refs.md

## Target Files — Decisions (10)
1. architecture.md
2. data-model.md
3. dependencies.md
4. e2e-server-lifecycle-hardening.md
5. error-handling.md
6. interface.md
7. local-dev-deployment.md
8. manifest.md
9. security.md
10. testing.md

## Reference Files
- `docs/proposals/global-doc-code-audit/proposal.md` — L3 Knowledge Base Audit Flow, Success Criteria, L3 Validity Rules

## Affected Files

### Create
| File | Description |
|------|-------------|
| docs/features/global-doc-code-audit/audit/l3-final-batch-report.md | L3 final batch (remaining lessons + all decisions) audit report |

### Modify
| File | Changes |
|------|---------|
| (none) | Audit reads only |

### Delete
| File | Reason |
|------|--------|
| (none) | No deletions |

## Acceptance Criteria
- [ ] All 23 target items classified (code-reference, process-standard, experience-summary)
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
- Decisions (docs/decisions/) are architectural records — assess whether each decision's rationale still holds and whether the described implementation still matches current code
