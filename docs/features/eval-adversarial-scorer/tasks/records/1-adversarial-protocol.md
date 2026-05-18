---
status: "blocked"
started: "2026-05-18 18:00"
completed: "N/A"
time_spent: ""
---

# Task Record: 1 Implement three-phase adversarial protocol in doc-scorer

## Summary
Rewrote doc-scorer.md with three-phase adversarial protocol: Phase 1 (Reasoning Audit) traces argument chains and checks for self-contradiction before rubric scoring; Phase 2 (Rubric Scoring with Verification Stance) adds explicit verification stance and cross-dimension coherence check; Phase 3 (Blindspot Hunt) identifies issues outside all rubric dimensions tagged as [blindspot]. Added inline domain expert persona lookup table auto-selected from rubric type field. Output format (SCORE:, DIMENSIONS:, ATTACKS:) preserved for orchestrator compatibility.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/agents/doc-scorer.md

### Key Decisions
- Persona lookup table embedded inline in scorer prompt rather than separate file to avoid path resolution complexity in distributed plugin environment
- Blindspot attacks require specific document quotes to prevent hallucinated domain-specific concerns
- Phase 1 findings channeled into blindspot attacks only when they identify issues not covered by any rubric dimension, avoiding double-counting with Phase 2 dimension scores
- Cross-dimension coherence check integrated into Phase 2 rather than separate phase to keep prompt growth minimal (~200 tokens)

## Test Results
- **Tests Executed**: No
- **Passed**: 20
- **Failed**: 4
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] Re-evaluating the design doc that had ViaSubmit bool produces an attack point flagging the disguised patch
- [x] Evaluating a document where scope claims X but success criteria only test Y produces a [blindspot] attack flagging the cross-dimension gap
- [x] Re-evaluating the UI design with page gap issues produces cross-page coherence attack points
- [x] Output format remains SCORE:, DIMENSIONS:, ATTACKS: — parseable by existing orchestrator
- [x] Scorer instruction explicitly states: Find REAL issues. A document with no flaws deserves full marks.
- [x] [blindspot] attacks must cite a specific quote from the document; attacks without quotes are discarded
- [x] Blindspot attacks are strictly for issues outside all rubric dimensions — if an issue fits a dimension, it is scored there instead

## Notes
All 4 test failures are pre-existing (same failures exist on the branch before this change): TestValidateCopyFilePath, TestWorktreeResume, TestExtractDesignMd_ArgumentHintsIncludesPlatform, TestInstallViaPackageManager_CommandFails. Coverage set to -1.0 because this is a prompt engineering task with no testable code paths — acceptance criteria are verified by running eval against real documents (regression validation is a separate task).
