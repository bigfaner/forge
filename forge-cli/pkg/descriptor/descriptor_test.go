package descriptor

import (
	"strings"
	"testing"
)

// --- FactEntry fixtures ---

func sampleFacts() []FactEntry {
	return []FactEntry{
		{
			Key:    "CLI_FEATURE_CREATE_OUTPUT",
			Value:  "Feature my-feature created successfully",
			Source: "internal/cmd/feature.go:45",
		},
		{
			Key:    "CLI_TASK_CLAIM_OUTPUT",
			Value:  "claimed task <task_id>",
			Source: "internal/cmd/claim.go:42",
		},
		{
			Key:    "CLI_TASK_SUBMIT_OUTPUT",
			Value:  "submitted task <task_id> with result: success",
			Source: "internal/cmd/submit.go:58",
		},
		{
			Key:    "CLI_TASK_CLAIM_NO_TASKS",
			Value:  "no tasks available for claiming",
			Source: "internal/cmd/claim.go:50",
		},
	}
}

// --- Test: BuildRegex resolves semantic descriptor via Fact Table ---

func TestBuildRegex(t *testing.T) {
	t.Run("resolves feature creation output", func(t *testing.T) {
		result := BuildRegex("success confirmation containing feature-slug", sampleFacts())
		if result.Unresolved {
			t.Fatal("expected resolved result")
		}
		if result.Pattern == "" {
			t.Fatal("expected non-empty pattern")
		}
		// Pattern should contain the escaped literal text from the fact
		if !strings.Contains(result.Pattern, "Feature") {
			t.Fatalf("pattern should contain 'Feature', got %q", result.Pattern)
		}
		if !strings.Contains(result.Pattern, "created successfully") {
			t.Fatalf("pattern should contain 'created successfully', got %q", result.Pattern)
		}
	})

	t.Run("resolves task claim output with placeholder", func(t *testing.T) {
		result := BuildRegex("claimed task with identifier", sampleFacts())
		if result.Unresolved {
			t.Fatal("expected resolved result")
		}
		// Should have a named capture group for task_id
		if len(result.Captures) == 0 {
			t.Fatalf("expected captures, got pattern %q", result.Pattern)
		}
		if _, ok := result.Captures["task_id"]; !ok {
			t.Fatalf("expected task_id capture, got captures %v", result.Captures)
		}
	})

	t.Run("resolves task submit output with placeholder", func(t *testing.T) {
		result := BuildRegex("submitted task with result", sampleFacts())
		if result.Unresolved {
			t.Fatal("expected resolved result")
		}
		if _, ok := result.Captures["task_id"]; !ok {
			t.Fatalf("expected task_id capture, got captures %v", result.Captures)
		}
	})

	t.Run("returns unresolved when no fact matches", func(t *testing.T) {
		result := BuildRegex("unknown semantic descriptor about unicorns", sampleFacts())
		if !result.Unresolved {
			t.Fatal("expected unresolved result for non-matching descriptor")
		}
		if result.Placeholder != "unknown semantic descriptor about unicorns" {
			t.Fatalf("expected original descriptor as placeholder, got %q", result.Placeholder)
		}
	})

	t.Run("returns unresolved for empty facts", func(t *testing.T) {
		result := BuildRegex("success confirmation", nil)
		if !result.Unresolved {
			t.Fatal("expected unresolved for empty facts")
		}
	})

	t.Run("no tasks available matches correctly", func(t *testing.T) {
		result := BuildRegex("no tasks available error message", sampleFacts())
		if result.Unresolved {
			t.Fatal("expected resolved result")
		}
		if !strings.Contains(result.Pattern, "no tasks available") {
			t.Fatalf("pattern should contain 'no tasks available', got %q", result.Pattern)
		}
	})
}

// --- Test: valueToRegex ---

