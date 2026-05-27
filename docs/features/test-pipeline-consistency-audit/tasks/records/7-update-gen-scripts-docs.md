---
status: "completed"
started: "2026-05-27 14:51"
completed: "2026-05-27 15:32"
time_spent: "~41m"
---

# Task Record: 7 更新 gen-contracts 和 gen-test-scripts Skill 文档术语

## Summary
Updated gen-contracts and gen-test-scripts Skill documentation terminology: 'e2e 测试管道' -> 'Forge 测试管道', tests/e2e/ example paths -> tests/<journey>/, 'e2e tests' reference -> 'test files' in convention-guide

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-contracts/rules/journey-contract-model.md
- plugins/forge/skills/gen-test-scripts/rules/step-1-contract-loading.md
- plugins/forge/skills/gen-test-scripts/rules/convention-guide.md

### Key Decisions
无

## Document Metrics
3 files modified, 4 terminology replacements

## Referenced Documents
- docs/proposals/test-pipeline-consistency-audit/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] gen-contracts/ 下所有 Skill 文档中 'e2e 测试管道' 替换为 'Forge 测试管道'
- [x] gen-test-scripts/ 下所有 Skill 文档中 tests/e2e/ 旧路径替换为 tests/<journey>/
- [x] step-1-contract-loading.md 中 tests/e2e/step1_test.go 示例路径已更新
- [x] convention-guide.md 中 'e2e tests' 引用已替换
- [x] journey-contract-model.md 第 159 行 'language profile' 未被修改

## Notes
Old Structure code blocks in migration guide (journey-contract-model.md lines 166+) and HARD-RULE negation references in SKILL.md line 238 intentionally retained as they show what NOT to do.
