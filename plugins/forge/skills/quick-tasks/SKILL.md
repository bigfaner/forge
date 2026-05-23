---
name: quick-tasks
description: Use for features (1-15 coding tasks, doc tasks unlimited) to generate tasks directly from proposal. No PRD or design needed.
---

# Quick Tasks

Generate executable tasks directly from a proposal document. For features (1-15 coding tasks, doc-type tasks unlimited) that don't need PRD or tech design.

## Prerequisites

| Artifact | Missing? Run |
|----------|-------------|
| `docs/proposals/<slug>/proposal.md` | `/brainstorm` or `/quick` |

<HARD-GATE>
Maximum 15 coding tasks (`coding.*` type). Doc-type tasks (`doc*` type prefix) are unlimited. If the proposal requires >15 coding tasks, STOP and recommend the full pipeline: `/write-prd` → `/tech-design` → `/breakdown-tasks`.
</HARD-GATE>

## Docs-Only Fast Path

When all tasks are `type: "doc"` (non-compilable, non-runnable output), skip **Step 0** (language) and **Step 4** (test tasks). **Step 5** (`forge task index`) is always mandatory — without `index.json`, `forge task claim` fails.

**Detection**: Step 1 extracts In Scope items → if every item targets non-compilable files only, the feature is docs-only.

```mermaid
graph LR
    S1["Step 1"] --> D{"docs-only?"}
    D -->|"skip 0, 4"| B["Steps 2–3"]
    D -->|No| A["Steps 0→2→3→4"]
    B --> S5
    A --> S5
    S5["Step 5: forge task index<br>(mandatory)"] --> S67["Steps 6–7"]
    style S5 fill:#fff3cd,stroke:#856404,stroke-width:2px
```

## Step 0: Resolve Language

1. Load Convention files from `docs/conventions/` by `domains` frontmatter (match `testing`, `go`, `typescript`, etc.). Extract language from `Framework` section.
2. Fallback: scan existing source/test files (`go.mod`, `package.json`, `*_test.go`, etc.). Also check subdirectories for monorepo.
3. On failure: ask user.

<HARD-RULE>
Do NOT silently default to any language.
</HARD-RULE>

Language info is used as context for task content (e.g., test framework selection). Test pipeline tasks are driven by the `interfaces` config field in `.forge/config.yaml`, not by language count.


## Step 1: Read Proposal

Determine the feature slug from the proposal directory name. Read `docs/proposals/<slug>/proposal.md` — the sole input document. Extract:

- **Problem** → task context and motivation
- **Proposed Solution** → task scope and boundaries
- **Scope > In Scope** → one task per bullet (split if >2h, merge if <30min)
- **Success Criteria** → acceptance criteria for each task
- **Key Risks** → implementation notes and risk mitigations

<HARD-RULE>
Enforce maximum 15 coding tasks (`coding.*` type). Doc-type tasks (`doc*` type prefix) are unlimited. If the In Scope section implies >15 coding tasks, STOP and recommend the full pipeline (`/write-prd` → `/tech-design` → `/breakdown-tasks`).
</HARD-RULE>

## Step 2: Derive Tasks

For each In Scope bullet: estimate effort (1-2h), derive acceptance criteria from Success Criteria, classify type (see Step 3 Template Selection), set scope via Scope Inference, fill Reference Files with directory paths from proposal context.

**Split by functional steps**: multiple independently verifiable steps in one bullet → separate tasks (coding tasks still ≤ 15, doc tasks unlimited).

**Dependencies**: linear chain unless parallel work implied. Simple integer IDs: `1`, `2`, `3`.

**Scope Inference** (from task description semantics): UI/pages/components → `scope: "frontend"`, API/server/database/CLI → `scope: "backend"`, mixed/unclear → `scope: "all"`.

## Step 3: Create Task Files

Read the appropriate template (see Template Selection below) for the task content structure. Create one task file per derived task in `docs/features/<slug>/tasks/`.

<HARD-RULE>
Naming & ID conventions:
- Business task: file `<seq>-<slug>.md`, ID `<seq>` (e.g., file `1-add-command.md`, ID `1`)
- Quick test: file `quick-<name>.md`, ID `T-quick-<N>`
- No phase prefixes, no sub-IDs, no summary/gate tasks
</HARD-RULE>

For each task, fill from proposal context: Description (Problem + Solution), Acceptance Criteria (Success Criteria), Implementation Notes (Key Risks). Fill Hard Rules only for critical constraints (specific recipes, hidden env deps, scope restrictions). Set `breaking: true` for tasks modifying shared interfaces/models/APIs.

### Type Assignment

