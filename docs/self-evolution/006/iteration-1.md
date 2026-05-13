---
date: "2026-05-13"
plugin_version: "3.0.0-beta-3"
iteration: "1"
target: "900"
evaluator: Claude (structural audit)
---

# Forge Plugin Audit -- Iteration 1

**Score: 875/1000** (target: 900)

```
+---------------------------------------------------------------+
|                  PLUGIN CONSISTENCY SCORECARD                  |
+------------------------------+----------+----------+-----------+
| Dimension                    | Score    | Max      | Status    |
+------------------------------+----------+----------+-----------+
| 1. Directory-Name Alignment  |  40      |  40      | OK        |
|    Skill name matches dir    |  25/25   |          |           |
|    Command name matches file |  15/15   |          |           |
+------------------------------+----------+----------+-----------+
| 2. Agent Reference Integrity |  100     |  100     | OK        |
|    Referenced agents exist   |  70/70   |          |           |
|    No orphan agents          |  30/30   |          |           |
+------------------------------+----------+----------+-----------+
| 3. Reference Integrity       |  70      |  80      | WARN      |
|    Template refs valid       |  25/25   |          |           |
|    Cross-skill refs valid    |  30/30   |          |           |
|    No orphan templates       |  10/15   |          |           |
|    No cross-file duplication |  5/10    |          |           |
+------------------------------+----------+----------+-----------+
| 4. Frontmatter Completeness  |  110     |  110     | OK        |
|    Skill frontmatter         |  45/45   |          |           |
|    Command frontmatter       |  35/35   |          |           |
|    Agent frontmatter         |  30/30   |          |           |
+------------------------------+----------+----------+-----------+
| 5. Eval Template Convention  |  100     |  100     | OK        |
|    rubric.md exists          |  30/30   |          |           |
|    report.md exists          |  30/30   |          |           |
|    Rubric->report chain      |  20/20   |          |           |
|    Rubric totals correct     |  20/20   |          |           |
+------------------------------+----------+----------+-----------+
| 6. Orchestrator Convention   |  40      |  40      | OK        |
|    Iron Laws present         |  25/25   |          |           |
|    Hard Gate present         |  15/15   |          |           |
+------------------------------+----------+----------+-----------+
| 7. Task CLI Alignment        |  165     |  240     | FAIL      |
|    Command existence (7a)    |  10/25   |          |           |
|    Flag correctness (7b)     |  10/25   |          |           |
|    Output field parsing (7c) |  10/15   |          |           |
|    Status machine align (7d) |  35/35   |          |           |
|    Claim scheduling (7e)     |  35/35   |          |           |
|    Record validation (7f)    |  20/35   |          |           |
|    Dynamic task add (7g)     |  15/25   |          |           |
|    Schema-code align (7h)    |  15/20   |          |           |
|    All-completed hook (7i)   |  5/10    |          |           |
|    Template existence (7j)   |  10/10   |          |           |
+------------------------------+----------+----------+-----------+
| 8. Hook Wiring Integrity     |  55      |  70      | WARN      |
|    hooks.json valid JSON     |  10/10   |          |           |
|    Hook scripts exist        |  25/25   |          |           |
|    Hook CLI commands valid   |  15/15   |          |           |
|    Hook event names valid    |  5/20    |          |           |
+------------------------------+----------+----------+-----------+
| 9. Guide Coverage+Concise    |  60      |  70      | WARN      |
|    Guide references valid    |  30/30   |          |           |
|    Core workflow skills doc  |  25/25   |          |           |
|    Conciseness / no redund.  |  5/15    |          |           |
+------------------------------+----------+----------+-----------+
| 10. Command Metadata         |  60      |  60      | OK        |
|    allowed_tools declared    |  35/35   |          |           |
|    argument-hints declared   |  25/25   |          |           |
+------------------------------+----------+----------+-----------+
| 11. Plugin Metadata          |  30      |  40      | WARN      |
|    keywords coverage         |  15/25   |          |           |
|    description accurate      |  15/15   |          |           |
+------------------------------+----------+----------+-----------+
| 12. Safety Marker Consist.   |  45      |  50      | WARN      |
|    Command/agent markers     |  25/30   |          |           |
|    Dispatch cmd coverage     |  20/20   |          |           |
+------------------------------+----------+----------+-----------+
| TOTAL                        |  875     |  1000    |           |
+------------------------------+----------+----------+-----------+
```

---

## Deductions

