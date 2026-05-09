// Package task provides task-related types and operations.
package task

import (
	"encoding/json"
	"time"
)

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
	Breaking      bool     `json:"breaking,omitempty"`
	// Scope indicates the task's domain: "frontend", "backend", or "all".
	// Default is "all" (enforced by consumers, not a Go zero value).
	// Omitempty allows existing index.json files without scope to remain valid.
	Scope string `json:"scope,omitempty"`
	// SourceTaskID records which task spawned this task (e.g. fix-task -> source task).
	// Empty for original tasks. Used by record auto-restore to unblock source when all deps complete.
	SourceTaskID string `json:"sourceTaskID,omitempty"`
	// MainSession indicates this task must run in the main session (not dispatched to task-executor).
	// Used by tasks that need to spawn subagents (e.g., eval-test-cases spawns doc-scorer/doc-reviser).
	MainSession bool `json:"mainSession,omitempty"`
}

// TaskIndex represents the index.json structure for a feature.
// The tasks map is unexported; use ByID, SetTask, TasksMap, etc.
type TaskIndex struct {
	Feature      string
	PRD          string   // e.g. "prd/prd-spec.md"
	Design       string   // e.g. "design/tech-design.md"
	Proposal     string   // e.g. "docs/proposals/<slug>/proposal.md" (quick mode alternative to PRD+Design)
	Created      string
	Status       string
	tasks        map[string]Task
	StatusEnum   []string
	PriorityEnum []string
	TestCommand  string
	E2ERound     int // current fix-e2e round (0 = no failures yet)
}

// taskIndexJSON mirrors TaskIndex for JSON serialization.
type taskIndexJSON struct {
	Feature      string          `json:"feature"`
	PRD          string          `json:"prd,omitempty"`
	Design       string          `json:"design,omitempty"`
	Proposal     string          `json:"proposal,omitempty"`
	Created      string          `json:"created,omitempty"`
	Status       string          `json:"status,omitempty"`
	Tasks        map[string]Task `json:"tasks"`
	StatusEnum   []string        `json:"statusEnum,omitempty"`
	PriorityEnum []string        `json:"priorityEnum,omitempty"`
	TestCommand  string          `json:"testCommand,omitempty"`
	E2ERound     int             `json:"e2eRound,omitempty"`
}

func (ti TaskIndex) MarshalJSON() ([]byte, error) {
	return json.Marshal(taskIndexJSON{
		Feature:      ti.Feature,
		PRD:          ti.PRD,
		Design:       ti.Design,
		Proposal:     ti.Proposal,
		Created:      ti.Created,
		Status:       ti.Status,
		Tasks:        ti.tasks,
		StatusEnum:   ti.StatusEnum,
		PriorityEnum: ti.PriorityEnum,
		TestCommand:  ti.TestCommand,
		E2ERound:     ti.E2ERound,
	})
}

func (ti *TaskIndex) UnmarshalJSON(data []byte) error {
	var j taskIndexJSON
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	ti.Feature = j.Feature
	ti.PRD = j.PRD
	ti.Design = j.Design
	ti.Proposal = j.Proposal
	ti.Created = j.Created
	ti.Status = j.Status
	ti.tasks = j.Tasks
	ti.StatusEnum = j.StatusEnum
	ti.PriorityEnum = j.PriorityEnum
	ti.TestCommand = j.TestCommand
	ti.E2ERound = j.E2ERound
	return nil
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
	Breaking      bool     `json:"breaking,omitempty"`
	// Scope mirrors Task.Scope for the claimed task.
	Scope string `json:"scope,omitempty"`
	// MainSession mirrors Task.MainSession for the claimed task.
	MainSession bool `json:"mainSession,omitempty"`
}

// RecordData represents the JSON input for record generation.
type RecordData struct {
	TaskID             string                `json:"taskId"`
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
		tasks:        make(map[string]Task),
		Created:      time.Now().Format("2006-01-02"),
	}
}
