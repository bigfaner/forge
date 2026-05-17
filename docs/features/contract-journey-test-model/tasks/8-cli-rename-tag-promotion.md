---
id: "8"
title: "CLI 重命名 + Tag-Based Promotion"
priority: "P1"
estimated_time: "2h"
dependencies: ["5", "6"]
scope: "backend"
breaking: true
type: "refactor"
mainSession: false
---

# 8: CLI 重命名 + Tag-Based Promotion

## Description

将 forge-cli 的 `testing` 子命令重命名为 `test`，所有子命令行为不变（仅命令前缀变更）。废弃 graduate-tests，替换为 Tag-Based Promotion——通过标签管理测试生命周期，消除文件迁移的复杂度。

来源：Scope "forge-cli `testing` 命令重命名和适配"和"废弃 graduate-tests，替换为 Tag-Based Promotion"。

## Reference Files
- `docs/proposals/contract-journey-test-model/proposal.md` — Source proposal
- `forge-cli/internal/cmd/testing*.go` — 现有 testing 子命令
- `forge-cli/internal/cmd/graduate*.go` — 现有 graduate-tests 命令
- `docs/conventions/forge-distribution.md` — Forge 分发模型

## Acceptance Criteria

- [ ] `forge testing` 重命名为 `forge test`，所有子命令行为不变（`test detect`、`test interfaces`、`test run` 等功能等价，仅命令前缀变更）
- [ ] `forge test promote <journey>` 将该 Journey 下所有 `@feature` 标签替换为 `@regression`
- [ ] promote 替换前后用 `git diff` 确认仅标签变更、无其他代码改动
- [ ] promote 执行前自动运行该 Journey 的所有测试，全部通过才更新标签；存在失败测试时拒绝晋升并输出失败报告
- [ ] CI 选择通过 `forge test run --tags regression` 或 `--tags feature` 实现，框架原生标签过滤机制执行（Go `-run`、pytest `-m`、Playwright `--grep`）
- [ ] 现有单步 TC 作为单步骤 Journey 退化形式继续工作；所有现有 e2e 测试通过（`just e2e-compile` 编译通过 + 现有 126+ 测试用例输出与改动前 diff 为空）

## Hard Rules

- `breaking: true`——`testing` → `test` 是破坏性变更，需要更新所有文档和 CI 脚本中的引用
- graduate-tests 完全移除，不保留兼容 shim
- 标签以语言框架原生方式嵌入（Go `//go:build`、Python `@pytest.mark`）

## Implementation Notes

- Tag-Based Promotion 消除 graduate-tests 中的文件移动、import 重写、去重、分类复杂度
- 标签生命周期：`@feature`（新生成，验证中）→ `@regression`（已验证，回归测试）
- 向后兼容：现有 e2e 测试提供回归安全网，`just e2e-compile` 编译通过确保无破坏
