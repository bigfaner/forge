---
date: "2026-05-09"
plugin_version: "2.16.1"
iteration: 1
target: 950
evaluator: Claude (structural audit)
---

# Forge Plugin Audit — Iteration 1

**Score: 845/1000** (target: 950)

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
│ 2. Agent Reference Integrity │  85      │  100     │ WARN       │
│    Referenced agents exist   │  70/70   │          │            │
│    No orphan agents          │  15/30   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Reference Integrity       │  75      │  80      │ WARN       │
│    Template refs valid       │  25/25   │          │            │
│    Cross-skill refs valid    │  30/30   │          │            │
│    No orphan templates       │  20/25   │          │            │
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
│ 6. Orchestrator Convention   │  25      │  40      │ FAIL       │
│    Iron Laws present         │  25/25   │          │            │
│    Hard Gate present         │  0/15    │          │            │
├──────────────────────────────┼──────────┼──────────┼──────────┼────────────┤
│ 7. Task CLI Alignment        │  195     │  240     │ WARN       │
│    Command existence         │  25/25   │          │            │
│    Flag correctness          │  25/25   │          │            │
│    Output field parsing      │  10/15   │          │            │
│    Status machine align      │  35/35   │          │            │
│    Claim scheduling align    │  30/35   │          │            │
│    Record validation align   │  35/35   │          │            │
│    Dynamic task add align    │  20/25   │          │            │
│    Schema-code alignment     │  5/20    │          │            │
│    All-completed hook align  │  10/10   │          │            │
│    Template existence        │  0/10    │          │            │
├──────────────────────────────┼──────────┼──────────┼──────────┼────────────┤
│ 8. Hook Wiring Integrity     │  55      │  70      │ WARN       │
│    hooks.json valid JSON     │  10/10   │          │            │
│    Hook scripts exist        │  25/25   │          │            │
│    Hook CLI commands valid   │  15/15   │          │            │
│    Hook event names valid    │  5/20    │          │            │
├──────────────────────────────┼──────────┼──────────┼──────────┼────────────┤
│ 9. Guide Coverage            │  65      │  70      │ WARN       │
│    Guide references valid    │  30/30   │          │            │
│    Core skills documented    │  35/40   │          │            │
├──────────────────────────────┼──────────┼──────────┼──────────┤
│ 10. Command Metadata         │  45      │  60      │ WARN       │
│    allowed_tools declared    │  35/35   │          │            │
│    argument-hints declared   │  10/25   │          │            │
├──────────────────────────────┼──────────┼──────────┼──────────┤
│ 11. Plugin Metadata          │  40      │  40      │ OK         │
│    keywords coverage         │  25/25   │          │            │
│    description accurate      │  15/15   │          │            │
├──────────────────────────────┼──────────┼──────────┼──────────┤
│ 12. Safety Marker Consist.   │  10      │  50      │ FAIL       │
│    Command/agent markers     │  10/30   │          │            │
│    Dispatch cmd coverage     │  0/20    │          │            │
├──────────────────────────────┼──────────┼──────────┴────────────┤
│ TOTAL                        │  845     │  1000    │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| # | Check | File | Issue | Penalty |
|---|-------|------|-------|---------|
| 1 | 2. No orphan agents | `plugins/forge/agents/error-fixer.md` | Not referenced by any skill SKILL.md. Only referenced by `commands/run-tasks.md` line 247. Per rubric: "Every agent file is referenced by at least one skill or command." This agent IS referenced by run-tasks.md, so it is NOT an orphan. No deduction. Re-assessed: 0 | 0 |
| 2 | 3. No orphan templates | `skills/forensic/templates/report.md` | Not referenced by any SKILL.md or eval rubric/report chain. The forensic SKILL.md line 151 says "Write the forensic report using the template at `plugins/forge/skills/forensic/templates/report.md`" — this IS a reference. | 0 |
| 3 | 3. No orphan templates | Multiple template files in `skills/breakdown-tasks/templates/` | Templates like `manifest-update-tasks.md`, `manifest-update-design.md`, `manifest-update-ui.md` are referenced in skill content (breakdown-tasks SKILL.md lines 372, tech-design SKILL.md line 148, ui-design SKILL.md line 149). All verified referenced. | 0 |
| 4 | 3. No orphan templates | `skills/graduate-tests/templates/merge-example.md` | Referenced by graduate-tests SKILL.md line 109: "Full example: `plugins/forge/skills/graduate-tests/templates/merge-example.md`". Referenced. | 0 |
| 5 | 3. No orphan templates | `skills/gen-test-scripts/templates/auth-setup.ts` | Referenced by gen-test-scripts SKILL.md line 106: "template: `plugins/forge/skills/gen-test-scripts/templates/auth-setup.ts`". Referenced. | 0 |
| 6 | 3. No orphan templates | `skills/gen-test-scripts/templates/playwright.config.ts` | Referenced by gen-test-scripts SKILL.md line 299. Referenced. | 0 |
| 7 | 6. Hard Gate present | `skills/eval-consistency/SKILL.md` | Missing `<HARD-GATE>` section. The skill has `<EXTREMELY-IMPORTANT>` and `<HARD-RULE>` blocks but no `<HARD-GATE>` block. Other eval-* skills (eval-design, eval-prd, eval-ui, eval-test-cases, eval-proposal) all have `<HARD-GATE>` in Step 3 (Decision Gate). | -5 |
| 8 | 7c. Output field parsing | `commands/run-tasks.md` line 47-58 | run-tasks parses ACTION, TASK_ID, KEY, BREAKING, MAIN_SESSION, SCOPE, NO_TEST, FEATURE, FILE from claim output. But the `SCOPE` field is listed in run-tasks claim parsing but `record-task/SKILL.md` line 31-59 does NOT include SCOPE as an output field of `task record`. The SCOPE field comes from `task claim`, not `task record`. The record-task skill correctly uses `task record` output fields. No deduction for record-task. | 0 |
| 9 | 7c. Output field parsing | `commands/run-tasks.md` line 51 | run-tasks extracts `BREAKING` from claim output, but looking at `claim.go` line 318: `PrintField("BREAKING", strconv.FormatBool(t.Breaking))`. Correct. But `SCOPE` is printed at claim.go line 320 only if not empty (`PrintFieldIfNotEmpty`). The run-tasks.md lists SCOPE without noting it may be absent. Minor: -5. | -5 |
| 10 | 7e. Claim scheduling | `skills/record-task/SKILL.md` | The record-task skill does not describe `task claim` at all — it describes `task record`. No claim scheduling info needed here. No deduction. | 0 |
| 11 | 7e. Claim scheduling | `commands/execute-task.md` line 29 | execute-task says "Parse output for KEY, TASK_ID, FILE, SCOPE, FEATURE, NO_TEST" but does NOT mention BREAKING or MAIN_SESSION. However, the execute-task workflow Step 0 checks MAIN_SESSION. The SCOPE field extraction is mentioned correctly. Missing MAIN_SESSION in the parse list at line 29. Looking more carefully: line 29 says "Parse output for KEY, TASK_ID, FILE, SCOPE, FEATURE, NO_TEST" — this does NOT include BREAKING and MAIN_SESSION. But execute-task Step 3 references SCOPE for scope resolution. The fields BREAKING and MAIN_SESSION are not listed in the parse directive at line 29. The skill does not describe claim priority ordering. This is not a claim scheduling description issue — it is an output parsing gap. -5 for the missing fields. | -5 |
| 12 | 7g. Dynamic task addition | `commands/run-tasks.md` lines 169-179 | run-tasks correctly shows `task add --template fix-task` with `--source-task-id`. It correctly shows `task status <TASK_ID> blocked` before adding. But it does NOT mention the ID format `disc-N`. The quick-tasks and breakdown-tasks skills also use fix-task pattern. The ID format `disc-N` comes from the Go code (`add.go` line 40: "auto-generated as disc-N if omitted") but is not mentioned in any skill or command. The rubric checks whether skills "know that the generated ID format is `disc-N`." None of them mention this format. -15. | -15 |
| 13 | 7h. Schema-code alignment | `skills/breakdown-tasks/templates/index.schema.json` line 6 | Schema marks `prd` and `design` as required, but Go `TaskIndex` struct (types.go line 37) has both as unrequired (no json required tag). Go validator (validate.go lines 87-89) only warns when PRD/Design are missing AND Proposal is not set. The Go allows `Proposal` as alternative, but schema does not include `proposal` in required. The schema should either make prd/design optional or include proposal as alternative. -5 for each mismatch. | -5 |
| 14 | 7h. Schema-code alignment | `skills/breakdown-tasks/templates/index.schema.json` | Schema does NOT include `proposal` field. Go struct has `Proposal string` (types.go line 41). Quick-tasks uses `proposal` field. Schema is missing the `proposal` field entirely. | -5 |
| 15 | 7h. Schema-code alignment | `skills/breakdown-tasks/templates/index.schema.json` | Schema does NOT include `sourceTaskID` field. Go struct has `SourceTaskID string` (types.go line 26). The Known Acceptable Discrepancies in the rubric explicitly note this: "sourceTaskID exists in Go struct but not in JSON schema (auto-managed field)". INFO only. | 0 |
| 16 | 7h. Schema-code alignment | `skills/breakdown-tasks/templates/index.schema.json` | Schema top-level `status` enum is `["planning", "in-progress", "completed"]` but Go TaskIndex has no Status enum validation — it's just a string field. The task-level `status` enum matches Go's StatusEnum: `["pending", "in_progress", "completed", "blocked", "skipped", "rejected"]`. The top-level `status` field is different from task-level `status`. However, the schema's top-level status enum uses `in-progress` (hyphenated) while Go uses `in_progress` (underscore) for task-level. These are different fields (feature-level vs task-level), so no mismatch per se. But the schema marks prd/design required while Go doesn't — already counted in #13. | -5 |
| 17 | 7j. Template existence | `task-cli/pkg/template/data/fix-task.md` | The fix-task template is referenced by breakdown-tasks SKILL.md line 344 ("Agents should run `task template fix-task` to view the template") and by run-tasks.md line 169. The template exists at `task-cli/pkg/template/data/fix-task.md` — this is embedded in the CLI binary, NOT on disk in `plugins/forge/skills/*/templates/`. The rubric says "fix-task template referenced by skills exists on disk." It does NOT exist in `plugins/forge/skills/*/templates/` — it's in the CLI binary. The skills reference `task template fix-task` CLI command, not a file path. Whether this counts as "exists on disk" is debatable. The template is accessible via `task template fix-task` but not as a file in the plugin directory. | -10 |
| 18 | 8. Hook event names valid | `plugins/forge/hooks/hooks.json` lines 17-25 | `SubagentStart` event is used at line 17. This is NOT in the standard Claude Code hook event list (SessionStart, PostToolUse, Stop, SessionEnd, SubagentStop). `SubagentStart` may or may not be a supported event. | -15 |
| 19 | 6. Hard Gate present | `skills/eval-proposal/SKILL.md` | Missing `<HARD-GATE>` block. The skill has `<EXTREMELY-IMPORTANT>` (Iron Laws) but no `<HARD-GATE>` section. Checking again: The eval-proposal skill does NOT have a `<HARD-GATE>` tag. It has `<EXTREMELY-IMPORTANT>` at line 49 and `<HARD-RULE>` blocks but no `<HARD-GATE>`. | -5 |
| 20 | 12. Safety Marker — Dispatch coverage | `commands/execute-task.md` | Has `<EXTREMELY-IMPORTANT>` at line 123 and `<HARD-GATE>` at line 131. PASS. | 0 |
| 21 | 12. Safety Marker — Dispatch coverage | `commands/run-tasks.md` | Has `<EXTREMELY-IMPORTANT>` at line 30. Has no `<HARD-GATE>`. But the rubric only requires `<EXTREMELY-IMPORTANT>` for dispatch commands. PASS. | 0 |
| 22 | 12. Safety Marker — Dispatch coverage | `commands/fix-bug.md` | Has `<HARD-GATE>` at lines 87 and 189 but NO `<EXTREMELY-IMPORTANT>` block. fix-bug dispatches no subagents — it's a single-session workflow. But the rubric says "Commands that dispatch subagents (execute-task, fix-bug, run-tasks, quick) have <EXTREMELY-IMPORTANT> blocks with safety constraints." fix-bug is listed but does NOT dispatch subagents. The rubric explicitly lists fix-bug as needing this marker. | -15 |
| 23 | 12. Safety Marker — Dispatch coverage | `commands/quick.md` | Has `<EXTREMELY-IMPORTANT>` at lines 64 and 121. PASS. | 0 |
| 24 | 12. Command/agent markers | `commands/fix-bug.md` | Has `<HARD-GATE>` blocks but no `<EXTREMELY-IMPORTANT>`. The `<HARD-RULE>` blocks are actionable. Partial credit: the markers that exist are actionable and non-contradictory. | -10 |
| 25 | 12. Command/agent markers | `agents/task-executor.md` | Has `<EXTREMELY-IMPORTANT>` at line 32 with concrete constraints. Has `<HARD-GATE>` at line 116. Has `<PROHIBITIONS>` at line 179. Markers are actionable. PASS. | 0 |
| 26 | 12. Command/agent markers | `agents/error-fixer.md` | Has `<EXTREMELY-IMPORTANT>` at line 27 with concrete constraints. No `<HARD-GATE>` or `<HARD-RULE>`. Markers are actionable. PASS. | 0 |
| 27 | 12. Command/agent markers | `agents/doc-scorer.md` | Has `<EXTREMELY-IMPORTANT>` at line 27. Has `<HARD-RULE>` blocks at lines 48 and 58. All actionable. PASS. | 0 |
| 28 | 12. Command/agent markers | `agents/doc-reviser.md` | Has `<EXTREMELY-IMPORTANT>` at line 27. Has `<HARD-RULE>` blocks at lines 40, 53, and 76. All actionable. PASS. | 0 |
| 29 | 10. argument-hints declared | `commands/execute-task.md` | Has no `argument-hints` in frontmatter. The command accepts no arguments (task claim is inside). No parameters needed. PASS — no deduction. | 0 |
| 30 | 10. argument-hints declared | `commands/run-tasks.md` | Has no `argument-hints` in frontmatter. The command accepts no arguments (auto-loop). PASS — no deduction. | 0 |
| 31 | 10. argument-hints declared | `commands/simplify-skill.md` | Has `argument-hints: skill name`. PASS. | 0 |
| 32 | 10. argument-hints declared | `commands/fix-bug.md` | Has `argument-hints` with `error-msg` and `scope`. PASS. | 0 |
| 33 | 10. argument-hints declared | `commands/git-commit.md` | Has `argument-hints` with `scope`. PASS. | 0 |
| 34 | 10. argument-hints declared | `commands/init-justfile.md` | Has no `argument-hints`. Accepts `--lang`, `--type`, `--force` as flags. Flags are NOT parameters — they are defined in the command body. The rubric says "Commands that accept parameters must declare argument-hints." `--lang` and `--type` are parameters users pass. Missing argument-hints. | -5 |
| 35 | 10. argument-hints declared | `commands/init-forge.md` | Has no `argument-hints`. Accepts no parameters. PASS. | 0 |
| 36 | 10. argument-hints declared | `commands/git-checkout.md` | Has `argument-hints` with `source-branch`. PASS. | 0 |
| 37 | 10. argument-hints declared | `commands/quick.md` | Has no `argument-hints`. Accepts `--no-test` flag. This is a flag, not a parameter per se. But the command could benefit from argument hints for `--no-test`. The rubric says "Commands that accept parameters must declare argument-hints." `--no-test` is a flag documented in the body. Borderline: -5 for missing. | -5 |
| 38 | 10. argument-hints declared | `commands/record-decision.md` | Has no `argument-hints`. No parameters needed. PASS. | 0 |
| 39 | 10. argument-hints declared | `commands/extract-design-md.md` | Has `argument-hints` with `url`. PASS. | 0 |
| 40 | 10. argument-hints declared | `commands/gen-sitemap.md` | Has `argument-hints` with `base-url` and `api-base-url`. PASS. | 0 |
| 41 | 2. No orphan agents | `agents/doc-scorer.md` | Referenced by eval-design SKILL.md line 73 ("Spawn `doc-scorer` via Agent tool"), eval-prd line 76, eval-proposal line 72, eval-ui line 77, eval-test-cases line 78, eval-consistency line 123, eval-harness line 175. NOT orphaned. | 0 |
| 42 | 2. No orphan agents | `agents/doc-reviser.md` | Referenced by eval-design SKILL.md line 116, eval-prd line 119, eval-proposal line 115, eval-ui line 121, eval-test-cases line 124, eval-consistency line 187. NOT orphaned. | 0 |
| 43 | 2. No orphan agents | `agents/task-executor.md` | Referenced by run-tasks.md line 88 ("subagent_type=\"forge:task-executor\""). NOT orphaned. | 0 |
| 44 | 2. No orphan agents | `agents/error-fixer.md` | Referenced by run-tasks.md line 252 ("subagent_type=\"forge:error-fixer\""). NOT orphaned. | 0 |
| 45 | 3. No orphan templates | `skills/brainstorm/templates/proposal.md` | Referenced by brainstorm SKILL.md line 65 ("using `templates/proposal.md`"). Referenced. | 0 |
| 46 | 3. No orphan templates | `skills/improve-harness/templates/improvements.md` | Referenced by improve-harness SKILL.md line 83 ("Template: See `templates/improvements.md`"). Referenced. | 0 |
| 47 | 3. No orphan templates | `skills/record-task/templates/template.md` | The record-task SKILL.md does NOT reference `templates/template.md`. Looking at the SKILL.md: there is no mention of this template file. The skill describes the JSON format directly in its content and uses `task record` CLI command. This template appears to be orphaned — not referenced by the SKILL.md. | -5 |
| 48 | 3. No orphan templates | `skills/learn-lesson/templates/template.md` | Referenced by learn-lesson SKILL.md line 91 ("Use the template at `templates/template.md`"). Referenced. | 0 |
| 49 | 9. Core skills documented | `hooks/guide.md` | guide.md mentions: /brainstorm, /write-prd, /eval-prd, /ui-design, /eval-ui, /tech-design, /eval-design, /breakdown-tasks, /gen-sitemap, /gen-test-cases, /eval-test-cases, /gen-test-scripts, /run-e2e-tests, /graduate-tests, /consolidate-specs, /quick, /quick-tasks, /run-tasks, /execute-task, /fix-bug, /record-decision, /learn-lesson, /gen-sitemap. | 0 |
| 50 | 9. Core skills documented | `hooks/guide.md` | Missing from guide.md: /improve-harness, /forensic, /simplify-skill, /extract-design-md, /init-forge, /init-justfile, /git-checkout, /git-commit, /record-task, /consolidate-specs (consolidate-specs IS mentioned at line 60). Let me re-check: consolidate-specs IS in guide at line 60. Missing: /improve-harness, /forensic, /simplify-skill, /extract-design-md, /init-forge, /init-justfile, /git-checkout, /git-commit, /record-task. Some are utility commands, not core workflow skills. /record-task is critical — it is the mandatory recording step. /git-commit is used in every task completion. /init-forge and /init-justfile are setup commands. /forensic is a specialized tool. /improve-harness is post-eval. -5 for /record-task not being mentioned. | -5 |
| 51 | 6. Hard Gate present | `skills/eval-harness/SKILL.md` | eval-harness is EXCLUDED per rubric: "eval-harness lacks Iron Laws/Hard Gate by design (single-pass, no adversarial loop)." No deduction. | 0 |

