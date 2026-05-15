---
created: "2026-05-15"
tags: [architecture, testing, local-dev-deployment]
---

# Stop hook quality-gate 对 docs-only 特性的任务污染

## Problem

一个 docs-only 特性（只改 markdown 文件，无代码改动）完成所有任务后，Stop hook 触发 `forge quality-gate`，该命令运行项目级全量测试（compile → fmt → lint → test → e2e），在某个预存在的、与本特性完全无关的 Go 测试失败（`pkg/just/TestEnsureJust_UserAccepts_PkgManagerSuccess`）后，向该特性的 `index.json` 添加了多个任务：

- fix-2（P0, breaking）— 修复预存在的测试失败
- T-quick-1 到 T-quick-5（5 个测试流程任务）— 因 fix task 引入 breaking=true 后重新索引生成

原本 3 个任务的特性膨胀到 10 个，且 fix task 的根因与当前特性无关。

## Root Cause

因果链（4 层）：

1. **触发条件**：`forge task submit` 在最后一个任务完成时写入 `.forge/state.json` 的 `allCompleted: true`
2. **直接原因**：Stop hook 读取该信号，无条件运行 `forge quality-gate`
3. **设计缺陷**：quality-gate 运行项目级全量测试（`just test`），不区分当前特性是否涉及代码改动。它对 docs-only 特性和 full-stack 特性执行相同的门禁
4. **级联效应**：`addFixTask()` 将 fix task 加入当前特性的 index.json。fix task 有 `breaking: true`，改变了任务的类型组成（从纯文档变为混合），后续索引重建时生成了 T-quick 测试任务替代原来的 T-eval-doc

**核心问题**：quality-gate 的测试范围（project-wide）与 fix task 的归属范围（per-feature）不匹配。预存在的不相关测试失败会"污染"当前特性的任务列表。

## Solution

**短期**：确认测试失败与当前特性无关后，手动标记 fix task 为 skipped：
```bash
forge task status <fix-task-id> skipped
```

**长期方向**（需评估）：
- quality-gate 在运行前检查特性的任务类型组成，如果全部是 documentation 类型，跳过 test/e2e 步骤（仅保留 lint 检查文档格式）
- 或：fix task 根据 git diff 确认失败测试是否在当前特性的改动范围内，不在则不阻塞当前特性

## Reusable Pattern

当 Stop hook 报告 quality-gate 失败时，先确认两件事再决定是否修复：
1. **失败是否与当前特性相关**：运行 `git diff --name-only` 看改动的文件，与测试失败的包对比
2. **失败是否为预存在问题**：`git stash` 后运行 `just test <failing-package>`，如果也失败则说明是预存在的

如果两者都是"否"，标记 fix task 为 skipped 而不是花时间修复不相关的测试。

## Example

```bash
# 确认是预存在的、不相关的失败
git stash
just test ./forge-cli/pkg/just/...  # 如果仍失败 → 预存在
git stash pop

# 确认与当前特性无关
git diff --name-only  # 只看到 plugins/forge/skills/quick-tasks/ 下的 .md 文件

# 标记 fix task 为 skipped
forge task status fix-2 skipped
```

## Related Files

- `forge-cli/internal/cmd/quality_gate.go` — quality-gate 实现
- `forge-cli/internal/cmd/submit.go` — `saveIndexAndSignalCompletion()` 设置 allCompleted
- `plugins/forge/hooks/hooks.json` — Stop hook 配置

## References

- `docs/lessons/gotcha-quality-gate-fix-task-loop.md` — 相关的 fix-task 循环问题
- Forge Guide > All-Completed Hook 章节
