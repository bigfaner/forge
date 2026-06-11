---
id: "1"
title: "Create embed template files for auto-gen task types and update GenerateTestTaskMD()"
priority: "P1"
estimated_time: "1-2h"
dependencies: []
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 1: Create embed template files for auto-gen task types

## Description

Create `forge-cli/pkg/task/data/` directory with one .md file per auto-gen task type (13 types total). Each file contains the Main Instructions body content for that task type. Update `GenerateTestTaskMD()` in `autogen.go` to use `embed.FS` to load the appropriate template file as the task body, replacing the current sparse body generation logic.

## Reference Files
- `docs/proposals/auto-task-main-instructions/proposal.md` — Source proposal

## Acceptance Criteria

- [ ] `forge-cli/pkg/task/data/` directory created with 13 .md template files
- [ ] Each template file contains meaningful Main Instructions for its task type
- [ ] Type-to-filename mapping added (convention: type with '.' replaced by '-', e.g., `test.gen-cases` → `test-gen-cases.md`)
- [ ] `GenerateTestTaskMD()` uses `embed.FS` to load template file as body
- [ ] When `StrategyContent` is non-empty, it's appended AFTER the template content
- [ ] When `TestType` is non-empty, it's noted in the body (as current behavior)
- [ ] `forge task index --feature <slug>` generates task files with the new body content
- [ ] Existing tests pass (backward compatible frontmatter)

## Hard Rules

- MUST NOT change frontmatter generation logic in `GenerateTestTaskMD()`
- MUST NOT change the `AutoGenTaskDef` struct
- MUST use `embed.FS` (same pattern as `forge-cli/pkg/prompt/prompt.go`)
- MUST handle missing template file gracefully (fallback to current behavior)

## Implementation Notes

### Auto-gen task types needing templates:

| Type constant | Value | Template filename |
|---------------|-------|-------------------|
| TypeTestGenCases | test.gen-cases | test-gen-cases.md |
| TypeTestEvalCases | test.eval-cases | test-eval-cases.md |
| TypeTestGenScripts | test.gen-scripts | test-gen-scripts.md |
| TypeTestGenAndRun | test.gen-and-run | test-gen-and-run.md |
| TypeTestRun | test.run | test-run.md |
| TypeTestGraduate | test.graduate | test-graduate.md |
| TypeTestVerifyRegression | test.verify-regression | test-verify-regression.md |
| TypeValidationCode | validation.code | validation-code.md |
| TypeValidationUx | validation.ux | validation-ux.md |
| TypeDocEval | doc.eval | doc-eval.md |
| TypeDocConsolidate | doc.consolidate | doc-consolidate.md |
| TypeDocDrift | doc.drift | doc-drift.md |
| TypeCleanCode | code-quality.simplify | code-quality-simplify.md |

### Embed FS pattern (same as prompt.go):

```go
//go:embed data/*.md
var autogenTemplateFS embed.FS
```

### Body generation logic:

```
1. Load template from embed FS via type→filename mapping
2. If template found: use as body
3. If TestType non-empty: append type line
4. If StrategyContent non-empty: append after template content
5. If template not found: fallback to current behavior
```
