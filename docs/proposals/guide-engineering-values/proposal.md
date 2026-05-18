---
created: 2026-05-18
author: "fanhuifeng"
status: Draft
---

# Proposal: guide.md 工程价值观层

## Problem

guide.md 作为 Forge 注入每次 Claude 会话的核心参考，目前是纯技术规范（目录约定、skill 工作流、quality gate 等），缺少对 agent 思维方式和工程行为的指导。这导致 agent 在执行任务时倾向于"拿到指令就动手"，而非先思考、保持简洁、精准变更、以目标驱动。

### Evidence

- Andrej Karpathy 观察到 LLM 编码 agent 的四大常见问题：不做假设检验、过度复杂化代码、不必要地改动相邻代码、不验证目标达成。他将这些提炼为四个工程原则，在 GitHub 上获得数万星。
- Forge 的 task executor subagent 在执行任务时，偶尔出现过度重构、引入未请求的功能、改动与任务无关的代码等行为——正是 Karpathy 描述的问题。

### Urgency

guide.md 是每个 session 的第一入口。加入价值观层的投入极小（~40 行），但对 agent 行为的改善是持续性的、每 session 生效的。越早加入，收益越大。

## Proposed Solution

在 guide.md 的 `# Forge Guide` 标题之后、`## Directory Conventions` 之前，插入一个 `<ENGINEERING-VALUES>` 价值观 preamble。内容基于 Karpathy 四原则，用词和例子适配 Forge 场景（task executor、quality gate、subagent 等）。现有技术内容完全不动。

### Innovation Highlights

直接采纳社区验证的 Karpathy 四原则模型，创新点在于 Forge 场景适配：
- "Simplicity First" 对应 task executor 的实现行为约束
- "Surgical Changes" 对应 quality gate 和 git diff 的变更范围
- "Goal-Driven Execution" 对应 Forge 的 task success criteria 和 verify 机制
- "Think Before Coding" 对应 subagent 的推理行为

## Requirements Analysis

### Key Scenarios

- **Task executor 执行任务时**：价值观层指导 subagent 在动手前先理解任务、保持最小实现、只改必要的代码、定义并验证成功标准
- **用户交互时**：价值观层指导主 agent 在模糊需求下主动澄清、在发现更简方案时回推
- **Quality gate 失败时**：价值观层确保修复是精准的，不引入无关改动

### Non-Functional Requirements

- Preamble 长度控制在 40 行以内，避免显著增加每次 session 的 context 负担
- 内容必须是 agent 可直接执行的行为准则，不是抽象口号

### Constraints & Dependencies

- guide.md 通过 SessionStart hook 注入，格式需兼容现有 hook 的 JSON 转义逻辑
- 现有技术内容零改动

## Alternatives & Industry Benchmarking

### Industry Solutions

Andrej Karpathy 的 `andrej-karpathy-skills` 仓库（GitHub 数万星）提供了一个 ~70 行的 CLAUDE.md 文件，包含四个工程原则。社区广泛采纳，证明简洁的价值观层对 agent 行为有显著改善。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | agent 缺乏行为指导，过度复杂化/不必要改动等问题持续存在 | Rejected: 成本低但问题持续 |
| 独立 philosophy.md 文件 | 自研 | 解耦 | 增加文件管理复杂度，需额外 hook 或引用机制 | Rejected: 过度工程化 |
| 原封翻译 Karpathy | andrej-karpathy-skills | 社区验证 | 语境不匹配 Forge 概念 | Rejected: 缺少场景适配 |
| **Forge 适配版 preamble** | Karpathy + Forge 定制 | 社区验证 + 场景适配，~40 行，零改动现有内容 | 轻微增加 context 长度 | **Selected: 最佳性价比** |

## Feasibility Assessment

### Technical Feasibility

纯文本编辑，无技术依赖。guide.md 已有的 SessionStart hook 无需修改。

### Resource & Timeline

单人可完成，预计 1 个 task。

### Dependency Readiness

无外部依赖。

## Scope

### In Scope

- guide.md 新增 `<ENGINEERING-VALUES>` preamble（~40 行）
- 四原则的 Forge 场景适配（用词、例子）
- preamble 放置于标题之后、Directory Conventions 之前

### Out of Scope

- 现有技术规则的任何修改
- 其他文件的变更
- 独立 philosophy/values 文件
- 超出四原则之外的额外哲学内容

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Preamble 过长导致 context 膨胀 | L | M | 控制在 40 行以内，每行有实际指导意义 |
| 原则过于抽象，agent 无法执行 | L | M | 每条原则附带 Forge 场景的具体行为指引 |
| 与 CLAUDE.md 中已有"思维方式"冲突 | M | L | 检查并协调两者关系，避免重复或矛盾 |

## Success Criteria

- [ ] guide.md 包含 `<ENGINEERING-VALUES>` preamble，位于标题和 Directory Conventions 之间
- [ ] Preamble 包含四个原则，每个原则有 Forge 场景适配的行为指引
- [ ] Preamble 总行数 ≤ 40 行
- [ ] 现有技术内容零改动（git diff 仅包含 preamble 新增）
- [ ] SessionStart hook 正常注入（无 JSON 转义错误）

## Next Steps

- Proceed to `/quick-tasks` for task generation
