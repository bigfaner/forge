package task

import (
	"sort"
	"strconv"
	"strings"

	"forge-cli/pkg/task"
)

// fallbackSortPriority is assigned to task IDs that cannot be parsed,
// ensuring they sort after all valid business IDs.
const fallbackSortPriority = 99999

// naturalSortTaskIDs sorts task IDs in natural order:
// business IDs grouped by numeric prefix (1, 1.gate, 2, 2.summary, ...),
// then test pipeline IDs (T-1, T-2, ...).
// Within the same numeric prefix, pure numeric comes before compound.
func naturalSortTaskIDs(ids []string) []string {
	sorted := make([]string, len(ids))
	copy(sorted, ids)

	sort.SliceStable(sorted, func(i, j int) bool {
		ki := sortKey(sorted[i])
		kj := sortKey(sorted[j])

		// Primary: T-prefixed IDs sort after all business IDs
		if ki.isTestPipeline != kj.isTestPipeline {
			return !ki.isTestPipeline
		}

		// Secondary: numeric prefix (groups 1, 1.gate together before 2)
		if ki.numPrefix != kj.numPrefix {
			return ki.numPrefix < kj.numPrefix
		}

		// Tertiary: same prefix — pure numeric before compound (1 < 1.gate)
		if ki.isPureNumeric != kj.isPureNumeric {
			return ki.isPureNumeric
		}

		// Quaternary: string comparison for compound suffixes
		return sorted[i] < sorted[j]
	})

	return sorted
}

type idSortKey struct {
	isTestPipeline bool
	isPureNumeric  bool
	numPrefix      int
}

func sortKey(id string) idSortKey {
	// Test pipeline IDs: T-1, T-2, etc.
	if numStr, ok := strings.CutPrefix(id, task.IDPrefixTestPipeline); ok {
		num, _ := strconv.Atoi(numStr)
		return idSortKey{isTestPipeline: true, numPrefix: num}
	}

	// Try pure numeric: "1", "2", "10"
	if num, err := strconv.Atoi(id); err == nil {
		return idSortKey{isPureNumeric: true, numPrefix: num}
	}

	// Compound IDs: "1.gate", "1.summary", "1.1" etc.
	dotIdx := strings.Index(id, ".")
	if dotIdx > 0 {
		prefix := id[:dotIdx]
		if num, err := strconv.Atoi(prefix); err == nil {
			return idSortKey{numPrefix: num}
		}
	}

	// Fallback: high numPrefix to sort last among business IDs
	return idSortKey{numPrefix: fallbackSortPriority}
}
