package task

import (
	"os"
	"strings"
	"testing"
)

// --- Task 1.2b: SurfaceKey/SurfaceType frontmatter tests ---

func TestParseFrontmatter_SurfaceFields(t *testing.T) {
	t.Run("surface-key and surface-type present", func(t *testing.T) {
		input := `---
id: "1.1"
title: "Frontend Task"
type: "coding.feature"
surface-key: admin-panel
surface-type: web
---

Body here`
		fm, body, err := ParseFrontmatter([]byte(input))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if fm.SurfaceKey != "admin-panel" {
			t.Errorf("SurfaceKey = %q, want %q", fm.SurfaceKey, "admin-panel")
		}
		if fm.SurfaceType != "web" {
			t.Errorf("SurfaceType = %q, want %q", fm.SurfaceType, "web")
		}
		if fm.ID != "1.1" {
			t.Errorf("ID = %q, want %q", fm.ID, "1.1")
		}
		_ = body // body not under test
	})

	t.Run("surface fields absent — empty values", func(t *testing.T) {
		input := `---
id: "1.2"
title: "No Surface Task"
type: "coding.feature"
---

Body`
		fm, _, err := ParseFrontmatter([]byte(input))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if fm.SurfaceKey != "" {
			t.Errorf("SurfaceKey = %q, want empty when absent", fm.SurfaceKey)
		}
		if fm.SurfaceType != "" {
			t.Errorf("SurfaceType = %q, want empty when absent", fm.SurfaceType)
		}
	})

	t.Run("old frontmatter with scope but no surface-key — no error", func(t *testing.T) {
		input := `---
id: "1.3"
title: "Legacy Task"
scope: all
type: "coding.feature"
---

Body`
		fm, _, err := ParseFrontmatter([]byte(input))
		if err != nil {
			t.Fatalf("should not error on legacy scope field: %v", err)
		}
		if fm.Scope != "all" {
			t.Errorf("Scope = %q, want %q", fm.Scope, "all")
		}
		if fm.SurfaceKey != "" {
			t.Errorf("SurfaceKey should be empty, got %q", fm.SurfaceKey)
		}
	})
}

func TestWriteFrontmatter_SurfaceFields(t *testing.T) {
	t.Run("surface-key and surface-type serialized", func(t *testing.T) {
		dir := t.TempDir()
		path := dir + "/test-task.md"
		fm := FrontmatterData{
			ID:          "1.1",
			Title:       "Frontend Task",
			SurfaceKey:  "admin-panel",
			SurfaceType: "web",
		}
		body := []byte("\n# Task body\n")

		if err := WriteFrontmatter(path, fm, body); err != nil {
			t.Fatalf("WriteFrontmatter failed: %v", err)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("ReadFile failed: %v", err)
		}
		content := string(data)
		if !strings.Contains(content, "surface-key: admin-panel") {
			t.Errorf("expected 'surface-key: admin-panel' in output, got:\n%s", content)
		}
		if !strings.Contains(content, "surface-type: web") {
			t.Errorf("expected 'surface-type: web' in output, got:\n%s", content)
		}
	})

	t.Run("empty surface fields omitted", func(t *testing.T) {
		dir := t.TempDir()
		path := dir + "/test-task.md"
		fm := FrontmatterData{
			ID:    "1.2",
			Title: "No Surface Task",
		}
		body := []byte("\n# Task body\n")

		if err := WriteFrontmatter(path, fm, body); err != nil {
			t.Fatalf("WriteFrontmatter failed: %v", err)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("ReadFile failed: %v", err)
		}
		content := string(data)
		if strings.Contains(content, "surface-key:") {
			t.Errorf("empty surface-key should be omitted, got:\n%s", content)
		}
		if strings.Contains(content, "surface-type:") {
			t.Errorf("empty surface-type should be omitted, got:\n%s", content)
		}
	})
}

