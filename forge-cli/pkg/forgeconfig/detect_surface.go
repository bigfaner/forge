package forgeconfig

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"forge-cli/pkg/types"
)

// DefaultDetectDepth is the default directory traversal depth for surface detection.
const DefaultDetectDepth = 3

// MaxDetectDepth is the maximum allowed detection depth.
const MaxDetectDepth = 10

// ErrInvalidDepth is returned when FORGE_DETECT_DEPTH is set to an invalid value (0 or negative).
type ErrInvalidDepth struct {
	Value int
}

func (e *ErrInvalidDepth) Error() string {
	return fmt.Sprintf("FORGE_DETECT_DEPTH=%d is invalid; valid range is 1-%d", e.Value, MaxDetectDepth)
}

// excludedDirs lists directory names that are skipped during detection traversal.
// These directories cannot contain valid surface signals.
var excludedDirs = map[string]bool{
	"node_modules": true,
	".git":         true,
	"vendor":       true,
	"dist":         true,
	"build":        true,
	"__pycache__":  true,
	".next":        true,
	"target":       true,
	"tests":        true,
}

// surfacePriority defines the conflict resolution priority order.
// Lower number = higher priority.
var surfacePriority = map[types.SurfaceType]int{
	types.SurfaceWeb:    1,
	types.SurfaceMobile: 2,
	types.SurfaceAPI:    3,
	types.SurfaceTUI:    4,
	types.SurfaceCLI:    5,
}

// packageJSONDeps maps dependency names to their detected surface types.
var packageJSONSignals = map[string]types.SurfaceType{
	// web signals
	"react":     types.SurfaceWeb,
	"react-dom": types.SurfaceWeb,
	"vue":       types.SurfaceWeb,
	"svelte":    types.SurfaceWeb,
	"next":      types.SurfaceWeb,
	"nuxt":      types.SurfaceWeb,
	"angular":   types.SurfaceWeb,
	// mobile signals
	"react-native": types.SurfaceMobile,
	"expo":         types.SurfaceMobile,
	// api signals
	"express":      types.SurfaceAPI,
	"fastify":      types.SurfaceAPI,
	"koa":          types.SurfaceAPI,
	"hapi":         types.SurfaceAPI,
	"@hapi/hapi":   types.SurfaceAPI,
	"nestjs":       types.SurfaceAPI,
	"@nestjs/core": types.SurfaceAPI,
	// cli signals
	"commander": types.SurfaceCLI,
	"yargs":     types.SurfaceCLI,
	"oclif":     types.SurfaceCLI,
	"inquirer":  types.SurfaceCLI,
	// tui signals
	"blessed":     types.SurfaceTUI,
	"neo-blessed": types.SurfaceTUI,
	"ink":         types.SurfaceTUI,
}

// goModSignals maps go.mod require paths (prefix-matched) to surface types.
var goModSignals = map[string]types.SurfaceType{
	// api signals
	"github.com/gin-gonic/":     types.SurfaceAPI,
	"github.com/labstack/echo/": types.SurfaceAPI,
	"github.com/go-chi/":        types.SurfaceAPI,
	"github.com/gorilla/":       types.SurfaceAPI, // mux, csrf, etc.
	// cli signals
	"github.com/spf13/cobra": types.SurfaceCLI,
	"github.com/urfave/":     types.SurfaceCLI,
	// tui signals
	"github.com/charmbracelet/bubbletea": types.SurfaceTUI,
	"github.com/charmbracelet/lipgloss":  types.SurfaceTUI,
	"github.com/rivo/tview":              types.SurfaceTUI,
	"github.com/gdamore/tcell/":          types.SurfaceTUI,
}

// cargoTomlSignals maps Cargo.toml dependency names (prefix-matched) to surface types.
var cargoTomlSignals = map[string]types.SurfaceType{
	// api signals
	"actix-web": types.SurfaceAPI,
	"axum":      types.SurfaceAPI,
	"rocket":    types.SurfaceAPI,
	"warp":      types.SurfaceAPI,
	// cli signals
	"clap":      types.SurfaceCLI,
	"structopt": types.SurfaceCLI,
	// tui signals
	"ratatui": types.SurfaceTUI,
	"cursive": types.SurfaceTUI,
}

