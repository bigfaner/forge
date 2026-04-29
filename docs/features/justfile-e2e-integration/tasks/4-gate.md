---
id: "4.gate"
title: "Phase 4 Exit Gate"
priority: "P0"
estimated_time: "30min"
dependencies: ["4.summary"]
status: pending
breaking: true
---

# 4.gate: Phase 4 Exit Gate

## Description

Final exit verification gate. Confirms all 13 in-scope files have been updated and no raw commands remain anywhere in `plugins/forge/`.

## Verification Checklist

1. [ ] `grep -c 'just test-e2e' plugins/forge/skills/breakdown-tasks/templates/run-e2e-tests.md` >= 1
2. [ ] `grep -c 'just e2e-verify' plugins/forge/skills/breakdown-tasks/templates/gen-test-scripts.md` >= 1
3. [ ] `grep -c 'just test-e2e' plugins/forge/skills/breakdown-tasks/templates/fix-e2e.md` >= 1
4. [ ] Full sweep: `grep -rn 'npx tsx\|cd tests/e2e && npm\|project-test-command' plugins/forge/` = 0 lines
5. [ ] `grep -rn 'just e2e-setup\|just e2e-verify\|just test-e2e\|just test\|just build' plugins/forge/` >= 20 lines total

## Reference Files

- All 13 in-scope files (see design/tech-design.md)

## Acceptance Criteria

- [ ] All verification checklist items pass
- [ ] Record created via `/record-task` with `coverage: -1.0`

## Implementation Notes

Verification-only. This is the final gate before T-test tasks begin.
