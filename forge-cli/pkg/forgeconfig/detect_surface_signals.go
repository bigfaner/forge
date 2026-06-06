package forgeconfig

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"forge-cli/pkg/types"
)

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

// packageJSONSignals maps dependency names to their detected surface types.
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

	// Special rule: react-native + react -> mobile (react is a shared dependency of react-native,
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
