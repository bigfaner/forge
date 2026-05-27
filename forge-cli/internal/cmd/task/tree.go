package task

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"forge-cli/internal/cmd/base"

	"forge-cli/pkg/task"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// treeNode represents a task in the dependency tree hierarchy.
type treeNode struct {
	Task     task.Task
	Children []*treeNode
	Expanded bool
}

// flatItem is a flattened tree node used for rendering and navigation.
type flatItem struct {
	Node  *treeNode
	Depth int
}

// forestOption configures forest building.
type forestOption struct {
	sortByID bool
}

// forestOptionFunc is a functional option for buildForest.
type forestOptionFunc func(*forestOption)

// withSortByID sets whether siblings are sorted by natural ID.
func withSortByID(v bool) forestOptionFunc {
	return func(o *forestOption) {
		o.sortByID = v
	}
}

// buildForest constructs a tree from the task index using dependency edges.
// Root nodes are tasks with no dependencies (or whose deps are missing/cyclic).
func buildForest(idx *task.TaskIndex, opts ...forestOptionFunc) []*treeNode {
	o := forestOption{}
	for _, fn := range opts {
		fn(&o)
	}

	tasks := idx.TasksMap()

	// Expand all dependencies (handle wildcards, missing).
	// Build a child -> parent mapping.
	childToParents := make(map[string][]string)
	allIDs := make([]string, 0, len(tasks))
	for id := range tasks {
		allIDs = append(allIDs, id)
	}

	for _, id := range allIDs {
		t := tasks[id]
		deps := expandDepsForTree(idx, t.Dependencies, id)
		for _, dep := range deps {
			if _, found := idx.ByID(dep); found {
				childToParents[id] = append(childToParents[id], dep)
			}
		}
	}

	// Build nodes map.
	nodeMap := make(map[string]*treeNode, len(tasks))
	for id, t := range tasks {
		nodeMap[id] = &treeNode{Task: t}
	}

	// Attach children to parents.
	for childID, parents := range childToParents {
		child := nodeMap[childID]
		for _, parentID := range parents {
			parent := nodeMap[parentID]
			parent.Children = append(parent.Children, child)
		}
	}

	// Find roots: tasks that have no parents in the childToParents map.
	// A task is a root if no other task lists it as a child via its expanded deps.
	hasParent := make(map[string]bool)
	for childID := range childToParents {
		hasParent[childID] = true
	}

	// Detect cycle nodes: tasks in cycles have no reachable root.
	// Use TopologicalSort to find cycles.
	_, cycles, _ := task.TopologicalSort(idx)
	cycleSet := make(map[string]bool, len(cycles))
	for _, c := range cycles {
		cycleSet[c] = true
	}

	// For cycle nodes, detach their children to avoid infinite recursion.
	// Cycle nodes appear as roots without their cyclic children.
	for id := range cycleSet {
		nodeMap[id].Children = nil
	}

	// Remove cyclic edges from non-cycle nodes' children.
	for _, id := range allIDs {
		if !cycleSet[id] {
			filtered := make([]*treeNode, 0, len(nodeMap[id].Children))
			for _, child := range nodeMap[id].Children {
				if !cycleSet[child.Task.ID] {
					filtered = append(filtered, child)
				}
			}
			nodeMap[id].Children = filtered
		}
	}

	var roots []*treeNode
	for _, id := range allIDs {
		if !hasParent[id] || cycleSet[id] {
			roots = append(roots, nodeMap[id])
		}
	}

	// Sort roots.
	sortTreeNodes(roots, o.sortByID)

	// Sort all children recursively.
	sortChildrenRecursive(roots, o.sortByID)

	// Expand root nodes by default.
	for _, r := range roots {
		r.Expanded = len(r.Children) > 0
	}

	return roots
}

// expandDepsForTree expands a task's dependencies for tree building.
// Similar to pkg/task.expandDeps but simplified (no missingSet tracking).
func expandDepsForTree(idx *task.TaskIndex, deps []string, selfID string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, dep := range deps {
		matches, isWildcard := task.ResolveWildcardDep(idx, dep)
		if isWildcard {
			for _, matchID := range matches {
				if !seen[matchID] {
					seen[matchID] = true
					result = append(result, matchID)
				}
			}
		} else {
			if dep == selfID {
				// Self-reference: skip in tree (would create cycle)
				continue
			}
			if _, found := idx.ByID(dep); found {
				if !seen[dep] {
					seen[dep] = true
					result = append(result, dep)
				}
			}
		}
	}
	return result
}

