---
created: 2026-05-15
author: "fanhuifeng"
status: Draft
---

# Proposal: Spec Drift Detection in consolidate-specs

## Problem

`docs/conventions/` 和 `docs/business-rules/` 中的规格文档在 feature 迭代后会逐渐与实际代码不一致（drift），而当前 consolidate-specs 只做单向推送（feature 文档 → 项目级 spec），没有验证现有 spec 是否仍然准确的机制。

### Evidence

- `docs/business-rules/` 有 3 个文件（5 条规则），`docs/conventions/` 有 2 个文件（2 条规范），全部来自 `forge-cli-v3` feature（2026-05-14 集成）
- 后续 feature（如 `typed-task-dispatch`）可能修改了相关代码但未触发 spec 审计
- 没有任何自动化流程验证 "代码改了，spec 是否还能反映真实行为"

### Urgency

随着 feature 数量增加，spec 文件漂移的累积风险线性上升。过时的 spec 会误导 agent 在任务执行时做出错误决策。现在 spec 文件少（5个），修复成本低。

## Proposed Solution

扩展 consolidate-specs skill，在现有的 "提取→审核→集成" 流程之后增加 **"漂移审计→自动修复"** 阶段。检测 feature 改动涉及的 domain 对应的 spec 文件，逐条验证规则是否与当前代码一致，漂移的规则自动更新。

### Innovation Highlights

- **代码感知审计**：不只看文档，而是对照实际代码验证每条规则
- **域定向检测**：通过 feature 改动范围确定需要审计的 spec domain，避免全量扫描
- **闭环维护**：从 "只推送" 变为 "推送 + 验证 + 修复"，形成文档-代码一致性的闭环

## Requirements Analysis

### Key Scenarios

- **正常漂移**：feature 改了错误处理逻辑，`docs/business-rules/error-reporting.md` 中的规则需要更新
- **规则失效**：某个规则对应的代码被完全移除，spec 自动删除该规则（commit message 记录 ID 和原因）
- **新增隐含规则**：代码中新增了符合 project-level 级别的约定，但 feature 文档未提及，应提取并追加
- **无漂移**：所有 spec 文件与代码一致，跳过修复

### Non-Functional Requirements

- 漂移检测应在 2 分钟内完成（spec 文件数量有限）
- 不影响现有 consolidate-specs 的提取和集成流程

### Constraints & Dependencies

- consolidate-specs 的 HARD-GATE "不要覆盖现有 spec 文件（仅追加）" 需调整为 "检测到漂移时允许修改"
- Quick 模式当前跳过 consolidate-specs，需要新增 T-quick-6（仅漂移检测，无提取步骤）

## Alternatives & Industry Benchmarking

### Industry Solutions

- **Schema migration tools**（如 Flyway/Liquibase）：版本化追踪变更，但适用于结构化 schema 而非自然语言 spec
- **Doc-as-code**（如 Docusaurus + versioned docs）：文档跟随代码版本，但不自动验证内容准确性
- **Contract testing**（如 Pact）：验证接口契约一致性，但范围限定在 API 层

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | Spec 漂移累积，误导 agent | Rejected: 漂移风险随 feature 数量线性上升 |
| 独立 spec-audit skill | 自研 | 关注点分离 | 多一个 skill 维护，T-test-5 需协调两个 skill | Rejected: 复杂度不值得，漂移审计与 consolidate-specs 天然耦合 |
| **扩展 consolidate-specs** | 自研 | 流程自然，一个 skill 覆盖全部，改动集中 | skill 复杂度增加 | **Selected: 最小改动实现完整闭环** |

## Feasibility Assessment

### Technical Feasibility

完全可行。consolidate-specs 已有读取 spec 文件和 feature 文档的能力，增加代码审计步骤是自然扩展。

### Resource & Timeline

scope 小（修改 3-4 个文件），适合 quick mode。

### Dependency Readiness

所有前置条件就绪：spec 文件格式已标准化（project-global ID），consolidate-specs 已有成熟的流程框架。

## Scope

### In Scope

- 扩展 consolidate-specs SKILL.md，添加漂移检测 + 自动修复子步骤（Step 9-11）
- 调整 HARD-GATE：从 "仅追加" 改为 "检测到漂移时允许修改"
- 在 quick-tasks SKILL.md 中添加 T-quick-6（仅漂移检测）
- 更新 breakdown-tasks SKILL.md 中 T-test-5 的描述
- 更新 guide.md 反映流程变化

### Out of Scope

- 新建独立 skill
- 修改 docs/decisions/ 或 docs/lessons/ 的管理逻辑
- 修改现有 spec 文件的内容
- all-completed hook 的全局审计（可作为后续增强）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 漂移检测误报（代码实际未变但检测为漂移） | M | L | 逐条对比规则关键词与代码实现，而非简单文本匹配 |
| 自动修复引入错误 | M | M | 修复时保持 spec 格式一致（project-global ID 不变），仅更新描述和行为描述 |
| 误删仍有价值的规则 | L | M | commit message 记录被删规则 ID 和原因，可通过 git 历史回溯恢复 |
| spec 文件过少（5个），漂移检测收益不明显 | L | L | 随 spec 文件增多收益自动提升，当前改动成本极低 |

## Success Criteria

- [ ] consolidate-specs 执行后，所有 `docs/conventions/` 和 `docs/business-rules/` 中的规则与当前代码一致
- [ ] 漂移的规则被自动更新，失效的规则被自动删除（commit message 中记录被删规则 ID 和原因，可通过 git 历史回溯）
- [ ] Quick 模式也能执行漂移检测（T-quick-6）
- [ ] 现有 consolidate-specs 的提取和集成流程不受影响

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
