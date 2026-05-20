---
status: "completed"
started: "2026-05-20 13:14"
completed: "2026-05-20 13:18"
time_spent: "~4m"
---

# Task Record: 2 Inject coverage target into prompt synthesis

## Summary
在 prompt 合成流程中注入覆盖率目标。renderTemplate() 新增 {{COVERAGE_STRATEGY}} 和 {{COVERAGE_TARGET}} 占位符替换，按优先级解析：frontmatter coverage > config per-type > built-in default。coding 类型模板(feature/enhancement/fix/refactor/cleanup)添加覆盖率占位符，非 testable 类型不注入。

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/prompt/prompt.go
- forge-cli/pkg/prompt/prompt_test.go
- forge-cli/pkg/prompt/data/coding-feature.md
- forge-cli/pkg/prompt/data/coding-enhancement.md
- forge-cli/pkg/prompt/data/coding-fix.md
- forge-cli/pkg/prompt/data/coding-refactor.md
- forge-cli/pkg/prompt/data/coding-cleanup.md

### Key Decisions
- 新增 resolveCoverage() 函数集中处理覆盖率解析逻辑，按 frontmatter > config > default 优先级
- coding 类型模板在 CODING_PRINCIPLES 后注入 COVERAGE_STRATEGY/COVERAGE_TARGET 标签
- 非 testable 类型通过 IsTestableType() 判断，占位符替换为空字符串后由 cleanTemplateOutput() 清理残留标签
- SynthesizeOpts 不添加新字段，覆盖率从 config 和 task 数据实时读取（遵循 Hard Rules）

## Test Results
- **Tests Executed**: Yes
- **Passed**: 43
- **Failed**: 0
- **Coverage**: 89.4%

## Acceptance Criteria
- [x] renderTemplate() 新增 {{COVERAGE_TARGET}} 占位符替换
- [x] 新增 {{COVERAGE_STRATEGY}} 占位符，值为 percentage 或 maintain
- [x] {{COVERAGE_TARGET}} 对于 percentage 策略渲染为 "达到 N% 测试覆盖率"
- [x] {{COVERAGE_TARGET}} 对于 maintain 策略渲染为 "保持现有覆盖率，下降不超过 2%"
- [x] 优先级正确：frontmatter coverage > config per-type > built-in default
- [x] 对于非 testable 类型（doc*、gate 等），不注入覆盖率相关占位符
- [x] 现有测试通过

## Notes
无
