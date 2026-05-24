---
status: "completed"
started: "2026-05-24 10:26"
completed: "2026-05-24 10:42"
time_spent: "~16m"
---

# Task Record: 1 Auto-discovery: 替换 prompt.go 和 autogen.go 手写 map + init-time 校验 + clean-code.md 重命名

## Summary
Replaced hand-written typeToTemplate and autogenTypeToFile maps with naming convention auto-discovery (strings.ReplaceAll(typeName, '.', '-') + '.md'). Renamed prompt/data/clean-code.md to code-quality-simplify.md. Added ValidatePromptTemplates() and ValidateAutogenTemplates() init-time validation in CLI entry point (Run function). All existing tests pass.

## Changes

### Files Created
- forge-cli/pkg/prompt/data/code-quality-simplify.md

### Files Modified
- forge-cli/cmd/forge/run.go
- forge-cli/pkg/prompt/prompt.go
- forge-cli/pkg/task/autogen.go
- forge-cli/pkg/task/autogen_test.go

### Key Decisions
- Split validation into two functions (ValidatePromptTemplates + ValidateAutogenTemplates) because prompt and autogen packages have separate embed.FS instances covering different type subsets
- Kept ValidTypes check in Synthesize() to preserve 'unknown type' error message for backward compatibility
- Validation functions skip types not present in their respective FS (graceful handling of split FS coverage)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 409
- **Failed**: 0
- **Coverage**: 86.8%

## Acceptance Criteria
- [x] prompt.go typeToTemplate map removed, Synthesize() uses naming convention
- [x] autogen.go autogenTypeToFile map removed, GenerateTestTaskMD() uses naming convention
- [x] prompt/data/clean-code.md renamed to code-quality-simplify.md
- [x] CLI entry has ValidateTemplateConventions() via ValidatePromptTemplates + ValidateAutogenTemplates
- [x] Missing template or mapping collision causes CLI startup failure
- [x] grep -r 'clean-code' forge-cli/ has no residual clean-code.md references
- [x] All existing tests pass

## Notes
Type code-quality.simplify maps to code-quality-simplify.md via convention. Validation runs in Run() not init() per hard rule. Tests in cmd/forge now exercise the validation path.
