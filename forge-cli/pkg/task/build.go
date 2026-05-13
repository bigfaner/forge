package task

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// StrategyResolver returns strategy content for a profile+kind pair.
// Returns nil if not found (BuildIndex will use fallback body).
type StrategyResolver func(profileName, kind string) []byte

// BuildIndexOpts holds options for building the task index.
type BuildIndexOpts struct {
	FeatureSlug     string
	ProjectRoot     string
	TasksDir        string // absolute path to tasks/
	IndexPath       string // absolute path to index.json
	NoTest          bool
	TestProfiles    []string // flag > config.yaml > none
	ResolveStrategy StrategyResolver
}

// BuildIndexResult holds the result of a BuildIndex operation.
type BuildIndexResult struct {
	NewCount       int
	UpdatedCount   int
	PreservedCount int
	Warnings       []string
}

// BuildIndex scans .md files and generates/updates index.json.
// It is idempotent: re-running with no changes produces the same output.
func BuildIndex(opts BuildIndexOpts) (*BuildIndexResult, error) {
	result := &BuildIndexResult{}

	// 1. Load existing index or create new
	var index *TaskIndex
	if data, err := os.ReadFile(opts.IndexPath); err == nil {
		index, err = loadIndexFromBytes(data)
		if err != nil {
			return nil, fmt.Errorf("load existing index: %w", err)
		}
	} else {
		index = NewTaskIndex(opts.FeatureSlug)
	}

	// 2. Detect mode
	mode := detectMode(opts.ProjectRoot, opts.FeatureSlug)

	// 3. Set feature metadata
	setFeatureMetadata(index, opts.ProjectRoot, opts.FeatureSlug)

	// 4. Use profiles from opts (caller resolves from config/flags)
	profiles := opts.TestProfiles

	// 5. Scan .md files
	existingKeys := make(map[string]bool) // track which keys come from .md files

	entries, err := os.ReadDir(opts.TasksDir)
	if err != nil {
		return nil, fmt.Errorf("read tasks dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		// Skip special directories/files
		if shouldSkipFile(entry.Name()) {
			continue
		}

		filePath := filepath.Join(opts.TasksDir, entry.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("read %s: %v", entry.Name(), err))
			continue
		}

		fm, _, err := ParseFrontmatter(content)
		if err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("parse %s: %v", entry.Name(), err))
			continue
		}

		if fm.ID == "" {
			result.Warnings = append(result.Warnings, fmt.Sprintf("skip %s: no id in frontmatter", entry.Name()))
			continue
		}

		key := strings.TrimSuffix(entry.Name(), ".md")
		existingKeys[key] = true

		// Build task from frontmatter
		taskType := fm.Type
		if taskType == "" {
			taskType = InferType(fm.ID)
		}

		newTask := Task{
			ID:            fm.ID,
			Title:         fm.Title,
			Priority:      fm.Priority,
			EstimatedTime: fm.EstimatedTime,
			Dependencies:  fm.Dependencies,
			File:          entry.Name(),
			Record:        path.Join("records", entry.Name()),
			Breaking:      fm.Breaking,
			Scope:         fm.Scope,
			MainSession:   fm.MainSession,
			NoTest:        fm.NoTest,
			Type:          taskType,
			Profile:       fm.Profile,
		}

		// Merge with existing
		if existing, found := index.ByID(fm.ID); found {
			// Preserve runtime state
			newTask.Status = existing.Status
			newTask.SourceTaskID = existing.SourceTaskID
			newTask.BlockedReason = existing.BlockedReason
			index.SetTask(key, newTask)
			result.UpdatedCount++
		} else {
			newTask.Status = "pending"
			index.SetTask(key, newTask)
			result.NewCount++
		}
	}

	// 6. Detect orphans
	for key, t := range index.TasksMap() {
		if !existingKeys[key] && !isTestTaskID(t.ID) {
			// Only warn about non-test tasks; test tasks are generated below
			result.Warnings = append(result.Warnings, fmt.Sprintf("orphan: key %q (id=%s) has no .md file", key, t.ID))
			result.PreservedCount++
		}
	}

	// 7. Generate test tasks (unless --no-test)
	if !opts.NoTest && len(profiles) > 0 && mode != "" {
		testTasks := generateTestTasks(mode, profiles)
		if len(testTasks) > 0 {
			ResolveFirstTestDep(testTasks, index.TasksMap(), mode)
		}

		for i, td := range testTasks {
			key := td.Key
			existingKeys[key] = true

			// Resolve strategy content for per-profile tasks
			if td.ProfileName != "" && td.StrategyKind != "" && opts.ResolveStrategy != nil {
				testTasks[i].StrategyContent = opts.ResolveStrategy(td.ProfileName, td.StrategyKind)
			}

			// Generate .md if missing
			mdPath := filepath.Join(opts.TasksDir, key+".md")
			if _, err := os.Stat(mdPath); os.IsNotExist(err) {
				content, err := GenerateTestTaskMD(testTasks[i], opts.FeatureSlug)
				if err != nil {
					result.Warnings = append(result.Warnings, fmt.Sprintf("generate %s: %v", key, err))
					continue
				}
				if err := os.WriteFile(mdPath, content, 0644); err != nil {
					result.Warnings = append(result.Warnings, fmt.Sprintf("write %s: %v", key, err))
					continue
				}
			}

			task := td.TaskFromFile()

			// Merge with existing (preserve runtime state)
			if existing, found := index.ByID(td.ID); found {
				task.Status = existing.Status
				task.SourceTaskID = existing.SourceTaskID
				task.BlockedReason = existing.BlockedReason
				index.SetTask(key, task)
				result.UpdatedCount++
			} else {
				index.SetTask(key, task)
				result.NewCount++
			}
		}
	}

	// 8. Normalize task files (remove empty ## Hard Rules sections)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		if shouldSkipFile(entry.Name()) {
			continue
		}
		filePath := filepath.Join(opts.TasksDir, entry.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}
		normalized := NormalizeTaskMD(content)
		if !bytes.Equal(normalized, content) {
			os.WriteFile(filePath, normalized, 0644)
		}
	}

	// 9. Save index
	if err := SaveIndex(opts.IndexPath, index); err != nil {
		return nil, fmt.Errorf("save index: %w", err)
	}

	return result, nil
}

