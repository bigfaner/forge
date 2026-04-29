---
id: "2.gate"
title: "Phase 2 Exit Gate"
priority: "P0"
estimated_time: "30min"
dependencies: ["2.summary"]
status: pending
breaking: true
---

# 2.gate: Phase 2 Exit Gate

## Description

Exit verification gate for Phase 2. Confirms no raw e2e commands remain in gen-test-scripts and run-e2e-tests skill files.

## Verification Checklist

1. [ ] `grep -c 'npx tsx\|npx playwright\|npm install' plugins/forge/skills/run-e2e-tests/SKILL.md` = 0
2. [ ] `grep -c 'npx playwright install\|cd tests/e2e' plugins/forge/skills/gen-test-scripts/SKILL.md` = 0
3. [ ] `grep -c 'just e2e-setup\|just test-e2e' plugins/forge/skills/run-e2e-tests/SKILL.md` >= 2
4. [ ] `grep -c 'just e2e-verify\|just e2e-setup' plugins/forge/skills/gen-test-scripts/SKILL.md` >= 2
5. [ ] gen-test-scripts Step 4 prose contains "exit 1 = skill incomplete" note

## Reference Files

- `plugins/forge/skills/run-e2e-tests/SKILL.md`
- `plugins/forge/skills/gen-test-scripts/SKILL.md`

## Acceptance Criteria

- [ ] All verification checklist items pass
- [ ] Record created via `/record-task` with `coverage: -1.0`

## Implementation Notes

Verification-only. Fix inline if trivial. Document deviations.
