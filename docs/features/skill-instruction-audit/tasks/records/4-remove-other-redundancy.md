---
status: "completed"
started: "2026-05-28 22:58"
completed: "2026-05-28 23:00"
time_spent: "~2m"
---

# Task Record: 4 Remove frontmatter duplication and rule preview redundancy

## Summary
Removed frontmatter-body duplication in 4 commands and 4 skills: deleted opening sentences that repeated frontmatter descriptions, shortened gen-contracts Core principle, merged run-tests opening, deleted Mobile/TUI Overview subsections from extract-design-md (replaced with rules reference).

## Changes

### Files Created
无

### Files Modified
- plugins/forge/commands/clean-code.md
- plugins/forge/commands/git-commit.md
- plugins/forge/commands/git-checkout.md
- plugins/forge/commands/init-forge.md
- plugins/forge/skills/learn/SKILL.md
- plugins/forge/skills/gen-contracts/SKILL.md
- plugins/forge/skills/run-tests/SKILL.md
- plugins/forge/skills/extract-design-md/SKILL.md

### Key Decisions
无

## Document Metrics
8 files modified, ~30 lines removed, 0 lines added

## Referenced Documents
- docs/proposals/skill-instruction-audit/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] 4 commands (clean-code, git-commit, git-checkout, init-forge) body doesn't start with frontmatter sentence
- [x] learn/SKILL.md body doesn't duplicate frontmatter
- [x] gen-contracts/SKILL.md description <=1 sentence
- [x] run-tests/SKILL.md has single merged opening
- [x] extract-design-md/SKILL.md has no Overview paragraphs; only rules reference

## Notes
All 5 AC items pass. Only the 8 listed files were modified per Hard Rules.
