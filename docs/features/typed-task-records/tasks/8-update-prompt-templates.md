---
id: "8"
title: "Update prompt templates for type-specific record field awareness"
priority: "P2"
estimated_time: "1h"
dependencies: ["7"]
scope: "backend"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 8: Update prompt templates for type-specific record field awareness

## Description

Update the type-specific prompt templates in `forge-cli/pkg/prompt/data/` to inform agents about category-appropriate record fields. Currently, all prompts give uniform record guidance. After this change, each prompt tells the agent which RecordData fields to populate based on task category.

## Reference Files
- `docs/proposals/typed-task-records/proposal.md` — Source proposal
- `forge-cli/pkg/prompt/data/*.md` — 21 type-specific prompt templates
- `forge-cli/pkg/task/category.go` — CategoryForType (from task 1)

## Acceptance Criteria
- [ ] Coding prompt templates mention: testsPassed, testsFailed, coverage are required for completed tasks
- [ ] Doc prompt templates mention: referencedDocs, reviewStatus, docMetrics are recommended fields
- [ ] Test prompt templates mention: casesGenerated, casesEvaluated, scriptsCreated are relevant fields
- [ ] Validation prompt templates mention: validationPassed, issuesFound are relevant fields
- [ ] Gate prompt template mentions: gatePassed, gateChecks are relevant fields
- [ ] Changes are additive (append record-field hints) — no restructuring of existing prompt content

## Hard Rules
- Follow forge plugin distribution model for prompt template files
- Changes must be additive only — do not restructure or remove existing prompt template content
- The prompt template files are embedded via `//go:embed` — verify they parse correctly after changes

## Implementation Notes
- Each prompt template already has a section about record submission. Add a "Record Fields" hint block specific to the category.
- Group templates by category: all coding.* get the same hint, all doc* get the same hint, etc.
- Keep hints concise — 2-3 lines listing the key fields for that category.
