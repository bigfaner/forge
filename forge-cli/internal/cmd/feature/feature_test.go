package feature

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
	verbose = false
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
		Cmd.SetArgs([]string{})
		return Cmd.Execute()
	})
	if err != nil {
		t.Fatalf("feature command failed: %v", err)
	}

	if !strings.Contains(output, "test-feature") {
		t.Errorf("expected output to contain 'test-feature', got %q", output)
	}
}

func TestRunFeature_Set(t *testing.T) {
	verbose = false
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
		Cmd.SetArgs([]string{"my-new-feature"})
		return Cmd.Execute()
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
	verbose = false
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
		Cmd.SetArgs([]string{})
		return Cmd.Execute()
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
		Cmd.SetArgs([]string{"list"})
		return Cmd.Execute()
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
		Cmd.SetArgs([]string{"list"})
		return Cmd.Execute()
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
		Cmd.SetArgs([]string{"status", slug})
		return Cmd.Execute()
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
		Cmd.SetArgs([]string{"status", slug})
		return Cmd.Execute()
	})
	require.NoError(t, err)
	assert.Contains(t, output, "PRD: 850")
	assert.Contains(t, output, "DESIGN: 920")
	assert.Contains(t, output, "UI: —") // no UI design, should show em-dash
}

func TestFeatureList_SortedByCreatedDescending(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test-project\n"), 0644))

	// Create features with different created dates in frontmatter.
	// Lexicographic order differs from date order to verify proper sorting.
	type featureSpec struct {
		slug    string
		created string
	}
	specs := []featureSpec{
		{slug: "alpha-feature", created: "2026-01-15"},
		{slug: "beta-feature", created: "2026-03-10"},
		{slug: "gamma-feature", created: "2026-02-01"},
	}

	for _, spec := range specs {
		featureDir := filepath.Join(dir, feature.FeaturesDir, spec.slug)
		require.NoError(t, os.MkdirAll(featureDir, 0755))

		manifestContent := fmt.Sprintf("---\nfeature: %s\nstatus: in-progress\ncreated: %s\n---\n", spec.slug, spec.created)
		manifestPath := filepath.Join(featureDir, feature.ManifestFileName)
		require.NoError(t, os.WriteFile(manifestPath, []byte(manifestContent), 0644))

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
		Cmd.SetArgs([]string{"list"})
		return Cmd.Execute()
	})
	require.NoError(t, err)

	// Verify newest first: beta (Mar 10) > gamma (Feb 1) > alpha (Jan 15)
	betaIdx := strings.Index(output, "beta-feature")
	gammaIdx := strings.Index(output, "gamma-feature")
	alphaIdx := strings.Index(output, "alpha-feature")

	assert.True(t, betaIdx < gammaIdx, "beta-feature (Mar) should appear before gamma-feature (Feb)")
	assert.True(t, gammaIdx < alphaIdx, "gamma-feature (Feb) should appear before alpha-feature (Jan)")
}

func TestFeatureList_SortedByMtimeFallback(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test-project\n"), 0644))

	// Create features WITHOUT created field — should fall back to mtime.
	type featureSpec struct {
		slug  string
		mtime time.Time
	}
	specs := []featureSpec{
		{slug: "old-no-created", mtime: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)},
		{slug: "new-no-created", mtime: time.Date(2026, 5, 16, 10, 0, 0, 0, time.UTC)},
	}

	for _, spec := range specs {
		featureDir := filepath.Join(dir, feature.FeaturesDir, spec.slug)
		require.NoError(t, os.MkdirAll(featureDir, 0755))

		manifestContent := fmt.Sprintf("---\nfeature: %s\nstatus: in-progress\n---\n", spec.slug)
		manifestPath := filepath.Join(featureDir, feature.ManifestFileName)
		require.NoError(t, os.WriteFile(manifestPath, []byte(manifestContent), 0644))
		require.NoError(t, os.Chtimes(manifestPath, spec.mtime, spec.mtime))

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
		Cmd.SetArgs([]string{"list"})
		return Cmd.Execute()
	})
	require.NoError(t, err)

	newIdx := strings.Index(output, "new-no-created")
	oldIdx := strings.Index(output, "old-no-created")
	assert.True(t, newIdx < oldIdx, "feature with newer mtime should appear before older mtime when no created field")
}

