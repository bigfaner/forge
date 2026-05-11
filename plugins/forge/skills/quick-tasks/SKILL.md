---
name: quick-tasks
description: Use for small features (1-2h, 1-4 tasks) to generate tasks directly from proposal. No PRD or design needed. Supports --no-test to skip test tasks.
---

# Quick Tasks

Generate executable tasks directly from a proposal document. For small features (1-2h total, 1-4 tasks) that don't need PRD or tech design.

## Prerequisites

| Artifact | Missing? Run |
|----------|-------------|
| `docs/proposals/<slug>/proposal.md` | `/brainstorm` or `/quick` |

<HARD-GATE>
Maximum 4 business tasks. If the proposal requires more, STOP and recommend the full pipeline: `/write-prd` → `/tech-design` → `/breakdown-tasks`.
</HARD-GATE>

## Flags

- `--no-test`: Skip T-quick-1~5 test tasks. Use for non-code proposals or when tests are handled separately.

## Step 1: Read Proposal

Determine the feature slug from the proposal directory name. Read `docs/proposals/<slug>/proposal.md` — the sole input document. Extract:

- **Problem** → task context and motivation
- **Proposed Solution** → task scope and boundaries
- **Scope > In Scope** → one task per bullet (split if >2h, merge if <30min)
- **Success Criteria** → acceptance criteria for each task
- **Key Risks** → implementation notes and risk mitigations

<HARD-RULE>
Enforce maximum 4 business tasks. If the In Scope section implies >4 tasks, STOP and recommend the full pipeline (`/write-prd` → `/tech-design` → `/breakdown-tasks`).
</HARD-RULE>

## Step 2: Derive Tasks

For each In Scope bullet:

1. Estimate effort (1-2h default)
2. Derive acceptance criteria from matching Success Criteria items
3. Determine affected file paths from the solution description
4. Set scope via Scope Assignment algorithm (same as breakdown-tasks — inspect file paths)
5. Set `breaking: true` if task modifies shared interfaces, data models, or API contracts

**Dependencies**: Linear chain (task 2 depends on task 1, etc.) unless proposal implies parallel work. Use simple integer IDs: `1`, `2`, `3`, `4`.

**Scope Assignment** (reuses breakdown-tasks algorithm):

1. Classify each affected file path:
   - `frontend`: path starts with `ui/`, `src/`, `components/`, `pages/`, `styles/`, `public/`
   - `backend`: path starts with `cmd/`, `internal/`, `pkg/`, `api/`
   - `undetermined`: other paths (`docs/`, root config, `justfile`)
2. Compute scope:
   - ALL paths `frontend` → `scope: "frontend"`
   - ALL paths `backend` → `scope: "backend"`
   - Otherwise → `scope: "all"`

## Step 3: Create Task Files

Read `templates/task.md` for the task content structure. Create one task file per derived task in `docs/features/<slug>/tasks/`.

<HARD-RULE>
Naming & ID conventions:
- Business task: file `<seq>-<slug>.md`, ID `<seq>` (e.g., file `1-add-command.md`, ID `1`)
- Quick test: file `quick-<name>.md`, ID `T-quick-<N>`
- No phase prefixes, no sub-IDs, no summary/gate tasks
</HARD-RULE>

For each task:
- Fill Description from proposal's Problem and Solution context
- Fill Affected Files from the solution description (Create/Modify/Delete tables)
- Fill Acceptance Criteria from matching Success Criteria items
- Fill Implementation Notes from Key Risks and solution details
- Set `breaking: true` only when modifying shared interfaces/models/APIs

### Type Assignment

For each task, set the `type` field in `index.json` using the following rules. These rules mirror `InferType()` in `task-cli/pkg/prompt/prompt.go` — both must stay in sync.

| Condition | `type` value |
|-----------|-------------|
| Task ID ends with `.summary` | `"doc-generation.summary"` |
| Task ID ends with `.gate` | `"gate"` |
| Task ID is `T-test-1` | `"test-pipeline.gen-cases"` |
| Task ID is `T-test-1b` | `"test-pipeline.eval-cases"` |
| Task ID is `T-test-2` | `"test-pipeline.gen-scripts"` |
| Task ID is `T-test-3` | `"test-pipeline.run"` |
| Task ID is `T-test-4` | `"test-pipeline.graduate"` |
| Task ID is `T-test-4.5` | `"test-pipeline.verify-regression"` |
| Task ID is `T-test-5` | `"doc-generation.consolidate"` |
| Task ID starts with `fix-` or `disc-` | `"fix"` |
| No match (fallback) | `"implementation"` — emit warning: `warn: task <ID> type could not be inferred, defaulting to implementation` |

