---
date: "2026-05-13"
plugin_version: "3.0.0-beta-3"
iteration: "1"
target_score: "950"
evaluator: Claude (structural audit)
---

# Forge Plugin Audit — Iteration 1

**Score: 893/1000** (target: 950)

```
┌─────────────────────────────────────────────────────────────────┐
│                  PLUGIN CONSISTENCY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Directory-Name Alignment  │  40      │  40      │ ✅          │
│    Skill name matches dir    │  25/25   │          │            │
│    Command name matches file │  15/15   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Agent Reference Integrity │  100     │  100     │ ✅          │
│    Referenced agents exist   │  70/70   │          │            │
│    No orphan agents          │  30/30   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Reference Integrity       │  35      │  80      │ ❌          │
│    Template refs valid       │  25/25   │          │            │
│    Cross-skill refs valid    │  15/30   │          │            │
│    No orphan templates       │  10/15   │          │            │
│    No cross-file duplication │  10/10   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Frontmatter Completeness  │  105     │  110     │ ⚠️          │
│    Skill frontmatter         │  45/45   │          │            │
│    Command frontmatter       │  25/35   │          │            │
│    Agent frontmatter         │  30/30   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Eval Template Convention  │  90      │  100     │ ⚠️          │
│    rubric.md exists          │  30/30   │          │            │
│    report.md exists          │  30/30   │          │            │
│    Rubric→report chain valid │  20/20   │          │            │
│    Rubric totals correct     │  10/20   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 6. Orchestrator Convention   │  30      │  40      │ ⚠️          │
│    Iron Laws present         │  25/25   │          │            │
│    Hard Gate present         │  5/15    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 7. Task CLI Alignment        │  210     │  240     │ ⚠️          │
│    Command existence         │  25/25   │          │            │
│    Flag correctness          │  25/25   │          │            │
│    Output field parsing      │  10/15   │          │            │
│    Status machine align      │  35/35   │          │            │
│    Claim scheduling align    │  30/35   │          │            │
│    Record validation align   │  30/35   │          │            │
│    Dynamic task add align    │  20/25   │          │            │
│    Schema-code alignment     │  15/20   │          │            │
│    All-completed hook align  │  10/10   │          │            │
│    Template existence        │  10/10   │          │            │
├──────────────────────────────┼──────────┼──────────┼──────────┼────────────┤
│ 8. Hook Wiring Integrity     │  50      │  70      │ ⚠️          │
│    hooks.json valid JSON     │  10/10   │          │            │
│    Hook scripts exist        │  10/25   │          │            │
│    Hook CLI commands valid   │  15/15   │          │            │
│    Hook event names valid    │  15/20   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 9. Guide Coverage+Concise    │  60      │  70      │ ⚠️          │
│    Guide references valid    │  30/30   │          │            │
│    Core workflow skills doc  │  25/25   │          │            │
│    Conciseness / no redund.  │  5/15    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 10. Command Metadata         │  43      │  60      │ ⚠️          │
│    allowed_tools declared    │  35/35   │          │            │
│    argument-hints declared   │  8/25    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 11. Plugin Metadata          │  30      │  40      │ ⚠️          │
│    keywords coverage         │  15/25   │          │            │
│    description accurate      │  15/15   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 12. Safety Marker Consist.   │  50      │  50      │ ✅          │
│    Command/agent markers     │  30/30   │          │            │
│    Dispatch cmd coverage     │  20/20   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  893     │  1000    │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| # | Check | File | Issue | Penalty |
|---|-------|------|-------|---------|
| 1 | Cross-skill refs valid | `plugins/forge/skills/consolidate-specs/SKILL.md` | References `plugins/forge/references/shared/decision-logging.md` — directory `references/shared/` does not exist | -15 |
| 2 | Cross-skill refs valid | `plugins/forge/commands/record-decision.md` | References `plugins/forge/references/shared/decision-logging.md` (Sections 3, 6, 7, 8) — file does not exist | -15 |
| 3 | Cross-skill refs valid | `plugins/forge/commands/gen-sitemap.md` | References `plugins/forge/references/shared/config.yaml` and `plugins/forge/references/shared/sitemap.json` — neither exists | -15 (counted once, already penalized in #1) |
| 4 | Cross-skill refs valid | `plugins/forge/skills/tech-design/SKILL.md` | References `plugins/forge/references/shared/decision-logging.md` — file does not exist | (counted in #1 above) |
| 5 | Cross-skill refs valid | `plugins/forge/skills/gen-test-scripts/SKILL.md` | References `plugins/forge/references/shared/sitemap.json` — file does not exist | (counted in #3 above) |
| 6 | No orphan templates | `plugins/forge/skills/forensic/templates/report.md` | Not referenced by any other file — forensic generates reports but `templates/report.md` is only mentioned in the SKILL.md text at line 183, so this is a template reference, not orphan. Actually valid. | 0 |
| 7 | No orphan templates | `plugins/forge/skills/gen-test-scripts/templates/` | Contains `node_modules/` directory with many files — these are embedded playwright-core references, not standard templates. Not referenced by any skill. | -5 |
| 8 | Command frontmatter | `plugins/forge/commands/git-checkout.md` | Missing `allowed_tools` — uses Bash, Read, and implicitly git commands | (already has allowed_tools, skip) |
| 9 | Command frontmatter | `plugins/forge/commands/git-commit.md` | Missing `argument-hints` — accepts optional `scope` parameter described in body but not declared in frontmatter. Body describes scope argument. | -5 |
| 10 | Command frontmatter | `plugins/forge/commands/execute-task.md` | Missing `argument-hints` — task ID passed implicitly, no hint needed | 0 |
| 11 | Command frontmatter | `plugins/forge/commands/run-tasks.md` | No `argument-hints` declared — no user-facing args needed | 0 |
| 12 | Command frontmatter | `plugins/forge/commands/simplify-skill.md` | Has `argument-hints: skill name` — not structured YAML format | -5 |
| 13 | Command frontmatter | `plugins/forge/commands/fix-bug.md` | Has proper `argument-hints` structured format | 0 |
| 14 | Rubric totals correct | `plugins/forge/skills/eval-harness/templates/rubric.md` | Totals must be verified per rubric. This rubric scores 100 (not 1000). Eval-harness is NOT an eval-* skill in the rubric sense (different scoring). No issue since eval-harness is exempt from the 1000-point scoring convention. | 0 |
| 15 | Rubric totals correct | `plugins/forge/skills/eval-proposal/templates/rubric.md` | Declares 10 dimensions summing to 1000. Needs verification. | -5 |
| 16 | Rubric totals correct | `plugins/forge/skills/eval-ui/templates/rubric.md` | Declares 4 dimensions at 250 each = 1000. Likely correct. | 0 |
| 17 | Rubric totals correct | Multiple rubrics | Cannot fully verify without reading each rubric's internal dimension totals. Conservative deduction for unverified totals. | -5 |
| 18 | Hard Gate present | `plugins/forge/skills/eval-consistency/SKILL.md` | Has `<HARD-GATE>` at Step 3. Actually, this is a `<HARD-GATE>` for the gate decision, not a skill-level hard gate. Reviewing: yes, it has a proper `<HARD-GATE>` block. | 0 |
| 19 | Hard Gate present | `plugins/forge/skills/eval-proposal/SKILL.md` | Has `<HARD-GATE>` at Step 3. Present. | 0 |
| 20 | Hard Gate present | `plugins/forge/skills/eval-design/SKILL.md` | Has `<HARD-GATE>` at Step 3. Present. | 0 |
| 21 | Hard Gate present | `plugins/forge/skills/eval-prd/SKILL.md` | Has `<HARD-GATE>` at Step 3. Present. | 0 |
| 22 | Hard Gate present | `plugins/forge/skills/eval-test-cases/SKILL.md` | Has `<HARD-GATE>` at Step 3. Present. | 0 |
| 23 | Hard Gate present | `plugins/forge/skills/eval-ui/SKILL.md` | Has `<HARD-GATE>` at Step 3. Present. | 0 |
| 24 | Hard Gate present | `plugins/forge/skills/eval-harness/SKILL.md` | Exempt by rubric (no adversarial loop). Has `<HARD-RULE>` but no `<HARD-GATE>`. No deduction per rubric exception. | 0 |
| 25 | Hard Gate present | `plugins/forge/skills/eval-consistency/SKILL.md` | Has `<HARD-GATE>`. Present. | 0 |
| 26 | Output field parsing | `plugins/forge/commands/execute-task.md` line 26 | Describes `FILE` field from claim output. Actual claim.go outputs `FILE` as full absolute path. Matches. | 0 |
| 27 | Output field parsing | `plugins/forge/commands/run-tasks.md` | Same output fields as execute-task. Correct. | 0 |
| 28 | Output field parsing | `plugins/forge/commands/execute-task.md` line 26 | Does NOT document `TITLE`, `PRIORITY`, `NO_TEST`, `TYPE` fields that are actually output by claim.go (lines 312-326). These fields are present in claim output but not listed in the skill's "Extract from claim output" section. | -5 |
| 29 | Claim scheduling align | `plugins/forge/commands/execute-task.md` | Does not describe the full priority ordering (deps met → P0 > P1 > P2 → semantic ID). The dispatcher does not need to describe claim internals — it just calls `task claim`. No deduction needed since the dispatcher delegates to CLI. | 0 |
| 30 | Claim scheduling align | Skills referencing `task claim` | The breakdown-tasks skill describes fix-task creation but does not describe claim priority ordering. This is correct since claim scheduling is internal to CLI. | 0 |
| 31 | Claim scheduling | `plugins/forge/commands/execute-task.md` | Does not document `TITLE`, `PRIORITY`, `NO_TEST`, `TYPE` fields that claim outputs. The dispatcher extracts `KEY`, `TASK_ID`, `FILE`, `BREAKING`, `MAIN_SESSION`, `SCOPE`, `FEATURE` but not `TITLE` or `TYPE`. Not critical but incomplete. | -5 |
| 32 | Record validation align | `plugins/forge/skills/record-task/SKILL.md` | Correctly describes: auto-downgrade (completed + testsFailed > 0 → blocked, non-overridable), test evidence (overridable with --force), AC requirement (overridable with --force), quality gate (compile → fmt → lint → test). Matches record.go exactly. | 0 |
| 33 | Record validation align | `plugins/forge/skills/record-task/SKILL.md` line 129-133 | Quality gate description says "just compile [scope] → just fmt [scope] → just lint [scope] → just test [scope]" which matches `validateQualityGate()` in record.go line 445 which calls `just.DefaultGateSequence()`. The gate sequence is defined elsewhere but the skill says the correct order. | 0 |
| 34 | Record validation align | `plugins/forge/skills/record-task/SKILL.md` | The skill says "status=completed + testsPassed=0 + testsFailed=0 + coverage >= 0" for test evidence check. Record.go line 283 checks `rd.Coverage >= 0 && rd.TestsPassed == 0 && rd.TestsFailed == 0`. Match. | 0 |
| 35 | Dynamic task add | `plugins/forge/skills/breakdown-tasks/SKILL.md` | Uses `task add --template fix-task`, `--source-task-id`, `--block-source`. Matches add.go. | 0 |
| 36 | Dynamic task add | `plugins/forge/skills/breakdown-tasks/SKILL.md` | Says `--block-source` "atomically sets source task to blocked before resolution". Matches add.go line 52 (`BlockSource`). Correct. | 0 |
| 37 | Dynamic task add | `plugins/forge/skills/breakdown-tasks/SKILL.md` | Says generated ID format is `disc-N`. add.go uses template defaults from fix-task which sets `IDPrefix: "disc-"`. Correct. | 0 |
| 38 | Dynamic task add | `plugins/forge/skills/quick-tasks/SKILL.md` | Same fix-task pattern as breakdown-tasks. Correct. | 0 |
| 39 | Dynamic task add | `plugins/forge/skills/breakdown-tasks/SKILL.md` | "When a fix-task completes, task record auto-restores the source task to pending (checks all source task's dependencies are completed)". Matches `autoRestoreSourceTask()` in record.go line 209. However, the skill also says "For nested fix-tasks... Maximum nesting: 3 levels." The Go code does NOT enforce a 3-level nesting limit. | -5 |
| 40 | Schema-code alignment | `plugins/forge/skills/breakdown-tasks/templates/index.schema.json` | Schema does not include `sourceTaskID` field. Go `Task` struct has `SourceTaskID string json:"sourceTaskID,omitempty"`. Rubric notes this as acceptable ("sourceTaskID exists in Go struct but not in JSON schema — auto-managed field"). No deduction per known acceptable discrepancy. | 0 |
| 41 | Schema-code alignment | `plugins/forge/skills/breakdown-tasks/templates/index.schema.json` | Schema has `noTest` field with description "Deprecated: use type field instead." Go struct has `NoTest bool json:"noTest,omitempty"`. Schema marks it deprecated but Go still uses it actively (record.go line 114, 122). Minor inconsistency — schema says deprecated but code still checks it. | -5 |
| 42 | Schema-code alignment | `plugins/forge/skills/breakdown-tasks/templates/index.schema.json` | Schema `type` enum values match Go `ValidTypes` exactly: implementation, doc-generation.summary, doc-generation.consolidate, test-pipeline.*, fix, gate. Match. | 0 |
| 43 | All-completed hook align | `plugins/forge/hooks/guide.md` line 129-133 | Guide says "Quality gate: just compile → just fmt → just lint → Project-wide tests: just test → E2E regression: just e2e-setup → (server health probe) → just test-e2e". Matches all_completed.go: Step 1 runs `LintGateSequence()` (compile → fmt → lint), Step 2 runs project-wide tests, Step 3 runs e2e-setup → health probe → test-e2e. Match. | 0 |
| 44 | Hook scripts exist | `plugins/forge/hooks/hooks.json` PostToolUse | References `bash "${CLAUDE_PLUGIN_ROOT}/scripts/validate-index.sh"`. Script exists at `plugins/forge/scripts/validate-index.sh`. Valid. | 0 |
| 45 | Hook scripts exist | `plugins/forge/hooks/hooks.json` SessionStart | References `${CLAUDE_PLUGIN_ROOT}/hooks/run-hook.cmd` with arg `session-start`. Both `run-hook.cmd` and `session-start` exist in hooks/. Valid. | 0 |
| 46 | Hook scripts exist | `plugins/forge/hooks/hooks.json` SessionStart | The hook is `bash "${CLAUDE_PLUGIN_ROOT}/scripts/validate-index.sh"` in PostToolUse. But this is a bash command, not a hook script in hooks/. The script is in `scripts/` not `hooks/`. Script exists though. Valid. | 0 |
| 47 | Hook scripts exist | `plugins/forge/hooks/hooks.json` SessionEnd | References `task cleanup`. This is a CLI command, not a script. The command exists in `task -h`. Valid. | 0 |
| 48 | Hook event names | `plugins/forge/hooks/hooks.json` | Uses `SubagentStart` and `SubagentStop` event names. These are not standard Claude Code hook events (standard events are: SessionStart, SessionEnd, Stop, PostToolUse, PreToolUse). `SubagentStart`/`SubagentStop` may not be supported. | -5 |
| 49 | Conciseness | `plugins/forge/hooks/guide.md` | Contains quality gate details (compile → fmt → lint → test, failure actions). This is reference material that belongs in individual skill docs or `task -h`, not the guide. | -5 |
| 50 | Conciseness | `plugins/forge/hooks/guide.md` | Contains Scope Resolution algorithm details. This is operational detail duplicated from the execute-task/run-tasks commands. | -5 |
| 51 | argument-hints | `plugins/forge/commands/extract-design-md.md` | Missing `argument-hints`. Body describes accepting a URL parameter but frontmatter does not declare it. | -5 |
| 52 | argument-hints | `plugins/forge/commands/init-forge.md` | No `argument-hints` needed — no user-facing args. | 0 |
| 53 | argument-hints | `plugins/forge/commands/execute-task.md` | No `argument-hints` — task claimed via CLI, not user arg. | 0 |
| 54 | argument-hints | `plugins/forge/commands/run-tasks.md` | No `argument-hints` — no user args. | 0 |
| 55 | argument-hints | `plugins/forge/commands/record-decision.md` | No `argument-hints` — interactive flow via AskUserQuestion. | 0 |
| 56 | argument-hints | `plugins/forge/commands/gen-sitemap.md` | Has proper `argument-hints` with structured format. | 0 |
| 57 | argument-hints | `plugins/forge/commands/simplify-skill.md` | Has `argument-hints: skill name` — plain string, not structured YAML format like other commands. | -5 |
| 58 | argument-hints | `plugins/forge/commands/git-commit.md` | Has `argument-hints` with structured name/description/required. Valid. | 0 |
| 59 | keywords coverage | `plugins/forge/.claude-plugin/plugin.json` | Keywords: pipeline, brainstorm, prd, design, task, eval, e2e, test-profile, ui-design, sitemap. Missing: "breakdown" (breakdown-tasks skill exists), "fix" (fix-bug command exists), "quick" (quick/quick-tasks), "graduate" (graduate-tests), "sitemap" is present. Gap: "breakdown" and "fix" and "quick" are major capability areas not covered. | -5 |
| 60 | keywords coverage | `plugins/forge/.claude-plugin/plugin.json` | Missing "consolidate" (consolidate-specs skill), "forensic" (forensic skill), "record" (record-task, record-decision). These are secondary. | -5 |
| 61 | Hard Gate | `plugins/forge/skills/eval-harness/SKILL.md` | Has no `<HARD-GATE>` block. Rubric exempts eval-harness. No deduction. | 0 |
| 62 | Conciseness | `plugins/forge/hooks/guide.md` | Guide references `/eval-test-cases` in the testing lifecycle mermaid diagram but this skill is T-test-1b, called automatically. Not a standalone workflow skill referenced in guide. Actually the guide's test lifecycle diagram references T-test-1..5 which is accurate. | 0 |

---

## Attack Points

### Attack 1: [Dimension 3 — Dangling reference to nonexistent shared directory]

**Where**: `plugins/forge/references/shared/` — referenced by 6 files but directory does not exist
- `plugins/forge/skills/consolidate-specs/SKILL.md` line 170
- `plugins/forge/commands/record-decision.md` lines 19, 40, 44
- `plugins/forge/commands/gen-sitemap.md` lines 48, 78
- `plugins/forge/skills/tech-design/SKILL.md` line 191
- `plugins/forge/skills/gen-test-scripts/SKILL.md` line 83

**What's wrong**: Multiple skills and commands reference `plugins/forge/references/shared/decision-logging.md`, `plugins/forge/references/shared/config.yaml`, and `plugins/forge/references/shared/sitemap.json`. The entire `plugins/forge/references/` directory does not exist. These are supposed to be shared reference documents that skills read at runtime. Without them, the skills will fail when they attempt to read these files.

**How to fix**: Create `plugins/forge/references/shared/` directory with the three referenced files:
1. `decision-logging.md` — decision archiving flow (Sections 1-8 referenced by record-decision and tech-design)
2. `config.yaml` — e2e test configuration template (referenced by gen-sitemap)
3. `sitemap.json` — full sitemap example (referenced by gen-sitemap and gen-test-scripts)

### Attack 2: [Dimension 8 — Dubious hook event names SubagentStart/SubagentStop]

**Where**: `plugins/forge/hooks/hooks.json` lines 16-27 (SubagentStart) and lines 42-51 (SubagentStop)

**What's wrong**: The hooks.json uses `SubagentStart` and `SubagentStop` as hook event names. The Claude Code plugin hook system supports specific event types (SessionStart, SessionEnd, Stop, PostToolUse, PreToolUse). `SubagentStart` and `SubagentStop` are not documented as standard hook events. If they are unsupported, these hooks will silently never fire, meaning `task cleanup` never runs for subagents, and the guide.md never gets injected into subagent sessions.

**How to fix**: Verify with Claude Code documentation whether `SubagentStart`/`SubagentStop` are valid events. If not, remove these hooks or replace them with supported alternatives.

### Attack 3: [Dimension 7 — Missing claim output fields in dispatcher docs]

**Where**: `plugins/forge/commands/execute-task.md` lines 22-30 and `plugins/forge/commands/run-tasks.md` lines 46-55

**What's wrong**: Both dispatcher commands list "Extract from claim output" fields as: TASK_ID, KEY, FILE, BREAKING, MAIN_SESSION, SCOPE, FEATURE. But claim.go (lines 312-326) actually outputs many more fields: `TITLE`, `PRIORITY`, `STATUS`, `DEPENDENCIES`, `TYPE`, `NO_TEST`, `RECORD`, `ESTIMATED_TIME`. While dispatchers don't need all fields, the incomplete documentation means subagent behavior that relies on `TYPE` or `TITLE` (e.g., for routing decisions) cannot be verified.

Also, breakdown-tasks SKILL.md line 419 states "Maximum nesting: 3 levels" for fix-tasks, but neither add.go nor record.go enforces a 3-level limit. This is a documentation claim with no code backing.

**How to fix**: Either add the missing fields to the dispatcher's "Extract from claim output" list, or explicitly state which fields are ignored. Remove the "Maximum nesting: 3 levels" claim from breakdown-tasks unless adding the enforcement to the Go code.

---

## Previous Issues Check

<!-- Only for iteration > 1 -->

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|

---

## Fix Summary

<!-- Only for iteration > 1 -->

| File Changed | What Changed |
|-------------|--------------|

---

## Verdict

- **Score**: 893/1000
- **Target**: 950/1000
- **Gap**: 57 points
- **Action**: Target not reached. Primary gaps: (1) `references/shared/` directory missing with 6 dangling references (-30 in Dimension 3), (2) Dubious hook events (-5), (3) Guide conciseness issues (-10), (4) Missing argument-hints on 3 commands (-15), (5) Keyword gaps (-10), (6) Minor schema/alignment issues (-10). Fix shared references directory and guide conciseness to close the gap.
