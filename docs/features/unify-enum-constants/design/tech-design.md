---
created: "2026-05-28"
prd: prd/prd-spec.md
status: Draft
---

# Technical Design: unify-enum-constants

## Overview

将 forge-cli 中 250+ 处枚举魔法值替换为 typed constants（`type Status string` 等），集中定义在 `pkg/types/` 叶包中，并修改所有相关结构体字段和函数签名。实现编译期类型安全，同时保持零行为变更。

核心策略：**按 package 分批迁移**（而非按枚举类别），避免重复修改同一文件。

## Architecture

### Layer Placement

单层重构——仅涉及 `pkg/` 和 `internal/` 层。`pkg/types/` 作为新增叶包被所有上层包导入。

依赖方向（严格遵守）：

```
internal/cmd/* → pkg/task, pkg/forgeconfig, pkg/types
pkg/task       → pkg/types
pkg/forgeconfig → pkg/types
pkg/feature    → pkg/types
pkg/types      → (无内部依赖，叶包)
```

### Component Diagram

```
┌─────────────────────────────────────────────────────┐
│                   internal/cmd/*                     │
│   (submit, claim, validate_index, quality_gate, ...) │
│   surfaces.go (重导出 KnownSurfaceTypes)              │
│         使用 types.Status / types.SurfaceType         │
│              / types.Priority                         │
└────────────┬──────────────┬──────────────┬───────────┘
             │              │              │
             ▼              ▼              ▼
┌─────────────────┐ ┌──────────────┐ ┌──────────────┐
│   pkg/task/     │ │pkg/forgeconfig│ │ pkg/feature/ │
│ statemachine    │ │ detect.go     │ │ constants.go │
│ add, state,     │ │ (KnownSurface │ │ (重导出)      │
│ deps, build,    │ │  Types 定义)  │ │              │
│ autogen, types  │ │ detect_surface│ │              │
│                 │ │ execution_    │ │              │
│                 │ │ order         │ │              │
└────────┬────────┘ └──────┬───────┘ └──────┬───────┘
         │                 │                │
         ▼                 ▼                ▼
      ┌─────────────────────────────────────────┐
      │            pkg/types/ (NEW)              │
      │  status.go    │ surface.go │ priority.go │
      │  (叶包：零内部依赖)                       │
      └─────────────────────────────────────────┘
```

### Dependencies

- 无新增外部依赖
- `pkg/types/` 是纯 Go 标准库代码（不导入任何第三方包）

## Interfaces

### Interface 1: Status Type

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

func AllStatuses() []Status
func IsTerminalStatus(s Status) bool
```

> **行为变更声明**：`IsTerminalStatus` 包含 `completed`、`skipped`、`rejected` 三种终态，与业务规则 `BIZ-task-lifecycle-001` 一致。但当前 `statemachine.go` 中的私有函数 `isTerminalStatus` 仅检查 `completed` 和 `rejected`（不含 `skipped`）。`add.go` 中的 `terminalStatuses` map 则包含全部三种。迁移后统一为 `types.IsTerminalStatus()`，将 `skipped` 纳入终态，影响 `isActiveFixTask`/`canAutoUnblock` 逻辑——`skipped` 的 fix-task 将被视为非活跃。这是**对齐业务规则的 bug 修复**，不是纯重构。需在 Phase 3 单独编写测试验证 `canAutoUnblock` 在 source task 为 skipped 时的行为。

### Interface 2: SurfaceType Type

```go
// pkg/types/surface.go
type SurfaceType string

const (
    SurfaceWeb    SurfaceType = "web"
    SurfaceAPI    SurfaceType = "api"
    SurfaceCLI    SurfaceType = "cli"
    SurfaceTUI    SurfaceType = "tui"
    SurfaceMobile SurfaceType = "mobile"
)

func AllSurfaceTypes() []SurfaceType
```

### Interface 3: Priority Type

```go
// pkg/types/priority.go
type Priority string

const (
    PriorityP0 Priority = "P0"
    PriorityP1 Priority = "P1"
    PriorityP2 Priority = "P2"
)

func AllPriorities() []Priority
```

### Interface 4: Re-export (pkg/feature/constants.go)

```go
// 向后兼容重导出
type Status = types.Status
type Priority = types.Priority

const (
    StatusPending    = types.StatusPending
    StatusInProgress = types.StatusInProgress
    StatusCompleted  = types.StatusCompleted
    StatusBlocked    = types.StatusBlocked
    StatusSuspended  = types.StatusSuspended
    StatusSkipped    = types.StatusSkipped
    StatusRejected   = types.StatusRejected
)

