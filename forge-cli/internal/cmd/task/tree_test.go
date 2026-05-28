package task

import (
	"strings"
	"testing"

	"forge-cli/pkg/task"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

// --- TreeNode construction tests ---

func TestBuildForest_BasicHierarchy(t *testing.T) {
	idx := task.NewTestIndex("test", map[string]task.Task{
		"1": {ID: "1", Dependencies: nil},
		"2": {ID: "2", Dependencies: []string{"1"}},
		"3": {ID: "3", Dependencies: []string{"1"}},
		"4": {ID: "4", Dependencies: []string{"2"}},
	})

	roots := buildForest(idx)
	assert.Len(t, roots, 1, "single root (task 1)")
	assert.Equal(t, "1", roots[0].Task.ID)
	assert.Len(t, roots[0].Children, 2, "task 1 has 2 children")
	// Children should be in natural order: 2, 3
	assert.Equal(t, "2", roots[0].Children[0].Task.ID)
	assert.Equal(t, "3", roots[0].Children[1].Task.ID)
	// Task 4 depends on 2
	assert.Len(t, roots[0].Children[0].Children, 1)
	assert.Equal(t, "4", roots[0].Children[0].Children[0].Task.ID)
}

func TestBuildForest_MultipleRoots(t *testing.T) {
	idx := task.NewTestIndex("test", map[string]task.Task{
		"1": {ID: "1", Dependencies: nil},
		"2": {ID: "2", Dependencies: nil},
		"3": {ID: "3", Dependencies: []string{"1"}},
	})

	roots := buildForest(idx)
	assert.Len(t, roots, 2, "two roots (tasks 1 and 2)")
	// Natural order: 1, 2
	assert.Equal(t, "1", roots[0].Task.ID)
	assert.Equal(t, "2", roots[1].Task.ID)
	// Task 3 is child of 1
	assert.Len(t, roots[0].Children, 1)
	assert.Equal(t, "3", roots[0].Children[0].Task.ID)
}

func TestBuildForest_CycleNodesAtEnd(t *testing.T) {
	// 1 -> 2 -> 3 -> 1 (cycle)
	idx := task.NewTestIndex("test", map[string]task.Task{
		"1": {ID: "1", Dependencies: []string{"3"}},
		"2": {ID: "2", Dependencies: []string{"1"}},
		"3": {ID: "3", Dependencies: []string{"2"}},
	})

	roots := buildForest(idx)
	// All are in cycle, so all appear as roots (no parent found)
	assert.Len(t, roots, 3, "all cycle nodes appear as roots")
}

func TestBuildForest_MissingDepsIgnored(t *testing.T) {
	idx := task.NewTestIndex("test", map[string]task.Task{
		"1": {ID: "1", Dependencies: nil},
		"2": {ID: "2", Dependencies: []string{"1", "999"}},
	})

	roots := buildForest(idx)
	assert.Len(t, roots, 1, "single root (task 1)")
	assert.Equal(t, "1", roots[0].Task.ID)
	assert.Len(t, roots[0].Children, 1)
	assert.Equal(t, "2", roots[0].Children[0].Task.ID)
}

func TestBuildForest_SortIDSiblings(t *testing.T) {
	// With sortByID=true, siblings should be sorted by natural ID
	idx := task.NewTestIndex("test", map[string]task.Task{
		"1":  {ID: "1", Dependencies: nil},
		"3":  {ID: "3", Dependencies: []string{"1"}},
		"2":  {ID: "2", Dependencies: []string{"1"}},
		"10": {ID: "10", Dependencies: []string{"1"}},
	})

	roots := buildForest(idx, withSortByID(true))
	assert.Len(t, roots, 1)
	assert.Equal(t, "1", roots[0].Task.ID)
	assert.Len(t, roots[0].Children, 3)
	// Natural sort: 2, 3, 10
	assert.Equal(t, "2", roots[0].Children[0].Task.ID)
	assert.Equal(t, "3", roots[0].Children[1].Task.ID)
	assert.Equal(t, "10", roots[0].Children[2].Task.ID)
}

// --- Status encoding tests ---

func TestStatusSymbol(t *testing.T) {
	tests := []struct {
		status string
		want   string
	}{
		{"completed", "✓"},
		{"in_progress", "~"},
		{"blocked", "✗"},
		{"failed", "✗"},
		{"pending", "○"},
		{"skipped", "○"},
		{"suspended", "○"},
		{"rejected", "✗"},
		{"unknown", "○"},
	}
	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			got := statusSymbol(tt.status)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStatusColor(t *testing.T) {
	tests := []struct {
		status string
		color  string
	}{
		{"completed", "green"},
		{"in_progress", "yellow"},
		{"blocked", "red"},
		{"failed", "red"},
		{"pending", "gray"},
		{"skipped", "gray"},
		{"suspended", "gray"},
		{"rejected", "red"},
		{"unknown", "gray"},
	}
	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			got := statusColor(tt.status)
			assert.Equal(t, tt.color, got)
		})
	}
}

