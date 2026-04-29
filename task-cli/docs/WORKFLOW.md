# claude-task-cli Key Workflows

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
│                        task claim                               │
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
                              │ Sort by Phase → Priority│
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
        if all tasks in that phase are completed:
            dependency satisfied
        else:
            dependency not satisfied
    else:                        # Exact dependency (e.g. "1.1")
        if dep task status == completed:
            dependency satisfied
        else:
            dependency not satisfied
```

---

## 2. Task Record Generation Workflow

```
┌─────────────────────────────────────────────────────────────────┐
│              task record <task-id> < input.json                 │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ Read from stdin │
                    │ JSON data       │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ Parse RecordData│
                    │ Validate        │
                    │ required fields │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ Generate from   │
                    │ template        │
                    │ Markdown content│
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ Write to        │
                    │ records/        │
                    │ <task-id>.md    │
                    └─────────────────┘
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
    "notes": "Optional notes"
}
```

---

## 3. verifyCompletion Workflow

```
┌─────────────────────────────────────────────────────────────────┐
│                   task verifyCompletion                         │
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

Note: verifyCompletion only validates status; it does not delete any files.
```

---

## 4. Cleanup Workflow

```
┌─────────────────────────────────────────────────────────────────┐
│                        task cleanup                             │
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
                    │ Task completed   │             │ Task not         │
                    │ Delete state file│             │ completed       │
                    └────────┬────────┘             │ Keep state file  │
                              │                      └─────────────────┘
                              ▼
                    ┌─────────────────┐
                    │ Delete:         │
                    │ - state.json    │
                    │ - record.json   │
                    │   (if exists)   │
                    └─────────────────┘
```

---

## 5. All-Completed Workflow

```
┌─────────────────────────────────────────────────────────────────┐
│                     task all-completed                          │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ Load index.json │
                    │ Get all tasks   │
                    └────────┬────────┘
                              │
              ┌───────────────┴───────────────┐
              │                               │
              ▼                               ▼
    ┌─────────────────┐             ┌─────────────────┐
    │ All completed   │             │ Has incomplete   │
    │ or skipped      │             │ tasks            │
    └────────┬────────┘             │ Silent exit(0)   │
              │                      └─────────────────┘
              ▼
    ┌─────────────────┐
    │ Warn if e2e     │
    │ scripts exist   │
    │ but not         │
    │ graduated       │
    └────────┬────────┘
              │
              ▼
    ┌─────────────────┐
    │ Project-wide    │
    │ unit/integration│
    │ tests           │
    └────────┬────────┘
              │
              ▼
    ┌─────────────────┐
    │ E2E regression  │
    │ (just test-e2e) │
    │ if available    │
    └────────┬────────┘
              │
    ┌─────────┴──────────┐
    │                    │
    ▼                    ▼
┌──────────┐      ┌──────────────────────────┐
│ PASS:    │      │ FAIL: Save raw output    │
│ exit 0   │      │ Block hook → Agent reads │
└──────────┘      │ raw output → task add    │
                  │ → claim fix tasks        │
                  └──────────────────────────┘
```

**Note**: Feature e2e tests are NOT run by this hook.
They are owned by T-test-3 (`run-e2e-tests` task) in the task chain.
This hook is the project health gate: unit/integration tests + regression suite.

**Test command detection order:**
1. `testCommand` field in `index.json`
2. `justfile`/`Justfile` contains `test` recipe -> `just test`
3. `Makefile` (with test: target) -> `make test`
4. `go.mod` → `go test ./...`
5. `package.json` (with scripts.test) → `npm test`
6. `pytest.ini` / `pyproject.toml` → `pytest`

---

## 6. Validation Workflow

```
┌─────────────────────────────────────────────────────────────────┐
│                      task validate [file]                       │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ Load index.json │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ JSON syntax     │
                    │ validation      │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ Required field  │
                    │ check           │
                    │ (id, title)     │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ Dependency      │
                    │ reference       │
                    │ validation      │
                    │ (refs must have │
                    │  existing IDs)  │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ Circular        │
                    │ dependency      │
                    │ detection       │
                    │ (DFS topological│
                    │  sort)          │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ File existence  │
                    │ check           │
                    │ (tasks/*.md)    │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ Output          │
                    │ validation      │
                    │ results         │
                    └─────────────────┘
```

---

## 5. Circular Dependency Detection Algorithm

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

## 7. Typical Development Workflow

### Option 1: Using Git Branch (Recommended)

```bash
# 1. Create feature branch
$ git checkout -b feature/auth-login

# 2. Claim task (auto-detects feature: auth-login)
$ task claim
> Claimed task 1.1: Implement user authentication

# 3. Develop task
# ... write code, tests ...

# 4. Generate record
$ task record 1.1 < record.json

# 5. Update status
$ task status 1.1 completed

# 6. Commit code (verifyCompletion auto-validates)
$ git commit -m "feat(auth): implement login"
> verifyCompletion: Task completed with record → commit allowed

# 7. Loop
$ task claim
> Claimed task 1.2: Implement permission check
```

### Option 2: Using Git Worktree

```bash
# 1. Create worktree (auto-detects feature)
$ git worktree add ../auth-login feature/auth-login

# 2. Work in the worktree
$ cd ../auth-login
$ task claim
> Claimed task 1.1: Implement user authentication

# 3. Develop, record, commit ...
```

### Option 3: Manual Feature Setup

```bash
# 1. Manually set feature
$ task feature auth-login

# 2. Claim task
$ task claim

# 3. Develop, record, commit ...
```

---

## 8. Error Handling Workflow

```
Error Type              Handling
─────────────────────────────────────────────────
Feature not found       Return error, suggest running: task feature <slug>
Multiple active         Return error, list active features,
Features                suggest switching
Task-state corrupted    Return error, suggest manual deletion
index.json syntax error Return detailed error location
Dependency not found    Return error, list invalid dependencies
Circular dependency     Return error, show cycle path
File not found          Return warning, does not block operation
```

---

## 9. Feature State Management

### Set Feature

```bash
$ task feature <slug>
```

Creates the `docs/features/<slug>/tasks/process/` directory as the feature's runtime state storage.

### Show Current Feature

```bash
$ task feature
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

## 10. Dynamic Task Addition Workflow

```
┌─────────────────────────────────────────────────────────────────┐
│                         task add                                │
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

**Auto-ID generation (gap-filling):**
- Scan existing tasks for `disc-*` keys
- Find the lowest unused integer N (starting from 1)
- Return `disc-{N}`

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
