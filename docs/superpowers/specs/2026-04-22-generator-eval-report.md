# 生成器 Skill 评估报告

> 日期：2026-04-22
> 方法：先独立判断"好文档应该做到什么"，再交叉参考评分器 rubric，标注异同

---

## 1. brainstorm（提案生成器）

### 独立判断 — 好提案应该做到什么

- 证明问题是真实存在的（而非"我觉得"）
- 说明为什么现在要做（紧迫性）
- 诚实对比"不做"的代价
- 成功标准能让人不靠主观判断就能验证

### 当前弱点

| # | 问题 | 严重度 | 来源 |
|---|------|--------|------|
| B1 | `## Problem` 只有一行注释 "What problem? Why now?"，不引导提供证据和紧迫性 | 高 | 独立判断 + rubric 一致 |
| B2 | Alternatives 表格没有 "Do nothing" 行，容易变成自说自话 | 中 | 独立判断 + rubric 一致 |
| B3 | Success Criteria 提示 "measurable outcome" 但没强调必须包含量化指标，容易产出"体验更好"类虚词 | 中 | 独立判断 + rubric 一致 |
| B4 | Step 5 Write Proposal 无质量标准，写完直接呈现 | 低 | 独立判断 |

### Rubric 未覆盖但可能重要的点

无。Rubric 6 个维度（Problem Definition / Solution Clarity / Alternatives / Scope / Risk / Success Criteria）与独立判断高度一致。

---

## 2. write-prd（PRD 生成器）

### 独立判断 — 好 PRD 应该做到什么

- 任何人读完不需要追问就能开始设计/开发
- 流程图覆盖正常路径和异常分支
- 用户故事覆盖所有用户角色，每条有客观验收标准
- 功能描述完整到可测试

### 当前状态

**模板本身非常扎实** — 背景三要素、量化目标、Mermaid 流程图示例（含决策点和异常分支）、功能描述表格（列表页/按钮/表单/关联改动）、非功能性需求、质量检查清单都已具备。

### 当前弱点

| # | 问题 | 严重度 | 来源 |
|---|------|--------|------|
| P1 | Step 7 Write User Stories 没有显式要求"每个用户角色至少一个故事" | 中 | 独立判断 + rubric 一致 |
| P2 | 无自检步骤 — SKILL.md 从 Step 9 (Create Manifest) 直接到 Step 10 (Review & Commit)，模板底部有质量检查清单但 SKILL.md 不引用 | 中 | 独立判断 |
| P3 | Step 5 表格 Key Points 没有强调"校验规则必填" | 低 | rubric 维度（Functional Specs 6 分在 validation rules） |

### Rubric 维度是否有冗余或遗漏

5 个维度（Background & Goals / Flow Diagrams / Functional Specs / User Stories / Scope Clarity）均合理。没有明显冗余或遗漏。

---

## 3. design-tech（设计生成器）

### 独立判断 — 好技术设计应该做到什么

- 开发者能直接编码，不需要猜测类型或结构
- 错误处理具体到错误码和 HTTP 状态码
- 测试策略说清楚用什么工具、测什么
- **最重要：每个 PRD 需求都能在设计里找到对应实现** — 否则需求会静默丢失

### 当前状态

**三个生成器中差距最大的。** 模板有所有正确的章节名，但内部引导不够具体。

### 当前弱点

| # | 问题 | 严重度 | 来源 |
|---|------|--------|------|
| D1 | **无 PRD 需求覆盖映射** — 设计不引用 PRD 验收标准，需求可能静默丢失 | 高 | 独立判断（rubric Breakdown-Readiness 20 分也指向此） |
| D2 | **Interfaces/Data Models 是占位符** — `interface SomeInterface { method(arg: Type): Result }` 不引导写真实类型签名 | 高 | 独立判断（rubric "directly implementable" 也扣分于此） |
| D3 | **Error Handling 无结构** — "Define custom error types" + prose，无错误码表、HTTP 映射、传播策略 | 中 | 独立判断 + rubric 一致 |
| D4 | **Testing 缺工具和分层** — "Unit Tests" / "Integration Tests" 是散文提示，无 per-layer 表格和覆盖率数字 | 中 | 独立判断 + rubric 一致 |
| D5 | Step 1 只读 prd-spec.md，不读 prd-user-stories.md 的 AC 列表 | 中 | 独立判断 |

### Rubric 是否合理

Breakdown-Readiness 作为 20 分维度并设 12 分门槛合理 — 设计文档的终极目标就是驱动 task 拆解。Security (10分) 标注为条件评分（有 auth/privacy 需求时才评）务实。

---

## 改进优先级

### 高优先级（独立判断认为真有价值，rubric 也指向）

- **D1 + D5**: design-tech 加入 PRD AC 覆盖映射（Step 1 读 user-stories 提取 AC，模板加 Coverage Map 表格）
- **D2**: design-tech 模板 Interfaces/Data Models 引导写真实类型签名（改注释措辞，不增加新结构）
- **B1**: brainstorm Problem 拆分为陈述 + 证据 + 紧迫性

### 中优先级

- **D3 + D4**: design-tech Error Handling（错误码表 + HTTP 映射 + 传播策略）、Testing（per-layer 表格 + 工具 + 覆盖率）
- **P1 + P2**: write-prd 用户故事角色覆盖 + 自检步骤
- **B2**: brainstorm 加入 "Do nothing" 替代方案

### 低优先级 / 存疑

- **B4**: brainstorm 质量标准表 — 有价值但优先级低
- **P3**: write-prd 校验规则强调 — 模板已有，SKILL.md 可不重复

### 存疑点

design-tech 加入 PRD Coverage Map 对小功能是否过重？结论：不会。即使是 3 行的 AC 映射表，"防止需求丢失"的价值远大于填写成本。

---

## Rubric 与独立判断的一致性总结

三个 rubric 的维度设定均与独立判断一致，没有发现为了评分而人为设置的维度。Rubric 的主要价值在于将模糊的"好文档"拆解为可检查的子项（如"证据 vs 紧迫性"、"typed params vs prose-only"），这些拆分本身是合理的。
