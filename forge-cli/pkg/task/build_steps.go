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

// buildContext holds shared state across BuildIndex step functions.
type buildContext struct {
	opts         BuildIndexOpts
	result       *BuildIndexResult
	index        *TaskIndex
	intent       string
	mode         string
	surfaces     map[string]string
	capabilities []string
	execOrder    []string
	bodyCtx      BodyContext
	existingKeys map[string]bool
	needsTest    bool
	needsEval    bool
	migrated     *migratedRunTestState
	entries      []os.DirEntry
}

// loadOrCreateIndex loads an existing index from disk or creates a new one.
func (bc *buildContext) loadOrCreateIndex() error {
	if data, err := os.ReadFile(bc.opts.IndexPath); err == nil {
		idx, err := loadIndexFromBytes(data)
		if err != nil {
			return fmt.Errorf("load existing index: %w", err)
		}
		bc.index = idx
	} else {
		bc.index = NewTaskIndex(bc.opts.FeatureSlug)
	}
	return nil
}

// resolveContext detects mode, sets metadata, and builds body context.
func (bc *buildContext) resolveContext() {
	// Resolve intent (default to new-feature if not set)
	bc.intent = bc.opts.Intent
	if bc.intent == "" {
		bc.intent = "new-feature"
	}

	// Detect mode
	bc.mode = detectMode(bc.opts.ProjectRoot, bc.opts.FeatureSlug, bc.intent)

	// Set feature metadata
	setFeatureMetadata(bc.index, bc.opts.ProjectRoot, bc.opts.FeatureSlug)

	// Build BodyContext from planning-time data
	bc.surfaces, _ = forgeconfig.ReadSurfaces(bc.opts.ProjectRoot)
	forgeconfig.ValidateSurfaceTypes(bc.surfaces)
	bc.capabilities = forgeconfig.SurfaceTypes(bc.surfaces)
	bc.execOrder, _ = forgeconfig.ReadExecutionOrder(bc.opts.ProjectRoot)
	bc.bodyCtx = extractBodyContext(bc.opts.ProjectRoot, bc.opts.FeatureSlug, bc.mode, bc.capabilities)
}

// scanTaskFiles reads .md files from the tasks directory and upserts them into the index.
func (bc *buildContext) scanTaskFiles() error {
	bc.existingKeys = make(map[string]bool)

	entries, err := os.ReadDir(bc.opts.TasksDir)
	if err != nil {
		return fmt.Errorf("read tasks dir: %w", err)
	}
	bc.entries = entries

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		if shouldSkipFile(entry.Name()) {
			continue
		}
		if err := bc.upsertTaskFromFile(entry.Name()); err != nil {
			bc.result.Warnings = append(bc.result.Warnings, err.Error())
		}
	}
	return nil
}

// upsertTaskFromFile parses a single .md file and inserts or updates the task in the index.
func (bc *buildContext) upsertTaskFromFile(filename string) error {
	filePath := filepath.Join(bc.opts.TasksDir, filename)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read %s: %v", filename, err)
	}

	fm, _, err := ParseFrontmatter(content)
	if err != nil {
		return fmt.Errorf("parse %s: %v", filename, err)
	}

	if fm.ID == "" {
		return fmt.Errorf("skip %s: no id in frontmatter", filename)
	}

	key := strings.TrimSuffix(filename, ".md")
	bc.existingKeys[key] = true

	taskType := fm.Type
	if taskType == "" {
		taskType = InferType(fm.ID, bc.surfaces)
	}

	newTask := Task{
		ID:            fm.ID,
		Title:         fm.Title,
		Priority:      types.Priority(fm.Priority),
		EstimatedTime: fm.EstimatedTime,
		Dependencies:  fm.Dependencies,
		File:          filename,
		Record:        path.Join("records", filename),
		Breaking:      fm.Breaking,
		SurfaceKey:    fm.SurfaceKey,
		SurfaceType:   fm.SurfaceType,
		MainSession:   fm.MainSession,
		Type:          taskType,
		Coverage:      fm.Coverage,
		Complexity:    fm.Complexity,
	}

	if existing, found := bc.index.ByID(fm.ID); found {
		PreserveRuntimeFields(&existing, &newTask)
		bc.index.SetTask(key, newTask)
		bc.result.UpdatedCount++
	} else {
		newTask.Status = types.StatusPending
		bc.index.SetTask(key, newTask)
		bc.result.NewCount++
	}
	return nil
}

