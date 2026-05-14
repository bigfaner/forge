---
date: "2026-05-14"
doc_dir: "docs/proposals/typed-verification-strategies/"
iteration: 3
target_score: "800"
evaluator: Claude (automated, adversarial)
---

# Proposal Eval — Iteration 3

**Score: 833/1000** (target: 800)

```
┌──────────────────────────────────────────────────────────────────────────┐
│                     PROPOSAL QUALITY SCORECARD (1000 pts)                │
├─────────────────────────────────────┬──────────┬──────────┬─────────────┤
│ Dimension                           │ Score    │ Max      │ Status      │
├─────────────────────────────────────┼──────────┼──────────┼─────────────┤
│ 1. Problem Definition               │  92      │  110     │ ✅           │
│    Problem clarity                  │  36/40   │          │             │
│    Evidence provided                │  32/40   │          │             │
│    Urgency justified                │  24/30   │          │             │
├─────────────────────────────────────┼──────────┼──────────┬─────────────┤
│ 2. Solution Clarity                 │  104     │  120     │ ✅           │
│    Approach concrete                │  36/40   │          │             │
│    User-facing behavior             │  36/45   │          │             │
│    Technical direction clear        │  32/35   │          │             │
├─────────────────────────────────────┼──────────┼──────────┬─────────────┤
│ 3. Industry Benchmarking            │  90      │  120     │ ⚠️           │
│    Industry solutions referenced    │  30/40   │          │             │
│    3+ meaningful alternatives       │  22/30   │          │             │
│    Honest trade-off comparison      │  18/25   │          │             │
│    Justified against benchmarks     │  20/25   │          │             │
├─────────────────────────────────────┼──────────┼──────────┬─────────────┤
│ 4. Requirements Completeness        │  92      │  110     │ ✅           │
│    Scenario coverage                │  32/40   │          │             │
│    Non-functional requirements      │  33/40   │          │             │
│    Constraints & dependencies       │  27/30   │          │             │
├─────────────────────────────────────┼──────────┼──────────┬─────────────┤
│ 5. Solution Creativity              │  72      │  100     │ ⚠️           │
│    Novelty over industry baseline   │  28/40   │          │             │
│    Cross-domain inspiration         │  22/35   │          │             │
│    Simplicity of insight            │  22/25   │          │             │
├─────────────────────────────────────┼──────────┼──────────┬─────────────┤
│ 6. Feasibility                      │  87      │  100     │ ✅           │
│    Technical feasibility            │  36/40   │          │             │
│    Resource & timeline feasibility  │  25/30   │          │             │
│    Dependency readiness             │  26/30   │          │             │
├─────────────────────────────────────┼──────────┼──────────┬─────────────┤
│ 7. Scope Definition                 │  72      │  80      │ ✅           │
│    In-scope concrete                │  27/30   │          │             │
│    Out-of-scope explicit            │  23/25   │          │             │
│    Scope bounded                    │  22/25   │          │             │
├─────────────────────────────────────┼──────────┼──────────┬─────────────┤
│ 8. Risk Assessment                  │  76      │  90      │ ✅           │
│    Risks identified (>=3)           │  26/30   │          │             │
│    Likelihood + impact rated        │  25/30   │          │             │
│    Mitigations actionable           │  25/30   │          │             │
├─────────────────────────────────────┼──────────┼──────────┬─────────────┤
│ 9. Success Criteria                 │  65      │  80      │ ⚠️           │
│    Measurable and testable          │  42/55   │          │             │
│    Coverage complete                │  23/25   │          │             │
├─────────────────────────────────────┼──────────┼──────────┬─────────────┤
│ 10. Logical Consistency             │  83      │  90      │ ✅           │
│     Solution <-> Problem            │  33/35   │          │             │
│     Scope <-> Solution <-> Criteria │  28/30   │          │             │
│     Requirements <-> Solution       │  22/25   │          │             │
├─────────────────────────────────────┼──────────┼──────────┬─────────────┤
│ TOTAL                               │  833     │  1000    │             │
└─────────────────────────────────────┴──────────┴──────────┴─────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Problem > Evidence (line 23) | "不生成契约测试、边界值测试、参数组合测试等集成层面的测试" — no concrete bug count or failure data for API/CLI. Only TUI has quantified evidence (11 bugs). The API/CLI claim remains assertion without data. | -8 pts (Evidence) |
| Problem > Urgency (line 27) | "漏掉的 bug 比发现的多" — vague superlative without a ratio or timeline. How many features, over what period? | -6 pts (Urgency) |
| Solution > User-facing (lines 87-91) | gen-test-scripts behavior change described in prose ("e2e 用例生成渲染+golden file 对比函数...integration 用例生成 HTTP 断言或子进程退出码检查函数") but no concrete generated code example shown. The before/after only demonstrates test-cases.md, not actual test script output. | -9 pts (User-facing behavior) |
| Benchmarking > Industry (lines 135-139) | Each tool (Pact, Bats/Cram, VHS, Appium/Espresso) gets 2-3 sentences. No concrete output format comparison. E.g., Bats: "本提案直接借鉴此模式，差异化在于 golden 内容由 AI 自动生成而非人工编写" — the only differentiation cited is "AI generates it," which is a tooling difference, not a conceptual one. | -10 pts (Industry solutions) |
| Benchmarking > Alternatives (line 147) | "Strategy hardcoded in SKILL.md": "新增 profile 需改 skill 代码；策略与 skill 版本强耦合" — two sentences, no depth. What does "改 skill 代码" mean in practice? How many lines? What coupling risks? Straw-man-adjacent. | -8 pts (Alternatives) |
| Benchmarking > Trade-offs (line 148) | "一致性漂移可通过 lint 规则检测（见 NFR），且框架差异使得强制一致性不可取" — "framework differences make forced consistency undesirable" is asserted without demonstration. What specific dimension would differ between go-test and rust-test TUI strategies? No example given. | -7 pts (Trade-offs) |
| Requirements > NFR (line 120) | "策略文件读取增加延迟 < 2s（gen-test-cases 总耗时 < 5%）" — these two thresholds are inconsistent. If gen-test-cases takes 10s total, 5% = 0.5s, not 2s. If it takes 100s, 5% = 5s, which violates the 2s cap. The two constraints conflict depending on baseline duration. | -7 pts (NFR) |
| Requirements > Scenarios (lines 100-105) | No error scenario coverage. What happens when gen-test-cases encounters an unknown capability type? What happens when a golden file mismatches due to a legitimate output change? What happens when strategy file validation fails mid-generation? | -8 pts (Scenario coverage) |
| Creativity > Novelty (line 93) | "策略与 profile 解耦但就近定义" — the core innovation. However, the mechanism is "策略文件在 profile 目录内...gen-test-cases 通过统一的 capability key 查询策略" — this is a standard plugin/configuration pattern (like ESLint config per project). The novelty claim is overstated. | -12 pts (Novelty) |
| Feasibility > Timeline (line 166) | Total 8.5h with zero contingency or uncertainty range. No buffer for the first-ever strategy file taking longer than expected. No mention of review cycles for the 6 strategy files. | -5 pts (Timeline) |
| Risk > Mitigations (line 192) | Risk 2 mitigation: "每次 profile 发布时运行跨 profile lint...同一 capability 验证维度差异 > 50% 触发人工 review" — the 50% threshold is high. If a TUI strategy has 8 dimensions and a profile drops 3, that's 37.5% and passes silently. The threshold should be justified. | -5 pts (Mitigations) |
| Success Criteria (line 199) | "每个 capability section 计数 dimensions/boundary 条目，低于阈值则报错" — what counts as a "dimension 条目"? A bullet point? A numbered item? A sentence with a keyword? The counting method is undefined. A strategy file with "1. 正确性检查" could be counted as 1 dimension or 3 depending on interpretation. | -8 pts (Measurable) |
| Success Criteria (line 203) | "差异体现在 import 列表、assertion 库选择、测试目录结构三方面" — prescribes how the difference must manifest, but does not define a minimum difference threshold. What if imports differ by one line? Is that sufficient? | -5 pts (Measurable) |
| Logical Consistency (line 120) | NFR performance target "< 2s" vs "< 5%" are two different constraints that can conflict (see deduction above). This creates an internal inconsistency between the requirement and itself. | -3 pts (Requirements <-> Solution) |

---

## Attack Points

### Attack 1: Success Criteria — counting method ambiguity undermines measurability

**Where**: Success Criteria (lines 199, 203)
**Quote**: "每个 capability section 计数 dimensions/boundary 条目，低于阈值则报错" and "差异体现在 import 列表、assertion 库选择、测试目录结构三方面"
**Why it's weak**: The proposal demands automated validation of strategy file quality but never defines what constitutes a countable "dimension 条目" or "boundary 条目." In Markdown, a dimension could be a bullet point, a numbered item, a table row, or a sentence. The lint script must parse unstructured Markdown (the strategy files are Markdown, not YAML or JSON), and the counting heuristic is unspecified. Similarly, the e2e vs integration code difference check says "无差异则失败" but does not define the minimum difference. One different import line would pass. The rubric awards 55 pts for "measurable and testable" criteria, and these criteria have quantified thresholds but unquantified measurement methods — the what is specified but the how is not. A test that counts items without knowing what an item is cannot be written.
**What must improve**: Define the counting method explicitly. For strategy file validation: "每个 capability section 下 `### 验证维度` 子标题后的无序列表项（`- ` 开头）计为一个 dimension 条目；`### 边界场景` 子标题后的无序列表项计为一个 boundary 条目。" For e2e/integration difference: "e2e 和 integration 生成的测试代码必须在 import 列表、assertion 函数调用、目录结构三方面均存在至少一处差异。"

