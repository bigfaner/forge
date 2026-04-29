---
id: "3.gate"
title: "Phase 3 Exit Gate"
priority: "P0"
estimated_time: "30min"
dependencies: ["3.summary"]
status: pending
breaking: true
---

# 3.gate: Phase 3 Exit Gate

## Description

Exit verification gate for Phase 3. Confirms no raw language-specific commands remain in the 7 build/test files.

## Verification Checklist

1. [ ] `grep -c 'project-test-command\|npx tsx' plugins/forge/commands/fix-bug.md` = 0
2. [ ] `grep -c 'go test\|npm test\|pytest' plugins/forge/commands/run-tasks.md` = 0
3. [ ] `grep -c 'go test\|npm test\|pytest\|npm run build' plugins/forge/agents/task-executor.md` = 0
4. [ ] `grep -c 'go test\|npm test\|pytest\|npm run build' plugins/forge/agents/error-fixer.md` = 0
5. [ ] `grep -c 'project-specific verification' plugins/forge/commands/execute-task.md` = 0
6. [ ] `grep -c 'go test\|npm test.*coverage\|pytest --cov' plugins/forge/skills/record-task/SKILL.md` = 0
7. [ ] `grep -c 'just test' plugins/forge/skills/improve-harness/SKILL.md` >= 1

## Reference Files

- All 7 Phase 3 files (see design/tech-design.md Phase 3)

## Acceptance Criteria

- [ ] All verification checklist items pass
- [ ] Record created via `/record-task` with `coverage: -1.0`

## Implementation Notes

Verification-only. Fix inline if trivial.
