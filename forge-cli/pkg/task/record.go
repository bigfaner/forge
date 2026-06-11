package task

import (
	"embed"
	"fmt"
	"strings"
	"text/template"
	"time"

	"forge-cli/pkg/types"
)

//go:embed records/*.md
var recordTemplateFS embed.FS

// RecordTemplateData combines all fields needed by record templates.
type RecordTemplateData struct {
	Status                      string
	Started                     string
	Completed                   string
	TimeSpent                   string
	TaskID                      string
	TaskTitle                   string
	Summary                     string
	FilesCreatedFormatted       string
	FilesModifiedFormatted      string
	KeyDecisionsFormatted       string
	TestsExecuted               string
	TestsPassed                 int
	TestsFailed                 int
	CoverageFormatted           string
	AcceptanceCriteriaFormatted string
	Notes                       string
	TypeReclassification        *TypeReclassification

	// Doc fields (used by record-doc.md template)
	DocMetricsFormatted     string
	ReferencedDocsFormatted string
	ReviewStatusFormatted   string

	// Test fields (used by record-test.md template)
	CasesGeneratedFormatted string
	CasesEvaluatedFormatted string
	ScriptsCreatedFormatted string
	TestResultsFormatted    string

	// Validation fields (used by record-validation.md template)
	ValidationPassedFormatted string
	IssuesFoundFormatted      string

	// Gate fields (used by record-gate.md template)
	GateChecksFormatted string
	GatePassedFormatted string

	// Eval fields (used by record-eval.md template)
	ScoreFormatted    string
	FindingsFormatted string
	SeverityFormatted string
	PassedFormatted   string
}

// NewRecordTemplateData creates a RecordTemplateData from task, record data, and started time.
func NewRecordTemplateData(t *Task, rd *RecordData, startedTime string) *RecordTemplateData {
	status := rd.Status
	started := startedTime
	if started == "" {
		started = time.Now().Format("2006-01-02 15:04")
	}
	completed := time.Now().Format("2006-01-02 15:04")
	if status != string(types.StatusCompleted) {
		completed = "N/A"
	}

	timeSpent := ""
	startedT, err1 := time.Parse("2006-01-02 15:04", started)
	completedT, err2 := time.Parse("2006-01-02 15:04", completed)
	if err1 == nil && err2 == nil && completedT.After(startedT) {
		dur := completedT.Sub(startedT)
		timeSpent = FormatDuration(dur)
	}

	notes := rd.Notes
	if notes == "" {
		notes = "无"
	}

	return &RecordTemplateData{
		Status:                      status,
		Started:                     started,
		Completed:                   completed,
		TimeSpent:                   timeSpent,
		TaskID:                      t.ID,
		TaskTitle:                   t.Title,
		Summary:                     rd.Summary,
		FilesCreatedFormatted:       FormatList(rd.FilesCreated),
		FilesModifiedFormatted:      FormatList(rd.FilesModified),
		KeyDecisionsFormatted:       FormatList(rd.KeyDecisions),
		TestsExecuted:               FormatTestsExecuted(rd.Coverage),
		TestsPassed:                 rd.TestsPassed,
		TestsFailed:                 rd.TestsFailed,
		CoverageFormatted:           FormatCoverage(rd.Coverage),
		AcceptanceCriteriaFormatted: FormatCriteria(rd.AcceptanceCriteria),
		Notes:                       notes,
		TypeReclassification:        rd.TypeReclassification,
		DocMetricsFormatted:         formatWithFallback(rd.DocMetrics, "N/A"),
		ReferencedDocsFormatted:     FormatList(rd.ReferencedDocs),
		ReviewStatusFormatted:       formatWithFallback(rd.ReviewStatus, "N/A"),
		CasesGeneratedFormatted:     formatIntWithFallback(rd.CasesGenerated),
		CasesEvaluatedFormatted:     formatIntWithFallback(rd.CasesEvaluated),
		ScriptsCreatedFormatted:     FormatList(rd.ScriptsCreated),
		TestResultsFormatted:        formatWithFallback(rd.TestResults, "N/A"),
		ValidationPassedFormatted:   FormatBool(rd.ValidationPassed, "Passed", "Failed"),
		IssuesFoundFormatted:        FormatList(rd.IssuesFound),
		GateChecksFormatted:         FormatList(rd.GateChecks),
		GatePassedFormatted:         FormatBool(rd.GatePassed, "Yes", "No"),
		ScoreFormatted:              formatScore(rd.Score),
		FindingsFormatted:           FormatList(rd.Findings),
		SeverityFormatted:           formatWithFallback(rd.Severity, "N/A"),
		PassedFormatted:             FormatBool(rd.Passed, "Yes", "No"),
	}
}

