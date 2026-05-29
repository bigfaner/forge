---
created: "2026-05-27"
updated: "2026-05-28"
author: "faner"
status: Draft
---

# Proposal: 统一枚举常量，消除魔法值

## Problem

forge-cli 中 Status（117 处）、Surface Type（~97 处）、Priority（36 处）三类枚举值大量以字符串字面量形式直接使用，而非引用已定义的常量。Surface Type 甚至无常量定义。这导致拼写错误无编译期保障，且常量定义分散在多个包中缺乏统一归属。

### Evidence

| 类别 | 魔法值数 | 常量引用数 | 常量定义位置 |
|------|---------|-----------|------------|
| Status | 117 | 5 | `pkg/feature/constants.go` |
| Surface Type | ~97 | 0 | 无 |
| Priority | 36 | 0 | `pkg/feature/constants.go` |

**Status 详解**（117 处，22 个非测试文件）：

| 文件 | 魔法值数 |
|------|---------|
| `pkg/task/statemachine.go` | 28（整张状态转移表 + isTerminalStatus） |
| `internal/cmd/task/submit.go` | 17 |
| `internal/cmd/task/validate_index.go` | 10 |
| `internal/cmd/task/claim.go` | 7 |
| `pkg/task/types.go` (NewTaskIndex) | 7 |
| `pkg/task/add.go` | 7 |
| `internal/cmd/task/tree.go` | 8 |
| `pkg/task/state.go` | 5 |
| `internal/cmd/cleanup.go` | 4 |
| `internal/cmd/feature/feature.go` | 6 |
| `internal/cmd/worktree/helpers.go` | 4 |
| 其他 11 个文件 | 1-3 各 |

**Surface Type 详解**（~97 处，6 个生产代码文件）：

| 文件 | 魔法值数 |
|------|---------|
| `pkg/forgeconfig/detect_surface.go` | ~50+（映射表 + 推断返回值 + 内联比较） |
| `pkg/forgeconfig/detect.go` | 5（KnownSurfaceTypes map） |
| `pkg/forgeconfig/execution_order.go` | 5（defaultExecutionOrder 切片） |
| `pkg/task/types.go` | 5（TestTypeTitle switch） |
| `pkg/task/autogen.go` | 3（uiSurfaceTypes map） |
| `internal/cmd/quality_gate.go` | 4（needsFullLifecycle + mobile 检查） |

**Priority 详解**（36 处，7 个生产代码文件）：

| 文件 | 魔法值数 |
|------|---------|
| `pkg/task/autogen.go` | 19（硬编码的 P1/P2 赋值） |
| `pkg/task/add.go` | 4 |
| `internal/cmd/task/validate_index.go` | 3 |
| `internal/cmd/task/claim.go` | 3 |
| `pkg/task/tasktemplate.go` | 2 |
| `pkg/task/types.go` | 3 |
| `internal/cmd/task/add.go` | 1 |
| `internal/cmd/quality_gate.go` | 1 |

### Urgency

中。代码库处于 v3.0.0 分支，此时统一枚举可避免魔法值在后续开发中继续扩散。

## Proposed Solution

新建 `pkg/types/` 包，定义 **typed constants**（`type Status string`、`type SurfaceType string`、`type Priority string`），并在全量替换魔法值的同时，修改所有相关结构体字段和函数签名，从 `string` 升级为具名类型。实现真正的编译期类型安全。

### Approach

**Typed Constants** 而非 untyped string constants。理由：
- 函数签名声明为 `Status` 类型后，编译器阻止传入任意 `string`
- IDE 自动补全和类型检查
- `type Status string` 在 Go 中仍兼容 JSON/YAML 序列化（行为零变更）

### Design Sketch

```go
// pkg/types/status.go
type Status string

const (
    StatusPending    Status = "pending"
    StatusInProgress Status = "in_progress"
    StatusCompleted  Status = "completed"
    StatusBlocked    Status = "blocked"
    StatusSuspended  Status = "suspended"
    StatusSkipped    Status = "skipped"
    StatusRejected   Status = "rejected"
)

func AllStatuses() []Status { ... }
func IsTerminalStatus(s Status) bool { ... }

// pkg/types/surface.go
type SurfaceType string

const (
    SurfaceWeb    SurfaceType = "web"
    SurfaceAPI    SurfaceType = "api"
    SurfaceCLI    SurfaceType = "cli"
    SurfaceTUI    SurfaceType = "tui"
    SurfaceMobile SurfaceType = "mobile"
)

func AllSurfaceTypes() []SurfaceType { ... }

// pkg/types/priority.go
type Priority string

const (
    PriorityP0 Priority = "P0"
    PriorityP1 Priority = "P1"
    PriorityP2 Priority = "P2"
)

func AllPriorities() []Priority { ... }
```

### Signature Migration Examples

