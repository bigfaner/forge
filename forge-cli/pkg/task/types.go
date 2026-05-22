// Package task provides task-related types and operations.
package task

import (
	"encoding/json"
	"sort"
	"strings"
	"time"
)

// Task ID suffix constants used for special task roles.
const (
	// IDSuffixGate is the suffix for quality gate tasks (e.g. "1.gate").
	IDSuffixGate = ".gate"
	// IDSuffixSummary is the suffix for phase summary tasks (e.g. "1.summary").
	IDSuffixSummary = ".summary"
	// IDSuffixWildcard is the suffix for wildcard dependencies (e.g. "1.x").
	IDSuffixWildcard = ".x"
	// IDPrefixTestPipeline is the prefix for auto-generated test pipeline tasks.
	IDPrefixTestPipeline = "T-"
)

// IsBusinessTask returns true for task IDs that are regular business tasks:
// not gate tasks, not summary tasks, and not auto-generated test pipeline tasks.
func IsBusinessTask(id string) bool {
	if strings.HasPrefix(id, IDPrefixTestPipeline) {
		return false
	}
	if strings.HasSuffix(id, IDSuffixGate) {
		return false
	}
	if strings.HasSuffix(id, IDSuffixSummary) {
		return false
	}
	return true
}

// Task type constants define the valid execution types.
// Naming convention: prefix-based categories (coding.*, doc*, test.*, validation.*).
const (
	TypeCodingFeature        = "coding.feature"
	TypeCodingEnhancement    = "coding.enhancement"
	TypeCodingCleanup        = "coding.cleanup"
	TypeCodingRefactor       = "coding.refactor"
	TypeCodingFix            = "coding.fix"
	TypeDoc                  = "doc"
	TypeDocEval              = "doc.eval"
	TypeDocSummary           = "doc.summary"
	TypeDocConsolidate       = "doc.consolidate"
	TypeDocDrift             = "doc.drift"
	TypeTestGenCases         = "test.gen-cases"
	TypeTestEvalCases        = "test.eval-cases"
	TypeTestGenScripts       = "test.gen-scripts"
	TypeTestRun              = "test.run"
	TypeTestGenAndRun        = "test.gen-and-run"
	TypeTestGraduate         = "test.graduate"
	TypeTestVerifyRegression = "test.verify-regression"
	TypeValidationCode       = "validation.code"
	TypeValidationUx         = "validation.ux"
	TypeGate                 = "gate"
	TypeCleanCode            = "code-quality.simplify"
)

// TaskTypeInfo describes a single task type for display and discovery.
type TaskTypeInfo struct { //nolint:revive // intentional naming for API clarity
	Name        string
	Description string
}

// TaskTypeRegistry is the centralized source of truth for all supported task types.
// Each entry's Description follows verb+object format and is <= 60 chars.
var TaskTypeRegistry = []TaskTypeInfo{
	{Name: TypeCodingFeature, Description: "implement new runtime behavior"},
	{Name: TypeCodingEnhancement, Description: "enhance existing behavior"},
	{Name: TypeCodingCleanup, Description: "remove dead code or fix technical debt"},
	{Name: TypeCodingRefactor, Description: "restructure code without behavior change"},
	{Name: TypeCodingFix, Description: "fix a bug or issue"},
	{Name: TypeDoc, Description: "write or update documentation"},
	{Name: TypeDocEval, Description: "evaluate documentation quality"},
	{Name: TypeDocSummary, Description: "generate documentation summary"},
	{Name: TypeDocConsolidate, Description: "consolidate documentation files"},
	{Name: TypeDocDrift, Description: "detect and fix spec drift against codebase"},
	{Name: TypeTestGenCases, Description: "generate test cases from acceptance criteria"},
	{Name: TypeTestEvalCases, Description: "evaluate generated test cases for quality"},
	{Name: TypeTestGenScripts, Description: "generate executable test scripts"},
	{Name: TypeTestRun, Description: "run test scripts and collect results"},
	{Name: TypeTestGenAndRun, Description: "generate and run test scripts in one session"},
	{Name: TypeTestGraduate, Description: "graduate tests to regression suite"},
	{Name: TypeTestVerifyRegression, Description: "verify regression suite after graduation"},
	{Name: TypeValidationCode, Description: "validate code quality and correctness"},
	{Name: TypeValidationUx, Description: "validate user experience quality"},
	{Name: TypeGate, Description: "validate quality gate before proceeding"},
	{Name: TypeCleanCode, Description: "simplify and clean up code quality"},
}

