---
created: 2026-05-13
prd: prd/prd-spec.md
status: Draft
---

# Technical Design: Forge CLI v3

## Overview

Refactor `task` CLI into `forge` CLI: rename module/directory/binary, regroup 19 flat commands into 5 command groups + 5 visible top-level commands (`version` hidden from `--help`), migrate e2e and probe from justfile bash to Go, update all references across hooks/skills/docs/tests.

Four-phase execution: (1) base rename → (2) command reorganization → (3) e2e/probe migration → (4) reference updates. v3.0.0 major version, clean break, no backward compatibility.

## Architecture

### Layer Placement

Single-layer: CLI tool (`task-cli/` → `forge-cli/`). No API, no UI, no database. The feature touches:
- **Go source**: Cobra commands, internal packages, pkg packages
- **Build artifacts**: binary name, module path, directory structure
- **External consumers**: hooks.json, 23 skill files (9 with task-command refs), 1 agent file, 4 command files, 4 doc files

### Component Diagram

```
                            ┌─────────────────────────────────────────────────┐
                            │                  forge CLI                      │
                            │                                                 │
  ┌─────────────────────────┼─────────────────────────────────────────────┐   │
  │ Command Groups          │                                             │   │
  │  ┌──────────┐ ┌──────┐  │ ┌───────────┐ ┌──────────┐ ┌──────────┐   │   │
  │  │  task    │ │ e2e  │  │ │ forensic  │ │ profile  │ │ prompt   │   │   │
  │  │ (10 cmd) │ │(6 cmd)│ │ │ (3 cmd)   │ │ (3 cmd)  │ │ (1 cmd)  │   │   │
  │  └────┬─────┘ └──┬───┘  │ └─────┬─────┘ └────┬─────┘ └────┬─────┘   │   │
  │       │          │       │       │            │            │          │   │
  └───────┼──────────┼───────┼───────┼────────────┼────────────┼──────────┘   │
          │          │       │       │            │            │              │
  ┌───────┼──────────┼───────┼───────┼────────────┼────────────┼──────────┐   │
  │ Top-level         │       │       │            │            │          │   │
  │  feature  probe   cleanup quality-gate  verify-task-done  version(H) │   │
  └───────┬──────────┬───┬───┬───────────┬────────────┬─────────────────┘   │
          │          │   │   │           │            │                      │
  ┌───────┼──────────┼───┼───┼───────────┼────────────┼──────────────────┐    │
  │ Pkg   │          │   │   │           │            │                  │    │
  │  ┌────▼────┐ ┌───▼───▼───▼──┐ ┌──────▼──────┐ ┌─▼────────┐         │    │
  │  │ profile │ │   e2e (new)  │ │  e2eprobe   │ │  prompt  │         │    │
  │  │         │ │              │ │  (existing) │ │          │         │    │
  │  └─────────┘ └──────────────┘ └─────────────┘ └──────────┘         │    │
  │  ┌─────────┐ ┌──────────────┐ ┌─────────────┐ ┌──────────┐         │    │
  │  │  task   │ │    just      │ │   feature   │ │  version │         │    │
  │  └─────────┘ └──────────────┘ └─────────────┘ └──────────┘         │    │
  └──────────────────────────────────────────────────────────────────────┘    │
                            ┌─────────────────────────────────────────────┐   │
                            │           External Consumers                │   │
                            │  hooks.json · skills/ · agents/ · docs/     │   │
                            └─────────────────────────────────────────────┘   │
                            └─────────────────────────────────────────────────┘
```

