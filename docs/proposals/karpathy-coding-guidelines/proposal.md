---
created: 2026-05-20
author: "faner"
status: Draft
---

# Proposal: Karpathy Coding Guidelines — 四原则注入 Coding 提示词模板

## Problem

Forge 的 coding.* 系列任务提示词模板缺少系统性的行为守则，导致 task-executor agent 在执行任务时容易犯 Karpathy 所总结的 LLM 编码通病：不做假设验证就开干、过度工程化、顺手修改无关代码、缺乏可验证的成功标准。

### Evidence

分析现有 5 个 coding 模板（feature/enhancement/fix/refactor/cleanup），发现以下具体问题：

1. **feature/enhancement 从"读任务"直接跳到 TDD 实现** — 没有要求 agent 先陈述对需求的理解、显式列出假设、或在歧义处停下来澄清。如果 agent 误解了需求，TDD cycle 只会验证错误的理解。
2. **feature/enhancement 没有防范过度工程** — TDD 的 "implement minimal code to pass" 只约束了代码量，没有约束 agent 添加未被请求的功能、抽象或配置项。
3. **fix 模板的 "MINIMAL CHANGES" 太简略** — 只有两行规则，缺乏具体的行为指导（如"不要顺手改善相邻代码"）。
4. **所有模板缺乏显式的成功标准定义** — "Coverage >= 80%" 是间接指标，不是任务本身的可验证目标。

### Urgency

coding.* 模板是 Forge 任务执行的核心。每执行一个 coding 任务都依赖这些模板。当前模板的不足会导致 agent 频繁产生需要人工介入的"过度修改"——影响整个流水线的效率。

## Proposed Solution

将 Karpathy 四原则（Think Before Coding、Simplicity First、Surgical Changes、Goal-Driven Execution）作为行为守则注入 coding.* 模板和 fix-bug command。每个文件根据自身特点裁剪适用的原则子集，以 `<CODING_PRINCIPLES>` XML 标签（大写）包裹后内联写入。XML 标签提高 LLM 对守则的关注度和遵循率。

### 结构设计原则

注入不是简单粘贴，需要精心组织每个文件的整体结构：

1. **无时序冲突** — 守则中的指令不能与工作流步骤产生矛盾或执行顺序混乱。例如 "Think Before Coding" 守则指导 Step 1（Read Task）的行为，而非在 Step 2 之后才出现要求"先思考"
2. **消除语义重叠** — 与现有规则重叠时，合并为统一的守则表述，而非两段相同含义的文本并存。例如 coding-fix.md 的 "MINIMAL CHANGES" 和 "NO REFACTORING" 应被 Simplicity First + Surgical Changes 原则吸收替换
3. **位置精准** — 守则放在角色描述之后、工作流步骤之前，使其作为全局行为约束贯穿整个执行流程
4. **指令层次清晰** — XML 标签划分守则边界，与 `<IMPORTANT>`（任务硬规则）、`<HARD-GATE>`（不可绕过的检查点）形成清晰的指令层次：`<CODING_PRINCIPLES>` = 行为指南（自律遵循），`<IMPORTANT>` = 任务级硬约束（必须遵循），`<HARD-GATE>` = 流程级强制检查点

### Innovation Highlights

Karpathy 的四原则源自对 LLM 编码行为缺陷的观察总结，不是通用编码规范。这四条原则直接映射到 LLM 最容易犯的四类错误（错误假设、过度复杂、误触无关代码、缺乏验证闭环），是针对 agent 行为的精准矫正，而非人类编码规范的照搬。

## Requirements Analysis

### Key Scenarios

- agent 执行 feature 任务时，先验证对需求的理解再动手
- agent 在 enhancement 任务中不添加超出任务范围的"灵活性"
- agent 在 fix 任务中严格限制修改范围
- agent 在 refactor 任务中不顺便"改善"相邻代码
- agent 在 cleanup 任务中只清理目标内容

### Non-Functional Requirements

- 模板可读性不降低 — 原则文本简洁，不膨胀模板
- 不增加 agent 执行步骤 — 行为守则，非结构化步骤
- 向后兼容 — 不改变现有工作流的步骤编号和输出格式