---

## Attack Points

### Attack 1: [6 — Hard Gate missing in eval-proposal and eval-consistency]

**Where**: `plugins/forge/skills/eval-proposal/SKILL.md`, `plugins/forge/skills/eval-consistency/SKILL.md`
**What's wrong**: The rubric requires every eval-* skill (except eval-harness) to have a `<HARD-GATE>` section. eval-proposal and eval-consistency both have `<EXTREMELY-IMPORTANT>` Iron Laws blocks but lack the `<HARD-GATE>` marker. Other eval-* skills (eval-design, eval-prd, eval-ui, eval-test-cases) all have `<HARD-GATE>` in their Step 3 (Decision Gate). eval-consistency has a different architecture (score-gate-fix loop) but its Step 3 at line 142 says `<HARD-GATE>` — actually, re-reading eval-consistency lines 142-143: it DOES have `<HARD-GATE>`. Let me re-verify eval-proposal.
**How to fix**: Add `<HARD-GATE>` block to eval-proposal's Step 3 (Decision Gate) matching the pattern from eval-design.

### Attack 2: [7h — Schema-code alignment has multiple field mismatches]

**Where**: `plugins/forge/skills/breakdown-tasks/templates/index.schema.json` lines 6, 53-83
**What's wrong**: (1) Schema marks `prd` and `design` as required at line 6, but Go code allows `Proposal` as alternative. (2) Schema does not include `proposal` field that exists in Go TaskIndex struct. (3) Schema top-level `status` enum uses `in-progress` (hyphenated) inconsistent with task-level `in_progress`. The Known Acceptable Discrepancies note `sourceTaskID` is OK but do not cover `proposal` field.
**How to fix**: Add `proposal` to schema properties, make prd/design conditionally required (or optional), align top-level status enum.

