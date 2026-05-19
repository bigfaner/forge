---
id: "2"
title: "Update expert path references in SKILL.md"
priority: "P1"
estimated_time: "20m"
dependencies: ["1"]
type: "documentation"
mainSession: false
---

# 2: Update expert path references in SKILL.md

## Description
Update the 4 path references in `skills/eval/SKILL.md` that point to the old `agents/experts/` location. Since expert files now live within the eval skill's own directory, references use relative paths (same convention as `rubrics/<type>.md`) — no `${CLAUDE_SKILL_DIR}` prefix needed.

## Reference Files
- `docs/proposals/inline-experts-to-eval/proposal.md` — Source proposal (Path Updates table)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/eval/SKILL.md` | Replace 4 absolute path references with relative paths |

## Acceptance Criteria
- [ ] Dispatch table header: `${CLAUDE_SKILL_DIR}/../../agents/experts/scorer/` → `experts/scorer/`
- [ ] Step 2.1 scorer protocol path: `${CLAUDE_SKILL_DIR}/../../agents/experts/protocol/scorer-protocol.md` → `experts/protocol/scorer-protocol.md`
- [ ] Step 2.1 expert file example: `${CLAUDE_SKILL_DIR}/../../agents/experts/scorer/pm.md` → `experts/scorer/pm.md`
- [ ] Step 4.1 reviser protocol path: `${CLAUDE_SKILL_DIR}/../../agents/experts/protocol/reviser-protocol.md` → `experts/protocol/reviser-protocol.md`
- [ ] No remaining references to `agents/experts/` in SKILL.md

## Hard Rules
- Only change path strings — do not modify any surrounding logic or prose
- Use bare relative paths (`experts/scorer/...`), not `${CLAUDE_SKILL_DIR}` prefixed paths

## Implementation Notes
- The proposal provides an exact mapping table under "Path Updates in SKILL.md (Part A)"
- Search for `${CLAUDE_SKILL_DIR}/../../agents/experts/` to find all occurrences — there should be exactly 4
