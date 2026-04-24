---
date: "2026-04-24"
doc_dir: "docs/proposals/agent-browser-to-playwright/"
iteration: 3
target_score: 90
evaluator: Claude (automated, adversarial)
---

# Proposal Eval — Iteration 3

**Score: 94/100** (target: 90)

```
┌─────────────────────────────────────────────────────────────────┐
│                    PROPOSAL QUALITY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Problem Definition        │  17      │  20      │ ✅         │
│    Problem clarity           │   6/7    │          │            │
│    Evidence provided         │   6/7    │          │            │
│    Urgency justified         │   5/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Solution Clarity          │  18      │  20      │ ✅         │
│    Approach concrete         │   7/7    │          │            │
│    User-facing behavior      │   6/7    │          │            │
│    Differentiated            │   5/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Alternatives Analysis     │  15      │  15      │ ✅         │
│    Alternatives listed (≥2)  │   5/5    │          │            │
│    Pros/cons honest          │   5/5    │          │            │
│    Rationale justified       │   5/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Scope Definition          │  15      │  15      │ ✅         │
│    In-scope concrete         │   5/5    │          │            │
│    Out-of-scope explicit     │   5/5    │          │            │
│    Scope bounded             │   5/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Risk Assessment           │  15      │  15      │ ✅         │
│    Risks identified (≥3)     │   5/5    │          │            │
│    Likelihood + impact rated │   5/5    │          │            │
│    Mitigations actionable    │   5/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 6. Success Criteria          │  14      │  15      │ ✅         │
│    Measurable                │   4/5    │          │            │
│    Coverage complete         │   5/5    │          │            │
│    Testable                  │   5/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  94      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Problem quantification, "不可复现性" bullet | "可能产生不同的编号" is speculative — uses hedging "可能" without a single observed instance of local-vs-CI divergence. An evidence section should cite observed events, not hypothetical scenarios. | -1 (Problem clarity) |
| Problem quantification, CI incidents | Three CI incidents (04-03, 04-12, 04-18) are described with root causes but no CI run IDs or dashboard links. The commit reference (`857c9cc`) is verifiable, but the CI failures remain unsourced assertions. | -1 (Evidence) |
| "按项目路线图，Q3 前预计扩展至 12 个 UI 页面" | "项目路线图" is referenced but not linked or named. Urgency projection rests on this roadmap claim — without knowing which roadmap, the reader cannot verify the growth assumption. | -1 (Urgency) |
| Developer Workflow, entire section | The developer interaction flow is now present (major improvement from iteration 2), but error iteration is underspecified: after a test fails, the doc says "调整 page-map.json 中对应元素的映射" but never explains whether the developer re-runs the full `/zcode:gen-test-scripts` command or can invoke individual steps. The "重新运行即可" phrasing in the degradation block is equally vague — re-run what? | -1 (User-facing behavior) |
| Core Insight, "无歧义场景" | "映射规则是确定性的一对一转换...没有判断空间，没有歧义场景" — this overstates the case. Two buttons with identical `name` on the same page would produce identical locators. The mapping is deterministic only when elements are unique by (role, name) pair. This edge case is not acknowledged. | -1 (Differentiated) |
| Success Criteria #2 vs Developer Workflow Step 3 | Success criterion #2 requires degradation rate < 10%, but the Developer Workflow example shows 14% degradation with "Degradation rate 14% (< 30% threshold, proceeding)". The document uses two different thresholds (10% and 30%) without explaining their relationship — is 10% the final quality bar while 30% is the process gate? The reader is left guessing. | -1 (Measurable) |

---

## Attack Points

### Attack 1: Success Criteria — dual degradation thresholds create confusion

**Where**: Success Criteria #2 states "降级 locator（优先级 5，使用 CSS 选择器）占比 < 10%" while Developer Workflow Step 3 shows "Degradation rate 14% (< 30% threshold, proceeding)" and the Locator Mapping Rules section says "降级率 > 30% 时 SKILL 输出警告"

**Why it's weak**: The document contains two unstated, potentially conflicting thresholds: 30% is the "process gate" (block exploration if exceeded) and 10% is the "quality bar" (success criterion). The developer workflow example shows a 14% degradation rate that passes the 30% gate but fails the 10% success criterion. The document never reconciles this — a developer reading the workflow would conclude 14% is acceptable, while a reviewer checking success criteria would conclude it fails. This ambiguity undermines both the measurability of the success criteria and the clarity of the developer workflow.

**What must improve**: Either (a) align both thresholds to the same number, or (b) explicitly define both thresholds and their relationship: "30% is the hard block during exploration (the SKILL will refuse to proceed). 10% is the quality target for the final output — the SKILL will proceed but flag the result for review." Add a note in the Developer Workflow example explaining that 14% exceeds the quality target and will be flagged accordingly.

### Attack 2: Problem Definition — "不可复现性" evidence remains hypothetical

**Where**: "不可复现性：@eN ref 依赖运行时 DOM 快照顺序，不同环境（本地 vs CI）或不同数据集可能产生不同的编号，导致'本地通过 CI 失败'的不可预测行为。"

**Why it's weak**: This bullet uses "可能" (might/could) without citing a single observed instance. The other three bullets in the evidence section cite concrete events (specific commits, specific dates, specific root causes). This bullet breaks the pattern — it is a theoretical concern masquerading as evidence. If local-vs-CI divergence has never been observed, it belongs in the Risk section, not in the evidence section. If it has been observed, the observation should be cited like the others.

**What must improve**: Either (a) replace with an observed instance ("2026-04-XX: test passed locally but failed in CI due to different DOM ordering caused by X"), or (b) move this point to the Risk Assessment section where hypothetical concerns belong, and replace it in the evidence section with another observed data point.

### Attack 3: Solution Clarity — developer error recovery workflow is incomplete

**Where**: Developer Workflow > "开发者迭代" subsection: "开发者修复方式：调整 page-map.json 中对应元素的映射，或为页面元素补充 aria-label，然后重新生成。"

**Why it's weak**: "重新生成" (regenerate) is vague — does the developer re-run the entire `/zcode:gen-test-scripts` command from Step 1 (re-exploring all pages), or can they skip exploration and only re-run locator mapping + generation? The workflow implies that `page-map.json` can be manually edited, but never specifies whether the SKILL detects an existing `page-map.json` and skips exploration, or always re-explores. Similarly, the degradation block says "修复后重新运行即可" — re-run what exactly? The developer has no way to know whether their manual `page-map.json` edits will be preserved or overwritten. This is the key developer interaction question: "what happens on the second run?" and it is unanswered.

**What must improve**: Add a brief subsection or paragraph explaining re-run behavior: (1) Does the SKILL detect an existing `page-map.json` and offer to skip/reuse it? (2) Can the developer invoke individual steps (e.g., only re-map locators without re-exploring)? (3) Are manual edits to `page-map.json` preserved on re-run? One sentence like "On re-run, if page-map.json exists, the SKILL will ask whether to reuse it or re-explore" would close this gap.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1 (Iter 2): Problem Definition — evidence provenance is missing | ✅ Yes | Evidence now includes specific commit hash (`857c9cc`), branch name, date (2026-04-15), file affected (`ui.spec.ts`), and 3 CI incidents with exact dates and root causes. CI run IDs/URLs still absent, but the improvement from "约 40%" to sourced incidents with dates and specific `@eN` values is substantial. Score improved from 16/20 to 17/20. |
| Attack 2 (Iter 2): Solution Clarity — developer experience is underspecified | ✅ Yes | New "开发者工作流" section added with step-by-step console output examples for all 4 steps, including degradation warning display, blocking behavior at 30%, and error iteration example. Developer sees exactly what the SKILL outputs at each stage. Minor gap remains on re-run behavior. Score improved from 17/20 to 18/20. |
| Attack 3 (Iter 2): Problem Definition — urgency consequence is speculative | ✅ Yes | Previous version used vague "失去信心" projections. Current version provides concrete quantification: "当前项目包含 4 个 UI 页面...Q3 前预计扩展至 12 个 UI 页面...单次 DOM 变更的测试维护成本将从 15-30 分钟增长至约 45-90 分钟". The growth projection is now grounded in specific numbers (4 -> 12 pages, with cost scaling). Score improved from 16/20 (combined with Attack 1) to 17/20. |

---

## Verdict

- **Score**: 94/100
- **Target**: 90/100
- **Gap**: 0 points (target exceeded by 4)
- **Action**: Target reached. Proposal quality is strong across all dimensions. Remaining deductions are minor (threshold naming, hypothetical evidence point, re-run behavior). No further iterations required.
