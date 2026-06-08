---
created: "2026-06-08"
tags: [testing, gotcha, pipeline, template]
---

# Auto-gen task templates missing AC rendering cause validation failure

## Problem

`forge task validate` reports "has 0 acceptance criteria" for auto-generated tasks (T-test-gen-scripts, T-test-run, T-quick-doc-drift, T-review-doc), causing quick-tasks and breakdown pipelines to fail validation at Step 7.

## Root Cause

1. **Template body missing AC section**: Auto-gen task templates in `pkg/task/templates/` declare `AcceptanceCriteria` in their frontmatter `context:` metadata (making AC data available to prompt templates), but 4 templates (`test-gen-scripts.md`, `test-run.md`, `doc-drift.md`, `doc-review.md`) omit `## Acceptance Criteria\n\n{{.AcceptanceCriteria}}` from their body.

2. **Default AC mechanism bypassed**: `autogen.go` has a fallback `defaultAcceptanceCriteria = "- [ ] All acceptance criteria met"` in `buildAutogenTemplateData`, but since the template body never renders `{{.AcceptanceCriteria}}`, the default value is computed but discarded.

3. **Validation reads .md not template**: `forge task validate` parses the generated `.md` files' `## Acceptance Criteria` section, not the template metadata. When the section is absent, it reports 0 AC.

## Solution

Add `## Acceptance Criteria\n\n{{.AcceptanceCriteria}}` to the body of each missing template. The `buildAutogenTemplateData` default fallback ensures at least one AC is rendered even when no PRD-derived criteria exist.

## Reusable Pattern

When creating auto-gen task templates that declare `AcceptanceCriteria` in `context:`, always include the rendering directive `{{.AcceptanceCriteria}}` in the template body. Declaring it in metadata without rendering it means the AC data reaches the prompt template (via `context`) but never appears in the generated `.md` task file that validation checks.
