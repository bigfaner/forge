---
status: "completed"
started: "2026-05-21 23:13"
completed: "2026-05-21 23:16"
time_spent: "~3m"
---

# Task Record: 3 Restructure ui.md and rewrite mobile.md with golden rules

## Summary
Restructured ui.md into Golden Rules + Reconnaissance Hints dual-zone structure with Session Reuse, Network Interception, and Viewport Management golden rules. Rewrote mobile.md from Maestro YAML tutorial into framework-agnostic principles with App State Reset, Permission Handling, Deep Link Pattern, Element Location Strategy, Touch and Gesture Principles, Screen Transition Assertions, and Application Lifecycle. Both files now reference _shared.md, use tests/<journey>/ output path, and have correct step number references.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-test-scripts/types/ui.md
- plugins/forge/skills/gen-test-scripts/types/mobile.md

### Key Decisions
- mobile.md Golden Rules section contains zero Maestro-specific syntax -- all Maestro patterns extracted into abstract principles (e.g., tapOn priority table became 'Accessibility ID > Resource ID > Text content' priority chain)
- Maestro YAML reference examples moved to Reconnaissance Hints with HTML comment markers to prevent LLM from treating them as generation instructions
- ui.md Generation Patterns section removed entirely -- its contents were framework-specific (Playwright code examples) that belong in Convention files, not type-level golden rules
- Element Location Strategy expressed as abstract priority chains in both files rather than concrete selector syntax, per Hard Rules

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] Each file has Golden Rules section (framework-agnostic constraints)
- [x] Each file has Reconnaissance Hints section with discovery-only annotation
- [x] Golden Rules reference _shared.md principles instead of redefining
- [x] Output path changed to tests/<journey>/
- [x] All step number references corrected to match SKILL.md
- [x] UI: Session Reuse golden rule added
- [x] UI: Network Interception golden rule added
- [x] UI: Viewport Management golden rule added
- [x] Mobile: App State Reset golden rule added
- [x] Mobile: Permission Handling golden rule added
- [x] Mobile: Deep Link Pattern golden rule added
- [x] Mobile: Framework-agnostic generation patterns (no Maestro YAML in Golden Rules)
- [x] Mobile: Maestro YAML examples in Reconnaissance Hints marked as reference
- [x] Mobile: Element Location Strategy as abstract priority chain
- [x] Hard Rule: mobile.md Golden Rules contain ZERO Maestro-specific syntax
- [x] Hard Rule: ui.md Golden Rules contain no Playwright-specific code
- [x] Hard Rule: Element location strategies are abstract priority chains

## Notes
mobile.md was the heaviest rewrite -- the original file was 208 lines of essentially a Maestro tutorial. The rewritten version is 208 lines but entirely framework-agnostic in Golden Rules. Maestro patterns are preserved as reference examples in Reconnaissance Hints for reconnaissance value.
