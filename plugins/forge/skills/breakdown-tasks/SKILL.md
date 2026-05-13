---
name: breakdown-tasks
description: Use when the technical design is finalized to break down into executable tasks. Creates task files based on technical design.
---

# Breakdown Tasks

Break a technical design into executable tasks (1-4h each, clear dependencies, testable acceptance criteria).

## Prerequisites

**Conditional Tags**: `<HAS_UI>`, `<NO_UI>`, `<UI_ONLY>`, `<HAS_PLACEMENT>`, `<RULE>` are inclusion markers.
- If `ui/ui-design.md` exists → include `<HAS_UI>` and `<UI_ONLY>` blocks, exclude `<NO_UI>` blocks.
- If `ui/ui-design.md` does NOT exist → include `<NO_UI>` blocks, exclude `<HAS_UI>`/`<UI_ONLY>` blocks.
- If `prd/prd-ui-functions.md` exists → include `<HAS_PLACEMENT>` and `<RULE>` blocks (independent of ui-design.md). `<RULE>` has no independent activation — it is always co-activated with `<HAS_PLACEMENT>`.
- If `prd/prd-ui-functions.md` does NOT exist → exclude `<HAS_PLACEMENT>` and `<RULE>` blocks.
- If `design/er-diagram.md` exists → include `<HAS_DB>` blocks.
- If `design/er-diagram.md` does NOT exist → exclude `<HAS_DB>` blocks.

| Artifact                | Missing? Run                    |
| ----------------------- | ------------------------------- |
| `prd/prd-spec.md`       | `/write-prd`                    |
| `design/tech-design.md` | `/tech-design` |

## Step 0: Resolve Profile

1. **Resolve profile**: Run `task profile` to get the active test profile(s). This reads `.forge/config.yaml`, falls back to project structure detection.
2. **On failure** (output shows `PROFILE: (none)`): ask the user to choose from known profiles (`web-playwright`, `go-test`, `maestro`, `java-junit`, `rust-test`, `pytest`). Run `task profile set <name>` to persist their choice.
3. **Load profile manifest**: Run `task profile get <profile-name> --manifest` for each resolved profile.

The resolved profiles drive per-profile task expansion in Step 4d. If only one profile is active, tasks use plain IDs (e.g., `T-test-2`). If multiple profiles are active, affected tasks are expanded with letter suffixes (e.g., `T-test-2a`, `T-test-2b`).

<HARD-RULE>
Do NOT silently default to any profile. If `task profile` returns no result and the user cannot decide, abort the skill.
</HARD-RULE>

## Step 1: Read All Documents

Read `manifest.md` to locate documents, then read all available files:

- `prd/prd-spec.md` — WHAT to build
- `design/tech-design.md` — HOW to build it
- `design/api-handbook.md` — interfaces (if exists)
- `design/er-diagram.md` — entity relationships (if exists)
- `design/schema.sql` — SQL DDL (if exists)
- `prd/prd-user-stories.md` — user scenarios with Given/When/Then AC (if exists)
- `prd/prd-ui-functions.md` — UI function requirements (if exists)
- `ui/ui-design.md` — UI component specs (if exists)

<HAS_UI>
If `ui/ui-design.md` exists, also list `ui/prototype/` files and read `ui/prototype/index.html` for page inventory (skip if no prototype directory).
</HAS_UI>

<HAS_PLACEMENT>
**Placement validation** (mandatory):
1. Read the Page Composition table from `prd/prd-ui-functions.md`
2. Check if `docs/sitemap/sitemap.json` exists. If not → WARN: `"sitemap.json not found — cannot verify existing-page routes. Run /gen-sitemap for full validation."` and proceed without route verification (skip step 3).
3. For each `existing-page:<route>` entry, verify the route exists in `docs/sitemap/sitemap.json`
4. If route not found in sitemap → ERROR: abort with message `"Route <route> not found in sitemap.json. Run /gen-sitemap first or verify the route is correct."`
5. If no Placement sections found in any UI Function → ERROR: `"Missing Placement declarations. All UI Functions must have a Placement section. Edit prd/prd-ui-functions.md to add Placement sections, or re-run /write-prd."`
</HAS_PLACEMENT>

