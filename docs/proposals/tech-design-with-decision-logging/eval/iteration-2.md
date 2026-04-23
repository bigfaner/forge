---
date: "2026-04-22"
doc_dir: "docs/proposals/tech-design-with-decision-logging/"
iteration: 2
target_score: 80
evaluator: Claude (automated, adversarial)
---

# Proposal Eval — Iteration 2

**Score: 77/100** (target: 80)

```
┌─────────────────────────────────────────────────────────────────┐
│                    PROPOSAL QUALITY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 1. Problem Definition        │  15      │  20      │ ⚠️         │
│    Problem clarity           │  6/7     │          │            │
│    Evidence provided         │  4/7     │          │            │
│    Urgency justified         │  5/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 2. Solution Clarity          │  17      │  20      │ ⚠️         │
│    Approach concrete         │  7/7     │          │            │
│    User-facing behavior      │  5/7     │          │            │
│    Differentiated            │  5/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 3. Alternatives Analysis     │  13      │  15      │ ✅         │
│    Alternatives listed (≥2)  │  5/5     │          │            │
│    Pros/cons honest          │  4/5     │          │            │
│    Rationale justified       │  4/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 4. Scope Definition          │  13      │  15      │ ✅         │
│    In-scope concrete         │  4/5     │          │            │
│    Out-of-scope explicit     │  4/5     │          │            │
│    Scope bounded             │  5/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 5. Risk Assessment           │  12      │  15      │ ✅         │
│    Risks identified (≥3)     │  4/5     │          │            │
│    Likelihood + impact rated │  4/5     │          │            │
│    Mitigations actionable    │  4/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 6. Success Criteria          │  10      │  15      │ ⚠️         │
│    Measurable                │  3/5     │          │            │
│    Coverage complete         │  4/5     │          │            │
│    Testable                  │  3/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ TOTAL (incl. inconsistency)  │  77*     │  100     │ ⚠️         │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

*Inconsistency penalty: -3 (in-scope item 8 has no matching success criterion)

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Problem/Evidence | "决策追溯成本持续上升" is asserted without data — how many features exist? How many times has someone needed cross-feature decision lookup? | -2 pts |
| Problem/Evidence | "容易混淆" remains author assertion with no user feedback or incident report to substantiate | -1 pt |
| Solution/User-facing | `/zcode:record-decision` says "通过 AskUserQuestion 收集决策信息：类型（8选1）、决策描述、决策理由、关联feature" but never shows the actual prompt text, error handling for invalid input, or what the user types in response | -2 pts |
| Solution/User-facing | Decision archival says "提示用户确认要归档的决策列表，确认后写入" — confirmation UX is undefined: is it a Yes/No prompt? A numbered list where user selects which decisions to archive? Can the user edit before confirming? | -2 pts |
| Solution/Differentiated | Rationale point 3 ("为什么表格行而非结构化模板") claims "单条记录约 1-2 行 markdown" and "填写成本" arguments but provides no user testing or time measurement — these are design predictions, not validated differentiators | -1 pt |
| Alternatives/Pros-cons | Alternative B ("仅重命名") con: "搜索某技术选型需遍历所有 feature 的 tech-design.md" — while true, no quantification of effort: how many features? how long does a search take today? | -1 pt |
| Alternatives/Rationale | The "Chosen Approach Rationale" section is a significant improvement, but the granularity argument ("8 个分类文件") conflates file count with category count. The rationale never explains why exactly 8 categories is the right number vs. 4 or 12 — what principle guided the split? | -1 pt |
| Scope/In-scope | Items 5 and 6 are rated "M" effort but their deliverables are described as single sentences — "tech-design 流程新增可选'决策归档'步骤" could mean editing one SKILL.md section or rewriting the entire flow; the reader cannot tell scope | -1 pt |
| Scope/Out-of-scope | "docs/DECISIONS.md 的创建" is listed as out-of-scope, but the proposal also says in section 6 "更新 hooks guide 和 exploration 示例中对 docs/DECISIONS.md 的引用，改为 docs/decisions/ 目录" — so DECISIONS.md migration is partially in scope? The boundary is unclear | -1 pt |
| Risk/Identified | Risk 5 ("决策记录与 tech-design.md 内容重复") is a Low/Low filler risk. The proposal's own design already addresses this by keeping records as summaries — this risk is already designed away and adds no analytical value | -1 pt |
| Risk/Likelihood+impact | Impact column uses free-text descriptions ("决策查找变慢", "索引计数错误") rather than a consistent rating scale (High/Medium/Low). The Likelihood column uses H/M/L but Impact does not — asymmetric | -1 pt |
| Risk/Mitigations | Risk 3 mitigation ("命令设计为 3 个必填字段 + 1 个可选字段，预计交互时间 < 30 秒；tech-design 流程结束后自动提示是否归档，减少主动调用的依赖") is a design description, not a risk response. The risk is "low usage" and the "mitigation" is "we made it easy" — this is aspirational, not actionable | -1 pt |
| Success/Measurable | Criterion 1: "可正常调用，流程完整走通" — "完整走通" is subjective. Does it mean no errors? Completes in under N minutes? Produces all expected output files? | -1 pt |
| Success/Measurable | Criterion 4: "docs/decisions/manifest.md 始终反映所有类型文件的决策计数和最近决策" — "始终反映" has no verification method. No acceptance test or validation command specified | -1 pt |
| Success/Coverage | In-scope item 8 ("更新 zcode/CLAUDE.md 中的 skill 列表和 plugins/zcode/SKILLS.md（如存在）中的 skill 注册信息") has no matching success criterion — no criterion verifies CLAUDE.md or SKILLS.md were actually updated | -2 pts (also triggers -3 inconsistency penalty) |
| Success/Testable | Criterion 2 ("若有决策则提示用户确认并归档；无决策则跳过") requires testing two conditional branches. No test strategy is given: how to produce the "no decision" scenario? What test data or mock is needed? | -2 pts |

---

## Attack Points

### Attack 1: Success Criteria — coverage gap and unmeasurable verbiage

**Where**: Success criteria section: criterion 1 ("可正常调用，流程完整走通"), criterion 4 ("始终反映"), and the absence of any criterion covering in-scope item 8.

**Why it's weak**: Three distinct failures here. (a) Criterion 1 uses "完整走通" — a vague phrase that two testers would interpret differently (one might say "no crash", another "all output files present with correct content"). (b) Criterion 4 says "始终反映" — this is an aspirational invariant with no verification method. If a CI script is being proposed as mitigation in the Risk section ("validate-manifest CI 脚本"), why not make that same CI script a success criterion? (c) In-scope item 8 explicitly lists deliverables ("更新 zcode/CLAUDE.md 中的 skill 列表和 plugins/zcode/SKILLS.md（如存在）中的 skill 注册信息") but no success criterion checks whether these files were updated. This is a scope-criteria inconsistency.

**What must improve**: Replace "完整走通" with specific verifiable outcomes (e.g., "produces tech-design.md in the feature directory, decision-entry written to correct category file, manifest.md counters updated"). Replace "始终反映" with "validate-manifest CI script passes" or equivalent automated check. Add a criterion: "zcode/CLAUDE.md skill list contains tech-design entry; SKILLS.md (if exists) registers record-decision command."

### Attack 2: Risk Assessment — one mitigation is still aspirational, not actionable

**Where**: Risk table, risk 3: "`/zcode:record-decision` 使用频率低 | High | 用户仍依赖手动记录或跳过归档 | 命令设计为 3 个必填字段 + 1 个可选字段，预计交互时间 < 30 秒；tech-design 流程结束后自动提示是否归档，减少主动调用的依赖"

**Why it's weak**: The risk is "users won't use the command" (rated High likelihood). The mitigation is "we designed it to be easy and we prompt users." This is not a risk mitigation — it is the product design itself. A true mitigation for low adoption would be: tracking usage metrics, adding a nudge after N unarchived decisions, creating a periodic reminder, or making archival the default with opt-out. "We made it simple so people will use it" is exactly the optimistic assumption that risk assessment is supposed to challenge, not reinforce. Meanwhile, the Impact column for this risk uses free text ("用户仍依赖手动记录或跳过归档") instead of a severity rating, breaking consistency with the Likelihood column's H/M/L scale.

**What must improve**: Replace the aspirational mitigation with a concrete contingency: e.g., "If after 10 features no decisions are archived via the command, add auto-archival prompt as default step (not optional) in tech-design flow" or "Add usage tracking to determine if the command is being invoked." Standardize the Impact column to use High/Medium/Low ratings instead of free text.

### Attack 3: Solution — user-facing behavior for two key interactions remains underspecified

**Where**: Section 3 ("决策归档为可选步骤"): "有决策时：tech-design 文档获用户批准后，提示用户确认要归档的决策列表，确认后写入 docs/decisions/" and Section 5 ("新增 /zcode:record-decision 命令"): "通过 AskUserQuestion 收集决策信息：类型（8选1）、决策描述、决策理由、关联feature"

**Why it's weak**: These are the two primary user interactions in the entire proposal, and neither is specified concretely enough for implementation. For decision archival: what does "确认" mean? Does the user see a numbered list and select which decisions to archive? Is it a single Yes/No for all decisions? Can the user edit the decision text before archiving? For record-decision: what does the AskUserQuestion flow look like? Is it one multi-field prompt or sequential single-question prompts? What happens if the user provides an invalid category? These are not edge cases — they are the core user experience, and an implementer would have to guess.

**What must improve**: For the archival confirmation, specify: "Display numbered list of candidate decisions, user types comma-separated numbers to select which to archive, or 'all' or 'none'." For record-decision, show at least one example exchange: the prompt text and a sample user response. Alternatively, add an "Interaction Examples" subsection with concrete prompt/response pairs for both flows.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1 (Risk: missing likelihood + circular mitigations) | ✅ Mostly | Likelihood column added with H/M/L ratings. Circular mitigations replaced with concrete actions: "validate-manifest CI 脚本", "按年份自动分割为 architecture-2026.md". One residual: risk 3 mitigation remains aspirational. |
| Attack 2 (Alternatives: no rationale) | ✅ Yes | New "Chosen Approach Rationale" section added with 3 numbered points covering granularity, centralization, and format tradeoffs. |
| Attack 3 (Scope: unbounded, no phasing) | ✅ Yes | Scope table now includes Effort (S/M), Depends on, and Phase columns. Three phases defined with descriptions. |

---

## Verdict

- **Score**: 77/100
- **Target**: 80/100
- **Gap**: 3 points
- **Action**: Continue to iteration 3 — address success criteria measurability (+3 pts for criterion 1/4 precision + item 8 coverage), specify the two key user interactions (+2 pts for user-facing behavior), and fix the one remaining aspirational risk mitigation (+1 pt). Total recoverable: ~6 pts.
