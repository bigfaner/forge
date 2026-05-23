package task

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// goldenRecordInput is a comprehensive test input that exercises all fields
// of the coding record template, including TypeReclassification.
func goldenRecordInput() (*Task, *RecordData, string) {
	return &Task{
			ID:    "2.1",
			Title: "Implement template engine",
		},
		&RecordData{
			Status:        "completed",
			Summary:       "Introduced text/template for record generation",
			FilesCreated:  []string{"pkg/task/data/record-coding.md", "pkg/task/record.go"},
			FilesModified: []string{"internal/cmd/task/submit.go"},
			KeyDecisions:  []string{"Used text/template over string concatenation"},
			TestsPassed:   10,
			TestsFailed:   0,
			Coverage:      85.5,
			AcceptanceCriteria: []AcceptanceCriterion{
				{Criterion: "Template file created", Met: true},
				{Criterion: "Byte-identical output", Met: true},
			},
			Notes: "Infrastructure step",
			TypeReclassification: &TypeReclassification{
				OriginalType: "coding.fix",
				ActualType:   "coding.feature",
				Reason:       "scope expanded beyond initial fix",
			},
		},
		"2026-05-23 10:00"
}

// goldenMinimalInput has no reclassification and minimal fields.
func goldenMinimalInput() (*Task, *RecordData, string) {
	return &Task{
			ID:    "1.1",
			Title: "Write PRD",
		},
		&RecordData{
			Status:   "completed",
			Summary:  "Created PRD",
			Coverage: -1.0,
		},
		"2026-05-23 10:00"
}

// goldenBlockedInput tests a non-completed status with partial data.
func goldenBlockedInput() (*Task, *RecordData, string) {
	return &Task{
			ID:    "3.2",
			Title: "Fix compile error",
		},
		&RecordData{
			Status:      "blocked",
			Summary:     "Blocked due to failing tests",
			TestsPassed: 3,
			TestsFailed: 2,
			Coverage:    60.0,
			Notes:       "Auto-downgraded",
		},
		"2026-05-23 10:00"
}

func TestRenderCodingRecord_MatchesFillRecordTemplate(t *testing.T) {
	// This test verifies that the template-based RenderCodingRecord
	// produces byte-identical output to the current fillRecordTemplate.
	// Once the template engine is implemented, this test will pass.

	t.Run("golden input with full fields", func(t *testing.T) {
		task, rd, startedTime := goldenRecordInput()
		expected := FillRecordTemplate(task, rd, startedTime)
		got := RenderCodingRecord(task, rd, startedTime)
		assert.Equal(t, expected, got, "template output must be byte-identical to current output")
	})

	t.Run("minimal input with coverage=-1", func(t *testing.T) {
		task, rd, startedTime := goldenMinimalInput()
		expected := FillRecordTemplate(task, rd, startedTime)
		got := RenderCodingRecord(task, rd, startedTime)
		assert.Equal(t, expected, got, "template output must be byte-identical for minimal input")
	})

	t.Run("blocked status", func(t *testing.T) {
		task, rd, startedTime := goldenBlockedInput()
		expected := FillRecordTemplate(task, rd, startedTime)
		got := RenderCodingRecord(task, rd, startedTime)
		assert.Equal(t, expected, got, "template output must be byte-identical for blocked status")
	})

	t.Run("empty startedTime uses current time", func(t *testing.T) {
		task := &Task{ID: "1.1", Title: "Test"}
		rd := &RecordData{
			Status:      "completed",
			Summary:     "Done",
			TestsPassed: 1,
			Coverage:    50.0,
		}
		expected := FillRecordTemplate(task, rd, "")
		got := RenderCodingRecord(task, rd, "")
		assert.Equal(t, expected, got, "empty startedTime should produce same output")
	})

	t.Run("non-completed sets completed to N/A", func(t *testing.T) {
		task := &Task{ID: "1.2", Title: "WIP"}
		rd := &RecordData{
			Status:  "in_progress",
			Summary: "Work in progress",
		}
		expected := FillRecordTemplate(task, rd, "2026-05-23 10:00")
		got := RenderCodingRecord(task, rd, "2026-05-23 10:00")
		assert.Equal(t, expected, got)
		assert.Contains(t, got, `completed: "N/A"`)
	})
}