// pyProjectSignals maps Python dependency names to surface types.
var pyProjectSignals = map[string]types.SurfaceType{
	// api signals
	"flask":     types.SurfaceAPI,
	"fastapi":   types.SurfaceAPI,
	"django":    types.SurfaceAPI,
	"starlette": types.SurfaceAPI,
	// cli signals
	"click":    types.SurfaceCLI,
	"typer":    types.SurfaceCLI,
	"argparse": types.SurfaceCLI,
}

// PathConflict records conflict metadata for a single detected path.
// When multiple signals (e.g., web + api) are detected at the same path,
// the conflict is auto-resolved by priority but the conflicting signals
// are preserved for TUI annotation.
type PathConflict struct {
	Path        string              // relative path (or "." for root)
	Resolved    types.SurfaceType   // the chosen surface type after priority resolution
	Conflicting []types.SurfaceType // all surface types that conflicted (2+ entries)
}

// DetectResult holds the output of surface detection, including conflict metadata.
type DetectResult struct {
	Surfaces  SurfacesMap    // the resolved surfaces map
	Conflicts []PathConflict // entries that had signal conflicts
	IsScalar  bool           // true when single-type project (surfaces has only "." key)
	Sources   SourcesMap     // how each surface was determined (e.g., "inference:cmd-dir", "dependency:cobra")
}

// SourcesMap tracks how each detected surface was determined.
// Keys are relative paths (or "." for root), values are source annotations
// following the pattern "inference:<rule-id>" or "dependency:<signal>".
type SourcesMap map[string]string

// DetectSurfaces scans a project directory for surface signals and returns a
// SurfacesMap with detected surface types. For single-type projects the key is "."
// (scalar form). For monorepo/workspace projects, keys are relative paths.
//
// Detection rules:
//   - Default depth: 3 levels, configurable via FORGE_DETECT_DEPTH (1-10)
//   - FORGE_DETECT_DEPTH=0 or negative produces an error (Hard Rule)
//   - Workspace detection: pnpm-workspace.yaml or package.json#workspaces
//   - Root deps are skipped in workspace mode
//   - Signal conflicts resolved via priority: web > mobile > api > cli > tui
func DetectSurfaces(projectRoot string) (SurfacesMap, error) {
	result, err := DetectSurfacesWithConflicts(projectRoot)
	if err != nil {
		return nil, err
	}
	return result.Surfaces, nil
}

// DetectSurfacesWithConflicts scans for surfaces and returns full detection
// metadata including conflict information for TUI annotation.
func DetectSurfacesWithConflicts(projectRoot string) (*DetectResult, error) {
	depth, err := resolveDetectDepth()
	if err != nil {
		return nil, err
	}

	isWorkspace := detectWorkspaceMode(projectRoot)
	result := make(SurfacesMap)
	sources := make(SourcesMap)
	var conflicts []PathConflict

	if isWorkspace {
		// Skip root deps, scan subdirs (sources not tracked in workspace mode)
		scanSubdirsWithSources(projectRoot, projectRoot, 0, depth, result, nil, nil)
	} else {
		// Scan root for signals
		if surface, source, conflict := detectSurfaceAtDirWithSources(projectRoot); surface != "" {
			result["."] = string(surface)
			sources["."] = source
			if len(conflict) > 1 {
				conflicts = append(conflicts, PathConflict{
					Path: ".", Resolved: surface, Conflicting: conflict,
				})
			}
		}
		// Also scan subdirs for mixed signals
		scanSubdirsWithSources(projectRoot, projectRoot, 0, depth, result, sources, &conflicts)
	}

	// Collapse to scalar form if only one surface detected.
	// Non-workspace: normalize subdirectory key to "." for consistent scalar handling.
	// Workspace: preserve subdirectory path keys.
	if len(result) == 1 {
		if !isWorkspace {
			for k, v := range result {
				if k != "." {
					result = SurfacesMap{".": v}
					if s, ok := sources[k]; ok {
						sources = SourcesMap{".": s}
					}
				}
			}
		}
		return &DetectResult{
			Surfaces:  result,
			Conflicts: conflicts,
			IsScalar:  true,
			Sources:   sources,
		}, nil
	}

	// If all paths have the same surface type and one of them is ".", collapse
	// (This handles non-workspace projects that found signals in subdirs too)
	if _, hasDot := result["."]; hasDot {
		seen := make(map[string]bool)
		for _, v := range result {
			seen[v] = true
		}
		if len(seen) == 1 {
			for _, v := range result {
				dotSource := sources["."]
				return &DetectResult{
					Surfaces:  SurfacesMap{".": v},
					Conflicts: nil, // collapsed, no conflicts to show
					IsScalar:  true,
					Sources:   SourcesMap{".": dotSource},
				}, nil
			}
		}
	}

	return &DetectResult{
		Surfaces:  result,
		Conflicts: conflicts,
		IsScalar:  false,
		Sources:   sources,
	}, nil
}

