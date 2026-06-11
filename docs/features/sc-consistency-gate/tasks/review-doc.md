---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["4"]
type: "doc.review"
scope: "all"
---

Review documentation quality for the sc-consistency-gate feature (quick mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 1-create-sc-consistency-rule

- [ ] Rule file exists at `plugins/forge/skills/brainstorm/rules/sc-consistency.md`
- [ ] Contains clustering protocol: group SC and InScope entries by affected area (file/directory/module)
- [ ] Contains intra-group satisfiability check protocol: for each pair within a group, execute bidirectional proof (assume A true → derive B state; assume B true → derive A state)
- [ ] Contains fallback cross-group direction check: after intra-group checks, run a lightweight all-pair scan for ADD vs SUBTRACT on same symbol across groups
- [ ] References the pipeline-integration-stitch contradiction case as an example (grep zero-result vs preserve migration prompt)
- [ ] Includes explicit rule: contradiction-free SC sets produce zero output (empty report)
- [ ] Includes handling for ambiguous contradictions: mark as "ambiguous — requires user confirmation" instead of forcing a binary choice
- [ ] Structured output format: for each contradiction, output conflict pair, type (mutual exclusion / direction conflict / resource competition), and suggested resolution


### 2-add-skill-reference

- [ ] SKILL.md Step 5 contains an explicit reference to `rules/sc-consistency.md`
- [ ] The reference is positioned after SC and InScope writing, before the quality standards table
- [ ] Consistency check is described as a mandatory step (not optional), aligning with the "hard protection" strategy from the proposal


### 3-expand-scorer-protocol

- [ ] scorer-protocol.md Phase 1 Step 4 (self-contradiction check) contains explicit clustering instruction: group SC entries by affected area (file/directory/module)
- [ ] Contains intra-group satisfiability check instruction: for each cluster, execute bidirectional SC↔SC and SC↔InScope satisfiability derivation
- [ ] References the gen-and-run contradiction scenario (grep zero-result vs preserve migration prompt) as an example use case
- [ ] Contradictions found are tagged as attack points requiring reviser revision
- [ ] Revised SC must re-pass consistency check (re-cluster + intra-group check) to avoid introducing new contradictions
- [ ] Eval layer differentiation: uses broader search prompt (not limited to area clustering) and optionally higher temperature for reasoning diversity


### 4-adjust-proposal-rubric-d9

- [ ] D9 contains "SC internal consistency" criterion worth 25pts with clear evaluation guidance (check SC↔SC and SC↔InScope for logical contradictions within clusters)
- [ ] "Criteria are measurable and testable" reduced from 55pts to 30pts
- [ ] "Coverage is complete" reduced from 25pts to 25pts (unchanged, as proposal says 40→25 for coverage but current rubric shows 25 — verify against actual current value)
- [ ] D9 total remains 80pts
- [ ] New criterion description checks SC internal satisfiability (intra-group SC↔SC and SC↔InScope), distinct from D10 which checks SC ↔ Scope/Solution alignment (no overlap)


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/sc-consistency-gate/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/sc-consistency-gate/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.