// ValidTypes is the complete set of valid task type values.
var ValidTypes = map[string]bool{
	TypeCodingFeature:        true,
	TypeCodingEnhancement:    true,
	TypeCodingCleanup:        true,
	TypeCodingRefactor:       true,
	TypeCodingFix:            true,
	TypeDoc:                  true,
	TypeDocEval:              true,
	TypeDocSummary:           true,
	TypeDocConsolidate:       true,
	TypeDocDrift:             true,
	TypeTestGenCases:         true,
	TypeTestEvalCases:        true,
	TypeTestGenScripts:       true,
	TypeTestRun:              true,
	TypeTestGenAndRun:        true,
	TypeTestGraduate:         true,
	TypeTestVerifyRegression: true,
	TypeValidationCode:       true,
	TypeValidationUx:         true,
	TypeGate:                 true,
	TypeCleanCode:            true,
}

// SystemTypes is the set of auto-generated system task types (13 total).
// These types are created by the forge pipeline, not by users.
// Dual-identity types (doc.consolidate, doc.drift) are excluded because
// they can also serve as business tasks.
var SystemTypes = map[string]bool{
	TypeGate:                 true,
	TypeTestGenCases:         true,
	TypeTestEvalCases:        true,
	TypeTestGenScripts:       true,
	TypeTestRun:              true,
	TypeTestGenAndRun:        true,
	TypeTestGraduate:         true,
	TypeTestVerifyRegression: true,
	TypeValidationCode:       true,
	TypeValidationUx:         true,
	TypeDocEval:              true,
	TypeDocSummary:           true,
	TypeCleanCode:            true,
}

// IsSystemType returns true if the given type is an auto-generated system type.
// Business types and dual-identity types (doc.consolidate, doc.drift) return false.
func IsSystemType(typ string) bool {
	return SystemTypes[typ]
}

// FormatSystemTypes returns a comma-separated list of all system type names for error messages.
func FormatSystemTypes() string {
	names := make([]string, 0, len(SystemTypes))
	for name := range SystemTypes {
		names = append(names, name)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
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
	// Type is the task execution type (e.g. "coding.feature", "coding.fix", "gate").
	// Required for all tasks after migration; validated by task validate.
	// omitempty allows existing index.json files to load without error.
	Type string `json:"type,omitempty"`
	// BlockedReason records why a task entered blocked state.
	// Written by run-tasks when task prompt exits non-zero.
	BlockedReason string `json:"blockedReason,omitempty"`
	// ManualBlock marks a task as manually blocked by the operator (forge task transition).
	// Auto-unblock skips these tasks, preventing the claim→unblock→claim infinite loop.
	ManualBlock bool `json:"manualBlock,omitempty"`
	// Coverage is an optional per-task coverage override from frontmatter.
	// nil means use the global coverage config default for the task type.
	// Non-nil (including 0) means the task has an explicit coverage target.
	Coverage *int `json:"coverage,omitempty"`
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
	// Type mirrors Task.Type for the claimed task (same pattern as MainSession).
	Type string `json:"type,omitempty"`
}

// RecordData represents the JSON input for record generation.
type RecordData struct {
	TaskID               string                `json:"taskId"`
	Status               string                `json:"status"`
	Summary              string                `json:"summary"`
	FilesCreated         []string              `json:"filesCreated"`
	FilesModified        []string              `json:"filesModified"`
	KeyDecisions         []string              `json:"keyDecisions"`
	TestsPassed          int                   `json:"testsPassed"`
	TestsFailed          int                   `json:"testsFailed"`
	Coverage             float64               `json:"coverage"`
	AcceptanceCriteria   []AcceptanceCriterion `json:"acceptanceCriteria"`
	Notes                string                `json:"notes"`
	TypeReclassification *TypeReclassification `json:"typeReclassification,omitempty"`
}

// TypeReclassification documents when an executor changes a task's type during execution.
type TypeReclassification struct {
	OriginalType string `json:"originalType"` // e.g. "fix"
	ActualType   string `json:"actualType"`   // e.g. "cleanup"
	Reason       string `json:"reason"`       // why the type was changed
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
