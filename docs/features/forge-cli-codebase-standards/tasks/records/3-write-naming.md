---
status: "completed"
started: "2026-05-30 22:07"
completed: "2026-05-30 22:11"
time_spent: "~4m"
---

# Task Record: 3 Write naming.md convention

## Summary
Created docs/conventions/naming.md covering file names, function names, constant names, and package names with normative target-state rules, module-level deviation summary, and executable verification commands.

## Changes

### Files Created
- docs/conventions/naming.md

### Files Modified
无

### Key Decisions
无

## Document Metrics
~220 lines across 7 sections covering 4 naming domains and 3 deviation categories

## Referenced Documents
- docs/conventions/package-organization.md
- docs/conventions/enum-constants.md
- docs/conventions/code-structure.md

## Review Status
final

## Acceptance Criteria
- [x] docs/conventions/naming.md exists, covers file names, function names, constant names, package names
- [x] Contains target state definitions (normative rules)
- [x] Contains module-level deviation summary
- [x] Rules are executable (verifiable via grep, go vet, or manual review)

## Notes
Analyzed naming patterns across forge-cli/internal/cmd/ (15+ top-level files, 7 subpackages) and forge-cli/pkg/ (17 packages). Identified 3 deviation categories: compound-word package names (N1: infocmd, forgeconfig), descriptive compound packages (N2: facttable, serverprobe, testrunner), and inconsistent run-function prefixes in worktree subpackage (N3).