// --- Terminal detection tests ---

func TestCanUseTUI(t *testing.T) {
	t.Run("returns false when isTerminal is false", func(t *testing.T) {
		assert.False(t, canUseTUI(false, ""))
	})
	t.Run("returns true when isTerminal and TERM is set", func(t *testing.T) {
		assert.True(t, canUseTUI(true, "xterm-256color"))
	})
	t.Run("returns false when TERM is dumb", func(t *testing.T) {
		assert.False(t, canUseTUI(true, "dumb"))
	})
	t.Run("returns false when TERM is empty", func(t *testing.T) {
		assert.False(t, canUseTUI(true, ""))
	})
}

// --- Tree rendering tests (plain text, for verification) ---

func TestRenderTreePlain_BasicTree(t *testing.T) {
	root := &treeNode{
		Task: task.Task{ID: "1", Title: "Root task", Status: "completed"},
		Children: []*treeNode{
			{
				Task: task.Task{ID: "2", Title: "Child task", Status: "pending"},
			},
		},
	}

	output := renderTreePlain([]*treeNode{root})
	lines := strings.Split(output, "\n")

	// Should contain root and child
	assert.True(t, strings.Contains(lines[0], "1"), "first line should contain ID 1")
	assert.True(t, strings.Contains(output, "2"), "output should contain ID 2")
	assert.True(t, strings.Contains(output, "Root task"), "output should contain title")
	assert.True(t, strings.Contains(output, "Child task"), "output should contain child title")
}

func TestRenderTreePlain_StatusSymbols(t *testing.T) {
	root := &treeNode{
		Task: task.Task{ID: "1", Title: "Done", Status: "completed"},
		Children: []*treeNode{
			{Task: task.Task{ID: "2", Title: "Working", Status: "in_progress"}},
			{Task: task.Task{ID: "3", Title: "Blocked", Status: "blocked"}},
			{Task: task.Task{ID: "4", Title: "Waiting", Status: "pending"}},
		},
	}

	output := renderTreePlain([]*treeNode{root})
	assert.Contains(t, output, "✓", "completed should show ✓")
	assert.Contains(t, output, "~", "in_progress should show ~")
	assert.Contains(t, output, "✗", "blocked should show ✗")
	assert.Contains(t, output, "○", "pending should show ○")
}

func TestRenderTreePlain_Indentation(t *testing.T) {
	root := &treeNode{
		Task: task.Task{ID: "1", Title: "Root", Status: "completed"},
		Children: []*treeNode{
			{
				Task: task.Task{ID: "2", Title: "Child", Status: "pending"},
				Children: []*treeNode{
					{Task: task.Task{ID: "3", Title: "Grandchild", Status: "pending"}},
				},
			},
		},
	}

	output := renderTreePlain([]*treeNode{root})
	lines := strings.Split(output, "\n")

	// Count indentation levels
	rootLine := lines[0]
	childLine := lines[1]
	grandchildLine := lines[2]

	// renderTreePlain uses "  " per depth level
	assert.True(t, strings.HasPrefix(rootLine, "✓ 1:"), "root has no indent")
	assert.True(t, strings.HasPrefix(childLine, "  "), "child is indented")
	assert.False(t, strings.HasPrefix(childLine, "    "), "child has exactly 1 level of indent")
	assert.True(t, strings.HasPrefix(grandchildLine, "    "), "grandchild is more indented")
}

