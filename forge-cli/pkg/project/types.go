// Package project provides project-level utilities.
package project

// RootType categorizes the detected project root.
type RootType int

const (
	// RootTypeUnknown indicates no root was found.
	RootTypeUnknown RootType = iota
	// RootTypeVCS indicates a version control boundary (.git, .hg).
	RootTypeVCS
	// RootTypeWorkspace indicates a monorepo/workspace root (go.work, pnpm-workspace.yaml).
	RootTypeWorkspace
	// RootTypeProject indicates a language-specific project root (go.mod, package.json).
	RootTypeProject
)

// String returns a human-readable representation of the root type.
func (t RootType) String() string {
	switch t {
	case RootTypeVCS:
		return "vcs"
	case RootTypeWorkspace:
		return "workspace"
	case RootTypeProject:
		return "project"
	default:
		return "unknown"
	}
}

// RootInfo contains full detection result.
type RootInfo struct {
	// Path is the absolute path to the detected root.
	Path string
	// Type indicates what kind of root was detected.
	Type RootType
	// Marker is the file/directory name that identified this root.
	Marker string
	// Languages contains detected project languages (e.g., "go", "java").
	Languages []string
}

// Marker represents a single root marker file or directory.
type Marker struct {
	// Name is the file or directory name (e.g., "go.mod", ".git").
	Name string
	// Type indicates the category of this marker.
	Type RootType
	// Languages lists applicable languages (e.g., ["go"], ["node", "javascript"]).
	Languages []string
	// IsDirectory indicates if this marker must be a directory.
	// If false, both files and directories are accepted (e.g., .git can be either).
	IsDirectory bool
	// IsFileGlob indicates if Name contains a glob pattern (e.g., "build.gradle*").
	IsFileGlob bool
}