// resolveDetectDepth reads FORGE_DETECT_DEPTH env var and validates it.
// Returns DefaultDetectDepth if not set. Returns ErrInvalidDepth for 0 or negative.
func resolveDetectDepth() (int, error) {
	val := os.Getenv("FORGE_DETECT_DEPTH")
	if val == "" {
		return DefaultDetectDepth, nil
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return 0, &ErrInvalidDepth{Value: 0}
	}
	if n <= 0 {
		return 0, &ErrInvalidDepth{Value: n}
	}
	if n > MaxDetectDepth {
		n = MaxDetectDepth
	}
	return n, nil
}

// detectWorkspaceMode returns true if the project uses workspace/monorepo setup.
func detectWorkspaceMode(root string) bool {
	// Check for pnpm-workspace.yaml
	if _, err := os.Stat(filepath.Join(root, "pnpm-workspace.yaml")); err == nil {
		return true
	}

	// Check for package.json with workspaces field
	pkgPath := filepath.Join(root, "package.json")
	data, err := os.ReadFile(pkgPath)
	if err != nil {
		return false
	}

	var pkg struct {
		Workspaces interface{} `json:"workspaces"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return false
	}
	return pkg.Workspaces != nil
}

// scanSubdirsWithSources recursively scans subdirectories for surface signals
// and records both surface types and source annotations.
// When sources is nil, source tracking is skipped (used by workspace mode).
func scanSubdirsWithSources(root, currentDir string, currentDepth, maxDepth int, result SurfacesMap, sources SourcesMap, conflicts *[]PathConflict) {
	if currentDepth >= maxDepth {
		return
	}

	entries, err := os.ReadDir(currentDir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()

		if isExcludedDir(name) {
			continue
		}

		subdirPath := filepath.Join(currentDir, name)

		surface, source, conflict := detectSurfaceAtDirWithSources(subdirPath)
		if surface != "" {
			rel, err := filepath.Rel(root, subdirPath)
			if err != nil {
				continue
			}
			rel = filepath.ToSlash(rel)
			result[rel] = string(surface)
			if sources != nil {
				sources[rel] = source
			}
			if len(conflict) > 1 && conflicts != nil {
				*conflicts = append(*conflicts, PathConflict{
					Path: rel, Resolved: surface, Conflicting: conflict,
				})
			}
		}

		scanSubdirsWithSources(root, subdirPath, currentDepth+1, maxDepth, result, sources, conflicts)
	}
}

// detectSurfaceAtDirWithSources detects the surface type, source annotation, and conflict metadata.
// Returns the resolved surface type, source annotation, and the list of all conflicting signal types.
// Priority chain: dependency signals first; inference only when ALL dependency signals return empty.
func detectSurfaceAtDirWithSources(dir string) (types.SurfaceType, string, []types.SurfaceType) {
	// Collect all unique signal types across all manifest files
	seen := make(map[types.SurfaceType]bool)
	var allSignals []types.SurfaceType

	collect := func(signals []types.SurfaceType) {
		for _, s := range signals {
			if !seen[s] {
				seen[s] = true
				allSignals = append(allSignals, s)
			}
		}
	}

	// Dependency signal detection
	nodeSignals := detectPackageJSONSignals(dir)
	collect(nodeSignals)

	goSignals := detectGoModSignals(dir)
	collect(goSignals)

	cargoSignals := detectCargoTomlSignals(dir)
	collect(cargoSignals)

	// Check mobile signals
	if detectMobile(dir) {
		if !seen[types.SurfaceMobile] {
			allSignals = append(allSignals, types.SurfaceMobile)
		}
	}

	pySignals := detectPyProjectSignals(dir)
	collect(pySignals)

	// Dependency signals found — skip inference entirely
	if len(allSignals) > 0 {
		resolved := resolveConflict(allSignals)

		// Determine which dependency caused the resolved signal
		source := detectDependencySource(dir, resolved, nodeSignals, goSignals, cargoSignals, pySignals)

		return resolved, source, allSignals
	}

	// All dependency signals empty — try structural inference
	inferredType, inferredSource := runInference(dir)
	if inferredType != "" {
		return inferredType, inferredSource, nil
	}

	return "", "", nil
}

// detectDependencySource identifies which dependency caused the resolved surface type.
func detectDependencySource(dir string, resolved types.SurfaceType, nodeSignals, goSignals, cargoSignals, pySignals []types.SurfaceType) string {
	// Check if the resolved type came from Node.js deps
	for _, s := range nodeSignals {
		if s == resolved {
			if dep := findDepForSignal(dir, "node", resolved); dep != "" {
				return "dependency:" + dep
			}
		}
	}
	// Check Go deps
	for _, s := range goSignals {
		if s == resolved {
			if dep := findDepForSignal(dir, "go", resolved); dep != "" {
				return "dependency:" + dep
			}
		}
	}
	// Check Python deps
	for _, s := range pySignals {
		if s == resolved {
			if dep := findDepForSignal(dir, "python", resolved); dep != "" {
				return "dependency:" + dep
			}
		}
	}
	// Check Cargo deps
	for _, s := range cargoSignals {
		if s == resolved {
			if dep := findDepForSignal(dir, "cargo", resolved); dep != "" {
				return "dependency:" + dep
			}
		}
	}
	// Fallback for mobile
	if resolved == types.SurfaceMobile {
		return "dependency:mobile-manifest"
	}
	return "dependency:unknown"
}

// findDepForSignal finds the first dependency name that maps to the given surface type.
func findDepForSignal(dir string, ecosystem string, targetSurface types.SurfaceType) string {
	switch ecosystem {
	case "node":
		pkgPath := filepath.Join(dir, "package.json")
		data, err := os.ReadFile(pkgPath)
		if err != nil {
			return ""
		}
		var pkg struct {
			Dependencies    map[string]interface{} `json:"dependencies"`
			DevDependencies map[string]interface{} `json:"devDependencies"`
		}
		if err := json.Unmarshal(data, &pkg); err != nil {
			return ""
		}
		for dep := range pkg.Dependencies {
			if surface, ok := packageJSONSignals[dep]; ok && surface == targetSurface {
				return dep
			}
		}
		for dep := range pkg.DevDependencies {
			if surface, ok := packageJSONSignals[dep]; ok && surface == targetSurface {
				return dep
			}
		}
	case "go":
		modPath := filepath.Join(dir, "go.mod")
		data, err := os.ReadFile(modPath)
		if err != nil {
			return ""
		}
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "//") {
				continue
			}
			for prefix, surface := range goModSignals {
				if surface == targetSurface &&
					(strings.HasPrefix(line, prefix) || strings.Contains(line, " "+prefix)) {
					// Extract short name: last segment of the import path
					trimmed := strings.TrimSuffix(prefix, "/")
					if idx := strings.LastIndex(trimmed, "/"); idx >= 0 {
						return trimmed[idx+1:]
					}
					return trimmed
				}
			}
		}
	case "python":
		tomlPath := filepath.Join(dir, "pyproject.toml")
		data, err := os.ReadFile(tomlPath)
		if err != nil {
			return ""
		}
		content := string(data)
		for dep, surface := range pyProjectSignals {
			if surface == targetSurface && strings.Contains(content, dep) {
				return dep
			}
		}
	case "cargo":
		tomlPath := filepath.Join(dir, "Cargo.toml")
		data, err := os.ReadFile(tomlPath)
		if err != nil {
			return ""
		}
		content := string(data)
		for dep, surface := range cargoTomlSignals {
			if surface == targetSurface && strings.Contains(content, dep) {
				return dep
			}
		}
	}
	return ""
}

// runInference determines which ecosystem's inference function to call based on
// manifest file presence, then calls the matching function.
func runInference(dir string) (types.SurfaceType, string) {
	// Ecosystem determined by manifest file presence; only one ecosystem's
	// inference function is called. Priority: go.mod > package.json > pyproject.toml
	// to match the order dependencies were checked.
	if fileExists(filepath.Join(dir, "go.mod")) {
		return inferGoSurface(dir)
	}
	if fileExists(filepath.Join(dir, "package.json")) {
		return inferNodeSurface(dir)
	}
	if fileExists(filepath.Join(dir, "pyproject.toml")) || fileExists(filepath.Join(dir, "setup.py")) {
		return inferPythonSurface(dir)
	}
	return "", ""
}

// detectPackageJSONSignals reads a package.json and returns all detected surface signal types.
func detectPackageJSONSignals(dir string) []types.SurfaceType {
	pkgPath := filepath.Join(dir, "package.json")
	data, err := os.ReadFile(pkgPath)
	if err != nil {
		return nil
	}

	var pkg struct {
		Dependencies    map[string]interface{} `json:"dependencies"`
		DevDependencies map[string]interface{} `json:"devDependencies"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil
	}

	// Merge both dep maps for signal detection
	allDeps := make(map[string]bool)
	for k := range pkg.Dependencies {
		allDeps[k] = true
	}
	for k := range pkg.DevDependencies {
		allDeps[k] = true
	}

	var signals []types.SurfaceType
	for dep := range allDeps {
		if surface, ok := packageJSONSignals[dep]; ok {
			signals = append(signals, surface)
		}
	}

	// Special rule: react-native + react → mobile (react is a shared dependency of react-native,
	// not an independent web signal). When react-native is present, suppress the web signal from react.
	if allDeps["react-native"] {
		filtered := make([]types.SurfaceType, 0, len(signals))
		hasMobile := false
		for _, s := range signals {
			if s == types.SurfaceMobile {
				hasMobile = true
			}
			if s != types.SurfaceWeb {
				filtered = append(filtered, s)
			}
		}
		// Only apply the suppression when both react-native and react exist
		// (mobile signal present alongside web)
		if hasMobile {
			signals = filtered
		}
	}

	return dedupSignals(signals)
}

