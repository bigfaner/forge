---
id: "3"
title: "Doc record template (record-doc.md)"
priority: "P0"
estimated_time: "1h"
dependencies: ["2"]
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 3: Doc record template (record-doc.md)

## Description

Create `record-doc.md` template for doc-category tasks (`doc`, `doc.eval`, `doc.summary`, `doc.consolidate`, `doc.drift`). Replaces test-related sections with document-specific sections: Document Metrics, Referenced Documents, Review Status.

## Reference Files
- `docs/proposals/typed-task-records/proposal.md` — Source proposal
- `forge-cli/pkg/task/data/record-coding.md` — Reference template structure (from task 2)

## Acceptance Criteria
- [ ] `record-doc.md` template file created in `forge-cli/pkg/task/data/`
- [ ] Template renders "Document Metrics" section using `.DocMetrics` field (fallback "N/A")
- [ ] Template renders "Referenced Documents" section using `.ReferencedDocs` list (fallback "无")
- [ ] Template renders "Review Status" section using `.ReviewStatus` field (fallback "N/A")
- [ ] Template has **zero** test-related sections (no "Test Results", no "Coverage", no "Tests Executed")
- [ ] Template includes: Summary, Changes (Created/Modified), Key Decisions, Acceptance Criteria, Notes sections (shared with coding template)
- [ ] `fillRecordTemplate()` routes doc-category types to this template
- [ ] Unit tests verify doc template output for: populated fields, empty fields, mixed populated/empty

## Hard Rules
- No test metrics in doc template — not even "N/A" placeholders. The entire "Test Results" block is omitted.
- Shared sections (Summary, Changes, Key Decisions, Criteria, Notes) should use the same format as the coding template for consistency.
- Follow the forge plugin distribution model for any template file changes.

## Implementation Notes
- The template uses the same `recordTemplateData` struct from task 2 — just accesses different fields.
- Consider using `{{define "shared-sections"}}...{{end}}` blocks if Go templates support is available, to reduce duplication between templates. Otherwise, duplicate shared sections is acceptable for Phase 2.
