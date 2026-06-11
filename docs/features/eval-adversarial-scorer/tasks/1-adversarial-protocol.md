---
id: "1"
title: "Implement three-phase adversarial protocol in doc-scorer"
priority: "P0"
estimated_time: "1.5h"
dependencies: []
scope: "backend"
breaking: false
type: "enhancement"
mainSession: false
---

# 1: Implement three-phase adversarial protocol in doc-scorer

## Description

Rewrite `plugins/forge/agents/doc-scorer.md` with a layered adversarial protocol that wraps around existing rubric scoring. The current scorer mechanically applies rubric dimensions without independently challenging the document's core reasoning. This produces high scores on documents with fundamental flaws (e.g., 958/1000 on a refactoring that reintroduced the pattern it aimed to eliminate).

The protocol has three phases:
1. **Phase 1 — Reasoning Audit**: Before rubric, trace the argument chain (problem → solution → evidence → success criteria). Check if the solution reintroduces what it claims to eliminate. Record findings as pre-score anchors.
2. **Phase 2 — Rubric Scoring with Verification Stance**: Existing dimension-by-dimension evaluation, but every assertion is treated as unverified until proven. Cross-dimension coherence check: verify scope, solution, and success criteria are internally consistent.
3. **Phase 3 — Blindspot Hunt**: After scoring, ask "what did the rubric miss?" Tag findings as `[blindspot]` — strictly for issues outside all rubric dimensions.

## Reference Files
- `docs/proposals/eval-adversarial-scorer/proposal.md` — Source proposal
- `plugins/forge/agents/doc-scorer.md` — Current scorer to rewrite
- `docs/conventions/forge-distribution.md` — MUST read before modifying plugins/forge/ files

## Acceptance Criteria

- [ ] Re-evaluating the design doc that had `ViaSubmit bool` produces an attack point flagging the disguised patch
- [ ] Evaluating a document where scope claims X but success criteria only test Y produces a `[blindspot]` attack flagging the cross-dimension gap
- [ ] Re-evaluating the UI design with page gap issues produces cross-page coherence attack points
- [ ] Output format remains `SCORE:`, `DIMENSIONS:`, `ATTACKS:` — parseable by existing orchestrator
- [ ] Scorer instruction explicitly states: "Find REAL issues. A document with no flaws deserves full marks. Manufacturing issues wastes the reviser's time and yours."
- [ ] `[blindspot]` attacks must cite a specific quote from the document; attacks without quotes are discarded
- [ ] Blindspot attacks are strictly for issues outside all rubric dimensions — if an issue fits a dimension, it is scored there instead

## Hard Rules

- MUST read `docs/conventions/forge-distribution.md` before modifying `plugins/forge/agents/doc-scorer.md`
- Output format (`SCORE:`, `DIMENSIONS:`, `ATTACKS:`) must remain parseable by the existing orchestrator — no structural changes to the return format
- No rubric file changes — all 17 rubrics stay as-is

## Implementation Notes

- The prompt grows by ~200 tokens — negligible compared to document reading cost
- Phase 2 rubric scores take precedence for `DIMENSIONS:` output; Phase 1 findings channel into `[blindspot]` attacks only if they identify issues not covered by any dimension
- If Phase 1 flags a fundamental flaw but Phase 2 scores that dimension well, the Phase 1 finding appears as a `[blindspot]` attack with notation: "Reasoning audit flagged this independently of dimension scoring."
- The adversarial protocol must work even when context is compacted across eval iterations — the prompt is in the system instructions, not conversation history
