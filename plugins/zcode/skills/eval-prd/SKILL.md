---
name: eval-prd
description: Evaluate a PRD document with 100-point scoring, then run adversarial iterations until target score is met. Main session orchestrates doc-scorer and doc-reviser subagents.
---

# Eval PRD

评估 PRD 文档质量（百分制），通过 doc-scorer / doc-reviser subagent 多轮对抗迭代，直到达到目标分数或耗尽迭代次数。

**架构**：主会话调度循环，独立 subagent 负责评分和修订。

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

Parse from user input. Examples:

- `/eval-prd` → target=80, iterations=3
- `/eval-prd --target 90` → target=90, iterations=3
- `/eval-prd --target 85 --iterations 5` → target=85, iterations=5

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
Scorer and reviser are **independent subagents** defined in `plugins/zcode/agents/`. Do NOT inline their prompts. Invoke them via Agent tool with the correct inputs.
</EXTREMELY-IMPORTANT>

## Step 1: Locate Documents

Check in order:
1. Path provided by user
2. Read `docs/features/<current-feature>/manifest.md` → locate PRD documents
3. Fall back to `docs/features/<current-feature>/prd/prd-spec.md` + `prd/prd-user-stories.md`
4. Ask user for path if not found

Determine `<feature-slug>` from the path. The PRD directory is `docs/features/<slug>/prd/`.

## Step 2: Invoke Scorer Subagent

Spawn `doc-scorer` agent via **Agent tool** (subagent_type: `zcode:doc-scorer` if registered, otherwise `general-purpose`).

<HARD-RULE>
Pass these inputs to the scorer:
- `DOC_DIR` = `docs/features/<slug>/prd/`
- `RUBRIC_PATH` = `plugins/zcode/skills/eval-prd/templates/rubric.md`
- `REPORT_PATH` = `docs/features/<slug>/prd-eval-iteration-{{N}}.md`
- `ITERATION` = current iteration number (1-based)
- `PREVIOUS_REPORT_PATH` = previous iteration report path (only if iteration > 1)
</HARD-RULE>

After the scorer returns, **parse its output in the main session**:

1. Extract `SCORE: X/100` line
2. Extract per-dimension scores from `DIMENSIONS:` section
3. Extract attack points from `ATTACKS:` section
4. Record score in iteration tracker

## Step 3: Decision Gate (Main Session)

<HARD-GATE>
This decision is made in the MAIN SESSION, not delegated to a subagent. The orchestrator (you) controls the loop.
</HARD-GATE>

| Condition                                  | Action                          |
| ------------------------------------------ | ------------------------------- |
| Score >= target                            | Skip to Step 6 (final report)   |
| Score < target AND iterations remaining    | Proceed to Step 4               |
| Score < target AND no iterations remaining | Skip to Step 6 (report failure) |

Report current status to user:

```
Iteration {{N}}/{{MAX}}: scored {{SCORE}}/100 (target: {{TARGET}}). Revision subagent starting...
```

## Step 4: Invoke Reviser Subagent

Spawn `doc-reviser` agent via **Agent tool** (subagent_type: `zcode:doc-reviser` if registered, otherwise `general-purpose`).

<HARD-RULE>
Pass these inputs to the reviser:
- `DOC_DIR` = `docs/features/<slug>/prd/`
- `RUBRIC_PATH` = `plugins/zcode/skills/eval-prd/templates/rubric.md`
- `EVAL_REPORT_PATH` = `docs/features/<slug>/prd-eval-iteration-{{N}}.md`
- `ATTACK_POINTS` = the 3 attack points extracted from scorer output
</HARD-RULE>

The reviser will overwrite the PRD file(s) in place.

## Step 5: Loop

Increment iteration counter. Return to Step 2.

<EXTREMELY-IMPORTANT>
The scorer must NEVER be told what changes the reviser made. It evaluates the PRD as-is. Only the `PREVIOUS_REPORT_PATH` input carries forward for "previous issues addressed" checking.
</EXTREMELY-IMPORTANT>

## Step 6: Final Report (Main Session)

When the loop ends, assemble and report to the user:

```
## Eval-PRD Complete

**Final Score**: {{SCORE}}/100 (target: {{TARGET}})
**Iterations Used**: {{N}}/{{MAX}}

### Score Progression
| Iteration | Score | Delta |
|-----------|-------|-------|
| 1 | {{s1}} | - |
| 2 | {{s2}} | +{{d2}} |
| ... | ... | ... |

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

Save the final report to `docs/features/<slug>/prd-eval.md`.

## Report Path Convention

| File               | Path                                                |
| ------------------ | --------------------------------------------------- |
| Iteration N report | `docs/features/<slug>/prd-eval-iteration-{{N}}.md` |
| Final report       | `docs/features/<slug>/prd-eval.md`                 |

## Related

- `/write-prd` — Create or revise the PRD
- `/design-tech` — Next step after PRD passes evaluation
- `/ui-design` — Next step (optional) if prd-ui-functions.md exists
- `/breakdown-tasks` — After design docs are finalized
