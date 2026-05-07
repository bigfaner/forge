---
date: "2026-05-07"
doc_dir: "docs/proposals/tech-design-db-schema/"
iteration: 2
target_score: 90
evaluator: Claude (automated, adversarial)
---

# Proposal Eval — Iteration 2

**Score: 89/100** (target: 90)

```
┌─────────────────────────────────────────────────────────────────┐
│                    PROPOSAL QUALITY SCORECARD                     │
├──────────────────────┬──────────┬──────────┬────────────────────┤
│ Dimension            │ Score    │ Max      │ Status             │
├──────────────────────┼──────────┼──────────┼────────────────────┤
│ 1. Problem Definition│   19     │  20      │ ✅                 │
│    Problem clarity   │  7/7     │          │                    │
│    Evidence provided │  7/7     │          │                    │
│    Urgency justified │  5/6     │          │                    │
├──────────────────────┼──────────┼──────────┼────────────────────┤
│ 2. Solution Clarity  │   18     │  20      │ ✅                 │
│    Approach concrete │  7/7     │          │                    │
│    User-facing behav │  6/7     │          │                    │
│    Differentiated    │  5/6     │          │                    │
├──────────────────────┼──────────┼──────────┼────────────────────┤
│ 3. Alternatives Anal │   13     │  15      │ ⚠️                 │
│    Alternatives ≥2   │  5/5     │          │                    │
│    Pros/cons honest  │  4/5     │          │                    │
│    Rationale justif  │  4/5     │          │                    │
├──────────────────────┼──────────┼──────────┼────────────────────┤
│ 4. Scope Definition  │   14     │  15      │ ✅                 │
│    In-scope concrete │  5/5     │          │                    │
│    Out-of-scope expl │  5/5     │          │                    │
│    Scope bounded     │  4/5     │          │                    │
├──────────────────────┼──────────┼──────────┼────────────────────┤
│ 5. Risk Assessment   │   14     │  15      │ ✅                 │
│    Risks identified  │  5/5     │          │                    │
│    Likelihood+impact │  5/5     │          │                    │
│    Mitigations act.  │  4/5     │          │                    │
├──────────────────────┼──────────┼──────────┼────────────────────┤
│ 6. Success Criteria  │   14     │  15      │ ✅                 │
│    Measurable        │  5/5     │          │                    │
│    Coverage complete │  4/5     │          │                    │
│    Testable          │  5/5     │          │                    │
├──────────────────────┼──────────┼──────────┼────────────────────┤
│ TOTAL (raw)          │  92      │  100     │                    │
│ Inconsistency penalty│  -3      │          │                    │
│ TOTAL                │  89      │  100     │                    │
└──────────────────────┴──────────┴──────────┴────────────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Risk table row 1 vs Scope | Inconsistency: Risk mitigation for detection misjudgment includes a fallback mechanism ("若 frontmatter 为 `none` 但检测到表名引用，发出提醒让用户确认") that is not present in any of the 9 in-scope items, nor is there a success criterion verifying this fallback works | -3 pts |
| Urgency section | No quantified cost data for the cited `jlc-schema-alignment` incident — the postmortem is referenced but specific time/cost impact is omitted | -1 pt (urgency evidence not fully quantified) |
| Alternatives table | Subdirectory alternative rejection cites manifest path logic changes as a con, but Scope #7 already covers manifest modifications — this con is already an in-scope effort, not a genuine incremental cost | -1 pt (analytical depth) |
| Success criteria | No criterion verifying the fallback detection mechanism described in Risk row 1 — a committed mitigation has no verification gate | -1 pt (coverage gap) |

---

## Attack Points

### Attack 1: Scope/Solution inconsistency — fallback detection mechanism is committed but not scoped

**Where**: Risk table, row 1: "增加防御机制：`tech-design` 在 Data Models 章节起草时，若 frontmatter 为 `none` 但检测到表名引用（如 `REFERENCES`、`TABLE` 关键词），发出提醒让用户确认"
**Why it's weak**: This is a concrete behavioral commitment — the system will scan for table references and emit a warning. But none of the 9 in-scope items describe implementing this detection logic, and no success criterion verifies it works. If this mitigation is real, it needs a scope item (e.g., "Modify `tech-design/SKILL.md`: add fallback keyword scan when `db-schema` is `none`") and a success criterion (e.g., "When `db-schema` is `none` but the PRD body contains `REFERENCES` or `TABLE` keywords, `tech-design` emits a warning prompt to the user"). If it's aspirational, it should be documented as a future enhancement, not a committed mitigation.
**What must improve**: Either add the fallback detection to Scope (with effort estimate and success criterion) or remove it from the risk mitigation and replace with a less ambitious but honest mitigation (e.g., "Rely on frontmatter declaration; document the `db-schema` field as a required review point in the PRD approval step").

### Attack 2: Alternatives Analysis — subdirectory rejection relies on circular reasoning

**Where**: Alternatives table, "独立子目录" row: con ② "现有 `manifest-update-design.md` 模板基于平级文件路径注册，引入子目录需要额外修改 manifest 路径逻辑" and verdict "两个文件不足以 justify 目录层级，需要时再拆分（YAGNI）"
**Why it's weak**: The manifest-path argument is not a genuine incremental cost of the subdirectory approach — Scope #7 already commits to modifying `manifest-update-design.md`. The con argues that subdirectories make manifest changes harder, but manifest changes are already planned regardless. The real rejection reason is YAGNI, which is valid, but packaging a non-incremental cost as a con inflates the case against this alternative. The alternatives analysis should present each option on its genuine merits rather than stacking redundant cons.
**What must improve**: Remove or rephrase con ② to reflect the actual incremental cost (e.g., "manifest paths would need to reference `design/db/` instead of `design/`, a trivial change that doesn't add complexity but couples the manifest to a directory structure that may not stabilize"). Alternatively, add a con that is genuinely unique to the subdirectory approach (e.g., "adds a navigation hop for developers who currently open files from the flat `design/` directory").

### Attack 3: Problem Definition — urgency evidence remains partially unquantified

**Where**: Urgency section, "过往教训" bullet: "在 `jlc-schema-alignment` 项目中...`TaskIndex` 表的 `ByID` 查询缺少预期索引直到实施阶段才暴露，导致额外一轮 DDL 变更和数据迁移（postmortem 记录了 6 项改进措施）"
**Why it's weak**: The incident is named and the failure mode is specific — this is a genuine improvement from iteration 1. But the cost remains unquantified. "额外一轮 DDL 变更和数据迁移" describes what happened, not how much it cost. Was it 2 hours of rework? 2 days? Did it block other work? The postmortem's "6 项改进措施" is cited but not summarized — the reader must leave the document to understand the severity. A single data point (e.g., "the missing index caused ~4 hours of unplanned rework across DDL revision, ORM migration, and test updates") would close the gap.
**What must improve**: Add one concrete cost metric from the `jlc-schema-alignment` incident — either time lost, number of files changed in the rework, or a summary of the 6 postmortem items. The evidence is strong qualitatively but missing the quantitative anchor that would make the urgency argument airtight.

---

## Previous Issues Check

| Previous Attack (Iteration 1) | Addressed? | Evidence |
|-------------------------------|------------|----------|
| Attack 1: Urgency is asserted, not demonstrated | ✅ Partially | Urgency section now has 3 structured bullets: rework chain mechanism, `jlc-schema-alignment` incident with specific failure mode (missing `ByID` index on `TaskIndex`), and signal-noise ratio analysis. The incident is named and specific. Still missing: quantified cost (time/hours). |
| Attack 2: Alternatives pros/cons lack analytical depth | ✅ Substantially | Each alternative now has numbered pros/cons with 2-3 items each, including specific mechanisms (e.g., "模板仅支持 `{ fieldName: Type }` 伪代码" for Do nothing con ①). Merit improved from 3/5 to 4/5. Remaining gap: subdirectory con relies on non-incremental cost argument (see Attack 2 above). |
| Attack 3: Optimistic risk likelihood ratings | ✅ Fully | Detection risk re-rated from Low to Medium. New risk added for DDL-to-migration gap (Medium/Medium). eval-design bias re-rated from Low to Medium. 3 of 5 risks now Medium, 1 Low, 1 Medium — no longer optimistic. |
| Deduction: Scope #9 vs SC #9 inconsistency ("独立" and "自动" qualifiers) | ✅ | The wording has been aligned — SC #9 no longer uses qualifiers absent from the scope item. Replaced with specific behavioral description ("任务内容引用 schema.sql 作为设计输入"). |
| Deduction: Single-sentence urgency justification | ✅ | Expanded from 1 sentence to 3 structured bullets with evidence chain. |
| Deduction: Shallow alternatives pros/cons | ✅ | Each alternative now has multi-point pros/cons with specific mechanisms. |

---

## Verdict

- **Score**: 89/100
- **Target**: 90/100
- **Gap**: 1 point
- **Action**: Continue to iteration 3 — primary fix is resolving the scope/mitigation inconsistency (fallback detection mechanism committed in risk table but absent from scope and success criteria). Either scope it or downgrade the mitigation. Secondary: quantify the `jlc-schema-alignment` incident cost with one data point.
