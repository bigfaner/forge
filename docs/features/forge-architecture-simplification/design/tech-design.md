---
created: 2026-05-18
updated: 2026-05-21
prd: prd/prd-spec.md
status: Draft
---

# Technical Design: Forge Architecture Simplification

## Overview

系统性重构 Forge CLI，通过 4 个渐进式阶段消除缺陷模式。核心策略：

1. **Phase 0**: Characterization tests 锁定当前行为
2. **Phase 1**: 纯重命名/删除，零行为变更
3. **Phase 2**: 行为修正——集中状态机、统一写入路径、完善错误处理
4. **Phase 3** (stretch goals): 结构优化——CLI UX、Config CRUD、包拆分。非阻塞里程碑，部分可推迟

所有变更限制在 `forge-cli/internal/` 和 `forge-cli/pkg/` 内，无新外部依赖。

## Architecture

### Current State Summary

| 组件 | 当前状态 | 目标状态 |
|------|---------|---------|
| `pkg/index/atomic.go` | **已有** `SaveIndexAtomic` (temp+rename) | 无需新建，可直接使用 |
| `pkg/index/lock.go` | **已有** `LockFile`/`UnlockFile` (5s advisory lock, Unix+Windows) | 新增 `WithLock` 回调封装 |
| `internal/cmd/errors.go` | **已有** 完整 AIError 体系 + `Exit()` | 新增 exit code 区分 (1 vs 2) + 新 error codes |
| `pkg/forgeconfig/config.go` | **已有** `GetConfigValue` (dot-notation) + `ReadAutoConfig` | 新增 `SetConfigValue` + `config set` CLI 命令 |
| `pkg/task/testgen.go` | 旧命名 | 重命名为 `autogen.go` |
| 状态验证 | 散布在 submit/claim/status 各命令 | 集中到 `pkg/task/statemachine.go` |
| `--force` flag | submit 和 status 都有 | 移除，用 `task reopen` 替代 |
| 写入安全 | submit 有锁+原子; claim/status/build/quality-gate 无锁无原子 | 统一 `WithLock` 模式 |

### Layer Placement

```
forge-cli/
├── cmd/forge/                    # Application entry point (仅 main.go, run.go)
├── internal/
│   └── cmd/                      # ALL CLI command implementations
│       ├── errors.go             # [扩展] AIError + Exit() + 新 error codes + exit code 区分
│       ├── output.go             # 输出格式化 (PrintBlock)
│       ├── submit.go             # [重构] 用 WithLock + 状态机，移除 --force
│       ├── claim.go              # [重构] 加锁 + 原子写入 + 状态机
│       ├── status.go             # [重构] 只读化，移除状态变更 + --force
│       ├── quality_gate.go       # [修正] addFixTask 加锁 + 真实 SourceTaskID + 活跃 cap
│       ├── reopen.go             # [新增] task reopen 子命令
│       ├── config.go             # [扩展] 新增 config set 子命令
│       ├── init.go               # huh TUI 配置初始化
│       ├── test_promote.go       # [加固] 路径遍历校验
│       ├── test_verify.go        # [加固] 解析失败处理 + unverifiable 标记
│       └── ...
├── pkg/
│   ├── task/
│   │   ├── types.go              # Task, TaskIndex, TaskState, RecordData
│   │   ├── index.go              # LoadIndex, SaveIndex (非原子 — 逐步淘汰)
│   │   ├── build.go              # BuildIndex — 改用原子保存
│   │   ├── add.go                # AddTask
│   │   ├── state.go              # SaveState — 改用原子写入
│   │   ├── testgen.go            # [重命名→autogen.go] AutoGenTaskDef
│   │   ├── statemachine.go       # [新增] ValidateTransition + CheckTransitionDeps
│   │   ├── preserve.go           # [新增] PreserveRuntimeFields
│   │   └── ...
│   ├── index/
│   │   ├── atomic.go             # [已有] SaveIndexAtomic
│   │   ├── lock.go               # [扩展] 新增 WithLock 回调
│   │   ├── lock_unix.go          # [已有] syscall.Flock
│   │   └── lock_windows.go       # [已有] LockFileEx
│   ├── feature/
│   │   ├── constants.go          # [扩展] 新增共享常量
│   │   ├── forge_state.go        # [重构] 改用原子写入
│   │   └── ...
│   ├── forgeconfig/
│   │   ├── config.go             # [扩展] 新增 SetConfigValue
│   │   └── ...
│   └── ...
```

