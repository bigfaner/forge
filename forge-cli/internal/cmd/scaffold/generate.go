package scaffold

import (
	"fmt"
	"sort"
	"strings"

	"forge-cli/pkg/types"
)

// recipeName returns the just recipe name with optional key prefix.
// Scalar surfaces (no key) return the verb directly (e.g. "test").
// Named surfaces return "<key>-<verb>" (e.g. "backend-test").
func recipeName(key, verb string) string {
	if key == "" {
		return verb
	}
	return key + "-" + verb
}

// Generate produces justfile recipes for the given surface type and key.
// Returns the generated justfile content as a string.
// For service surface types (api/web/mobile), an orchestration recipe
// (<key> or bare verb) is appended after the standard recipe list.
func Generate(surfaceType types.SurfaceType, key string) (string, error) {
	spec, ok := surfaceSpecs[surfaceType]
	if !ok {
		return "", fmt.Errorf("unknown surface type: %q", surfaceType)
	}

	var b strings.Builder
	for i, r := range spec.Recipes {
		if i > 0 {
			b.WriteString("\n")
		}
		writeRecipe(&b, key, r)
	}

	// For service surfaces, append the orchestration recipe
	if needsKey(surfaceType) {
		if b.Len() > 0 {
			b.WriteString("\n")
		}
		writeOrchestrationRecipe(&b, key, surfaceType)
	}

	return b.String(), nil
}

// writeRecipe writes a single recipe with dual-platform variants.
// All lifecycle and quality recipes are marked with "# user-customized".
// When IsScript is true, the body is indented line-by-line (shebang recipes).
// When IsScript is false, the body is a simple one-liner.
func writeRecipe(b *strings.Builder, key string, r RecipeSpec) {
	name := recipeName(key, r.Name)

	// user-customized marker for lifecycle and quality recipes
	b.WriteString("# user-customized\n")

	// [unix] variant
	fmt.Fprintf(b, "%s [unix]:\n", name)
	writeBody(b, r.UnixBody, r.IsScript)

	// [windows] variant
	fmt.Fprintf(b, "%s [windows]:\n", name)
	writeBody(b, r.WindowsBody, r.IsScript)
}

// writeBody writes a recipe body with proper indentation.
// Script bodies are indented line-by-line; one-liner bodies get a single indent.
func writeBody(b *strings.Builder, body string, isScript bool) {
	if !isScript {
		fmt.Fprintf(b, "    %s\n", body)
		return
	}
	for _, line := range strings.Split(body, "\n") {
		fmt.Fprintf(b, "    %s\n", line)
	}
}

// writeOrchestrationRecipe generates the <key> aggregate recipe that chains
// lifecycle steps in sequence. The recipe stops on the first failure,
// always runs teardown, and propagates the original exit code.
func writeOrchestrationRecipe(b *strings.Builder, key string, surfaceType types.SurfaceType) {
	steps := orchestrationSteps(surfaceType)

	// user-customized marker
	b.WriteString("# user-customized\n")

	// Build the orchestration command: step1 && step2 && ... ; rc=$?; teardown; exit $rc
	// The orchestration recipe name is the bare key (e.g. "backend" for --key backend).
	prefix := key
	var chainParts []string
	for _, step := range steps {
		chainParts = append(chainParts, "just "+recipeName(key, step))
	}
	teardownName := "just " + recipeName(key, "teardown")
	chain := strings.Join(chainParts, " && ")

	// [unix] variant
	fmt.Fprintf(b, "%s [unix]:\n", prefix)
	fmt.Fprintf(b, "    %s; rc=$?; %s; exit $rc\n", chain, teardownName)

	// [windows] variant
	fmt.Fprintf(b, "%s [windows]:\n", prefix)
	fmt.Fprintf(b, "    %s; rc=$?; %s; exit $rc\n", chain, teardownName)
}

// orchestrationSteps returns the lifecycle step names in execution order,
// excluding teardown (which is handled separately in the orchestration recipe).
func orchestrationSteps(surfaceType types.SurfaceType) []string {
	switch surfaceType {
	case types.SurfaceMobile:
		return []string{"test-setup", "dev", "probe", "test"}
	case types.SurfaceAPI, types.SurfaceWeb:
		return []string{"dev", "probe", "test"}
	default:
		return nil
	}
}

// ValidateArgs checks the --type and --key combination.
// Returns an error for:
//   - unknown surface type (not in the 5 known types)
//   - scalar surface (cli, tui) with --key set
//   - named surface (api, web, mobile) without --key
//   - surface type not yet supported for recipe generation
func ValidateArgs(surfaceType types.SurfaceType, key string) error {
	// Check against all known surface types first
	if !types.AllSurfaceTypesSet()[surfaceType] {
		return fmt.Errorf("unknown surface type: %q; valid types: cli, tui, api, web, mobile", surfaceType)
	}

	isNamed := needsKey(surfaceType)
	if !isNamed && key != "" {
		return fmt.Errorf("surface type %q is scalar and does not accept --key; remove --key flag", surfaceType)
	}
	if isNamed && key == "" {
		return fmt.Errorf("surface type %q is named and requires --key", surfaceType)
	}

	return nil
}

