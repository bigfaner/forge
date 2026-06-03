---
id: "15"
title: "Generate Fix Tasks from Audit Findings"
priority: "P0"
estimated_time: "1h"
dependencies: [14]
type: "doc"
complexity: "high"
mainSession: false
---

# 15: Generate Fix Tasks from Audit Findings

## Description
Convert all audit findings from the consolidated report into executable fix tasks using three task templates: fix-type (task-executor independent), review-type (human confirmation required), and cross-layer-verification-type (depends on prior audit reports).

## Reference Files
- `docs/proposals/global-doc-code-audit/proposal.md` — Success Criteria, Non-Functional Requirements, Constraints & Dependencies
- docs/features/global-doc-code-audit/audit/consolidated-report.md: source findings (ref: Success Criteria)
- docs/conventions/forge-distribution.md: distribution model for task template compliance (ref: Constraints & Dependencies)

## Affected Files

### Create
| File | Description |
|------|-------------|
| docs/features/global-doc-code-audit/audit/fix-tasks.md | Generated fix tasks using three template types |

### Modify
| File | Changes |
|------|---------|
| (none) | Reads consolidated report only |

### Delete
| File | Reason |
|------|--------|
| (none) | No deletions |

## Acceptance Criteria
- [ ] All findings converted to executable fix tasks using appropriate template: fix-type, review-type, or cross-layer-verification-type
- [ ] Knowledge base cleanup tasks (deletion/merge recommendations) marked as requiring human confirmation
- [ ] Fix tasks are self-contained: include full context, do not depend on other fix tasks
- [ ] All output written in English

## Hard Rules
- Knowledge base deletion/merge tasks MUST be marked as human-confirmation-required
- All fix task templates written in English
- Fix tasks must be self-contained (include context, not depend on other fix tasks)

## Implementation Notes
- Fix-type tasks: standard task-executor template, independent execution
- Review-type tasks: include human confirmation checkpoints for knowledge base changes
- Cross-layer-verification-type tasks: reference cross-layer influence lists as input dependency
- Human confirmation SLA: proposal author checks daily, escalates to P1 after 3 working days