func TestFrontmatter_SurfaceKeyRoundTrip(t *testing.T) {
	t.Run("write then read preserves surface fields", func(t *testing.T) {
		dir := t.TempDir()
		path := dir + "/roundtrip.md"

		original := FrontmatterData{
			ID:            "2.1",
			Title:         "API Task",
			Priority:      "P1",
			Type:          "coding.feature",
			SurfaceKey:    "payment-service",
			SurfaceType:   "api",
			EstimatedTime: "2h",
			Breaking:      true,
		}
		body := []byte("\n## Description\n\nImplement the payment API.\n")

		if err := WriteFrontmatter(path, original, body); err != nil {
			t.Fatalf("WriteFrontmatter failed: %v", err)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("ReadFile failed: %v", err)
		}

		parsed, parsedBody, err := ParseFrontmatter(data)
		if err != nil {
			t.Fatalf("ParseFrontmatter failed: %v", err)
		}

		if parsed.SurfaceKey != original.SurfaceKey {
			t.Errorf("SurfaceKey roundtrip = %q, want %q", parsed.SurfaceKey, original.SurfaceKey)
		}
		if parsed.SurfaceType != original.SurfaceType {
			t.Errorf("SurfaceType roundtrip = %q, want %q", parsed.SurfaceType, original.SurfaceType)
		}
		if parsed.ID != original.ID {
			t.Errorf("ID roundtrip = %q, want %q", parsed.ID, original.ID)
		}
		if parsed.Title != original.Title {
			t.Errorf("Title roundtrip = %q, want %q", parsed.Title, original.Title)
		}
		if !parsed.Breaking {
			t.Error("Breaking should be preserved in roundtrip")
		}
		_ = parsedBody
	})
}

func TestParseFrontmatter_CoverageField(t *testing.T) {
	t.Run("coverage present", func(t *testing.T) {
		input := `---
id: "1"
title: "Task"
type: "coding.feature"
coverage: 95
---

Body`
		fm, _, err := ParseFrontmatter([]byte(input))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if fm.Coverage == nil {
			t.Fatal("expected Coverage non-nil")
		}
		if *fm.Coverage != 95 {
			t.Errorf("Coverage = %d, want 95", *fm.Coverage)
		}
	})

	t.Run("coverage absent is nil", func(t *testing.T) {
		input := `---
id: "1"
title: "Task"
type: "coding.feature"
---

Body`
		fm, _, err := ParseFrontmatter([]byte(input))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if fm.Coverage != nil {
			t.Errorf("expected Coverage nil, got %d", *fm.Coverage)
		}
	})

	t.Run("coverage zero is valid", func(t *testing.T) {
		input := `---
id: "1"
title: "Task"
type: "coding.feature"
coverage: 0
---

Body`
		fm, _, err := ParseFrontmatter([]byte(input))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if fm.Coverage == nil {
			t.Fatal("expected Coverage non-nil even for zero value")
		}
		if *fm.Coverage != 0 {
			t.Errorf("Coverage = %d, want 0", *fm.Coverage)
		}
	})
}

func TestParseFrontmatter_ComplexityField(t *testing.T) {
	t.Run("complexity present", func(t *testing.T) {
		input := `---
id: "1"
title: "Task"
type: "coding.feature"
complexity: low
---

Body`
		fm, _, err := ParseFrontmatter([]byte(input))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if fm.Complexity != "low" {
			t.Errorf("Complexity = %q, want %q", fm.Complexity, "low")
		}
	})

	t.Run("complexity absent is empty", func(t *testing.T) {
		input := `---
id: "1"
title: "Task"
type: "coding.feature"
---

Body`
		fm, _, err := ParseFrontmatter([]byte(input))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if fm.Complexity != "" {
			t.Errorf("expected Complexity empty, got %q", fm.Complexity)
		}
	})

	t.Run("complexity roundtrip", func(t *testing.T) {
		dir := t.TempDir()
		path := dir + "/roundtrip.md"

		original := FrontmatterData{
			ID:         "2.1",
			Title:      "Complex Task",
			Priority:   "P1",
			Type:       "coding.feature",
			Complexity: "high",
		}
		body := []byte("\n## Description\n\nImplement the feature.\n")

		if err := WriteFrontmatter(path, original, body); err != nil {
			t.Fatalf("WriteFrontmatter failed: %v", err)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("ReadFile failed: %v", err)
		}

		parsed, _, err := ParseFrontmatter(data)
		if err != nil {
			t.Fatalf("ParseFrontmatter failed: %v", err)
		}

		if parsed.Complexity != original.Complexity {
			t.Errorf("Complexity roundtrip = %q, want %q", parsed.Complexity, original.Complexity)
		}
	})
}

