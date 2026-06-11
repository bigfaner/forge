# forge CLI Feature Overview

> A task management CLI tool based on the features directory structure, providing intelligent task claiming and dependency management for Claude Code workflows.

## Core Features

### 1. Intelligent Task Claiming (`forge task claim`)

Automatically selects the next available task based on a multi-dimensional strategy:

| Dimension | Priority Rule |
|-----------|--------------|
| Priority | P0 > P1 > P2 |
| Dependencies | Only claim tasks whose dependencies are satisfied |
| In-Progress | Automatically resume in-progress tasks |

**Dependency syntax support:**
- Exact match: `1.1`, `1.2`
- Wildcard match: `1.x` (prefix-level dependency)

### 2. Task Record Generation (`forge task submit`)

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

Override with `--force`: `forge task submit <id> --data record.json --force`

### 3. Status Management

| Command | Function |
|---------|----------|
| `forge task status <id>` | Query task status |
| `forge task status <id> <status>` | Update task status |
| `forge task query <id>` | Query task details |
| `forge feature [slug]` | Set/display current feature |

**Status values:** `pending`, `in_progress`, `completed`, `blocked`, `skipped`, `rejected`

**State machine guards:**
- `completed` is terminal — cannot transition out without `--force`
- `in_progress → completed` blocked — use `forge task submit` instead
- `pending`/`in_progress` transitions require all dependencies to be completed or skipped
- `--force` flag on `forge task status` bypasses all guards

**Auto-restore:** When a fix-task (task with `sourceTaskID`) is recorded as `completed` or `skipped`, the source task is automatically restored to `pending` if all its dependencies are completed or skipped.

### 4. Validation and Verification

| Command | Function |
|---------|----------|
| `forge task validate-index [file]` | Validate index.json structure and lifecycle |
| `forge task check-deps` | Check all task dependencies |

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
| `forge verify-task-done` | PreToolUse (git commit) | Verify task completion status, block commits for incomplete tasks |
| `forge cleanup` | Stop | Clean up state files for completed, blocked, or rejected tasks |
| `forge quality-gate` | Stop hook | Check if all tasks are completed, and if so, automatically run tests |
| `forge task index` | Index generation | Build or rebuild index.json from .md files with test task generation |
| `forge prompt get-by-task-id <id>` | Prompt synthesis | Generate agent prompt for a task based on its type; `--fix-record-missed` for recovery |
| `forge forensic` | Session analysis | Search/extract past session transcripts for deviation analysis |
| `forge task migrate` | Data migration | Infer type fields for all tasks in index.json |
| `forge version` | Version info | Print CLI version |

**Surface type detection signals:**

| Signal | Surface Type |
|--------|-------------|
| `package.json` + `@playwright/test` or `playwright.config.*` | `web` |
| `go.mod` | `cli` |
| `android/` or `ios/` directory | `mobile` |
| `pom.xml` or `build.gradle(.kts)` | `api` |
| `Cargo.toml` | `cli` |
| `requirements.txt`/`pyproject.toml` + pytest | `cli` |
| `package.json` without Playwright | `web` (fallback) |

**all-completed behavior:**
- Guard: only proceeds when `.forge/state.json` has `allCompleted=true`
- Any task is `pending`/`in_progress`/`blocked` → silent exit, exit 0
- No feature or no project root → silent exit, exit 0
- When all tasks done, runs a multi-step pipeline in order:
  1. **Quality gate**: `just compile → just fmt (non-blocking) → just lint`
  2. **Project-wide unit tests**: `just test` (or detected test command)
  3. **Full test regression**: `just test` (test scripts in `tests/`)

**Auto-fix task creation on failure:**
- At any step failure, `addFixTask()` creates a P0 fix-task using the `fix-task` template
- Extracts source file paths from error output, saves raw output to disk
- Prints hook JSON block reason → agent picks up fix task via `forge task claim`

