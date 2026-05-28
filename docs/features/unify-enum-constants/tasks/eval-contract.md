---
id: "T-eval-contract"
title: "Evaluate Contract Quality"
priority: "P1"
estimated_time: "20-30min"
dependencies: ["T-test-gen-contracts"]
type: "eval.contract"
surface-key: ""
surface-type: ""
mainSession: true
---

Evaluate Contract quality for the unify-enum-constants feature using the 6-dimension rubric (1000-point scale).

## Discovery Strategy
Scan `tests/<journey>/_contracts/` for all Contract files per Journey.

For each Journey's Contracts:
1. Run `/eval --type contract` using the contract rubric (`eval/rubrics/contract.md`)
2. Scoring dimensions: Completeness, Semantic Purity, Precondition Exclusivity, Fact Alignment, Surface Fitness, Internal Consistency
3. Target score: 850/1000 with all dimensions above min thresholds

If any Contract fails evaluation after max iterations, report the failure and abort. Do not proceed to gen-test-scripts with low-quality Contracts.

## Acceptance Criteria
- [ ] All Contracts scored >= 850/1000
- [ ] All dimensions above min threshold per rubric
- [ ] Eval report written to `tests/<journey>/.eval-report.md`
