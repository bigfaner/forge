---
created: 2026-05-20
author: "faner"
status: Approved
---

# Proposal: 逐 Skill 瘦身——拆分、精简、消歧

## Problem

Forge plugin 的 22 个 SKILL.md 文件总计 4421 行，其中多个文件（consolidate-specs 348 行、gen-test-scripts 325 行、init-justfile 327 行）虽在 350 行阈值以内，但混合了流程指令、业务规则和内联模板，导致 LLM 上下文浪费且维护困难。以 consolidate-specs 为例：348 行 SKILL.md + 282 行辅助文件，SKILL.md 中仍含可进一步拆出的规则解释文本。同时部分 skill 存在指令歧义（如 noTest vs doc* 概念混淆），增加 agent 执行偏差风险。

### Evidence

- 实际行数（2026-05-20 统计）：consolidate-specs 348 行、gen-test-scripts 325 行、init-justfile 327 行为前三大文件，均接近 350 行阈值
- 所有 22 个 SKILL.md 均在 350 行以下，但多个文件（eval 200 行 + 2295 行辅助文件、gen-test-cases 136 行 + 843 行辅助文件、ui-design 228 行 + 1759 行辅助文件）的内容结构仍有精简和消歧空间
- `guide.md` 和多个 SKILL.md 对 `noTest`/`doc*` 的描述产生歧义（例如 `noTest` 一词在 guide.md 中指"跳过测试生成"，但在部分 skill 中被解读为"该 skill 不涉及测试"）

### Urgency

v3.0.0 重构窗口期。已有 5 个瘦身相关提案均未执行——方向分散、范围过大是主因。需要一个可立即落地的增量方案。

## Proposed Solution

按大小分层、逐组处理：大文件（400+ 行）独立拆分，中/小文件按领域分组合并处理。每个任务聚焦一组 skill，依次完成拆分结构、精简行数、消除歧义三项目标。

### Splitting Heuristic

拆分决策遵循以下具体规则，而非依赖隐喻或执行时临时判断：

**留在 SKILL.md 的内容**（流程骨架层）：
1. 所有步骤编号及其描述（如 "步骤 1: 收集输入 → 步骤 2: 分析 → ..."）
2. 条件分支逻辑（如 "如果 X 则执行 A，否则执行 B"）
3. 输入/输出契约定义（skill 接受什么、产出什么）
4. 对 rules/ 和 templates/ 的引用指令（如 "参照 rules/disambiguation.md 中的术语定义"）

**移至 rules/ 的内容**（规则细节层）：
1. 超过 5 行的规则定义和解释性文本（如业务规则详细说明、约束条件列表）
2. 术语定义和消歧文档
3. 命名约定、路径规范等参考性内容

**移至 templates/ 的内容**（模板资源层）：
1. 超过 10 行的输出模板（如 markdown 输出格式模板、报告结构模板）
2. 可复用的代码片段或配置模板
3. 示例输入/输出

**边界规则**：当一段内容同时包含流程指令和规则细节时，流程指令保留在 SKILL.md，规则细节移至 rules/ 并在原位置添加引用路径。

### Worked Example: consolidate-specs

**拆分前**（348 行 SKILL.md + 282 行辅助文件）：
```
SKILL.md (348 行) — 混合流程指令和规则细节
├── 已有 10 个辅助文件（282 行），但 SKILL.md 仍含可拆出的规则解释文本

辅助文件 (10 files, 282 lines):
├── rules/ 目录下的规则定义
└── templates/ 目录下的输出模板
```

**拆分后目标**（SKILL.md 精简至 ~200 行，辅助文件扩充）：
```
SKILL.md (~200 行)
├── 元数据 + 触发条件 (保留)
├── 输入/输出契约 (保留)
├── 步骤流程描述 (保留，精简至核心指令)
├── 引用: 辅助文件的路径指令
└── 错误处理摘要 (精简)

辅助文件 (扩充) — 从 SKILL.md 移出的规则细节和模板
```

SKILL.md 从 348 行降至约 200 行，减少约 42%。重点是将 SKILL.md 中的规则解释文本和内联模板移至已有的辅助文件结构中。