| # | Check | File | Issue | Penalty |
|---|-------|------|-------|---------|
| 1 | No orphan templates (3) | `plugins/forge/skills/gen-test-scripts/templates/node_modules/` | Entire node_modules directory (100+ files) in templates/ is not referenced by any SKILL.md. Not a template -- runtime dependencies accidentally committed. | -5 |
| 2 | No cross-file duplication (3) | `plugins/forge/commands/execute-task.md` + `plugins/forge/commands/run-tasks.md` | Steps 1-3 are nearly identical across both files (200+ lines duplicated). Canonical location should be one of them; the other should reference it. | -5 |
| 3 | Command existence (7a) | `plugins/forge/skills/gen-test-cases/SKILL.md` lines 18-19; `gen-test-scripts/SKILL.md` lines 29-33; `run-e2e-tests/SKILL.md` lines 19-24; `graduate-tests/SKILL.md` lines 18-27; `breakdown-tasks/SKILL.md` lines 26-30; `quick-tasks/SKILL.md` lines 25-30; `eval-test-cases/SKILL.md` lines 9-13; `tech-design/SKILL.md` Step 0 | All reference `task profile`, `task profile set <name>`, `task profile get <profile-name> --manifest/--generate/--run/--graduate` -- but `task profile` command does NOT exist in the CLI (`task -h` shows no "profile" subcommand). Every invocation of `task profile` will produce "unknown command 'profile' for 'task'" at runtime. This is the single most severe alignment gap: the entire test profile resolution workflow documented across 8 skills is fictional. | -15 |
| 4 | Flag correctness (7b) | Same files as #3 | Flags `--manifest`, `--generate`, `--run`, `--graduate` on the non-existent `task profile get` command are all fictional. The CLI has no such command or flags. | -15 |
| 5 | Output field parsing (7c) | `plugins/forge/commands/run-tasks.md` line 53 | SCOPE field documented as "may be omitted entirely by claim output when not set", but `claim.go` printTaskDetails (line 322) uses `PrintFieldIfNotEmpty("SCOPE", t.Scope)` which always prints it when non-empty. Minor description inconsistency. | -5 |
| 6 | Record validation align (7f) | `plugins/forge/commands/execute-task.md`; `plugins/forge/commands/run-tasks.md` | The auto-downgrade rule in `record.go` line 271-274 (`completed + testsFailed > 0 -> blocked`, non-overridable) is not documented in the dispatcher commands. The quality gate (`just compile -> fmt -> lint -> test`) enforced by `record.go` line 122-124 before allowing completion is also not described in the dispatchers -- they delegate to task-executor which calls record-task which calls `task record`, so the gate IS enforced, but the dispatcher commands are unaware of these rules. The rubric requires skills to match these rules. | -15 |
| 7 | Dynamic task add (7g) | `plugins/forge/commands/execute-task.md`; `plugins/forge/commands/run-tasks.md` | Skills don't document that the generated ID format is `disc-N` (from `add.go` line 43). The `--source-task-id` auto-resolution to root ancestor is documented in run-tasks.md line 139, but the pre-add pattern (block source -> add fix -> claim picks up) is not explicitly spelled out step-by-step. The `--block-source` flag atomically handles blocking, which is correctly documented. | -10 |
| 8 | Schema-code alignment (7h) | `plugins/forge/skills/breakdown-tasks/templates/index.schema.json` | Schema index-level `status` field has enum `["planning", "in-progress", "completed"]` but Go TaskIndex has no validation for this field. Schema marks `prd` and `design` as required, but Go allows `Proposal` as alternative (Known Acceptable Discrepancy). Schema includes `noTest` (deprecated) which still matches Go struct. Minor: schema has `e2eRound` field with min/max constraints not enforced in Go. | -5 |
| 9 | All-completed hook align (7i) | `plugins/forge/hooks/guide.md` lines 128-133 | Guide says all-completed runs: "1. Quality gate: just compile -> just fmt -> just lint  2. Project-wide tests: just test  3. E2E regression: just e2e-setup -> just test-e2e". But `all_completed.go` lines 167-172 runs `e2eprobe.ProbeServers()` between e2e-setup and test-e2e -- the server health probe step is not mentioned in the guide. Also guide says "On failure at any step, a P0 fix-task is automatically created" but e2e-setup failure just skips e2e with a warning, not a fix task. | -5 |
| 10 | Hook event names (8) | `plugins/forge/hooks/hooks.json` lines 13-21 | Uses `SubagentStart` event. Claude Code supported hook events are: `PreToolUse`, `PostToolUse`, `Notification`, `Stop`, `SubagentStop`. `SubagentStart` is NOT in the standard list. If Claude Code does not fire this event, subagents never receive guide.md injection. | -15 |
| 11 | Conciseness / no redundancy (9) | `plugins/forge/hooks/guide.md` lines 153-166 | The "Task-CLI" section contains a "Key Commands" reference table with 5 CLI commands and their descriptions. The rubric states guide.md "contains only workflow rules and conventions -- no CLI output format tables, no API reference, no information that belongs in task -h or individual SKILL.md files." This table duplicates information available via `task -h`. | -10 |
| 12 | Keywords coverage (11) | `plugins/forge/.claude-plugin/plugin.json` | Keywords are `["pipeline", "brainstorm", "prd", "design", "task", "eval", "e2e", "test-profile"]`. Missing major capability areas: "ui-design" (ui-design skill exists), "sitemap" (gen-sitemap command exists). Each major gap is -5. | -10 |
| 13 | Command/agent markers (12a) | `plugins/forge/agents/doc-scorer.md` line 3 | EXTREMELY-IMPORTANT rule 3: "Never give full marks unless content is genuinely excellent" -- "genuinely excellent" is subjective and non-actionable. Markers should define concrete, verifiable constraints. | -5 |

