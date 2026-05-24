# forge CLI Key Workflows

## 1. Feature Identification Workflow

```
┌─────────────────────────────────────────────────────────────────┐
│                   GetCurrentFeature()                            │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ Get Git context │
                    │ (worktree/branch)│
                    └────────┬────────┘
                              │
              ┌───────────────┴───────────────┐
              │                               │
              ▼                               ▼
    ┌─────────────────┐             ┌─────────────────┐
    │ Git context      │             │ No Git context   │
    │ exists: check    │             │ Scan process/    │
    │ feature dir      │             │ directory        │
    └────────┬────────┘             └────────┬────────┘
              │                               │
              ▼                               ▼
    ┌─────────────────┐             ┌─────────────────┐
    │ Exists: return   │             │ Has task-state:  │
    │ Not exists:      │             │ Return that      │
    │ create & return  │             │ feature          │
    └─────────────────┘             └─────────────────┘

Feature identification priority:
1. Git Worktree name (e.g.: feature-auth-login)
2. Git branch name (extract xxx from feature/xxx)
3. Feature in the features directory that has tasks/process/state.json
4. Only feature in the features directory that has index.json
```

### Git Branch → Feature Mapping

```
Branch Name                 → Feature Slug
─────────────────────────────────────────────
feature/auth-login         → auth-login
feat/user-registration     → user-registration
fix/null-pointer           → null-pointer
bugfix/memory-leak         → memory-leak
hotfix/security-issue      → security-issue
chore/update-deps          → update-deps
main/master/HEAD           → (ignored, fallback to directory scan)
custom-branch              → custom-branch
```

---

## 2. Task Claim Workflow

```
┌─────────────────────────────────────────────────────────────────┐
│                     forge task claim                            │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ Call            │
                    │ GetCurrentFeature│
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ Load task-state │
                    │ Check active    │
                    │ task            │
                    └────────┬────────┘
                              │
              ┌───────────────┴───────────────┐
              │                               │
              ▼                               ▼
    ┌─────────────────┐             ┌─────────────────┐
    │ Has active task  │             │ No active task   │
    │ Return it        │             │ Search next task │
    └─────────────────┘             └────────┬────────┘
                                              │
                                              ▼
                                    ┌─────────────────┐
                                    │ Load index.json │
                                    │ Get all tasks   │
                                    └────────┬────────┘
                                              │
                                              ▼
                                    ┌─────────────────┐
                                    │ Filter pending  │
                                    │ status tasks    │
                                    └────────┬────────┘
                                              │
                                              ▼
                                    ┌─────────────────┐
                                    │ Exclude tasks   │
                                    │ with unmet deps │
                                    └────────┬────────┘
                                              │
                                              ▼
                              ┌─────────────────────────┐
                              │ Sort by Priority → ID   │
                              │                         │
                              └────────────┬────────────┘
                                              │
                                              ▼
                                    ┌─────────────────┐
                                    │ Select top-ranked│
                                    │ Update status    │
                                    └────────┬────────┘
                                              │
                                              ▼
                                    ┌─────────────────┐
                                    │ Save state.json │
                                    │ to tasks/process│
                                    └─────────────────┘
```

### Dependency Check Logic

```
Check if task T's dependencies are satisfied:

for each dep in T.Dependencies:
    if dep contains ".x":           # Wildcard dependency (e.g. "1.x")
        phase = extract phase number
        if all tasks in that phase are completed OR skipped:
            dependency satisfied
        else:
            dependency not satisfied
    else:                        # Exact dependency (e.g. "1.1")
        if dep task status == completed OR skipped:
            dependency satisfied
        else:
            dependency not satisfied
```

---

## 3. Task Record Generation Workflow