// detectGoModSignals reads a go.mod and returns all detected surface signal types.
func detectGoModSignals(dir string) []types.SurfaceType {
	modPath := filepath.Join(dir, "go.mod")
	data, err := os.ReadFile(modPath)
	if err != nil {
		return nil
	}

	seen := make(map[types.SurfaceType]bool)
	var signals []types.SurfaceType
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		// Skip comments and non-require lines
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		for prefix, surface := range goModSignals {
			if strings.HasPrefix(line, prefix) || strings.Contains(line, " "+prefix) {
				if !seen[surface] {
					seen[surface] = true
					signals = append(signals, surface)
				}
			}
		}
	}

	return signals
}

// detectCargoTomlSignals reads a Cargo.toml and returns all detected surface signal types.
func detectCargoTomlSignals(dir string) []types.SurfaceType {
	tomlPath := filepath.Join(dir, "Cargo.toml")
	data, err := os.ReadFile(tomlPath)
	if err != nil {
		return nil
	}

	seen := make(map[types.SurfaceType]bool)
	var signals []types.SurfaceType
	content := string(data)

	for dep, surface := range cargoTomlSignals {
		// Simple prefix match in the dependencies section
		if strings.Contains(content, dep) {
			if !seen[surface] {
				seen[surface] = true
				signals = append(signals, surface)
			}
		}
	}

	return signals
}

