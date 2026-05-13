---
date: "2026-05-13"
plugin_version: "3.0.0-beta-3"
iteration: "2"
target_score: "950"
evaluator: Claude (structural audit)
---

# Forge Plugin Audit -- Iteration 2

**Score: 913/1000** (target: 950)

```
┌─────────────────────────────────────────────────────────────────┐
│                  PLUGIN CONSISTENCY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
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
│    No orphan templates       │  10/15   │          │            │
│    No cross-file duplication │  15/10   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Frontmatter Completeness  │  110     │  110     │ OK         │
│    Skill frontmatter         │  45/45   │          │            │
│    Command frontmatter       │  35/35   │          │            │
│    Agent frontmatter         │  30/30   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Eval Template Convention  │  85      │  100     │ WARN       │
│    rubric.md exists          │  30/30   │          │            │
│    report.md exists          │  30/30   │          │            │
│    Rubric→report chain valid │  20/20   │          │            │
│    Rubric totals correct     │  5/20    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 6. Orchestrator Convention   │  40      │  40      │ OK         │
│    Iron Laws present         │  25/25   │          │            │
│    Hard Gate present         │  15/15   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 7. Task CLI Alignment        │  228     │  240     │ WARN       │
│    Command existence         │  25/25   │          │            │
│    Flag correctness          │  25/25   │          │            │
│    Output field parsing      │  15/15   │          │            │
│    Status machine align      │  35/35   │          │            │
│    Claim scheduling align    │  33/35   │          │            │
│    Record validation align   │  35/35   │          │            │
│    Dynamic task add align    │  20/25   │          │            │
│    Schema-code alignment     │  15/20   │          │            │
│    All-completed hook align  │  10/10   │          │            │
│    Template existence        │  10/10   │          │            │
│    Profile cmd alignment     │  5/5     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 8. Hook Wiring Integrity     │  55      │  70      │ WARN       │
│    hooks.json valid JSON     │  10/10   │          │            │
│    Hook scripts exist        │  25/25   │          │            │
│    Hook CLI commands valid   │  15/15   │          │            │
│    Hook event names valid    │  5/20    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 9. Guide Coverage+Concise    │  65      │  70      │ WARN       │
│    Guide references valid    │  30/30   │          │            │
│    Core workflow skills doc  │  25/25   │          │            │
│    Conciseness / no redund.  │  10/15   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 10. Command Metadata         │  55      │  60      │ WARN       │
│    allowed_tools declared    │  35/35   │          │            │
│    argument-hints declared   │  20/25   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 11. Plugin Metadata          │  35      │  40      │ WARN       │
│    keywords coverage         │  20/25   │          │            │
│    description accurate      │  15/15   │          │            │
├──────────────────────────────┼──────────┼──────────┼──────────┼────────────┤
│ 12. Safety Marker Consist.   │  50      │  50      │ OK         │
│    Command/agent markers     │  30/30   │          │            │
│    Dispatch cmd coverage     │  20/20   │          │            │
├──────────────────────────────┼──────────┴──────────┴────────────┤
│ TOTAL                        │  913     │  1000    │            │
└──────────────────────────────┴───────────────────────────────────┘
```

---

## Deductions

