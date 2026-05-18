---
created: 2026-05-18
prd: prd/prd-spec.md
status: Draft
---

# Technical Design: Forge Architecture Simplification

## Overview

系统性重构 Forge CLI，通过 4 个渐进式阶段消除 19 个缺陷模式。核心策略：

1. **Phase 0**: Characterization tests 锁定当前行为
2. **Phase 1**: 纯重命名/删除，零行为变更
3. **Phase 2**: 行为修正——集中状态机、统一写入路径、完善错误处理
4. **Phase 3** (stretch goals): 结构优化——CLI UX、Config CRUD、包拆分。非阻塞里程碑，部分可推迟

所有变更限制在 `forge-cli/internal/` 和 `forge-cli/pkg/` 内，无新外部依赖。

## Architecture

### Layer Placement

```
forge-cli/
├── cmd/                          # 应用层：cobra 命令处理（输入/输出）
│   ├── errors.go                 # AIError 工厂函数（扩展）
│   ├── output.go                 # 输出格式化
│   ├── reopen.go                 # [新增] task reopen 子命令
│   └── {command}.go              # 各命令入口
├── internal/                     # (当前未使用，Phase 3 可能引入)
└── pkg/                          # 领域层：业务逻辑
    ├── task/
    │   ├── statemachine.go       # [新增] 状态机验证
    │   ├── build_index.go        # [重命名] BuildIndex
    │   ├── autogen.go            # [重命名] 自动生成任务
    │   ├── preserve.go           # [新增] PreserveRuntimeFields
    │   └── ...
    ├── index/
    │   ├── atomic.go             # [扩展] AtomicWrite 原语
    │   └── lock.go               # [扩展] WithLock 回调 + advisory lock
    ├── feature/
    │   ├── forge_state.go        # [修改] SaveStateAtomic
    │   └── constants.go          # [删除] 常量移到 pkg/constants/
    ├── constants/                 # [新增] 统一常量 (LockTimeoutSeconds, 其他魔法值)
    │   └── forge.go
    ├── e2eprobe/                  # [修改] typed YAML unmarshal
    └── profile/
        └── config.go             # [修改] config set/get
```

依赖方向严格保持：`cmd -> pkg`，禁止反向。

### Component Diagram

```
                    ┌─────────────────────────────────────────┐
                    │             Cobra Commands              │
                    │  (submit, claim, reopen, add, ...)      │
                    └────────────┬────────────────────┬───────┘
                                 │                    │
                    ┌────────────▼────────┐  ┌───────▼────────┐
                    │  Transition Validation│ │  AIError        │
                    │  ValidateTransition  │  │  Exit()         │
                    │  CheckTransitionDeps │  │  Factory funcs  │
                    └────────────┬────────┘  └────────────────┘
                                 │
              ┌──────────────────┼──────────────────┐
              │                  │                  │
    ┌─────────▼──────┐ ┌───────▼────────┐ ┌──────▼─────────┐
    │  WithLock +     │ │ Task Operations│ │ BuildIndex     │
    │  AtomicWrite    │ │ Add/Query/     │ │ Orphan cleanup │
    │                 │ │ Validate       │ │ Preserve       │
    └─────────────────┘ └────────────────┘ └────────────────┘
```

### Dependencies

| Dependency | Type | Status | Usage |
|------------|------|--------|-------|
| `spf13/cobra` | Direct | Existing | CLI framework, RunE migration |
| `gopkg.in/yaml.v3` | Direct | Existing | Config YAML parsing (replace hand-written parser) |
| `stretchr/testify` | Direct | Existing | Test assertions |
| `charmbracelet/huh` | Indirect | Existing | Config init TUI (merge two paths) |
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
// pkg/task/statemachine.go

// TransitionRole identifies which command is requesting the transition.
// Each role gets explicit rules in the transition table.
type TransitionRole string

const (
    RoleSubmit TransitionRole = "submit"  // forge task submit
    RoleClaim  TransitionRole = "claim"   // forge task claim
    RoleReopen TransitionRole = "reopen"  // forge task reopen
    RoleAuto   TransitionRole = "auto"    // auto-downgrade, auto-unblock
)

// ValidateTransition validates a state transition (pure, no data lookup).
// Phase 1 of validation: checks terminal state protection and role-based rules.
// All state-changing commands must call this as their single validation entry point.
func ValidateTransition(current, target string, role TransitionRole) error

