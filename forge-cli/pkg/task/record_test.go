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

// --- Doc record template tests ---

// goldenDocInput returns a task with doc type and fully populated doc fields.
func goldenDocInput() (*Task, *RecordData, string) {
	return &Task{
			ID:    "3",
			Title: "Doc record template",
		},
		&RecordData{
			Status:        "completed",
			Summary:       "Created doc-specific record template",
			FilesCreated:  []string{"pkg/task/data/record-doc.md"},
			FilesModified: []string{"pkg/task/record.go"},
			KeyDecisions:  []string{"Separate doc template from coding template"},
			Coverage:      -1.0,
			AcceptanceCriteria: []AcceptanceCriterion{
				{Criterion: "Template renders Document Metrics", Met: true},
				{Criterion: "No test-related sections", Met: true},
			},
			Notes:          "Doc tasks need no test metrics",
			DocMetrics:     "5 docs reviewed, 2 updated",
			ReferencedDocs: []string{"docs/guide.md", "docs/api.md"},
			ReviewStatus:   "Approved by tech lead",
		},
		"2026-05-23 10:00"
}

// goldenDocEmptyInput returns a doc task with all doc fields empty.
func goldenDocEmptyInput() (*Task, *RecordData, string) {
	return &Task{
			ID:    "5",
			Title: "Write README",
		},
		&RecordData{
			Status:   "completed",
			Summary:  "Added README",
			Coverage: -1.0,
		},
		"2026-05-23 10:00"
}

// goldenDocMixedInput returns a doc task with some fields populated, some empty.
func goldenDocMixedInput() (*Task, *RecordData, string) {
	return &Task{
			ID:    "6",
			Title: "Update API docs",
		},
		&RecordData{
			Status:         "completed",
			Summary:        "Updated API reference docs",
			FilesModified:  []string{"docs/api.md"},
			Coverage:       -1.0,
			DocMetrics:     "3 endpoints documented",
			ReferencedDocs: []string{"docs/architecture.md"},
		},
		"2026-05-23 10:00"
}

