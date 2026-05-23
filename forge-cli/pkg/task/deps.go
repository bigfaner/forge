package task

import "strings"

// ResolveWildcardDep resolves a dependency string against the task index.
// For wildcard dependencies (e.g. "1.x"), returns all matching business task IDs.
// For exact dependencies, returns a single-element slice containing the dep string.
// The second return value indicates whether the dep is a wildcard pattern.
func ResolveWildcardDep(index *TaskIndex, dep string) ([]string, bool) {
	if !strings.HasSuffix(dep, IDSuffixWildcard) {
		return []string{dep}, false
	}

	prefix := strings.TrimSuffix(dep, IDSuffixWildcard)
	prefixWithDot := prefix + "."

	var matches []string
	for _, t := range index.tasks {
		if strings.HasPrefix(t.ID, prefixWithDot) && IsBusinessTask(t.ID) {
			matches = append(matches, t.ID)
		}
	}
	return matches, true
}

// satisfiedStatuses are the statuses that satisfy dependency checks.
var satisfiedStatuses = map[string]bool{
	"completed": true,
	"skipped":   true,
}

// IsDepSatisfied returns true if the given status satisfies a dependency.
func IsDepSatisfied(status string) bool {
	return satisfiedStatuses[status]
}

// GetUnmetDeps returns the concrete dependency IDs that are not satisfied.
// Wildcard deps are expanded to their matching business task IDs.
// selfID is excluded from wildcard matches to prevent self-dependency.
// An exact dep that doesn't exist in the index is reported as unmet.
func GetUnmetDeps(index *TaskIndex, selfID string, deps []string) []string {
	var unmet []string
	for _, dep := range deps {
		matches, isWildcard := ResolveWildcardDep(index, dep)
		if isWildcard {
			for _, matchID := range matches {
				if matchID == selfID {
					continue
				}
				t, found := index.ByID(matchID)
				if !found || !IsDepSatisfied(t.Status) {
					unmet = append(unmet, matchID)
				}
			}
		} else {
			t, found := index.ByID(dep)
			if !found || !IsDepSatisfied(t.Status) {
				unmet = append(unmet, dep)
			}
		}
	}
	return unmet
}