Every task receives a `type` field in its frontmatter. The type controls quality-gate routing.

| Type | When to assign |
|------|----------------|
| `coding.feature` | Task adds new runtime behavior, new user-facing capability, or new files |
| `coding.enhancement` | Task improves existing behavior without adding new capabilities |
| `coding.cleanup` | Task removes dead code, fixes technical debt, or improves code hygiene |
| `coding.refactor` | Task restructures code without changing behavior (rename, reorganize, extract) |
| `doc` | Tasks producing only markdown, specs, or templates (non-compilable, non-runnable) |
| `doc.consolidate` | User manually creates a consolidation task for legacy projects — merging scattered spec files into `docs/business-rules/` or `docs/conventions/` |
| `doc.drift` | User manually creates a drift audit task — detecting inconsistencies between existing specs and current code |

Fallback: `coding.feature`. **Classify by output artifact, not intent.**

Test pipeline tasks are auto-generated by `forge task index`.

**Rule: classify by output artifact, not by intent.** The type determines quality-gate behavior. Quality-gate (compile, fmt, lint, test) only makes sense for compilable or runnable output. Therefore, the decisive factor is *what the task produces*, not *what the task intends to accomplish*.

| Category | Types | Quality-gate |
|----------|-------|-------------|
| Code | `coding.feature`, `coding.enhancement`, `coding.cleanup`, `coding.refactor`, `coding.fix` | Run (compile + fmt + lint + test) |
| Doc | `doc`, `doc.consolidate`, `doc.drift` | Skip entirely |

How to apply:

1. Look at the **affected files** listed in the task definition.
2. If all affected files are non-compilable, non-runnable artifacts (`.md`, `.yaml`, `.json` under `skills/`, `docs/`, etc.), the type **must** be `doc`.
3. If any affected file is compilable or runnable source code, use the appropriate Code type from the table above.

### Intent Propagation

If `proposal.md` frontmatter has `intent` (e.g., `intent: cleanup`), use as default type for all tasks. Individual task `type` overrides. Missing intent → per-task Type Assignment. 1:1 mapping.

### Template Selection

| Condition | Template |
|-----------|----------|
| All affected files non-compilable, non-runnable | `templates/task-doc.md` |
| Any compilable or runnable file | `templates/task.md` |

## Step 4: Test Tasks (auto-generated)

Test tasks are auto-generated by `forge task index` based on the `interfaces` field in `.forge/config.yaml`. **Do NOT create test task `.md` files manually.**

To add a fix task for a failing test: `forge task add --template fix-task --title "Fix: <desc>" --source-task-id <TASK_ID> --block-source --var SOURCE_FILES="<paths>" --var TEST_SCRIPT="<test>" --var TEST_RESULTS="<results>" --description "<root cause>"`

## Step 5: Generate index.json via CLI

After all business task `.md` files (Step 3) are written, run:

```bash
forge task index --feature <slug>
```

This auto-generates stage-gate files, test task `.md` files, and `index.json` (runs validation automatically). Existing files are preserved on re-run.

## Step 6: Create Manifest

Read `templates/manifest-quick.md` for the format. Write to `docs/features/<slug>/manifest.md`. Replace `{{DATE}}` with today's date in `YYYY-MM-DD` format.

## Step 7: Validate

```bash
forge task validate-index docs/features/<slug>/tasks/index.json
```

## Step 8: Commit Planning Artifacts

Only execute if Step 7 validation passed. If validation failed, fix issues first.

<HARD-RULE>
Stage only planning artifact paths — never use `git add -A` or `git add .`.
</HARD-RULE>

```bash
git add docs/features/<slug>/tasks/*.md docs/features/<slug>/tasks/index.json docs/features/<slug>/manifest.md
git commit -m "docs(<slug>): add quick-tasks planning artifacts"
```

Other uncommitted changes remain unstaged.

## Output Checklist

- [ ] `docs/features/<slug>/tasks/` contains ≤15 coding task files + any number of doc task files
- [ ] `index.json` valid per schema, `forge task validate-index` passes
- [ ] Stage-gate files (`.summary.md`, `.gate.md`) auto-generated by `forge task index` for phases with >=2 business tasks (if using `<phase>.<sub>` IDs)
- [ ] Every Success Criterion covered by ≥1 task
- [ ] Dependency graph is a DAG (no cycles)
- [ ] `docs/features/<slug>/manifest.md` written with `mode: quick`
- [ ] Planning artifacts committed (task .md files, index.json, manifest.md)

## Integration

- `/brainstorm` → generate proposal before quick-tasks
- `/quick` → full pipeline: brainstorm → quick-tasks → run-tasks
- `/run-tasks` → execute generated tasks