// detectMobile checks for mobile project signals (AndroidManifest.xml, *.xcodeproj).
func detectMobile(dir string) bool {
	// Check for AndroidManifest.xml anywhere under this directory (depth 3)
	if found := findFile(dir, "AndroidManifest.xml", 0, 3); found {
		return true
	}

	// Check for *.xcodeproj entries (files or directories)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".xcodeproj") {
			return true
		}
	}

	// Check for pubspec.yaml with flutter dependency
	pubspecPath := filepath.Join(dir, "pubspec.yaml")
	data, err := os.ReadFile(pubspecPath)
	if err == nil && strings.Contains(string(data), "flutter") {
		return true
	}

	return false
}

// detectPyProjectSignals reads a pyproject.toml and returns all detected surface signal types.
func detectPyProjectSignals(dir string) []types.SurfaceType {
	tomlPath := filepath.Join(dir, "pyproject.toml")
	data, err := os.ReadFile(tomlPath)
	if err != nil {
		return nil
	}

	seen := make(map[types.SurfaceType]bool)
	var signals []types.SurfaceType
	content := string(data)

	for dep, surface := range pyProjectSignals {
		if strings.Contains(content, dep) {
			if !seen[surface] {
				seen[surface] = true
				signals = append(signals, surface)
			}
		}
	}

	return signals
}

