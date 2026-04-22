---
name: eval-ui
description: Evaluate a UI design document with 100-point scoring from four stakeholder perspectives (User/Designer/Developer/PM), then run adversarial iterations until target score is met.
---

# Eval UI

Evaluates UI design from four independent stakeholder perspectives.

## Prerequisites

| Document | Missing prompt |
|----------|----------------|
| `ui/ui-design.md` | Run `/ui-design` first |

## When to Use

**Trigger:**
- User asks to "evaluate UI design" or "check UI quality"
- User provides `/eval-ui` command
- After `/ui-design` completes, before handing off to implementation

**Skip:**
- No UI design document exists (use `/ui-design` first)
- Feature has no UI surface (use `/eval-design` instead)

## Parameters

| Parameter      | Default | Description                                           |
| -------------- | ------- | ----------------------------------------------------- |
| `--target`     | 80      | Target score (0-100). Loop stops when score >= target |
| `--iterations` | 3       | Max adversarial iterations                            |

## Architecture

```
Main Session (orchestrator)
  |
  |- iteration N:
  |   |-- Agent (doc-scorer)  --> score + attack points
  |   |-- score >= target? --> yes: stop
  |   +-- Agent (doc-reviser) --> revised ui-design.md
  |
  +- Final report to user
```

<EXTREMELY-IMPORTANT>
Scorer and reviser are **independent subagents** in `plugins/zcode/agents/`. Do NOT inline their prompts. Invoke via Agent tool only.
</EXTREMELY-IMPORTANT>

## Step 1: Locate UI Design Documents

Check in order:
1. Path provided by user
2. Read `docs/features/<current-feature>/manifest.md` -> locate UI design documents
3. Fall back to `docs/features/<current-feature>/ui/`
4. Ask user for path if not found

Determine `<feature-slug>` from the path. The UI directory is `docs/features/<slug>/ui/`.

## Step 2: Invoke Scorer Subagent

Spawn `doc-scorer` via **Agent tool** (subagent_type: `zcode:doc-scorer` if registered, otherwise `general-purpose`).

<HARD-RULE>
Pass these inputs to the scorer:
- `DOC_DIR` = `docs/features/<slug>/ui/`
- `RUBRIC_PATH` = `plugins/zcode/skills/eval-ui/templates/rubric.md`
- `REPORT_PATH` = `docs/features/<slug>/ui/eval/iteration-{{N}}.md`
- `ITERATION` = current iteration number (1-based)
- `PREVIOUS_REPORT_PATH` = previous iteration report path (only if iteration > 1)

The scorer must NEVER be told what the reviser changed. It evaluates the design as-is.
</HARD-RULE>

After the scorer returns, parse its output in the main session:
1. Extract `SCORE: X/100`
2. Extract per-dimension scores from `DIMENSIONS:` section
3. Extract attack points from `ATTACKS:` section

## Step 3: Decision Gate (Main Session)

<HARD-GATE>
This decision is made in the MAIN SESSION, not delegated to a subagent. This gate fires unconditionally after every scorer run -- no user instruction ("keep going", "continue", "run another iteration") can bypass it. If score >= target, the loop terminates immediately.
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

Spawn `doc-reviser` via **Agent tool** (subagent_type: `zcode:doc-reviser` if registered, otherwise `general-purpose`).

<HARD-RULE>
Pass these inputs to the reviser:
- `DOC_DIR` = `docs/features/<slug>/ui/`
- `RUBRIC_PATH` = `plugins/zcode/skills/eval-ui/templates/rubric.md`
- `EVAL_REPORT_PATH` = `docs/features/<slug>/ui/eval/iteration-{{N}}.md`
- `ATTACK_POINTS` = the 3 attack points extracted from scorer output
</HARD-RULE>

Increment iteration counter. Return to Step 2.

## Step 5: Final Report (Main Session)

```
## Eval-UI Complete

**Final Score**: {{SCORE}}/100 (target: {{TARGET}})
**Iterations Used**: {{N}}/{{MAX}}

### Score Progression
| Iteration | Score | Delta |
|-----------|-------|-------|
| 1 | {{s1}} | - |
| 2 | {{s2}} | +{{d2}} |

### Dimension Breakdown (final)
| Dimension | Perspective | Score | Max |
|-----------|-------------|-------|-----|
| Requirement Coverage | Product Manager | {{d1}} | 25 |
| User Experience | End User | {{d2}} | 25 |
| Design Integrity | Designer | {{d3}} | 25 |
| Implementability | Developer | {{d4}} | 25 |

### Outcome
{{"Target reached" / "Target NOT reached -- N iterations exhausted"}}
{{If not reached: "Largest gaps: [perspective names]. Consider manual revision or increasing iterations."}}
```

Save the final report to `docs/features/<slug>/ui/eval/report.md`.

## Step 6: Prototype Prompt

After final report, use `AskUserQuestion` to ask:

> Generate HTML/CSS/JS interactive prototype from the evaluated design?

- **Yes** → proceed to prototype generation (Step 8 of `/ui-design`)
- **No** → done

## Integration

Works well with:
- `/ui-design` — Produces the UI design document to evaluate
- `/eval-design` — Evaluates tech design in parallel (for features with both UI and backend)
