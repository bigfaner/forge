package task

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeTaskMD writes a minimal task .md file with frontmatter.
// Defaults type to "feature" for tasks whose ID is not auto-inferable.
func writeTaskMD(t *testing.T, dir, filename, id, title string, deps []string) {
	t.Helper()
	var depLine string
	if len(deps) > 0 {
		quoted := make([]string, len(deps))
		for i, d := range deps {
			quoted[i] = `"` + d + `"`
		}
		depLine = "dependencies:\n  - " + joinStrings(quoted, "\n  - ") + "\n"
	}
	// InferType for business IDs returns "", so we must set type explicitly.
	// Auto-gen IDs (gates, summaries, test tasks) get their type from InferType.
	taskType := InferType(id)
	if taskType == "" {
		taskType = TypeFeature
	}
	content := "---\nid: " + `"` + id + `"` + "\ntitle: " + `"` + title + `"` +
		"\npriority: \"P1\"\nestimated_time: \"1h\"\ntype: " + `"` + taskType + `"` +
		"\n" + depLine + "scope: \"all\"\n---\n\n# " + title + "\n"
	if err := os.WriteFile(filepath.Join(dir, filename), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func joinStrings(ss []string, sep string) string {
	return strings.Join(ss, sep)
}

// writeTaskMDWithType writes a minimal task .md file with frontmatter including a type field.
func writeTaskMDWithType(t *testing.T, dir, filename, id, title, taskType string, deps []string) {
	t.Helper()
	var depLine string
	if len(deps) > 0 {
		quoted := make([]string, len(deps))
		for i, d := range deps {
			quoted[i] = `"` + d + `"`
		}
		depLine = "dependencies:\n  - " + joinStrings(quoted, "\n  - ") + "\n"
	}
	content := "---\nid: " + `"` + id + `"` + "\ntitle: " + `"` + title + `"` +
		"\npriority: \"P1\"\nestimated_time: \"1h\"\ntype: " + `"` + taskType + `"` +
		"\n" + depLine + "scope: \"all\"\n---\n\n# " + title + "\n"
	if err := os.WriteFile(filepath.Join(dir, filename), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

// setupBuildEnv creates a temp project root with feature dirs.
// mode: "breakdown" creates prd/prd-spec.md, "quick" creates proposal, "" creates nothing.
func setupBuildEnv(t *testing.T, mode string) (projectRoot, tasksDir, indexPath string) {
	t.Helper()
	projectRoot = t.TempDir()
	featureSlug := "test-feature"
	featureDir := filepath.Join(projectRoot, "docs", "features", featureSlug)
	tasksDir = filepath.Join(featureDir, "tasks")
	indexPath = filepath.Join(tasksDir, "index.json")

	if err := os.MkdirAll(tasksDir, 0755); err != nil {
		t.Fatal(err)
	}

	switch mode {
	case "breakdown":
		prdDir := filepath.Join(featureDir, "prd")
		if err := os.MkdirAll(prdDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(prdDir, "prd-spec.md"), []byte("# PRD"), 0644); err != nil {
			t.Fatal(err)
		}
	case "quick":
		propDir := filepath.Join(projectRoot, "docs", "proposals", featureSlug)
		if err := os.MkdirAll(propDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(propDir, "proposal.md"), []byte("# Proposal"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	return projectRoot, tasksDir, indexPath
}

func TestBuildIndex_FreshBuild(t *testing.T) {
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "")

	writeTaskMD(t, tasksDir, "1-foo.md", "1", "Foo Task", nil)
	writeTaskMD(t, tasksDir, "2-bar.md", "2", "Bar Task", []string{"1"})

	opts := BuildIndexOpts{
		FeatureSlug: "test-feature",
		ProjectRoot: projectRoot,
		TasksDir:    tasksDir,
		IndexPath:   indexPath,
	}

	result, err := BuildIndex(opts)
	if err != nil {
		t.Fatalf("BuildIndex error: %v", err)
	}
	if result.NewCount != 2 {
		t.Errorf("NewCount = %d, want 2", result.NewCount)
	}
	if result.UpdatedCount != 0 {
		t.Errorf("UpdatedCount = %d, want 0", result.UpdatedCount)
	}

	// Verify index.json written
	data, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("read index.json: %v", err)
	}
	var idx taskIndexJSON
	if err := json.Unmarshal(data, &idx); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(idx.Tasks) != 2 {
		t.Errorf("tasks count = %d, want 2", len(idx.Tasks))
	}
	if idx.Tasks["1-foo"].Status != "pending" {
		t.Errorf("1-foo status = %q, want pending", idx.Tasks["1-foo"].Status)
	}
	if idx.Feature != "test-feature" {
		t.Errorf("feature = %q, want test-feature", idx.Feature)
	}
}

func TestBuildIndex_IdempotentRebuild(t *testing.T) {
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "")

	writeTaskMD(t, tasksDir, "1-foo.md", "1", "Foo Task", nil)

	opts := BuildIndexOpts{
		FeatureSlug: "test-feature",
		ProjectRoot: projectRoot,
		TasksDir:    tasksDir,
		IndexPath:   indexPath,
	}

	// First build
	_, err := BuildIndex(opts)
	if err != nil {
		t.Fatalf("first build: %v", err)
	}

	// Read index to capture created date
	data1, _ := os.ReadFile(indexPath)
	var idx1 taskIndexJSON
	_ = json.Unmarshal(data1, &idx1)

	// Second build (no changes)
	result2, err := BuildIndex(opts)
	if err != nil {
		t.Fatalf("second build: %v", err)
	}
	if result2.UpdatedCount != 1 {
		t.Errorf("UpdatedCount = %d, want 1", result2.UpdatedCount)
	}
	if result2.NewCount != 0 {
		t.Errorf("NewCount = %d, want 0", result2.NewCount)
	}

	// Created date preserved
	data2, _ := os.ReadFile(indexPath)
	var idx2 taskIndexJSON
	_ = json.Unmarshal(data2, &idx2)
	if idx2.Created != idx1.Created {
		t.Errorf("created date changed: %q -> %q", idx1.Created, idx2.Created)
	}
}

func TestBuildIndex_StatusPreservation(t *testing.T) {
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "")

	writeTaskMD(t, tasksDir, "1-foo.md", "1", "Foo Task", nil)

	// Build first time
	opts := BuildIndexOpts{
		FeatureSlug: "test-feature",
		ProjectRoot: projectRoot,
		TasksDir:    tasksDir,
		IndexPath:   indexPath,
	}
	if _, err := BuildIndex(opts); err != nil {
		t.Fatalf("first build: %v", err)
	}

	// Manually modify status
	data, _ := os.ReadFile(indexPath)
	var raw map[string]json.RawMessage
	_ = json.Unmarshal(data, &raw)
	var tasksMap map[string]json.RawMessage
	_ = json.Unmarshal(raw["tasks"], &tasksMap)

	var task1 map[string]any
	_ = json.Unmarshal(tasksMap["1-foo"], &task1)
	task1["status"] = "in_progress"
	task1["sourceTaskID"] = "some-source"
	task1["blockedReason"] = "waiting for review"
	updated, _ := json.Marshal(task1)
	tasksMap["1-foo"] = updated
	raw["tasks"], _ = json.Marshal(tasksMap)
	finalData, _ := json.MarshalIndent(raw, "", "  ")
	if err := os.WriteFile(indexPath, append(finalData, '\n'), 0644); err != nil {
		t.Fatal(err)
	}

	// Rebuild
	result, err := BuildIndex(opts)
	if err != nil {
		t.Fatalf("rebuild: %v", err)
	}

	// Verify status preserved
	data2, _ := os.ReadFile(indexPath)
	var idx taskIndexJSON
	_ = json.Unmarshal(data2, &idx)
	if idx.Tasks["1-foo"].Status != "in_progress" {
		t.Errorf("status = %q, want in_progress", idx.Tasks["1-foo"].Status)
	}
	if idx.Tasks["1-foo"].SourceTaskID != "some-source" {
		t.Errorf("sourceTaskID = %q, want some-source", idx.Tasks["1-foo"].SourceTaskID)
	}
	if idx.Tasks["1-foo"].BlockedReason != "waiting for review" {
		t.Errorf("blockedReason = %q, want waiting for review", idx.Tasks["1-foo"].BlockedReason)
	}
	_ = result
}

func TestBuildIndex_NewMDAdded(t *testing.T) {
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "")

	writeTaskMD(t, tasksDir, "1-foo.md", "1", "Foo Task", nil)

	opts := BuildIndexOpts{
		FeatureSlug: "test-feature",
		ProjectRoot: projectRoot,
		TasksDir:    tasksDir,
		IndexPath:   indexPath,
	}
	if _, err := BuildIndex(opts); err != nil {
		t.Fatalf("first build: %v", err)
	}

	// Add new task
	writeTaskMD(t, tasksDir, "2-bar.md", "2", "Bar Task", []string{"1"})

	result, err := BuildIndex(opts)
	if err != nil {
		t.Fatalf("rebuild: %v", err)
	}
	if result.NewCount != 1 {
		t.Errorf("NewCount = %d, want 1", result.NewCount)
	}
}

func TestBuildIndex_FrontmatterUpdate(t *testing.T) {
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "")

	writeTaskMD(t, tasksDir, "1-foo.md", "1", "Old Title", nil)

	opts := BuildIndexOpts{
		FeatureSlug: "test-feature",
		ProjectRoot: projectRoot,
		TasksDir:    tasksDir,
		IndexPath:   indexPath,
	}
	if _, err := BuildIndex(opts); err != nil {
		t.Fatalf("first build: %v", err)
	}

	// Update the .md with new title
	writeTaskMD(t, tasksDir, "1-foo.md", "1", "New Title", nil)

	if _, err := BuildIndex(opts); err != nil {
		t.Fatalf("rebuild: %v", err)
	}

	data, _ := os.ReadFile(indexPath)
	var idx taskIndexJSON
	_ = json.Unmarshal(data, &idx)
	if idx.Tasks["1-foo"].Title != "New Title" {
		t.Errorf("title = %q, want New Title", idx.Tasks["1-foo"].Title)
	}
}

func TestBuildIndex_OrphanDetection(t *testing.T) {
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "")

	// Create index with a task
	writeTaskMD(t, tasksDir, "1-foo.md", "1", "Foo", nil)
	opts := BuildIndexOpts{
		FeatureSlug: "test-feature",
		ProjectRoot: projectRoot,
		TasksDir:    tasksDir,
		IndexPath:   indexPath,
	}
	if _, err := BuildIndex(opts); err != nil {
		t.Fatalf("first build: %v", err)
	}

	// Remove the .md file
	if err := os.Remove(filepath.Join(tasksDir, "1-foo.md")); err != nil {
		t.Fatal(err)
	}

	// Add a different task so dir isn't empty
	writeTaskMD(t, tasksDir, "2-bar.md", "2", "Bar", nil)

	result, err := BuildIndex(opts)
	if err != nil {
		t.Fatalf("rebuild: %v", err)
	}
	if result.PreservedCount != 1 {
		t.Errorf("PreservedCount = %d, want 1", result.PreservedCount)
	}
	found := false
	for _, w := range result.Warnings {
		if len(w) > 6 && w[:6] == "orphan" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected orphan warning, got %v", result.Warnings)
	}
}

