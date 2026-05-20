//go:build e2e

package e2e

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// createTempForgeProject creates a temp directory with go.mod and .forge/state.json
// pointing to the given feature slug. Returns the project dir path.
func createTempForgeProject(t *testing.T, slug string) string {
	t.Helper()
	dir := t.TempDir()
	err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test-project\n\ngo 1.26\n"), 0644)
	assert.NoError(t, err, "failed to create go.mod")

	// Create .forge/state.json to set feature context
	forgeDir := filepath.Join(dir, ".forge")
	err = os.MkdirAll(forgeDir, 0755)
	assert.NoError(t, err, "failed to create .forge directory")

	state := map[string]interface{}{
		"feature":      slug,
		"allCompleted": false,
		"updatedAt":    "2026-05-17T00:00:00Z",
	}
	stateData, err := json.Marshal(state)
	assert.NoError(t, err, "failed to marshal forge state")
	err = os.WriteFile(filepath.Join(forgeDir, "state.json"), stateData, 0644)
	assert.NoError(t, err, "failed to write state.json")

	return dir
}

// createFeatureIndex creates a minimal feature with index.json in the project dir.
func createFeatureIndex(t *testing.T, projectDir, slug string, indexJSON string) {
	t.Helper()
	tasksDir := filepath.Join(projectDir, "docs", "features", slug, "tasks")
	err := os.MkdirAll(tasksDir, 0755)
	assert.NoError(t, err, "failed to create feature tasks dir")

	err = os.WriteFile(filepath.Join(tasksDir, "index.json"), []byte(indexJSON), 0644)
	assert.NoError(t, err, "failed to write index.json")
}

// createRecordFile creates a record file in the feature's records directory.
func createRecordFile(t *testing.T, projectDir, slug, recordFile, content string) {
	t.Helper()
	recordsDir := filepath.Join(projectDir, "docs", "features", slug, "tasks", "records")
	err := os.MkdirAll(recordsDir, 0755)
	assert.NoError(t, err, "failed to create records dir")

	err = os.WriteFile(filepath.Join(recordsDir, recordFile), []byte(content), 0644)
	assert.NoError(t, err, "failed to write record file")
}

// writeSubmitJSON writes a minimal record JSON file and returns its path.
func writeSubmitJSON(t *testing.T, dir string, summary string) string {
	t.Helper()
	record := map[string]interface{}{
		"summary":     summary,
		"status":      "completed",
		"testsPassed": 1,
		"testsFailed": 0,
		"coverage":    100.0,
		"acceptanceCriteria": []map[string]interface{}{
			{"criterion": "works", "met": true},
		},
	}
	data, err := json.Marshal(record)
	assert.NoError(t, err, "failed to marshal record JSON")
	path := filepath.Join(dir, "record.json")
	err = os.WriteFile(path, data, 0644)
	assert.NoError(t, err, "failed to write record JSON")
	return path
}

// minimalIndexJSON returns a valid index.json with one task.
func minimalIndexJSON(slug, taskKey, taskID, taskTitle string) string {
	index := map[string]interface{}{
		"feature": slug,
		"tasks": map[string]interface{}{
			taskKey: map[string]interface{}{
				"id":       taskID,
				"title":    taskTitle,
				"priority": "P1",
				"status":   "pending",
				"file":     taskKey + ".md",
				"record":   "records/" + taskKey + ".md",
			},
		},
		"statusEnum":   []string{"pending", "in_progress", "completed", "blocked", "skipped", "rejected"},
		"priorityEnum": []string{"P0", "P1", "P2"},
	}
	data, _ := json.Marshal(index)
	return string(data)
}

// envBlacklist lists env vars that override project root detection,
// causing forge to resolve to the real project instead of the test's temp dir.
var envBlacklist = []string{"CLAUDE_PROJECT_DIR", "PROJECT_ROOT"}

