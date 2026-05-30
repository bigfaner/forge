package cmd

import "time"

// Retry and tuning parameters for quality gate operations.
const (
	// maxProbeRetries is the maximum number of probe attempts during surface lifecycle.
	maxProbeRetries = 3

	// probeRetryInterval is the delay between probe retries.
	probeRetryInterval = 5 * time.Second

	// conciseErrorMaxLines is the maximum number of lines extracted for concise error output.
	conciseErrorMaxLines = 5

	// maxSourceFiles caps the number of source files extracted from error output.
	maxSourceFiles = 10
)
