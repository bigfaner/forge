---
created: 2026-05-20
author: "faner + Claude"
status: Draft
---

# Proposal: CLI `created` 字段统一与表格显示优化

## Problem

Forge CLI 的三种文档类型（lesson / proposal / feature）在时间戳处理和排序上不一致，且列表中 slug 列被截断导致不可辨认。

### Evidence

- **时间戳不一致**: Proposal 按 frontmatter `created` 排序；Lesson 按文件 mtime 排序；Feature 按 manifest mtime 排序。三种策略，三种行为。
- **mtime 不可靠**: 从 git clone/pull 获取的文件，mtime 是 clone 时间而非创建时间。Lesson 已有 `created` frontmatter 却未用于排序。
- **Feature 无 `created` 字段**: 标准 manifest 模板和 quick manifest 模板均无 `created` frontmatter。
- **slug 截断**: 最长 slug 达 42 字符（`profile-aware-shared-infra-precise-staging`），而 proposal/feature 列宽固定为 30，lesson 为 35。

### Urgency

低优先级。不阻塞任何工作流，但影响 CLI 可用性。每次用 `forge feature list` 或 `forge proposal` 时都会遇到截断问题。

## Proposed Solution

两处改动：

1. **统一 `created` 字段**: 所有文档类型的 frontmatter 都包含 `created`，排序统一使用 `created` 降序，mtime 仅作 fallback。
2. **动态列宽**: CLI 表格的 slug 列宽度根据实际数据动态计算（取最长 slug 长度），设置最小 30 字符、最大 60 字符边界。

### Innovation Highlights

无创新。这是标准的"补齐字段 + 统一排序"维护工作。

## Requirements Analysis

### Key Scenarios

- `forge lesson` / `forge proposal` / `forge feature list` 列表按创建时间降序排列，行为一致可预测
- `forge proposal` 列表显示完整的 `profile-aware-shared-infra-precise-staging` 而非 `profile-aware-shared-infra-p...`
- 已有 feature manifest 无 `created` 字段时，fallback 到 mtime，不影响存量数据

### Non-Functional Requirements

- 向后兼容：缺少 `created` 的旧文档 fallback 到 mtime，不报错

## Alternatives & Industry Benchmarking

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 排序不一致、slug 截断 | Rejected: 每次使用都受影响 |
| `created` 统一 + 动态列宽 | — | 一致性强、显示完整 | 需改模板和 Go 代码 | **Selected: 改动小，收益明确** |
| 只补 `created`，不改显示 | — | 最小改动 | 截断问题仍在 | Rejected: 截断是直接痛点 |

## Feasibility Assessment

### Technical Feasibility

完全可行。改动集中在：
- 2 个 manifest 模板（加 `created` 字段）
- 3 个 Go pkg（lesson/proposal/feature 的排序逻辑）
- 3 个 CLI cmd（表格列宽逻辑）

### Resource & Timeline

小改动，1-2 个 coding task。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 需要内置模糊搜索 | Occam's Razor + XY Detection | **Overturned**: `forge proposal \| grep <keyword>` 已满足过滤需求，无需在 CLI 内部重写 |
| mtime 可以作为排序依据 | 5 Whys | **Overturned**: git clone 会重置 mtime，frontmatter `created` 是唯一可靠时间源 |
| slug 需要命名规范来缩短 | XY Detection | **Overturned**: 问题不是 slug 太长，而是显示列宽不够 |

## Scope

### In Scope

- 标准 manifest 模板（`write-prd/templates/manifest.md`）加 `created` frontmatter
- Quick manifest 模板（`quick-tasks/templates/manifest-quick.md`）加 `created` frontmatter
- `forge feature list` 排序改为 `created` 降序（mtime fallback）
- `forge lesson` 排序改为 `created` 降序（mtime fallback）
- `forge proposal` 排序逻辑保持不变（已是 `created` 降序）
- 三个列表命令的 slug 列改为动态宽度（min 30, max 60）

### Out of Scope

- 内置模糊搜索（`forge X | grep` 足够）
- 跨类型统一搜索命令
- 存量 feature manifest 的 `created` 回填（新 feature 才加，旧 feature fallback mtime）
- FZF 风格模糊匹配
- Slug 命名规范调整

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 存量 feature 无 `created` | H | L | mtime fallback 已覆盖 |
| 动态列宽导致表格过宽 | L | L | max 60 字符上限控制 |

## Success Criteria

- [ ] `forge feature list` 输出中无 slug 被截断（最长 42 字符 slug 完整显示）
- [ ] `forge lesson` 列表按 `created` 降序排列
- [ ] `forge feature list` 列表按 `created` 降序排列
- [ ] 新建的 feature manifest 包含 `created` frontmatter
- [ ] 缺少 `created` 的旧文档 fallback 到 mtime，不报错

## Next Steps

- 本提案范围小，适合 `/quick` 流程