// dedupSignals deduplicates a slice of signal types while preserving order.
func dedupSignals(signals []types.SurfaceType) []types.SurfaceType {
	seen := make(map[types.SurfaceType]bool)
	var result []types.SurfaceType
	for _, s := range signals {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}

// resolveConflict takes a list of detected surface signals and returns the
// highest priority one. Priority: web > mobile > api > cli > tui.
// Returns "" if the list is empty.
func resolveConflict(signals []types.SurfaceType) types.SurfaceType {
	if len(signals) == 0 {
		return ""
	}
	if len(signals) == 1 {
		return signals[0]
	}

	best := signals[0]
	bestPriority := surfacePriority[best]

	for _, s := range signals[1:] {
		p := surfacePriority[s]
		if p < bestPriority {
			best = s
			bestPriority = p
		}
	}

	return best
}

// isExcludedDir returns true if the directory name should be skipped during traversal.
func isExcludedDir(name string) bool {
	return excludedDirs[name]
}

// findFile recursively searches for a file with the given name up to maxDepth levels.
func findFile(dir, filename string, currentDepth, maxDepth int) bool {
	if currentDepth > maxDepth {
		return false
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if !entry.IsDir() && entry.Name() == filename {
			return true
		}
		if entry.IsDir() && !isExcludedDir(entry.Name()) {
			if findFile(filepath.Join(dir, entry.Name()), filename, currentDepth+1, maxDepth) {
				return true
			}
		}
	}
	return false
}

// --- Structural inference functions ---

// hasSubdirs returns true if the named directory exists and contains at least
// one subdirectory (i.e., it has child directories, not just files).
func hasSubdirs(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if e.IsDir() {
			return true
		}
	}
	return false
}

// dirExists returns true if path exists and is a directory.
func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// fileExists returns true if path exists and is a regular file.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// inferGoSurface inspects a Go project directory structure to infer surface type.
// Rules:
//   - cmd/ with subdirectories → cli (inference:cmd-dir)
//   - api/ directory → api (inference:api-dir)
//   - handler/ directory → api (inference:handler-dir)
//   - Both cmd/ and api/handler/ → api wins, cli discarded
//
// Returns empty strings if no inference possible.
// Only called when go.mod exists and all dependency signals are empty.
func inferGoSurface(dir string) (types.SurfaceType, string) {
	// Verify go.mod exists (ecosystem gate)
	if !fileExists(filepath.Join(dir, "go.mod")) {
		return "", ""
	}

	hasCmdSubdirs := false
	hasAPI := false
	hasHandler := false

	cmdDir := filepath.Join(dir, "cmd")
	if dirExists(cmdDir) && hasSubdirs(cmdDir) {
		hasCmdSubdirs = true
	}

	if dirExists(filepath.Join(dir, "api")) {
		hasAPI = true
	}
	if dirExists(filepath.Join(dir, "handler")) {
		hasHandler = true
	}

	// api/handler wins over cmd when both present
	if hasAPI {
		return types.SurfaceAPI, "inference:api-dir"
	}
	if hasHandler {
		return types.SurfaceAPI, "inference:handler-dir"
	}
	if hasCmdSubdirs {
		return types.SurfaceCLI, "inference:cmd-dir"
	}

	return "", ""
}

