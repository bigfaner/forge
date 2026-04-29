---
date: "2026-04-30"
doc_dir: "docs/features/justfile-standard-vocabulary/prd/"
iteration: 3
target_score: 90
evaluator: Claude (automated, adversarial)
---

# PRD Eval — Iteration 3

**Score: 92/100** (target: 90)

```
┌─────────────────────────────────────────────────────────────────┐
│                       PRD QUALITY SCORECARD                      │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 1. Background & Goals        │  19      │  20      │ ✅         │
│    Background three elements │  7/7     │          │            │
│    Goals quantified          │  7/7     │          │            │
│    Logical consistency       │  5/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 2. Flow Diagrams             │  18      │  20      │ ✅         │
│    Mermaid diagram exists    │  7/7     │          │            │
│    Main path complete        │  7/7     │          │            │
│    Decision + error branches │  4/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 3. Functional Specs          │  19      │  20      │ ✅         │
│    Tables complete           │  7/7     │          │            │
│    Field descriptions clear  │  6/7     │          │            │
│    Validation rules explicit │  6/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 4. User Stories              │  18      │  20      │ ✅         │
│    Coverage per user type    │  7/7     │          │            │
│    Format correct            │  7/7     │          │            │
│    AC per story              │  4/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 5. Scope Clarity             │  18      │  20      │ ✅         │
│    In-scope concrete         │  7/7     │          │            │
│    Out-of-scope explicit     │  5/7     │          │            │
│    Consistent with specs     │  6/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ TOTAL                        │  92      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| prd-spec.md:42 vs prd-spec.md:210-224 | "消除原始命令" goal says "当前至少 8 处散落的原始命令" but section 5.4 migration table shows only 3 rows with actual command changes (npx serve -> just run, just build -> just compile x2). The remaining 8 rows say "保持不变". The number "8 处" is not substantiated by the migration data in the PRD itself. The background says "散落着原始 shell 命令" but only 1 non-just raw command (`npx serve`) is identified in the migration table. This inflates the scope claim without evidence. | -1 pt (logical consistency) |
| prd-spec.md:152-173 | init-justfile diagram: `ErrorNoProject` node says "错误: 未检测到已知项目标记文件" but does not specify exit code, error channel (stdout vs stderr), or suggested remediation. Compare with the scope-resolution diagram (line 88-89) which specifies exact message format and fallback behavior. Inconsistency in error-handling specificity across the two diagrams. | -1 pt (error branches incomplete) |
| prd-spec.md:159-163 | init-justfile diagram: `ConfirmOverwrite` -> user confirms or aborts. This is an interactive prompt, but section "Agent 友好性需求" (line 252) states "所有 recipe 不使用交互式输入（如 read、select），agent 无法处理交互". If init-justfile is invoked by a skill/agent, this interactive prompt blocks execution. The document does not specify whether init-justfile is agent-facing or human-only, creating a conflict between the flow and the agent-friendliness requirements. | -1 pt (error branches incomplete) |
| prd-spec.md:114-130 | Command table column "必需" has values "否" for all commands except `project-type` which has "是". The column header "必需" is ambiguous — required for what? Required to exist in every justfile? Required as a mandatory positional argument? The scope-resolution flow diagram (line 87) calls `just project-type` and handles its failure, treating it as optional at the skill layer. If `project-type` is "required" but skills must handle its absence, the column semantics are contradictory. No footnote or header description clarifies the column meaning. | -1 pt (field descriptions) |
| prd-user-stories.md:14-16 | Story 1 AC line 1: "Given 一个使用 forge 插件的项目，When 任意 skill 执行构建/测试/编译操作时，Then 通过 just <verb> 调用，不直接调用 go、npm、cargo 等语言工具链". "任意 skill" is unbounded — there is no enumeration linking "任意 skill" to the specific 8 skills in the migration checklist (5.4). An acceptance tester cannot verify "任意 skill" without a finite list. The AC should reference the migration checklist explicitly. | -1 pt (AC quality) |
| prd-user-stories.md:54-57 | Story 4 has 4 ACs but AC 4 ("Given agent 连续执行 just install、just compile、just test，When 全部成功时，Then 每步退出码均为 0，agent 无需人工介入") tests a sequential pipeline, not a single command behavior. This AC overlaps with Story 1 (standardized execution) and tests runtime behavior (no human intervention) that has no observable criterion within the PRD's control — it depends on the underlying tools working correctly, not on the justfile vocabulary. | -1 pt (AC quality) |
| prd-spec.md:60-64 | Out-of-scope has 3 items. Still missing: (a) what about commands that do not fit the 15-command model (e.g., `just migrate`, `just deploy`, `just seed`)? (b) Windows/non-Unix shell compatibility — justfile recipes often assume bash, but forge runs on Windows per the environment. (c) Recipe-level documentation generation or help text standards. The out-of-scope list is thin relative to the in-scope ambition. | -2 pts |

---

## Attack Points

### Attack 1: Flow Diagrams — init-justfile flow contradicts agent-friendliness requirements

**Where**: prd-spec.md:159-163 "CheckExisting -> 是 -> ConfirmOverwrite[提示用户确认覆盖]" vs prd-spec.md:252 "所有 recipe 不使用交互式输入（如 read、select），agent 无法处理交互"
**Why it's weak**: The init-justfile flow includes an interactive confirmation prompt when a justfile already exists. But the "Agent 友好性需求" section explicitly prohibits interactive prompts because "agent 无法处理交互". If init-justfile is ever called by an agent (e.g., a skill that auto-initializes a project), the prompt blocks execution. The document does not distinguish between agent-facing and human-only commands, nor does it specify a `--force` or `--yes` flag to skip confirmation. This is a design inconsistency that would surface during implementation.
**What must improve**: Either (a) add a non-interactive mode flag (e.g., `/init-justfile --force`) to the init-justfile flow and document it, or (b) explicitly scope init-justfile as human-only in the user stories and add a note that agents must never invoke it, or (c) change the overwrite behavior to non-interactive (backup old file, generate new one).

### Attack 2: User Stories — Story 1 AC is untestable due to "任意 skill"

**Where**: prd-user-stories.md:14 "When 任意 skill 执行构建/测试/编译操作时，Then 通过 just <verb> 调用，不直接调用 go、npm、cargo 等语言工具链"
**Why it's weak**: "任意 skill" is an unbounded universal quantifier. The PRD has a concrete migration checklist (section 5.4) with exactly 11 files, of which only 3 require actual command changes. An acceptance tester reading this AC would need to check every skill that has ever been written or will ever be written — an impossible task. The AC should constrain itself to the finite set listed in the migration table: "When the 3 skills listed in section 5.4 with actual command changes execute build/test/compile operations, Then they call through `just <verb>` instead of direct toolchain commands." This was flagged in iteration 2 (deduction line 6) and remains unfixed.
**What must improve**: Rewrite Story 1 AC 1 to reference the specific migration table: "Given a project using forge plugins, When the skills listed in section 5.4 migration table execute build/test/compile operations, Then they invoke `just <verb>` and do not directly call `go`, `npm`, `cargo` or other language toolchain executables."

### Attack 3: Scope Clarity — out-of-scope list is too thin for the feature's ambition

**Where**: prd-spec.md:60-64 "Out of Scope" lists only 3 items: other projects' justfiles, project-type in CI, and parallel execution optimization.
**Why it's weak**: The feature introduces 15 standard commands, scope parameters, project-type detection, and a migration of 14 files. The out-of-scope list should address foreseeable confusion points: (a) custom recipes that do not map to any of the 15 verbs — what should a developer do with `just migrate`, `just deploy`, `just seed`? (b) Windows shell compatibility — the environment runs on Windows, but justfile recipes typically assume Unix shells; the document says nothing about cross-platform recipe requirements. (c) Help text or documentation generation for the standard vocabulary. (d) Versioning of the vocabulary itself — what happens when command 16 is needed? The 3-item out-of-scope list was thin in iteration 1, flagged in iteration 2, and remains unchanged in iteration 3.
**What must improve**: Add at least 2-3 more out-of-scope items: "Non-standard recipes (deploy, migrate, seed) — projects may add these as custom recipes outside the standard block", "Cross-platform shell compatibility (Windows/cmd/powershell) — recipes assume bash; Windows users need WSL or Git Bash", "Standard vocabulary versioning or deprecation process — future commands will be added via PRD amendment".

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1 (Iter 2): Goals 3 and 4 lack quantified metrics | ✅ Yes | prd-spec.md:44 now reads "11 个带 scope 参数的命令支持前端/后端选择性操作（15 个命令中 11 个接受 scope 位置参数)" with concrete numbers. prd-spec.md:45 now reads "支持 3 种项目类型（frontend/backend/mixed）自动探测与 justfile 生成" with a count. Both goals now have numeric metrics. |
| Attack 2 (Iter 2): Migration deliverables not traced to section 5.5 | ✅ Yes | prd-spec.md:226-234 section 5.5 now has 5 rows including row 3 "skill 文件 — 8 个 SKILL.md", row 4 "agent 文件 — 2 个 .md agent 定义", row 5 "task 模板 — 4 个 task 模板文件". The migration work is now traceable in the related-changes table. |
| Attack 3 (Iter 2): Story 3 ACs 2-4 break Given/When/Then format | ✅ Yes | prd-user-stories.md:41 now reads "Given 一个任务仅涉及前端目录文件... When /breakdown-tasks 处理该任务时，Then scope 标记为 frontend". Lines 42-43 similarly include explicit When clauses. All Story 3 ACs now follow the format. |

---

## Verdict

- **Score**: 92/100
- **Target**: 90/100
- **Gap**: 0 points (target exceeded by 2)
- **Action**: Target reached. All three iteration 2 attack points were addressed. Remaining deductions are minor: one untestable AC in Story 1, one design inconsistency between init-justfile flow and agent-friendliness requirements, and a thin out-of-scope list. These do not block implementation.
