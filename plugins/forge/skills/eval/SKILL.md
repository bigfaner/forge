---
name: eval
description: Generic document evaluation with scorer→gate→revise loop. Parameterized by rubric file. Supports 100-point and 1000-point scales. Detects UI platform for eval-ui. Skips reviser when iterations ≤ 1.
argument-hint: "[--type <type>] [--target 900] [--iterations 3]"
effort: high
---

# Eval

## Prerequisites

| Type | Required Artifact |
|------|-------------------|
| `proposal` | `docs/proposals/<slug>/proposal.md` |
| `prd` | `prd/prd-spec.md` + `prd/prd-user-stories.md` |
| `design` | `design/tech-design.md` |
| `ui-web`, `ui-mobile`, `ui-tui` | `ui/ui-design.md` |
| `test-cases` | `testing/test-cases.md` |
| `ui-test-cases` | `testing/ui-test-cases.md` |
| `tui-test-cases` | `testing/tui-test-cases.md` |
| `mobile-test-cases` | `testing/mobile-test-cases.md` |
| `api-test-cases` | `testing/api-test-cases.md` |
| `cli-test-cases` | `testing/cli-test-cases.md` |
| `consistency` | `manifest.md` + `prd/prd-spec.md` + at least one other doc |
| `harness` | Project has CLAUDE.md or AGENTS.md |
| `validate-code` | PRD (`prd/prd-spec.md` + `prd/prd-user-stories.md`) + git diff against base branch |
| `validate-ux` | PRD + compilable project (binary or web server) |

If missing, tell user to create it first.

## Parameters

| Parameter | Default | Description |
|-----------|---------|-------------|
| `--type` | (required) | `proposal`, `prd`, `design`, `ui`, `ui-web`, `ui-mobile`, `ui-tui`, `test-cases`, `ui-test-cases`, `tui-test-cases`, `mobile-test-cases`, `api-test-cases`, `cli-test-cases`, `consistency`, `harness`, `validate-code`, `validate-ux` |
| `--target` | rubric frontmatter | Override target score |
| `--iterations` | rubric frontmatter | Override max iterations |
| `--scope` | `docs` | `consistency` only: `docs` or `full` |

Resolution: explicit `--type` in `<command-args>` → command name `/eval-<type>` → ask user.

### Rubric Context Frontmatter (optional)

Rubrics may declare a `context` frontmatter field to inject project reality files into the scorer prompt. Rubrics without `context` continue to work unchanged.

```yaml
context:
  conventions: [api, naming, ux]  # list of category strings (optional)
  business-rules: auto            # "auto" or list of filenames (optional)
```

| Sub-field | Type | Description |
|-----------|------|-------------|
| `conventions` | list of strings | Each string matches filenames in `docs/conventions/` by prefix. E.g., `api` matches `api*.md`. Non-matching strings are skipped silently. |
| `business-rules` | `"auto"` or list of strings | `auto` loads all `.md` files from `docs/business-rules/`. A list specifies exact filenames. Missing files are skipped silently. |

At least one sub-field must be present for context injection to activate.

## Architecture

```mermaid
flowchart TD
    A([Start]) --> B["1. Resolve Type & Load Rubric"]
    B --> C{"iterations ≤ 1?"}
    C -->|"yes"| D["2a. Score (subagent)"]
    D --> E["3a. Final Report"]
    C -->|"no"| F["2b. Score (subagent)"]
    F --> G{"3b. Gate (main session)"}
    G -->|"score >= target"| E
    G -->|"score < target, no iterations left"| E
    G -->|"score < target, iterations remaining"| H["4. Revise (subagent)"]
    H --> F
```

## Orchestrator Iron Laws

<EXTREMELY-IMPORTANT>
- Main session owns the loop. NEVER delegate the full eval to a single agent.
- Per iteration: score (subagent) → gate (main session) → revise (subagent).
- Scorer and reviser are ALWAYS invoked via Agent tool, never inline.
</EXTREMELY-IMPORTANT>

## Step 1: Resolve Type, Rubric, and Locate Documents

### 1.1 Resolve Rubric Path

Load: `rubrics/<type>.md`
Exception: type `ui` → detect platform first (see 1.3), then load `ui-<platform>.md`.

Parse rubric frontmatter: `scale`, `target`, `iterations`, `context`. CLI `--target`/`--iterations` override frontmatter. Store `context` declaration for use in Step 1.4 and Step 2.

### 1.2 Locate Documents

1. User-provided path
2. `docs/features/<current-feature>/manifest.md`
3. Default paths:

