package serverprobe

import "time"

// defaultProbeTimeout is the maximum time to wait for an HTTP probe response.
const defaultProbeTimeout = 5 * time.Second
