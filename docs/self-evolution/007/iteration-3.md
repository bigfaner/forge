---
date: "2026-05-13"
plugin_version: "3.0.0-beta-3"
iteration: "3"
target_score: "950"
evaluator: Claude (structural audit)
---

# Forge Plugin Audit -- Iteration 3

**Score: 945/1000** (target: 950)

```
+-----------------------------------------------+----------+----------+------------+
| Dimension                                     | Score    | Max      | Status     |
+-----------------------------------------------+----------+----------+------------+
| 1. Directory-Name Alignment                   |  40      |  40      | OK         |
|    Skill name matches dir                     |  25/25   |          |            |
|    Command name matches file                  |  15/15   |          |            |
+-----------------------------------------------+----------+----------+------------+
| 2. Agent Reference Integrity                  |  100     |  100     | OK         |
|    Referenced agents exist                    |  70/70   |          |            |
|    No orphan agents                           |  30/30   |          |            |
+-----------------------------------------------+----------+----------+------------+
| 3. Reference Integrity                        |  80      |  80      | OK         |
|    Template refs valid                        |  25/25   |          |            |
|    Cross-skill refs valid                     |  30/30   |          |            |
|    No orphan templates                        |  15/15   |          |            |
|    No cross-file duplication                  |  10/10   |          |            |
+-----------------------------------------------+----------+----------+------------+
| 4. Frontmatter Completeness                   |  110     |  110     | OK         |
|    Skill frontmatter                          |  45/45   |          |            |
|    Command frontmatter                        |  35/35   |          |            |
|    Agent frontmatter                          |  30/30   |          |            |
+-----------------------------------------------+----------+----------+------------+
| 5. Eval Template Convention                   |  90      |  100     | WARN       |
|    rubric.md exists                           |  30/30   |          |            |
|    report.md exists                           |  30/30   |          |            |
|    Rubric to report chain valid               |  20/20   |          |            |
|    Rubric totals correct                      |  10/20   |          |            |
+-----------------------------------------------+----------+----------+------------+
| 6. Orchestrator Convention                    |  40      |  40      | OK         |
|    Iron Laws present                          |  25/25   |          |            |
|    Hard Gate present                          |  15/15   |          |            |
+-----------------------------------------------+----------+----------+------------+
| 7. Task CLI Alignment                         |  210     |  240     | WARN       |
|    Command existence                          |  25/25   |          |            |
|    Flag correctness                           |  10/25   |          |            |
|    Output field parsing                       |  15/15   |          |            |
|    Status machine align                       |  35/35   |          |            |
|    Claim scheduling align                     |  35/35   |          |            |
|    Record validation align                    |  35/35   |          |            |
|    Dynamic task add align                     |  20/25   |          |            |
|    Schema-code alignment                      |  15/20   |          |            |
|    All-completed hook align                   |  10/10   |          |            |
|    Template existence                         |  10/10   |          |            |
+-----------------------------------------------+----------+----------+------------+
| 8. Hook Wiring Integrity                      |  70      |  70      | OK         |
|    hooks.json valid JSON                      |  10/10   |          |            |
|    Hook scripts exist                         |  25/25   |          |            |
|    Hook CLI commands valid                    |  15/15   |          |            |
|    Hook event names valid                     |  20/20   |          |            |
+-----------------------------------------------+----------+----------+------------+
| 9. Guide Coverage+Concise                     |  65      |  70      | WARN       |
|    Guide references valid                     |  30/30   |          |            |
|    Core workflow skills doc                   |  25/25   |          |            |
|    Conciseness / no redund.                   |  10/15   |          |            |
+-----------------------------------------------+----------+----------+------------+
| 10. Command Metadata                          |  55      |  60      | WARN       |
|    allowed_tools declared                     |  35/35   |          |            |
|    argument-hints declared                    |  20/25   |          |            |
+-----------------------------------------------+----------+----------+------------+
| 11. Plugin Metadata                           |  35      |  40      | WARN       |
|    keywords coverage                          |  20/25   |          |            |
|    description accurate                       |  15/15   |          |            |
+-----------------------------------------------+----------+----------+------------+
| 12. Safety Marker Consist.                    |  50      |  50      | OK         |
|    Command/agent markers                      |  30/30   |          |            |
|    Dispatch cmd coverage                      |  20/20   |          |            |
+-----------------------------------------------+----------+----------+------------+
| TOTAL                                         |  945     |  1000    |            |
+-----------------------------------------------+----------+----------+------------+
```