func TestBuildIndex_NoProfilesSkipsTestGen(t *testing.T) {
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "breakdown")

	writeTaskMD(t, tasksDir, "1-foo.md", "1", "Foo", nil)

	opts := BuildIndexOpts{
		FeatureSlug: "test-feature",
		ProjectRoot: projectRoot,
		TasksDir:    tasksDir,
		IndexPath:   indexPath,
	}

	result, err := BuildIndex(opts)
	if err != nil {
		t.Fatalf("BuildIndex: %v", err)
	}

	// Only 1 business task, no test tasks (no profiles provided)
	if result.NewCount != 1 {
		t.Errorf("NewCount = %d, want 1", result.NewCount)
	}

	// No test task .md files generated
	entries, _ := os.ReadDir(tasksDir)
	for _, e := range entries {
		if e.Name() != "1-foo.md" && e.Name() != "index.json" {
			t.Errorf("unexpected file: %s", e.Name())
		}
	}
}

func TestBuildIndex_ModeDetection(t *testing.T) {
	tests := []struct {
		name         string
		mode         string
		wantFeature  string
		wantProposal string
	}{
		{"breakdown sets PRD", "breakdown", "prd/prd-spec.md", ""},
		{"quick sets proposal", "quick", "", "docs/proposals/test-feature/proposal.md"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectRoot, tasksDir, indexPath := setupBuildEnv(t, tt.mode)

			writeTaskMD(t, tasksDir, "1-foo.md", "1", "Foo", nil)

			opts := BuildIndexOpts{
				FeatureSlug: "test-feature",
				ProjectRoot: projectRoot,
				TasksDir:    tasksDir,
				IndexPath:   indexPath,
			}

			if _, err := BuildIndex(opts); err != nil {
				t.Fatalf("BuildIndex: %v", err)
			}

			data, _ := os.ReadFile(indexPath)
			var idx taskIndexJSON
			_ = json.Unmarshal(data, &idx)

			if tt.mode == "breakdown" {
				if idx.PRD != "prd/prd-spec.md" {
					t.Errorf("PRD = %q, want prd/prd-spec.md", idx.PRD)
				}
				if strings.Contains(idx.PRD, "\\") {
					t.Errorf("PRD path contains backslash: %q", idx.PRD)
				}
			}
			if tt.wantProposal != "" {
				if idx.Proposal != tt.wantProposal {
					t.Errorf("Proposal = %q, want %q", idx.Proposal, tt.wantProposal)
				}
				if strings.Contains(idx.Proposal, "\\") {
					t.Errorf("Proposal path contains backslash: %q", idx.Proposal)
				}
			}
		})
	}
}