// detectMode determines the feature mode from file existence.
func detectMode(projectRoot, slug string) string {
	featureDir := filepath.Join(projectRoot, "docs", "features", slug)
	if _, err := os.Stat(filepath.Join(featureDir, "prd", "prd-spec.md")); err == nil {
		return "breakdown"
	}
	if _, err := os.Stat(filepath.Join(projectRoot, "docs", "proposals", slug, "proposal.md")); err == nil {
		return "quick"
	}
	return ""
}

// setFeatureMetadata sets PRD/Design/Proposal paths on the index.
func setFeatureMetadata(index *TaskIndex, projectRoot, slug string) {
	featureDir := filepath.Join(projectRoot, "docs", "features", slug)

	if _, err := os.Stat(filepath.Join(featureDir, "prd", "prd-spec.md")); err == nil {
		index.PRD = "prd/prd-spec.md"
	}
	if _, err := os.Stat(filepath.Join(featureDir, "design", "tech-design.md")); err == nil {
		index.Design = "design/tech-design.md"
	}
	if _, err := os.Stat(filepath.Join(projectRoot, "docs", "proposals", slug, "proposal.md")); err == nil {
		index.Proposal = path.Join("docs", "proposals", slug, "proposal.md")
	}
}

// generateTestTasks returns test task definitions for the given mode and profiles.
func generateTestTasks(mode string, profiles []string) []TestTaskDef {
	switch mode {
	case "breakdown":
		return GetBreakdownTestTasks(profiles)
	case "quick":
		return GetQuickTestTasks(profiles)
	default:
		return nil
	}
}

// shouldSkipFile returns true for files that should not be parsed as task files.
func shouldSkipFile(name string) bool {
	switch {
	case strings.HasPrefix(name, "_"):
		return true
	case name == "index.json":
		return true
	}
	return false
}

// isTestTaskID returns true for test pipeline task IDs.
func isTestTaskID(id string) bool {
	return strings.HasPrefix(id, "T-test-") || strings.HasPrefix(id, "T-quick-")
}

// loadIndexFromBytes deserializes index JSON.
func loadIndexFromBytes(data []byte) (*TaskIndex, error) {
	var idx TaskIndex
	if err := idx.UnmarshalJSON(data); err != nil {
		return nil, err
	}
	return &idx, nil
}
