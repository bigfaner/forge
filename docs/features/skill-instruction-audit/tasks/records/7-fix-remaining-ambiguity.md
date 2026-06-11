---
status: "completed"
started: "2026-05-28 23:08"
completed: "2026-05-28 23:11"
time_spent: "~3m"
---

# Task Record: 7 Fix ambiguity and logic issues in remaining skills/commands

## Summary
Fixed ambiguity and logic issues in 7 non-pipeline skills/commands: consolidate-specs (scan dirs + safety qualification), gen-sitemap (@latest/pinning clarity), clean-code (config exception + offline branch detection), deep-research (Q4 dimension selection + review wait step), ui-design (web+mobile output rule + eval-skip labels), simplify-skill (user vs plugin scope), forensic (project-hash discovery + skills parent directory resolution)

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/consolidate-specs/SKILL.md
- plugins/forge/skills/gen-sitemap/SKILL.md
- plugins/forge/skills/clean-code/SKILL.md
- plugins/forge/skills/deep-research/SKILL.md
- plugins/forge/skills/ui-design/SKILL.md
- plugins/forge/commands/simplify-skill.md
- plugins/forge/skills/forensic/SKILL.md

### Key Decisions
无

## Document Metrics
7 files modified, 7 AC items resolved, 0 contradictions remaining

## Referenced Documents
无

## Review Status
final

## Acceptance Criteria
- [x] consolidate-specs Steps 9-11 specify scan directories (docs/business-rules/, docs/conventions/)
- [x] gen-sitemap has no @latest vs pinning contradiction
- [x] clean-code has clear config file exception or none
- [x] deep-research has explicit wait for user review between report presentation and proposal conversion ask
- [x] ui-design has web+mobile output rule; eval-skip option is clearly labeled
- [x] simplify-skill states whether it targets user skills or plugin skills
- [x] forensic explains how to obtain project-hash path

## Notes
All 7 AC items addressed. clean-code also fixed offline default branch detection (network-free fallback before git remote show). forensic also clarified skills parent directory resolution for both plugin and user skill locations.