func TestRecordTemplateData(t *testing.T) {
	t.Run("all fields populated", func(t *testing.T) {
		task, rd, startedTime := goldenRecordInput()
		data := NewRecordTemplateData(task, rd, startedTime)

		assert.Equal(t, "completed", data.Status)
		assert.Equal(t, "2026-05-23 10:00", data.Started)
		assert.Equal(t, "2.1", data.TaskID)
		assert.Equal(t, "Implement template engine", data.TaskTitle)
		assert.Equal(t, "- pkg/task/data/record-coding.md\n- pkg/task/record.go", data.FilesCreatedFormatted)
		assert.Equal(t, "- internal/cmd/task/submit.go", data.FilesModifiedFormatted)
		assert.Equal(t, "- Used text/template over string concatenation", data.KeyDecisionsFormatted)
		assert.Equal(t, 10, data.TestsPassed)
		assert.Equal(t, 0, data.TestsFailed)
		assert.Equal(t, "85.5%", data.CoverageFormatted)
		assert.NotNil(t, data.TypeReclassification)
		assert.Equal(t, "Infrastructure step", data.Notes)
	})

	t.Run("nil reclassification", func(t *testing.T) {
		task, rd, startedTime := goldenMinimalInput()
		data := NewRecordTemplateData(task, rd, startedTime)
		assert.Nil(t, data.TypeReclassification)
	})

	t.Run("default notes when empty", func(t *testing.T) {
		task := &Task{ID: "1.1", Title: "Test"}
		rd := &RecordData{Status: "completed", Summary: "Done", TestsPassed: 1, Coverage: 50.0}
		data := NewRecordTemplateData(task, rd, "2026-05-23 10:00")
		assert.Equal(t, "无", data.Notes)
	})

	t.Run("non-completed completed time is N/A", func(t *testing.T) {
		task, rd, _ := goldenBlockedInput()
		data := NewRecordTemplateData(task, rd, "2026-05-23 10:00")
		assert.Equal(t, "N/A", data.Completed)
	})

	t.Run("timeSpent computed when valid", func(t *testing.T) {
		task := &Task{ID: "1.1", Title: "Test"}
		rd := &RecordData{Status: "completed", Summary: "Done", TestsPassed: 1, Coverage: 50.0}
		data := NewRecordTemplateData(task, rd, "2026-05-23 10:00")
		// completed time is "now" which is after 10:00, so timeSpent should be non-empty
		assert.NotEmpty(t, data.TimeSpent)
	})
}

func TestTemplateHelperFunctions(t *testing.T) {
	t.Run("formatList", func(t *testing.T) {
		assert.Equal(t, "无", templateFormatList(nil))
		assert.Equal(t, "无", templateFormatList([]string{}))
		assert.Equal(t, "- item1", templateFormatList([]string{"item1"}))
		assert.Equal(t, "- a\n- b", templateFormatList([]string{"a", "b"}))
	})

	t.Run("formatCoverage", func(t *testing.T) {
		assert.Equal(t, "N/A (task has no tests)", templateFormatCoverage(-1.0))
		assert.Equal(t, "85.5%", templateFormatCoverage(85.5))
		assert.Equal(t, "0.0%", templateFormatCoverage(0.0))
	})

	t.Run("formatTestsExecuted", func(t *testing.T) {
		assert.Equal(t, "No", templateFormatTestsExecuted(-1.0))
		assert.Equal(t, "Yes", templateFormatTestsExecuted(0.0))
		assert.Equal(t, "Yes", templateFormatTestsExecuted(85.5))
	})

	t.Run("formatCriteria", func(t *testing.T) {
		assert.Equal(t, "无", templateFormatCriteria(nil))
		assert.Equal(t, "无", templateFormatCriteria([]AcceptanceCriterion{}))
		assert.Equal(t, "- [x] Pass", templateFormatCriteria([]AcceptanceCriterion{{Criterion: "Pass", Met: true}}))
		assert.Equal(t, "- [ ] Fail", templateFormatCriteria([]AcceptanceCriterion{{Criterion: "Fail", Met: false}}))
	})

	t.Run("formatDuration", func(t *testing.T) {
		assert.Equal(t, "~45m", templateFormatDuration(45*time.Minute))
		assert.Equal(t, "~2h", templateFormatDuration(2*time.Hour))
		assert.Equal(t, "~2h 30m", templateFormatDuration(2*time.Hour+30*time.Minute))
		assert.Equal(t, "~0m", templateFormatDuration(0))
	})
}

func TestFillRecordTemplate_Unchanged(t *testing.T) {
	// Sanity: FillRecordTemplate still works as before (we're not removing it yet).
	t.Run("still produces output", func(t *testing.T) {
		task, rd, startedTime := goldenRecordInput()
		got := FillRecordTemplate(task, rd, startedTime)
		assert.Contains(t, got, "2.1")
		assert.Contains(t, got, "Implement template engine")
		assert.Contains(t, got, "## Type Reclassification")
		assert.Contains(t, got, "- Original: coding.fix")
	})
}
