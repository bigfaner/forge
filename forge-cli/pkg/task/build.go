package task

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"forge-cli/pkg/forgeconfig"
	indexPkg "forge-cli/pkg/index"
)

// BuildIndexOpts holds options for building the task index.
type BuildIndexOpts struct {
	FeatureSlug string
	ProjectRoot string
	TasksDir    string                 // absolute path to tasks/
	IndexPath   string                 // absolute path to index.json
	AutoConfig  forgeconfig.AutoConfig // auto-behavior config (defaults filled by caller)
}

// BuildIndexResult holds the result of a BuildIndex operation.
type BuildIndexResult struct {
	NewCount            int
	UpdatedCount        int
	PreservedCount      int
	StageGatesGenerated int
	Warnings            []string
	Index               *TaskIndex
}

// BuildIndex scans .md files and generates/updates index.json.
// It is idempotent: re-running with no changes produces the same output.
func BuildIndex(opts BuildIndexOpts) (*BuildIndexResult, error) {
	result := &BuildIndexResult{}

	// Apply defaults only when AutoConfig is completely zero (nothing was loaded from config).
	// WithDefaults() is NOT called here because it cannot distinguish "user explicitly set
	// both fields to false" from "field was never set" — both equal ModeToggle{}.
	// The caller (forgeconfig.ReadAutoConfig) already applies per-field defaults correctly
	// using raw YAML field tracking. When config is loaded, all fields are explicitly set.
	if opts.AutoConfig.IsZero() {
		opts.AutoConfig = forgeconfig.AutoConfigDefaults()
	}

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

	// 3.5 Build BodyContext from planning-time data (proposal/PRD + config)
	surfaces, _ := forgeconfig.ReadSurfaces(opts.ProjectRoot)
	// Validate surfaces: log warnings for unknown types, filter them out
	forgeconfig.ValidateSurfaceTypes(surfaces)
	capabilities := forgeconfig.SurfaceTypes(surfaces)
	bodyCtx := extractBodyContext(opts.ProjectRoot, opts.FeatureSlug, mode, capabilities)

	// 4. Profiles and interfaces resolved by caller (task 1.4)
	// BuildIndex no longer holds Languages/TestInterfaces; caller injects them into generateTestTasks.

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
			Type:          taskType,
			Coverage:      fm.Coverage,
		}

		// Merge with existing
		if existing, found := index.ByID(fm.ID); found {
			PreserveRuntimeFields(&existing, &newTask)
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
		if IsAutoGenTaskID(t.ID) {
			continue // auto-generated tasks use InferType
		}
		if t.Type == "" {
			return nil, fmt.Errorf("task %s has empty type (set type in frontmatter or use a recognizable ID pattern)", t.File)
		}
	}

	// 5.5.0 Validate: non-auto-gen tasks must not use system types
	for key, t := range index.TasksMap() {
		if !existingKeys[key] {
			continue // auto-generated tasks are not from .md files
		}
		if IsAutoGenTaskID(t.ID) {
			continue // auto-generated tasks are allowed to use system types
		}
		if IsSystemType(t.Type) {
			return nil, fmt.Errorf("task '%s': type '%s' is a system-reserved type (reserved: %s)", t.ID, t.Type, FormatSystemTypes())
		}
	}

	// 5.5.1 Detect pipeline needs
	needsTest := needsTestPipeline(index.TasksMap())
	needsEval := needsReviewDoc(index.TasksMap())

	// 6. Detect and clean up orphans
	// Test tasks (T-test-*) are also cleaned up when they have no .md file AND
	// are not regenerated by Step 7.5 (config disabled test pipeline).
	// Fix tasks (fix-*) are excluded — they are business tasks without auto-generated .md files.
	for key, t := range index.TasksMap() {
		if !existingKeys[key] && !isFixTaskID(t.ID) {
			result.Warnings = append(result.Warnings, fmt.Sprintf("orphan: key %q (id=%s) has no .md file", key, t.ID))
			delete(index.TasksMap(), key)
			result.PreservedCount++
		}
	}

	// 6.5 Generate stage-gates (skip for features without testable types)
	var generated int
	if needsTest {
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
				Type:          taskType,
				Coverage:      fm.Coverage,
			}
			// Preserve runtime state if task already exists in index
			if existing, found := index.ByID(fm.ID); found {
				PreserveRuntimeFields(&existing, &newTask)
				index.SetTask(key, newTask)
				result.UpdatedCount++
			} else {
				newTask.Status = "pending"
				index.SetTask(key, newTask)
				result.NewCount++
			}
		}
	}

	// 7. Generate test tasks or T-review-doc
	if needsEval {
		// Docs-only: generate T-review-doc instead of test pipeline
		evalTask := GetReviewDocTask()
		ResolveReviewDocDep(&evalTask, index.TasksMap())

		evalKey := evalTask.Key
		existingKeys[evalKey] = true

		// Generate .md if missing
		mdPath := filepath.Join(opts.TasksDir, evalKey+".md")
		if _, err := os.Stat(mdPath); os.IsNotExist(err) {
			evalContent, err := GenerateTestTaskMD(evalTask, bodyCtx)
			if err != nil {
				result.Warnings = append(result.Warnings, fmt.Sprintf("generate %s: %v", evalKey, err))
			} else if err := os.WriteFile(mdPath, evalContent, 0644); err != nil {
				result.Warnings = append(result.Warnings, fmt.Sprintf("write %s: %v", evalKey, err))
			}
		}

		task := evalTask.TaskFromFile()
		if existing, found := index.ByID(evalTask.ID); found {
			PreserveRuntimeFields(&existing, &task)
			index.SetTask(evalKey, task)
			result.UpdatedCount++
		} else {
			index.SetTask(evalKey, task)
			result.NewCount++
		}
	}

	// 7.5 Generate test pipeline tasks
	if needsTest && mode != "" {
		if len(capabilities) == 0 {
			return nil, fmt.Errorf("no surfaces configured in .forge/config.yaml. Run `forge init` to configure surfaces")
		}
		testTasks := GenerateTestTasks(mode, capabilities, opts.AutoConfig)
		for _, td := range testTasks {
			ttKey := td.Key
			existingKeys[ttKey] = true

			// Generate .md if missing
			mdPath := filepath.Join(opts.TasksDir, ttKey+".md")
			if _, err := os.Stat(mdPath); os.IsNotExist(err) {
				content, genErr := GenerateTestTaskMD(td, bodyCtx)
				if genErr != nil {
					result.Warnings = append(result.Warnings, fmt.Sprintf("generate %s: %v", ttKey, genErr))
					continue
				}
				if writeErr := os.WriteFile(mdPath, content, 0644); writeErr != nil {
					result.Warnings = append(result.Warnings, fmt.Sprintf("write %s: %v", ttKey, writeErr))
					continue
				}
			}

			t := td.TaskFromFile()
			if existing, found := index.ByID(td.ID); found {
				PreserveRuntimeFields(&existing, &t)
				index.SetTask(ttKey, t)
				result.UpdatedCount++
			} else {
				index.SetTask(ttKey, t)
				result.NewCount++
			}
		}

		// Resolve first-test-task dependency to the highest business gate.
		// This MUST run after PreserveRuntimeFields, which otherwise overwrites
		// the correctly computed deps with stale values from the previous index.
		ResolveFirstTestDep(testTasks, index.TasksMap(), mode)

		// For mixed features (needsEval + needsTest), inject T-review-doc as a
		// dependency of the first test pipeline task. This ensures review-doc
		// executes before test generation, so tests are based on reviewed docs.
		firstTestIdx := findFirstTestTaskIdx(testTasks)
		if needsEval && firstTestIdx >= 0 {
			testTasks[firstTestIdx].Dependencies = append(
				[]string{"T-review-doc"}, testTasks[firstTestIdx].Dependencies...)
		}

		// Write the modified first-test-task deps back to the index.
		if firstTestIdx >= 0 {
			firstKey := testTasks[firstTestIdx].Key
			if t, found := index.ByID(testTasks[firstTestIdx].ID); found {
				t.Dependencies = testTasks[firstTestIdx].Dependencies
				index.SetTask(firstKey, t)
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
	if err := indexPkg.SaveIndexAtomic(opts.IndexPath, index); err != nil {
		return nil, fmt.Errorf("save index: %w", err)
	}

	result.Index = index
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

// GenerateTestTasks returns test task definitions for the given mode and interfaces.
// Exported for use by caller (task 1.4).
func GenerateTestTasks(mode string, capabilities []string, auto forgeconfig.AutoConfig) []AutoGenTaskDef {
	switch mode {
	case "breakdown":
		return GetBreakdownTestTasks(capabilities, auto)
	case "quick":
		return GetQuickTestTasks(capabilities, auto)
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
	return strings.HasPrefix(id, "T-test-") ||
		strings.HasPrefix(id, "T-quick-") ||
		strings.HasPrefix(id, "T-specs-") ||
		strings.HasPrefix(id, "T-clean-") ||
		strings.HasPrefix(id, "T-validate-") ||
		strings.HasPrefix(id, "T-eval-")
}

// isFixTaskID returns true for fix-task IDs (auto-generated by quality gate pipeline).
func isFixTaskID(id string) bool {
	return strings.HasPrefix(id, "fix-")
}

// IsTestableType returns true if the given task type has testable runtime behavior.
// Covers coding.* prefix and code-quality.simplify (explicit match).
func IsTestableType(typ string) bool {
	return strings.HasPrefix(typ, "coding.") || typ == TypeCleanCode
}

// needsTestPipeline returns true when any non-auto-gen task has a testable
// runtime behavior type (feature, enhancement, or fix).
// An empty task map returns false.
func needsTestPipeline(tasks map[string]Task) bool {
	for _, t := range tasks {
		if IsAutoGenTaskID(t.ID) {
			continue
		}
		if IsTestableType(t.Type) {
			return true
		}
	}
	return false
}

// needsReviewDoc returns true when ANY non-auto-gen task has a doc-category type.
// Uses CategoryForType to check: this covers TypeDoc, TypeDocConsolidate, TypeDocDrift, etc.
// Doc subtypes that are system types (doc.review, doc.summary) cannot appear as business tasks
// due to system-type validation, so they are effectively excluded.
// An empty task map returns false.
func needsReviewDoc(tasks map[string]Task) bool {
	for _, t := range tasks {
		if IsAutoGenTaskID(t.ID) {
			continue
		}
		if CategoryForType(t.Type) == CategoryDoc {
			return true
		}
	}
	return false
}

// findFirstTestTaskIdx returns the index of the first test pipeline task in the
// generated task list. For breakdown mode this is T-eval-journey; for quick mode
// this is the first T-quick-gen-and-run task.
func findFirstTestTaskIdx(tasks []AutoGenTaskDef) int {
	// Check breakdown mode first (T-eval-journey)
	for i, t := range tasks {
		if t.ID == "T-eval-journey" {
			return i
		}
	}
	// Quick mode: T-quick-gen-and-run*
	for i, t := range tasks {
		if strings.HasPrefix(t.ID, "T-quick-gen-and-run") {
			return i
		}
	}
	// Fallback: first task
	if len(tasks) > 0 {
		return 0
	}
	return -1
}

// IsAutoGenTaskID returns true for task IDs that are auto-generated
// (test pipeline, gates, summaries, T-review-doc).
// fix- and disc- tasks are NOT auto-generated; they are business tasks (they modify code).
func IsAutoGenTaskID(id string) bool {
	if isTestTaskID(id) {
		return true
	}
	if id == "T-review-doc" {
		return true
	}
	if strings.HasSuffix(id, IDSuffixGate) || strings.HasSuffix(id, IDSuffixSummary) {
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
