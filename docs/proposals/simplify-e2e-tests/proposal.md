---
created: 2026-05-15
author: faner
status: Draft
---

# Proposal: Simplify Forge-CLI E2E Tests — Remove Plugin Text Verification

## Problem

The forge-cli e2e test suite contains **32 test cases** across 2 files that verify **static text content** of plugin/template files rather than testing **forge-cli code behavior**. These text-verification tests are brittle, high-maintenance, and do not test CLI functionality.

### Evidence

- `tui-ui-design/tui_ui_design_cli_test.go` — 31 test functions reading plugin template/SKILL.md files and asserting string presence (entire file is text-verification)
- `justfile-canonical-e2e/justfile_canonical_e2e_cli_test.go` — 1 text-verification test (`TC-020`) reading 6 static `manifest.yaml` files and asserting YAML key presence/absence; remaining 19 tests are CLI behavior tests
- Text tests break whenever plugin documentation is reworded, even when CLI behavior is unchanged

### Urgency

Current test suite has significant maintenance overhead — any SKILL.md wording change triggers e2e test failures unrelated to code regressions.

## Proposed Solution

Delete the entire `tui-ui-design/` directory (31 text-verification tests + local helpers). Remove the single text-verification test (`TC-020`) from the bridged `justfile-canonical-e2e/justfile_canonical_e2e_cli_test.go`, preserving its 19 CLI behavior tests.

### Innovation Highlights

Straightforward cleanup — no novel approach. The insight is separating "does the CLI work correctly" from "does the plugin documentation say the right things."

## Requirements Analysis

### Key Scenarios

- Delete `tui-ui-design/` directory and its sole test file
- Remove `TC-020` from `justfile-canonical-e2e/justfile_canonical_e2e_cli_test.go` (text-verification test reading static manifest.yaml files)
- Verify remaining test suite still compiles and passes

### Constraints & Dependencies

- Must not reduce coverage of forge-cli runtime behavior
- Test suite must still compile (`go test -tags=e2e ./tests/e2e/...`)

## Alternatives & Industry Benchmarking

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| Do nothing | No effort | Continued maintenance burden, false failures | Rejected |
| Move text tests to plugin-level test suite | Preserves coverage | Requires new test infrastructure | Deferred: out of scope |
| **Delete text-verification file** | Immediate cleanup, focused test suite | Loses plugin text coverage | **Selected** |

## Feasibility Assessment

### Technical Feasibility

Directory deletion + surgical removal of 1 test function from a bridged file. The `tui-ui-design/` package is self-contained — its helpers are local, not shared.

### Resource & Timeline

Trivial scope: 1 deletion + 1 edit + verification.

## Scope

### In Scope

- Delete `tests/e2e/tui-ui-design/` directory (contains `tui_ui_design_cli_test.go` — 31 text-verification tests)
- Remove `TestTC_020_AllManifestsContainZeroRunAndGraduateFields` from `tests/e2e/justfile-canonical-e2e/justfile_canonical_e2e_cli_test.go` (reads static manifest.yaml files — text-verification)
- Verify remaining suite compiles and passes

### Out of Scope

- Creating replacement plugin-level text test suite
- Changes to non-e2e tests

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Removing test coverage for plugin template changes | Low | Low | Template changes are caught by code review; runtime behavior is still tested |
| TC-020 implicitly verified CLI behavior | Low | Medium | TC-020 only reads static YAML files — no CLI invocation; profile schema changes are caught by `go test` in forge-cli package |

## Success Criteria

- [ ] `tests/e2e/tui-ui-design/` directory deleted
- [ ] `TestTC_020` removed from `justfile-canonical-e2e/justfile_canonical_e2e_cli_test.go`
- [ ] `go test -tags=e2e ./tests/e2e/...` compiles without errors
- [ ] Remaining CLI-focused tests pass

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
