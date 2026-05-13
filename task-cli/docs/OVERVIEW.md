# claude-task-cli Feature Overview

> A task management CLI tool based on the features directory structure, providing intelligent task claiming and dependency management for Claude Code workflows.

## Core Features

### 1. Intelligent Task Claiming (`task claim`)

Automatically selects the next available task based on a multi-dimensional strategy:

| Dimension | Priority Rule |
|-----------|--------------|
| Priority | P0 > P1 > P2 |
| Dependencies | Only claim tasks whose dependencies are satisfied |
| In-Progress | Automatically resume in-progress tasks |

**Dependency syntax support:**
- Exact match: `1.1`, `1.2`
- Wildcard match: `1.x` (prefix-level dependency)

### 2. Task Record Generation (`task record`)

Generate structured markdown execution records from JSON input, including:

- Task summary and status
- Time tracking
- List of created/modified files
- Key decisions
- Test results
- Acceptance criteria confirmation

**Validation rules (hard validation):**

| Condition | Error | Fix |
|-----------|-------|-----|
| `status=completed` + `testsPassed=0` + `testsFailed=0` + `coverage >= 0` | No test evidence | Run tests, or set `coverage: -1.0` |
| `status=completed` + any `acceptanceCriteria.met=false` | Unmet AC | Fix issue, or set `status: "blocked"` |
| `summary` empty or whitespace | Missing summary | Provide a summary |

Override with `--force`: `task record <id> --data record.json --force`

### 3. Status Management

| Command | Function |
|---------|----------|
| `task status <id>` | Query task status |
| `task status <id> <status>` | Update task status |
| `task query <id>` | Query task details |
| `task feature [slug]` | Set/display current feature |

**Status values:** `pending`, `in_progress`, `completed`, `blocked`, `skipped`, `rejected`

**State machine guards:**
- `completed` is terminal — cannot transition out without `--force`
- `in_progress → completed` blocked — use `task record` instead
- `pending`/`in_progress` transitions require all dependencies to be completed or skipped
- `--force` flag on `task status` bypasses all guards

**Auto-restore:** When a fix-task (task with `sourceTaskID`) is recorded as `completed` or `skipped`, the source task is automatically restored to `pending` if all its dependencies are completed or skipped.

### 4. Validation and Verification

| Command | Function |
|---------|----------|
| `task validate [file]` | Validate index.json structure and lifecycle |
| `task check` | Check all task dependencies |

**Validation rules:**
- JSON syntax check
- Required field validation
- Dependency reference validity
- Circular dependency detection
- File existence check
- Lifecycle liveness: orphaned blocked tasks, stale blocking, deadlock detection

### 5. Claude Code Integration Commands

| Command | Purpose | Function |
|---------|---------|----------|
| `task verify-completion` | PreToolUse (git commit) | Verify task completion status, block commits for incomplete tasks |
| `task cleanup` | Stop | Clean up state files for completed, blocked, or rejected tasks |
| `task all-completed` | Stop hook | Check if all tasks are completed, and if so, automatically run tests |
| `task profile` | Profile resolution | Resolve active test profile from config or project structure |
| `task index` | Index generation | Build or rebuild index.json from .md files with test task generation |
| `task prompt <id>` | Prompt synthesis | Generate agent prompt for a task based on its type; `--fix-record-missed` for recovery |
| `task template [name]` | Template viewer | List available templates or display template content |
| `task forensic` | Session analysis | Search/extract past session transcripts for deviation analysis |
| `task migrate` | Data migration | Infer type fields for all tasks in index.json |
| `task validate-specs` | Spec validation | AST validation against generated Playwright spec files |
| `task version` | Version info | Print CLI version |

**`task profile` subcommands:**

| Command | Description |
|---------|-------------|
| `task profile` | Resolve active profile(s): reads `.forge/config.yaml`, falls back to file-signal detection |
| `task profile set <name>` | Persist a profile choice to `.forge/config.yaml` |
| `task profile detect` | Run detection only (ignores existing config) |
| `task profile get <name> --<flag>` | Get profile data: `--manifest`, `--generate`, `--run`, `--graduate`, `--justfile`, `--template <file>` |
| `task profile --json` | Machine-readable JSON output |

**Profile detection signals:**

| Signal | Profile |
|--------|---------|
| `package.json` + `@playwright/test` or `playwright.config.*` | `web-playwright` |
| `go.mod` | `go-test` |
| `android/` or `ios/` directory | `maestro` |
| `pom.xml` or `build.gradle(.kts)` | `java-junit` |
| `Cargo.toml` | `rust-test` |
| `requirements.txt`/`pyproject.toml` + pytest | `pytest` |
| `package.json` without Playwright | `web-playwright` (fallback) |

