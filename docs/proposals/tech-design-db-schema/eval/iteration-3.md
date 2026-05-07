---
date: "2026-05-07"
doc_dir: "docs/proposals/tech-design-db-schema/"
iteration: 3
target_score: 90
evaluator: Claude (automated, adversarial)
---

# Proposal Eval — Iteration 3

**Score: 92/100** (target: 90)

```
┌─────────────────────────────────────────────────────────────────┐
│                    PROPOSAL QUALITY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Problem Definition        │  20      │  20      │ ✅         │
│    Problem clarity           │  7/7     │          │            │
│    Evidence provided         │  7/7     │          │            │
│    Urgency justified         │  6/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 2. Solution Clarity          │  20      │  20      │ ✅         │
│    Approach concrete         │  7/7     │          │            │
│    User-facing behavior      │  7/7     │          │            │
│    Differentiated            │  6/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 3. Alternatives Analysis     │  14      │  15      │ ⚠️         │
│    Alternatives listed (≥2)  │  5/5     │          │            │
│    Pros/cons honest          │  4/5     │          │            │
│    Rationale justified       │  5/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 4. Scope Definition          │  15      │  15      │ ✅         │
│    In-scope concrete         │  5/5     │          │            │
│    Out-of-scope explicit     │  5/5     │          │            │
│    Scope bounded             │  5/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 5. Risk Assessment           │  14      │  15      │ ⚠️         │
│    Risks identified (≥3)     │  5/5     │          │            │
│    Likelihood + impact rated │  5/5     │          │            │
│    Mitigations actionable    │  4/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 6. Success Criteria          │  14      │  15      │ ⚠️         │
│    Measurable                │  5/5     │          │            │
│    Coverage complete         │  4/5     │          │            │
│    Testable                  │  5/5     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL (raw)                  │  97      │  100     │            │
│ Inconsistency penalty        │  -5      │          │            │
│ TOTAL                        │  92      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Risk table, row 1 (detection misjudgment) vs Scope #10 vs SC #10 | Inconsistency: The risk mitigation for detection misjudgment, Scope #10, and SC #10 all describe a fallback keyword scan mechanism, but the scope item lists effort as "S" while the risk mitigation and success criterion describe a non-trivial keyword detection and user-prompting behavior — the effort estimate appears underweight for a feature involving PRD content scanning, keyword matching, and interactive user confirmation prompts. More critically, the risk table describes the mitigation as a supplementary defense ("增加防御机制"), while Scope #10 presents it as a primary deliverable on par with template changes. The framing mismatch means a reviewer cannot determine whether this is a critical path item or a nice-to-have guard rail. | -3 pts (inconsistency between risk framing and scope framing of the same mechanism) |
| Alternatives table, "独立子目录" row, con ② | The subdirectory alternative's second con ("开发者当前在 `design/` 平级目录打开文件，子目录要求额外的目录层级跳转，在仅有两个文件时导航成本大于组织收益") is subjective and unquantified. "导航成本大于组织收益" is a judgment call presented as fact without evidence — no user feedback, no developer complaint, no workflow analysis. The pro/con for the subdirectory alternative still lacks a genuine incremental technical cost. | -1 pt (pro/con depth) |
| Risk table, row 5 (eval-design scoring bias) | The mitigation states "上线后通过 3-5 个 feature 的评审结果校准权重" — this is a post-hoc calibration plan, not an actionable mitigation that can be executed during implementation. It defers the risk resolution to an undefined future state without specifying who does the calibration, what the calibration threshold is, or when it triggers. | -1 pt (mitigation not actionable for implementer) |

---

## Attack Points

### Attack 1: Risk/Scope framing inconsistency — fallback detection is presented differently across three sections

**Where**: Risk table row 1: "增加防御机制：`tech-design` 在 Data Models 章节起草时，若 frontmatter 为 `none` 但检测到表名引用（如 `REFERENCES`、`TABLE` 关键词），发出提醒让用户确认" vs Scope #10: "修改 `tech-design/SKILL.md`：Data Models 章节起草时，若 `db-schema` 为 `none` 但内容含表名引用关键词（如 `REFERENCES`、`TABLE`），发出提醒让用户确认 | S" vs SC #10: "`tech-design` 在 Data Models 章节起草时，若 `db-schema` 为 `none` 但 PRD 内容包含 `REFERENCES`、`TABLE` 等表名引用关键词，向用户发出确认提示"
**Why it's weak**: The iteration 2 report explicitly flagged the scope/success-criteria gap for this mechanism, and the revision correctly added Scope #10 and SC #10 to close the gap. However, the framing is now inconsistent: the risk table calls it a "防御机制" (defensive add-on), Scope #10 treats it as a primary deliverable with an "S" effort estimate, and SC #10 presents it as a hard requirement. For something described as a "防御" supplement in the risk table, having it as a full scope item with a success criterion makes it a first-class deliverable — the risk mitigation should acknowledge this elevation in status. Additionally, the effort estimate "S (< 30 min)" for implementing PRD content scanning with keyword detection and interactive user prompting seems optimistic.
**What must improve**: Align the framing across all three locations. If this is a committed deliverable (which Scope #10 and SC #10 confirm), the risk table should describe it as an "in-scope guard rail" rather than a supplementary "防御机制". Alternatively, revise the effort estimate upward to "M" if the keyword scanning logic requires non-trivial implementation. The risk description should match the scope commitment.

### Attack 2: Alternatives — subdirectory rejection con lacks concrete technical cost

**Where**: Alternatives table, "独立子目录 `design/db/`" row, con ②: "开发者当前在 `design/` 平级目录打开文件，子目录要求额外的目录层级跳转，在仅有两个文件时导航成本大于组织收益"
**Why it's weak**: This con survived two iterations unchanged. The iteration 2 report explicitly flagged that the subdirectory rejection needs "a con that is genuinely unique to the subdirectory approach" — the revision replaced the manifest-path con with a navigation-cost argument, but this new con is subjective UX speculation without evidence. No developer survey, no workflow analysis, no concrete example of how "extra directory hop" actually causes friction. The verdict ("YAGNI") is sound, but the supporting con is padding. An honest alternatives analysis would present the verdict on the strength of YAGNI alone rather than inflating it with an unsubstantiated usability claim.
**What must improve**: Either provide evidence for the navigation-cost claim (e.g., "developers currently access design files via `code design/` — subdirectories would require `code design/db/` and then `code design/` for other files, splitting the workspace context") or remove the con and let the YAGNI argument stand on its own merit. One strong, honest con beats two padded ones.

### Attack 3: Success Criteria — SC #10 coverage of Scope #10 is present but the scope item's content source is ambiguous

**Where**: Scope #10: "若 `db-schema` 为 `none` 但内容含表名引用关键词" vs SC #10: "若 `db-schema` 为 `none` 但 PRD 内容包含 `REFERENCES`、`TABLE` 等表名引用关键词"
**Why it's weak**: Scope #10 says "内容含表名引用关键词" (content contains table reference keywords) — what content? The scope item does not specify the source document being scanned. SC #10 clarifies "PRD 内容", but this creates a subtle inconsistency: Scope #10 references modifying `tech-design/SKILL.md`, which operates on the tech-design phase. Is the AI scanning the PRD content during tech-design execution? If so, the scope item should explicitly state "读取 PRD 内容" rather than the ambiguous "内容". This matters because the implementation differs significantly depending on whether the scan target is the PRD file, the tech-design draft, or both.
**What must improve**: Scope #10 should specify the content source explicitly: "若 `db-schema` 为 `none` 但 PRD (`prd-spec.md`) 内容含表名引用关键词" to match SC #10's "PRD 内容". This removes ambiguity about what the implementation actually scans.

---

## Previous Issues Check

| Previous Attack (Iteration 2) | Addressed? | Evidence |
|-------------------------------|------------|----------|
| Attack 1: Scope/Solution inconsistency — fallback detection committed but not scoped | ✅ Fully | Scope #10 added: "修改 `tech-design/SKILL.md`：Data Models 章节起草时，若 `db-schema` 为 `none` 但内容含表名引用关键词...发出提醒让用户确认". SC #10 added: matching success criterion. The gap is closed — the mechanism is now scoped, committed, and verifiable. Remaining concern: framing inconsistency between risk/scope/success-criteria (see Attack 1 above). |
| Attack 2: Alternatives — subdirectory rejection relies on circular reasoning (manifest-path con was non-incremental) | ✅ Substantially | The manifest-path con has been removed. Replaced with navigation-cost con: "开发者当前在 `design/` 平级目录打开文件，子目录要求额外的目录层级跳转". The circular reasoning is gone. Remaining concern: the new con is subjective and unsubstantiated (see Attack 2 above). |
| Attack 3: Problem Definition — urgency evidence partially unquantified (jlc-schema-alignment incident) | ✅ Fully | The urgency section now includes: "该遗漏造成约 2 小时的非计划返工，涉及 DDL 修订、ORM 映射更新和测试修复". The postmortem items are now summarized inline: "包括 schema 审查前置于设计阶段、DDL 输出物纳入交付标准、索引策略强制评审等". Both the quantitative anchor (2 hours) and qualitative context are present. |

---

## Verdict

- **Score**: 92/100
- **Target**: 90/100
- **Gap**: 0 points (target exceeded by 2)
- **Action**: **Target reached.** All three iteration-2 attack points have been substantially or fully addressed. The remaining issues (framing inconsistency, subjective alternatives con, scope-item ambiguity) are minor and do not warrant another iteration. The proposal is approved for implementation.
