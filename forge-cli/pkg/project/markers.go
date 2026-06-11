package project

// VCS markers identify version control boundaries.
// These are checked first and serve as the ultimate fallback.
var vcsMarkers = []Marker{
	{Name: ".git", Type: RootTypeVCS, IsDirectory: false}, // Can be file (worktree) or directory
	{Name: ".hg", Type: RootTypeVCS, IsDirectory: true},
}

// Workspace markers identify monorepo/workspace roots.
// These take priority over project markers when present.
var workspaceMarkers = []Marker{
	// Forge project root — created by task claim on first use.
	// Subagents walking up from subdirectories (e.g., backend/) find .forge/
	// and resolve to the correct project root (workspace > project priority).
	{Name: ".forge", Type: RootTypeWorkspace, IsDirectory: true},
	// Go multi-module workspace
	{Name: "go.work", Type: RootTypeWorkspace, Languages: []string{"go"}},
	// Node.js monorepo tools
	{Name: "pnpm-workspace.yaml", Type: RootTypeWorkspace, Languages: []string{"node"}},
	{Name: "lerna.json", Type: RootTypeWorkspace, Languages: []string{"node"}},
	{Name: "turbo.json", Type: RootTypeWorkspace, Languages: []string{"node"}},
	{Name: "nx.json", Type: RootTypeWorkspace, Languages: []string{"node"}},
	// Bazel workspace
	{Name: "WORKSPACE", Type: RootTypeWorkspace, Languages: []string{"bazel"}},
	{Name: "WORKSPACE.bazel", Type: RootTypeWorkspace, Languages: []string{"bazel"}},
	// Gradle multi-project (settings.gradle defines root)
	{Name: "settings.gradle", Type: RootTypeWorkspace, Languages: []string{"java", "groovy"}},
	{Name: "settings.gradle.kts", Type: RootTypeWorkspace, Languages: []string{"java", "kotlin"}},
}

// Project markers identify language-specific project roots.
var projectMarkers = []Marker{
	// Go
	{Name: "go.mod", Type: RootTypeProject, Languages: []string{"go"}},
	// Node.js / JavaScript / TypeScript
	{Name: "package.json", Type: RootTypeProject, Languages: []string{"node", "javascript", "typescript"}},
	// Rust
	{Name: "Cargo.toml", Type: RootTypeProject, Languages: []string{"rust"}},
	// Python
	{Name: "pyproject.toml", Type: RootTypeProject, Languages: []string{"python"}},
	{Name: "setup.py", Type: RootTypeProject, Languages: []string{"python"}},
	{Name: "setup.cfg", Type: RootTypeProject, Languages: []string{"python"}},
	{Name: "requirements.txt", Type: RootTypeProject, Languages: []string{"python"}},
	// Java - Maven
	{Name: "pom.xml", Type: RootTypeProject, Languages: []string{"java"}},
	// Java/Kotlin - Gradle (build files without settings are subprojects)
	{Name: "build.gradle", Type: RootTypeProject, Languages: []string{"java", "groovy"}},
	{Name: "build.gradle.kts", Type: RootTypeProject, Languages: []string{"java", "kotlin"}},
}

// allMarkers returns all markers in priority order for detection.
func allMarkers() []Marker {
	markers := make([]Marker, 0, len(workspaceMarkers)+len(projectMarkers)+len(vcsMarkers))
	markers = append(markers, workspaceMarkers...)
	markers = append(markers, projectMarkers...)
	markers = append(markers, vcsMarkers...)
	return markers
}
