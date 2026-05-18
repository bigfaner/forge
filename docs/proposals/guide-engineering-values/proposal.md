---
created: 2026-05-18
author: "fanhuifeng"
status: Draft
---

# Proposal: guide.md 价值观层 + 结构精简

## Problem

guide.md 作为 Forge 注入每次 Claude 会话的核心参考（241 行），存在两个问题：

**P1: 缺少工程价值观指导。** 当前是纯技术规范（目录约定、skill 工作流、quality gate 等），缺少对 agent 思维方式和工程行为的指导。这导致 agent 在执行任务时倾向于"拿到指令就动手"，而非先思考、保持简洁、精准变更、以目标驱动。

**P2: 结构冗余。** 内容堆叠无清晰分组，部分信息重复，heading 层级不一致。

### Evidence

- Andrej Karpathy 观察到 LLM 编码 agent 的四大常见问题：不做假设检验、过度复杂化代码、不必要地改动相邻代码、不验证目标达成。他将这些提炼为四个工程原则，在 GitHub 上获得数万星。
- Forge 的 task executor subagent 在执行任务时，偶尔出现过度重构、引入未请求的功能、改动与任务无关的代码等行为——正是 Karpathy 描述的问题。
- guide.md 具体结构问题（241 行审计）：
  1. **H3 孤儿节点**：`Evaluation Parameter Exceptions`、`Knowledge Accumulation`、`Other Auxiliary Skills` 三个 `###` 悬空在 `## Testing Lifecycle` 和 `## Task-CLI` 之间，无父级 `##`
  2. **信息重复**：`Auto-Behavior Configuration` 的 YAML 示例（14 行）和表格（9 行）表达完全相同的信息，YAML 即是 defaults 的权威源
  3. **过度详细**：Quick Mode 差异列表提及具体 task 编号（T-quick-1~4, T-quick-specs-1 等），guide.md 层面只需说"简化测试管线"
  4. **关联内容分离**：`Testing Lifecycle`（3 层测试表格）与 `Quality Gate Protocol`（测试执行规则）紧密相关但分属不同 `##` section

### Urgency

guide.md 是每个 session 的第一入口。价值观层投入极小（~40 行），对 agent 行为的改善是持续性、每 session 生效的。结构精简减少 context 占用，提升注入效率。两项一起做避免后续二次修改。

## Proposed Solution

两项改动：

**S1: 新增价值观 preamble。** 在 `# Forge Guide` 标题之后、`## Directory Conventions` 之前，插入 `<ENGINEERING-VALUES>` 价值观层（~40 行）。基于 Karpathy 四原则，Forge 场景适配。

**S2: 结构精简。** 修复 4 个结构问题，零信息损失：
- 修复 H3 孤儿节点 → 重新归组到合理的 `##` 下
- 合并 Auto-Behavior 的 YAML+表格为纯 YAML（删除冗余表格）
- Quick Mode 差异列表简化（去掉具体 task 编号，保留语义描述）
- Testing Lifecycle 合并到 Quality Gate Protocol

### Innovation Highlights

价值观层：直接采纳社区验证的 Karpathy 四原则模型，Forge 场景适配：
- "Think Before Coding" → subagent 推理行为 + brainstorm 交互
- "Simplicity First" → task executor 实现、quick mode 简化管线
- "Surgical Changes" → quality gate 变更范围、git diff 审查
- "Goal-Driven Execution" → task success criteria、verify 机制

结构精简：以信息论视角审查——每条信息在 guide.md 中只出现一次。

## Requirements Analysis

### Key Scenarios

- **Task executor 执行任务时**：价值观层指导 subagent 在动手前先理解任务、保持最小实现、只改必要的代码、定义并验证成功标准
- **用户交互时**：价值观层指导主 agent 在模糊需求下主动澄清、在发现更简方案时回推
- **Quality gate 失败时**：价值观层确保修复是精准的，不引入无关改动

### Non-Functional Requirements

