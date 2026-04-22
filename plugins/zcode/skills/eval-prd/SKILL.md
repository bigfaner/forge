---
name: eval-prd
description: Evaluate a PRD document with 100-point scoring, then run adversarial iterations until target score is met. Main session orchestrates doc-scorer and doc-reviser subagents.
---

# Eval PRD

## Prerequisites

检查上一阶段产物，缺失则中止并提示用户：

| 产物 | 缺失时提示 |
|------|-----------|
| `prd/prd-spec.md` | 先执行 `/write-prd` |
| `prd/prd-user-stories.md` | 先执行 `/write-prd` |

## When to Use

**Trigger:**
- User asks to "evaluate PRD" or "check PRD quality"
- User provides `/eval-prd` command
- Before handing off PRD to `/design-tech` or `/ui-design`

**Skip:**
- PRD doesn't exist yet (use `/write-prd` first)

## Parameters

| Parameter      | Default | Description                                           |
| -------------- | ------- | ----------------------------------------------------- |
| `--target`     | 80      | Target score (0-100). Loop stops when score >= target |
| `--iterations` | 3       | Max adversarial iterations                            |

## Architecture

```
Main Session (orchestrator)
  │
  ├─ iteration N:
  │   ├── Agent (doc-scorer)  ──→ score + attack points
  │   ├── score >= target? ──→ yes: stop
  │   └── Agent (doc-reviser) ──→ revised PRD doc(s)
  │
  └─ Final report to user
```

<EXTREMELY-IMPORTANT>
Scorer and reviser are **independent subagents** in `plugins/zcode/agents/`. Do NOT inline their prompts. Invoke via Agent tool only.
</EXTREMELY-IMPORTANT>

## Step 1: Locate Documents

Check in order:
1. Path provided by user
2. Read `docs/features/<current-feature>/manifest.md` → locate PRD documents
3. Fall back to `docs/features/<current-feature>/prd/`
4. Ask user for path if not found

Determine `<feature-slug>` from the path. The PRD directory is `docs/features/<slug>/prd/`.

## Step 2: Invoke Scorer Subagent

Spawn `doc-scorer` via **Agent tool** (subagent_type: `zcode:doc-scorer` if registered, otherwise `general-purpose`).

<HARD-RULE>
Pass these inputs to the scorer:
- `DOC_DIR` = `docs/features/<slug>/prd/`
- `RUBRIC_PATH` = `plugins/zcode/skills/eval-prd/templates/rubric.md`
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
This decision is made in the MAIN SESSION, not delegated to a subagent.
</HARD-GATE>

| Condition                                  | Action                          |
| ------------------------------------------ | ------------------------------- |
| Score >= target                            | Skip to Step 5 (final report)   |
| Score < target AND iterations remaining    | Proceed to Step 4 (revise)      |
| Score < target AND no iterations remaining | Skip to Step 5 (report failure) |

Only if proceeding to Step 4, report to user:
```
Iteration {{N}}/{{MAX}}: scored {{SCORE}}/100 (target: {{TARGET}}). Revision subagent starting...
```

## Step 4: Invoke Reviser Subagent

Spawn `doc-reviser` via **Agent tool** (subagent_type: `zcode:doc-reviser` if registered, otherwise `general-purpose`).

<HARD-RULE>
Pass these inputs to the reviser:
- `DOC_DIR` = `docs/features/<slug>/prd/`
- `RUBRIC_PATH` = `plugins/zcode/skills/eval-prd/templates/rubric.md`
- `EVAL_REPORT_PATH` = `docs/features/<slug>/prd/eval/iteration-{{N}}.md`
- `ATTACK_POINTS` = the 3 attack points extracted from scorer output
</HARD-RULE>

Increment iteration counter. Return to Step 2.

## Step 5: Final Report (Main Session)

```
## Eval-PRD Complete

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
| Background & Goals | {{d1}} | 20 |
| Flow Diagrams | {{d2}} | 20 |
| Functional Specs | {{d3}} | 20 |
| User Stories | {{d4}} | 20 |
| Scope Clarity | {{d5}} | 20 |

### Outcome
{{"Target reached" / "Target NOT reached — N iterations exhausted"}}
{{If not reached: "Largest gaps: [dimension names]. Consider manual revision or increasing iterations."}}
```

Save the final report to `docs/features/<slug>/prd/eval/report.md`.

## Related

- `/write-prd` — Create or revise the PRD
- `/design-tech` — Next step after PRD passes evaluation
- `/ui-design` — Next step (optional) if prd-ui-functions.md exists
- `/breakdown-tasks` — After design docs are finalized