| Type | Default Doc Dir |
|------|----------------|
| `proposal` | `docs/proposals/<slug>/` |
| `prd` | `docs/features/<slug>/prd/` |
| `design` | `docs/features/<slug>/design/` |
| `ui-*` | `docs/features/<slug>/ui/` |
| `test-cases` | `docs/features/<slug>/testing/` |
| `ui-test-cases` | `docs/features/<slug>/testing/` |
| `tui-test-cases` | `docs/features/<slug>/testing/` |
| `mobile-test-cases` | `docs/features/<slug>/testing/` |
| `api-test-cases` | `docs/features/<slug>/testing/` |
| `cli-test-cases` | `docs/features/<slug>/testing/` |
| `consistency` | `docs/features/<slug>/` |
| `harness` | `docs/harness-reports/` |
| `validate-code` | `docs/features/<slug>/prd/` |
| `validate-ux` | `docs/features/<slug>/prd/` |

4. Ask user if not found

### 1.3 UI Platform Detection (type `ui` only)

1. Check UI doc frontmatter for `platform` field
2. If absent, infer: ASCII mockups/terminal keybindings → `tui`; touch targets/safe areas → `mobile`; else → `web`
3. Load rubric `ui-<platform>.md`

Multi-platform: run independent score→gate→revise loops per platform.

### 1.4 Pre-Processing by Type

| Type | Before Scoring |
|------|---------------|
| **All types** | If rubric has `context` frontmatter, load filtered context files: (1) for each string in `conventions`, glob `docs/conventions/<string>*.md` and read matching files; (2) if `business-rules: auto`, glob `docs/business-rules/*.md` and read all, else read listed filenames. Concatenate into `CONTEXT_CONTENT` for Step 2 injection. Skip missing files silently (no error, no abort). |
| `harness` | Gather project context, write snapshot. Scorer evaluates snapshot, not raw files. |
| `consistency` | Assemble document bundle — copy relevant docs into flat directory for scorer. |
| `test-cases` | Resolve test language via `forge test detect`. Pass project interfaces to scorer. |
| `ui-test-cases`, `tui-test-cases`, `mobile-test-cases`, `api-test-cases`, `cli-test-cases` | Resolve test language via `forge test detect`. Pass project interfaces to scorer. |
| `prd` | Detect mode: `prd-ui-functions.md` exists → Mode A (with UI), else Mode B (no UI). |
| `validate-code` | 1) Read PRD → extract user scenarios list (from prd-spec.md flow descriptions and prd-user-stories.md acceptance criteria). 2) Run `git diff <base-branch>...HEAD` to get changed files and diff hunks. 3) Compile changed file list. 4) Pass PRD scenarios + diff + file list to scorer as assembled input. |
| `validate-ux` | **Two-phase pre-processing** (must execute in git worktree or temp dir). Full sub-pipeline: `${CLAUDE_SKILL_DIR}/rubrics/validate-ux-pipeline.md`. |

## Expert Dispatch Table

Expert files are located at `experts/scorer/`.

| type | scorer experts |
|------|---------------|
| `proposal` | `[cto]` |
| `prd` | `[pm, qa]` |
| `design` | `[architect]` |
| `ui-web`, `ui-mobile`, `ui-tui` | `[ux-engineer]` |
| `test-cases`, `ui-test-cases`, `tui-test-cases`, `mobile-test-cases`, `api-test-cases`, `cli-test-cases` | `[qa]` |
| `consistency` | `[editor]` |
| `harness` | `[harness-engineer]` |
| `validate-code` | `[code-reviewer]` |
| `validate-ux` | `[ux-auditor]` |

Fallback for unmapped types: use the generic inline prompt below (no expert file loaded).

```
You are a senior reviewer evaluating the document at {{DOC_DIR}} against the rubric at {{RUBRIC_PATH}}. Apply the rubric strictly and identify all weaknesses.
```

## Iteration Initialization

Set `ITERATION = 1`, `MAX_ITERATIONS = resolved value from rubric or CLI`.

## Step 2: Invoke Scorer Subagent(s)

### 2.1 Compose Scorer Prompts

Resolve the eval type to its scorer expert(s) from the dispatch table above.

Read the scorer protocol at `experts/protocol/scorer-protocol.md`.

For each expert, compose a scorer prompt by concatenating:
1. Scorer protocol content (with template variables replaced: `{{DOC_DIR}}`, `{{RUBRIC_PATH}}`, `{{REPORT_PATH}}`, `{{ITERATION}}`, `{{PREVIOUS_REPORT_PATH}}`)
2. Expert file content (e.g., `experts/scorer/pm.md`)
3. Context injection (if `CONTEXT_CONTENT` was loaded in Step 1.4 — see below)

**Context Injection**: If `CONTEXT_CONTENT` was loaded in Step 1.4, append the following section after the expert content in every composed prompt:

```
<injected-context>
The following project reference material is provided for reality-checking the evaluated document. Use it to detect contradictions, violations, or gaps — do not evaluate the reference material itself.

{{CONTEXT_CONTENT}}
</injected-context>
```

