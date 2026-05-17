---
id: "7"
title: "forge test verify 契约断裂检测"
priority: "P1"
estimated_time: "3h"
dependencies: ["4"]
scope: "backend"
breaking: false
type: "feature"
mainSession: false
---

# 7: forge test verify 契约断裂检测

## Description

实现 `forge test verify` 命令：扫描 Contract 规范，对每个 Output 断言重新执行代码侦察，比较 Fact Table 中的实际值与 Contract 中声明的语义描述是否匹配。当 CLI 命令的输出格式、API endpoint schema 或 TUI Model 字段变更时，自动检测受影响的下游 Contract。

来源：proposal "契约断裂检测机制"和 Scope "`forge test verify`"。

## Reference Files
- `docs/proposals/contract-journey-test-model/proposal.md` — Source proposal（含 verify 报告输出示例）
- `forge-cli/pkg/e2e/` — Fact Table 机制
- `forge-cli/internal/cmd/testing*.go` — 现有 testing 子命令

## Acceptance Criteria

- [ ] `forge test verify` 扫描所有 Contract 规范文件（`tests/<journey>/_contracts/*.md`），对每个 Output 断言重新执行代码侦察
- [ ] 检测到断裂时，报告受影响的 Contract（列出 Contract 文件路径、断裂维度、期望值、实际值）
- [ ] 对 20+ 个未变更的 Contract 运行 verify，误报数 = 0
- [ ] 每次运行时自动重新采集 Fact Table（基于当前代码库），不依赖历史快照
- [ ] bootstrap 策略：Phase 1 结束时用 126+ 已知正确输出生成 Fact Table 快照作为 verify 自身的首次正确性验证

## Hard Rules

- verify 不修改任何文件，只读取和报告
- 每次运行时重新采集 Fact Table（不依赖缓存快照）
- 误报率是核心质量指标，必须 ≤ 0（20+ 未变更 Contract 上零误报）

## Implementation Notes

- verify 输出格式参考 proposal 中的示例：`BROKEN (N): <path> <dimension> expected X → actual Y` + `OK (M): ...` + `Summary`
- Fact Table 采集：运行 forge-cli 的 126+ 测试用例，收集 stdout/stderr 快照
- verify 自身准确性验证：Phase 1 结束时用已知正确映射做 bootstrap，后续每次运行都基于最新代码库
