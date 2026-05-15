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
	TasksDir        string   // absolute path to tasks/
	IndexPath       string   // absolute path to index.json
	TestProfiles    []string // flag > config.yaml > none
	ResolveStrategy StrategyResolver
}

// BuildIndexResult holds the result of a BuildIndex operation.
type BuildIndexResult struct {
	NewCount            int
	UpdatedCount        int
	PreservedCount      int
	StageGatesGenerated int
	Warnings            []string
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

	// 5.5 Validate: business tasks from .md files must have a type
	for key, t := range index.TasksMap() {
		if !existingKeys[key] {
			continue // auto-generated tasks are not from .md files
		}
		if isAutoGenTaskID(t.ID) {
			continue // auto-generated tasks use InferType
		}
		if t.Type == "" {
			return nil, fmt.Errorf("task %s has empty type (set type in frontmatter or use a recognizable ID pattern)", t.File)
		}
	}

	// 5.5.1 Detect docs-only feature
	docsOnly := isDocsOnlyFeature(index.TasksMap())

	// 6. Detect orphans
	for key, t := range index.TasksMap() {
		if !existingKeys[key] && !isTestTaskID(t.ID) {
			// Only warn about non-test tasks; test tasks are generated below
			result.Warnings = append(result.Warnings, fmt.Sprintf("orphan: key %q (id=%s) has no .md file", key, t.ID))
			result.PreservedCount++
		}
	}

	// 6.5 Generate stage-gates (skip for docs-only features)
	var generated int
	if !docsOnly {
		// Collect all task IDs from the current index for phase detection.
		var allTaskIDs []string
		for _, t := range index.TasksMap() {
			allTaskIDs = append(allTaskIDs, t.ID)
		}
		generated, err = GenerateStageGates(allTaskIDs, opts.TasksDir, opts.FeatureSlug)
		if err != nil {
			return nil, fmt.Errorf("generate stage-gates: %w", err)
		}
		result.StageGatesGenerated = generated
	}

	// Index any newly generated stage-gate files
	if generated > 0 {
		stageEntries, err := os.ReadDir(opts.TasksDir)
		if err != nil {
			return nil, fmt.Errorf("read tasks dir for stage-gates: %w", err)
		}
		for _, entry := range stageEntries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
				continue
			}
			if shouldSkipFile(entry.Name()) {
				continue
			}
			key := strings.TrimSuffix(entry.Name(), ".md")
			if existingKeys[key] {
				continue // already indexed
			}

			filePath := filepath.Join(opts.TasksDir, entry.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				continue
			}
			fm, _, err := ParseFrontmatter(content)
			if err != nil || fm.ID == "" {
				continue
			}

			existingKeys[key] = true
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
			// Preserve runtime state if task already exists in index
			if existing, found := index.ByID(fm.ID); found {
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
	}

	// 7. Generate test tasks or T-eval-doc
	if docsOnly && hasBusinessTasks(index.TasksMap()) {
		// Docs-only: generate T-eval-doc instead of test pipeline
		evalTask := GetDocEvalTask()
		ResolveDocEvalDep(&evalTask, index.TasksMap())

		evalKey := evalTask.Key
		existingKeys[evalKey] = true

		// Generate .md if missing
		mdPath := filepath.Join(opts.TasksDir, evalKey+".md")
		if _, err := os.Stat(mdPath); os.IsNotExist(err) {
			evalContent, err := GenerateTestTaskMD(evalTask, opts.FeatureSlug)
			if err != nil {
				result.Warnings = append(result.Warnings, fmt.Sprintf("generate %s: %v", evalKey, err))
			} else if err := os.WriteFile(mdPath, evalContent, 0644); err != nil {
				result.Warnings = append(result.Warnings, fmt.Sprintf("write %s: %v", evalKey, err))
			}
		}

		task := evalTask.TaskFromFile()
		if existing, found := index.ByID(evalTask.ID); found {
			task.Status = existing.Status
			task.SourceTaskID = existing.SourceTaskID
			task.BlockedReason = existing.BlockedReason
			index.SetTask(evalKey, task)
			result.UpdatedCount++
		} else {
			index.SetTask(evalKey, task)
			result.NewCount++
		}
	} else if len(profiles) > 0 && mode != "" {
		// Detect test types from test-cases.md
		var detectedTypes []string
		testCasesPath := filepath.Join(opts.ProjectRoot, "docs", "features", opts.FeatureSlug, "testing", "test-cases.md")
		if tcData, err := os.ReadFile(testCasesPath); err == nil {
			detectedTypes = DetectTypesFromTestCases(tcData)
		}

		testTasks := generateTestTasks(mode, profiles, detectedTypes)
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
			_ = os.WriteFile(filePath, normalized, 0644)
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
func generateTestTasks(mode string, profiles []string, detectedTypes []string) []TestTaskDef {
	switch mode {
	case "breakdown":
		return GetBreakdownTestTasks(profiles, detectedTypes)
	case "quick":
		return GetQuickTestTasks(profiles, detectedTypes)
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

// isDocsOnlyFeature returns true when all business tasks have types that are
// neither implementation nor fix. Auto-generated tasks (gates, summaries, test
// pipeline tasks) are excluded from the check.
// An empty task map returns true.
func isDocsOnlyFeature(tasks map[string]Task) bool {
	for _, t := range tasks {
		if isAutoGenTaskID(t.ID) {
			continue
		}
		if t.Type == TypeImplementation || t.Type == TypeFix {
			return false
		}
	}
	return true
}

// hasBusinessTasks returns true if there is at least one non-auto-generated task.
func hasBusinessTasks(tasks map[string]Task) bool {
	for _, t := range tasks {
		if !isAutoGenTaskID(t.ID) {
			return true
		}
	}
	return false
}

// isAutoGenTaskID returns true for task IDs that are auto-generated
// (test pipeline, gates, summaries, T-eval-doc).
// fix- and disc- tasks are NOT auto-generated; they are business tasks (they modify code).
func isAutoGenTaskID(id string) bool {
	if isTestTaskID(id) {
		return true
	}
	if id == "T-eval-doc" {
		return true
	}
	if strings.HasSuffix(id, ".gate") || strings.HasSuffix(id, ".summary") {
		return true
	}
	return false
}

// loadIndexFromBytes deserializes index JSON.
func loadIndexFromBytes(data []byte) (*TaskIndex, error) {
	var idx TaskIndex
	if err := idx.UnmarshalJSON(data); err != nil {
		return nil, err
	}
	return &idx, nil
}