// CheckTransitionDeps validates dependency satisfaction for blocked → pending/in_progress.
// Phase 2 of validation: call after ValidateTransition succeeds, only for unblock transitions.
// Returns unmet dependency IDs, or nil if all deps are met.
func CheckTransitionDeps(index *TaskIndex, taskID string) ([]string, error)

// TransitionRule defines a single allowed/disallowed transition
type TransitionRule struct {
    From       string
    To         string
    Role       TransitionRole // "" = any role
    Allowed    bool
    GuardMsg   string  // Message shown when transition is blocked
}
```

**`canAutoUnblock`** is an unexported helper inside `statemachine.go`, called by `CheckTransitionDeps`. Not a standalone interface — it's internal guard logic for the `blocked → pending/in_progress` transition, not a separate contract.

**State transition table** (the single authority):

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

**Key simplification vs. original design:**
- `Forceable` column removed — no override mechanism
- `RoleStatus` removed — `task status` is read-only, no state mutations
- `RoleReopen` added — explicit intent to re-activate, only rejected/skipped → pending
- `force bool` parameter removed from `ValidateTransition` signature

**Caller pattern:**

```go
// In submit.go:
if err := task.ValidateTransition(current, "completed", task.RoleSubmit); err != nil {
    return err
}

// In reopen.go:
if err := task.ValidateTransition(current, "pending", task.RoleReopen); err != nil {
    return err
}

// In claim.go (needs dep check):
if err := task.ValidateTransition(current, "pending", task.RoleClaim); err != nil {
    return err
}
unmet, err := task.CheckTransitionDeps(index, taskID)
// handle unmet deps...
```

### Interface 2: Atomic Write

Single `AtomicWrite` primitive for all file writes. Lock management via `WithLock` callback (absorbed from [encapsulated-index-lock proposal](../../../../proposals/encapsulated-index-lock/proposal.md)).

```go
// pkg/index/lock.go (extended)

// WithLock acquires an advisory lock, calls fn, then releases the lock.
// Wraps the entire read-modify-write cycle: LoadIndex → mutation → SaveIndexAtomic.
// Returns ErrLockConflict if lock cannot be acquired within LockTimeoutSeconds (5s).
// All index writers must use this — never call LockFile/UnlockFile directly.
func WithLock(indexPath string, fn func() error) error
```

```go
// pkg/index/atomic.go (extended)

// AtomicWrite writes data to path atomically (temp + rename).
// The single primitive for all file writes that need crash safety.
func AtomicWrite(path string, data []byte, perm os.FileMode) error

// SaveIndexAtomic atomically saves index (temp + rename, no lock).
// Must be called inside a WithLock callback for concurrent safety.
func SaveIndexAtomic(indexPath string, index *task.TaskIndex) error

// SaveStateAtomic atomically writes forge state.
// Composes AtomicWrite + state marshaling.
func SaveStateAtomic(statePath string, state []byte) error
```

**Why `WithLock` callback instead of `SaveIndexLocked`**: The lock must wrap the entire read-modify-write cycle (LoadIndex → mutation → SaveIndexAtomic), not just the write. `SaveIndexLocked` only protects the write step, leaving a TOCTOU gap between load and lock. `WithLock` keeps `pkg/index` as a leaf package (zero internal deps) because LoadIndex/SaveIndexAtomic stay in the caller's callback.

**Standard pattern for all index writers:**

```go
err = index.WithLock(indexPath, func() error {
    idx, err := task.LoadIndex(indexPath)
    if err != nil { return err }
    // ... command-specific mutation ...
    return index.SaveIndexAtomic(indexPath, idx)
})
```

### Interface 3: Preserve Runtime Fields

```go
// pkg/task/preserve.go (new)

// PreserveRuntimeFields copies runtime-only fields from existing task to new task.
// Called during BuildIndex re-index to preserve state that isn't in .md frontmatter.
// Implementation uses explicit switch statement (not reflection) for compile-time safety.
func PreserveRuntimeFields(existing, newTask *Task)
```

**Implementation uses explicit field assignment** — not reflection-based name matching, not a string slice. This ensures compile-time safety when fields are renamed. To add a new preserved field, add one `case` line:

```go
func PreserveRuntimeFields(existing, newTask *Task) {
    if existing == nil {
        return
    }
    // Preserve fields that exist only in runtime state, not in .md frontmatter.
    newTask.Status = existing.Status
    newTask.SourceTaskID = existing.SourceTaskID
    newTask.BlockedReason = existing.BlockedReason
}
```

### Interface 4: Config Set/Get

```go
// pkg/profile/config.go (extended)

