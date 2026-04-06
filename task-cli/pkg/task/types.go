// Package task provides task-related types and operations.
package task

import "time"

// Task represents a single task in the feature index.
type Task struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Priority      string   `json:"priority"`
	EstimatedTime string   `json:"estimatedTime,omitempty"`
	Dependencies  []string `json:"dependencies,omitempty"`
	Status        string   `json:"status"`
	File          string   `json:"file"`
	Record        string   `json:"record"`
}

// TaskIndex represents the index.json structure for a feature.
type TaskIndex struct {
	Feature      string          `json:"feature"`
	PRD          string          `json:"prd,omitempty"`
	Design       string          `json:"design,omitempty"`
	Created      string          `json:"created,omitempty"`
	Status       string          `json:"status,omitempty"`
	Tasks        map[string]Task `json:"tasks"`
	StatusEnum   []string        `json:"statusEnum,omitempty"`
	PriorityEnum []string        `json:"priorityEnum,omitempty"`
}

// TaskState represents the current in-progress task state.
type TaskState struct {
	TaskID        string   `json:"task_id"`
	Key           string   `json:"key"`
	Title         string   `json:"title"`
	Priority      string   `json:"priority"`
	EstimatedTime string   `json:"estimatedTime,omitempty"`
	Dependencies  []string `json:"dependencies,omitempty"`
	File          string   `json:"file"`
	Record        string   `json:"record"`
	StartedTime   string   `json:"startedTime"`
}

// RecordData represents the JSON input for record generation.
type RecordData struct {
	Status             string                `json:"status"`
	Summary            string                `json:"summary"`
	FilesCreated       []string              `json:"filesCreated"`
	FilesModified      []string              `json:"filesModified"`
	KeyDecisions       []string              `json:"keyDecisions"`
	TestsPassed        int                   `json:"testsPassed"`
	TestsFailed        int                   `json:"testsFailed"`
	Coverage           float64               `json:"coverage"`
	AcceptanceCriteria []AcceptanceCriterion `json:"acceptanceCriteria"`
	Notes              string                `json:"notes"`
}

// AcceptanceCriterion represents a single acceptance criterion.
type AcceptanceCriterion struct {
	Criterion string `json:"criterion"`
	Met       bool   `json:"met"`
}

// NewTaskIndex creates a new TaskIndex with default enum values.
func NewTaskIndex(feature string) *TaskIndex {
	return &TaskIndex{
		Feature: feature,
		StatusEnum: []string{
			"pending",
			"in_progress",
			"completed",
			"blocked",
			"skipped",
		},
		PriorityEnum: []string{"P0", "P1", "P2"},
		Tasks:        make(map[string]Task),
		Created:      time.Now().Format("2006-01-02"),
	}
}
