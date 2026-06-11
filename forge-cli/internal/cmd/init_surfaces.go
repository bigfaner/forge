package cmd

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"forge-cli/pkg/forgeconfig"
	"forge-cli/pkg/forgelog"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// conflictStyle highlights conflict annotation text.
var conflictStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorConflict)).Bold(true)

// surfaceStyle highlights surface type values.
var surfaceStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorModeHighlight)).Bold(true)

// sourceStyle styles source annotation text.
var sourceStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorSource)).Italic(true)

// askSurfaceConfirmation is the function variable for TUI surface confirmation.
// Variable for testability — huh requires a real TTY, so tests override this.
var askSurfaceConfirmation = askSurfaceConfirmationImpl

// askSurfaceConfirmationImpl runs the TUI surface confirmation flow.
// It detects surfaces, displays them, and lets the user confirm/edit.
// Returns the confirmed SurfacesMap, SourcesMap, or nil if cancelled.
func askSurfaceConfirmationImpl(projectRoot string) (forgeconfig.SurfacesMap, forgeconfig.SourcesMap, bool) {
	result, err := forgeconfig.DetectSurfacesWithConflicts(projectRoot)
	if err != nil {
		forgelog.Warn("WARNING: surface detection failed: %v\n", err)
		surfaces, cancelled := manualSurfaceEntry()
		return surfaces, nil, cancelled
	}

	if len(result.Surfaces) == 0 {
		surfaces, cancelled := manualSurfaceEntry()
		return surfaces, nil, cancelled
	}

	// Build conflict lookup
	conflictMap := make(map[string]*forgeconfig.PathConflict)
	for i := range result.Conflicts {
		conflictMap[result.Conflicts[i].Path] = &result.Conflicts[i]
	}

	// Display detected surfaces
	if result.IsScalar {
		surfaces, cancelled := askScalarConfirmation(result, conflictMap)
		return surfaces, result.Sources, cancelled
	}
	surfaces, cancelled := askMapConfirmation(result, conflictMap)
	return surfaces, result.Sources, cancelled
}

// askScalarConfirmation handles the single-type (scalar) TUI flow.
// Hard Rule: single-type detection should NOT show path column.
func askScalarConfirmation(result *forgeconfig.DetectResult, conflictMap map[string]*forgeconfig.PathConflict) (forgeconfig.SurfacesMap, bool) {
	surfaceType := result.Surfaces["."]

	// Build display description with source annotation
	desc := fmt.Sprintf("Detected surface type: %s", surfaceStyle.Render(surfaceType))

	if source := result.Sources["."]; source != "" {
		desc += " " + sourceStyle.Render(formatSourceAnnotation(source))
	}

	if c, ok := conflictMap["."]; ok {
		desc += "\n" + conflictStyle.Render(formatConflictAnnotation(c))
	}

	// Add hint text for inferred surfaces
	if isInferred(result.Sources["."]) {
		desc += "\n\nThis was inferred from project structure. Edit to correct if needed."
	}

	confirm := true
	form := huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title("Confirm detected surface").
			Description(desc).
			Affirmative("Confirm").
			Negative("Edit").
			Value(&confirm),
	))

	if err := form.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil, true
		}
		return nil, true
	}

	if confirm {
		return result.Surfaces, false
	}

	// Edit mode: let user enter a surface type
	var edited string
	editForm := huh.NewForm(huh.NewGroup(
		huh.NewInput().
			Title("Surface type").
			Description(fmt.Sprintf("Enter surface type (was: %s)", surfaceType)).
			Placeholder("web, api, cli, tui, mobile").
			Value(&edited),
	))
	if err := editForm.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil, true
		}
		return nil, true
	}

	edited = strings.TrimSpace(edited)
	if edited == "" {
		return result.Surfaces, false
	}
	return forgeconfig.SurfacesMap{".": edited}, false
}

// askMapConfirmation handles the multi-type (map) TUI flow.
func askMapConfirmation(result *forgeconfig.DetectResult, conflictMap map[string]*forgeconfig.PathConflict) (forgeconfig.SurfacesMap, bool) {
	surfaces := result.Surfaces
	sources := result.Sources

	for {
		// Display current surfaces and ask what to do
		lines := buildDisplayLines(surfaces, conflictMap, sources)
		desc := strings.Join(lines, "\n")

		action, cancelled := askMapAction(desc)
		if cancelled {
			return nil, true
		}

		switch action {
		case "confirm":
			return surfaces, false
		case "edit":
			var ok bool
			surfaces, ok = editMapEntry(surfaces, conflictMap)
			if !ok {
				return nil, true
			}
		case "add":
			var ok bool
			surfaces, ok = addMapEntry(surfaces)
			if !ok {
				return nil, true
			}
		case "delete":
			var ok bool
			surfaces, ok = deleteMapEntry(surfaces, conflictMap)
			if !ok {
				return nil, true
			}
		}
	}
}

