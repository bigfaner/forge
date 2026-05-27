---
status: "completed"
started: "2026-05-27 19:21"
completed: "2026-05-27 19:37"
time_spent: "~16m"
---

# Task Record: 4 清理废弃模板、死代码和 TypeTestVerifyRegression 类型

## Summary
清理废弃模板、死代码、TypeTestVerifyRegression 类型和修复依赖链

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/infer.go
- forge-cli/pkg/task/types.go
- forge-cli/pkg/task/autogen.go
- forge-cli/pkg/template/data/coding.fix.md
- forge-cli/pkg/prompt/prompt_test.go
- forge-cli/pkg/task/autogen_test.go
- forge-cli/pkg/task/types_test.go
- forge-cli/pkg/task/build_test.go
- forge-cli/pkg/task/autoconfig_test.go
- forge-cli/pkg/task/infer_test.go
- forge-cli/pkg/task/category_test.go
- forge-cli/pkg/template/template_test.go

### Key Decisions
- 删除 TypeTestVerifyRegression 常量及所有相关注册项（从未被 autogen 生成，运行时无引用）
- autogen.go 中 T-validate-ux 添加对最后一个 run-test 的依赖（与 T-validate-code 同级）
- coding.fix.md 中 'E2E Fix Boundaries' 替换为 'Fix Boundaries'，使用 surface-neutral 术语

## Test Results
- **Tests Executed**: Yes
- **Passed**: 1151
- **Failed**: 0
- **Coverage**: 87.7%

## Acceptance Criteria
- [x] pkg/prompt/data/test-verify-regression.md 已删除
- [x] pkg/task/data/test-verify-regression.md 已删除
- [x] infer.go 中 T-quick-verify-regression 相关代码已移除
- [x] types.go 中 TypeTestVerifyRegression 常量已删除
- [x] types.go 中 ValidTypes/SystemTypes 注册项不含 TypeTestVerifyRegression
- [x] pkg/template/data/coding.fix.md 中 'E2E Fix Boundaries' section 已更新
- [x] pkg/template/data/coding.fix.md 中 just test-e2e 引用已更新
- [x] autogen.go 中 T-validate-ux 依赖最后一个 run-test 任务
- [x] docsync_test.go 中旧常量名引用已更新
- [x] grep -rn 'test-e2e|E2E Fix' forge-cli/pkg/template/data/ 返回 0 结果
- [x] go build ./... 和 go test ./... 通过

## Notes
所有测试编译和运行通过。共修改 3 个源码文件和 9 个测试文件。删除 2 个废弃模板文件 (test-verify-regression.md) 和 TypeTestVerifyRegression 常量及其在 4 处注册项。
