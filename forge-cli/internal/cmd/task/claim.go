package task

import (
	"fmt"
	"forge-cli/internal/cmd/base"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/forgelog"
	indexPkg "forge-cli/pkg/index"
	"forge-cli/pkg/project"
	"forge-cli/pkg/task"
	"forge-cli/pkg/types"

	"github.com/spf13/cobra"
)

// unreachableDepth is assigned to tasks in dependency cycles,
// indicating they are not reachable from any root in BFS traversal.
const unreachableDepth = 99999

var claimCmd = &cobra.Command{
	Use:   "claim",
	Short: "Claim the next available task",
	Long: `Claim the next available task from the current feature's task list.

The task is selected based on:
1. All dependencies must be met
2. Topological order (shallower depth first)
3. Priority as tiebreaker within same depth (P0 > P1 > P2)
4. Task ID (semantic version ordering)`,
	Args: cobra.NoArgs,
	RunE: runClaim,
}

func runClaim(_ *cobra.Command, _ []string) error {
	result, err := executeClaim()
	if err != nil {
		base.Exit(err)
	}

	// Print result
	if result.Action == "CONTINUE" {
		printContinueTask(result.State, result.Task, result.ProjectRoot, result.FeatureSlug)
	} else {
		printNewTask(result.Key, result.Task, result.ProjectRoot, result.FeatureSlug)
	}
	return nil
}

// ClaimResult represents the result of a claim operation
type ClaimResult struct {
	Action      string // "CLAIMED" or "CONTINUE"
	Key         string
	Task        *task.Task
	State       *task.TaskState
	StartedTime string // for CONTINUE action
	ProjectRoot string
	FeatureSlug string
}

// executeClaim contains the core logic for the claim command, testable
func executeClaim() (*ClaimResult, error) {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return nil, base.ErrProjectNotFound()
	}

	featureSlug, err := feature.RequireFeature(projectRoot)
	if err != nil {
		return nil, base.ErrFeatureNotSet()
	}

	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))
	statePath := feature.GetTaskStatePath(projectRoot, featureSlug)

	// Read-only check: no lock needed.
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		return nil, base.NewAIError(base.ErrNotFound, "Failed to load task index", err.Error(), "Run `forge task index --feature "+featureSlug+"` to generate it", "forge task index --feature "+featureSlug)
	}

	// Check for existing task state
	continueTask, hasIssues, issues := task.CheckExistingTaskState(projectRoot, index, statePath)

	if hasIssues {
		return nil, base.ErrDataIntegrity(issues)
	}

	if continueTask {
		state, _ := task.LoadState(statePath)
		t, _ := index.ByID(state.Key)
		return &ClaimResult{
			Action:      "CONTINUE",
			Key:         state.Key,
			Task:        &t,
			State:       state,
			StartedTime: state.StartedTime,
			ProjectRoot: projectRoot,
			FeatureSlug: featureSlug,
		}, nil
	}

	// Claim new task — wrapped in WithLock for atomic read-modify-write.
	var key string
	var t *task.Task
	if err := indexPkg.WithLock(indexPath, func() error {
		// Re-load index inside lock (it may have changed since our read-only check).
		idx, err := task.LoadIndex(indexPath)
		if err != nil {
			return base.NewAIError(base.ErrNotFound, "Failed to load task index", err.Error(), "Run `forge task index --feature "+featureSlug+"` to generate it", "forge task index --feature "+featureSlug)
		}

		k, claimedTask, err := claimNextTask(idx)
		if err != nil {
			return err
		}

		if err := indexPkg.SaveIndexAtomic(indexPath, idx); err != nil {
			return base.NewAIError(base.ErrConflict, "Failed to save index", err.Error(), "Check index.json is writable", "cat "+indexPath)
		}

		key, t = k, claimedTask
		return nil
	}); err != nil {
		return nil, err
	}

	// Bootstrap .forge/state.json as workspace marker so subagents in subdirectories
	// (e.g., backend/) can find the project root via FindProjectRoot().
	// Placed after all validation to avoid creating artifacts in error paths.
	if err := feature.EnsureForgeState(projectRoot, featureSlug); err != nil {
		forgelog.Warn("WARNING: failed to write .forge/state.json: %v\n", err)
	}

	// Save state
	state := &task.TaskState{
		TaskID:        t.ID,
		Key:           key,
		Title:         t.Title,
		Priority:      string(t.Priority),
		EstimatedTime: t.EstimatedTime,
		Dependencies:  t.Dependencies,
		File:          t.File,
		Record:        t.Record,
		StartedTime:   time.Now().Format("2006-01-02 15:04"),
		Breaking:      t.Breaking,
		SurfaceKey:    t.SurfaceKey,
		SurfaceType:   t.SurfaceType,
		MainSession:   t.MainSession,
		Type:          t.Type,
	}
	if err := task.SaveState(statePath, state); err != nil {
		return nil, err
	}

	return &ClaimResult{
		Action:      "CLAIMED",
		Key:         key,
		Task:        t,
		State:       state,
		ProjectRoot: projectRoot,
		FeatureSlug: featureSlug,
	}, nil
}

