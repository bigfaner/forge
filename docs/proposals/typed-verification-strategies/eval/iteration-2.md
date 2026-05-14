---
date: "2026-05-14"
doc_dir: "docs/proposals/typed-verification-strategies/"
iteration: 2
target_score: "800"
evaluator: Claude (automated, adversarial)
---

# Proposal Eval — Iteration 2

**Score: 735/1000** (target: 800)

```
┌──────────────────────────────────────────────────────────────────────────┐
│                     PROPOSAL QUALITY SCORECARD (1000 pts)                │
├─────────────────────────────────────┬──────────┬──────────┬─────────────┤
│ Dimension                           │ Score    │ Max      │ Status      │
├─────────────────────────────────────┼──────────┼──────────┼─────────────┤
│ 1. Problem Definition               │  84      │  110     │ ⚠️           │
│    Problem clarity                  │  34/40   │          │             │
│    Evidence provided                │  28/40   │          │             │
│    Urgency justified                │  22/30   │          │             │
├─────────────────────────────────────┼──────────┼──────────┼─────────────┤
│ 2. Solution Clarity                 │  95      │  120     │ ✅           │
│    Approach concrete                │  33/40   │          │             │
│    User-facing behavior             │  32/45   │          │             │
│    Technical direction              │  30/35   │          │             │
├─────────────────────────────────────┼──────────┼──────────┼─────────────┤
│ 3. Industry Benchmarking            │  88      │  120     │ ⚠️           │
│    Industry solutions referenced    │  28/40   │          │             │
│    3+ meaningful alternatives       │  24/30   │          │             │
│    Honest trade-off comparison      │  18/25   │          │             │
│    Justified against benchmarks     │  18/25   │          │             │
├─────────────────────────────────────┼──────────┼──────────┼─────────────┤
│ 4. Requirements Completeness        │  85      │  110     │ ⚠️           │
│    Scenario coverage                │  33/40   │          │             │
│    Non-functional requirements      │  30/40   │          │             │
│    Constraints & dependencies       │  22/30   │          │             │
├─────────────────────────────────────┼──────────┼──────────┼─────────────┤
│ 5. Solution Creativity              │  72      │  100     │ ⚠️           │
│    Novelty over industry baseline   │  28/40   │          │             │
│    Cross-domain inspiration         │  22/35   │          │             │
│    Simplicity of insight            │  22/25   │          │             │
├─────────────────────────────────────┼──────────┼──────────┼─────────────┤
│ 6. Feasibility                      │  85      │  100     │ ✅           │
│    Technical feasibility            │  35/40   │          │             │
│    Resource & timeline feasibility  │  26/30   │          │             │
│    Dependency readiness             │  24/30   │          │             │
├─────────────────────────────────────┼──────────┼──────────┼─────────────┤
│ 7. Scope Definition                 │  71      │  80      │ ✅           │
│    In-scope concrete                │  27/30   │          │             │
│    Out-of-scope explicit            │  22/25   │          │             │
│    Scope bounded                    │  22/25   │          │             │
├─────────────────────────────────────┼──────────┼──────────┼─────────────┤
│ 8. Risk Assessment                  │  67      │  90      │ ⚠️           │
│    Risks identified (>=3)           │  25/30   │          │             │
│    Likelihood + impact rated        │  22/30   │          │             │
│    Mitigations actionable           │  20/30   │          │             │
├─────────────────────────────────────┼──────────┼──────────┼─────────────┤
│ 9. Success Criteria                 │  56      │  80      │ ⚠️           │
│    Measurable and testable          │  36/55   │          │             │
│    Coverage complete                │  20/25   │          │             │
├─────────────────────────────────────┼──────────┼──────────┼─────────────┤
│ 10. Logical Consistency             │  81      │  90      │ ✅           │
│     Solution <-> Problem            │  32/35   │          │             │
│     Scope <-> Solution <-> Criteria │  26/30   │          │             │
│     Requirements <-> Solution       │  23/25   │          │             │
├─────────────────────────────────────┼──────────┼──────────┼─────────────┤
│ TOTAL                               │  735     │  1000    │             │
└─────────────────────────────────────┴──────────┴──────────┴─────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Problem > Evidence (line 23) | "不生成契约测试、边界值测试、参数组合测试等集成层面的测试" — still no concrete bug count or failure data for API/CLI. Only TUI has quantified evidence (11 bugs). The API/CLI claim remains assertion without data. | -12 pts (Evidence) |
| Problem > Urgency (line 27) | "漏掉的 bug 比发现的多" — vague. No metric, no ratio, no timeline of incidents. How many features were affected? Over what period? | -8 pts (Urgency) |
| Solution > User-facing (line 91) | "开发者可观察到：test-cases.md 多出 Level 和 Interface 字段；测试脚本按级别分目录" — the gen-test-scripts behavior change is described in prose but no concrete code example is shown. What does the generated test function look like for e2e vs integration? The before/after only shows test-cases.md, not the actual generated test code output. | -13 pts (User-facing behavior) |
| Benchmarking > Industry (lines 130-136) | Pact, Bats/Cram, VHS, Appium/Espresso are referenced but each gets only 2-3 sentences. No concrete comparison of how forge's approach differs in output format, assertion style, or test lifecycle. For example: how does forge's golden file pattern differ from Bats' `$ command` + `expected output` syntax beyond "AI generates it"? | -12 pts (Industry solutions) |
| Benchmarking > Alternatives (lines 143-144) | "Centralized config registry" alternative: "6 个 profile 框架差异大（Go 用 testdata/，Rust 用 include_str!），集中策略无法表达框架特异手段" — this is a conclusory dismissal. No exploration of whether a base strategy + per-profile override could work. "Strategy hardcoded in SKILL.md": "新增 profile 需改 skill 代码" — not analyzed with any depth. | -7 pts (Alternatives) |
| Benchmarking > Trade-offs (line 145) | "一致性漂移可通过 lint 规则检测（见 NFR），且框架差异使得强制一致性不可取" — the "framework differences make forced consistency undesirable" argument is asserted, not demonstrated. What specific dimension would differ between go-test and rust-test TUI strategies? No example given. | -7 pts (Trade-offs) |
| Requirements > Constraints (line 111) | "不引入新的外部依赖" — still vague. What existing dependencies does gen-test-cases rely on? What is the parsing mechanism for verification-strategies.md? Markdown parsing is non-trivial — is there an existing parser? | -8 pts (Constraints) |
| Risk > Mitigations (line 189) | "策略文件列出的是最小验证维度集，agent 可根据 PRD 追加额外维度" — who validates the agent's additions? What ensures completeness? This mitigation is "the agent will handle it," which is the same as no mitigation. | -10 pts (Mitigations) |
| Risk > Likelihood (line 190) | "6 个 profile 策略文件维护成本" rated Low likelihood — 6 files across different profile maintainers with no synchronization mechanism is at minimum Medium. | -8 pts (Likelihood+Impact) |
| Success Criteria (line 196) | "定义了各 capability 的验证维度和边界场景" — not measurable. "定义了" is a binary checklist item but there is no quality gate. A strategy file that lists one dimension and zero boundary scenarios technically passes. | -8 pts (Measurable) |
| Success Criteria (line 200) | "gen-test-scripts 按 Level 字段选择不同的代码生成策略" — "不同的" is vague. What constitutes "different"? Must the generated code differ in structure, imports, assertions, or just variable names? No acceptance threshold. | -6 pts (Measurable) |
| Success Criteria (line 201) | "eval-test-cases rubric 包含'类型化验证完整度'维度" — what does this dimension score? What is the weight? What score constitutes passing? No definition. | -5 pts (Measurable) |

---

## Attack Points

### Attack 1: Success Criteria — measurability remains the weakest link

**Where**: Success Criteria section (lines 196-201). Six checkbox items, none with quantitative acceptance thresholds.
**Quote**: "定义了各 capability 的验证维度和边界场景" — the verb "定义了" is a existence check, not a quality check. A file with one line per capability would pass. "gen-test-scripts 按 Level 字段选择不同的代码生成策略" — "不同的" is unbounded. How different is different enough?
**Why it's weak**: The rubric awards 55 pts for "measurable and testable" criteria. This proposal's criteria are checklists, not metrics. No criterion specifies: a minimum number of verification dimensions per capability type, a minimum number of boundary scenarios, a percentage of test cases that must have Level fields, or a minimum accuracy score for the "类型化验证完整度" rubric dimension. The gap between "has a file" and "the file is good enough" is entirely undefined.
**What must improve**: Add quantitative thresholds. Examples: "每个 capability 至少包含 5 个验证维度和 3 个边界场景"; "gen-test-cases 生成的 test-cases.md 中 Level 字段覆盖率 >= 95%"; "类型化验证完整度维度权重 >= 15%，通过阈值为 700/1000"。Each criterion should be verifiable by a test or automated check.

### Attack 2: Risk Assessment — mitigations are not actionable

**Where**: Key Risks table (lines 187-193). Four risks listed with mitigations that describe the current design rather than addressing the risk.
**Quote**: "策略文件列出的是最小验证维度集，agent 可根据 PRD 追加额外维度" — this is a restatement of the design, not a mitigation. The risk is "verification dimensions don't match the PRD"; the mitigation is "the agent adds more if needed." But who validates the agent's additions? What happens when the agent's additions are wrong? The same design that caused 11 TUI bugs to pass (agent didn't know to check rendering) is now supposed to self-correct.
**Why it's weak**: Risk #1 (dimension-PRD mismatch) has the highest combined severity (Medium x Medium) yet its mitigation is essentially "trust the agent." Risk #4 (CI terminal differences) is the only risk with a concrete technical mitigation (fixed terminal size). The other three mitigations are either design descriptions or size constraints, not countermeasures. The proposal does not describe any review, validation, or feedback loop for when a mitigation fails.
**What must improve**: Replace design-describing mitigations with actionable countermeasures. For risk #1: "eval-test-cases rubric 增加维度覆盖率检查：对比 PRD 中提到的 interface 行为与策略文件维度，差异 > 30% 则标记为 incomplete." For risk #3: "token 开销监控：gen-test-cases 输出末尾追加 token count，超过阈值时 warning 并建议精简策略。" Each mitigation must name an owner, a trigger, and an action.

### Attack 3: Requirements Completeness — constraints remain underspecified

**Where**: Constraints & Dependencies subsection (lines 109-111). Three bullet points, none with depth.
**Quote**: "不引入新的外部依赖" — this constraint is stated but never operationalized. What does gen-test-cases use to parse verification-strategies.md? Markdown parsing is non-trivial. If it uses the existing LLM-based skill, the "parsing" is prompt engineering — that should be stated explicitly. "策略文件是 profile 目录的一部分，随 profile 版本更新" — what versioning scheme? What happens when a profile major-version bumps — does the strategy file need to change? No analysis.
**Why it's weak**: The NFR subsection was added (good) but the Constraints subsection was not improved from iteration 1. Three vague bullets do not constitute a proper dependency analysis. The proposal depends on: (a) gen-test-cases correctly reading and interpreting markdown strategy files via LLM prompting, (b) the capability key in manifest.yaml being consistent with the capability names used in verification-strategies.md, (c) the 200-line limit being enforced. None of these dependencies are explicitly stated or analyzed.
**What must improve**: Expand constraints to include: parsing mechanism (LLM prompt-based vs. structured format like YAML frontmatter), capability key consistency requirements between manifest.yaml and strategy files, strategy file validation rules (what makes a strategy file "valid"?), version coupling between profile and strategy, and who is responsible for lint enforcement.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1 (Iter 1): Solution Clarity — user-facing behavior entirely absent | ✅ Yes | Added "Developer Workflow — Before & After" section (lines 54-91) with concrete test-cases.md before/after, Level and Interface fields, and gen-test-scripts behavioral description. Score improved from 20/45 to 32/45. |
| Attack 2 (Iter 1): Industry Benchmarking — superficial references, straw-man alternatives | ✅ Partially | Added Pact, Bats/Cram, VHS, Appium/Espresso references with 2-3 sentences each. "Centralized config registry" replaces one straw man. But depth remains shallow and "hardcoded in SKILL.md" alternative is still under-analyzed. Score improved from 65/120 to 88/120. |
| Attack 3 (Iter 1): Requirements Completeness — NFRs missing entirely | ✅ Yes | Full NFR table added (lines 114-122) covering Performance, Compatibility (2 cases), Security, Consistency, and Observability with verification methods. Score improved from 58/110 to 85/110. |

---

## Verdict

- **Score**: 735/1000
- **Target**: 800/1000
- **Gap**: 65 points
- **Action**: Continue to iteration 3. Priority improvements: (1) add quantitative thresholds to all Success Criteria, (2) replace design-describing risk mitigations with actionable countermeasures that name owners and triggers, (3) expand Constraints & Dependencies with parsing mechanism, capability key consistency, and validation rules.
