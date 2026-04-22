# Design: Adversarial Evaluation Loop for eval-design and eval-prd

**Date:** 2026-04-22  
**Status:** Approved

---

## Summary

Bring `eval-design` and `eval-prd` to full parity with `eval-proposal`: numeric 100-point scoring, adversarial iteration loop, and dedicated subagents. Simultaneously refactor `eval-proposal` to use a shared `doc-scorer` / `doc-reviser` agent pair driven by per-skill `rubric.md` files, replacing the current proposal-specific agents.

---

## Architecture

All three eval skills share the same orchestration pattern:

```
Main Skill (orchestrator)
  │
  ├─ iteration N:
  │   ├── Agent (doc-scorer)   → SCORE + DIMENSIONS + ATTACKS
  │   ├── score >= target? → yes: stop
  │   └── Agent (doc-reviser)  → revised source doc(s)
  │
  └─ Final report to user
```

### Shared Agents

Two new generic agents replace `proposal-scorer` and `proposal-reviser`:

**`doc-scorer`**
- Inputs: `DOC_PATHS` (comma-separated), `RUBRIC_PATH`, `REPORT_PATH`, `ITERATION`, `PREVIOUS_REPORT_PATH` (optional)
- Reads the rubric at `RUBRIC_PATH` to understand scoring dimensions and point allocations
- Scores the document(s), writes the report, returns structured output the orchestrator parses
- Output format identical to current proposal-scorer (SCORE / DIMENSIONS / ATTACKS blocks)

**`doc-reviser`**
- Inputs: `DOC_PATHS` (comma-separated), `RUBRIC_PATH`, `EVAL_REPORT_PATH`, `ATTACK_POINTS`
- Reads the rubric to understand what "good" looks like before revising
- Overwrites source doc(s) in place
- Returns `REVISED: / CHANGES:` block

### Per-Skill Rubric Files

Each skill owns its rubric as a standalone file:

| Skill | Rubric path |
|-------|-------------|
| eval-proposal | `skills/eval-proposal/templates/rubric.md` |
| eval-design | `skills/eval-design/templates/rubric.md` |
| eval-prd | `skills/eval-prd/templates/rubric.md` |

Rubric files contain: dimension names, point allocations, per-criterion scoring tables, and deduction rules. Agents read the rubric at runtime — no rubric content is hardcoded in agent prompts.

---

## Scoring Model

### eval-design (100 pts)

| Dimension | Max | Notes |
|-----------|-----|-------|
| Architecture Clarity | 20 | Layer placement, diagram, dependencies |
| Interface & Model Definitions | 20 | Typed sigs, concrete models, implementable |
| Error Handling | 15 | Error types, propagation, HTTP mapping |
| Testing Strategy | 15 | Per-layer, coverage target, tooling |
| Breakdown-Readiness ★ | 20 | Critical gate — weighted higher than current |
| Security Considerations | 10 | N/A if no security surface |

★ Breakdown-Readiness raised from equal weight to 20 pts to reflect its gate role.

### eval-prd (100 pts)

| Dimension | Max | Notes |
|-----------|-----|-------|
| Background & Goals | 20 | Three elements, quantified targets |
| Flow Diagrams | 20 | Mermaid, main path, decision points, error branches |
| Functional Specs | 20 | Table completeness, field clarity, validation rules |
| User Stories | 20 | Coverage, format, AC per story |
| Scope Clarity | 20 | In/out concrete, consistent with specs |

### eval-proposal (unchanged rubric, migrated to rubric.md)

Existing 6-dimension / 100-point rubric extracted verbatim from `proposal-scorer.md` into `skills/eval-proposal/templates/rubric.md`.

---

## Skill Changes

### eval-proposal/SKILL.md

- Replace `subagent_type: zcode:proposal-scorer` → `zcode:doc-scorer`
- Replace `subagent_type: zcode:proposal-reviser` → `zcode:doc-reviser`
- Add `RUBRIC_PATH = plugins/zcode/skills/eval-proposal/templates/rubric.md` to both agent invocations
- No behavior change — same loop, same parameters, same report paths

### eval-design/SKILL.md

- Add `--target` and `--iterations` parameters (defaults: 80, 3)
- Replace single-pass agent invocation with orchestration loop (Steps 2–5 from eval-proposal)
- Pass `DOC_PATHS` = design.md + api-handbook.md + ui-design.md (comma-separated, skip missing)
- Pass `RUBRIC_PATH = plugins/zcode/skills/eval-design/templates/rubric.md`
- Report paths: `docs/features/<slug>/design-eval-iteration-{{N}}.md`, final: `design-eval.md`
- Final report includes score progression table (same format as eval-proposal)

### eval-prd/SKILL.md

- Add `--target` and `--iterations` parameters (defaults: 80, 3)
- Replace single-pass agent invocation with orchestration loop
- Pass `DOC_PATHS` = prd-spec.md + prd-user-stories.md + prd-ui-functions.md (skip missing)
- Pass `RUBRIC_PATH = plugins/zcode/skills/eval-prd/templates/rubric.md`
- Report paths: `docs/features/<slug>/prd-eval-iteration-{{N}}.md`, final: `prd-eval.md`
- Final report includes score progression table

---

## Report Templates

Existing `eval-design/templates/report.md` and `eval-prd/templates/report.md` are updated to match the `eval-proposal` scorecard format: ASCII scorecard table with per-dimension numeric scores, Deductions section, Attack Points section, Previous Issues Check (iteration > 1), and Verdict.

---

## File Changeset

| Action | Path |
|--------|------|
| Create | `plugins/zcode/agents/doc-scorer.md` |
| Create | `plugins/zcode/agents/doc-reviser.md` |
| Create | `plugins/zcode/skills/eval-proposal/templates/rubric.md` |
| Create | `plugins/zcode/skills/eval-design/templates/rubric.md` |
| Create | `plugins/zcode/skills/eval-prd/templates/rubric.md` |
| Update | `plugins/zcode/skills/eval-proposal/SKILL.md` |
| Update | `plugins/zcode/skills/eval-design/SKILL.md` |
| Update | `plugins/zcode/skills/eval-prd/SKILL.md` |
| Update | `plugins/zcode/skills/eval-design/templates/report.md` |
| Update | `plugins/zcode/skills/eval-prd/templates/report.md` |
| Delete | `plugins/zcode/agents/proposal-scorer.md` |
| Delete | `plugins/zcode/agents/proposal-reviser.md` |

---

## Constraints

- `doc-scorer` and `doc-reviser` must be registered in `plugin.json` as named agents so skills can reference them via `subagent_type`
- Rubric files are plain markdown — no frontmatter required, agents read them as reference text
- The scorer must never be told what the reviser changed (same isolation rule as eval-proposal)
- `DOC_PATHS` is comma-separated; agents skip paths that don't exist on disk