For unmapped types (not in dispatch table), compose a single prompt using the generic inline fallback above plus the scorer protocol (with variables replaced) plus context injection.

### 2.2 Spawn Scorer Agents

Spawn each composed prompt as a `general-purpose` agent via the Agent tool with `model: "sonnet"`.

- **Single-expert types**: spawn one agent.
- **Multi-expert types** (e.g., `prd` → `[pm, qa]`): spawn multiple agents **in parallel** (multiple Agent tool calls in a single message). Each agent receives its own composed prompt and writes to its own report path.

Common inputs for all agents:
- `DOC_DIR` = document directory
- `RUBRIC_PATH` = resolved rubric file
- `ITERATION` = current iteration (1-based)
- `PREVIOUS_REPORT_PATH` = previous report (only if iteration > 1)

Report paths per expert (for multi-expert types, each expert writes to a separate report):
- `REPORT_PATH` = `<doc_dir>/eval/iteration-{{N}}.md` (single-expert)
- `REPORT_PATH` = `<doc_dir>/eval/iteration-{{N}}-{{expert}}.md` (multi-expert)

Type-specific report path overrides:
- `harness`: `docs/harness-reports/YYYY-MM-DD.md`
- `consistency`: `docs/features/<slug>/eval-consistency/eval/iteration-{{N}}.md`
- `proposal`: `docs/proposals/<slug>/eval/iteration-{{N}}.md`
- `validate-code`: `docs/features/<slug>/eval/validate-code.md`
- `validate-ux`: `docs/features/<slug>/eval/validate-ux.md`

Type-specific inputs:
- `ui-*`: add `PRD_PATH` = `docs/features/<slug>/prd/prd-ui-functions.md` (if exists)
- `test-cases`, `ui-test-cases`, `tui-test-cases`, `mobile-test-cases`, `api-test-cases`, `cli-test-cases`: add `PRD_FILES` = paths to prd-spec.md and prd-user-stories.md
- `consistency`: add `SCOPE` = value from `--scope`
- `validate-ux`: add `UX_SNAPSHOT_PATH` = path to generated `ux-snapshot.md`

Do NOT pass reviser change summaries to the scorer.

### 2.3 Collect and Merge Results

After all scorer agents return:

**For single-expert types**: extract using robust score extraction:
1. Extract score using regex `/SCORE:\s*(\d+)\/(\d+)/`. If pattern not found, scan the scorer agent's output for the last line matching a `number/number` pattern. If still not found, report error and stop.
2. Per-dimension scores from `DIMENSIONS:` section
3. Attack points from `ATTACKS:` section

**For multi-expert types**:
1. Extract score and attacks from each expert's output
2. **Gate score**: average the total scores across all experts (rounded to nearest integer)
3. **Attack points merge**: LLM-merge attack points from all experts in the main session using this prompt:

```
Merge overlapping attack points from {{N}} expert evaluations. Keep unique attacks from each. Combine duplicates into single attacks preserving the strongest prescription. Do not remove any unique attack. Output the merged list in the same format:

1. [dimension]: [specific weakness] — [quote from document] — [what must improve]
```

4. **Write merged report**: Write the merged attacks + averaged scores to `<doc_dir>/eval/iteration-{{N}}-merged.md`. This file serves as `EVAL_REPORT_PATH` for the reviser (Step 4.1). Single-expert types continue using `iteration-{{N}}.md` directly.

`test-cases`, `ui-test-cases`, `tui-test-cases`, `mobile-test-cases`, `api-test-cases`, `cli-test-cases`: If Step Actionability < 200, warn that gen-test-scripts is blocked.

## Step 3a: Single-Pass (iterations ≤ 1)

Skip gate and reviser. Go directly to Step 5.

## Step 3b: Decision Gate (Main Session)

Use the averaged score (for multi-expert types) or single score (for single-expert types) from Step 2.3.

| Condition | Action |
|-----------|--------|
| Score >= target | Go to Step 5 |
| Score < target, iterations remaining | Go to Step 4 |
| Score < target, no iterations remaining | Go to Step 5 (report failure) |

If proceeding to Step 4, report: `Iteration {{N}}/{{MAX}}: scored {{SCORE}}/{{SCALE}} (target: {{TARGET}}). Revising...`

## Step 4: Invoke Reviser Subagent (only when Step 3b routes here)

### 4.1 Compose Reviser Prompt

Read the reviser protocol at `experts/protocol/reviser-protocol.md`.

Resolve `EVAL_REPORT_PATH`:
- **Single-expert types**: `<doc_dir>/eval/iteration-{{N}}.md`
- **Multi-expert types**: `<doc_dir>/eval/iteration-{{N}}-merged.md` (written in Step 2.3)