func TestBuildIndex_SkipNoID(t *testing.T) {
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "")

	// Write .md without id in frontmatter
	content := "---\ntitle: \"No ID\"\npriority: \"P1\"\n---\n\n# No ID task\n"
	if err := os.WriteFile(filepath.Join(tasksDir, "no-id.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	writeTaskMD(t, tasksDir, "1-foo.md", "1", "Foo", nil)

	opts := BuildIndexOpts{
		FeatureSlug: "test-feature",
		ProjectRoot: projectRoot,
		TasksDir:    tasksDir,
		IndexPath:   indexPath,
	}

	result, err := BuildIndex(opts)
	if err != nil {
		t.Fatalf("BuildIndex: %v", err)
	}
	if result.NewCount != 1 {
		t.Errorf("NewCount = %d, want 1 (no-id skipped)", result.NewCount)
	}
	// Should have warning about no id
	found := false
	for _, w := range result.Warnings {
		if len(w) >= 4 && w[:4] == "skip" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected skip warning, got %v", result.Warnings)
	}
}

func TestBuildIndex_SkipUnderscoreFiles(t *testing.T) {
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "")

	writeTaskMD(t, tasksDir, "1-foo.md", "1", "Foo", nil)

	// Create _template.md (should be skipped)
	if err := os.WriteFile(filepath.Join(tasksDir, "_template.md"), []byte("---\nid: \"skip-me\"\n---\n"), 0644); err != nil {
		t.Fatal(err)
	}

	opts := BuildIndexOpts{
		FeatureSlug: "test-feature",
		ProjectRoot: projectRoot,
		TasksDir:    tasksDir,
		IndexPath:   indexPath,
	}

	result, err := BuildIndex(opts)
	if err != nil {
		t.Fatalf("BuildIndex: %v", err)
	}
	if result.NewCount != 1 {
		t.Errorf("NewCount = %d, want 1 (_template.md skipped)", result.NewCount)
	}
}

func TestBuildIndex_TypeInference(t *testing.T) {
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "")

	// Task with explicit type
	content := "---\nid: \"1\"\ntitle: \"Gate\"\npriority: \"P1\"\ntype: \"gate\"\nscope: \"all\"\n---\n\n# Gate\n"
	if err := os.WriteFile(filepath.Join(tasksDir, "1-gate.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// Task without type (should infer)
	writeTaskMD(t, tasksDir, "2-bar.md", "2", "Bar", nil)

	opts := BuildIndexOpts{
		FeatureSlug: "test-feature",
		ProjectRoot: projectRoot,
		TasksDir:    tasksDir,
		IndexPath:   indexPath,
	}

	if _, err := BuildIndex(opts); err != nil {
		t.Fatalf("BuildIndex: %v", err)
	}

	data, _ := os.ReadFile(indexPath)
	var idx taskIndexJSON
	_ = json.Unmarshal(data, &idx)

	if idx.Tasks["1-gate"].Type != "gate" {
		t.Errorf("explicit type = %q, want gate", idx.Tasks["1-gate"].Type)
	}
	if idx.Tasks["2-bar"].Type != TypeFeature {
		t.Errorf("inferred type = %q, want %q", idx.Tasks["2-bar"].Type, TypeFeature)
	}
}

func TestBuildIndex_EmptyTasksDir(t *testing.T) {
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "")

	opts := BuildIndexOpts{
		FeatureSlug: "test-feature",
		ProjectRoot: projectRoot,
		TasksDir:    tasksDir,
		IndexPath:   indexPath,
	}

	result, err := BuildIndex(opts)
	if err != nil {
		t.Fatalf("BuildIndex: %v", err)
	}
	if result.NewCount != 0 {
		t.Errorf("NewCount = %d, want 0", result.NewCount)
	}
	if result.UpdatedCount != 0 {
		t.Errorf("UpdatedCount = %d, want 0", result.UpdatedCount)
	}
}

func TestBuildIndex_WithTestTasks(t *testing.T) {
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "breakdown")

	// Create a business task (feature type) and gate tasks for dep resolution
	writeTaskMDWithType(t, tasksDir, "1-feat.md", "1.1", "Feature Task", TypeFeature, nil)
	writeTaskMD(t, tasksDir, "1-gate.md", "1.gate", "Phase 1 Gate", nil)
	writeTaskMD(t, tasksDir, "2-gate.md", "2.gate", "Phase 2 Gate", nil)

	opts := BuildIndexOpts{
		FeatureSlug:      "test-feature",
		ProjectRoot:      projectRoot,
		TasksDir:         tasksDir,
		IndexPath:        indexPath,
		TestProfiles:     []string{"go-test"},
		TestCapabilities: []string{"cli"},
	}

	result, err := BuildIndex(opts)
	if err != nil {
		t.Fatalf("BuildIndex: %v", err)
	}

	// 1 business + 2 gates + 7 test tasks (gen-cases, eval-cases, gen-scripts-cli, run, graduate, verify-regression, consolidate) = 10
	total := result.NewCount + result.UpdatedCount
	if total != 10 {
		t.Errorf("total tasks = %d (new=%d, updated=%d), want 10", total, result.NewCount, result.UpdatedCount)
	}

	// Verify test task .md files were generated
	if _, err := os.Stat(filepath.Join(tasksDir, "gen-test-cases.md")); os.IsNotExist(err) {
		t.Error("gen-test-cases.md not generated")
	}
	if _, err := os.Stat(filepath.Join(tasksDir, "gen-test-scripts-go-test-cli.md")); os.IsNotExist(err) {
		t.Error("gen-test-scripts-go-test-cli.md not generated")
	}

	// Verify index contains test tasks
	data, _ := os.ReadFile(indexPath)
	var idx taskIndexJSON
	_ = json.Unmarshal(data, &idx)

	if _, ok := idx.Tasks["gen-test-cases"]; !ok {
		t.Error("missing gen-test-cases in index")
	}
	if _, ok := idx.Tasks["gen-test-scripts-go-test-cli"]; !ok {
		t.Error("missing gen-test-scripts-go-test in index")
	}

	// T-test-1 should depend on highest gate
	if len(idx.Tasks["gen-test-cases"].Dependencies) == 0 {
		t.Error("T-test-1 has no dependencies")
	} else if idx.Tasks["gen-test-cases"].Dependencies[0] != "2.gate" {
		t.Errorf("T-test-1 dep = %v, want [2.gate]", idx.Tasks["gen-test-cases"].Dependencies)
	}

	// Verify Profile field: per-profile test tasks have profile set
	if idx.Tasks["gen-test-scripts-go-test-cli"].Profile != "go-test" {
		t.Errorf("gen-test-scripts-go-test profile = %q, want go-test", idx.Tasks["gen-test-scripts-go-test-cli"].Profile)
	}
	// Shared test task (gen-test-cases in breakdown mode) has empty profile
	if idx.Tasks["gen-test-cases"].Profile != "" {
		t.Errorf("gen-test-cases profile = %q, want empty", idx.Tasks["gen-test-cases"].Profile)
	}
	// Business task has empty profile
	if idx.Tasks["1-gate"].Profile != "" {
		t.Errorf("1-gate profile = %q, want empty", idx.Tasks["1-gate"].Profile)
	}
}