### Attack 2: Industry Benchmarking — depth remains shallow despite breadth

**Where**: Alternatives & Industry Benchmarking (lines 129-148)
**Quote**: "Bats / Cram（CLI Golden Testing）：`$ command` + `expected output` 模式。本提案直接借鉴此模式，差异化在于 golden 内容由 AI 自动生成而非人工编写。" and "charmbracelet/vhs（TUI Testing）：固定终端尺寸+输出比对。forge 借鉴此思路但不依赖 VHS 的录制机制"
**Why it's weak**: Five industry tools are cited, each getting 2-3 sentences. The "differentiation" for every tool reduces to "forge uses AI to generate instead of humans writing it." This is not a meaningful technical comparison — it is a single observation repeated five times. No tool gets a concrete output format comparison, assertion style analysis, or test lifecycle mapping. The proposal scores 120 pts for this dimension but treats benchmarking as a citation exercise rather than a design input. The "Centralized config registry" alternative is analyzed with one sentence of downside ("6 个 profile 框架差异大") — this is conclusory, not analytical. A base strategy + per-profile override layer (which would address the framework-difference concern) is not explored.
**What must improve**: Pick the two most relevant benchmarks (Bats for CLI golden testing, VHS for TUI testing) and do a deep comparison: show a concrete example of what Bats expects as input vs what forge's generated test would look like. Explain the specific output format, assertion mechanism, and failure reporting differences. For the alternatives analysis, explore at least one hybrid approach (e.g., base strategy with override patches) before dismissing centralized configuration.

