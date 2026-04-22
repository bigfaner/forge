---
name: eval-design
description: Evaluate a tech design document with 100-point scoring, then run adversarial iterations until target score is met. Main session orchestrates doc-scorer and doc-reviser subagents.
---

# Eval Design

评估 tech-design.md 文档质量（百分制），通过 doc-scorer / doc-reviser subagent 多轮对抗迭代，直到达到目标分数或耗尽迭代次数。重点检查能否直接驱动 `/breakdown-tasks`。

**架构**：主会话调度循环，独立 subagent 负责评分和修订。

## When to Use

**Trigger:**
- User asks to "evaluate design" or "check design quality"
- User provides `/eval-design` command
- Before handing off tech-design.md to `/breakdown-tasks`

**Skip:**
- design.md doesn't exist yet (use `/design-tech` first)

## Parameters

| Parameter      | Default | Description                                           |
| -------------- | ------- | ----------------------------------------------------- |
| `--target`     | 80      | Target score (0-100). Loop stops when score >= target |
| `--iterations` | 3       | Max adversarial iterations                            |

Parse from user input. Examples:

- `/eval-design` → target=80, iterations=3
- `/eval-design --target 90` → target=90, iterations=3
- `/eval-design --target 85 --iterations 5` → target=85, iterations=5

## Architecture

```
Main Session (orchestrator)
  │
  ├─ iteration N:
  │   ├── Agent (doc-scorer)  ──→ score + attack points
  │   ├── score >= target? ──→ yes: stop
  │   └── Agent (doc-reviser) ──→ revised design doc(s)
  │
  └─ Final report to user
```

<EXTREMELY-IMPORTANT>
Scorer and reviser are **independent subagents** defined in `plugins/zcode/agents/`. Do NOT inline their prompts. Invoke them via Agent tool with the correct inputs.
</EXTREMELY-IMPORTANT>

## Step 1: Locate Design Documents

Check in order:
1. Path provided by user
2. Read `docs/features/<current-feature>/manifest.md` → locate design documents
3. Fall back to `design/tech-design.md`, `design/api-handbook.md`, `ui/ui-design.md`
4. Ask user for path if not found

Determine `<feature-slug>` from the path. The design directory is `docs/features/<slug>/design/` (or wherever the design files live).

## Step 2: Invoke Scorer Subagent

Spawn `doc-scorer` agent via **Agent tool** (subagent_type: `zcode:doc-scorer` if registered, otherwise `general-purpose`).

<HARD-RULE>
Pass these inputs to the scorer:
- `DOC_DIR` = `docs/features/<slug>/design/`
- `RUBRIC_PATH` = `plugins/zcode/skills/eval-design/templates/rubric.md`
- `REPORT_PATH` = `docs/features/<slug>/design/eval/iteration-{{N}}.md`
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
- `DOC_DIR` = `docs/features/<slug>/design/`
- `RUBRIC_PATH` = `plugins/zcode/skills/eval-design/templates/rubric.md`
- `EVAL_REPORT_PATH` = `docs/features/<slug>/design/eval/iteration-{{N}}.md`
- `ATTACK_POINTS` = the 3 attack points extracted from scorer output
</HARD-RULE>

The reviser will overwrite the design file(s) in place.

## Step 5: Loop

Increment iteration counter. Return to Step 2.

<EXTREMELY-IMPORTANT>
The scorer must NEVER be told what changes the reviser made. It evaluates the design as-is. Only the `PREVIOUS_REPORT_PATH` input carries forward for "previous issues addressed" checking.
</EXTREMELY-IMPORTANT>

## Step 6: Final Report (Main Session)

When the loop ends, assemble and report to the user:

```
## Eval-Design Complete

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

## Report Path Convention

| File               | Path                                                   |
| ------------------ | ------------------------------------------------------ |
| Iteration N report | `docs/features/<slug>/design/eval/iteration-{{N}}.md` |
| Final report       | `docs/features/<slug>/design/eval/report.md`          |

## Related

- `/design-tech` — Create or revise the design.md
- `/eval-prd` — Evaluate PRD before design starts
- `/breakdown-tasks` — Next step after design passes evaluation