func TestBuildForest_MapKeyDiffersFromTaskID(t *testing.T) {
	// Regression: when the map key in index.json differs from Task.ID,
	// buildForest used map keys for nodeMap but Task.ID values for dependency
	// lookups, causing nil pointer dereference on parent.Children.
	idx := task.NewTestIndex("test", map[string]task.Task{
		"1-setup-env":   {ID: "1", Dependencies: nil},
		"2-write-tests": {ID: "2", Dependencies: []string{"1"}},
		"3-run-tests":   {ID: "3", Dependencies: []string{"2"}},
	})

	roots := buildForest(idx)
	assert.Len(t, roots, 1, "single root (task 1)")
	assert.Equal(t, "1", roots[0].Task.ID)
	assert.Len(t, roots[0].Children, 1, "task 1 has child 2")
	assert.Equal(t, "2", roots[0].Children[0].Task.ID)
	assert.Len(t, roots[0].Children[0].Children, 1, "task 2 has child 3")
	assert.Equal(t, "3", roots[0].Children[0].Children[0].Task.ID)
}

func TestBuildForest_MapKeyDiffersFromTaskID_WildcardDeps(t *testing.T) {
	// Wildcard deps resolve to Task.ID values which may differ from map keys.
	idx := task.NewTestIndex("test", map[string]task.Task{
		"1-impl":   {ID: "1.1", Dependencies: nil},
		"1-gate":   {ID: "1.2", Dependencies: nil},
		"2-review": {ID: "2", Dependencies: []string{"1.x"}},
	})

	roots := buildForest(idx)
	// 1.1 and 1.2 are roots, 2 is child of both
	assert.Len(t, roots, 2)
	found2 := false
	for _, r := range roots {
		for _, c := range r.Children {
			if c.Task.ID == "2" {
				found2 = true
			}
		}
	}
	assert.True(t, found2, "task 2 should be a child of one of the roots")
}

func TestBuildForest_WildcardDeps(t *testing.T) {
	idx := task.NewTestIndex("test", map[string]task.Task{
		"1.1": {ID: "1.1", Dependencies: nil},
		"1.2": {ID: "1.2", Dependencies: nil},
		"2":   {ID: "2", Dependencies: []string{"1.x"}},
	})

	roots := buildForest(idx)
	// 1.1 and 1.2 are roots, 2 is child of both
	assert.Len(t, roots, 2)
	// Both 1.1 and 1.2 should have 2 as child
	found2 := false
	for _, r := range roots {
		for _, c := range r.Children {
			if c.Task.ID == "2" {
				found2 = true
			}
		}
	}
	assert.True(t, found2, "task 2 should be a child of one of the roots")
}

func TestRenderTreeFallback(t *testing.T) {
	root := &treeNode{
		Task: task.Task{ID: "1", Title: "Root task", Status: "completed"},
		Children: []*treeNode{
			{Task: task.Task{ID: "2", Title: "Child task", Status: "pending"}},
		},
	}

	output := renderTreeFallback([]*treeNode{root}, "my-feature")
	assert.Contains(t, output, "2 found", "should show count")
	assert.Contains(t, output, "feature: my-feature", "should show feature name")
	assert.Contains(t, output, "✓", "should show completed symbol")
	assert.Contains(t, output, "○", "should show pending symbol")
	assert.Contains(t, output, "Root task", "should show task title")
}

// --- --tree flag integration test ---

func TestListCmd_TreeFlag(t *testing.T) {
	t.Run("--tree flag is registered", func(t *testing.T) {
		f := listCmd.Flags().Lookup("tree")
		assert.NotNil(t, f, "--tree flag should be registered")
		assert.Equal(t, "false", f.DefValue, "--tree default should be false")
	})
}

