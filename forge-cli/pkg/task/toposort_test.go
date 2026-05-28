package task

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TopologicalSort should return ordered IDs, cycle IDs, and missing dep IDs.

func TestTopologicalSort_Empty(t *testing.T) {
	idx := NewTestIndex("test", map[string]Task{})
	ordered, cycles, missing := TopologicalSort(idx)
	assert.Empty(t, ordered)
	assert.Empty(t, cycles)
	assert.Empty(t, missing)
}

func TestTopologicalSort_SingleTaskNoDeps(t *testing.T) {
	idx := NewTestIndex("test", map[string]Task{
		"1": {ID: "1", Dependencies: nil},
	})
	ordered, cycles, missing := TopologicalSort(idx)
	assert.Equal(t, []string{"1"}, ordered)
	assert.Empty(t, cycles)
	assert.Empty(t, missing)
}

func TestTopologicalSort_LinearChain(t *testing.T) {
	// 3 → 2 → 1
	idx := NewTestIndex("test", map[string]Task{
		"1": {ID: "1", Dependencies: nil},
		"2": {ID: "2", Dependencies: []string{"1"}},
		"3": {ID: "3", Dependencies: []string{"2"}},
	})
	ordered, cycles, missing := TopologicalSort(idx)
	assert.Equal(t, []string{"1", "2", "3"}, ordered)
	assert.Empty(t, cycles)
	assert.Empty(t, missing)
}

func TestTopologicalSort_DiamondDeps(t *testing.T) {
	//     1
	//    / \
	//   2   3
	//    \ /
	//     4
	idx := NewTestIndex("test", map[string]Task{
		"1": {ID: "1", Dependencies: nil},
		"2": {ID: "2", Dependencies: []string{"1"}},
		"3": {ID: "3", Dependencies: []string{"1"}},
		"4": {ID: "4", Dependencies: []string{"2", "3"}},
	})
	ordered, cycles, missing := TopologicalSort(idx)
	assert.Empty(t, cycles)
	assert.Empty(t, missing)
	// 1 must come before 2 and 3; 2 and 3 must come before 4
	assertIndex := func(list []string, target string) int {
		for i, v := range list {
			if v == target {
				return i
			}
		}
		return -1
	}
	i1 := assertIndex(ordered, "1")
	i2 := assertIndex(ordered, "2")
	i3 := assertIndex(ordered, "3")
	i4 := assertIndex(ordered, "4")
	assert.True(t, i1 < i2, "1 before 2")
	assert.True(t, i1 < i3, "1 before 3")
	assert.True(t, i2 < i4, "2 before 4")
	assert.True(t, i3 < i4, "3 before 4")
}

func TestTopologicalSort_SameLevelNaturalIDOrder(t *testing.T) {
	// All tasks are independent → sorted by natural ID
	idx := NewTestIndex("test", map[string]Task{
		"3":   {ID: "3"},
		"1":   {ID: "1"},
		"2":   {ID: "2"},
		"10":  {ID: "10"},
		"1.1": {ID: "1.1"},
	})
	ordered, cycles, missing := TopologicalSort(idx)
	assert.Empty(t, cycles)
	assert.Empty(t, missing)
	// Natural sort: 1, 1.1, 2, 3, 10
	assert.Equal(t, []string{"1", "1.1", "2", "3", "10"}, ordered)
}

func TestTopologicalSort_Cycle(t *testing.T) {
	// 1 → 2 → 3 → 1 (cycle)
	idx := NewTestIndex("test", map[string]Task{
		"1": {ID: "1", Dependencies: []string{"3"}},
		"2": {ID: "2", Dependencies: []string{"1"}},
		"3": {ID: "3", Dependencies: []string{"2"}},
	})
	ordered, cycles, missing := TopologicalSort(idx)
	assert.Empty(t, ordered)
	assert.ElementsMatch(t, []string{"1", "2", "3"}, cycles)
	assert.Empty(t, missing)
}

func TestTopologicalSort_SelfReference(t *testing.T) {
	idx := NewTestIndex("test", map[string]Task{
		"1": {ID: "1", Dependencies: []string{"1"}},
	})
	ordered, cycles, missing := TopologicalSort(idx)
	assert.Empty(t, ordered)
	assert.ElementsMatch(t, []string{"1"}, cycles)
	assert.Empty(t, missing)
}

func TestTopologicalSort_MissingDep(t *testing.T) {
	idx := NewTestIndex("test", map[string]Task{
		"1": {ID: "1", Dependencies: nil},
		"2": {ID: "2", Dependencies: []string{"1", "999"}},
	})
	ordered, cycles, missing := TopologicalSort(idx)
	assert.Equal(t, []string{"1", "2"}, ordered)
	assert.Empty(t, cycles)
	assert.ElementsMatch(t, []string{"999"}, missing)
}

func TestTopologicalSort_DisconnectedComponents(t *testing.T) {
	// Component A: 1 → 2
	// Component B: 3 → 4
	// Kahn's processes nodes in natural ID order when in-degree is 0.
	// Valid topological order: 1 before 2, 3 before 4. Exact order among
	// independent nodes is determined by the min-heap.
	idx := NewTestIndex("test", map[string]Task{
		"1": {ID: "1", Dependencies: nil},
		"2": {ID: "2", Dependencies: []string{"1"}},
		"3": {ID: "3", Dependencies: nil},
		"4": {ID: "4", Dependencies: []string{"3"}},
	})
	ordered, cycles, missing := TopologicalSort(idx)
	assert.Empty(t, cycles)
	assert.Empty(t, missing)
	assertIndex := func(list []string, target string) int {
		for i, v := range list {
			if v == target {
				return i
			}
		}
		return -1
	}
	i1 := assertIndex(ordered, "1")
	i2 := assertIndex(ordered, "2")
	i3 := assertIndex(ordered, "3")
	i4 := assertIndex(ordered, "4")
	assert.True(t, i1 < i2, "1 before 2")
	assert.True(t, i3 < i4, "3 before 4")
}

