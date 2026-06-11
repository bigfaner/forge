---
created: "2026-05-30"
author: "faner"
status: Draft
intent: "refactor"
---

# Proposal: Remove CLAUDE.md Generation from forge init

## Problem

`forge init` 生成的 CLAUDE.md 模板是纯通用 LLM 行为指南（60 行，4 个章节），存在三层冗余：

1. **与 prompt 模板重复**：4 条行为准则（Think Before Coding、Simplicity First、Surgical Changes、Goal-Driven Execution）已逐字融入所有 coding prompt 模板（`forge-cli/pkg/prompt/templates/coding-*.md`）。
2. **与 hook 注入的 guide.md 重复**：Forge 工作流上下文已通过 plugin hook 以 `<forge-guide>` 形式注入到每次对话的 system prompt 中。
3. **无独特价值**：模板既不含 Forge 上下文（由 guide.md 覆盖），行为准则也不独特（由 prompt 模板覆盖）。

### Evidence

- `coding-feature.md`、`coding-enhancement.md`、`coding-fix.md`、`coding-cleanup.md`、`coding-refactor.md` 均包含与 `claudemd_template.md` 完全相同的行为准则表述
- `plugins/forge/hooks/` 中的 hook 在每次对话时注入完整 Forge 参考到 system prompt
- 当前模板无任何 Forge 关键字，对 Forge 项目无附加价值

### Urgency

三层冗余意味着每次对话都在上下文窗口中注入 3 份相同或重叠的内容。移除 `forge init` 的 CLAUDE.md 生成步骤是最干净的解决方案。

## Proposed Solution

**从 `forge init` 中移除 CLAUDE.md 生成步骤**（Step 2）。具体变更：

1. 删除 `forge-cli/internal/embedded/claudemd_template.md` 嵌入模板文件
2. 删除 `forge-cli/internal/embedded/claudemd.go` 中的 embed 声明和导出变量
3. 从 `forge-cli/internal/cmd/init.go` 中移除 `createCLAUDEmd()` 函数及其调用
4. 更新 `initCmd` 的 Long 描述，移除 "generates CLAUDE.md from embedded template" 语句
5. 审查 .gitignore 条目和 init 流程其余步骤是否需要调整

### Innovation Highlights

这是基于冗余分析的简化操作。核心洞察：当内容已被两个独立通道（prompt 模板 + hook 注入）完整覆盖时，第三份重复文件应被移除而非修补。

## Requirements Analysis

### Key Scenarios

- 新用户运行 `forge init`，不再自动生成 CLAUDE.md；Forge 上下文由 hook 注入提供
- 已有项目的 CLAUDE.md 不受影响（不会被删除或修改）
- 行为准则仍通过 prompt 模板在任务执行时注入，不丢失任何能力
- 用户可自行创建 CLAUDE.md 添加项目特定指令（Forge 不干预）

### Non-Functional Requirements

- init 流程更简洁（少一个步骤）
- 上下文窗口不再浪费在冗余内容上

### Constraints & Dependencies

- 二进制体积微减（移除嵌入的 ~60 行文本）
- 非破坏性变更：已有 CLAUDE.md 不受影响

## Alternatives & Industry Benchmarking

### Industry Solutions

多数 AI 编码工具不在 init 时生成项目级指令文件，留由用户自行决定。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 三层冗余持续浪费上下文 | Rejected: 问题已确认 |
| 替换为 Forge 上下文 | — | 消除行为准则冗余 | 与 guide.md 仍然重复 | Rejected: 第二层冗余未解决 |
| 极简占位符 | — | 最小改动 | 仍有文件存在，用户可能困惑 | Rejected: 不彻底 |
| **移除生成步骤** | — | 彻底消除三层冗余 | 用户需自行创建 CLAUDE.md | **Selected: 最干净** |

## Feasibility Assessment

### Technical Feasibility

完全可行。删除代码比修改代码更简单，无技术风险。

### Resource & Timeline

单次变更，涉及 3 个文件的删除/修改。

### Dependency Readiness

无外部依赖。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| CLAUDE.md 模板有独特价值 | XY Problem Detection | **Overturned**: 行为准则由 prompt 模板覆盖，Forge 上下文由 guide.md 覆盖，模板无独特价值 |
| 应保留模板并优化内容 | Assumption Flip | **Overturned**: 优化重复内容不如移除重复源 |
| .gitignore 条目可能已过时 | Codebase Analysis | Confirmed: 7 条条目与项目实际 .gitignore 一致，无需变更 |
| Init 流程可能缺少步骤 | Codebase Analysis | Confirmed: 其余 5 步流程完整，移除 CLAUDE.md 后仍有 .forge/、.gitignore、just、config、surfaces |

## Scope

### In Scope

- 移除 `forge init` 中的 CLAUDE.md 生成步骤及相关嵌入资源
- 更新 init 命令描述
- 审查 .gitignore 条目完整性并记录结论
- 确保变更后 Go 编译和测试通过

### Out of Scope

- Justfile 模板更新（Playwright 硬编码 → test profile 对齐）
- Hook 注入的 forge-guide 变更
- Init 交互流程其余部分调整
- Prompt/task 模板或其他嵌入资源更新

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 用户期望 init 生成 CLAUDE.md | M | L | init summary 报告明确显示未生成；用户可自行创建 |
| 遗漏引用 claudemd_template 的其他代码 | L | M | 全局搜索 `CLAUDEmdTemplate` 和 `claudemd` 确认无其他引用 |

## Success Criteria

- [ ] `forge init` 不再生成 CLAUDE.md
- [ ] `claudemd_template.md` 和 `claudemd.go` 已删除
- [ ] `init.go` 中 `createCLAUDEmd()` 函数已移除
- [ ] 无残留的 `CLAUDEmdTemplate` 引用
- [ ] .gitignore 条目经审查确认为最新
- [ ] `go build ./...` 和 `go test ./...` 在 forge-cli 目录通过
