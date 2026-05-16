---
feature: "cli-list-reverse-chronological"
sources:
  - docs/proposals/cli-list-reverse-chronological/proposal.md
  - docs/features/cli-list-reverse-chronological/tasks/1-sort-proposal-list.md
  - docs/features/cli-list-reverse-chronological/tasks/2-sort-feature-list.md
generated: "2026-05-16"
---

# Test Cases: cli-list-reverse-chronological

## Summary

| Type | Count |
|------|-------|
| UI   | 0     |
| **Integration** | **0** |
| API  | 0     |
| CLI  | 6     |
| **Total** | **6** |

---

## CLI Test Cases

### TC-001: forge proposal lists proposals sorted by created date descending
- **Source**: Proposal Success Criterion "forge proposal lists proposals newest-first by created date" + Task 1 AC "runProposalList() sorts proposals by Created date descending (newest first)"
- **Type**: CLI
- **Target**: cli/proposal-list
- **Test ID**: cli/proposal-list/forge-proposal-lists-proposals-sorted-by-created-date-descending
- **Pre-conditions**: Project has 3+ proposals with distinct `created` frontmatter dates spanning different months
- **Steps**:
  1. Create a temporary project with `go.mod`
  2. Create proposals with `created` dates: `2026-01-15`, `2026-03-10`, `2026-02-01`
  3. Run `forge proposal`
  4. Parse output for slug order
- **Expected**: Proposals appear in order: `2026-03-10` (newest), `2026-02-01`, `2026-01-15` (oldest). Slugs in output respect this descending date order.
- **Priority**: P0

### TC-002: forge proposal handles proposals without created frontmatter (mtime fallback)
- **Source**: Task 1 AC "Proposals without created frontmatter still sort correctly (fallback mtime)" + Proposal Key Scenario "Proposals without created frontmatter fall back to file modification time"
- **Type**: CLI
- **Target**: cli/proposal-list
- **Test ID**: cli/proposal-list/forge-proposal-handles-proposals-without-created-frontmatter-mtime-fallback
- **Pre-conditions**: Project has proposals, some missing the `created` frontmatter field
- **Steps**:
  1. Create a temporary project with `go.mod`
  2. Create one proposal with `created: 2026-05-01` in frontmatter
  3. Create one proposal without `created` field (empty or missing)
  4. Run `forge proposal`
  5. Verify both proposals appear and command does not error
- **Expected**: Both proposals appear in output. Command completes without error. The proposal with missing `created` uses file mtime as fallback for sort positioning.
- **Priority**: P1

### TC-003: forge proposal with empty proposals directory
- **Source**: Proposal scope — implicit negative case from Task 1 AC "Existing tests continue to pass"
- **Type**: CLI
- **Target**: cli/proposal-list
- **Test ID**: cli/proposal-list/forge-proposal-with-empty-proposals-directory
- **Pre-conditions**: Project exists with no proposals
- **Steps**:
  1. Create a temporary project with `go.mod` but no proposals
  2. Run `forge proposal`
- **Expected**: Command outputs "no proposals found" to stderr without error exit.
- **Priority**: P1

### TC-004: forge feature list sorts features by manifest mtime descending
- **Source**: Proposal Success Criterion "forge feature list lists features newest-first by manifest mtime" + Task 2 AC "runFeatureList() sorts features by manifest.md mtime descending (newest first)"
- **Type**: CLI
- **Target**: cli/feature-list
- **Test ID**: cli/feature-list/forge-feature-list-sorts-features-by-manifest-mtime-descending
- **Pre-conditions**: Project has 3+ features with manifest.md files having distinct modification times
- **Steps**:
  1. Create a temporary project with `go.mod`
  2. Create features with manifest.mtimes set via `os.Chtimes`: `2026-01-01`, `2026-03-15`, `2026-05-16`
  3. Run `forge feature list`
  4. Parse output for slug order
- **Expected**: Features appear in order: `2026-05-16` (newest), `2026-03-15`, `2026-01-01` (oldest). Slugs in output respect this descending mtime order.
- **Priority**: P0

### TC-005: forge feature list sorts features with missing manifest to the end
- **Source**: Task 2 AC "Features with missing/unreadable manifest sort to the end"
- **Type**: CLI
- **Target**: cli/feature-list
- **Test ID**: cli/feature-list/forge-feature-list-sorts-features-with-missing-manifest-to-the-end
- **Pre-conditions**: Project has features, some without `manifest.md`
- **Steps**:
  1. Create a temporary project with `go.mod`
  2. Create one feature with a valid `manifest.md` (recent mtime)
  3. Create one feature directory without `manifest.md`
  4. Run `forge feature list`
  5. Compare output positions of both features
- **Expected**: Feature with manifest appears before feature without manifest. Features with missing/unreadable manifest (mtime=0) sort to the end of the list.
- **Priority**: P0

### TC-006: forge feature list with empty features directory
- **Source**: Proposal scope — implicit negative case from Task 2 AC "Existing tests continue to pass"
- **Type**: CLI
- **Target**: cli/feature-list
- **Test ID**: cli/feature-list/forge-feature-list-with-empty-features-directory
- **Pre-conditions**: Project exists with no features
- **Steps**:
  1. Create a temporary project with `go.mod` but no features
  2. Run `forge feature list`
- **Expected**: Command outputs "no features found" to stderr without error exit.
- **Priority**: P1

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Proposal Success Criterion + Task 1 AC-1 | CLI | cli/proposal-list | P0 |
| TC-002 | Task 1 AC-2 + Proposal Key Scenario | CLI | cli/proposal-list | P1 |
| TC-003 | Task 1 AC-3 (existing tests pass, empty case) | CLI | cli/proposal-list | P1 |
| TC-004 | Proposal Success Criterion + Task 2 AC-1 | CLI | cli/feature-list | P0 |
| TC-005 | Task 2 AC-2 | CLI | cli/feature-list | P0 |
| TC-006 | Task 2 AC-3 (existing tests pass, empty case) | CLI | cli/feature-list | P1 |