## Step 2: Map → Tasks

### Element Mapping

| Design Element               | Source         | Task Type                |
| ---------------------------- | -------------- | ------------------------ |
| Interface definition         | tech-design.md | Interface task           |
| Data model                   | tech-design.md | Model task               |
| DB schema (er-diagram + schema.sql) | design/er-diagram.md, design/schema.sql | Schema task |
| Backend component            | tech-design.md | Implementation (Backend) |
| Error type                   | tech-design.md | Error handling task      |
| PRD flow gate (diamond node) | prd-spec.md    | Gate verification task   |

<UI_ONLY>
| UI Component (Layout + States + Interactions + Binding) | ui/ui-design.md | Implementation (UI) |
| Integration Spec (existing-page) | tech-design.md | Integration (UI) |
| Page composition (new-page) | prd-ui-functions.md Page Composition | Page Assembly task |
</UI_ONLY>

<HAS_PLACEMENT>
<RULE>
UI Task Split Rules — driven by PRD Placement:

1. For each UI Function with `placement: new-page`:
   - Create one "Build Component" task per component (existing behavior)
   - Create one "Page Assembly" task: create page file, register route, compose all components
   - Build tasks depend on interfaces + models
   - Page Assembly depends on all Build tasks for its page

2. For each UI Function with `placement: existing-page:<route>`:
   - Create one "Build Component" task (component implementation + unit tests)
   - Create one "Integrate Component" task (wire component into existing page)
   - Build task depends on interfaces + models
   - Integrate task depends on Build task
   - Integrate task's acceptance criteria MUST reference:
     a. Target page file from tech-design Integration Spec
     b. Insertion point from tech-design Integration Spec
     c. Component visible at correct position (verifiable by e2e)

3. For mixed scenarios (some new-page, some existing-page):
   - Apply rules 1 and 2 independently per UI Function

4. NO fallback to one-to-one rule. Every UI component MUST have explicit Placement.
</RULE>

**Placement format note**: The PRD template stores Placement as two separate fields (`Mode: new-page | existing-page` and `Target Page: <page route or name>`). Downstream consumers (including the rules above) use the combined canonical form: `<mode>:<target-page-value>`. For example, if Mode is `existing-page` and Target Page is `/dashboard`, the canonical placement is `existing-page:/dashboard`. If Mode is `new-page` and Target Page is `Analytics`, the canonical placement is `new-page:Analytics`.
</HAS_PLACEMENT>

### PRD Coverage Verification

Read the **PRD Coverage Map** from `tech-design.md`. Every PRD acceptance criterion must map to at least one task. UI-facing requirements → UI tasks, not generic Implementation.

Fallback: if Coverage Map is incomplete, use `prd/prd-user-stories.md` acceptance criteria directly.

### Phase & Gate Detection

Analyze the PRD and tech-design to identify the feature's natural execution phases and quality gates. This drives the phase structure in Step 3 and gate tasks in Step 4d.

**Explicit detection** (highest priority):

- Flow diagrams with diamond decision nodes (quality gates, phase transitions)
- PRD sections explicitly named "Round 1/2", "Phase 1/2", "Stage 1/2"

**Heuristic detection** (when no explicit structure is defined):
Scan `prd/prd-spec.md` and `design/tech-design.md` for these patterns:

- Sequential markers: "Round 1/2/3", "Phase/Stage 1/2/3", "第X阶段/轮", "Step 1/2/3", "第一轮/第二轮"
- Conditional transitions: "after X passes", "once X is verified", "X通过后", "确认X后再进行"
- Go/no-go checkpoints: "verify all tests pass", "confirm X before proceeding", "全部通过"
- Gated prose: "第一阶段...第二阶段...", "first pass...second pass...", "先X再Y"

**Fallback** (no phases detected): The skill will use artifact-driven decomposition in Step 3.

