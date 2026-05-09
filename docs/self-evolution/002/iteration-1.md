---
date: "2026-05-09"
plugin_version: "2.16.1"
iteration: "1"
target_score: "800"
evaluator: Claude (structural audit)
---

# Forge Plugin Audit -- Iteration 1

**Score: 910/1000** (target: 800)

```
+----------------------------------+----------+----------+------------+
| Dimension                        | Score    | Max      | Status     |
+----------------------------------+----------+----------+------------+
| 1. Directory-Name Alignment      |  40      |  40      | OK         |
|    Skill name matches dir        |  25/25   |          |            |
|    Command name matches file     |  15/15   |          |            |
+----------------------------------+----------+----------+------------+
| 2. Agent Reference Integrity     |  100     |  100     | OK         |
|    Referenced agents exist       |  70/70   |          |            |
|    No orphan agents              |  30/30   |          |            |
+----------------------------------+----------+----------+------------+
| 3. Reference Integrity           |  80      |  80      | OK         |
|    Template refs valid           |  25/25   |          |            |
|    Cross-skill refs valid        |  30/30   |          |            |
|    No orphan templates           |  25/25   |          |            |
+----------------------------------+----------+----------+------------+
| 4. Frontmatter Completeness      |  110     |  110     | OK         |
|    Skill frontmatter             |  45/45   |          |            |
|    Command frontmatter           |  35/35   |          |            |
|    Agent frontmatter             |  30/30   |          |            |
+----------------------------------+----------+----------+------------+
| 5. Eval Template Convention      |  100     |  100     | OK         |
|    rubric.md exists              |  30/30   |          |            |
|    report.md exists              |  30/30   |          |            |
|    Rubric->report chain valid    |  20/20   |          |            |
|    Rubric totals correct         |  20/20   |          |            |
+----------------------------------+----------+----------+------------+
| 6. Orchestrator Convention       |  40      |  40      | OK         |
|    Iron Laws present             |  25/25   |          |            |
|    Hard Gate present             |  15/15   |          |            |
+----------------------------------+----------+----------+------------+
| 7. Task CLI Alignment            |  205     |  240     | WARN       |
|    Command existence             |  25/25   |          |            |
|    Flag correctness              |  25/25   |          |            |
|    Output field parsing          |  15/15   |          |            |
|    Status machine align          |  35/35   |          |            |
|    Claim scheduling align        |  35/35   |          |            |
|    Record validation align       |  35/35   |          |            |
|    Dynamic task add align        |  15/25   |          |            |
|    Schema-code alignment         |  5/20    |          |            |
|    All-completed hook align      |  5/10    |          |            |
|    Template existence            |  10/10   |          |            |
+----------------------------------+----------+----------+------------+
| 8. Hook Wiring Integrity         |  70      |  70      | OK         |
|    hooks.json valid JSON         |  10/10   |          |            |
|    Hook scripts exist            |  25/25   |          |            |
|    Hook CLI commands valid       |  15/15   |          |            |
|    Hook event names valid        |  20/20   |          |            |
+----------------------------------+----------+----------+------------+
| 9. Guide Coverage                |  30      |  70      | FAIL       |
|    Guide references valid        |  30/30   |          |            |
|    Core skills documented        |  0/40    |          |            |
+----------------------------------+----------+----------+------------+
| 10. Command Metadata             |  60      |  60      | OK         |
|    allowed_tools declared        |  35/35   |          |            |
|    argument-hints declared       |  25/25   |          |            |
+----------------------------------+----------+----------+------------+
| 11. Plugin Metadata              |  40      |  40      | OK         |
|    keywords coverage             |  25/25   |          |            |
|    description accurate          |  15/15   |          |            |
+----------------------------------+----------+----------+------------+
| 12. Safety Marker Consist.       |  35      |  50      | WARN       |
|    Command/agent markers         |  30/30   |          |            |
|    Dispatch cmd coverage         |  5/20    |          |            |
+----------------------------------+----------+----------+------------+
| TOTAL                            |  910     |  1000    |            |
+----------------------------------+----------+----------+------------+
```

---

## Deductions

| # | Check | File | Issue | Penalty |
|---|-------|------|-------|---------|
| 1 | 7g. Dynamic task addition | `skills/quick-tasks/SKILL.md:96-108` | The fix-task reference in quick-tasks uses `task add --template fix-task` but does NOT include required `--var` template variables (SOURCE_FILES, TEST_SCRIPT, TEST_RESULTS). The breakdown-tasks SKILL.md (line 344-353) and run-tasks command properly include these flags. | -5 |
| 2 | 7g. Dynamic task addition | `skills/quick-tasks/SKILL.md:96-108` | Quick-tasks fix-task pattern does not mention running `task template fix-task` first to view the template and required variables. Breakdown-tasks SKILL.md (line 344) does include this prerequisite step. | -5 |
| 3 | 7h. Schema-code alignment | `skills/breakdown-tasks/templates/index.schema.json:6` | Schema declares `"required": ["feature", "prd", "design", ...]` but Go types.go allows `Proposal` as an alternative to PRD+Design (quick mode). Quick mode index.json files would fail schema validation since `prd` and `design` are required but absent. | -5 |
| 4 | 7h. Schema-code alignment | `skills/breakdown-tasks/templates/index.schema.json` | Schema has no `proposal` field, but Go struct `TaskIndex` has `Proposal string` with `json:"proposal,omitempty"`. Schema is missing this field entirely. | -5 |
| 5 | 7h. Schema-code alignment | `skills/breakdown-tasks/templates/index.schema.json` (task properties) | Schema marks `scope` as required in task objects, but Go struct has `Scope string` with `json:"scope,omitempty"` and comment "Omitempty allows existing index.json files without scope to remain valid." Schema and code conflict. | -5 |
| 6 | 7i. All-completed hook align | `hooks/guide.md:127-129` | Guide says all-completed runs "just e2e-setup -> just probe -> just test-e2e" but actual all_completed.go uses programmatic `e2eprobe.ProbeServers()` (Go function), not `just probe`. The probe step is not a just recipe -- it is embedded in the Go binary. | -5 |
| 7 | 9. Core skills documented | `hooks/guide.md` | 9 skills/commands completely undocumented: `/simplify-skill`, `/git-commit`, `/git-checkout`, `/init-forge`, `/init-justfile`, `/record-decision`, `/extract-design-md`, `/forensic`, `/improve-harness`. At -5 each, total deduction = 45, capped to criterion max of 40. | -40 |
| 8 | 12. Dispatch cmd coverage | `commands/fix-bug.md` | fix-bug is explicitly listed in the rubric as a dispatch command requiring `<EXTREMELY-IMPORTANT>` ("execute-task, fix-bug, run-tasks, quick"), but the file has only `<HARD-GATE>` blocks and no `<EXTREMELY-IMPORTANT>` safety marker. | -15 |

