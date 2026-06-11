---
status: "completed"
started: "2026-05-17 13:18"
completed: "2026-05-17 13:21"
time_spent: "~3m"
---

# Task Record: 3 Create types/ui.md instruction file (includes Step 2-3 Sitemap/Locators)

## Summary
Created types/ui.md instruction file for gen-test-scripts, extracting UI-specific logic from monolithic SKILL.md. Includes Reconnaissance Strategy, Fact Table Required Keys, Sitemap Resolution (Step 2 verbatim move), Locator Mapping (Step 3 verbatim move with priority chain preserved), Verification Method, Generation Patterns, Integration Test Scripts, and UI Antipattern Guards. All hard rules verified: locator priority chain (Fact Table > Sitemap > Semantic inference) preserved exactly, data-testid derivation rules require Fact Table sourcing, sitemap/locator content moved verbatim from SKILL.md.

## Changes

### Files Created
- plugins/forge/skills/gen-test-scripts/types/ui.md

### Files Modified
无

### Key Decisions
- Sitemap Resolution and Locator Mapping content moved verbatim from SKILL.md lines 292-327 as structural move, not paraphrase
- Integration Test Scripts section (SKILL.md lines 410-428) moved to ui.md as UI-specific logic
- UI-specific HARD-RULEs about source-code-first derivation moved into Locator Mapping section
- Antipattern guards extended beyond generic 6 with UI-specific patterns (CSS selectors, screenshot-only assertions, missing testid fallback rules)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 0
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] plugins/forge/skills/gen-test-scripts/types/ui.md exists
- [x] Frontmatter declares type: ui and conventions: [testing-ui.md, frontend.md]
- [x] Contains Reconnaissance Strategy section with UI-specific search patterns
- [x] Contains Fact Table Required Keys section listing minimum keys for UI type
- [x] Contains Sitemap Resolution section (Step 2 equivalent)
- [x] Contains Locator Mapping section (Step 3 equivalent) with priority chain
- [x] Contains Verification Method section
- [x] Contains Generation Patterns section
- [x] Contains Integration Test Scripts section
- [x] Contains UI Antipattern Guards section
- [x] At least 3 section headings unique to this file

## Notes
Documentation-only task (type: documentation, mainSession: false). 128 lines, 8 sections, 3 unique headings (Sitemap Resolution, Locator Mapping, Integration Test Scripts) not found in other type files.