- Preamble 长度控制在 40 行以内，避免显著增加每次 session 的 context 负担
- 内容必须是 agent 可直接执行的行为准则，不是抽象口号
- 结构精简后总行数 ≤ 当前 241 行 + preamble 40 行（即净减少来自精简的行数应 ≥ preamble 新增的行数）

### Constraints & Dependencies

- guide.md 通过 SessionStart hook 注入，格式需兼容现有 hook 的 JSON 转义逻辑
- 结构精简为零信息损失：不删除任何语义内容，只做重组和去重

## Alternatives & Industry Benchmarking

### Industry Solutions

Andrej Karpathy 的 `andrej-karpathy-skills` 仓库（GitHub 数万星）提供了一个 ~70 行的 CLAUDE.md 文件，包含四个工程原则。社区广泛采纳，证明简洁的价值观层对 agent 行为有显著改善。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | agent 缺乏行为指导，guide.md 结构问题持续存在 | Rejected: 两项问题均未解决 |
| 只加 preamble，不改结构 | Karpathy 采纳 | 最小变更 | 结构冗余问题未解决，context 继续膨胀 | Rejected: 放弃了精简的机会 |
| 独立 philosophy.md 文件 | 自研 | 解耦 | 增加文件管理复杂度，需额外 hook 或引用机制 | Rejected: 过度工程化 |
| 原封翻译 Karpathy | andrej-karpathy-skills | 社区验证 | 语境不匹配 Forge 概念 | Rejected: 缺少场景适配 |
| **价值观 preamble + 结构精简** | Karpathy + Forge 定制 + 信息论 | 社区验证 + 场景适配 + 减少冗余，净行数接近零增长 | 两项改动需要一起验证 | **Selected: 最佳性价比，一次修改到位** |

## Feasibility Assessment

### Technical Feasibility

纯文本编辑，无技术依赖。guide.md 已有的 SessionStart hook 无需修改。

### Resource & Timeline

单人可完成，预计 2 个 task（价值观 preamble + 结构精简）。

### Dependency Readiness

无外部依赖。

## Scope

### In Scope

- guide.md 新增 `<ENGINEERING-VALUES>` preamble（~40 行），位于标题和 Directory Conventions 之间
- 四原则的 Forge 场景适配（用词、例子）
- 修复 H3 孤儿节点（Evaluation Parameter Exceptions、Knowledge Accumulation、Other Auxiliary Skills）→ 重新归组
- Auto-Behavior Configuration 去重（删除冗余表格，保留 YAML）
- Quick Mode 差异列表简化（去掉具体 task 编号）
- Testing Lifecycle 合并到 Quality Gate Protocol

### Out of Scope

- 其他文件的变更
- 独立 philosophy/values 文件
- 超出四原则之外的额外哲学内容
- guide.md 中任何语义内容的删除（只做重组和去重）
- SessionStart hook 的修改

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Preamble 过长导致 context 膨胀 | L | M | 控制在 40 行以内，每行有实际指导意义 |
| 原则过于抽象，agent 无法执行 | L | M | 每条原则附带 Forge 场景的具体行为指引 |
| 与 CLAUDE.md 中已有"思维方式"冲突 | M | L | 检查并协调两者关系，避免重复或矛盾 |
| 结构精简丢失语义 | L | H | 精简前后做 diff 审查，确保零信息损失 |

## Success Criteria

- [ ] guide.md 包含 `<ENGINEERING-VALUES>` preamble，位于标题和 Directory Conventions 之间
- [ ] Preamble 包含四个原则，每个原则有 Forge 场景适配的行为指引
- [ ] Preamble 总行数 ≤ 40 行
- [ ] H3 孤儿节点全部归组到合理的 `##` section 下
- [ ] Auto-Behavior Configuration 无冗余（YAML 或表格二选一）
- [ ] Quick Mode 差异列表无具体 task 编号
- [ ] Testing Lifecycle 内容合并到 Quality Gate Protocol
- [ ] guide.md 总行数 ≤ 270 行（241 原始 + 40 preamble - ≥11 精简）
- [ ] SessionStart hook 正常注入（无 JSON 转义错误）

## Next Steps

- Proceed to `/quick-tasks` for task generation
