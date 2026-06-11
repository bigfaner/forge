---
created: 2026-05-19
author: faner
status: Draft
intent: enhancement
---

# Proposal: quick-tasks 自动提交并无缝衔接 run-tasks

## Problem

`/quick` 管线中 quick-tasks (Step 3) 生成规划产物后、run-tasks (Step 4) 执行前，存在两个断裂：

1. **规划产物未提交**——task .md 文件、index.json、manifest.md 在 git 中保持 untracked 状态。run-tasks 的 task executor 只提交自己修改的代码文件，不提交预先存在的规划文档。管线结束后，规划文件仍然 untracked。
2. **Step 3→4 衔接不够明确**——/quick 命令已顺序调用两个 skill，但缺少"生成后立即执行，无中间停顿"的显式约束，agent 可能在 Step 3 完成后停顿或输出摘要。

### Evidence

- `docs/lessons/gotcha-quick-tasks-no-commit.md` 详细记录了此问题：auto-consolidate-specs feature 中 3 个 task 文件有 2 个未被提交
- quick-tasks SKILL.md 有 Step 0-7 但没有 commit step
- /quick 命令 Step 3→4 之间无衔接指令

### Urgency

每次使用 /quick 都会产生 untracked 文件，污染 git 状态。低风险但高频。

## Proposed Solution

两处修改：

1. **quick-tasks 新增 Step 8 (Commit)**：在 Step 7 验证通过后，提交所有生成的规划产物
2. **/quick 命令收紧 Step 3→4 衔接**：明确指示 agent 在 quick-tasks 完成后立即启动 run-tasks，不输出中间摘要

### Innovation Highlights

无创新，纯运维修复。遵循"谁创建谁持久化"原则。

## Requirements Analysis

### Key Scenarios

- **Happy path**: /quick 管线完成后，所有规划产物已提交，代码变更已提交，git status 干净
- **Dirty working tree**: quick-tasks Step 8 commit 时存在其他未提交变更——只 add 规划产物，不 stage 其他文件
- **Validation failure**: Step 7 失败时不执行 commit，先修复

### Constraints & Dependencies

- 修改 plugin 文件需遵循 `docs/conventions/forge-distribution.md`
- commit 消息遵循 Conventional Commits

## Alternatives & Industry Benchmarking

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | untracked 文件持续积累 | Rejected: 已知 gotcha |
| run-tasks 提交规划文件 | — | 集中在执行器 | 违反"谁创建谁持久化"原则 | Rejected: 职责错位 |
| **quick-tasks Step 8 提交 + /quick 收紧衔接** | gotcha 教训 | 最小改动，职责清晰 | 不覆盖 breakdown-tasks（同问题） | **Selected** |
| breakdown-tasks 也加 commit | — | 全面修复 | 超出 /quick 管线范围 | Deferred: 独立处理 |

## Feasibility Assessment

### Technical Feasibility

完全可行。两处文档修改，不涉及代码逻辑变更。

### Resource & Timeline

1 个任务，<30min。

## Scope

### In Scope

- quick-tasks SKILL.md 新增 Step 8 (Commit)
- /quick 命令 Step 3→4 衔接收紧

### Out of Scope

- breakdown-tasks commit gap（同问题但不同 skill）
- 独立调用 /quick-tasks 时的 commit 行为
- run-tasks 修改

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Commit 时 working tree 有冲突文件 | Low | Low——git add 只 stage 指定文件 | Step 8 使用 `git add` 指定路径，不用 `git add -A` |
| /quick 衔接指令被 agent 忽略 | Low | Medium——回到现状 | 使用 EXTREMELY-IMPORTANT 标记 |

## Success Criteria

- [ ] quick-tasks SKILL.md 包含 Step 8，执行 `git add` 指定路径 + commit
- [ ] /quick 命令 Step 3→4 之间有明确"立即执行"指令
- [ ] /quick 管线完成后 `git status` 无 untracked 规划文件

## Next Steps

- Proceed to `/quick-tasks` to generate tasks
