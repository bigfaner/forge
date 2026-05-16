package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/task"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// captureOutput captures stdout and stderr during a function execution
func captureOutput(f func() error) (string, error) {
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	rOut, wOut, err := os.Pipe()
	if err != nil {
		return "", err
	}
	rErr, wErr, err := os.Pipe()
	if err != nil {
		return "", err
	}

	os.Stdout = wOut
	os.Stderr = wErr

	outCh := make(chan string)
	errCh := make(chan string)

	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, rOut)
		outCh <- buf.String()
	}()

	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, rErr)
		errCh <- buf.String()
	}()

	runErr := f()

	_ = wOut.Close()
	_ = wErr.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	stdout := <-outCh
	stderr := <-errCh

	return stdout + stderr, runErr
}

func TestRunFeature_Display(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	goMod := filepath.Join(dir, "go.mod")
	if err := os.WriteFile(goMod, []byte("module test-project\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create feature with proper structure
	if err := feature.EnsureFeatureDir(dir, "test-feature"); err != nil {
		t.Fatal(err)
	}

	// Create index.json
	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test-feature"))
	index := &task.TaskIndex{Feature: "test-feature"}
	indexData, _ := json.Marshal(index)
	if err := os.WriteFile(indexPath, indexData, 0644); err != nil {
		t.Fatal(err)
	}

	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origWd) }()

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"feature"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("feature command failed: %v", err)
	}

	if !strings.Contains(output, "test-feature") {
		t.Errorf("expected output to contain 'test-feature', got %q", output)
	}
}

func TestRunFeature_Set(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	goMod := filepath.Join(dir, "go.mod")
	if err := os.WriteFile(goMod, []byte("module test-project\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origWd) }()

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	_, err = captureOutput(func() error {
		rootCmd.SetArgs([]string{"feature", "my-new-feature"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("feature set command failed: %v", err)
	}

	// Verify feature directory structure was created
	featureProcessDir := filepath.Join(dir, feature.FeaturesDir, "my-new-feature", feature.TasksDirName, feature.ProcessDirName)
	if _, err := os.Stat(featureProcessDir); os.IsNotExist(err) {
		t.Errorf("feature process directory %s was not created", featureProcessDir)
	}
}

func TestRunFeature_NoFeatureSet(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	goMod := filepath.Join(dir, "go.mod")
	if err := os.WriteFile(goMod, []byte("module test-project\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create features directory but no feature subdirectory
	featuresDir := filepath.Join(dir, feature.FeaturesDir)
	if err := os.MkdirAll(featuresDir, 0755); err != nil {
		t.Fatal(err)
	}

	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origWd) }()

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	output, _ := captureOutput(func() error {
		rootCmd.SetArgs([]string{"feature"})
		return rootCmd.Execute()
	})

	if !strings.Contains(output, "FEATURE: (none)") {
		t.Errorf("expected output to contain 'FEATURE: (none)', got %q", output)
	}
}

func TestFeatureList_WithFeatures(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test-project\n"), 0644))

	// Create two features with manifests
	for _, slug := range []string{"feature-a", "feature-b"} {
		featureDir := filepath.Join(dir, feature.FeaturesDir, slug)
		require.NoError(t, os.MkdirAll(featureDir, 0755))

		manifestContent := fmt.Sprintf("---\nfeature: %s\nstatus: in-progress\n---\n", slug)
		require.NoError(t, os.WriteFile(filepath.Join(featureDir, feature.ManifestFileName), []byte(manifestContent), 0644))

		// Create task index
		tasksDir := filepath.Join(featureDir, feature.TasksDirName)
		require.NoError(t, os.MkdirAll(tasksDir, 0755))

		index := &task.TaskIndex{
			Feature:    slug,
			StatusEnum: []string{"pending", "in_progress", "completed"},
		}
		index.SetTasks(map[string]task.Task{
			"t1": {ID: "1", Title: "Task 1", Status: "completed"},
			"t2": {ID: "2", Title: "Task 2", Status: "pending"},
		})
		indexData, err := json.Marshal(index)
		require.NoError(t, err)
		require.NoError(t, os.WriteFile(filepath.Join(tasksDir, feature.IndexFileName), indexData, 0644))
	}

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"feature", "list"})
		return rootCmd.Execute()
	})
	require.NoError(t, err)
	assert.Contains(t, output, "FEATURES")
	assert.Contains(t, output, "feature-a")
	assert.Contains(t, output, "feature-b")
	assert.Contains(t, output, "1/2") // progress
}

func TestFeatureList_Empty(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test-project\n"), 0644))

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"feature", "list"})
		return rootCmd.Execute()
	})
	// No features prints to stderr
	_ = output
	_ = err
}

