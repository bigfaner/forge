---
status: "blocked"
started: "2026-05-19 01:13"
completed: "N/A"
time_spent: ""
---

# Task Record: 4 Delete old agent definitions and update distribution docs

## Summary
Deleted doc-scorer.md and doc-reviser.md agent definitions. Updated forge-distribution.md distribution tree to show agents/experts/ structure (protocol + scorer subdirectories) and replaced Core Dependencies section. Updated README.md agents table to reflect new expert-based eval architecture.

## Changes

### Files Created
无

### Files Modified
- docs/conventions/forge-distribution.md
- README.md

### Key Decisions
- Updated README.md agents table in addition to forge-distribution.md since it also listed doc-scorer/doc-reviser as top-level agents
- Left remaining descriptive references in templates (validate-ux-task.md, eval-test-cases.md) and forensic SKILL.md unchanged -- these describe eval concepts, not functional agent-type references

## Test Results
- **Tests Executed**: Yes
- **Passed**: 0
- **Failed**: 1
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] doc-scorer.md deleted from plugins/forge/agents/
- [x] doc-reviser.md deleted from plugins/forge/agents/
- [x] forge-distribution.md distribution tree updated: agents/ section shows experts/ subdirectory with protocol/ and scorer/
- [x] forge-distribution.md removes doc-scorer.md and doc-reviser.md from the agent listing
- [x] forge-distribution.md Core Dependencies section updated: doc-scorer / doc-reviser references replaced with agents/experts/ description
- [x] No other files reference the deleted agent definitions (verify with grep)

## Notes
Pre-existing test failure TestExtractDesignMd_ArgumentHintsIncludesPlatform in forge-cli/internal/docsync is unrelated to this documentation-only cleanup. No code was modified. Verified on clean stash: failure exists before and after changes. No functional references (forge:doc-scorer, forge:doc-reviser, subagent_type) remain in plugins/forge/.