func TestFeatureList_CreatedTakesPriorityOverMtime(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test-project\n"), 0644))

	// Feature with older mtime but newer created date should sort first.
	featureDir1 := filepath.Join(dir, feature.FeaturesDir, "old-mtime-new-created")
	require.NoError(t, os.MkdirAll(featureDir1, 0755))
	manifest1 := "---\nfeature: old-mtime-new-created\nstatus: in-progress\ncreated: 2026-05-01\n---\n"
	manifestPath1 := filepath.Join(featureDir1, feature.ManifestFileName)
	require.NoError(t, os.WriteFile(manifestPath1, []byte(manifest1), 0644))
	require.NoError(t, os.Chtimes(manifestPath1, time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)))
	tasksDir1 := filepath.Join(featureDir1, feature.TasksDirName)
	require.NoError(t, os.MkdirAll(tasksDir1, 0755))
	index1 := &task.TaskIndex{Feature: "old-mtime-new-created"}
	indexData1, _ := json.Marshal(index1)
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir1, feature.IndexFileName), indexData1, 0644))

	// Feature with newer mtime but older created date should sort second.
	featureDir2 := filepath.Join(dir, feature.FeaturesDir, "new-mtime-old-created")
	require.NoError(t, os.MkdirAll(featureDir2, 0755))
	manifest2 := "---\nfeature: new-mtime-old-created\nstatus: in-progress\ncreated: 2026-01-01\n---\n"
	manifestPath2 := filepath.Join(featureDir2, feature.ManifestFileName)
	require.NoError(t, os.WriteFile(manifestPath2, []byte(manifest2), 0644))
	require.NoError(t, os.Chtimes(manifestPath2, time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)))
	tasksDir2 := filepath.Join(featureDir2, feature.TasksDirName)
	require.NoError(t, os.MkdirAll(tasksDir2, 0755))
	index2 := &task.TaskIndex{Feature: "new-mtime-old-created"}
	indexData2, _ := json.Marshal(index2)
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir2, feature.IndexFileName), indexData2, 0644))

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		Cmd.SetArgs([]string{"list"})
		return Cmd.Execute()
	})
	require.NoError(t, err)

	newCreatedIdx := strings.Index(output, "old-mtime-new-created")
	oldCreatedIdx := strings.Index(output, "new-mtime-old-created")
	assert.True(t, newCreatedIdx < oldCreatedIdx, "feature with newer created date should sort first regardless of mtime")
}

func TestFeatureList_MissingCreatedSortsAfterCreated(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test-project\n"), 0644))

	// Create feature WITH created field
	featureDir := filepath.Join(dir, feature.FeaturesDir, "with-created")
	require.NoError(t, os.MkdirAll(featureDir, 0755))
	manifestContent := "---\nfeature: with-created\nstatus: in-progress\ncreated: 2026-05-16\n---\n"
	manifestPath := filepath.Join(featureDir, feature.ManifestFileName)
	require.NoError(t, os.WriteFile(manifestPath, []byte(manifestContent), 0644))

	tasksDir := filepath.Join(featureDir, feature.TasksDirName)
	require.NoError(t, os.MkdirAll(tasksDir, 0755))
	index := &task.TaskIndex{Feature: "with-created"}
	indexData, _ := json.Marshal(index)
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, feature.IndexFileName), indexData, 0644))

	// Create feature WITHOUT created field (fallback to mtime, which sorts after valid created)
	featureDir2 := filepath.Join(dir, feature.FeaturesDir, "no-created")
	require.NoError(t, os.MkdirAll(featureDir2, 0755))
	manifestContent2 := "---\nfeature: no-created\nstatus: in-progress\n---\n"
	manifestPath2 := filepath.Join(featureDir2, feature.ManifestFileName)
	require.NoError(t, os.WriteFile(manifestPath2, []byte(manifestContent2), 0644))
	require.NoError(t, os.Chtimes(manifestPath2, time.Date(2026, 5, 16, 10, 0, 0, 0, time.UTC), time.Date(2026, 5, 16, 10, 0, 0, 0, time.UTC)))

	tasksDir2 := filepath.Join(featureDir2, feature.TasksDirName)
	require.NoError(t, os.MkdirAll(tasksDir2, 0755))
	index2 := &task.TaskIndex{Feature: "no-created"}
	indexData2, _ := json.Marshal(index2)
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir2, feature.IndexFileName), indexData2, 0644))

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		Cmd.SetArgs([]string{"list"})
		return Cmd.Execute()
	})
	require.NoError(t, err)

	// Feature with created should appear before feature without created (fallback sorts to end)
	withIdx := strings.Index(output, "with-created")
	noIdx := strings.Index(output, "no-created")

	assert.True(t, withIdx < noIdx, "feature with created field should appear before feature without created")
}