func TestParseFrontmatter(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantID    string
		wantTitle string
		wantType  string
		wantBody  string
		wantErr   bool
	}{
		{
			name: "valid full frontmatter",
			input: `---
id: "1.1"
title: "Define Interfaces"
priority: "P0"
estimated_time: "2h"
dependencies:
  - "1.0"
scope: "all"
breaking: true
mainSession: false
type: "coding.feature"
---

# Task body here
Some content`,
			wantID:    "1.1",
			wantTitle: "Define Interfaces",
			wantType:  "coding.feature",
			wantBody:  "\n\n# Task body here\nSome content",
		},
		{
			name: "minimal frontmatter",
			input: `---
id: "T-test-gen-scripts-api"
title: "Generate Test Cases"
type: "test.gen-scripts"
---

Body here`,
			wantID:    "T-test-gen-scripts-api",
			wantTitle: "Generate Test Cases",
			wantType:  "test.gen-scripts",
			wantBody:  "\n\nBody here",
		},
		{
			name:     "no frontmatter",
			input:    "# Just a markdown file\nNo frontmatter here",
			wantID:   "",
			wantBody: "# Just a markdown file\nNo frontmatter here",
		},
		{
			name:     "empty content",
			input:    "",
			wantID:   "",
			wantBody: "",
		},
		{
			name: "frontmatter with dependencies list",
			input: `---
id: "2.1"
title: "Task"
dependencies:
  - "1.1"
  - "1.gate"
---

Body`,
			wantID:    "2.1",
			wantTitle: "Task",
			wantBody:  "\n\nBody",
		},
		{
			name: "frontmatter with profile-suffixed type",
			input: `---
id: "T-test-2a"
title: "Generate Scripts (go)"
type: "test-pipeline.gen-scripts"
---

Call gen-test-scripts`,
			wantID:    "T-test-2a",
			wantTitle: "Generate Scripts (go)",
			wantType:  "test-pipeline.gen-scripts",
			wantBody:  "\n\nCall gen-test-scripts",
		},
		{
			name: "unknown noTest field gracefully ignored",
			input: `---
id: "1.summary"
title: "Summary"
noTest: true
---

Summary body`,
			wantID:    "1.summary",
			wantTitle: "Summary",
			wantBody:  "\n\nSummary body",
		},
		{
			name:      "closing --- at end without trailing newline",
			input:     "---\nid: \"1\"\ntitle: \"Task\"\n---",
			wantID:    "1",
			wantTitle: "Task",
			wantBody:  "",
		},
		{
			name: "frontmatter with coverage field",
			input: `---
id: "1"
title: "Task with coverage"
type: "coding.feature"
coverage: 90
---

Body here`,
			wantID:    "1",
			wantTitle: "Task with coverage",
			wantType:  "coding.feature",
			wantBody:  "\n\nBody here",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, body, err := ParseFrontmatter([]byte(tt.input))
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if fm.ID != tt.wantID {
				t.Errorf("ID = %q, want %q", fm.ID, tt.wantID)
			}
			if fm.Title != tt.wantTitle {
				t.Errorf("Title = %q, want %q", fm.Title, tt.wantTitle)
			}
			if fm.Type != tt.wantType {
				t.Errorf("Type = %q, want %q", fm.Type, tt.wantType)
			}
			if string(body) != tt.wantBody {
				t.Errorf("Body = %q, want %q", string(body), tt.wantBody)
			}
		})
	}
}
