---
name: eval-proposal
description: Evaluate a proposal document with 100-point scoring, then run adversarial iterations until target score is met. Main session orchestrates proposal-scorer and proposal-reviser subagents. Specify target score and max iterations.
---

# Eval Proposal

评估 proposal 文档质量（百分制），通过 proposal-scorer / proposal-reviser subagent 多轮对抗迭代，直到达到目标分数或耗尽迭代次数。

**架构**：主会话调度循环，独立 subagent 负责评分和修订。

## When to Use

**Trigger:**

- User says yes to adversarial eval prompt after `/brainstorm`
- User provides `/eval-proposal` command
- User wants iterative refinement: `/eval-proposal --target 85 --iterations 5`

**Skip:**

- No proposal document exists (use `/brainstorm` first)
- Requirements are already in PRD form (use `/eval-prd` instead)

## Parameters

| Parameter      | Default | Description                                           |
| -------------- | ------- | ----------------------------------------------------- |
| `--target`     | 80      | Target score (0-100). Loop stops when score >= target |
| `--iterations` | 3       | Max adversarial iterations                            |

Parse from user input. Examples:

- `/eval-proposal` → target=80, iterations=3
- `/eval-proposal --target 90` → target=90, iterations=3
- `/eval-proposal --target 85 --iterations 5` → target=85, iterations=5

## Architecture

```
Main Session (orchestrator)
  │
  ├─ iteration 1:
  │   ├── Agent (proposal-scorer)  ──→ score + attack points
  │   ├── score >= target? ──→ yes: stop
  │   └── Agent (proposal-reviser) ──→ revised proposal
  │
  ├─ iteration 2:
  │   ├── Agent (proposal-scorer)  ──→ score + attack points  (blind to changes)
  │   ├── score >= target? ──→ yes: stop
  │   └── Agent (proposal-reviser) ──→ revised proposal
  │
  ├─ ... (loop)
  │
  └─ Final report to user
```

<EXTREMELY-IMPORTANT>
Scorer and reviser are **independent subagents** defined in `plugins/zcode/agents/`. Do NOT inline their prompts. Invoke them via Agent tool with the correct inputs.
</EXTREMELY-IMPORTANT>

## Step 1: Locate Proposal

Check in order:

1. Path provided by user
2. `docs/proposals/<slug>/proposal.md` — find latest by modification time
3. `Glob` `docs/proposals/*/proposal.md` and list options to user
4. Ask user for path if not found

Determine `<slug>` from path (e.g., `docs/proposals/eval-proposal/proposal.md` → slug is `eval-proposal`).

## Step 2: Invoke Scorer Subagent

Spawn `proposal-scorer` agent via **Agent tool** (subagent_type: `zcode:proposal-scorer` if registered, otherwise `general-purpose`).

<HARD-RULE>
Pass these inputs to the scorer:
- `PROPOSAL_PATH` = `docs/proposals/<slug>/proposal.md`
- `REPORT_PATH` = `docs/proposals/<slug>/eval-iteration-{{N}}.md`
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

Spawn `proposal-reviser` agent via **Agent tool** (subagent_type: `zcode:proposal-reviser` if registered, otherwise `general-purpose`).

<HARD-RULE>
Pass these inputs to the reviser:
- `PROPOSAL_PATH` = `docs/proposals/<slug>/proposal.md`
- `EVAL_REPORT_PATH` = `docs/proposals/<slug>/eval-iteration-{{N}}.md`
- `ATTACK_POINTS` = the 3 attack points extracted from scorer output
</HARD-RULE>

The reviser will overwrite the proposal file in place.

## Step 5: Loop

Increment iteration counter. Return to Step 2.

<EXTREMELY-IMPORTANT>
The scorer must NEVER be told what changes the reviser made. It evaluates the proposal as-is. Only the `PREVIOUS_REPORT_PATH` input carries forward for "previous issues addressed" checking.
</EXTREMELY-IMPORTANT>

## Step 6: Final Report (Main Session)

When the loop ends, assemble and report to the user:

```
## Eval-Proposal Complete

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
| Problem Definition | {{d1}} | 20 |
| Solution Clarity | {{d2}} | 20 |
| Alternatives Analysis | {{d3}} | 15 |
| Scope Definition | {{d4}} | 15 |
| Risk Assessment | {{d5}} | 15 |
| Success Criteria | {{d6}} | 15 |

### Outcome
{{"Target reached" / "Target NOT reached — N iterations exhausted"}}
{{If not reached: "Largest gaps: [dimension names]. Consider manual revision or increasing iterations."}}
```

Save the final report to `docs/proposals/<slug>/eval-report.md`.

## Report Path Convention

| File               | Path                                            |
| ------------------ | ----------------------------------------------- |
| Iteration N report | `docs/proposals/<slug>/eval-iteration-{{N}}.md` |
| Final report       | `docs/proposals/<slug>/eval-report.md`          |

## Related

- `/brainstorm` — Creates or revises the proposal document (runs in main session)
