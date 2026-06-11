---
date: "2026-05-14"
doc_dir: "docs/proposals/typed-verification-strategies/"
iteration: 4
target_score: "800"
evaluator: Claude (automated, adversarial)
---

# Proposal Eval — Iteration 4

**Score: 836/1000** (target: 800)

```
┌──────────────────────────────────────────────────────────────────────────┐
│                     PROPOSAL QUALITY SCORECARD (1000 pts)                │
├─────────────────────────────────────┬──────────┬──────────┬─────────────┤
│ Dimension                           │ Score    │ Max      │ Status      │
├─────────────────────────────────────┼──────────┼──────────┬─────────────┤
│ 1. Problem Definition               │  93      │  110     │ ✅           │
│    Problem clarity                  │  37/40   │          │             │
│    Evidence provided                │  33/40   │          │             │
│    Urgency justified                │  23/30   │          │             │
├─────────────────────────────────────┼──────────┼──────────┬─────────────┤
│ 2. Solution Clarity                 │  106     │  120     │ ✅           │
│    Approach concrete                │  37/40   │          │             │
│    User-facing behavior             │  37/45   │          │             │
│    Technical direction clear        │  32/35   │          │             │
├─────────────────────────────────────┼──────────┼──────────┬─────────────┤
│ 3. Industry Benchmarking            │  96      │  120     │ ⚠️           │
│    Industry solutions referenced    │  34/40   │          │             │
│    3+ meaningful alternatives       │  24/30   │          │             │
│    Honest trade-off comparison      │  19/25   │          │             │
│    Justified against benchmarks     │  19/25   │          │             │
├─────────────────────────────────────┼──────────┼──────────┬─────────────┤
│ 4. Requirements Completeness        │  98      │  110     │ ✅           │
│    Scenario coverage                │  37/40   │          │             │
│    Non-functional requirements      │  35/40   │          │             │
│    Constraints & dependencies       │  26/30   │          │             │
├─────────────────────────────────────┼──────────┼──────────┬─────────────┤
│ 5. Solution Creativity              │  72      │  100     │ ⚠️           │
│    Novelty over industry baseline   │  28/40   │          │             │
│    Cross-domain inspiration         │  22/35   │          │             │
│    Simplicity of insight            │  22/25   │          │             │
├─────────────────────────────────────┼──────────┼──────────┬─────────────┤
│ 6. Feasibility                      │  87      │  100     │ ✅           │
│    Technical feasibility            │  36/40   │          │             │
│    Resource & timeline feasibility  │  26/30   │          │             │
│    Dependency readiness             │  25/30   │          │             │
├─────────────────────────────────────┼──────────┼──────────┬─────────────┤
│ 7. Scope Definition                 │  72      │  80      │ ✅           │
│    In-scope concrete                │  27/30   │          │             │
│    Out-of-scope explicit            │  23/25   │          │             │
│    Scope bounded                    │  22/25   │          │             │
├─────────────────────────────────────┼──────────┼──────────┬─────────────┤
│ 8. Risk Assessment                  │  76      │  90      │ ✅           │
│    Risks identified (≥3)            │  26/30   │          │             │
│    Likelihood + impact rated        │  25/30   │          │             │
│    Mitigations actionable           │  25/30   │          │             │
├─────────────────────────────────────┼──────────┼──────────┬─────────────┤
│ 9. Success Criteria                 │  78      │  80      │ ✅           │
│    Measurable and testable          │  53/55   │          │             │
│    Coverage complete                │  25/25   │          │             │
├─────────────────────────────────────┼──────────┼──────────┬─────────────┤
│ 10. Logical Consistency             │  58      │  90      │ ⚠️           │
│     Solution ↔ Problem              │  30/35   │          │             │
│     Scope ↔ Solution ↔ Criteria     │  13/30   │          │             │
│     Requirements ↔ Solution         │  15/25   │          │             │
├─────────────────────────────────────┼──────────┼──────────┬─────────────┤
│ TOTAL                               │  836     │  1000    │             │
└─────────────────────────────────────┼──────────┴──────────┴─────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Problem > Evidence (line 23) | "当前 profile 声明了 api 和 cli capability，但 gen-test-cases 只生成功能性测试用例，不生成契约测试、边界值测试、参数组合测试等集成层面的测试" — assertion about API/CLI gap remains without concrete bug count or failure data. Only TUI has quantified evidence (11 bugs, 0/11 caught). API/CLI deficiency is stated as a gap in test type coverage but provides no evidence that this gap has caused any actual failures. | -7 pts (Evidence) |
| Problem > Urgency (line 27) | "漏掉的 bug 比发现的多" — vague superlative. How many features over what period? No ratio, no timeline, no data. If urgency rests on "bugs are being missed," the reader needs to know the scale. | -7 pts (Urgency) |
| Solution > User-facing (lines 87-91) | gen-test-scripts behavior change described in prose: "e2e 用例生成渲染+golden file 对比函数...integration 用例生成 HTTP 断言或子进程退出码检查函数." No concrete generated code example is shown. The before/after demonstrates test-cases.md output, but the actual test script output (what a developer would see in their editor) remains abstract. | -8 pts (User-facing behavior) |
| Benchmarking > Industry solutions (line 143) | Appium / Espresso analysis: "Appium 跨平台但慢，Espresso 仅 Android。forge 的 mobile-ui 策略不绑定工具——profile 层决定框架（如 Maestro），verification-strategies.md 只定义验证维度，与 Appium 互补而非竞争。" Two sentences. No concrete comparison of assertion style, test lifecycle, or output format. The "互补而非竞争" conclusion is asserted without analysis. This tool is listed but not analyzed at the same depth as Pact/Bats/VHS. | -6 pts (Industry solutions) |
| Benchmarking > Alternatives (line 152) | "Strategy hardcoded in SKILL.md": "新增 profile 需改 skill 代码；策略与 skill 版本强耦合" — two sentences of downside, no depth. What does "改 skill 代码" mean in practice? How many lines of a SKILL.md would need editing? This is not analyzed at the same depth as the Centralized config alternative. It functions more as a foil to make the chosen option look better. | -6 pts (Alternatives) |
| Benchmarking > Trade-offs (line 148) | Centralized config rejection: "一致性收益被框架差异抵消，混合方案引入的 override 机制在 AI-generation 场景中无价值" — the hybrid approach (base + override) is explored in one sentence and dismissed because "token 开销和 prompt 复杂度双重增加." But no data is offered. How many tokens would a base strategy consume vs a per-profile one? What is the current token budget for gen-test-cases? The rejection is analytical but not quantitative. | -6 pts (Trade-offs) |
| Requirements > NFR (line 124) | "策略文件读取增加延迟 < 2s（gen-test-cases 总耗时 < 5%）" — two thresholds that can conflict. If gen-test-cases takes 10s, 5% = 0.5s, which satisfies <2s. If it takes 100s, 5% = 5s, which violates <2s. These are independent constraints that must both hold, but the baseline duration is unspecified, making them ambiguous. | -5 pts (NFR) |
| Requirements > Scenarios (line 108) | Error scenarios were added (scenarios 6-9), which addresses iteration-3 Attack 3. However, scenario 8 (Golden file staleness) describes a recovery workflow that involves the agent making a judgment call ("若为 intentional behavior change... 若为 regression") but does not specify how this decision is made or what evidence the agent uses. The scenario assumes perfect agent judgment without defining criteria. | -3 pts (Scenario coverage) |
| Creativity > Novelty (line 95) | "策略与 profile 解耦但就近定义：策略文件在 profile 目录内...gen-test-cases 通过统一的 capability key 查询策略" — this is a standard per-project configuration pattern (ESLint per-project .eslintrc, Prettier per-project .prettierrc, tsconfig per-project). The claim of innovation is "proximal but decoupled," which is the standard location for per-project config. The actual novelty is the AI-generation context (strategy file as LLM prompt context rather than structured config), but this is mentioned only as a constraint, not framed as the innovation. | -12 pts (Novelty) |
| Feasibility > Timeline (line 170) | Total 8.5h with zero contingency buffer. First-ever strategy files for 6 profiles, each needing domain expertise in their framework's testing patterns. No review cycle time included. No uncertainty range. | -4 pts (Timeline) |
| Feasibility > Dependency (line 115) | "不引入新的外部依赖" — stated but not analyzed. The proposal depends on gen-test-cases correctly parsing Markdown strategy files via LLM prompt. What happens if the LLM fails to extract the section structure correctly? The parsing mechanism (LLM prompt) is itself a dependency with known reliability characteristics that are not discussed. | -5 pts (Dependency readiness) |
| Risk > Mitigations (line 197) | Risk 2 mitigation threshold: "同一 capability 验证维度差异 > 50% 触发人工 review." If a TUI strategy has 8 dimensions and a profile drops 3, that is 37.5% and passes silently. The 50% threshold is high enough to miss meaningful drift. No justification for why 50% was chosen over 30% or 40%. | -4 pts (Mitigations) |
| Logical Consistency > Scope-Criteria alignment | Success criterion 5 (line 207) specifies gen-test-scripts output validation: "e2e 生成渲染截获 + golden file 比对函数；integration 生成 HTTP 断言或子进程退出码检查函数." But the In Scope section (lines 176-181) lists "gen-test-scripts SKILL.md 增强：按测试级别和 interface 类型选择代码生成策略" — the scope item is vague ("选择代码生成策略") while the success criterion specifies exact output formats. The success criterion is more specific than the scope it measures, meaning the scope does not fully bound what the success criterion tests. | -9 pts (Scope-Solution-Criteria) |
| Logical Consistency > Scope-Criteria alignment | Success criterion 6 (line 208) specifies eval-test-cases rubric changes: "权重 ≥ 15%（总 rubric 1000 分中 ≥ 150 分），通过阈值为该维度得分 ≥ 70%." The In Scope section (line 180) says only "eval-test-cases rubric 更新：新增'类型化验证完整度'评分维度." The scope does not mention the 15% weight threshold or the 70% pass threshold. These are success criteria without corresponding scope items. | -8 pts (Scope-Solution-Criteria) |
| Logical Consistency > Requirements-Solution | Scenario 8 (line 108) describes golden file staleness recovery involving `gen-test-scripts --update-golden` flag. But gen-test-scripts is listed in the In Scope section only as "gen-test-scripts SKILL.md 增强：按测试级别和 interface 类型选择代码生成策略." The `--update-golden` flag is a new CLI behavior that is not listed in scope. This is a solution behavior (error recovery) that exists in requirements but has no scope boundary. | -10 pts (Requirements-Solution) |

---

## Attack Points

### Attack 1: Logical Consistency — scope-criteria-alignment breakdown between success criteria and scope definition

**Where**: Success Criteria (lines 203-208) vs Scope (lines 176-181)
**Why it's weak**: The proposal's scope section and success criteria are misaligned in both directions. Success criterion 5 (line 207) specifies exact output format requirements for gen-test-scripts (e2e must import `os/exec` + `golden`, integration must import `net/http` + `assert`), but the corresponding scope item (line 179) says only "gen-test-scripts SKILL.md 增强：按测试级别和 interface 类型选择代码生成策略" — the scope does not commit to producing specific output formats. Conversely, success criterion 6 (line 208) introduces quantified thresholds (rubric weight >= 15%, pass threshold >= 70%) that have no corresponding scope commitment. And scenario 8 (line 108) introduces `gen-test-scripts --update-golden` as a recovery mechanism, but this new CLI flag is not listed anywhere in the In Scope items. The rubric requires "no contradictions between what we're building, what's in scope, and how we measure success" (30 pts), and three separate alignment failures exist.
**What must improve**: (1) Expand the gen-test-scripts scope item to include "为 e2e/integration 生成结构不同的测试代码（含 import 列表、assertion 库、目录结构差异）及 `--update-golden` flag 支持 golden file 更新." (2) Add to the eval-test-cases scope item: "新增维度权重 ≥ 15%，通过阈值 ≥ 70%." (3) Ensure every success criterion traces to a scope item and every scope item has at least one success criterion validating it.

### Attack 2: Industry Benchmarking — hybrid alternative dismissed without evidence, depth still uneven

**Where**: Alternatives & Industry Benchmarking (lines 145-152)
**Why it's weak**: The centralized config alternative (line 148) now includes a hybrid approach analysis (base strategy + per-profile override patches), which was recommended in iteration-3 Attack 2. However, the dismissal is: "token 开销和 prompt 复杂度双重增加——而 AI 直接读取完整的 profile-local 策略文件无需合并逻辑." This is conclusory. No token count estimate is given for either approach. The current proposal already requires gen-test-cases to read strategy files via LLM prompt — how many tokens does a strategy file consume? What is the gen-test-cases token budget? Without this data, "token 开销增加" is a vague claim. Additionally, the Appium/Espresso entry (line 143) remains at 2 sentences, far less depth than the 8+ sentences given to Pact, Bats, and VHS. The rubric requires "each alternative must be a genuinely different approach, not a straw man" (30 pts) and "honest trade-off comparison based on actual project constraints" (25 pts). The hybrid dismissal lacks project-specific data.
**What must improve**: (1) Provide a token estimate: "当前 gen-test-cases 平均消耗 ~X tokens。一个策略文件约 ~Y tokens。base + override 方案需要读取 base (~Y tokens) + override (~Z tokens) + 合并指令 (~W tokens)，总计 ~Y+Z+W tokens vs profile-local 方案 ~Y tokens。" (2) Either expand Appium/Espresso to match the depth of Pact/Bats/VHS, or explain why mobile-ui is less relevant (e.g., "forge 目前无 mobile-ui profile，此 benchmark 为前瞻性参考").

### Attack 3: Requirements Completeness — golden file staleness recovery is underspecified

**Where**: Requirements Analysis > Key Scenarios (line 108)
**Why it's weak**: Scenario 8 describes golden file staleness: "agent review diff 后：(a) 若为 intentional behavior change → 运行 `gen-test-scripts --update-golden` 更新 golden file 并重新执行测试；(b) 若为 regression → 以 diff 作为证据 file bug 并修复代码。" This recovery workflow depends on the agent making a correct judgment between intentional change and regression, but provides no criteria for this decision. What does the agent compare? The PRD? The commit message? The test case description? In an automated pipeline, this is a critical branching point. Furthermore, `gen-test-scripts --update-golden` is a new CLI flag introduced in a scenario description but not listed in the scope, not mentioned in the solution description, and not covered by any success criterion. It is an orphan requirement — a behavior the system must support but that is not scoped, designed, or validated. The rubric requires "requirements map cleanly to the proposed solution — no orphan requirements" (25 pts for Requirements-Solution coherence).
**What must improve**: (1) Define the agent's decision criteria: "agent 比对 diff 与 PRD 中相关 interface 行为描述：若 diff 对应的 PRD 行为描述已变更 → intentional；若 PRD 未变更 → regression。" (2) Add `--update-golden` flag to the gen-test-scripts scope item. (3) Add a success criterion: "golden file 更新后重新执行测试，测试通过率 100%。"

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1 (Iter 3): Success Criteria — counting method ambiguity undermines measurability | ✅ Yes | Substantively addressed. Line 203 now defines the exact counting method: "每个 `## <capability-key>` 标题下，`### 验证维度` 子标题之后、下一个 `###` 或 `##` 标题之前的所有无序列表项（`- ` 开头）计为一个 dimension 条目" with a concrete example. Line 207 now requires "三方面均必须至少存在一处差异" with example imports. Score improved from 65/80 to 78/80. |
| Attack 2 (Iter 3): Industry Benchmarking — depth remains shallow despite breadth | ✅ Partial | Pact (line 137), Bats (line 139), and VHS (line 141) now have deep 3-aspect comparisons (assertion style, output format, lifecycle). This is a significant improvement. However, Appium/Espresso (line 143) remains at 2 sentences, and the hybrid alternative dismissal (line 148) still lacks quantitative token estimates. Score improved from 90/120 to 96/120. Partial credit — depth improved for 3/5 tools but uneven. |
| Attack 3 (Iter 3): Requirements Completeness — error scenarios remain unaddressed | ✅ Yes | Substantively addressed. Three error/edge scenarios added: scenario 6 (Unknown capability type — warning + fallback), scenario 7 (Capability key mismatch — abort with error), scenario 9 (Strategy file parse failure — warning + fallback). Each has a defined recovery path. Score improved from 92/110 to 98/110. |

---

## Verdict

- **Score**: 836/1000
- **Target**: 800/1000
- **Gap**: Target reached (+36 points above threshold)
- **Action**: Target reached. Proposal is ready for `/write-prd`. The remaining weaknesses are: (1) scope-criteria alignment gaps where success criteria specify more than scope commits to, (2) uneven benchmarking depth across all cited tools, and (3) underspecified error recovery decision criteria. These are acceptable for proposal stage and should be refined during PRD/design.