```
┌─────────────────────────────────────────────────────────────────┐
│            forge task submit <task-id> --data <path>            │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ Read & parse    │
                    │ JSON data       │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ Validate        │
                    │ required fields │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────────┐
                    │ Quality gate check  │
                    │ (completed only,    │
                    │  testable type,        │
                    │  not --force):      │
                    │ just compile → fmt  │
                    │ → lint → test      │
                    └────────┬────────────┘
                             │
                 ┌───────────┴───────────┐
                 │                       │
                 ▼                       ▼
       ┌──────────────────┐    ┌──────────────────┐
       │ Gate passed or   │    │ testsFailed > 0  │
       │ skipped          │    │ Auto-downgrade:  │
       └────────┬─────────┘    │ completed→blocked│
                │              │ (non-overridable)│
                │              └────────┬─────────┘
                │                       │
                ▼                       ▼
       ┌───────────────────────────────────────┐
       │ Generate markdown from template       │
       │ Write to records/<task-id>.md         │
       │ Update index.json status              │
       └───────────────────────────────────────┘
```

### RecordData Structure

```json
{
    "taskId": "1.1",
    "status": "completed",
    "summary": "Brief description of what was done",
    "filesCreated": ["path/to/new/file.go"],
    "filesModified": ["path/to/modified/file.go"],
    "keyDecisions": ["Decision 1", "Decision 2"],
    "testsPassed": 5,
    "testsFailed": 0,
    "coverage": 85.5,
    "acceptanceCriteria": [
        {"criterion": "Feature works", "met": true}
    ],
    "notes": "Optional notes",
    "typeReclassification": {
        "originalType": "fix",
        "actualType": "cleanup",
        "reason": "flaky test, not introduced by this feature"
    },
    "referencedDocs": ["doc1.md"],
    "reviewStatus": "approved",
    "docMetrics": "50% coverage",
    "casesGenerated": 12,
    "casesEvaluated": 10,
    "scriptsCreated": ["test1.sh"],
    "testResults": "10 passed",
    "validationPassed": true,
    "issuesFound": ["issue1"],
    "gatePassed": true,
    "gateChecks": ["lint"],
    "score": 850,
    "findings": ["issue1"],
    "severity": "major",
    "passed": true
}
```

---

## 4. verify-task-done Workflow

```
┌─────────────────────────────────────────────────────────────────┐
│                   forge verify-task-done                        │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ Check if        │
                    │ task-state      │
                    │ exists          │
                    └────────┬────────┘
                              │
              ┌───────────────┴───────────────┐
              │                               │
              ▼                               ▼
    ┌─────────────────┐             ┌─────────────────┐
    │ No task-state   │             │ Has task-state   │
    │ Return success(0)│             │ Check task status│
    └─────────────────┘             └────────┬────────┘
                                              │
                              ┌───────────────┴───────────────┐
                              │                               │
                              ▼                               ▼
                    ┌─────────────────┐             ┌─────────────────┐
                    │ Task completed   │             │ Task not         │
                    │ Check record file│             │ completed       │
                    └────────┬────────┘             │ Return failure(2)│
                              │                      └─────────────────┘
              ┌───────────────┴───────────────┐
              │                               │
              ▼                               ▼
    ┌─────────────────┐             ┌─────────────────┐
    │ Has record file  │             │ No record file   │
    │ Return success(0)│             │ Return failure(2)│
    └─────────────────┘             └─────────────────┘

Note: verify-task-done only validates status; it does not delete any files.
```

---

## 5. Cleanup Workflow

```
┌─────────────────────────────────────────────────────────────────┐
│                       forge cleanup                             │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ Check if        │
                    │ task-state      │
                    │ exists          │
                    └────────┬────────┘
                              │
              ┌───────────────┴───────────────┐
              │                               │
              ▼                               ▼
    ┌─────────────────┐             ┌─────────────────┐
    │ No task-state   │             │ Has task-state   │
    │ Exit(0)         │             │ Check task status│
    └─────────────────┘             └────────┬────────┘
                                              │
                              ┌───────────────┴───────────────┐
                              │                               │
                              ▼                               ▼
                    ┌─────────────────┐             ┌─────────────────┐
                    │ Task completed/  │             │ Task not         │
                    │ blocked/suspended│             │ completed       │
                    │ /rejected        │             │ Keep state file  │
                    │ Delete state file│             │ Keep state file  │
                    └────────┬────────┘             └─────────────────┘
                              ▼
                    ┌─────────────────┐
                    │ Delete:         │
                    │ - state.json    │
                    │ - record.json   │
                    │   (if exists)   │
                    └─────────────────┘
```