Collect results into a **Phase Inventory** and write it to `tasks/phase-inventory.json`:

```json
[
  {"phase": 1, "name": "...", "source": "PRD-explicit|PRD-heuristic|design|fallback", "gates": [{"afterPhase": 1, "description": "..."}]},
  {"phase": 2, "name": "...", "source": "...", "gates": []}
]
```

This file persists the planning output for cross-step reference and later review.

## Step 3: Derive Phases & Dependencies

Use the Phase Inventory from Step 2 to determine phase structure. Number phases sequentially (1.x, 2.x, ...).

### PRD-defined phases (preferred)

When the Phase Inventory contains PRD-explicit or PRD-heuristic phases, map each to a numbered phase. Within each phase, create tasks from the design elements belonging to that phase's scope. PRD-defined structure always takes priority over fixed templates.

Example: a 4-round cleanup PRD produces phases 1.x–4.x, each containing tasks for that round.

### Artifact-driven decomposition (fallback)

When the Phase Inventory source is "fallback", derive phases organically from the design artifacts:

1. List all design elements from the Element Mapping table
2. Determine dependency edges between elements (what builds on what)
3. Group into dependency layers — elements at the same depth form one phase
4. Number phases sequentially in dependency order (foundations first, consumers later)

<HAS_UI>
UI components form a natural dependency layer after data models and interfaces, but do NOT require backend implementation to be complete (can mock).
</HAS_UI>

### Dependency principles

- Tasks within the same phase: parallel unless they conflict on shared resources
- Cross-phase: a task depends on prerequisite phases' outputs
  - If a gate task exists at the boundary → depend on the gate
  - If no gate → depend on the prerequisite phase's summary or last task
    <HAS_UI>
- UI tasks depend on interfaces + models only (can mock backend; does NOT need the backend implementation phase)
  </HAS_UI>

### Task granularity

Split each design element into tasks of 1–4 hours, independently testable with clear acceptance criteria. Merge small elements (<1h combined) into one task. Split large elements (>4h) by sub-responsibility.

## Step 4: Create Task Files

<HARD-RULE>
Read the corresponding template before writing each task type.

**Naming & ID conventions:**
- Business task: file `<seq>.<sub>-<slug>.md`, ID `<seq>.<sub>`
- Phase summary: file `<phase>-summary.md`, ID `<phase>.summary`, depends on `["<phase>.x"]`
- Gate task: file `<phase>-gate.md`, ID `<phase>.gate`, `breaking: true`
- Standard test: file `<title-slug>.md` (e.g., `gen-test-cases.md`), ID `T-test-<N>`

**Gate attribution:**
- `N.gate` is phase N's exit verification gate — confirms phase N output is complete and consistent
- Depends on `N.summary` (e.g., `["1.summary"]`)
- Next phase's tasks depend on `N.gate` (explicit, not wildcard)

**Sort order within a phase** (alphabetic sub-ID): numeric < `gate` < `summary`
- e.g., `1.1` < `1.2` < `1.gate` < `1.summary`
- Execution is dependency-driven; sort order is for display only

**Frontmatter rules (propagated to index.json by `task index`):**
- `dependencies` arrays reference task IDs (`"1.1"`), not index keys (`"1.1-interface"`)
- Wildcard `"<phase>.x"` means "all tasks in phase <phase>" (resolved by task CLI, excludes .summary/.gate/self)
</HARD-RULE>

### 4a. Business Tasks

Read `templates/task.md` for task content structure. Create one task file per design element from the Element Mapping table, following dependencies from Step 3. For each task, set `breaking: true` if it modifies shared interfaces, data models, or API contracts (e.g., changing a schema column type, renaming a shared field). Additive changes to existing Go interfaces are breaking — all implementors (including mocks) must be updated.

#### Existing-Code Modification Split

When a task modifies **existing shared code** (interfaces, models, API contracts, utility functions) with multiple downstream consumers, split by dependency layers so each sub-task is independently compilable and testable:

