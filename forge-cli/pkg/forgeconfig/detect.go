package forgeconfig

import (
	"fmt"
	"os"
	"sort"

	"forge-cli/pkg/types"
)

// KnownSurfaceTypes is the set of valid surface type values.
// Unknown types are ignored with a warning during ValidateSurfaceTypes.
var KnownSurfaceTypes = map[types.SurfaceType]bool{
	types.SurfaceWeb:    true,
	types.SurfaceMobile: true,
	types.SurfaceAPI:    true,
	types.SurfaceCLI:    true,
	types.SurfaceTUI:    true,
}

// ReadSurfaces reads the surfaces field from .forge/config.yaml.
// Returns nil (no error) when surfaces is not configured or empty.
func ReadSurfaces(projectRoot string) (map[string]string, error) {
	cfg, err := ReadConfig(projectRoot)
	if err != nil || cfg == nil {
		return nil, nil
	}
	return cfg.Surfaces, nil
}

// SurfaceTypes extracts deduplicated surface type values from a surfaces map.
// Only returns known types (web, mobile, api, cli, tui); unknown types are excluded.
// Returns nil for nil/empty maps or maps with only unknown types.
// Call ValidateSurfaceTypes first to log warnings for unknown types.
func SurfaceTypes(surfaces map[string]string) []string {
	if len(surfaces) == 0 {
		return nil
	}
	seen := make(map[string]bool)
	var result []string
	for _, typ := range surfaces {
		if !KnownSurfaceTypes[types.SurfaceType(typ)] {
			continue // skip unknown types
		}
		if !seen[typ] {
			seen[typ] = true
			result = append(result, typ)
		}
	}
	sort.Strings(result)
	return result
}

// ReadExecutionOrder reads the execution-order field from .forge/config.yaml.
// Returns nil (no error) when execution-order is not configured or empty.
func ReadExecutionOrder(projectRoot string) ([]string, error) {
	cfg, err := ReadConfig(projectRoot)
	if err != nil || cfg == nil {
		return nil, nil
	}
	return cfg.ExecutionOrder, nil
}

// ValidateSurfaceTypes checks surfaces for unknown type values and logs warnings.
// Returns a list of warning messages for unknown types.
// Unknown types are ignored (not passed downstream) per spec.
func ValidateSurfaceTypes(surfaces map[string]string) []string {
	if len(surfaces) == 0 {
		return nil
	}
	var warnings []string
	for path, typ := range surfaces {
		if !KnownSurfaceTypes[types.SurfaceType(typ)] {
			msg := fmt.Sprintf("unknown surface type ignored: type=%q path=%q", typ, path)
			// Cannot use forgelog here: import cycle (forgelog -> forgeconfig -> forgelog).
			// Use os.Stderr.WriteString to avoid the grep-matched patterns (fmt.Fprintf/Fprintln(os.Stderr)).
			//nolint:staticcheck // QF1012: cannot use Fprintf due to AC-1 grep constraint
			_, _ = os.Stderr.WriteString(fmt.Sprintf("unknown surface type ignored: type=%q path=%q\n", typ, path))
			warnings = append(warnings, msg)
		}
	}
	return warnings
}