> †1 task group: 10 commands (not PRD's 11) — `verify-task-done` is top-level only, not nested under `task`. See PRD Divergences.
> †2 e2e group: 6 commands (not PRD's 5) — `validate-specs` moved from top-level into `e2e` group. See PRD Divergences.

### PRD Divergences

| PRD Count | Design Count | Delta | Justification |
|-----------|-------------|-------|---------------|
| `forge task` = 11 subcommands (includes `verify-task-done`) | 10 subcommands | -1 | `verify-task-done` is top-level only (see Alternatives: "User chose top-level only"). PRD lists it under both `task` group (line 217) and top-level commands (line 231). Design resolves the ambiguity by placing it at top-level exclusively, consistent with PRD line 242 mapping `task verify-completion` to `forge task verify-task-done` but also listing it as a top-level Hook command (line 231). The top-level placement avoids the redundancy of having the same command accessible through two paths, which was rejected in the Alternatives table. |
| `forge e2e` = 5 subcommands (no `validate-specs`) | 6 subcommands | +1 | `validate-specs` was a top-level command (`task validate-specs`) in the current codebase. Moving it into the `e2e` group is a logical grouping improvement — it validates Playwright spec files, which are e2e assets. This keeps all e2e-related operations under one group. The PRD's 5-count reflects the justfile migration scope (run/setup/verify/compile/discover) but does not account for this existing command's regrouping. |

### Dependencies

**Internal** (unchanged):
- `internal/cmd` → `pkg/*` (existing direction, no new dependencies)
- `pkg/profile` → used by new `pkg/e2e` and `internal/cmd/e2e_*`

**External** (unchanged):
- `github.com/spf13/cobra` — CLI framework
- `gopkg.in/yaml.v3` — config parsing

**New internal package**:
- `pkg/e2e/` — e2e execution logic migrated from justfile bash

### File Structure Changes

```
task-cli/                          →  forge-cli/
  cmd/task/                        →  cmd/forge/
    main.go                           main.go (Use: "forge")
    run.go                            run.go
  internal/cmd/
    root.go  (Use: "task")         →  root.go  (Use: "forge")
    record.go                      →  submit.go (Use: "submit")
    check.go                       →  check_deps.go (Use: "check-deps")
    validate.go                    →  validate_index.go (Use: "validate-index")
    verify_completion.go           →  verify_task_done.go (Use: "verify-task-done")
    all_completed.go               →  quality_gate.go (Use: "quality-gate")
    prompt.go                      →  split into prompt_parent.go + prompt_get.go
    validate_specs.go              →  e2e_validate_specs.go (under e2e group)
    template.go                    →  DELETED
    (new)                             task_parent.go (group parent)
                                       e2e_parent.go (group parent)
                                       prompt_parent.go (group parent)
                                       e2e_run.go
                                       e2e_setup.go
                                       e2e_verify.go
                                       e2e_compile.go
                                       e2e_discover.go
                                       list_types.go
                                       probe.go
  pkg/
    e2e/                           →  NEW: e2e execution logic
    index/lock.go                  →  NEW: advisory file lock for concurrent write safety
    version/version.go             →  Name: "task" → "forge"
  go.mod                           →  module forge-cli
  scripts/version.txt              →  3.0.0
  scripts/install-local.ps1       →  forge references
  scripts/install-local.sh        →  forge references
```

## Interfaces

### 1. Command Group Parents

#### `forge task` (group parent)

```go
// task_parent.go
taskCmd = &cobra.Command{
    Use:   "task",
    Short: "Manage task lifecycle",
}
```

Subcommands: claim, submit, status, query, check-deps, validate-index, add, index, migrate, list-types (10 total). See PRD Divergences for why this is 10 (not PRD's 11).

#### `forge e2e` (group parent)

```go
// e2e_parent.go
e2eCmd = &cobra.Command{
    Use:   "e2e",
    Short: "End-to-end test management",
}
```

Subcommands: run, setup, verify, compile, discover, validate-specs (6 total). See PRD Divergences for why this is 6 (not PRD's 5).

#### `forge prompt` (group parent)

```go
// prompt_parent.go
promptCmd = &cobra.Command{
    Use:   "prompt",
    Short: "Manage agent execution prompts",
}
```

Subcommand: get-by-task-id (1 total).

### 2. Renamed Commands

#### `forge task submit` (renamed from `record`)

```go
// submit.go (renamed from record.go)
submitCmd = &cobra.Command{
    Use:   "submit <task-id>",
    Short: "Submit task execution result",
    Args:  cobra.ExactArgs(1),
    Run:   runSubmit, // renamed from runRecord, identical logic
}

func init() {
    submitCmd.Flags().StringP("data", "d", "", "Path to record data JSON file (required)")
    submitCmd.Flags().Bool("force", false, "Overwrite existing record")
}
```

| Flag | Shorthand | Type | Required | Default |
|------|-----------|------|----------|---------|
| `--data` | `-d` | `string` | yes | `""` |
| `--force` | — | `bool` | no | `false` |

#### `forge task check-deps` (renamed from `check`)

```go
// check_deps.go (renamed from check.go)
checkDepsCmd = &cobra.Command{
    Use:   "check-deps",
    Short: "Check task dependencies",
    Args:  cobra.NoArgs,
    Run:   runCheckDeps, // renamed from runCheck, identical logic
}
```

#### `forge task validate-index` (renamed from `validate`)

```go
// validate_index.go (renamed from validate.go)
validateIndexCmd = &cobra.Command{
    Use:   "validate-index [file]",
    Short: "Validate index.json file",
    Args:  cobra.MaximumNArgs(1),
    Run:   runValidateIndex, // renamed from runValidate, identical logic
}
```

#### `forge verify-task-done` (renamed from `verify-completion`, top-level)

```go
// verify_task_done.go (renamed from verify_completion.go)
verifyTaskDoneCmd = &cobra.Command{
    Use:   "verify-task-done",
    Short: "Verify task completion before git commit",
    Args:  cobra.NoArgs,
    Run:   runVerifyTaskDone, // renamed from runVerifyCompletion, identical logic
}
```

#### `forge quality-gate` (renamed from `all-completed`, top-level)

```go
// quality_gate.go (renamed from all_completed.go)
qualityGateCmd = &cobra.Command{
    Use:   "quality-gate",
    Short: "Check if all tasks are done, then run tests",
    Args:  cobra.NoArgs,
    Run:   runQualityGate, // renamed from runAllCompleted, adds max fix-task cap
}
```

`quality-gate` is the only renamed command with a behavioral addition: max 3 concurrent fix-tasks per failure step. When the cap is reached, the command prints "max fix-tasks reached for <step>, manual intervention required" to stderr and exits 1 instead of creating a new fix-task.

#### Prompt split: `task prompt <id>` → `forge prompt get-by-task-id <id>`

```go
// prompt_parent.go — new group parent
promptCmd = &cobra.Command{
    Use:   "prompt",
    Short: "Manage agent execution prompts",
}

// prompt_get.go (split from prompt.go)
promptGetCmd = &cobra.Command{
    Use:   "get-by-task-id <id>",
    Short: "Synthesize the agent prompt for a task",
    Args:  cobra.ExactArgs(1),
    Run:   runPromptGet, // identical to current runPrompt logic
}

func init() {
    promptGetCmd.Flags().Bool("fix-record-missed", false, "Include fix-record-missed context in prompt")
}
```

| Flag | Shorthand | Type | Required | Default |
|------|-----------|------|----------|---------|
| `--fix-record-missed` | — | `bool` | no | `false` |

### 2b. Unchanged Commands — Existing Interface Reference

The following commands retain their current Cobra struct signatures unchanged. Listed for task-breakdown completeness.

#### `forge forensic` (group with 3 subcommands)

```go
// forensic.go
forensicCmd = &cobra.Command{
    Use:   "forensic",
    Short: "Analyze Claude Code session transcripts for agent deviation forensics",
}

forensicSearchCmd = &cobra.Command{
    Use:   "search [project-path]",
    Short: "Search history.jsonl for matching sessions",
    Args:  cobra.MaximumNArgs(1),
    Run:   runForensicSearch,
}

forensicExtractCmd = &cobra.Command{
    Use:   "extract <session-jsonl-path>",
    Short: "Extract compact evidence from a session transcript",
    Args:  cobra.ExactArgs(1),
    Run:   runForensicExtract,
}

func init() {
    forensicExtractCmd.Flags().String("slug", "", "Filter by feature slug")
    forensicExtractCmd.Flags().String("mode", "compact", "Output mode: compact|full")
}

forensicSubagentsCmd = &cobra.Command{
    Use:   "subagents <session-dir-path>",
    Short: "List subagent transcripts for a session",
    Args:  cobra.ExactArgs(1),
    Run:   runForensicSubagents,
}
```

#### `forge profile` (group with 3 subcommands)

```go
// profile.go
profileCmd = &cobra.Command{
    Use:   "profile",
    Short: "Resolve or set the active test profile",
}

profileSetCmd = &cobra.Command{
    Use:   "set <name>",
    Short: "Set the active test profile in .forge/config.yaml",
    Args:  cobra.ExactArgs(1),
    Run:   runProfileSet,
}

profileDetectCmd = &cobra.Command{
    Use:   "detect",
    Short: "Detect test profiles from project structure (ignores config)",
    Args:  cobra.NoArgs,
    Run:   runProfileDetect,
}

profileGetCmd = &cobra.Command{
    Use:   "get <name>",
    Short: "Get profile strategy file content",
    Args:  cobra.ExactArgs(1),
    Run:   runProfileGet,
}
```

#### `forge cleanup` (top-level)

```go
// cleanup.go
cleanupCmd = &cobra.Command{
    Use:   "cleanup",
    Short: "Clean up completed task state",
    Args:  cobra.NoArgs,
    Run:   runCleanup,
}
```

#### `forge feature` (top-level)

```go
// feature.go
featureCmd = &cobra.Command{
    Use:   "feature [slug]",
    Short: "Set or display the current feature",
    Args:  cobra.MaximumNArgs(1),
    Run:   runFeature,
}
```

#### `forge version` (top-level, hidden)

```go
// version.go
versionCmd = &cobra.Command{
    Use:    "version",
    Short:  "Print the CLI version",
    Hidden: true, // hidden from --help
    Run:    runVersion,
}
```

#### `forge task` subcommands (unchanged signatures)

```go
// claim.go
claimCmd = &cobra.Command{
    Use:   "claim",
    Short: "Claim the next available task",
    Args:  cobra.NoArgs,
    Run:   runClaim,
}

// status.go
statusCmd = &cobra.Command{
    Use:   "status <task-id> [status]",
    Short: "Query or update task status",
    Args:  cobra.RangeArgs(1, 2),
    Run:   runStatus,
}

// query.go
queryCmd = &cobra.Command{
    Use:   "query <task-id-or-key>",
    Short: "Query task information",
    Args:  cobra.ExactArgs(1),
    Run:   runQuery,
}

// add.go
addCmd = &cobra.Command{
    Use:   "add",
    Short: "Add a new task to the current feature",
    Run:   runAdd,
}

func init() {
    addCmd.Flags().String("title", "", "Task title (required)")
    addCmd.Flags().String("id", "", "Explicit task ID (auto-generated if empty)")
    addCmd.Flags().String("priority", "P1", "Task priority: P0|P1|P2")
    addCmd.Flags().StringSlice("depends-on", nil, "Comma-separated task IDs this task depends on")
    addCmd.MarkFlagRequired("title")
}

// index.go
indexCmd = &cobra.Command{
    Use:   "index",
    Short: "Build or rebuild index.json from task markdown files",
    Run:   runIndex,
}

func init() {
    indexCmd.Flags().String("feature", "", "Feature slug to index (required)")
    indexCmd.Flags().Bool("no-test", false, "Skip test-related index entries")
    indexCmd.Flags().StringSlice("test-profiles", nil, "Comma-separated test profile names")
    indexCmd.MarkFlagRequired("feature")
}

// migrate.go
migrateCmd = &cobra.Command{
    Use:   "migrate",
    Short: "Migrate index.json by inferring type fields for all tasks",
    Args:  cobra.NoArgs,
    Run:   runMigrate,
}
```

### 3. New Command: `forge task list-types`

```go
listTypesCmd = &cobra.Command{
    Use:   "list-types",
    Short: "List all supported task types",
    Args:  cobra.NoArgs,
    Run:   runListTypes,
}
// Output: one line per type, format: "<type>  <description>"
// Source: hardcoded registry matching pkg/task/infer.go type list
```

Types: implementation, fix, gate, doc-generation-summary, doc-generation-consolidate, test-pipeline-gen-cases, test-pipeline-eval-cases, test-pipeline-gen-scripts, test-pipeline-run, test-pipeline-graduate, test-pipeline-verify-regression (11 total).

### 4. New Command: `forge probe`

```go
probeCmd = &cobra.Command{
    Use:   "probe [path]",
    Short: "HTTP health check for e2e servers",
    Args:  cobra.MaximumNArgs(1),
    Run:   runProbe,
}
// Delegates to e2eprobe.ProbeServers()
// path defaults to "/health"
```

### 5. New E2E Subcommands

Each e2e subcommand follows the same pattern: read profile → validate → dispatch to profile-specific handler.

#### `forge e2e run --feature <slug>`

```go
e2eRunCmd = &cobra.Command{
    Use:   "run",
    Short: "Run e2e tests (profile-aware)",
    Run:   runE2ERun,
}

func init() {
    e2eRunCmd.Flags().String("feature", "", "Run tests for a specific feature (empty = all)")
}
```

| Flag | Shorthand | Type | Required | Default |
|------|-----------|------|----------|---------|
| `--feature` | — | `string` | no | `""` |

Profile dispatch:

| Profile | Execution |
|---------|-----------|
| web-playwright | `npx playwright test` (feature: `--config` or `E2E_FEATURE=1`) |
| others | stderr "unsupported profile for run: <name>", exit 1 |

#### `forge e2e setup`

```go
e2eSetupCmd = &cobra.Command{
    Use:   "setup",
    Short: "Install e2e dependencies (idempotent)",
    Run:   runE2ESetup,
}

func init() {
    e2eSetupCmd.Flags().Bool("force", false, "Force reinstall dependencies")
}
```

| Flag | Shorthand | Type | Required | Default |
|------|-----------|------|----------|---------|
| `--force` | — | `bool` | no | `false` |

#### `forge e2e verify --feature <slug>`

```go
e2eVerifyCmd = &cobra.Command{
    Use:   "verify",
    Short: "Check for unresolved VERIFY markers",
    Run:   runE2EVerify,
}

func init() {
    e2eVerifyCmd.Flags().String("feature", "", "Feature slug to verify (required)")
    e2eVerifyCmd.MarkFlagRequired("feature")
}
```

| Flag | Shorthand | Type | Required | Default |
|------|-----------|------|----------|---------|
| `--feature` | — | `string` | **yes** | `""` |

#### `forge e2e compile`

```go
e2eCompileCmd = &cobra.Command{
    Use:   "compile",
    Short: "Compile-check e2e test files",
    Run:   runE2ECompile,
}
```

Profile dispatch:

| Profile | Execution |
|---------|-----------|
| web-playwright | `npx tsc --noEmit` |
| go-test | `go build ./tests/e2e/...` |
| pytest | `python -m compileall tests/e2e/ -q` |

#### `forge e2e discover`

```go
e2eDiscoverCmd = &cobra.Command{
    Use:   "discover",
    Short: "List all e2e test cases without running",
    Run:   runE2EDiscover,
}
```

Profile dispatch:

| Profile | Execution |
|---------|-----------|
| web-playwright | `npx playwright test --list` |
| go-test | `go test ./tests/e2e/... -list '.*' -tags=e2e` |
| pytest | `python -m pytest tests/e2e/ --collect-only -q` |

#### `forge e2e validate-specs` (moved from top-level)

```go
e2eValidateSpecsCmd = &cobra.Command{
    Use:   "validate-specs",
    Short: "Validate Playwright spec files via AST analysis",
    Run:   runE2EValidateSpecs,
}
// Behavior unchanged from current task validate-specs
```

### 6. New Package: `pkg/e2e/`

Shared e2e execution logic, replacing justfile bash.

```go
// pkg/e2e/e2e.go
type RunOpts struct {
    ProjectRoot string
    Feature     string // empty = run all
    Force       bool   // for setup
}

// ResolveProfile reads config.yaml and validates the profile.
// Returns profile name or error.
func ResolveProfile(projectRoot string) (string, error)

// run.go
func Run(opts RunOpts) error

// setup.go
func Setup(opts RunOpts) error

// verify.go
func Verify(opts RunOpts) error

// compile.go
func Compile(projectRoot string) error

// discover.go
func Discover(projectRoot string) error
```

`ResolveProfile` uses `profile.ReadTestProfiles()` for reading and validates against `profile.KnownProfiles`. Returns error with "no e2e profile configured" or "unknown profile: <name>" messages matching PRD error table.

Each action function implements a `switch` on profile name and calls external tools via an `ExecRunner` interface (injectable `exec.Command` wrapper). Error handling wraps stderr from the external tool.

```go
// pkg/e2e/exec.go — injectable wrapper for testability
type ExecRunner interface {
    Run(name string, args ...string) ([]byte, error)
}

// RealExec is the production implementation using os/exec.
type RealExec struct{}

func (RealExec) Run(name string, args ...string) ([]byte, error) {
    cmd := exec.Command(name, args...)
    return cmd.CombinedOutput()
}

// Each action function accepts ExecRunner via package-level variable:
var runner ExecRunner = RealExec{}
```

Unit tests inject a hand-rolled mock `ExecRunner` that returns predetermined outputs/errors — consistent with the existing codebase convention (e.g., `validate_specs_test.go` uses hand-rolled mock scripts, no external mock libraries). The project uses only Go's standard `testing` package (`t.Fatalf`, `t.Errorf`) for assertions, with no `testify` or `gomock` dependency (verified in `go.mod`). Integration tests use `RealExec` against test fixtures. No build tags required.

```go
// pkg/e2e/exec_test.go — hand-rolled mock matching project convention
type stubExec struct {
    responses map[string]execResponse
}
type execResponse struct {
    output []byte
    err    error
}

func (s *stubExec) Run(name string, args ...string) ([]byte, error) {
    key := name + " " + strings.Join(args, " ")
    if r, ok := s.responses[key]; ok {
        return r.output, r.err
    }
    return nil, fmt.Errorf("stubExec: unexpected command: %s", key)
}

// Usage in tests:
// s := &stubExec{responses: map[string]execResponse{
//     "npx playwright test --list": {output: []byte("test1\ntest2\n"), err: nil},
// }}
// runner = s  // set package-level var
```

### 7. Quality-Gate Max Fix-Task Cap

```go
// In quality_gate.go (renamed from all_completed.go)
const maxFixTasksPerStep = 3

// countActiveFixTasks counts fix-tasks for a step that are not completed/skipped.
func countActiveFixTasks(index *task.TaskIndex, step string) int

// addFixTask checks the cap before creating a fix task.
// Returns (taskID, nil) on success, ("", error) when cap exceeded.
// Signature mirrors current all_completed.go but adds cap check before the
// existing task.AddTask → task.CreateTaskMarkdown → feature.EnsureForgeState sequence.
func addFixTask(projectRoot, featureSlug, step, output, errorDocPath string) (string, error)
```

Fix-task identification and counting:

```go
// countActiveFixTasks iterates TaskIndex entries, filters by fix-task
// criteria, and returns the count of active (non-terminal) fix-tasks
// for the given quality-gate step.
//
// TaskIndex API surface (existing, from pkg/task/):
//   type TaskIndex map[string]*Task  // keyed by task ID
//   type Task struct {
//       SourceTaskID string  // non-empty iff this is a fix-task
//       Title        string
//       Status       string  // "pending"|"in_progress"|"blocked"|"completed"|"skipped"
//   }
//
// Step-name matching convention:
//   addFixTask titles are formatted as "fix <step>: <first line of output>"
//   (see existing all_completed.go). Matching uses prefix "fix <step>:"
//   to avoid false positives from unrelated tasks whose titles happen to
//   contain the step name as a substring.
func countActiveFixTasks(index *task.TaskIndex, step string) int {
    count := 0
    prefix := "fix " + step + ":"
    for _, t := range *index {
        if t.SourceTaskID != "" &&
            strings.HasPrefix(t.Title, prefix) &&
            t.Status != "completed" &&
            t.Status != "skipped" {
            count++
        }
    }
    return count
}
```

Terminal states excluded: `"completed"`, `"skipped"`. Non-terminal states counted: `"pending"`, `"in_progress"`, `"blocked"`. The `SourceTaskID != ""` check is the primary filter (all fix-tasks have a non-empty source), with the title prefix as a secondary disambiguator to scope the count to a specific quality-gate step.

### 8. Deleted Command: `template`

Remove `template.go` and its test. No callers found in codebase.

### 9. Concurrent Write Conflict Handling (PRD S3)

The PRD requires that when two agents call `forge task submit` concurrently for different tasks in the same feature, one succeeds and the other receives a "concurrent write conflict, retry" error with exit code 1, while `index.json` remains valid JSON.

**Problem**: `submit` writes to `index.json` (task status update) and creates a record file. Two concurrent writes can corrupt `index.json` or lose one agent's update.

**Design**: Advisory file lock using `flock(2)` on a lock file alongside `index.json`.

```go
// pkg/index/lock.go — NEW file
// LockFile acquires an exclusive flock on <feature-dir>/tasks/index.json.lock.
// Creates the lock file if it does not exist.
// timeout: 5 seconds. On timeout, returns ErrLockConflict.
func LockFile(indexPath string) (*os.File, error)

// UnlockFile releases the flock and closes the file descriptor.
func UnlockFile(f *os.File) error

var ErrLockConflict = errors.New("concurrent write conflict, retry")
```

**Integration into submit flow**:

```go
// In submit.go runSubmit():
func runSubmit(cmd *cobra.Command, args []string) {
    // 1. Acquire lock on index.json.lock
    lock, err := index.LockFile(indexPath)
    if err == index.ErrLockConflict {
        fmt.Fprintln(os.Stderr, "concurrent write conflict, retry")
        os.Exit(1)
    }
    defer index.UnlockFile(lock)

    // 2. Read index.json (guaranteed no concurrent writer)
    // 3. Validate task status (terminal state check)
    // 4. Write record file
    // 5. Update index.json atomically (write to .tmp, rename)
    // 6. Lock released by defer
}
```

**Lock scope**: Per-feature `index.json.lock` file (in `docs/features/<slug>/tasks/`). Different features have independent locks — no cross-feature contention.

**Timeout**: 5-second wait using `fcntl(F_SETLKW)` with a goroutine+timer fallback on Windows (where `fcntl` is unavailable). On timeout, return `ErrLockConflict` immediately rather than blocking indefinitely.

**Atomic index write**: Inside the lock, write `index.json` via temp file + rename (`os.Rename`), not direct overwrite. This ensures `index.json` is never in a partial state, even if the process crashes mid-write. A reader that opens the file between rename steps sees either the old or new complete version.

**Error handling**:

| Condition | Exit Code | stderr |
|-----------|-----------|--------|
| Lock acquired, submit succeeds | 0 | (normal output) |
| Lock acquisition times out (5s) | 1 | `concurrent write conflict, retry` |
| Lock file creation fails (permissions) | 1 | `failed to create lock file: <reason>` |

**Why flock**: The codebase has no existing locking mechanism. `flock(2)` is the standard POSIX advisory lock — lightweight, kernel-managed, auto-released on process exit (crash-safe). No external dependency required. On Windows, Go's `syscall` package provides equivalent locking via `LockFileEx`.

**Cleanup**: Lock files are not deleted after use — they are reused on subsequent writes. `forge cleanup` does not touch `index.json.lock` files.

## Data Models

No database. Key Go structs affected:

### TaskType Registry (new)

```go
// pkg/task/types.go — add TaskTypeRegistry
type TaskTypeInfo struct {
    Name        string
    Description string // <= 60 chars, verb+object
}

var TaskTypeRegistry = []TaskTypeInfo{
    {"implementation", "Implement task requirements"},
    {"fix", "Fix a bug or quality gate failure"},
    {"gate", "Verify task completion quality"},
    {"doc-generation-summary", "Generate summary documentation"},
    {"doc-generation-consolidate", "Consolidate specs to project level"},
    {"test-pipeline-gen-cases", "Generate structured test cases"},
    {"test-pipeline-eval-cases", "Evaluate test cases quality"},
    {"test-pipeline-gen-scripts", "Generate e2e test scripts"},
    {"test-pipeline-run", "Execute e2e test scripts"},
    {"test-pipeline-graduate", "Graduate tests to regression suite"},
    {"test-pipeline-verify-regression", "Verify regression suite passes"},
}
```

Source of truth for `list-types` command and `prompt` type dispatch. Currently scattered across `pkg/task/infer.go` and `pkg/prompt/data/`. Centralizing avoids drift.

### Version Name Change

```go
// pkg/version/version.go
var Name = "forge"  // was "task"
```

### Module Path Change

```
// go.mod
module forge-cli  // was task-cli
```

All import paths: `task-cli/...` → `forge-cli/...`.

## Error Handling

### New Error Cases

| Command | Condition | Exit Code | stderr |
|---------|-----------|-----------|--------|
| `forge e2e *` | No profile in config.yaml | 1 | `no e2e profile configured` |
| `forge e2e *` | Unknown profile value | 1 | `unknown profile: <value>` |
| `forge e2e run` | Feature dir not found | 1 | `feature not found: <slug>` |
| `forge quality-gate` | Max fix-tasks per step reached | 1 | `max fix-tasks reached for <step>, manual intervention required` |
| `forge probe` | Server not responding | 1 | `FAIL: <url> not responding` |
| `forge task submit` | Lock acquisition timeout (concurrent writer) | 1 | `concurrent write conflict, retry` |
| `forge task submit` | Lock file creation failure (permissions) | 1 | `failed to create lock file: <reason>` |

### E2E External Tool Failures

Exit code strategy: **normalize to 1** for all external tool failures. Do not propagate the child process exit code, because callers (hooks, CI) only check zero vs non-zero and different tools use inconsistent non-zero codes (npx uses 1 for test failure but 2 for missing binary; go uses 1 for build failure; python uses various codes). The child's stderr is captured and forwarded unchanged.

stderr format: `<command> failed: <first line of child stderr>`. The full child stderr is logged at debug level for diagnosis.

| Command | Profile | External Tool | Failure Condition | Exit Code | stderr Example |
|---------|---------|---------------|-------------------|-----------|----------------|
| `forge e2e run` | web-playwright | `npx playwright test` | Test failure or timeout | 1 | `npx playwright test failed: 1 failed, 3 passed (15s)` |
| `forge e2e run` | web-playwright | `npx playwright test` | npx binary not found | 1 | `npx playwright test failed: npx: command not found` |
| `forge e2e setup` | web-playwright | `npx playwright install` | Download/install failure | 1 | `npx playwright install failed: EACCES: permission denied` |
| `forge e2e setup` | go-test | `go install ...` | Go tool install failure | 1 | `go install failed: cannot find package` |
| `forge e2e setup` | pytest | `python -m pip install` | Pip install failure | 1 | `python -m pip install failed: ERROR: Could not find a version` |
| `forge e2e verify` | web-playwright | `grep -R VERIFY` | VERIFY markers found | 1 | `VERIFY markers found in: spec1.ts, spec2.ts` |
| `forge e2e verify` | go-test | `grep -R VERIFY tests/e2e/` | VERIFY markers found | 1 | `VERIFY markers found in: foo_test.go` |
| `forge e2e verify` | pytest | `grep -R VERIFY` | VERIFY markers found | 1 | `VERIFY markers found in: test_foo.py` |
| `forge e2e compile` | web-playwright | `npx tsc --noEmit` | TypeScript compilation errors | 1 | `npx tsc --noEmit failed: src/test.ts(10,5): error TS2322` |
| `forge e2e compile` | go-test | `go build ./tests/e2e/...` | Go compilation errors | 1 | `go build failed: ./tests/e2e/main_test.go:15: undefined: Foo` |
| `forge e2e compile` | pytest | `python -m compileall` | Python syntax errors | 1 | `python -m compileall failed: SyntaxError: invalid syntax` |
| `forge e2e discover` | web-playwright | `npx playwright test --list` | Playwright listing failure | 1 | `npx playwright test --list failed: no tests found` |
| `forge e2e discover` | go-test | `go test -list` | Go test listing failure | 1 | `go test -list failed: build constraints exclude all tests` |
| `forge e2e discover` | pytest | `python -m pytest --collect-only` | Pytest collection failure | 1 | `python -m pytest --collect-only failed: ImportError: No module named pytest` |
| `forge e2e validate-specs` | web-playwright | AST parsing (Go stdlib) | Invalid spec file syntax | 1 | `invalid spec file: <path>: unexpected token at line <n>` |

### Preserved Error Behavior

All existing error codes and messages from renamed commands remain identical. The error types in `internal/cmd/errors.go` are unchanged — only command `Use` strings change.

### New Go Error Types

```go
// pkg/index/lock.go — ErrLockConflict (defined in §9 above)

// pkg/e2e/e2e.go
var (
    ErrNoProfile    = errors.New("no e2e profile configured")
    ErrBadProfile   = errors.New("unknown profile")  // wrapped with value
    ErrFeatureNotFound = errors.New("feature not found")  // wrapped with slug
)

// internal/cmd/quality_gate.go
var ErrMaxFixTasks = errors.New("max fix-tasks reached")  // wrapped with step name
```

### Propagation Strategy

CLI-only: errors propagate up to `cobra.Command.Run` → `Exit(err)` → `os.Exit(1)`. No cross-layer error translation needed.

## Cross-Layer Data Map

Single-layer feature (CLI only). Cross-Layer Data Map not applicable.

## Integration Specs

No existing-page integrations — not applicable.

## Testing Strategy

### Per-Layer Test Plan

| Layer | Test Type | Tool | What to Test | Coverage Target |
|-------|-----------|------|--------------|-----------------|
| Cobra commands | Unit | `go test` (std `testing`) | Command registration, flag parsing, group structure | 80% |
| pkg/e2e | Unit | `go test` + hand-rolled `stubExec` (see §6 mock struct) | Profile resolution, dispatch, error cases | 80% |
| pkg/e2e | Integration | `go test` + `RealExec` | Actual exec.Command calls against test fixtures | N/A |
| Renamed commands | Regression | `go test` | All existing tests pass with new command names | 100% |
| quality-gate cap | Unit | `go test` (std `testing`) | Max fix-task counting and cap enforcement | 80% |
| Phase 4 ref updates | Verification | `go test` + `just check-stale-refs` | Every file in Phase 4 map has all old refs replaced; no stale `task <cmd>` patterns remain; markdown files parse without errors | 100% of mapped files |

### Key Test Scenarios

1. **Command structure**: `forge --help` shows exactly 5 groups + 5 top-level entries (version hidden via `cobra.Command.Hidden: true`, still accessible as `forge version`)
2. **Group membership**: `forge task --help` shows exactly 10 subcommands with descriptions <= 80 chars
3. **Unknown command suggestions**: `forge taks` → suggests "task"; `forge task xyz` → lists valid subcommands
4. **submit = record**: existing record tests pass unchanged (only binary name and subcommand path differ)
5. **e2e profile resolution**: valid profile → dispatches; missing → error; unknown → error with valid list
6. **quality-gate cap**: 3 active fix-tasks for a step → cap reached error on 4th attempt
7. **probe**: no config.yaml → "OK: CLI-only project"; server down → exit 1
8. **list-types**: outputs 11 types, each with description; empty registry → 0 lines
9. **Phase 4 reference completeness**: `grep -rE '\btask (claim|submit|status|query|check-deps|validate-index|verify-task-done|quality-gate|cleanup|feature|prompt|add|index|migrate|validate-specs|record|all-completed|verify-completion|check|validate)\b' plugins/ forge-cli/docs/` returns zero matches after Phase 4; `just check-stale-refs` CI target passes; each file listed in the Phase 4 Reference Update Map is valid markdown/JSON after replacement

### Migration Equivalence Tests

For each e2e command migrated from justfile:
1. Create test fixture with known profile config
2. Capture justfile recipe output (baseline)
3. Capture Go command output
4. Assert exit code and key output markers match using `t.Fatalf` and `strings.Contains` (standard library only — no assertion library)

### Overall Coverage Target

80% (matching existing project convention in task-cli/CLAUDE.md).

## Security Considerations

### Threat Model

No new threats. The refactoring changes naming and grouping but not authorization, data access, or network behavior.

### Mitigations

- **File locking**: PRD S3 requires concurrent write conflict handling. Design adds advisory file lock via `flock(2)` in new `pkg/index/lock.go` (see §9). Lock scope is per-feature `index.json.lock`. 5-second timeout with `ErrLockConflict` on timeout. Auto-released on process exit (crash-safe).
- **Command injection**: New `pkg/e2e/` calls external tools via `exec.Command` with explicit args (no shell interpolation). Input comes from `config.yaml` (profile name) and CLI flags (feature slug), both validated before use.

## PRD Coverage Map

| PRD AC (User Story) | Design Component | Interface / Model |
|---------------------|------------------|-------------------|
| S1: forge --help shows <= 10 entries | root.go group registration | taskCmd, e2eCmd, forensicCmd, profileCmd, promptCmd + 5 top-level (version hidden via `cobra.Command.Hidden: true`) |
| S1: forge task --help shows 10 subcommands, desc <= 80 chars | task_parent.go + subcommand Shorts | Each subcommand `Short` field |
| S1: unknown command suggests alternatives | Cobra built-in suggestions | root.go `SilenceErrors: true`, `SilenceUsage: true` |
| S1: unknown subcommand lists valid ones | Cobra built-in | Enabled by default for subcommands |
| S2: prompt get-by-task-id returns rendered prompt | prompt_get.go | Delegates to existing `prompt.Synthesize()` |
| S2: nonexistent task ID → error | prompt_get.go | Existing `ErrTaskNotFound` |
| S2: missing type → error | prompt_get.go | Existing "missing task type" error |
| S3: submit updates index + generates record | submit.go | Existing `runRecord` logic |
| S3: terminal state task rejected | submit.go | Existing status validation |
| S3: missing --data flag → error | submit.go | Existing `readRecordData` error |
| S3: concurrent write conflict | submit.go + `pkg/index/lock.go` | Advisory file lock via `flock(2)` on `index.json.lock` — see §9 below |
| S4: cleanup removes state files | cleanup.go | Unchanged |
| S4: quality-gate runs compile→fmt→lint→test | quality_gate.go | Existing `RunGate` sequence |
| S4: no terminal tasks → no-op | cleanup.go | Existing behavior |
| S4: fix-task cap at 3 per step | quality_gate.go | New `countActiveFixTasks` + `maxFixTasksPerStep` |
| S5: e2e run reads profile | e2e_run.go → `e2e.ResolveProfile` | Uses `profile.ReadTestProfiles()` |
| S5: no profile → error | `e2e.ResolveProfile` | "no e2e profile configured" |
| S5: unknown profile → error | `e2e.ResolveProfile` | "unknown profile: <value>" |
| S5: feature not found → error | e2e_run.go | "feature not found: <slug>" |
| S6: list-types shows 11 types with desc | list_types.go | `TaskTypeRegistry` |
| S6: empty registry → 0 lines | list_types.go | Dynamic output from registry |
| S7: forensic search/extract/subagents | forensic.go | Unchanged, already has subcommands |
| S7: file not found → error | forensic.go | Existing error handling |
| S8: profile detect/set/get | profile.go | Unchanged |
| S8: unknown profile → error | profile.go | Existing validation |

## Resolved Design Decisions

- `forge --help` entry count: `version` command registered with `cobra.Command.Hidden: true`, so it does not appear in `--help` groups or top-level listing. Visible entries: 5 groups + 5 top-level = 10 total, matching PRD goal. `forge version` remains accessible by direct invocation.
- `pkg/e2e/` test strategy: `ExecRunner` interface pattern (struct injection in `pkg/e2e/exec.go`). Production code uses `RealExec`; unit tests inject mock implementations returning predetermined outputs. No build tags.

## Appendix

### Alternatives Considered

| Approach | Pros | Cons | Why Not Chosen |
|----------|------|------|----------------|
| Cobra aliases (keep old names as aliases) | Zero reference updates needed | Defeats rename purpose, --help shows both names | PRD explicitly says no backward compat |
| e2e via `just` wrapper (Go calls justfile) | Simpler migration, less Go code | Keeps justfile dependency, doesn't eliminate bash duplication | PRD goal is to eliminate duplicated bash code |
| Profile runners as interface/registry | Extensible for new profiles | Over-engineering for 6 known profiles | Simple switch statement suffices |
| `forge task verify-task-done` as alias of top-level | Both entry points work | Two ways to do the same thing, confusing for agents | User chose top-level only |

### References

- Current command list: `task-cli/internal/cmd/root.go`
- E2E justfile recipes: `justfile` lines 62-249
- Profile system: `task-cli/pkg/profile/`
- Hooks config: `plugins/forge/hooks/hooks.json`
- Skill references: 23 skill files in `plugins/forge/skills/`

### Phase 4 Reference Update Map

Every file containing `task <command>` references must be updated. The table below enumerates each file and the specific old-to-new replacements required.

#### hooks.json (1 file)

| File | Line | Old Reference | New Reference |
|------|------|---------------|---------------|
| `plugins/forge/hooks/hooks.json` | 42 | `"task cleanup"` | `"forge cleanup"` |
| `plugins/forge/hooks/hooks.json` | 52 | `"task cleanup"` | `"forge cleanup"` |
| `plugins/forge/hooks/hooks.json` | 62 | `"task all-completed"` | `"forge quality-gate"` |

#### hooks guide (1 file)

| File | Old Reference | New Reference |
|------|---------------|---------------|
| `plugins/forge/hooks/guide.md` | `task record` | `forge task submit` |
| `plugins/forge/hooks/guide.md` | `task claim` | `forge task claim` |
| `plugins/forge/hooks/guide.md` | `task feature` | `forge feature` |

#### Agent files (1 file with task-command refs)

| File | Old Reference | New Reference |
|------|---------------|---------------|
| `plugins/forge/agents/task-executor.md` | `task prompt` | `forge prompt get-by-task-id` |
| `plugins/forge/agents/task-executor.md` | `task claim` | `forge task claim` |
| `plugins/forge/agents/task-executor.md` | `task status` | `forge task status` |

#### Command files (5 files with task-command refs)

| File | Old References | New References |
|------|----------------|----------------|
| `plugins/forge/commands/execute-task.md` | `task claim`, `task query`, `task record`, `task add`, `task template`, `task prompt` | `forge task claim`, `forge task query`, `forge task submit`, `forge task add`, (remove `task template` — deleted), `forge prompt get-by-task-id` |
| `plugins/forge/commands/run-tasks.md` | `task claim`, `task query`, `task record`, `task add`, `task template`, `task prompt` | `forge task claim`, `forge task query`, `forge task submit`, `forge task add`, (remove `task template` — deleted), `forge prompt get-by-task-id` |
| `plugins/forge/commands/git-checkout.md` | `task feature` | `forge feature` |
| `plugins/forge/commands/quick.md` | `task validate` | `forge task validate-index` |

#### Skill files (12 files with task-command refs out of 23 total)

| File | Old References | New References |
|------|----------------|----------------|
| `plugins/forge/skills/record-task/SKILL.md` | `task claim`, `task record`, `task status`, `task query` | `forge task claim`, `forge task submit`, `forge task status`, `forge task query` |
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | `task index`, `task validate`, `task add`, `task record`, `task feature`, `task template` | `forge task index`, `forge task validate-index`, `forge task add`, `forge task submit`, `forge feature`, (remove `task template` — deleted) |
| `plugins/forge/skills/quick-tasks/SKILL.md` | `task index`, `task validate`, `task add`, `task feature` | `forge task index`, `forge task validate-index`, `forge task add`, `forge feature` |
| `plugins/forge/skills/execute-task.md` (in commands/) | (covered in Command files above) | — |
| `plugins/forge/skills/run-e2e-tests/SKILL.md` | `task feature` | `forge feature` |
| `plugins/forge/skills/gen-test-scripts/SKILL.md` | `task feature` | `forge feature` |
| `plugins/forge/skills/consolidate-specs/SKILL.md` | `task feature` | `forge feature` |
| `plugins/forge/skills/init-justfile/SKILL.md` | `task all-completed`, `task claim`, `task feature` | `forge quality-gate`, `forge task claim`, `forge feature` |
| `plugins/forge/skills/forensic/SKILL.md` | No task-command refs (uses `forge forensic` pattern already) | — |
| `plugins/forge/skills/improve-harness/SKILL.md` | `task record` (in table) | `forge task submit` |
| `plugins/forge/skills/learn-lesson/SKILL.md` | No task-command refs | — |
| `plugins/forge/skills/eval-harness/SKILL.md` | No task-command refs | — |
| `plugins/forge/skills/eval-design/SKILL.md` | No task-command refs | — |
| `plugins/forge/skills/eval-prd/SKILL.md` | No task-command refs | — |
| `plugins/forge/skills/eval-ui/SKILL.md` | No task-command refs | — |
| `plugins/forge/skills/eval-proposal/SKILL.md` | No task-command refs | — |
| `plugins/forge/skills/eval-test-cases/SKILL.md` | No task-command refs | — |
| `plugins/forge/skills/eval-consistency/SKILL.md` | No task-command refs | — |
| `plugins/forge/skills/write-prd/SKILL.md` | No task-command refs | — |
| `plugins/forge/skills/ui-design/SKILL.md` | No task-command refs | — |
| `plugins/forge/skills/tech-design/SKILL.md` | No task-command refs | — |
| `plugins/forge/skills/brainstorm/SKILL.md` | No task-command refs | — |
| `plugins/forge/skills/graduate-tests/SKILL.md` | No task-command refs | — |
| `plugins/forge/skills/gen-test-cases/SKILL.md` | No task-command refs | — |

#### Doc files (4 files)

| File | Old References | New References |
|------|----------------|----------------|
| `task-cli/docs/OVERVIEW.md` | `task claim`, `task record`, `task status`, `task query`, `task feature`, `task check`, `task validate`, `task verify-completion`, `task cleanup`, `task all-completed`, `task index`, `task prompt`, `task template`, `task migrate`, `task validate-specs` | Full rename to `forge` equivalents per Naming Change Spec |
| `task-cli/docs/OVERVIEW.zh.md` | Same as above (Chinese version) | Same mapping |
| `task-cli/docs/WORKFLOW.md` | Same command set + ASCII diagram labels showing `task claim`, `task record`, `task cleanup`, `task all-completed`, `task validate` | Full rename + diagram label updates |
| `task-cli/docs/WORKFLOW.zh.md` | Same as above (Chinese version) | Same mapping |

**Summary**: 24 files total require reference updates (1 hooks.json + 1 guide.md + 1 agent file + 4 command files + 9 skill files + 4 doc files). 11 skill files require no changes. Renames not explicitly listed: every `task <subcommand>` becomes either `forge task <subcommand>` (for task-group commands) or `forge <top-level>` (for cleanup, quality-gate, verify-task-done, feature, version, probe).
