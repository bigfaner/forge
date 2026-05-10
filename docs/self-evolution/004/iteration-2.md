---
date: "2026-05-10"
plugin_version: "2.18.0"
iteration: "2"
target_score: "950"
evaluator: Claude (structural audit)
---

# Forge Plugin Audit — Iteration 2

**Score: 955/1000** (target: 950)

```
┌─────────────────────────────────────────────────────────────────┐
│                  PLUGIN CONSISTENCY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Directory-Name Alignment  │  40      │  40      │ ✅         │
│    Skill name matches dir    │  25/25   │          │            │
│    Command name matches file │  15/15   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Agent Reference Integrity │  100     │  100     │ ✅         │
│    Referenced agents exist   │  70/70   │          │            │
│    No orphan agents          │  30/30   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Reference Integrity       │  75      │  80      │ ⚠️         │
│    Template refs valid       │  25/25   │          │            │
│    Cross-skill refs valid    │  30/30   │          │            │
│    No orphan templates       │  15/15   │          │            │
│    No cross-file duplication │  5/10    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Frontmatter Completeness  │  110     │  110     │ ✅         │
│    Skill frontmatter         │  45/45   │          │            │
│    Command frontmatter       │  35/35   │          │            │
│    Agent frontmatter         │  30/30   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Eval Template Convention  │  100     │  100     │ ✅         │
│    rubric.md exists          │  30/30   │          │            │
│    report.md exists          │  30/30   │          │            │
│    Rubric→report chain valid │  20/20   │          │            │
│    Rubric totals correct     │  20/20   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 6. Orchestrator Convention   │  40      │  40      │ ✅         │
│    Iron Laws present         │  25/25   │          │            │
│    Hard Gate present         │  15/15   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 7. Task CLI Alignment        │  230     │  240     │ ⚠️         │
│    Command existence         │  25/25   │          │            │
│    Flag correctness          │  25/25   │          │            │
│    Output field parsing      │  15/15   │          │            │
│    Status machine align      │  35/35   │          │            │
│    Claim scheduling align    │  35/35   │          │            │
│    Record validation align   │  25/35   │          │            │
│    Dynamic task add align    │  25/25   │          │            │
│    Schema-code alignment     │  20/20   │          │            │
│    All-completed hook align  │  10/10   │          │            │
│    Template existence        │  5/10    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 8. Hook Wiring Integrity     │  70      │  70      │ ✅         │
│    hooks.json valid JSON     │  10/10   │          │            │
│    Hook scripts exist        │  25/25   │          │            │
│    Hook CLI commands valid   │  15/15   │          │            │
│    Hook event names valid    │  20/20   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 9. Guide Coverage+Concise    │  70      │  70      │ ✅         │
│    Guide references valid    │  30/30   │          │            │
│    Core workflow skills doc  │  25/25   │          │            │
│    Conciseness / no redund.  │  15/15   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 10. Command Metadata         │  60      │  60      │ ✅         │
│    allowed_tools declared    │  35/35   │          │            │
│    argument-hints declared   │  25/25   │          │            │
├──────────────────────────────┼──────────┼──────────┼──────────┤
│ 11. Plugin Metadata          │  40      │  40      │ ✅         │
│    keywords coverage         │  25/25   │          │            │
│    description accurate      │  15/15   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 12. Safety Marker Consist.   │  50      │  50      │ ✅         │
│    Command/agent markers     │  30/30   │          │            │
│    Dispatch cmd coverage     │  20/20   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  955     │  1000    │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| # | Check | File | Issue | Penalty |
|---|-------|------|-------|---------|
| 1 | 7f-Record validation align | plugins/forge/skills/record-task/SKILL.md:146 | Skill says "Override any validation error with `--force`" but Go code (record.go:270-274) applies auto-downgrade (`completed + testsFailed > 0 → blocked`) unconditionally — `--force` is NOT checked before downgrade. The skill's claim that `--force` overrides "any" validation is inaccurate for auto-downgrade. | -10 |
| 2 | 3-No cross-file duplication | plugins/forge/commands/execute-task.md:59, plugins/forge/commands/fix-bug.md:176, plugins/forge/agents/task-executor.md:97-107, plugins/forge/agents/error-fixer.md:86-98 | Quality gate failure actions (compile→fix, fmt→blocked, lint→self-fix, test→fix) copy-pasted across 4 files. Task-executor and error-fixer are autonomous subagents (exception applies), but execute-task and fix-bug run in main session and could reference guide.md. | -5 |
| 3 | 7j-Template existence | task-cli (embedded binary) | The `fix-task` template is confirmed to exist via `task template fix-task` but no on-disk file exists in `plugins/forge/skills/*/templates/` — it is embedded in the Go binary. Partial alignment since templates should ideally be on disk for discoverability. | -5 |

---

## Attack Points

### Attack 1: [7f — --force does NOT override auto-downgrade]

**Where**: `plugins/forge/skills/record-task/SKILL.md:146`
**What's wrong**: The skill states "Override any validation error with `--force`" (line 146). However, Go code at `task-cli/internal/cmd/record.go:270-274` applies auto-downgrade (`completed + testsFailed > 0 → blocked`) **without checking the `force` flag**. The `force` flag is only checked later (line 281) for test evidence and acceptance criteria checks. Auto-downgrade is intentionally non-overridable, but the skill misleads by saying `--force` overrides "any" validation error. The skill should explicitly state that auto-downgrade for test failures is non-overridable, and `--force` only bypasses test evidence and AC checks.
**How to fix**: Update line 146 to: "Override test evidence and AC validation errors with `--force` (auto-downgrade for test failures is always enforced)." Update the table to distinguish overridable vs non-overridable validations.

### Attack 2: [3 — Quality gate failure actions duplicated across 4 files]

**Where**: `plugins/forge/commands/execute-task.md:59`, `plugins/forge/commands/fix-bug.md:176`, `plugins/forge/agents/task-executor.md:97-107`, `plugins/forge/agents/error-fixer.md:86-98`
**What's wrong**: The quality gate failure actions (compile→fix & retry, fmt→blocked, lint→self-fix then blocked, test→fix & retry) are copy-pasted across 4 files. The canonical location is guide.md's Quality Gate Protocol section. Task-executor and error-fixer are autonomous agents that cannot read other files at runtime (exception applies). But execute-task.md (line 59) and fix-bug.md (line 176) run in the main session and could reference guide.md instead of duplicating. Both currently use a single-line inline summary rather than a full table, which is better than iteration 1 but still duplicated.
**How to fix**: In execute-task.md and fix-bug.md, replace the inline quality gate failure action text with "Apply the Quality Gate Protocol from the Forge Guide." Keep the table in task-executor.md and error-fixer.md since they run as subagents.

### Attack 3: [7j — fix-task template exists only in binary, not on disk]

**Where**: `plugins/forge/skills/*/templates/` (missing), `task-cli/pkg/template/` (embedded)
**What's wrong**: The `fix-task` template is referenced by breakdown-tasks, quick-tasks, execute-task, run-tasks, and task-executor skills. While it exists and works via `task template fix-task`, no on-disk template file exists in the forge plugin's `skills/*/templates/` directory. Templates are embedded in the Go binary. This means on-disk template inventory is incomplete — agents inspecting the templates directory won't find it.
**How to fix**: Either (a) extract the fix-task template to `plugins/forge/skills/breakdown-tasks/templates/fix-task.md` for discoverability, or (b) document in guide.md that fix-task is embedded in the CLI binary (run `task template fix-task` to view).

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1: [7f — Record validation quality gate mismatch (4 vs 3 steps)] | ✅ | The iteration 1 audit incorrectly stated `DefaultGateSequence()` has 3 steps. Actually `just.go:21-28` shows 4 steps including test. The skill at line 130 now correctly shows 4 steps matching the code. This was a false positive in iteration 1. |
| Attack 2: [10 — simplify-skill missing Write/Edit in allowed_tools] | ✅ | `simplify-skill.md:5` now has `allowed_tools: ["Read", "Write", "Edit", "AskUserQuestion"]`. Write and Edit are included. |
| Attack 3: [7h — Schema-code scope required mismatch] | ✅ | `index.schema.json:51` now has `"required": ["id", "title", "priority", "status", "file"]` — `scope` is no longer required. Matches Go code where scope is `omitempty`. |
| Attack 4: [9 — Guide contains duplicated task record workflow detail] | ✅ | `guide.md:148-164` now has only a Key Commands table and a one-line reference: "For record workflow details, see the `/record-task` skill." The detailed workflow section was removed. |
| Attack 5: [3 — Quality gate failure table duplicated across 4 files] | ⚠️ Partial | execute-task.md and fix-bug.md reduced from full table to single-line inline summary. But the content is still duplicated text. Full resolution would replace with "Apply the Quality Gate Protocol from the Forge Guide." |
| Attack 6: [7i — All-completed hook guide missing auto-fix-task behavior] | ✅ | `guide.md:131` now says: "On failure at any step, a P0 fix-task is automatically created. Run `task claim` to pick it up." Matches `all_completed.go:128-129`. |
| Attack 7: [12 — error-fixer agent lacks ONE-TASK constraint equivalent] | ✅ | `error-fixer.md:167-181` now has a STOP section with `<HARD-RULE>` and `<PROHIBITIONS>` block matching task-executor's pattern. |

---

## Fix Summary

N/A (iteration 2 — no fixes applied, audit only)

---

## Verdict

- **Score**: 955/1000
- **Target**: 950/1000
- **Gap**: 0 points (target exceeded by 5)
- **Action**: Target reached. Remaining deductions are minor (misleading --force documentation, inline quality gate duplication, embedded template discoverability). No critical issues found.
