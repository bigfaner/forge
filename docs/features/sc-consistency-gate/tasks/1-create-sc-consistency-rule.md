---
id: "1"
title: "Create sc-consistency.md rule file with clustering + satisfiability check"
priority: "P0"
estimated_time: "1h"
dependencies: []
type: "doc"
mainSession: false
---

# 1: Create sc-consistency.md rule file with clustering + satisfiability check

## Description

Create `plugins/forge/skills/brainstorm/rules/sc-consistency.md` — the core SC Consistency Check rule file used by brainstorm Step 5. This rule implements the clustering + intra-group satisfiability check protocol described in the proposal, plus the fallback cross-group direction check.

The rule enables brainstorm agents to detect logical contradictions between SC entries (SC↔SC) and between SC and InScope entries (SC↔InScope) before the proposal is committed.

## Reference Files

- `proposal.md#Proposed-Solution` — defines the two-layer defense model and the clustering + intra-group check algorithm
- `proposal.md#Requirements-Analysis` — Key Scenarios (prevention, detection, normal pass, ambiguous contradiction, clustering error) that the rule must handle
- `proposal.md#Key-Risks` — risk of agent ignoring rules and the hard-protection strategy (structured output field)
- `proposal.md#Innovation-Highlights` — clustering heuristic rationale and performance characteristics

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/brainstorm/rules/sc-consistency.md` | SC consistency check rule file |

### Modify
| File | Changes |
|------|---------|

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria

- [ ] Rule file exists at `plugins/forge/skills/brainstorm/rules/sc-consistency.md`
- [ ] Contains clustering protocol: group SC and InScope entries by affected area (file/directory/module)
- [ ] Contains intra-group satisfiability check protocol: for each pair within a group, execute bidirectional proof (assume A true → derive B state; assume B true → derive A state)
- [ ] Contains fallback cross-group direction check: after intra-group checks, run a lightweight all-pair scan for ADD vs SUBTRACT on same symbol across groups
- [ ] References the pipeline-integration-stitch contradiction case as an example (grep zero-result vs preserve migration prompt)
- [ ] Includes explicit rule: contradiction-free SC sets produce zero output (empty report)
- [ ] Includes handling for ambiguous contradictions: mark as "ambiguous — requires user confirmation" instead of forcing a binary choice
- [ ] Structured output format: for each contradiction, output conflict pair, type (mutual exclusion / direction conflict / resource competition), and suggested resolution

## Hard Rules

- Follow `docs/conventions/forge-distribution.md` distribution constraints — use relative paths within the skill directory, no project-root paths
- Bidirectional proof prompt structure: "assume A true → derive B state; assume B true → derive A state" rather than directly asking "are these contradictory"

## Implementation Notes

- The rule file lives in `brainstorm/rules/` and is referenced by SKILL.md Step 5 (Task 2 adds the reference)
- The prompt strategy uses bidirectional proof to reduce LLM bias — avoid directly asking "are these contradictory?"
- For token overflow risk with large proposals (SC > 25): rule should specify serial per-cluster checking rather than loading all pairs at once
- The fallback direction check covers only ADD vs SUBTRACT on same symbol — non-directional cross-group contradictions are left to the eval layer
