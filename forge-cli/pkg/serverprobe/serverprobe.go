// Package serverprobe provides server health probing for functional and e2e tests.
package serverprobe

import (
	"net/http"
	"os"
	"strings"
	"time"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/just"

	"forge-cli/pkg/forgelog"
)

// ProbeEndpoint checks if an HTTP endpoint responds with status < 500.
func ProbeEndpoint(url string, timeout time.Duration) bool {
	client := http.Client{Timeout: timeout}
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer func() { _ = resp.Body.Close() }()
	return resp.StatusCode < 500
}

// ProbeServers reads tests/config.yaml and probes baseUrl/apiBaseUrl.
// Returns true if all configured endpoints respond, or if no config exists.
// path is the health check path appended to each URL (defaults to "/health").
func ProbeServers(projectRoot, path string) bool {
	if path == "" {
		path = defaultHealthPath
	}

	configPath := feature.GetTestConfigPath(projectRoot)
	if !just.FileExists(configPath) {
		forgelog.Info("OK: CLI-only project\n")
		return true
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		forgelog.Warn("  WARNING: cannot read config.yaml: %v\n", err)
		return true
	}

	baseURL := ExtractYAMLStringField(data, "baseUrl")
	apiBaseURL := ExtractYAMLStringField(data, "apiBaseUrl")

	endpoints := []string{}
	if baseURL != "" {
		endpoints = append(endpoints, baseURL)
	}
	if apiBaseURL != "" {
		endpoints = append(endpoints, apiBaseURL)
	}
	if len(endpoints) == 0 {
		forgelog.Info("OK: CLI-only project\n")
		return true
	}

	probeTimeout := defaultProbeTimeout
	for _, ep := range endpoints {
		probeURL := strings.TrimRight(ep, "/") + path
		if !ProbeEndpoint(probeURL, probeTimeout) {
			forgelog.Warn("FAIL: %s not responding\n", probeURL)
			return false
		}
		forgelog.Info("OK: %s\n", probeURL)
	}
	return true
}

// ExtractYAMLStringField extracts a top-level string field from simple YAML.
func ExtractYAMLStringField(data []byte, field string) string {
	for line := range strings.SplitSeq(string(data), "\n") {
		line = strings.TrimSpace(line)
		if after, ok := strings.CutPrefix(line, field+":"); ok {
			val := strings.TrimSpace(after)
			val = strings.Trim(val, `'"`)
			return val
		}
	}
	return ""
}