Write `type` into the task entry in `index.json` alongside `scope`. For quick-tasks, business tasks (IDs `1`–`4`) and `T-quick-*` test tasks all fall through to the fallback `"implementation"` — no warning needed for these expected cases.

## Step 4: Create Test Tasks (unless --no-test)

If `--no-test` flag is NOT set, append five test tasks. Read each template before writing:

- **T-quick-1**: read `templates/quick-test-cases.md`, generates test cases from proposal's Success Criteria
- **T-quick-2**: read `templates/quick-gen-scripts.md`, generates e2e test scripts from test cases
- **T-quick-3**: read `templates/quick-run-tests.md`, runs feature e2e tests
- **T-quick-4**: read `templates/quick-graduate.md`, graduates scripts to regression suite
- **T-quick-5**: read `templates/quick-verify-regression.md`, runs full e2e regression

Replace `{{T_QUICK_1_DEP}}` with the last business task ID (e.g., `"2"` if 2 business tasks).

**Responsibility chain:**
- T-quick-1: generate test case documentation from proposal Success Criteria (no sitemap, no eval)
- T-quick-2: generate test scripts from test cases
- T-quick-3: execute feature e2e tests; on failure, mark blocked, add fix tasks (P0)
- T-quick-4: graduate scripts to `tests/e2e/`
- T-quick-5: run full regression suite; on failure, mark blocked, add fix tasks (P0)

**Fix-task reference**:

```bash
task add --template fix-task --title "Fix: <description>" \
  --source-task-id <source-task-id> \
  --block-source \
  --var SOURCE_FILES="<affected paths>" \
  --var TEST_SCRIPT="<failing test>" \
  --var TEST_RESULTS="<results path>" \
  --description "<root cause>"
```

**`--block-source`**: atomically sets source task to blocked before resolution. `task add` automatically deduplicates — check output: `ACTION: ADDED` (new fix task) or `ACTION: SKIPPED` (active fix already exists).

## Step 5: Create index.json

Read `templates/index.json` before writing. Assemble all tasks from Steps 3-4.

<HARD-RULE>
index.json rules:
- Paths relative to `tasks/` directory
- `dependencies` arrays reference task IDs (`"1"`, `"T-quick-1"`)
- `proposal` field points to the proposal path (relative to feature dir)
- Copy all boolean flags from the task template's YAML frontmatter (`breaking`, `noTest`, `mainSession`) directly into the corresponding index.json entry
- If a quick task needs to spawn subagents (unlikely in quick mode), set `mainSession: true` and add `## Main Session Instructions` to the task file
</HARD-RULE>

Reference: [templates/index.json](templates/index.json) | Schema: [templates/index.schema.json](templates/index.schema.json)

## Step 6: Create Manifest

Read `templates/manifest-quick.md` for the simplified manifest format. Write to `docs/features/<slug>/manifest.md`.

The quick manifest contains:
- Documents table: proposal path + optional test-cases path
- Tasks table: ID, Title, Status for all tasks (business + test)
- No Traceability table

## Step 7: Validate

```bash
task validate docs/features/<slug>/tasks/index.json
```

## Output Checklist

- [ ] `docs/features/<slug>/tasks/` contains 1-4 business task files
- [ ] `index.json` valid per schema, `task validate` passes
- [ ] Every Success Criterion covered by ≥1 task
- [ ] Dependency graph is a DAG (no cycles)
- [ ] Each task file includes `## Affected Files` section with Create/Modify/Delete
- [ ] (if not --no-test) T-quick-1~5 appended with correct dependency chain
- [ ] `docs/features/<slug>/manifest.md` written with `mode: quick`

## Integration

Works well with:
- `/brainstorm` — Generate the proposal before running quick-tasks
- `/quick` — Full pipeline: brainstorm → quick-tasks → run-tasks
- `/run-tasks` — Execute generated tasks (index.json compatible)