// inferNodeSurface inspects a Node.js project directory structure to infer surface type.
// Rules:
//   - bin field in package.json → cli (inference:bin-field)
//   - index.html at project root (same dir as package.json) → web (inference:index-html)
//   - Does NOT scan subdirectories for index.html
//   - Both bin and index.html → web wins (higher priority)
//
// Returns empty strings if no inference possible.
func inferNodeSurface(dir string) (types.SurfaceType, string) {
	// Verify package.json exists (ecosystem gate)
	pkgPath := filepath.Join(dir, "package.json")
	if !fileExists(pkgPath) {
		return "", ""
	}

	hasBin := false
	hasIndexHTML := false

	// Parse package.json for bin field
	data, err := os.ReadFile(pkgPath)
	if err != nil {
		return "", ""
	}

	var pkg struct {
		Bin interface{} `json:"bin"`
	}
	// recover from malformed JSON
	func() {
		defer func() {
			_ = recover()
		}()
		_ = json.Unmarshal(data, &pkg)
	}()

	if pkg.Bin != nil {
		// bin can be string or object, both are valid
		hasBin = true
	}

	// Check index.html at project root only
	if fileExists(filepath.Join(dir, "index.html")) {
		hasIndexHTML = true
	}

	// web (index.html) wins over cli (bin) by priority
	if hasIndexHTML {
		return types.SurfaceWeb, "inference:index-html"
	}
	if hasBin {
		return types.SurfaceCLI, "inference:bin-field"
	}

	return "", ""
}

// inferPythonSurface inspects a Python project directory structure to infer surface type.
// Rules (checked in order, first match wins):
//  1. [project.scripts] in pyproject.toml → cli (inference:py-scripts)
//  2. entry_points in setup.py → cli (inference:py-scripts)
//  3. app.py or main.py at root → cli (inference:py-main), ONLY when:
//     - no setup.py with name matching directory name
//     - no [project.packages] or [tool.setuptools.packages.find] in pyproject.toml
//
// Returns empty strings if no inference possible.
func inferPythonSurface(dir string) (types.SurfaceType, string) {
	hasPyProject := fileExists(filepath.Join(dir, "pyproject.toml"))
	hasSetupPy := fileExists(filepath.Join(dir, "setup.py"))

	// Check [project.scripts] in pyproject.toml
	if hasPyProject {
		data, err := os.ReadFile(filepath.Join(dir, "pyproject.toml"))
		if err == nil {
			content := string(data)
			if strings.Contains(content, "[project.scripts]") {
				return types.SurfaceCLI, "inference:py-scripts"
			}
		}
	}

	// Check entry_points in setup.py
	if hasSetupPy {
		data, err := os.ReadFile(filepath.Join(dir, "setup.py"))
		if err == nil {
			content := string(data)
			if strings.Contains(content, "entry_points") {
				return types.SurfaceCLI, "inference:py-scripts"
			}
		}
	}

	// Check app.py / main.py with library exclusion
	hasAppPy := fileExists(filepath.Join(dir, "app.py"))
	hasMainPy := fileExists(filepath.Join(dir, "main.py"))

	if hasAppPy || hasMainPy {
		// Exclusion: setup.py with name matching directory name
		if hasSetupPy {
			data, err := os.ReadFile(filepath.Join(dir, "setup.py"))
			if err == nil {
				content := string(data)
				dirName := filepath.Base(dir)
				if strings.Contains(content, "name='"+dirName+"'") ||
					strings.Contains(content, "name=\""+dirName+"\"") {
					return "", ""
				}
			}
		}

		// Exclusion: [project.packages] or [tool.setuptools.packages.find] in pyproject.toml
		if hasPyProject {
			data, err := os.ReadFile(filepath.Join(dir, "pyproject.toml"))
			if err == nil {
				content := string(data)
				if strings.Contains(content, "[project.packages]") ||
					strings.Contains(content, "[tool.setuptools.packages.find]") {
					return "", ""
				}
			}
		}

		return types.SurfaceCLI, "inference:py-main"
	}

	return "", ""
}