// cleanForgeEnv returns os.Environ with project-root override vars removed.
func cleanForgeEnv() []string {
	env := os.Environ()
	filtered := make([]string, 0, len(env))
	for _, e := range env {
		skip := false
		for _, bl := range envBlacklist {
			if strings.HasPrefix(e, bl+"=") {
				skip = true
				break
			}
		}
		if !skip {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

// runForgeInDir runs forge CLI in a given working directory, fatalfing on non-zero exit.
func runForgeInDir(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command(ForgeBinary, args...)
	cmd.Dir = dir
	cmd.Env = cleanForgeEnv()
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("forge %s failed in %s: %s\noutput: %s", strings.Join(args, " "), dir, err, out)
	}
	return string(out)
}

// runForgeInDirRaw runs forge CLI in a given working directory, returning output and exit code.
func runForgeInDirRaw(t *testing.T, dir string, args ...string) (string, int) {
	t.Helper()
	cmd := exec.Command(ForgeBinary, args...)
	cmd.Dir = dir
	cmd.Env = cleanForgeEnv()
	out, err := cmd.CombinedOutput()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}
	return string(out), exitCode
}

// triParseBlock extracts lines between "---" separators from raw CLI output.
func triParseBlock(t *testing.T, raw string) []string {
	t.Helper()
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	if len(lines) < 2 || strings.TrimSpace(lines[0]) != "---" || strings.TrimSpace(lines[len(lines)-1]) != "---" {
		t.Fatalf("output must be wrapped in --- separators, got:\n%s", raw)
	}
	inner := lines[1 : len(lines)-1]
	result := make([]string, 0, len(inner))
	for _, l := range inner {
		result = append(result, strings.TrimSpace(l))
	}
	return result
}

// triHasField checks that a parsed block contains a "KEY: value" line.
func triHasField(lines []string, key, value string) bool {
	prefix := key + ": "
	for _, l := range lines {
		if strings.HasPrefix(l, prefix) {
			if value == "" {
				return true
			}
			return l == prefix+value
		}
	}
	return false
}

// triHasNoField checks that a parsed block does NOT contain any line starting with key.
func triHasNoField(lines []string, key string) bool {
	prefix := key + ": "
	for _, l := range lines {
		if strings.HasPrefix(l, prefix) {
			return false
		}
	}
	return true
}

// TC-001: Submit record succeeds when no record exists
// Traceability: TC-001 -> Proposal / Key Scenario 1, SC-1 (implicit happy path)
func TestTC_001_SubmitRecordSucceedsWhenNoRecordExists(t *testing.T) {
	slug := "test-feature-tri"
	projectDir := createTempForgeProject(t, slug)
	indexJSON := minimalIndexJSON(slug, "task-1", "1", "Test Task")
	createFeatureIndex(t, projectDir, slug, indexJSON)

	recordPath := writeSubmitJSON(t, projectDir, "Implemented feature successfully")
	output, exitCode := runForgeInDirRaw(t, projectDir, "task", "submit", "1", "--data", recordPath, "--force")

	assert.Equal(t, 0, exitCode, "submit should exit 0, got output:\n%s", output)

	// Verify record file was created
	recordFile := filepath.Join(projectDir, "docs", "features", slug, "tasks", "records", "task-1.md")
	_, err := os.Stat(recordFile)
	assert.NoError(t, err, "record file should exist at %s", recordFile)

	// Verify record content contains the summary
	content, err := os.ReadFile(recordFile)
	assert.NoError(t, err, "failed to read record file")
	assert.Contains(t, string(content), "Implemented feature successfully", "record should contain the submission summary")
}

// TC-002: Submit record blocked when record already exists
// Traceability: TC-002 -> Proposal / Key Scenario 2, SC-1
func TestTC_002_SubmitRecordBlockedWhenRecordAlreadyExists(t *testing.T) {
	slug := "test-feature-tri"
	projectDir := createTempForgeProject(t, slug)
	indexJSON := minimalIndexJSON(slug, "task-1", "1", "Test Task")
	createFeatureIndex(t, projectDir, slug, indexJSON)

	// Pre-create the record file
	recordFile := filepath.Join(projectDir, "docs", "features", slug, "tasks", "records", "task-1.md")
	createRecordFile(t, projectDir, slug, "task-1.md", "original content that must not change")

	recordPath := writeSubmitJSON(t, projectDir, "New submission attempt")
	output, exitCode := runForgeInDirRaw(t, projectDir, "task", "submit", "1", "--data", recordPath)

	assert.Equal(t, 1, exitCode, "submit without --force should exit 1 when record exists, got output:\n%s", output)
	assert.Contains(t, output, "Record for task 1 already exists", "stderr should contain record exists message")
	assert.Contains(t, output, "Use --force to overwrite, or create a fix task instead", "stderr should contain hint")

	// Verify original record is NOT modified
	content, err := os.ReadFile(recordFile)
	assert.NoError(t, err, "failed to read record file")
	assert.Equal(t, "original content that must not change", string(content), "record file should not be modified")
}

// TC-003: Submit record with --force overwrites existing record
// Traceability: TC-003 -> Proposal / Key Scenario 3, SC-2
func TestTC_003_SubmitWithForceOverwritesExistingRecord(t *testing.T) {
	slug := "test-feature-tri"
	projectDir := createTempForgeProject(t, slug)
	indexJSON := minimalIndexJSON(slug, "task-1", "1", "Test Task")
	createFeatureIndex(t, projectDir, slug, indexJSON)

	// Pre-create the record file with known content
	createRecordFile(t, projectDir, slug, "task-1.md", "old content")

	recordPath := writeSubmitJSON(t, projectDir, "Force-overwritten submission")
	output, exitCode := runForgeInDirRaw(t, projectDir, "task", "submit", "1", "--data", recordPath, "--force")

	assert.Equal(t, 0, exitCode, "submit --force should exit 0, got output:\n%s", output)
	assert.Contains(t, output, "WARNING: Overwriting existing record", "stderr should contain warning")

	// Verify record file is replaced
	recordFile := filepath.Join(projectDir, "docs", "features", slug, "tasks", "records", "task-1.md")
	content, err := os.ReadFile(recordFile)
	assert.NoError(t, err, "failed to read record file")
	assert.Contains(t, string(content), "Force-overwritten submission", "record should contain new submission content")
	assert.NotContains(t, string(content), "old content", "record should NOT contain old content")
}

// TC-004: Default query shows 4 fields unchanged
// Traceability: TC-004 -> Proposal / Key Scenario 4, SC-3
func TestTC_004_DefaultQueryShowsFourFieldsUnchanged(t *testing.T) {
	slug := "test-feature-tri"
	projectDir := createTempForgeProject(t, slug)

	indexData := map[string]interface{}{
		"feature": slug,
		"tasks": map[string]interface{}{
			"task-1": map[string]interface{}{
				"id":       "1",
				"title":    "Test Task",
				"priority": "P1",
				"status":   "pending",
				"file":     "task-1.md",
				"record":   "records/task-1.md",
				"scope":    "backend",
			},
		},
		"statusEnum":   []string{"pending", "in_progress", "completed", "blocked", "skipped", "rejected"},
		"priorityEnum": []string{"P0", "P1", "P2"},
	}
	indexJSON, _ := json.Marshal(indexData)
	createFeatureIndex(t, projectDir, slug, string(indexJSON))

	output := runForgeInDir(t, projectDir, "task", "query", "1")

	lines := triParseBlock(t, output)
	assert.True(t, triHasField(lines, "TASK_ID", "1"), "should contain TASK_ID")
	assert.True(t, triHasField(lines, "STATUS", "pending"), "should contain STATUS")
	assert.True(t, triHasField(lines, "SCOPE", "backend"), "should contain SCOPE")
	// Verify non-default fields are absent
	assert.True(t, triHasNoField(lines, "TITLE"), "should NOT contain TITLE in default mode")
	assert.True(t, triHasNoField(lines, "PRIORITY"), "should NOT contain PRIORITY in default mode")
	assert.True(t, triHasNoField(lines, "TYPE"), "should NOT contain TYPE in default mode")
	assert.True(t, triHasNoField(lines, "DEPENDENCIES"), "should NOT contain DEPENDENCIES in default mode")
	assert.True(t, triHasNoField(lines, "TASK_FILE"), "should NOT contain TASK_FILE in default mode")
	assert.True(t, triHasNoField(lines, "RECORD_FILE"), "should NOT contain RECORD_FILE in default mode")
	assert.True(t, triHasNoField(lines, "KEY"), "should NOT contain KEY in default mode")
	assert.True(t, triHasNoField(lines, "RELATED_FIXES"), "should NOT contain RELATED_FIXES in default mode")
}

// TC-005: Verbose query displays all task fields
// Traceability: TC-005 -> Proposal / Key Scenario 5, SC-4
func TestTC_005_VerboseQueryDisplaysAllTaskFields(t *testing.T) {
	slug := "test-feature-tri"
	projectDir := createTempForgeProject(t, slug)

	indexData := map[string]interface{}{
		"feature": slug,
		"tasks": map[string]interface{}{
			"task-1": map[string]interface{}{
				"id":           "1",
				"title":        "Test Task Title",
				"priority":     "P1",
				"status":       "pending",
				"type":         "feature",
				"file":         "task-1.md",
				"record":       "records/task-1.md",
				"scope":        "backend",
				"dependencies": []string{"2", "3"},
			},
		},
		"statusEnum":   []string{"pending", "in_progress", "completed", "blocked", "skipped", "rejected"},
		"priorityEnum": []string{"P0", "P1", "P2"},
	}
	indexJSON, _ := json.Marshal(indexData)
	createFeatureIndex(t, projectDir, slug, string(indexJSON))

	output := runForgeInDir(t, projectDir, "task", "query", "1", "--verbose")

	// Output may span multiple blocks (main fields + RELATED_FIXES if any)
	// The main block is bounded by --- markers
	assert.Contains(t, output, "---", "output should contain block separator")
	assert.Contains(t, output, "KEY:", "verbose should contain KEY")
	assert.Contains(t, output, "TASK_ID: 1", "verbose should contain TASK_ID")
	assert.Contains(t, output, "TITLE: Test Task Title", "verbose should contain TITLE")
	assert.Contains(t, output, "STATUS: pending", "verbose should contain STATUS")
	assert.Contains(t, output, "PRIORITY: P1", "verbose should contain PRIORITY")
	assert.Contains(t, output, "TYPE: feature", "verbose should contain TYPE")
	assert.Contains(t, output, "SCOPE: backend", "verbose should contain SCOPE")
	assert.Contains(t, output, "DEPENDENCIES:", "verbose should contain DEPENDENCIES")
	assert.Contains(t, output, "  2", "verbose should list dependency 2")
	assert.Contains(t, output, "  3", "verbose should list dependency 3")
	assert.Contains(t, output, "TASK_FILE:", "verbose should contain TASK_FILE")
	assert.Contains(t, output, "RECORD_FILE:", "verbose should contain RECORD_FILE")
}

// TC-006: Verbose query shows RELATED_FIXES for tasks with fix records
// Traceability: TC-006 -> Proposal / SC-5
func TestTC_006_VerboseQueryShowsRelatedFixes(t *testing.T) {
	slug := "test-feature-tri"
	projectDir := createTempForgeProject(t, slug)

	indexData := map[string]interface{}{
		"feature": slug,
		"tasks": map[string]interface{}{
			"task-2": map[string]interface{}{
				"id":       "2",
				"title":    "Original Task",
				"priority": "P1",
				"status":   "completed",
				"file":     "task-2.md",
				"record":   "records/task-2.md",
			},
			"fix-1": map[string]interface{}{
				"id":           "fix-1",
				"title":        "Fix the bug",
				"priority":     "P0",
				"status":       "completed",
				"file":         "fix-1.md",
				"record":       "records/fix-1.md",
				"sourceTaskID": "2",
			},
		},
		"statusEnum":   []string{"pending", "in_progress", "completed", "blocked", "skipped", "rejected"},
		"priorityEnum": []string{"P0", "P1", "P2"},
	}
	indexJSON, _ := json.Marshal(indexData)
	createFeatureIndex(t, projectDir, slug, string(indexJSON))

	output := runForgeInDir(t, projectDir, "task", "query", "2", "--verbose")

	assert.Contains(t, output, "RELATED_FIXES:", "should contain RELATED_FIXES field")
	assert.Contains(t, output, "fix-1 [completed] Fix the bug", "should show fix as '<id> [<status>] <title>'")
}

// TC-007: Verbose query omits RELATED_FIXES when no fixes exist
// Traceability: TC-007 -> Proposal / SC-6
func TestTC_007_VerboseQueryOmitsRelatedFixesWhenNoneExist(t *testing.T) {
	slug := "test-feature-tri"
	projectDir := createTempForgeProject(t, slug)

	indexData := map[string]interface{}{
		"feature": slug,
		"tasks": map[string]interface{}{
			"task-1": map[string]interface{}{
				"id":       "1",
				"title":    "Solo Task",
				"priority": "P1",
				"status":   "pending",
				"file":     "task-1.md",
				"record":   "records/task-1.md",
			},
		},
		"statusEnum":   []string{"pending", "in_progress", "completed", "blocked", "skipped", "rejected"},
		"priorityEnum": []string{"P0", "P1", "P2"},
	}
	indexJSON, _ := json.Marshal(indexData)
	createFeatureIndex(t, projectDir, slug, string(indexJSON))

	output := runForgeInDir(t, projectDir, "task", "query", "1", "--verbose")

	assert.NotContains(t, output, "RELATED_FIXES", "should NOT contain RELATED_FIXES when no fixes exist")
	// Verify other verbose fields are still present
	assert.Contains(t, output, "KEY:", "should contain KEY")
	assert.Contains(t, output, "TASK_ID: 1", "should contain TASK_ID")
	assert.Contains(t, output, "TITLE: Solo Task", "should contain TITLE")
}

// TC-008: Status command behavior unchanged
// Traceability: TC-008 -> Proposal / SC-7
func TestTC_008_StatusCommandBehaviorUnchanged(t *testing.T) {
	slug := "test-feature-tri"
	projectDir := createTempForgeProject(t, slug)

	indexData := map[string]interface{}{
		"feature": slug,
		"tasks": map[string]interface{}{
			"task-1": map[string]interface{}{
				"id":       "1",
				"title":    "Test Task",
				"priority": "P1",
				"status":   "in_progress",
				"file":     "task-1.md",
				"record":   "records/task-1.md",
			},
		},
		"statusEnum":   []string{"pending", "in_progress", "completed", "blocked", "skipped", "rejected"},
		"priorityEnum": []string{"P0", "P1", "P2"},
	}
	indexJSON, _ := json.Marshal(indexData)
	createFeatureIndex(t, projectDir, slug, string(indexJSON))

	output := runForgeInDir(t, projectDir, "task", "status", "1")

	lines := triParseBlock(t, output)
	assert.True(t, triHasField(lines, "TASK_ID", "1"), "status output should contain TASK_ID")
	assert.True(t, triHasField(lines, "STATUS", "in_progress"), "status output should contain STATUS")
}

// TC-009: Verbose query with short flag -v
// Traceability: TC-009 -> Proposal / Proposed Solution 2 ("--verbose / -v flag")
func TestTC_009_VerboseQueryWithShortFlag(t *testing.T) {
	slug := "test-feature-tri"
	projectDir := createTempForgeProject(t, slug)

	indexData := map[string]interface{}{
		"feature": slug,
		"tasks": map[string]interface{}{
			"task-1": map[string]interface{}{
				"id":       "1",
				"title":    "Test Task Title",
				"priority": "P1",
				"status":   "pending",
				"type":     "feature",
				"file":     "task-1.md",
				"record":   "records/task-1.md",
			},
		},
		"statusEnum":   []string{"pending", "in_progress", "completed", "blocked", "skipped", "rejected"},
		"priorityEnum": []string{"P0", "P1", "P2"},
	}
	indexJSON, _ := json.Marshal(indexData)
	createFeatureIndex(t, projectDir, slug, string(indexJSON))

	verboseOutput := runForgeInDir(t, projectDir, "task", "query", "1", "--verbose")
	shortOutput := runForgeInDir(t, projectDir, "task", "query", "1", "-v")

	assert.Equal(t, verboseOutput, shortOutput, "-v output should be identical to --verbose output")
}

// TC-010: Verbose query omits SCOPE when empty
// Traceability: TC-010 -> Proposal / Proposed Solution 2 ("SCOPE (omit if empty)")
func TestTC_010_VerboseQueryOmitsScopeWhenEmpty(t *testing.T) {
	slug := "test-feature-tri"
	projectDir := createTempForgeProject(t, slug)

	indexData := map[string]interface{}{
		"feature": slug,
		"tasks": map[string]interface{}{
			"task-1": map[string]interface{}{
				"id":       "1",
				"title":    "No Scope Task",
				"priority": "P1",
				"status":   "pending",
				"file":     "task-1.md",
				"record":   "records/task-1.md",
				// scope intentionally omitted
			},
		},
		"statusEnum":   []string{"pending", "in_progress", "completed", "blocked", "skipped", "rejected"},
		"priorityEnum": []string{"P0", "P1", "P2"},
	}
	indexJSON, _ := json.Marshal(indexData)
	createFeatureIndex(t, projectDir, slug, string(indexJSON))

	output := runForgeInDir(t, projectDir, "task", "query", "1", "--verbose")

	assert.NotContains(t, output, "SCOPE:", "should NOT contain SCOPE when empty")
	// Verify other verbose fields are still present
	assert.Contains(t, output, "KEY:", "should contain KEY")
	assert.Contains(t, output, "TASK_ID: 1", "should contain TASK_ID")
}

// TC-011: Verbose query omits BREAKING when false
// Traceability: TC-011 -> Proposal / Proposed Solution 2 (default mode: "BREAKING (if true)")
func TestTC_011_VerboseQueryOmitsBreakingWhenFalse(t *testing.T) {
	slug := "test-feature-tri"
	projectDir := createTempForgeProject(t, slug)

	indexData := map[string]interface{}{
		"feature": slug,
		"tasks": map[string]interface{}{
			"task-1": map[string]interface{}{
				"id":       "1",
				"title":    "Non-Breaking Task",
				"priority": "P1",
				"status":   "pending",
				"file":     "task-1.md",
				"record":   "records/task-1.md",
			},
		},
		"statusEnum":   []string{"pending", "in_progress", "completed", "blocked", "skipped", "rejected"},
		"priorityEnum": []string{"P0", "P1", "P2"},
	}
	indexJSON, _ := json.Marshal(indexData)
	createFeatureIndex(t, projectDir, slug, string(indexJSON))

	output := runForgeInDir(t, projectDir, "task", "query", "1", "--verbose")

	assert.NotContains(t, output, "BREAKING", "should NOT contain BREAKING when false/unset")
	assert.Contains(t, output, "KEY:", "should contain KEY")
	assert.Contains(t, output, "TASK_ID: 1", "should contain TASK_ID")
}

// TC-012: Verbose query displays multi-line DEPENDENCIES
// Traceability: TC-012 -> Proposal / Proposed Solution 2 ("DEPENDENCIES (multi-line if multiple)")
func TestTC_012_VerboseQueryDisplaysMultiLineDependencies(t *testing.T) {
	slug := "test-feature-tri"
	projectDir := createTempForgeProject(t, slug)

	indexData := map[string]interface{}{
		"feature": slug,
		"tasks": map[string]interface{}{
			"task-1": map[string]interface{}{
				"id":           "1",
				"title":        "Multi Dep Task",
				"priority":     "P1",
				"status":       "pending",
				"file":         "task-1.md",
				"record":       "records/task-1.md",
				"dependencies": []string{"2", "3"},
			},
			"task-2": map[string]interface{}{
				"id":       "2",
				"title":    "Dep A",
				"priority": "P1",
				"status":   "completed",
				"file":     "task-2.md",
				"record":   "records/task-2.md",
			},
			"task-3": map[string]interface{}{
				"id":       "3",
				"title":    "Dep B",
				"priority": "P1",
				"status":   "completed",
				"file":     "task-3.md",
				"record":   "records/task-3.md",
			},
		},
		"statusEnum":   []string{"pending", "in_progress", "completed", "blocked", "skipped", "rejected"},
		"priorityEnum": []string{"P0", "P1", "P2"},
	}
	indexJSON, _ := json.Marshal(indexData)
	createFeatureIndex(t, projectDir, slug, string(indexJSON))

	output := runForgeInDir(t, projectDir, "task", "query", "1", "--verbose")

	assert.Contains(t, output, "DEPENDENCIES", "should contain DEPENDENCIES field")

	// Verify each dependency appears on its own line (indented with 2 spaces)
	lines := strings.Split(output, "\n")
	var depLines []string
	inDeps := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "DEPENDENCIES") {
			inDeps = true
			continue
		}
		if inDeps && strings.HasPrefix(line, "  ") {
			depLines = append(depLines, strings.TrimSpace(line))
		} else if inDeps && trimmed != "" {
			inDeps = false
		}
	}
	assert.Equal(t, 2, len(depLines), "should have 2 dependency lines, got: %v", depLines)
	if len(depLines) >= 2 {
		assert.Equal(t, "2", depLines[0], "first dependency should be 2")
		assert.Equal(t, "3", depLines[1], "second dependency should be 3")
	}
}
