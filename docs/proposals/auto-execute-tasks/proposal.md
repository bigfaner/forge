---
created: 2026-05-20
author: "faner"
status: Draft
---

# Proposal: Auto-Execute Tasks After Generation

## Problem

`/quick` 流水线在 brainstorm 用户确认提案后，仍会在 Step 2 停下来展示提案摘要并再次确认"是否生成任务"。这个确认门在 quick 模式下是多余的——用户已经在 brainstorm 步骤中确认并提交了提案。

### Evidence

- `quick.md` Step 2 强制使用 `AskUserQuestion` 展示摘要并等待确认，即使 brainstorm 已有用户审批环节
- 每次使用 `/quick` 都必须手动点一次确认，无法实现真正的自动化流水线

### Urgency

Quick 模式定位是"精简流水线"，多一次手动确认降低了效率。对于熟练用户，这个确认门是纯粹的摩擦。

## Proposed Solution

在 `.forge/config.yaml` 的 `auto` 块中添加 `runTasks` 配置项（ModeToggle 类型，`quick`/`full` 子键），控制 `/quick` 流水线是否跳过 Step 2 确认门直接进入任务生成+执行。

默认值：`quick: true, full: false`——quick 模式默认自动执行，full 模式保留确认门。

### Innovation Highlights

直接复用已有的 `ModeToggle` 模式（与 `e2eTest`、`consolidateSpecs` 同构），无需引入新的配置结构。

## Requirements Analysis

### Key Scenarios

- **自动执行（默认）**: 用户执行 `/quick`，brainstorm 确认提案后，自动生成任务并执行，全程无额外确认
- **手动确认**: 用户设置 `auto.runTasks.quick: false`，brainstorm 后暂停展示摘要，用户确认后才继续

### Constraints & Dependencies

- 依赖 Go CLI 的 `AutoConfig` 结构体支持新字段
- 依赖 `forge config get` 命令能读取新配置值

## Alternatives & Industry Benchmarking

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 无开发成本 | 保留不必要的摩擦 | Rejected: 违背 quick 模式精简定位 |
| 硬编码跳过确认门 | — | 最简单 | 无法回退，剥夺用户选择权 | Rejected: 不够灵活 |
| **配置控制（ModeToggle）** | 已有模式 | 灵活、向后兼容、复用现有结构 | 需改 Go CLI + quick.md | **Selected: 最优平衡** |

## Feasibility Assessment

### Technical Feasibility

改动范围小且明确：Go CLI 加一个字段 + 两个文件的逻辑分支。

### Resource & Timeline

单个开发者 30 分钟内可完成。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| Step 2 确认门是必要的 | Occam's Razor | Refuted: brainstorm 已包含用户审批，Step 2 是冗余确认 |
| 所有用户都想要自动执行 | Assumption Flip | Refined: 保留配置项让用户自选，默认 true 符合 quick 定位 |

## Scope

### In Scope

- Go CLI `AutoConfig` 添加 `RunTasks ModeToggle` 字段及默认值
- `quick.md` 命令文件添加配置检查逻辑，条件性跳过 Step 2
- forge guide 文档更新

### Out of Scope

- full 流水线的确认门行为变更（当前默认 `full: false` 即可）
- `/run-tasks` 独立运行时的行为变更

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 用户误操作，不想执行时自动开始 | L | M | brainstorm 本身已有审批环节，且配置可关闭 |
| AI agent 未正确读取配置值 | L | L | quick.md 中用 EXTREMELY-IMPORTANT 标注配置检查逻辑 |

## Success Criteria

- [ ] `auto.runTasks.quick: true`（默认）时，`/quick` 在 brainstorm 确认后直接生成+执行任务，无 Step 2 暂停
- [ ] `auto.runTasks.quick: false` 时，保留 Step 2 确认门（当前行为）
- [ ] `forge config get auto.runTasks` 正确返回配置值
- [ ] 向后兼容：未配置时使用默认值 `quick: true, full: false`

## Next Steps

- Proceed to `/quick-tasks` for task generation (quick mode)