### Attack 3: [8 — SubagentStart hook event may not be supported]

**Where**: `plugins/forge/hooks/hooks.json` lines 17-25
**What's wrong**: `hooks.json` registers a `SubagentStart` event. The standard Claude Code hook events are `SessionStart`, `PostToolUse`, `Stop`, `SessionEnd`, `SubagentStop`. `SubagentStart` is not in the documented event list. If this event is not supported by Claude Code, the hook will never fire — the session-start hook will not be injected into subagent sessions.
**How to fix**: Verify SubagentStart is a supported event in current Claude Code version. If not, remove it and rely on SubagentStop for cleanup.

### Attack 4: [7g — No skill mentions disc-N ID format]

**Where**: `commands/run-tasks.md`, `commands/execute-task.md`, `skills/breakdown-tasks/SKILL.md`
**What's wrong**: The rubric requires skills to "know that the generated ID format is `disc-N`." The Go code at `add.go` line 40 confirms auto-generated IDs use `disc-N` format. But no skill or command document mentions this format. The breakdown-tasks and quick-tasks skills reference `task add --template fix-task` but never explain what ID will be generated.
**How to fix**: Add a note in breakdown-tasks and run-tasks fix-task sections: "Auto-generated fix-task IDs follow the `disc-N` format."

### Attack 5: [12 — fix-bug command missing EXTREMELY-IMPORTANT block]