// validateTaskTypes checks that business tasks have a type and do not use system types.
func (bc *buildContext) validateTaskTypes() error {
	for key, t := range bc.index.TasksMap() {
		if !bc.existingKeys[key] {
			continue
		}
		if IsAutoGenTaskID(t.ID) {
			continue
		}
		if t.Type == "" {
			return fmt.Errorf("task %s has empty type (set type in frontmatter or use a recognizable ID pattern)", t.File)
		}
		if IsSystemType(t.Type) {
			return fmt.Errorf("task '%s': type '%s' is a system-reserved type (reserved: %s)", t.ID, t.Type, FormatSystemTypes())
		}
	}
	return nil
}

// detectPipelineNeeds determines which pipelines are needed and extracts doc task AC.
func (bc *buildContext) detectPipelineNeedsAndAC() {
	bc.needsTest = needsTestPipeline(bc.index.TasksMap(), bc.intent)
	bc.needsEval = needsReviewDoc(bc.index.TasksMap())

	if !bc.needsEval {
		return
	}

	bc.bodyCtx.DocTaskCriteria = extractDocTaskCriteria(bc.opts.TasksDir)
	bc.validateDocTaskAC()
}

// validateDocTaskAC emits warnings for doc tasks missing Acceptance Criteria.
func (bc *buildContext) validateDocTaskAC() {
	for taskName, acContent := range bc.bodyCtx.DocTaskCriteria {
		if acContent == "" {
			bc.result.Warnings = append(bc.result.Warnings,
				fmt.Sprintf("[WARN] task %s has no Acceptance Criteria section", taskName))
		}
	}

	// Check that DocTaskCriteria keys cover all doc tasks
	docTaskKeys := make(map[string]bool)
	for key, t := range bc.index.TasksMap() {
		if IsAutoGenTaskID(t.ID) {
			continue
		}
		if CategoryForType(t.Type) == CategoryDoc {
			docTaskKeys[key] = true
		}
	}
	for key := range docTaskKeys {
		if _, ok := bc.bodyCtx.DocTaskCriteria[key]; !ok {
			bc.result.Warnings = append(bc.result.Warnings,
				fmt.Sprintf("[WARN] task %s has no Acceptance Criteria section", key))
		}
	}

	// Special warning when ALL doc tasks lack AC
	allEmpty := true
	for _, acContent := range bc.bodyCtx.DocTaskCriteria {
		if strings.TrimSpace(acContent) != "" {
			allEmpty = false
			break
		}
	}
	if len(bc.bodyCtx.DocTaskCriteria) > 0 && allEmpty {
		bc.result.Warnings = append(bc.result.Warnings,
			"[WARN] feature has no AC for any doc task")
	}
}

// cleanupOrphans migrates fix-task sources and removes orphaned index entries.
func (bc *buildContext) cleanupOrphans() {
	// Migrate fix-tasks from legacy T-test-run to per-surface-key
	if !isSingleSurface(bc.surfaces) {
		resolvedEO, _ := forgeconfig.ResolveExecutionOrder(bc.surfaces, bc.execOrder)
		bc.migrated = migrateFixTaskSources(bc.index, bc.surfaces, resolvedEO)
	}

	// Remove orphaned entries (no .md file, not a fix task)
	for key, t := range bc.index.TasksMap() {
		if !bc.existingKeys[key] && !isFixTaskID(t.ID) {
			bc.result.Warnings = append(bc.result.Warnings,
				fmt.Sprintf("orphan: key %q (id=%s) has no .md file", key, t.ID))
			delete(bc.index.TasksMap(), key)
			bc.result.PreservedCount++
		}
	}
}

// generateAndIndexStageGates creates stage-gate .md files and indexes them.
func (bc *buildContext) generateAndIndexStageGates() error {
	if !bc.needsTest {
		return nil
	}

	var allTaskIDs []string
	for _, t := range bc.index.TasksMap() {
		allTaskIDs = append(allTaskIDs, t.ID)
	}
	generated, err := GenerateStageGates(allTaskIDs, bc.opts.TasksDir, bc.opts.FeatureSlug)
	if err != nil {
		return fmt.Errorf("generate stage-gates: %w", err)
	}
	bc.result.StageGatesGenerated = generated

	if generated > 0 {
		return bc.indexNewStageGates()
	}
	return nil
}

// indexNewStageGates scans for newly generated stage-gate files and adds them to the index.
func (bc *buildContext) indexNewStageGates() error {
	stageEntries, err := os.ReadDir(bc.opts.TasksDir)
	if err != nil {
		return fmt.Errorf("read tasks dir for stage-gates: %w", err)
	}

	for _, entry := range stageEntries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		if shouldSkipFile(entry.Name()) {
			continue
		}
		key := strings.TrimSuffix(entry.Name(), ".md")
		if bc.existingKeys[key] {
			continue
		}

		filePath := filepath.Join(bc.opts.TasksDir, entry.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}
		fm, _, err := ParseFrontmatter(content)
		if err != nil || fm.ID == "" {
			continue
		}

		bc.existingKeys[key] = true
		bc.upsertAutoEntry(key, fm, entry.Name())
	}
	return nil
}

