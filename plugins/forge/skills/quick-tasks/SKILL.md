---
name: quick-tasks
description: Use for features (coding tasks and doc tasks unlimited) to generate tasks directly from proposal. No PRD or design needed.
---

# Quick Tasks

Generate executable tasks directly from a proposal document. For features (coding tasks and doc-type tasks unlimited) that don't need PRD or tech design.

## Prerequisites

| Artifact | Missing? Run |
|----------|-------------|
| `docs/proposals/<slug>/proposal.md` | `/brainstorm` or `/quick` |

<HARD-GATE>
Maximum 6 Acceptance Criteria per task. If a task has >6 AC, its scope is too large — split further by functional boundary. No overall task count cap; task volume is bounded by proposal scope and the AC max rule.
</HARD-GATE>

## Docs-Only Fast Path

When all tasks are `type: "doc"`, skip **Step 0** (language) and **Step 5** (test tasks). **Step 6** (`forge task index`) is always mandatory.

**Detection**: Step 1 extracts In Scope items → if every item targets non-compilable files only, the feature is docs-only.

## Step 0: Resolve Language

Discover via `docs/conventions/testing/index.md` (preferred) or scan existing source/test files. On failure: ask user.

<HARD-RULE>
Do NOT silently default to any language. Do NOT use `domains` frontmatter filtering — use index.md-based discovery.
</HARD-RULE>


## Step 1: Read Proposal

Determine the feature slug from the proposal directory name. Read `docs/proposals/<slug>/proposal.md` — the sole input document. Extract:

- **Problem** → task context and motivation
- **Proposed Solution** → task scope and boundaries
- **Scope > In Scope** → one task per bullet (split if not independently verifiable, merge if independently verifiable together)
- **Success Criteria** → acceptance criteria for each task
- **Key Risks** → implementation notes and risk mitigations

## Step 2: Derive Tasks

For each In Scope bullet: estimate effort (1-2h), derive acceptance criteria from Success Criteria, classify type (see Step 3), resolve surface-key/surface-type, fill Reference Files with section-level references from proposal context.

**Split Rules** (priority order): (1) Independently verifiable standard — separate tasks if outcomes require different verification contexts. (2) Multi-verb detection — split by functional boundary when verbs target different concerns. (3) Operational ceiling — split by file group when modifying >8 files with the same pattern; each sub-task targets ≤8 files.

**Complexity判定** (at task generation): `low` = AC ≤ 3 AND no Hard Rules AND Reference Files ≤ 1; `high` = AC ≥ 5 OR has Hard Rules; `medium` = everything else. LLM judgment override allowed with reason in Implementation Notes.

**Dependencies**: linear chain unless parallel work implied. Simple integer IDs: `1`, `2`, `3`.

**Surface-Key/Type Inference**: Run `forge surfaces` once. Single-surface project → set same key/type on all tasks (scalar surfaces: `surface-key` empty). Multi-surface → path-prefix match per task; ambiguous → call per file. On failure: leave empty, continue. Parsing rule: see Forge Guide → Surface Output Parsing.

**Reference Files Generation**:

<HARD-RULE>
1. BEFORE writing Reference Files, Grep `^#{1,4} ` on `docs/proposals/<slug>/proposal.md` to extract all headers. Only use headers that actually exist. If no match found, omit `(ref: ...)` — never fabricate headers.
2. First entry MUST be the full proposal path: `- docs/proposals/<slug>/proposal.md — <relevant sections>`
</HARD-RULE>

**Inline format**: `- <file-path>: <change description> (ref: <header>)`. Max 5 inline entries. Each coding/doc task needs ≥1 inline reference. Include `docs/lessons/`, `docs/conventions/`, `docs/decisions/` when relevant.

**Priority**: P0 = core mechanism or blocks others; P1 = maps to scope item or success criterion; P2 = polish/edge cases.

## Step 3: Create Task Files

Read the appropriate template for task content structure. Create one task file per derived task in `docs/features/<slug>/tasks/`.

### Task Template Placeholders

| Placeholder | Value Source |
|-------------|-------------|
| `{{ID}}` | Sequential integer (e.g., `1`, `2`) |
| `{{TITLE}}` | Derived from the In Scope bullet — concise imperative title |
| `{{PRIORITY}}` | P0 / P1 / P2 (see Step 2 Priority) |
| `{{ESTIMATED_TIME}}` | Effort estimate (e.g., `"1h"`, `"2h"`) |
| `{{DEPENDENCIES}}` | Comma-separated task IDs; empty `[]` for first task |
| `{{SLUG}}` | Feature slug (from proposal directory name) |
| `{{DESCRIPTION}}` | Task description from Problem + Solution context |
| `{{ACCEPTANCE_CRITERIA}}` | Derived from Success Criteria as `- [ ]` checklist items |
| `{{HARD_RULES}}` | Critical constraints only; leave empty for normal tasks |
| `{{NOTES}}` | Implementation notes from Key Risks |
| `{{SURFACE_KEY}}` | Surface key (coding tasks only) |
| `{{SURFACE_TYPE}}` | Surface type (coding tasks only) |
| `{{NEW_FILES}}` / `{{MODIFIED_FILES}}` / `{{DELETED_FILES}}` | File scope for doc tasks |