### Innovation Highlights

三层瘦身法：对每个 skill 按需施以精简（冗余文本）、消歧（模糊指令）、拆分（如有必要）三种操作，而非一刀切。安全增量：每个任务独立 commit，可逐个验证回滚。

本方案与常规文件提取重构的关键区别在于：不是机械地按行数切分，而是以 **agent 指令执行语义** 为边界——SKILL.md 保留 agent 理解流程所需的"最少充分指令集"，辅助文件存放仅在特定步骤需要的"按需参考内容"。这种区分确保精简后的 SKILL.md 对 agent 而言是语义完整的，而非需要频繁跳转引用的碎片。

**跨领域借鉴**：
- **数据库范式化（Database Normalization）**：拆分层对应 2NF（消除部分依赖——将非流程内容移出流程主表），精简层对应消除冗余依赖（去除重复的规则描述），消歧层对应消除多值依赖（统一模糊术语的唯一语义）。
- **Strangler Fig Pattern（Martin Fowler）**：增量替换而非一次性重写，每个 task group 是一棵"绞杀者"的缠绕步骤，逐步替换旧结构。

**消歧层的方法论**：消歧操作针对 SKILL.md 中语义模糊的指令术语，采用"识别 → 定义 → 替换"三步法：
1. **识别**：扫描 SKILL.md 中出现但未在当前文件内定义的术语（如 `noTest`、`doc*`），标记为歧义项
2. **定义**：为每个歧义项撰写唯一、明确的定义，写入 skill 的 rules/ 消歧文档
3. **替换**：将 SKILL.md 中的模糊引用替换为指向消歧文档的精确引用

**消歧范围边界**：本方案仅处理以下已识别的歧义项，不做开放性扫描发现。若执行过程中发现其他歧义项，记录到 backlog 但不纳入本次范围，避免范围蔓延。

**已识别的歧义项**（执行前需确认）：
- `noTest`：在 guide.md 中指"跳过测试生成"，在部分 skill 中被解读为"该 skill 不涉及测试" → 统一为"跳过测试生成"语义
- `doc*` 前缀通配：在不同 skill 中指代范围不一致 → 统一为 `docs/proposals/` 下的 markdown 文档

## Requirements Analysis

### Key Scenarios

- **Agent 行为变化**：精简和消歧后 agent 加载 skill 时上下文窗口占用进一步减少，agent 的指令遵循准确率预期提升——因为 SKILL.md 仅包含直接相关的流程指令，agent 不再在冗余的规则解释中"稀释注意力"。具体表现为：(1) agent 在多步骤流程中跳步的概率降低（指令密度提高后关键步骤更突出）；(2) 术语歧义消除后 agent 不再因一词多义而选择错误的执行路径（如将 noTest 误解为"该 skill 不涉及测试"而跳过测试相关逻辑）。
- 开发者维护 skill 时通过 SKILL.md 快速理解流程，通过 rules/templates 了解细节

### Non-Functional Requirements

- 每个 SKILL.md 行数不超过 350 行（拆分后）
- 拆分产生的辅助文件放在 skill 目录内的 rules/ 或 templates/ 子目录
- 不改变 skill 的输入/输出契约

### Constraints & Dependencies

- 遵守 `docs/conventions/forge-distribution.md` 分发模型
- 遵守 `docs/conventions/skill-self-containment.md` 自洽原则——SKILL.md 必须包含完整流程步骤，辅助文件仅存放规则和模板细节
- 不涉及 Go 源码修改
- 不合并同类 skill（那是 skill-rationalization 的范畴）

## Alternatives & Industry Benchmarking

### Industry Solutions

大型 prompt 工程项目中，指令拆分是常见实践，以下为代表性案例：

