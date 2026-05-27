---
created: "2026-05-27"
author: "faner"
status: Draft
---

# Proposal: 统一枚举常量，消除魔法值

## Problem

forge-cli 中 Status（89 处）、Surface Type（70+ 处）、Priority（~15 处）三类枚举值大量以字符串字面量形式直接使用，而非引用已定义的常量。Surface Type 甚至无常量定义。这导致拼写错误无编译期保障，且常量定义分散在多个包中缺乏统一归属。

### Evidence

| 类别 | 魔法值数 | 常量引用数 | 常量定义位置 |
|------|---------|-----------|------------|
| Status | 89 | 9 | `pkg/feature/constants.go` |
| Surface Type | 70+ | 0 | 无 |
| Priority | ~15 | 极少 | `pkg/feature/constants.go` |
| Task Type | 14 | 53 | `pkg/task/types.go` |

Status 常量已定义但几乎没人用。最严重的 `pkg/task/statemachine.go` 中整张状态转移表全是原始字符串。`internal/cmd/task/submit.go` 有 15+ 处魔法值。

### Urgency

中。代码库处于 v3.0.0 分支，此时统一枚举可避免魔法值在后续开发中继续扩散。与 `forge-cli-clean-code` 提案互补——常量统一后，重复的状态判断逻辑更容易合并。

## Proposed Solution

新建 `pkg/types/` 包作为所有共享枚举常量的唯一归属。将 `pkg/feature/constants.go` 中的 Status 和 Priority 常量迁移至此，新增 Surface Type 常量。然后逐一替换所有魔法值为常量引用。

### Innovation Highlights

无。标准的 Go 常量组织实践。

## Requirements Analysis

### Key Scenarios

- 开发者在 `statemachine.go` 中添加新状态转移：使用 `types.StatusCompleted` 而非手写 `"completed"`，拼写错误会被编译器捕获
- 开发者添加新的 Surface Type 支持：只需在 `pkg/types/` 中增加一个常量，所有引用点自动生效
- 新人阅读代码：看到 `types.StatusPending` 立即知道是状态枚举，而非不确定 `"pending"` 的含义

### Non-Functional Requirements

- **零行为变更**：所有 CLI 命令输入输出不变
- **构建稳定**：`go build ./...` 和 `go test ./...` 全部通过
- **无循环依赖**：`pkg/types/` 不导入任何 forge-cli 内部包（纯常量定义）

### Constraints & Dependencies

- `pkg/task` → `pkg/forgeconfig` 已存在依赖方向，`pkg/forgeconfig` 不能导入 `pkg/task`
- `pkg/types/` 必须是叶包（零内部依赖），被 `pkg/feature`、`pkg/task`、`pkg/forgeconfig`、`internal/cmd` 共同导入
- Task Type 常量保留在 `pkg/task/types.go`，不迁移（与 task 逻辑深度耦合）

## Alternatives & Industry Benchmarking

### Industry Solutions

Go 社区标准做法：将共享枚举常量放在独立的 `types/` 或 `constants/` 包中，所有业务包引用此包。`kubernetes` 项目即采用 `pkg/api/types/` 模式。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 魔法值持续累积，拼写错误风险 | Rejected: 与代码清理目标矛盾 |
| 只加 Surface Type 常量 | — | 改动小 | Status/Priority 仍分散，未解决主要问题 | Rejected: 覆盖面不足 |
| 在 pkg/feature/ 加 Surface Type | 现有模式 | 无新包 | pkg/feature 语义不符（surface 不是 feature） | Rejected: 语义混乱 |
| **新建 pkg/types/ + 全量迁移** | Go 社区实践 | 枚举集中、依赖清晰、可扩展 | 改动量大（~180 处替换） | **Selected: 唯一能彻底解决问题的方案** |

## Feasibility Assessment

### Technical Feasibility

纯机械替换。编译器会捕获所有遗漏（未定义的常量 → 编译错误）。

### Resource & Timeline

约 6-8 个任务，预计 2-3 小时。

### Dependency Readiness

无外部依赖。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "Status 常量放在 pkg/feature/ 是合理的" | 5 Whys | Overturned: Status 是任务状态概念，不是 feature 特有的。pkg/task 和 pkg/forgeconfig 都需要使用，放在 feature 包造成语义混乱和依赖困难 |
| "现有 9 处常量引用需要保持向后兼容" | Assumption Flip | Confirmed 但不适用: forge-cli 是 CLI 工具，不是库。外部插件（skills/commands）是 YAML+prompt 形式，不导入 Go 包。引用迁移无外部影响 |
| "Surface Type 无常量是因为用法太分散不值得定义" | Stress Test | Overturned: 70+ 处使用恰恰是最需要常量的场景。添加新 Surface Type 时需修改 6+ 文件，常量化后只需改 1 处 |

## Scope

### In Scope

**A. 新建 `pkg/types/` 包**
- 定义 Status 常量：`StatusPending`、`StatusInProgress`、`StatusCompleted`、`StatusBlocked`、`StatusSuspended`、`StatusSkipped`、`StatusRejected`
- 定义 Priority 常量：`PriorityP0`、`PriorityP1`、`PriorityP2`
- 定义 Surface Type 常量：`SurfaceWeb`、`SurfaceAPI`、`SurfaceCLI`、`SurfaceTUI`、`SurfaceMobile`
- 提供 `AllStatuses()`、`AllPriorities()`、`AllSurfaceTypes()` 辅助函数
- 提供 `IsTerminalStatus()` 辅助函数（统一 `add.go` 和 `statemachine.go` 中的重复定义）