**all-completed behavior:**
- Guard: only proceeds when `.forge/state.json` has `allCompleted=true`
- Any task is `pending`/`in_progress`/`blocked` → silent exit, exit 0
- No feature or no project root → silent exit, exit 0
- When all tasks done, runs a multi-step pipeline in order:
  1. **Quality gate**: `just compile → just fmt (non-blocking) → just lint`
  2. **Project-wide tests**: `just test` (or detected test command)
  3. **E2E setup**: `just e2e-setup` (if recipe exists)
  4. **Server health probe**: check e2e servers are responding before running tests
  5. **E2E regression**: `just test-e2e` (if recipe exists)

**Auto-fix task creation on failure:**
- At any step failure, `addFixTask()` creates a P0 fix-task using the `fix-task` template
- Extracts source file paths from error output, saves raw output to disk
- Prints hook JSON block reason → agent picks up fix task via `task claim`

**feature e2e tests (NOT run by this hook):**
- Feature e2e execution is owned by T-test-3 (`run-e2e-tests` task in the task chain)
- If `tests/e2e/features/<feature>/` exists but no graduation marker, hook prints a WARNING to guide migration

**e2e test script graduation model:**
- Graduation is agent-driven via T-test-4 (`graduate-tests` task) — not automatic
- T-test-4 checks `testing/results/latest.md` for PASS status before calling `/graduate-tests`
- Graduation marker: `tests/e2e/.graduated/<slug>` (YAML with schema_version, status, timestamp, source, targets, modules, testCount)
- Source scripts at `tests/e2e/features/<feature>/` are reorganized into `tests/e2e/<target>/` after graduation

**Test command auto-detection order (project-level):**
1. `testCommand` field in `index.json` (explicit configuration)
2. `justfile`/`Justfile` contains `test` recipe -> `just test`
3. `Makefile` contains `test:` target -> `make test`
4. `go.mod` exists -> `go test ./...`
5. `package.json` contains `scripts.test` -> `npm test`
6. `pytest.ini` / `pyproject.toml` exists -> `pytest`

**e2e test detection order:**
1. `justfile`/`Justfile` contains `test-e2e` recipe -> `just test-e2e`

---

## Directory Structure Convention

```
project-root/
├── docs/
│   ├── proposals/<slug>/           # /brainstorm output
│   │   └── proposal.md
│   └── features/<slug>/            # Feature workspace
│       ├── manifest.md             # Feature index & traceability mapping
│       ├── prd/
│       │   ├── prd-spec.md         # PRD Spec
│       │   ├── prd-user-stories.md # User stories
│       │   └── prd-ui-functions.md # UI function highlights (optional)
│       ├── design/
│       │   ├── tech-design.md      # Technical design
│       │   └── api-handbook.md     # API documentation
│       ├── ui/
│       │   └── ui-design.md        # UI design specification (optional)
│       ├── testing/
│       │   ├── test-cases.md      # Test cases (with target field)
│       │   └── results/
│       │       └── latest.md      # e2e test results report
│       └── tasks/
│           ├── index.json          # Task definitions
│           ├── process/            # Runtime state
│           │   ├── state.json
│           │   └── record.json
│           ├── 1.1-<title>.md     # Task details
│           └── records/            # Execution records
├── tests/
│   └── e2e/                       # Post-graduation regression test suite
│       ├── .graduated/            # Graduation marker files
│       │   └── <slug>             # YAML marker (schema_version, status, timestamp, source, targets, modules, testCount)
│       ├── ui/<page>/             # UI tests (aggregated by page)
│       │   └── ui.spec.ts
│       ├── api/<resource>/        # API tests (aggregated by resource)
│       │   └── api.spec.ts
│       └── cli/<command>/         # CLI tests (aggregated by command)
│           └── cli.spec.ts
```

### Project Root Detection

The tool automatically detects the project root, supporting multiple project types and monorepo structures:

**Detection priority** (from highest to lowest):
1. Environment variables: `CLAUDE_PROJECT_DIR` > `PROJECT_ROOT`
2. Workspace markers: `go.work`, `pnpm-workspace.yaml`, `lerna.json`, `turbo.json`, `nx.json`, `WORKSPACE`, `settings.gradle*`
3. Project markers: `go.mod`, `package.json`, `Cargo.toml`, `pyproject.toml`, `pom.xml`, `build.gradle*`
4. VCS boundary: `.git`, `.hg`

**Supported project types**:
- Go (`go.mod`, `go.work`)
- Node.js (`package.json`)
- Rust (`Cargo.toml`)
- Python (`pyproject.toml`, `setup.py`)
- Java/Maven (`pom.xml`)
- Java/Gradle (`build.gradle`, `settings.gradle`)
- Bazel (`WORKSPACE`)
- General Git repository (`.git`)

