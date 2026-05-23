package task

import (
	"testing"
)

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