func TestFeatureSet_CreatesDirAndState(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test-project\n"), 0644))

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		Cmd.SetArgs([]string{"set", "my-feature"})
		return Cmd.Execute()
	})
	require.NoError(t, err)

	// Verify stdout contains the feature slug
	assert.Contains(t, output, "FEATURE: my-feature")

	// Verify feature directory structure exists
	featureDir := filepath.Join(dir, feature.FeaturesDir, "my-feature")
	info, err := os.Stat(featureDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())

	// Verify .forge/state.json was written with correct values
	statePath := filepath.Join(dir, feature.ForgeDir, feature.ForgeStateFileName)
	data, err := os.ReadFile(statePath)
	require.NoError(t, err)

	var state feature.ForgeState
	require.NoError(t, json.Unmarshal(data, &state))
	assert.Equal(t, "my-feature", state.Feature)
	assert.False(t, state.AllCompleted)
}

func TestFeatureSet_EmptySlug(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test-project\n"), 0644))

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	_, err = captureOutput(func() error {
		Cmd.SetArgs([]string{"set", ""})
		return Cmd.Execute()
	})
	assert.Error(t, err)

	// Verify no state.json was created
	statePath := filepath.Join(dir, feature.ForgeDir, feature.ForgeStateFileName)
	_, statErr := os.Stat(statePath)
	assert.True(t, os.IsNotExist(statErr), "state.json should not be created for empty slug")
}

func TestFeatureSet_BackwardCompat_PositionalArg(t *testing.T) {
	// Verify that the existing `forge feature <slug>` positional arg still works
	// and does NOT write state.json (backward compatible).
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test-project\n"), 0644))

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		Cmd.SetArgs([]string{"legacy-feature"})
		return Cmd.Execute()
	})
	require.NoError(t, err)
	assert.Contains(t, output, "FEATURE: legacy-feature")

	// Verify feature directory was created
	featureDir := filepath.Join(dir, feature.FeaturesDir, "legacy-feature")
	_, err = os.Stat(featureDir)
	require.NoError(t, err)

	// Verify state.json was NOT written (old behavior preserved)
	statePath := filepath.Join(dir, feature.ForgeDir, feature.ForgeStateFileName)
	_, err = os.Stat(statePath)
	assert.True(t, os.IsNotExist(err), "state.json should not exist for positional arg")
}

func TestRunFeature_Verbose_FromStateJSON(t *testing.T) {
	verbose = false
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test-project\n"), 0644))

	// Create feature directory and index.json
	require.NoError(t, feature.EnsureFeatureDir(dir, "state-feature"))
	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("state-feature"))
	indexData, _ := json.Marshal(&task.TaskIndex{Feature: "state-feature"})
	require.NoError(t, os.WriteFile(indexPath, indexData, 0644))

	// Write .forge/state.json to set explicit feature
	require.NoError(t, feature.EnsureForgeState(dir, "state-feature"))

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		Cmd.SetArgs([]string{"-v"})
		return Cmd.Execute()
	})
	require.NoError(t, err)
	assert.Contains(t, output, "FEATURE: state-feature (from: state.json)")
}

func TestRunFeature_Verbose_FromFeaturesDir(t *testing.T) {
	verbose = false
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test-project\n"), 0644))

	// Create single feature directory with index.json
	require.NoError(t, feature.EnsureFeatureDir(dir, "dir-feature"))
	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("dir-feature"))
	indexData, _ := json.Marshal(&task.TaskIndex{Feature: "dir-feature"})
	require.NoError(t, os.WriteFile(indexPath, indexData, 0644))

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		Cmd.SetArgs([]string{"-v"})
		return Cmd.Execute()
	})
	require.NoError(t, err)
	assert.Contains(t, output, "FEATURE: dir-feature (from: features-dir)")
}