---

## 6. Quality Gate Workflow

```
┌─────────────────────────────────────────────────────────────────┐
│                      forge quality-gate                         │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ FindProjectRoot │
                    │ GetCurrentFeature│
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ .forge/state.json│
                    │ allCompleted=true│
                    │ guard check      │
                    └────────┬────────┘
                              │
              ┌───────────────┴───────────────┐
              │                               │
              ▼                               ▼
    ┌─────────────────┐             ┌─────────────────┐
    │ Guard passed:   │             │ No guard or     │
    │ consume state,  │             │ incomplete tasks │
    │ load index.json │             │ Silent exit(0)   │
    └────────┬────────┘             └─────────────────┘
              ▼
    ┌─────────────────┐
    │ Warn if e2e     │
    │ scripts exist   │
    │ but not         │
    │ promoted        │
    └────────┬────────┘
              │
              ▼
    ┌─────────────────┐      ┌──────────────────────────────┐
    │ Step 1: Quality │─────►│ FAIL: addFixTask() creates   │
    │ gate            │      │ P0 fix-task, save output,    │
    │ just compile    │      │ block hook → exit(0)         │
    │ just fmt        │      └──────────────────────────────┘
    │ (non-blocking)  │
    │ just lint       │
    └────────┬────────┘
              │ pass
              ▼
    ┌─────────────────┐      ┌──────────────────────────────┐
    │ Step 2: Project │─────►│ FAIL: addFixTask() creates   │
    │ -wide tests     │      │ P0 fix-task, save output,    │
    │ just test       │      │ block hook → exit(0)         │
    └────────┬────────┘      └──────────────────────────────┘
              │ pass
              ▼
    ┌─────────────────┐
    │ Step 3: E2E     │
    │ regression      │
    │ (if test-e2e    │
    │  recipe exists) │
    │                 │
    │ 3a. just        │──────► setup FAIL: skip e2e, warn
    │     e2e-setup   │
    │                 │
    │ 3b. Server      │──────► probe FAIL: skip e2e, warn
    │     health      │
    │     probe       │
    │                 │
    │ 3c. just        │      ┌──────────────────────────────┐
    │     test-e2e    │─────►│ FAIL: addFixTask() creates   │
    │                 │      │ P0 fix-task, save output,    │
    └────────┬────────┘      │ block hook → exit(0)         │
             │               └──────────────────────────────┘
             │ pass (or e2e unavailable)
             ▼
    ┌─────────────────┐
    │ ALL PASS:       │
    │ exit(0)         │
    └─────────────────┘
```

**addFixTask()**: On failure at any gate/test step, auto-creates a P0 fix-task using
the `fix-task` template. Extracts source files from error output, saves raw output
to `tests/results/` (unit) or `tests/e2e/results/` (e2e), updates `.forge/state.json`,
and prints a hook JSON block reason so the agent can `forge task claim` the fix.

**Note**: Feature e2e tests are NOT run by this hook.
They are owned by T-test-3 (`run-e2e-tests` task) in the task chain.
This hook is the project health gate: unit/integration tests + regression suite.

**Test command detection order:**
1. `justfile`/`Justfile` contains `test` recipe -> `just test`
2. `Makefile` (with test: target) -> `make test`
3. `go.mod` → `go test ./...`
4. `package.json` (with scripts.test) → `npm test`
5. `pytest.ini` / `pyproject.toml` → `pytest`

---

## 7. Validation Workflow

