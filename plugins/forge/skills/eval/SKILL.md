---
name: eval
description: Generic document evaluation with scorer→gate→revise loop. Parameterized by rubric file. Supports 100-point and 1000-point scales. Detects UI platform for eval-ui. Skips reviser when iterations ≤ 1.
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

If missing, tell user to create it first.

## Parameters

| Parameter | Default | Description |
|-----------|---------|-------------|
| `--type` | (required) | `proposal`, `prd`, `design`, `ui`, `ui-web`, `ui-mobile`, `ui-tui`, `test-cases`, `ui-test-cases`, `tui-test-cases`, `mobile-test-cases`, `api-test-cases`, `cli-test-cases`, `consistency`, `harness`, `validate-code` |
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

Load: `plugins/forge/skills/eval/rubrics/<type>.md`
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
| `test-cases` | Resolve test profile via `forge profile`. Pass profile capabilities to scorer. |
| `ui-test-cases`, `tui-test-cases`, `mobile-test-cases`, `api-test-cases`, `cli-test-cases` | Resolve test profile via `forge profile`. Pass profile capabilities to scorer. |
| `prd` | Detect mode: `prd-ui-functions.md` exists → Mode A (with UI), else Mode B (no UI). |
| `validate-code` | 1) Read PRD → extract user scenarios list (from prd-spec.md flow descriptions and prd-user-stories.md acceptance criteria). 2) Run `git diff <base-branch>...HEAD` to get changed files and diff hunks. 3) Compile changed file list. 4) Pass PRD scenarios + diff + file list to scorer as assembled input. |

## Step 2: Invoke Scorer Subagent

Spawn `doc-scorer` via Agent tool (subagent_type: `forge:doc-scorer` or `general-purpose`).

Inputs:
- `DOC_DIR` = document directory
- `RUBRIC_PATH` = resolved rubric file
- `REPORT_PATH` = `<doc_dir>/eval/iteration-{{N}}.md`
  - `harness`: `docs/harness-reports/YYYY-MM-DD.md`
  - `consistency`: `docs/features/<slug>/eval-consistency/eval/iteration-{{N}}.md`
  - `proposal`: `docs/proposals/<slug>/eval/iteration-{{N}}.md`
  - `validate-code`: `docs/features/<slug>/eval/validate-code.md`
- `ITERATION` = current iteration (1-based)
- `PREVIOUS_REPORT_PATH` = previous report (only if iteration > 1)

Type-specific inputs:
- `ui-*`: add `PRD_PATH` = `docs/features/<slug>/prd/prd-ui-functions.md` (if exists)
- `test-cases`, `ui-test-cases`, `tui-test-cases`, `mobile-test-cases`, `api-test-cases`, `cli-test-cases`: add `PRD_FILES` = paths to prd-spec.md and prd-user-stories.md
- `consistency`: add `SCOPE` = value from `--scope`

Do NOT pass reviser change summaries to the scorer.

**Context Injection**: If `CONTEXT_CONTENT` was loaded in Step 1.4, append the following section to the scorer prompt:

```
<injected-context>
The following project reference material is provided for reality-checking the evaluated document. Use it to detect contradictions, violations, or gaps — do not evaluate the reference material itself.

{{CONTEXT_CONTENT}}
</injected-context>
```

After scorer returns, extract:
1. `SCORE: X/{{scale}}`
2. Per-dimension scores from `DIMENSIONS:` section
3. Attack points from `ATTACKS:` section

`test-cases`, `ui-test-cases`, `tui-test-cases`, `mobile-test-cases`, `api-test-cases`, `cli-test-cases`: If Step Actionability < 200, warn that gen-test-scripts is blocked.

## Step 3a: Single-Pass (iterations ≤ 1)

Skip gate and reviser. Go directly to Step 5.

## Step 3b: Decision Gate (Main Session)

| Condition | Action |
|-----------|--------|
| Score >= target | Go to Step 5 |
| Score < target, iterations remaining | Go to Step 4 |
| Score < target, no iterations remaining | Go to Step 5 (report failure) |

On "continue"/"keep going": run scorer again (Step 2), then re-evaluate this gate.

If proceeding to Step 4, report: `Iteration {{N}}/{{MAX}}: scored {{SCORE}}/{{SCALE}} (target: {{TARGET}}). Revising...`

## Step 4: Invoke Reviser Subagent (only when Step 3b routes here)

Spawn `doc-reviser` via Agent tool (subagent_type: `forge:doc-reviser` or `general-purpose`).

Inputs:
- `DOC_DIR`, `RUBRIC_PATH`, `EVAL_REPORT_PATH`, `ATTACK_POINTS`

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

`ui-*` invoked as sub-step of `/ui-design`: return control to ui-design, do NOT prompt.

## Rubric Reference

All rubrics: `plugins/forge/skills/eval/rubrics/<type>.md`

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
