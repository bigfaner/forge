---
created: 2026-05-19
author: "fanhuifeng"
status: Draft
---

# Proposal: forge lesson 按时间倒序排列

## Problem

`forge lesson` 列表命令按文件系统默认顺序（文件名字母序）展示 lesson，用户无法快速找到最近添加的经验。

### Evidence

- `forge-cli/pkg/lesson/lesson.go` 的 `Discover()` 函数直接返回 `os.ReadDir` 的结果，无任何排序逻辑
- `forge-cli/internal/cmd/lesson.go` 的 `runLessonList()` 直接遍历输出
- 当前 78 个 lesson 文件按 `gotcha-*`/`lesson-*` 前缀字母序排列，与时间无关

### Urgency

低。功能缺失但不阻塞工作流。修复成本低（几行代码），顺手解决。

## Proposed Solution

在 `Discover()` 返回结果后，按文件修改时间（`os.FileInfo.ModTime()`）倒序排列。使用文件修改时间而非 frontmatter 日期，因为约 40 个文件使用 `created` 字段、20+ 个无 frontmatter、仅 1 个使用 `date` 字段，frontmatter 日期解析覆盖率太低。

同时在 Go 解析器中兼容 `created` 字段（模板使用的是 `created` 而非 `date`），修复字段名不一致的附带 bug。

### Innovation Highlights

无特别创新。标准的列表排序改进。

## Requirements Analysis

### Key Scenarios

- 用户运行 `forge lesson` 看到最新添加的 lesson 排在最前
- 无 frontmatter 或日期解析失败的 lesson 排在列表末尾（作为 fallback）
- `forge lesson <name>` 查看详情不受影响

### Non-Functional Requirements

- 性能：78 个文件的排序开销可忽略
- 兼容性：不改变 lesson 文件格式

### Constraints & Dependencies

- Go 标准库 `sort` 包即可完成
- 不引入新依赖

## Alternatives & Industry Benchmarking

### Industry Solutions

CLI 工具（`gh issue list`、`kubectl get events`）普遍默认按时间倒序列表。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 用户无法快速定位最新经验 | Rejected: 用户体验差 |
| frontmatter 日期排序 | — | 语义准确 | ~60% 文件无有效日期 | Rejected: 覆盖率低 |
| 文件名日期排序 | — | 无 IO 开销 | 文件名不含日期 | Rejected: 不可行 |
| **文件修改时间排序** | 行业惯例 | 100% 覆盖，零依赖 | 修改文件会改变排序 | **Selected: 覆盖率最高，实现最简** |

## Feasibility Assessment

### Technical Feisk

Go 标准库 `sort.Slice` + `os.FileInfo.ModTime()` 即可完成，无需额外依赖。

### Resource & Timeline

单任务，10 分钟内完成。

### Dependency Readiness

无外部依赖。

## Scope

### In Scope

- `forge lesson` 列表按文件修改时间倒序排列
- Go 解析器兼容 `created` frontmatter 字段（与模板一致）

### Out of Scope

- 批量修复已有 lesson 文件的 frontmatter 字段名
- 添加排序 flag（如 `--sort name`/`--sort date`）
- 其他 `forge` 命令的排序改进

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 文件修改时间不等于创建时间 | M | L | 对于 git 管理的项目可接受；未来可迁移到 frontmatter 日期 |

## Success Criteria

- [ ] `forge lesson` 输出按时间倒序（最新在前）
- [ ] `forge lesson <name>` 功能不受影响
- [ ] 现有测试全部通过
- [ ] 新增排序逻辑有对应测试覆盖

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