func TestBuildIndex_TestTasksIdempotent(t *testing.T) {
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "breakdown")

	writeTaskMDWithType(t, tasksDir, "1-feat.md", "1.1", "Feature Task", TypeFeature, nil)
	writeTaskMD(t, tasksDir, "1-gate.md", "1.gate", "Gate", nil)

	opts := BuildIndexOpts{
		FeatureSlug:      "test-feature",
		ProjectRoot:      projectRoot,
		TasksDir:         tasksDir,
		IndexPath:        indexPath,
		TestProfiles:     []string{"go-test"},
		TestCapabilities: []string{"cli"},
	}

	// First build
	if _, err := BuildIndex(opts); err != nil {
		t.Fatalf("first build: %v", err)
	}

	// Read first build
	data1, _ := os.ReadFile(indexPath)
	var idx1 taskIndexJSON
	_ = json.Unmarshal(data1, &idx1)

	// Second build
	result, err := BuildIndex(opts)
	if err != nil {
		t.Fatalf("second build: %v", err)
	}

	// No new tasks (all updated)
	if result.NewCount != 0 {
		t.Errorf("NewCount = %d, want 0 on second build", result.NewCount)
	}

	// Test task .md files should still exist (not regenerated)
	if _, err := os.Stat(filepath.Join(tasksDir, "gen-test-cases.md")); err != nil {
		t.Error("gen-test-cases.md missing after rebuild")
	}
}

func TestBuildIndex_MultiProfile(t *testing.T) {
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "quick")

	writeTaskMD(t, tasksDir, "1-foo.md", "1", "Foo Task", nil)

	opts := BuildIndexOpts{
		FeatureSlug:      "test-feature",
		ProjectRoot:      projectRoot,
		TasksDir:         tasksDir,
		IndexPath:        indexPath,
		TestProfiles:     []string{"go-test", "web-playwright"},
		TestCapabilities: []string{"api"},
	}

	result, err := BuildIndex(opts)
	if err != nil {
		t.Fatalf("BuildIndex: %v", err)
	}

	// 1 business + 3*2 per-profile + 2 shared = 9
	total := result.NewCount + result.UpdatedCount
	if total != 9 {
		t.Errorf("total = %d (new=%d, updated=%d), want 9", total, result.NewCount, result.UpdatedCount)
	}

	// Verify suffixed .md files exist
	for _, name := range []string{
		"quick-test-cases-go-test.md",
		"quick-test-cases-web-playwright.md",
	} {
		if _, err := os.Stat(filepath.Join(tasksDir, name)); os.IsNotExist(err) {
			t.Errorf("%s not generated", name)
		}
	}

	// Verify Profile field: per-profile test tasks have profile set
	data, _ := os.ReadFile(indexPath)
	var idx taskIndexJSON
	_ = json.Unmarshal(data, &idx)

	if idx.Tasks["quick-test-cases-go-test"].Profile != "go-test" {
		t.Errorf("quick-test-cases-go-test profile = %q, want go-test", idx.Tasks["quick-test-cases-go-test"].Profile)
	}
	if idx.Tasks["quick-test-cases-web-playwright"].Profile != "web-playwright" {
		t.Errorf("quick-test-cases-web-playwright profile = %q, want web-playwright", idx.Tasks["quick-test-cases-web-playwright"].Profile)
	}
	// Business task has empty profile
	if idx.Tasks["1-foo"].Profile != "" {
		t.Errorf("1-foo profile = %q, want empty", idx.Tasks["1-foo"].Profile)
	}
}