const (
    PriorityP0 = types.PriorityP0
    PriorityP1 = types.PriorityP1
    PriorityP2 = types.PriorityP2
)
```

使用 type alias（`=`）确保 `feature.Status` 与 `types.Status` 是同一类型，无需类型转换。

## Data Models

db-schema: no — 无数据库变更。

### Model 1: Status State Machine

`statemachine.go` 中 `TransitionRule` 结构体的 `From`/`To` 字段类型从 `string` 升级为 `types.Status`：

```go
// Before (pkg/task/statemachine.go)
type TransitionRule struct {
    From     string         // current status; "*" matches any
    To       string         // target status; "*" matches any
    Role     TransitionRole // required role; "" matches any
    Allowed  bool
    GuardMsg string
}

// After
type TransitionRule struct {
    From     types.Status   // current status; "*" matches any
    To       types.Status   // target status; "*" matches any
    Role     TransitionRole // required role; "" matches any
    Allowed  bool
    GuardMsg string
}
```

> **注**：通配符 `"*"` 是 `transitionTable` 中的特殊值，需保留为 `types.Status("*")` 或定义常量 `StatusAny`。`matchRule` 函数需同步调整比较逻辑。

### Model 2: Task Struct Fields

```go
// Before
Status   string
Priority string

// After
Status   types.Status
Priority types.Priority
```

### Model 3: SurfaceType in ForgeConfig

```go
// Before (map keys are string)
KnownSurfaceTypes map[string]bool
surfacePriority   map[string]int
defaultExecutionOrder []string