1. **Shared artifact update** (sub-ID `<seq>.<sub>a`, `breaking: true`):
   - Apply changes to the shared artifact (interface signature, model struct, API contract, etc.)
   - Reconcile ALL downstream consumers so existing code compiles and tests pass
   - No new business logic — only signature changes + stubs/adapters
   - AC: "All existing tests compile and pass"

2. **Feature implementation** (sub-ID `<seq>.<sub>b`, depends on `<seq>.<sub>a`):
   - Implement the actual feature logic using the updated shared artifact
   - All consumers already compile; agent focuses on business logic
   - Standard acceptance criteria from the design

**When to apply**: inspect tech-design for changes to artifacts that already exist in the codebase. If the change propagates to >5 files OR spans multiple architectural layers (e.g., repository → service → handler), apply this split. Purely additive new code (new files, new interfaces) does not need splitting.

**Example** (adding methods to a shared interface with 9 mock consumers):

    # Without split (agent stalls reconciling 9 mock types × 17 methods):
    1.2-add-milestone-queries

    # With split:
    1.2a-update-main-item-repo  (interface + mock stubs, breaking: true)
    1.2b-milestone-queries      (feature logic, depends: ["1.2a"])

<HAS_DB>

For each entity in `design/er-diagram.md`, create one Schema task:
- References `design/schema.sql` and `design/er-diagram.md` as input
- AC: "DDL executes without error", "all FK references resolve", "indexes created"
- `breaking: true` if it ALTERs an existing table; `breaking: false` if all CREATE TABLE are new
- Depends on interface tasks (if any) since the migration may need type information
- scope: "backend"

</HAS_DB>

<HAS_UI>

For each UI task, **Reference Files** must include:

1. Matching `ui/ui-design.md` Component section
2. Corresponding `ui/prototype/<page>.html` (or note "No HTML prototype available")
3. Data binding table for this component
4. Relevant `tech-design.md` interfaces

Example:

```
- ui/ui-design.md Component "Dashboard" — layout, states, interactions
- ui/prototype/dashboard.html — interactive prototype
- design/tech-design.md Interfaces — data contracts
```

For **Integration tasks** (existing-page), Reference Files must include:

1. `tech-design.md` Integration Spec section
2. `ui-design.md` Component Placement section
3. Target page file path (for file-diff verification)
4. Any relevant prototype file

For **Page Assembly tasks** (new-page), Reference Files must include:

1. `prd-ui-functions.md` Page Composition table
2. `ui-design.md` Components for this page
3. Route configuration file (for route registration)
4. Navigation component file (for adding nav links)

</HAS_UI>

For each task, populate the **User Stories** section with matching stories from `prd/prd-user-stories.md`. Include full Given/When/Then acceptance criteria. If no match, note "No direct user story mapping."

### Scope Assignment

For each task, determine the `scope` field for `index.json`:

**Algorithm**: inspect the task's affected file paths (listed in the task's "Files Created/Modified" section derived from the tech-design).

1. Classify each file path:
   - `frontend`: path starts with `ui/`, `src/`, `components/`, `pages/`, `styles/`, `public/`, or any directory containing `package.json` with no `go.mod`/`Cargo.toml` at the same level
   - `backend`: path starts with `cmd/`, `internal/`, `pkg/`, `api/`, or any directory containing `go.mod`/`Cargo.toml`/`pyproject.toml` with no `package.json` at the same level
   - `undetermined`: path does not match either pattern (e.g., `docs/`, root config files, `justfile`)

2. Compute scope:
   - If ALL paths are `frontend` → `scope: "frontend"`
   - If ALL paths are `backend` → `scope: "backend"`
   - Otherwise (mixed paths, `undetermined` paths, or no file paths) → `scope: "all"`

3. Write `scope` into the task `.md` frontmatter's `scope` field.

**Non-mixed projects**: when `init-justfile` detects a pure frontend or backend project, all tasks receive `scope: "all"` (scope distinction is irrelevant when `just project-type` does not return `"mixed"`).