// SetAutoKey sets a single auto configuration key.
// Returns ErrInvalidInput if key is not a recognized auto field.
func SetAutoKey(projectRoot, key, value string) error

// GetAutoKeyValue returns the value of any auto configuration key.
// Returns ("", false, nil) if key exists but is not set.
// Returns ("", false, ErrInvalidInput) if key is not a recognized auto field.
func GetAutoKeyValue(projectRoot, key string) (string, bool, error)
```

**Valid auto keys** (from `AutoConfig` struct): `e2eTest`, `consolidateSpecs`, `cleanCode`, `gitPush`. Any other key returns `ErrInvalidInput`.

## Data Models

*Single-layer CLI feature. No database. db-schema: "no".*

### Model 1: TransitionRule

```go
TransitionRule = {
    From:       string         // current status
    To:         string         // target status
    Role:       TransitionRole // "" = any role
    Allowed:    bool           // whether transition is permitted
    GuardMsg:   string         // human-readable reason when blocked
}
```

### Model 2: TransitionRole

```go
TransitionRole = string  // enum: "submit" | "claim" | "reopen" | "auto"
```

### Model 3: PreserveConfig

No separate model — `PreserveRuntimeFields` uses explicit field assignment. No reflection, no string-driven field lookup. Adding a preserved field = adding one line of code.

## Error Handling

### Exit Code Strategy

Forge commands operate in **two execution contexts** with different exit code semantics:

**Context 1: Hook execution** (only `verify-task-done`). Per [Claude Code hooks convention](../../official-references/hooks.md):
- Exit 0 → success, action proceeds
- Exit 2 → blocking error, tool call is prevented, stderr shown to Claude
- Exit 1 → non-blocking error, action still proceeds

**Context 2: Bash tool execution** (all other commands). Claude Code treats both exit 1 and 2 as "command failed" — there is no behavioral difference. The agent reads stderr content (AIError format with Code/Message/Cause/Hint/Action) to decide next steps.

**Why differentiate exit codes anyway?** The exit code distinction serves the *agent reading stderr*, not Claude Code's runtime. An AIError with exit code 2 signals "change your approach" in the error message semantics, while exit code 1 signals "retry the same approach". This is useful even without runtime behavioral differences because the structured AIError output already tells the agent what to do.

| Exit Code | Semantics | Forge Usage | Agent Guidance |
|-----------|-----------|-------------|----------------|
| **0** | Success | Command completed normally | Proceed |
| **2** | Policy violation | Invalid transition (terminal state, wrong role), path traversal, unverifiable contract | Change approach, do not retry same action |
| **1** | Retryable failure | Lock timeout, parse failure, no feature, not found | Retry same action or fix precondition |

**Exit code assignment by error type:**

| ErrorCode | Exit Code | Rationale |
|-----------|-----------|-----------|
| `ErrInvalidTransition` | 2 | Agent's action is fundamentally wrong — must change approach |
| `ErrInvalidPath` | 2 | Security policy violation — agent must not retry |
| `ErrContractUnverifiable` | 2 | Structural problem — agent needs human intervention |
| `ErrLockConflict` | 1 | Transient — agent can retry after brief wait |
| `ErrEvalParseFailure` | 1 | Infrastructure issue — retryable with different input |
| `ErrFeatureNotSet` (existing) | 1 | Configuration issue — agent can run `forge feature <slug>` then retry |
| All other existing AIError codes | 1 | Default: retryable |

**Reference**: `verify_task_done.go` already uses exit code 2 as a hook (blocking git commit). This is the only command that runs as a hook; all others run via Bash tool where both codes 1 and 2 are treated as failure.

### Error Types & Codes

| Error Code | Name | Description | Usage |
|------------|------|-------------|-------|
| `ERR_INVALID_TRANSITION` | `ErrInvalidTransition` | State transition not allowed by state machine | SM-1, SM-3 |
| `ERR_LOCK_CONFLICT` | `ErrLockConflict` | File lock acquisition timed out | MA-1 |
| `ERR_NO_FEATURE` | `ErrFeatureNotSet` (existing) | No feature configured for quality gate | QG-3 |
| `ERR_PARSE_FAILURE` | `ErrEvalParseFailure` | Eval scorer output could not be parsed | EP-1 |
| `ERR_INVALID_PATH` | `ErrInvalidPath` | Path traversal detected in input | TM-1 |
| `ERR_UNVERIFIABLE` | `ErrContractUnverifiable` | Contract has no Fact Table entry | TM-3 |

New factory functions to add:

```go
func ErrInvalidTransition(from, to string, hint string) *AIError
func ErrEvalParseFailure(raw string) *AIError
func ErrInvalidPath(input string) *AIError
func ErrContractUnverifiable(contractPath string) *AIError
```

### Propagation Strategy

**Current**: Mixed — some commands use `Run` + `os.Exit()`, others use `RunE` + `error` return. AIError used inconsistently. All errors exit with code 1. `--force` flag bypasses validation.

**Target**: All commands use `RunE` + return `error`. `Exit()` function is the single error output point with differentiated exit codes. No `--force` flag — `task reopen` is the explicit mechanism for re-activating rejected/skipped tasks.

```go
// Exit prints the AI-friendly error and exits with the appropriate code.
// AIError.ExitCode() returns 2 for blocking errors, 1 for soft failures.
func Exit(err error) {
    if aiErr, ok := err.(*AIError); ok {
        printAIError(aiErr)
        os.Exit(aiErr.ExitCode())
    }
    fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
    os.Exit(1)
}

