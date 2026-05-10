// Package e2eprobe provides end-to-end server health probing.
package e2eprobe

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"task-cli/pkg/just"
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

// ProbeServers reads tests/e2e/config.yaml and probes baseUrl/apiBaseUrl.
// Returns true if all configured endpoints respond, or if no config exists.
func ProbeServers(projectRoot string) bool {
	configPath := filepath.Join(projectRoot, "tests", "e2e", "config.yaml")
	if !just.FileExists(configPath) {
		fmt.Fprintln(os.Stderr, "  No tests/e2e/config.yaml found; skipping health check")
		return true
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  WARNING: cannot read config.yaml: %v\n", err)
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
		return true
	}

	probeTimeout := 5 * time.Second
	for _, ep := range endpoints {
		if !ProbeEndpoint(ep, probeTimeout) {
			fmt.Fprintf(os.Stderr, "  Endpoint unreachable: %s\n", ep)
			return false
		}
		fmt.Fprintf(os.Stderr, "  Endpoint OK: %s\n", ep)
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
