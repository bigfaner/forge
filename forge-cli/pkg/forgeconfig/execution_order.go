package forgeconfig

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"forge-cli/pkg/types"
)

// surfaceKeyPattern defines the valid surface-key format after normalization.
// Must start with a lowercase letter, followed by lowercase letters, digits, or hyphens.
var surfaceKeyPattern = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

// NormalizeSurfaceKey normalizes a raw surface key from YAML config and returns the normalized form.
// Normalization rules (aligned with surface-aware-justfile proposal):
//   - Convert to lowercase
//   - Replace spaces and special characters (non-alphanumeric, non-hyphen) with hyphens
//
// After normalization, the key must match [a-z][a-z0-9-]*.
// Returns the normalized key or an error if it is invalid.
func NormalizeSurfaceKey(raw string) (string, error) {
	normalized := strings.ToLower(raw)

	// Replace any character that is not [a-z0-9-] with a hyphen
	var b strings.Builder
	for _, r := range normalized {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		} else {
			b.WriteRune('-')
		}
	}
	normalized = b.String()

	// Validate: must start with lowercase letter, rest must be [a-z0-9-]*
	if !surfaceKeyPattern.MatchString(normalized) {
		return normalized, fmt.Errorf("invalid surface-key %q: normalized to %q which does not match [a-z][a-z0-9-]*", raw, normalized)
	}

	return normalized, nil
}

// normalizeSurfaceKeyValue returns the normalized form of a surface key.
// This is used during config load to rewrite map keys.
// It reuses the same normalization logic as NormalizeSurfaceKey but returns only the string.
// Validation is deferred to ValidateSurfaceKeys called in validateSurfacesConfig.
func normalizeSurfaceKeyValue(raw string) string {
	normalized, _ := NormalizeSurfaceKey(raw)
	// NormalizeSurfaceKey returns the normalized string even on error
	// (the error is about the key not matching the pattern, not about normalization failure)
	return normalized
}

// ValidateSurfaceKeys checks all surface keys for validity.
// The "." key (scalar form) is exempt from validation.
// Returns an error listing all invalid keys.
func ValidateSurfaceKeys(surfaces map[string]string) error {
	var invalid []string
	for key := range surfaces {
		if key == "." {
			continue // scalar form, exempt
		}
		if !surfaceKeyPattern.MatchString(key) {
			invalid = append(invalid, key)
		}
	}
	if len(invalid) > 0 {
		sort.Strings(invalid)
		return fmt.Errorf("invalid surface-key(s): %s (must match [a-z][a-z0-9-]*)", strings.Join(invalid, ", "))
	}
	return nil
}

// ValidateExecutionOrder validates the execution-order config against the surfaces map.
// Checks:
//  1. Each key in execution-order must exist in surfaces
//  2. When multiple surfaces share the same type, execution-order must be specified
//
// Returns nil if valid, or an error describing the problem.
func ValidateExecutionOrder(surfaces map[string]string, executionOrder []string) error {
	if len(surfaces) == 0 {
		return nil
	}

	// Scalar form or single surface: no ordering needed
	if len(surfaces) == 1 {
		return nil
	}

	// Validate execution-order references
	if len(executionOrder) > 0 {
		for _, key := range executionOrder {
			if _, ok := surfaces[key]; !ok {
				return fmt.Errorf("execution-order references non-existent surface-key %q", key)
			}
		}
		return nil
	}

	// No explicit execution-order: check for same-type conflicts
	typeToKeys := make(map[string][]string)
	for key, typ := range surfaces {
		typeToKeys[typ] = append(typeToKeys[typ], key)
	}

	var conflicts []string
	for typ, keys := range typeToKeys {
		if len(keys) > 1 {
			sort.Strings(keys)
			conflicts = append(conflicts, fmt.Sprintf("type %q has keys %s", typ, strings.Join(keys, ", ")))
		}
	}

	if len(conflicts) > 0 {
		sort.Strings(conflicts)
		return fmt.Errorf("same-type surfaces detected without execution-order: %s — add 'execution-order' to config to specify ordering", strings.Join(conflicts, "; "))
	}

	return nil
}

// defaultExecutionOrder defines the priority order for execution when no explicit
// execution-order is configured. api > web > cli > tui > mobile.
var defaultExecutionOrder = []types.SurfaceType{
	types.SurfaceAPI,
	types.SurfaceWeb,
	types.SurfaceCLI,
	types.SurfaceTUI,
	types.SurfaceMobile,
}

// ResolveExecutionOrder determines the execution order of surface keys.
// If executionOrder is provided, it is used as-is (already validated).
// Otherwise, surfaces are ordered by default priority (api > web > cli > tui > mobile).
// For types not in the default priority list, YAML map key order is preserved.
// Returns nil for nil/empty surfaces, scalar form, or single-surface configs.
func ResolveExecutionOrder(surfaces map[string]string, executionOrder []string) ([]string, error) {
	if len(surfaces) == 0 {
		return nil, nil
	}

	// Scalar form or single surface: no ordering needed
	if len(surfaces) == 1 {
		return nil, nil
	}

	// If explicit order provided, return as-is (already validated by ValidateExecutionOrder)
	if len(executionOrder) > 0 {
		return executionOrder, nil
	}

	// Default ordering: sort by type priority, with YAML order as tiebreaker for same type
	typePriority := make(map[types.SurfaceType]int)
	for i, typ := range defaultExecutionOrder {
		typePriority[typ] = i
	}

	keys := make([]string, 0, len(surfaces))
	for key := range surfaces {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		typI := types.SurfaceType(surfaces[keys[i]])
		typJ := types.SurfaceType(surfaces[keys[j]])
		priI, hasI := typePriority[typI]
		priJ, hasJ := typePriority[typJ]

		// Both have default priority: sort by priority
		if hasI && hasJ {
			return priI < priJ
		}
		// Only one has default priority: default-priority types come first
		if hasI && !hasJ {
			return true
		}
		if !hasI && hasJ {
			return false
		}
		// Neither has default priority: preserve insertion order (stable sort)
		return false
	})

	return keys, nil
}
