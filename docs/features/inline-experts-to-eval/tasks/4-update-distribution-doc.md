---
id: "4"
title: "Update forge-distribution.md for new expert location"
priority: "P2"
estimated_time: "20m"
dependencies: ["1"]
type: "documentation"
mainSession: false
---

# 4: Update forge-distribution.md for new expert location

## Description
Update `docs/conventions/forge-distribution.md` to reflect the new location of expert files under `skills/eval/experts/` instead of `agents/experts/`.

## Reference Files
- `docs/proposals/inline-experts-to-eval/proposal.md` — Source proposal (Distribution Doc Updates)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `docs/conventions/forge-distribution.md` | Update directory tree, component table, and section 3 |

## Acceptance Criteria
- [ ] Directory tree: `agents/experts/` subtree moved under `skills/eval/`, removed from `agents/`
- [ ] Component table: remove `agents/experts/` from agents row description; add `experts/` mention under skills row
- [ ] Section "3. 核心依赖 → agents/experts/": update title and all paths to reflect new location under `skills/eval/experts/`
- [ ] No remaining references to `agents/experts/` in the document

## Hard Rules
- Only update path references and descriptions — do not change the document's structure or add new sections

## Implementation Notes
- The agents/ component row should still mention `task-executor.md` but no longer reference `experts/`
- The skills/ component row should now mention that eval skill contains `experts/` subdirectory
