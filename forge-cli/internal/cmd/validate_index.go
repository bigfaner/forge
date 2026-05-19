package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/project"
	"forge-cli/pkg/task"

	"github.com/spf13/cobra"
)

var validateIndexCmd = &cobra.Command{
	Use:   "validate-index [file]",
	Short: "Validate index.json file",
	Long: `Validate an index.json file for structural and semantic correctness.

If no file is specified, validates the current feature's index.json.

Validations:
  - JSON syntax
  - Required fields present
  - Dependency references exist
  - No circular dependencies`,
	Run: runValidateIndex,
}

var (
	validStatus   = map[string]bool{"pending": true, "in_progress": true, "completed": true, "blocked": true, "skipped": true, "rejected": true}
	validPriority = map[string]bool{"P0": true, "P1": true, "P2": true}
)

func runValidateIndex(_ *cobra.Command, args []string) {
	var filePath string
	if len(args) > 0 {
		filePath = args[0]
	} else {
		projectRoot, err := project.FindProjectRoot()
		if err != nil {
			Exit(ErrProjectNotFound())
		}
		slug, err := feature.RequireFeature(projectRoot)
		if err != nil {
			Exit(ErrFeatureNotSet())
		}
		filePath = filepath.Join(projectRoot, feature.GetFeatureIndexFile(slug))
	}

	v := &validator{filePath: filePath}
	if err := v.run(); err != nil {
		Exit(err)
	}
}

type validator struct {
	filePath  string
	quickMode bool // true when Proposal is set instead of PRD+Design
	errors    []string
	warnings  []string
	info      []string
}

func (v *validator) run() error {
	data, err := os.ReadFile(v.filePath)
	if err != nil {
		return ErrFileNotFound(v.filePath)
	}

	var idx task.TaskIndex
	if err := json.Unmarshal(data, &idx); err != nil {
		return ErrInvalidJSON(v.filePath, err.Error())
	}

	if idx.Feature == "" {
		return NewAIError(ErrValidation, "Missing required field", "'feature' field is empty", "Add feature name to index.json", "Add \"feature\": \"<name>\" to index.json")
	}

	v.info = append(v.info, fmt.Sprintf("Feature: %s", idx.Feature))
	v.info = append(v.info, fmt.Sprintf("Tasks: %d", idx.TaskCount()))

	if idx.PRD == "" && idx.Proposal == "" {
		v.warnings = append(v.warnings, "Missing 'prd' field")
	}
	if idx.Design == "" && idx.Proposal == "" {
		v.warnings = append(v.warnings, "Missing 'design' field")
	}

	// Quick mode: Proposal replaces PRD+Design, flat tasks without phases/gates/summaries
	v.quickMode = idx.Proposal != "" && idx.PRD == "" && idx.Design == ""
	if len(idx.StatusEnum) == 0 {
		v.warnings = append(v.warnings, "Missing 'statusEnum' field — task record/status commands may fail")
	}

	v.validateTasks(idx.TasksMap())
	v.validateDependencies(idx.TasksMap())
	v.validateCircularDeps(idx.TasksMap())
	v.validateFilesExist(idx.Feature, idx.TasksMap())
	v.validateWildcardSelfDeps(idx.TasksMap())
	v.validateGateIntegrity(idx.TasksMap())
	v.validatePhaseOrder(idx.TasksMap())
	v.validatePhaseSummaries(idx.TasksMap())
	v.validateLiveness(idx.TasksMap())
	if !v.printResults() {
		return NewAIError(ErrValidation, "Validation failed", fmt.Sprintf("%d errors found", len(v.errors)), "Fix errors in index.json", "cat "+v.filePath)
	}
	return nil
}

func (v *validator) validateTasks(tasks map[string]task.Task) {
	for key, t := range tasks {
		if t.ID == "" {
			v.errors = append(v.errors, fmt.Sprintf("Task '%s': missing 'id'", key))
		}
		if t.Title == "" {
			v.errors = append(v.errors, fmt.Sprintf("Task '%s': missing 'title'", key))
		}
		if t.File == "" {
			v.errors = append(v.errors, fmt.Sprintf("Task '%s': missing 'file'", key))
		}
		if t.Status != "" && !validStatus[t.Status] {
			v.errors = append(v.errors, fmt.Sprintf("Task '%s': invalid status '%s'", key, t.Status))
		}
		if t.Priority != "" && !validPriority[t.Priority] {
			v.errors = append(v.errors, fmt.Sprintf("Task '%s': invalid priority '%s'", key, t.Priority))
		}
		switch {
		case t.Type == "":
			v.errors = append(v.errors, fmt.Sprintf("Task '%s': missing 'type'", key))
		case !task.ValidTypes[t.Type]:
			v.errors = append(v.errors, fmt.Sprintf("Task '%s': invalid type '%s'", key, t.Type))
		}
	}
}