// buildDisplayLines creates the display lines for map-form surfaces.
func buildDisplayLines(surfaces forgeconfig.SurfacesMap, conflictMap map[string]*forgeconfig.PathConflict, sources forgeconfig.SourcesMap) []string {
	// Sort paths for consistent display
	paths := make([]string, 0, len(surfaces))
	for p := range surfaces {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	var lines []string
	lines = append(lines, "Detected surfaces:")
	hasInferred := false
	for _, p := range paths {
		surfaceType := surfaces[p]
		line := fmt.Sprintf("  %s  →  %s", p, surfaceStyle.Render(surfaceType))
		if source, ok := sources[p]; ok && source != "" {
			line += " " + sourceStyle.Render(formatSourceAnnotation(source))
			if isInferred(source) {
				hasInferred = true
			}
		}
		if c, ok := conflictMap[p]; ok {
			line += " " + conflictStyle.Render(formatConflictAnnotation(c))
		}
		lines = append(lines, line)
	}

	// Add hint text if any surface was inferred
	if hasInferred {
		lines = append(lines, "", "Hint: Inferred entries are based on project structure. Edit to correct if needed.")
	}

	lines = append(lines, "", "Actions: Confirm / Edit / Add / Delete")
	return lines
}

// formatConflictAnnotation produces the conflict annotation string.
// Format: (冲突信号: web + api，已按优先级选择 web)
func formatConflictAnnotation(c *forgeconfig.PathConflict) string {
	// Convert []types.SurfaceType to []string for display
	conflicting := make([]string, len(c.Conflicting))
	for i, s := range c.Conflicting {
		conflicting[i] = string(s)
	}
	return fmt.Sprintf("(冲突信号: %s，已按优先级选择 %s)",
		strings.Join(conflicting, " + "), c.Resolved)
}

// formatSourceAnnotation converts a source annotation code into a human-readable string.
// Input format: "inference:cmd-dir" -> "inferred from cmd/ directory structure"
// Input format: "dependency:cobra" -> "detected from cobra dependency"
func formatSourceAnnotation(source string) string {
	if source == "" {
		return ""
	}

	parts := strings.SplitN(source, ":", 2)
	if len(parts) != 2 {
		return fmt.Sprintf("(%s)", source)
	}

	category := parts[0]
	detail := parts[1]

	switch category {
	case "inference":
		return fmt.Sprintf("(inferred from %s)", formatInferenceDetail(detail))
	case "dependency":
		return fmt.Sprintf("(detected from %s dependency)", detail)
	default:
		return fmt.Sprintf("(%s)", source)
	}
}

// formatCompactSourceAnnotation converts a source annotation code into a compact string
// suitable for the init summary display.
// Input format: "inference:cmd-dir" -> "(inferred:cmd-dir)"
// Input format: "dependency:cobra" -> "(from cobra)"
func formatCompactSourceAnnotation(source string) string {
	if source == "" {
		return ""
	}

	parts := strings.SplitN(source, ":", 2)
	if len(parts) != 2 {
		return fmt.Sprintf("(%s)", source)
	}

	category := parts[0]
	detail := parts[1]

	switch category {
	case "inference":
		return fmt.Sprintf("(inferred:%s)", detail)
	case "dependency":
		return fmt.Sprintf("(from %s)", detail)
	default:
		return fmt.Sprintf("(%s)", source)
	}
}

// formatInferenceDetail converts an inference rule ID into a human-readable description.
func formatInferenceDetail(ruleID string) string {
	switch ruleID {
	case "cmd-dir":
		return "cmd/ directory structure"
	case "api-dir":
		return "api/ directory"
	case "handler-dir":
		return "handler/ directory"
	case "bin-field":
		return "bin field in package.json"
	case "index-html":
		return "index.html at project root"
	case "py-scripts":
		return "project.scripts or entry_points"
	case "py-main":
		return "app.py/main.py at root"
	default:
		return ruleID
	}
}

// isInferred returns true if the source annotation indicates structural inference.
func isInferred(source string) bool {
	return strings.HasPrefix(source, "inference:")
}
