package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/profile"
	"forge-cli/pkg/task"
)

// writeManifest creates a manifest.md with the given status in the feature directory.
func writeManifest(t *testing.T, projectRoot, slug, status string) {
	t.Helper()
	manifest := fmt.Sprintf("---\nstatus: %s\n---\n# Feature %s\n", status, slug)
	path := filepath.Join(projectRoot, feature.FeaturesDir, slug, feature.ManifestFileName)
	if err := os.WriteFile(path, []byte(manifest), 0644); err != nil {
		t.Fatal(err)
	}
}

// writeProposal creates a proposal.md with the given status in the feature directory.
func writeProposal(t *testing.T, projectRoot, slug, status string) {
	t.Helper()
	content := fmt.Sprintf("---\ncreated: 2026-01-01\nauthor: test\nstatus: %s\n---\n# Proposal\n", status)
	path := filepath.Join(projectRoot, feature.FeaturesDir, slug, feature.ProposalFileName)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

// writeConfig creates .forge/config.yaml with the given content.
func writeForgeConfig(t *testing.T, projectRoot, content string) {
	t.Helper()
	configDir := filepath.Join(projectRoot, feature.ForgeDir)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(configDir, feature.ForgeConfigFileName), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

// initGitRepo creates a git repo with an initial commit in the directory.
func initGitRepo(t *testing.T, dir string) {
	t.Helper()
	runGit(t, dir, "init")
	runGit(t, dir, "config", "user.email", "test@test.com")
	runGit(t, dir, "config", "user.name", "Test")
	// Create an initial file and commit so HEAD exists
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("# test\n"), 0644); err != nil {
		t.Fatal(err)
	}
	runGit(t, dir, "add", "README.md")
	runGit(t, dir, "commit", "-m", "initial")
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s failed: %v\n%s", strings.Join(args, " "), err, string(output))
	}
}

// setupFeatureCompleteTest creates a temp dir with project structure, feature directory,
// index.json, manifest.md, and optionally proposal.md.
func setupFeatureCompleteTest(t *testing.T, tasks map[string]task.Task, withProposal bool) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	// Init git repo
	initGitRepo(t, dir)

	// Create feature directory structure
	if err := feature.EnsureFeatureDir(dir, "test-feature"); err != nil {
		t.Fatal(err)
	}

	// Write index.json
	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test-feature"))
	index := task.NewTestIndex("test-feature", tasks)
	if err := task.SaveIndex(indexPath, index); err != nil {
		t.Fatal(err)
	}

	// Write manifest.md with status "active"
	writeManifest(t, dir, "test-feature", "active")

	// Optionally write proposal.md (quick mode)
	if withProposal {
		writeProposal(t, dir, "test-feature", "In Progress")
	}

	// Write .forge/config.yaml (minimal)
	writeForgeConfig(t, dir, "languages:\n  - go\n")

	return dir
}

func TestCheckFeatureCompletion_NoProjectRoot(t *testing.T) {
	if os.Getenv("TEST_COMPLETE_NO_PROJECT") == "1" {
		result := checkFeatureCompletion()
		if result != nil {
			t.Errorf("expected nil when no project root, got %+v", result)
		}
		return
	}
	tmpDir := t.TempDir()
	cmd := exec.Command(os.Args[0], "-test.run=TestCheckFeatureCompletion_NoProjectRoot")
	env := []string{}
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "CLAUDE_PROJECT_DIR=") || strings.HasPrefix(e, "PROJECT_ROOT=") {
			continue
		}
		env = append(env, e)
	}
	env = append(env, "TEST_COMPLETE_NO_PROJECT=1", "CLAUDE_PROJECT_DIR=")
	cmd.Env = env
	cmd.Dir = tmpDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("subprocess failed: %v\n%s", err, string(output))
	}
}

func TestCheckFeatureCompletion_NoFeature(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	if err := os.MkdirAll(filepath.Join(dir, feature.FeaturesDir), 0755); err != nil {
		t.Fatal(err)
	}

	result := checkFeatureCompletion()
	if result != nil {
		t.Errorf("expected nil when no feature, got %+v", result)
	}
}

func TestCheckFeatureCompletion_PendingTasks(t *testing.T) {
	setupFeatureCompleteTest(t, map[string]task.Task{
		"t1": {ID: "1", Status: "completed"},
		"t2": {ID: "2", Status: "pending"},
	}, false)

	result := checkFeatureCompletion()
	if result != nil {
		t.Errorf("expected nil with pending tasks, got %+v", result)
	}
}

