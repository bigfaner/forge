package task

import "sort"

// TopologicalSort performs a topological sort on the tasks in the given TaskIndex
// using Kahn's algorithm. It returns three slices:
//   - ordered: task IDs in topological order (dependencies before dependents)
//   - cycles: task IDs that participate in dependency cycles
//   - missing: dependency IDs that reference non-existent tasks
//
// Same-level tasks (no dependency relationship) are sorted by natural ID for determinism.
// Wildcard dependencies (e.g. "1.x") are expanded via ResolveWildcardDep before building
// the adjacency list. The input TaskIndex is not modified.
func TopologicalSort(idx *TaskIndex) (ordered []string, cycles []string, missing []string) {
	tasks := idx.TasksMap()
	if len(tasks) == 0 {
		return nil, nil, nil
	}

	// Use Task.ID as canonical key (map key may differ from Task.ID).
	allIDs := make([]string, 0, len(tasks))
	for _, t := range tasks {
		allIDs = append(allIDs, t.ID)
	}

	// Build adjacency list and compute in-degrees.
	// adjacency[a] = [b, c] means a must come before b and c.
	adjacency := make(map[string][]string)
	inDegree := make(map[string]int)
	for _, id := range allIDs {
		inDegree[id] = 0
	}

	// Track all unique missing deps across tasks.
	missingSet := make(map[string]bool)

	for _, id := range allIDs {
		t, _ := idx.ByID(id)
		expandedDeps := expandDeps(idx, t.Dependencies, id, missingSet)
		for _, dep := range expandedDeps {
			adjacency[dep] = append(adjacency[dep], id)
			inDegree[id]++
		}
	}

	if len(missingSet) > 0 {
		missing = setToSortedSlice(missingSet)
	}

	// Kahn's algorithm with a min-heap by natural ID for deterministic ordering.
	queue := newIDQueue()
	for _, id := range allIDs {
		if inDegree[id] == 0 {
			queue.push(id)
		}
	}

	for queue.len() > 0 {
		node := queue.pop()
		ordered = append(ordered, node)
		for _, neighbor := range adjacency[node] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue.push(neighbor)
			}
		}
	}

	// Remaining nodes with in-degree > 0 are cycle members.
	for _, id := range allIDs {
		if inDegree[id] > 0 {
			cycles = append(cycles, id)
		}
	}
	if len(cycles) > 0 {
		sort.Slice(cycles, func(i, j int) bool {
			return CompareVersionIDs(cycles[i], cycles[j])
		})
	}

	return ordered, cycles, missing
}

// expandDeps expands a task's dependency list, resolving wildcards and
// identifying missing deps. Returns deduplicated concrete task IDs.
func expandDeps(idx *TaskIndex, deps []string, selfID string, missingSet map[string]bool) []string {
	seen := make(map[string]bool)
	var result []string

	for _, dep := range deps {
		matches, isWildcard := ResolveWildcardDep(idx, dep)
		if isWildcard {
			for _, matchID := range matches {
				// Self-reference via wildcard: include as edge so Kahn detects cycle.
				if !seen[matchID] {
					seen[matchID] = true
					result = append(result, matchID)
				}
			}
		} else {
			// Exact dep: check existence.
			if _, found := idx.ByID(dep); !found {
				// Self-referencing a non-existent ID is impossible, but handle gracefully.
				if dep != selfID {
					missingSet[dep] = true
				}
				continue
			}
			// Self-reference: include as edge so Kahn detects cycle (in-degree never reaches 0).
			if !seen[dep] {
				seen[dep] = true
				result = append(result, dep)
			}
		}
	}
	return result
}

// setToSortedSlice converts a string set to a naturally-sorted slice.
func setToSortedSlice(s map[string]bool) []string {
	result := make([]string, 0, len(s))
	for k := range s {
		result = append(result, k)
	}
	sort.Slice(result, func(i, j int) bool {
		return CompareVersionIDs(result[i], result[j])
	})
	return result
}

// idQueue is a priority queue that pops the smallest ID by natural sort order.
type idQueue struct {
	items []string
}

func newIDQueue() *idQueue {
	return &idQueue{}
}

func (q *idQueue) push(id string) {
	q.items = append(q.items, id)
	// Bubble up to maintain heap property.
	i := len(q.items) - 1
	for i > 0 {
		parent := (i - 1) / 2
		if CompareVersionIDs(q.items[i], q.items[parent]) {
			q.items[i], q.items[parent] = q.items[parent], q.items[i]
			i = parent
		} else {
			break
		}
	}
}

func (q *idQueue) pop() string {
	if len(q.items) == 0 {
		return ""
	}
	top := q.items[0]
	last := len(q.items) - 1
	q.items[0] = q.items[last]
	q.items = q.items[:last]
	// Sink down.
	i := 0
	for {
		left := 2*i + 1
		right := 2*i + 2
		smallest := i
		if left < len(q.items) && CompareVersionIDs(q.items[left], q.items[smallest]) {
			smallest = left
		}
		if right < len(q.items) && CompareVersionIDs(q.items[right], q.items[smallest]) {
			smallest = right
		}
		if smallest == i {
			break
		}
		q.items[i], q.items[smallest] = q.items[smallest], q.items[i]
		i = smallest
	}
	return top
}

func (q *idQueue) len() int {
	return len(q.items)
}
