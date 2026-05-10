---
date: "2026-05-10"
plugin_version: "2.18.0"
iteration: "3"
target_score: "950"
evaluator: Claude (structural audit)
---

# Forge Plugin Audit — Iteration 3

**Score: 960/1000** (target: 950)

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
│ 3. Reference Integrity       │  80      │  80      │ ✅         │
│    Template refs valid       │  25/25   │          │            │
│    Cross-skill refs valid    │  30/30   │          │            │
│    No orphan templates       │  15/15   │          │            │
│    No cross-file duplication │  10/10   │          │            │
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
│ 7. Task CLI Alignment        │  235     │  240     │ ⚠️         │
│    Command existence         │  25/25   │          │            │
│    Flag correctness          │  25/25   │          │            │
│    Output field parsing      │  15/15   │          │            │
│    Status machine align      │  35/35   │          │            │
│    Claim scheduling align    │  35/35   │          │            │
│    Record validation align   │  30/35   │          │            │
│    Dynamic task add align    │  25/25   │          │            │
│    Schema-code alignment     │  15/20   │          │            │
│    All-completed hook align  │  10/10   │          │            │
│    Template existence        │  10/10   │          │            │
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
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 11. Plugin Metadata          │  40      │  40      │ ✅         │
│    keywords coverage         │  25/25   │          │            │
│    description accurate      │  15/15   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 12. Safety Marker Consist.   │  50      │  50      │ ✅         │
│    Command/agent markers     │  30/30   │          │            │
│    Dispatch cmd coverage     │  20/20   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  960     │  1000    │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| # | Check | File | Issue | Penalty |
|---|-------|------|-------|---------|
| 1 | 7f-Record validation align | plugins/forge/skills/record-task/SKILL.md:146 | Skill says "Override quality gate, test evidence, and AC validation with `--force`" but the quality gate (`validateQualityGate` in record.go:442-453) is checked BEFORE the `--force` flag is consulted (record.go:121-124 shows `!recordForce` condition on line 122, but `validateQualityGate` is called independently). Actually, re-reading the code: record.go:121-124 shows `if rd.Status == "completed" && !recordForce && !t.NoTest { validateQualityGate(...) }` — so `--force` DOES skip the quality gate. However, the skill's table at line 141 says auto-downgrade (`completed + testsFailed > 0 → blocked`) is "non-overridable" but the `--force` description at line 146 could be misread as overriding "any" validation. The note at line 151 clarifies this correctly. Net: the description is borderline but the explicit note resolves ambiguity. Minor inconsistency in wording only. | -5 |
| 2 | 7h-Schema-code alignment | plugins/forge/skills/quick-tasks/templates/index.schema.json:45 | Quick-tasks schema requires `"scope"` in each task's `required` array (`["id", "title", "priority", "status", "file", "scope"]`), but the Go struct `Task` (types.go:23) has `scope` with `omitempty`. The validate.go (line 115-133) does not enforce `scope` as required. This means a valid task index (per Go code) with `scope` missing will fail JSON schema validation. The breakdown-tasks schema was fixed in iteration 1 but the quick-tasks copy was not updated. | -5 |

---

## Attack Points

### Attack 1: [7f — Record validation `--force` description could be more precise]

**Where**: `plugins/forge/skills/record-task/SKILL.md:146`
**What's wrong**: Line 146 says "Override quality gate, test evidence, and AC validation with `--force`" which is technically correct — the Go code confirms `--force` bypasses `validateQualityGate` (record.go:122), test evidence check (record.go:281-284), and AC check (record.go:287-298). The auto-downgrade at line 271-273 is correctly documented as non-overridable (SKILL.md line 151). The wording is acceptable but the skill uses "Override any validation error" in the table header context which could confuse. The explicit note at line 151 ("auto-downgrade... is **never** overridden by `--force`") resolves this. This is a very minor wording issue.
**How to fix**: Change line 146 header from "Override quality gate, test evidence, and AC validation with `--force`" to "Override quality gate, test evidence, and AC validation errors with `--force` (auto-downgrade for test failures is always enforced):" to make it crystal clear.

