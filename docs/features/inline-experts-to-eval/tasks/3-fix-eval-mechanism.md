---
id: "3"
title: "Fix 5 eval mechanism issues in SKILL.md"
priority: "P1"
estimated_time: "1h"
dependencies: ["2"]
type: "documentation"
mainSession: false
---

# 3: Fix 5 eval mechanism issues in SKILL.md

## Description
Five logical inconsistencies and sequencing gaps in the eval SKILL.md that affect multi-expert reviser dispatch, iteration initialization, gate override, context injection, and score extraction. All five fixes target `skills/eval/SKILL.md`.

## Reference Files
- `docs/proposals/inline-experts-to-eval/proposal.md` — Source proposal (Part B: Fix Eval Mechanism Issues)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/eval/SKILL.md` | Apply 5 mechanism fixes across Steps 1–4 |

## Acceptance Criteria
- [ ] **Fix 1 — Multi-expert reviser EVAL_REPORT_PATH**: After Step 2.3's LLM merge, main session writes merged report to `<doc_dir>/eval/iteration-{{N}}-merged.md`. This merged report serves as `EVAL_REPORT_PATH` for the reviser. Single-expert types continue using `iteration-{{N}}.md` directly.
- [ ] **Fix 2 — Iteration counter initialization**: After Step 1, add explicit initialization: "Set `ITERATION = 1`, `MAX_ITERATIONS = resolved value from rubric or CLI`."
- [ ] **Fix 3 — Remove ambiguous "continue" override**: Remove the "On 'continue'/'keep going'" line from Step 3b. Gate decisions are purely score-driven.
- [ ] **Fix 4 — Context injection for reviser**: Apply the same context injection block from Step 2.1 to Step 4.1. Append `<injected-context>...</injected-context>` to the reviser prompt when `CONTEXT_CONTENT` was loaded.
- [ ] **Fix 5 — Score extraction robustness**: In Step 2.3, add extraction instruction: "Extract score using regex `/SCORE:\s*(\d+)\/(\d+)/`. If pattern not found, scan the scorer agent's output for the last line matching a `number/number` pattern. If still not found, report error and stop."

## Hard Rules
- Each fix is scoped to the exact location described — do not refactor surrounding text
- Fix 1 must preserve single-expert behavior unchanged

## Implementation Notes
- The proposal provides exact change descriptions for each fix under "Mechanism Fixes in SKILL.md (Part B)"
- Fix 2: add after Step 1's existing content, before Step 2 header
- Fix 3: single line deletion in Step 3b
- Fix 4: copy the `<injected-context>` block from Step 2.1 and adapt for reviser context
- Fix 5: replace the "extract directly" instruction in Step 2.3 with the robust extraction logic
