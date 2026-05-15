---
name: quick-tasks
description: Use for features (1-10 tasks) to generate tasks directly from proposal. No PRD or design needed.
---

# Quick Tasks

Generate executable tasks directly from a proposal document. For features (1-10 tasks) that don't need PRD or tech design.

## Prerequisites

| Artifact | Missing? Run |
|----------|-------------|
| `docs/proposals/<slug>/proposal.md` | `/brainstorm` or `/quick` |

<HARD-GATE>
Maximum 10 business tasks. If the proposal requires more, STOP and recommend the full pipeline: `/write-prd` → `/tech-design` → `/breakdown-tasks`.
</HARD-GATE>

## Step 0: Resolve Profile

1. **Resolve profile**: Run `forge profile` to get the active test profile(s). This reads `.forge/config.yaml`, falls back to project structure detection.
2. **On failure** (output shows `PROFILE: (none)`): ask the user to choose from known profiles (`web-playwright`, `go-test`, `maestro`, `java-junit`, `rust-test`, `pytest`). Run `forge profile set <name>` to persist their choice.
3. **Load profile manifest**: Run `forge profile get <profile-name> --manifest` for each resolved profile.

**Profile resolution outcome**:
- **Single profile**: one active profile (default behavior, no per-profile suffixing needed)
- **Multiple profiles**: two or more active profiles (triggers per-profile task suffixing in Step 4)

<HARD-RULE>
Do NOT silently default to any profile. If `forge profile` returns no result and the user cannot decide, abort the skill.
</HARD-RULE>

## Step 1: Read Proposal

Determine the feature slug from the proposal directory name. Read `docs/proposals/<slug>/proposal.md` — the sole input document. Extract:

- **Problem** → task context and motivation
- **Proposed Solution** → task scope and boundaries
- **Scope > In Scope** → one task per bullet (split if >2h, merge if <30min)
- **Success Criteria** → acceptance criteria for each task
- **Key Risks** → implementation notes and risk mitigations

<HARD-RULE>
Enforce maximum 10 business tasks. If the In Scope section implies >10 tasks, STOP and recommend the full pipeline (`/write-prd` → `/tech-design` → `/breakdown-tasks`).
</HARD-RULE>

## Step 2: Derive Tasks

For each In Scope bullet:

1. Estimate effort (1-2h default)
2. Derive acceptance criteria from matching Success Criteria items
3. Classify task type by description semantics (see Template Selection in Step 3)
4. For implementation tasks: set scope via Scope Inference and set `breaking: true` if task modifies shared interfaces, data models, or API contracts
5. Fill Reference Files with directory-level paths inferred from the proposal's In Scope context (e.g., proposal says "Add --type argument to gen-test-scripts" → `plugins/forge/skills/gen-test-scripts/`)

**Split by functional steps**: If one In Scope bullet contains multiple independently verifiable functional steps, split into separate tasks. Each task is an independently verifiable functional unit. Total tasks must still be ≤ 10.

**Dependencies**: Linear chain (task 2 depends on task 1, etc.) unless proposal implies parallel work. Use simple integer IDs: `1`, `2`, `3`, `4`.

**Scope Inference** (from task description semantics):

- Description mentions UI, pages, components, styles → `scope: "frontend"`
- Description mentions API, server, database, CLI → `scope: "backend"`
- Mixed or unclear → `scope: "all"`

## Step 3: Create Task Files

Read the appropriate template (see Template Selection below) for the task content structure. Create one task file per derived task in `docs/features/<slug>/tasks/`.

<HARD-RULE>
Naming & ID conventions:
- Business task: file `<seq>-<slug>.md`, ID `<seq>` (e.g., file `1-add-command.md`, ID `1`)
- Quick test: file `quick-<name>.md`, ID `T-quick-<N>`
- No phase prefixes, no sub-IDs, no summary/gate tasks
</HARD-RULE>

For each task:
- Fill Description from proposal's Problem and Solution context
- Fill Acceptance Criteria from matching Success Criteria items
- Fill Implementation Notes from Key Risks and solution details
- Fill Hard Rules only when the task has critical constraints the agent must not override (must use specific justfile recipe, commands have hidden dependencies like env vars/server lifecycle, or explicit file scope restrictions). Leave empty for normal tasks.
- For implementation tasks: set `breaking: true` when modifying shared interfaces/models/APIs
- Pass resolved profile name(s) as context when the task involves test generation or execution (see Step 4 for per-profile template variable usage)