func TestFeatureStatus_Found(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test-project\n"), 0644))

	slug := "test-feature"
	featureDir := filepath.Join(dir, feature.FeaturesDir, slug)
	require.NoError(t, os.MkdirAll(featureDir, 0755))

	// Create manifest
	manifestContent := fmt.Sprintf("---\nfeature: %s\nstatus: tasks\n---\n", slug)
	require.NoError(t, os.WriteFile(filepath.Join(featureDir, feature.ManifestFileName), []byte(manifestContent), 0644))

	// Create task index
	tasksDir := filepath.Join(featureDir, feature.TasksDirName)
	require.NoError(t, os.MkdirAll(tasksDir, 0755))

	index := &task.TaskIndex{
		Feature:    slug,
		StatusEnum: []string{"pending", "in_progress", "completed", "blocked"},
	}
	index.SetTasks(map[string]task.Task{
		"t1": {ID: "1", Title: "Task 1", Status: "completed"},
		"t2": {ID: "2", Title: "Task 2", Status: "in_progress"},
		"t3": {ID: "3", Title: "Task 3", Status: "pending"},
	})
	indexData, err := json.Marshal(index)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, feature.IndexFileName), indexData, 0644))

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"feature", "status", slug})
		return rootCmd.Execute()
	})
	require.NoError(t, err)
	assert.Contains(t, output, "SLUG: test-feature")
	assert.Contains(t, output, "STATUS: tasks")
	assert.Contains(t, output, "[TASKS]")
	assert.Contains(t, output, "completed: 1")
	assert.Contains(t, output, "in_progress: 1")
	assert.Contains(t, output, "pending: 1")
	assert.Contains(t, output, "TOTAL: 3")
	assert.Contains(t, output, "[ARTIFACTS]")
}

func TestFeatureStatus_WithScores(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test-project\n"), 0644))

	slug := "scored-feature"
	featureDir := filepath.Join(dir, feature.FeaturesDir, slug)
	require.NoError(t, os.MkdirAll(featureDir, 0755))

	manifestContent := fmt.Sprintf("---\nfeature: %s\nstatus: completed\n---\n", slug)
	require.NoError(t, os.WriteFile(filepath.Join(featureDir, feature.ManifestFileName), []byte(manifestContent), 0644))

	// Create task index
	tasksDir := filepath.Join(featureDir, feature.TasksDirName)
	require.NoError(t, os.MkdirAll(tasksDir, 0755))

	index := &task.TaskIndex{Feature: slug, StatusEnum: []string{"completed"}}
	index.SetTasks(map[string]task.Task{
		"t1": {ID: "1", Title: "Done", Status: "completed"},
	})
	indexData, _ := json.Marshal(index)
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, feature.IndexFileName), indexData, 0644))

	// Create PRD with score
	prdDir := filepath.Join(featureDir, feature.PRDDirName)
	require.NoError(t, os.MkdirAll(prdDir, 0755))
	prdContent := "---\nscore: 850\n---\n# PRD\n"
	require.NoError(t, os.WriteFile(filepath.Join(prdDir, feature.PRDSpecFile), []byte(prdContent), 0644))

	// Create design with score
	designDir := filepath.Join(featureDir, feature.DesignDirName)
	require.NoError(t, os.MkdirAll(designDir, 0755))
	designContent := "---\nscore: 920\n---\n# Design\n"
	require.NoError(t, os.WriteFile(filepath.Join(designDir, feature.TechDesignFile), []byte(designContent), 0644))

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"feature", "status", slug})
		return rootCmd.Execute()
	})
	require.NoError(t, err)
	assert.Contains(t, output, "PRD: 850")
	assert.Contains(t, output, "DESIGN: 920")
	assert.Contains(t, output, "UI: —") // no UI design, should show em-dash
}

