---
status: "completed"
started: "2026-05-19 01:26"
completed: "2026-05-19 01:27"
time_spent: "~1m"
---

# Task Record: 4 Inline config.yaml and sitemap.json examples into gen-sitemap command

## Summary
Inlined config.yaml template and sitemap.json example directly into gen-sitemap.md, replacing two external references to references/shared/ directory files with fenced code blocks (yaml and json language tags respectively).

## Changes

### Files Created
无

### Files Modified
- plugins/forge/commands/gen-sitemap.md

### Key Decisions
- Replaced 'copy from template' instruction with inline YAML code block so the template content is self-contained
- Replaced 'See ... for a full example' with inline JSON code block preserving the complete sitemap.json content

## Test Results
- **Tests Executed**: Yes
- **Passed**: 0
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] No occurrence of references/shared/config.yaml in gen-sitemap.md
- [x] No occurrence of references/shared/sitemap.json in gen-sitemap.md
- [x] Config template and sitemap example are fully inline as code blocks

## Notes
Also verified no remaining ${CLAUDE_SKILL_DIR} references in the file. Both code blocks use correct language tags (yaml, json) per Hard Rules.
