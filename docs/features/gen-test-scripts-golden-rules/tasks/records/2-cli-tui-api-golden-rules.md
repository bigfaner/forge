---
status: "completed"
started: "2026-05-21 23:09"
completed: "2026-05-21 23:12"
time_spent: "~3m"
---

# Task Record: 2 Restructure cli.md, tui.md, api.md with golden rules

## Summary
Restructured cli.md, tui.md, and api.md into Golden Rules + Reconnaissance Hints dual-zone structure. Added CLI-specific rules (Two-Level Timeout, Binary Isolation, Environment Hermeticity), TUI-specific rules (Terminal Size Contract, ANSI Sanitization, Stable State Detection), and API-specific rules (Idempotency Check, Request Timeout, Content-Type Verification). Fixed Output path from tests/e2e/features/<feature>/ to tests/<journey>/. Fixed step number references from 'Step 1.5' to 'Step 1.3 Fact Table build'. Removed shared antipattern guards replaced by _shared.md references. Added Discovery hints annotation to Reconnaissance Hints sections.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-test-scripts/types/cli.md
- plugins/forge/skills/gen-test-scripts/types/tui.md
- plugins/forge/skills/gen-test-scripts/types/api.md

### Key Decisions
- Golden Rules section contains zero language-specific code, import paths, or grep commands per Hard Rules
- Key Sequence Encoding table retained in Golden Rules (defines WHAT to encode, not HOW)
- Classification Indicators, Fact Table Required Keys, Verification Method sections preserved as-is
- Generation Patterns preserved but framework-specific code examples removed from Golden Rules zone

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] Each file has Golden Rules section (framework-agnostic)
- [x] Each file has Reconnaissance Hints section with annotation
- [x] Golden Rules reference _shared.md principles
- [x] Shared antipattern guards removed, replaced by _shared.md reference
- [x] Type-specific antipattern guards retained in Golden Rules
- [x] Output path changed to tests/<journey>/
- [x] Step number references corrected to SKILL.md Step 1.3
- [x] CLI: Two-Level Timeout (test function + process-level SIGKILL guard)
- [x] CLI: Binary Isolation (dedicated binary, no go run or PATH)
- [x] CLI: Environment Hermeticity (explicit env inheritance + override)
- [x] TUI: Terminal Size Contract (TERM=dumb, fixed LINES/COLUMNS)
- [x] TUI: ANSI Sanitization (strip escape sequences before matching)
- [x] TUI: Stable State Detection (observable signals, not time-based)
- [x] API: Idempotency Check (repeated PUT/DELETE identical results)
- [x] API: Request Timeout (connection + read/write timeouts)
- [x] API: Content-Type Verification (Accept header + response Content-Type)

## Notes
无
