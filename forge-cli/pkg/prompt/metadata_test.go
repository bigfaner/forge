package prompt

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseMetadataFrontmatter_WithFrontmatter(t *testing.T) {
	input := `---
type: coding.feature
category: coding
variables:
  - TaskID
  - TaskFile
  - SurfaceKey
---
TASK_ID: {{.TaskID}}
TASK_FILE: {{.TaskFile}}
`
	body, meta := parseMetadataFrontmatter(input)
	if meta == nil {
		t.Fatal("expected non-nil metadata")
	}
	if meta.Type != "coding.feature" {
		t.Errorf("Type = %q, want coding.feature", meta.Type)
	}
	if meta.Category != "coding" {
		t.Errorf("Category = %q, want coding", meta.Category)
	}
	if len(meta.Variables) != 3 {
		t.Fatalf("Variables = %v, want 3 items", meta.Variables)
	}
	wantVars := []string{"TaskID", "TaskFile", "SurfaceKey"}
	for i, v := range wantVars {
		if meta.Variables[i] != v {
			t.Errorf("Variables[%d] = %q, want %q", i, meta.Variables[i], v)
		}
	}
	if !strings.Contains(body, "TASK_ID") {
		t.Errorf("body should contain template content, got: %q", body)
	}
	if strings.Contains(body, "---") {
		t.Errorf("body should not contain frontmatter markers, got: %q", body)
	}
}

func TestParseMetadataFrontmatter_NoFrontmatter(t *testing.T) {
	input := "TASK_ID: {{.TaskID}}\n"
	body, meta := parseMetadataFrontmatter(input)
	if meta != nil {
		t.Error("expected nil metadata for content without frontmatter")
	}
	if body != input {
		t.Errorf("body should be unchanged, got: %q", body)
	}
}

func TestParseMetadataFrontmatter_EmptyVariables(t *testing.T) {
	input := `---
type: test.run
category: test
variables:
---
Body content
`
	body, meta := parseMetadataFrontmatter(input)
	if meta == nil {
		t.Fatal("expected non-nil metadata")
	}
	if len(meta.Variables) != 0 {
		t.Errorf("Variables = %v, want empty", meta.Variables)
	}
	if !strings.Contains(body, "Body content") {
		t.Errorf("body should contain content after frontmatter, got: %q", body)
	}
}

func TestStripMetadataFrontmatter_WithFrontmatter(t *testing.T) {
	input := "---\ntype: gate\ncategory: gate\n---\nBody"
	result := stripMetadataFrontmatter(input)
	if result != "Body" {
		t.Errorf("expected %q, got %q", "Body", result)
	}
}

func TestStripMetadataFrontmatter_NoFrontmatter(t *testing.T) {
	input := "Body without frontmatter"
	result := stripMetadataFrontmatter(input)
	if result != input {
		t.Errorf("expected %q, got %q", input, result)
	}
}

func TestValidateMetadataVariables_AllMatch(t *testing.T) {
	meta := &TemplateMetadata{
		Variables: []string{"TaskID", "TaskFile", "SurfaceKey"},
	}
	structType := reflect.TypeOf(promptTemplateData{})
	if err := validateMetadataVariables(meta, structType); err != nil {
		t.Errorf("expected no error for valid variables, got: %v", err)
	}
}

func TestValidateMetadataVariables_Mismatch(t *testing.T) {
	meta := &TemplateMetadata{
		Variables: []string{"TaskID", "NonExistentField"},
	}
	structType := reflect.TypeOf(promptTemplateData{})
	err := validateMetadataVariables(meta, structType)
	if err == nil {
		t.Error("expected error for invalid variable name")
	}
	if !strings.Contains(err.Error(), "NonExistentField") {
		t.Errorf("error should mention mismatched variable, got: %v", err)
	}
}

func TestValidateMetadataVariables_NilMeta(t *testing.T) {
	structType := reflect.TypeOf(promptTemplateData{})
	if err := validateMetadataVariables(nil, structType); err != nil {
		t.Errorf("expected no error for nil metadata, got: %v", err)
	}
}

func TestValidateMetadataVariables_EmptyVariables(t *testing.T) {
	meta := &TemplateMetadata{
		Variables: []string{},
	}
	structType := reflect.TypeOf(promptTemplateData{})
	if err := validateMetadataVariables(meta, structType); err != nil {
		t.Errorf("expected no error for empty variables, got: %v", err)
	}
}
