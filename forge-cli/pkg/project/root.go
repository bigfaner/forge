// Package project provides project-level utilities.
package project

import (
	"fmt"
	"os"
	"path/filepath"
)

// FindProjectRoot locates the project root by searching for project markers.
// It maintains backward compatibility with the original Go-specific implementation
// while now supporting multiple languages and monorepo structures.
func FindProjectRoot() (string, error) {
	info, err := FindRootInfo()
	if err != nil {
		return "", err
	}
	return info.Path, nil
}

// FindProjectRootFrom locates the project root starting from a specific directory.
func FindProjectRootFrom(startDir string) (string, error) {
	info, err := FindRootInfoFrom(startDir)
	if err != nil {
		return "", err
	}
	return info.Path, nil
}

// FindRootInfo returns detailed information about the detected project root.
func FindRootInfo() (*RootInfo, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("get working directory: %w", err)
	}
	return FindRootInfoFrom(cwd)
}

// FindRootInfoFrom returns detailed information about the detected project root,
// starting the search from the specified directory.
func FindRootInfoFrom(startDir string) (*RootInfo, error) {
	// Step 1: Check environment variable override (highest priority)
	if envRoot := GetProjectRootFromEnv(); envRoot != "" {
		return &RootInfo{
			Path:   envRoot,
			Type:   RootTypeUnknown,
			Marker: "ENV",
		}, nil
	}

	// Step 2: Resolve symlinks and get absolute path
	absPath, err := filepath.Abs(startDir)
	if err != nil {
		return nil, fmt.Errorf("resolve path: %w", err)
	}

	// Step 3: Walk up the directory tree, collecting markers
	var foundWorkspace *RootInfo
	var foundProject *RootInfo
	var foundVCS *RootInfo

	dir := absPath
	for {
		// Check all markers in this directory
		for _, marker := range allMarkers() {
			if matchesMarker(dir, marker) {
				info := &RootInfo{
					Path:      dir,
					Type:      marker.Type,
					Marker:    marker.Name,
					Languages: marker.Languages,
				}

				switch marker.Type {
				case RootTypeWorkspace:
					// First workspace marker wins (closest to cwd)
					if foundWorkspace == nil {
						foundWorkspace = info
					}
				case RootTypeProject:
					// First project marker wins (closest to cwd)
					if foundProject == nil {
						foundProject = info
					}
				case RootTypeVCS:
					// Keep the VCS marker (we want the furthest, but we'll collect all)
					// Since we're walking up, the last VCS we find is the boundary
					foundVCS = info
				}
			}
		}

		// Move to parent directory
		parent := filepath.Dir(dir)
		if parent == dir {
			break // Reached filesystem root
		}
		dir = parent
	}

	// Step 4: Return based on priority
	// Workspace > Project > VCS
	if foundWorkspace != nil {
		return foundWorkspace, nil
	}
	if foundProject != nil {
		return foundProject, nil
	}
	if foundVCS != nil {
		return foundVCS, nil
	}

	return nil, fmt.Errorf("could not find project root (no markers found)")
}

// FindVCSRoot returns the VCS boundary, useful for monorepo-wide operations.
func FindVCSRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}
	return FindVCSRootFrom(cwd)
}

// FindVCSRootFrom returns the VCS boundary starting from a specific directory.
func FindVCSRootFrom(startDir string) (string, error) {
	absPath, err := filepath.Abs(startDir)
	if err != nil {
		return "", fmt.Errorf("resolve path: %w", err)
	}

	dir := absPath
	for {
		for _, marker := range vcsMarkers {
			if matchesMarker(dir, marker) {
				return dir, nil
			}
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("could not find VCS root (no .git or .hg found)")
}

// GetProjectRootFromEnv returns the project root from environment variables.
// It checks CLAUDE_PROJECT_DIR first (existing convention), then PROJECT_ROOT.
func GetProjectRootFromEnv() string {
	if dir := os.Getenv("CLAUDE_PROJECT_DIR"); dir != "" {
		return filepath.Clean(dir)
	}
	if dir := os.Getenv("PROJECT_ROOT"); dir != "" {
		return filepath.Clean(dir)
	}
	return ""
}

// matchesMarker checks if a marker exists in the given directory.
func matchesMarker(dir string, marker Marker) bool {
	path := filepath.Join(dir, marker.Name)

	// Handle glob patterns (e.g., build.gradle*)
	if marker.IsFileGlob {
		matches, _ := filepath.Glob(path)
		return len(matches) > 0
	}

	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	// For markers that can be either file or directory (e.g., .git in worktrees)
	if !marker.IsDirectory {
		return true // Both files and directories are acceptable
	}

	// For markers that must be directories
	return info.IsDir()
}