func TestListCmd_TreeFallbackToTable(t *testing.T) {
	t.Run("--tree in non-TTY falls back to table mode", func(t *testing.T) {
		tasks := map[string]task.Task{
			"1": {ID: "1", Title: "Root task", Type: "coding.feature", Status: "completed"},
			"2": {ID: "2", Title: "Child task", Type: "coding.feature", Status: "pending", Dependencies: []string{"1"}},
		}
		_ = setupFullProject(t, SetupOpts{Tasks: tasks})

		// Force non-TTY
		orig := listIsTerminalFunc
		listIsTerminalFunc = func() bool { return false }
		defer func() { listIsTerminalFunc = orig }()

		cmd := helperListCmd("topo")
		cmd.Flags().Bool("tree", true, "")
		_ = cmd.Flags().Set("tree", "true")

		output := captureStdout(func() {
			err := runList(cmd, []string{})
			if err != nil {
				t.Fatalf("runList returned error: %v", err)
			}
		})

		// Should fall back to table mode (contain table headers)
		assert.Contains(t, output, "ID", "should contain table header ID")
		assert.Contains(t, output, "STATUS", "should contain table header STATUS")
	})
}

func TestTreeModel_View(t *testing.T) {
	t.Run("tree model view contains task info", func(t *testing.T) {
		root := &treeNode{
			Task: task.Task{ID: "1", Title: "Root task", Status: "completed"},
			Children: []*treeNode{
				{Task: task.Task{ID: "2", Title: "Child task", Status: "pending"}},
			},
		}
		root.Expanded = true

		model := newTreeModel([]*treeNode{root}, false)
		view := model.View()

		assert.Contains(t, view, "1", "view should contain task ID 1")
		assert.Contains(t, view, "Root task", "view should contain root title")
		assert.Contains(t, view, "Task Dependency Tree", "view should contain header")
	})
}

func TestTreeModel_Navigation(t *testing.T) {
	t.Run("up/down navigation updates cursor", func(t *testing.T) {
		root := &treeNode{
			Task: task.Task{ID: "1", Title: "Root", Status: "completed"},
			Children: []*treeNode{
				{Task: task.Task{ID: "2", Title: "Child A", Status: "pending"}},
				{Task: task.Task{ID: "3", Title: "Child B", Status: "pending"}},
			},
		}
		root.Expanded = true

		model := newTreeModel([]*treeNode{root}, false)
		assert.Equal(t, 0, model.cursor)

		// Move down
		updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
		m := updated.(treeModel)
		assert.Equal(t, 1, m.cursor)

		// Move down again
		updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = updated.(treeModel)
		assert.Equal(t, 2, m.cursor)

		// Move up
		updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
		m = updated.(treeModel)
		assert.Equal(t, 1, m.cursor)
	})

	t.Run("q quits", func(t *testing.T) {
		root := &treeNode{
			Task: task.Task{ID: "1", Title: "Root", Status: "completed"},
		}
		model := newTreeModel([]*treeNode{root}, false)
		updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
		m := updated.(treeModel)
		assert.True(t, m.quitting)
		assert.NotNil(t, cmd)
	})

	t.Run("Ctrl+C quits", func(t *testing.T) {
		root := &treeNode{
			Task: task.Task{ID: "1", Title: "Root", Status: "completed"},
		}
		model := newTreeModel([]*treeNode{root}, false)
		updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
		m := updated.(treeModel)
		assert.True(t, m.quitting)
		assert.NotNil(t, cmd)
	})

	t.Run("cursor cannot go below zero", func(t *testing.T) {
		root := &treeNode{Task: task.Task{ID: "1", Title: "Root", Status: "completed"}}
		model := newTreeModel([]*treeNode{root}, false)
		updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyUp})
		m := updated.(treeModel)
		assert.Equal(t, 0, m.cursor, "cursor should stay at 0")
	})

	t.Run("cursor cannot exceed items length", func(t *testing.T) {
		root := &treeNode{Task: task.Task{ID: "1", Title: "Root", Status: "completed"}}
		model := newTreeModel([]*treeNode{root}, false)
		updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
		m := updated.(treeModel)
		assert.Equal(t, 0, m.cursor, "cursor should stay at 0 with single item")
	})

	t.Run("expand right arrow", func(t *testing.T) {
		root := &treeNode{
			Task: task.Task{ID: "1", Title: "Root", Status: "completed"},
			Children: []*treeNode{
				{Task: task.Task{ID: "2", Title: "Child", Status: "pending"}},
			},
		}
		model := newTreeModel([]*treeNode{root}, false)
		assert.False(t, root.Expanded, "root should start collapsed")

		updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRight})
		m := updated.(treeModel)
		assert.True(t, m.roots[0].Expanded, "right arrow should expand root")
		assert.Len(t, m.items, 2, "expanded root should show 2 items")
	})

	t.Run("collapse left arrow", func(t *testing.T) {
		root := &treeNode{
			Task:     task.Task{ID: "1", Title: "Root", Status: "completed"},
			Expanded: true,
			Children: []*treeNode{
				{Task: task.Task{ID: "2", Title: "Child", Status: "pending"}},
			},
		}
		model := newTreeModel([]*treeNode{root}, false)
		assert.True(t, root.Expanded)
		assert.Len(t, model.items, 2)

		updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyLeft})
		m := updated.(treeModel)
		assert.False(t, m.roots[0].Expanded, "left arrow should collapse root")
		assert.Len(t, m.items, 1, "collapsed root should show 1 item")
	})

	t.Run("window resize updates dimensions", func(t *testing.T) {
		root := &treeNode{Task: task.Task{ID: "1", Title: "Root", Status: "completed"}}
		model := newTreeModel([]*treeNode{root}, false)

		updated, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
		m := updated.(treeModel)
		assert.Equal(t, 120, m.width)
		assert.Equal(t, 40, m.height)
	})

	t.Run("View with color", func(t *testing.T) {
		root := &treeNode{
			Task: task.Task{ID: "1", Title: "Root task", Status: "completed"},
		}
		model := newTreeModel([]*treeNode{root}, true)
		view := model.View()
		assert.Contains(t, view, "1", "colored view should contain task ID")
		assert.Contains(t, view, "Root task", "colored view should contain title")
	})

	t.Run("View with empty items", func(t *testing.T) {
		model := newTreeModel(nil, false)
		view := model.View()
		assert.Equal(t, "", view, "empty model should return empty view")
	})

	t.Run("View quitting returns empty", func(t *testing.T) {
		root := &treeNode{Task: task.Task{ID: "1", Title: "Root", Status: "completed"}}
		model := newTreeModel([]*treeNode{root}, false)
		model.quitting = true
		view := model.View()
		assert.Equal(t, "", view)
	})
}