依赖方向严格保持：`internal/cmd → pkg/*`，禁止反向。`pkg/index` 是叶子包（零内部依赖）。

### Component Diagram

```
                    ┌─────────────────────────────────────────┐
                    │         internal/cmd/ (Commands)        │
                    │  submit, claim, reopen, status, ...     │
                    └────────────┬────────────────────┬───────┘
                                 │                    │
                    ┌────────────▼────────┐  ┌───────▼────────┐
                    │  pkg/task/           │  │  AIError        │
                    │  ValidateTransition  │  │  Exit()         │
                    │  CheckTransitionDeps │  │  Factory funcs  │
                    └────────────┬────────┘  └────────────────┘
                                 │
              ┌──────────────────┼──────────────────┐
              │                  │                  │
    ┌─────────▼──────┐ ┌───────▼────────┐ ┌──────▼─────────┐
    │  pkg/index/     │ │ Task Operations│ │ BuildIndex     │
    │  WithLock +     │ │ Add/Query/     │ │ Orphan cleanup │
    │  SaveIndexAtomic│ │ Validate       │ │ Preserve       │
    └─────────────────┘ └────────────────┘ └────────────────┘
```

### Dependencies

| Dependency | Type | Status | Usage |
|------------|------|--------|-------|
| `spf13/cobra` | Direct | Existing | CLI framework, RunE migration |
| `gopkg.in/yaml.v3` | Direct | Existing | Config YAML parsing |
| `stretchr/testify` | Direct | Existing | Test assertions |
| `charmbracelet/huh` | Direct | Existing | Config init TUI (init.go 已使用) |
| Go 1.25 | — | Existing | Minimum Go version |
| **No new dependencies** | — | — | 所有变更使用 Go stdlib |

## Interfaces

### Interface 1: Centralized Transition Validation

Two-phase design: pure state validation (no dependencies needed) + optional dependency check. No `--force` override — terminal states are absolute. Dedicated `reopen` subcommand handles re-activation of rejected/skipped tasks.

**Design decisions:**
- `task status` is **read-only** — no state mutation via status command
- `--force` flag is **removed** — no escape hatch for terminal state protection
- `task reopen <id>` is the **only** way to leave rejected/skipped (not completed)
- Completed tasks are **irreversible** — if work needs re-doing, create a new subtask

```go
// pkg/task/statemachine.go (new)

type TransitionRole string

const (
    RoleSubmit TransitionRole = "submit"  // forge task submit
    RoleClaim  TransitionRole = "claim"   // forge task claim
    RoleReopen TransitionRole = "reopen"  // forge task reopen
    RoleAuto   TransitionRole = "auto"    // auto-downgrade, auto-unblock
)

// ValidateTransition validates a state transition (pure, no data lookup).
// Phase 1 of validation: checks terminal state protection and role-based rules.
func ValidateTransition(current, target string, role TransitionRole) error

// CheckTransitionDeps validates dependency satisfaction for blocked → pending/in_progress.
// Phase 2 of validation: call after ValidateTransition succeeds.
// Returns unmet dependency IDs, or nil if all deps are met.
func CheckTransitionDeps(index *TaskIndex, taskID string) ([]string, error)

type TransitionRule struct {
    From     string
    To       string
    Role     TransitionRole // "" = any role
    Allowed  bool
    GuardMsg string
}
```

**State transition table** (single authority):

| From | To | Role | Allowed | Guard |
|------|----|------|---------|-------|
| completed | * | * | **No** | "task already completed, create a subtask if re-work needed" |
| rejected | * | reopen | **Yes** (→ pending only) | — |
| rejected | * | other | **No** | "task rejected, use forge task reopen" |
| skipped | * | reopen | **Yes** (→ pending only) | — |
| skipped | * | other | **No** | "task skipped, use forge task reopen" |
| * | completed | submit | **Yes** | — |
| * | completed | non-submit | **No** | "use forge task submit" |
| in_progress | blocked | submit | **Yes** | sets BlockedReason |
| blocked | pending/in_progress | * | Dep check (phase 2) | `canAutoUnblock` |
| pending | blocked | * | **Yes** | block-source, dependency wait |
| * | * (same) | * | **Yes** | no-op |

**Caller pattern:**