1. **OpenAI GPT Best Practices**（openai.com/chatgpt-best-practices）：核心建议包括：(a) 将 system prompt 限制为"关键指令 + 行为定义"，避免与详细规则混在一起；(b) 使用"分段标记"（section markers）将不同类别的指令分离。本方案对应实践：SKILL.md 仅包含流程步骤（关键指令），规则和模板通过引用路径分离（分段标记）。
2. **Claude Tool Use Patterns**（docs.anthropic.com/en/docs/build-with-claude/tool-use）：Anthropic 的 tool-use 架构将工具定义（description + parameters）与主 prompt 分离，工具描述作为独立上下文块注入。具体模式：每个 tool 的 `description` 字段控制在 1-3 句话，详细参数说明放在 `parameters` 的 JSON Schema 中。本方案的 rules/ 子目录与此模式等价——SKILL.md 是主 prompt 的 description（流程概要），rules/ 是 parameters（详细规则）。
3. **Cursor Rules 架构**（github.com/getcursor/cursor）：将项目规则拆分为 `.cursor/rules/` 目录下的独立文件，通过 glob 模式（如 `*.ts`、`src/**`）按需加载，而非维护单一巨型规则文件。每个规则文件聚焦单一关注点（如 TypeScript 规范、测试约定）。本方案的 rules/ 结构与此一致：每个文件存放单一类型的规则或模板。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 5 个提案已证明现状不可持续 | Rejected: 债务积累 |
| LLM 自动压缩 prompt | DSPy (github.com/stanfordnlp/dspy) 的 prompt optimizer 通过 LLM 自动精简指令 | 无需人工判断，可批量处理 | 压缩后指令语义可能漂移；需额外 LLM 调用成本；不适合需要精确术语定义的场景 | Rejected: 不可控的语义风险 |
| 按需懒加载规则（仅当 agent 进入特定步骤时加载对应 rules 文件） | LangChain 的 dynamic prompt loading、Cursor Rules 的 glob 匹配按需加载 | 减少 token 消耗，仅加载必要上下文 | 需要修改 Forge 的 skill 加载机制（Go 源码），超出本方案范围 | Rejected: 需 Go 源码改动 |
| **按大小分层逐组处理** | 增量重构最佳实践（Martin Fowler: Strangler Fig Pattern） | 安全可控，立即可执行，无需工具链改动 | 小组内 skill 可能需不同策略 | **Selected: 平衡效率与安全** |

## Feasibility Assessment

### Technical Feasibility

纯文本修改 + 文件拆分。

**回归检测与回滚机制**：
- **回归定义**：拆分后 agent 对同一 skill 的执行输出与拆分前不一致（遗漏步骤、错误引用路径、产生格式偏差），视为回归。
- **检测方法**：每个 task 完成后，对涉及的 skill 运行一次 agent 测试。确定性 skill 使用固定输入 prompt 对比输出 diff，交互式 skill 验证核心功能可达。核心检查项：(1) SKILL.md 中引用的所有 rules/templates 路径存在且可读；(2) agent 仍按流程步骤顺序执行，无遗漏；(3) 输出格式与拆分前结构一致（确定性 skill）或核心功能可完成（交互式 skill）。
- **回滚触发**：若 agent 测试发现上述任一检查项失败，立即 `git revert` 该 task 的 commit，记录失败原因后重新执行该 task。

### Resource & Timeline

22 个 skill 分 9 组，预计 9 个任务。每个任务包含：分析 → 拆分/精简/消歧 → 验证引用完整性 → commit。

**时间估算**：
- Tier 1（独立任务，3 tasks）：每个 task 约 1 小时（含精简 + 消歧 + 验证），共 3 小时
- Tier 2（分组任务，3 tasks）：每个 task 约 45-60 分钟，共 2-3 小时
- Tier 3（分组任务，3 tasks）：每个 task 约 30-45 分钟，共 1.5-2 小时
- **总计**：约 6.5-8 小时工作量，可在 2 个工作日内完成

### Dependency Readiness

无外部依赖。所有文件已在本地。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "逐个 skill 清理效率低" | XY Detection | 用户核心需求是安全增量，不是效率最优；逐组处理可在安全与效率间平衡 |
| "SKILL.md 必须 100% 自包含" | Assumption Flip | 自洽不等于单文件；SKILL.md 包含完整流程 + 辅助文件包含细节规则，仍是自洽的。现有 eval skill 已用此模式 |
| "需要先审计再行动" | Occam's Razor | 5 个提案已做过充分分析；直接动手 + 逐个验证更简单有效 |

