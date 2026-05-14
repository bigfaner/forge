# Lesson: Vibe Coding 阶段的 Scope 控制与反思阈值

## 问题

deep-drill-analytics feature 在 Forge pipeline 完成后，产生了 23 个 vibe coding 提交。其中 3 个 `feat` 提交引入了新功能，尤其是 `375d28a` (+718 行) 添加了完整的 hook analysis section — 这本身就是一个独立的 feature，但作为 "fix 阶段的一部分" 直接提交了。

同时，23 个连续提交没有任何暂停反思点。5 个 fix 提交集中在 `detail.go`，5 个 fix/style 集中在 `dashboard.go`，这种聚集说明存在结构性问题，但直到全部修完才提取 convention。

**Why:** Vibe coding 模式下没有 scope guard 和 reflection checkpoint。agent 持续修复/增强，不区分"修复当前 feature"与"构建新 feature"，也不在问题聚集时暂停分析根因。

## 规则 1: 功能性提交的 scope guard

vibe coding 阶段（feature tasks 完成后的 fix/style 迭代），任何满足以下条件的提交**必须**走 Forge pipeline（至少创建 task，最好创建 proposal）：

- 新增 >50 行功能代码（不含 test、golden file）
- 引入新的 UI section 或面板
- 修改 parser/stats 层的数据结构（不仅是修 bug）

```bash
# 提交前自检
git diff --stat HEAD
# 如果某个 .go 文件 +50 lines 且不是 _test.go → 停下来，创建 task
```

**例外**: 纯 fix（修正已有功能的错误行为）和 style（调整已有元素的视觉表现）不受此限制。

## 规则 2: 连续修复反思阈值

当 vibe coding 阶段出现以下信号时，**暂停编码**，先提取 convention/lesson：

- **同文件聚集**: 同一个文件累积 3 个以上 fix/style 提交
- **同类型重复**: 同一类 bug（如对齐、溢出）在不同文件中重复出现 2 次以上
- **连续提交数**: vibe coding 阶段连续提交达到 5 个

```
fix(detail): path truncation      ← detail.go #1
fix(detail): scrollbar width      ← detail.go #2
fix(detail): column alignment     ← detail.go #3 → 停！提取 convention
```

**Why**: 连续修同一文件说明实现时缺乏规则指导。停下来提取 convention 后，后续同类 fix 可以避免。deep-drill-analytics 的 23 个提交中，如果在前 5 个时暂停提取规则，后 18 个中至少一半可以避免。

## 规则 3: Vibe coding 不修改 parser 层

vibe coding 阶段**禁止**修改 `internal/parser/` 和 `internal/stats/` 的核心逻辑。这些层的数据结构变更属于 feature scope，需要经过 design review。

参考: `4432aa6` 修改了 `jsonl.go` 的 `ScanSubAgentsDir`，`cb93d0d` 修改了 `stats.go` 的 hook 搜索逻辑 — 这些是 parser 层的 bug fix，但也可能是数据契约问题，应该在 parser 层有 unit test 覆盖。

## How to apply

### 手动（当前）

1. 每次 vibe coding 提交前，检查 `git diff --stat`
2. 如果出现 scope guard 触发条件，暂停并创建 `task add`
3. 每 5 个连续 fix/style 提交后，运行 `/learn-lesson` 提取 convention

### Forge 集成（后续）

1. `/execute-task` 在 task 完成后的 fix 迭代中，跟踪连续 fix 次数
2. 达到阈值时自动建议 "提取 convention 后继续"
3. 检测到 feat 提交时，提示 "是否创建新 task/proposal?"

## 预期效果

| 指标 | 改进前 | 改进后（预期） |
|------|--------|---------------|
| Vibe coding 提交总数 | 23 | 8-10 |
| Scope creep (不应出现的 feat) | 3 | 0 |
| Convention 提取时机 | 全部修完后 | 前 5 个提交时 |

## Tags

`scope-control`, `vibe-coding`, `process`, `forge-improvement`, `reflection`
