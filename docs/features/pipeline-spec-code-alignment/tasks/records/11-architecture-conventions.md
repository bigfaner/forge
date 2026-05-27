---
status: "completed"
started: "2026-05-27 01:25"
completed: "2026-05-27 01:28"
time_spent: "~3m"
---

# Task Record: 11 Fix architecture and conventions docs

## Summary
Fixed architecture and conventions docs: removed fictitious agent descriptions from ARCHITECTURE.md, corrected scope resolution algorithm to match actual ResolveScope code, fixed forge forge typo, updated dispatcher-quality.md to use just abstractions and add coding.cleanup, replaced interfaces with surfaces in gen-contracts docs, fixed just test references in clean-code and fix-bug, corrected execute-task.md frontmatter description

## Changes

### Files Created
无

### Files Modified
- docs/ARCHITECTURE.md
- docs/conventions/dispatcher-quality.md
- plugins/forge/skills/gen-contracts/SKILL.md
- plugins/forge/skills/gen-contracts/rules/validation.md
- plugins/forge/skills/clean-code/SKILL.md
- plugins/forge/commands/fix-bug.md
- plugins/forge/commands/execute-task.md

### Key Decisions
无

## Document Metrics
7 files modified, 8 acceptance criteria met

## Referenced Documents
- docs/proposals/pipeline-spec-code-alignment/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] ARCHITECTURE.md does not describe doc-scorer/doc-reviser as independent agents
- [x] ARCHITECTURE.md scope resolution algorithm matches actual Go code
- [x] No forge forge duplicate in ARCHITECTURE.md
- [x] dispatcher-quality.md uses just abstractions and mentions coding.cleanup
- [x] gen-contracts docs use surfaces terminology (not interfaces)
- [x] clean-code/SKILL.md references just unit-test
- [x] fix-bug.md references correct just target for running tests
- [x] execute-task.md frontmatter description is accurate

## Notes
Preserved document structure and heading hierarchy per Hard Rules. Historical references in Why sections left untouched.