func TestRenderDocRecord(t *testing.T) {
	t.Run("populated fields", func(t *testing.T) {
		task, rd, startedTime := goldenDocInput()
		got := RenderDocRecord(task, rd, startedTime)

		// Shared sections
		assert.Contains(t, got, "# Task Record: 3 Doc record template")
		assert.Contains(t, got, "## Summary")
		assert.Contains(t, got, "Created doc-specific record template")
		assert.Contains(t, got, "### Files Created")
		assert.Contains(t, got, "- pkg/task/data/record-doc.md")
		assert.Contains(t, got, "### Files Modified")
		assert.Contains(t, got, "- pkg/task/record.go")
		assert.Contains(t, got, "### Key Decisions")
		assert.Contains(t, got, "- Separate doc template from coding template")
		assert.Contains(t, got, "## Acceptance Criteria")
		assert.Contains(t, got, "- [x] Template renders Document Metrics")
		assert.Contains(t, got, "## Notes")
		assert.Contains(t, got, "Doc tasks need no test metrics")

		// Doc-specific sections
		assert.Contains(t, got, "## Document Metrics")
		assert.Contains(t, got, "5 docs reviewed, 2 updated")
		assert.Contains(t, got, "## Referenced Documents")
		assert.Contains(t, got, "- docs/guide.md\n- docs/api.md")
		assert.Contains(t, got, "## Review Status")
		assert.Contains(t, got, "Approved by tech lead")

		// NO test-related sections
		assert.NotContains(t, got, "## Test Results")
		assert.NotContains(t, got, "Tests Executed")
		assert.NotContains(t, got, "Coverage")
	})

	t.Run("empty fields use fallbacks", func(t *testing.T) {
		task, rd, startedTime := goldenDocEmptyInput()
		got := RenderDocRecord(task, rd, startedTime)

		assert.Contains(t, got, "# Task Record: 5 Write README")

		// Empty fields should show fallbacks
		assert.Contains(t, got, "## Document Metrics")
		assert.Contains(t, got, "N/A")
		assert.Contains(t, got, "## Referenced Documents")
		assert.Contains(t, got, "无")
		assert.Contains(t, got, "## Review Status")
		assert.Contains(t, got, "N/A")

		// Notes fallback
		assert.Contains(t, got, "## Notes\n无")

		// Shared sections show fallbacks
		assert.Contains(t, got, "### Files Created\n无")
		assert.Contains(t, got, "### Files Modified\n无")
		assert.Contains(t, got, "### Key Decisions\n无")
		assert.Contains(t, got, "## Acceptance Criteria\n无")

		// Still no test sections
		assert.NotContains(t, got, "## Test Results")
		assert.NotContains(t, got, "Tests Executed")
	})

	t.Run("mixed populated and empty fields", func(t *testing.T) {
		task, rd, startedTime := goldenDocMixedInput()
		got := RenderDocRecord(task, rd, startedTime)

		assert.Contains(t, got, "# Task Record: 6 Update API docs")

		// Populated fields
		assert.Contains(t, got, "## Document Metrics")
		assert.Contains(t, got, "3 endpoints documented")
		assert.Contains(t, got, "## Referenced Documents")
		assert.Contains(t, got, "- docs/architecture.md")

		// Empty ReviewStatus should fallback to N/A
		assert.Contains(t, got, "## Review Status")
		assert.Contains(t, got, "N/A")

		// Files created empty
		assert.Contains(t, got, "### Files Created\n无")
		// Files modified populated
		assert.Contains(t, got, "### Files Modified\n- docs/api.md")
		// Key decisions empty
		assert.Contains(t, got, "### Key Decisions\n无")

		// Still no test sections
		assert.NotContains(t, got, "## Test Results")
	})

	t.Run("blocked doc task", func(t *testing.T) {
		task := &Task{ID: "7", Title: "Draft proposal"}
		rd := &RecordData{
			Status:       "blocked",
			Summary:      "Blocked on missing reference",
			DocMetrics:   "Draft in progress",
			ReviewStatus: "Pending review",
		}
		got := RenderDocRecord(task, rd, "2026-05-23 10:00")

		assert.Contains(t, got, `completed: "N/A"`)
		assert.Contains(t, got, "## Document Metrics")
		assert.Contains(t, got, "Draft in progress")
		assert.Contains(t, got, "## Review Status")
		assert.Contains(t, got, "Pending review")
		assert.NotContains(t, got, "## Test Results")
	})

	t.Run("type reclassification in doc record", func(t *testing.T) {
		task := &Task{ID: "8", Title: "Update docs"}
		rd := &RecordData{
			Status:   "completed",
			Summary:  "Updated docs",
			Coverage: -1.0,
			TypeReclassification: &TypeReclassification{
				OriginalType: "coding.feature",
				ActualType:   "doc",
				Reason:       "scope was documentation-only",
			},
		}
		got := RenderDocRecord(task, rd, "2026-05-23 10:00")

		assert.Contains(t, got, "## Type Reclassification")
		assert.Contains(t, got, "- Original: coding.feature")
		assert.Contains(t, got, "- Actual: doc")
		assert.Contains(t, got, "- Reason: scope was documentation-only")
		assert.NotContains(t, got, "## Test Results")
	})
}

func TestFormatWithFallback(t *testing.T) {
	t.Run("non-empty value", func(t *testing.T) {
		assert.Equal(t, "hello", formatWithFallback("hello", "fallback"))
	})
	t.Run("empty string", func(t *testing.T) {
		assert.Equal(t, "fallback", formatWithFallback("", "fallback"))
	})
	t.Run("whitespace only", func(t *testing.T) {
		assert.Equal(t, "fallback", formatWithFallback("   ", "fallback"))
	})
}

func TestRecordTemplateData_DocFields(t *testing.T) {
	t.Run("doc fields populated", func(t *testing.T) {
		task := &Task{ID: "3", Title: "Doc task"}
		rd := &RecordData{
			Status:         "completed",
			Summary:        "Doc work",
			Coverage:       -1.0,
			DocMetrics:     "5 docs",
			ReferencedDocs: []string{"a.md", "b.md"},
			ReviewStatus:   "Approved",
		}
		data := NewRecordTemplateData(task, rd, "2026-05-23 10:00")
		assert.Equal(t, "5 docs", data.DocMetricsFormatted)
		assert.Equal(t, "- a.md\n- b.md", data.ReferencedDocsFormatted)
		assert.Equal(t, "Approved", data.ReviewStatusFormatted)
	})

	t.Run("doc fields empty use fallbacks", func(t *testing.T) {
		task := &Task{ID: "4", Title: "Doc task"}
		rd := &RecordData{
			Status:   "completed",
			Summary:  "Doc work",
			Coverage: -1.0,
		}
		data := NewRecordTemplateData(task, rd, "2026-05-23 10:00")
		assert.Equal(t, "N/A", data.DocMetricsFormatted)
		assert.Equal(t, "无", data.ReferencedDocsFormatted)
		assert.Equal(t, "N/A", data.ReviewStatusFormatted)
	})
}