```go
// In internal/cmd/submit.go:
if err := task.ValidateTransition(current, "completed", task.RoleSubmit); err != nil {
    return err
}

// In internal/cmd/reopen.go:
if err := task.ValidateTransition(current, "pending", task.RoleReopen); err != nil {
    return err
}

// In internal/cmd/claim.go (needs dep check):
if err := task.ValidateTransition(current, "pending", task.RoleClaim); err != nil {
    return err
}
unmet, err := task.CheckTransitionDeps(index, taskID)
```

### Interface 2: WithLock Callback + SaveStateAtomic

Atomic write primitives already exist in `pkg/index/atomic.go` and `pkg/index/lock.go`. This interface adds the `WithLock` callback pattern and extends atomic writes to state files.

```go
// pkg/index/lock.go (extended)

// WithLock acquires an advisory lock, calls fn, then releases the lock.
// Wraps the entire read-modify-write cycle: LoadIndex → mutation → SaveIndexAtomic.
// Returns ErrLockConflict if lock cannot be acquired within 5s.
// All index writers must use this — never call LockFile/UnlockFile directly.
func WithLock(indexPath string, fn func() error) error
```

```go
// pkg/feature/forge_state.go (modified)

// SaveStateAtomic atomically writes forge state via temp+rename.
// Replaces current os.WriteFile calls in WriteForgeState/EnsureForgeState/MarkFeatureCompleted.
func SaveStateAtomic(statePath string, data []byte) error
```

```go
// pkg/task/state.go (modified)

// SaveState is modified to use atomic write (temp+rename).
// Replaces current os.WriteFile call.
```

**Why `WithLock` callback instead of manual lock management**: The lock must wrap the entire read-modify-write cycle (LoadIndex → mutation → SaveIndexAtomic), not just the write. Manual lock management in submit.go shows the pattern works but is error-prone — claim.go and status.go omit it entirely. `WithLock` guarantees lock release even on panic.

**Standard pattern for all index writers:**

```go
err = index.WithLock(indexPath, func() error {
    idx, err := task.LoadIndex(indexPath)
    if err != nil { return err }
    // ... command-specific mutation ...
    return index.SaveIndexAtomic(indexPath, idx)
})
```

**Already-existing primitives (no changes needed):**
- `pkg/index/atomic.go` — `SaveIndexAtomic(path, data)` — temp file + `os.Rename`
- `pkg/index/lock.go` — `LockFile(indexPath)` / `UnlockFile(f)` — 5s timeout advisory lock
- `pkg/index/lock_unix.go` — `syscall.Flock` with `LOCK_EX|LOCK_NB`
- `pkg/index/lock_windows.go` — `LockFileEx`/`UnlockFileEx`

### Interface 3: Preserve Runtime Fields

```go
// pkg/task/preserve.go (new)

// PreserveRuntimeFields copies runtime-only fields from existing task to new task.
// Called during BuildIndex re-index to preserve state that isn't in .md frontmatter.
func PreserveRuntimeFields(existing, newTask *Task)
```

**Implementation uses explicit field assignment** — not reflection-based:

```go
func PreserveRuntimeFields(existing, newTask *Task) {
    if existing == nil {
        return
    }
    newTask.Status = existing.Status
    newTask.SourceTaskID = existing.SourceTaskID
    newTask.BlockedReason = existing.BlockedReason
}
```

### Interface 4: Config Set/Get

```go
// pkg/forgeconfig/config.go (extended)

// SetConfigValue sets a single config key using dot-notation.
// Returns ErrInvalidInput if key is not recognized.
// Supports: auto.runTasks.quick, auto.runTasks.full, auto.gitPush, etc.
func SetConfigValue(projectRoot, key, value string) error
```

**Existing capabilities (no changes needed):**
- `GetConfigValue(projectRoot, key)` — dot-notation reader, already supports all auto keys
- `ReadAutoConfig(projectRoot)` — reads full AutoConfig struct
- `WriteConfig(projectRoot, config)` — writes full Config struct
- `config get` CLI command — already functional

**New additions:**
- `SetConfigValue` — single-key setter using dot-notation
- `forge config set <key> <value>` CLI command

**Supported keys** (from `AutoConfig` struct): `auto.e2eTest.quick`, `auto.e2eTest.full`, `auto.consolidateSpecs.quick`, `auto.consolidateSpecs.full`, `auto.cleanCode.quick`, `auto.cleanCode.full`, `auto.validation.quick`, `auto.validation.full`, `auto.runTasks.quick`, `auto.runTasks.full`, `auto.knowledgeSave.quick`, `auto.knowledgeSave.full`, `auto.gitPush`, `worktree.source-branch`, `worktree.copy-files`.

