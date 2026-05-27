---
id: "3"
title: "Align claimNextTask with topological ordering"
priority: "P0"
estimated_time: "1h"
dependencies: ["1"]
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 3: Align claimNextTask with topological ordering

## Description

Change `claimNextTask` internal sorting from priority+version to topological ordering, ensuring `forge task claim` picks tasks in the same order displayed by `forge task list`. This eliminates the ordering discrepancy between the two commands.

## Reference Files
- `proposal.md#Proposed-Solution` — defines claimNextTask alignment: "内部排序由优先级+版本号改为拓扑排序，与列表显示顺序一致"
- `proposal.md#Evidence` — existing claim logic uses priority+version sort, disconnected from list display order
- `forge-cli/internal/cmd/task/claim.go` — current `claimNextTask` at lines 168-226, target for modification

## Acceptance Criteria

- [ ] `claimNextTask` selects eligible tasks in topological order (first eligible task in topo sequence is claimed)
- [ ] Priority still acts as a tiebreaker among tasks at the same topological level (P0 before P1 before P2 within same depth)
- [ ] Existing `claim_test.go` tests pass after modification
- [ ] New test: verify topo-ordered claiming for a multi-dependency task graph
- [ ] New test: verify priority tiebreaker within same topological level

## Hard Rules

- Do not remove priority as a secondary sort — it still matters for tasks at the same topological depth
- Must not break the lazy unblock scan or fix-task blocking logic in `claimNextTask`

## Implementation Notes

- Current sorting at claim.go:212-219 uses `sort.SliceStable` with priority then `compareVersionIDs` as tiebreaker
- Replace with: (1) compute topo order from eligible tasks, (2) priority as tiebreaker within same topo level
- The eligible task set is already dependency-filtered (only tasks with all deps met), so the topo sort among them is a subset of the full graph topo sort
