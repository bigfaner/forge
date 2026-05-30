package cmd

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"forge-cli/pkg/forgeconfig"

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
		fmt.Fprintf(os.Stderr, "WARNING: surface detection failed: %v\n", err)
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

// askMapAction presents the action selection for map-form surfaces.
func askMapAction(desc string) (string, bool) {
	var action string
	form := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().
			Title("Surface configuration").
			Description(desc).
			Options(
				huh.NewOption("Confirm (save surfaces)", "confirm"),
				huh.NewOption("Edit an entry", "edit"),
				huh.NewOption("Add new mapping", "add"),
				huh.NewOption("Delete an entry", "delete"),
			).
			Value(&action),
	))

	if err := form.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return "", true
		}
		return "", true
	}
	return action, false
}

// editMapEntry lets the user select and edit a single entry.
func editMapEntry(surfaces forgeconfig.SurfacesMap, conflictMap map[string]*forgeconfig.PathConflict) (forgeconfig.SurfacesMap, bool) {
	selected, ok := selectSurfaceEntry("Select entry to edit", surfaces)
	if !ok {
		return surfaces, false
	}
	if selected == "" {
		return surfaces, true // cancelled
	}

	// Edit the selected entry
	var newPath, newSurface string
	editForm := huh.NewForm(huh.NewGroup(
		huh.NewInput().
			Title("Path").
			Description(fmt.Sprintf("Edit path (was: %s)", selected)).
			Value(&newPath),
		huh.NewInput().
			Title("Surface type").
			Description(fmt.Sprintf("Edit surface type (was: %s)", surfaces[selected])).
			Value(&newSurface),
	))
	if err := editForm.Run(); err != nil {
		return surfaces, !errors.Is(err, huh.ErrUserAborted)
	}

	newPath = strings.TrimSpace(newPath)
	newSurface = strings.TrimSpace(newSurface)

	if newPath == "" {
		newPath = selected
	}
	if newSurface == "" {
		newSurface = surfaces[selected]
	}

	// Remove old key if path changed
	if newPath != selected {
		delete(surfaces, selected)
		// Remove conflict metadata for old path
		delete(conflictMap, selected)
	}
	surfaces[newPath] = newSurface
	// Clear conflict for edited entry
	delete(conflictMap, newPath)

	return surfaces, false
}

// addMapEntry lets the user add a new path -> surface mapping.
func addMapEntry(surfaces forgeconfig.SurfacesMap) (forgeconfig.SurfacesMap, bool) {
	var path, surfaceType string
	form := huh.NewForm(huh.NewGroup(
		huh.NewInput().
			Title("Path").
			Description("Enter path relative to project root (e.g., frontend, backend)").
			Value(&path),
		huh.NewInput().
			Title("Surface type").
			Description("Enter surface type (web, api, cli, tui, mobile)").
			Value(&surfaceType),
	))

	if err := form.Run(); err != nil {
		return surfaces, !errors.Is(err, huh.ErrUserAborted)
	}

	path = strings.TrimSpace(path)
	surfaceType = strings.TrimSpace(surfaceType)
	if path != "" && surfaceType != "" {
		surfaces[path] = surfaceType
	}
	return surfaces, false
}

// deleteMapEntry lets the user select and delete an entry.
func deleteMapEntry(surfaces forgeconfig.SurfacesMap, conflictMap map[string]*forgeconfig.PathConflict) (forgeconfig.SurfacesMap, bool) {
	selected, ok := selectSurfaceEntry("Select entry to delete", surfaces)
	if !ok {
		return surfaces, false
	}
	if selected == "" {
		return surfaces, true // cancelled
	}

	delete(surfaces, selected)
	delete(conflictMap, selected)
	return surfaces, false
}

// selectSurfaceEntry shows a TUI select for a surface entry and returns the selected path.
// Returns ("", true) on cancel, (path, false) on selection.
func selectSurfaceEntry(title string, surfaces forgeconfig.SurfacesMap) (string, bool) {
	paths := sortedPaths(surfaces)
	opts := make([]huh.Option[string], len(paths))
	for i, p := range paths {
		opts[i] = huh.NewOption(fmt.Sprintf("%s: %s", p, surfaces[p]), p)
	}

	var selected string
	form := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().
			Title(title).
			Options(opts...).
			Value(&selected),
	))
	if err := form.Run(); err != nil {
		return "", errors.Is(err, huh.ErrUserAborted)
	}
	return selected, false
}