---

## Attack Points

### Attack 1: [7 -- task profile is a phantom command affecting 8 skills]

**Where**: `plugins/forge/skills/gen-test-cases/SKILL.md` lines 18-19; `gen-test-scripts/SKILL.md` lines 29-33; `run-e2e-tests/SKILL.md` lines 19-24; `graduate-tests/SKILL.md` lines 18-27; `breakdown-tasks/SKILL.md` lines 26-30; `quick-tasks/SKILL.md` lines 25-30; `eval-test-cases/SKILL.md` lines 9-13; `tech-design/SKILL.md` Step 0

**What's wrong**: Eight skills reference `task profile`, `task profile set <name>`, and `task profile get <name> --manifest/--generate/--run/--graduate` as Step 0 of their workflow. Running `task -h` confirms there is NO `profile` subcommand. Running `task profile` produces "unknown command 'profile' for 'task'". Every single invocation at runtime will fail. The entire test profile resolution workflow documented across 8+ skills is fictional. This is the largest alignment gap and costs -30 points alone (7a: -15, 7b: -15).

**How to fix**: Either (1) implement `task profile` in task-cli with the documented subcommands and flags, or (2) replace all Step 0 profile resolution blocks with direct `.forge/config.yaml` reading and remove the fictional CLI commands from all 8 skills.

### Attack 2: [8 -- SubagentStart hook event may not be supported]

**Where**: `plugins/forge/hooks/hooks.json` lines 13-21

**What's wrong**: hooks.json defines a `SubagentStart` event that runs `session-start` for subagents. The Claude Code plugin documentation lists supported hook events as: `PreToolUse`, `PostToolUse`, `Notification`, `Stop`, `SubagentStop`. `SubagentStart` is NOT in this list. If Claude Code never fires this event, subagents will never receive the guide.md injection, silently operating without project conventions. Cost: -15.

**How to fix**: Verify with Claude Code documentation whether `SubagentStart` is supported. If not, remove the `SubagentStart` block from hooks.json. Consider alternative approaches (e.g., having the session-start hook output apply to subagent contexts automatically).

### Attack 3: [9 -- guide.md contains CLI reference material violating conciseness rule]

**Where**: `plugins/forge/hooks/guide.md` lines 153-166

**What's wrong**: The Task-CLI section contains a "Key Commands" table with 5 rows listing CLI commands and descriptions. The rubric explicitly states guide.md should contain "only workflow rules and conventions -- no CLI output format tables, no API reference, no information that belongs in `task -h` or individual SKILL.md files." This table is API reference material. Cost: -10.

**How to fix**: Replace the Key Commands table with a single line: "For full command reference, run `task -h` or `task [command] -h`." Keep the "Typical flow" one-liner as workflow context.

---

## Previous Issues Check

<!-- First iteration, no previous issues -->

---

## Fix Summary

<!-- First iteration, no fixes applied yet -->

---

## Verdict

- **Score**: 875/1000
- **Target**: 900/1000
- **Gap**: 25 points
- **Action**: The highest-impact fix is implementing or removing the `task profile` command family (+30 points). Second priority: validate or remove `SubagentStart` hook event (+15 points). Third: clean guide.md reference table (+10 points). These three fixes alone would bring the score to 930/1000, exceeding the target.
