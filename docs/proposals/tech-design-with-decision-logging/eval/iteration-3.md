---
date: "2026-04-22"
doc_dir: "docs/proposals/tech-design-with-decision-logging/"
iteration: 3
target_score: 80
evaluator: Claude (automated, adversarial)
---

# Proposal Eval — Iteration 3

**Score: 85/100** (target: 80)

```
┌─────────────────────────────────────────────────────────────────┐
│                    PROPOSAL QUALITY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 1. Problem Definition        │  16      │  20      │ ✅         │
│    Problem clarity           │  6/7     │          │            │
│    Evidence provided         │  5/7     │          │            │
│    Urgency justified         │  5/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 2. Solution Clarity          │  18      │  20      │ ✅         │
│    Approach concrete         │  7/7     │          │            │
│    User-facing behavior      │  6/7     │          │            │
│    Differentiated            │  5/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 3. Alternatives Analysis     │  13      │  15      │ ✅         │
│    Alternatives listed (≥2)  │  5/5     │          │            │
│    Pros/cons honest          │  4/5     │          │            │
│    Rationale justified       │  4/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 4. Scope Definition          │  14      │  15      │ ✅         │
│    In-scope concrete         │  5/5     │          │            │
│    Out-of-scope explicit     │  4/5     │          │            │
│    Scope bounded             │  5/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 5. Risk Assessment           │  12      │  15      │ ✅         │
│    Risks identified (≥3)     │  4/5     │          │            │
│    Likelihood + impact rated │  4/5     │          │            │
│    Mitigations actionable    │  4/5     │          │            │
├──────────────────────────────┼──────────┬────────────┤
│ 6. Success Criteria          │  12      │  15      │ ✅         │
│    Measurable                │  4/5     │          │            │
│    Coverage complete         │  4/5     │          │            │
│    Testable                  │  4/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ TOTAL                        │  85      │  100     │ ✅         │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Problem/Evidence | "容易混淆" is still an author assertion — no user feedback, issue report, or incident log substantiates that anyone was actually confused by the naming | -1 pt |
| Problem/Evidence | "决策追溯成本持续上升" has no baseline data — how many features exist today? How many cross-feature lookups have been needed? Without a number, the reader cannot calibrate urgency | -1 pt |
| Solution/User-facing | The `edit:<编号>` feature in the archival interaction says "输入 edit:<编号> 可重新编辑该条决策的 Decision 或 Rationale 字段后再归档" but does not specify how the edit works: does it open the field inline? Present a new prompt? What happens if the user cancels the edit? This is a user-facing behavior gap | -1 pt |
| Solution/Differentiated | Rationale point 3 claims "单条记录约 1-2 行 markdown" and "填写成本" arguments — these are design predictions, not validated measurements. A genuine differentiator would cite actual time saved vs. ADR format from a trial run | -1 pt |
| Alternatives/Pros-cons | Alternative B con: "搜索某技术选型需遍历所有 feature 的 tech-design.md" — true but unquantified. How many features exist? How long does a grep take? Without numbers, this reads as a straw-man argument | -1 pt |
| Alternatives/Rationale | The "Chosen Approach Rationale" explains why 8 categories but never justifies why *these particular 8*. Why is "Error Handling" a top-level category separate from "Architecture"? Why is "Local Dev & Deployment" not split into two? The category design lacks a stated principle | -1 pt |
| Scope/Out-of-scope | "docs/DECISIONS.md 的创建" is out-of-scope, but scope item 7 says "更新 hooks guide 和 exploration 示例中对 docs/DECISIONS.md 的引用，改为 docs/decisions/ 目录" — so migration of DECISIONS.md references *is* in scope but creating it is not. The boundary between "migration of references" and "dealing with the file itself" could be cleaner | -1 pt |
| Risk/Identified | Risk 5 ("决策记录与 tech-design.md 内容重复 | Low | L") is a designed-away risk. The mitigation is literally the design itself: "表格记录为摘要级...字段设计上避免复制长文本". This adds no analytical value — it is self-congratulatory risk theater | -1 pt |
| Risk/Likelihood+impact | Impact column uses inconsistent formatting: "M — 单文件超过 200 行后可读性下降", "H — 用户仍依赖手动记录或跳过归档", "L — 同一决策在两处维护". The descriptive text after the rating is helpful but varies in specificity — risk 2 has a concrete trigger ("索引计数错误") while risk 5 has a vague one ("更新时遗漏一处") | -1 pt |
| Risk/Mitigations | Risk 1 mitigation includes "当任意类型文件超过 50 条时，reference 流程按年份自动分割为 architecture-2026.md 子文件" — this is a future feature not listed in scope. The proposal's scope table has no item for implementing year-based file splitting, making this an unfunded mitigation | -1 pt |
| Success/Measurable | Criterion 9: "若 plugins/zcode/SKILLS.md 存在，则注册了 tech-design 和 record-decision 两个 skill" — the "if exists" conditional makes this criterion non-deterministic. Two runs of the acceptance test could yield different pass/fail results depending on whether SKILLS.md exists in the repo at test time | -1 pt |
| Success/Coverage | No success criterion explicitly verifies that the `manifest.md` Recent Decisions table is populated correctly. Criterion 4 checks counts but not the recent-decisions listing. The Risk section mentions "Recent Decisions 表包含最近 5 条记录" but this verification is buried in a risk mitigation, not surfaced as a standalone criterion | -1 pt |
| Success/Testable | Criterion 5: "决策记录包含 Date、Feature、Decision、Rationale、Source 五个字段" — this checks field presence but not field validity. A record with empty Decision or a Source pointing to a non-existent file would pass this criterion. The testability is superficial | -1 pt |

---

## Attack Points

### Attack 1: Risk Assessment — designed-away risks and unfunded mitigations

**Where**: Risk table, risk 1 mitigation ("当任意类型文件超过 50 条时，reference 流程按年份自动分割为 architecture-2026.md 子文件") and risk 5 ("决策记录与 tech-design.md 内容重复").

**Why it's weak**: Two problems. (a) Risk 1's mitigation proposes implementing year-based auto-splitting of decision files when they exceed 50 entries. This is a feature — it requires changes to the decision-logging reference logic to detect file size, rename/split files, and update the manifest. Yet the scope table has zero items covering this work. An unfunded mitigation is not a mitigation; it is a wish. (b) Risk 5 (content duplication with tech-design.md) is already solved by the design itself — records are "摘要级" with "Decision + Rationale 各一句话" while tech-design.md keeps full analysis. A risk that your own design already eliminates is not a risk; it is self-validation dressed up as analysis. It displaces space that could be used for a genuine risk, such as: what happens when two features make contradictory architecture decisions? The manifest has no conflict-detection mechanism, yet this is listed as out-of-scope without acknowledging it as a risk.

**What must improve**: Either add the year-based auto-split to the scope table as a Phase 3+ item with effort estimate, or demote it from "mitigation" to "future consideration." Replace risk 5 with a genuine uncovered risk — e.g., "conflicting decisions across features with no detection mechanism" or "decision records become stale when tech-design.md is revised but the archived decision is not updated."

### Attack 2: Alternatives — category count is unjustified and alternative B con remains unquantified

**Where**: "Chosen Approach Rationale" point 1 ("为什么 8 个分类文件而非 per-decision 文件"), and Alternative B con ("搜索某技术选型需遍历所有 feature 的 tech-design.md").

**Why it's weak**: The rationale explains why 8 files is better than 1 file or N files, but never explains why the number is 8 specifically. The categories listed are: Architecture, Interface, Data Model, Dependencies, Error Handling, Testing, Security, Local Dev & Deployment. These appear to be copied from the tech-design template's sections rather than derived from actual decision patterns. Is "Error Handling" frequent enough to warrant its own file? Are "Dependencies" decisions common enough to fill a standalone document? Without evidence that these 8 categories match actual decision distribution, the reader must trust the author's intuition. Meanwhile, Alternative B's con ("遍历所有 feature") would be far more persuasive with a concrete number: "Today with 12 features, a cross-feature decision search requires opening 12 files and scanning ~200 lines each" vs. "With the proposed system, one manifest lookup."

**What must improve**: Add a one-sentence principle for the category split (e.g., "categories correspond to the tech-design template sections, ensuring every decision has a natural home") and acknowledge this is a bootstrap choice that may evolve. For Alternative B, add a concrete quantification of the current search cost.

### Attack 3: Success Criteria — conditional criterion and shallow field validation

**Where**: Criterion 5 ("决策记录包含 Date、Feature、Decision、Rationale、Source 五个字段") and Criterion 9 ("若 plugins/zcode/SKILLS.md 存在，则注册了 tech-design 和 record-decision 两个 skill").

**Why it's weak**: (a) Criterion 9 uses a conditional ("if exists") that makes the acceptance test non-deterministic. If SKILLS.md does not exist, the criterion trivially passes without verifying anything about the skill registration. This should either be two criteria (one unconditional for CLAUDE.md, one conditional for SKILLS.md with a stated precondition) or rephrased to always be verifiable: "The skill registration is updated in all existing documentation files that list available skills." (b) Criterion 5 checks that five fields are *present* but says nothing about their content. A record with `| 2026-04-22 | | | | |` (all fields empty) would pass this criterion. A useful test would specify that Decision and Rationale are non-empty strings, Source points to a file that exists, and Date matches ISO format.

**What must improve**: Split criterion 9 into an unconditional CLAUDE.md check and a separate conditional SKILLS.md check, or rephrase to be always-verifiable. For criterion 5, add content validation: "Decision and Rationale fields are non-empty; Source field references an existing file path; Date field matches YYYY-MM-DD format."

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1 (Success Criteria: coverage gap, unmeasurable verbiage, missing item 8 criterion) | ✅ Yes | Criteria now specific: #1 describes exact deliverables, #4 references CI script, #7 uses `grep` verification, #8 and #9 explicitly cover CLAUDE.md/SKILLS.md. The scope-criteria inconsistency from iteration 2 is resolved. |
| Attack 2 (Risk: aspirational mitigation for low command usage) | ✅ Yes | Risk 3 mitigation now includes concrete escalation trigger: "若连续 10 个 feature 未产生任何归档记录，则将归档步骤从'可选'提升为'默认执行'". This is an actionable contingency, not just "we made it easy." |
| Attack 3 (Solution: underspecified user-facing interactions) | ✅ Yes | Decision archival now shows numbered candidate list with comma-separated selection and `all/none/edit:<编号>` options. `record-decision` shows full 4-round interaction example with sample prompts and responses. |

---

## Verdict

- **Score**: 85/100
- **Target**: 80/100
- **Gap**: -5 (above target)
- **Action**: Target reached. The proposal addresses all three prior attack points substantively. Remaining deductions are about depth of analysis (unquantified assertions, designed-away risks, conditional test criteria) rather than structural gaps.

SCORE: 85/100
DIMENSIONS:
  Problem Definition: 16/20
  Solution Clarity: 18/20
  Alternatives Analysis: 13/15
  Scope Definition: 14/15
  Risk Assessment: 12/15
  Success Criteria: 12/15
ATTACKS:
1. Risk Assessment: unfunded mitigation and designed-away risk — "当任意类型文件超过 50 条时，reference 流程按年份自动分割为 architecture-2026.md 子文件" is a feature not listed in scope, and risk 5 ("决策记录与 tech-design.md 内容重复") is already eliminated by the design itself — replace with a genuine uncovered risk like cross-feature decision conflicts
2. Alternatives Analysis: category count unjustified and alternative B con unquantified — "为什么 8 个分类文件" never explains why these specific 8 categories vs. 4 or 12, and Alternative B con "搜索某技术选型需遍历所有 feature 的 tech-design.md" lacks concrete numbers for current search cost
3. Success Criteria: conditional criterion and shallow validation — criterion 9's "若 plugins/zcode/SKILLS.md 存在" makes the test non-deterministic, and criterion 5 only checks field presence without validating content (empty fields would pass)
