---
status: "completed"
started: "2026-06-05 00:52"
completed: "2026-06-05 01:00"
time_spent: "~8m"
---

# Task Record: 3 Update 7 eval-* commands to read config and pass to skill

## Summary
Updated all 7 eval-* commands (proposal, prd, design, ui, journey, contract, consistency) to resolve target/iterations from config via forge config get before invoking eval skill. CLI arguments take priority over config values; when neither is set, args are omitted and eval skill uses rubric defaults.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/commands/eval-proposal.md
- plugins/forge/commands/eval-prd.md
- plugins/forge/commands/eval-design.md
- plugins/forge/commands/eval-ui.md
- plugins/forge/commands/eval-journey.md
- plugins/forge/commands/eval-contract.md
- plugins/forge/commands/eval-consistency.md

### Key Decisions
- Config resolution happens in command layer (not skill), keeping eval skill config-agnostic
- CLI args checked via grep on $ARGUMENTS before falling back to config
- Uniform bash template across all 7 commands, differing only in eval type name

## Test Results
- **Tests Executed**: Yes
- **Passed**: 25
- **Failed**: 0
- **Coverage**: 85.5%

## Acceptance Criteria
- [x] Each eval-* command reads eval.<type>.target and eval.<type>.iterations via forge config get
- [x] When config not set, --target/--iterations arg is omitted (eval skill uses rubric default)
- [x] CLI --target/--iterations arguments take priority over config values
- [x] All 7 eval-* commands follow the same config resolution pattern
- [x] Existing auto.eval.* auto-run behavior unchanged

## Notes
No Go code changes required -- task 1 already implemented the EvalSettings struct and reflection routing. This task only updated the 7 markdown command files.