// ExitCode returns 2 for blocking errors, 1 for soft failures.
func (e *AIError) ExitCode() int {
    switch e.Code {
    case ErrInvalidInput, ErrConflict, ErrValidation:
        return 2 // blocking: agent must change approach
    default:
        return 1 // soft failure: agent can retry
    }
}
```

**Flow:**

```
Command (RunE) → return error OR return nil
    ↓
rootCmd.Execute() catches error
    ↓
Exit(err) → AIError.ExitCode() → prints to stderr + os.Exit(1|2)
```

**Per-file migration:**

| File | Current | Target |
|------|---------|--------|
| `worktree.go` | `fmt.Errorf` everywhere | AIError factory functions + `Exit()` |
| `submit.go` | `fmt.Fprintln(os.Stderr) + os.Exit(1)`, `--force` bypass | AIError + `Exit()` (lock conflict → code 1, invalid transition → code 2). Remove `--force` flag. |
| `status.go` | `isTransitionAllowed` + state mutation | Read-only: show task status only. Remove state mutation logic. |
| `reopen.go` | (new) | `ValidateTransition(current, "pending", RoleReopen)`. Only rejected/skipped → pending. |
| `quality_gate.go` | Returns `nil` on infrastructure failure | Returns error (soft failure → code 1) |
| `test_promote.go` | `os.Exit(1)` on path traversal | AIError + `Exit()` (path traversal → code 2) |
| `test_verify.go` | `os.Exit(1)` on parse/verify failure | AIError + `Exit()` (parse failure → code 1, unverifiable → code 2) |
| `verify_task_done.go` | Already uses exit code 2 | Keep as-is (pattern reference) |

## Cross-Layer Data Map

Single-layer feature (CLI only). Cross-Layer Data Map not applicable.

## Integration Specs

No existing-page integrations — not applicable (CLI-only feature).

## Testing Strategy

### Per-Layer Test Plan

| Layer | Test Type | Tool | What to Test | Coverage Target |
|-------|-----------|------|--------------|-----------------|
| pkg/task | Unit | testing + testify | StateMachine transitions (all 6×6 combinations) | 100% transition matrix |
| pkg/task | Unit | testing + testify | PreserveRuntimeFields | All DefaultPreservedFields |
| pkg/index | Unit | testing + testify | WithLock concurrent access | Lock conflict detection |
| pkg/index | Unit | testing + testify | SaveStateAtomic crash safety | Temp file cleanup on error |
| cmd/* | Characterization | testing + testify | Current behavior locked (Phase 0) | All SM/QG paths |
| cmd/* | Integration | forge test | End-to-end command behavior | 126+ existing tests |
| Plugin | Skill eval | Manual | eval/SKILL.md backup+rollback | Manual verification |

### Key Test Scenarios

**Characterization Tests (Phase 0) — must pass before Phase 2:**

1. `TestSubmit_CurrentBehavior_AllowsCompletedResubmit` — submit on completed task currently succeeds
2. `TestAdd_BlockSource_CurrentBehavior_AllowsCompletedToBlocked` — block-source on completed task currently succeeds
3. `TestClaim_AutoUnblock_CurrentBehavior` — auto-unblock with/without BlockedReason
4. `TestStatus_BlockedToCompleted_CurrentlyBlocked` — already fixed in code, verify
5. `TestQualityGate_SourceTaskID_IsSentinel` — verify current sentinel behavior
6. `TestQualityGate_CountFixTasks_CountsAll` — verify counts completed fix-tasks
7. `TestBuildIndex_Orphan_WarningOnly` — verify orphan is only warned, not cleaned

**State Machine Tests (Phase 2):**

1. `TestValidateTransition_TerminalStateProtection` — completed/rejected/skipped cannot leave (except rejected/skipped via RoleReopen)
2. `TestValidateTransition_SubmitOnlyPath` — only RoleSubmit can set completed
3. `TestValidateTransition_BlockedUnblock` — pure validation passes, deps checked separately
4. `TestValidateTransition_CompletedIrreversible` — completed → * always blocked, no override
5. `TestValidateTransition_ReopenOnlyRejectedSkipped` — RoleReopen allows rejected/skipped → pending, not completed → pending
6. `TestCheckTransitionDeps_FixTaskAware` — blocked while fix-task is active
7. `TestValidateTransition_RoleIsolation` — RoleClaim/RoleReopen/RoleAuto each have correct rules

**Atomic Write Tests (Phase 2):**

1. `TestWithLock_ConcurrentAccess` — two goroutines, only one succeeds first
2. `TestSaveIndexAtomic_CrashSafety` — verify temp file cleanup on error
3. `TestSaveStateAtomic_Atomic` — verify temp+rename behavior
4. `TestAtomicWrite_Primitive` — verify AtomicWrite with various data sizes

**Exit Code Tests (Phase 2):**

1. `TestExitCode_InvalidTransition_Returns2` — ErrInvalidTransition → exit 2
2. `TestExitCode_LockConflict_Returns1` — ErrLockConflict → exit 1
3. `TestExitCode_InvalidPath_Returns2` — ErrInvalidPath → exit 2
4. `TestExitCode_Default_Returns1` — unknown ErrorCode → exit 1

### Overall Coverage Target

- New code (statemachine.go, preserve.go, constants/): 90%+
- Modified code (claim.go, submit.go, etc.): maintain or improve existing coverage
- Characterization tests: cover all SM-1~SM-8, QG-1~QG-3, GI-1 scenarios

## Security Considerations

### Threat Model

| Threat | Vector | Impact |
|--------|--------|--------|
| Path traversal | `forge test promote <journey>` with `../` in journeyName | File system access outside tests/ |
| Concurrent data corruption | Multiple agents writing index.json simultaneously | Data loss, invalid JSON |
| Eval document loss | Reviser modifies files in-place, no backup on failure | Loss of original documents |
| Stale lock after crash | Process holding advisory lock crashes without unlock | All subsequent writers blocked until lock timeout (5s) |

### Mitigations

| Threat | Mitigation | Implementation |
|--------|-----------|----------------|
| Path traversal | Input validation | `filepath.Base()` + reject `..` components in test_promote.go |
| Concurrent corruption | Advisory file locking | `WithLock` with 5s timeout (from `pkg/constants/`) |
| Eval document loss | Step 1 backup before reviser runs | `cp -r DOC_DIR DOC_DIR.bak` + restore on failure |
| Stale lock | POSIX flock auto-release on process death | Verify on Linux + macOS; advisory locks are per-fd, released on close/exit |

## PRD Coverage Map

| PRD Requirement | Design Component | Interface / Model |
|-----------------|------------------|-------------------|
| **DR-1** Index atomicity | WithLock + AtomicWrite | Interface 2 |
| **DR-2** Index locking | WithLock callback | Interface 2 |
| **DR-3** State atomicity | SaveStateAtomic | Interface 2 |
| **DR-4** State consistency | Write false instead of delete | Interface 2 (SaveStateAtomic) |
| **BC-1** Submit state check | ValidateTransition | Interface 1 |
| **BC-2** Block-source terminal guard | ValidateTransition | Interface 1 |
| **BC-3** Dependency check consistency | ValidateTransition + CheckTransitionDeps | Interface 1 |
| **BC-4** Auto-downgrade BlockedReason | ValidateTransition (sets BlockedReason) | Interface 1 |
| **BC-5** QG real SourceTaskID | addFixTask uses real ID | W4 implementation |
| **BC-6** Active-only cap | countFixTasks status filter | W4 implementation |
| **BC-7** No-feature error | AIError ErrFeatureNotSet | Error Handling |
| **BC-8** Orphan cleanup | BuildIndex default cleanup | W6 implementation |
| **BC-9** Fix-task orphan exempt | isAutoGenTaskID extended | W6 implementation |
| **BC-10** Preserve extensibility | PreserveRuntimeFields | Interface 3 |
| **EC-1~5** AIError unification | Factory functions + Exit() with exit codes 1/2 | Error Handling |
| **EC-6** Path traversal | ErrInvalidPath (exit 2) | Error Handling |
| **EC-7** Parse failure | ErrEvalParseFailure (exit 1) | Error Handling |
| **EC-8** Unverifiable marking | ErrContractUnverifiable (exit 2) | Error Handling |
| **ES-1** Parse recovery | Eval SKILL.md backup step | Plugin changes |
| **ES-2** Rollback | Eval SKILL.md Step 1 backup | Plugin changes |
| **ES-3** Reviser context | CONTEXT_CONTENT injection | Plugin changes |
| **ES-4** Scope validation | doc-reviser scope check | Plugin changes |
| **CE-1** Config set | SetAutoKey | Interface 4 |
| **CE-2** Config get coverage | GetAutoKeyValue extension | Interface 4 |
| **CE-3** Config init complete | Merge to huh TUI | W8 implementation |
| **CE-4** Schema version | Version field in ForgeConfig | W10 implementation |
| **CE-5** Schema enum alignment | Schema + CLI sync | W10 implementation |
| **CC-1** RunE unification | Run → RunE migration | W8 implementation |
| **CC-2** Config init merge | Delete bufio path | W8 implementation |
| **CC-3** Args validation | cobra.NoArgs / MaximumNArgs | W8 implementation |
| **CC-4** PrintBlock | Output standardization | W8 implementation |
| **CH-1** Naming accuracy | testgen→autogen rename | W1 implementation |
| **CH-2** Dead code deletion | Remove 9 unused exports | W1 implementation |
| **CH-3** Magic value extraction | pkg/constants/ package | W2 implementation |
| **CH-4** Code dedup | isBusinessTask to pkg/task | W2 implementation |

## Open Questions

- [ ] W5 Windows lock: Does `LockFileEx` on Windows NTFS provide the same exclusive-lock semantics as `flock` on POSIX? (Go/No-Go checkpoint)
- [ ] W5 lock crash recovery: If a process holding the advisory lock crashes, does POSIX `flock` release automatically? Verify on Linux and macOS. If not, stale lock will block all writers.
- [ ] W12 package split: How to handle unexported helpers that are shared between commands moving to different packages? (export vs extract to `internal/`)
- [ ] Import cycle resolution: `pkg/constants/` as new package may create cycles with `pkg/feature/` and `pkg/profile/` — verify with `go build` before Phase 1. If cycle found, move constants to `pkg/task/constants.go`.
- [ ] CH-4 `isBusinessTask` dedup: Export as `IsBusinessTask` to `pkg/task/` in Phase 1. Three duplicates (`cmd/validate_index.go`, `pkg/prompt/prompt.go`, `pkg/task/add.go`) → single canonical location.
- [ ] `--force` removal: `--force` flag removed from submit/claim/status. `task reopen` is the only mechanism for leaving rejected/skipped states. Completed tasks are irreversible (create subtask instead).

## Appendix

### Alternatives Considered

| Approach | Pros | Cons | Why Not Chosen |
|----------|------|------|----------------|
| Go FSM library (looplab/state) | Formal state machine, guard/action pattern, DOT graph output | Overkill for 6 states/~10 transitions; new dependency; learning curve; callback-based API doesn't fit CLI command flow | Table-driven `ValidateTransition` is simpler and matches existing codebase style |
| Event sourcing for state | Full audit trail, replay capability | Massive overengineering for a CLI task manager; no event bus infrastructure | Current index.json + state.json approach is sufficient |
| Middleware pattern (cobra middleware) | Centralized pre/post hooks for all commands | Requires invasive cobra changes; RunE already provides error propagation | RunE + Exit() pattern achieves same result with less disruption |

### References

- Michael Feathers, *Working Effectively with Legacy Code* — Characterization Tests pattern
- `spf13/cobra` documentation — RunE vs Run
- Go `os.Rename` atomicity guarantees on POSIX and Windows NTFS
- Kubernetes `cmd/` sub-package organization pattern
- [Claude Code Hooks Reference](../../official-references/hooks.md) — Exit code 0/1/2 semantics for agent-CLI interaction