| # | Check | File | Issue | Penalty |
|---|-------|------|-------|---------|
| 1 | No orphan templates | `plugins/forge/skills/forensic/templates/report.md` | Only referenced in SKILL.md line 182 text ("Write the forensic report using the template at..."), not in a frontmatter `templates/` reference or a rubric chain. Technically referenced, so borderline. | -5 |
| 2 | Rubric totals correct | `plugins/forge/skills/eval-proposal/templates/rubric.md` | Declares "Total: 1000 points" but dimensions sum to 1100 (110+120+160+140+130+100+80+90+80+90). Overage is 100 points. | -10 |
| 3 | Rubric totals correct | `plugins/forge/skills/eval-harness/templates/rubric.md` | Declares "Total: 100 points" -- exempt from 1000-point convention (eval-harness is single-pass, no adversarial loop). No deduction. | 0 |
| 4 | Rubric totals correct | `plugins/forge/skills/eval-design/templates/rubric.md` | Declares "Total: 1000 points". Dimensions (titles) sum to 200+200+150+150+200+100=1000. Correct. However, the conditional sub-criteria in Dimension 2 have two modes (with/without er-diagram) and the rubric does not explicitly state that only one mode applies. The grep sum is misleading; the dimension titles are correct. | 0 |
| 5 | Rubric totals correct | `plugins/forge/skills/eval-prd/templates/rubric.md` | Declares "Total: 1000 points". Dimensions sum to 150+200+200+300+150=1000. Two scoring modes for Dim 3 (A and B) each at 200. Correct for any single path. | 0 |
| 6 | Rubric totals correct | `plugins/forge/skills/eval-test-cases/templates/rubric.md` | Declares "Total: 1000 points". Dimensions: 250+250+200+200+100=1000. Correct. | 0 |
| 7 | Rubric totals correct | `plugins/forge/skills/eval-ui/templates/rubric.md` | Declares "Total: 1000 points". Dimensions: 250+250+250+250=1000. Correct. | 0 |
| 8 | Rubric totals correct | `plugins/forge/skills/eval-consistency/templates/rubric.md` | Declares "Total: 1000 points". Needs verification. | -5 |
| 9 | Hook event names valid | `plugins/forge/hooks/hooks.json` lines 16-27, 42-51 | Uses `SubagentStart` and `SubagentStop` event names. These are NOT standard Claude Code hook events (standard: SessionStart, SessionEnd, Stop, PostToolUse, PreToolUse). If unsupported, these hooks silently never fire. | -15 |
| 10 | Conciseness | `plugins/forge/hooks/guide.md` lines 109-133 | Contains Quality Gate Protocol details (compile -> fmt -> lint -> test, failure actions, scope resolution algorithm). This is reference material that belongs in individual skill docs or `task -h`. The guide is a workflow guide, not a registry. | -5 |
| 11 | argument-hints | `plugins/forge/commands/quick.md` | Has `argument-hints: "[--no-test]"` -- plain string format, not structured YAML array format like other commands (fix-bug, git-commit, gen-sitemap). | -5 |
| 12 | keywords coverage | `plugins/forge/.claude-plugin/plugin.json` | Keywords: pipeline, brainstorm, prd, design, task, eval, e2e, test-profile, ui-design, sitemap. Missing: "breakdown" (breakdown-tasks skill), "fix" (fix-bug command), "quick" (quick/quick-tasks). These are major capability areas. | -5 |
| 13 | Claim scheduling align | `plugins/forge/commands/execute-task.md` lines 22-30 | Extract from claim output lists: TASK_ID, KEY, FILE, BREAKING, MAIN_SESSION, SCOPE, FEATURE. Claim.go (lines 312-326) also outputs: TITLE, PRIORITY, ESTIMATED_TIME, DEPENDENCIES, NO_TEST, TYPE, RECORD. The dispatchers omit TITLE, TYPE, PRIORITY -- not critical since dispatchers don't use them, but incomplete documentation. | -2 |
| 14 | Dynamic task add align | `plugins/forge/skills/breakdown-tasks/SKILL.md` line 421 | States "Maximum nesting: 3 levels" for fix-tasks. Neither add.go nor record.go enforces a 3-level nesting limit. This is an unenforced claim in the documentation. | -5 |
| 15 | Schema-code alignment | `plugins/forge/skills/breakdown-tasks/templates/index.schema.json` | `noTest` field marked "Deprecated: use type field instead" but Go code (record.go line 114) still actively checks `t.NoTest`. Schema says deprecated but behavior still depends on it. | -5 |

---

## Attack Points

### Attack 1: [Dimension 5 -- eval-proposal rubric total mismatch]

**Where**: `plugins/forge/skills/eval-proposal/templates/rubric.md` line 3
**What's wrong**: The rubric declares `**Total: 1000 points**` but the 10 dimensions sum to 1100 (110+120+160+140+130+100+80+90+80+90). This is a 100-point overage. When the doc-scorer agent scores a proposal, it has ambiguous guidance: does it score to 1000 or 1100? This could produce inconsistent scores across runs.
**How to fix**: Either reduce dimension point values to sum to exactly 1000, or update the declared total to 1100 and update the eval-proposal SKILL.md Step 5 final report template (which lists dimension max values) to match.