Compose the reviser prompt by concatenating:
1. Reviser protocol content (with template variables replaced: `{{DOC_DIR}}`, `{{EVAL_REPORT_PATH}}`)
2. The merged `ATTACK_POINTS` from Step 2.3 (replacing the `{{ATTACK_POINTS}}` placeholder in the protocol)
3. Context injection (if `CONTEXT_CONTENT` was loaded in Step 1.4 — see below)

**Context Injection**: If `CONTEXT_CONTENT` was loaded in Step 1.4, append the following section after the attack points in the reviser prompt:

```
<injected-context>
The following project reference material is provided for reality-checking the evaluated document. Use it to detect contradictions, violations, or gaps — do not evaluate the reference material itself.

{{CONTEXT_CONTENT}}
</injected-context>
```

The reviser receives **only** the protocol + merged attacks + optional context. No rubric, no expert file.

### 4.2 Spawn Reviser Agent

Spawn as a `general-purpose` agent via the Agent tool with `model: "sonnet"`.

Inputs:
- `DOC_DIR`
- `EVAL_REPORT_PATH`
- `ATTACK_POINTS` (merged)

Type-specific constraints:
- `consistency`: Do NOT modify `prd/`. Classify attack points by fix target before invoking.
- `test-cases`: ONLY modify `test-cases.md`.
- `ui-test-cases`, `tui-test-cases`, `mobile-test-cases`, `api-test-cases`, `cli-test-cases`: ONLY modify `{type}-test-cases.md`.

After reviser completes:
- `consistency`: re-assemble document bundle
- Increment iteration counter, return to Step 2

## Step 5: Final Report

```
## Eval-{{TYPE}} Complete
**Final Score**: {{SCORE}}/{{SCALE}} (target: {{TARGET}})
**Iterations Used**: {{N}}/{{MAX}}

### Score Progression
| Iteration | Score | Delta |
|-----------|-------|-------|

### Dimension Breakdown (final)
{{from rubric}}

### Outcome
{{"Target reached" / "Target NOT reached — N iterations exhausted"}}
```

Type-specific additions:
- `harness`: priority improvement table (P0/P1/P2)
- `consistency`: "Files Modified" and "Residual Issues"
- `test-cases`, `ui-test-cases`, `tui-test-cases`, `mobile-test-cases`, `api-test-cases`, `cli-test-cases`: Step Actionability blocking warning if < 200
- `design`: Breakdown-Readiness gate status

Save report to type-specific report path.

## Step 6: Next Step

Ask user via `AskUserQuestion`:

| Type | Next Skill |
|------|-----------|
| `proposal` | `/write-prd` |
| `prd` | `/ui-design` or `/tech-design` |
| `design` | `/breakdown-tasks` |
| `ui-*` | `/tech-design` |
| `test-cases` | `/gen-test-scripts` |
| `ui-test-cases`, `tui-test-cases`, `mobile-test-cases`, `api-test-cases`, `cli-test-cases` | `/gen-test-scripts` |
| `consistency` | `/run-tasks` or re-eval |
| `harness` | `/improve-harness` |
| `validate-code` | `/run-tasks` (proceed to test pipeline) |
| `validate-ux` | `/run-tasks` (feature complete) |

`ui-*` invoked as sub-step of `/ui-design`: return control to ui-design, do NOT prompt.

## Rubric Reference

All rubrics: `rubrics/<type>.md`

| Rubric | Scale | Target | Iterations | Notes |
|--------|-------|--------|------------|-------|
| `proposal` | 1000 | 900 | 3 | |
| `prd` | 1000 | 900 | 3 | Mode A/B detection |
| `design` | 1000 | 900 | 3 | Breakdown-Readiness gate |
| `ui-web` | 1000 | 950 | 3 | |
| `ui-mobile` | 1000 | 950 | 3 | |
| `ui-tui` | 1000 | 950 | 3 | |
| `test-cases` | 1000 | 900 | 6 | Step Actionability blocking threshold |
| `ui-test-cases` | 1000 | 900 | 6 | Step Actionability blocking threshold |
| `tui-test-cases` | 1000 | 900 | 6 | Step Actionability blocking threshold |
| `mobile-test-cases` | 1000 | 900 | 6 | Step Actionability blocking threshold |
| `api-test-cases` | 1000 | 900 | 6 | Step Actionability blocking threshold |
| `cli-test-cases` | 1000 | 900 | 6 | Step Actionability blocking threshold |
| `consistency` | 1000 | 900 | 3 | docs/full scope modes |
| `harness` | 100 | 70 | 1 | Single-pass; no reviser |
| `validate-code` | 1000 | 700 | 1 | Single-pass; scenario tracing; no reviser |
| `validate-ux` | 1000 | 700 | 1 | Single-pass; two-phase (snapshot + score); no reviser |