// --- Test record template tests ---

// goldenTestInput returns a test task with fully populated test-specific fields.
func goldenTestInput() (*Task, *RecordData, string) {
	return &Task{
			ID:    "T-1",
			Title: "Generate test cases",
		},
		&RecordData{
			Status:        "completed",
			Summary:       "Generated test cases from acceptance criteria",
			FilesCreated:  []string{"tests/cases/login.md", "tests/cases/api.md"},
			FilesModified: []string{"tests/index.json"},
			KeyDecisions:  []string{"Used journey-based case organization"},
			Coverage:      -1.0,
			AcceptanceCriteria: []AcceptanceCriterion{
				{Criterion: "Cases generated for all journeys", Met: true},
				{Criterion: "Each case has traceability to PRD", Met: true},
			},
			Notes:          "Auto-generated by test pipeline",
			CasesGenerated: 15,
			CasesEvaluated: 12,
			ScriptsCreated: []string{"tests/journeys/login_test.go", "tests/journeys/api_test.go"},
			TestResults:    "All scripts passed (12/12)",
		},
		"2026-05-23 10:00"
}

// goldenTestEmptyInput returns a test task with all test-specific fields empty.
func goldenTestEmptyInput() (*Task, *RecordData, string) {
	return &Task{
			ID:    "T-3",
			Title: "Run tests",
		},
		&RecordData{
			Status:   "completed",
			Summary:  "Ran test suite",
			Coverage: -1.0,
		},
		"2026-05-23 10:00"
}

// goldenTestPartialInput returns a test task with some fields populated.
func goldenTestPartialInput() (*Task, *RecordData, string) {
	return &Task{
			ID:    "T-5",
			Title: "Generate and run scripts",
		},
		&RecordData{
			Status:         "completed",
			Summary:        "Generated and ran scripts",
			FilesCreated:   []string{"tests/journeys/checkout_test.go"},
			Coverage:       -1.0,
			CasesGenerated: 8,
			TestResults:    "7 passed, 1 failed",
		},
		"2026-05-23 10:00"
}

