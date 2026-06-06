package task

import (
	"fmt"

	"forge-cli/pkg/task"
	"forge-cli/pkg/types"
)

// validateLiveness checks for lifecycle anomalies in blocked tasks.
func (v *validator) validateLiveness(index *task.TaskIndex) {
	for key, t := range index.TasksMap() {
		if t.Status != types.StatusBlocked {
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
			matches, isWildcard := task.ResolveWildcardDep(index, dep)
			if isWildcard {
				for _, matchID := range matches {
					if matchID == t.ID {
						continue
					}
					other, _ := index.ByID(matchID)
					if !task.IsDepSatisfied(string(other.Status)) {
						allDepsCompleted = false
						if other.Status == types.StatusPending || other.Status == types.StatusInProgress {
							hasActiveDep = true
						}
					}
				}
				continue
			}
			depTask, found := index.ByID(dep)
			if !found {
				v.errors = append(v.errors,
					fmt.Sprintf("Task '%s' (%s): blocked on missing dependency '%s'", key, t.ID, dep))
				allDepsCompleted = false
				continue
			}
			if task.IsDepSatisfied(string(depTask.Status)) {
				continue
			}
			allDepsCompleted = false
			if depTask.Status == types.StatusPending || depTask.Status == types.StatusInProgress {
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