func (v *validator) validateDependencies(tasks map[string]task.Task) {
	taskIDs := make(map[string]bool)
	for _, t := range tasks {
		taskIDs[t.ID] = true
	}

	for key, t := range tasks {
		for _, dep := range t.Dependencies {
			isWildcard := strings.HasSuffix(dep, ".x")

			if isWildcard {
				prefix := strings.TrimSuffix(dep, ".x") + "."
				var matches []string
				for id := range taskIDs {
					if strings.HasPrefix(id, prefix) && isBusinessTask(id) {
						matches = append(matches, id)
					}
				}
				if len(matches) == 0 {
					v.errors = append(v.errors, fmt.Sprintf("Task '%s': wildcard '%s' matches no business tasks", key, dep))
				}
			} else if !taskIDs[dep] {
				v.errors = append(v.errors, fmt.Sprintf("Task '%s': dependency '%s' not found", key, dep))
			}
		}
	}
}

func (v *validator) validateCircularDeps(tasks map[string]task.Task) {
	graph := make(map[string][]string)
	taskIDs := make(map[string]bool)
	for _, t := range tasks {
		taskIDs[t.ID] = true
	}
	for _, t := range tasks {
		for _, dep := range t.Dependencies {
			if !strings.HasSuffix(dep, ".x") && taskIDs[dep] {
				graph[t.ID] = append(graph[t.ID], dep)
			}
		}
	}

	visited := make(map[string]bool)
	inStack := make(map[string]bool)

	var dfs func(id string, path []string) bool
	dfs = func(id string, path []string) bool {
		visited[id] = true
		inStack[id] = true
		for _, neighbor := range graph[id] {
			if !visited[neighbor] {
				if dfs(neighbor, append(path, neighbor)) {
					return true
				}
			} else if inStack[neighbor] {
				v.errors = append(v.errors, fmt.Sprintf("Circular: %s", strings.Join(append(path, neighbor), " -> ")))
				return true
			}
		}
		inStack[id] = false
		return false
	}

	for id := range taskIDs {
		if !visited[id] && dfs(id, []string{id}) {
			break
		}
	}
}

func (v *validator) validateFilesExist(featureSlug string, tasks map[string]task.Task) {
	// Get project root from filePath (which is .../docs/features/<slug>/tasks/index.json)
	// Go up 4 levels: index.json -> tasks -> <slug> -> features -> docs -> projectRoot
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(v.filePath)))))
	tasksDir := filepath.Join(projectRoot, feature.GetFeatureTasksDir(featureSlug))

	for key, t := range tasks {
		if t.File == "" {
			continue
		}
		taskFile := filepath.Join(tasksDir, t.File)
		if _, err := os.Stat(taskFile); os.IsNotExist(err) {
			v.warnings = append(v.warnings, fmt.Sprintf("Task '%s': file '%s' missing", key, t.File))
		}

		// Check for unresolved template placeholders in first-test task files
		if t.ID == "T-test-gen-cases" {
			v.validateFirstTestTaskTemplate(taskFile, "T-test-gen-cases", []string{"{{LAST_BUSINESS_TASK_ID}}", "{{T_TEST_1_DEP}}"})
		}
		if t.ID == "T-quick-gen-cases" {
			v.validateFirstTestTaskTemplate(taskFile, "T-quick-gen-cases", []string{"{{T_QUICK_1_DEP}}"})
		}
	}
}

// validateFirstTestTaskTemplate checks if a first-test task file has unresolved placeholders.
func (v *validator) validateFirstTestTaskTemplate(taskFile string, taskID string, placeholders []string) {
	data, err := os.ReadFile(taskFile)
	if err != nil {
		return // File existence already checked above
	}

	content := string(data)
	for _, ph := range placeholders {
		if strings.Contains(content, ph) {
			v.errors = append(v.errors,
				fmt.Sprintf("Task '%s': file contains unresolved placeholder %s — replace with actual dependency ID", taskID, ph))
		}
	}
}