## Data Models

*Single-layer CLI feature. No database. db-schema: "no".*

### Model 1: TransitionRule

```go
TransitionRule = {
    From:     string         // current status
    To:       string         // target status
    Role:     TransitionRole // "" = any role
    Allowed:  bool           // whether transition is permitted
    GuardMsg: string         // human-readable reason when blocked
}
```

### Model 2: TransitionRole

```go
TransitionRole = string  // enum: "submit" | "claim" | "reopen" | "auto"
```

## Error Handling

### Current State

`internal/cmd/errors.go` 已有成熟的 AIError 体系：
- `ErrorCode` 类型 + 6 个错误码常量 (`NO_PROJECT`, `NO_FEATURE`, `INVALID_INPUT`, `NOT_FOUND`, `CONFLICT`, `VALIDATION_ERROR`)
- `AIError` 结构体：`Code`, `Message`, `Cause`, `Hint`, `Action` 字段
- `Exit(err)` 函数：打印结构化输出 + `os.Exit(1)`
- 17+ 工厂函数

**缺失部分：**
- `Exit()` 不区分 exit code — 始终 `os.Exit(1)`
- 缺少 `ExitCode()` 方法
- 缺少 Phase 2 需要的新 error codes

### Exit Code Strategy

Forge commands operate in **two execution contexts**:

**Context 1: Hook execution** (only `verify-task-done`). Per Claude Code hooks convention:
- Exit 0 → success, action proceeds
- Exit 2 → blocking error, tool call is prevented, stderr shown to Claude
- Exit 1 → non-blocking error, action still proceeds

**Context 2: Bash tool execution** (all other commands). Claude Code treats both exit 1 and 2 as "command failed". The exit code distinction serves the *agent reading stderr* — code 2 signals "change your approach", code 1 signals "retry".

| Exit Code | Semantics | Forge Usage | Agent Guidance |
|-----------|-----------|-------------|----------------|
| **0** | Success | Command completed normally | Proceed |
| **2** | Policy violation | Invalid transition, path traversal, unverifiable contract | Change approach |
| **1** | Retryable failure | Lock timeout, parse failure, not found | Retry or fix precondition |

**Exit code assignment by error type:**

| ErrorCode | Exit Code | Rationale |
|-----------|-----------|-----------|
| `ErrInvalidTransition` | 2 | Agent must change approach |
| `ErrInvalidPath` | 2 | Security policy violation |
| `ErrContractUnverifiable` | 2 | Structural problem |
| `ErrLockConflict` | 1 | Transient — retryable |
| `ErrEvalParseFailure` | 1 | Retryable with different input |
| `ErrFeatureNotSet` (existing) | 1 | Agent can configure and retry |
| All other existing AIError codes | 1 | Default: retryable |

### Changes to errors.go

**New error codes and factory functions:**

```go
const (
    // ... existing codes ...
    ErrInvalidTransition ErrorCode = "INVALID_TRANSITION"
    ErrInvalidPath       ErrorCode = "INVALID_PATH"
    ErrEvalParseFailure  ErrorCode = "EVAL_PARSE_FAILURE"
    ErrContractUnverif   ErrorCode = "CONTRACT_UNVERIFIABLE"
)

func NewErrInvalidTransition(from, to string, hint string) *AIError
func NewErrEvalParseFailure(raw string) *AIError
func NewErrInvalidPath(input string) *AIError
func NewErrContractUnverifiable(contractPath string) *AIError
```

**Exit code differentiation:**

```go
func (e *AIError) ExitCode() int {
    switch e.Code {
    case ErrInvalidInput, ErrConflict, ErrValidation,
         ErrInvalidTransition, ErrInvalidPath, ErrContractUnverif:
        return 2
    default:
        return 1
    }
}

func Exit(err error) {
    if aiErr, ok := err.(*AIError); ok {
        printAIError(aiErr)
        os.Exit(aiErr.ExitCode())
    }
    fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
    os.Exit(1)
}
```

### Propagation Strategy

**Current**: Mixed — submit uses manual lock + `SaveIndexAtomic`; claim/status use non-atomic `SaveIndex`; `--force` bypasses validation.

**Target**: All commands use `RunE` + `WithLock` + state machine validation. `Exit()` differentiates exit codes.

**Per-file migration:**