## Scope

### In Scope

- 22 个 `skills/*/SKILL.md` 文件的拆分、精简、消歧
- 在各 skill 目录内新建 rules/ 或 templates/ 子目录（按需）
- 清理过时标签、路径引用和歧义描述

### Out of Scope

- Go 源码修改
- skill 输入/输出契约变更
- 合并同类 skill（skill-rationalization 范畴）
- commands/ 和 agents/ 目录的文件
- hooks/、references/、scripts/ 目录

### Task Grouping (9 tasks)

**Tier 1: 大文件独立任务（3 tasks）**
1. consolidate-specs (348 行) → 精简 + 消歧
2. tech-design (190 行) → 精简 + 消歧
3. write-prd (231 行) → 精简 + 消歧

**Tier 2: 中文件按领域分组（3 tasks）**
4. eval (200) + gen-contracts (186) + test-guide (186) → 评测/质量域
5. gen-sitemap (229) + gen-journeys (211) + gen-test-cases (136) + gen-test-scripts (325) → 生成域
6. init-justfile (327) + ui-design (228) + extract-design-md (132) → 基础设施/设计域

**Tier 3: 小文件按领域分组（3 tasks）**
7. breakdown-tasks (144) + quick-tasks (173) + submit-task (156) → 任务管线域
8. brainstorm (106) + learn (183) + forensic (184) + improve-harness (163) → 元分析域
9. clean-code (190) + run-e2e-tests (193) → 工具域

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 拆分后 SKILL.md 丢失关键指令 | M | H | 每个 task 完成后执行 diff 审查清单：(1) 原文所有步骤编号在新 SKILL.md 中均有对应；(2) 所有条件分支和约束条件保留在 SKILL.md 或被正确引用；自动化检查：`grep -c` 对比拆分前后步骤关键字数量 |
| 辅助文件命名不统一 | L | L | 约定 rules/ 放规则、templates/ 放模板，不新建其他子目录类型；每个 task commit 前检查目录结构符合约定 |
| 消歧时引入新歧义 | L | M | 每处消歧需在 commit message 中注明原文和修改理由；消歧后的术语定义需在对应 rules/ 文档中可检索 |
| 拆分风格跨 task 不一致 | M | M | 在本提案中定义明确的拆分启发式规则（见 Splitting Heuristic 节），所有 task 遵循同一套规则而非参照标杆。最终做一次全局 review 确保一致性 |

## Success Criteria

- [ ] 每个 SKILL.md 行数不超过 350 行
- [ ] 22 个 SKILL.md 总行数减少 15%+（当前 4421 行 → 目标 3758 行以下）
- [ ] 无内部引用断裂（所有 SKILL.md 中引用的文件路径均存在）
- [ ] 每个 commit 仅涉及 1 组 skill，可独立回滚
- [ ] **功能正确性**：每个修改后的 skill 经 agent 测试验证，输出与拆分前功能等价。验证方法分类：
  - **确定性 skill**（如 consolidate-specs、gen-contracts、write-prd、tech-design）：使用固定输入 prompt 执行，对比拆分前后输出。等价标准：相同步骤按相同顺序执行，输出格式结构一致（允许措辞差异，不允许步骤遗漏或格式偏差），所有引用路径有效。
  - **交互式/非确定性 skill**（如 brainstorm、learn、forensic）：执行完整交互会话，验证拆分后的 skill 仍能完成其声明的核心功能（如 brainstorm 仍能产出结构化的想法列表，forensic 仍能定位问题根因）。等价标准：核心功能可完成，流程步骤无遗漏，不要求输出逐字一致。
- [ ] **消歧验证**：所有已识别的歧义项（`noTest`、`doc*`）在对应的 rules/ 文档中有明确定义，SKILL.md 中的原始模糊引用已替换为精确引用。记录 before/after 定义对照表