func TestDetectMode(t *testing.T) {
	tests := []struct {
		name string
		mode string
		want string
	}{
		{"breakdown", "breakdown", "breakdown"},
		{"quick", "quick", "quick"},
		{"neither", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectRoot, _, _ := setupBuildEnv(t, tt.mode)
			got := detectMode(projectRoot, "test-feature")
			if got != tt.want {
				t.Errorf("detectMode = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestShouldSkipFile(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"_template.md", true},
		{"_private.md", true},
		{"index.json", true},
		{"1-foo.md", false},
		{"gen-test-cases.md", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldSkipFile(tt.name); got != tt.want {
				t.Errorf("shouldSkipFile(%q) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestIsTestTaskID(t *testing.T) {
	tests := []struct {
		id   string
		want bool
	}{
		{"T-test-1", true},
		{"T-test-2a", true},
		{"T-quick-1", true},
		{"1", false},
		{"1.gate", false},
		{"fix-1", false},
	}
	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			if got := isTestTaskID(tt.id); got != tt.want {
				t.Errorf("isTestTaskID(%q) = %v, want %v", tt.id, got, tt.want)
			}
		})
	}
}

func TestBuildIndex_StageGatesGenerated(t *testing.T) {
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "")

	// Create tasks in 2 phases: phase 1 has 2 tasks, phase 2 has 1 task
	writeTaskMD(t, tasksDir, "1-foo.md", "1.1", "Phase 1 Task 1", nil)
	writeTaskMD(t, tasksDir, "2-bar.md", "1.2", "Phase 1 Task 2", []string{"1.1"})
	writeTaskMD(t, tasksDir, "3-baz.md", "2.1", "Phase 2 Task 1", []string{"1.2"})

	opts := BuildIndexOpts{
		FeatureSlug: "test-feature",
		ProjectRoot: projectRoot,
		TasksDir:    tasksDir,
		IndexPath:   indexPath,
	}

	result, err := BuildIndex(opts)
	if err != nil {
		t.Fatalf("BuildIndex error: %v", err)
	}

	// Phase 1 has 2 business tasks -> should generate summary + gate = 2 files
	// Phase 2 has 1 business task -> below threshold, no generation
	if result.StageGatesGenerated != 2 {
		t.Errorf("StageGatesGenerated = %d, want 2", result.StageGatesGenerated)
	}

	// Verify files exist
	if _, err := os.Stat(filepath.Join(tasksDir, "1.summary.md")); os.IsNotExist(err) {
		t.Error("1.summary.md not generated")
	}
	if _, err := os.Stat(filepath.Join(tasksDir, "1.gate.md")); os.IsNotExist(err) {
		t.Error("1.gate.md not generated")
	}
	if _, err := os.Stat(filepath.Join(tasksDir, "2.summary.md")); err == nil {
		t.Error("2.summary.md should NOT be generated (single-task phase)")
	}

	// Verify tasks appear in index.json with correct types
	data, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("read index.json: %v", err)
	}
	var idx taskIndexJSON
	if err := json.Unmarshal(data, &idx); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if task, ok := idx.Tasks["1.summary"]; !ok {
		t.Error("1.summary not in index")
	} else if task.Type != TypeDocGenerationSummary {
		t.Errorf("1.summary type = %q, want %q", task.Type, TypeDocGenerationSummary)
	}
	if task, ok := idx.Tasks["1.gate"]; !ok {
		t.Error("1.gate not in index")
	} else if task.Type != TypeGate {
		t.Errorf("1.gate type = %q, want %q", task.Type, TypeGate)
	}

	// Verify gate depends on summary
	if len(idx.Tasks["1.gate"].Dependencies) == 0 || idx.Tasks["1.gate"].Dependencies[0] != "1.summary" {
		t.Errorf("1.gate deps = %v, want [1.summary]", idx.Tasks["1.gate"].Dependencies)
	}
	// Verify summary depends on business tasks
	summaryDeps := idx.Tasks["1.summary"].Dependencies
	if len(summaryDeps) != 2 {
		t.Errorf("1.summary deps count = %d, want 2", len(summaryDeps))
	}
}

func TestBuildIndex_StageGatesWithNoProfiles(t *testing.T) {
	// Stage gates should be generated even when no profiles are configured
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "quick")

	writeTaskMD(t, tasksDir, "1-foo.md", "1.1", "Task 1", nil)
	writeTaskMD(t, tasksDir, "2-bar.md", "1.2", "Task 2", []string{"1.1"})

	opts := BuildIndexOpts{
		FeatureSlug: "test-feature",
		ProjectRoot: projectRoot,
		TasksDir:    tasksDir,
		IndexPath:   indexPath,
	}

	result, err := BuildIndex(opts)
	if err != nil {
		t.Fatalf("BuildIndex error: %v", err)
	}

	// Stage gates generated even without profiles
	if result.StageGatesGenerated != 2 {
		t.Errorf("StageGatesGenerated = %d, want 2", result.StageGatesGenerated)
	}

	// No test tasks generated (no profiles)
	entries, _ := os.ReadDir(tasksDir)
	for _, e := range entries {
		name := e.Name()
		if name == "1.summary.md" || name == "1.gate.md" || name == "1-foo.md" || name == "2-bar.md" || name == "index.json" {
			continue
		}
		t.Errorf("unexpected file %s (test tasks should not be generated without profiles)", name)
	}
}

func TestBuildIndex_StageGatesIdempotent(t *testing.T) {
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "")

	writeTaskMD(t, tasksDir, "1-foo.md", "1.1", "Task 1", nil)
	writeTaskMD(t, tasksDir, "2-bar.md", "1.2", "Task 2", []string{"1.1"})

	opts := BuildIndexOpts{
		FeatureSlug: "test-feature",
		ProjectRoot: projectRoot,
		TasksDir:    tasksDir,
		IndexPath:   indexPath,
	}

	// First build
	result1, err := BuildIndex(opts)
	if err != nil {
		t.Fatalf("first build: %v", err)
	}
	if result1.StageGatesGenerated != 2 {
		t.Errorf("first build StageGatesGenerated = %d, want 2", result1.StageGatesGenerated)
	}

	// Second build -- no new files should be generated
	result2, err := BuildIndex(opts)
	if err != nil {
		t.Fatalf("second build: %v", err)
	}
	if result2.StageGatesGenerated != 0 {
		t.Errorf("second build StageGatesGenerated = %d, want 0 (idempotent)", result2.StageGatesGenerated)
	}
}

func TestBuildIndex_StageGatesTestTaskExclusion(t *testing.T) {
	// T-test/T-quick tasks should not count toward phase threshold
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "")

	// Phase 1: 1 business task + 1 T-test task = only 1 business task (below threshold)
	writeTaskMD(t, tasksDir, "1-foo.md", "1.1", "Business Task", nil)
	writeTaskMD(t, tasksDir, "test-1.md", "T-test-1", "Test Task", nil)

	// Phase 2: 2 business tasks (qualifies)
	writeTaskMD(t, tasksDir, "2-bar.md", "2.1", "Phase 2 Task 1", nil)
	writeTaskMD(t, tasksDir, "3-baz.md", "2.2", "Phase 2 Task 2", []string{"2.1"})

	opts := BuildIndexOpts{
		FeatureSlug: "test-feature",
		ProjectRoot: projectRoot,
		TasksDir:    tasksDir,
		IndexPath:   indexPath,
	}

	result, err := BuildIndex(opts)
	if err != nil {
		t.Fatalf("BuildIndex error: %v", err)
	}

	// Only phase 2 should generate (2 business tasks); phase 1 has only 1 business task
	if result.StageGatesGenerated != 2 {
		t.Errorf("StageGatesGenerated = %d, want 2 (only phase 2)", result.StageGatesGenerated)
	}

	// Phase 1 files should NOT exist
	if _, err := os.Stat(filepath.Join(tasksDir, "1.summary.md")); err == nil {
		t.Error("1.summary.md should not exist (phase 1 has only 1 business task)")
	}
	// Phase 2 files should exist
	if _, err := os.Stat(filepath.Join(tasksDir, "2.summary.md")); os.IsNotExist(err) {
		t.Error("2.summary.md should exist")
	}
}