```
┌─────────────────────────────────────────────────────────────────┐
│                forge task validate-index [file]                 │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ Load index.json │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ 1. JSON syntax  │
                    │ check           │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ 2. Required     │
                    │ field check     │
                    │ (id, title,     │
                    │  file, type)    │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ 3. Dependency   │
                    │ reference       │
                    │ validity        │
                    │ (incl. wildcard │
                    │  match check)   │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ 4. Circular     │
                    │ dependency      │
                    │ detection (DFS) │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ 5. Wildcard     │
                    │ self-dependency │
                    │ check           │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ 6. Gate         │
                    │ integrity       │
                    │ (gate→summary,  │
                    │  next-phase→    │
                    │  gate deps)     │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ 7. Phase order  │
                    │ validation      │
                    │ (cross-phase    │
                    │  deps exist)    │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ 8. Phase        │
                    │ summary         │
                    │ validation      │
                    │ (each phase has │
                    │  .summary task) │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ 9. Liveness     │
                    │ check:          │
                    │ orphaned, stale,│
                    │ deadlock        │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ 10. Type field  │
                    │ validation      │
                    │ (must be in     │
                    │  ValidTypes)    │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ 11. File        │
                    │ existence check │
                    │ (tasks/*.md)    │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ 12. T-test-1 /  │
                    │ T-quick-1       │
                    │ template        │
                    │ placeholder     │
                    │ check           │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ Output results  │
                    │ (PASS/FAIL)     │
                    └─────────────────┘
```

---

## 8. Circular Dependency Detection Algorithm

```go
// Depth-first search to detect cycles
func detectCycle(tasks map[string]Task) []string {
    visited := make(map[string]bool)
    recStack := make(map[string]bool)

    var cycle []string

    var dfs func(id string) bool
    dfs = func(id string) bool {
        visited[id] = true
        recStack[id] = true

        for _, dep := range tasks[id].Dependencies {
            if !visited[dep] {
                if dfs(dep) {
                    cycle = append(cycle, dep)
                    return true
                }
            } else if recStack[dep] {
                cycle = append(cycle, dep)
                return true
            }
        }

        recStack[id] = false
        return false
    }

    for id := range tasks {
        if !visited[id] {
            dfs(id)
        }
    }

    return cycle
}
```

---

## 9. Typical Development Workflow

### Option 1: Using Git Branch (Recommended)

```bash
# 1. Create feature branch
$ git checkout -b feature/auth-login

# 2. Claim task (auto-detects feature: auth-login)
$ forge task claim
> Claimed task 1.1: Implement user authentication

# 3. Develop task
# ... write code, tests ...

# 4. Generate record
$ forge task submit 1.1 --data record.json

# 5. Update status
$ forge task status 1.1 completed

# 6. Commit code (verify-task-done auto-validates)
$ git commit -m "feat(auth): implement login"
> verify-task-done: Task completed with record → commit allowed

# 7. Loop
$ forge task claim
> Claimed task 1.2: Implement permission check
```

### Option 2: Using Git Worktree (Recommended for parallel features)

The `forge worktree` command group manages git worktrees tailored to the Forge feature workflow. Forge-cli's existing `GetWorktreeName()` automatically detects the feature from the worktree name, so no manual feature setup is needed.

#### Commands

| Command | Description |
|---------|-------------|
| `forge worktree start <slug>` | Create a worktree at `../<slug>` and launch `claude` in it |
| `forge worktree list` | Show all worktrees with name, branch, path, and feature status |
| `forge worktree resume <slug>` | Re-launch `claude` in an existing worktree |
| `forge worktree remove <slug>` | Remove the worktree (branch is preserved for manual merge) |

#### Typical Workflow

```bash
# 1. Start a worktree — creates branch <slug> from HEAD, launches claude
$ forge worktree start auth-login

# 2. Inside the claude session (auto-detects feature: auth-login)
$ forge task claim
> Claimed task 1.1: Implement user authentication

# 3. Develop, record, commit ...

# 4. Later, resume an existing session if you closed it
$ forge worktree resume auth-login

# 5. List all worktrees to see status
$ forge worktree list

# 6. Remove when done (branch is kept for merge)
$ forge worktree remove auth-login
```

#### Multiple Features in Parallel

```bash
# Terminal 1
$ forge worktree start feature-a

# Terminal 2
$ forge worktree start feature-b

# Both features develop independently in separate worktrees
```

#### Manual Git Worktree (Fallback)

If `forge worktree` is not available, the equivalent manual steps are:

```bash
# 1. Create worktree
$ git worktree add ../auth-login auth-login

# 2. Work in the worktree
$ cd ../auth-login
$ forge task claim
> Claimed task 1.1: Implement user authentication

# 3. Develop, record, commit ...

# 4. Clean up
$ cd ..
$ git worktree remove auth-login
```

### Option 3: Manual Feature Setup

