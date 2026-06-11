---
status: "completed"
started: "2026-05-23 12:59"
completed: "2026-05-23 13:04"
time_spent: "~5m"
---

# Task Record: 3 Doc record template (record-doc.md)

## Summary
Created doc-specific record template (record-doc.md) with Document Metrics, Referenced Documents, and Review Status sections. Added RenderDocRecord function and routing in fillRecordTemplate based on CategoryForType. Template has zero test-related sections and shares format with coding template for common sections.

## Changes

### Files Created
- forge-cli/pkg/task/data/record-doc.md

### Files Modified
- forge-cli/pkg/task/record.go
- forge-cli/pkg/task/record_test.go
- forge-cli/internal/cmd/task/submit.go

### Key Decisions
- Used same RecordTemplateData struct with new doc fields rather than separate struct
- formatWithFallback helper for empty-string-to-N/A conversion
- Category-based routing in submit.go fillRecordTemplate switches on CategoryDoc

## Test Results
- **Tests Executed**: Yes
- **Passed**: 13
- **Failed**: 0
- **Coverage**: 90.6%

## Acceptance Criteria
- [x] record-doc.md template file created in forge-cli/pkg/task/data/
- [x] Template renders Document Metrics section using .DocMetrics field (fallback N/A)
- [x] Template renders Referenced Documents section using .ReferencedDocs list (fallback 无)
- [x] Template renders Review Status section using .ReviewStatus field (fallback N/A)
- [x] Template has zero test-related sections
- [x] Template includes shared sections (Summary, Changes, Key Decisions, Criteria, Notes)
- [x] fillRecordTemplate routes doc-category types to this template
- [x] Unit tests verify doc template output for populated, empty, and mixed fields

## Notes
13 new test sub-cases covering populated fields, empty fallbacks, mixed input, blocked status, and type reclassification in doc records