// isBusinessTask returns true for tasks that are neither gate nor summary.
func isBusinessTask(id string) bool {
	return !strings.HasSuffix(id, ".gate") && !strings.HasSuffix(id, ".summary")
}

// V1: Wildcard self-dependency detection
func (v *validator) validateWildcardSelfDeps(tasks map[string]task.Task) {
	for key, t := range tasks {
		for _, dep := range t.Dependencies {
			if !strings.HasSuffix(dep, ".x") {
				continue
			}
			prefix := strings.TrimSuffix(dep, ".x") + "."
			if !strings.HasPrefix(t.ID, prefix) || !isBusinessTask(t.ID) {
				continue
			}
			// This task's own ID matches the wildcard. Check if other business tasks also match.
			others := 0
			for _, other := range tasks {
				if other.ID != t.ID && strings.HasPrefix(other.ID, prefix) && isBusinessTask(other.ID) {
					others++
				}
			}
			if others == 0 {
				v.errors = append(v.errors, fmt.Sprintf("Task '%s' (%s): wildcard '%s' only matches itself (self-dependency deadlock)", key, t.ID, dep))
			} else {
				v.warnings = append(v.warnings, fmt.Sprintf("Task '%s' (%s): wildcard '%s' matches itself plus %d others (self excluded at runtime, but verify intent)", key, t.ID, dep, others))
			}
		}
	}
}

// V2: Gate integrity
func (v *validator) validateGateIntegrity(tasks map[string]task.Task) {
	type gateInfo struct {
		key string
		id  string
	}

	// Find all gate tasks
	var gates []gateInfo
	for key, t := range tasks {
		if strings.HasSuffix(t.ID, ".gate") && t.Breaking {
			gates = append(gates, gateInfo{key: key, id: t.ID})
		}
	}

	for _, g := range gates {
		phase := getTaskPhase(g.id)
		if phase <= 0 {
			continue
		}

		// Gate N.gate is phase N's exit gate — must depend on N.summary
		ownSummary := fmt.Sprintf("%d.summary", phase)
		hasOwnSummary := false
		for _, t := range tasks {
			if t.ID == ownSummary {
				hasOwnSummary = true
				break
			}
		}
		if hasOwnSummary {
			found := false
			for _, dep := range tasks[g.key].Dependencies {
				if dep == ownSummary {
					found = true
					break
				}
			}
			if !found {
				v.errors = append(v.errors, fmt.Sprintf("Gate '%s' (%s): must depend on own phase summary '%s'", g.key, g.id, ownSummary))
			}
		}

		// Next phase's business tasks must depend on this gate
		gateID := g.id
		nextPhase := phase + 1
		for key, t := range tasks {
			if !isBusinessTask(t.ID) {
				continue
			}
			if getTaskPhase(t.ID) != nextPhase {
				continue
			}
			// Check if this business task depends on the gate
			dependsOnGate := false
			for _, dep := range t.Dependencies {
				if dep == gateID {
					dependsOnGate = true
					break
				}
			}
			if !dependsOnGate {
				v.errors = append(v.errors, fmt.Sprintf("Task '%s' (%s): must depend on gate '%s'", key, t.ID, gateID))
			}
		}
	}
}

// V3: Phase order sanity
func (v *validator) validatePhaseOrder(tasks map[string]task.Task) {
	// Quick mode: flat tasks without phases
	if v.quickMode {
		return
	}
	// Build a lookup for gate tasks (used to recognize transitive cross-phase deps)
	gateIDs := make(map[string]bool)
	for _, t := range tasks {
		if strings.HasSuffix(t.ID, ".gate") && t.Breaking {
			gateIDs[t.ID] = true
		}
	}

	for key, t := range tasks {
		if !isBusinessTask(t.ID) {
			continue
		}
		phase := getTaskPhase(t.ID)
		if phase <= 1 {
			continue // Phase 1 has no previous phase
		}
		hasCrossPhaseDep := false
		for _, dep := range t.Dependencies {
			// Depending on a gate satisfies cross-phase ordering because gate N.gate
			// (phase N) depends on N.summary, and the current task is in a later phase.
			if gateIDs[dep] {
				hasCrossPhaseDep = true
				break
			}
			depPhase := getTaskPhase(dep)
			if depPhase > 0 && depPhase < phase {
				hasCrossPhaseDep = true
				break
			}
			// Check wildcard deps for cross-phase ordering
			if strings.HasSuffix(dep, ".x") {
				if wp, err := strconv.Atoi(strings.TrimSuffix(dep, ".x")); err == nil && wp < phase {
					hasCrossPhaseDep = true
					break
				}
			}
		}
		if !hasCrossPhaseDep {
			v.warnings = append(v.warnings, fmt.Sprintf("Task '%s' (%s): no dependency on previous phase (may be claimed too early)", key, t.ID))
		}
	}
}

