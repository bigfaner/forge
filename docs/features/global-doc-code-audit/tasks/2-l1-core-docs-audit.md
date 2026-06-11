---
id: "2"
title: "L1 Core User Docs Audit"
priority: "P1"
estimated_time: "2h"
dependencies: [1]
type: "doc"
complexity: "high"
mainSession: false
---

# 2: L1 Core User Docs Audit

## Description
Audit core user-facing documentation files for consistency with actual code behavior. Target files: docs/ARCHITECTURE.md, DESIGN.md, and docs/user-guide/ (4 files: architecture-overview.md, environment-setup.md, initialization.md, usage-guide.md).

Follow the L1 audit methodology. Apply cross-layer influence tracking: record any code structure inconsistencies found (e.g., hook execution order, module paths) in the cross-layer influence list for L3 reference.

## Reference Files
- `docs/proposals/global-doc-code-audit/proposal.md` — Audit Execution Flow, Proposed Solution, Key Risks
- docs/ARCHITECTURE.md: system architecture claims vs actual hook/module structure (ref: Proposed Solution)
- DESIGN.md: design document claims vs current implementation (ref: Proposed Solution)
- docs/user-guide/: user-facing behavior descriptions vs CLI actual behavior (ref: Proposed Solution)

## Affected Files

### Create
| File | Description |
|------|-------------|
| docs/features/global-doc-code-audit/audit/l1-core-docs-report.md | L1 core docs audit report |

### Modify
| File | Changes |
|------|---------|
| (none) | Audit reads only |

### Delete
| File | Reason |
|------|--------|
| (none) | No deletions |

## Acceptance Criteria
- [ ] All 6 target files audited with complete declaration extraction
- [ ] Each claim verified against codebase (paths via `find`/`grep`, behaviors via code reading)
- [ ] Every inconsistency recorded: file path, line range, severity (P0-P3), suggested action
- [ ] Cross-layer influence items identified and recorded for L3 reference (e.g., hook names, module paths mentioned in docs)
- [ ] Audit report follows unified template: baseline commit, issue summary, issue details, quality review

## Hard Rules
- Do NOT modify any code or documentation — audit only
- Record baseline commit hash in report header
- All audit output written in English

## Implementation Notes
- ARCHITECTURE.md is expected to be the most complex file — allocate sufficient time for behavioral comparison
- user-guide/ files describe user-facing CLI behavior; verify against actual CLI command signatures and output
- Per proposal risk: watch for document complexity variance — ARCHITECTURE.md may need deeper code tracing than usage guide
