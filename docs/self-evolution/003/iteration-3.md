---
date: "2026-05-09"
plugin_version: "2.16.1"
iteration: 3
target: 950
evaluator: Claude (structural audit)
---

# Forge Plugin Audit — Iteration 3

**Score: 965/1000** (target: 950)

```
┌─────────────────────────────────────────────────────────────────┐
│                  PLUGIN CONSISTENCY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 1. Directory-Name Alignment  │  40      │  40      │ OK         │
│    Skill name matches dir    │  25/25   │          │            │
│    Command name matches file │  15/15   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Agent Reference Integrity │  100     │  100     │ OK         │
│    Referenced agents exist   │  70/70   │          │            │
│    No orphan agents          │  30/30   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Reference Integrity       │  80      │  80      │ OK         │
│    Template refs valid       │  25/25   │          │            │
│    Cross-skill refs valid    │  30/30   │          │            │
│    No orphan templates       │  25/25   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Frontmatter Completeness  │  110     │  110     │ OK         │
│    Skill frontmatter         │  45/45   │          │            │
│    Command frontmatter       │  35/35   │          │            │
│    Agent frontmatter         │  30/30   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Eval Template Convention  │  100     │  100     │ OK         │
│    rubric.md exists          │  30/30   │          │            │
│    report.md exists          │  30/30   │          │            │
│    Rubric→report chain valid │  20/20   │          │            │
│    Rubric totals correct     │  20/20   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 6. Orchestrator Convention   │  40      │  40      │ OK         │
│    Iron Laws present         │  25/25   │          │            │
│    Hard Gate present         │  15/15   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 7. Task CLI Alignment        │  240     │  240     │ OK         │
│    Command existence         │  25/25   │          │            │
│    Flag correctness          │  25/25   │          │            │
│    Output field parsing      │  15/15   │          │            │
│    Status machine align      │  35/35   │          │            │
│    Claim scheduling align    │  35/35   │          │            │
│    Record validation align   │  35/35   │          │            │
│    Dynamic task add align    │  25/25   │          │            │
│    Schema-code alignment     │  20/20   │          │            │
│    All-completed hook align  │  10/10   │          │            │
│    Template existence        │  10/10   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 8. Hook Wiring Integrity     │  70      │  70      │ OK         │
│    hooks.json valid JSON     │  10/10   │          │            │
│    Hook scripts exist        │  25/25   │          │            │
│    Hook CLI commands valid   │  15/15   │          │            │
│    Hook event names valid    │  20/20   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 9. Guide Coverage            │  35      │  70      │ WARN       │
│    Guide references valid    │  30/30   │          │            │
│    Core skills documented    │  5/40    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 10. Command Metadata         │  60      │  60      │ OK         │
│    allowed_tools declared    │  35/35   │          │            │
│    argument-hints declared   │  25/25   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 11. Plugin Metadata          │  40      │  40      │ OK         │
│    keywords coverage         │  25/25   │          │            │
│    description accurate      │  15/15   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 12. Safety Marker Consist.   │  50      │  50      │ OK         │
│    Command/agent markers     │  30/30   │          │            │
│    Dispatch cmd coverage     │  20/20   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  965     │  1000    │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| # | Check | File | Issue | Penalty |
|---|-------|------|-------|---------|
| 1 | 9. Core skills documented | `plugins/forge/hooks/guide.md` | `/forensic` skill not mentioned anywhere in guide.md. Forensic is a specialized debugging tool for analyzing agent deviations. | -5 |
| 2 | 9. Core skills documented | `plugins/forge/hooks/guide.md` | `/eval-consistency` skill not mentioned. This eval-* skill evaluates cross-document consistency and is part of the eval suite. | -5 |
| 3 | 9. Core skills documented | `plugins/forge/hooks/guide.md` | `/simplify-skill` command not mentioned. Maintenance/refactoring tool for skill files. | -5 |
| 4 | 9. Core skills documented | `plugins/forge/hooks/guide.md` | `/extract-design-md` command not mentioned. Tool for extracting visual style from web apps. | -5 |
| 5 | 9. Core skills documented | `plugins/forge/hooks/guide.md` | `/init-forge` command not mentioned. Setup command to build and install task-cli. | -5 |
| 6 | 9. Core skills documented | `plugins/forge/hooks/guide.md` | `/init-justfile` command not mentioned. Setup command to scaffold Justfile targets. | -5 |
| 7 | 9. Core skills documented | `plugins/forge/hooks/guide.md` | `/git-checkout` command not mentioned. Utility for creating feature branches. | -5 |

---

## Attack Points

### Attack 1: [9 — Guide coverage incomplete for utility commands and specialized skills]

**Where**: `plugins/forge/hooks/guide.md`
**What's wrong**: Seven skills/commands have zero mentions in guide.md: `/forensic`, `/eval-consistency`, `/simplify-skill`, `/extract-design-md`, `/init-forge`, `/init-justfile`, `/git-checkout`. The rubric requires every skill and command directory to have "at least a mention in guide.md (workflow reference, usage example, or description)." While these are specialized/utility tools rather than core workflow skills, the rubric applies uniformly.
**How to fix**: Add brief mentions to guide.md. Options: (1) Add a "Utility Commands" section listing setup/debug/maintenance tools, or (2) add parenthetical mentions in relevant workflow sections (e.g., "Run `/init-justfile` to scaffold build targets" in the setup context, "Use `/forensic` to investigate agent deviations" in a troubleshooting section).

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1: Hard Gate missing in eval-proposal | YES | `plugins/forge/skills/eval-proposal/SKILL.md` line 92 has `<HARD-GATE>` |
| Attack 2: Schema-code alignment — missing proposal field | YES | `plugins/forge/skills/breakdown-tasks/templates/index.schema.json` line 17 now has `proposal` field with description |
| Attack 3: Orphan template in record-task | YES | `plugins/forge/skills/record-task/templates/` directory no longer exists — orphan removed |
| Iter 2 Attack 1: Schema missing proposal + feature-level status mismatch | YES | `proposal` field added to schema. Feature-level status `in-progress` is acknowledged as a different field from task-level `in_progress` |
| Iter 2 Attack 2: Orphan template in record-task | YES | Template file removed |
| Iter 2 Attack 3: Missing guide coverage + argument-hints | PARTIAL | `/record-task`, `/git-commit`, `/improve-harness` now mentioned in guide.md (lines 111, 133). `init-justfile.md` now has `argument-hints` (line 5). `quick.md` now has `argument-hints` (line 5). Remaining: 7 utility/specialized commands still undocumented in guide.md |
| SubagentStart hook event | YES | Instructions confirm SubagentStart is valid; no deduction |
| fix-task template on disk | YES | Instructions confirm embedded template is acceptable; no deduction |
| prd/design vs Proposal | YES | Known acceptable discrepancy per rubric; noted as INFO |
| sourceTaskID in schema | YES | Known acceptable discrepancy per rubric; noted as INFO |
| eval-harness lacking Iron Laws/Hard Gate | YES | Known acceptable discrepancy per rubric; no deduction |

---

## Fix Summary

| File Changed | What Changed |
|-------------|--------------|
| `plugins/forge/skills/eval-proposal/SKILL.md` | Added `<HARD-GATE>` block (iteration 1 fix) |
| `plugins/forge/commands/fix-bug.md` | Added `<EXTREMELY-IMPORTANT>` safety block (iteration 1 fix) |
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | Added disc-N format documentation (iteration 1 fix) |
| `plugins/forge/commands/run-tasks.md` | Added disc-N format documentation (iteration 1 fix) |
| `plugins/forge/commands/execute-task.md` | Added disc-N format documentation (iteration 1 fix) |
| `plugins/forge/agents/task-executor.md` | Added disc-N format documentation (iteration 1 fix) |
| `plugins/forge/skills/breakdown-tasks/templates/index.schema.json` | Added `proposal` field with description (iteration 2-3 fix) |
| `plugins/forge/skills/record-task/templates/template.md` | Removed orphan template file (iteration 2-3 fix) |
| `plugins/forge/commands/init-justfile.md` | Added `argument-hints` to frontmatter (iteration 2-3 fix) |
| `plugins/forge/commands/quick.md` | Added `argument-hints` to frontmatter (iteration 2-3 fix) |
| `plugins/forge/hooks/guide.md` | Added `/record-task`, `/git-commit`, `/improve-harness` mentions (iteration 2-3 fix) |

---

## Known Acceptable Discrepancies (INFO, no deduction)

1. **Schema marks prd/design required, Go allows Proposal as alternative** — validate.go (lines 87-89) only warns when missing AND Proposal is not set. Quick-tasks uses Proposal instead of PRD+Design.

2. **sourceTaskID exists in Go struct but not in JSON schema** — auto-managed field, injected by `--source-task-id` flag in add.go. Per rubric Known Acceptable Discrepancies.

3. **eval-harness lacks Iron Laws/Hard Gate by design** — single-pass evaluation, no adversarial loop. Per rubric Known Acceptable Discrepancies.

4. **Feature-level status enum uses hyphens (`in-progress`) vs task-level underscores (`in_progress`)** — these are different fields at different levels. The feature-level status tracks pipeline stage (planning/in-progress/completed), task-level status tracks execution state.

5. **SubagentStart is a valid Claude Code hook event** — confirmed per instructions. No deduction.

6. **fix-task template is embedded in task-cli binary** — accessible via `task template fix-task`. Per instructions, does not need to exist as a separate file on disk.

---

## Verdict

- **Score**: 965/1000
- **Target**: 950/1000
- **Gap**: 0 points (target exceeded by 15)
- **Action**: Target reached. Remaining gap is guide coverage for 7 utility/specialized commands (forensic, eval-consistency, simplify-skill, extract-design-md, init-forge, init-justfile, git-checkout). These are non-critical — core workflow skills are fully documented.