**Feature-specific tests (NOT run by this hook):**
- Feature test execution is owned by the test run task in the task chain
- If scripts exist in the feature's tests directory but lack promotion markers, hook prints a WARNING to guide migration

**Test tag-based promotion model:**
- Promotion is done via the test pipeline — after tests pass, `@feature` tags are replaced with `@regression`
- Tag-based lifecycle: `@feature` (newly generated) -> `@regression` (verified, promoted)
- Source scripts are organized under `tests/<journey>/` and promoted via tag replacement

**Test command auto-detection order (project-level):**
1. `justfile`/`Justfile` contains `test` recipe -> `just test`
2. `Makefile` contains `test:` target -> `make test`
3. `go.mod` exists -> `go test ./...`
4. `package.json` contains `scripts.test` -> `npm test`
5. `pytest.ini` / `pyproject.toml` exists -> `pytest`

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
│       │       └── latest.md      # test results report
│       └── tasks/
│           ├── index.json          # Task definitions
│           ├── process/            # Runtime state
│           │   ├── state.json
│           │   └── record.json
│           ├── 1.1-<title>.md     # Task details
│           └── records/            # Execution records
├── tests/
│   ├── results/                   # Test results output
│   │   └── raw-output.txt
│   ├── ui/<page>/                 # UI tests (aggregated by page)
│   │   └── ui.spec.ts
│   ├── api/<resource>/            # API tests (aggregated by resource)
│   │   └── api.spec.ts
│   └── cli/<command>/             # CLI tests (aggregated by command)
│       └── cli.spec.ts
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
    Status        string   `json:"status"`                  // pending/in_progress/completed/blocked/suspended/skipped/rejected
    File          string   `json:"file"`                    // Task file
    Record        string   `json:"record"`                  // Record file
    Breaking      bool     `json:"breaking,omitempty"`      // Global change flag; triggers full test suite on completion
    SurfaceKey    string   `json:"surface-key,omitempty"`   // User-defined surface identifier (e.g. "admin-panel")
    SurfaceType   string   `json:"surface-type,omitempty"`  // Surface type enumeration (e.g. "web", "api", "cli")
    SourceTaskID  string   `json:"sourceTaskID,omitempty"`  // ID of the task that spawned this task (e.g. fix-task -> source)
    MainSession   bool     `json:"mainSession,omitempty"`   // Task must run in main session (not dispatched to task-executor)
    Type          string   `json:"type,omitempty"`          // Task execution type; required after migration
    BlockedReason string   `json:"blockedReason,omitempty"` // Why this task entered blocked state; written by run-tasks when task prompt fails
    Coverage      *int     `json:"coverage,omitempty"`      // Per-task coverage override from frontmatter; nil = use global default
    Complexity    string   `json:"complexity,omitempty"`    // Task complexity level: "low", "medium", or "high"; empty defaults to "medium"
    Scope         string   `json:"scope,omitempty"`         // Legacy scope field (deprecated, retained for migration detection)
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
| `test.run` | Run test scripts and collect results |
| `test.gen-journeys` | Generate test journeys from specs |
| `test.gen-contracts` | Generate test contracts from journeys |
| `test.gen-scripts` | Generate executable test scripts |

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
forge task claim              # Claim the next task
forge task submit 1.1         # Generate task record
forge task submit 1.1 --force # Generate task record (skip validation)
forge task add --title "Fix: ..." --priority P0 --breaking  # Add a new task dynamically
forge task add --type coding.fix --title "Fix ..." --source-task-id 1.1 --block-source  # Add fix-task, block source
forge task add --title "..." --var KEY=VALUE  # Add task with template variables
forge task status 1.1         # Query task status
forge task status 1.1 done    # Update status
forge task query 1.1          # Query task details
forge feature auth            # Switch feature
forge task check-deps         # Dependency check
forge task validate-index     # Validate index.json
forge task index --feature <slug>  # Build index.json from .md files
forge verify-task-done        # Verify task completion (git commit hook)
forge cleanup                 # Clean up completed task state (stop hook)
```