func TestRunFeature_Verbose_NoFeature(t *testing.T) {
	verbose = false
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test-project\n"), 0644))

	// Create empty features directory
	require.NoError(t, os.MkdirAll(filepath.Join(dir, feature.FeaturesDir), 0755))

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, _ := captureOutput(func() error {
		Cmd.SetArgs([]string{"-v"})
		return Cmd.Execute()
	})
	assert.Contains(t, output, "FEATURE: (none)")
}

func TestRunFeature_NonVerboseUnchanged(t *testing.T) {
	verbose = false
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test-project\n"), 0644))

	require.NoError(t, feature.EnsureFeatureDir(dir, "plain-feature"))
	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("plain-feature"))
	indexData, _ := json.Marshal(&task.TaskIndex{Feature: "plain-feature"})
	require.NoError(t, os.WriteFile(indexPath, indexData, 0644))

	// Write state to ensure it resolves from state.json
	require.NoError(t, feature.EnsureForgeState(dir, "plain-feature"))

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		Cmd.SetArgs([]string{})
		return Cmd.Execute()
	})
	require.NoError(t, err)
	// Without -v, output should be just the slug (no source info)
	assert.Contains(t, output, "FEATURE: plain-feature")
	assert.NotContains(t, output, "(from:")
}

func TestRunFeature_VerboseFlagNotLeakedToSubcommands(t *testing.T) {
	// Verify -v is a local flag, not persistent — subcommands should not recognize it
	f := Cmd.Flags().Lookup("verbose")
	require.NotNil(t, f, "verbose flag should exist on Cmd")

	p := Cmd.PersistentFlags().Lookup("verbose")
	assert.Nil(t, p, "verbose flag should NOT be a persistent flag")
}

func TestFeatureStatus_CompletedManifest(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test-project\n"), 0644))

	slug := "completed-feature"
	featureDir := filepath.Join(dir, feature.FeaturesDir, slug)
	require.NoError(t, os.MkdirAll(featureDir, 0755))

	manifestContent := fmt.Sprintf("---\nfeature: %s\nstatus: completed\n---\n", slug)
	require.NoError(t, os.WriteFile(filepath.Join(featureDir, feature.ManifestFileName), []byte(manifestContent), 0644))

	// Create task index
	tasksDir := filepath.Join(featureDir, feature.TasksDirName)
	require.NoError(t, os.MkdirAll(tasksDir, 0755))

	index := &task.TaskIndex{
		Feature:    slug,
		StatusEnum: []string{"completed"},
	}
	index.SetTasks(map[string]task.Task{
		"t1": {ID: "1", Title: "Task 1", Status: "completed"},
	})
	indexData, err := json.Marshal(index)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, feature.IndexFileName), indexData, 0644))

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		Cmd.SetArgs([]string{"status", slug})
		return Cmd.Execute()
	})
	require.NoError(t, err)
	assert.Contains(t, output, "SLUG: completed-feature")
	assert.Contains(t, output, "STATUS: completed")
	assert.Contains(t, output, "completed: 1")
	assert.Contains(t, output, "TOTAL: 1")
}

func TestFeatureList_ApprovedStatus(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test-project\n"), 0644))

	slug := "approved-feature"
	featureDir := filepath.Join(dir, feature.FeaturesDir, slug)
	require.NoError(t, os.MkdirAll(featureDir, 0755))

	manifestContent := fmt.Sprintf("---\nfeature: %s\nstatus: approved\n---\n", slug)
	require.NoError(t, os.WriteFile(filepath.Join(featureDir, feature.ManifestFileName), []byte(manifestContent), 0644))

	// Create task index
	tasksDir := filepath.Join(featureDir, feature.TasksDirName)
	require.NoError(t, os.MkdirAll(tasksDir, 0755))

	index := &task.TaskIndex{Feature: slug}
	indexData, err := json.Marshal(index)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, feature.IndexFileName), indexData, 0644))

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		Cmd.SetArgs([]string{"list"})
		return Cmd.Execute()
	})
	require.NoError(t, err)
	assert.Contains(t, output, "approved-feature")
	assert.Contains(t, output, "approved")
}

func TestScoreDisplay(t *testing.T) {
	assert.Equal(t, "850", scoreDisplay("850"))
	assert.Equal(t, "—", scoreDisplay(""))
}