// manualSurfaceEntry is the function variable for manual surface entry.
// Variable for testability -- huh requires a real TTY, so tests override this.
// Hard Rule: the re-run Edit flow must call this same variable -- no separate code path.
var manualSurfaceEntry = manualSurfaceEntryImpl

// manualSurfaceEntryImpl is the actual implementation of manual surface entry.
func manualSurfaceEntryImpl() (forgeconfig.SurfacesMap, bool) {
	var path, surfaceType string
	form := huh.NewForm(huh.NewGroup(
		huh.NewInput().
			Title("Path (leave empty for project root)").
			Description("Enter path relative to project root, or leave empty for single-surface project").
			Value(&path),
		huh.NewInput().
			Title("Surface type").
			Description("Enter surface type (web, api, cli, tui, mobile)").
			Value(&surfaceType),
	))

	if err := form.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil, true
		}
		return nil, true
	}

	path = strings.TrimSpace(path)
	surfaceType = strings.TrimSpace(surfaceType)
	if surfaceType == "" {
		return nil, false
	}

	if path == "" {
		path = "."
	}
	return forgeconfig.SurfacesMap{path: surfaceType}, false
}

// sortedPaths returns sorted keys of a SurfacesMap.
func sortedPaths(surfaces forgeconfig.SurfacesMap) []string {
	paths := make([]string, 0, len(surfaces))
	for p := range surfaces {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	return paths
}

// formatSurfacesSummary builds a human-readable summary string for the init output.
// Shows actual surface types with compact source annotations instead of opaque "N mappings".
// Uses compact format: (inferred:cmd-dir), (from cobra) — suitable for single-line summary.
func formatSurfacesSummary(surfaces forgeconfig.SurfacesMap, sources forgeconfig.SourcesMap) string {
	if len(surfaces) == 0 {
		return ""
	}

	// Scalar form: single type
	if surfaces["."] != "" && len(surfaces) == 1 {
		detail := surfaces["."]
		if source := sources["."]; source != "" {
			detail += " " + formatCompactSourceAnnotation(source)
		}
		return detail
	}

	// Map form: show each path=type
	paths := sortedPaths(surfaces)
	parts := make([]string, 0, len(paths))
	for _, p := range paths {
		entry := fmt.Sprintf("%s=%s", p, surfaces[p])
		if source, ok := sources[p]; ok && source != "" {
			entry += " " + formatCompactSourceAnnotation(source)
		}
		parts = append(parts, entry)
	}
	return strings.Join(parts, ", ")
}

// askRerunPrompt is the function variable for the re-run prompt.
// Variable for testability -- huh requires a real TTY, so tests override this.
var askRerunPrompt = askRerunPromptImpl

// askRerunPromptImpl presents the user with options when surfaces are already configured.
// Returns: action ("confirm", "redetect", "edit"), cancelled bool.
func askRerunPromptImpl(currentSurfaces forgeconfig.SurfacesMap) (string, bool) {
	// Build summary of current surfaces
	var summary string
	if currentSurfaces["."] != "" && len(currentSurfaces) == 1 {
		summary = currentSurfaces["."]
	} else {
		paths := sortedPaths(currentSurfaces)
		parts := make([]string, 0, len(paths))
		for _, p := range paths {
			parts = append(parts, currentSurfaces[p])
		}
		summary = strings.Join(parts, ", ")
	}

	desc := fmt.Sprintf("Surfaces already configured: %s. Re-detect?", summary)

	var action string
	form := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().
			Title("Surface configuration").
			Description(desc).
			Options(
				huh.NewOption("Confirm (keep existing)", "confirm"),
				huh.NewOption("Re-detect", "redetect"),
				huh.NewOption("Edit (manual entry)", "edit"),
			).
			Value(&action),
	))

	if err := form.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return "", true
		}
		return "", true
	}
	return action, false
}