func claimNextTask(index *task.TaskIndex) (string, *task.Task, error) {
	type taskWithKey struct {
		key string
		t   task.Task
	}
	var eligibleTasks []taskWithKey

	// Lazy unblock scan: check blocked tasks and auto-transition eligible ones to pending.
	// Runs before the hasPending check so newly-unblocked tasks are visible.
	// Suspended tasks are naturally excluded (they have status "suspended", not "blocked").
	for key, t := range index.TasksMap() {
		if t.Status != types.StatusBlocked {
			continue
		}
		if met, _ := checkDependenciesMet(index, t.ID, t); met {
			t.Status = types.StatusPending
			index.SetTask(key, t)
			fmt.Printf("Auto-unblocked task %s\n", t.ID)
		}
	}

	hasPending := false
	for _, t := range index.TasksMap() {
		if t.Status == types.StatusPending {
			hasPending = true
			break
		}
	}
	if !hasPending {
		return "", nil, base.ErrNoPendingTasks()
	}

	for key, t := range index.TasksMap() {
		if t.Status == types.StatusPending {
			if met, _ := checkDependenciesMet(index, t.ID, t); met {
				eligibleTasks = append(eligibleTasks, taskWithKey{key: key, t: t})
			}
		}
	}

	if len(eligibleTasks) == 0 {
		return "", nil, fmt.Errorf("no task available with met dependencies")
	}

	priorityOrder := map[types.Priority]int{types.PriorityP0: 0, types.PriorityP1: 1, types.PriorityP2: 2}

	// Compute topological depths for all tasks in the index.
	depths := computeTopoDepths(index)

	sort.Slice(eligibleTasks, func(i, j int) bool {
		di, dj := depths[eligibleTasks[i].t.ID], depths[eligibleTasks[j].t.ID]
		if di != dj {
			return di < dj
		}
		pi, pj := priorityOrder[eligibleTasks[i].t.Priority], priorityOrder[eligibleTasks[j].t.Priority]
		if pi != pj {
			return pi < pj
		}
		return task.CompareVersionIDs(eligibleTasks[i].t.ID, eligibleTasks[j].t.ID)
	})

	twk := eligibleTasks[0]
	t, _ := index.ByID(twk.key)
	t.Status = types.StatusInProgress
	index.SetTask(twk.key, t)
	return twk.key, &t, nil
}

func checkDependenciesMet(index *task.TaskIndex, selfID string, t task.Task) (bool, []string) {
	rawUnmet := task.GetUnmetDeps(index, selfID, t.Dependencies)

	// GetUnmetDeps reports missing exact deps as unmet, but for claim purposes
	// unknown deps are vacuously satisfied (they don't block claiming).
	var unmet []string
	for _, id := range rawUnmet {
		if _, found := index.ByID(id); !found {
			continue
		}
		unmet = append(unmet, id)
	}

	// Check for pending fix tasks whose SourceTaskID matches any dependency.
	// If task depends on X and a fix task with sourceTaskID "X" is still
	// pending/in_progress, the dependency is not truly met.
	for _, dep := range t.Dependencies {
		for _, other := range index.TasksMap() {
			if other.ID != selfID && other.Type == task.TypeCodingFix && other.SourceTaskID == dep &&
				(other.Status == types.StatusPending || other.Status == types.StatusInProgress) {
				unmet = append(unmet, other.ID)
			}
		}
	}

	// Check for active fix tasks targeting this task itself (SourceTaskID == selfID).
	// If a fix task with sourceTaskID == selfID is still pending/in_progress,
	// this task should not be claimed (--block-source scenario).
	for _, other := range index.TasksMap() {
		if other.Type == task.TypeCodingFix && other.SourceTaskID == selfID &&
			(other.Status == types.StatusPending || other.Status == types.StatusInProgress) {
			unmet = append(unmet, other.ID)
		}
	}

	return len(unmet) == 0, unmet
}