// --- Flat tree render test for TUI content ---

func TestFlattenTree(t *testing.T) {
	root := &treeNode{
		Task: task.Task{ID: "1", Title: "Root", Status: "completed"},
		Children: []*treeNode{
			{Task: task.Task{ID: "2", Title: "Child A", Status: "pending"}},
			{Task: task.Task{ID: "3", Title: "Child B", Status: "in_progress"}},
		},
	}

	items := flattenTree([]*treeNode{root}, true)
	assert.Len(t, items, 3)
	assert.Equal(t, "1", items[0].Node.Task.ID)
	assert.Equal(t, 0, items[0].Depth)
	assert.Equal(t, "2", items[1].Node.Task.ID)
	assert.Equal(t, 1, items[1].Depth)
	assert.Equal(t, "3", items[2].Node.Task.ID)
	assert.Equal(t, 1, items[2].Depth)
}

func TestFlattenTree_Collapsed(t *testing.T) {
	root := &treeNode{
		Task:     task.Task{ID: "1", Title: "Root", Status: "completed"},
		Expanded: true,
		Children: []*treeNode{
			{
				Task:     task.Task{ID: "2", Title: "Child", Status: "pending"},
				Expanded: true,
				Children: []*treeNode{
					{Task: task.Task{ID: "3", Title: "Grandchild", Status: "pending"}},
				},
			},
		},
	}

	// Fully expanded (allExpanded=true overrides Expanded flag)
	items := flattenTree([]*treeNode{root}, true)
	assert.Len(t, items, 3)

	// Collapse child (use allExpanded=false to respect Expanded flag)
	root.Children[0].Expanded = false
	items = flattenTree([]*treeNode{root}, false)
	assert.Len(t, items, 2, "collapsed child should hide grandchild")

	// Collapse root
	root.Expanded = false
	items = flattenTree([]*treeNode{root}, false)
	assert.Len(t, items, 1, "collapsed root should hide all children")
}
