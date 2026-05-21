---
created: 2026-05-21
author: faner
status: Draft
---

# Proposal: Testing Domain Separation

## Problem

编码任务执行时，agent 加载了所有测试相关的约定文件（包括 e2e 管道知识），浪费 context window 并可能误导单元测试编写。

### Evidence

当前 `docs/conventions/` 下 6 个测试约定文件全部以 `testing` 作为 domains 标签。编码模板（coding-*.md）的约定加载指令是：`Load files whose domains overlap with the task context`。当任务涉及测试时，所有 6 个文件都会被加载，包括：

- `testing-journey-contract.md` — 纯 e2e 管道知识（Journey-Contract 模型、六维声明）
- `testing-conventions.md` — 测试管道元参考（文件结构、验证规则、合并语义）
- `testing-go/vitest/ginkgo.md` — 框架知识（但内容全是 e2e 导向：build tag、TC 命名、e2e helpers）

这些内容对编码任务的单元测试编写没有帮助，反而可能误导 agent 添加 `//go:build e2e` 标签或使用 TC 命名风格写单元测试。

### Urgency

e2e 测试知识持续增长（Journey-Contract 模型、isolation 规则等），每次编码任务都加载这些知识是不必要的 context 消耗。越早修正，浪费越少。

## Proposed Solution

将所有测试约定文件的 domains 标签从 `testing` 改为 `e2e`，明确表达这些文件服务于 e2e 测试管道而非编码任务。同步更新所有引用 `testing` 域的技能 文件。

编码模板不需要修改 — 现有的 "load files whose domains overlap" 指令会自然跳过不再含 `testing` 的文件。

### Innovation Highlights

无特别创新。这是一个语义精确化改动：将模糊的 `testing` 标签替换为精确的 `e2e`，让约定文件的用途与标签一致。

## Requirements Analysis

### Key Scenarios

- **编码任务执行**：agent 写单元测试时，不加载任何 e2e 约定文件（零噪音）
- **e2e 管道执行**：gen-test-scripts / run-e2e-tests / init-justfile 正常加载框架约定（行为不变）
- **新建约定文件**：test-guide 生成的文件自动使用 `e2e` 域

### Non-Functional Requirements

- 向后兼容：已有的约定文件通过一次性迁移更新
- 一致性：所有引用 `testing` 域的地方统一改为 `e2e`

### Constraints & Dependencies

- gen-test-scripts Step 0.1.3 硬编码了 `domains contain testing` 过滤
- run-e2e-tests Step 0 硬编码了 `domains containing testing` 过滤
- init-justfile Step 硬编码了 `domains frontmatter containing testing` 过滤
- test-guide 模板生成了 `domains: [testing, {{SCOPE}}]`

## Alternatives & Industry Benchmarking

### Industry Solutions

约定文件的域标记（domain tagging）是常见的知识分类方法。关键在于标签的语义精确性。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 无改动成本 | 持续浪费 context、可能误导 agent | Rejected: 问题真实存在且会恶化 |
| test-scope 字段 | 前述讨论 | 声明式分类、可扩展 | 新增字段、所有文件和模板都要更新 | Rejected: 更重的方案，当前只需区分 e2e/非 e2e |
| 精细化 domains（`testing` → `e2e`） | 当前方案 | 最小改动、语义精确、编码模板无需修改 | 一次性迁移所有文件 | **Selected: 最直接、最少的副作用** |

## Feasibility Assessment

### Technical Feasibility

纯文本修改（frontmatter + SKILL.md 文档），无代码变更。技术上零风险。

### Resource & Timeline

约 10 个文件修改，1 人时内可完成。

### Dependency Readiness

无外部依赖。所有修改都在 forge plugin 和 docs 目录内。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 框架文件（testing-go 等）包含单元测试知识 | 代码审查 | Overturned: 所有框架文件内容 100% e2e 导向（build tag、TC 命名、e2e helpers） |
| 编码任务需要加载测试约定 | 5 Whys | Refined: 编码任务需要写单元测试，但当前约定文件没有单元测试知识，agent 依赖自身知识即可 |
| 移除 `testing` 域会破坏 gen-test-scripts | 依赖分析 | Confirmed: gen-test-scripts 等 3 个 skill 硬编码了 `testing` 域过滤，需同步修改 |

## Scope

### In Scope

- 更新 5 个约定文件的 frontmatter domains（`testing` → `e2e`）
- 更新 3 个 skill SKILL.md 的域过滤描述（`testing` → `e2e`）
- 更新 test-guide 的约定模板和结构规则（`testing` → `e2e`）
- 更新 testing-conventions.md 元参考文档中的域标签说明和示例

### Out of Scope

- 编码模板（coding-*.md）修改 — 不需要，现有加载指令自然过滤
- 新增单元测试约定文件 — 当前 agent 依赖自身框架知识 + 项目代码侦察足够
- testing-isolation.md 修改 — 已无 `testing` 域，不涉及

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 遗漏某处 `testing` 域引用 | L | M | 全局 grep `testing` 在 skill 文件中，确保不遗漏 |
| 未来新增单元测试约定文件时使用旧 `testing` 标签 | M | L | testing-conventions.md 文档明确说明 `e2e` 域的语义 |
| 第三方用户项目已有自定义约定文件用 `testing` | M | M | 在 testing-conventions.md 中添加迁移说明 |

## Success Criteria

- [ ] 编码任务执行时不再加载任何 `docs/conventions/testing-*.md` 文件
- [ ] gen-test-scripts / run-e2e-tests / init-justfile 仍能正确发现并加载框架约定文件
- [ ] test-guide 生成的约定文件 frontmatter 使用 `e2e` 域
- [ ] testing-conventions.md 文档反映更新后的域标签语义

## Next Steps

- Proceed to `/quick-tasks` to generate tasks