### Attack 3: Requirements Completeness — error scenarios remain unaddressed

**Where**: Requirements Analysis > Key Scenarios (lines 100-105)
**Quote**: The five key scenarios are all happy-path: "gen-test-cases 检测到 `tui` capability → 自动生成 golden file 测试用例..." Each scenario follows the pattern "detection → generation → labeling." There are no error scenarios.
**Why it's weak**: The proposal introduces several new failure modes (strategy file missing, strategy file malformed, capability key mismatch, golden file stale, token budget exceeded, lint failure) but none appear as key scenarios. What happens when gen-test-cases encounters a capability not covered by any strategy file? What happens when the golden file comparison fails in CI because of a legitimate output change (not a bug)? What happens when the "维度覆盖率 < 70%" check fires — what is the developer's recovery workflow? The NFR section addresses compatibility for missing/malformed strategy files (warning + fallback), but the key scenarios section — which should describe the primary interaction paths — only covers the success case. An AI agent following these requirements would not know how to handle the error paths.
**What must improve**: Add at least 3 error/edge scenarios to Key Scenarios: (1) "Strategy file missing for a declared capability: gen-test-cases outputs warning, generates generic test cases without Level field, continues execution." (2) "Capability key mismatch: gen-test-cases aborts with error listing missing/extra keys." (3) "Golden file stale after legitimate output change: developer updates golden file via `gen-test-scripts --update-golden`, CI re-runs."

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1 (Iter 2): Success Criteria — measurability remains the weakest link | ✅ Yes | Substantively addressed. All 6 success criteria now have quantitative thresholds: >=3 dimensions, >=2 boundary scenarios, >=95% Level field coverage, >=1 golden file assertion, >=15% rubric weight with >=70% pass threshold. Score improved from 56/80 to 65/80. Remaining gap is counting method ambiguity (see Attack 1 above). |
| Attack 2 (Iter 2): Risk Assessment — mitigations are not actionable | ✅ Yes | Substantively addressed. All 4 risks now have Owner, Trigger, and Corrective action columns. Risk 1: eval-test-cases rubric auto-check with <70% threshold. Risk 3: token count logging with 15% warning. Risk 4: hardcoded terminal size. Score improved from 67/90 to 76/90. |
| Attack 3 (Iter 2): Requirements Completeness — constraints remain underspecified | ✅ Yes | Substantively addressed. Constraints now specify: Markdown format parsed via LLM prompt context (not structured YAML), section structure requirements (## capability-key + ### sub-headings), capability key consistency validation mechanism, strategy file validity rules with concrete thresholds. Score improved from 85/110 to 92/110. |

---

## Verdict

- **Score**: 833/1000
- **Target**: 800/1000
- **Gap**: Target reached (+33 points above threshold)
- **Action**: Target reached. Proposal is ready for `/write-prd`. Remaining weaknesses (benchmarking depth, error scenario coverage, counting method ambiguity) are acceptable for proposal stage and can be refined during PRD/design.