func printTaskDetails(key string, t *task.Task, projectRoot, featureSlug string) {
	_ = key // key is still used internally for routing, but no longer emitted
	base.PrintField("TASK_ID", t.ID)
	base.PrintFieldIfNotEmpty("TYPE", t.Type)
	base.PrintFieldIfNotEmpty("TASK_CATEGORY", task.CategoryForType(t.Type))
	base.PrintFieldIfNotEmpty("FEATURE", featureSlug)
	base.PrintField("FILE", filepath.Join(projectRoot, feature.GetTaskFile(featureSlug, t.File)))
	base.PrintFieldIfNotEmpty("SURFACE_KEY", t.SurfaceKey)
	base.PrintFieldIfNotEmpty("SURFACE_TYPE", t.SurfaceType)
	if t.MainSession {
		base.PrintField("MAIN_SESSION", "true")
	}
}

func printContinueTask(state *task.TaskState, t *task.Task, projectRoot, featureSlug string) {
	base.PrintBlockStart()
	base.PrintField("ACTION", "CONTINUE")
	printTaskDetails(state.Key, t, projectRoot, featureSlug)
	base.PrintField("STARTED_AT", state.StartedTime)
	base.PrintBlockEnd()
}

func printNewTask(key string, t *task.Task, projectRoot, featureSlug string) {
	base.PrintBlockStart()
	base.PrintField("ACTION", "CLAIMED")
	printTaskDetails(key, t, projectRoot, featureSlug)
	base.PrintBlockEnd()
}

// computeTopoDepths computes the topological depth of each task in the index
// using BFS (Kahn's algorithm style). Root tasks (no dependencies) have depth 0.
// The returned map keys are task IDs (not map keys).
// Tasks in dependency cycles receive a large depth value (they won't be eligible anyway).
func computeTopoDepths(index *task.TaskIndex) map[string]int {
	tasks := index.TasksMap()
	depths := make(map[string]int, len(tasks))

	// Build adjacency list: dep task ID -> dependent task IDs.
	// Keys are task IDs, not map keys.
	adj := make(map[string][]string)
	inDeg := make(map[string]int)
	for _, t := range tasks {
		inDeg[t.ID] = 0
	}

	for _, t := range tasks {
		for _, dep := range t.Dependencies {
			if strings.HasSuffix(dep, ".x") {
				matches, _ := task.ResolveWildcardDep(index, dep)
				for _, m := range matches {
					if m != t.ID { // skip self-edges
						adj[m] = append(adj[m], t.ID)
						inDeg[t.ID]++
					}
				}
			} else {
				if _, found := index.ByID(dep); found && dep != t.ID {
					adj[dep] = append(adj[dep], t.ID)
					inDeg[t.ID]++
				}
			}
		}
	}

	// BFS from roots (in-degree 0).
	var queue []string
	for _, t := range tasks {
		if inDeg[t.ID] == 0 {
			depths[t.ID] = 0
			queue = append(queue, t.ID)
		}
	}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		for _, next := range adj[curr] {
			d := depths[curr] + 1
			if d > depths[next] {
				depths[next] = d
			}
			inDeg[next]--
			if inDeg[next] == 0 {
				queue = append(queue, next)
			}
		}
	}

	// Tasks in cycles get a large depth (unreachable from BFS).
	for _, t := range tasks {
		if _, ok := depths[t.ID]; !ok {
			depths[t.ID] = unreachableDepth
		}
	}

	return depths
}