func TestFeatureList_SortedByManifestMtime(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test-project\n"), 0644))

	// Create features with different manifest mtimes
	// We control mtime via os.Chtimes to guarantee ordering regardless of filesystem timing

	type featureSpec struct {
		slug  string
		mtime time.Time
	}
	specs := []featureSpec{
		{slug: "oldest-feature", mtime: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)},
		{slug: "middle-feature", mtime: time.Date(2026, 3, 15, 12, 0, 0, 0, time.UTC)},
		{slug: "newest-feature", mtime: time.Date(2026, 5, 16, 10, 0, 0, 0, time.UTC)},
	}

	for _, spec := range specs {
		featureDir := filepath.Join(dir, feature.FeaturesDir, spec.slug)
		require.NoError(t, os.MkdirAll(featureDir, 0755))

		manifestContent := fmt.Sprintf("---\nfeature: %s\nstatus: in-progress\n---\n", spec.slug)
		manifestPath := filepath.Join(featureDir, feature.ManifestFileName)
		require.NoError(t, os.WriteFile(manifestPath, []byte(manifestContent), 0644))

		// Set manifest mtime explicitly
		require.NoError(t, os.Chtimes(manifestPath, spec.mtime, spec.mtime))

		// Create task index
		tasksDir := filepath.Join(featureDir, feature.TasksDirName)
		require.NoError(t, os.MkdirAll(tasksDir, 0755))
		index := &task.TaskIndex{Feature: spec.slug}
		indexData, err := json.Marshal(index)
		require.NoError(t, err)
		require.NoError(t, os.WriteFile(filepath.Join(tasksDir, feature.IndexFileName), indexData, 0644))
	}

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"feature", "list"})
		return rootCmd.Execute()
	})
	require.NoError(t, err)

	// Verify features appear in reverse chronological order (newest first)
	newestIdx := strings.Index(output, "newest-feature")
	middleIdx := strings.Index(output, "middle-feature")
	oldestIdx := strings.Index(output, "oldest-feature")

	assert.True(t, newestIdx < middleIdx, "newest-feature should appear before middle-feature")
	assert.True(t, middleIdx < oldestIdx, "middle-feature should appear before oldest-feature")
}

func TestFeatureList_MissingManifestSortsToEnd(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test-project\n"), 0644))

	// Create feature with a manifest
	featureDir := filepath.Join(dir, feature.FeaturesDir, "with-manifest")
	require.NoError(t, os.MkdirAll(featureDir, 0755))
	manifestContent := "---\nfeature: with-manifest\nstatus: in-progress\n---\n"
	manifestPath := filepath.Join(featureDir, feature.ManifestFileName)
	require.NoError(t, os.WriteFile(manifestPath, []byte(manifestContent), 0644))
	require.NoError(t, os.Chtimes(manifestPath, time.Date(2026, 5, 16, 10, 0, 0, 0, time.UTC), time.Date(2026, 5, 16, 10, 0, 0, 0, time.UTC)))

	tasksDir := filepath.Join(featureDir, feature.TasksDirName)
	require.NoError(t, os.MkdirAll(tasksDir, 0755))
	index := &task.TaskIndex{Feature: "with-manifest"}
	indexData, _ := json.Marshal(index)
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, feature.IndexFileName), indexData, 0644))

	// Create feature WITHOUT a manifest (no manifest.md file)
	featureDir2 := filepath.Join(dir, feature.FeaturesDir, "no-manifest")
	require.NoError(t, os.MkdirAll(featureDir2, 0755))
	tasksDir2 := filepath.Join(featureDir2, feature.TasksDirName)
	require.NoError(t, os.MkdirAll(tasksDir2, 0755))
	index2 := &task.TaskIndex{Feature: "no-manifest"}
	indexData2, _ := json.Marshal(index2)
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir2, feature.IndexFileName), indexData2, 0644))

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"feature", "list"})
		return rootCmd.Execute()
	})
	require.NoError(t, err)

	// Feature with manifest should appear before feature without
	withIdx := strings.Index(output, "with-manifest")
	noIdx := strings.Index(output, "no-manifest")

	assert.True(t, withIdx < noIdx, "feature with manifest should appear before feature without manifest")
}

func TestScoreDisplay(t *testing.T) {
	assert.Equal(t, "850", scoreDisplay("850"))
	assert.Equal(t, "—", scoreDisplay(""))
}

func TestRunQuery(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	goMod := filepath.Join(dir, "go.mod")
	if err := os.WriteFile(goMod, []byte("module test-project\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	featureDir := filepath.Join(dir, "docs", "features", "test-feature")
	tasksDir := filepath.Join(featureDir, "tasks")
	indexPath := filepath.Join(tasksDir, "index.json")

	if err := os.MkdirAll(tasksDir, 0755); err != nil {
		t.Fatal(err)
	}

	index := &task.TaskIndex{
		Feature:      "test-feature",
		PRD:          "prd/prd-spec.md",
		Design:       "design/tech-design.md",
		StatusEnum:   []string{"pending", "in_progress", "completed"},
		PriorityEnum: []string{"P0", "P1", "P2"},
	}
	index.SetTasks(map[string]task.Task{
		"task1": {ID: "1.1", Title: "Task 1", Status: "pending", Priority: "P0", File: "1.1.md", Record: "1.1.md", EstimatedTime: "30m", Dependencies: []string{"1.0"}},
	})

	if err := task.SaveIndex(indexPath, index); err != nil {
		t.Fatal(err)
	}

	// Ensure feature directory structure exists
	if err := feature.EnsureFeatureDir(dir, "test-feature"); err != nil {
		t.Fatal(err)
	}

	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origWd) }()

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"task", "query", "1.1"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("query command failed: %v", err)
	}

	if !strings.Contains(output, "TASK_ID: 1.1") {
		t.Errorf("expected output to contain 'TASK_ID: 1.1', got %q", output)
	}
	if !strings.Contains(output, "STATUS: pending") {
		t.Errorf("expected output to contain 'STATUS: pending', got %q", output)
	}
	// Removed fields must NOT appear
	if strings.Contains(output, "KEY:") {
		t.Errorf("KEY should not appear in query output, got %q", output)
	}
	if strings.Contains(output, "TITLE:") {
		t.Errorf("TITLE should not appear in query output, got %q", output)
	}
	if strings.Contains(output, "DEPENDENCIES:") {
		t.Errorf("DEPENDENCIES should not appear in query output, got %q", output)
	}
}

