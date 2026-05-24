package forgeconfig

import (
	"fmt"
	"strings"
)

// NormalizePath applies path normalization rules:
//  1. Strip leading "./"
//  2. Strip trailing "/"
//  3. Convert "\" to "/" (Windows compatibility)
//  4. Reject paths containing ".." segments (security: no path traversal)
//  5. No symlink resolution — purely string-based normalization
func NormalizePath(path string) (string, error) {
	// Rule 3: convert backslashes to forward slashes
	path = strings.ReplaceAll(path, `\`, "/")

	// Rule 4: reject ".." segments
	// Check after backslash conversion so `foo\..` is caught
	for _, seg := range strings.Split(path, "/") {
		if seg == ".." {
			return "", fmt.Errorf("path contains '..'")
		}
	}

	// Rule 1: strip leading "./"
	path = strings.TrimPrefix(path, "./")

	// Rule 2: strip trailing "/"
	path = strings.TrimRight(path, "/")

	return path, nil
}

// MatchSurface finds the surface type for a query path using segment prefix matching.
//
// Matching rules:
//   - Scalar form (single key "."): any path returns the value directly, no matching
//   - Map form: segment prefix matching — longest segment match wins
//   - No partial segment match: "frontend" does NOT match "frontend-new"
//   - No match: returns error with manual config hint
//
// The query path is normalized before matching.
func MatchSurface(surfaces map[string]string, query string) (string, error) {
	if len(surfaces) == 0 {
		return "", fmt.Errorf("no surface configured; run `forge init` to configure surfaces")
	}

	// Scalar form: single key "." means any path returns the value
	if len(surfaces) == 1 {
		if v, ok := surfaces["."]; ok {
			return v, nil
		}
	}

	// Normalize the query path
	normalized, err := NormalizePath(query)
	if err != nil {
		return "", err
	}

	querySegments := strings.Split(normalized, "/")

	// Find the longest segment-prefix match
	bestMatch := ""
	bestDepth := -1

	for configPath, surfaceType := range surfaces {
		configSegments := strings.Split(configPath, "/")

		if !isSegmentPrefix(configSegments, querySegments) {
			continue
		}

		if len(configSegments) > bestDepth {
			bestDepth = len(configSegments)
			bestMatch = surfaceType
		}
	}

	if bestDepth < 0 {
		return "", fmt.Errorf(
			"no surface found for path %q; run `forge init` to configure surfaces",
			query,
		)
	}

	return bestMatch, nil
}

// isSegmentPrefix checks whether configSegs is a segment-level prefix of querySegs.
// "frontend" matches "frontend/src" but NOT "frontend-new".
func isSegmentPrefix(configSegs, querySegs []string) bool {
	if len(configSegs) > len(querySegs) {
		return false
	}
	for i, seg := range configSegs {
		if seg != querySegs[i] {
			return false
		}
	}
	return true
}
