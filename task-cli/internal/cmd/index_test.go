package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"task-cli/pkg/task"
)

func TestIndexCmd_Integration(t *testing.T) {
	projectRoot := t.TempDir()

	// Create feature dirs
	slug := "test-feat"
	featureDir := filepath.Join(projectRoot, "docs", "features", slug)
	tasksDir := filepath.Join(featureDir, "tasks")
	if err := os.MkdirAll(tasksDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a task .md
	taskMD := "---\nid: \"1\"\ntitle: \"Task One\"\npriority: \"P1\"\nestimated_time: \"1h\"\nscope: \"all\"\n---\n\n# Task One\n"
	if err := os.WriteFile(filepath.Join(tasksDir, "1-task-one.md"), []byte(taskMD), 0644); err != nil {
		t.Fatal(err)
	}

	// Create proposal for quick mode detection
	propDir := filepath.Join(projectRoot, "docs", "proposals", slug)
	if err := os.MkdirAll(propDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(propDir, "proposal.md"), []byte("# Proposal"), 0644); err != nil {
		t.Fatal(err)
	}

	indexPath := filepath.Join(tasksDir, "index.json")

	opts := task.BuildIndexOpts{
		FeatureSlug: slug,
		ProjectRoot: projectRoot,
		TasksDir:    tasksDir,
		IndexPath:   indexPath,
		NoTest:      true,
	}

	result, err := task.BuildIndex(opts)
	if err != nil {
		t.Fatalf("BuildIndex: %v", err)
	}
	if result.NewCount != 1 {
		t.Errorf("NewCount = %d, want 1", result.NewCount)
	}

	// Verify index.json exists
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Error("index.json not created")
	}
}