### Attack 2: [7h — quick-tasks schema still requires scope when Go code does not]

**Where**: `plugins/forge/skills/quick-tasks/templates/index.schema.json:45`
**What's wrong**: The quick-tasks schema has `"required": ["id", "title", "priority", "status", "file", "scope"]` (line 45), meaning `scope` is mandatory. But the Go struct `Task.Scope` is `omitempty` (types.go:23), and `validate.go` does not check for missing `scope`. The breakdown-tasks schema was fixed in iteration 1 to remove `scope` from required, but the quick-tasks schema copy was not updated to match. A quick-mode task index with `scope` omitted would pass `task validate` (Go) but fail JSON schema validation.
**How to fix**: In `plugins/forge/skills/quick-tasks/templates/index.schema.json:45`, change `"required": ["id", "title", "priority", "status", "file", "scope"]` to `"required": ["id", "title", "priority", "status", "file"]` to match the breakdown-tasks schema and Go behavior.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1: [7f — Record validation quality gate mismatch (4 vs 3 steps)] | ✅ | Iteration 1 was a false positive. `DefaultGateSequence()` has 4 steps including test (just.go:21-28). SKILL.md line 130 correctly shows 4 steps. Confirmed correct in iteration 2 and remains correct. |
| Attack 2: [10 — simplify-skill missing Write/Edit in allowed_tools] | ✅ | `simplify-skill.md:5` has `allowed_tools: ["Read", "Write", "Edit", "AskUserQuestion"]`. Write and Edit present. |
| Attack 3: [7h — Schema-code scope required mismatch] | ⚠️ Partial | `breakdown-tasks/templates/index.schema.json:51` fixed to `["id", "title", "priority", "status", "file"]` (scope removed). But `quick-tasks/templates/index.schema.json:45` still has scope in required array. See Attack 2 above. |
| Attack 4: [9 — Guide contains duplicated task record workflow detail] | ✅ | `guide.md` has a Key Commands table and one-line reference to `/record-task` skill. No duplicated workflow detail. |
| Attack 5: [3 — Quality gate failure table duplicated across 4 files] | ✅ | execute-task.md and fix-bug.md now use single-line inline references ("See Forge Guide Quality Gate Protocol for failure actions"). task-executor.md and error-fixer.md retain the full table (autonomous agent exception applies). No cross-file duplication remaining for main-session files. |
| Attack 6: [7i — All-completed hook guide missing auto-fix-task behavior] | ✅ | `guide.md:133` says: "On failure at any step, a P0 fix-task is automatically created. Run `task claim` to pick it up." Matches `all_completed.go:128-129`. |
| Attack 7: [12 — error-fixer agent lacks ONE-TASK constraint equivalent] | ✅ | `error-fixer.md:167-181` has STOP section with `<HARD-RULE>` and `<PROHIBITIONS>` block matching task-executor's pattern. |
| Iter-2 Attack 1: [7f — --force does NOT override auto-downgrade] | ✅ | SKILL.md line 151 explicitly states: "auto-downgrade (`completed` + `testsFailed > 0` → `blocked`) is **never** overridden by `--force`." Accurate per Go code (record.go:270-273 applies before force check). |
| Iter-2 Attack 2: [3 — Quality gate failure actions duplicated] | ✅ | Main-session commands (execute-task, fix-bug) now use single-line references. Resolved. |
| Iter-2 Attack 3: [7j — fix-task template exists only in binary] | ✅ | `plugins/forge/skills/breakdown-tasks/templates/fix-task.md` now exists on disk with full template content (47 lines). Template includes E2E Fix Boundaries section. Matches the embedded template from `task template fix-task`. |

---

## Fix Summary

| File Changed | What Changed |
|-------------|--------------|
| N/A | No fixes applied this iteration — audit only |

---

## Verdict

- **Score**: 960/1000
- **Target**: 950/1000
- **Gap**: 0 points (target exceeded by 10)
- **Action**: Target reached. Remaining deductions are minor: (1) borderline `--force` wording in record-task SKILL.md, (2) quick-tasks JSON schema still requires `scope` when Go code does not. No critical or high-severity issues found.
