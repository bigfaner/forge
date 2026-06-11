---
id: "4"
title: "清理废弃模板、死代码和 TypeTestVerifyRegression 类型"
priority: "P0"
estimated_time: "2h"
dependencies: [1]
surface-key: "."
surface-type: "cli"
breaking: true
type: "coding.cleanup"
mainSession: false
---

# 4: 清理废弃模板、死代码和 TypeTestVerifyRegression 类型

## Description
删除废弃的模板文件和未使用的任务类型注册项，移除 `infer.go` 中 `T-quick-verify-regression` 死代码，清理 `types.go` 中 `TypeTestVerifyRegression` 常量及描述 "after graduation"，修复 `coding.fix.md` 中的 "E2E Fix Boundaries" section 和 `just test-e2e` 引用，修复 `autogen.go` 中 `T-validate-ux` 依赖链（应依赖最后一个 `run-test`），更新 `docsync_test.go` 中旧常量名引用。

## Reference Files
- `proposal.md#Layer-1-Go-代码层术语路径统一` — 第 6 项定义了所有需删除的模板、类型和死代码
- `proposal.md#Scope` — In Scope 第 5、9、10 项覆盖废弃代码清理、validate-ux 依赖修复、TypeTestVerifyRegression 清理
- `proposal.md#Success-Criteria` — 验证条件：TypeTestVerifyRegression 不存在、T-quick-verify-regression 死代码已移除、validate-ux 依赖正确

## Acceptance Criteria
- [ ] `pkg/prompt/data/test-verify-regression.md` 已删除
- [ ] `pkg/task/data/test-verify-regression.md` 已删除
- [ ] `infer.go` 中 `T-quick-verify-regression` 相关代码已移除
- [ ] `types.go` 中 `TypeTestVerifyRegression` 常量已删除
- [ ] `types.go` 中 `ValidTypes` / `SystemTypes` 注册项不含 `TypeTestVerifyRegression`
- [ ] `pkg/template/data/coding.fix.md` 中 "E2E Fix Boundaries" section 已更新
- [ ] `pkg/template/data/coding.fix.md` 中 `just test-e2e` 引用已更新
- [ ] `autogen.go` 中 `T-validate-ux` 依赖最后一个 `run-test` 任务（与 `T-validate-code` 同级）
- [ ] `docsync_test.go` 中旧常量名引用已更新
- [ ] `grep -rn "test-e2e\|E2E Fix" forge-cli/pkg/template/data/` 返回 0 结果
- [ ] `go build ./...` 和 `go test ./...` 通过

## Implementation Notes
- `TypeTestVerifyRegression` 从未被 autogen 生成，运行时无引用，可安全删除
- `coding.fix.md` 中的 "E2E Fix Boundaries" section 需替换为 surface-neutral 术语
- `T-validate-ux` 依赖修复需确认 `autogen.go` 中 run-test 任务的位置

### Integration Test Impact
- Affected test suite(s): `forge-cli/pkg/feature/`, `forge-cli/internal/`
- Expected fixture changes: `docsync_test.go` 断言更新、类型系统相关测试
- Risk level: medium
