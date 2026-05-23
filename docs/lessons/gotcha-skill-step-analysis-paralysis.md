---
created: "2026-05-23"
tags: [architecture, testing]
---

# Skill Step Analysis Paralysis

## Problem

执行 `/breakdown-tasks` 时，在完成文档读取和规则加载后，陷入长时间的思考循环，反复推敲次要细节（Convention 文件分发路径、distribution model 细节、scope 分配策略），迟迟不进入"写文件"阶段。最终在大量上下文消耗后仍未产出一个任务文件。

## Root Cause

**Level 1 — 过度分析次要细节**：花了大量时间推敲 Convention 文件到底应该在 `docs/conventions/testing/` 还是 plugin 分发目录，以及 scope 应该是 `all` 还是 `backend`，而不是基于 tech-design 明确给出的文件路径直接写任务。

**Level 2 — 违反 skill 步骤顺序**：skill 定义了清晰的 Step 0→1→2→...→6 流程，但实际执行时试图在 Step 4（写任务文件）之前解决所有 Step 5/6/7 才需要验证的问题。`forge task index` 和 `forge task validate-index` 正是用来兜底校验的工具，应该在写完文件后运行，而不是在写之前预设完美。

**Level 3 — 完美主义陷阱**：试图让每个任务文件在第一次写入时就完全正确，忽略了 skill 的设计意图——先写后验，用 CLI 工具校验。

## Solution

严格按照 skill 定义的步骤顺序执行。完成 Step 1-3 的分析后，直接进入 Step 4 写任务文件。tech-design 给出的文件路径是权威来源，不要因为"想确认"而反复探查代码库。写完后用 `forge task index` 和 `forge task validate-index` 校验。

## Reusable Pattern

**执行多步骤 skill 时：**
1. 按步骤顺序执行，不跳步、不提前解决后续步骤的问题
2. 设计文档给出的路径和结构是权威来源，直接使用
3. 利用 skill 提供的 CLI 校验工具（index、validate-index）做兜底，而不是在写之前追求完美
4. 当发现自己在思考块中反复推敲同一类细节超过 2 轮时，停下来，按照当前最佳判断写文件，后续校验会捕获问题

## Related Files

- `docs/features/test-capability-v2/design/tech-design.md` — PRD Coverage Map 提供了权威的文件路径
- `C:\Users\panda\.claude\plugins\cache\forge\forge\3.0.0-rc.18\skills\breakdown-tasks\SKILL.md` — 步骤定义
