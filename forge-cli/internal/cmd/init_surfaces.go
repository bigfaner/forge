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
var conflictStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF8700")).Bold(true)

// surfaceStyle highlights surface type values.
var surfaceStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7DCFFF")).Bold(true)

// askSurfaceConfirmation runs the TUI surface confirmation flow.
// It detects surfaces, displays them, and lets the user confirm/edit.
// Returns the confirmed SurfacesMap, or nil if cancelled.
func askSurfaceConfirmation(projectRoot string) (forgeconfig.SurfacesMap, bool) {
	result, err := forgeconfig.DetectSurfacesWithConflicts(projectRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: surface detection failed: %v\n", err)
		return manualSurfaceEntry()
	}

	if len(result.Surfaces) == 0 {
		return manualSurfaceEntry()
	}

	// Build conflict lookup
	conflictMap := make(map[string]*forgeconfig.PathConflict)
	for i := range result.Conflicts {
		conflictMap[result.Conflicts[i].Path] = &result.Conflicts[i]
	}

	// Display detected surfaces
	if result.IsScalar {
		return askScalarConfirmation(result, conflictMap)
	}
	return askMapConfirmation(result, conflictMap)
}

// askScalarConfirmation handles the single-type (scalar) TUI flow.
// Hard Rule: single-type detection should NOT show path column.
func askScalarConfirmation(result *forgeconfig.DetectResult, conflictMap map[string]*forgeconfig.PathConflict) (forgeconfig.SurfacesMap, bool) {
	surfaceType := result.Surfaces["."]

	// Build display description
	desc := fmt.Sprintf("Detected surface type: %s", surfaceStyle.Render(surfaceType))
	if c, ok := conflictMap["."]; ok {
		desc += "\n" + conflictStyle.Render(formatConflictAnnotation(c))
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

	for {
		// Display current surfaces and ask what to do
		lines := buildDisplayLines(surfaces, conflictMap)
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
func buildDisplayLines(surfaces forgeconfig.SurfacesMap, conflictMap map[string]*forgeconfig.PathConflict) []string {
	// Sort paths for consistent display
	paths := make([]string, 0, len(surfaces))
	for p := range surfaces {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	var lines []string
	lines = append(lines, "Detected surfaces:")
	for _, p := range paths {
		surfaceType := surfaces[p]
		line := fmt.Sprintf("  %s: %s", p, surfaceStyle.Render(surfaceType))
		if c, ok := conflictMap[p]; ok {
			line += " " + conflictStyle.Render(formatConflictAnnotation(c))
		}
		lines = append(lines, line)
	}
	lines = append(lines, "", "Actions: Confirm / Edit / Add / Delete")
	return lines
}

// formatConflictAnnotation produces the conflict annotation string.
// Format: (冲突信号: web + api，已按优先级选择 web)
func formatConflictAnnotation(c *forgeconfig.PathConflict) string {
	return fmt.Sprintf("(冲突信号: %s，已按优先级选择 %s)",
		strings.Join(c.Conflicting, " + "), c.Resolved)
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
	// Select which entry to edit
	paths := sortedPaths(surfaces)
	opts := make([]huh.Option[string], len(paths))
	for i, p := range paths {
		opts[i] = huh.NewOption(fmt.Sprintf("%s: %s", p, surfaces[p]), p)
	}

	var selected string
	selectForm := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().
			Title("Select entry to edit").
			Options(opts...).
			Value(&selected),
	))
	if err := selectForm.Run(); err != nil {
		return surfaces, !errors.Is(err, huh.ErrUserAborted)
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
	paths := sortedPaths(surfaces)
	opts := make([]huh.Option[string], len(paths))
	for i, p := range paths {
		opts[i] = huh.NewOption(fmt.Sprintf("%s: %s", p, surfaces[p]), p)
	}

	var selected string
	selectForm := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().
			Title("Select entry to delete").
			Options(opts...).
			Value(&selected),
	))
	if err := selectForm.Run(); err != nil {
		return surfaces, !errors.Is(err, huh.ErrUserAborted)
	}

	delete(surfaces, selected)
	delete(conflictMap, selected)
	return surfaces, false
}

// manualSurfaceEntry is the fallback when detection finds nothing.
func manualSurfaceEntry() (forgeconfig.SurfacesMap, bool) {
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