// recordFuncMap provides helper functions available in record templates.
var recordFuncMap = template.FuncMap{
	"formatList":          templateFormatList,
	"formatCoverage":      templateFormatCoverage,
	"formatTestsExecuted": templateFormatTestsExecuted,
	"formatCriteria":      templateFormatCriteria,
	"formatDuration":      templateFormatDuration,
}

// templateFormatList formats a string slice as a markdown list.
func templateFormatList(items []string) string {
	return FormatList(items)
}

// templateFormatCoverage formats a coverage value.
func templateFormatCoverage(c float64) string {
	return FormatCoverage(c)
}

// templateFormatTestsExecuted returns "Yes" or "No" based on coverage.
func templateFormatTestsExecuted(c float64) string {
	return FormatTestsExecuted(c)
}

// templateFormatCriteria formats acceptance criteria.
func templateFormatCriteria(criteria []AcceptanceCriterion) string {
	return FormatCriteria(criteria)
}

// templateFormatDuration formats a duration.
func templateFormatDuration(dur time.Duration) string {
	return FormatDuration(dur)
}

// formatWithFallback returns the value if non-empty, otherwise returns the fallback.
func formatWithFallback(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

// formatIntWithFallback returns the string representation of n if positive,
// otherwise returns "N/A".
func formatIntWithFallback(n int) string {
	if n > 0 {
		return fmt.Sprintf("%d", n)
	}
	return "N/A"
}

// FormatBool returns trueVal if cond is true, otherwise falseVal.
func FormatBool(cond bool, trueVal, falseVal string) string {
	if cond {
		return trueVal
	}
	return falseVal
}

// formatScore formats an eval score for display.
func formatScore(score float64) string {
	if score <= 0 {
		return "N/A"
	}
	return fmt.Sprintf("%.0f/1000", score)
}

// FormatCoverage formats coverage value for display.
func FormatCoverage(c float64) string {
	if c < 0 {
		return "N/A (task has no tests)"
	}
	return fmt.Sprintf("%.1f%%", c)
}

// FormatTestsExecuted returns "Yes" or "No" based on coverage.
func FormatTestsExecuted(c float64) string {
	if c < 0 {
		return "No"
	}
	return "Yes"
}

// FormatList formats a string slice as a markdown bullet list.
func FormatList(items []string) string {
	if len(items) == 0 {
		return "无"
	}
	lines := make([]string, len(items))
	for i, item := range items {
		lines[i] = "- " + item
	}
	return strings.Join(lines, "\n")
}

// FormatDuration formats a duration as a human-readable string.
func FormatDuration(dur time.Duration) string {
	d := int(dur.Hours())
	m := int(dur.Minutes()) % 60
	switch {
	case d > 0 && m > 0:
		return fmt.Sprintf("~%dh %dm", d, m)
	case d > 0:
		return fmt.Sprintf("~%dh", d)
	default:
		return fmt.Sprintf("~%dm", m)
	}
}

// FormatCriteria formats acceptance criteria as a markdown checklist.
func FormatCriteria(criteria []AcceptanceCriterion) string {
	if len(criteria) == 0 {
		return "无"
	}
	lines := make([]string, len(criteria))
	for i, c := range criteria {
		check := "[ ]"
		if c.Met {
			check = "[x]"
		}
		lines[i] = "- " + check + " " + c.Criterion
	}
	return strings.Join(lines, "\n")
}

// RenderCodingRecord renders the coding record template with the given data.
func RenderCodingRecord(t *Task, rd *RecordData, startedTime string) string {
	return renderRecordTemplate("records/coding.md", t, rd, startedTime)
}

// RenderRecord dispatches to the appropriate category-specific renderer.
func RenderRecord(t *Task, rd *RecordData, startedTime string) string {
	switch CategoryForType(t.Type) {
	case CategoryDoc:
		return RenderDocRecord(t, rd, startedTime)
	case CategoryTest:
		return RenderTestRecord(t, rd, startedTime)
	case CategoryValidation:
		return RenderValidationRecord(t, rd, startedTime)
	case CategoryGate:
		return RenderGateRecord(t, rd, startedTime)
	case CategoryEval:
		return RenderEvalRecord(t, rd, startedTime)
	default:
		return RenderCodingRecord(t, rd, startedTime)
	}
}

// RenderDocRecord renders the doc record template with the given data.
func RenderDocRecord(t *Task, rd *RecordData, startedTime string) string {
	return renderRecordTemplate("records/doc.md", t, rd, startedTime)
}

// RenderTestRecord renders the test record template with the given data.
func RenderTestRecord(t *Task, rd *RecordData, startedTime string) string {
	return renderRecordTemplate("records/test.md", t, rd, startedTime)
}

// RenderValidationRecord renders the validation record template with the given data.
func RenderValidationRecord(t *Task, rd *RecordData, startedTime string) string {
	return renderRecordTemplate("records/validation.md", t, rd, startedTime)
}

// RenderGateRecord renders the gate record template with the given data.
func RenderGateRecord(t *Task, rd *RecordData, startedTime string) string {
	return renderRecordTemplate("records/gate.md", t, rd, startedTime)
}

// RenderEvalRecord renders the eval record template with the given data.
func RenderEvalRecord(t *Task, rd *RecordData, startedTime string) string {
	return renderRecordTemplate("records/eval.md", t, rd, startedTime)
}

// renderRecordTemplate renders a named record template with the given data.
// Metadata frontmatter is stripped before parsing; the output frontmatter
// (status, started, etc.) remains in the template body and is rendered by Go template.
func renderRecordTemplate(templateName string, t *Task, rd *RecordData, startedTime string) string {
	data, err := recordTemplateFS.ReadFile(templateName)
	if err != nil {
		// Fallback: should never happen with embedded templates
		return fmt.Sprintf("ERROR: template %s not found: %v", templateName, err)
	}

	// Strip metadata frontmatter before parsing (metadata is not part of rendered output)
	content := stripAutogenMetadata(string(data))

	tmpl, err := template.New("record").Funcs(recordFuncMap).Parse(content)
	if err != nil {
		return fmt.Sprintf("ERROR: parse template %s: %v", templateName, err)
	}

	td := NewRecordTemplateData(t, rd, startedTime)

	var buf strings.Builder
	if err := tmpl.Execute(&buf, td); err != nil {
		return fmt.Sprintf("ERROR: execute template %s: %v", templateName, err)
	}

	return buf.String()
}

// FillRecordTemplate generates a record using the current string-concatenation method.
// This is the original implementation preserved for backward compatibility and testing.
func FillRecordTemplate(t *Task, rd *RecordData, startedTime string) string {
	status := rd.Status
	started := startedTime
	if started == "" {
		started = time.Now().Format("2006-01-02 15:04")
	}
	completed := time.Now().Format("2006-01-02 15:04")
	if status != string(types.StatusCompleted) {
		completed = "N/A"
	}

	timeSpent := ""
	startedT, err1 := time.Parse("2006-01-02 15:04", started)
	completedT, err2 := time.Parse("2006-01-02 15:04", completed)
	if err1 == nil && err2 == nil && completedT.After(startedT) {
		dur := completedT.Sub(startedT)
		timeSpent = FormatDuration(dur)
	}

	notes := rd.Notes
	if notes == "" {
		notes = "无"
	}

	var reclassBlock string
	if rd.TypeReclassification != nil {
		reclassBlock = fmt.Sprintf(`## Type Reclassification
- Original: %s
- Actual: %s
- Reason: %s

`, rd.TypeReclassification.OriginalType, rd.TypeReclassification.ActualType, rd.TypeReclassification.Reason)
	}

	return fmt.Sprintf(`---
status: "%s"
started: "%s"
completed: "%s"
time_spent: "%s"
---

# Task Record: %s %s

## Summary
%s

%s## Changes

### Files Created
%s

### Files Modified
%s

### Key Decisions
%s

## Test Results
- **Tests Executed**: %s
- **Passed**: %d
- **Failed**: %d
- **Coverage**: %s

## Acceptance Criteria
%s

## Notes
%s
`,
		status, started, completed, timeSpent,
		t.ID, t.Title,
		rd.Summary,
		reclassBlock,
		FormatList(rd.FilesCreated),
		FormatList(rd.FilesModified),
		FormatList(rd.KeyDecisions),
		FormatTestsExecuted(rd.Coverage), rd.TestsPassed, rd.TestsFailed, FormatCoverage(rd.Coverage),
		FormatCriteria(rd.AcceptanceCriteria),
		notes,
	)
}
