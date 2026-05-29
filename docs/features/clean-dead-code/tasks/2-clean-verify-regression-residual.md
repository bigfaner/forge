---
id: "2"
title: "清理 verify-regression 残留推断规则和陈旧测试引用"
priority: "P1"
estimated_time: "30m"
complexity: "medium"
dependencies: []
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.cleanup"
mainSession: false
---

# 2: 清理 verify-regression 残留推断规则和陈旧测试引用

## Description
`test.verify-regression` 任务类型的常量和注册已从 `types.go` 清除，但 3 处残留代码仍存在：(1) `infer.go` 中 2 条推断规则仍返回 `"test.verify-regression"` 字符串；(2) `autoconfig_test.go` 中 2 处错误消息引用 `T-test-verify-regression`；(3) `quality_gate_test.go` 中测试数据使用 `T-quick-verify-regression`。这些残留代码产生不可用的类型字符串和误导性测试断言。

## Reference Files
- `forge-cli/pkg/task/infer.go`: 移除返回 `"test.verify-regression"` 的 case 分支（T-test-verify-regression 和 T-quick-verify-regression）(source: proposal.md#Scope-B)
- `forge-cli/pkg/task/autoconfig_test.go`: 清理行 264、280 中 T-test-verify-regression 陈旧错误消息 (source: proposal.md#Scope-B)
- `forge-cli/internal/cmd/quality_gate_test.go`: 清理行 1330 中 T-quick-verify-regression 测试数据 (source: proposal.md#Scope-B)

## Acceptance Criteria
- [ ] `infer.go` 中 2 条 verify-regression 推断规则已移除，`grep "verify-regression" forge-cli/pkg/task/infer.go` 无结果
- [ ] `autoconfig_test.go` 中 `T-test-verify-regression` 陈旧引用已清理
- [ ] `quality_gate_test.go` 中 `T-quick-verify-regression` 测试数据已清理
- [ ] `go test ./pkg/task/... ./internal/cmd/...` 全部通过

## Implementation Notes
- infer.go 推断规则返回的 `"test.verify-regression"` 不在 TaskTypeRegistry 中，任何下游代码都会拒绝该类型——移除规则不会影响任何有效行为
- 清理测试引用时需确保不破坏测试断言逻辑——如果断言检查的是错误消息中包含该字符串，移除后需同步更新期望值
- grep 全量搜索 `verify-regression` 确保无遗漏：`grep -r "verify-regression" forge-cli/ --include="*.go"`