| File | Current | Target |
|------|---------|--------|
| `internal/cmd/submit.go` | Manual `LockFile`/`UnlockFile`, `SaveIndexAtomic`, `--force` | `WithLock` + state机 + 移除 `--force` |
| `internal/cmd/claim.go` | No lock, `task.SaveIndex` (non-atomic) | `WithLock` + `SaveIndexAtomic` + 状态机 |
| `internal/cmd/status.go` | Update mode: no lock, `task.SaveIndex`, `--force` | Read-only only, 移除 update mode + `--force` |
| `internal/cmd/reopen.go` | (new) | `ValidateTransition(current, "pending", RoleReopen)` |
| `internal/cmd/quality_gate.go` | `addFixTask`: no lock, no atomic | `WithLock` + `SaveIndexAtomic` |
| `internal/cmd/test_promote.go` | No path traversal check | `filepath.Base()` + reject `..` |
| `internal/cmd/test_verify.go` | Silent zero-value on parse failure, OK on no Fact Table | AIError on parse failure, unverifiable on no Fact Table |
| `pkg/task/build.go` | `task.SaveIndex` (non-atomic) | `SaveIndexAtomic` |
| `pkg/task/state.go` | `os.WriteFile` (non-atomic) | temp+rename atomic write |
| `pkg/feature/forge_state.go` | `os.WriteFile` (non-atomic) for all writes | temp+rename atomic write |

## Cross-Layer Data Map

Single-layer feature (CLI only). Not applicable.

## Integration Specs

No existing-page integrations — not applicable (CLI-only feature).

## Testing Strategy

### Per-Layer Test Plan

| Layer | Test Type | Tool | What to Test | Coverage Target |
|-------|-----------|------|--------------|-----------------|
| pkg/task | Unit | testing + testify | StateMachine transitions (all combinations) | 100% transition matrix |
| pkg/task | Unit | testing + testify | PreserveRuntimeFields | All preserved fields |
| pkg/index | Unit | testing + testify | WithLock concurrent access | Lock conflict detection |
| internal/cmd | Characterization | testing + testify | Current behavior locked (Phase 0) | All SM/QG paths |
| internal/cmd | Integration | forge test | End-to-end command behavior | All existing tests |

### Key Test Scenarios

**Characterization Tests (Phase 0):**

1. `TestSubmit_CurrentBehavior_AllowsCompletedResubmit` — completed task submit currently succeeds
2. `TestAdd_BlockSource_CurrentBehavior_AllowsCompletedToBlocked` — block-source on completed succeeds
3. `TestClaim_AutoUnblock_CurrentBehavior` — auto-unblock behavior
4. `TestQualityGate_SourceTaskID_IsSentinel` — sentinel ID behavior
5. `TestQualityGate_CountFixTasks_CountsAll` — counts completed fix-tasks
6. `TestBuildIndex_Orphan_WarningOnly` — orphan only warned
7. `TestStatus_AllowsMutation` — status command allows state changes
8. `TestSubmit_AutoDowngrade_NoBlockedReason` — no BlockedReason set

**State Machine Tests (Phase 2):**

1. Terminal state protection
2. Submit-only path to completed
3. Blocked unblock with dependency check
4. Completed irreversible
5. Reopen only rejected/skipped
6. Role isolation

**WithLock Tests (Phase 2):**

1. Concurrent access — two goroutines
2. Lock conflict returns ErrLockConflict
3. Lock released on panic

**Exit Code Tests (Phase 2):**

1. ErrInvalidTransition → exit 2
2. ErrLockConflict → exit 1
3. Default → exit 1

### Overall Coverage Target

- New code (statemachine.go, preserve.go): 90%+
- Modified code: maintain or improve existing coverage
- Characterization tests: cover all SM/QG/GI scenarios

## Security Considerations

### Threat Model

| Threat | Vector | Impact |
|--------|--------|--------|
| Path traversal | `forge test promote <journey>` with `../` | File system access outside tests/ |
| Concurrent data corruption | Multiple agents writing index.json | Data loss, invalid JSON |
| Eval document loss | Reviser modifies files in-place, no backup | Loss of original documents |
| Stale lock after crash | Process crashes without unlock | Writers blocked until timeout (5s) |

### Mitigations