// sortTreeNodes sorts a slice of tree nodes.
// By default, uses natural ID order (topo-based from buildForest).
// If sortByID is true, uses CompareVersionIDs for natural ordering.
func sortTreeNodes(nodes []*treeNode, sortByID bool) {
	sort.SliceStable(nodes, func(i, j int) bool {
		if sortByID {
			return task.CompareVersionIDs(nodes[i].Task.ID, nodes[j].Task.ID)
		}
		// Default: same as topological order (already sorted by buildForest)
		return task.CompareVersionIDs(nodes[i].Task.ID, nodes[j].Task.ID)
	})
}

// sortChildrenRecursive sorts children at all levels.
func sortChildrenRecursive(roots []*treeNode, sortByID bool) {
	for _, r := range roots {
		sortTreeNodes(r.Children, sortByID)
		sortChildrenRecursive(r.Children, sortByID)
	}
}

// statusSymbol returns the symbol for a task status.
func statusSymbol(status string) string {
	switch status {
	case "completed":
		return "✓"
	case "in_progress":
		return "~"
	case "blocked", "failed", "rejected":
		return "✗"
	default: // pending, skipped, suspended, unknown
		return "○"
	}
}

// statusColor returns the lipgloss color name for a task status.
func statusColor(status string) string {
	switch status {
	case "completed":
		return "green"
	case "in_progress":
		return "yellow"
	case "blocked", "failed", "rejected":
		return "red"
	default: // pending, skipped, suspended, unknown
		return "gray"
	}
}

// canUseTUI determines whether the terminal supports TUI mode.
func canUseTUI(isTerminal bool, termEnv string) bool {
	if !isTerminal {
		return false
	}
	if termEnv == "" || termEnv == "dumb" {
		return false
	}
	return true
}

// flattenTree returns a flat list of visible tree items for rendering.
func flattenTree(roots []*treeNode, allExpanded bool) []flatItem {
	var items []flatItem
	for _, root := range roots {
		flattenNode(root, 0, &items, allExpanded)
	}
	return items
}

// flattenNode recursively flattens a tree node and its visible children.
func flattenNode(node *treeNode, depth int, items *[]flatItem, allExpanded bool) {
	*items = append(*items, flatItem{Node: node, Depth: depth})
	if node.Expanded || allExpanded {
		for _, child := range node.Children {
			flattenNode(child, depth+1, items, allExpanded)
		}
	}
}

// renderTreePlain renders the tree as plain text (no ANSI codes).
// Used for non-TTY output and testing.
func renderTreePlain(roots []*treeNode) string {
	var sb strings.Builder
	items := flattenTree(roots, true)
	for _, item := range items {
		indent := strings.Repeat("  ", item.Depth)
		sym := statusSymbol(item.Node.Task.Status)
		title := item.Node.Task.Title
		id := item.Node.Task.ID
		fmt.Fprintf(&sb, "%s%s %s: %s\n", indent, sym, id, title)
	}
	return sb.String()
}

// --- bubbletea TUI model ---

// treeModel is the bubbletea Model for the task tree TUI.
type treeModel struct {
	roots    []*treeNode
	items    []flatItem
	cursor   int
	useColor bool
	quitting bool
	width    int
	height   int
}

// newTreeModel creates a new TUI tree model.
func newTreeModel(roots []*treeNode, useColor bool) treeModel {
	items := flattenTree(roots, false)
	if len(items) == 0 {
		return treeModel{
			roots:    roots,
			items:    items,
			useColor: useColor,
		}
	}
	return treeModel{
		roots:    roots,
		items:    items,
		cursor:   0,
		useColor: useColor,
		width:    80,
		height:   24,
	}
}

// Init implements bubbletea.Model.
func (m treeModel) Init() tea.Cmd {
	return nil
}

