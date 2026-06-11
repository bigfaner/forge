---
id: "1"
title: "Add complexity data pipeline and conditional template rendering"
priority: "P0"
estimated_time: "2h"
dependencies: []
surface-key: "."
surface-type: "cli"
breaking: true
type: "coding.enhancement"
mainSession: false
---

# 1: Add complexity data pipeline and conditional template rendering

## Description

Implement the full data pipeline to pass the `complexity` field from task frontmatter through to prompt template rendering, including conditional paragraph support in `cleanTemplateOutput()`.

The pipeline chain: FrontmatterData struct → Task struct → index.json serialization → renderTemplate() → cleanTemplateOutput() conditional paragraph deletion.

## Reference Files
- `docs/proposals/task-pipeline-precision/proposal.md#Constraints-&-Dependencies` — defines the 4-step data pipeline chain and cleanTemplateOutput() conditional paragraph mechanism
- `docs/proposals/task-pipeline-precision/proposal.md#Feasibility-Assessment` — describes FrontmatterData → Task → index.json → renderTemplate synchronization requirement
- `docs/proposals/task-pipeline-precision/proposal.md#Success-Criteria` — SC-6 and SC-7 define verification: forge prompt get-by-task-id output and backward compatibility

## Acceptance Criteria

- [ ] `FrontmatterData` struct has `Complexity` field (string: "low"/"medium"/"high", optional)
- [ ] `Task` struct has `Complexity` field
- [ ] `index.json` serialization/deserialization handles the `complexity` field — existing tasks without it default to `"medium"`
- [ ] `renderTemplate()` adds `{{COMPLEXITY}}` placeholder replacement
- [ ] `cleanTemplateOutput()` supports conditional paragraph deletion: paragraphs wrapped in `<!-- IF NOT_LOW -->...<!-- END_IF -->` are removed when complexity is "low"
- [ ] `forge prompt get-by-task-id <task-with-complexity-low>` output does NOT contain Step 1.5 spec-code scan paragraph
- [ ] `forge prompt get-by-task-id <task-without-complexity>` output includes Step 1.5 (default medium behavior unchanged)

## Hard Rules

- The 4-step pipeline (FrontmatterData → Task → index.json → renderTemplate) MUST be modified synchronously. Missing any step causes the field to be silently lost.
- The `<!-- IF NOT_LOW -->` marker format must be documented in a comment in `cleanTemplateOutput()` so future template authors know the convention.

## Implementation Notes

### Test Impact
- Affected test suite(s): `forge-cli/internal/...`, `forge-cli/pkg/prompt/...`
- Expected fixture changes: prompt synthesis test fixtures may need complexity field added
- Risk level: medium

Key files to modify (in order):
1. `forge-cli/internal/frontmatter/frontmatter.go` — add Complexity to FrontmatterData
2. `forge-cli/internal/types/task.go` — add Complexity to Task struct
3. `forge-cli/pkg/prompt/prompt.go` — add {{COMPLEXITY}} to renderTemplate, extend cleanTemplateOutput for conditional paragraphs
4. Any index.json serialization paths that need updating
