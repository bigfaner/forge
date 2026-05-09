---
date: "2026-05-09"
plugin_version: "2.16.1"
iteration: 2
target: 950
evaluator: Claude (structural audit)
---

# Forge Plugin Audit — Iteration 2

**Score: 920/1000** (target: 950)

```
┌─────────────────────────────────────────────────────────────────┐
│                  PLUGIN CONSISTENCY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 1. Directory-Name Alignment  │  40      │  40      │ OK         │
│    Skill name matches dir    │  25/25   │          │            │
│    Command name matches file │  15/15   │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 2. Agent Reference Integrity │  100     │  100     │ OK         │
│    Referenced agents exist   │  70/70   │          │            │
│    No orphan agents          │  30/30   │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 3. Reference Integrity       │  75      │  80      │ WARN       │
│    Template refs valid       │  25/25   │          │            │
│    Cross-skill refs valid    │  30/30   │          │            │
│    No orphan templates       │  20/25   │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 4. Frontmatter Completeness  │  110     │  110     │ OK         │
│    Skill frontmatter         │  45/45   │          │            │
│    Command frontmatter       │  35/35   │          │            │
│    Agent frontmatter         │  30/30   │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 5. Eval Template Convention  │  100     │  100     │ OK         │
│    rubric.md exists          │  30/30   │          │            │
│    report.md exists          │  30/30   │          │            │
│    Rubric→report chain valid │  20/20   │          │            │
│    Rubric totals correct     │  20/20   │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 6. Orchestrator Convention   │  40      │  40      │ OK         │
│    Iron Laws present         │  25/25   │          │            │
│    Hard Gate present         │  15/15   │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 7. Task CLI Alignment        │  195     │  240     │ WARN       │
│    Command existence         │  25/25   │          │            │
│    Flag correctness          │  25/25   │          │            │
│    Output field parsing      │  10/15   │          │            │
│    Status machine align      │  35/35   │          │            │
│    Claim scheduling align    │  35/35   │          │            │
│    Record validation align   │  35/35   │          │            │
│    Dynamic task add align    │  25/25   │          │            │
│    Schema-code alignment     │  5/20    │          │            │
│    All-completed hook align  │  10/10   │          │            │
│    Template existence        │  0/10    │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 8. Hook Wiring Integrity     │  70      │  70      │ OK         │
│    hooks.json valid JSON     │  10/10   │          │            │
│    Hook scripts exist        │  25/25   │          │          │
│    Hook CLI commands valid   │  15/15   │          │            │
│    Hook event names valid    │  20/20   │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 9. Guide Coverage            │  60      │  70      │ WARN       │
│    Guide references valid    │  30/30   │          │            │
│    Core skills documented    │  30/40   │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 10. Command Metadata         │  50      │  60      │ WARN       │
│    allowed_tools declared    │  35/35   │          │            │
│    argument-hints declared   │  15/25   │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 11. Plugin Metadata          │  40      │  40      │ OK         │
│    keywords coverage         │  25/25   │          │            │
│    description accurate      │  15/15   │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 12. Safety Marker Consist.   │  40      │  50      │ WARN       │
│    Command/agent markers     │  25/30   │          │            │
│    Dispatch cmd coverage     │  15/20   │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ TOTAL                        │  920     │  1000    │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| # | Check | File | Issue | Penalty |
|---|-------|------|-------|---------|
| 1 | 3. No orphan templates | `plugins/forge/skills/record-task/templates/template.md` | Not referenced by `record-task/SKILL.md`. The SKILL.md uses `task record` CLI and describes JSON format inline. The template exists at `skills/record-task/templates/template.md` but no line in the SKILL.md references it. | -5 |
| 2 | 7c. Output field parsing | `plugins/forge/commands/run-tasks.md` line 51 | run-tasks extracts `SCOPE` from claim output without noting it may be absent. In `claim.go` line 320, SCOPE is printed with `PrintFieldIfNotEmpty` (omitted when empty). Guide.md line 179 correctly notes "No" in "Always present" column for SCOPE. run-tasks does not handle the absent case. | -5 |
| 3 | 7h. Schema-code alignment | `plugins/forge/skills/breakdown-tasks/templates/index.schema.json` line 6 | Schema marks `prd` and `design` as required. Go `TaskIndex` struct (types.go lines 37-42) has both as optional (no required JSON tag). Go validator (validate.go lines 87-89) only warns when missing AND Proposal is not set. KNOWN ACCEPTABLE DISCREPANCY per rubric: "Schema marks prd/design as required, but Go allows Proposal as alternative." No deduction per instructions. | 0 (INFO) |
| 4 | 7h. Schema-code alignment | `plugins/forge/skills/breakdown-tasks/templates/index.schema.json` | Schema does NOT include `proposal` field. Go struct has `Proposal string` (types.go line 41). Quick-tasks uses `proposal` field in index.json. Schema is missing the `proposal` property entirely. | -5 |
| 5 | 7h. Schema-code alignment | `plugins/forge/skills/breakdown-tasks/templates/index.schema.json` | Schema top-level `status` enum is `["planning", "in-progress", "completed"]` at line 23. This is a feature-level status (not task-level). The values `planning` and `in-progress` are used by manifest.md but the enum uses hyphenated `in-progress` while task-level uses underscored `in_progress`. These are different fields (feature-level vs task-level). However, the schema does not validate against the Go code's feature status values. Go TaskIndex has `Status string` without enum constraint. The manifest uses these values. Low severity but inconsistent naming convention. | -5 |
| 6 | 7h. Schema-code alignment | `plugins/forge/skills/breakdown-tasks/templates/index.schema.json` | `sourceTaskID` exists in Go struct (types.go line 26) but not in JSON schema. KNOWN ACCEPTABLE DISCREPANCY per rubric: "sourceTaskID exists in Go struct but not in JSON schema (auto-managed field)." No deduction. | 0 (INFO) |
| 7 | 7j. Template existence | `task-cli/pkg/template/data/fix-task.md` | The fix-task template is embedded in the CLI binary, accessible via `task template fix-task`. It does NOT exist as a file in `plugins/forge/skills/*/templates/`. The rubric criterion says "fix-task template referenced by skills exists on disk." The instructions state: "The fix-task template is embedded in the task-cli binary, accessible via `task template fix-task`. It does NOT need to exist as a separate file on disk." Per instructions, no deduction. | 0 (INFO) |
| 8 | 10. argument-hints declared | `plugins/forge/commands/init-justfile.md` | Has no `argument-hints` in frontmatter. Accepts `--lang`, `--type`, `--force` as user-facing flags described in the Parameters table (lines 32-36). These are parameters users pass. Missing argument-hints. | -5 |
| 9 | 10. argument-hints declared | `plugins/forge/commands/quick.md` | Has no `argument-hints` in frontmatter. Accepts `--no-test` flag described at line 24. This is a user-facing parameter. Missing argument-hints. | -5 |
| 10 | 12. Command/agent markers | `plugins/forge/commands/fix-bug.md` | fix-bug now has `<EXTREMELY-IMPORTANT>` at line 34 (added since iteration 1). However, the `<HARD-RULE>` at line 127 ("If the new unit test passes before any fix, the test does not capture the bug") is actionable and correct. The `<HARD-GATE>` at lines 87 and 189 are actionable. But the `<EXTREMELY-IMPORTANT>` block at lines 34-40 contains constraint "Tests and fix must be committed together in a single atomic commit" — this contradicts the task-executor agent which does Step 4 (record) then Step 5 (commit) separately. In fix-bug, Step 6 is commit-only (no separate record step). The contradiction is between fix-bug's "one atomic commit" and task-executor's "record then commit" pattern. However, fix-bug is a standalone command, not dispatched through task-executor. The contexts are different: fix-bug is manual, task-executor is automated. Not a true contradiction. The markers are actionable. No deduction. | 0 |
| 11 | 9. Core skills documented | `plugins/forge/hooks/guide.md` | Missing: `/record-task`, `/git-commit`, `/init-forge`, `/init-justfile`, `/git-checkout`, `/improve-harness`, `/forensic`, `/simplify-skill`, `/extract-design-md`. The guide.md mentions `task record` in the CLI section (line 158) but not the `/record-task` skill by name. `/record-task` is mandatory for task completion and should be documented. `/git-commit` is used in every task completion. `/improve-harness` is the fix counterpart to `/eval-harness`. `/forensic` is a specialized tool. `/init-forge`, `/init-justfile`, `/git-checkout`, `/simplify-skill`, `/extract-design-md` are utility commands. Deducting for `/record-task` and `/git-commit` as they are core workflow skills: -5 each. `/improve-harness` is also important as the fix counterpart to eval-harness: -5. Other utilities are lower priority: no deduction. | -10 |

---

## Attack Points

### Attack 1: [7h — Schema missing `proposal` field and feature-level status mismatch]

**Where**: `plugins/forge/skills/breakdown-tasks/templates/index.schema.json` lines 6, 23
**What's wrong**: (1) The `proposal` field exists in Go `TaskIndex` struct (types.go line 41) and is used by quick-tasks, but the JSON schema does not include it at all. (2) The top-level `status` enum uses `["planning", "in-progress", "completed"]` which is a different naming convention from the task-level `in_progress` (underscore vs hyphen). The Go code does not enforce feature-level status values, making the schema's enum somewhat aspirational.
**How to fix**: Add `"proposal": { "type": "string" }` to the schema properties. Consider aligning the feature-level status enum to use underscores matching the task-level convention, or document the intentional difference.

### Attack 2: [3 — Orphan template in record-task]

**Where**: `plugins/forge/skills/record-task/templates/template.md`
**What's wrong**: The `record-task/SKILL.md` does not reference `templates/template.md` anywhere in its content. The skill describes JSON format inline (lines 29-59) and delegates output generation entirely to `task record` CLI. The template appears to be a legacy artifact from before the CLI took over record generation. The actual record template is embedded in Go code (`record.go` lines 314-384, `fillRecordTemplate`).
**How to fix**: Either remove `skills/record-task/templates/template.md` as it is superseded by the CLI's embedded template, or add a reference to it in the SKILL.md for documentation purposes.

### Attack 3: [9+10 — Missing guide coverage and argument-hints for utility commands]

**Where**: `plugins/forge/hooks/guide.md`, `plugins/forge/commands/init-justfile.md`, `plugins/forge/commands/quick.md`
**What's wrong**: (1) guide.md does not mention `/record-task` (the mandatory recording step), `/git-commit` (used in every task), or `/improve-harness` (fix counterpart to eval-harness). These are core workflow skills. (2) `init-justfile.md` and `quick.md` lack `argument-hints` for their user-facing flags (`--lang`/`--type`/`--force` and `--no-test` respectively).
**How to fix**: (1) Add `/record-task`, `/git-commit`, and `/improve-harness` mentions to guide.md. (2) Add `argument-hints` frontmatter to init-justfile.md and quick.md.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1: Hard Gate missing in eval-proposal | YES | eval-proposal SKILL.md now has `<HARD-GATE>` at line 92-94 |
| Attack 1: Hard Gate missing in eval-consistency | YES (was already present) | eval-consistency SKILL.md has `<HARD-GATE>` at line 143-145 |
| Attack 2: Schema-code alignment — missing proposal field | NO | `index.schema.json` still does not include `proposal` field |
| Attack 3: SubagentStart hook event | YES (confirmed valid) | Instructions confirm SubagentStart is valid; no deduction in iter 2 |
| Attack 4: disc-N not mentioned in skills | YES | breakdown-tasks SKILL.md line 344, run-tasks.md line 169, execute-task.md line 70, task-executor.md line 202 all now mention disc-N format |
| Attack 5: fix-bug missing EXTREMELY-IMPORTANT | YES | fix-bug.md now has `<EXTREMELY-IMPORTANT>` at lines 34-40 |
| Attack 6: Missing argument-hints for init-justfile and quick | NO | init-justfile.md and quick.md still lack `argument-hints` |
| Attack 7: Orphan template in record-task | NO | `skills/record-task/templates/template.md` still not referenced by SKILL.md |
| Attack 8: record-task not documented in guide.md | NO | guide.md still does not mention `/record-task` skill |

---

## Fix Summary

| File Changed | What Changed |
|-------------|--------------|
| `plugins/forge/skills/eval-proposal/SKILL.md` | Added `<HARD-GATE>` block to Step 3 (Decision Gate) |
| `plugins/forge/commands/fix-bug.md` | Added `<EXTREMELY-IMPORTANT>` safety block |
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | Added disc-N format documentation |
| `plugins/forge/commands/run-tasks.md` | Added disc-N format documentation |
| `plugins/forge/commands/execute-task.md` | Added disc-N format documentation |
| `plugins/forge/agents/task-executor.md` | Added disc-N format documentation |

---

## Verdict

- **Score**: 920/1000
- **Target**: 950/1000
- **Gap**: 30 points
- **Action**: Remaining issues are schema alignment (proposal field), orphan template, guide coverage, and argument-hints. Fix applied, re-auditing
