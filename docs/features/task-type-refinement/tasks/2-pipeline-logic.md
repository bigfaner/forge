---
id: "2"
title: "Replace isDocsOnlyFeature with needsTestPipeline and needsDocEval"
priority: "P0"
estimated_time: "1h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 2: Replace isDocsOnlyFeature with needsTestPipeline and needsDocEval

## Description

Refactor `isDocsOnlyFeature()` into two focused functions: `needsTestPipeline()` (returns true when any business task has a testable runtime behavior type) and `needsDocEval()` (returns true when all business tasks are documentation-only). This enables the three-tier decision: docs-only → T-eval-doc, testable → test pipeline, cleanup/refactor-only → no pipeline.

## Reference Files
- `docs/proposals/task-type-refinement/proposal.md` — Source proposal (D2: Test pipeline logic)
- `forge-cli/pkg/task/build.go` — `isDocsOnlyFeature()`, `BuildIndex()` uses it
- `forge-cli/internal/cmd/quality_gate.go` — Second `isDocsOnly()` function (lines 112-119)

## Acceptance Criteria
- [ ] `needsTestPipeline(tasks) bool` returns true if any non-auto-gen task has type `feature`, `enhancement`, or `fix`
- [ ] `needsDocEval(tasks) bool` returns true if ALL non-auto-gen tasks have type `documentation`
- [ ] `isDocsOnlyFeature()` removed from `build.go`
- [ ] `BuildIndex()` updated: uses `needsTestPipeline()` for test pipeline generation, `needsDocEval()` for T-eval-doc generation
- [ ] `quality_gate.go` `isDocsOnly()` updated to use `needsTestPipeline()` logic (skip gate only when no testable types)
- [ ] Three-tier behavior verified: docs-only feature gets T-eval-doc; feature/enhancement/fix feature gets test pipeline; cleanup/refactor-only feature gets neither

## Hard Rules
- The `testableTypes` check must use the new constants from task 1 (`TypeFeature`, `TypeEnhancement`, `TypeFix`), not hardcoded strings.

## Implementation Notes
- The proposal D2 defines the exact `testableTypes` map and function signatures. Follow that design.
- `quality_gate.go` has a separate `isDocsOnly` that does NOT exclude auto-gen tasks — this is intentional for the quality gate skip logic. Keep this behavior: if ANY task (including auto-gen) is non-doc, run the gate. Only change the type check to use the new testable types.
