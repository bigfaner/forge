---
name: eval-design
description: Evaluate a tech design document with 100-point scoring, then run adversarial iterations until target score is met. Main session orchestrates doc-scorer and doc-reviser subagents.
---

# Eval Design

Focus on whether the design can directly drive `/breakdown-tasks`.

## Prerequisites

Check previous stage artifacts. Abort and prompt user if missing:

| Artifact | Missing prompt |
|----------|----------------|
| `design/tech-design.md` | Run `/tech-design` first |

## When to Use

**Trigger:**
- User asks to "evaluate design" or "check design quality"
- User provides `/eval-design` command
- Before handing off tech-design.md to `/breakdown-tasks`

**Skip:**
- tech-design.md doesn't exist yet (use `/tech-design` first)

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

❌ Wrong: `Agent(general-purpose, "evaluate this design and iterate until score >= 80")`
✅ Right: Main session calls scorer → parses score → gates → calls reviser → loops
</EXTREMELY-IMPORTANT>

## Step 1: Locate Design Documents

Check in order:
1. Path provided by user
2. Read `docs/features/<current-feature>/manifest.md` → locate design documents
3. Fall back to `docs/features/<current-feature>/design/`
4. Ask user for path if not found

Determine `<feature-slug>` from the path. The design directory is `docs/features/<slug>/design/`.

## Step 2: Invoke Scorer Subagent

Spawn `doc-scorer` via **Agent tool** (subagent_type: `forge:doc-scorer` if registered, otherwise `general-purpose`).

<HARD-RULE>
Pass these inputs to the scorer:
- `DOC_DIR` = `docs/features/<slug>/design/`
- `RUBRIC_PATH` = `plugins/forge/skills/eval-design/templates/rubric.md`
- `REPORT_PATH` = `docs/features/<slug>/design/eval/iteration-{{N}}.md`
- `ITERATION` = current iteration number (1-based)
- `PREVIOUS_REPORT_PATH` = previous iteration report path (only if iteration > 1)
- `HAS_DB_SCHEMA` = `true` if `er-diagram.md` exists in DOC_DIR, `false` otherwise. When true, scorer uses the db-schema variant of Dimension 2 sub-criteria from the rubric.

The scorer must NEVER be told what the reviser changed. It evaluates the design as-is.
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
- `DOC_DIR` = `docs/features/<slug>/design/`
- `RUBRIC_PATH` = `plugins/forge/skills/eval-design/templates/rubric.md`
- `EVAL_REPORT_PATH` = `docs/features/<slug>/design/eval/iteration-{{N}}.md`
- `ATTACK_POINTS` = the 3 attack points extracted from scorer output
</HARD-RULE>

Increment iteration counter. Return to Step 2.

## Step 5: Final Report (Main Session)

```
## Eval-Design Complete

**Final Score**: {{SCORE}}/100 (target: {{TARGET}})
**Iterations Used**: {{N}}/{{MAX}}

### Score Progression
| Iteration | Score | Delta |
|-----------|-------|-------|
| 1 | {{s1}} | - |
| 2 | {{s2}} | +{{d2}} |

### Dimension Breakdown (final)
| Dimension | Score | Max |
|-----------|-------|-----|
| Architecture Clarity | {{d1}} | 20 |
| Interface & Model Definitions | {{d2}} | 20 |
| Error Handling | {{d3}} | 15 |
| Testing Strategy | {{d4}} | 15 |
| Breakdown-Readiness ★ | {{d5}} | 20 |
| Security Considerations | {{d6}} | 10 |

### Outcome
{{"Target reached" / "Target NOT reached — N iterations exhausted"}}
{{Breakdown-Readiness: {{score}}/20 — can/cannot proceed to /breakdown-tasks}}
{{If not reached: "Largest gaps: [dimension names]. Consider manual revision or increasing iterations."}}
```

Save the final report to `docs/features/<slug>/design/eval/report.md`.
