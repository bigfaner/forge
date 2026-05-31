// Package qualitygate provides quality gate command functionality.
// Package qualitygate provides quality gate command functionality.
package qualitygate

import "time"

const (
	maxProbeRetries      = 3
	probeRetryInterval   = 5 * time.Second
	conciseErrorMaxLines = 5
	maxSourceFiles       = 10
)
