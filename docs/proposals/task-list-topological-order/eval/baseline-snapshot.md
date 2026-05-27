---
created: "2026-05-27"
author: "brainstorm"
status: Draft
---

# Proposal: forge task list 按依赖顺序拓扑排序 + TUI 树视图

## Problem

`forge task list` 目前按自然 ID 顺序（1, 2, 3...）以表格展示任务。开发者无法一眼看出任务的执行顺序和依赖关系链。

当前查看依赖路径的替代方式：
- `forge task check-deps` — 逐任务查看依赖是否满足，但无全局视图
- `forge task query <id>` — 查看单个任务的 `dependencies` 字段
- 手动阅读 `index.json` — 需要理解 JSON 结构

对于有 10+ 个任务且依赖链复杂的 feature，缺乏全局依赖视图导致开发者频繁在多个命令间切换，效率低下。

### Evidence

- `forge task list` 已是最常用的查看命令，但其输出与执行流程脱节
- 任务 Dependencies 字段已存在于数据模型（`Task.Dependencies []string`），但列表从未利用这些信息
- 在 `forge task claim` 的 `claimNextTask` 中存在自动排序逻辑（按优先级 + 版本号），但列表本身不反映这个执行顺序

### Urgency

低。这是一个增量改进，不阻塞现有工作流。但作为 `forge-cli v3.0.0` 周期内的体验优化，时机合适。

## Proposed Solution

### 核心变更

1. **变更默认排序**：`forge task list` 默认从自然 ID 排序改为拓扑排序（Kahn's algorithm）
2. **保留旧排序**：新增 `--sort id` 标志恢复自然 ID 排序
3. **新增 TUI 树视图**：新增 `--tree` 标志进入交互式 TUI 树视图

### 拓扑排序规则

- 使用 Kahn's algorithm，按 `Dependencies []string` 字段构建 DAG
- 同层（无依赖关系的任务）按自然 ID 排序以保持确定性
- 检测到环时输出警告，列表中标记 cyclic 节点
- 仅依赖直系依赖（非传递依赖）

### TUI 树视图（--tree）

- 基础交互：键盘上下左右导航、展开/折叠节点
- 颜色区分状态：完成=绿，进行中=黄，阻塞/失败=红，待处理=灰
- 作用域：当前 feature

### Usage

```bash
# 默认：拓扑排序表格
forge task list

# 恢复自然 ID 排序
forge task list --sort id

# TUI 树视图
forge task list --tree
```

### Innovation Highlights

本方案是对标准拓扑排序的常规使用，没有创新性技术突破。TUI 树视图从 `go-task`（Taskfile）和 `kubectl tree` 等工具的树状输出借鉴模式。

## Requirements Analysis

### Key Scenarios

**Happy path：**
- 5 个串行任务（1→2→3→4→5）：按 1,2,3,4,5 顺序显示
- 10 个并行+串行混合任务：先显示无依赖的任务组，依次向下游展开
- 点击 `--tree` 展开/折叠查看依赖子树

**Edge cases：**
- 无依赖的任务（孤立节点）：放在拓扑序最前面
- 通配符依赖（`1.x`）：展开为 1.1, 1.2, 1.3...
- Phase gate / summary 任务（`.gate`, `.summary`）：按语义依赖参与拓扑排序
- 空 feature：显示 "no tasks found"

**Error scenarios：**
- 依赖环：Kahn's algorithm 检测剩余节点 → 输出警告并在环中节点标记 `[cycle]`
- 缺失依赖：依赖指向不存在的任务 ID → 标记为 `[missing: <id>]`

### Non-Functional Requirements

- 拓扑排序在 100 个任务内应在 O(V+E) 时间内完成（Kahn's algorithm）
- `--tree` 模式使用纯 Go TUI 库（如 `bubbletea`），不依赖外部工具
- 颜色输出在非 TTY 环境自动禁用

### Constraints & Dependencies

- 依赖数据模型已存在于 `Task.Dependencies []string`，无需 schema 变更
- TUI 树需要引入新依赖（Go TUI library）
- 必须向后兼容管道输出（`forge task list | grep foo`）

## Alternatives & Industry Benchmarking

### Industry Solutions