func TestRenderTestRecord(t *testing.T) {
	t.Run("populated test fields", func(t *testing.T) {
		task, rd, startedTime := goldenTestInput()
		got := RenderTestRecord(task, rd, startedTime)

		// Shared sections present
		assert.Contains(t, got, "# Task Record: T-1 Generate test cases")
		assert.Contains(t, got, "## Summary")
		assert.Contains(t, got, "Generated test cases from acceptance criteria")
		assert.Contains(t, got, "### Files Created")
		assert.Contains(t, got, "- tests/cases/login.md\n- tests/cases/api.md")
		assert.Contains(t, got, "### Files Modified")
		assert.Contains(t, got, "- tests/index.json")
		assert.Contains(t, got, "### Key Decisions")
		assert.Contains(t, got, "- Used journey-based case organization")
		assert.Contains(t, got, "## Acceptance Criteria")
		assert.Contains(t, got, "- [x] Cases generated for all journeys")
		assert.Contains(t, got, "## Notes")
		assert.Contains(t, got, "Auto-generated by test pipeline")

		// Test-specific sections
		assert.Contains(t, got, "## Cases Generated")
		assert.Contains(t, got, "15")
		assert.Contains(t, got, "## Cases Evaluated")
		assert.Contains(t, got, "12")
		assert.Contains(t, got, "## Scripts Created")
		assert.Contains(t, got, "- tests/journeys/login_test.go\n- tests/journeys/api_test.go")
		assert.Contains(t, got, "## Test Results")
		assert.Contains(t, got, "All scripts passed (12/12)")

		// NO coding-specific sections
		assert.NotContains(t, got, "## Test Results\n- **Tests Executed**")
		assert.NotContains(t, got, "**Passed**")
		assert.NotContains(t, got, "**Failed**")
		assert.NotContains(t, got, "**Coverage**")
	})

	t.Run("empty fields use fallbacks", func(t *testing.T) {
		task, rd, startedTime := goldenTestEmptyInput()
		got := RenderTestRecord(task, rd, startedTime)

		assert.Contains(t, got, "# Task Record: T-3 Run tests")

		// Test-specific sections with fallbacks
		assert.Contains(t, got, "## Cases Generated\nN/A")
		assert.Contains(t, got, "## Cases Evaluated\nN/A")
		assert.Contains(t, got, "## Scripts Created\n无")
		assert.Contains(t, got, "## Test Results\nN/A")

		// Shared fallbacks
		assert.Contains(t, got, "### Files Created\n无")
		assert.Contains(t, got, "### Files Modified\n无")
		assert.Contains(t, got, "### Key Decisions\n无")
		assert.Contains(t, got, "## Acceptance Criteria\n无")
		assert.Contains(t, got, "## Notes\n无")
	})

	t.Run("partial fields", func(t *testing.T) {
		task, rd, startedTime := goldenTestPartialInput()
		got := RenderTestRecord(task, rd, startedTime)

		assert.Contains(t, got, "# Task Record: T-5 Generate and run scripts")

		// Populated fields
		assert.Contains(t, got, "## Cases Generated")
		assert.Contains(t, got, "8")
		assert.Contains(t, got, "## Test Results")
		assert.Contains(t, got, "7 passed, 1 failed")
		assert.Contains(t, got, "### Files Created")
		assert.Contains(t, got, "- tests/journeys/checkout_test.go")

		// Empty fields fallback
		assert.Contains(t, got, "## Cases Evaluated\nN/A")
		assert.Contains(t, got, "## Scripts Created\n无")
	})

	t.Run("blocked test task", func(t *testing.T) {
		task := &Task{ID: "T-7", Title: "Run regression"}
		rd := &RecordData{
			Status:      "blocked",
			Summary:     "Blocked on failing tests",
			TestResults: "3 failed, need investigation",
		}
		got := RenderTestRecord(task, rd, "2026-05-23 10:00")

		assert.Contains(t, got, `completed: "N/A"`)
		assert.Contains(t, got, "## Test Results")
		assert.Contains(t, got, "3 failed, need investigation")
		assert.NotContains(t, got, "**Coverage**")
	})

	t.Run("type reclassification in test record", func(t *testing.T) {
		task := &Task{ID: "T-9", Title: "Generate cases"}
		rd := &RecordData{
			Status:         "completed",
			Summary:        "Generated cases",
			Coverage:       -1.0,
			CasesGenerated: 5,
			TypeReclassification: &TypeReclassification{
				OriginalType: "test.gen-cases",
				ActualType:   "test.gen-and-run",
				Reason:       "scope expanded to include run step",
			},
		}
		got := RenderTestRecord(task, rd, "2026-05-23 10:00")

		assert.Contains(t, got, "## Type Reclassification")
		assert.Contains(t, got, "- Original: test.gen-cases")
		assert.Contains(t, got, "- Actual: test.gen-and-run")
		assert.Contains(t, got, "- Reason: scope expanded to include run step")
		assert.NotContains(t, got, "**Coverage**")
	})
}

func TestRecordTemplateData_TestFields(t *testing.T) {
	t.Run("test fields populated", func(t *testing.T) {
		task := &Task{ID: "T-1", Title: "Test task"}
		rd := &RecordData{
			Status:         "completed",
			Summary:        "Test work",
			Coverage:       -1.0,
			CasesGenerated: 10,
			CasesEvaluated: 8,
			ScriptsCreated: []string{"a_test.go", "b_test.go"},
			TestResults:    "All passed",
		}
		data := NewRecordTemplateData(task, rd, "2026-05-23 10:00")
		assert.Equal(t, "10", data.CasesGeneratedFormatted)
		assert.Equal(t, "8", data.CasesEvaluatedFormatted)
		assert.Equal(t, "- a_test.go\n- b_test.go", data.ScriptsCreatedFormatted)
		assert.Equal(t, "All passed", data.TestResultsFormatted)
	})

	t.Run("test fields empty use fallbacks", func(t *testing.T) {
		task := &Task{ID: "T-2", Title: "Test task"}
		rd := &RecordData{
			Status:   "completed",
			Summary:  "Test work",
			Coverage: -1.0,
		}
		data := NewRecordTemplateData(task, rd, "2026-05-23 10:00")
		assert.Equal(t, "N/A", data.CasesGeneratedFormatted)
		assert.Equal(t, "N/A", data.CasesEvaluatedFormatted)
		assert.Equal(t, "无", data.ScriptsCreatedFormatted)
		assert.Equal(t, "N/A", data.TestResultsFormatted)
	})
}

func TestFormatIntWithFallback(t *testing.T) {
	t.Run("positive value", func(t *testing.T) {
		assert.Equal(t, "10", formatIntWithFallback(10))
	})
	t.Run("zero value", func(t *testing.T) {
		assert.Equal(t, "N/A", formatIntWithFallback(0))
	})
	t.Run("negative value", func(t *testing.T) {
		assert.Equal(t, "N/A", formatIntWithFallback(-1))
	})
}
