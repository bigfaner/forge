---
id: "5"
title: "Validation and Gate record templates"
priority: "P1"
estimated_time: "1h"
dependencies: ["2"]
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 5: Validation and Gate record templates

## Description

Create two template files:

1. **`record-validation.md`** for `validation.code` and `validation.ux` tasks. Sections: Pass/Fail Verdict, Issues Found.

2. **`record-gate.md`** for `gate` tasks. Minimal template: gate checks list + overall pass status.

## Reference Files
- `docs/proposals/typed-task-records/proposal.md` — Source proposal
- `forge-cli/pkg/task/data/record-coding.md` — Reference template structure (from task 2)

## Acceptance Criteria
- [ ] `record-validation.md` template created with "Pass/Fail Verdict" (`.ValidationPassed`) and "Issues Found" (`.IssuesFound` list) sections
- [ ] `record-gate.md` template created with "Gate Checks" (`.GateChecks` list) and "Gate Status" (`.GatePassed`) sections
- [ ] Gate template is minimal: Summary + Gate Checks + Gate Status + Notes. No Changes, no Criteria sections.
- [ ] Validation template includes shared sections: Summary, Changes, Key Decisions, Criteria, Notes
- [ ] `fillRecordTemplate()` routes validation-category and gate-category types to respective templates
- [ ] Unit tests for both templates

## Hard Rules
- Gate template must be minimal — gate tasks are quality checkpoints, not implementation tasks. They produce verdicts, not artifacts.
- Validation template includes Changes section because `validation.code` and `validation.ux` may produce artifacts (reports, screenshots).

## Implementation Notes
- Both templates are small. They can be created in one task since they share the same pattern and similar complexity.
