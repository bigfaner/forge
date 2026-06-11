---
id: "1"
title: "Implement Kahn's algorithm topological sort"
priority: "P0"
estimated_time: "2h"
dependencies: []
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.feature"
mainSession: false
---

# 1: Implement Kahn's algorithm topological sort

## Description

Implement the core topological sorting algorithm using Kahn's algorithm in `pkg/task/`. This is the foundation for both `forge task list` default ordering and `claimNextTask` alignment. The algorithm must handle wildcard dependency expansion (reusing existing `ResolveWildcardDep`), cycle detection, and missing dependency identification.

## Reference Files
- `proposal.md#Proposed-Solution` — defines Kahn's algorithm as the core approach, cycle detection, and wildcard expansion requirements
- `proposal.md#拓扑排序规则` — detailed sorting rules: same-level natural ID ordering, cycle detection with `[cycle]` marking, direct-only dependencies
- `proposal.md#通配符依赖规范` — wildcard expansion rules: `N.x` prefix matching, natural sort of expanded results, stability guarantee, unresolved warning
- `forge-cli/pkg/task/deps.go` — existing `ResolveWildcardDep` and `GetUnmetDeps` to reuse for wildcard expansion

## Acceptance Criteria

- [ ] `TopologicalSort(tasks TaskIndex) ([]string, []string, []string)` returns ordered IDs, cycle IDs, and missing dep IDs
- [ ] Kahn's algorithm produces correct ordering: if B depends on A, A appears before B
- [ ] Same-level tasks (no dependency relationship) sorted by natural ID for determinism
- [ ] Wildcard deps (`1.x`) expanded via `ResolveWildcardDep` before building adjacency list; expansion is stable across runs
- [ ] Cycle detection: remaining nodes after Kahn's completes are identified as cycle members
- [ ] Missing deps: dependencies referencing non-existent task IDs are identified
- [ ] Empty task list returns empty results (no panic)
- [ ] Single task with no deps returns `[that-task-id]`
- [ ] Self-referencing dependency detected as cycle
- [ ] Test coverage for: linear chain, diamond deps, disconnected components, wildcards, mixed exact+wildcard deps, cycle scenarios

## Hard Rules

- Reuse `ResolveWildcardDep` from `deps.go` for wildcard expansion — do not reimplement
- Algorithm must be O(V+E) — no quadratic shortcuts
- No side effects on input TaskIndex

## Implementation Notes

- `ResolveWildcardDep` already handles `N.x` prefix matching; call it for each dep that matches the wildcard pattern
- For mixed deps (e.g., `["1.1", "1.x"]`), expand wildcards then deduplicate against exact IDs before building the adjacency list
- Kahn's algorithm naturally prevents infinite loops — remaining in-degree > 0 nodes after processing are cycle members
