---
id: "fix-3"
title: "Fix: skill-ops TestCleanCode_SkillFile_ContainsFivePrinciples pre-existing failure"
priority: "P0"
estimated_time: "30min"
dependencies: []
status: pending
breaking: true
type: "coding.fix"
---

# Fix: skill-ops TestCleanCode_SkillFile_ContainsFivePrinciples pre-existing failure

## Root Cause

TestCleanCode_SkillFile_ContainsFivePrinciples fails because clean-code SKILL.md no longer contains 'Focus Scope' principle. Pre-existing failure also present in original forge-cli/tests/skill-ops/. Fix assertion to match current SKILL.md content.

## Reference Files

- Source: tests/skill-ops/clean_code_skill_test.go
- Test script: cd tests && go test -tags=cli_functional ./skill-ops/... -run TestCleanCode_SkillFile_ContainsFivePrinciples
- Test results: SKILL.md missing 'Focus Scope' principle

## Surface Inference

This fix-task was created by the quality-gate hook. If `surface-key` and `surface-type` above are empty, infer them at execution time:

1. Parse `tests/skill-ops/clean_code_skill_test.go` to extract the first file path (comma-separated).
2. Run `forge surfaces --json <file-path>` to resolve surface-key/type.
3. Use the resolved surface-type to load the appropriate `rules/surfaces/<type>.md` for test orchestration guidance.

If `forge surfaces --json` fails (no surfaces configured, command not found), proceed without surface information — this does not block the fix.

## Fix Boundaries

When fixing test failures, observe these boundaries:

**Forbidden:**
- Starting dev server (`npx expo start`, `npm run dev`, etc.)
- Running `npm install` more than 3 times — mark task as blocked if dependency installation fails 3 times
- Running full test suite — regression is verified by the dispatcher after fix completes
- Manually opening browser to verify rendering

**Correct workflow:**
1. Read failing test + corresponding component source
2. Compare test's expected testID/selectors vs actual DOM structure
3. Modify component (add testID) or test (adjust selectors/assertions)
4. Run targeted tests on affected packages — unit tests must pass
5. Record completion

## Verification

After fixing, verify the fix works:
1. Run targeted tests on changed packages: `go test -race ./affected/package/...`
2. Replace the path with the actual packages you modified

> **Note:** Full project-wide tests run at CLI submit (`forge task submit`) — agent runs targeted tests only.

Full regression is verified by the dispatcher, not by this fix task.

When this task is recorded as completed via `task record`, the source task 5 is automatically restored to pending if all its dependencies are completed.