```go
// Before
func (t *Task) SetStatus(status string) error
func needsFullLifecycle(surfaceType string) bool
terminalStatuses := map[string]bool{...}

// After
func (t *Task) SetStatus(status types.Status) error
func needsFullLifecycle(surfaceType types.SurfaceType) bool
terminalStatuses := map[types.Status]bool{...}
```

### Boundary Conversion

CLI flags、config 解析等外部接口仍使用 `string`，在边界处做一次性转换：

```go
status := types.Status(viper.GetString("status"))
surfaceType := types.SurfaceType(flagArg)
```

## Requirements Analysis

### Key Scenarios

- 开发者在 `statemachine.go` 中添加新状态转移：使用 `types.StatusCompleted`，编译器同时验证类型和拼写
- 开发者添加新的 Surface Type：只需在 `pkg/types/surface.go` 中增加常量，所有引用点自动生效
- 新人阅读代码：`types.StatusPending` 比 `"pending"` 更明确；`func SetStatus(status types.Status)` 比 `func SetStatus(status string)` 文档性更强

### Non-Functional Requirements

- **零行为变更**：所有 CLI 命令输入输出不变（JSON 序列化行为不变）
- **构建稳定**：`go build ./...` 和 `go test ./...` 全部通过
- **无循环依赖**：`pkg/types/` 不导入任何 forge-cli 内部包（纯类型定义）

### Constraints & Dependencies

- `pkg/types/` 必须是叶包（零内部依赖），被 `pkg/feature`、`pkg/task`、`pkg/forgeconfig`、`internal/cmd` 共同导入
- Task Type 常量保留在 `pkg/task/types.go`，不迁移（与 task 逻辑深度耦合）
- `type Status string` 在 Go 中 JSON/YAML 序列化行为与 `string` 完全一致

## Alternatives & Industry Benchmarking

### Comparison Table

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| Do nothing | 零成本 | 250+ 魔法值持续累积 | Rejected |
| 只替换魔法值，不改签名 | 改动较小 | 仍是 `string`，无类型安全 | Rejected: 治标不治本 |
| **Typed constants + 全量签名升级** | 真正的类型安全，IDE 支持 | 改动量大（~250 替换 + 签名变更） | **Selected** |
| 分步走：先 untyped 后 typed | 每次变更最小 | 两轮大规模替换，中间状态不理想 | Rejected: 不如一次到位 |

## Feasibility Assessment

### Technical Feasibility

机械替换 + 签名变更。编译器会捕获所有遗漏（类型不匹配 → 编译错误）。Go 的 `type X string` 保证 JSON 序列化兼容。

### Resource & Timeline

估计 10-15 个任务，5-8 小时。建议走 full pipeline（PRD → tech-design → breakdown-tasks）。

## Assumptions Challenged

| Assumption | Finding |
|------------|---------|
| "Status 魔法值约 89 处" | **Overturned**: 实际 117 处，分散在 22 个文件中 |
| "Priority 魔法值约 15 处" | **Overturned**: 实际 36 处，分布在 7 个文件 |
| "Surface Type 约 70+ 处" | **Confirmed but underestimated**: 实际 ~97 处 |
| "现有 9 处常量引用" | **Overturned**: 仅 5 处常量引用，Priority 常量引用为 0 |
| "pkg/template/template.go 存在" | **Overturned**: 该文件不存在，实际为 `pkg/task/tasktemplate.go` |
| "untyped constants 足够" | **Challenged**: typed constants 提供真正的编译期类型安全，值得额外代价 |

## Scope

### In Scope

**A. 新建 `pkg/types/` 包（typed constants）**
- `status.go`：`type Status string` + 7 个常量 + `AllStatuses()` + `IsTerminalStatus()`
- `surface.go`：`type SurfaceType string` + 5 个常量 + `AllSurfaceTypes()`
- `priority.go`：`type Priority string` + 3 个常量 + `AllPriorities()`

**B. 迁移现有常量**
- `pkg/feature/constants.go`：移除 Status 和 Priority 常量定义，改为从 `pkg/types/` 重导出
- 所有 `feature.StatusXxx` / `feature.PriorityXxx` 引用改为 `types.StatusXxx` / `types.PriorityXxx`

