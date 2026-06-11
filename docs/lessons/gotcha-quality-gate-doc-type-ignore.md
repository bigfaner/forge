---
created: "2026-05-22"
tags: [testing, architecture, local-dev-deployment]
---

# Gotcha: Quality Gate Protocol 忽略 task type，doc-only 任务被执行编译+测试

## Problem

`enforce-forge-task-add-in-loop` 特性的所有任务都是 `type: "doc"`（只修改 markdown 文件，无代码改动）。完成所有任务后，Stop hook 触发了 `task all-completed`，该命令运行了完整的 quality gate（`just compile → just fmt → just lint → just test → just test-e2e`），其中 e2e regression 因预存的无关问题而失败。

纯文档特性本应完全跳过编译和测试流程。

事后排查发现 `forge quality-gate`（新命令）已内置 `isDocsOnly()` 检查可正确跳过，但 Stop hook 调用的是旧命令 `task all-completed`，其实现缺少该跳过逻辑。

## Root Cause

因果链（4 层）：

1. **症状**: Stop hook 运行 `task all-completed` 后，触发了 e2e regression 并失败，即使所有任务都是 `type: doc`（纯文档变更）
2. **直接原因**: `task all-completed` 命令（`~/.zcode-task-cli/task` 二进制）内部没有 `isDocsOnly()` 跳过逻辑——它对所有 feature 都无条件执行 compile → test → e2e 三步门禁
3. **深层原因**: `task all-completed` 后来被重命名为 `forge quality-gate`，新命令（`~/.forge/bin/forge`）中新增了 `isDocsOnly()` 检查（基于 `IsTestableType()`），但旧 `task` 二进制的 `all-completed` 实现从未同步更新
4. **根因**: 两个独立二进制维护了同样的功能入口——`task all-completed` 和 `forge quality-gate` 是两套独立的 Go 编译产物，代码不同步。`task` 二进制没有 `isDocsOnly()`，`forge` 二进制有。Stop hook 在不经意间调用了旧命令

## Solution

方案分两层：

1. **Short-term**: Stop hook 统一使用 `forge quality-gate` 而非 `task all-completed`。`forge quality-gate` 内已内置 `isDocsOnly()` 检查，对纯文档任务直接跳过所有测试门禁（打印 "Feature is docs-only — skipping quality gate" 并 exit 0）
2. **Long-term**: 废弃 `task all-completed` 命令，确保所有调用路径都经过 `forge quality-gate`。从 hooks.json 确认已使用正确的命令名

## Reusable Pattern

当需要为 CI/pipeline 添加基于任务类型的条件分支时，遵循以下检查清单：

1. **检查所有入口点**：不仅仅是命令入口，还包括所有二进制中的同名/相似命令
2. **检查所有 layer**：`task` CLI 二进制 + `forge` CLI 二进制 + guide 文档 + hook 配置——每一层都可能有独立的 type 判断逻辑
3. **对 doc-only feature 的期望行为**：不执行编译（无可编译文件）、不执行单元测试（无代码变更）、不执行 e2e（无功能变更）
4. **命令重命名时，检查旧名称是否仍存在于其他二进制中**：`task all-completed` → `forge quality-gate` 的重命名遗漏了旧 `task` 二进制中的独立实现

## Example

```
# 当前 all-completed hook（无条件执行）：
just compile → just fmt → just lint → just test → just test-e2e

# 期望的 doc-only 行为：
# 跳过所有测试，只做最小验证
```

## Related Files

- `plugins/forge/hooks/hooks.json` — Stop hook 配置（`forge quality-gate` vs `task all-completed`）
- `forge-cli/internal/cmd/quality_gate.go` — `forge quality-gate` 实现（含 `isDocsOnly()`）
- `forge-cli/pkg/task/build.go` — `testableTypes` map + `IsTestableType()`
- `forge-cli/internal/cmd/submit.go` — `IsTestableType()` 使用处
- `~/.forge/bin/forge` — 有 `isDocsOnly()` 的新 binary
- `~/.zcode-task-cli/task` — **无** `isDocsOnly()` 的旧 binary

## References

- [Gotcha: docs-only proposals need code path audit](gotcha-docs-only-needs-code-audit.md)
- [Gotcha: quality gate cross-feature pollution](gotcha-quality-gate-cross-feature-pollution.md)