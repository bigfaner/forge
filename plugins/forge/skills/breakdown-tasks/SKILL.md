---
name: breakdown-tasks
description: Use when the technical design is finalized to break down into executable tasks. Creates task files based on technical design.
---

# Breakdown Tasks

Break a technical design into executable tasks (1-4h each, clear dependencies, testable acceptance criteria).

## Prerequisites

**Conditional Tags**: `<HAS_UI>`, `<NO_UI>`, `<UI_ONLY>`, `<RULE>` are inclusion markers. If `ui/ui-design.md` exists, include `<HAS_UI>` and `<UI_ONLY>` blocks and exclude `<NO_UI>` blocks. If not, include `<NO_UI>` blocks and exclude `<HAS_UI>`/`<UI_ONLY>`/`<RULE>` blocks.

| Artifact                | Missing? Run                    |
| ----------------------- | ------------------------------- |
| `prd/prd-spec.md`       | `/write-prd`                    |
| `design/tech-design.md` | `/tech-design` → `/eval-design` |

## Step 1: Read All Documents

Read `manifest.md` to locate documents, then read all available files:

- `prd/prd-spec.md` — WHAT to build
- `design/tech-design.md` — HOW to build it
- `design/api-handbook.md` — interfaces (if exists)
- `prd/prd-user-stories.md` — user scenarios with Given/When/Then AC (if exists)
- `prd/prd-ui-functions.md` — UI function requirements (if exists)
- `ui/ui-design.md` — UI component specs (if exists)

<HAS_UI>
If `ui/ui-design.md` exists, also list `ui/prototype/` files and read `ui/prototype/index.html` for page inventory (skip if no prototype directory).
</HAS_UI>

## Step 2: Map → Tasks

### Element Mapping

| Design Element               | Source         | Task Type                |
| ---------------------------- | -------------- | ------------------------ |
| Interface definition         | tech-design.md | Interface task           |
| Data model                   | tech-design.md | Model task               |
| Backend component            | tech-design.md | Implementation (Backend) |
| Error type                   | tech-design.md | Error handling task      |
| PRD flow gate (diamond node) | prd-spec.md    | Gate verification task   |

<UI_ONLY>
| UI Component (Layout + States + Interactions + Binding) | ui/ui-design.md | Implementation (UI) |
</UI_ONLY>

<RULE>
- Each `ui/ui-design.md` Component → **one** UI task (split only if >4h)
</RULE>

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

**index.json rules:**
- Paths relative to `tasks/` directory
- `dependencies` arrays reference task IDs (`"1.1"`), not index keys (`"1.1-interface"`)
- Wildcard `"<phase>.x"` means "all tasks in phase <phase>" (resolved by task CLI, excludes .summary/.gate/self)
</HARD-RULE>

### 4a. Business Tasks

Read `templates/task.md` for task content structure. Create one task file per design element from the Element Mapping table, following dependencies from Step 3. For each task, set `breaking: true` if it modifies shared interfaces, data models, or API contracts (e.g., changing a schema column type, renaming a shared field). Additive changes are non-breaking.

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

3. Write `scope` into the task entry in `index.json`.

**Non-mixed projects**: when `init-justfile` detects a pure frontend or backend project, all tasks receive `scope: "all"` (scope distinction is irrelevant when `just project-type` does not return `"mixed"`).

**Examples**:

| Task file paths | scope | Reason |
|----------------|-------|--------|
| `ui/components/Button.tsx`, `src/styles.css` | `frontend` | All frontend paths |
| `cmd/server/main.go`, `pkg/handler/api.go` | `backend` | All backend paths |
| `ui/App.tsx`, `cmd/server/main.go` | `all` | Mixed frontend + backend |
| `docs/WORKFLOW.md`, `justfile` | `all` | Undetermined paths |
| Any task in a pure backend project | `all` | Non-mixed project, scope distinction is irrelevant |

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

Append five fixed test tasks:

- **T-test-1**: read `templates/gen-test-cases.md`, calls `/gen-sitemap` first (if `sitemap.json` missing) then `/gen-test-cases`, file `gen-test-cases.md`
- **T-test-2**: read `templates/gen-test-scripts.md`, calls `/gen-test-scripts`, depends on T-test-1, file `gen-test-scripts.md`
- **T-test-3**: read `templates/run-e2e-tests.md`, calls `/run-e2e-tests`, depends on T-test-2, file `run-e2e-tests.md`
- **T-test-4**: read `templates/graduate-tests.md`, calls `/graduate-tests`, depends on T-test-3, file `graduate-tests.md`
- **T-test-5**: read `templates/consolidate-specs.md`, calls `/consolidate-specs`, depends on T-test-4, file `consolidate-specs.md`

Replace `{{T_TEST_1_DEP}}` with the last phase's gate ID if a gate exists (e.g., `"2.gate"`), otherwise the last phase's summary ID.

**Responsibility chain:**
- T-test-1/2: generate test artifacts (documentation + scripts)
- T-test-3: execute feature e2e tests; on failure, `task add` fix tasks (P0) and mark completed
- T-test-4: verify e2e passed (check `latest.md`), then graduate scripts to `tests/e2e/`
- T-test-5: extract business rules and tech specs, user reviews and confirms integration

## Step 5: Create index.json

Read `templates/index.json` before writing. Assemble all tasks from Steps 4a–4d. Populate `dependencies` per Step 3 rules. Each task's `breaking` field should already be set from Step 4a.

Reference: [templates/index.json](templates/index.json) | Schema: [templates/index.schema.json](templates/index.schema.json)

## Step 6: Validate

```bash
task validate docs/features/<slug>/tasks/index.json
```

## Step 7: Update Manifest

Read `templates/manifest-update-tasks.md` for the traceability table format and frontmatter update instructions.

- Fill traceability table (4-column: PRD Section | Design Section | UI Component | Tasks); use "—" for UI Component when no UI
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
- [ ] UI tasks reference prototype files (if applicable)
- [ ] User Stories populated from `prd-user-stories.md`
- [ ] `index.json` ends with T-test-1 through T-test-5
- [ ] `manifest.md` updated with traceability + `status: tasks`