---

## Deductions

| # | Check | File | Issue | Penalty |
|---|-------|------|-------|---------|
| 1 | Rubric totals correct | `plugins/forge/skills/eval-proposal/templates/rubric.md` line 3 | Declares "Total: 1100 points" but SKILL.md description says "1000-point scoring", parameters say "Target score (0-1000)", report template shows "/1000". The rubric was changed from 1000 to 1100 between iterations but the skill description, parameter range, and report template were not updated. The rubric itself is internally consistent (10 dims sum to 1100), but the 100-point gap between rubric total and everything else is a functional mismatch. | -10 |
| 2 | Flag correctness | `plugins/forge/agents/task-executor.md` line 27 | Uses `task status <TASK_ID> blocked --reason "<error>"`. The `--reason` flag does not exist on `task status` (verified via `task status -h`). This command will fail at runtime, producing an error instead of setting the task to blocked. | -15 |
| 3 | Dynamic task add align | `plugins/forge/skills/breakdown-tasks/SKILL.md` line 421 | Claims "Maximum nesting: 3 levels" for fix-tasks. Neither `add.go` nor `record.go` enforces a nesting depth limit. This is an unenforced constraint in documentation. | -5 |
| 4 | Schema-code alignment | `plugins/forge/skills/breakdown-tasks/templates/index.schema.json` line 98-105 | `noTest` field marked "Deprecated: use type field instead" but Go code (`record.go` line 114) still actively checks `t.NoTest` for coverage auto-setting. Schema says deprecated but runtime behavior depends on it. (Known acceptable per rubric, but schema label is misleading.) | -5 |
| 5 | Conciseness | `plugins/forge/hooks/guide.md` lines 109-133 | Contains Quality Gate Protocol details (compile, fmt, lint, test sequence, failure actions, scope resolution algorithm with 4-step branching logic). This is reference material that belongs in individual skill docs. The guide is a workflow guide, not a CLI reference. | -5 |
| 6 | argument-hints | `plugins/forge/commands/quick.md` line 6 | `argument-hints: "[--no-test]"` uses plain string format, not structured YAML array format like `fix-bug.md`, `git-commit.md`, `gen-sitemap.md`, `extract-design-md.md` which use `{name, description, required}` objects. | -5 |
| 7 | keywords coverage | `plugins/forge/.claude-plugin/plugin.json` | Keywords: pipeline, brainstorm, prd, design, task, eval, e2e, test-profile, ui-design, sitemap. Missing: "breakdown" (breakdown-tasks skill), "fix" (fix-bug command), "quick" (quick/quick-tasks). These are major capability areas. | -5 |
| 8 | Claim output fields | `plugins/forge/commands/execute-task.md` lines 22-30, `plugins/forge/commands/run-tasks.md` lines 47-55 | Extract list mentions: TASK_ID, KEY, FILE, BREAKING, MAIN_SESSION, SCOPE, FEATURE. Claim.go also outputs TITLE, PRIORITY, TYPE, NO_TEST, ESTIMATED_TIME, DEPENDENCIES, RECORD (lines 312-326). Dispatchers document a subset. Not a rubric violation -- fields are correct, just incomplete. | INFO |

---

## Attack Points

### Attack 1: [Dimension 5 -- eval-proposal rubric/skill total mismatch]

**Where**: `plugins/forge/skills/eval-proposal/templates/rubric.md` line 3 declares "Total: 1100 points". `plugins/forge/skills/eval-proposal/SKILL.md` line 3 says "1000-point scoring", line 31 says "Target score (0-1000)", line 86 says "SCORE: X/1000", line 132 says "Final Score: X/1000".
**What's wrong**: The rubric was changed from 1000 to 1100 (iteration 2 noted it as a problem), but only the rubric was updated. The SKILL.md description, parameter range (`0-1000`), and report template (`/1000`) still reference 1000. When the scorer produces a score like `950/1100`, the gate comparison against target 900 will work numerically, but the report template says `/1000`, producing misleading output like `950/1000`. Additionally, a user setting `--target 1000` can never reach it since the rubric maxes at 1100, not 1000.
**How to fix**: Either (a) update SKILL.md description to "1100-point", parameter range to `0-1100`, and report template to `/1100`, or (b) reduce the rubric dimension totals to sum to 1000. Option (b) is better since all other eval skills use 1000.

