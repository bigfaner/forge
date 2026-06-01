---
status: "completed"
started: "2026-06-01 21:34"
completed: "2026-06-01 21:35"
time_spent: "~1m"
---

# Task Record: 4 Create /release-cli slash command

## Summary
Created /release-cli slash command that automates CLI release workflow: version bump, commit, tag, and push

## Changes

### Files Created
- .claude/commands/release-cli.md

### Files Modified
无

### Key Decisions
无

## Document Metrics
5 steps, ~20 lines, covers all 5 acceptance criteria

## Referenced Documents
- .claude/commands/upgrade-forge.md
- forge-cli/scripts/version.txt

## Review Status
final

## Acceptance Criteria
- [x] Command reads current version from forge-cli/scripts/version.txt
- [x] Prompts developer for new version number with semver bump suggestion
- [x] Updates forge-cli/scripts/version.txt with new version
- [x] Executes commit, tag (forge-cli/v{version}), and push
- [x] Documents GitHub Actions auto-trigger from tag push

## Notes
Followed upgrade-forge.md pattern for prompt structure. Tag format uses forge-cli/v prefix per Hard Rules. CLI version independent from Plugin version.
