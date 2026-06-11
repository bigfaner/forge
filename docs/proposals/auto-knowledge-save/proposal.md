---
created: 2026-05-20
author: faner
status: Draft
---

# Proposal: Auto Knowledge Save

## Problem

自动知识提取触发点（fix-bug、write-prd、tech-design）在提取知识后总是弹出 `AskUserQuestion` 让用户确认是否保存，打断工作流连贯性。

### Evidence

4 个触发点中有 3 个（fix-bug、write-prd、tech-design）在提取流程的 Step 5 使用 `AskUserQuestion` 呈现候选知识，要求用户逐条确认。在快速迭代场景下（如连续修 bug、快速输出 PRD），这种确认步骤频繁中断用户注意力。

### Urgency

每次触发点的知识确认都是一次上下文切换。用户已通过 `/quick` 模式表达了对效率的偏好，但知识提取的交互确认与这一偏好矛盾。

## Proposed Solution

在 `auto` 配置块中添加 `knowledgeSave` ModeToggle，控制 3 个触发点（fix-bug、write-prd、tech-design）的知识保存行为：

- **`true`**：自动提取并静默保存，不弹确认框。对标 `/consolidate-specs` 非交互模式。
- **`false`**：保留当前的交互确认流程。

默认值：`quick: true, full: false`。

同时**移除** `run-tasks` 的知识审查章节——它是任务调度器，真正的知识提取由 `doc.consolidate` 任务覆盖。

### Innovation Highlights

无特别创新。对标 `auto.consolidateSpecs` 的非交互模式设计，将"自动批准 `[CROSS]` 项"的模式推广到内联提取触发点。

## Requirements Analysis

### Key Scenarios

1. **Quick 模式连续修 bug**：`/fix-bug` 提取的 root cause 和调试模式自动写入 `docs/lessons/`，无需确认
2. **Full 模式写 PRD**：`/write-prd` 提取的业务规则保留交互确认，用户审查后决定是否保存
3. **Quick 模式输出 PRD + 设计**：`/write-prd` 和 `/tech-design` 提取的知识自动保存
4. **配置关闭时**：`auto.knowledgeSave: {quick: false, full: false}`，所有触发点恢复交互确认

### Non-Functional Requirements

- 向后兼容：现有项目未配置该字段时，Go 代码的 `WithDefaults()` 填充默认值，full 模式行为不变
- 一致性：静默保存的知识文件使用 `[auto-knowledge]` git commit 标签（对标 consolidate-specs 的 `[auto-specs]`），便于事后审查

### Constraints & Dependencies

- 依赖 `forge config get` CLI 命令读取配置
- 触发点在 skill/command markdown 中，由 AI agent 读取配置后决定行为
- Go 代码变更限于 `forge-cli/pkg/forgeconfig/config.go` 和 JSON Schema

## Alternatives & Industry Benchmarking

### Industry Solutions

自动化工具通常提供 "silent mode" 或 "auto-accept" 开关，允许跳过交互确认。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零改动 | 每次确认打断工作流 | Rejected: 核心问题未解决 |
| 完全移除确认 | — | 最简单 | 可能存入不想要的知识 | Rejected: 失去质量控制 |
| **ModeToggle 配置** | Forge `auto.*` 模式 | 灵活、向后兼容、与现有设计一致 | 需改动 Go + 3 个 skill 文件 | **Selected: 最小侵入性方案** |

## Feasibility Assessment

### Technical Feasibility

完全可行。Go 侧添加 ModeToggle 字段有成熟的模式（`consolidateSpecs`、`e2eTest`、`cleanCode`）。Skill 文件侧添加配置读取逻辑是条件分支。

### Resource & Timeline

小型变更：Go 代码 ~30 行，JSON Schema ~5 行，3 个 skill 文件各 ~10 行，run-tasks 删除 ~80 行。预计 1 个任务可完成。

### Dependency Readiness

`forge config get` 已支持读取 auto 配置块。无需新增基础设施。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "知识提取后必须让用户确认" | XY Detection | Confirmed: 确认是手段不是目的。目的是防止低质量知识入库，可通过事后 git diff 审查替代 |
| "run-tasks 需要知识审查" | 5 Whys | Overturned: run-tasks 是调度器，不产生新知识。consolidate-specs 任务已覆盖知识提取 |
| "配置应该控制是否提取" | Occam's Razor | Refined: 控制是否自动保存更准确。`false` 时仍提取但保留确认，`true` 时静默保存 |

## Scope

### In Scope

- Go 代码：`AutoConfig` 添加 `KnowledgeSave ModeToggle`，默认 `{Quick: true, Full: false}`
- Go 代码：更新 JSON Schema 添加 `knowledgeSave` 字段
- 移除 `run-tasks.md` 的知识审查章节（Knowledge Review + 后续的 Notable Knowledge Heuristics）
- 更新 `fix-bug.md`：读取 `auto.knowledgeSave` 配置，`true` 时跳过确认直接保存
- 更新 `write-prd/rules/knowledge-extraction.md`：同上
- 更新 `tech-design/rules/knowledge-extraction.md`：同上

### Out of Scope

- `/consolidate-specs` 行为变更（已有独立配置 `auto.consolidateSpecs`）
- `/learn` 行为变更（手动操作，始终可用）
- task-executor agent 变更
- 知识质量评分或过滤机制

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 静默保存低质量知识 | M | M | `[auto-knowledge]` git 标签便于事后审查和回退 |
| 用户不知道知识被自动保存 | L | L | run-tasks 结束摘要中列出已保存的知识条目 |
| 默认值改变现有行为 | L | H | full 模式默认 false，保持向后兼容；quick 模式默认 true 是新行为 |

## Success Criteria

- [ ] `auto.knowledgeSave` 配置项在 `.forge/config.yaml` 中生效，`forge config get auto.knowledgeSave` 正确返回值
- [ ] `quick: true` 时，3 个触发点的知识提取+保存流程无 `AskUserQuestion` 交互
- [ ] `full: false` 时，3 个触发点保留现有确认流程
- [ ] run-tasks.md 不再包含知识审查章节
- [ ] 未配置该字段的项目（零值），Go 代码 `WithDefaults()` 正确填充默认值

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
