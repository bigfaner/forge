---
status: "completed"
started: "2026-05-28 01:59"
completed: "2026-05-28 02:14"
time_spent: "~15m"
---

# Task Record: 7 统一模板 metadata frontmatter

## Summary
Added unified metadata frontmatter (type, category, variables) to all 41 template files across 3 Go packages (prompt, task, template). Refactored 6 record templates to dual-frontmatter structure. Implemented metadata stripping in template loaders before template.Parse(). Implemented ValidateTemplates() with reflection-based cross-validation of metadata variables against struct fields. Created new metadata.go and metadata_test.go in pkg/prompt.

## Changes

### Files Created
- forge-cli/pkg/prompt/metadata.go
- forge-cli/pkg/prompt/metadata_test.go

### Files Modified
- forge-cli/pkg/prompt/prompt.go
- forge-cli/pkg/task/autogen.go
- forge-cli/pkg/task/record.go
- forge-cli/pkg/template/template.go
- forge-cli/scripts/version.txt
- forge-cli/pkg/prompt/data/agent-executor.md
- forge-cli/pkg/prompt/data/code-review.md
- forge-cli/pkg/prompt/data/coding.cleanup.md
- forge-cli/pkg/prompt/data/coding.fix.md
- forge-cli/pkg/prompt/data/coding.feature.md
- forge-cli/pkg/prompt/data/coding.refactor.md
- forge-cli/pkg/prompt/data/commit.md
- forge-cli/pkg/prompt/data/doc.architecture.md
- forge-cli/pkg/prompt/data/doc.business-rules.md
- forge-cli/pkg/prompt/data/doc.conventions.md
- forge-cli/pkg/prompt/data/doc.overview.md
- forge-cli/pkg/prompt/data/doc.reference.md
- forge-cli/pkg/prompt/data/doc.workflow.md
- forge-cli/pkg/prompt/data/eval-proposal.md
- forge-cli/pkg/prompt/data/eval.md
- forge-cli/pkg/prompt/data/graduate-tests.md
- forge-cli/pkg/prompt/data/gate.md
- forge-cli/pkg/prompt/data/git-commit.md
- forge-cli/pkg/prompt/data/smart-commit.md
- forge-cli/pkg/prompt/data/submit-task.md
- forge-cli/pkg/prompt/data/task-plan.md
- forge-cli/pkg/prompt/data/test.functional.md
- forge-cli/pkg/prompt/data/test.run.md
- forge-cli/pkg/prompt/data/validation.md
- forge-cli/pkg/prompt/data/code-quality.simplify.md
- forge-cli/pkg/template/data/coding.fix.md
- forge-cli/pkg/template/data/coding.cleanup.md
- forge-cli/pkg/task/data/agent-executor.md
- forge-cli/pkg/task/data/commit.md
- forge-cli/pkg/task/data/eval.md
- forge-cli/pkg/task/data/eval-proposal.md
- forge-cli/pkg/task/data/gate.md
- forge-cli/pkg/task/data/git-commit.md
- forge-cli/pkg/task/data/graduate-tests.md
- forge-cli/pkg/task/data/smart-commit.md
- forge-cli/pkg/task/data/submit-task.md
- forge-cli/pkg/task/data/task-plan.md
- forge-cli/pkg/task/data/test.functional.md
- forge-cli/pkg/task/data/test.run.md
- forge-cli/pkg/task/data/validation.md
- forge-cli/pkg/task/data/record-validation.md
- forge-cli/pkg/task/data/record-test.md
- forge-cli/pkg/task/data/record-gate.md
- forge-cli/pkg/task/data/record-eval.md
- forge-cli/pkg/task/data/record-coding.md
- forge-cli/pkg/task/data/record-doc.md

### Key Decisions
- Simple line-based YAML parser for metadata frontmatter instead of full YAML library to avoid new dependencies
- Metadata is optional for backward compatibility - templates without frontmatter are returned unchanged
- Record templates use dual-frontmatter: metadata frontmatter (stripped before parse) + output frontmatter (rendered by Go template)
- Inline reflect.Ptr/reflect.Struct constants to satisfy govet inline-const rule
- Removed unused error return from parseMetadataFrontmatter/parseAutogenMetadata to satisfy unparam linter

## Test Results
- **Tests Executed**: Yes
- **Passed**: 10
- **Failed**: 0
- **Coverage**: 74.1%

## Acceptance Criteria
- [x] All 41 template files have metadata frontmatter with type, category, and variables fields
- [x] Record templates use dual-frontmatter structure (metadata + output)
- [x] Template loaders strip metadata before template.Parse()
- [x] ValidateTemplates() cross-validates metadata variables against struct fields via reflection
- [x] All existing tests continue to pass

## Notes
Version bumped to 5.14.0. Coverage is 74.1% for pkg/prompt (new metadata.go), 86.0% for pkg/task, 37.0% for pkg/template.
