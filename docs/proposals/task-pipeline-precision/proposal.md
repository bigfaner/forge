---
created: "2026-05-27"
author: "faner"
status: Approved
---

# Proposal: Task Pipeline Precision Tuning

## Problem

Task executor 在简单任务上浪费大量时间——87% 时间花在冗余探索上，25 分钟做一个纯文本替换——因为任务生成和 prompt 模板机制缺乏精度控制：任务粒度过粗、模板忽略复杂度差异、Reference Files 指向过宽范围导致越界修改。

### Evidence

4 个独立 but 相互关联的 lesson 记录了具体故障：

| Lesson | 故障现象 | 根因 |
|--------|---------|------|
| gotcha-task-executor-thinking-overhead | 18.3 min 完成一个 constant rename（87% 思考时间） | 任务粒度过粗（4 个独立步骤合并为 1），无搜索策略引导，重复 grep 同一模式 |
| gotcha-prompt-template-complexity-agnostic | 25 min 完成一个 6 AC 的纯文本替换 | `coding-enhancement.md` 对所有 enhancement 任务强制完整 Step 1 + Step 1.5 spec-code scan，无视复杂度 |
| gotcha-quick-tasks-merge-threshold | 11 条 AC 中 2 条未完成却被标为完成 | 合并标准是时间估算（<30min），不是可独立验证性 |
| gotcha-task-reference-files-scope-creep | executor 越界修改 task 4 的文件 | Reference Files 指向 `proposal.md#Layer-1-...`，executor 读取后发现"spec vs code 不一致"就主动修复，无视任务边界 |

### Urgency

每次 quick mode 执行都受这些问题影响。executor 效率直接决定开发者的等待时间和 API 成本。4 个 lesson 在同一次 feature 执行中发现，说明问题是系统性的而非偶发。

## Proposed Solution

从任务生成到执行的 3 个环节注入精度控制：

1. **任务生成环节**：拆分规则从"时间估算"改为"可独立验证"标准；Reference Files 从 proposal section 指针改为内联精确信息
2. **任务元数据环节**：task frontmatter 加 `complexity: low/medium/high` 字段，生成时自动判定
3. **执行环节**：prompt template 根据 complexity 字段分支（low 跳过 spec-code scan、简化探索）；加搜索策略引导——具体内容为在 implementation 步骤前插入指令段落："在修改任何文件前，先用 Grep/Glob 搜索所有需要修改的位置，收集完整清单后再执行修改。禁止边搜边改。"该引导位于 Step 1 之后、Step 2 之前，与 Step 1.5（spec-code scan）互补：Step 1.5 是规范一致性扫描，搜索策略引导是操作顺序约束

### Innovation Highlights

**Complexity-aware prompt routing**：不是业界标准做法——大多数 AI coding agent 对所有任务使用同一 prompt。这里通过在任务生成阶段就标注复杂度，让执行阶段自适应调整探索深度，减少简单任务的 token 浪费。

**Self-contained task documents**：Reference Files 内联化使每个 task 文件成为自包含的执行单元，executor 不需要（也不应该）访问 proposal 全文。这是最小权限原则在任务执行中的应用。

**Scope boundary declaration**：每个 task 文件头部嵌入显式 scope 边界声明（`## Scope Boundary` 段落），列出允许修改的文件路径列表。当 spec-code scan 发现的"不一致"涉及 scope 边界外的文件时，该声明直接覆盖 scan 结果的修复倾向。这是信息隔离之外的主动防御机制——即使 executor 读到了越界信息，scope 声明提供了明确的不修改指令。

## Requirements Analysis

### Key Scenarios