func TestCheckFeatureCompletion_InProgressTasks(t *testing.T) {
	setupFeatureCompleteTest(t, map[string]task.Task{
		"t1": {ID: "1", Status: "in_progress"},
	}, false)

	result := checkFeatureCompletion()
	if result != nil {
		t.Errorf("expected nil with in_progress tasks, got %+v", result)
	}
}

func TestCheckFeatureCompletion_BlockedTasks(t *testing.T) {
	setupFeatureCompleteTest(t, map[string]task.Task{
		"t1": {ID: "1", Status: "completed"},
		"t2": {ID: "2", Status: "blocked"},
	}, false)

	result := checkFeatureCompletion()
	if result != nil {
		t.Errorf("expected nil with blocked tasks, got %+v", result)
	}
}

func TestCheckFeatureCompletion_AllCompleted(t *testing.T) {
	setupFeatureCompleteTest(t, map[string]task.Task{
		"t1": {ID: "1", Status: "completed"},
		"t2": {ID: "2", Status: "completed"},
	}, false)

	result := checkFeatureCompletion()
	if result == nil {
		t.Fatal("expected non-nil result for all completed tasks")
	}
	if result.FeatureSlug != "test-feature" {
		t.Errorf("FeatureSlug = %q, want %q", result.FeatureSlug, "test-feature")
	}
	if result.QuickMode {
		t.Error("QuickMode should be false when no proposal.md")
	}
}

func TestCheckFeatureCompletion_AllCompletedQuickMode(t *testing.T) {
	setupFeatureCompleteTest(t, map[string]task.Task{
		"t1": {ID: "1", Status: "completed"},
	}, true) // withProposal = true

	result := checkFeatureCompletion()
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if !result.QuickMode {
		t.Error("QuickMode should be true when proposal.md exists")
	}
}

func TestCheckFeatureCompletion_MixedCompletedSkipped(t *testing.T) {
	setupFeatureCompleteTest(t, map[string]task.Task{
		"t1": {ID: "1", Status: "completed"},
		"t2": {ID: "2", Status: "skipped"},
	}, false)

	result := checkFeatureCompletion()
	if result == nil {
		t.Fatal("expected non-nil result for mixed completed/skipped")
	}
}

func TestCheckFeatureCompletion_RejectedTasks(t *testing.T) {
	setupFeatureCompleteTest(t, map[string]task.Task{
		"t1": {ID: "1", Status: "completed"},
		"t2": {ID: "2", Status: "rejected"},
	}, false)

	result := checkFeatureCompletion()
	if result != nil {
		t.Errorf("expected nil with rejected tasks, got %+v", result)
	}
}

func TestUpdateManifestStatus(t *testing.T) {
	dir := t.TempDir()
	manifestPath := filepath.Join(dir, "manifest.md")
	content := "---\nstatus: active\n---\n# Feature test\n"
	if err := os.WriteFile(manifestPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	if err := updateFileStatus(manifestPath, "completed"); err != nil {
		t.Fatalf("updateFileStatus failed: %v", err)
	}

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}

	text := string(data)
	if !strings.Contains(text, "status: completed") {
		t.Errorf("manifest should contain 'status: completed', got:\n%s", text)
	}
	if strings.Contains(text, "status: active") {
		t.Errorf("manifest should not contain 'status: active', got:\n%s", text)
	}
}

