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
	"forge-cli/pkg/types"
)

// BuildIndexOpts holds options for building the task index.
type BuildIndexOpts struct {
	FeatureSlug string
	ProjectRoot string
	TasksDir    string                 // absolute path to tasks/
	IndexPath   string                 // absolute path to index.json
	AutoConfig  forgeconfig.AutoConfig // auto-behavior config (defaults filled by caller)
	Intent      string                 // feature intent: "new-feature" (default), "refactor", "cleanup"
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

	// 1.5 Resolve intent (default to new-feature if not set)
	intent := opts.Intent
	if intent == "" {
		intent = "new-feature"
	}

	// 2. Detect mode
	mode := detectMode(opts.ProjectRoot, opts.FeatureSlug, intent)

	// 3. Set feature metadata
	setFeatureMetadata(index, opts.ProjectRoot, opts.FeatureSlug)

	// 3.5 Build BodyContext from planning-time data (proposal/PRD + config)
	surfaces, _ := forgeconfig.ReadSurfaces(opts.ProjectRoot)
	// Validate surfaces: log warnings for unknown types, filter them out
	forgeconfig.ValidateSurfaceTypes(surfaces)
	capabilities := forgeconfig.SurfaceTypes(surfaces)
	executionOrder, _ := forgeconfig.ReadExecutionOrder(opts.ProjectRoot)
	bodyCtx := extractBodyContext(opts.ProjectRoot, opts.FeatureSlug, mode, capabilities)

	// 4. Profiles and surfaces resolved by caller (task 1.4)
	// BuildIndex no longer holds Languages/Surfaces; caller injects them into generateTestTasks.

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
			taskType = InferType(fm.ID, surfaces)
		}

		newTask := Task{
			ID:            fm.ID,
			Title:         fm.Title,
			Priority:      types.Priority(fm.Priority),
			EstimatedTime: fm.EstimatedTime,
			Dependencies:  fm.Dependencies,
			File:          entry.Name(),
			Record:        path.Join("records", entry.Name()),
			Breaking:      fm.Breaking,
			SurfaceKey:    fm.SurfaceKey,
			SurfaceType:   fm.SurfaceType,
			MainSession:   fm.MainSession,
			Type:          taskType,
			Coverage:      fm.Coverage,
			Complexity:    fm.Complexity,
		}

		// Merge with existing
		if existing, found := index.ByID(fm.ID); found {
			PreserveRuntimeFields(&existing, &newTask)
			index.SetTask(key, newTask)
			result.UpdatedCount++
		} else {
			newTask.Status = types.StatusPending
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
	needsTest := needsTestPipeline(index.TasksMap(), intent)
	needsEval := needsReviewDoc(index.TasksMap())

	// 5.5.2 Extract AC from doc tasks for review-doc template
	if needsEval {
		bodyCtx.DocTaskCriteria = extractDocTaskCriteria(opts.TasksDir)

		// Validate: warn about doc tasks without AC section
		for taskName, acContent := range bodyCtx.DocTaskCriteria {
			if acContent == "" {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("[WARN] task %s has no Acceptance Criteria section", taskName))
			}
		}

		// Validate: check that DocTaskCriteria keys cover all doc tasks
		// Collect doc task keys from the index (non-auto-gen, doc-category)
		docTaskKeys := make(map[string]bool)
		for key, t := range index.TasksMap() {
			if IsAutoGenTaskID(t.ID) {
				continue
			}
			if CategoryForType(t.Type) == CategoryDoc {
				docTaskKeys[key] = true
			}
		}
		// Verify every doc task has an entry in DocTaskCriteria
		for key := range docTaskKeys {
			if _, ok := bodyCtx.DocTaskCriteria[key]; !ok {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("[WARN] task %s has no Acceptance Criteria section", key))
			}
		}

		// Special warning when ALL doc tasks lack AC (zero-AC feature)
		allEmpty := true
		for _, acContent := range bodyCtx.DocTaskCriteria {
			if strings.TrimSpace(acContent) != "" {
				allEmpty = false
				break
			}
		}
		if len(bodyCtx.DocTaskCriteria) > 0 && allEmpty {
			result.Warnings = append(result.Warnings,
				"[WARN] feature has no AC for any doc task")
		}
	}

	// 5.9 Migrate fix-tasks from legacy T-test-run to per-surface-key T-test-run-{key}
	// Only applies to multi-surface projects. Single-surface projects keep T-test-run unchanged.
	// Must run BEFORE orphan cleanup (step 6) so the old T-test-run entry is still available.
	var migratedState *migratedRunTestState
	if !isSingleSurface(surfaces) {
		resolvedEO, _ := forgeconfig.ResolveExecutionOrder(surfaces, executionOrder)
		migratedState = migrateFixTaskSources(index, surfaces, resolvedEO)
	}

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
				taskType = InferType(fm.ID, surfaces)
			}
			newTask := Task{
				ID:            fm.ID,
				Title:         fm.Title,
				Priority:      types.Priority(fm.Priority),
				EstimatedTime: fm.EstimatedTime,
				Dependencies:  fm.Dependencies,
				File:          entry.Name(),
				Record:        path.Join("records", entry.Name()),
				Breaking:      fm.Breaking,
				SurfaceKey:    fm.SurfaceKey,
				SurfaceType:   fm.SurfaceType,
				MainSession:   fm.MainSession,
				Type:          taskType,
				Coverage:      fm.Coverage,
				Complexity:    fm.Complexity,
			}
			// Preserve runtime state if task already exists in index
			if existing, found := index.ByID(fm.ID); found {
				PreserveRuntimeFields(&existing, &newTask)
				index.SetTask(key, newTask)
				result.UpdatedCount++
			} else {
				newTask.Status = types.StatusPending
				index.SetTask(key, newTask)
				result.NewCount++
			}
		}
	}

	// 7. Generate auto-gen tasks via registry (T-review-doc + test pipeline + validation + consolidation + clean-code)
	// The registry handles mode/config/intent/condition filtering, expansion, and dependency resolution.
	// Step 7 (T-review-doc only) and Step 7.5 (test pipeline only) are now unified.
	if mode != "" && (needsTest || needsEval) {
		// Note: surfaces may be empty for docs-only features — that's fine.
		// GenerateTestTasks generates non-surface tasks (T-review-doc, T-clean-code,
		// T-validate-code) even with empty surfaces; surface-dependent nodes produce zero tasks.

		// Collect business tasks for GenerateCondition and dependency resolvers
		var businessTasks []Task
		for _, t := range index.TasksMap() {
			if !IsAutoGenTaskID(t.ID) {
				businessTasks = append(businessTasks, t)
			}
		}

		resolvedExecOrder, _ := forgeconfig.ResolveExecutionOrder(surfaces, executionOrder)
		testTasks := GenerateTestTasks(mode, surfaces, resolvedExecOrder, opts.AutoConfig, intent, businessTasks, index.TasksMap())

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

		// Phase 2 of migration: apply saved T-test-run state to the first run-test task
		if migratedState != nil {
			applyMigratedRunTestState(index, migratedState)
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

	// 8.5 Check for legacy scope fields that require migration
	var allTasks []Task
	for _, t := range index.TasksMap() {
		allTasks = append(allTasks, t)
	}
	if err := CheckLegacyScope(allTasks); err != nil {
		return nil, err
	}

	// 9. Save index
	if err := indexPkg.SaveIndexAtomic(opts.IndexPath, index); err != nil {
		return nil, fmt.Errorf("save index: %w", err)
	}

	result.Index = index
	return result, nil
}

// detectMode determines the feature mode from file existence and intent.
// When intent is "cleanup", it forces Quick mode regardless of document existence.
func detectMode(projectRoot, slug, intent string) string {
	// cleanup intent always forces Quick mode, ignoring document existence
	if intent == "cleanup" {
		return "quick"
	}

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

// isTestTaskID returns true for auto-generated pipeline task IDs (excluding
// gate/summary/T-review-doc). Derived from PipelineRegistry: all registry nodes
// with a "T-" prefix are test pipeline tasks.
func isTestTaskID(id string) bool {
	if !strings.HasPrefix(id, "T-") {
		return false
	}
	if id == "T-review-doc" {
		return false // review-doc is auto-gen but not a test task
	}
	if strings.HasSuffix(id, IDSuffixGate) || strings.HasSuffix(id, IDSuffixSummary) {
		return false // gate/summary handled separately
	}
	// Check if the ID matches any registry node pattern
	return matchRegistryID(id, nil) != ""
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
// When intent is "refactor" or "cleanup", it returns false immediately
// without iterating tasks — these intents skip the test pipeline entirely.
// An empty task map returns false.
func needsTestPipeline(tasks map[string]Task, intent string) bool {
	// Intent short-circuit: refactor/cleanup skip test pipeline entirely
	if intent == "refactor" || intent == "cleanup" {
		return false
	}

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

// migratedRunTestState holds the saved state from a legacy T-test-run entry
// that needs to be applied to the new T-test-run-{first-surface-key} task.
type migratedRunTestState struct {
	FirstRunTestID string // new target ID (e.g. "T-test-run-backend")
	Status         string
	BlockedReason  string
}

// migrateFixTaskSources remaps fix-tasks with SourceTaskID "T-test-run" to the new
// per-surface-key task ID (T-test-run-{first-surface-key}) in multi-surface projects.
// It also removes the old T-test-run entry from the index and returns the saved state
// so it can be applied to the first run-test task after test pipeline generation.
// Returns nil if no migration is needed.
func migrateFixTaskSources(index *TaskIndex, _ map[string]string, executionOrder []string) *migratedRunTestState {
	if len(executionOrder) == 0 {
		return nil
	}

	// Find all fix-tasks with SourceTaskID "T-test-run"
	var fixKeys []string
	for key, t := range index.TasksMap() {
		if t.SourceTaskID == "T-test-run" {
			fixKeys = append(fixKeys, key)
		}
	}

	if len(fixKeys) == 0 {
		return nil
	}

	// Find and save the old T-test-run entry's state
	firstSurfaceKey := executionOrder[0]
	newRunTestID := "T-test-run-" + firstSurfaceKey

	var state migratedRunTestState
	state.FirstRunTestID = newRunTestID

	// Look for the old T-test-run entry by iterating all tasks
	for key, t := range index.TasksMap() {
		if t.ID == "T-test-run" {
			state.Status = string(t.Status)
			state.BlockedReason = t.BlockedReason
			// Remove the old entry
			delete(index.TasksMap(), key)
			break
		}
	}

	// Remap fix-task SourceTaskIDs
	for _, key := range fixKeys {
		task, exists := index.TasksMap()[key]
		if !exists {
			continue
		}
		task.SourceTaskID = newRunTestID
		index.SetTask(key, task)
	}

	return &state
}

// applyMigratedRunTestState applies the saved state from a migrated T-test-run
// to the first per-surface-key run-test task in the index.
func applyMigratedRunTestState(index *TaskIndex, state *migratedRunTestState) {
	if state == nil || state.FirstRunTestID == "" {
		return
	}

	// Find the first run-test task by ID
	for key, t := range index.TasksMap() {
		if t.ID == state.FirstRunTestID {
			// Apply saved state only when the new task is still in pending status
			// (don't override if it already has a runtime state)
			if t.Status == types.StatusPending && state.Status != "" {
				t.Status = types.Status(state.Status)
				t.BlockedReason = state.BlockedReason
				index.SetTask(key, t)
			}
			return
		}
	}
}

// loadIndexFromBytes deserializes index JSON.
func loadIndexFromBytes(data []byte) (*TaskIndex, error) {
	var idx TaskIndex
	if err := idx.UnmarshalJSON(data); err != nil {
		return nil, err
	}
	return &idx, nil
}