func TestBuildIndex_StageGatesMultiPhase(t *testing.T) {
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "")

	// 3 phases, each with 2 tasks
	writeTaskMD(t, tasksDir, "1-1.md", "1.1", "P1 T1", nil)
	writeTaskMD(t, tasksDir, "1-2.md", "1.2", "P1 T2", []string{"1.1"})
	writeTaskMD(t, tasksDir, "2-1.md", "2.1", "P2 T1", []string{"1.2"})
	writeTaskMD(t, tasksDir, "2-2.md", "2.2", "P2 T2", []string{"2.1"})
	writeTaskMD(t, tasksDir, "3-1.md", "3.1", "P3 T1", []string{"2.2"})
	writeTaskMD(t, tasksDir, "3-2.md", "3.2", "P3 T2", []string{"3.1"})

	opts := BuildIndexOpts{
		FeatureSlug: "test-feature",
		ProjectRoot: projectRoot,
		TasksDir:    tasksDir,
		IndexPath:   indexPath,
	}

	result, err := BuildIndex(opts)
	if err != nil {
		t.Fatalf("BuildIndex error: %v", err)
	}

	// 3 phases * 2 files = 6
	if result.StageGatesGenerated != 6 {
		t.Errorf("StageGatesGenerated = %d, want 6", result.StageGatesGenerated)
	}

	data, _ := os.ReadFile(indexPath)
	var idx taskIndexJSON
	_ = json.Unmarshal(data, &idx)

	// Verify all 6 generated tasks in index
	for _, key := range []string{"1.summary", "1.gate", "2.summary", "2.gate", "3.summary", "3.gate"} {
		if _, ok := idx.Tasks[key]; !ok {
			t.Errorf("%s not in index", key)
		}
	}
}

// --- Task 2: docs-only detection and conditional pipeline tests ---