### Template Selection

Choose the task template based on task description:

| Condition | Template |
|-----------|----------|
| Task produces only documentation, specs, or templates (non-compilable, non-runnable) | `templates/task-doc.md` (type: `"documentation"`) |
| Task modifies or creates source code, build configs, or runtime configs | `templates/task.md` (type: `"implementation"`) |

## Step 4: Test Tasks (auto-generated)

Test tasks are auto-generated by `forge task index` based on the profiles resolved in Step 0. **Do NOT create test task `.md` files manually.**

**Responsibility chain (reference for task-executor agents):**

- T-quick-1[letter]: generate test case documentation from proposal Success Criteria (no sitemap, no eval)
- T-quick-2[letter]: generate test scripts from test cases
- T-quick-3[letter]: execute feature e2e tests; on failure, mark blocked, add fix tasks (P0)
- T-quick-4[letter]: graduate scripts to `tests/e2e/`
- T-quick-5: run full regression suite; on failure, mark blocked, add fix tasks (P0)

> **Note**: Quick pipeline intentionally skips eval-test-cases (T-test-1b in the full pipeline) and consolidate-specs (T-test-5 in the full pipeline). T-quick-5 corresponds to T-test-4.5 (verify-regression) in the full pipeline, not T-test-5.

**Fix-task reference** (applies to both single and multiple profile modes):

```bash
forge task add --template fix-task --title "Fix: <description>" \
  --source-task-id <source-task-id> \
  --block-source \
  --var SOURCE_FILES="<affected paths>" \
  --var TEST_SCRIPT="<failing test>" \
  --var TEST_RESULTS="<results path>" \
  --description "<root cause>"
```

**`--block-source`**: atomically sets source task to blocked before resolution. `forge task add` automatically deduplicates — check output: `ACTION: ADDED` (new fix task) or `ACTION: SKIPPED` (active fix already exists).

## Step 5: Generate index.json via CLI

After all business task `.md` files (Step 3) are written, run:

```bash
forge task index --feature <slug>
```

This command:
1. Scans all `.md` files in `tasks/`, parses YAML frontmatter
2. Auto-generates stage-gate files (`<N>.summary.md` and `<N>.gate.md`) for phases with >=2 business tasks — detects numbered phases from task IDs using `<digit>.<digit>` pattern
3. Auto-generates test task `.md` files from embedded profiles
4. Produces `index.json` with all business + stage-gate + test tasks
5. Runs validation automatically

**Stage-gate auto-generation**: `forge task index` detects numbered phases from business task IDs (e.g., tasks `1.1`, `1.2`, `2.1` yield phases 1 and 2) and generates `.summary` and `.gate` files for phases with >=2 business tasks. Single-task phases are skipped. This is idempotent — existing files are preserved on re-run.

If the profile was not set in Step 0, pass it explicitly: `forge task index --feature <slug> --test-profiles <p1>,<p2>`.

## Step 6: Create Manifest

Read `templates/manifest-quick.md` for the simplified manifest format. Write to `docs/features/<slug>/manifest.md`.

The quick manifest contains:
- Documents table: proposal path + optional test-cases path
- Tasks table: ID, Title, Status for all tasks (business + test)
- No Traceability table

## Step 7: Validate

```bash
forge task validate-index docs/features/<slug>/tasks/index.json
```

## Output Checklist

- [ ] `docs/features/<slug>/tasks/` contains 1-10 business task files
- [ ] `index.json` valid per schema, `forge task validate-index` passes
- [ ] Stage-gate files (`.summary.md`, `.gate.md`) auto-generated by `forge task index` for phases with >=2 business tasks (if using `<phase>.<sub>` IDs)
- [ ] Every Success Criterion covered by ≥1 task
- [ ] Dependency graph is a DAG (no cycles)
- [ ] Test tasks appended with correct dependency chain:
  - Single profile: T-quick-1 through T-quick-5
  - Multiple profiles: T-quick-1a/1b/... through T-quick-4a/4b/... all per profile, T-quick-5 shared
- [ ] Per-profile test tasks include `profile` field in index.json entries
- [ ] `docs/features/<slug>/manifest.md` written with `mode: quick`

## Integration

Works well with:
- `/brainstorm` — Generate the proposal before running quick-tasks
- `/quick` — Full pipeline: brainstorm → quick-tasks → run-tasks
- `/run-tasks` — Execute generated tasks (index.json compatible)
