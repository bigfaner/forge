package journey

import (
	"strings"
	"testing"
)

// --- Test: TestOutputDir ---

func TestTestOutputDir(t *testing.T) {
	got := TestOutputDir("/project", "task-lifecycle")
	if !strings.Contains(got, "tests") || !strings.Contains(got, "task-lifecycle") {
		t.Fatalf("expected path containing tests/task-lifecycle, got %q", got)
	}
}

// --- Test: SmokeTestName ---

func TestSmokeTestName(t *testing.T) {
	tests := []struct {
		journey  string
		expected string
	}{
		{"task-lifecycle", "TestJourneyTaskLifecycleSmoke"},
		{"session-diagnostics", "TestJourneySessionDiagnosticsSmoke"},
		{"simple", "TestJourneySimpleSmoke"},
	}

	for _, tt := range tests {
		t.Run(tt.journey, func(t *testing.T) {
			got := SmokeTestName(tt.journey)
			if got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

// --- Test: sanitizeName ---

func TestSanitizeName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"forge feature my-feature", "forge_feature_my_feature"},
		{"forge task claim", "forge_task_claim"},
		{"camelCase", "camelcase"},
		{"with/slash", "with_slash"},
		{"with.dot", "with_dot"},
		{"with:colon", "with_colon"},
		{"extra  spaces", "extra_spaces"},
		{"trailing_", "trailing"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := sanitizeName(tt.input)
			if got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}