// generateAutoGenTasks creates auto-generated tasks via the pipeline registry.
func (bc *buildContext) generateAutoGenTasks() {
	if bc.mode == "" || (!bc.needsTest && !bc.needsEval) {
		return
	}

	// Collect business tasks for condition evaluation
	var businessTasks []Task
	for _, t := range bc.index.TasksMap() {
		if !IsAutoGenTaskID(t.ID) {
			businessTasks = append(businessTasks, t)
		}
	}

	// Docs-only features: suppress surfaces to avoid surface-dependent tasks
	effectiveSurfaces := bc.surfaces
	if !bc.needsTest {
		effectiveSurfaces = nil
	}

	resolvedExecOrder, _ := forgeconfig.ResolveExecutionOrder(effectiveSurfaces, bc.execOrder)
	testTasks := GenerateTestTasks(bc.mode, effectiveSurfaces, resolvedExecOrder, bc.opts.AutoConfig, bc.intent, businessTasks, bc.index.TasksMap())

	for _, td := range testTasks {
		ttKey := td.Key
		bc.existingKeys[ttKey] = true

		// Generate .md if missing
		mdPath := filepath.Join(bc.opts.TasksDir, ttKey+".md")
		if _, err := os.Stat(mdPath); os.IsNotExist(err) {
			content, genErr := GenerateTestTaskMD(td, bc.bodyCtx)
			if genErr != nil {
				bc.result.Warnings = append(bc.result.Warnings, fmt.Sprintf("generate %s: %v", ttKey, genErr))
				continue
			}
			if writeErr := os.WriteFile(mdPath, content, 0o644); writeErr != nil {
				bc.result.Warnings = append(bc.result.Warnings, fmt.Sprintf("write %s: %v", ttKey, writeErr))
				continue
			}
		}

		t := td.TaskFromFile()
		if existing, found := bc.index.ByID(td.ID); found {
			PreserveRuntimeFields(&existing, &t)
			bc.index.SetTask(ttKey, t)
			bc.result.UpdatedCount++
		} else {
			bc.index.SetTask(ttKey, t)
			bc.result.NewCount++
		}
	}

	// Phase 2 of migration: apply saved state to the first run-test task
	if bc.migrated != nil {
		applyMigratedRunTestState(bc.index, bc.migrated)
	}
}

// normalizeAndSave normalizes task files, checks legacy scope, and saves the index.
func (bc *buildContext) normalizeAndSave() error {
	// Normalize task files (remove empty ## Hard Rules sections)
	for _, entry := range bc.entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		if shouldSkipFile(entry.Name()) {
			continue
		}
		filePath := filepath.Join(bc.opts.TasksDir, entry.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}
		normalized := NormalizeTaskMD(content)
		if !bytes.Equal(normalized, content) {
			_ = os.WriteFile(filePath, normalized, 0o644)
		}
	}

	// Check for legacy scope fields
	var allTasks []Task
	for _, t := range bc.index.TasksMap() {
		allTasks = append(allTasks, t)
	}
	if err := CheckLegacyScope(allTasks); err != nil {
		return err
	}

	// Save index atomically
	if err := indexPkg.SaveIndexAtomic(bc.opts.IndexPath, bc.index); err != nil {
		return fmt.Errorf("save index: %w", err)
	}
	return nil
}

// upsertAutoEntry inserts or updates a stage-gate entry parsed from frontmatter.
func (bc *buildContext) upsertAutoEntry(key string, fm FrontmatterData, filename string) {
	taskType := fm.Type
	if taskType == "" {
		taskType = InferType(fm.ID, bc.surfaces)
	}

	newTask := Task{
		ID:            fm.ID,
		Title:         fm.Title,
		Priority:      types.Priority(fm.Priority),
		EstimatedTime: fm.EstimatedTime,
		Dependencies:  fm.Dependencies,
		File:          filename,
		Record:        path.Join("records", filename),
		Breaking:      fm.Breaking,
		SurfaceKey:    fm.SurfaceKey,
		SurfaceType:   fm.SurfaceType,
		MainSession:   fm.MainSession,
		Type:          taskType,
		Coverage:      fm.Coverage,
		Complexity:    fm.Complexity,
	}

	if existing, found := bc.index.ByID(fm.ID); found {
		PreserveRuntimeFields(&existing, &newTask)
		bc.index.SetTask(key, newTask)
		bc.result.UpdatedCount++
	} else {
		newTask.Status = types.StatusPending
		bc.index.SetTask(key, newTask)
		bc.result.NewCount++
	}
}
