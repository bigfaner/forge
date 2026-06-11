---
id: "1"
title: "Restructure guide.md into 3 thematic sections"
priority: "P1"
estimated_time: "1h"
dependencies: []
type: "documentation"
mainSession: false
---

# 1: Restructure guide.md into 3 thematic sections

## Description
Simplify `plugins/forge/hooks/guide.md` from 241 lines to ~100-120 lines. Remove 6 sections of reference content (duplicated in skill files or guide-only reference docs). Restructure remaining content into 3 thematic sections for better agent navigation.

## Reference Files
- `docs/proposals/guide-simplify/proposal.md` — Source proposal with scope and decisions

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/hooks/guide.md` | Remove 6 sections, restructure into 3 themes |

## Acceptance Criteria
- [ ] guide.md is ~100-120 lines (from 241)
- [ ] 6 sections removed: Skill Workflow mermaid diagrams, Quick Mode details, Testing Lifecycle, Evaluation Parameter Exceptions, Knowledge Accumulation details, Auxiliary Skills table
- [ ] 3 thematic sections present: Directory Conventions, Execution Rules, Automation Config
- [ ] All remaining rules preserved accurately — quality gate, scope resolution, auto-config, all-completed hook, task-CLI flow
- [ ] No functional behavior change

## Hard Rules
- Only modify `plugins/forge/hooks/guide.md` — do not touch hooks.json, skill files, or any other file
- Preserve every factual rule in the remaining sections verbatim or with equivalent precision

## Implementation Notes
**Sections to remove entirely** (content exists in skill files or is reference-only):
1. "Skill Workflow" — two mermaid diagrams + prerequisite note (lines 32-66)
2. "Quick Mode" — mermaid + when-to-use comparison + differences list (lines 68-99)
3. "Testing Lifecycle" — 3-layer table + flow diagram (lines 169-182)
4. "Evaluation Parameter Exceptions" — table (lines 184-193)
5. "Knowledge Accumulation" — /learn usage + auto-extract triggers + /consolidate-specs (lines 195-222)
6. "Other Auxiliary Skills" — table (lines 224-232)

**Sections to keep and reorganize**:
- Directory Conventions (Rules + Project-Level Documents)
- Manifest description (simplify)
- Quality Gate Protocol (sequence, failure handling, docs-only skip)
- Scope Resolution (just verb scope logic)
- All-Completed Hook (steps, docs-only skip, fix-task creation)
- Auto-Behavior Configuration (auto block + table + defaults note)
- Task-CLI (simplified flow)

**Target structure**:
1. **Directory Conventions** — where things live (rules + project-level docs)
2. **Execution Rules** — quality gate + scope resolution + all-completed hook + task-CLI
3. **Automation Config** — auto-behavior block + table
