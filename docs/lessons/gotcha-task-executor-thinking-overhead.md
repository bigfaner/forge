---
created: "2026-05-27"
tags: [architecture, testing]
---

# Task Executor 87% 时间消耗在 Thinking 上

## Problem

task-executor 执行 task 1（常量重命名）耗时 18.3 分钟，其中 thinking 占 967s（87%），工具调用仅 151s（13%）。88 次工具调用中有大量重复 grep/find/read。

## Root Cause

因果链（3 层）：

1. **表面现象**：18.3 分钟执行一个中等复杂度的重命名任务，远超预估的 2h
2. **直接原因**：
   - agent 做了大量冗余探索——多次搜索相同模式（`E2ETestsBaseDir`、`tests/e2e` 等），每次 grep 后再 thinking 推理下一步
   - `forge task submit` 被调用 2 次（42.6s + 37.9s），第一次因预先存在的 `TestDetectModeFromPath` 失败被 quality gate 拦截
   - 任务范围过大，合并了常量重命名、路径扁平化、目录移动、testrunner 路径确认等多个独立步骤
3. **根因**：
   - **任务粒度过粗**：一个 task 涵盖 4 个可独立验证的步骤，agent 需要理解每个步骤的完整影响范围
   - **quality gate 无法区分预先存在的失败**：`TestDetectModeFromPath` 与 task 1 变更无关，但 submit 质量门无法区分，导致额外的 submit 重试（80.5s）
   - **executor prompt 缺乏搜索策略指引**：没有 "先批量收集引用点，再批量修改" 的约束，agent 倾向于逐个 grep→think→grep 的循环模式

## Solution

1. 拆分粗粒度任务为独立可验证的步骤（常量重命名 vs 目录移动 vs testrunner 路径）
2. 在 task executor prompt 中加入搜索策略指引：先一次性收集所有引用点，再批量修改，避免重复搜索
3. quality gate 通过 `git stash` 对比基线来区分 "当前变更引入的失败" 和 "预先存在的失败"

## Reusable Pattern

- **任务拆分原则**：如果一个 In Scope 条目包含 "并" 或多个独立动词（重命名 + 扁平化 + 确认），应拆成独立任务
- **搜索效率**：对代码库做 grep/find 时，先列出所有需要搜索的模式，一次性批量搜索，避免 "搜索→发现→再搜索" 的串行循环
- **quality gate 基线**：提交前先跑一次基线测试记录已知失败，变更后再跑一次，仅对比增量失败

## Related Files

- `docs/forensics/task1-duration/evidence/evidence.json` — 完整证据数据
- `forge-cli/internal/cmd/quality_gate.go` — quality gate 实现
