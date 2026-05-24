---
id: "1"
title: "Implement Phase 0.5 Pre-Revision in SKILL.md"
priority: "P0"
estimated_time: "4h"
dependencies: []
type: "doc"
mainSession: false
---

# 1: Implement Phase 0.5 Pre-Revision in SKILL.md

## Description

After Phase 0 (freeform review + findings extraction) completes and before the Scorer cycle starts, add a Phase 0.5 Pre-Revision step. Format freeform findings as ATTACK_POINTS, construct a synthetic eval report (iteration: 0 + ATTACK_POINTS + empty rubric), invoke the existing Reviser subagent, annotate modified paragraphs with `<!-- pre-revised: {severity} -->` markers. Includes iteration initialization changes (ITERATION = 0, increment to 1 after pre-revision), BASELINE_SCORE evaluation (single Scorer call before pre-revision for informational comparison), two-level rollback (Scorer loop → pre-revised checkpoint, overall → Phase 0 baseline snapshot), tag lifecycle management (generation → survival → cleanup in Step 5 and Step 1.4), error handling for 4 failure scenarios with degradation to standard Scorer flow, and `--iterations 2` warning.

## Reference Files

- `proposal.md#Proposed-Solution` — defines the new information flow chain and Phase 0.5 architecture
- `proposal.md#Design-Decisions` — Decision 1 (pre-revision before Scorer), Decision 4 (reuse Reviser with synthetic report), Decision 5 (iteration 0 budget)
- `proposal.md#Scope` — 改动文件 table for SKILL.md, Phase 0.5 failure handling table, Non-Functional Requirements
- `proposal.md#Success-Criteria` — SC #1, #3, #4, #5, #6, #7 define acceptance criteria
- `proposal.md#Key-Risks` — pre-reviser mechanical response, pre-revision damage, LLM hallucination, INITIAL_SCORE baseline drift

## Affected Files

### Create
| File | Description |
|------|-------------|
| _(none)_ | |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/eval/SKILL.md` | Phase 0.5 Pre-Revision step after P0.4; Iteration Initialization change (ITERATION=0→1); architecture diagram update; Step 5 tag cleanup + Step 1.4 startup residual cleanup; BASELINE_SCORE single Scorer evaluation; two-level rollback checkpoint management; `--iterations 2` warning |

### Delete
| File | Reason |
|------|--------|
| _(none)_ | |

## Acceptance Criteria

1. **SC #1**: Phase 0 completes → findings auto-converted to ATTACK_POINTS → Reviser revision triggered before Scorer starts
2. **SC #3**: Pre-revision changes logged in iteration-0 report with title "Pre-Revision (Freeform Findings)"
3. **SC #4**: Final eval report contains "Pre-Revision" independent section listing each finding's status (accepted/partially-accepted/deferred/skipped) with edit summaries; skipped findings include classification rationale and original finding summary
4. **SC #5**: Existing degradation paths unaffected — Phase 0 failure or Phase 0.5 exception skips pre-revision, enters Scorer directly
5. **SC #6**: BASELINE_SCORE obtained via single Scorer call on original proposal before pre-revision (informational metric, not gate)
6. **SC #7**: High-severity findings triage rate >= 80% (accepted + partially-accepted + deferred), with accepted + partially-accepted >= 60%

## Hard Rules

- Only modify `plugins/forge/skills/eval/SKILL.md` in this task — scorer-composition.md and freeform-injection.md changes belong to Task 2
- Reviser protocol (`experts/protocol/reviser-protocol.md`) and composition (`rules/reviser-composition.md`) are NOT modified — construct synthetic eval report to satisfy Reviser's `EVAL_REPORT_PATH` dependency (see Decision 4)
- Use relative paths per forge-distribution.md (not project root paths like `plugins/forge/...`)
- `--iterations 1` skips pre-revision entirely (behavior unchanged)
- Pre-revision occupies iteration 0 from MAX_ITERATIONS budget

## Implementation Notes

- **Findings format for ATTACK_POINTS**: `- **[severity]** summary | 原文引用: "quote" | 期望改进方向: <动词短语>`
- **Synthetic eval report**: `iteration: 0` + ATTACK_POINTS list + rubric all dimensions N/A — satisfies Reviser protocol Step 1 minimum input format
- **Error handling (4 scenarios)**:
  - Findings formatting failure → skip pre-revision, enter Scorer directly
  - Pre-reviser returns error → skip, log warning
  - Empty report produced → log iteration-0 "no changes", Scorer starts normally
  - Format anomaly → discard pre-revision results, use original proposal for Scorer
- **Two-level rollback**: (1) Scorer loop rollback → pre-revised checkpoint; (2) overall rollback → Phase 0 baseline snapshot saved at `<DOC_DIR>/eval/baseline-snapshot/`
- **Tag lifecycle**: Generation (SKILL.md inserts after Reviser edits) → Survival (visible during Scorer review) → Cleanup (Step 5 strips all `<!-- pre-revised -->` after eval; Step 1.4 clears residuals at startup)
- **BASELINE_SCORE**: single Scorer subagent call (no Reviser) on original proposal before pre-revision — informational only, does not consume iteration budget
- **`--iterations 2` warning**: output when freeform review is active and iterations <= 2 — "Pre-revision 占用 1 个 iteration，Scorer 仅执行 1 轮评估。建议使用 `--iterations 3` 保证 Scorer 有修正机会。"
- **Baseline drift detection**: if iteration-1 INITIAL_SCORE < BASELINE_SCORE by >50 points (1000-scale), annotate eval report with "基线漂移告警" for manual review — does not auto-trigger rollback
- **Three-layer finding triage**: factual correction (direct edit) → structural suggestion (edit only on contradiction, else defer) → subjective preference (mark not actionable, log in iteration-0 report)
- **Borderline handling**: uncertain findings marked "borderline" and deferred (not silently classified), listed separately in iteration-0 report
- **Architecture diagram**: update mermaid flowchart to show Phase 0.5 between Phase 0 and Expert Dispatch
