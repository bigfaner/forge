---
created: 2026-05-26
author: "faner"
status: Draft
---

# Proposal: Auto-Eval Configuration for Document Evaluation Stages

## Problem

4 个文档评估 skill（eval-proposal、eval-prd、eval-ui、eval-design）在对应文档生成后，通过 `AskUserQuestion` 手动询问用户是否运行评估。每次交互都需要用户手动选择，增加了流水线的交互成本。项目已有 `auto.runTasks` 等 ModeToggle 配置驱动自动化的成功模式，eval 阶段应采用相同机制。

### Evidence

- brainstorm、write-prd、tech-design 三个 skill 均在文档提交后通过 `AskUserQuestion` 询问是否运行 eval
- ui-design 是唯一例外——无条件自动运行 eval-ui，与其他 skill 行为不一致
- `auto.runTasks` 已证明 ModeToggle 配置模式可以成功消除冗余交互

### Urgency

随着 quick 流水线不断优化自动化流程（auto.runTasks、auto.consolidateSpecs 等），eval 的手动确认成为剩余的主要交互摩擦点。统一配置模式可降低流水线使用成本。

## Proposed Solution

在 `.forge/config.yaml` 的 `auto` 块中新增 `eval` 嵌套结构体，包含 4 个独立的 ModeToggle 字段：

- `auto.eval.proposal` — 控制 eval-proposal 是否自动运行
- `auto.eval.prd` — 控制 eval-prd 是否自动运行
- `auto.eval.uiDesign` — 控制 eval-ui 是否自动运行
- `auto.eval.techDesign` — 控制 eval-design 是否自动运行

每个 ModeToggle 支持 `quick`/`full` 子键，区分 quick 和 full 流水线的行为。

默认值：
- `proposal`: `quick: true, full: true` — proposal eval 默认自动运行
- `prd`: `quick: false, full: false` — 默认询问用户
- `uiDesign`: `quick: false, full: false` — 默认询问用户（改变现有无条件自动行为）
- `techDesign`: `quick: false, full: false` — 默认询问用户

### Innovation Highlights

复用已有 ModeToggle 模式，无新概念。与 `auto.runTasks` 等现有配置完全同构，降低学习成本。

## Requirements Analysis

### Key Scenarios

- **自动评估 proposal（默认）**: 用户执行 `/quick`，brainstorm 提交 proposal 后自动运行 eval-proposal，无需手动确认
- **手动确认 prd eval**: 用户执行 `/write-prd`，提交 PRD 后询问是否运行 eval-prd（默认行为）
- **配置驱动的灵活控制**: 用户通过 `forge config set auto.eval.prd.full true` 开启 PRD 自动评估
- **ui-design 行为统一**: ui-design 从无条件自动运行改为读取配置，与其他 skill 行为一致

### Non-Functional Requirements

- 向后兼容：未配置时使用默认值，不影响现有行为（除 ui-design）
- 配置热生效：修改配置后无需重启即可生效

### Constraints & Dependencies

- 依赖 Go CLI 的 `AutoConfig` 结构体支持嵌套字段
- 依赖 `forge config get/set` 命令支持嵌套路径（`auto.eval.proposal.quick`）
- 4 个 skill 文件需增加 config check 逻辑

## Alternatives & Industry Benchmarking

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 无开发成本 | 保留交互摩擦，ui-design 行为不一致 | Rejected: 与自动化方向矛盾 |
| 单一 auto.eval 开关 | — | 配置简单 | 无法针对不同阶段单独控制 | Rejected: 粒度不足 |
| **4 个独立 ModeToggle（嵌套）** | auto.runTasks 模式 | 灵活、复用现有模式、支持 quick/full 分控 | 配置项较多 | **Selected: 最优平衡** |
| 扁平结构（evalProposal 等） | — | 实现简单 | 不符合 `auto.eval.xxx` 命名空间预期 | Rejected: 命名不直观 |

## Feasibility Assessment

### Technical Feasibility

完全可行。Go CLI 已有 ModeToggle 基础设施，新增嵌套结构体和 4 个 skill 的 config check 是确定性改动。

### Resource & Timeline

预计 1-2 小时完成所有改动（Go CLI + 4 个 skill + 测试）。

### Dependency Readiness

所有依赖（Go CLI config 系统、ModeToggle、forge config 命令）均已就绪。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 每次评估都需要用户确认 | Occam's Razor | Refuted: 熟练用户知道何时需要评估，手动确认是摩擦 |
| ui-design 的无条件自动评估是正确的 | Assumption Flip | Refined: 应与其他 skill 一致，由配置控制 |
| proposal 默认应询问用户 | 5 Whys | Refuted: proposal 是流水线入口，自动评估可尽早发现问题 |

## Scope

### In Scope

- Go CLI `AutoConfig` 新增 `Eval` 嵌套结构体（proposal/prd/uiDesign/techDesign 各为 ModeToggle）
- `AutoConfigDefaults()` 设置默认值（proposal: ON，其余: OFF）
- `forge config get/set auto.eval.*` 命令支持
- 4 个 skill（brainstorm、write-prd、tech-design、ui-design）增加 config check 逻辑
- JSON schema 和 example config 更新
- 单元测试更新

### Out of Scope

- eval skill 本身（eval/SKILL.md）的行为修改
- 新的 eval 类型添加
- quick/full 管道流程的其他变更
- forge guide 文档更新（文档变更随代码一起完成）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 嵌套结构改变 auto 字段解析路径，影响现有 flat 字段 | L | H | 分离 eval 子路径解析逻辑，不触碰现有 flat 字段解析 |
| ui-design 从无条件自动变为配置驱动，现有用户习惯被打破 | M | M | 默认 full:false 保持询问；在 release notes 中说明 |
| 4 个 skill 的 config check 逻辑不一致 | L | M | 使用统一的配置检查模板，在 skill 中用 EXTREMELY-IMPORTANT 标注 |

## Success Criteria

- [ ] `forge config set auto.eval.proposal.quick true` 正确写入 config
- [ ] `forge config get auto.eval.proposal` 返回 `quick:true full:true`
- [ ] brainstorm 在 `auto.eval.proposal` 对应模式为 true 时跳过 AskUserQuestion，直接运行 eval-proposal
- [ ] brainstorm 在 `auto.eval.proposal` 对应模式为 false 时保持现有 AskUserQuestion 行为
- [ ] write-prd/tech-design/ui-design 同理遵循配置驱动
- [ ] 未配置时（config missing）：proposal 默认自动运行，其余默认询问用户
- [ ] 现有配置测试（config_test.go、config_schema_test.go）通过

## Next Steps

- Proceed to `/quick-tasks` for task generation (quick mode)
