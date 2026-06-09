package scaffold

import (
	"fmt"
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

	return b.String(), nil
}

// writeRecipe writes a single recipe with dual-platform variants.
// All lifecycle and quality recipes are marked with "# user-customized".
func writeRecipe(b *strings.Builder, key string, r RecipeSpec) {
	name := recipeName(key, r.Name)

	// user-customized marker for lifecycle and quality recipes
	b.WriteString("# user-customized\n")

	// [unix] variant
	fmt.Fprintf(b, "%s [unix]:\n", name)
	fmt.Fprintf(b, "    %s\n", r.UnixBody)

	// [windows] variant
	fmt.Fprintf(b, "%s [windows]:\n", name)
	fmt.Fprintf(b, "    %s\n", r.WindowsBody)
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

	// Check if recipe generation is supported for this type
	if _, ok := surfaceSpecs[surfaceType]; !ok {
		return fmt.Errorf("surface type %q is not yet supported for scaffold generation", surfaceType)
	}

	return nil
}
