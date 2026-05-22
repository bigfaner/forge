package task

// PreserveRuntimeFields copies runtime-only fields from existing task to new task.
// Called during BuildIndex re-index to preserve state that isn't in .md frontmatter.
// Uses explicit field assignment (not reflection) for compile-time safety.
// Adding a new preserved field = adding one line of code.
func PreserveRuntimeFields(existing, newTask *Task) {
	if existing == nil {
		return
	}
	newTask.Status = existing.Status
	newTask.SourceTaskID = existing.SourceTaskID
	newTask.BlockedReason = existing.BlockedReason
}