- **简单机械替换**（low complexity）：6 个 grep+sed 级 AC，无 Hard Rules，Reference Files ≤ 1 → 跳过 spec-code scan，简化 Step 1，直接批量修改
- **中等增强任务**（medium complexity）：3-6 AC，有 Reference Files → 标准 Step 1 + 跳过 Step 1.5 spec-code scan
- **复杂新功能**（high complexity）：>6 AC 或有 Hard Rules → 完整流程（Step 1 + Step 1.5 + Step 2）
- **多动词 In Scope 条目**：一个 proposal bullet 包含"rename + flatten + confirm"→ 按 functional boundary 拆分为独立任务
- **Reference Files 生成**：quick-tasks/breakdown-tasks 为每个 coding task 生成内联精确信息（文件路径 + 具体修改描述）

### Non-Functional Requirements

- **向后兼容**：现有 index.json 中无 complexity 字段的任务默认为 `medium`，不受影响
- **零运行时开销**：complexity 判定在任务生成时完成，不增加执行时开销
- **探索效率**：complexity: low 任务从收到 prompt 到开始修改文件的探索阶段应 < 30s（软目标，非硬性门控）
- **模板一致性**：5 个 coding template 的复杂度分支逻辑保持统一结构

### Constraints & Dependencies

- prompt template 位于 `forge-cli/pkg/prompt/data/`，通过 Go `embed.FS` 加载
- `prompt.go` 的 `renderTemplate()` 使用 `strings.ReplaceAll` 做占位符替换，需新增 `{{COMPLEXITY}}` 占位符
- task type → template 映射通过命名约定（`coding.enhancement` → `coding-enhancement.md`），无需修改
- `renderTemplate()` 是纯文本替换引擎，不支持条件块。实现"low 跳过 Step 1.5"需扩展后处理：在 template 中用标记注释（如 `<!-- IF NOT_LOW -->...<!-- END_IF -->`）包裹 Step 1.5 段落，由 `cleanTemplateOutput()` 根据 complexity 值删除标记块。这比引入模板引擎更轻量，且保持向后兼容（无标记的 template 不受影响）

## Alternatives & Industry Benchmarking

### Industry Solutions

大多数 AI coding agent（Cursor, Copilot, Aider）对任务不区分复杂度——统一 prompt 或完全依赖模型自行判断。本方案的 complexity-aware routing 是一个更结构化的方法。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零改动成本 | executor 效率持续低下，token 浪费 | Rejected: 4 个 lesson 证明问题系统性存在 |
| Prompt 内置启发式规则 | lesson suggestion B | 不改 task 文件格式 | 模板膨胀；每次执行都做判定 | Rejected: 判定在生成时做更高效 |
| **显式 complexity 字段 + 内联 Reference** | lesson suggestion A | 精确、可审计、零运行时开销 | 需改 task 模板和生成逻辑 | **Selected: 判定一次，执行时零成本** |

## Feasibility Assessment

### Technical Feasibility

改动涉及配置层和代码层两部分。配置层（SKILL.md、template .md）为纯文本改动，无风险。代码层需修改完整数据管道以传递 complexity 字段：① `FrontmatterData` struct 加 `Complexity` 字段 → ② `Task` struct 加 `Complexity` 字段 → ③ `index.json` schema 更新 → ④ `renderTemplate()` 新增 `{{COMPLEXITY}}` 占位符。这 4 处改动链路固定、改动量小，但需同步修改，遗漏任一环节会导致字段丢失。

### Resource & Timeline

~10 个文件改动，均为配置/模板层。预计 6-10 个 coding task，适合 quick mode。

### Dependency Readiness

无外部依赖。所有修改的文件均在 plugins/forge/ 和 forge-cli/ 内部。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 任务类型决定探索深度 | XY Detection | Overturned: 同一 type 内复杂度差异巨大（纯文本替换 vs 多文件重构）。复杂度，而非类型，才是探索深度的决定因素 |
| Reference Files 指向 proposal section 提供完整上下文 | 5 Whys | Overturned: 指针式引用让 executor 看到超出任务范围的 spec 需求，导致越界修改。自包含的内联信息更安全 |
| <30min 的步骤应该合并 | Assumption Flip | Overturned: 合并标准应是"能否独立验证"而非时间估算。4 个独立步骤各 <30min 但合并后不可独立验证 |

## Scope