| Threat | Mitigation | Implementation |
|--------|-----------|----------------|
| Path traversal | Input validation | `filepath.Base()` + reject `..` in test_promote.go |
| Concurrent corruption | Advisory file locking | `WithLock` with 5s timeout |
| Eval document loss | Backup before reviser | `cp -r DOC_DIR DOC_DIR.bak` + restore on failure |
| Stale lock | POSIX flock auto-release | Per-fd lock, released on close/exit |

## PRD Coverage Map

| PRD Requirement | Design Component | Interface |
|-----------------|------------------|-----------|
| **DR-1** Index atomicity | SaveIndexAtomic (已有) + WithLock | Interface 2 |
| **DR-2** Index locking | WithLock callback | Interface 2 |
| **DR-3** State atomicity | SaveStateAtomic | Interface 2 |
| **DR-4** State consistency | Write false instead of delete | Interface 2 |
| **BC-1** Submit state check | ValidateTransition | Interface 1 |
| **BC-2** Block-source terminal guard | ValidateTransition | Interface 1 |
| **BC-3** Dependency check consistency | CheckTransitionDeps | Interface 1 |
| **BC-4** Auto-downgrade BlockedReason | ValidateTransition | Interface 1 |
| **BC-5** QG real SourceTaskID | addFixTask fix | Task 2.8 |
| **BC-6** Active-only fix-task cap | countFixTasks fix | Task 2.8 |
| **BC-7** No-feature error | AIError (已有) | Error Handling |
| **BC-8** Orphan cleanup | BuildIndex fix | Task 2.7 |
| **BC-9** Fix-task orphan exempt | isAutoGenTaskID | Task 2.7 |
| **BC-10** Preserve extensibility | PreserveRuntimeFields | Interface 3 |
| **BC-11** Reopen command | reopen.go | Task 2.5 |
| **EC-1~5** AIError unification | Exit code differentiation | Error Handling |
| **EC-6** Path traversal | ErrInvalidPath | Task 2.10 |
| **EC-7** Parse failure | ErrEvalParseFailure | Task 2.10 |
| **EC-8** Unverifiable marking | ErrContractUnverifiable | Task 2.10 |
| **ES-1~4** Eval safety | Backup + context injection | Task 2.9 |
| **CE-1** Config set | SetConfigValue | Interface 4 |
| **CE-2** Config get coverage | GetConfigValue (已有) | Interface 4 |
| **CE-3** Config init complete | Merge to huh TUI | Task 3.2 |
| **CE-4** Schema version | Version field | Task 3.4 |
| **CE-5** Schema enum alignment | Schema + CLI sync | Task 3.4 |
| **CC-1** RunE unification | Run → RunE migration | Task 3.1 |
| **CC-2** Config init merge | Delete bufio path | Task 3.2 |
| **CC-3** Args validation | cobra args validators | Task 3.1 |
| **CC-4** PrintBlock | Output standardization | Task 3.1 |
| **CH-1** Naming accuracy | testgen→autogen rename | Task 1.2 |
| **CH-2** Dead code deletion | Remove unused exports | Task 1.1 |
| **CH-3** Magic value extraction | Constants in pkg/feature/ | Task 1.3 |
| **CH-4** Code dedup | IsBusinessTask to pkg/task | Task 1.3 |

## Open Questions

- [ ] Windows lock: `LockFileEx` on NTFS provides same semantics as `flock` on POSIX? (Go/No-Go — cross-platform lock 已实现，需验证语义一致性)
- [ ] `--force` removal: submit 和 status 都有 `--force`。移除后 `task reopen` 是唯一重新激活机制。确保下游 skill 无 `--force` 依赖。
- [ ] Config init 合并: `forge init` 用 huh TUI，`forge config init` 用 bufio。合并后是否保留 `forge config init` 作为 `forge init --config-only` 的别名？

## Appendix

### Alternatives Considered

| Approach | Why Not Chosen |
|----------|---------------|
| Go FSM library (looplab/state) | Table-driven `ValidateTransition` is simpler, matches existing codebase |
| Event sourcing for state | Massive overengineering for CLI task manager |
| Cobra middleware pattern | RunE + Exit() achieves same result with less disruption |
| New `pkg/constants/` package | Import cycle risk; constants fit in existing `pkg/feature/constants.go` |

### References

- `pkg/index/atomic.go` — existing atomic write primitive
- `pkg/index/lock.go` — existing file locking primitive
- `internal/cmd/errors.go` — existing AIError system
- `pkg/forgeconfig/config.go` — existing config read/write
- `design/state-transition-diagram.md` — 状态流转图
