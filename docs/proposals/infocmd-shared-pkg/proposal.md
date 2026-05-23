---
created: 2026-05-23
author: faner
status: Draft
---

# Proposal: 提取 Info-Command 共享工具包

## Problem

`research`、`proposal`、`lesson` 三个 CLI 命令的代码高度同构，但各自维护独立副本。`parseFrontmatter` 复制 3 份、`Discover`/`FindBySlug` 模式重复 3 次、排序逻辑复制 3 份、表格渲染模式雷同。每次 bug fix 或改进需改三处，新增类似命令需再次复制粘贴。

### Evidence

- `pkg/research/research.go`、`pkg/proposal/proposal.go`、`pkg/lesson/lesson.go` 各含一份 `parseFrontmatter`（逐字节相同）
- 三个 `pkg/` 包均实现 `Discover()` + `FindBySlug()`/`FindByName()`，结构一致仅类型不同
- `internal/cmd/research.go`、`proposal.go`、`lesson.go` 的 `runXList`/`runXDetail` 渲染逻辑几乎相同
- 即将新增更多 info-command，复制粘贴问题将进一步扩大

### Urgency

即将新增更多同类命令。如果先新增再抽象，改造范围更大、回归风险更高。现在提取成本最低。

## Proposed Solution

创建 `pkg/infocmd/` 共享工具包，用 Go 泛型封装重复模式。每个 info-command 只需定义数据模型 + 列配置 + 几行胶水代码，框架负责 frontmatter 解析、目录扫描、排序、表格渲染。

### Innovation Highlights

无创新，纯工程实践。Go 1.18+ 泛型使"一份 Discover[T]"成为可能，消除三份相同代码是标准的 DRY 重构。

## Requirements Analysis

### Key Scenarios

- **列表视图**：`forge <cmd>` → Discover → 排序 → 渲染表格，列头和列内容由调用方定义
- **详情视图**：`forge <cmd> <slug>` → FindBySlug → 渲染键值对，字段由调用方定义
- **跨引用扩展**：`proposal` 详情需检查 PRD/Feature 状态，框架需支持自定义字段渲染
- **标识符差异**：`lesson` 用 `Name` 而非 `Slug`，框架需适配两种命名

### Non-Functional Requirements

- CLI 输出格式与重构前完全一致（零行为变更）
- 现有测试全部通过（允许适配测试以使用新 API，但不降低覆盖率）

### Constraints & Dependencies

- Go 1.18+ 泛型
- 不改变 `cobra.Command` 的注册方式和 CLI 接口
- 不涉及 skill 层（SKILL.md）的修改

## Alternatives & Industry Benchmarking

### Industry Solutions

Go 生态中 CLI 命令复用通常通过 generic repository 或 shared package 实现，无成熟框架可参考。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 持续复制粘贴，维护负担线性增长 | Rejected: 即将新增命令 |
| 声明式框架（注册即生成 cobra.Command） | 自研 | 新增命令极度简洁 | 过度设计，proposal 的跨引用需要 hook 机制 | Rejected: YAGNI |
| 最小提取（仅 parseFrontmatter + sort） | — | 改动最小 | Discover/FindBySlug/渲染仍重复 | Rejected: 解决重复不彻底 |
| **共享工具包（泛型 Discover[T] + 渲染）** | 自研 | 消除 60-70% 重复，灵活性高 | 每个命令仍有 ~30 行胶水 | **Selected: 最优性价比** |

## Feasibility Assessment

### Technical Feasibility

Go 泛型完全支持 `Discover[T any]` 模式。现有三个命令的数据模型已是独立 struct，无需改动。

### Resource & Timeline

纯重构，无外部依赖。预计 4-6 个 coding task。

### Dependency Readiness

无外部依赖。`base/output.go` 已有 `PrintBlockStart/End`、`PrintField`、`CalcSlugColWidth`、`PadRight`、`TruncateSlug` 等工具函数可直接复用。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "声明式框架是最优解" | Occam's Razor | Overturned: 共享工具包更简单且足以满足需求，声明式框架增加复杂度但收益有限 |
| "三个命令完全同构" | Stress Test | Refined: `proposal` 有跨引用特例（PRD/Feature 状态检查），`lesson` 用 Name 而非 Slug，框架需预留灵活性 |

## Scope

### In Scope

- 创建 `pkg/infocmd/` 包，含泛型 `Discover[T]`、`FindBySlug[T]`、`RenderTable[T]`、`RenderDetail[T]`
- 共享 `parseFrontmatter` 函数
- 共享排序工具（按 created 降序，mtime 回退）
- 改造 `research` 命令使用新包
- 改造 `proposal` 命令使用新包（保留跨引用逻辑）
- 改造 `lesson` 命令使用新包
- 删除三个 `pkg/` 中的重复代码

### Out of Scope

- 新增 info-command（由未来工作完成）
- 改变 CLI 输出格式或接口
- 改变 skill 层（SKILL.md）
- `forge config`、`forge init`、`forge feature` 命令（属于 forge-info-commands 提案）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 泛型抽象无法覆盖 proposal 的跨引用特例 | Low | Medium | `RenderDetail` 接受自定义字段函数，proposal 在胶水代码中手动添加跨引用字段 |
| 重构引入输出格式差异 | Medium | High | 现有测试作为快照守卫，重构后逐命令对比输出 |
| `lesson` 用 Name 而非 Slug 导致 API 不统一 | Low | Low | 配置化的标识符提取函数（`IDKey func(T) string`） |
| `proposal` 的 sort 在 command 层而非 data 层，改造需移动 | Low | Low | 统一 sort 到 Discover 内部，通过 `SortKey func(T) string` 配置 |

## Success Criteria

- [ ] `parseFrontmatter` 仅存在于 `pkg/infocmd/` 一处
- [ ] `Discover[T]` + `FindBySlug[T]` 泛型函数可复用
- [ ] `research`、`proposal`、`lesson` 三个命令均使用 `pkg/infocmd/`
- [ ] CLI 输出与重构前逐字节一致
- [ ] 现有测试全部通过
- [ ] 新增 info-command 时只需定义 struct + 列配置 + ~30 行胶水代码

## Next Steps

- Proceed to task generation via `/quick-tasks`
