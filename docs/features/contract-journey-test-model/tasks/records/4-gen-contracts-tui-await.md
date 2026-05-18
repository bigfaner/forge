---
status: "completed"
started: "2026-05-18 01:07"
completed: "2026-05-18 01:29"
time_spent: "~22m"
---

# Task Record: 4 gen-contracts skill + TUI await 形式化

## Summary
Implemented gen-contracts skill (SKILL.md + templates) and Contract validation logic in forge-cli/pkg/contract/. The skill generates Contract specifications with six dimensions (Preconditions/Input/Output/State/Side-effect/Invariants) from Journey documents + Fact Table code reconnaissance. TUI async Cmd await semantics are formalized with fail-fast timeout behavior and tea.Batch() completion semantics. Contract validation enforces: mandatory dimension completeness, semantic descriptor purity (no regex), Outcome Preconditions mutual exclusivity, Journey Invariants presence, Outcome count checkpoint (>5 triggers review), and TUI async timeout Outcome requirements.

## Changes

### Files Created
- plugins/forge/skills/gen-contracts/SKILL.md
- plugins/forge/skills/gen-contracts/templates/contract.md
- plugins/forge/skills/gen-contracts/templates/outcome-block.md
- forge-cli/pkg/contract/contract.go
- forge-cli/pkg/contract/validate.go
- forge-cli/pkg/contract/render.go
- forge-cli/pkg/contract/validate_test.go
- forge-cli/pkg/contract/render_test.go

### Files Modified
- forge-cli/scripts/version.txt

### Key Decisions
- Preconditions mutual exclusivity uses substring containment check: if one precondition string contains another, they are considered overlapping (subset relationship). This is a conservative heuristic that catches the most common case of overlapping preconditions.
- TUI await validation requires a corresponding timeout Outcome for every async TUI Outcome. The timeout Outcome must have TimedOutCmd field populated with the name of the Cmd that would time out.
- State verification level is determined automatically from Fact Table reconnaissance: full (state query interface exists), partial (state inferable from output), deferred (cannot infer).
- Batch splitting threshold: 15 Contracts. First batch = success Outcomes, second batch = edge case Outcomes.
- Outcome fields IsAsyncTUI, AwaitTimeout, TimedOutCmd added to Outcome struct for TUI await semantics tracking.

## Test Results
- **Tests Executed**: Yes
- **Passed**: 73
- **Failed**: 0
- **Coverage**: 95.1%

## Acceptance Criteria
- [x] gen-contracts derives Contracts from Journey + Fact Table with 4 mandatory dimensions per Outcome
- [x] All generated Contracts pass six-dimension completeness validation (mandatory non-empty, no regex in descriptors)
- [x] Each Journey generates Journey-level Invariants (at least 1 per Contract)
- [x] Multi-Outcome Contract correctness: Preconditions are mutually exclusive
- [x] TUI await semantics: fail-fast on timeout, report timed-out Cmd name, tea.Batch waits for all
- [x] Contract specs stored as structured markdown in tests/<journey>/_contracts/

## Notes
The gen-contracts SKILL.md is the most complex skill in the pipeline. It references code reconnaissance (Fact Table) for grounding semantic descriptors in real code. The validation logic is implemented as a separate Go package (forge-cli/pkg/contract/) that can be used both by the skill and by potential future CLI commands like 'forge test verify'. Version bumped to 3.21.0.
