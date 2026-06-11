package types

// Priority represents the urgency level of a task.
type Priority string

// Task priority constants.
const (
	PriorityP0 Priority = "P0"
	PriorityP1 Priority = "P1"
	PriorityP2 Priority = "P2"
)

// AllPriorities returns all defined Priority constants.
func AllPriorities() []Priority {
	return []Priority{
		PriorityP0,
		PriorityP1,
		PriorityP2,
	}
}
