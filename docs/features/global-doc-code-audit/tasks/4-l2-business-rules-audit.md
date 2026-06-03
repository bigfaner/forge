---
id: "4"
title: "L2 Business Rules + CLAUDE.md Audit"
priority: "P1"
estimated_time: "1.5h"
dependencies: [1]
type: "doc"
complexity: "medium"
mainSession: false
---

# 4: L2 Business Rules + CLAUDE.md Audit

## Description
Audit docs/business-rules/ (4 files) and root CLAUDE.md for consistency with actual implementation. Business rules govern domain constraints consumed by AI agents during task execution. CLAUDE.md is the first entry point for AI agents understanding project conventions.

## Reference Files
- `docs/proposals/global-doc-code-audit/proposal.md` — Audit Execution Flow, Constraints & Dependencies, Scope
- docs/business-rules/: domain constraint docs vs actual enforcement in code (ref: Audit Execution Flow)
- CLAUDE.md: AI agent instructions vs actual project structure and conventions (ref: Scope)

## Affected Files

### Create
| File | Description |
|------|-------------|
| docs/features/global-doc-code-audit/audit/l2-business-rules-report.md | L2 business rules + CLAUDE.md audit report |

### Modify
| File | Changes |
|------|---------|
| (none) | Audit reads only |

### Delete
| File | Reason |
|------|--------|
| (none) | No deletions |

## Acceptance Criteria
- [ ] All 4 business-rules files + CLAUDE.md audited with declaration extraction
- [ ] Each business rule claim verified against actual code enforcement (e.g., naming rules vs code constants)
- [ ] CLAUDE.md claims verified against actual project structure, file paths, and conventions
- [ ] Every inconsistency recorded: file path, line range, severity (P0-P3), suggested action
- [ ] Cross-layer influence items recorded for L3 reference
- [ ] Audit report follows unified template

## Hard Rules
- Do NOT modify any code or documentation — audit only
- All audit output written in English

## Implementation Notes
- Business rules may reference specific code paths or constants — verify these exist in current codebase
- CLAUDE.md is critical: incorrect instructions here cause cascading AI agent errors
- Distribution model: CLAUDE.md exists in source repo only; verify claims about plugin distribution paths