**B. 迁移现有常量**
- `pkg/feature/constants.go`：移除 Status 和 Priority 常量定义，改为从 `pkg/types/` 重导出（`type alias` 或 `var` 赋值），保持 `feature.StatusXxx` 可用
- 所有现有 `feature.StatusXxx` / `feature.PriorityXxx` 引用点改为 `types.StatusXxx` / `types.PriorityXxx`

**C. 替换 Status 魔法值（89 处，~18 文件）**
主要受影响文件：
- `pkg/task/statemachine.go`：整张状态转移表
- `pkg/task/state.go`：switch/case
- `pkg/task/add.go`：`terminalStatuses` map + 默认赋值
- `pkg/task/deps.go`：`satisfiedStatuses` map
- `pkg/task/record.go`：状态检查
- `pkg/task/build.go`：默认状态 + 路径常量
- `pkg/task/autogen.go`：默认状态
- `pkg/task/types.go`：`NewTaskIndex()` 中的 StatusEnum
- `internal/cmd/task/submit.go`：15+ 处
- `internal/cmd/task/claim.go`：8+ 处
- `internal/cmd/task/reopen.go`：3 处
- `internal/cmd/task/transition.go`：2 处
- `internal/cmd/task/validate_index.go`：重定义的 validStatus map → 改用 `types.AllStatuses()`
- `internal/cmd/task/add.go`：默认状态
- `internal/cmd/cleanup.go`：状态过滤
- `internal/cmd/verify_task_done.go`：状态检查
- `internal/cmd/quality_gate.go`：状态过滤 + Surface Type 判断
- `internal/cmd/feature/feature.go`：内联状态数组
- `internal/cmd/feature/feature_complete.go`：状态赋值
- `internal/cmd/worktree/helpers.go`：状态比较

**D. 替换 Priority 魔法值（~15 处，~8 文件）**
- `pkg/task/types.go`：`NewTaskIndex()` PriorityEnum
- `pkg/task/add.go`：优先级验证
- `pkg/task/autogen.go`：多处默认赋值
- `internal/cmd/task/validate_index.go`：重定义的 validPriority map → 改用 `types.AllPriorities()`
- `internal/cmd/task/claim.go`：优先级排序
- `internal/cmd/task/add.go`：默认优先级
- `internal/cmd/quality_gate.go`：优先级判断
- `pkg/template/template.go`：默认优先级

**E. 替换 Surface Type 魔法值（70+ 处，~6 文件）**
- `pkg/forgeconfig/detect_surface.go`：库检测映射表（最大量）
- `pkg/forgeconfig/detect.go`：`KnownSurfaceTypes` map
- `pkg/forgeconfig/execution_order.go`：`defaultExecutionOrder` 切片
- `pkg/task/autogen.go`：`uiSurfaceTypes` map
- `pkg/task/types.go`：`TestTypeTitle()` switch
- `internal/cmd/quality_gate.go`：`needsFullLifecycle()` 判断

### Out of Scope

- Task Type 常量迁移（保留在 `pkg/task/types.go`，与 task 逻辑深度耦合）
- `pkg/forgeconfig/config.go` 中的 Coverage Config 默认键（`"coding.feature"` 等）—— 这些是 task type 字符串，非 surface type/status，且存在循环依赖约束
- 路径常量（`"prd"`、`"design"` 等）—— 属于 `forge-cli-clean-code` 提案范围
- Config dotpath 键（`"eval.proposal"` 等）—— 非枚举值，是嵌套配置路径

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 迁移常量时遗漏引用导致编译失败 | L | L | `go build ./...` 立即捕获 |
| 重导出层（pkg/feature）遗漏导致下游编译失败 | L | M | 迁移后全局 grep `feature.Status` 确认无遗漏 |
| Surface Type 常量化后 detect_surface.go 映射表行为变化 | L | H | 常量值与原始字符串完全相同，仅引用方式变化，行为零变更 |
| 误将非枚举字符串当作魔法值替换 | M | M | 只替换明确匹配枚举值的字面量，不碰相似但含义不同的字符串 |

## Success Criteria

- [ ] `pkg/types/` 包存在，定义 Status（7 个）、Priority（3 个）、Surface Type（5 个）常量
- [ ] `pkg/feature/constants.go` 中 Status 和 Priority 常量已迁移，保留重导出兼容
- [ ] Status 魔法值从 89 处降至 0
- [ ] Surface Type 魔法值从 70+ 处降至 0
- [ ] Priority 魔法值从 ~15 处降至 0
- [ ] `pkg/task/statemachine.go` 中状态转移表全部使用 `types.StatusXxx`
- [ ] `internal/cmd/task/validate_index.go` 中 `validStatus`/`validPriority` map 改用 `types.AllStatuses()`/`types.AllPriorities()`
- [ ] `pkg/types/` 不导入任何 forge-cli 内部包（纯常量定义）
- [ ] `go build ./...` 零错误
- [ ] `go test ./...` 全部通过

## Next Steps

- Proceed to `/quick-tasks` to generate task breakdown
