---
status: "completed"
started: "2026-06-06 17:18"
completed: "2026-06-06 17:24"
time_spent: "~6m"
---

# Task Record: 6 拆分 runExtract 304 行函数

## Summary
Refactored runExtract (304 lines, 7+ nesting) into 14 focused functions across parse/aggregate/output phases. All functions <= 80 lines, all nesting <= 4 levels, file reduced from 384 to 427 lines (within 500-line limit). Zero behavior change verified by all 60+ forensic tests passing.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/forensic/extract.go

### Key Decisions
- Extracted parseJSONLEntries as main parsing loop with parseState struct for cross-line timestamp tracking
- Split entry-type handling into parseAssistantEntry, parseUserEntry, parseAttachmentEntry
- Extracted recordToolUse and matchPendingToolUse for tool_use/tool_result block processing
- Extracted aggregateTimings and computeTimeRange for post-parse aggregation phase
- Extracted writeExtractOutput for output phase (JSON serialization + file writing)
- Extracted resolveOutDir and newExtractResult for initialization
- Added updateTimestamps helper to centralize timestamp state updates

## Test Results
- **Tests Executed**: Yes
- **Passed**: 63
- **Failed**: 0
- **Coverage**: 90.3%

## Acceptance Criteria
- [x] runExtract and all extracted subfunctions <= 80 lines
- [x] All functions nesting <= 4 levels (early return / guard clause)
- [x] go test ./... all green, zero behavior change
- [x] File <= 500 lines

## Notes
Original runExtract was 304 lines with 7+ nesting. Refactored into 14 functions (max 43 lines, max nesting 4). The pendingCall type was moved from a local closure to package-level since it's shared across parseAssistantEntry and parseUserEntry.
