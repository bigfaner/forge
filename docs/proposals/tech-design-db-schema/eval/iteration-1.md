---
date: "2026-05-07"
doc_dir: "docs/proposals/tech-design-db-schema/"
iteration: 1
target_score: 90
evaluator: Claude (automated, adversarial)
---

# Proposal Eval — Iteration 1

**Score: 77/100** (target: 90)

```
┌─────────────────────────────────────────────────────────────────┐
│                    PROPOSAL QUALITY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Problem Definition        │   15     │  20      │ ⚠️         │
│    Problem clarity           │  6/7     │          │            │
│    Evidence provided         │  5/7     │          │            │
│    Urgency justified         │  4/6     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Solution Clarity          │  17      │  20      │ ✅         │
│    Approach concrete         │  7/7     │          │            │
│    User-facing behavior      │  5/7     │          │            │
│    Differentiated            │  5/6     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Alternatives Analysis     │  12      │  15      │ ⚠️         │
│    Alternatives listed (≥2)  │  5/5     │          │            │
│    Pros/cons honest          │  3/5     │          │            │
│    Rationale justified       │  4/5     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Scope Definition          │  13      │  15      │ ⚠️         │
│    In-scope concrete         │  5/5     │          │            │
│    Out-of-scope explicit     │  5/5     │          │            │
│    Scope bounded             │  3/5     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Risk Assessment           │  12      │  15      │ ⚠️         │
│    Risks identified (≥3)     │  4/5     │          │            │
│    Likelihood + impact rated │  4/5     │          │            │
│    Mitigations actionable    │  4/5     │          │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 6. Success Criteria          │  13      │  15      │ ⚠️         │
│    Measurable                │  5/5     │          │            │
│    Coverage complete         │  4/5     │          │            │
│    Testable                  │  4/5     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  77      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Scope #9 vs Success Criteria #9 | Inconsistency: scope says "Element Mapping 表新增数据库变更 → schema 任务映射" but success criterion says "自动生成独立的 schema 执行任务" — "独立" and "自动" are qualifiers not present in the scope item | -3 pts |
| Urgency section | Single-sentence justification with no concrete past incidents or data points backing the claim about schema rework cost | -2 pts (evidence weakness) |
| Alternatives table | Pros/cons are one-phrase each, insufficient depth for trade-off analysis | -2 pts (shallow analysis) |

---

## Attack Points

### Attack 1: Problem Definition — urgency is asserted, not demonstrated

**Where**: "数据库 schema 变更是成本最高的返工类型之一（涉及数据迁移、兼容性）。越早建立独立审批机制，越能减少实施阶段的 schema 返工。"
**Why it's weak**: The entire urgency argument is a single sentence making a broad claim with zero supporting evidence. No past incidents of schema rework, no cost data, no examples of bugs caught late, no user feedback about review quality. The "cost最高的返工类型之一" assertion is stated as fact without backing.
**What must improve**: Add 1-2 concrete examples of schema rework that occurred in this project, or cite specific user/developer feedback about database review quality problems. Even a hypothetical scenario with estimated time cost would strengthen this.

### Attack 2: Alternatives Analysis — pros/cons lack analytical depth

**Where**: The entire alternatives table — e.g., Do nothing has pros "零改动成本" and cons "数据库设计评审质量低，schema 返工成本高"
**Why it's weak**: Every pro and con is a single short phrase. There is no quantitative or even qualitative depth. "评审质量低" is vague — low in what way? What specific quality dimension is lacking? The "Do nothing" alternative's con doesn't even describe the mechanism of failure, just asserts it. The "独立子目录" alternative's con ("子目录过度设计") is a judgment call without justification for why two files don't warrant a subdirectory.
**What must improve**: Expand each alternative to include at least 2-3 pros and cons with specific mechanisms. For "Do nothing", explain the failure mode: e.g., " reviewers miss missing indexes or incorrect constraint names because they are buried in a 200-line markdown section." For "独立子目录", explain why flat is better with a specific reason.

### Attack 3: Risk Assessment — optimistic likelihood ratings and missing migration risk

**Where**: "数据库变更检测误判（漏检或误检）" rated Low likelihood; "Mermaid erDiagram 语法限制" rated Low; "eval-design 评审维度扩展引入新的评分偏差" rated Low
**Why it's weak**: Three out of four risks are rated "Low" likelihood — this is optimism bias. The "检测误判" risk is particularly questionable: the entire detection mechanism relies on a human (or AI) correctly filling `db-schema` in the PRD frontmatter during the `write-prd` stage. If the PRD author forgets or misjudges whether a feature involves DB changes, the system silently skips the entire DB review track. This is a single point of failure in a human/AI step, which should rate at least Medium likelihood. Additionally, there is no risk about migration path: the proposal generates DDL but explicitly excludes migration scripts (out of scope), yet provides no risk analysis for the gap between approved DDL and actual migration execution.
**What must improve**: Re-rate the detection risk to Medium likelihood and add a concrete mitigation for the case where `db-schema` is mislabeled (e.g., a fallback check in `tech-design` that scans for table references). Add a risk entry for the DDL-to-migration gap: approved DDL may not directly translate to migration scripts, creating a handoff risk.

---

## Previous Issues Check

<!-- First iteration — no previous issues -->

---

## Verdict

- **Score**: 77/100
- **Target**: 90/100
- **Gap**: 13 points
- **Action**: Continue to iteration 2 — focus on strengthening urgency evidence, deepening alternatives analysis, and rebalancing risk likelihood ratings
