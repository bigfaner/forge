---
date: "2026-05-10"
plugin_version: "2.18.0"
iteration: "4"
target_score: "950"
evaluator: Claude (structural audit)
---

# Forge Plugin Audit — Iteration 4

**Score: 965/1000** (target: 950)

```
┌─────────────────────────────────────────────────────────────────┐
│                  PLUGIN CONSISTENCY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 1. Directory-Name Alignment  │  40      │  40      │ ✅         │
│    Skill name matches dir    │  25/25   │          │            │
│    Command name matches file │  15/15   │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 2. Agent Reference Integrity │  100     │  100     │ ✅         │
│    Referenced agents exist   │  70/70   │          │            │
│    No orphan agents          │  30/30   │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 3. Reference Integrity       │  80      │  80      │ ✅         │
│    Template refs valid       │  25/25   │          │            │
│    Cross-skill refs valid    │  30/30   │          │            │
│    No orphan templates       │  15/15   │          │            │
│    No cross-file duplication │  10/10   │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 4. Frontmatter Completeness  │  110     │  110     │ ✅         │
│    Skill frontmatter         │  45/45   │          │            │
│    Command frontmatter       │  35/35   │          │            │
│    Agent frontmatter         │  30/30   │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 5. Eval Template Convention  │  100     │  100     │ ✅         │
│    rubric.md exists          │  30/30   │          │            │
│    report.md exists          │  30/30   │          │            │
│    Rubric→report chain valid │  20/20   │          │            │
│    Rubric totals correct     │  20/20   │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 6. Orchestrator Convention   │  40      │  40      │ ✅         │
│    Iron Laws present         │  25/25   │          │            │
│    Hard Gate present         │  15/15   │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 7. Task CLI Alignment        │  235     │  240     │ ⚠️         │
│    Command existence         │  25/25   │          │            │
│    Flag correctness          │  25/25   │          │            │
│    Output field parsing      │  15/15   │          │            │
│    Status machine align      │  35/35   │          │            │
│    Claim scheduling align    │  35/35   │          │            │
│    Record validation align   │  35/35   │          │            │
│    Dynamic task add align    │  25/25   │          │            │
│    Schema-code alignment     │  15/20   │          │            │
│    All-completed hook align  │  10/10   │          │            │
│    Template existence        │  10/10   │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 8. Hook Wiring Integrity     │  70      │  70      │ ✅         │
│    hooks.json valid JSON     │  10/10   │          │            │
│    Hook scripts exist        │  25/25   │          │            │
│    Hook CLI commands valid   │  15/15   │          │            │
│    Hook event names valid    │  20/20   │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 9. Guide Coverage+Concise    │  70      │  70      │ ✅         │
│    Guide references valid    │  30/30   │          │            │
│    Core workflow skills doc  │  25/25   │          │            │
│    Conciseness / no redund.  │  15/15   │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 10. Command Metadata         │  60      │  60      │ ✅         │
│    allowed_tools declared    │  35/35   │          │            │
│    argument-hints declared   │  25/25   │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 11. Plugin Metadata          │  40      │  40      │ ✅         │
│    keywords coverage         │  25/25   │          │            │
│    description accurate      │  15/15   │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 12. Safety Marker Consist.   │  50      │  50      │ ✅         │
│    Command/agent markers     │  30/30   │          │            │
│    Dispatch cmd coverage     │  20/20   │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ TOTAL                        │  965     │  1000    │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| # | Check | File | Issue | Penalty |
|---|-------|------|-------|---------|
| 1 | 7h-Schema-code alignment | `plugins/forge/skills/breakdown-tasks/templates/index.schema.json` and `quick-tasks/templates/index.schema.json` | Both schemas define `status` per-task enum as `["pending", "in_progress", "completed", "blocked", "skipped", "rejected"]` and top-level `status` enum as `["planning", "in-progress", "completed"]`. The Go code `TaskIndex` struct (types.go:37-48) has `Status string` with no enum constraint, and the `Task` struct (types.go:10-32) has `Status string` also unconstrained. The `validate.go` hardcodes `validStatus` (line 35) as a map, not reading from the schema. The schema enums are documentation-only — they do not match a Go enum type since Go lacks native string enums. This is a documentation-vs-code gap, not a functional bug, but the schema declares strict enums that Go does not enforce. The `scope` field is now correctly optional (not required) in both schemas, matching Go's `omitempty`. The `sourceTaskID` field is correctly absent from both schemas (auto-managed). | -5 |

---

## Attack Points

### Attack 1: [7h — Schema declares enums not enforced by Go code]

**Where**: `plugins/forge/skills/breakdown-tasks/templates/index.schema.json:61` and `plugins/forge/skills/quick-tasks/templates/index.schema.json:55`
**What's wrong**: Both JSON schemas declare strict per-task `status` enum: `["pending", "in_progress", "completed", "blocked", "skipped", "rejected"]` and top-level `status` enum: `["planning", "in-progress", "completed"]`. The Go code validates status via `validStatus` map in `validate.go:35` and `index.StatusEnum` in `types.go:148-155` (which provides the same 6 values). The Go code is consistent with the schema enums. However, the top-level `status` field in the schema uses `["planning", "in-progress", "completed"]` while Go has no enum constraint on `TaskIndex.Status` — it's just a `string`. The task-level status is effectively enforced via `validate.go:35` which matches the schema. This is a documentation-only gap for the top-level status. Net: very minor inconsistency.
**How to fix**: Either add a top-level status enum to the Go `TaskIndex` validation or change the schema to document that top-level status is advisory.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Iter-1 Attack 1: [7f — Record validation quality gate mismatch (4 vs 3 steps)] | ✅ | Verified: `DefaultGateSequence()` in `just.go` includes 4 steps (compile, fmt, lint, test). SKILL.md line 130 shows all 4. Correct. |
| Iter-1 Attack 2: [10 — simplify-skill missing Write/Edit in allowed_tools] | ✅ | `simplify-skill.md:5` has `allowed_tools: ["Read", "Write", "Edit", "AskUserQuestion"]`. |
| Iter-1 Attack 3: [7h — Schema-code scope required mismatch] | ✅ | Both schemas (`breakdown-tasks/templates/index.schema.json:51` and `quick-tasks/templates/index.schema.json:45`) now have `"required": ["id", "title", "priority", "status", "file"]` — scope removed from required. Matches Go `omitempty` on `Task.Scope`. |
| Iter-1 Attack 4: [9 — Guide duplicated task record workflow detail] | ✅ | `guide.md` Key Commands table is concise; record detail delegates to `/record-task` skill. |
| Iter-1 Attack 5: [3 — Quality gate table duplicated across 4 files] | ✅ | `execute-task.md` and `fix-bug.md` use single-line "See Forge Guide Quality Gate Protocol". Autonomous agents (task-executor, error-fixer) retain full tables per exception. |
| Iter-1 Attack 6: [7i — All-completed hook missing auto-fix behavior] | ✅ | `guide.md:133` describes auto fix-task creation. `all_completed.go:128` implements it. |
| Iter-1 Attack 7: [12 — error-fixer lacks ONE-TASK constraint] | ✅ | `error-fixer.md:167-181` has STOP section with HARD-RULE and PROHIBITIONS. |
| Iter-2 Attack 1: [7f — --force does NOT override auto-downgrade] | ✅ | SKILL.md line 151 explicitly states auto-downgrade is never overridden. `record.go:270-273` applies before force check. |
| Iter-2 Attack 2: [3 — Quality gate failure actions duplicated] | ✅ | Main-session commands now use inline references. Resolved. |
| Iter-2 Attack 3: [7j — fix-task template exists only in binary] | ✅ | `plugins/forge/skills/breakdown-tasks/templates/fix-task.md` exists on disk with full template. Matches `task template fix-task` output. |
| Iter-3 Attack 1: [7f — Record validation --force description precision] | ✅ | SKILL.md line 146 and line 151 are consistent. Line 151 explicitly states auto-downgrade is never overridden by `--force`. Go code confirms. Acceptable wording. |
| Iter-3 Attack 2: [7h — quick-tasks schema still requires scope] | ✅ | `quick-tasks/templates/index.schema.json:45` now has `"required": ["id", "title", "priority", "status", "file"]` — scope removed. Fixed since iteration 3. |

---

## Fix Summary

| File Changed | What Changed |
|-------------|--------------|
| N/A | No fixes applied this iteration — final verification audit. The quick-tasks schema scope issue from iteration 3 was already fixed prior to this audit. |

---

## Verdict

- **Score**: 965/1000
- **Target**: 950/1000
- **Gap**: 0 points (target exceeded by 15)
- **Action**: Target reached. Remaining deduction is a single minor schema-vs-code documentation gap (top-level status enum in JSON schema not enforced by Go). No critical, high, or medium-severity issues found. Plugin is structurally consistent across all 12 dimensions.
