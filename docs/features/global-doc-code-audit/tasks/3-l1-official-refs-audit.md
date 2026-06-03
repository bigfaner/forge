---
id: "3"
title: "L1 Official References Audit"
priority: "P1"
estimated_time: "1.5h"
dependencies: [1]
type: "doc"
complexity: "medium"
mainSession: false
---

# 3: L1 Official References Audit

## Description
Audit docs/official-references/ (5 files: hooks.md, plugin-marketplace.md, plugin.md, skills-ref.md, worktree.md) for consistency with actual code implementation. These reference documents describe plugin hooks, plugin marketplace, plugin structure, skill definitions, and worktree management.

## Reference Files
- `docs/proposals/global-doc-code-audit/proposal.md` — Audit Execution Flow, Proposed Solution, Constraints & Dependencies
- docs/official-references/hooks.md: hook definitions vs actual hook registry in code (ref: Proposed Solution)
- docs/official-references/skills-ref.md: skill structure claims vs actual SKILL.md files (ref: Proposed Solution)
- docs/official-references/worktree.md: worktree management claims vs forge CLI code (ref: Proposed Solution)

## Affected Files

### Create
| File | Description |
|------|-------------|
| docs/features/global-doc-code-audit/audit/l1-official-refs-report.md | L1 official references audit report |

### Modify
| File | Changes |
|------|---------|
| (none) | Audit reads only |

### Delete
| File | Reason |
|------|--------|
| (none) | No deletions |

## Acceptance Criteria
- [ ] All 5 target files audited with complete declaration extraction
- [ ] Each claim verified: hook names/parameters vs code, plugin structure vs actual templates, skill definitions vs actual SKILL.md files
- [ ] Every inconsistency recorded: file path, line range, severity (P0-P3), suggested action
- [ ] Cross-layer influence items recorded for L3 reference
- [ ] Audit report follows unified template

## Hard Rules
- Do NOT modify any code or documentation — audit only
- All audit output written in English

## Implementation Notes
- hooks.md is high-priority: incorrect hook documentation directly causes AI agent errors
- plugin.md and skills-ref.md should be cross-checked against plugins/forge/ actual structure
- Distribution model: these docs only exist in source repo, not distributed to users