### Constraints & Dependencies

- 模板文件通过 Go `//go:embed` 编译进 CLI 二进制 — 修改后需要重新编译
- 模板中的占位符（`{{TASK_ID}}` 等）不能被原则文本干扰

## Alternatives & Industry Benchmarking

### Industry Solutions

andrej-karpathy-skills 仓库提供了 CLAUDE.md 全局注入方式（所有任务统一生效）。Cursor 有类似的 rules 机制。但这些方案缺乏针对不同任务类型的裁剪能力。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | agent 继续犯同样的错误 | Rejected: 问题已验证存在 |
| CLAUDE.md 全局注入 | karpathy-skills | 简单、零代码 | 不区分任务类型，全部加载不相关原则 | Rejected: 粒度太粗 |
| 按模板裁剪 + 内联 | 本提案 | 精准、DRY 感知好 | 需要逐模板手动裁剪 | **Selected: 精准匹配每个模板的行为需求** |

## Feasibility Assessment

### Technical Feasibility

纯文本修改 5 个 `.md` 文件。无代码变更、无架构变更、无新依赖。

### Resource & Timeline

单一 PR，1 个 commit，5 个文件修改。

### Dependency Readiness

无外部依赖。

## Scope

### In Scope

- 修改 5 个 coding.* 模板 + 1 个 command 文件：
  - `forge-cli/pkg/prompt/data/coding-feature.md`
  - `forge-cli/pkg/prompt/data/coding-enhancement.md`
  - `forge-cli/pkg/prompt/data/coding-fix.md`
  - `forge-cli/pkg/prompt/data/coding-refactor.md`
  - `forge-cli/pkg/prompt/data/coding-cleanup.md`
  - `plugins/forge/commands/fix-bug.md`
- 原则到文件的映射（每个文件用 `<CODING_PRINCIPLES>` 大写 XML 标签包裹裁剪后的原则子集）：
  - **feature**: Think Before Coding + Simplicity First + Surgical Changes + Goal-Driven Execution
  - **enhancement**: Think Before Coding + Simplicity First + Surgical Changes + Goal-Driven Execution
  - **fix**: Think Before Coding + Simplicity First + Surgical Changes（替换现有简略规则）
  - **refactor**: Surgical Changes（已有 impact mapping 覆盖 Think）
  - **cleanup**: Simplicity First + Surgical Changes
  - **fix-bug**: Think Before Coding + Simplicity First + Surgical Changes + Goal-Driven Execution
- 更新 CLI 版本号（patch bump）

### Out of Scope

- 不修改 task-executor agent 定义
- 不修改其他 command 文件（execute-task、run-tasks 等）
- 不修改 task 文件模板（breakdown-tasks/quick-tasks 的 task.md）
- 不修改 clean-code.md（委托型模板，非 coding.*）
- 不修改非 coding 类型模板（doc.*、test.*、validation.*）
- 不修改 Go 合成引擎代码

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 原则文本增加 agent 上下文消耗 | L | L | 原则文本控制在 200 词以内每条 |
| 行为守则约束太强导致 agent 过于保守 | M | M | 原则文本中加入 "trivial tasks use judgment" 豁免 |
| fix 模板替换现有规则后行为变化 | L | M | 新原则覆盖旧规则语义，确保无损替换 |

## Success Criteria

- [ ] 5 个 coding 模板 + fix-bug command 均包含裁剪后的 Karpathy 原则文本，用 `<CODING_PRINCIPLES>` 大写 XML 标签包裹
- [ ] 原则与现有工作流步骤无时序冲突、无语义重叠（重叠处合并替换，非并存）
- [ ] coding-fix.md 的原有 "MINIMAL CHANGES" + "NO REFACTORING" 规则被 Simplicity First + Surgical Changes 原则替换，语义等价或更强
- [ ] fix-bug command 的原有 "Fix only what the failing tests require" 等规则与原则合并，不重复
- [ ] 模板步骤编号、输出格式、占位符均不变
- [ ] CLI 版本号 patch bump
- [ ] `forge prompt get-by-task-id` 对各 coding 类型生成的提示词包含原则文本

## Next Steps

- Proceed to `/quick-tasks` to generate tasks from this proposal
