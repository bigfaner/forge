---
status: "completed"
started: "2026-05-29 18:01"
completed: "2026-05-29 18:06"
time_spent: "~5m"
---

# Task Record: T-quick-doc-drift Detect Spec Drift

## Summary
Detected and fixed 1 spec drift in forge-distribution.md: Quick Pipeline description was outdated ('1-15 coding tasks') and missing intent-driven pipeline branching documentation. Added Intent-Driven Pipeline Branching section with mode/test-pipeline behavior table for new-feature/refactor/cleanup intents.

## Changes

### Files Created
无

### Files Modified
- docs/conventions/forge-distribution.md

### Key Decisions
无

## Document Metrics
14 specs checked, 1 drift found and fixed, 0 false positives

## Referenced Documents
- docs/business-rules/task-lifecycle.md
- docs/business-rules/quality-gate.md
- docs/business-rules/surface-orchestration.md
- docs/conventions/forge-distribution.md
- docs/conventions/prompt-template-hierarchy.md
- docs/conventions/enum-constants.md
- docs/conventions/skill-self-containment.md
- docs/conventions/skill-structure.md
- docs/conventions/forge-cli-reference.md
- docs/conventions/surface-cli.md
- docs/conventions/surface-rules.md
- docs/conventions/code-structure.md
- docs/conventions/dispatcher-quality.md
- docs/conventions/error-handling.md

## Review Status
final

## Acceptance Criteria
- [x] Run git diff --name-only main...HEAD to identify changed files
- [x] List all spec files in docs/business-rules/ and docs/conventions/
- [x] Read domains frontmatter for each spec file
- [x] Only verify specs whose domains overlap with changed files
- [x] Skip specs with no domain overlap
- [x] Auto-fix drifted specs

## Notes
Drift found: forge-distribution.md Quick Pipeline section was stale (said '1-15 coding tasks' but quick-tasks now supports unlimited coding+doc tasks). Also missing intent-driven pipeline branching documentation. Verified no drift in: task-lifecycle.md (SystemTypes=12 correct, IsAutoGenTaskID patterns match), quality-gate.md (3-phase pipeline correct), enum-constants.md (pkg/types/ leaf package correct), prompt-template-hierarchy.md, skill-self-containment.md, skill-structure.md, surface-orchestration.md, forge-cli-reference.md. Note: write-prd/SKILL.md is 366 lines exceeding the 350-line convention but that is a code violation, not a spec drift.
