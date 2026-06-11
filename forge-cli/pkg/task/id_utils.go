package task

import (
	"strconv"
	"strings"
)

// GetTaskPhase returns the phase number from a task ID.
// For example, "1.2" returns 1, "3.gate" returns 3.
// Returns -1 if the ID has no numeric prefix.
func GetTaskPhase(id string) int {
	parts := strings.Split(id, ".")
	if len(parts) > 0 {
		phase, err := strconv.Atoi(parts[0])
		if err == nil {
			return phase
		}
	}
	return -1
}

// ParseSegment returns the numeric value of a segment and whether it's numeric.
// Numeric segments (e.g., "1", "12") return their int value with true.
// Alphabetic segments (e.g., "summary", "gate") return a lexicographic rank with false.
func ParseSegment(parts []string, i int) (int, bool) {
	if i >= len(parts) {
		return -1, true // missing segments sort before everything
	}
	if n, err := strconv.Atoi(parts[i]); err == nil {
		return n, true
	}
	// Alphabetic segments: sort after all numeric, with deterministic order
	switch parts[i] {
	case "gate":
		return 1, false
	case "summary":
		return 2, false
	default:
		return 0, false
	}
}

// CompareVersionIDs compares two task version IDs for ordering.
// Returns true if a should come before b in sorted order.
func CompareVersionIDs(a, b string) bool {
	partsA := strings.Split(a, ".")
	partsB := strings.Split(b, ".")
	maxLen := len(partsA)
	if len(partsB) > maxLen {
		maxLen = len(partsB)
	}
	for i := 0; i < maxLen; i++ {
		na, aIsNum := ParseSegment(partsA, i)
		nb, bIsNum := ParseSegment(partsB, i)
		if aIsNum != bIsNum {
			return aIsNum // numeric sorts before alphabetic
		}
		if na != nb {
			return na < nb
		}
	}
	return false
}
