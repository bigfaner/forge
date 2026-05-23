---
created: 2026-05-23
type: validation
---

# eval-freeform-expert-review: 验证与设计复盘

## 1. Subagent 模拟验证

### 测试场景

对一份假设的「AI 驱动自动化测试生成系统」提案执行 Phase 0 模拟，验证动态专家推断和自由评审的实际效果。

### 推断结果

**生成专家**：测试基础设施与 AI 工程交叉领域首席架构师

| 维度 | 内容 |
|------|------|
| Domain | automated testing, LLM-based code generation, static analysis, developer tooling, ML system reliability |
| Background | 日均 50 万次测试的企业级测试平台经验、Babel/SWC 级别的 AST 分析工具开发、多个 LLM 代码生成管线投产经验、IntelliTest/EvoSuite/Randoop 理论基础、亲历类似项目因假阳性率过高被弃用 |
| Review Style | 以「生产环境连续运行 6 个月后会怎样」为核心视角，关注边界条件、失败模式、长期可维护性，特别追问 LLM 输出不可控性的驯服方案 |

### CTO rubric 无法覆盖的领域特有发现

| # | 发现 | 类别 | CTO rubric 覆盖？ |
|---|------|------|-------------------|
| 1 | 通过率 80% 是假指标，需用变异测试杀死率替代 | 度量方法 | 否——rubric 只问「可测吗」，不追问度量是否正确 |
| 2 | 反馈循环有 5 种不可区分的失败归因，prompt 优化是在噪声上做梯度下降 | 架构可行性 | 否——rubric 会说「缺少错误处理」，但无法分析 LLM 归因的不可解性 |
| 3 | AST「数据流分析」轻描淡写，实际需完整类型系统参与 | 工程深度 | 否——rubric 无法判断 AST 分析的实现难度 |
| 4 | LangChain 的抽象层阻碍调试，应直接用 OpenAI SDK | 技术选型 | 否——rubric 无框架选型评判维度 |
| 5 | 自动化悖论：开发者停止写测试，覆盖率反而下降 | 行为适应 | 否——rubric 风险评估无「人的行为适应」维度 |
| 6 | 生成的测试六个月后变成技术债务——命名约定和目录隔离是必需品 | 长期维护 | 否——纯工程实践问题，rubric 不覆盖 |
| 7 | PRD 与代码不一致时 oracle 问题无解——按 PRD 生成会失败，按代码生成只验证现状 | 测试理论 | 否——超出 rubric 所有 10 个维度的范围 |

### 结论

7 个发现中有 6 个是 CTO rubric 10 个维度完全无法覆盖的领域特有风险。验证了动态专家评审的核心假设：**领域专家能发现 rubric 视而不见的问题**。

## 2. 设计复盘：逐级审批机制

### 原始设计

```
eval-proposal（默认 CTO rubric）  [--freeform-expert 可选] → Phase 0 自由评审 → rubric
```

自由专家是 opt-in 的可选功能，CTO rubric 是默认路径。

### 用户的挑战

> 审核提案的专家：为什么是 CTO？

CTO 作为固定专家是一个惯性选择——「提案通常给高管审批」。但 CTO 关注的是通用管理问题（隐藏成本、回滚计划），对任何提案都问一样的话，导致评审深度肤浅。

### 设计演进

```
逐级审批：领域专家（动态，自由叙事）→ CTO rubric（固定，结构化评分）
```

- 领域专家**默认启用**，不再是 opt-in
- CTO rubric 紧随其后，接收领域专家的发现注入
- 移除 `--freeform-expert` 参数

### 两层各自的价值

| 层级 | 职责 | 提供 | 不提供 |
|------|------|------|--------|
| 领域专家 | 技术深度 | 领域特有的风险发现、架构挑战、技术选型批评 | 结构化评分、系统性覆盖、迭代修订 |
| CTO rubric | 质量兜底 | 成功标准可测性、替代方案公平性、范围边界、评分闸门、reviser 迭代循环 | 领域特有的技术洞察 |

**一句话总结**：领域专家提供洞察的**上限**，CTO rubric 提供质量的**下限**。

### 已落地的代码变更

- `plugins/forge/skills/eval/SKILL.md` — Phase 0 改为 proposal 类型默认行为，移除 `--freeform-expert` 参数，流程图条件从 `--freeform-expert && type == proposal` 改为 `type == proposal`
- `plugins/forge/commands/eval-proposal.md` — 描述更新为「Two-tier sequential evaluation」，移除参数
