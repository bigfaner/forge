---
id: "3"
title: "Refactor InferType and lookup functions to derive from registry"
priority: "P1"
estimated_time: "1.5h"
complexity: "medium"
dependencies: [1]
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 3: Refactor InferType and Lookup Functions to Derive from Registry

## Description
Replace the 15-case InferType switch with registry iteration + prefix/suffix fallback. Refactor `isTestTaskID`, `isAutoGenForDep`, `IsAutoGenTaskID` to derive their results from the PipelineRegistry instead of hardcoded string lists.

## Reference Files
- `forge-cli/pkg/task/infer.go:16-70`: InferType 15-case switch — replaced by registry iteration (source: proposal.md#Derived-Functions)
- `forge-cli/pkg/task/build.go:516-530`: `isTestTaskID` 6-prefix hardcoded list — derive from registry expanded IDs (source: proposal.md#isTestTaskID)
- `forge-cli/pkg/task/autogen.go:1080-1091`: `isAutoGenForDep` — derive from registry (source: proposal.md#isAutoGenForDep)
- `forge-cli/pkg/task/build.go:614-626`: `IsAutoGenTaskID` — derive from registry (source: proposal.md#IsAutoGenTaskID)

## Acceptance Criteria
- [ ] InferType iterates PipelineRegistry, matches ID patterns with wildcard support for `{surface-key}`/`{surface-type}` placeholders
- [ ] Single surface degenerate IDs (e.g., `T-test-run` without suffix) matched by template when only one surface exists
- [ ] Prefix/suffix fallback covers runtime tasks (fix-*, doc-fix-*, disc-*) and stage-gate tasks (*.gate, *.summary)
- [ ] `isTestTaskID`, `isAutoGenForDep`, `IsAutoGenTaskID` build lookup sets from registry expanded IDs at init time
- [ ] All existing task IDs correctly typed: T-review-doc → TypeDocReview, T-test-gen-scripts-api → TypeTestGenScripts, T-test-run-cli → TypeTestRun, etc.
- [ ] `go build ./...` passes

## Implementation Notes
- The expanded ID set is built once at init time from PipelineRegistry + surfaces config
- InferType keeps current two-phase signature: `InferType(taskID string, surfaces map[string]string) string`
- The registry iteration approach auto-covers new task types added to PipelineRegistry — no manual switch case needed