**Examples**:

| Task file paths | scope | Reason |
|----------------|-------|--------|
| `ui/components/Button.tsx`, `src/styles.css` | `frontend` | All frontend paths |
| `cmd/server/main.go`, `pkg/handler/api.go` | `backend` | All backend paths |
| `ui/App.tsx`, `cmd/server/main.go` | `all` | Mixed frontend + backend |
| `docs/WORKFLOW.md`, `justfile` | `all` | Undetermined paths |
| Any task in a pure backend project | `all` | Non-mixed project, scope distinction is irrelevant |

### Type Assignment

`task index` auto-infers `type` from task ID patterns via `InferType()` in `task-cli/pkg/task/infer.go`. The frontmatter `type` field takes precedence when set; otherwise the following patterns apply:

| Condition | `type` value |
|-----------|-------------|
| Task ID ends with `.summary` | `"doc-generation.summary"` |
| Task ID ends with `.gate` | `"gate"` |
| Task ID is `T-test-1` | `"test-pipeline.gen-cases"` |
| Task ID is `T-test-1b` | `"test-pipeline.eval-cases"` |
| Task ID matches `T-test-2` or `T-test-2[a-z]` | `"test-pipeline.gen-scripts"` |
| Task ID matches `T-test-3` or `T-test-3[a-z]` | `"test-pipeline.run"` |
| Task ID matches `T-test-4` or `T-test-4[a-z]` | `"test-pipeline.graduate"` |
| Task ID is `T-test-4.5` | `"test-pipeline.verify-regression"` |
| Task ID is `T-test-5` | `"doc-generation.consolidate"` |
| Task ID starts with `fix-` | `"fix"` |
| No match (fallback) | `"implementation"` |

No need to write `type` manually — `task index` handles it.

### 4b. Phase Summary Tasks

For each phase in the decomposition (from Step 3), insert a phase summary task at the end of that phase. Read `templates/phase-summary-task.md` for task content.

Example for phase 1:

```
"1.summary": {
  "id": "1.summary",
  "title": "Phase 1 Summary",
  "priority": "P0",
  "estimatedTime": "15min",
  "dependencies": ["1.x"],
  "status": "pending",
  "file": "1-summary.md",
  "record": "records/1-summary.md"
}
```

### 4c. Gate Tasks

Create a gate for every phase (including the last) when EITHER condition is met:

1. **Cross-layer**: The feature spans multiple layers (detected from the Cross-Layer Data Map or architecture diagram)
2. **PRD-defined phases**: The Phase Inventory (from Step 2) contains detected quality gates between phases

Read `templates/gate-task.md` for task content. `N.gate` is phase N's exit verification — it confirms phase N's output is complete. It depends on `N.summary`, and the next phase's tasks depend on `N.gate`. The last phase's gate verifies final output before T-test tasks begin.

Example dependency chain:

```
Phase 1: 1.1, 1.2                 (dependencies: none or earlier phases)
Phase 1 summary: 1.summary         (dependencies: ["1.x"])
Phase 1 gate: 1.gate               (dependencies: ["1.summary"])
Phase 2: 2.1, 2.2                  (dependencies: ["1.gate"])
Phase 2 summary: 2.summary         (dependencies: ["2.x"])
Phase 2 gate: 2.gate               (dependencies: ["2.summary"])
```

### 4d. Standard Test Tasks

Test tasks are auto-generated by `task index` based on the profiles resolved in Step 0. **Do NOT create test task `.md` files manually** — `task index` handles the following:

**Profile-suffix convention**: When multiple profiles are active, per-profile tasks use a lowercase letter suffix (a, b, c, ...) matching the order listed in `.forge/config.yaml`'s `test-profiles` array. Single profile produces no suffix.

#### Shared tasks (always single, regardless of profile count)

- **T-test-1**: calls `/gen-sitemap` first (if `sitemap.json` missing) then `/gen-test-cases`
- **T-test-1b**: calls `/eval-test-cases`, depends on T-test-1, main session task

