---
created: 2026-05-23
author: "faner"
status: Draft
---

# Proposal: Type-Differentiated Task Record Templates

## Problem

All task types share a single execution record template (`fillRecordTemplate()` in `submit.go`). Doc tasks get Test Results/Coverage sections that are always "N/A"; test tasks lack fields for generated scripts and case counts; validation tasks lack pass/fail verdict fields. The record is the primary artifact for auditing what happened — when it contains irrelevant fields and omits relevant ones, it fails its purpose.

### Evidence

- `fillRecordTemplate()` (submit.go:367-448): one function generates records for all 21 task types via string concatenation
- Doc task records always render `Tests Executed: No`, `Coverage: N/A` — pure noise
- Test tasks have no fields for cases generated, scripts created, or regression results
- `RecordData` struct (types.go:276-289): all fields are coding-oriented (`testsPassed`, `testsFailed`, `coverage`)
- `submit-task` SKILL.md: one set of validation tiers applies to all types — warns on missing `testsPassed` even for doc tasks

### Urgency

Medium. The current records are functional but noisy. As more non-coding task types are added (test pipeline expansion, validation gates), the gap between record content and actual work done widens. Fixing it now avoids accumulating misleading records.

## Proposed Solution

Introduce type-specific record templates: 5 Go `text/template` files (coding, doc, test, validation, gate) selected at submit time based on task type category. Extend `RecordData` with optional field groups for each category. Update `submit-task` SKILL.md with type-specific record.json instructions and validation tiers.

### Innovation Highlights

Follows the established pattern already used by `prompt/data/*.md` (21 type-specific prompt templates) and `task/data/*.md` (13 auto-gen task templates). The innovation is extending this pattern to the record generation layer, completing the type-differentiation trifecta: prompt → execution → record.

## Requirements Analysis

### Key Scenarios

- Agent completes a doc task: record.json includes `referencedDocs`, `reviewStatus`, `docMetrics`, `notes`; generated record.md has "Document Metrics" section instead of "Test Results"
- Agent completes a coding task: unchanged from current behavior (backward compatible)
- Agent completes a test task (e.g., test.gen-cases): record.json includes `casesGenerated`, `casesEvaluated`; record.md shows test-pipeline-specific metrics
- Agent completes a validation task: record.json includes `validationPassed`, `issuesFound`; record.md shows pass/fail verdict
- Agent completes a gate task: minimal record with gate checks and pass status
- Existing records remain readable — no format breaking changes

### Non-Functional Requirements

- **Backward compatibility**: existing record files and index.json entries unchanged
- **Performance**: template rendering must be faster than current string concatenation (marginal)
- **Extensibility**: adding a new type category requires only a new template file + struct fields, no code changes to `fillRecordTemplate()`

### Constraints & Dependencies

- Must follow the forge plugin distribution model (see `docs/conventions/forge-distribution.md`)
- Template files embedded via `//go:embed` (same pattern as existing template dirs)
- `submit-task` SKILL.md changes affect agent behavior — must be tested with task-executor

## Alternatives & Industry Benchmarking

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | Zero cost | Records remain noisy, incomplete for non-coding types | Rejected: gap widens as test/validation types expand |
| Conditional branches in fillRecordTemplate() | — | Single file change | Function bloat, inconsistent with template-based architecture | Rejected: already have 2 template-based systems |
| **Template files per type category** | prompt/data/*.md pattern | Architecture-consistent, extensible, maintainable | More files, larger change set | **Selected: consistent with existing patterns** |

## Feasibility Assessment

### Technical Feasibility

Straightforward. Go `text/template` is already in use. `RecordData` struct extension is additive. No external dependencies.

### Resource & Timeline

Estimated 4-6 coding tasks. All changes are in forge-cli (Go) + skill docs (markdown).

### Dependency Readiness

No external dependencies. Internal dependency on existing type system (`pkg/task/types.go`).

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "All tasks need test results in records" | XY Detection | Confirmed as false: doc/test/validation tasks need different metrics. The real need (Y) is "records must accurately reflect what happened" |
| "One template is simpler" | Occam's Razor | Overturned: one template with growing conditional branches is more complex than N focused templates |
| "RecordData must be backward compatible" | Stress Test | Confirmed: adding optional fields is backward compatible; changing existing field semantics is not |

## Scope

### Phase 1 (Current): Doc Type De-noising

Minimal change to remove test-related noise from doc task records. No template engine — conditional branches in existing `fillRecordTemplate()`.

- `RecordData` struct: add doc-specific optional fields (`referencedDocs`, `reviewStatus`, `docMetrics`)
- `fillRecordTemplate()`: add category-based conditional rendering (doc type renders Document Metrics section instead of Test Results)
- `CategoryForType()` utility function: map task types to categories
- `validateRecordData()`: type-aware field validation (doc tasks don't require `testsPassed`)
- `submit-task` SKILL.md: type-specific record.json instructions (doc vs coding)
- Unit tests for all new logic

### Phase 2 (Future): Template Engine Migration

Introduce `text/template` + per-category template files for test/validation/gate types.

- 5 Go template files (record-coding.md, record-doc.md, record-test.md, record-validation.md, record-gate.md)
- `fillRecordTemplate()` refactored to template engine rendering
- Test task records: Cases Generated/Evaluated, Scripts Created sections
- Validation task records: Pass/Fail Verdict, Issues Found sections
- Gate task records: minimal gate checks + pass status
- Prompt templates updated to inform agents about type-specific record fields

### Out of Scope

- Task definition templates (task.md / task-doc.md — already differentiated)
- index.json structure changes
- Record file naming or path changes
- Breaking format changes to existing records

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Agent fills wrong fields for a type (e.g., testsPassed for doc task) | M | M | SKILL.md conditional instructions + type-aware validation |
| Backward incompatibility in RecordData serialization | L | H | All new fields are optional; existing fields unchanged |
| Phase 1 conditional branches become unmaintainable before Phase 2 | L | M | Phase 1 limited to doc vs non-doc split (2 branches only) |
| `code-quality.simplify` type unclassified | L | L | Phase 1 maps it to coding category (same as coding.*) |

## Success Criteria (Phase 1)

- [ ] Doc task records contain zero test-related sections (no "Test Results", no "Coverage")
- [ ] Doc task records include Document Metrics, Referenced Documents, Review Status sections
- [ ] Coding task records are identical to current format (backward compatible)
- [ ] `submit-task` SKILL.md has distinct record.json instructions for doc vs coding
- [ ] `forge task submit` validates required fields per type category
- [ ] `CategoryForType()` covers all 21 task types
- [ ] Unit tests for template rendering per category and type-aware validation

## Next Steps

- Proceed to `/quick-tasks` for task generation (via quick pipeline)