func TestValueToRegex(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		contains string // pattern must contain this substring
		exact    string // if set, must match exactly
	}{
		{
			name:     "plain text is escaped",
			value:    "Feature my-feature created successfully",
			contains: `Feature my-feature created successfully`,
		},
		{
			name:     "single placeholder becomes capture group",
			value:    "claimed task <task_id>",
			contains: `(?P<task_id>`,
		},
		{
			name:     "hyphenated placeholder name normalized",
			value:    "created feature <feature-slug> successfully",
			contains: `(?P<feature_slug>`,
		},
		{
			name:     "multiple placeholders",
			value:    "task <task_id> in feature <feature-slug>",
			contains: `(?P<task_id>`,
		},
		{
			name:     "text with regex metacharacters is escaped",
			value:    "output: [done]. (100%)",
			contains: `\[done\]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern := valueToRegex(tt.value)
			if tt.exact != "" && pattern != tt.exact {
				t.Fatalf("expected %q, got %q", tt.exact, pattern)
			}
			if tt.contains != "" && !strings.Contains(pattern, tt.contains) {
				t.Fatalf("expected pattern to contain %q, got %q", tt.contains, pattern)
			}
		})
	}
}

// --- Test: extractCaptures ---

func TestExtractCaptures(t *testing.T) {
	t.Run("no placeholders returns empty map", func(t *testing.T) {
		captures := extractCaptures("plain text no placeholders")
		if len(captures) != 0 {
			t.Fatalf("expected empty captures, got %v", captures)
		}
	})

	t.Run("single placeholder", func(t *testing.T) {
		captures := extractCaptures("claimed task <task_id>")
		if len(captures) != 1 {
			t.Fatalf("expected 1 capture, got %d", len(captures))
		}
		if captures["task_id"] != 1 {
			t.Fatalf("expected task_id at index 1, got %d", captures["task_id"])
		}
	})

	t.Run("multiple placeholders", func(t *testing.T) {
		captures := extractCaptures("<feature-slug>/<task_id>/status")
		if len(captures) != 2 {
			t.Fatalf("expected 2 captures, got %d", len(captures))
		}
		if captures["feature_slug"] != 1 {
			t.Fatalf("expected feature_slug at index 1, got %d", captures["feature_slug"])
		}
		if captures["task_id"] != 2 {
			t.Fatalf("expected task_id at index 2, got %d", captures["task_id"])
		}
	})

	t.Run("hyphenated placeholder names normalized", func(t *testing.T) {
		captures := extractCaptures("result: <my-long-field-name>")
		if _, ok := captures["my_long_field_name"]; !ok {
			t.Fatalf("expected my_long_field_name in captures, got %v", captures)
		}
	})
}

// --- Test: extractKeywords ---

func TestExtractKeywords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "extracts significant words",
			input:    "success confirmation containing feature-slug",
			expected: []string{"success", "confirmation", "containing", "feature-slug"},
		},
		{
			name:     "filters stop words",
			input:    "the output contains a success message",
			expected: []string{"output", "contains", "success", "message"},
		},
		{
			name:     "filters short words",
			input:    "no tasks available for claiming",
			expected: []string{"tasks", "available", "claiming"},
		},
		{
			name:     "empty input",
			input:    "",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keywords := extractKeywords(tt.input)
			if len(keywords) != len(tt.expected) {
				t.Fatalf("expected %v, got %v", tt.expected, keywords)
			}
			for i, kw := range keywords {
				if kw != tt.expected[i] {
					t.Fatalf("expected keyword %q at index %d, got %q", tt.expected[i], i, kw)
				}
			}
		})
	}
}

// --- Test: BuildAssertionRegex ---

func TestBuildAssertionRegex(t *testing.T) {
	t.Run("resolved descriptor returns regex", func(t *testing.T) {
		pattern := BuildAssertionRegex("success confirmation containing feature-slug", sampleFacts())
		if pattern == "" {
			t.Fatal("expected non-empty pattern")
		}
		// Should not be escaped literal of the original descriptor
		if pattern == "success confirmation containing feature\\-slug" {
			t.Fatal("should return Fact Table based regex, not escaped descriptor")
		}
	})

	t.Run("unresolved descriptor returns escaped literal", func(t *testing.T) {
		pattern := BuildAssertionRegex("unknown descriptor here", sampleFacts())
		if pattern == "" {
			t.Fatal("expected non-empty pattern")
		}
	})
}

// --- Test: Sensitive field handling ---

func TestBuildSensitiveFieldPlaceholder(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		expected string
	}{
		{"token", "token", "<from-env>"},
		{"api_key", "api_key", "<from-env>"},
		{"secret", "secret", "<from-env>"},
		{"password", "password", "<from-env>"},
		{"normal field", "username", "username"},
		{"mixed case", "Token", "<from-env>"},
		{"access_token", "access_token", "<from-env>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildSensitiveFieldPlaceholder(tt.field)
			if got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestIsSensitiveField(t *testing.T) {
	tests := []struct {
		name      string
		field     string
		sensitive bool
	}{
		{"token", "token", true},
		{"api_key", "api_key", true},
		{"secret", "secret", true},
		{"password", "password", true},
		{"credential", "credential", true},
		{"username", "username", false},
		{"feature_slug", "feature_slug", false},
		{"task_id", "task_id", false},
		{"mixed case token", "AuthToken", true},
		{"mixed case key", "APIKey", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSensitiveField(tt.field)
			if got != tt.sensitive {
				t.Fatalf("IsSensitiveField(%q) = %v, want %v", tt.field, got, tt.sensitive)
			}
		})
	}
}
