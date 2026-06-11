---
id: "3"
title: "Expand scorer-protocol self-contradiction check with clustering + satisfiability"
priority: "P0"
estimated_time: "1h"
dependencies: []
type: "doc"
mainSession: false
---

# 3: Expand scorer-protocol self-contradiction check with clustering + satisfiability

## Description

Modify `plugins/forge/skills/eval/experts/protocol/scorer-protocol.md` Phase 1 (Reasoning Audit) Step 4 self-contradiction check to include explicit clustering + intra-group SC↔SC and SC↔InScope satisfiability check instructions. This is Layer 2 (eval detection) of the two-layer defense.

## Reference Files

- `proposal.md#Proposed-Solution` — Layer 2 (eval detection): expands scorer-protocol Phase 1 self-contradiction check
- `proposal.md#Requirements-Analysis` — Scenario 2 (eval detection): scorer detects contradictions via clustering and generates attack points; Scenario 5 (clustering error): fallback mechanism
- `proposal.md#Key-Risks` — LLM false negative risk; eval layer differentiation strategy (broader search prompt, higher temperature 0.7 vs brainstorm 0.3)
- `proposal.md#Success-Criteria` — SC requiring scorer-protocol to contain explicit clustering + satisfiability check instructions with gen-and-run contradiction as example

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/eval/experts/protocol/scorer-protocol.md` | Expand Phase 1 Step 4 self-contradiction check with clustering + satisfiability protocol |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria

- [ ] scorer-protocol.md Phase 1 Step 4 (self-contradiction check) contains explicit clustering instruction: group SC entries by affected area (file/directory/module)
- [ ] Contains intra-group satisfiability check instruction: for each cluster, execute bidirectional SC↔SC and SC↔InScope satisfiability derivation
- [ ] References the gen-and-run contradiction scenario (grep zero-result vs preserve migration prompt) as an example use case
- [ ] Contradictions found are tagged as attack points requiring reviser revision
- [ ] Revised SC must re-pass consistency check (re-cluster + intra-group check) to avoid introducing new contradictions
- [ ] Eval layer differentiation: uses broader search prompt (not limited to area clustering) and optionally higher temperature for reasoning diversity

## Hard Rules

- Follow `docs/conventions/forge-distribution.md` — use relative paths within the skill directory
- Extend the existing self-contradiction check sub-step, do NOT replace or restructure the overall Phase 1 workflow

## Implementation Notes

- The eval layer is independent from the brainstorm layer — scorer-protocol should not reference brainstorm rules
- The eval layer can use a broader scanning approach (not limited to clustering) since it runs less frequently and serves as a safety net
- The key differentiator from brainstorm layer: eval uses full-pair scanning as the primary method, with clustering as an optimization hint, rather than the brainstorm layer's clustering-first approach