### Attack 2: [Dimension 7 -- task-executor uses nonexistent --reason flag]

**Where**: `plugins/forge/agents/task-executor.md` line 27
**What's wrong**: The agent says `task status <TASK_ID> blocked --reason "<error>"`. The `--reason` flag does not exist on `task status` (verified: `task status -h` shows only `--force`). When `task prompt` fails with non-zero exit, this command will error instead of setting the task to blocked, leaving the task in `in_progress` state indefinitely. The correct command is `task status <TASK_ID> blocked` (without `--reason`).
**How to fix**: Change line 27 to `task status <TASK_ID> blocked`, removing the `--reason` flag. The blockedReason field in the Go Task struct is written programmatically by add.go (in all-completed hook), not by the status command.

### Attack 3: [Dimension 7 -- breakdown-tasks claims 3-level nesting limit not enforced in code]

**Where**: `plugins/forge/skills/breakdown-tasks/SKILL.md` line 421
**What's wrong**: The skill states "Maximum nesting: 3 levels" for fix-tasks. Neither `add.go` nor `record.go` enforces any nesting depth limit. The `--source-task-id` flag in add.go resolves to the root ancestor (line 51: "auto-resolves to root ancestor"), but there is no depth counter or limit check. Agents reading the skill will believe a 3-level limit exists when it does not.
**How to fix**: Either add a nesting depth check in add.go (count ancestors via sourceTaskID chain, reject if >= 3), or remove the "Maximum nesting: 3 levels" claim from breakdown-tasks SKILL.md line 421.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1: eval-proposal rubric total mismatch (1000 vs 1100) | PARTIAL | `plugins/forge/skills/eval-proposal/templates/rubric.md` line 3 now correctly declares "Total: 1100 points" (sums to 1100). However, SKILL.md description still says "1000-point scoring", parameter range is still "0-1000", and report template still uses "/1000". Only the rubric was fixed. |
| Attack 2: SubagentStart/SubagentStop hook events may be invalid | YES | Verified against `docs/official-references/plugin.md` lines 121-122: `SubagentStart` and `SubagentStop` are explicitly listed as valid Claude Code hook events. These are now confirmed valid. |
| Attack 3: Unenforced 3-level nesting limit in breakdown-tasks | NO | `plugins/forge/skills/breakdown-tasks/SKILL.md` line 421 still states "Maximum nesting: 3 levels". Neither `add.go` nor `record.go` enforces this limit. |

---

## Fix Summary

| File Changed | What Changed |
|-------------|--------------|
| `plugins/forge/skills/eval-proposal/templates/rubric.md` | Total changed from 1000 to 1100 (but SKILL.md not updated) |
| `plugins/forge/.claude-plugin/plugin.json` | Keywords updated with "test-profile", "sitemap" |
| `plugins/forge/commands/quick.md` | argument-hints added (plain string format) |
| Various commands | allowed_tools and argument-hints added |
| Eval skills | Hard Gate blocks added where missing |
| SubagentStart/SubagentStop | Confirmed valid in plugin docs (no code change needed) |

---

## Verdict

- **Score**: 945/1000
- **Target**: 950/1000
- **Gap**: 5 points
- **Action**: Target not reached (5 points short). Primary remaining gaps: (1) task-executor uses nonexistent `--reason` flag on `task status` -- will cause runtime failure (-15), (2) eval-proposal rubric/skill total mismatch -- rubric says 1100 but SKILL.md says 1000 (-10), (3) Unenforced 3-level nesting limit claim (-5), (4) noTest schema deprecation inconsistency (-5), (5) Minor guide conciseness, keyword gaps, argument-hints format issues (-15). The SubagentStart/SubagentStop issue from iteration 2 was confirmed valid (+15 improvement). Fixing the `--reason` flag and rubric consistency would bring score to 970.