```bash
# 1. Manually set feature
$ forge feature auth-login

# 2. Claim task
$ forge task claim

# 3. Develop, record, commit ...
```

---

## 10. Error Handling Workflow

```
Error Type              Handling
─────────────────────────────────────────────────
Feature not found       Return error, suggest running: forge feature <slug>
Multiple active         Return error, list active features,
Features                suggest switching
Task-state corrupted    Return error, suggest manual deletion
index.json syntax error Return detailed error location
Dependency not found    Return error, list invalid dependencies
Circular dependency     Return error, show cycle path
File not found          Return warning, does not block operation
```

---

## 11. Feature State Management

### Set Feature

```bash
$ forge feature <slug>
```

Creates the `docs/features/<slug>/tasks/process/` directory as the feature's runtime state storage.

### Show Current Feature

```bash
$ forge feature
> Current feature: auth-login
```

### Feature Identification Priority

```
Priority  Source                           Example
─────────────────────────────────────────────────────────────────
1         Git Worktree                     worktrees/auth-login → auth-login
2         Git branch name                  feature/auth-login → auth-login
3         State file                       docs/features/auth-login/tasks/process/state.json
4         Unique feature directory         Used when only one feature has index.json

```

### Rules for Inferring Feature from Git

```
Branch Prefix       → Action
───────────────────────────────────
feature/            → Remove prefix
feat/               → Remove prefix
fix/                → Remove prefix
bugfix/             → Remove prefix
hotfix/             → Remove prefix
chore/              → Remove prefix
main/master/HEAD    → Ignore, use directory scan
Other               → Replace / with -
```

Examples:
- `feature/user-auth` → `user-auth`
- `custom/branch/name` → `custom-branch-name`
- `main` → use directory scan
```

## 12. Dynamic Task Addition Workflow

```
┌─────────────────────────────────────────────────────────────────┐
│                       forge task add                            │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ FindProjectRoot │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ RequireFeature  │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ Validate title  │
                    │ is non-empty    │
                    └────────┬────────┘
                              │
              ┌───────────────┴───────────────┐
              │                               │
              ▼                               ▼
    ┌─────────────────┐             ┌─────────────────┐
    │ --id provided:  │             │ --id omitted:   │
    │ validate unique │             │ auto-generate   │
    └────────┬────────┘             │ disc-N          │
             │                      └────────┬────────┘
             └──────────┬────────────────────┘
                        │
                        ▼
              ┌─────────────────┐
              │ Validate deps   │
              │ exist in index  │
              └────────┬────────┘
                        │
                        ▼
              ┌─────────────────┐
              │ Add to index    │
              │ Create .md file │
              └────────┬────────┘
                        │
                        ▼
              ┌─────────────────┐
              │ Reset forge     │
              │ state           │
              │ (allCompleted=  │
              │  false)         │
              └────────┬────────┘
                        │
                        ▼
              ┌─────────────────┐
              │ Print ADDED     │
              │ output block    │
              └─────────────────┘
```

**Auto-ID generation (max+1):**
- Without template: scan existing tasks for `disc-*` keys → `disc-{max+1}`
- With template: use template's `IDPrefix` (e.g. `fix-task` → `fix-{max+1}`)
- Returns max(existing prefix-N) + 1; starts from 1 when none exist

**Flags:**

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--title` | Yes | - | Task title |
| `--id` | No | auto `disc-N` | Custom task ID |
| `--priority` | No | P1 | P0/P1/P2 |
| `--depends-on` | No | none | Comma-separated task IDs |
| `--estimated-time` | No | - | Time estimate |
| `--breaking` | No | false | Triggers full test suite |
| `--description` | No | - | Task body content |

| `--template` | No | - | Template name (e.g. fix-task) |
| `--var` | No | - | Template variable (key=value, repeatable) |
| `--source-task-id` | No | - | Source task ID for fix-tasks |

**SourceTaskID behavior:**
- Persists `sourceTaskID` on the new Task in index.json
- Auto-adds new task as dependency of the source task (reverse dependency injection)
- Template variable `{{SOURCE_TASK_ID}}` is auto-populated