<HARD-RULE>
Naming & ID conventions:
- Business task: file `<seq>-<slug>.md`, ID `<seq>` (e.g., file `1-add-command.md`, ID `1`)
- Auto-generated tasks: `T-test-*`, `T-quick-doc-drift`, `T-validate-*`, `T-clean-*` (by `forge task index`; do NOT create manually)
- No phase prefixes, no sub-IDs, no summary/gate tasks
</HARD-RULE>

### File Scope Boundary

For tasks involving multiple files: enumerate exact file names in Implementation Notes (never "all templates" or vague terms). When operational ceiling triggers split, add Hard Rule: `仅修改以下文件：<file list>`.

### Breaking Task Test Impact Assessment

When `breaking: true`, add to Implementation Notes:
```
### Test Impact
- Affected test suite(s): <test directory paths>
- Expected fixture changes: <which test fixtures need updating>
- Risk level: low/medium/high
```
Fix-tasks in the same test directory are merged into one.

### Type Assignment

| Type | When to assign |
|------|----------------|
| `coding.feature` | New runtime behavior, new user-facing capability, or new files |
| `coding.enhancement` | Improves existing behavior without new capabilities |
| `coding.cleanup` | Removes dead code, fixes tech debt, improves hygiene |
| `coding.refactor` | Restructures code without behavior change |
| `coding.fix` | Auto-generated for test failures; do not assign manually |
| `doc` | Non-compilable, non-runnable output only (e.g., `.md`, `.yaml`, `.json`, `.sql`, `.toml`, `.graphql`) |
| `doc.consolidate` | User-created consolidation task (legacy projects) |
| `doc.drift` | User-created drift audit task |

Fallback: `coding.feature`. **Classify by output artifact, not intent.** If the task produces no compilable or runnable files, type must be `doc`.

<HARD-RULE>
Non-compilable files (`.md`, `.sql`, `.yaml`, `.json`, `.toml`, `.graphql`, etc.) are always non-compilable regardless of directory location — even under `pkg/`, `src/`, `internal/`. If output is ONLY non-compilable files, type **must** be `doc`, not `coding.*`. Decision test: "Does the output include any file that needs compilation or runtime testing?" If NO → `doc`.
</HARD-RULE>

| Category | Quality-gate |
|----------|-------------|
| Code (`coding.*`) | Run (compile + fmt + lint + test) |
| Doc (`doc`, `doc.consolidate`, `doc.drift`) | Skip |

### Intent Propagation

If `proposal.md` has `intent`, use as default type. 1:1 mapping: `new-feature`→`coding.feature`, `enhancement`→`coding.enhancement`, `refactor`→`coding.refactor`, `cleanup`→`coding.cleanup`, `fix`→`coding.fix`, `doc`→`doc`. Individual task `type` overrides. `doc.consolidate` and `doc.drift` are auto-generated, unified under `doc`.

### Template Selection

All files non-compilable → `templates/task-doc.md`; any compilable → `templates/task.md`.

## Step 4: Task Sizing Audit

After all task files are written, self-audit every task: (1) Multi-verb detection — split if title links independent actions; (2) AC cross-domain — split if AC covers unrelated domains; (3) Operational ceiling — split if modifying >8 files with same pattern. Split → re-assign IDs, re-wire dependencies, output audit report.

<HARD-GATE>
If any task still has >6 AC after splitting, split further. Do not proceed to Step 6 until all tasks pass.
</HARD-GATE>

## Step 5: Test Tasks (auto-generated)

Auto-generated by `forge task index` from `.forge/config.yaml` surfaces. **Do NOT create test `.md` files manually.**

Fix task: `forge task add --type <fix-type> --title "Fix: <desc>" --source-task-id <TASK_ID> --block-source --var SOURCE_FILES="<paths>" --var TEST_SCRIPT="<test>" --var TEST_RESULTS="<results>" --description "<cause>"`

Fix-Type: `doc`/`eval` → `doc.fix`; `coding`/`test`/`validation`/`gate` → `coding.fix`.

## Step 6: Generate index.json via CLI

```bash
forge task index --feature <slug>
```

Auto-generates test tasks and `index.json` (validates automatically). Quick mode uses simple integer IDs — no stage-gates.

## Step 7: Validate

```bash
forge task validate docs/features/<slug>/tasks/index.json
```

## Step 8: Create Manifest

Read `templates/manifest-quick.md`. Write `docs/features/<slug>/manifest.md` with placeholders: `{{SLUG}}`, `{{DATE}}`, `{{TASK_ROWS}}`.

## Step 9: Commit Planning Artifacts

Only if Step 7 passed.

<HARD-RULE>
Stage only planning artifact paths — never `git add -A` or `git add .`.
</HARD-RULE>

```bash
git add docs/features/<slug>/tasks/*.md docs/features/<slug>/tasks/index.json docs/features/<slug>/manifest.md
git commit -m "docs(<slug>): add quick-tasks planning artifacts"
```

## Output Checklist

- [ ] Task files with ≤6 AC each
- [ ] `index.json` valid, `forge task validate` passes
- [ ] Every Success Criterion covered by ≥1 task
- [ ] DAG (no cycles)
- [ ] `manifest.md` written with `mode: quick`
- [ ] Planning artifacts committed
