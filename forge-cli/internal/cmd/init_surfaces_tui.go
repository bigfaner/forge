package cmd

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"forge-cli/pkg/forgeconfig"

	"github.com/charmbracelet/huh"
)

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
