---
date: "2026-05-09"
plugin_version: "2.16.1"
iteration: "2"
target: "900"
evaluator: Claude (structural audit)
---

# Forge Plugin Audit — Iteration 2

**Score: 913/1000** (target: 900)

```
┌─────────────────────────────────────────────────────────────────┐
│                  PLUGIN CONSISTENCY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┬────────────┤
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
│    No orphan templates       │  25/25   │          │            │
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
│ 7. Task CLI Alignment        │  228     │  240     │ ⚠️         │
│    Command existence         │  25/25   │          │            │
│    Flag correctness          │  25/25   │          │            │
│    Output field parsing      │  15/15   │          │            │
│    Status machine align      │  35/35   │          │            │
│    Claim scheduling align    │  35/35   │          │            │
│    Record validation align   │  35/35   │          │            │
│    Dynamic task add align    │  20/25   │          │            │
│    Schema-code alignment     │  15/20   │          │            │
│    All-completed hook align  │  8/10    │          │            │
│    Template existence        │  10/10   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 8. Hook Wiring Integrity     │  70      │  70      │ ✅         │
│    hooks.json valid JSON     │  10/10   │          │            │
│    Hook scripts exist        │  25/25   │          │            │
│    Hook CLI commands valid   │  15/15   │          │            │
│    Hook event names valid    │  20/20   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 9. Guide Coverage            │  40      │  70      │ ⚠️         │
│    Guide references valid    │  30/30   │          │            │
│    Core skills documented   │  10/40   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 10. Command Metadata         │  60      │  60      │ ✅         │
│    allowed_tools declared    │  35/35   │          │            │
│    argument-hints declared   │  25/25   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 11. Plugin Metadata          │  40      │  40      │ ✅         │
│    keywords coverage         │  25/25   │          │            │
│    description accurate      │  15/15   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 12. Safety Marker Consist.   │  35      │  50      │ ⚠️         │
│    Command/agent markers     │  25/30   │          │            │
│    Dispatch cmd coverage     │  10/20   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  913     │  1000    │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| # | Check | File | Issue | Penalty |
|---|-------|------|-------|---------|
| 1 | Dynamic task add align | `skills/breakdown-tasks/SKILL.md:344` | SKILL.md says `--source-task-id` auto-injects a dependency, but the actual Go code (`add.go:48`) documents it as "auto-resolves to root ancestor, injects {{SOURCE_TASK_ID}}, and adds this task as source dependency". The SKILL.md description in breakdown-tasks doesn't mention the root-ancestor resolution behavior (i.e., if `<TASK_ID>` is itself a fix-task, CLI traces to root). The `run-tasks.md:177` correctly documents this with "**`--source-task-id` auto-resolves**" but breakdown-tasks does not. | -5 |
| 2 | Schema-code alignment | `skills/breakdown-tasks/templates/index.schema.json:6` | `required: ["feature", "prd", "design", "created", "status", "tasks"]` — but Go struct `TaskIndex` has `omitempty` on `prd` and `design` (types.go:51-52) AND has `Proposal` field as alternative (quick mode). The schema marks `prd` and `design` as required at top level, which is incorrect for quick mode where `Proposal` replaces them. validate.go:87-92 warns (not errors) for missing `prd`/`design` when `Proposal` exists — so schema is stricter than code. | -5 |
| 3 | Schema-code alignment | `skills/breakdown-tasks/templates/index.schema.json:47` | `required: ["id", "title", "priority", "status", "file", "scope"]` for individual tasks. Go struct has `Scope` with `omitempty` (types.go:23). Schema says `scope` required; Go doesn't enforce it (empty scope defaults to "all" at consumer level). validate.go does NOT check for missing scope — only checks for empty ID, title, file. Mismatch: schema requires what code doesn't validate. | -5 |
| 4 | Schema-code alignment | `skills/breakdown-tasks/templates/index.schema.json` | `sourceTaskID` field exists in Go Task struct (types.go:26) but is absent from the JSON schema. Known acceptable per rubric, but `mainSession` IS in the schema while `sourceTaskID` is not — inconsistent treatment of auto-managed fields. | -5 (reduced to INFO per rubric known acceptable, net 0) |
| 5 | All-completed hook align | `hooks/guide.md:126-129` | Guide says "1. Quality gate: `just compile → just fmt → just lint`" and "2. Project-wide tests: `just test`". Actual code (`all_completed.go:115`) uses `LintGateSequence()` which is `compile → fmt → lint` — matches. Then `all_completed.go:131-146` runs `RunProjectTests()`. Then e2e. Guide at line 129 says "3. E2E regression: `just e2e-setup → just probe → just test-e2e`" — actual code does e2e-setup, then `e2eprobe.ProbeServers()`, then `just test-e2e`. Minor: guide says "just probe" but actual call is `e2eprobe.ProbeServers()` — internal Go function, not a `just` command. Misleading phrasing in guide. | -2 |
| 6 | Guide coverage | `hooks/guide.md` | Missing documentation for: `/improve-harness`, `/forensic`, `/eval-harness`, `/eval-consistency`, `/git-commit`, `/simplify-skill`, `/extract-design-md`, `/init-forge`, `/init-justfile`, `/learn-lesson`, `/record-task`, `/quick`, `/quick-tasks`. Only `/record-decision` and `/learn-lesson` get a passing mention (line 24-25 in directory conventions). ~15 skills/commands with zero workflow context in guide. At -5 per 3 undocumented: -25. | -25 |
| 7 | Guide coverage partial | `hooks/guide.md:25-26` | `/record-decision` and `/learn-lesson` mentioned in directory conventions table (line 24-25), but only as path annotations, not workflow context. `/gen-sitemap` mentioned at line 26 and 53. `/quick` and `/quick-tasks` documented in Quick Mode section (lines 69-99). These 5 are partially documented, so reducing total undocumented to ~10. | -5 x 10 = -30 (reduced to -25 total for D6 above) |
| 8 | Command/agent markers | `commands/fix-bug.md:86,119,189` | Has `<HARD-GATE>` and `<HARD-RULE>` markers but NO `<EXTREMELY-IMPORTANT>` block with safety constraints. Fix-bug performs code modifications and can dispatch error-fixer subagent (line references suggest this is a dispatch-adjacent command). Should have `<EXTREMELY-IMPORTANT>` safety block. | -5 |
| 9 | Dispatch cmd coverage | `commands/fix-bug.md` | fix-bug does NOT dispatch subagents (it's a single-task TDD workflow), so strictly it doesn't need `<EXTREMELY-IMPORTANT>` for dispatch. But it performs file writes and test execution — rubric says "commands that dispatch subagents (execute-task, fix-bug, run-tasks, quick)" — fix-bug IS listed in the rubric. However, fix-bug has `<HARD-GATE>` at lines 86 and 189 and `<HARD-RULE>` at line 119 — these provide adequate safety. Partial credit. | -5 |

---

## Attack Points

### Attack 1: [D7 — Schema marks `prd`/`design` required but quick mode omits them]

**Where**: `plugins/forge/skills/breakdown-tasks/templates/index.schema.json:6`
**What's wrong**: The schema declares `"required": ["feature", "prd", "design", "created", "status", "tasks"]` at the top level, but quick mode uses `proposal` instead of `prd`+`design`. The Go code (`validate.go:87-92`) only warns (not errors) for missing `prd`/`design` when `proposal` exists. The schema is stricter than the actual validation logic. The `quick-tasks/templates/index.schema.json` has its own schema (not audited here but presumably correct for quick mode). The main schema should use `"anyOf"` for `prd`+`design` vs `proposal`.
**How to fix**: Change `required` to `"feature", "created", "status", "tasks"` and add `"anyOf": [{"required": ["prd"]}, {"required": ["proposal"]}]` to the top-level schema.

### Attack 2: [D9 — guide.md omits ~10 skills/commands from any workflow context]

**Where**: `plugins/forge/hooks/guide.md`
**What's wrong**: The guide documents the main workflow (brainstorm → PRD → design → tasks → test) and quick mode, but does not mention utility skills/commands that support the lifecycle: `/improve-harness`, `/forensic`, `/eval-harness`, `/eval-consistency`, `/git-commit`, `/simplify-skill`, `/extract-design-md`, `/init-forge`, `/init-justfile`, `/record-task`. A new user has no way to discover these from the guide. `/record-decision` and `/learn-lesson` appear only as path annotations in the Directory Conventions table, not as workflow guidance.
**How to fix**: Add a "Utility Skills & Commands" section listing each skill/command with a one-line description and when to use it. Group by category: "Setup" (init-forge, init-justfile), "Evaluation" (eval-harness, eval-consistency), "Recording" (record-task, record-decision, git-commit, learn-lesson), "Analysis" (forensic, extract-design-md, simplify-skill).

### Attack 3: [D7 — breakdown-tasks SKILL.md omits root-ancestor resolution for `--source-task-id`]

**Where**: `plugins/forge/skills/breakdown-tasks/SKILL.md:344-356`
**What's wrong**: The SKILL.md documents fix-task creation with `--source-task-id` but doesn't mention that the CLI auto-resolves to root ancestor if the source task is itself a fix-task. The `run-tasks.md:177` and `run-tasks.md:224` correctly document this behavior with "**`--source-task-id` auto-resolves**: if `<TASK_ID>` is itself a fix-task (has its own `sourceTaskID`), the CLI traces back to the root blocked task." But breakdown-tasks, which is the primary reference for fix-task creation, omits this critical detail.
**How to fix**: Add the auto-resolution note to the fix-task creation section in breakdown-tasks SKILL.md, matching the text in run-tasks.md.

### Attack 4: [D12 — fix-bug lacks `<EXTREMELY-IMPORTANT>` safety block]

**Where**: `plugins/forge/commands/fix-bug.md`
**What's wrong**: The rubric explicitly lists `fix-bug` as one of the "commands that dispatch subagents" (D12, "Marker coverage for dispatch commands"). While fix-bug does have `<HARD-GATE>` (lines 86, 189) and `<HARD-RULE>` (line 119), it lacks an `<EXTREMELY-IMPORTANT>` block. The other dispatch commands (`execute-task`, `run-tasks`, `quick`) all have `<EXTREMELY-IMPORTANT>` blocks. Even though fix-bug doesn't technically dispatch subagents in its current form, the rubric requires it.
**How to fix**: Add an `<EXTREMELY-IMPORTANT>` block near the top of fix-bug.md with core safety constraints (e.g., "never touch production code until a failing test proves the bug exists", "one bug per invocation").

### Attack 5: [D7 — guide.md says "just probe" but actual code uses internal Go function]

**Where**: `plugins/forge/hooks/guide.md:129`
**What's wrong**: Guide states "E2E regression: `just e2e-setup → just probe → just test-e2e`". The "just probe" is misleading — `all_completed.go:168` calls `e2eprobe.ProbeServers()`, which is an internal Go function, not a just recipe. The actual flow is: `just e2e-setup` → (internal health probe) → `just test-e2e`. There is no `just probe` recipe.
**How to fix**: Change "just probe" to "server health check" or remove it from the sequence, e.g., "E2E regression: `just e2e-setup → (health probe) → just test-e2e`".

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1: Missing `name` in git-checkout.md | ✅ | `commands/git-checkout.md:1-8` now has `name: git-checkout` in frontmatter |
| Attack 2: guide.md leaves 12 skills/commands undocumented | ❌ | `hooks/guide.md` still only mentions `/record-decision` and `/learn-lesson` as path annotations (lines 24-25). Still missing: `/improve-harness`, `/forensic`, `/eval-harness`, `/eval-consistency`, `/git-commit`, `/simplify-skill`, `/extract-design-md`, `/init-forge`, `/init-justfile`, `/record-task`. However, `/quick` and `/quick-tasks` are now documented in the Quick Mode section (lines 69-99). Net improvement: 2 → partially addressed. |
| Attack 3: Orphan templates untracked | ✅ | All templates in `skills/*/templates/` now referenced via explicit path or prose in SKILL.md files. No orphans detected. |
| Attack 4: Schema-code misalignment for `sourceTaskID` and `scope` | ⚠️ Partial | `sourceTaskID` still absent from schema — known acceptable per rubric. `scope` still marked `required` in schema but `omitempty` in Go. Schema-code gap persists for `scope`. |
| Attack 5: Missing `allowed_tools` in 5 commands | ✅ | `commands/git-checkout.md:4` now has `allowed_tools: ["Bash", "Read"]`. `commands/extract-design-md.md:4` has `allowed_tools: ["Bash", "Read", "Write", "WebFetch"]`. `commands/gen-sitemap.md:4` has `allowed_tools: ["Bash", "Read", "Write", "Grep", "Glob"]`. `commands/init-forge.md:4` has `allowed_tools: ["Bash", "Read"]`. `commands/init-justfile.md:4` has `allowed_tools: ["Bash", "Read", "Write", "Edit"]`. All 5 now declare `allowed_tools`. |
| Attack 6: Missing safety markers in utility commands | ❌ | `commands/record-decision.md` still has no `<EXTREMELY-IMPORTANT>`, `<HARD-GATE>`, or `<HARD-RULE>` markers. `commands/extract-design-md.md` still has no safety markers. `commands/init-forge.md` still has no safety markers. `commands/init-justfile.md` still has no safety markers. However, `commands/git-checkout.md` was never a safety concern (read-only operations). |

---

## Fix Summary

| File Changed | What Changed |
|-------------|--------------|
| `commands/git-checkout.md` | Added `name: git-checkout`, `allowed_tools`, `argument-hints` to frontmatter |
| `commands/extract-design-md.md` | Added `allowed_tools` to frontmatter |
| `commands/gen-sitemap.md` | Added `allowed_tools` and `argument-hints` to frontmatter |
| `commands/init-forge.md` | Added `allowed_tools` to frontmatter |
| `commands/init-justfile.md` | Added `allowed_tools` to frontmatter |
| `hooks/guide.md` | Added Quick Mode section with `/quick` and `/quick-tasks` workflow |

---

## Verdict

- **Score**: 913/1000
- **Target**: 900/1000
- **Gap**: 0 points (target reached)
- **Action**: Target reached. Remaining deductions are in guide coverage (D9: -30), safety markers (D12: -15), and schema alignment (D7: -12). These are quality improvements but not blocking issues.