**Template defaults:**
- `fix-task`: Priority=P0, Breaking=true, EstimatedTime=30min, IDPrefix=fix
- Defaults are applied unless the corresponding flag is explicitly set

## 13. Fix-Task Lifecycle

```
Source task (in_progress)
         │
         ▼  test fails
forge task status <id> blocked
         │
         ▼
forge task add --template fix-task --source-task-id <id>
         │
         ▼  fix-N (P0, pending, auto-ID from template prefix)
   forge task claim → picks P0 first
         │
         ▼  fix-task executes
   forge task submit → fix-task completed or skipped
         │
         ▼  auto-restore checks:
   - fix-task has SourceTaskID?
   - Source is blocked?
   - ALL source deps completed or skipped?
         │
    YES  │  NO
    ┌────┘  └─── source stays blocked
    ▼           (other fix-tasks still pending)
Source → pending
         │
         ▼
   forge task claim → source re-claimed
```

**Multi-fix scenario:** When multiple fix-tasks are created for one source, the source is auto-restored only when the LAST fix-task completes.

**Nested fix-tasks:** When a fix-task itself fails, `--source-task-id` points to the FAILED fix-task (not the original source). Auto-restore chains: fix-B completes → fix-A restored → fix-A completes → source restored. Max nesting: 3 levels.

## 14. State Machine

```
                    ┌──────────┐
          claim     │ pending   │
      ┌────────────│           │◄─────────────────┐
      │            └──────────┘                   │
      │                 │                         │
      │                 │ forge task status blocked │ auto-restore
      │                 ▼                         │ (via forge task submit)
      │            ┌──────────┐                   │
      │            │ blocked   │───────────────────┘
      │            └──────────┘   (all deps completed)
      │                 │
      │                 │ (all deps completed +
      │                 │  validated by forge task status)
      │                 ▼
      │            ┌──────────┐
      ├───────────►│in_progress│
      │            └──────────┘
      │                 │
      │                 │ forge task status blocked
      │                 ▼
      │            ┌──────────┐
      │            │ blocked   │───────────────────┐
      │            └──────────┘                   │
      │                                           │
      │                 │ forge task submit         │
      │                 ▼                         │
      │            ┌──────────┐                   │
      │            │ completed │◄──────────────────┘
      │            └──────────┘  (terminal, no exit without --force)
      │
      │            ┌──────────┐
      └───────────►│ skipped   │
                   └──────────┘

      │            ┌──────────┐
      └───────────►│ rejected  │
                   └──────────┘
```

**Guards:**
- `completed → *`: Blocked (terminal state). Use `--force` to override.
- `rejected → *`: Blocked (terminal state). Use `--force` to override.
- `skipped → *`: Blocked (terminal state). Use `--force` to override.
- `in_progress → completed`: Blocked. Use `forge task submit` instead.
- `* → pending` / `* → in_progress`: Requires all dependencies to be completed or skipped.
- `--force` flag bypasses all state machine guards.

## 15. Lifecycle Liveness Validation

`forge task validate-index` detects lifecycle anomalies:
- **Orphaned**: blocked task with no dependencies
- **Stale**: blocked task whose deps are all completed or skipped (should be pending)
- **Deadlock**: blocked task whose deps are all blocked or missing (no path to resolution)

## 16. Index Build Workflow

```
┌─────────────────────────────────────────────────────────────────┐
│            forge task index --feature <slug>                    │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ Load existing   │
                    │ index.json or   │
                    │ create new      │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ Detect mode:    │
                    │ prd→breakdown   │
                    │ proposal→quick  │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ Scan tasks/     │
                    │ *.md files      │
                    │ parse frontmatter│
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ Merge: preserve │
                    │ status/         │
                    │ sourceTaskID/   │
                    │ blockedReason   │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ Generate test   │
                    │ tasks from      │
                    │ embedded        │
                    │ profiles        │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ Save index.json │
                    │ + run validate  │
                    └─────────────────┘
```

**Flags:**

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--feature` | Yes | - | Feature slug |
| `--test-profiles` | No | from config | Override test profiles (comma-separated) |

**Idempotent:** re-running produces the same output. Runtime state (status, sourceTaskID, blockedReason) is always preserved.