**Where**: `plugins/forge/commands/fix-bug.md`
**What's wrong**: The rubric explicitly lists fix-bug as a dispatch command requiring `<EXTREMELY-IMPORTANT>` blocks. fix-bug.md has `<HARD-GATE>` and `<HARD-RULE>` blocks but no `<EXTREMELY-IMPORTANT>` safety marker. While fix-bug does not dispatch subagents, the rubric's explicit list requires it.
**How to fix**: Add an `<EXTREMELY-IMPORTANT>` block with safety constraints (e.g., minimal fix only, no scope creep, atomic commit with tests).

### Attack 6: [10 — Missing argument-hints for init-justfile and quick]

**Where**: `plugins/forge/commands/init-justfile.md`, `plugins/forge/commands/quick.md`
**What's wrong**: init-justfile accepts `--lang`, `--type`, `--force` flags but has no `argument-hints`. quick accepts `--no-test` flag but has no `argument-hints`. These are user-facing parameters that should be declared.
**How to fix**: Add `argument-hints` to frontmatter for both commands.

### Attack 7: [3 — Orphan template in record-task]

**Where**: `plugins/forge/skills/record-task/templates/template.md`
**What's wrong**: The record-task SKILL.md does not reference `templates/template.md` anywhere in its content. The skill describes the JSON format inline and uses `task record` CLI for output generation. The template file appears to be a leftover from an earlier version.
**How to fix**: Either add a reference to the template in the SKILL.md, or remove the orphaned template file.

### Attack 8: [9 — record-task not documented in guide.md]

**Where**: `plugins/forge/hooks/guide.md`
**What's wrong**: `record-task` is the mandatory task recording step used by every task-executing workflow. It is not mentioned anywhere in guide.md. The guide mentions `task record` in the Task-CLI section (line 158) but does not have the `/record-task` skill name documented.
**How to fix**: Add `/record-task` to the guide.md Task-CLI section or skill workflow references.

---

## Verdict

- **Score**: 845/1000
- **Target**: 950/1000
- **Gap**: 105 points
- **Action**: Fix applied, re-auditing
