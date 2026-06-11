---
created: 2026-05-19
author: "faner"
status: Draft
---

# Proposal: Forge Skill 指令层技术债清理

## Problem

Forge plugin 的 21 个 skills + 18 个 commands + 1 个 agent 共计 ~7600 行指令文本，经多轮迭代后积累了废弃标识、冗余描述和模糊边界，导致 LLM 上下文浪费且难以维护。

### Evidence

- `consolidate-specs/SKILL.md` 607 行，`tech-design/SKILL.md` 472 行——单个 skill 指令过长
- 多个 skill 内部存在冗余的流程描述（如质量门禁、scope 解析在单个 SKILL.md 中重复定义）。注意：根据 `docs/conventions/skill-self-containment.md` 的自洽原则，跨 skill 的指令重复是预期行为，不属于技术债
- 部分标识（如废弃的 `<TAG>` 格式、过时的引用路径）仍在文件中残留
- skill 之间职责边界模糊（如 `quick-tasks` vs `breakdown-tasks` 的重叠区域）
- **描述冗余造成概念混淆**：以 `noTest` 为例——`guide.md` 和多个 SKILL.md 将 `doc*` type prefix 与 `noTest: true` 并列描述，给人一种"noTest 是旧方案、应被 doc* 替代"的错觉。实际上两者分工不同：`doc*` 是 docs-only 任务的主要跳过触发器，`noTest: true` 是自动生成任务（testgen.go 产出的 T-test-1、T-eval-doc 等）的显式覆盖。历史上 `--no-test` CLI 标志确实已移除（commit `2d82321`），但 `Task.NoTest` 字段在 Go 代码中有 5 个消费者（submit.go、claim.go、build.go）仍在活跃使用。文档未区分这两个概念，导致 agent 理解偏差

### Urgency

v3.0.0 重构窗口期。当前分支已在做 skill 生态优化（test-profile-system、simplify-breakdown-tasks-prompt 等），趁此机会系统性清理技术债，避免债务随新功能继续积累。

## Proposed Solution

分两阶段执行：先审计全部指令文件识别共性技术债模式，再按模式分批清理。

### Innovation Highlights

模式驱动的批量清理策略——不是逐个 skill 盲改，而是先识别跨文件的 5 类技术债模式（废弃标识、单文件内冗余流程、概念混淆性描述、边界重叠、过时模板引用），再按模式统一修改。每个模式一个 commit，便于 review 和回滚。

#### 核心约束：Skill 自洽原则

根据 `docs/conventions/skill-self-containment.md`：每个 skill/command 必须逻辑自洽，读者通过单文件即可理解完整工作流，无需交叉引用其他 skill 或共享文档。因此：
- 跨 skill 的指令重复是**预期行为**，不属于本次清理目标
- 清理只针对单文件内的冗余（同一流程在文件内重复定义、解释性膨胀、废弃残留）
- 不将共享流程抽取到外部引用文件——这会破坏自洽性

#### 技术债模式示例：概念混淆性描述

`noTest` / `doc*` 的文档描述是典型案例。`noTest` 字段在 Go 代码中仍活跃（5 个消费者），但文档将其与 `doc*` type prefix 并列描述，暗示"noTest 是旧方案"。清理方式：在每个 SKILL.md 内明确两者的不同职责场景，而非写成"A 或 B"的替代关系。

## Requirements Analysis

### Key Scenarios

- Agent 加载 skill 后获得精简、无歧义的指令，无需自行过滤废弃内容
- 开发者维护 skill 时能快速理解职责边界，不与其他 skill 混淆
- 新 skill 创建时有清晰的职责定位参考

### Non-Functional Requirements

- 清理后每个 SKILL.md 行数减少 15-30%（高行数 skill 降幅更大）
- 清理后 agent 执行行为不偏离——只删冗余，不删关键指令

### Constraints & Dependencies

- 必须遵守 `docs/conventions/forge-distribution.md` 的分发模型约束
- 不改变 skill 的输入/输出契约（manifest、index.json 格式不变）
- 不涉及 Go 源码修改

## Alternatives & Industry Benchmarking

### Industry Solutions

技术债清理在工程领域是标准实践。通常分为：lint-based（自动化检测）、manual audit（人工审计）、pattern-driven（模式驱动批量处理）。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 债务持续积累，上下文浪费加剧 | Rejected: 成本会随迭代增长 |
| 逐个 skill 清理 | 常见做法 | 简单直接 | 同一模式重复修改，commit 碎片化 | Rejected: 效率低 |
| **模式审计 + 批量清理** | 重构最佳实践 | 全局视角，commit 按模式组织 | 前期审计耗时 | **Selected: 效率与可追溯性最优** |

## Feasibility Assessment

### Technical Feasibility

纯文本修改，无技术障碍。git 提供完整回滚能力。

### Resource & Timeline

审计阶段约需扫描 40 个文件，清理阶段按模式分批执行。预计 10 个以内任务可完成。

### Dependency Readiness

无外部依赖。

## Scope

### In Scope

- 21 个 `skills/*/SKILL.md` 主文件
- 18 个 `commands/*.md` 文件
- 1 个 `agents/task-executor.md` 文件
- skill 目录下的 `templates/` 和 `examples/` 中的指令性内容（不含纯输出模板）
- 跨 skill 的职责边界厘清

### Out of Scope

- Go 源码修改（`internal/`、`cmd/`）
- skill 的输入/输出契约变更
- `hooks/`、`references/`、`scripts/` 目录
- 新功能开发

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 过度精简导致 agent 执行偏差 | M | M | git 回滚；实际执行时自然暴露问题 |
| 边界厘清时误删跨 skill 协作指令 | L | M | 审计时标注协作点，清理时保留交叉引用 |
| templates 中的指令性内容被误判为纯模板 | L | L | 只清理明确冗余的指令部分 |

## Success Criteria

- [ ] 每个 SKILL.md 无废弃标识（如过时的 `<TAG>`、废弃的路径引用）
- [ ] 每个 skill 的职责在文件开头一句话可描述，且不与其他 skill 重叠
- [ ] 全部 SKILL.md 总行数减少 15%+（当前 6129 行 → 目标 5200 行以下）
- [ ] 清理后 `forge eval-forge` 审计无新增结构性问题

## Next Steps

- Proceed to `/quick-tasks` to generate tasks