func TestTopologicalSort_WildcardDeps(t *testing.T) {
	// Task "2" depends on "1.x" → expands to 1.1, 1.2
	idx := NewTestIndex("test", map[string]Task{
		"1.1": {ID: "1.1", Dependencies: nil},
		"1.2": {ID: "1.2", Dependencies: nil},
		"2":   {ID: "2", Dependencies: []string{"1.x"}},
	})
	ordered, cycles, missing := TopologicalSort(idx)
	assert.Empty(t, cycles)
	assert.Empty(t, missing)
	assertIndex := func(list []string, target string) int {
		for i, v := range list {
			if v == target {
				return i
			}
		}
		return -1
	}
	i11 := assertIndex(ordered, "1.1")
	i12 := assertIndex(ordered, "1.2")
	i2 := assertIndex(ordered, "2")
	assert.True(t, i11 < i2, "1.1 before 2")
	assert.True(t, i12 < i2, "1.2 before 2")
}

func TestTopologicalSort_MixedExactAndWildcardDeps(t *testing.T) {
	// Task "3" depends on ["1.1", "1.x"] → expands to 1.1, 1.2, dedupes 1.1
	idx := NewTestIndex("test", map[string]Task{
		"1.1": {ID: "1.1", Dependencies: nil},
		"1.2": {ID: "1.2", Dependencies: nil},
		"2":   {ID: "2", Dependencies: nil},
		"3":   {ID: "3", Dependencies: []string{"1.1", "1.x"}},
	})
	ordered, cycles, missing := TopologicalSort(idx)
	assert.Empty(t, cycles)
	assert.Empty(t, missing)
	assertIndex := func(list []string, target string) int {
		for i, v := range list {
			if v == target {
				return i
			}
		}
		return -1
	}
	i11 := assertIndex(ordered, "1.1")
	i12 := assertIndex(ordered, "1.2")
	i3 := assertIndex(ordered, "3")
	assert.True(t, i11 < i3, "1.1 before 3")
	assert.True(t, i12 < i3, "1.2 before 3")
}

func TestTopologicalSort_WildcardNoMatch(t *testing.T) {
	// "2" depends on "5.x" which matches nothing → no missing (wildcard unmatched is not missing dep)
	idx := NewTestIndex("test", map[string]Task{
		"2": {ID: "2", Dependencies: []string{"5.x"}},
	})
	ordered, cycles, missing := TopologicalSort(idx)
	// Wildcard with no matches: no edges created, 2 has no deps → it appears alone
	assert.Equal(t, []string{"2"}, ordered)
	assert.Empty(t, cycles)
	// Unresolved wildcard is not treated as missing dep per spec
	assert.Empty(t, missing)
}

func TestTopologicalSort_CycleWithNonCycleNodes(t *testing.T) {
	// 1 → 2 (happy), 3 → 4 → 3 (cycle)
	idx := NewTestIndex("test", map[string]Task{
		"1": {ID: "1", Dependencies: nil},
		"2": {ID: "2", Dependencies: []string{"1"}},
		"3": {ID: "3", Dependencies: []string{"4"}},
		"4": {ID: "4", Dependencies: []string{"3"}},
	})
	ordered, cycles, missing := TopologicalSort(idx)
	assert.Equal(t, []string{"1", "2"}, ordered)
	assert.ElementsMatch(t, []string{"3", "4"}, cycles)
	assert.Empty(t, missing)
}

func TestTopologicalSort_NoSideEffectsOnInput(t *testing.T) {
	idx := NewTestIndex("test", map[string]Task{
		"1": {ID: "1", Dependencies: nil},
		"2": {ID: "2", Dependencies: []string{"1"}},
	})
	// Capture original deps slice
	originalDeps := idx.TasksMap()["2"].Dependencies
	TopologicalSort(idx)
	// Verify deps not mutated
	assert.Equal(t, originalDeps, idx.TasksMap()["2"].Dependencies)
}

func TestTopologicalSort_MultipleMissingDeps(t *testing.T) {
	idx := NewTestIndex("test", map[string]Task{
		"1": {ID: "1", Dependencies: []string{"99", "100"}},
	})
	ordered, cycles, missing := TopologicalSort(idx)
	assert.Equal(t, []string{"1"}, ordered)
	assert.Empty(t, cycles)
	assert.ElementsMatch(t, []string{"99", "100"}, missing)
}

func TestTopologicalSort_WildcardStability(t *testing.T) {
	// Run multiple times to verify wildcard expansion is stable
	idx := NewTestIndex("test", map[string]Task{
		"1.1": {ID: "1.1"},
		"1.2": {ID: "1.2"},
		"1.3": {ID: "1.3"},
		"2":   {ID: "2", Dependencies: []string{"1.x"}},
	})
	var results [][]string
	for i := 0; i < 5; i++ {
		ordered, _, _ := TopologicalSort(idx)
		results = append(results, ordered)
	}
	for i := 1; i < len(results); i++ {
		assert.Equal(t, results[0], results[i], "topo sort must be stable across runs")
	}
}
