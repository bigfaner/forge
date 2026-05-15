// Package task provides task-related types and operations.
package task

import (
	"encoding/json"
	"time"
)

// Task type constants define the 14 valid execution types.
const (
	TypeImplementation               = "implementation"
	TypeDocumentation                = "documentation"
	TypeDocEvaluation                = "doc-evaluation"
	TypeDocGenerationSummary         = "doc-generation.summary"
	TypeDocGenerationConsolidate     = "doc-generation.consolidate"
	TypeDocGenerationDrift           = "doc-generation.drift"
	TypeTestPipelineGenCases         = "test-pipeline.gen-cases"
	TypeTestPipelineEvalCases        = "test-pipeline.eval-cases"
	TypeTestPipelineGenScripts       = "test-pipeline.gen-scripts"
	TypeTestPipelineRun              = "test-pipeline.run"
	TypeTestPipelineGraduate         = "test-pipeline.graduate"
	TypeTestPipelineVerifyRegression = "test-pipeline.verify-regression"
	TypeFix                          = "fix"
	TypeGate                         = "gate"
)

// TaskTypeInfo describes a single task type for display and discovery.
type TaskTypeInfo struct { //nolint:revive // intentional naming for API clarity
	Name        string
	Description string
}

// TaskTypeRegistry is the centralized source of truth for all supported task types.
// Each entry's Description follows verb+object format and is <= 60 chars.
var TaskTypeRegistry = []TaskTypeInfo{
	{Name: TypeImplementation, Description: "implement feature task"},
	{Name: TypeDocumentation, Description: "write or update documentation"},
	{Name: TypeDocEvaluation, Description: "evaluate documentation quality"},
	{Name: TypeFix, Description: "fix a bug or issue"},
	{Name: TypeGate, Description: "validate quality gate before proceeding"},
	{Name: TypeDocGenerationSummary, Description: "generate documentation summary"},
	{Name: TypeDocGenerationConsolidate, Description: "consolidate documentation files"},
	{Name: TypeDocGenerationDrift, Description: "detect and fix spec drift against codebase"},
	{Name: TypeTestPipelineGenCases, Description: "generate test cases from acceptance criteria"},
	{Name: TypeTestPipelineEvalCases, Description: "evaluate generated test cases for quality"},
	{Name: TypeTestPipelineGenScripts, Description: "generate executable test scripts"},
	{Name: TypeTestPipelineRun, Description: "run test scripts and collect results"},
	{Name: TypeTestPipelineGraduate, Description: "graduate tests to regression suite"},
	{Name: TypeTestPipelineVerifyRegression, Description: "verify regression suite after graduation"},
}

// ValidTypes is the complete set of valid task type values.
var ValidTypes = map[string]bool{
	TypeImplementation:               true,
	TypeDocumentation:                true,
	TypeDocEvaluation:                true,
	TypeDocGenerationSummary:         true,
	TypeDocGenerationConsolidate:     true,
	TypeDocGenerationDrift:           true,
	TypeTestPipelineGenCases:         true,
	TypeTestPipelineEvalCases:        true,
	TypeTestPipelineGenScripts:       true,
	TypeTestPipelineRun:              true,
	TypeTestPipelineGraduate:         true,
	TypeTestPipelineVerifyRegression: true,
	TypeFix:                          true,
	TypeGate:                         true,
}

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
	// NoTest indicates this task does not require tests (e.g., documentation-only tasks).
	// When true, quality gate and test evidence checks are skipped, and coverage is auto-set to -1.0.
	NoTest bool `json:"noTest,omitempty"`
	// Type is the task execution type (e.g. "implementation", "fix", "gate").
	// Required for all tasks after migration; validated by task validate.
	// omitempty allows existing index.json files to load without error.
	Type string `json:"type,omitempty"`
	// Profile indicates the test profile associated with this task (e.g. "web-playwright", "go-test").
	// Set by task index for per-profile test tasks; empty for business tasks and shared test tasks.
	Profile string `json:"profile,omitempty"`
	// BlockedReason records why a task entered blocked state.
	// Written by run-tasks when task prompt exits non-zero.
	BlockedReason string `json:"blockedReason,omitempty"`
}

// TaskIndex represents the index.json structure for a feature.
// The tasks map is unexported; use ByID, SetTask, TasksMap, etc.
type TaskIndex struct { //nolint:revive // intentional naming for API clarity
	Feature      string
	PRD          string // e.g. "prd/prd-spec.md"
	Design       string // e.g. "design/tech-design.md"
	Proposal     string // e.g. "docs/proposals/<slug>/proposal.md" (quick mode alternative to PRD+Design)
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

// MarshalJSON serializes the TaskIndex to JSON.
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

// UnmarshalJSON deserializes JSON data into the TaskIndex.
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
type TaskState struct { //nolint:revive // intentional naming for API clarity
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
	// NoTest mirrors Task.NoTest for the claimed task.
	NoTest bool `json:"noTest,omitempty"`
	// Type mirrors Task.Type for the claimed task (same pattern as MainSession, NoTest).
	Type string `json:"type,omitempty"`
	// Profile mirrors Task.Profile for the claimed task.
	Profile string `json:"profile,omitempty"`
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
			"rejected",
		},
		PriorityEnum: []string{"P0", "P1", "P2"},
		tasks:        make(map[string]Task),
		Created:      time.Now().Format("2006-01-02"),
	}
}
