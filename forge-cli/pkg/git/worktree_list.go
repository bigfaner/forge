package git

import (
	"fmt"
	"path/filepath"
	"strings"
)

// WorktreeEntry represents a single git worktree parsed from porcelain output.
type WorktreeEntry struct {
	Path   string // Absolute path to the worktree directory
	HEAD   string // Current HEAD commit hash
	Branch string // Branch name (empty if detached HEAD)
	IsMain bool   // True for the first (main) worktree
}

// Name returns the worktree name derived from the directory basename.
func (e WorktreeEntry) Name() string {
	return filepath.Base(e.Path)
}

// ParsePorcelainWorktrees parses the output of `git worktree list --porcelain`.
// Each block is separated by a blank line. The first worktree is the main worktree.
// Bare repository markers ("bare") are skipped.
func ParsePorcelainWorktrees(output string) ([]WorktreeEntry, error) {
	output = strings.TrimSpace(output)
	if output == "" {
		return nil, nil
	}

	blocks := strings.Split(output, "\n\n")
	var entries []WorktreeEntry

	for i, block := range blocks {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}

		var entry WorktreeEntry
		lines := strings.Split(block, "\n")
		for _, line := range lines {
			key, value, found := strings.Cut(line, " ")
			if !found {
				// "bare" marker — skip this entry
				if key == "bare" {
					entry = WorktreeEntry{}
					goto nextBlock
				}
				continue
			}
			switch key {
			case "worktree":
				entry.Path = value
			case "HEAD":
				entry.HEAD = value
			case "branch":
				// Strip refs/heads/ prefix
				entry.Branch = strings.TrimPrefix(value, "refs/heads/")
			}
		}

		// Skip entries without a worktree path (e.g., bare repos)
		if entry.Path == "" {
			continue
		}

		if i == 0 {
			entry.IsMain = true
		}
		entries = append(entries, entry)

	nextBlock:
	}

	return entries, nil
}

// ListWorktrees returns all git worktrees for the repository at projectRoot.
// It uses `git worktree list --porcelain` for reliable machine-readable output.
func ListWorktrees(projectRoot string) ([]WorktreeEntry, error) {
	output, err := Run(projectRoot, "worktree", "list", "--porcelain")
	if err != nil {
		return nil, fmt.Errorf("list worktrees: %w", err)
	}

	return ParsePorcelainWorktrees(output)
}