// V4: Phase summary existence
func (v *validator) validatePhaseSummaries(tasks map[string]task.Task) {
	// Quick mode: flat tasks without phases, summaries, or gates
	if v.quickMode {
		return
	}
	// Collect phases that have business tasks
	phasesWithBusiness := make(map[int]bool)
	for _, t := range tasks {
		if isBusinessTask(t.ID) {
			if p := getTaskPhase(t.ID); p > 0 {
				phasesWithBusiness[p] = true
			}
		}
	}

	// Check each such phase has a .summary task
	for phase := range phasesWithBusiness {
		summaryID := fmt.Sprintf("%d.summary", phase)
		found := false
		for _, t := range tasks {
			if t.ID == summaryID {
				found = true
				break
			}
		}
		if !found {
			v.warnings = append(v.warnings, fmt.Sprintf("Phase %d has business tasks but no '%d.summary' task", phase, phase))
		}
	}
}

func (v *validator) printResults() bool {
	if len(v.info) > 0 {
		PrintSection("INFO")
		sort.Strings(v.info)
		for _, i := range v.info {
			PrintListItem(i)
		}
	}

	if len(v.warnings) > 0 {
		PrintSection("WARNINGS")
		sort.Strings(v.warnings)
		for _, w := range v.warnings {
			PrintListItem(w)
		}
	}

	if len(v.errors) > 0 {
		PrintSection("ERRORS")
		sort.Strings(v.errors)
		for _, e := range v.errors {
			PrintListItem(e)
		}
	}

	if len(v.errors) == 0 {
		PrintResult("PASS", v.filePath)
		return true
	}
	PrintResult("FAIL", fmt.Sprintf("%s (%d errors)", v.filePath, len(v.errors)))
	return false
}

// validateLiveness checks for lifecycle anomalies in blocked tasks.
func (v *validator) validateLiveness(tasks map[string]task.Task) {
	for key, t := range tasks {
		if t.Status != "blocked" {
			continue
		}

		if len(t.Dependencies) == 0 {
			v.warnings = append(v.warnings,
				fmt.Sprintf("Task '%s' (%s): blocked with no dependencies (orphaned)", key, t.ID))
			continue
		}

		allDepsCompleted := true
		hasActiveDep := false
		for _, dep := range t.Dependencies {
			if strings.HasSuffix(dep, ".x") {
				prefix := strings.TrimSuffix(dep, ".x")
				prefixWithDot := prefix + "."
				wildcardHasMatch := false
				for _, other := range tasks {
					if other.ID == t.ID {
						continue
					}
					if strings.HasPrefix(other.ID, prefixWithDot) && isBusinessTask(other.ID) {
						wildcardHasMatch = true
						if other.Status != "completed" && other.Status != "skipped" {
							allDepsCompleted = false
							if other.Status == "pending" || other.Status == "in_progress" {
								hasActiveDep = true
							}
						}
					}
				}
				// Wildcard matches no tasks — vacuously satisfied
				_ = wildcardHasMatch
				continue
			}
			depTask, found := tasks[dep]
			if !found {
				v.errors = append(v.errors,
					fmt.Sprintf("Task '%s' (%s): blocked on missing dependency '%s'", key, t.ID, dep))
				allDepsCompleted = false
				continue
			}
			if depTask.Status == "completed" || depTask.Status == "skipped" {
				continue
			}
			allDepsCompleted = false
			if depTask.Status == "pending" || depTask.Status == "in_progress" {
				hasActiveDep = true
			}
		}

		if allDepsCompleted {
			v.warnings = append(v.warnings,
				fmt.Sprintf("Task '%s' (%s): blocked but all dependencies resolved (stale, should be pending)", key, t.ID))
		} else if !hasActiveDep {
			v.warnings = append(v.warnings,
				fmt.Sprintf("Task '%s' (%s): blocked with no path to resolution (all deps blocked or missing)", key, t.ID))
		}
	}
}