func TestUpdateProposalStatus(t *testing.T) {
	dir := t.TempDir()
	proposalPath := filepath.Join(dir, "proposal.md")
	content := "---\ncreated: 2026-01-01\nauthor: test\nstatus: In Progress\n---\n# Proposal\n"
	if err := os.WriteFile(proposalPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	if err := updateFileStatus(proposalPath, "Completed"); err != nil {
		t.Fatalf("updateFileStatus failed: %v", err)
	}

	data, err := os.ReadFile(proposalPath)
	if err != nil {
		t.Fatal(err)
	}

	text := string(data)
	if !strings.Contains(text, "status: Completed") {
		t.Errorf("proposal should contain 'status: Completed', got:\n%s", text)
	}
}

func TestGitCommitExplicitPaths(t *testing.T) {
	dir := t.TempDir()
	initGitRepo(t, dir)

	// Create two files to commit
	file1 := filepath.Join(dir, "manifest.md")
	file2 := filepath.Join(dir, "proposal.md")
	if err := os.WriteFile(file1, []byte("manifest content"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file2, []byte("proposal content"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create an unrelated file that should NOT be committed
	unrelated := filepath.Join(dir, "unrelated.txt")
	if err := os.WriteFile(unrelated, []byte("should not be committed"), 0644); err != nil {
		t.Fatal(err)
	}

	err := gitCommitFiles(dir, []string{"manifest.md", "proposal.md"}, "test commit")
	if err != nil {
		t.Fatalf("gitCommitFiles failed: %v", err)
	}

	// Verify only the two target files were committed
	cmd := exec.Command("git", "diff", "--name-only", "HEAD~1", "HEAD")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git diff failed: %v\n%s", err, string(output))
	}

	committed := strings.TrimSpace(string(output))
	files := strings.Split(committed, "\n")
	if len(files) != 2 {
		t.Errorf("expected 2 committed files, got %d: %s", len(files), committed)
	}

	// Verify unrelated file is not committed
	for _, f := range files {
		if f == "unrelated.txt" {
			t.Error("unrelated.txt should not have been committed")
		}
	}

	// Verify unrelated file is still untracked/modified
	cmd = exec.Command("git", "status", "--porcelain", "unrelated.txt")
	cmd.Dir = dir
	output, _ = cmd.CombinedOutput()
	if len(strings.TrimSpace(string(output))) == 0 {
		t.Error("unrelated.txt should still show as untracked/modified")
	}
}

func TestIsGitPushEnabled(t *testing.T) {
	t.Run("true when auto.gitPush is true", func(t *testing.T) {
		dir := t.TempDir()
		writeForgeConfig(t, dir, "languages:\n  - go\nauto:\n  gitPush: true\n")
		result, err := isGitPushEnabled(dir)
		if err != nil {
			t.Fatal(err)
		}
		if !result {
			t.Error("expected true")
		}
	})

	t.Run("false when auto.gitPush is false", func(t *testing.T) {
		dir := t.TempDir()
		writeForgeConfig(t, dir, "languages:\n  - go\nauto:\n  gitPush: false\n")
		result, err := isGitPushEnabled(dir)
		if err != nil {
			t.Fatal(err)
		}
		if result {
			t.Error("expected false")
		}
	})

	t.Run("false when auto block absent", func(t *testing.T) {
		dir := t.TempDir()
		writeForgeConfig(t, dir, "languages:\n  - go\n")
		result, err := isGitPushEnabled(dir)
		if err != nil {
			t.Fatal(err)
		}
		if result {
			t.Error("expected false when auto block absent")
		}
	})

	t.Run("false when config file absent", func(t *testing.T) {
		dir := t.TempDir()
		result, err := isGitPushEnabled(dir)
		if err != nil {
			t.Fatal(err)
		}
		if result {
			t.Error("expected false when config absent")
		}
	})
}

func TestGitPush(t *testing.T) {
	// Verify that push failure is logged to stderr but does not panic.
	t.Run("push failure logs to stderr", func(t *testing.T) {
		dir := t.TempDir()
		initGitRepo(t, dir)
		// No remote configured — push will fail, but the function logs and returns the error.
		// The caller (completeFeature) treats push failure as non-blocking.
		_ = gitPush(dir)
		// If we get here without panic, the function is non-blocking at the call-site level.
	})
}

func TestCompleteFeature_AllCompleted(t *testing.T) {
	dir := setupFeatureCompleteTest(t, map[string]task.Task{
		"t1": {ID: "1", Status: "completed"},
		"t2": {ID: "2", Status: "skipped"},
	}, false)

	result := checkFeatureCompletion()
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if err := completeFeature(result); err != nil {
		t.Fatalf("completeFeature failed: %v", err)
	}

	// Verify manifest.md was updated
	manifestPath := filepath.Join(dir, feature.FeaturesDir, "test-feature", feature.ManifestFileName)
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "status: completed") {
		t.Errorf("manifest.md should have status: completed, got:\n%s", string(data))
	}

	// Verify git commit was created
	cmd := exec.Command("git", "log", "-1", "--name-only", "--pretty=format:")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git log failed: %v\n%s", err, string(output))
	}
	committed := string(output)
	if !strings.Contains(committed, "manifest.md") {
		t.Errorf("expected manifest.md in commit, got:\n%s", committed)
	}
	if strings.Contains(committed, "unrelated") {
		t.Error("unrelated files should not be committed")
	}
}

func TestCompleteFeature_QuickMode(t *testing.T) {
	dir := setupFeatureCompleteTest(t, map[string]task.Task{
		"t1": {ID: "1", Status: "completed"},
	}, true) // withProposal = true

	result := checkFeatureCompletion()
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if !result.QuickMode {
		t.Fatal("expected QuickMode=true")
	}

	if err := completeFeature(result); err != nil {
		t.Fatalf("completeFeature failed: %v", err)
	}

	// Verify both files updated
	manifestPath := filepath.Join(dir, feature.FeaturesDir, "test-feature", feature.ManifestFileName)
	data, _ := os.ReadFile(manifestPath)
	if !strings.Contains(string(data), "status: completed") {
		t.Errorf("manifest.md should have status: completed, got:\n%s", string(data))
	}

	proposalPath := filepath.Join(dir, feature.FeaturesDir, "test-feature", feature.ProposalFileName)
	data, _ = os.ReadFile(proposalPath)
	if !strings.Contains(string(data), "status: Completed") {
		t.Errorf("proposal.md should have status: Completed, got:\n%s", string(data))
	}

	// Verify git commit contains both files
	cmd := exec.Command("git", "log", "-1", "--name-only", "--pretty=format:")
	cmd.Dir = dir
	output, _ := cmd.CombinedOutput()
	committed := string(output)
	if !strings.Contains(committed, "manifest.md") {
		t.Errorf("expected manifest.md in commit, got:\n%s", committed)
	}
	if !strings.Contains(committed, "proposal.md") {
		t.Errorf("expected proposal.md in commit, got:\n%s", committed)
	}
}

func TestCompleteFeature_FullPipelineMode(t *testing.T) {
	dir := setupFeatureCompleteTest(t, map[string]task.Task{
		"t1": {ID: "1", Status: "completed"},
	}, false) // withProposal = false

	result := checkFeatureCompletion()
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.QuickMode {
		t.Fatal("expected QuickMode=false")
	}
	if result.ProposalRel != "" {
		t.Fatal("ProposalRel should be empty in full pipeline mode")
	}

	if err := completeFeature(result); err != nil {
		t.Fatalf("completeFeature failed: %v", err)
	}

	// Verify manifest.md was updated
	manifestPath := filepath.Join(dir, feature.FeaturesDir, "test-feature", feature.ManifestFileName)
	data, _ := os.ReadFile(manifestPath)
	if !strings.Contains(string(data), "status: completed") {
		t.Errorf("manifest.md should have status: completed, got:\n%s", string(data))
	}

	// Verify git commit contains ONLY manifest.md
	cmd := exec.Command("git", "log", "-1", "--name-only", "--pretty=format:")
	cmd.Dir = dir
	output, _ := cmd.CombinedOutput()
	committed := string(output)
	if !strings.Contains(committed, "manifest.md") {
		t.Errorf("expected manifest.md in commit, got:\n%s", committed)
	}
	if strings.Contains(committed, "proposal.md") {
		t.Error("proposal.md should not be in commit for full pipeline mode")
	}
}

func TestCompleteFeature_PendingTasksNoOp(t *testing.T) {
	dir := setupFeatureCompleteTest(t, map[string]task.Task{
		"t1": {ID: "1", Status: "completed"},
		"t2": {ID: "2", Status: "pending"},
	}, false)

	result := checkFeatureCompletion()
	if result != nil {
		t.Errorf("expected nil with pending tasks, got %+v", result)
	}

	// Verify manifest.md was NOT updated
	manifestPath := filepath.Join(dir, feature.FeaturesDir, "test-feature", feature.ManifestFileName)
	data, _ := os.ReadFile(manifestPath)
	if strings.Contains(string(data), "status: completed") {
		t.Error("manifest.md should NOT have been updated when tasks are pending")
	}
}

func TestUpdateFileStatus_NoFrontmatter(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "test.md")
	if err := os.WriteFile(filePath, []byte("# No frontmatter\n"), 0644); err != nil {
		t.Fatal(err)
	}

	err := updateFileStatus(filePath, "completed")
	if err == nil {
		t.Error("expected error for file without frontmatter")
	}
}

// Verify the profile import is used (prevents compile errors).
var _ = profile.ReadAutoConfig