#### Per-profile tasks (expanded when multiple profiles active)

| Base ID | Skill | Depends on |
|---------|-------|------------|
| T-test-2 | `/gen-test-scripts` | T-test-1b |
| T-test-3 | `/run-e2e-tests` | T-test-2\* |
| T-test-4 | `/graduate-tests` | T-test-3\* |

\* Per-profile: T-test-3a depends on T-test-2a, T-test-3b depends on T-test-2b, etc.

#### Shared tasks (continued)

- **T-test-4.5**: runs full e2e regression across all profiles, depends on ALL T-test-4 variants
- **T-test-5**: calls `/consolidate-specs`, depends on T-test-4.5

**Responsibility chain:**
- T-test-1: generate test case documentation (shared)
- T-test-1b: evaluate test cases for downstream executability, main session task (shared)
- T-test-2[letter]: generate test scripts from evaluated test cases (per profile)
- T-test-3[letter]: execute feature e2e tests; on failure, mark blocked, add fix tasks (P0) with unblock instruction — re-runs after fix (per profile)
- T-test-4[letter]: verify e2e passed (check `latest.md`), then graduate scripts to regression suite (per profile)
- T-test-4.5: run full regression suite across all profiles; on failure, mark blocked, add fix tasks (P0) with unblock instruction — re-runs after fix (shared)
- T-test-5: extract business rules and tech specs, user reviews and confirms integration (shared)

**Fix-task reference**: Templates are managed by task-cli and embedded in the binary. Auto-generated fix-task IDs follow the `fix-N` format (e.g., `fix-1`, `fix-2`). Agents should run `task template fix-task` to view the template and required variables before creating fix tasks:

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

When a fix-task completes, `task record` auto-restores the source task to `pending` (checks all source task's dependencies are completed). For nested fix-tasks (fix-task itself fails), `--source-task-id` must point to the FAILED fix-task, not the original source.

## Step 5: Generate index.json via CLI

After all business task `.md` files (Steps 4a–4c) are written, run:

```bash
task index --feature <slug>
```

This command:
1. Scans all `.md` files in `tasks/`, parses YAML frontmatter (id, title, priority, dependencies, scope, type, breaking, noTest, mainSession)
2. Auto-generates test task `.md` files from embedded profiles (Step 4d tasks)
3. Produces `index.json` with all business + test tasks
4. Runs validation automatically

If the profile was not set in Step 0, pass it explicitly: `task index --feature <slug> --test-profiles <p1>,<p2>`.

## Step 6: Validate

```bash
task validate docs/features/<slug>/tasks/index.json
```

## Step 7: Update Manifest

Read `templates/manifest-update-tasks.md` for the traceability table format and frontmatter update instructions.

- Fill traceability table (5-column: PRD Section | Design Section | UI Component | Placement | Tasks); use "—" for UI Component when no UI, use "—" for Placement when no UI Functions
- Advance status to `tasks`

## Output Checklist

- [ ] `tasks/phase-inventory.json` written with detected phases and gates
- [ ] All task files follow naming conventions from HARD-RULE
- [ ] `index.json` valid per schema, `task validate` passes
- [ ] Every PRD AC covered by ≥1 task
- [ ] Dependency graph is a DAG (no cycles) — verify with `task validate`
- [ ] Every Phase Inventory gate has a corresponding gate task
- [ ] Gate tasks: correct phase attribution, `breaking: true`, explicit dependency chains
- [ ] `breaking: true` set on tasks that modify shared contracts
- [ ] Tasks modifying existing shared code (>5 downstream files or cross-layer) split into artifact-update + feature sub-tasks
- [ ] UI tasks reference prototype files (if applicable)
- [ ] User Stories populated from `prd-user-stories.md`
- [ ] `index.json` ends with test tasks matching the profile expansion from Step 0 (shared T-test-1, T-test-1b, per-profile T-test-2/3/4, shared T-test-4.5, T-test-5)
- [ ] `manifest.md` updated with traceability + `status: tasks`