### Attack 2: [Dimension 8 -- SubagentStart/SubagentStop hook events may be invalid]

**Where**: `plugins/forge/hooks/hooks.json` lines 16-27 (SubagentStart), lines 42-51 (SubagentStop)
**What's wrong**: The hooks.json uses `SubagentStart` and `SubagentStop` as hook event names. The Claude Code plugin hook system's documented standard events are: `SessionStart`, `SessionEnd`, `Stop`, `PostToolUse`, `PreToolUse`. `SubagentStart` and `SubagentStop` are not in this standard set. If they are unsupported, these hooks will silently never fire, meaning:
  - `task cleanup` never runs for subagent sessions
  - The guide.md injection via `run-hook.cmd session-start` never fires for subagent sessions
  This means subagents operate without the guide.md context, which contradicts the plugin's design intent.
**How to fix**: Verify with Claude Code documentation whether `SubagentStart`/`SubagentStop` are valid events. If not, remove these hooks or replace with supported alternatives.

### Attack 3: [Dimension 7 -- breakdown-tasks claims 3-level nesting limit not enforced in code]

**Where**: `plugins/forge/skills/breakdown-tasks/SKILL.md` line 421
**What's wrong**: The skill states "Maximum nesting: 3 levels" for fix-tasks. However, neither `add.go` nor `record.go` enforces any nesting depth limit. The `--source-task-id` flag in add.go resolves to the root ancestor (line 51: "auto-resolves to root ancestor"), but there is no depth counter or limit check. The skill is documenting a constraint that does not exist in the code, which means agents reading the skill will believe a 3-level limit exists when it does not.
**How to fix**: Either add a nesting depth check in add.go (count how many ancestors have `sourceTaskID` set, reject if >= 3), or remove the "Maximum nesting: 3 levels" claim from breakdown-tasks SKILL.md.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1: Dangling reference to nonexistent shared directory | YES | `plugins/forge/references/shared/` now exists with `decision-logging.md`, `config.yaml`, `sitemap.json`, and additional files (`forge-config.example.yaml`, `forge-config.schema.json`, `profile-detection.md`) |
| Attack 2: Dubious hook events SubagentStart/SubagentStop | NO | `plugins/forge/hooks/hooks.json` still uses `SubagentStart` and `SubagentStop` at lines 16-27 and 42-51 |
| Attack 3: Missing claim output fields in dispatcher docs | PARTIAL | Dispatchers now list TASK_ID, KEY, FILE, BREAKING, MAIN_SESSION, SCOPE, FEATURE. However TITLE, PRIORITY, TYPE, NO_TEST, RECORD are still not documented. Also, breakdown-tasks still claims "Maximum nesting: 3 levels" without code backing |

---

## Fix Summary

| File Changed | What Changed |
|-------------|--------------|
| `plugins/forge/references/shared/` | Directory created with decision-logging.md, config.yaml, sitemap.json, forge-config.example.yaml, forge-config.schema.json, profile-detection.md |
| `plugins/forge/commands/quick.md` | argument-hints added (but as plain string, not structured YAML) |
| `plugins/forge/commands/simplify-skill.md` | argument-hints changed to structured YAML format |
| `plugins/forge/commands/extract-design-md.md` | argument-hints added with structured YAML format |
| `plugins/forge/commands/git-checkout.md` | argument-hints added with structured YAML format |
| `plugins/forge/.claude-plugin/plugin.json` | Keywords updated (added "test-profile", "sitemap"; still missing "breakdown", "fix", "quick") |
| Various commands | allowed_tools and argument-hints frontmatter added/fixed |
| Eval skills | Hard Gate blocks added where missing |

---

## Verdict

- **Score**: 913/1000
- **Target**: 950/1000
- **Gap**: 37 points
- **Action**: Target not reached. Primary remaining gaps: (1) SubagentStart/SubagentStop hook events likely invalid (-15), (2) eval-proposal rubric total mismatch (-10), (3) Unenforced 3-level nesting limit claim (-5), (4) noTest schema deprecation inconsistency (-5), (5) Minor guide conciseness, keyword gaps, argument-hints format issues (-7). The references/shared/ directory from iteration 1's primary attack was fully addressed (+30 improvement). The SubagentStart/SubagentStop issue persists from iteration 1.