func TestRunStatus(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	goMod := filepath.Join(dir, "go.mod")
	if err := os.WriteFile(goMod, []byte("module test-project\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	featureDir := filepath.Join(dir, "docs", "features", "test-feature")
	tasksDir := filepath.Join(featureDir, "tasks")
	indexPath := filepath.Join(tasksDir, "index.json")

	if err := os.MkdirAll(tasksDir, 0755); err != nil {
		t.Fatal(err)
	}

	index := &task.TaskIndex{
		Feature:      "test-feature",
		PRD:          "prd/prd-spec.md",
		Design:       "design/tech-design.md",
		StatusEnum:   []string{"pending", "in_progress", "completed", "blocked", "skipped"},
		PriorityEnum: []string{"P0", "P1", "P2"},
	}
	index.SetTasks(map[string]task.Task{
		"task1": {ID: "1.1", Title: "Task 1", Status: "pending", Priority: "P0", File: "1.1.md", Record: "1.1.md"},
	})

	if err := task.SaveIndex(indexPath, index); err != nil {
		t.Fatal(err)
	}

	// Ensure feature directory structure exists
	if err := feature.EnsureFeatureDir(dir, "test-feature"); err != nil {
		t.Fatal(err)
	}

	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origWd) }()

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	_, err = captureOutput(func() error {
		rootCmd.SetArgs([]string{"task", "status", "1.1", "blocked"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("status command failed: %v", err)
	}

	updatedIndex, err := task.LoadIndex(indexPath)
	if err != nil {
		t.Fatalf("failed to load updated index: %v", err)
	}
	if updatedIndex.TasksMap()["task1"].Status != "blocked" {
		t.Errorf("expected status 'blocked', got %q", updatedIndex.TasksMap()["task1"].Status)
	}
}

func TestRunCheck(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	goMod := filepath.Join(dir, "go.mod")
	if err := os.WriteFile(goMod, []byte("module test-project\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	featureDir := filepath.Join(dir, "docs", "features", "test-feature")
	tasksDir := filepath.Join(featureDir, "tasks")
	indexPath := filepath.Join(tasksDir, "index.json")

	if err := os.MkdirAll(tasksDir, 0755); err != nil {
		t.Fatal(err)
	}

	index := &task.TaskIndex{
		Feature:      "test-feature",
		PRD:          "prd/prd-spec.md",
		Design:       "design/tech-design.md",
		StatusEnum:   []string{"pending", "in_progress", "completed"},
		PriorityEnum: []string{"P0", "P1", "P2"},
	}
	index.SetTasks(map[string]task.Task{
		"task1": {ID: "1.1", Title: "Task 1", Status: "pending", Priority: "P0", File: "1.1.md", Record: "1.1.md"},
		"task2": {ID: "1.2", Title: "Task 2", Status: "pending", Priority: "P1", File: "1.2.md", Record: "1.2.md", Dependencies: []string{"1.1"}},
	})

	if err := task.SaveIndex(indexPath, index); err != nil {
		t.Fatal(err)
	}

	// Ensure feature directory structure exists
	if err := feature.EnsureFeatureDir(dir, "test-feature"); err != nil {
		t.Fatal(err)
	}

	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origWd) }()

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"task", "check-deps"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("check command failed: %v", err)
	}

	if !strings.Contains(output, "[TASKS]") {
		t.Errorf("expected output to contain '[TASKS]', got %q", output)
	}
	if !strings.Contains(output, "[DEPENDENCIES]") {
		t.Errorf("expected output to contain '[DEPENDENCIES]', got %q", output)
	}
	if !strings.Contains(output, "RESULT: PASS") {
		t.Errorf("expected output to contain 'RESULT: PASS', got %q", output)
	}
}
