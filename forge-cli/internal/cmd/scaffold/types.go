// Package scaffold implements the `forge justfile scaffold` command.
//
// It generates just recipes for a given surface type, outputting
// placeholder-templated justfile code to stdout. Only the recipe
// skeleton is produced — the agent fills project-specific values.
package scaffold

import "forge-cli/pkg/types"

// RecipeKind categorizes recipes into lifecycle and quality groups.
type RecipeKind string

const (
	// KindLifecycle marks recipes that manage surface lifecycle (test, teardown, dev, probe).
	KindLifecycle RecipeKind = "lifecycle"
	// KindQuality marks recipes for code quality checks (compile, fmt, lint, unit-test).
	KindQuality RecipeKind = "quality"
)

// RecipeSpec describes a single recipe to generate for a surface type.
type RecipeSpec struct {
	// Name is the recipe verb (e.g. "test", "teardown", "compile").
	Name string
	// Kind is the recipe category: lifecycle or quality.
	Kind RecipeKind
	// UnixBody is the [unix] recipe body template with <<PLACEHOLDER>> slots.
	UnixBody string
	// WindowsBody is the [windows] recipe body template with <<PLACEHOLDER>> slots.
	WindowsBody string
}

// SurfaceSpec defines the recipe set for a surface type.
type SurfaceSpec struct {
	// Type is the surface type this spec covers.
	Type types.SurfaceType
	// RequiresKey is true when this surface type must have a --key (named surface).
	// When false, passing --key is an error (scalar surface).
	RequiresKey bool
	// Recipes is the ordered list of recipe specs for this surface.
	Recipes []RecipeSpec
}

// needsKey returns true if the surface type requires a named key (api, web, mobile).
func needsKey(surfaceType types.SurfaceType) bool {
	switch surfaceType {
	case types.SurfaceAPI, types.SurfaceWeb, types.SurfaceMobile:
		return true
	default:
		return false
	}
}

// simpleRecipes returns the shared recipe set for cli and tui surface types.
// These are simple surfaces: no dev/probe lifecycle, only test + teardown + quality.
func simpleRecipes() []RecipeSpec {
	return []RecipeSpec{
		{
			Name:        "test",
			Kind:        KindLifecycle,
			UnixBody:    `cd <<PROJECT_DIR>> && <<TEST_CMD>>`,
			WindowsBody: `cd <<PROJECT_DIR>> && <<TEST_CMD>>`,
		},
		{
			Name:        "teardown",
			Kind:        KindLifecycle,
			UnixBody:    `echo "No teardown needed for CLI/TUI surface"`,
			WindowsBody: `echo No teardown needed for CLI/TUI surface`,
		},
		{
			Name:        "compile",
			Kind:        KindQuality,
			UnixBody:    `cd <<PROJECT_DIR>> && <<COMPILE_CMD>>`,
			WindowsBody: `cd <<PROJECT_DIR>> && <<COMPILE_CMD>>`,
		},
		{
			Name:        "fmt",
			Kind:        KindQuality,
			UnixBody:    `cd <<PROJECT_DIR>> && <<FMT_CMD>>`,
			WindowsBody: `cd <<PROJECT_DIR>> && <<FMT_CMD>>`,
		},
		{
			Name:        "lint",
			Kind:        KindQuality,
			UnixBody:    `cd <<PROJECT_DIR>> && <<LINT_CMD>>`,
			WindowsBody: `cd <<PROJECT_DIR>> && <<LINT_CMD>>`,
		},
		{
			Name:        "unit-test",
			Kind:        KindQuality,
			UnixBody:    `cd <<PROJECT_DIR>> && <<UNIT_TEST_CMD>>`,
			WindowsBody: `cd <<PROJECT_DIR>> && <<UNIT_TEST_CMD>>`,
		},
	}
}

// surfaceSpecs maps each surface type to its spec.
// cli and tui share simpleRecipes; api/web/mobile will have richer specs in future tasks.
var surfaceSpecs = map[types.SurfaceType]SurfaceSpec{
	types.SurfaceCLI: {
		Type:    types.SurfaceCLI,
		Recipes: simpleRecipes(),
	},
	types.SurfaceTUI: {
		Type:    types.SurfaceTUI,
		Recipes: simpleRecipes(),
	},
}
