---
id: "5"
title: "L2 Conventions Audit Batch 1"
priority: "P1"
estimated_time: "1.5h"
dependencies: [1]
type: "doc"
complexity: "high"
mainSession: false
---

# 5: L2 Conventions Audit Batch 1

## Description
Audit first batch of docs/conventions/ top-level files (8 files): code-structure.md, constants.md, dead-code.md, dispatcher-quality.md, enum-constants.md, error-handling.md, forge-cli-reference.md, forge-distribution.md.

Conventions are consumed by AI agents during task execution — incorrect conventions lead to wrong code generation.

## Reference Files
- `docs/proposals/global-doc-code-audit/proposal.md` — Audit Execution Flow, Constraints & Dependencies, Key Risks
- docs/conventions/code-structure.md: directory structure claims vs actual layout (ref: Audit Execution Flow)
- docs/conventions/constants.md: constant naming rules vs actual Go code (ref: Audit Execution Flow)
- docs/conventions/forge-distribution.md: distribution model constraints vs actual plugin paths (ref: Constraints & Dependencies)

## Affected Files

### Create
| File | Description |
|------|-------------|
| docs/features/global-doc-code-audit/audit/l2-conventions-batch1-report.md | L2 conventions batch 1 audit report |

### Modify
| File | Changes |
|------|---------|
| (none) | Audit reads only |

### Delete
| File | Reason |
|------|--------|
| (none) | No deletions |

## Acceptance Criteria
- [ ] All 8 target files audited with declaration extraction
- [ ] Each convention claim verified: file paths via `find`, code constants via `grep`, structural rules vs actual codebase
- [ ] Every inconsistency recorded: file path, line range, severity (P0-P3), suggested action
- [ ] Cross-layer influence items recorded for L3 reference
- [ ] Audit report follows unified template

## Hard Rules
- Do NOT modify any code or documentation — audit only
- Only audit the following 8 files: code-structure.md, constants.md, dead-code.md, dispatcher-quality.md, enum-constants.md, error-handling.md, forge-cli-reference.md, forge-distribution.md
- All audit output written in English

## Implementation Notes
- forge-distribution.md is particularly important: it defines path resolution rules that affect all plugin components
- constants.md and enum-constants.md: verify each constant name mentioned actually exists in current Go code
- code-structure.md: verify directory structure claims against actual filesystem