---

## Attack Points

### Attack 1: [Dimension 9 -- guide.md missing 9 utility skills/commands]

**Where**: `plugins/forge/hooks/guide.md` (entire file)
**What's wrong**: The guide documents the pipeline workflow comprehensively but omits all utility/support skills and commands: `/simplify-skill`, `/git-commit`, `/git-checkout`, `/init-forge`, `/init-justfile`, `/record-decision`, `/extract-design-md`, `/forensic`, `/improve-harness`. Agents operating from the guide have no way to discover these tools.
**How to fix**: Add a "Utility Commands & Skills" section listing each with one-line description and trigger. Example: `/git-commit` -- "Format commits with Conventional Commits convention", `/forensic` -- "Analyze past session transcripts for agent deviation", etc.

### Attack 2: [Dimension 7h -- index.schema.json stale vs Go code]

**Where**: `plugins/forge/skills/breakdown-tasks/templates/index.schema.json:6` (required fields) and missing proposal field
**What's wrong**: Three schema-code mismatches: (1) `prd` and `design` are required in schema but Go allows `proposal` as alternative (quick mode), (2) `proposal` field is entirely absent from schema, (3) `scope` is required in schema task objects but Go uses `omitempty`. Quick-mode generated index.json files fail schema validation.
**How to fix**: Add `"proposal"` as optional property. Remove `prd` and `design` from the top-level `required` array (or use `oneOf` patterns). Move `scope` out of task-level `required` array. This is noted as a Known Acceptable Discrepancy for `sourceTaskID` but the `proposal` and `scope` issues are NOT listed as acceptable.

### Attack 3: [Dimension 7g -- quick-tasks incomplete fix-task reference]

**Where**: `plugins/forge/skills/quick-tasks/SKILL.md:96-108`
**What's wrong**: The fix-task code block in quick-tasks is missing: (1) the `task template fix-task` prerequisite step, (2) the required `--var SOURCE_FILES`, `--var TEST_SCRIPT`, `--var TEST_RESULTS` flags. Breakdown-tasks (line 344-353) and run-tasks both include these correctly. Quick-tasks agents will create malformed fix tasks.
**How to fix**: Replace lines 96-108 with the complete pattern from breakdown-tasks SKILL.md, including `task template fix-task` and all `--var` flags.

### Attack 4: [Dimension 12 -- fix-bug missing EXTREMELY-IMPORTANT safety marker]

**Where**: `plugins/forge/commands/fix-bug.md` (top of file, after frontmatter)
**What's wrong**: The rubric explicitly lists fix-bug as a dispatch command requiring `<EXTREMELY-IMPORTANT>` ("Commands that dispatch subagents (execute-task, fix-bug, run-tasks, quick) have `<EXTREMELY-IMPORTANT>` blocks with safety constraints"). fix-bug has `<HARD-GATE>` blocks but no `<EXTREMELY-IMPORTANT>`. While fix-bug does not technically dispatch subagents, the rubric considers it in scope for this check.
**How to fix**: Add `<EXTREMELY-IMPORTANT>` block at the top of the workflow: "1. NEVER fix without a failing test first 2. Minimal fix only -- no refactoring or scope creep 3. Fix + tests in ONE atomic commit".

### Attack 5: [Dimension 7i -- guide.md all-completed hook description inaccurate]

**Where**: `plugins/forge/hooks/guide.md:129`
**What's wrong**: Guide says e2e regression step runs "just e2e-setup -> just probe -> just test-e2e" but the actual code (all_completed.go:168) uses `e2eprobe.ProbeServers()`, a Go function that reads config.yaml and probes programmatically. There is no `just probe` recipe -- the probe is embedded in the Go binary. The guide description could mislead developers into looking for a nonexistent just recipe.
**How to fix**: Change "just e2e-setup -> just probe -> just test-e2e" to "just e2e-setup -> server health probe (built-in) -> just test-e2e".

---

## Previous Issues Check

N/A (iteration 1)

---

## Fix Summary

N/A (iteration 1)

---

## Verdict

- **Score**: 910/1000
- **Target**: 800/1000
- **Gap**: target exceeded by 110 points
- **Action**: Target reached. Remaining issues are non-blocking: guide.md utility section (cosmetic but improves discoverability), schema freshness (functional for full pipeline, broken for quick mode), fix-bug safety marker (defense-in-depth), quick-tasks fix-task reference (only affects error recovery path).
