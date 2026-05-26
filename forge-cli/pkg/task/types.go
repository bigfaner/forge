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
	TypeDocReview            = "doc.review"
	TypeDocSummary           = "doc.summary"
	TypeDocConsolidate       = "doc.consolidate"
	TypeDocDrift             = "doc.drift"
	TypeTestGenContracts     = "test.gen-contracts"
	TypeTestGenJourneys      = "test.gen-journeys"
	TypeTestGenScripts       = "test.gen-scripts"
	TypeTestRun              = "test.run"
	TypeTestVerifyRegression = "test.verify-regression"
	TypeEvalJourney          = "eval.journey"
	TypeEvalContract         = "eval.contract"
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
	{Name: TypeDocReview, Description: "review documentation against acceptance criteria"},
	{Name: TypeDocSummary, Description: "generate documentation summary"},
	{Name: TypeDocConsolidate, Description: "consolidate documentation files"},
	{Name: TypeDocDrift, Description: "detect and fix spec drift against codebase"},
	{Name: TypeTestGenContracts, Description: "generate test contracts from journeys"},
	{Name: TypeTestGenJourneys, Description: "generate test journeys from specs"},
	{Name: TypeTestGenScripts, Description: "generate executable test scripts"},
	{Name: TypeTestRun, Description: "run test scripts and collect results"},
	{Name: TypeTestVerifyRegression, Description: "verify regression suite after graduation"},
	{Name: TypeEvalJourney, Description: "evaluate Journey quality with rubric scoring"},
	{Name: TypeEvalContract, Description: "evaluate Contract quality with rubric scoring"},
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
	TypeDocReview:            true,
	TypeDocSummary:           true,
	TypeDocConsolidate:       true,
	TypeDocDrift:             true,
	TypeTestGenContracts:     true,
	TypeTestGenJourneys:      true,
	TypeTestGenScripts:       true,
	TypeTestRun:              true,
	TypeTestVerifyRegression: true,
	TypeEvalJourney:          true,
	TypeEvalContract:         true,
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
	TypeTestGenContracts:     true,
	TypeTestGenJourneys:      true,
	TypeTestGenScripts:       true,
	TypeTestRun:              true,
	TypeTestVerifyRegression: true,
	TypeEvalJourney:          true,
	TypeEvalContract:         true,
	TypeValidationCode:       true,
	TypeValidationUx:         true,
	TypeDocReview:            true,
	TypeDocSummary:           true,
	TypeCleanCode:            true,
}

// IsSystemType returns true if the given type is an auto-generated system type.
// Business types and dual-identity types (doc.consolidate, doc.drift) return false.
// Surface-specific variants (e.g. "test.gen-scripts.cli") are also recognized
// by checking if the type with the last segment stripped matches a system type.
func IsSystemType(typ string) bool {
	if SystemTypes[typ] {
		return true
	}
	// Check surface-specific variants: strip last ".xxx" segment and check base type
	if idx := strings.LastIndex(typ, "."); idx >= 0 {
		base := typ[:idx]
		return SystemTypes[base]
	}
	return false
}

// TestTypeTitle returns the human-readable test type name for a given surface type.
// Maps surface types to test type names per docs/reference/test-type-model.md.
// Returns "Functional Test" as fallback for unknown or empty surface types.
func TestTypeTitle(surfaceType string) string {
	switch surfaceType {
	case "cli":
		return "CLI Functional Test"
	case "tui":
		return "Terminal Functional Test"
	case "api":
		return "API Functional Test"
	case "web":
		return "Web E2E Test"
	case "mobile":
		return "Mobile E2E Test"
	default:
		return "Functional Test"
	}
}

// GenSurfaceTestType generates a surface-specific test type name by appending
// the surface segment to the base type. Returns the base type unchanged when
// surface is empty.
func GenSurfaceTestType(baseType, surface string) string {
	if surface == "" {
		return baseType
	}
	return baseType + "." + surface
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
	// SurfaceKey is the user-defined surface identifier (e.g. "admin-panel").
	// Empty value means cross-surface (no specific surface).
	SurfaceKey string `json:"surface-key,omitempty"`
	// SurfaceType is the surface type enumeration (e.g. "web", "api", "cli").
	// Empty value means unknown or not yet resolved.
	SurfaceType string `json:"surface-type,omitempty"`
	// SourceTaskID records which task spawned this task (e.g. fix-task -> source task).
	// Empty for original tasks. Used by record auto-restore to unblock source when all deps complete.
	SourceTaskID string `json:"sourceTaskID,omitempty"`
	// MainSession indicates this task must run in the main session (not dispatched to task-executor).
	// Used by tasks that need to spawn subagents (e.g., eval-journey/eval-contract spawns doc-scorer/doc-reviser).
	MainSession bool `json:"mainSession,omitempty"`
	// Type is the task execution type (e.g. "coding.feature", "coding.fix", "gate").
	// Required for all tasks after migration; validated by task validate.
	// omitempty allows existing index.json files to load without error.
	Type string `json:"type,omitempty"`
	// BlockedReason records why a task entered blocked state.
	// Written by run-tasks when task prompt exits non-zero.
	BlockedReason string `json:"blockedReason,omitempty"`
	// Coverage is an optional per-task coverage override from frontmatter.
	// nil means use the global coverage config default for the task type.
	// Non-nil (including 0) means the task has an explicit coverage target.
	Coverage *int `json:"coverage,omitempty"`
	// Scope is the legacy scope field (deprecated, replaced by SurfaceKey/SurfaceType).
	// Retained solely for migration detection via CheckLegacyScope.
	// Loaded from both index.json (json tag) and task frontmatter (yaml tag).
	Scope string `json:"scope,omitempty" yaml:"scope"`
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
	// SurfaceKey mirrors Task.SurfaceKey for the claimed task.
	SurfaceKey string `json:"surface-key,omitempty"`
	// SurfaceType mirrors Task.SurfaceType for the claimed task.
	SurfaceType string `json:"surface-type,omitempty"`
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

	// Doc fields
	ReferencedDocs []string `json:"referencedDocs,omitempty"`
	ReviewStatus   string   `json:"reviewStatus,omitempty"`
	DocMetrics     string   `json:"docMetrics,omitempty"`

	// Test fields
	CasesGenerated int      `json:"casesGenerated,omitempty"`
	CasesEvaluated int      `json:"casesEvaluated,omitempty"`
	ScriptsCreated []string `json:"scriptsCreated,omitempty"`
	TestResults    string   `json:"testResults,omitempty"`

	// Validation fields
	ValidationPassed bool     `json:"validationPassed,omitempty"`
	IssuesFound      []string `json:"issuesFound,omitempty"`

	// Gate fields
	GatePassed bool     `json:"gatePassed,omitempty"`
	GateChecks []string `json:"gateChecks,omitempty"`

	// Eval fields
	Score    float64  `json:"score,omitempty"`
	Findings []string `json:"findings,omitempty"`
	Severity string   `json:"severity,omitempty"`
	Passed   bool     `json:"passed,omitempty"`
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
			"suspended",
			"skipped",
			"rejected",
		},
		PriorityEnum: []string{"P0", "P1", "P2"},
		tasks:        make(map[string]Task),
		Created:      time.Now().Format("2006-01-02"),
	}
}
