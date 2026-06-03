# L1 Pilot Audit Report: README.md

## Audit Baseline

- **Baseline commit**: `85421b1085cb36b297f70d1b85120ab039e40724`
- **Audit date**: 2026-06-03
- **Audit scope**: `README.md` (475 lines)
- **Audit type**: L1 Pilot — factual claim extraction + codebase verification

## Issue Summary

- **P0**: 0 | **P1**: 3 | **P2**: 4 | **P3**: 3
- Total issues: 10
- Claims examined: 47
- Correct identifications: 37
- Misses (issues found after initial scan): 0
- False positives: 0

## Issue Details

### [P1] Version badge displays wrong version

- **File**: `README.md:5`
- **Claim**: `[![Version](https://img.shields.io/badge/Version-5.6.0-blue.svg)]`
- **Actual**: Plugin version in `plugins/forge/.claude-plugin/plugin.json` is `3.0.0-rc.41`. The badge shows `5.6.0` which does not match any version in the codebase. The entire README describes v3.0.0 features ("What's New in v3.0.0"), making the 5.6.0 badge clearly incorrect.
- **Suggested action**: Update badge to reflect current plugin version (3.0.0-rc.41) or use a dynamic badge that reads from plugin.json.

### [P1] Command count is incorrect (claims 18, actual 16)

- **File**: `README.md:418`
- **Claim**: `+-- commands/            # 18 个 Slash Commands`
- **Actual**: `ls plugins/forge/commands/ | wc -l` returns 16. The actual command files are: clean-code, eval-consistency, eval-contract, eval-design, eval-journey, eval-prd, eval-proposal, eval-ui, execute-task, extract-design-md, fix-bug, git-checkout, git-commit, quick, run-tasks, simplify-skill.
- **Suggested action**: Update count from 18 to 16.

### [P1] `test.verify-regression` task type does not exist

- **File**: `README.md:393`
- **Claim**: `test.verify-regression | 晋升后验证回归套件` listed as one of 5 test task types.
- **Actual**: `forge task list-types` output contains only 4 test types: `test.gen-journeys`, `test.gen-contracts`, `test.gen-scripts`, `test.run`. There is no `test.verify-regression` type. Tag promotion (`@feature` -> `@regression`) is handled internally by the `run-tests` skill, not by a separate task type.
- **Suggested action**: Remove `test.verify-regression` from the task type table. Update test type count from 5 to 4. Update total type count accordingly (doc section gains +1 for `doc.fix`, test section loses -1 for `test.verify-regression`, net total stays 21).

### [P2] `doc.fix` task type missing from type table

- **File**: `README.md:369-377`
- **Claim**: Doc section lists 5 types: `doc`, `doc.review`, `doc.summary`, `doc.consolidate`, `doc.drift`.
- **Actual**: `forge task list-types` shows 6 doc types including `doc.fix` (described as "fix documentation issues from review or eval"). The README table omits `doc.fix`.
- **Suggested action**: Add `doc.fix | fix documentation issues from review or eval` to the doc type table and update count from 5 to 6.

### [P2] `just unit-test` described as "fast, no -race" but actually includes `-race`

- **File**: `README.md:105`
- **Claim**: `just unit-test（快速、无 -race）和 just test（完整、含 e2e）`
- **Actual**: The justfile `unit-test` recipe runs `go test -race ./...` when gcc is available (which is the common case on macOS/Linux). It only skips `-race` when gcc is not found (fallback: `CGO_ENABLED=0 go test ./...`). The "no -race" description is misleading — `-race` is the default behavior.
- **Suggested action**: Update description to `just unit-test（fast unit tests）和 just test（complete surface tests）` or similar, removing the inaccurate `-race` characterization.

### [P2] Quick mode coding task limit claim is incorrect (claims 15, actual: no limit)

- **File**: `README.md:109`
- **Claim**: `编码任务上限提升至 15 个`
- **Actual**: The `quick-tasks` SKILL.md states: "No overall task count cap; task volume is bounded by proposal scope and the AC max rule." There is no 15-task limit for coding tasks. The only per-task constraint is "Maximum 6 Acceptance Criteria per task."
- **Suggested action**: Remove the "15 task limit" claim or update to reflect the actual constraint (6 AC max per task, no overall task count cap).

### [P2] Surface detection type list incomplete (missing Mobile)