- **go-task/Taskfile**: `task --list` 按文件名排列，`task --list-all` 显示依赖但无拓扑排序
- **Make**: `make -p` 显示数据库但不按拓扑序排列
- **kubectl tree**: 类似 `--tree` 模式的参考，显示 K8s 资源归属关系
- **npm/gradle**: `npm ls --all` / `gradle dependencies` 显示树状依赖但侧重包管理

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 开发者缺乏全局依赖视野 | Rejected: 低成本改进，ROI 高 |
| 新增 `--topo` 标志保持默认不变 | 包容性设计 | 不破坏现有习惯 | 用户需要额外输入参数才能获得更好的默认体验 | Rejected: 该选项已在 Scope 中作为 `--sort id` 保留 |
| 额外新增 ASCII DAG | 可视化深度 | 一目了然依赖拓扑 | 实现复杂，终端宽度有限时易错乱 | Deferred: 用户选择暂不纳入 |
| **变更默认排序 + --tree** | 本提案 | 默认展示执行流程，TUI 提供交互深度 | 少量习惯冲击 | **Selected: 根本解决依赖可见性问题** |

## Feasibility Assessment

### Technical Feasibility

技术可行：
- 拓扑排序算法已有多轮正确运行的 GO 标准实现
- TUI 树可使用 `bubbletea` 或 `gocui`（Go 社区成熟库）
- 当前任务模型已包含依赖信息

### Resource & Timeline

- 单工程攻坚：拓扑排序逻辑 1 天，TUI 树视图 2-3 天
- 无需跨团队协调

### Dependency Readiness

- 依赖数据已就位（`Task.Dependencies` 字段在所有 task .md 文件中）
- 无需外部服务依赖

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "用户需要默认看到拓扑排序" | Assumption Flip | 开发者的核心需求是"列出的任务反映执行顺序"而非"看到 ID 顺序"。但 `--sort id` 作为回退选项降低冲击风险 |
| "依赖信息足够构建完整 DAG" | XY Detection | `Task.Dependencies` 确实包含所有依赖信息，但注意通配符（`1.x`）需要在构建 DAG 时展开为具体任务 |
| "TUI 树视图是必要的" | Occam's Razor | 拓扑排序表格解决了 80% 的查看需求，TUI 树是额外的交互增强。用户明确选择两者都要 |

## Scope

### In Scope

- 实现 Kahn's algorithm 拓扑排序，支持环检测和告警
- `forge task list` 默认输出拓扑排序表格
- `--sort id` 标志恢复自然 ID 排序
- `--tree` 标志进入 TUI 树视图（基础交互：导航、折叠、颜色状态）
- 通配符依赖（`1.x`）展开参与拓扑排序
- 颜色渲染：完成=绿，进行中=黄，阻塞=红，待处理=灰
- 非 TTY 环境自动禁用颜色
- 环警告 + `[cycle]` 标记
- 缺失依赖 `[missing: <id>]` 标记

### Out of Scope

- 跨 feature 的全局 DAG 视图
- 传递依赖展示（仅直系依赖）
- ASCII DAG 纯文本连线图（─┬─ │ └─）
- TUI 中的编辑操作（修改任务状态等）
- 导出为图片/文件
- 管道输出中的颜色控制

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 改变默认排序破坏现有脚本依赖 ID 顺序 | Low | Medium | 非 TTY（管道）环境仍按 ID 排序，仅 TTY 启用拓扑排序；`--sort id` 始终可用 |
| TUI 库在 SSH/远程终端渲染异常 | Medium | Medium | 回退到表格模式，TUI 启动时检测终端能力 |
| 大数量任务（50+）拓扑排序不明显 | Low | Low | 分层缩进 + 同层按 ID 排序保持可扫描性 |
| 环检测不完善导致死循环 | Low | High | Kahn's algorithm 天然防环——剩余节点即为环成员，不会死循环 |

## Success Criteria

- [ ] `forge task list` 在 TTY 下默认输出拓扑排序表格，管道输出保持 ID 排序
- [ ] 拓扑排序结果正确：任意任务 B 依赖 A，A 出现在 B 之前
- [ ] `forge task list --sort id` 恢复自然 ID 排序
- [ ] `forge task list --tree` 进入 TUI 树视图，支持上下导航和展开/折叠
- [ ] 环检测：有环时输出警告，环中节点显示 `[cycle]` 标记
- [ ] 缺失依赖：依赖指向不存在的任务时显示 `[missing: <id>]` 标记
- [ ] 颜色正确区分 4 种状态（完成/进行中/阻塞/待处理）
- [ ] 非 TTY 环境不输出颜色
- [ ] 通配符依赖（`1.x`）正确展开参与拓扑排序
- [ ] 空 feature 输出 "no tasks found"
- [ ] 现有 `forge task list` 测试全部通过，新增覆盖拓扑排序和 TUI 树的测试

consistency_check_result:
  status: pass
  pairs_checked: 6
  conflicts_found: 0

## Next Steps

- Proceed to `/write-prd` to formalize requirements