// serviceOrder defines the dependency-based startup order for service surface types.
// Index 0 starts first; teardown runs in reverse order.
var serviceOrder = map[types.SurfaceType]int{
	types.SurfaceAPI:    0,
	types.SurfaceWeb:    1,
	types.SurfaceMobile: 2,
}

// isServiceType returns true if the surface type is a service type (api, web, mobile).
func isServiceType(t types.SurfaceType) bool {
	_, ok := serviceOrder[t]
	return ok
}

// GenerateAggregate produces cross-surface aggregate recipes (install, ci, clean)
// and optionally a test-setup aggregate recipe when multiple service-type surfaces exist.
// Aggregate recipes do NOT have the "# user-customized" marker.
func GenerateAggregate(surfaces []SurfaceEntry) (string, error) {
	var b strings.Builder

	// install recipe: aggregates all <key>-install
	writeAggregateRecipe(&b, "install", aggregateCommands(surfaces, "install"))

	// ci recipe: aggregates all <key>-lint + <key>-compile + <key>-unit-test (no surface-level test)
	var ciCmds []string
	for _, s := range surfaces {
		prefix := recipeName(s.Key, "")
		for _, verb := range []string{"lint", "compile", "unit-test"} {
			ciCmds = append(ciCmds, "just "+prefix+verb)
		}
	}
	writeAggregateRecipe(&b, "ci", ciCmds)

	// clean recipe: aggregates all <key>-clean
	writeAggregateRecipe(&b, "clean", aggregateCommands(surfaces, "clean"))

	// test-setup: only when multiple service-type surfaces exist
	if hasMultipleServiceSurfaces(surfaces) {
		writeAggregateTestSetup(&b, surfaces)
	}

	return b.String(), nil
}

// writeAggregateRecipe writes a single aggregate recipe with dual-platform variants.
// Aggregate recipes do NOT carry the "# user-customized" marker.
func writeAggregateRecipe(b *strings.Builder, name string, commands []string) {
	body := strings.Join(commands, " && ")
	if body == "" {
		body = "echo \"No surfaces configured\""
	}

	fmt.Fprintf(b, "%s [unix]:\n", name)
	fmt.Fprintf(b, "    %s\n", body)
	fmt.Fprintf(b, "%s [windows]:\n", name)
	fmt.Fprintf(b, "    %s\n", body)
}

// writeAggregateTestSetup generates the test-setup aggregate recipe for multi-service
// orchestration. Services start in dependency order (api → web → mobile) and
// teardown runs in reverse. The recipe stops on first failure, always runs teardown,
// and propagates the original exit code.
func writeAggregateTestSetup(b *strings.Builder, surfaces []SurfaceEntry) {
	// Sort service surfaces by dependency order
	services := filterServiceSurfaces(surfaces)
	sortServicesByOrder(services)

	// Build startup chain: just <key>-dev for each in order
	var startupCmds []string
	for _, s := range services {
		startupCmds = append(startupCmds, "just "+recipeName(s.Key, "dev"))
	}

	// Build teardown chain: reverse order
	var teardownCmds []string
	for i := len(services) - 1; i >= 0; i-- {
		teardownCmds = append(teardownCmds, "just "+recipeName(services[i].Key, "teardown"))
	}

	startup := strings.Join(startupCmds, " && ")
	teardown := strings.Join(teardownCmds, " && ")

	fmt.Fprintf(b, "test-setup [unix]:\n")
	fmt.Fprintf(b, "    %s; rc=$?; %s; exit $rc\n", startup, teardown)
	fmt.Fprintf(b, "test-setup [windows]:\n")
	fmt.Fprintf(b, "    %s; rc=$?; %s; exit $rc\n", startup, teardown)
}

// aggregateCommands builds a list of "just <key>-<verb>" commands for all surfaces.
func aggregateCommands(surfaces []SurfaceEntry, verb string) []string {
	cmds := make([]string, 0, len(surfaces))
	for _, s := range surfaces {
		cmds = append(cmds, "just "+recipeName(s.Key, verb))
	}
	return cmds
}

// hasMultipleServiceSurfaces returns true when there are 2+ service-type surfaces.
func hasMultipleServiceSurfaces(surfaces []SurfaceEntry) bool {
	count := 0
	for _, s := range surfaces {
		if isServiceType(s.Type) {
			count++
		}
	}
	return count >= 2
}

// filterServiceSurfaces returns only service-type surfaces (api, web, mobile).
func filterServiceSurfaces(surfaces []SurfaceEntry) []SurfaceEntry {
	var result []SurfaceEntry
	for _, s := range surfaces {
		if isServiceType(s.Type) {
			result = append(result, s)
		}
	}
	return result
}

// sortServicesByOrder sorts service surfaces by dependency order (api=0, web=1, mobile=2).
func sortServicesByOrder(services []SurfaceEntry) {
	sort.Slice(services, func(i, j int) bool {
		return serviceOrder[services[i].Type] < serviceOrder[services[j].Type]
	})
}
