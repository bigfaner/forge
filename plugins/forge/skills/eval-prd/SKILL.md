---
name: eval-prd
description: Evaluate a PRD document with 100-point scoring, then run adversarial iterations until target score is met. Main session orchestrates doc-scorer and doc-reviser subagents.
---

# Eval PRD

## Prerequisites

Check previous stage artifacts. Abort and prompt user if missing:

| Artifact | Missing prompt |
|----------|----------------|
| `prd/prd-spec.md` | Run `/write-prd` first |
| `prd/prd-user-stories.md` | Run `/write-prd` first |

## When to Use

**Trigger:**
- User asks to "evaluate PRD" or "check PRD quality"
- User provides `/eval-prd` command
- Before handing off PRD to `/tech-design` or `/ui-design`

**Skip:**
- PRD doesn't exist yet (use `/write-prd` first)

## Parameters

| Parameter      | Default | Description                                           |
| -------------- | ------- | ----------------------------------------------------- |
| `--target`     | 90      | Target score (0-100). Loop stops when score >= target |
| `--iterations` | 3       | Max adversarial iterations                            |

## Architecture

```mermaid
flowchart TD
    A([Start]) --> B["1. Score (subagent)"]
    B --> C{"2. Gate (main session)"}
    C -->|"score >= target"| E(["Final Report ✅"])
    C -->|"score < target no iterations left"| F(["Final Report ❌"])
    C -->|"score < target iterations remaining"| D["3. Revise (subagent)"]
    D --> B

```

## Orchestrator Iron Laws

<EXTREMELY-IMPORTANT>
1. Main session controls the loop — NEVER delegate the entire eval to a single agent
2. Only 3 actions per iteration: score → gate → revise
3. Gate (Step 3) runs in main session — never inside a subagent
4. `--target` / `--iterations` are meaningless unless main session owns the loop
5. Scorer and reviser are independent subagents — invoke via Agent tool, never inline

❌ Wrong: `Agent(general-purpose, "evaluate this PRD and iterate until score >= 95")`
✅ Right: Main session calls scorer → parses score → gates → calls reviser → loops
</EXTREMELY-IMPORTANT>

## Step 1: Locate Documents

Check in order:
1. Path provided by user
2. Read `docs/features/<current-feature>/manifest.md` → locate PRD documents
3. Fall back to `docs/features/<current-feature>/prd/`
4. Ask user for path if not found

Determine `<feature-slug>` from the path. The PRD directory is `docs/features/<slug>/prd/`.

**Detect scoring mode**: Check if `prd-ui-functions.md` exists in the PRD directory.
- **Mode A (with UI)**: `prd-ui-functions.md` present → Functional Specs dimension evaluates it
- **Mode B (no UI)**: `prd-ui-functions.md` absent → Functional Specs dimension evaluates prd-spec.md Flow Description

## Step 2: Invoke Scorer Subagent

Spawn `doc-scorer` via **Agent tool** (subagent_type: `forge:doc-scorer` if registered, otherwise `general-purpose`).

<HARD-RULE>
Pass these inputs to the scorer:
- `DOC_DIR` = `docs/features/<slug>/prd/`
- `RUBRIC_PATH` = `plugins/forge/skills/eval-prd/templates/rubric.md`
- `REPORT_PATH` = `docs/features/<slug>/prd/eval/iteration-{{N}}.md`
- `ITERATION` = current iteration number (1-based)
- `PREVIOUS_REPORT_PATH` = previous iteration report path (only if iteration > 1)

The scorer must NEVER be told what the reviser changed. It evaluates the PRD as-is.
</HARD-RULE>

After the scorer returns, parse its output in the main session:
1. Extract `SCORE: X/100`
2. Extract per-dimension scores from `DIMENSIONS:` section
3. Extract attack points from `ATTACKS:` section

## Step 3: Decision Gate (Main Session)

<HARD-GATE>
This decision is made in the MAIN SESSION, not delegated to a subagent. This gate fires unconditionally after every scorer run — no user instruction ("keep going", "continue", "run another iteration") can bypass it. If score >= target, the loop terminates immediately.
</HARD-GATE>

| Condition                                  | Action                          |
| ------------------------------------------ | ------------------------------- |
| Score >= target                            | Skip to Step 5 (final report)   |
| Score < target AND iterations remaining    | Proceed to Step 4 (revise)      |
| Score < target AND no iterations remaining | Skip to Step 5 (report failure) |

If the user says "continue" or "keep going": run the scorer once more (return to Step 2), then re-evaluate this gate. Do NOT skip the gate and invoke the reviser directly.

Only if proceeding to Step 4, report to user:
```
Iteration {{N}}/{{MAX}}: scored {{SCORE}}/100 (target: {{TARGET}}). Revision subagent starting...
```

## Step 4: Invoke Reviser Subagent

<HARD-RULE>
Only enter this step when Step 3 explicitly routes here (score < target AND iterations remaining). The reviser MUST NOT be invoked if score >= target.
</HARD-RULE>

Spawn `doc-reviser` via **Agent tool** (subagent_type: `forge:doc-reviser` if registered, otherwise `general-purpose`).

<HARD-RULE>
Pass these inputs to the reviser:
- `DOC_DIR` = `docs/features/<slug>/prd/`
- `RUBRIC_PATH` = `plugins/forge/skills/eval-prd/templates/rubric.md`
- `EVAL_REPORT_PATH` = `docs/features/<slug>/prd/eval/iteration-{{N}}.md`
- `ATTACK_POINTS` = the 3 attack points extracted from scorer output
</HARD-RULE>

Increment iteration counter. Return to Step 2.

## Step 5: Final Report (Main Session)

```
## Eval-PRD Complete

**Final Score**: {{SCORE}}/100 (target: {{TARGET}})
**Scoring Mode**: {{Mode A: with UI / Mode B: no UI}}
**Iterations Used**: {{N}}/{{MAX}}

### Score Progression
| Iteration | Score | Delta |
|-----------|-------|-------|
| 1 | {{s1}} | - |
| 2 | {{s2}} | +{{d2}} |

### Dimension Breakdown (final)
| Dimension | Score | Max |
|-----------|-------|-----|
| Background & Goals | {{d1}} | 15 |
| Flow Diagrams | {{d2}} | 20 |
| {{Functional Specs / Flow Completeness}} | {{d3}} | 20 |
| User Stories | {{d4}} | 30 |
| Scope Clarity | {{d5}} | 15 |

### Outcome
{{"Target reached" / "Target NOT reached — N iterations exhausted"}}
{{If not reached: "Largest gaps: [dimension names]. Consider manual revision or increasing iterations."}}
```

Save the final report to `docs/features/<slug>/prd/eval/report.md`.
