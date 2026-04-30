---
date: "2026-04-30"
doc_dir: "docs/features/justfile-standard-vocabulary/prd/"
iteration: 2
target_score: 90
evaluator: Claude (automated, adversarial)
---

# PRD Eval — Iteration 2

**Score: 88/100** (target: 90)

```
┌─────────────────────────────────────────────────────────────────┐
│                       PRD QUALITY SCORECARD                      │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 1. Background & Goals        │  17      │  20      │ ⚠️         │
│    Background three elements │  7/7     │          │            │
│    Goals quantified          │  5/7     │          │            │
│    Logical consistency       │  5/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 2. Flow Diagrams             │  19      │  20      │ ✅         │
│    Mermaid diagram exists    │  7/7     │          │            │
│    Main path complete        │  7/7     │          │            │
│    Decision + error branches │  5/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 3. Functional Specs          │  18      │  20      │ ⚠️         │
│    Tables complete           │  7/7     │          │            │
│    Field descriptions clear  │  6/7     │          │            │
│    Validation rules explicit │  5/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 4. User Stories              │  18      │  20      │ ⚠️         │
│    Coverage per user type    │  7/7     │          │            │
│    Format correct            │  7/7     │          │            │
│    AC per story              │  4/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 5. Scope Clarity             │  16      │  20      │ ⚠️         │
│    In-scope concrete         │  7/7     │          │            │
│    Out-of-scope explicit     │  5/7     │          │            │
│    Consistent with specs     │  4/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ TOTAL                        │  88      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| prd-spec.md:44 | "支持混合项目" goal says "混合项目可通过 scope 参数选择性操作前端/后端" — this is a boolean capability, not a quantified metric. No percentage, count, or measurable threshold is provided. | -2 pts |
| prd-spec.md:45 | "自适应生成" goal says "init-justfile 根据项目结构自动生成匹配类型的 justfile" — no metric (e.g., "% of project types correctly classified", "3 project type templates"). Vague. | -2 pts |
| prd-spec.md:42 vs prd-spec.md:210-224 | "消除原始命令" goal claims "当前至少 8 处散落的原始命令" but migration table 5.4 shows only 3 rows with actual command changes (`npx serve` → `just run`, `just build && just test` → `just compile && just test` x2). The remaining 8 rows say "保持不变". If only 3 actual migrations are needed, claiming "at least 8 scattered raw commands" inflates the scope without evidence. | -2 pts (inconsistency) |
| prd-spec.md:152-173 | init-justfile diagram shows `ErrorNoProject` branch ending at `[结束]` with no guidance — should indicate what the user sees (error message content, exit code, suggested next steps). | -1 pt |
| prd-spec.md:114-130 | Command table column "位置参数" uses `frontend`/`backend` for scoped commands but column "必需" says "否" — yet `project-type` has "是". The "必需" header is ambiguous: required for what? Required to exist in the justfile? Required as a positional argument? The column meaning should be clarified in a footnote or header description. | -1 pt |
| prd-spec.md:200-208 | Validation rules table covers scope validation well, but the entry "混合项目 + 非法 scope" specifies stderr message `[forge] invalid scope 'foo'; expected frontend/backend` — this message format is a justfile-internal concern, not a skill-layer concern. The document does not specify who produces this message: is it the justfile recipe itself (via a shell guard), or the skill wrapper? The implementation boundary is unclear. | -1 pt |
| prd-user-stories.md:14-16 | Story 1 AC line 1: "Given 一个使用 forge 插件的项目，When 任意 skill 执行构建/测试/编译操作时" — "任意 skill" is unbounded. This AC cannot be verified because there is no enumeration of which skills exist and which ones execute build/test/compile. Should reference the migration checklist in prd-spec 5.4 to make the scope finite and testable. | -1 pt |
| prd-user-stories.md:39-43 | Story 3 AC lines 2-3 use "Then scope 标记为 frontend" but do not follow strict Given/When/Then format. Lines 41-43 read "Given ... Then ..." without a When clause. The AC mixes the format. | -1 pt |
| prd-spec.md:62-64 | Out-of-scope lists only 3 items. Missing: what about recipes that don't fit the 15-command model (e.g., `just migrate`, `just deploy`)? What about Windows/non-Unix shell compatibility? What about recipe-level documentation generation? The out-of-scope is thin. | -2 pts |
| prd-spec.md:51-58 vs prd-spec.md:226-231 | In-scope item "更新 8 个 skill 文件使用新词汇" and "更新 2 个 agent 文件" but section 5.5 关联性需求改动 lists only 2 items (task-cli schema + e2e tests). The skill/agent/command file migrations are listed as in-scope deliverables but are absent from the related changes table — they should appear there as well for cross-reference traceability. | -2 pts |

---

## Attack Points

### Attack 1: Background & Goals — two goals lack quantified metrics

**Where**: prd-spec.md:44 "混合项目可通过 scope 参数选择性操作前端/后端" and prd-spec.md:45 "init-justfile 根据项目结构自动生成匹配类型的 justfile"
**Why it's weak**: The rubric requires "at least one numeric target (% , count, time)" per goal. Goals 3 and 4 have none. Goal 3 is a binary capability statement (either supported or not). Goal 4 describes behavior without any success criterion — what counts as "自适应"? If init-justfile generates a justfile that compiles but misses `project-type`, does that count? The first two goals are well-quantified ("0 处" and "15 个"), making the gap in goals 3-4 conspicuous. This was flagged in iteration 1 and the goals remain unchanged.
**What must improve**: Replace "混合项目可通过 scope 参数选择性操作前端/后端" with a measurable target, e.g., "混合项目支持 3 种 scope 操作（全部/前端/后端），覆盖 15 个标准命令中的 11 个带 scope 参数的命令". For goal 4, add a metric like "支持 3 种项目类型（前端/后端/混合）的自动探测，每种生成独立模板".

### Attack 2: Scope Clarity — migration deliverables not traced to related changes

**Where**: prd-spec.md:51-58 in-scope items "更新 8 个 skill 文件使用新词汇", "更新 2 个 agent 文件", "更新 4 个 task 模板" vs prd-spec.md:226-231 section 5.5 which lists only 2 items
**Why it's weak**: The in-scope checklist commits to updating 14 files across 3 categories (8 skills + 2 agents + 4 task templates). Section 5.5 "关联性需求改动" — which is supposed to map cross-module impact — lists only task-cli schema and e2e tests. The 14 file migrations are the largest work items in the entire PRD, yet they have zero traceability in the related-changes table. This means a developer reading section 5.5 to understand downstream impact would miss the bulk of the implementation work. The migration checklist in 5.4 enumerates the files but does not state what module each file belongs to, what the acceptance test is, or how to verify each migration. The scope section and the related-changes section are disconnected.
**What must improve**: Expand section 5.5 to include rows for each migration category: "skill 文件迁移 — 更新 8 个 SKILL.md 文件 — 将原始 shell 命令替换为标准词汇", "agent 文件迁移 — 更新 2 个 .md 文件 — 同上", "task 模板迁移 — 更新 4 个模板 — 同上". Alternatively, add a traceability column to the migration table (5.4) linking each file to its module and acceptance criterion.

### Attack 3: User Stories — Story 3 acceptance criteria break Given/When/Then format

**Where**: prd-user-stories.md:41-43 "Given 一个任务仅涉及前端目录文件（如 `web/src/components/Button.tsx`），Then scope 标记为 `frontend`"
**Why it's weak**: The AC for Story 3 has 4 criteria. The first (line 40) correctly uses Given/When/Then. But lines 41-43 each skip the When clause, going directly from Given to Then. For example, line 41 says "Given 一个任务仅涉及前端目录文件... Then scope 标记为 `frontend`" — the When clause (presumably "When breakdown-tasks processes this task") is missing. Line 43 similarly omits When. This makes the AC ambiguous: is the scope assigned during task breakdown (When = `/breakdown-tasks` runs)? Or is it assigned by some other process? The rubric requires "every story has at least one AC in Given/When/Then format" — half of Story 3's ACs break the format. Story 5's two ACs are properly formatted, making Story 3's inconsistency stand out.
**What must improve**: Rewrite Story 3 ACs 2-4 to include explicit When clauses. Example: "Given 一个任务仅涉及前端目录文件（如 `web/src/components/Button.tsx`），When `/breakdown-tasks` 处理该任务时，Then scope 标记为 `frontend`". Apply the same pattern to ACs 3 and 4.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1 (Iter 1): 16 vs 15 command count inconsistency | ✅ Yes | prd-spec.md:24 now reads "扩展至 15 个", line 43 reads "15 个标准命令". The command table (lines 114-130) has exactly 15 rows. All references are consistent. |
| Attack 2 (Iter 1): Missing init-justfile generation flow diagram | ✅ Yes | prd-spec.md:152-173 adds a complete Mermaid flowchart for init-justfile generation with start-to-end path, decision points (`HasSignals`, `CheckExisting`, `MixedCheck`), and error/abort branches (`ErrorNoProject`, `Abort`). |
| Attack 2 (Iter 1): Missing error branches for `just project-type` failure in scope-resolution diagram | ✅ Yes | prd-spec.md:82-99 now includes `PTExit` decision node with three error paths: non-zero exit code, missing recipe, unexpected output — all routing to `PTError` then `Fallback`. |
| Attack 3 (Iter 1): Missing validation rules for invalid scope arguments | ✅ Yes | prd-spec.md:200-208 adds a 5-row validation rules table covering: valid scope on mixed project, invalid scope (with exact stderr message), scope on non-mixed project (with fallback behavior), missing `project-type` recipe, and unexpected `project-type` output. |
| Attack 3 (Iter 1): Missing JSON schema for scope field in index.json | ✅ Yes | prd-spec.md:187-196 adds a JSON Schema snippet with `type: string`, `enum: ["frontend", "backend", "all"]`, `default: "all"`, plus a note that the field is optional and consumers should treat missing values as `"all"`. |

---

## Verdict

- **Score**: 88/100
- **Target**: 90/100
- **Gap**: 2 points
- **Action**: Continue to iteration 3. Priority fixes: (1) quantify goals 3 and 4 with numeric metrics, (2) add migration work items to section 5.5 related-changes table or add traceability to section 5.4, (3) fix Story 3 ACs 2-4 to include explicit When clauses in Given/When/Then format.
