---
status: "completed"
started: "2026-05-17 21:50"
completed: "2026-05-17 21:55"
time_spent: "~5m"
---

# Task Record: 6 Add auto-extract trigger to run-tasks

## Summary
Added knowledge auto-extraction trigger to run-tasks Post-Completion section. After the main task loop ends and e2e suggestions are printed, a new Knowledge Review step reads the shared extraction routine (knowledge-extraction.md) and executes its flow with trigger=run-tasks and artifacts covering task outcomes, code changes (git diff), and manifest. The extraction flow handles scanning for notable knowledge, silent exit on routine tasks, user confirmation via AskUserQuestion, and writing confirmed knowledge to appropriate directories.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/commands/run-tasks.md

### Key Decisions
- Placed Knowledge Review as a subsection under Post-Completion, after loop summary and e2e suggestions, preserving existing flow
- Included shared extraction routine by reference (read file + execute flow) per Hard Rule, not by copying content
- Added guard: skip knowledge review if loop ended due to 3 consecutive failures (incomplete feature)

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] After the main task loop ends (Post-Completion section), a new knowledge review step runs
- [x] Step reads plugins/forge/references/shared/knowledge-extraction.md for extraction logic
- [x] Scans task outcomes, code changes, and manifest for notable knowledge
- [x] Looks for: architectural decisions, novel patterns, gotchas, business rules
- [x] Silent when no notable knowledge detected (routine tasks)
- [x] Presents extracted knowledge via AskUserQuestion for user confirmation
- [x] Writes confirmed knowledge to appropriate directories using shared formats
- [x] Does not interfere with existing Post-Completion flow (test/e2e suggestions)

## Notes
Prompt-only task (Markdown command file modification). No executable code produced. Quality gate (compile/fmt/lint/test) all passed. Coverage set to -1.0 since no testable code was produced.
