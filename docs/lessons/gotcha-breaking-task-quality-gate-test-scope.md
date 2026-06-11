---
created: "2026-05-29"
tags: [architecture, testing, interface]
---

# breaking=true 任务的 quality gate 与 deferred test fixes 存在结构性矛盾

## Problem

Task 2（`breaking: true`）执行器被卡住。它在执行核心重构之外，就地修复了 25 个测试失败（autogen_test.go 1089 行 diff），而不是按照 task 文件的 "deferred to task 6" 指示跳过测试修复。最终仍未通过所有测试，无法提交。

## Root Cause

因果链（5 层）：

1. **L1 症状**: Task 2 executor 就地修复测试而非 defer，导致范围膨胀（1089 行测试修改）并卡住

2. **L2 直接原因**: task-executor 的 quality gate 要求 `go test` 全部通过才能提交。删除 `GetBreakdownTestTasks`/`GetQuickTestTasks` 后，autogen_test.go 和 autoconfig_test.go 中的所有测试引用了不存在的函数，产生编译错误和测试失败。Executor 面临两个选择：(a) 接受测试失败但无法通过 quality gate 提交；(b) 就地修复测试以通过 quality gate。它选了 (b)。

3. **L3 设计矛盾**: Task planning 模型和 task execution 模型存在结构性冲突：
   - **Planning 模型**承认 breaking 变更会导致测试失败，通过 `breaking: true` 标记 + 专用 test-fix task（task 6）处理
   - **Execution 模型**的 quality gate 不支持"预期内测试失败"——它要求 ALL tests pass
   - Task 文件写了 "deferred to task 6" 但这只是文本提示，不是执行约束——quality gate 不读取也不尊重这个指示

4. **L4 缺失机制**: task-executor 没有 "expected failures" 的概念。缺少：
   - 排除特定测试文件的 quality gate 配置（如 `--exclude-test-files=autogen_test.go,autoconfig_test.go`）
   - `breaking: true` 标记对 quality gate 行为的影响（breaking tasks 应该只要求 `go build` 通过，不要求 `go test` 全过）
   - 对 "deferred to task N" 指示的机器可读编码（当前只是 Implementation Notes 中的自由文本）

5. **L5 根因**: Quality gate 的语义是"代码质量保证"，但对于 breaking refactor tasks，正确的质量标准是"编译通过 + 核心逻辑正确"，而非"所有测试通过"。当前 gate 将这两者混为一谈。

## Solution

**短期（task 2 当前状态）**: 丢弃 executor 的测试修改，只保留核心重构（autogen.go 删除 + pipeline.go GenerateTestTasks + build.go 调用方更新）。测试修复留给 task 6。

**中期（task-executor 改进）**: `breaking: true` 的 task 应有两种 quality gate 模式：
- `go build ./...` 必须通过（编译正确性）
- `go test ./...` 的失败仅创建 fix task（不阻塞提交），且 fix task 指向 deferred test-fix task

**长期（planning ↔ execution 契约）**: 在 task frontmatter 中添加 `exclude-tests` 字段，作为 quality gate 的机器可读排除列表。这使 "deferred to task 6" 从文本提示升级为执行约束。

## Reusable Pattern

当规划 breaking refactor 时：

1. **Breaking task 的 quality gate 应仅要求编译通过**，不要求测试全过。测试修复是独立 task 的职责
2. **Task 文件中的 "deferred" 指示是文本提示，不是执行约束**。如果需要 executor 遵守，必须编码为 frontmatter 字段或 quality gate 配置
3. **Executor 遇到大量测试失败时会默认修复**——这是"让所有测试变绿"的惯性。必须在 task executor prompt 中显式声明 breaking tasks 的测试修复策略
4. **范围隔离**：breaking task 只改生产代码，test-fix task 只改测试代码。两者不应由同一个 executor 同时处理

## Example

```
# 当前行为（有问题）
Task 2: 删除 GetBreakdownTestTasks → 25 tests fail
→ Executor 修 1089 行测试代码 → 卡住 → 无法提交

# 正确行为
Task 2: 删除 GetBreakdownTestTasks → 25 tests fail (expected)
→ Executor 只确保 go build 通过 → 提交（status: completed）
→ Task 6: 专门修复所有测试 → 全绿 → 提交
```

## Related Files

- `docs/features/pipeline-topology-registry/tasks/2-refactor-task-generation.md` — task 文件，含 "deferred to task 6" 指示
- `plugins/forge/agents/task-executor/` — task-executor prompt，缺少 breaking task 的测试策略
- `forge-cli/pkg/task/autogen_test.go` — 被 executor 过度修改的测试文件（1089 行 diff）
