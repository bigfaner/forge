package forgeconfig

import (
	"fmt"
	"log/slog"
)

// KnownSurfaceTypes is the set of valid surface type values.
// Unknown types are ignored with a warning during ValidateSurfaceTypes.
var KnownSurfaceTypes = map[string]bool{
	"web":    true,
	"mobile": true,
	"api":    true,
	"cli":    true,
	"tui":    true,
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
	var types []string
	for _, typ := range surfaces {
		if !KnownSurfaceTypes[typ] {
			continue // skip unknown types
		}
		if !seen[typ] {
			seen[typ] = true
			types = append(types, typ)
		}
	}
	return types
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
		if !KnownSurfaceTypes[typ] {
			msg := fmt.Sprintf("unknown surface type ignored: type=%q path=%q", typ, path)
			slog.Warn("unknown surface type ignored", "type", typ, "path", path)
			warnings = append(warnings, msg)
		}
	}
	return warnings
}
