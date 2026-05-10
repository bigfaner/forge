---
date: "2026-05-10"
plugin_version: "2.18.0"
iteration: "1"
target_score: "950"
evaluator: Claude (structural audit)
---

# Forge Plugin Audit — Iteration 1

**Score: 885/1000** (target: 950)

```
┌─────────────────────────────────────────────────────────────────┐
│                  PLUGIN CONSISTENCY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Directory-Name Alignment  │  40      │  40      │ ✅         │
│    Skill name matches dir    │  25/25   │          │            │
│    Command name matches file │  15/15   │          │            │
├──────────────────────────────┼──────────┼──────────┼──────────┤
│ 2. Agent Reference Integrity │  100     │  100     │ ✅         │
│    Referenced agents exist   │  70/70   │          │            │
│    No orphan agents          │  30/30   │          │            │
├──────────────────────────────┼──────────┼──────────┼──────────┤
│ 3. Reference Integrity       │  75      │  80      │ ⚠️         │
│    Template refs valid       │  25/25   │          │            │
│    Cross-skill refs valid    │  30/30   │          │            │
│    No orphan templates       │  15/15   │          │            │
│    No cross-file duplication │  5/10    │          │            │
├──────────────────────────────┼──────────┼──────────┼──────────┤
│ 4. Frontmatter Completeness  │  110     │  110     │ ✅         │
│    Skill frontmatter         │  45/45   │          │            │
│    Command frontmatter       │  35/35   │          │            │
│    Agent frontmatter         │  30/30   │          │            │
├──────────────────────────────┼──────────┼──────────┼──────────┤
│ 5. Eval Template Convention  │  100     │  100     │ ✅         │
│    rubric.md exists          │  30/30   │          │            │
│    report.md exists          │  30/30   │          │            │
│    Rubric→report chain valid │  20/20   │          │            │
│    Rubric totals correct     │  20/20   │          │            │
├──────────────────────────────┼──────────┼──────────┼──────────┤
│ 6. Orchestrator Convention   │  40      │  40      │ ✅         │
│    Iron Laws present         │  25/25   │          │            │
│    Hard Gate present         │  15/15   │          │            │
├──────────────────────────────┼──────────┼──────────┼──────────┤
│ 7. Task CLI Alignment        │  190     │  240     │ ❌         │
│    Command existence         │  25/25   │          │            │
│    Flag correctness          │  25/25   │          │            │
│    Output field parsing      │  15/15   │          │            │
│    Status machine align      │  35/35   │          │            │
│    Claim scheduling align    │  35/35   │          │            │
│    Record validation align   │  10/35   │          │            │
│    Dynamic task add align    │  25/25   │          │            │
│    Schema-code alignment     │  10/20   │          │            │
│    All-completed hook align  │  5/10    │          │            │
│    Template existence        │  5/10    │          │            │
├──────────────────────────────┼──────────┼──────────┼──────────┤
│ 8. Hook Wiring Integrity     │  65      │  70      │ ⚠️         │
│    hooks.json valid JSON     │  10/10   │          │            │
│    Hook scripts exist        │  25/25   │          │            │
│    Hook CLI commands valid   │  10/15   │          │            │
│    Hook event names valid    │  20/20   │          │            │
├──────────────────────────────┼──────────┼──────────┼──────────┤
│ 9. Guide Coverage+Concise    │  65      │  70      │ ⚠️         │
│    Guide references valid    │  30/30   │          │            │
│    Core workflow skills doc  │  25/25   │          │            │
│    Conciseness / no redund.  │  10/15   │          │            │
├──────────────────────────────┼──────────┼──────────┼──────────┤
│ 10. Command Metadata         │  50      │  60      │ ⚠️         │
│    allowed_tools declared    │  25/35   │          │            │
│    argument-hints declared   │  25/25   │          │            │
├──────────────────────────────┼──────────┼──────────┼──────────┤
│ 11. Plugin Metadata          │  40      │  40      │ ✅         │
│    keywords coverage         │  25/25   │          │            │
│    description accurate      │  15/15   │          │            │
├──────────────────────────────┼──────────┼──────────┼──────────┤
│ 12. Safety Marker Consist.   │  40      │  50      │ ⚠️         │
│    Command/agent markers     │  25/30   │          │            │
│    Dispatch cmd coverage     │  15/20   │          │            │
├──────────────────────────────┼──────────┼──────────┴────────────┤
│ TOTAL                        │  885     │  1000    │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| # | Check | File | Issue | Penalty |
|---|-------|------|-------|---------|
| 1 | 3-No cross-file duplication | plugins/forge/commands/execute-task.md, plugins/forge/agents/task-executor.md | Quality gate failure table (compile→fmt→lint→test) copy-pasted verbatim across 4 files (execute-task, fix-bug, task-executor, error-fixer). The canonical location is guide.md Scope Resolution. | -5 |
| 2 | 7f-Record validation align | plugins/forge/skills/record-task/SKILL.md:130 | Skill says quality gate runs `just compile → just fmt → just lint → just test` before record. Go code (record.go:445) runs `just.DefaultGateSequence()` which is compile→fmt→lint only (no test). The `just test` is not part of the pre-record gate — it is the task executor's responsibility before recording. | -25 |
| 3 | 7f-Record validation align | plugins/forge/skills/record-task/SKILL.md:140 | Skill says `status=completed` + `testsPassed=0` + `testsFailed=0` + `coverage >= 0` is an error (no test evidence). Go code (record.go:283) checks `rd.Coverage >= 0 && rd.TestsPassed == 0 && rd.TestsFailed == 0`. Skill description omits the `coverage >= 0` condition — the skill says "no test evidence" without explaining the coverage gate. Minor misalignment. | -0 |
| 4 | 7h-Schema-code alignment | plugins/forge/skills/breakdown-tasks/templates/index.schema.json | Schema has `"required": ["id", "title", "priority", "status", "file", "scope"]` but Go struct `Task` has no `scope` required validation in validate.go — scope is optional in the Go code. Also `sourceTaskID` exists in Go struct but is absent from JSON schema. Known acceptable per rubric, but `record` field is in schema as optional but not in Go required validation either. Net: schema is stricter than code. | -5 |
| 5 | 7h-Schema-code alignment | plugins/forge/skills/breakdown-tasks/templates/index.schema.json | Schema requires `scope` in task properties but the Go struct marks it `omitempty`. The validator (validate.go) does not check for missing `scope`. Schema says scope is required, code silently accepts missing scope. | -5 |
| 6 | 7i-All-completed hook align | plugins/forge/hooks/guide.md:127 | Guide says all-completed runs `just compile → just fmt → just lint` (quality gate) then `just test` then e2e. Go code (all_completed.go:119-134) runs `just.LintGateSequence()` which is compile→fmt→lint, then separate unit tests, then e2e. But guide.md line 128 says "Quality gate: `just compile → just fmt → just lint`" which matches code. However guide does not mention the auto-fix-task creation behavior on failure. Minor gap. | -5 |
| 7 | 7j-Template existence | task-cli (embedded binary) | The `fix-task` template is referenced in skills (breakdown-tasks, quick-tasks, execute-task, run-tasks, task-executor) and confirmed to exist via `task template fix-task`. But no on-disk file exists in `plugins/forge/skills/*/templates/` — it is embedded in the Go binary. No deduction since it exists in the binary. | -0 |
| 8 | 7j-Template existence | plugins/forge/skills/breakdown-tasks/templates/index.json | Template index.json referenced by breakdown-tasks and quick-tasks exists on disk. No fix-task template on disk in forge plugin — embedded in Go binary. Partial alignment. | -5 |
| 9 | 8-Hook CLI commands valid | plugins/forge/hooks/hooks.json:39 | `task cleanup` in SessionEnd hook — exists in `task -h` output. `task all-completed` in Stop hook — exists. But `bash "${CLAUDE_PLUGIN_ROOT}/scripts/validate-index.sh"` in PostToolUse — this is a script, not a task CLI command. The hook wiring mixes task CLI commands with direct bash scripts. The validate-index.sh script runs `task validate` internally, which IS a valid command. No deduction for scripts. | -0 |
| 10 | 8-Hook CLI commands valid | plugins/forge/hooks/hooks.json:33 | PostToolUse hook runs `bash "${CLAUDE_PLUGIN_ROOT}/scripts/validate-index.sh"`. The script itself is a bash wrapper. While this works, the hook references a script that calls `task validate`, but `task validate` is NOT listed in the "Hook CLI commands valid" check directly. The command chain is indirect. | -5 |
| 11 | 10-allowed_tools declared | plugins/forge/commands/fix-bug.md | fix-bug.md has `allowed_tools: ["Bash", "Read", "Write", "Edit", "Grep", "Glob", "Agent", "LSP"]` and uses Bash, Read, Write, Edit, Grep, Glob, LSP. No tool usage without declaration. CORRECT. | -0 |
| 12 | 10-allowed_tools declared | plugins/forge/commands/record-decision.md | record-decision.md uses `Read`, `Write`, `Edit`, `Bash`, `AskUserQuestion`. Has `allowed_tools: ["Read", "Write", "Edit", "Bash", "AskUserQuestion"]`. CORRECT. | -0 |
| 13 | 10-allowed_tools declared | plugins/forge/commands/simplify-skill.md | simplify-skill.md uses `Read`, `AskUserQuestion`. Has `allowed_tools: ["Read", "AskUserQuestion"]`. CORRECT. | -0 |
| 14 | 10-allowed_tools declared | plugins/forge/commands/init-forge.md | init-forge.md uses `Bash`, `Read`. Has `allowed_tools: ["Bash", "Read"]`. CORRECT. | -0 |
| 15 | 10-allowed_tools declared | plugins/forge/commands/quick.md | quick.md invokes Skill tool (`Skill(skill="forge:brainstorm")`) and uses `AskUserQuestion`. Has `allowed_tools: ["Bash", "Read", "Write", "Edit", "Grep", "Glob", "Agent", "Skill", "AskUserQuestion"]`. CORRECT. | -0 |
| 16 | 10-allowed_tools declared | plugins/forge/commands/run-tasks.md | run-tasks.md uses `Bash`, `Read`, `Agent`. Has `allowed_tools: ["Bash", "Read", "Agent", "TaskOutput", "Skill"]`. CORRECT. | -0 |
| 17 | 10-allowed_tools declared | plugins/forge/commands/execute-task.md | execute-task.md uses `Bash`, `Read`, invokes Skill tool. Has `allowed_tools: ["Bash", "Read", "Write", "Edit", "Grep", "Glob", "Agent", "LSP"]`. CORRECT. | -0 |
| 18 | 10-allowed_tools declared | plugins/forge/commands/git-checkout.md | git-checkout.md uses `Bash`, `Read`. Has `allowed_tools: ["Bash", "Read"]`. CORRECT. | -0 |
| 19 | 10-allowed_tools declared | plugins/forge/commands/git-commit.md | git-commit.md uses `Bash`, `Read`. Has `allowed_tools: ["Bash", "Read"]`. CORRECT. | -0 |
| 20 | 10-allowed_tools declared | plugins/forge/commands/extract-design-md.md | extract-design-md.md uses `Bash`, `Read`, `Write`, `WebFetch`. Has `allowed_tools: ["Bash", "Read", "Write", "WebFetch"]`. CORRECT. | -0 |
| 21 | 10-allowed_tools declared | plugins/forge/commands/gen-sitemap.md | gen-sitemap.md uses `Bash`, `Read`, `Write`, `Grep`, `Glob`. Has `allowed_tools: ["Bash", "Read", "Write", "Grep", "Glob"]`. CORRECT. | -0 |
| 22 | 10-allowed_tools missing | plugins/forge/skills/brainstorm/SKILL.md | brainstorm skill invokes `/eval-proposal` via Skill tool. The SKILL.md does not have `allowed_tools` frontmatter — skills are not required to declare this. No deduction for skills. | -0 |
| 23 | 10-allowed_tools missing | plugins/forge/commands/simplify-skill.md:1 | No `allowed_tools` for the `Skill` tool, but simplify-skill only uses Read and AskUserQuestion. However, it says "extract content to new files" which would need Write/Edit. The allowed_tools list does not include Write/Edit. | -10 |
| 24 | 9-Conciseness | plugins/forge/hooks/guide.md:150-170 | Guide contains a `task record` workflow section with detailed command usage (`task record <id> --data docs/features/{slug}/tasks/process/record.json`). This duplicates the detailed workflow in record-task/SKILL.md. Guide should reference, not replicate. | -5 |
| 25 | 12-Safety marker consistency | plugins/forge/agents/doc-reviser.md | doc-reviser agent has `<EXTREMELY-IMPORTANT>` and `<HARD-RULE>` markers — all actionable. CORRECT. | -0 |
| 26 | 12-Safety marker consistency | plugins/forge/agents/doc-scorer.md | doc-scorer agent has `<EXTREMELY-IMPORTANT>` and `<HARD-RULE>` markers — all actionable. CORRECT. | -0 |
| 27 | 12-Safety marker consistency | plugins/forge/agents/error-fixer.md | error-fixer agent has `<EXTREMELY-IMPORTANT>` with 3 rules. Actionable. CORRECT. | -0 |
| 28 | 12-Dispatch cmd coverage | plugins/forge/commands/fix-bug.md | fix-bug is a dispatch-style command (executes autonomous bug fix workflow). Has `<EXTREMELY-IMPORTANT>` block with 5 safety constraints. CORRECT. | -0 |
| 29 | 12-Dispatch cmd coverage | plugins/forge/commands/execute-task.md | execute-task dispatches to subagents indirectly (via task-executor). Has `<EXTREMELY-IMPORTANT>` with 5 rules. CORRECT. | -0 |
| 30 | 12-Dispatch cmd coverage | plugins/forge/commands/run-tasks.md | run-tasks dispatches task-executor subagents. Has `<EXTREMELY-IMPORTANT>` with 5 iron laws. CORRECT. | -0 |
| 31 | 12-Dispatch cmd coverage | plugins/forge/commands/quick.md | quick dispatches brainstorm, quick-tasks, run-tasks. Has `<EXTREMELY-IMPORTANT>` with 3 rules. CORRECT. | -0 |
| 32 | 12-Command/agent markers | plugins/forge/skills/record-task/SKILL.md | Two `<EXTREMELY-IMPORTANT>` blocks and a `<HARD-GATE>`. All actionable. CORRECT. | -0 |
| 33 | 12-Command/agent markers | plugins/forge/skills/gen-test-scripts/SKILL.md | Two `<HARD-RULE>` blocks, two `<EXTREMELY-IMPORTANT>` blocks, one `<HARD-GATE>`. All actionable with specific rules. CORRECT. | -0 |
| 34 | 12-Command/agent markers | plugins/forge/skills/consolidate-specs/SKILL.md | `<HARD-GATE>` block present. Actionable. CORRECT. | -0 |
| 35 | 12-Safety marker consistency | plugins/forge/skills/init-justfile/SKILL.md | Has `<EXTREMELY-IMPORTANT>` block — contains 4 rules, all actionable (manual-only, prompt before overwrite, boundary markers, run verification). CORRECT. | -0 |
| 36 | 12-Safety marker consistency | plugins/forge/skills/run-e2e-tests/SKILL.md | Has `<HARD-GATE>`, `<PRINCIPLE>`, `<HARD-RULE>` blocks. All actionable. CORRECT. | -0 |
| 37 | 12-Safety marker issue | plugins/forge/agents/task-executor.md | task-executor has `<PROHIBITIONS>` block — specific forbidden actions. CORRECT. However, the error-fixer agent does NOT have a STOP/HARD-RULE block equivalent to task-executor's "ONE TASK PER INVOCATION" rule. While error-fixer is not a dispatcher, it could benefit from clearer boundaries. Minor. | -5 |

---

## Attack Points

### Attack 1: [7f — Record validation quality gate mismatch]

**Where**: `plugins/forge/skills/record-task/SKILL.md:129-133`
**What's wrong**: The skill states the quality gate runs `just compile [scope] → just fmt [scope] → just lint [scope] → just test [scope]` (4 steps including test). The Go source code (`task-cli/internal/cmd/record.go:445`) calls `just.DefaultGateSequence()` which runs compile→fmt→lint only (3 steps, NO test step). The skill misleads agents into thinking `just test` is part of the pre-record gate, when actually the task executor runs tests separately before calling record.
**How to fix**: Update record-task/SKILL.md Quality Gate section to match actual code: `just compile [scope] → just fmt [scope] → just lint [scope]` (3 steps, no test). Note that `just test` is the task executor's responsibility, not part of `task record`'s gate.

### Attack 2: [10 — simplify-skill missing Write/Edit in allowed_tools]

**Where**: `plugins/forge/commands/simplify-skill.md:1`
**What's wrong**: The command's Phase 4 says "Extract content to new files" and "Replace in skill.md with reference", which requires Write and Edit tools. But `allowed_tools` only lists `["Read", "AskUserQuestion"]`. The agent cannot create new files or edit existing ones without these tools.
**How to fix**: Add `"Write"` and `"Edit"` to the `allowed_tools` list: `["Read", "Write", "Edit", "AskUserQuestion"]`.

### Attack 3: [7h — Schema-code scope required mismatch]

**Where**: `plugins/forge/skills/breakdown-tasks/templates/index.schema.json:51` vs `task-cli/internal/cmd/validate.go`
**What's wrong**: JSON schema declares `scope` as required in each task (`"required": ["id", "title", "priority", "status", "file", "scope"]`). But Go validation code (validate.go:115-133) does NOT check for missing `scope` — it only validates ID, title, file, status, and priority. Tasks without `scope` pass Go validation but fail schema validation. Additionally, `sourceTaskID` exists in the Go struct (`pkg/task/types.go:27`) but is absent from the JSON schema entirely (known acceptable per rubric).
**How to fix**: Either add scope validation to validate.go (check `t.Scope == ""` produces a warning) or remove `scope` from the schema's required list and rely on the omitempty default.

### Attack 4: [9 — Guide contains duplicated task record workflow detail]

**Where**: `plugins/forge/hooks/guide.md:150-170`
**What's wrong**: The guide's Task-CLI section contains a detailed `task record` workflow with specific command syntax (`task record <id> --data docs/features/{slug}/tasks/process/record.json`), file locations, and a "One command does 2 things" explanation. This duplicates the authoritative content in `plugins/forge/skills/record-task/SKILL.md` lines 108-118. The guide is a workflow reference, not an API doc — CLI usage details belong in the skill file.
**How to fix**: Replace the detailed task record workflow in guide.md with a one-line reference: "For record workflow details, see `/record-task` skill." Keep the Key Commands table but remove the detailed `task record` Workflow subsection.

### Attack 5: [3 — Quality gate failure table duplicated across 4 files]

**Where**: `plugins/forge/commands/execute-task.md:61-66`, `plugins/forge/commands/fix-bug.md:174-183`, `plugins/forge/agents/task-executor.md:96-107`, `plugins/forge/agents/error-fixer.md:86-98`
**What's wrong**: The quality gate failure table (compile→fix→retry, fmt→blocked, lint→self-fix, test→fix→retry) is copy-pasted verbatim across 4 files. The canonical location should be guide.md's Scope Resolution section, with each file referencing it. The rubric grants an exception for autonomous agents that cannot read other files at runtime — task-executor and error-fixer qualify. But execute-task and fix-bug (commands running in the main session) could reference guide.md.
**How to fix**: In execute-task.md and fix-bug.md, replace the inline table with "Apply the Quality Gate Protocol from the Forge Guide." Keep the table in task-executor.md and error-fixer.md since they run as subagents without access to guide.md.

### Attack 6: [7i — All-completed hook guide missing auto-fix-task behavior]

**Where**: `plugins/forge/hooks/guide.md:127-129` vs `task-cli/internal/cmd/all_completed.go:119-134`
**What's wrong**: Guide says the all-completed hook runs quality gate then tests then e2e. But Go code shows that on any gate failure, it automatically creates a fix-task (P0, breaking) via `addFixTask()`. The guide does not mention this auto-fix-task creation. Agents reading the guide would not know that failures trigger automatic fix-task creation.
**How to fix**: Add a note to guide.md's All-Completed Hook section: "On failure at any step, a P0 fix-task is automatically created. Run `task claim` to pick it up."

### Attack 7: [12 — error-fixer agent lacks ONE-TASK constraint equivalent]

**Where**: `plugins/forge/agents/error-fixer.md`
**What's wrong**: The error-fixer agent has a 5-step workflow but no explicit `<HARD-RULE>` or `<PROHIBITIONS>` block limiting it to a single fix. The task-executor agent has a clear STOP section with `<PROHIBITIONS>` preventing task claim or reading next task. Error-fixer lacks equivalent boundaries — nothing prevents it from attempting multiple fixes or claiming new tasks after completing one.
**How to fix**: Add a STOP section with `<HARD-RULE>` to error-fixer.md: "ONE FIX PER INVOCATION. After Step 5, STOP immediately. FORBIDDEN: run task claim, read other task files, or attempt additional fixes."

---

## Previous Issues Check

N/A (iteration 1)

---

## Fix Summary

N/A (iteration 1)

---

## Verdict

- **Score**: 885/1000
- **Target**: 950/1000
- **Gap**: 65 points
- **Action**: Top fixes needed: (1) Correct record-task quality gate description (-25), (2) Add Write/Edit to simplify-skill allowed_tools (-10), (3) Fix schema-code scope alignment (-10), (4) Remove duplicated guide.md task record detail (-5), (5) Deduplicate quality gate tables (-5), (6) Update guide for auto-fix-task behavior (-5), (7) Add STOP rule to error-fixer (-5)