// After (map keys are typed)
KnownSurfaceTypes map[types.SurfaceType]bool
surfacePriority   map[types.SurfaceType]int
defaultExecutionOrder []types.SurfaceType
```

## Error Handling

### Error Types & Codes

无新增错误类型。类型不匹配在编译期被捕获（`go build`），非运行时错误。

### Propagation Strategy

`TransitionError` 结构体（`statemachine.go`）持有 `From string` 和 `To string` 字段。迁移后这两个字段将变为 `From types.Status` 和 `To types.Status`。由于 `type Status string`，`fmt.Sprintf("... %s -> %s ...", e.From, e.To, ...)` 的输出格式不变——`%s` 格式化仍输出原始字符串值。

但需注意：现有测试中若有 `assert.Contains(t, err.Error(), "completed -> pending")` 之类的字符串断言，由于 `TransitionError.From`/`To` 字段类型变更，赋值方式从字面量变为常量，值不变但编译器类型检查更严格。建议在 Phase 3 迁移 `statemachine.go` 后立即运行 `go test ./pkg/task/...` 验证错误格式兼容性。

### Error Structs Holding Enum Values

| Struct | File | Fields | Impact |
|--------|------|--------|--------|
| `TransitionError` | `pkg/task/statemachine.go` | `From string`, `To string` | 字段类型变为 `types.Status`；`Error()` 输出格式不变 |
| `ActiveFixExistsError` | `pkg/task/add.go` | `SourceTaskID string` | 不涉及枚举字段，无影响 |

## Cross-Layer Data Map

Single-layer feature — 不涉及跨层数据流。所有变更都在 Go 代码层内。

## Integration Specs

No existing-page integrations — 不适用（纯后端重构）。

## Testing Strategy

### Per-Layer Test Plan

| Layer | Test Type | Tool | What to Test | Coverage Target |
|-------|-----------|------|--------------|-----------------|
| `pkg/types/` | Unit | go test | 常量值正确性、AllXxx 返回完整列表、IsTerminalStatus 判断正确 | 100% |
| `pkg/types/` | Unit | go test | JSON 序列化/反序列化兼容性（`type X string` 仍输出为 plain string） | 100% |
| 全局 | Build | go build | 所有类型签名一致，零编译错误 | 零错误 |
| 全局 | Regression | go test ./... | 行为零变更——所有现有测试通过 | 现有覆盖率不降 |

### Key Test Scenarios

1. **JSON 序列化兼容**：`types.Status("pending")` 序列化为 `"pending"`（无额外包装）
2. **AllStatuses 完整性**：返回 7 个常量，与 `pkg/feature/constants.go` 原始定义一致
3. **AllSurfaceTypes 完整性**：返回 5 个常量，与 `KnownSurfaceTypes` map 一致
4. **IsTerminalStatus**：`completed`、`skipped`、`rejected` 返回 true；其余返回 false
5. **Re-export 等价性**：`feature.StatusPending == types.StatusPending` 为 true（type alias 保证）

### Overall Coverage Target

`pkg/types/` 新增代码 100%；现有测试全部通过（不降覆盖率）。

## Security Considerations

### Threat Model

无安全风险。纯代码组织重构，不改变外部接口或数据流。

### Mitigations

不适用。

## PRD Coverage Map

| PRD Requirement / AC | Design Component | Interface / Model |
|----------------------|------------------|-------------------|
| US1: typed Status constants 防止拼写错误 | `pkg/types/status.go` | Interface 1: Status |
| US1: statemachine.go 使用 types.StatusXxx | `pkg/task/statemachine.go` 签名升级 | Model 1: TransitionRule |
| US2: SurfaceType 集中定义 | `pkg/types/surface.go` | Interface 2: SurfaceType |
| US2: 新增 Surface Type 只改一处 | `AllSurfaceTypes()` + 常量定义 | Interface 2 |
| US3: go build 验证枚举引用完整性 | 全局编译 | Testing Strategy |
| US3: pkg/types/ 是叶包 | pkg/types/ 零导入 | Architecture |
| US3: 魔法值降为 0 | 全局替换 | 所有 Interface + Model |
| US4: validate_index.go 用 AllStatuses() | `internal/cmd/task/validate_index.go` 重构 | Interface 1: AllStatuses() |
| SC: 重导出兼容 | `pkg/feature/constants.go` | Interface 4: Re-export |
| SC: JSON 序列化不变 | type X string 保持兼容 | Testing: JSON 序列化兼容 |

## Migration Plan

按 package 分批，每个文件只修改一次。**每个 Phase 完成后必须通过验证关卡再进入下一 Phase**。

| Phase | Package | Files | Changes |
|-------|---------|-------|---------|
| 1 | `pkg/types/` | 3 新文件 | 定义 typed constants + helpers |
| **Gate** | | | `go build ./pkg/types/... && go test ./pkg/types/...` |
| 2 | `pkg/feature/` | constants.go | 移除原始定义，添加重导出 |
| **Gate** | | | `go build ./pkg/feature/... && go test ./pkg/feature/...` |
| 3 | `pkg/task/` | statemachine.go, add.go, state.go, deps.go, build.go, autogen.go, types.go, record.go, index.go, tasktemplate.go | Status + Priority + SurfaceType 魔法值替换 + 签名升级 |
| **Gate** | | | `go build ./pkg/task/... && go test ./pkg/task/...` |
| 4 | `pkg/forgeconfig/` | detect_surface.go, detect.go, execution_order.go | SurfaceType 魔法值替换 + 签名升级 |
| **Gate** | | | `go build ./pkg/forgeconfig/... && go test ./pkg/forgeconfig/...` |
| 5 | `internal/cmd/` | task/submit.go, task/claim.go, task/validate_index.go, task/add.go, task/reopen.go, task/transition.go, task/tree.go, task/migrate.go, quality_gate.go, cleanup.go, verify_task_done.go, feature/feature.go, feature/feature_complete.go, worktree/helpers.go, **surfaces.go** | Status + Priority + SurfaceType 魔法值替换 + 签名升级。**surfaces.go** 重导出 `KnownSurfaceTypes = forgeconfig.KnownSurfaceTypes`，map key 类型从 `string` 变为 `types.SurfaceType`，需同步更新 `runSurfacesTypes` 中的 lookup 逻辑（`KnownSurfaceTypes[typ]` 的 `typ` 需转换为 `types.SurfaceType`） |
| **Gate** | | | `go build ./internal/cmd/... && go test ./internal/cmd/...` |
| 6 | 全局验证 | 全局 | `go build ./... && go test -race -cover ./... && go vet ./...` |

## Open Questions

- [ ] `internal/cmd/task/migrate.go` 中使用了 `feature.StatusInProgress`（仅有的常量引用之一），迁移后是否直接改为 `types.StatusInProgress`？

## Appendix

### Alternatives Considered

| Approach | Pros | Cons | Why Not Chosen |
|----------|------|------|----------------|
| 按枚举类别分批（先 Status，后 Surface，再 Priority） | 每批改动范围清晰 | 同一文件被修改 3 次，merge 冲突多 | 按包分批更高效 |
| 大爆炸：一个 commit 全部替换 | 最快完成 | PR 过大难审查，回滚粒度粗 | 分 phase 提交更安全 |
| 使用 `stringer` 代码生成 | 自动生成 String() 方法 | 引入 build step 复杂度，过度工程 | 当前不需要 String() 方法 |

### References

- `docs/conventions/enum-constants.md` — 枚举常量组织规范
- `docs/business-rules/task-lifecycle.md` — Status 状态机定义
- `docs/proposals/unify-enum-constants/proposal.md` — 原始提案
