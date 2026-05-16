---
status: "completed"
started: "2026-05-16 21:19"
completed: "2026-05-16 21:28"
time_spent: "~9m"
---

# Task Record: 5 Migration logic and proposal intent field

## Summary
Two changes: (1) Changed migration fallback in migrate.go from TypeImplementation to TypeFeature as the conservative default. (2) Added Intent field to proposal Metadata struct and Proposal struct, wired through Discover() so proposals can declare their dominant intent.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/migrate.go
- forge-cli/pkg/proposal/proposal.go
- forge-cli/internal/cmd/migrate_test.go
- forge-cli/pkg/proposal/proposal_test.go
- forge-cli/scripts/version.txt

### Key Decisions
- TypeFeature is the conservative migration default because it is the broadest business type -- any implementation task could plausibly be a feature
- Intent field on proposal Metadata uses yaml:"intent" tag; empty/missing intent does not break existing proposals due to Go zero-value behavior
- Version bumped to 3.16.0 (minor) since this adds a new field to a public struct

## Test Results
- **Tests Executed**: Yes
- **Passed**: 9
- **Failed**: 0
- **Coverage**: 90.2%

## Acceptance Criteria
- [x] migrate.go maps TypeImplementation to TypeFeature as the conservative default fallback
- [x] forge task migrate on an old index.json with type: implementation produces type: feature
- [x] proposal.Metadata struct has Intent string yaml:"intent" field
- [x] parseFrontmatter() correctly parses intent from proposal frontmatter
- [x] Empty/missing intent field does not break existing proposals

## Notes
Migration is idempotent as verified by TestRunMigrate_Idempotent. All existing migrate tests updated to expect TypeFeature instead of TypeImplementation for unknown IDs.
