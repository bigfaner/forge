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

**Status values:** `pending`, `in_progress`, `completed`, `blocked`, `skipped`

### 4. Validation and Verification

| Command | Function |
|---------|----------|
| `task validate [file]` | Validate index.json structure |
| `task check` | Check all task dependencies |

**Validation rules:**
- JSON syntax check
- Required field validation
- Dependency reference validity
- Circular dependency detection
- File existence check

### 5. Claude Code Integration Commands

| Command | Purpose | Function |
|---------|---------|----------|
| `task verifyCompletion` | PreToolUse (git commit) | Verify task completion status, block commits for incomplete tasks |
| `task cleanup` | Stop | Clean up state files for completed tasks |
| `task all-completed` | Stop hook | Check if all tasks are completed, and if so, automatically run tests |

**all-completed behavior:**
- All tasks are `completed` or `skipped` -> run feature e2e tests + project-level tests, exit 0
- Any task is `pending`/`in_progress`/`blocked` -> silent exit, exit 1
- No feature or no project root -> silent exit, exit 1

**e2e test failure recovery:**
- When e2e tests fail, automatically append a `fix-e2e-N` task to `index.json` (N starts from 1)
- fix-e2e task format:
  - id: `fix-e2e-N`
  - title: "Fix e2e test failures"
  - priority: `P0`
  - file: `testing/results/latest.md` (points to failure details)
- If a pending fix-e2e task already exists, skip appending (avoid duplicates)
- fix-e2e task limit is 3; after exceeding, print a warning and exit 0 (avoid infinite loops)
- After appending, exit 1, triggering the agent to continue working and claim the fix-e2e task

**e2e test script graduation model:**
- On first successful e2e test, migrate scripts to `tests/e2e/<type>/<target>/` based on the test case's `target` field
- Graduation marker: `tests/e2e/.graduated/<slug>` (content is a timestamp)
- If graduation marker already exists, skip migration (not first success)
- `docs/features/<slug>/testing/scripts/` is preserved (as traceability record)

**Test command auto-detection order (project-level):**
1. `testCommand` field in `index.json` (explicit configuration)
2. `justfile`/`Justfile` contains `test` recipe -> `just test`
3. `Makefile` contains `test:` target -> `make test`
4. `go.mod` exists -> `go test ./...`
5. `package.json` contains `scripts.test` -> `npm test`
6. `pytest.ini` / `pyproject.toml` exists -> `pytest`

**e2e test detection order:**
1. `justfile`/`Justfile` contains `test-e2e` recipe -> `just test-e2e`
2. `Makefile` contains `test-e2e:` target -> `make test-e2e`
3. `testing/scripts/package.json` exists -> `npm run test:all --if-present`

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
│       │   ├── scripts/           # Development-phase test scripts
│       │   │   ├── ui.spec.ts
│       │   │   ├── api.spec.ts
│       │   │   └── cli.spec.ts
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
│       │   └── <slug>             # Timestamp
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
}
```

### TaskIndex

```go
type TaskIndex struct {
    Feature      string          `json:"feature"`
    PRD          string          `json:"prd,omitempty"`
    Design       string          `json:"design,omitempty"`
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
task status 1.1         # Query task status
task status 1.1 done    # Update status
task query 1.1          # Query task details
task feature auth       # Switch feature
task check              # Dependency check
task validate           # Validate index.json
task verifyCompletion   # Verify task completion (git commit hook)
task cleanup            # Clean up completed task state (stop hook)
```
