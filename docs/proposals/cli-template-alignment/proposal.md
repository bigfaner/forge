---
created: "2026-05-30"
author: "faner"
status: Draft
intent: "refactor"
---

# Proposal: CLI Built-in Template Alignment with Forge v3

## Problem

`forge init` 生成的 CLAUDE.md 模板是纯通用 LLM 行为指南（60 行，4 个章节），完全没有 Forge 相关内容。用户在一个 Forge 管理的项目中，AI 助手从 CLAUDE.md 层面无法感知 Forge 的存在和用法，完全依赖 hook 注入的 forge-guide。

### Evidence

- 当前模板 `forge-cli/internal/embedded/claudemd_template.md` 无任何 Forge 关键字
- Forge 已演进至 v3.0.0（Surface types、Test profiles、pluggable test framework），但模板仍停留在初始通用版本
- .gitignore 条目和 init 流程也需审查是否与当前 Forge 运行时产物对齐

### Urgency

Forge 每次版本迭代都可能引入新的运行时产物或流程变化。模板长期未同步会导致新用户初始化后得到过时的项目配置。成本较低（单文件编辑 + 审查），应尽快完成。

## Proposed Solution

在现有 4 章行为准则之后，新增 "Forge Workflow" 章节，包含：
1. **项目结构**：Forge 创建的目录及其用途（精简版）
2. **CLI 命令速查**：核心 `forge` 命令（task claim/submit/status、feature set/complete）
3. **Surface & Test 类型**：5 种 Surface → Test 类型映射
4. **任务生命周期**：pending → in_progress → completed 状态流转
5. **详细参考指引**：指向 hook 注入的 forge-guide 获取完整文档

同时审查 .gitignore 条目和 init 流程是否需要调整。

### Innovation Highlights

这是常规的模板维护对齐，无技术创新。核心价值在于让 CLAUDE.md 作为项目级指令文件包含足够的 Forge 上下文，即使 hook 加载失败也能提供基本指引。

## Requirements Analysis

### Key Scenarios

- 新用户运行 `forge init`，得到的 CLAUDE.md 包含 Forge 工作流参考
- 已有项目的 CLAUDE.md 不受影响（init 跳过已存在的文件）
- AI 助手阅读 CLAUDE.md 即可理解项目使用 Forge 管理任务

### Non-Functional Requirements

- 模板大小控制在 120 行以内（现有 60 行 + Forge 章节 ~60 行）
- 不与 hook 注入的 forge-guide 大量重复

### Constraints & Dependencies

- 模板通过 Go embed 嵌入二进制，更新需要重新编译
- .gitignore 条目硬编码在 Go 代码中，修改需要重新编译

## Alternatives & Industry Benchmarking

### Industry Solutions

多数 AI 编码工具（Cursor、Windsurf）提供项目级配置模板，通常包含工具链特定的上下文。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 模板与 Forge 现状脱节 | Rejected: 对齐成本低，不做的风险高 |
| 全面重写 | — | 自包含 | 与 forge-guide 重复，维护负担 | Rejected: 重复内容违反 DRY |
| **精准增强** | — | 最小改动，聚焦关键信息 | 依赖 hook 提供完整参考 | **Selected: 最优平衡点** |

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
| .gitignore 条目可能已过时 | Codebase Analysis | Confirmed: 7 条条目与项目实际 .gitignore 一致，无需变更 |
| Init 流程可能缺少步骤 | Codebase Analysis | Confirmed: 6 步流程完整（.forge/ → CLAUDE.md → .gitignore → just → config → surfaces），无需变更 |
| CLAUDE.md 需要 Forge 完整参考 | Assumption Flip | Refined: 完整参考由 hook 注入，模板只需精简速查 + 指引 |

## Scope

### In Scope

- 更新 `claudemd_template.md`：新增 Forge Workflow 章节
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
| 模板过长导致 AI 上下文浪费 | L | L | 控制在 120 行以内，仅保留精简参考 |
| 模板内容与 forge-guide 漂移 | M | M | 模板只放稳定的骨架信息，详细内容指向 forge-guide |
| 遗漏新增的 Forge 运行时产物 | L | L | 基于 .gitignore 和代码审查交叉验证 |

## Success Criteria

- [ ] `claudemd_template.md` 包含 "Forge Workflow" 章节，涵盖项目结构、CLI 速查、Surface/Test 概念、任务生命周期
- [ ] 模板总行数 ≤ 120 行
- [ ] 现有 4 章行为准则完整保留、无修改
- [ ] .gitignore 条目经审查确认为最新（或已补充遗漏项）
- [ ] Init 流程经审查确认无问题
- [ ] `go build ./...` 和 `go test ./...` 在 forge-cli 目录通过
