---
id: "fix-4"
title: "Fix: TestCheckAllCompleted_NoProject env leakage"
priority: "P0"
estimated_time: "30min"
dependencies: []
status: pending
breaking: true
---

# Fix: TestCheckAllCompleted_NoProject env leakage

## Root Cause

Test expects nil when no project root, but picks up forge-testing-optimization feature from ancestor directory Z:\project\ai\forge. The fix-3 changes to NoProjectRoot handling in integration_test.go may need the same treatment in all_completed_test.go.

## Reference Files

- Source: {{SOURCE_FILES}}
- Test script: {{TEST_SCRIPT}}
- Test results: {{TEST_RESULTS}}

## Verification

After fixing, verify the fix works:
1. `just test [scope]` — must pass
2. If UI/page related: `just test-e2e --feature <slug>` — must also pass

When this task is recorded as completed via `task record`, the source task fix-3 is automatically restored to pending if all its dependencies are completed.
