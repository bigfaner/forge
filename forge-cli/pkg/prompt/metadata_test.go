package prompt

import (
	"reflect"
	"strings"
	"testing"
)

// --- Grouped frontmatter parsing tests (Task 4) ---

func TestParseMetadataFrontmatter_GroupedFormat(t *testing.T) {
	input := `---
type: coding.feature
category: coding
identity:
  TaskID: true
  TaskFile: true
context:
  FeatureSlug: true
  SurfaceKey: true
conditional:
  PhaseSummary: true
  CoverageStrategy: true
variables:
  - TaskCategory
  - SurfaceType
---
Body content here
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

	// Identity group
	if len(meta.Identity) != 2 {
		t.Fatalf("Identity = %v, want 2 entries", meta.Identity)
	}
	if !meta.Identity["TaskID"] {
		t.Error("Identity should contain TaskID")
	}
	if !meta.Identity["TaskFile"] {
		t.Error("Identity should contain TaskFile")
	}

	// Context group
	if len(meta.Context) != 2 {
		t.Fatalf("Context = %v, want 2 entries", meta.Context)
	}
	if !meta.Context["FeatureSlug"] {
		t.Error("Context should contain FeatureSlug")
	}

	// Conditional group
	if len(meta.Conditional) != 2 {
		t.Fatalf("Conditional = %v, want 2 entries", meta.Conditional)
	}
	if !meta.Conditional["PhaseSummary"] {
		t.Error("Conditional should contain PhaseSummary")
	}

	// Variables list (for backward compatibility)
	if len(meta.Variables) != 2 {
		t.Fatalf("Variables = %v, want 2 items", meta.Variables)
	}
	if meta.Variables[0] != "TaskCategory" {
		t.Errorf("Variables[0] = %q, want TaskCategory", meta.Variables[0])
	}

	if !strings.Contains(body, "Body content here") {
		t.Errorf("body should contain content after frontmatter, got: %q", body)
	}
}

func TestParseMetadataFrontmatter_GroupedWithEmptyVariables(t *testing.T) {
	input := `---
type: coding.feature
category: coding
identity:
  TaskID: true
context:
  FeatureSlug: true
---
Body content
`
	body, meta := parseMetadataFrontmatter(input)
	if meta == nil {
		t.Fatal("expected non-nil metadata")
	}
	if len(meta.Identity) != 1 {
		t.Errorf("Identity = %v, want 1 entry", meta.Identity)
	}
	if len(meta.Variables) != 0 {
		t.Errorf("Variables = %v, want empty", meta.Variables)
	}
	if !strings.Contains(body, "Body content") {
		t.Errorf("body should contain content, got: %q", body)
	}
}

func TestParseMetadataFrontmatter_BackwardCompat_OldFlatFormat(t *testing.T) {
	input := `---
type: coding.feature
category: coding
variables:
  - TaskID
  - TaskFile
  - SurfaceKey
---
Body content
`
	body, meta := parseMetadataFrontmatter(input)
	if meta == nil {
		t.Fatal("expected non-nil metadata")
	}
	if meta.Type != "coding.feature" {
		t.Errorf("Type = %q, want coding.feature", meta.Type)
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
	// No groups in old format
	if len(meta.Identity) != 0 {
		t.Errorf("Identity should be empty for old format, got %v", meta.Identity)
	}
	if len(meta.Context) != 0 {
		t.Errorf("Context should be empty for old format, got %v", meta.Context)
	}
	if len(meta.Conditional) != 0 {
		t.Errorf("Conditional should be empty for old format, got %v", meta.Conditional)
	}
	if !strings.Contains(body, "Body content") {
		t.Errorf("body should contain content, got: %q", body)
	}
}

func TestAllFields_GroupedFormat(t *testing.T) {
	meta := &TemplateMetadata{
		Identity:    map[string]bool{"TaskID": true, "TaskFile": true},
		Context:     map[string]bool{"FeatureSlug": true},
		Conditional: map[string]bool{"PhaseSummary": true},
		Variables:   []string{"TaskCategory", "SurfaceType"},
	}
	all := meta.AllFields()

	expected := map[string]bool{
		"TaskID": true, "TaskFile": true, "FeatureSlug": true,
		"PhaseSummary": true, "TaskCategory": true, "SurfaceType": true,
	}
	if len(all) != len(expected) {
		t.Fatalf("AllFields() returned %d items, want %d: %v", len(all), len(expected), all)
	}
	for _, f := range all {
		if !expected[f] {
			t.Errorf("unexpected field in AllFields(): %q", f)
		}
	}
}

func TestAllFields_OldFormat(t *testing.T) {
	meta := &TemplateMetadata{
		Variables: []string{"TaskID", "TaskFile", "SurfaceKey"},
	}
	all := meta.AllFields()
	if len(all) != 3 {
		t.Fatalf("AllFields() returned %d items, want 3: %v", len(all), all)
	}
	for _, v := range []string{"TaskID", "TaskFile", "SurfaceKey"} {
		found := false
		for _, f := range all {
			if f == v {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("AllFields() missing %q", v)
		}
	}
}

func TestAllFields_Empty(t *testing.T) {
	meta := &TemplateMetadata{}
	all := meta.AllFields()
	if len(all) != 0 {
		t.Errorf("AllFields() should return empty for zero metadata, got %v", all)
	}
}

func TestValidateMetadataVariables_GroupedFields(t *testing.T) {
	meta := &TemplateMetadata{
		Identity:    map[string]bool{"TaskID": true, "TaskFile": true},
		Context:     map[string]bool{"FeatureSlug": true},
		Conditional: map[string]bool{"PhaseSummary": true},
		Variables:   []string{"SurfaceKey", "SurfaceType"},
	}
	structType := reflect.TypeOf(promptTemplateData{})
	if err := validateMetadataVariables(meta, structType); err != nil {
		t.Errorf("expected no error for valid grouped variables, got: %v", err)
	}
}

func TestValidateMetadataVariables_GroupedMismatch(t *testing.T) {
	meta := &TemplateMetadata{
		Identity: map[string]bool{"TaskID": true, "NonExistentField": true},
	}
	structType := reflect.TypeOf(promptTemplateData{})
	err := validateMetadataVariables(meta, structType)
	if err == nil {
		t.Error("expected error for invalid identity variable name")
	}
	if !strings.Contains(err.Error(), "NonExistentField") {
		t.Errorf("error should mention mismatched variable, got: %v", err)
	}
}

func TestValidateMetadataVariables_ContextMismatch(t *testing.T) {
	meta := &TemplateMetadata{
		Context: map[string]bool{"BogusContextField": true},
	}
	structType := reflect.TypeOf(promptTemplateData{})
	err := validateMetadataVariables(meta, structType)
	if err == nil {
		t.Error("expected error for invalid context variable name")
	}
	if !strings.Contains(err.Error(), "BogusContextField") {
		t.Errorf("error should mention mismatched variable, got: %v", err)
	}
}

func TestValidateMetadataVariables_ConditionalMismatch(t *testing.T) {
	meta := &TemplateMetadata{
		Conditional: map[string]bool{"FakeConditional": true},
	}
	structType := reflect.TypeOf(promptTemplateData{})
	err := validateMetadataVariables(meta, structType)
	if err == nil {
		t.Error("expected error for invalid conditional variable name")
	}
	if !strings.Contains(err.Error(), "FakeConditional") {
		t.Errorf("error should mention mismatched variable, got: %v", err)
	}
}

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
