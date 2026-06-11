---
status: "completed"
started: "2026-06-10 21:09"
completed: "2026-06-10 21:11"
time_spent: "~2m"
---

# Task Record: 13 Fix extract-design-md cross-references and templates (MEDIUM-A4, MINOR-C3)

## Summary
Fixed extract-design-md cross-skill reference in SKILL.md by adding path resolution context for ui-design skill directory, and standardized all placeholders in 3 design templates (design-web.md, design-tui.md, design-mobile.md) to UPPER_CASE format

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/extract-design-md/SKILL.md
- plugins/forge/skills/extract-design-md/templates/design-web.md
- plugins/forge/skills/extract-design-md/templates/design-tui.md
- plugins/forge/skills/extract-design-md/templates/design-mobile.md

### Key Decisions
无

## Document Metrics
4 files modified, 1 cross-skill reference clarified, ~40 placeholders standardized across 3 templates

## Referenced Documents
- docs/conventions/forge-distribution.md

## Review Status
final

## Acceptance Criteria
- [x] SKILL.md cross-skill reference includes path resolution context explaining how to resolve ui-design skill directory in Forge distribution
- [x] All 3 design templates use UPPER_CASE placeholder format exclusively

## Notes
SKILL.md line 124 now explains: resolve path relative to skills parent directory. Templates converted mixed-case placeholders like {{App Name or Domain}}, {{YYYY-MM-DD}}, {{value}}, {{key}}, {{action}} to UPPER_CASE equivalents like {{APP_NAME_OR_DOMAIN}}, {{DATE}}, {{CARD_BACKGROUND}}, {{KEY}}, {{KEY_ACTION}}.
