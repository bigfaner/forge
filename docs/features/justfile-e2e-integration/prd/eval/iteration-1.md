---
date: "2026-04-29"
doc_dir: "docs/features/justfile-e2e-integration/prd/"
iteration: "1"
target_score: "80"
evaluator: Claude (automated, adversarial)
---

# PRD Eval — Iteration 1

**Score: 75/100** (target: 80)

```
┌─────────────────────────────────────────────────────────────────┐
│                       PRD QUALITY SCORECARD                      │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Background & Goals        │  16      │  20      │ ⚠️          │
│    Background three elements │   6/7    │          │            │
│    Goals quantified          │   5/7    │          │            │
│    Logical consistency       │   5/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Flow Diagrams             │  16      │  20      │ ⚠️          │
│    Mermaid diagram exists    │   7/7    │          │            │
│    Main path complete        │   5/7    │          │            │
│    Decision + error branches │   4/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Functional Specs          │  13      │  20      │ ❌          │
│    Tables complete           │   4/7    │          │            │
│    Field descriptions clear  │   5/7    │          │            │
│    Validation rules explicit │   4/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. User Stories              │  18      │  20      │ ✅          │
│    Coverage per user type    │   7/7    │          │            │
│    Format correct            │   7/7    │          │            │
│    AC per story              │   4/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Scope Clarity             │  17      │  20      │ ⚠️          │
│    In-scope concrete         │   7/7    │          │            │
│    Out-of-scope explicit     │   7/7    │          │            │
│    Consistent with specs     │   3/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL (before deductions)    │  80      │  100     │            │
│ Deductions                   │  -5      │          │            │
│ TOTAL                        │  75      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| 需求目标表：说明列 | 第二个目标的说明列写"工具链变化只需修改 justfile"——这是模糊的收益描述，无量化指标（节省多少时间？减少多少错误？） | -2 pts |
| Scope In vs. User Stories | Scope 列出 13 个文件，但 User Stories 只覆盖 4 个场景；fix-bug、run-tasks、task-executor、error-fixer、execute-task、record-task、improve-harness 均在 Scope 内但无对应 User Story | -3 pts |

---

## Attack Points

### Attack 1: Functional Specs — 表格结构不完整，关键错误路径缺失

**Where**: `5.1 新增 just 目标` 中 `just e2e-setup` 表格：`前置条件 | tests/e2e/package.json 存在`

**Why it's weak**: 表格只描述了"前置条件存在时"的正常路径，完全没有说明前置条件不满足时的行为。`package.json` 不存在时 `just e2e-setup` 应该 exit 1 并输出什么？还是静默失败？此外，表格缺少"输出格式"字段——`e2e-setup` 成功时输出什么，失败时输出什么，调用方如何判断执行结果？`e2e-verify` 有明确的成功/失败输出描述，但 `e2e-setup` 没有。两个目标的表格结构不一致。

**What must improve**: 为 `e2e-setup` 补充"失败条件"和"输出格式"两行；明确 `package.json` 不存在时的 exit code 和错误信息；统一两个目标的表格字段结构。

---

### Attack 2: Flow Diagrams — 三条流程只画了一张图，核心替换流程无可视化

**Where**: `业务流程说明` 描述了三条流程："新增目标流程"、"命令替换流程"、"VERIFY 硬门控流程"，但 `业务流程图` 只有一张 Mermaid 图

**Why it's weak**: "命令替换流程"是本 PRD 工作量最大的部分（13 个文件的批量修改），却完全没有流程图。文字描述是："每个 skill/agent 文件中的原始命令（代码块 + 文字描述）被替换为对应 just 命令"——这个流程涉及谁来执行替换、如何验证替换完整性、替换后如何回归测试，这些决策点在图中完全缺失。现有的单张图把"新增目标流程"和"VERIFY 硬门控流程"混合在一起，且图中 `N → K` 的修复循环没有退出条件（无限循环风险）。

**What must improve**: 为"命令替换流程"单独绘制流程图，包含：触发条件 → 逐文件替换 → 验证无原始命令残留 → 完成；为现有图的修复循环添加最大重试次数或人工介入的退出分支。

---

### Attack 3: Scope Clarity — Scope 与 User Stories 严重脱节，9 个文件无故事覆盖

**Where**: Scope In-Scope 列表包含 `fix-bug command`、`run-tasks command`、`task-executor agent`、`error-fixer agent`、`execute-task command`、`record-task SKILL.md`、`improve-harness SKILL.md`（共 7 个文件），以及 3 个 `breakdown-tasks` 模板；User Stories 只有 4 个故事

**Why it's weak**: Story 1 只覆盖 Skill 维护者查看 `run-e2e-tests` 的场景；Stories 2-4 只覆盖 AI Agent 执行 `task-executor`、`gen-test-scripts`、`fix-e2e` 的场景。`fix-bug`、`run-tasks`、`error-fixer`、`execute-task`、`record-task`、`improve-harness` 这 6 个文件的修改没有任何 User Story 说明"谁在什么场景下受益"。这意味着这些修改的验收标准只有 AC 表格中的文字检查，没有业务场景支撑，无法判断修改是否真正解决了用户问题。

**What must improve**: 要么为缺失的文件补充 User Stories（至少合并为"AI Agent 执行构建/测试任务时获得统一命令"一个故事），要么在 Scope 中明确说明这些文件的修改是纯技术一致性工作、不需要独立 User Story，并在 Out-of-Scope 的对立面给出理由。

---

## Previous Issues Check

<!-- Only for iteration > 1 — N/A for iteration 1 -->

---

## Verdict

- **Score**: 75/100
- **Target**: 80/100
- **Gap**: -5 points
- **Action**: Continue to iteration 2 — fix Functional Specs table completeness (Attack 1) and add missing flow diagram for command replacement (Attack 2) to close the gap
