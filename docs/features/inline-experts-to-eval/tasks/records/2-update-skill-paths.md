---
status: "completed"
started: "2026-05-19 11:13"
completed: "2026-05-19 11:18"
time_spent: "~5m"
---

# Task Record: 2 Update expert path references in SKILL.md

## Summary
Updated 4 absolute path references in skills/eval/SKILL.md to use bare relative paths (experts/scorer/...) instead of ${CLAUDE_SKILL_DIR}/../../agents/experts/... prefix. Dispatch table header, scorer protocol path, expert file example, and reviser protocol path all updated. No remaining references to agents/experts/ in SKILL.md.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/eval/SKILL.md

### Key Decisions
无

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] Dispatch table header: ${CLAUDE_SKILL_DIR}/../../agents/experts/scorer/ -> experts/scorer/
- [x] Step 2.1 scorer protocol path: ${CLAUDE_SKILL_DIR}/../../agents/experts/protocol/scorer-protocol.md -> experts/protocol/scorer-protocol.md
- [x] Step 2.1 expert file example: ${CLAUDE_SKILL_DIR}/../../agents/experts/scorer/pm.md -> experts/scorer/pm.md
- [x] Step 4.1 reviser protocol path: ${CLAUDE_SKILL_DIR}/../../agents/experts/protocol/reviser-protocol.md -> experts/protocol/reviser-protocol.md
- [x] No remaining references to agents/experts/ in SKILL.md

## Notes
Only path strings were changed; no surrounding logic or prose modified. The remaining ${CLAUDE_SKILL_DIR} reference on line 136 points to rubrics/validate-ux-pipeline.md, unrelated to experts.
