---
created: "2026-05-30"
author: "faner"
status: Draft
intent: "refactor"
---

# Proposal: CLI Built-in Template Alignment with Forge v3

## Problem

`forge init` 生成的 CLAUDE.md 模板是纯通用 LLM 行为指南（60 行，4 个章节），存在两个问题：

1. **内容冗余**：模板中的 4 条行为准则（Think Before Coding、Simplicity First、Surgical Changes、Goal-Driven Execution）已逐字融入所有 coding prompt 模板（`forge-cli/pkg/prompt/templates/coding-*.md`），在 CLAUDE.md 中重复存在。
2. **缺失 Forge 上下文**：模板无任何 Forge 相关内容。用户在一个 Forge 管理的项目中，AI 助手从 CLAUDE.md 层面无法感知 Forge 的存在和用法。

### Evidence

- `coding-feature.md`、`coding-enhancement.md`、`coding-fix.md`、`coding-cleanup.md`、`coding-refactor.md` 均包含与 `claudemd_template.md` 完全相同的行为准则表述
- 当前模板无任何 Forge 关键字
- Forge 已演进至 v3.0.0（Surface types、Test profiles、pluggable test framework），但模板仍停留在初始通用版本

### Urgency

行为准则的重复导致每次对话都在 CLAUDE.md 和 prompt 模板中双重注入相同内容，浪费上下文窗口。应尽快替换为有价值的 Forge 上下文。

## Proposed Solution

**替换** CLAUDE.md 模板内容：移除冗余的通用行为准则，用纯 Forge 项目上下文替代。新模板包含：

1. **项目结构**：Forge 创建的目录及其用途（精简版）
2. **CLI 命令速查**：核心 `forge` 命令（task claim/submit/status、feature set/complete）
3. **Surface & Test 类型**：5 种 Surface → Test 类型映射
4. **任务生命周期**：pending → in_progress → completed 状态流转
5. **详细参考指引**：指向 hook 注入的 forge-guide 获取完整文档

同时审查 .gitignore 条目和 init 流程是否需要调整。

### Innovation Highlights

这是常规的模板维护对齐，无技术创新。核心洞察在于识别出行为准则已在 prompt 模板层覆盖，CLAUDE.md 应专注于 prompt 模板不覆盖的领域：Forge 项目级上下文。

## Requirements Analysis

### Key Scenarios

- 新用户运行 `forge init`，得到的 CLAUDE.md 包含 Forge 工作流参考
- 已有项目的 CLAUDE.md 不受影响（init 跳过已存在的文件）
- AI 助手阅读 CLAUDE.md 即可理解项目使用 Forge 管理任务
- 行为准则仍通过 prompt 模板在任务执行时注入，不丢失任何能力

### Non-Functional Requirements

- 模板大小控制在 80 行以内（纯 Forge 上下文，无冗余行为准则）
- 不与 hook 注入的 forge-guide 大量重复

### Constraints & Dependencies

- 模板通过 Go embed 嵌入二进制，更新需要重新编译
- .gitignore 条目硬编码在 Go 代码中，修改需要重新编译

## Alternatives & Industry Benchmarking

### Industry Solutions

多数 AI 编码工具（Cursor、Windsurf）提供项目级配置模板，通常包含工具链特定的上下文而非通用行为准则。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 冗余 + 缺失 Forge 上下文 | Rejected: 双重问题 |
| 保留 + 追加 | — | 最小改动 | 行为准则双重注入，浪费上下文 | Rejected: 冗余内容应移除 |
| **替换为 Forge 上下文** | — | 消除冗余，聚焦有价值内容 | 无行为准则（已由 prompt 覆盖） | **Selected: 最优方案** |

## Feasibility Assessment

### Technical Feasibility

完全可行。单文件编辑 + 审查，无技术风险。

### Resource & Timeline

单次编辑，预计 1 个任务即可完成。

### Dependency Readiness

无外部依赖。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 通用行为准则应保留在 CLAUDE.md | Codebase Analysis | **Overturned**: 准则已融入 `coding-*.md` prompt 模板，CLAUDE.md 中的版本完全冗余 |
| .gitignore 条目可能已过时 | Codebase Analysis | Confirmed: 7 条条目与项目实际 .gitignore 一致，无需变更 |
| Init 流程可能缺少步骤 | Codebase Analysis | Confirmed: 6 步流程完整，无需变更 |
| CLAUDE.md 需要 Forge 完整参考 | Assumption Flip | Refined: 完整参考由 hook 注入，模板只需精简速查 + 指引 |

## Scope

### In Scope

- 替换 `claudemd_template.md` 内容：移除冗余行为准则，替换为 Forge 项目上下文
- 审查 .gitignore 条目完整性并记录结论
- 审查 init 流程正确性并记录结论
- 确保变更后 Go 编译和测试通过

### Out of Scope

- Justfile 模板更新（Playwright 硬编码 → test profile 对齐，属于独立 feature）
- Hook 注入的 forge-guide 变更
- Init 交互流程（TUI 提示、配置生成）调整
- Prompt/task 模板或其他嵌入资源更新

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 模板内容与 forge-guide 漂移 | M | M | 模板只放稳定的骨架信息，详细内容指向 forge-guide |
| 遗漏新增的 Forge 运行时产物 | L | L | 基于 .gitignore 和代码审查交叉验证 |

## Success Criteria

- [ ] `claudemd_template.md` 包含 Forge 项目上下文（项目结构、CLI 速查、Surface/Test 概念、任务生命周期）
- [ ] 模板中不再包含已融入 prompt 模板的冗余行为准则
- [ ] 模板总行数 ≤ 80 行
- [ ] .gitignore 条目经审查确认为最新（或已补充遗漏项）
- [ ] Init 流程经审查确认无问题
- [ ] `go build ./...` 和 `go test ./...` 在 forge-cli 目录通过