**Feature auto-detection**:
- Git worktree name -> feature slug
- Git branch name (e.g. `feature/auth-login`) -> auth-login
- Directory scan (features with `tasks/process/state.json` take priority)

**State isolation**: Each feature's runtime state is stored in its own `docs/features/<slug>/tasks/process/` directory, avoiding state conflicts across multiple features.

---

## Data Models

### Task

```go
type Task struct {
    ID            string   `json:"id"`                      // Task ID (e.g. "1.1")
    Title         string   `json:"title"`                   // Task title
    Priority      string   `json:"priority"`                // P0/P1/P2
    EstimatedTime string   `json:"estimatedTime,omitempty"` // Estimated time
    Dependencies  []string `json:"dependencies,omitempty"`  // List of dependency task IDs
    Status        string   `json:"status"`                  // pending/in_progress/completed/blocked/skipped
    File          string   `json:"file"`                    // Task file
    Record        string   `json:"record"`                  // Record file
    Breaking      bool     `json:"breaking,omitempty"`      // Global change flag; triggers full test suite on completion
    Scope         string   `json:"scope,omitempty"`         // Task scope: frontend/backend/all (default: all)
    SourceTaskID  string   `json:"sourceTaskID,omitempty"`  // ID of the task that spawned this task (e.g. fix-task -> source)
    MainSession   bool     `json:"mainSession,omitempty"`   // Task must run in main session (not dispatched to task-executor)
    NoTest        bool     `json:"noTest,omitempty"`        // Task does not require tests (e.g. documentation-only); skips quality gate and test evidence check
    Type          string   `json:"type,omitempty"`          // Task execution type; required after migration
    Profile       string   `json:"profile,omitempty"`       // Test profile name (e.g. "web-playwright"); set by task index for per-profile test tasks
    BlockedReason string   `json:"blockedReason,omitempty"` // Why this task entered blocked state; written by run-tasks when task prompt fails
}
```

**Valid task types (enforced by validation):**

| Type | Description |
|------|-------------|
| `implementation` | Standard development task |
| `fix` | Bug fix task (auto-created by all-completed on failure) |
| `gate` | Quality gate between phases |
| `doc-generation.summary` | Phase summary document generation |
| `doc-generation.consolidate` | Specification consolidation |
| `test-pipeline.gen-cases` | Test case generation |
| `test-pipeline.eval-cases` | Test case evaluation |
| `test-pipeline.gen-scripts` | Test script generation |
| `test-pipeline.run` | E2E test execution |
| `test-pipeline.graduate` | Test graduation to regression suite |
| `test-pipeline.verify-regression` | Regression verification |

### TaskIndex

```go
type TaskIndex struct {
    Feature      string          `json:"feature"`
    PRD          string          `json:"prd,omitempty"`
    Design       string          `json:"design,omitempty"`
    Proposal     string          `json:"proposal,omitempty"`
    Created      string          `json:"created,omitempty"`
    Status       string          `json:"status,omitempty"`
    Tasks        map[string]Task `json:"tasks"`
    StatusEnum   []string        `json:"statusEnum,omitempty"`
    PriorityEnum []string        `json:"priorityEnum,omitempty"`
    TestCommand  string          `json:"testCommand,omitempty"`
    E2ERound     int             `json:"e2eRound,omitempty"` // current fix-e2e round (0 = no failures yet)
}
```

---

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Language | Go 1.21 |
| CLI Framework | github.com/spf13/cobra |
| External Dependencies | Only cobra (minimal footprint) |

---

## Architecture Constraints

```
Dependency direction: cmd -> internal -> pkg (reverse strictly forbidden)
Module interaction: via interfaces/type definitions, no direct dependency on internal implementations
```

## Command Quick Reference

```bash
task claim              # Claim the next task
task record 1.1         # Generate task record
task record 1.1 --force # Generate task record (skip validation)
task add --title "Fix: ..." --priority P0 --breaking  # Add a new task dynamically
task add --template fix-task --title "Fix ..." --source-task-id 1.1 --block-source  # Add fix-task, block source
task add --title "..." --var KEY=VALUE  # Add task with template variables
task status 1.1         # Query task status
task status 1.1 done    # Update status
task query 1.1          # Query task details
task feature auth       # Switch feature
task check              # Dependency check
task validate           # Validate index.json
task index --feature <slug>  # Build index.json from .md files
task verify-completion   # Verify task completion (git commit hook)
task cleanup            # Clean up completed task state (stop hook)
```