func TestNeedsTestPipeline(t *testing.T) {
	tests := []struct {
		name  string
		tasks map[string]Task
		want  bool
	}{
		{
			name: "feature type needs test pipeline",
			tasks: map[string]Task{
				"1-feat": {ID: "1.1", Type: TypeFeature},
				"2-doc":  {ID: "1.2", Type: TypeDocumentation},
			},
			want: true,
		},
		{
			name: "enhancement type needs test pipeline",
			tasks: map[string]Task{
				"1-enh": {ID: "1.1", Type: TypeEnhancement},
			},
			want: true,
		},
		{
			name: "fix type needs test pipeline",
			tasks: map[string]Task{
				"1-fix": {ID: "fix-1", Type: TypeFix},
			},
			want: true,
		},
		{
			name: "documentation-only does NOT need test pipeline",
			tasks: map[string]Task{
				"1-doc": {ID: "1.1", Type: TypeDocumentation},
				"2-doc": {ID: "1.2", Type: TypeDocumentation},
			},
			want: false,
		},
		{
			name: "cleanup-only does NOT need test pipeline",
			tasks: map[string]Task{
				"1-clean": {ID: "1.1", Type: TypeCleanup},
			},
			want: false,
		},
		{
			name: "refactor-only does NOT need test pipeline",
			tasks: map[string]Task{
				"1-ref": {ID: "1.1", Type: TypeRefactor},
			},
			want: false,
		},
		{
			name: "only auto-generated tasks (no business tasks) returns false",
			tasks: map[string]Task{
				"1.summary": {ID: "1.summary", Type: TypeDocGenerationSummary},
				"1.gate":    {ID: "1.gate", Type: TypeGate},
			},
			want: false,
		},
		{
			name: "business tasks with auto-generated tasks mixed in",
			tasks: map[string]Task{
				"1-doc":    {ID: "1.1", Type: TypeDocumentation},
				"2-doc":    {ID: "1.2", Type: TypeDocumentation},
				"1.gate":   {ID: "1.gate", Type: TypeGate},
				"T-test-1": {ID: "T-test-1", Type: TypeTestPipelineGenCases},
			},
			want: false,
		},
		{
			name:  "empty tasks returns false",
			tasks: map[string]Task{},
			want:  false,
		},
		{
			name: "mixed feature and cleanup needs test pipeline",
			tasks: map[string]Task{
				"1-feat":  {ID: "1.1", Type: TypeFeature},
				"2-clean": {ID: "1.2", Type: TypeCleanup},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := needsTestPipeline(tt.tasks)
			if got != tt.want {
				t.Errorf("needsTestPipeline() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNeedsDocEval(t *testing.T) {
	tests := []struct {
		name  string
		tasks map[string]Task
		want  bool
	}{
		{
			name: "documentation-only needs doc eval",
			tasks: map[string]Task{
				"1-doc": {ID: "1.1", Type: TypeDocumentation},
				"2-doc": {ID: "1.2", Type: TypeDocumentation},
			},
			want: true,
		},
		{
			name: "feature type does NOT need doc eval",
			tasks: map[string]Task{
				"1-feat": {ID: "1.1", Type: TypeFeature},
			},
			want: false,
		},
		{
			name: "fix type does NOT need doc eval",
			tasks: map[string]Task{
				"1-fix": {ID: "fix-1", Type: TypeFix},
			},
			want: false,
		},
		{
			name: "cleanup-only does NOT need doc eval",
			tasks: map[string]Task{
				"1-clean": {ID: "1.1", Type: TypeCleanup},
			},
			want: false,
		},
		{
			name: "mixed documentation and feature does NOT need doc eval",
			tasks: map[string]Task{
				"1-doc":  {ID: "1.1", Type: TypeDocumentation},
				"2-feat": {ID: "1.2", Type: TypeFeature},
			},
			want: false,
		},
		{
			name: "doc-evaluation type does NOT need doc eval (not documentation)",
			tasks: map[string]Task{
				"1-doc":  {ID: "1.1", Type: TypeDocumentation},
				"2-eval": {ID: "1.2", Type: TypeDocEvaluation},
			},
			want: false,
		},
		{
			name: "only auto-generated tasks returns false (no business tasks)",
			tasks: map[string]Task{
				"1.summary": {ID: "1.summary", Type: TypeDocGenerationSummary},
				"1.gate":    {ID: "1.gate", Type: TypeGate},
			},
			want: false,
		},
		{
			name:  "empty tasks returns false",
			tasks: map[string]Task{},
			want:  false,
		},
		{
			name: "auto-generated tasks mixed with documentation still returns true",
			tasks: map[string]Task{
				"1-doc":    {ID: "1.1", Type: TypeDocumentation},
				"1.gate":   {ID: "1.gate", Type: TypeGate},
				"T-test-1": {ID: "T-test-1", Type: TypeTestPipelineGenCases},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := needsDocEval(tt.tasks)
			if got != tt.want {
				t.Errorf("needsDocEval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildIndex_MissingTypeHardError(t *testing.T) {
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "quick")

	// Write a task .md without type field (and InferType returns "" for plain numeric IDs)
	// Must use raw content since writeTaskMD now auto-sets type.
	content := "---\nid: \"1\"\ntitle: \"Foo Task\"\npriority: \"P1\"\nestimated_time: \"1h\"\nscope: \"all\"\n---\n\n# Foo Task\n"
	if err := os.WriteFile(filepath.Join(tasksDir, "1-foo.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	opts := BuildIndexOpts{
		FeatureSlug:      "test-feature",
		ProjectRoot:      projectRoot,
		TasksDir:         tasksDir,
		IndexPath:        indexPath,
		TestProfiles:     []string{"go-test"},
		TestCapabilities: []string{"cli"},
	}

	_, err := BuildIndex(opts)
	if err == nil {
		t.Fatal("expected error for missing type, got nil")
	}
	if !strings.Contains(err.Error(), "1-foo.md") {
		t.Errorf("error should name the file, got: %v", err)
	}
	if !strings.Contains(err.Error(), "type") {
		t.Errorf("error should mention type, got: %v", err)
	}
}

func TestBuildIndex_DocsOnlySkipsGatesAndTests(t *testing.T) {
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "quick")

	// Create 2 documentation tasks in same phase (would normally trigger gate generation)
	writeTaskMDWithType(t, tasksDir, "1-doc.md", "1.1", "Doc Task 1", TypeDocumentation, nil)
	writeTaskMDWithType(t, tasksDir, "2-doc.md", "1.2", "Doc Task 2", TypeDocumentation, []string{"1.1"})

	opts := BuildIndexOpts{
		FeatureSlug:      "test-feature",
		ProjectRoot:      projectRoot,
		TasksDir:         tasksDir,
		IndexPath:        indexPath,
		TestCapabilities: []string{"cli"},
		TestProfiles:     []string{"go-test"},
	}

	result, err := BuildIndex(opts)
	if err != nil {
		t.Fatalf("BuildIndex error: %v", err)
	}

	// Stage gates should NOT be generated for docs-only
	if result.StageGatesGenerated != 0 {
		t.Errorf("StageGatesGenerated = %d, want 0 (docs-only)", result.StageGatesGenerated)
	}

	// No gate or summary files should exist
	if _, err := os.Stat(filepath.Join(tasksDir, "1.summary.md")); err == nil {
		t.Error("1.summary.md should NOT exist for docs-only feature")
	}
	if _, err := os.Stat(filepath.Join(tasksDir, "1.gate.md")); err == nil {
		t.Error("1.gate.md should NOT exist for docs-only feature")
	}

	// No test task files should exist
	entries, _ := os.ReadDir(tasksDir)
	for _, e := range entries {
		name := e.Name()
		if name == "1-doc.md" || name == "2-doc.md" || name == "index.json" || name == "eval-doc.md" {
			continue
		}
		t.Errorf("unexpected file %s (docs-only should not generate gates or tests)", name)
	}
}

func TestBuildIndex_DocsOnlyGeneratesEvalDoc(t *testing.T) {
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "quick")

	writeTaskMDWithType(t, tasksDir, "1-doc.md", "1.1", "Doc Task 1", TypeDocumentation, nil)
	writeTaskMDWithType(t, tasksDir, "2-doc.md", "1.2", "Doc Task 2", TypeDocumentation, []string{"1.1"})

	opts := BuildIndexOpts{
		FeatureSlug:      "test-feature",
		ProjectRoot:      projectRoot,
		TasksDir:         tasksDir,
		TestCapabilities: []string{"cli"},
		IndexPath:        indexPath,
		TestProfiles:     []string{"go-test"},
	}

	result, err := BuildIndex(opts)
	if err != nil {
		t.Fatalf("BuildIndex error: %v", err)
	}

	// eval-doc.md should have been generated
	if _, err := os.Stat(filepath.Join(tasksDir, "eval-doc.md")); os.IsNotExist(err) {
		t.Error("eval-doc.md not generated for docs-only feature")
	}

	// Verify T-eval-doc is in the index
	data, _ := os.ReadFile(indexPath)
	var idx taskIndexJSON
	_ = json.Unmarshal(data, &idx)

	evalTask, ok := idx.Tasks["eval-doc"]
	if !ok {
		t.Fatal("eval-doc not in index")
	}
	if evalTask.ID != "T-eval-doc" {
		t.Errorf("eval-doc ID = %q, want T-eval-doc", evalTask.ID)
	}
	if evalTask.Type != TypeDocEvaluation {
		t.Errorf("eval-doc type = %q, want %q", evalTask.Type, TypeDocEvaluation)
	}
	if !evalTask.NoTest {
		t.Error("eval-doc should have noTest=true")
	}
	// Should depend on last business task
	if len(evalTask.Dependencies) == 0 {
		t.Error("eval-doc has no dependencies")
	} else {
		lastDep := evalTask.Dependencies[len(evalTask.Dependencies)-1]
		if lastDep != "1.2" {
			t.Errorf("eval-doc last dep = %q, want 1.2", lastDep)
		}
	}

	// Count: 2 business + 1 eval-doc = 3
	total := result.NewCount + result.UpdatedCount
	if total != 3 {
		t.Errorf("total tasks = %d (new=%d, updated=%d), want 3", total, result.NewCount, result.UpdatedCount)
	}
}

func TestBuildIndex_CodeFeatureUnchanged(t *testing.T) {
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "quick")

	// Feature with feature-type tasks should generate test pipeline and stage gates
	writeTaskMDWithType(t, tasksDir, "1-feat.md", "1.1", "Feature Task 1", TypeFeature, nil)
	writeTaskMDWithType(t, tasksDir, "2-feat.md", "1.2", "Feature Task 2", TypeFeature, []string{"1.1"})

	opts := BuildIndexOpts{
		FeatureSlug:      "test-feature",
		ProjectRoot:      projectRoot,
		TestCapabilities: []string{"cli"},
		TasksDir:         tasksDir,
		IndexPath:        indexPath,
		TestProfiles:     []string{"go-test"},
	}

	result, err := BuildIndex(opts)
	if err != nil {
		t.Fatalf("BuildIndex error: %v", err)
	}

	// Stage gates SHOULD be generated for code feature
	if result.StageGatesGenerated != 2 {
		t.Errorf("StageGatesGenerated = %d, want 2 (code feature)", result.StageGatesGenerated)
	}

	// Test tasks should be generated (quick mode with 1 profile = 4+1 = 5 test tasks)
	data, _ := os.ReadFile(indexPath)
	var idx taskIndexJSON
	_ = json.Unmarshal(data, &idx)

	if _, ok := idx.Tasks["quick-test-cases-go-test"]; !ok {
		t.Error("quick-test-cases-go-test missing (code feature should have test tasks)")
	}

	// eval-doc should NOT be generated for code feature
	if _, ok := idx.Tasks["eval-doc"]; ok {
		t.Error("eval-doc should NOT exist for code feature")
	}
}

func TestBuildIndex_MissingTypeAllowedForAutoGenTasks(t *testing.T) {
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "quick")

	// Create a feature task (has type) and a gate task file without explicit type
	// (gates use InferType)
	writeTaskMDWithType(t, tasksDir, "1-feat.md", "1.1", "Feature Task", TypeFeature, nil)

	// Write a gate file without type in frontmatter - should be OK since InferType handles it
	gateContent := "---\nid: \"1.gate\"\ntitle: \"Phase 1 Gate\"\npriority: \"P0\"\nestimated_time: \"1h\"\nscope: \"all\"\n---\n\n# Gate\n"
	if err := os.WriteFile(filepath.Join(tasksDir, "1.gate.md"), []byte(gateContent), 0644); err != nil {
		t.Fatal(err)
	}

	opts := BuildIndexOpts{
		FeatureSlug:      "test-feature",
		TestCapabilities: []string{"cli"},
		ProjectRoot:      projectRoot,
		TasksDir:         tasksDir,
		IndexPath:        indexPath,
		TestProfiles:     []string{"go-test"},
	}

	// Should NOT error - gate is an auto-generated task, InferType handles it
	_, err := BuildIndex(opts)
	if err != nil {
		t.Fatalf("BuildIndex should not error for auto-gen task without type: %v", err)
	}
}

func TestBuildIndex_EmptyTestCapabilities_NoTestTasks(t *testing.T) {
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "breakdown")

	writeTaskMDWithType(t, tasksDir, "1-feat.md", "1.1", "Feature Task", TypeFeature, nil)
	writeTaskMD(t, tasksDir, "1-gate.md", "1.gate", "Phase 1 Gate", nil)
	writeTaskMD(t, tasksDir, "2-gate.md", "2.gate", "Phase 2 Gate", nil)

	opts := BuildIndexOpts{
		FeatureSlug:      "test-feature",
		ProjectRoot:      projectRoot,
		TasksDir:         tasksDir,
		IndexPath:        indexPath,
		TestProfiles:     []string{"go-test"},
		TestCapabilities: []string{},
	}

	result, err := BuildIndex(opts)
	if err != nil {
		t.Fatalf("BuildIndex: %v", err)
	}

	// Only business + gate tasks, no test pipeline tasks
	total := result.NewCount + result.UpdatedCount
	if total != 3 {
		t.Errorf("total tasks = %d, want 3 (no test tasks with empty capabilities)", total)
	}

	// No test task .md files generated
	entries, _ := os.ReadDir(tasksDir)
	for _, e := range entries {
		name := e.Name()
		if name == "1-feat.md" || name == "1-gate.md" || name == "2-gate.md" || name == "index.json" {
			continue
		}
		t.Errorf("unexpected file %s (no test tasks expected with empty capabilities)", name)
	}
}

func TestBuildIndex_WithCapabilities_ProducesPerTypeTasks(t *testing.T) {
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "quick")

	writeTaskMD(t, tasksDir, "1-foo.md", "1", "Foo Task", nil)

	opts := BuildIndexOpts{
		FeatureSlug:      "test-feature",
		ProjectRoot:      projectRoot,
		TasksDir:         tasksDir,
		IndexPath:        indexPath,
		TestProfiles:     []string{"go-test"},
		TestCapabilities: []string{"cli"},
	}

	result, err := BuildIndex(opts)
	if err != nil {
		t.Fatalf("BuildIndex: %v", err)
	}

	// 1 business + 5 test tasks (gen-cases, gen-and-run-cli, graduate, verify-regression, drift) = 6
	total := result.NewCount + result.UpdatedCount
	if total != 6 {
		t.Errorf("total = %d (new=%d, updated=%d), want 6", total, result.NewCount, result.UpdatedCount)
	}

	// Verify per-type task .md exists
	if _, err := os.Stat(filepath.Join(tasksDir, "quick-gen-and-run-go-test-cli.md")); os.IsNotExist(err) {
		t.Error("quick-gen-and-run-go-test-cli.md not generated")
	}

	// Verify per-type task in index with correct type
	data, _ := os.ReadFile(indexPath)
	var idx taskIndexJSON
	_ = json.Unmarshal(data, &idx)

	task, ok := idx.Tasks["quick-gen-and-run-go-test-cli"]
	if !ok {
		t.Fatal("quick-gen-and-run-go-test-cli not in index")
	}
	if task.ID != "T-quick-2-cli" {
		t.Errorf("per-type task ID = %q, want T-quick-2-cli", task.ID)
	}
	if task.Profile != "go-test" {
		t.Errorf("per-type task profile = %q, want go-test", task.Profile)
	}
}

func TestBuildIndex_DeterministicOutput(t *testing.T) {
	// BuildIndex produces identical output regardless of whether test-cases.md exists
	projectRoot, tasksDir, indexPath := setupBuildEnv(t, "quick")

	writeTaskMD(t, tasksDir, "1-foo.md", "1", "Foo Task", nil)

	opts := BuildIndexOpts{
		FeatureSlug:      "test-feature",
		ProjectRoot:      projectRoot,
		TasksDir:         tasksDir,
		IndexPath:        indexPath,
		TestProfiles:     []string{"go-test"},
		TestCapabilities: []string{"api"},
	}

	_, err := BuildIndex(opts)
	if err != nil {
		t.Fatalf("first build: %v", err)
	}

	data1, _ := os.ReadFile(indexPath)

	// Create a test-cases.md file in the feature's testing dir (should NOT affect output)
	testingDir := filepath.Join(projectRoot, "docs", "features", "test-feature", "testing")
	if err := os.MkdirAll(testingDir, 0755); err != nil {
		t.Fatal(err)
	}
	testCasesContent := `## Summary

| Type | Count |
|------|-------|
| API  | 5     |
| CLI  | 3     |
| **Total** | **8** |`
	if err := os.WriteFile(filepath.Join(testingDir, "test-cases.md"), []byte(testCasesContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Rebuild
	_, err = BuildIndex(opts)
	if err != nil {
		t.Fatalf("second build: %v", err)
	}

	data2, _ := os.ReadFile(indexPath)

	// Output must be identical
	if string(data1) != string(data2) {
		t.Errorf("BuildIndex output differs with/without test-cases.md\nfirst=%s\nsecond=%s", string(data1), string(data2))
	}

	// The key invariant is that output JSON is identical (verified above).
}