### In Scope

- quick-tasks SKILL.md 拆分规则优化：合并标准改为"可独立验证"、AC 上限 6 条、多动词检测规则
- breakdown-tasks SKILL.md 同步修改拆分规则（具体段落：Task Splitting Rules 段落加"可独立验证"标准、AC 上限 6 条、多动词检测；task.md frontmatter 模板加 complexity 字段；Reference Files 生成段落改为内联格式）
- task frontmatter 加 `complexity` 字段（quick-tasks 和 breakdown-tasks 的 task.md 模板）
- quick-tasks/breakdown-tasks 复杂度自动判定逻辑：以静态指标（AC 数量 + Hard Rules + Reference Files 数量）作为默认启发式，同时在 SKILL.md 中提供 LLM 判断指引（如"涉及多文件架构变更时提升 complexity 等级"），允许任务生成阶段的 LLM 在静态指标与认知判断冲突时覆盖默认值
- Reference Files 生成策略改为内联精确信息（quick-tasks 和 breakdown-tasks）
- 4 个 coding prompt templates（enhancement/refactor/new-feature/test-gen）加复杂度分支（low 跳过 Step 1.5，简化 Step 1）和搜索策略引导。`coding.fix` 模板不纳入 complexity routing——fix-task 由 dispatcher 自动生成且具有固定的 5-step 流程，不需要复杂度分支
- prompt.go 传递 complexity 字段到模板渲染
- 移除 quick-tasks 15 coding task 上限（SKILL.md 和 /quick 命令中的 15 task 限制）

### Out of Scope

- Template 合并/体系重组（适合单独 proposal）
- 质量门基线测试改进（thinking-overhead lesson 第三方案，独立关注点）
- task-executor agent 定义本身（搜索引导在 prompt template 层）
- 现有 proposal 合并或替代（slim-task-prompt-templates, prompt-template-audit）
- /quick 命令本身的 15 task 上限（仅改 quick-tasks skill）
- ~~移除 quick-tasks 15 coding task 上限（属于独立架构决策，与精度控制无关，应另提 proposal）~~

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| complexity 判定不准确，简单任务被标为 high | M | M | 使用保守启发式（AC>6 或有 Hard Rules 才标 high），默认 medium |
| 内联 Reference Files 使 task 文件过长 | L | L | 限制每个 task 的 Reference Files 条目数（≤5），且内联信息应简洁（1-2 行/条） |
| 内联 Reference Files 在 proposal 变更时产生 stale reference | M | M | 每条内联信息附带溯源标注（`source: proposal.md#section-id`），执行时如发现内联信息与实际代码不符，以代码为准并标记 stale |
| prompt template 复杂度分支增加模板维护成本 | M | L | 分支逻辑统一模板化（一段条件块，4 个 template 复制相同结构；fix 模板不参与） |
| breakdown-tasks 同步修改引入不一致 | M | M | 两个 SKILL.md 使用相同的判定规则描述，确保逻辑一致 |

## Success Criteria

- [ ] quick-tasks 生成的任务中，AC ≤ 3 且无 Hard Rules 且 Reference Files ≤ 1 的任务被标记为 `complexity: low`
- [ ] complexity: low 的任务执行时跳过 Step 1.5 spec-code scan（通过 `forge prompt get-by-task-id` 输出验证不包含 Step 1.5 段落）
- [ ] 所有 coding task 的 Reference Files 为内联格式（文件路径 + 修改描述），无 `proposal.md#` 指针
- [ ] 无任务包含 > 6 条 AC（如果 proposal bullet 自然产生 >6 AC，必须拆分）
- [ ] 搜索策略引导出现在所有 4 个非 fix coding template 的 implementation 步骤前
- [ ] `forge prompt get-by-task-id` 输出包含 complexity 对应的流程分支内容
- [ ] 现有无 complexity 字段的 index.json 任务执行时默认为 medium，行为不变

## Next Steps

- Proceed to `/quick-tasks` to generate tasks from this proposal
