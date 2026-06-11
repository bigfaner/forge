---
status: "completed"
started: "2026-06-03 19:22"
completed: "2026-06-03 19:28"
time_spent: "~6m"
---

# Task Record: 6 L2 Conventions Audit Batch 2

## Summary
Audited 10 docs/conventions/ files (naming.md, package-organization.md, prompt-template-hierarchy.md, skill-self-containment.md, skill-structure.md, surface-cli.md, surface-rules.md, testing/index.md, testing/cli/index.md, testing/cli/core.md) against codebase. Found 16 issues: 0 P0, 7 P1, 6 P2, 3 P3. Key findings: 3 non-existent pkg/ packages in naming.md (version, lesson, research), missing qualitygate subpackage in package-organization.md, skill-structure.md constraints violated by auxiliary directories and 5 oversized SKILL.md files, testing/index.md only covers 1 of 5 surface types.

## Changes

### Files Created
- docs/features/global-doc-code-audit/audit/l2-conventions-batch2-report.md

### Files Modified
无

### Key Decisions
无

## Document Metrics
10 files audited, 16 issues found (P0:0 P1:7 P2:6 P3:3), 5 cross-layer influence items identified

## Referenced Documents
- docs/proposals/global-doc-code-audit/proposal.md
- docs/conventions/naming.md
- docs/conventions/skill-structure.md
- docs/conventions/testing/index.md

## Review Status
final

## Acceptance Criteria
- [x] All 10 target files audited with declaration extraction
- [x] Each convention verified against codebase: naming patterns vs actual code, skill structure vs actual files, test conventions vs actual test setup
- [x] Every inconsistency recorded: file path, line range, severity (P0-P3), suggested action
- [x] Cross-layer influence items recorded for L3 reference
- [x] Audit report follows unified template

## Notes
No modifications to code or documentation per Hard Rules. naming.md had 3 non-existent packages. skill-structure.md had 5 SKILL.md files exceeding 350-line limit and 6 skills using non-standard auxiliary directories. testing/ conventions had incomplete surface type coverage and incorrect directory paths.
