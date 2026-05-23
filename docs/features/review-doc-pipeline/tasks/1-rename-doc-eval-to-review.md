---
id: "1"
title: "Rename doc.eval → doc.review and rewrite templates"
priority: "P0"
estimated_time: "1.5h"
dependencies: []
scope: "backend"
breaking: true
type: "coding.refactor"
mainSession: false
---

# 1: Rename doc.eval → doc.review and rewrite templates

## Description

Rename the `doc.eval` type to `doc.review` across the entire Go CLI codebase, and rewrite the associated templates from the broken 8-dimension rubric model to an AC checklist + direct-fix model.

The current `doc.eval` system is a shell: the autogen template references a non-existent rubric, and the agent prompt uses a 1000-point scoring system that is overkill for doc tasks. The new `doc.review` model uses AC checklist核对 with direct fixes — matching the actual workflow of doc tasks (check deliverables against acceptance criteria, fix non-conformances).

## Reference Files
- `docs/proposals/review-doc-pipeline/proposal.md` — Source proposal

## Acceptance Criteria

- [ ] `TypeDocEval` constant renamed to `TypeDocReview` with value `"doc.review"` in `types.go`
- [ ] `doc.eval` removed from `ValidTypes` and `SystemTypes`, replaced with `doc.review`
- [ ] `TaskTypeRegistry` entry updated: `"doc.review"` with description "review documentation against acceptance criteria"
- [ ] `GetDocEvalTask()` renamed to `GetReviewDocTask()` in `autogen.go`, with `Key: "review-doc"`, `ID: "T-review-doc"`, `Title: "Review Documentation Quality"`, `Type: TypeDocReview`
- [ ] `ResolveDocEvalDep` renamed to `ResolveReviewDocDep` in `autogen.go`
- [ ] `autogenTypeToFile` mapping updated: `TypeDocReview → "data/doc-review.md"`
- [ ] `needsDocEval()` renamed to `needsReviewDoc()` in `build.go`
- [ ] All call sites updated in `build.go`
- [ ] `IsAutoGenTaskID` updated: `"T-eval-doc"` → `"T-review-doc"` in `build.go` (or `autogen.go` wherever defined)
- [ ] `infer.go` type inference updated: `"T-review-doc" → TypeDocReview`
- [ ] `prompt.go` template mapping updated: `TypeDocReview → "data/doc-review.md"`
- [ ] `forge-cli/pkg/task/data/doc-eval.md` renamed to `doc-review.md`, content rewritten to describe AC checklist model (no 8-dimension rubric reference)
- [ ] `forge-cli/pkg/prompt/data/doc-eval.md` renamed to `doc-review.md`, content rewritten to AC核对 + direct fix workflow (read task deliverables → check each AC → fix non-conformances → report)
- [ ] `go build ./...` passes with zero errors after all renames
- [ ] No remaining references to `doc.eval`, `DocEval`, `T-eval-doc`, or `doc-eval` in Go source or templates

## Hard Rules

- Keep `category.go` logic unchanged — `strings.HasPrefix(typ, "doc")` already covers `doc.review`
- Do NOT add `doc.review` to user-facing type tables in skills (it remains system-internal, auto-generated only)
- The `doc-eval.md` files must be deleted, not left as dead files

## Implementation Notes

Files to modify (in order of dependency):
1. `forge-cli/pkg/task/types.go` — constant, ValidTypes, SystemTypes, TaskTypeRegistry
2. `forge-cli/pkg/task/autogen.go` — function renames, map updates
3. `forge-cli/pkg/task/build.go` — function rename + call sites
4. `forge-cli/pkg/task/infer.go` — ID → type mapping
5. `forge-cli/pkg/prompt/prompt.go` — template path mapping
6. `forge-cli/pkg/task/data/doc-eval.md` → `doc-review.md` — rewrite autogen template
7. `forge-cli/pkg/prompt/data/doc-eval.md` → `doc-review.md` — rewrite agent prompt

The autogen template (`task/data/doc-review.md`) should instruct:
- Scan docs/features/{{FEATURE_SLUG}}/ for all documents created/modified
- For each doc task, read its acceptance criteria from the task .md file
- Check each deliverable against its AC
- Report pass/fail per AC item

The agent prompt (`prompt/data/doc-review.md`) should implement:
- Step 1: Load task definition, identify all doc tasks in the feature
- Step 2: For each doc task, read its deliverables and acceptance criteria
- Step 3: Check each AC — if not met, directly modify the document to fix
- Step 4: Report summary (which ACs passed, which were fixed)
- No scoring — just pass/fail per AC with direct fixes