// Update implements bubbletea.Model.
func (m treeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, quitKeys):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, upKeys):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, downKeys):
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case key.Matches(msg, expandKeys):
			if m.cursor < len(m.items) {
				node := m.items[m.cursor].Node
				if len(node.Children) > 0 {
					node.Expanded = true
					m.items = flattenTree(m.roots, false)
				}
			}
		case key.Matches(msg, collapseKeys):
			if m.cursor < len(m.items) {
				node := m.items[m.cursor].Node
				if len(node.Children) > 0 {
					node.Expanded = false
					m.items = flattenTree(m.roots, false)
				}
			}
		}
	}
	return m, nil
}

// View implements bubbletea.Model.
func (m treeModel) View() string {
	if m.quitting || len(m.items) == 0 {
		return ""
	}

	var sb strings.Builder

	// Header
	fmt.Fprintf(&sb, "Task Dependency Tree (q: quit, ↑↓: navigate, ←→: collapse/expand)\n\n")

	// Calculate visible window
	visibleHeight := m.height - 4 // reserve for header + footer
	if visibleHeight < 5 {
		visibleHeight = 5
	}

	start := 0
	end := len(m.items)
	if end > visibleHeight {
		// Scroll window around cursor
		start = m.cursor - visibleHeight/2
		if start < 0 {
			start = 0
		}
		end = start + visibleHeight
		if end > len(m.items) {
			end = len(m.items)
			start = end - visibleHeight
			if start < 0 {
				start = 0
			}
		}
	}

	for i := start; i < end; i++ {
		item := m.items[i]
		indent := strings.Repeat("  ", item.Depth)
		sym := statusSymbol(item.Node.Task.Status)
		id := item.Node.Task.ID
		title := item.Node.Task.Title

		// Expand/collapse indicator
		expandChar := " "
		if len(item.Node.Children) > 0 {
			if item.Node.Expanded {
				expandChar = "▾"
			} else {
				expandChar = "▸"
			}
		}

		// Cursor indicator
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}

		line := fmt.Sprintf("%s %s%s %s %s", cursor, indent, expandChar, sym, formatNodeLine(id, title))

		if m.useColor {
			color := statusColor(item.Node.Task.Status)
			styled := lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(sym+" "+id) + ": " + title
			line = fmt.Sprintf("%s %s%s %s", cursor, indent, expandChar, styled)
		}

		sb.WriteString(line)
		sb.WriteString("\n")
	}

	// Footer
	fmt.Fprintf(&sb, "\n%d/%d tasks", m.cursor+1, len(m.items))

	return sb.String()
}

// formatNodeLine formats a tree node line as "ID: Title".
func formatNodeLine(id, title string) string {
	return id + ": " + title
}

// Key bindings.
var (
	quitKeys = key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	)
	upKeys = key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	)
	downKeys = key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	)
	expandKeys = key.NewBinding(
		key.WithKeys("right", "l", "enter"),
		key.WithHelp("→/l/enter", "expand"),
	)
	collapseKeys = key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "collapse"),
	)
)

// runTreeTUI launches the interactive TUI tree view.
func runTreeTUI(roots []*treeNode) error {
	_, termEnv := getTerminalInfo()
	useColor := termEnv != "dumb" && termEnv != ""

	p := tea.NewProgram(
		newTreeModel(roots, useColor),
		tea.WithAltScreen(),
	)
	_, err := p.Run()
	return err
}

// renderTreeFallback renders the tree as plain text for non-TTY environments.
// Falls back to table-like output with tree structure.
func renderTreeFallback(roots []*treeNode, featureSlug string) string {
	var sb strings.Builder
	items := flattenTree(roots, true)
	fmt.Fprintf(&sb, "%d found  (feature: %s)\n\n", len(items), featureSlug)

	for _, item := range items {
		indent := strings.Repeat("  ", item.Depth)
		sym := statusSymbol(item.Node.Task.Status)
		id := item.Node.Task.ID
		title := base.TruncateSlug(item.Node.Task.Title, 50)
		status := item.Node.Task.Status
		fmt.Fprintf(&sb, "%s%s %s  %-50s  %s\n", indent, sym, id, title, status)
	}

	return sb.String()
}

// getTerminalInfo returns whether stdout is a terminal and the TERM env value.
func getTerminalInfo() (isTerminal bool, termEnv string) {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false, ""
	}
	isTerminal = fi.Mode()&os.ModeCharDevice != 0
	termEnv = os.Getenv("TERM")
	return isTerminal, termEnv
}