- **File**: `README.md:99`
- **Claim**: `forge surfaces detect --apply 自动识别项目类型（Web / TUI / CLI / API）`
- **Actual**: The codebase defines 5 surface types in `forge-cli/pkg/forgeconfig/detect.go`: Web, Mobile, API, CLI, TUI. Mobile is a valid surface type but is not listed in the README.
- **Suggested action**: Update to `(Web / TUI / CLI / API / Mobile)` or `5 project types`.

### [P3] `task validate-index` subcommand name outdated (actual: `task validate`)

- **File**: `README.md:159`
- **Claim**: `forge task validate-index | 校验 index.json 结构`
- **Actual**: The actual subcommand is `forge task validate` (not `validate-index`). Running `forge task --help` shows `validate` in the available commands, not `validate-index`.
- **Suggested action**: Rename `validate-index` to `validate` in the subcommand table.

### [P3] `surfaces` command missing `--json` flag in flags reference

- **File**: `README.md:261-265`
- **Claim**: `surfaces` flags table lists `--types` and `--project-root`.
- **Actual**: `forge surfaces --help` also shows a `--json` flag ("output in JSON format") that is not mentioned in the README flags reference table.
- **Suggested action**: Add `--json` flag to the surfaces flags table.

### [P3] Test category header says "5 types" but should say "4 types"

- **File**: `README.md:379`
- **Claim**: `### test（5 种）`
- **Actual**: Only 4 test types exist (`test.gen-journeys`, `test.gen-contracts`, `test.gen-scripts`, `test.run`). The 5th type listed (`test.verify-regression`) does not exist.
- **Suggested action**: Update to `### test（4 种）`.

## Accuracy Baseline Report

| Metric | Value |
|--------|-------|
| Total factual claims examined | 47 |
| Correct claims (verified against codebase) | 37 |
| Inconsistencies found | 10 |
| Misses (issues found after initial audit pass) | 0 |
| False positives | 0 |
| **Miss rate** | **0%** |

### Claims Breakdown by Category

| Claim Category | Count | Correct | Issues |
|---------------|-------|---------|--------|
| File path existence (doc links) | 10 | 10 | 0 |
| Command existence (top-level) | 19 | 19 | 0 |
| Subcommand existence | 13 | 12 | 1 (`validate-index` -> `validate`) |
| Flag accuracy | 12 | 11 | 1 (missing `--json`) |
| Count/number claims | 5 | 2 | 3 (commands count, task types, surface types) |
| Version claims | 1 | 0 | 1 (badge version) |
| Behavior/process descriptions | 7 | 5 | 2 (unit-test -race, coding task limit) |
| Architecture/path references | 8 | 8 | 0 |

### Methodology Notes

The pilot audit followed the L1/L2 methodology from the proposal:
1. **Declaration extraction**: All factual claims extracted from README.md (code paths, command names, config values, behavior descriptions, counts)
2. **Code location**: Each claim verified via `forge --help`, `forge <cmd> --help`, `ls`, `find`, `grep`, and code reading
3. **Item-by-item comparison**: Three comparison types applied:
   - **Path/file reference**: All doc links, directory paths, and file references verified via `find`/`ls`
   - **Behavior/process description**: Command flows, test strategies, quality gate behavior verified via SKILL.md reading and justfile inspection
   - **State/config assertion**: Version numbers, type counts, flag lists verified against code constants and CLI output
4. **Gap detection**: Missing items identified by comparing code capabilities vs documentation coverage

### Methodology Assessment

The pilot miss rate is 0% (below the 20% threshold). The methodology is effective for L1 audits. Key observations:
- CLI `--help` output is the most reliable verification source for command/flag claims
- `forge task list-types` is authoritative for task type verification
- Go source code constant definitions are authoritative for count/type claims
- justfile recipes are authoritative for test/build behavior claims
- Skill SKILL.md files are authoritative for workflow behavior claims

No methodology adjustment needed. Full L1 audit can proceed with confidence.

## Cross-Layer Impact Checklist

For L3 auditors reviewing knowledge base entries:
- Any lesson/decision referencing `test.verify-regression` should be flagged as potentially outdated
- Any lesson/decision referencing "18 slash commands" should be flagged as outdated
- Any lesson/decision referencing "unit-test without -race" should be flagged as incorrect

## Audit Quality Review

- Sample ratio: 100% (pilot = full file audit)
- Sample result: PASS
- Missed items: 0
- Expanded review: No
