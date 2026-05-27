package template

import (
	"strings"
	"testing"
)

func TestList(t *testing.T) {
	names := List()
	if len(names) == 0 {
		t.Fatal("List() returned no templates")
	}
	found := false
	for _, n := range names {
		if n == "coding.fix" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("List() missing coding.fix, got %v", names)
	}
}

func TestGet_FixTask(t *testing.T) {
	content, err := Get("coding.fix")
	if err != nil {
		t.Fatalf("Get(coding.fix) error: %v", err)
	}
	checks := []string{
		"{{ID}}",
		"{{TITLE}}",
		"{{DESCRIPTION}}",
		"{{SOURCE_TASK_ID}}",
		"{{SOURCE_FILES}}",
		"{{TEST_SCRIPT}}",
		"{{TEST_RESULTS}}",
		"## Reference Files",
		"## Verification",
		"targeted tests",
		"just test-e2e",
		"breaking: true",
		`priority: "P0"`,
	}
	for _, want := range checks {
		if !strings.Contains(content, want) {
			t.Errorf("coding.fix template missing %q", want)
		}
	}
}

func TestGet_NotFound(t *testing.T) {
	_, err := Get("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent template")
	}
}

func TestGetDefaults_FixTask(t *testing.T) {
	defs, err := GetDefaults("coding.fix")
	if err != nil {
		t.Fatalf("GetDefaults(coding.fix) error: %v", err)
	}
	if defs.Priority != "P0" {
		t.Errorf("Priority = %q, want P0", defs.Priority)
	}
	if !defs.Breaking {
		t.Error("Breaking = false, want true")
	}
	if defs.EstimatedTime != "30min" {
		t.Errorf("EstimatedTime = %q, want 30min", defs.EstimatedTime)
	}
}

func TestGetDefaults_CleanupTask(t *testing.T) {
	defs, err := GetDefaults("coding.cleanup")
	if err != nil {
		t.Fatalf("GetDefaults(coding.cleanup) error: %v", err)
	}
	if defs.Priority != "P0" {
		t.Errorf("Priority = %q, want P0", defs.Priority)
	}
	if defs.Breaking {
		t.Error("Breaking = true, want false (cleanup tasks are non-breaking)")
	}
	if defs.EstimatedTime != "15min" {
		t.Errorf("EstimatedTime = %q, want 15min", defs.EstimatedTime)
	}
}

func TestGetDefaults_NotFound(t *testing.T) {
	_, err := GetDefaults("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent template defaults")
	}
}

func TestTemplateType_ValidTypes(t *testing.T) {
	// Templates must use a type value that exists in ValidTypes.
	// Historically, fix-task.md used bare "fix" instead of "coding.fix".
	templateTypes := map[string]string{
		"coding.fix":     "",
		"coding.cleanup": "",
	}
	for name := range templateTypes {
		content, err := Get(name)
		if err != nil {
			t.Fatalf("Get(%q) error: %v", name, err)
		}
		// Extract type value from frontmatter
		for _, line := range strings.Split(content, "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "type:") {
				typ := strings.Trim(strings.TrimSpace(strings.TrimPrefix(line, "type:")), `"`)
				templateTypes[name] = typ
				break
			}
		}
		if templateTypes[name] == "" {
			t.Fatalf("%s template: could not extract type from frontmatter", name)
		}
	}

	// Verify each extracted type is a known valid type
	validTypes := map[string]bool{
		"coding.fix":            true,
		"coding.cleanup":        true,
		"code-quality.simplify": true,
	}
	for name, typ := range templateTypes {
		if !validTypes[typ] {
			t.Errorf("bug: %s template has invalid type %q — not in ValidTypes", name, typ)
		}
	}
}