**C. 替换 Status 魔法值 + 升级签名（117 处，22 文件）**
主要受影响文件（完整列表）：
- `pkg/task/statemachine.go`：28 处（状态转移表 + isTerminalStatus）
- `internal/cmd/task/submit.go`：17 处
- `internal/cmd/task/validate_index.go`：10 处（validStatus map → `types.AllStatuses()`）
- `internal/cmd/task/claim.go`：7 处
- `pkg/task/types.go`：7 处（StatusEnum → `[]types.Status{}`）
- `pkg/task/add.go`：7 处（terminalStatuses map → `map[types.Status]bool{}`）
- `internal/cmd/task/tree.go`：8 处
- `pkg/task/state.go`：5 处
- `internal/cmd/cleanup.go`：4 处
- `internal/cmd/feature/feature.go`：6 处
- `internal/cmd/worktree/helpers.go`：4 处
- `internal/cmd/task/reopen.go`：3 处
- `pkg/task/build.go`：3 处
- `pkg/task/deps.go`：2 处（satisfiedStatuses map → `map[types.Status]bool{}`）
- `pkg/task/record.go`：2 处
- `internal/cmd/task/transition.go`：2 处
- `pkg/task/autogen.go`：1 处
- `pkg/task/index.go`：1 处
- `internal/cmd/task/add.go`：1 处
- `internal/cmd/quality_gate.go`：3 处
- `internal/cmd/verify_task_done.go`：1 处
- `internal/cmd/feature/feature_complete.go`：1 处

**D. 替换 Priority 魔法值 + 升级签名（36 处，7 文件）**
- `pkg/task/autogen.go`：19 处
- `pkg/task/add.go`：4 处
- `internal/cmd/task/validate_index.go`：3 处（validPriority map → `types.AllPriorities()`）
- `internal/cmd/task/claim.go`：3 处（priorityOrder map → `map[types.Priority]int{}`）
- `pkg/task/tasktemplate.go`：2 处
- `pkg/task/types.go`：3 处（PriorityEnum → `[]types.Priority{}`）
- `internal/cmd/task/add.go`：1 处
- `internal/cmd/quality_gate.go`：1 处

**E. 替换 Surface Type 魔法值 + 升级签名（~97 处，6 文件）**
- `pkg/forgeconfig/detect_surface.go`：~50+ 处（映射表 + 推断函数返回值）
- `pkg/forgeconfig/detect.go`：5 处（KnownSurfaceTypes → `map[types.SurfaceType]bool{}`）
- `pkg/forgeconfig/execution_order.go`：5 处（defaultExecutionOrder → `[]types.SurfaceType{}`）
- `pkg/task/types.go`：5 处（TestTypeTitle 签名 + switch）
- `pkg/task/autogen.go`：3 处（uiSurfaceTypes → `map[types.SurfaceType]bool{}`）
- `internal/cmd/quality_gate.go`：4 处（needsFullLifecycle 签名）

### Out of Scope

- Task Type 常量迁移（保留在 `pkg/task/types.go`，与 task 逻辑深度耦合）
- `pkg/forgeconfig/config.go` 中的 Coverage Config 默认键（`"coding.feature"` 等）—— 循环依赖约束
- 路径常量（`"prd"`、`"design"` 等）—— 属于其他提案范围
- Config dotpath 键（`"eval.proposal"` 等）—— 非枚举值
- 测试文件中的魔法值替换（可在主任务完成后单独处理）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| typed constants 导致 JSON 序列化行为变化 | L | H | Go `type X string` 保持 string 序列化行为，但需在 CI 中验证 |
| 签名变更遗漏导致编译失败 | L | L | `go build ./...` 立即捕获，类型系统是编译期检查 |
| 重导出层（pkg/feature）遗漏 | L | M | 迁移后全局 grep `feature.Status` 确认无遗漏 |
| 大规模改动引入 merge conflict | M | M | 基于 main 创建干净分支，完成后尽快合并 |
| detect_surface.go 映射表改动导致 surface 检测逻辑回归 | L | H | 常量值与原始字符串完全相同，仅类型包装，行为零变更 |

## Success Criteria

- [ ] `pkg/types/` 包存在，定义 `type Status string`（7 常量）、`type SurfaceType string`（5 常量）、`type Priority string`（3 常量）
- [ ] `pkg/types/` 不导入任何 forge-cli 内部包（纯类型定义）
- [ ] `pkg/feature/constants.go` 中 Status 和 Priority 常量已迁移，保留重导出兼容
- [ ] 所有 Status 相关结构体字段和函数参数使用 `types.Status` 类型
- [ ] 所有 Surface Type 相关结构体字段和函数参数使用 `types.SurfaceType` 类型
- [ ] 所有 Priority 相关结构体字段和函数参数使用 `types.Priority` 类型
- [ ] Status 魔法值从 117 处降至 0
- [ ] Surface Type 魔法值从 ~97 处降至 0
- [ ] Priority 魔法值从 36 处降至 0
- [ ] `pkg/task/statemachine.go` 中状态转移表全部使用 `types.StatusXxx`
- [ ] `internal/cmd/task/validate_index.go` 中 validStatus/validPriority map 改用 `types.AllStatuses()`/`types.AllPriorities()`
- [ ] `go build ./...` 零错误
- [ ] `go test ./...` 全部通过

## Next Steps

- Proceed to `/write-prd` → `/tech-design` → `/breakdown-tasks` (full pipeline)
