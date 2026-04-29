---
date: "2026-04-29"
doc_dir: "docs/features/justfile-e2e-integration/prd/"
iteration: "2"
target_score: "80"
evaluator: Claude (automated, adversarial)
---

# PRD Eval — Iteration 2

**Score: 83/100** (target: 80)

```
┌─────────────────────────────────────────────────────────────────┐
│                       PRD QUALITY SCORECARD                      │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Background & Goals        │  17      │  20      │ ⚠️          │
│    Background three elements │   6/7    │          │            │
│    Goals quantified          │   6/7    │          │            │
│    Logical consistency       │   5/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Flow Diagrams             │  19      │  20      │ ✅          │
│    Mermaid diagram exists    │   7/7    │          │            │
│    Main path complete        │   7/7    │          │            │
│    Decision + error branches │   5/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Functional Specs          │  15      │  20      │ ⚠️          │
│    Tables complete           │   5/7    │          │            │
│    Field descriptions clear  │   6/7    │          │            │
│    Validation rules explicit │   4/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. User Stories              │  19      │  20      │ ✅          │
│    Coverage per user type    │   7/7    │          │            │
│    Format correct            │   7/7    │          │            │
│    AC per story              │   5/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Scope Clarity             │  18      │  20      │ ⚠️          │
│    In-scope concrete         │   7/7    │          │            │
│    Out-of-scope explicit     │   7/7    │          │            │
│    Consistent with specs     │   4/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL (before deductions)    │  88      │  100     │            │
│ Deductions                   │  -5      │          │            │
│ TOTAL                        │  83      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| 需求目标表：第二行说明列 | "工具链变化只需修改 justfile" — 同 iteration 1 的问题，仍是模糊收益描述，无量化指标（节省多少时间？减少多少错误？） | -2 pts |
| Scope In vs. 功能描述 5.2 | `breakdown-tasks` 三个模板（run-e2e-tests.md、gen-test-scripts.md、fix-e2e.md）在 In-Scope 中列出，但 Section 5.2 替换规则表中无对应条目，功能描述与 Scope 直接矛盾 | -3 pts |

---

## Attack Points

### Attack 1: Functional Specs — breakdown-tasks 模板替换规则在功能描述中完全缺失

**Where**: Scope In-Scope 列出 `breakdown-tasks 模板 run-e2e-tests.md：Implementation Notes → just test-e2e --feature <slug>`、`breakdown-tasks 模板 gen-test-scripts.md：Implementation Notes → just e2e-verify --feature <slug>`、`breakdown-tasks 模板 fix-e2e.md：新增修复后验证步骤 → just test-e2e --feature <slug>`

**Why it's weak**: Section 5.2 的两张替换规则表（E2E 命令替换、单元测试/构建命令替换）完全没有 breakdown-tasks 模板的条目。这三个文件的替换规则是什么？原始命令是什么？替换后是什么？影响哪些具体行？这些信息在功能描述中一概缺失。Scope 承诺了交付物，功能描述却没有说明如何交付，验收时无法判断这三个文件是否正确修改。

**What must improve**: 在 Section 5.2 补充 breakdown-tasks 模板的替换条目，明确每个模板文件的原始命令、替换目标和影响位置；或者将这三个文件的替换规则单独列为 Section 5.2.3。

---

### Attack 2: Functional Specs — `<slug>` 参数从未定义，调用方无法确定传入值

**Where**: `just e2e-verify --feature <slug>` 的参数说明：`--feature <slug>（必填，缺省时 exit 1 并提示用法）`

**Why it's weak**: `<slug>` 是整个 PRD 中使用频率最高的参数（出现在 e2e-verify、test-e2e、流程图、User Stories 中），但从未定义它的格式和来源。它是目录名？feature 的短标识符？由谁生成？Agent 从哪里获取这个值？如果 feature 名称包含空格或大写字母，slug 是否需要转换？这个参数的模糊性直接影响 Agent 能否正确调用命令。

**What must improve**: 在 Section 5.1 的 e2e-verify 表格中增加"参数格式"行，明确 `<slug>` 的定义（如：`tests/e2e/` 下的子目录名，小写字母和连字符，由 gen-test-scripts 生成时确定）；并说明 Agent 从何处获取该值。

---

### Attack 3: Flow Diagrams — 命令替换流程图的修复循环无退出条件

**Where**: 命令替换流程图中节点 `N[记录残留位置，返回修复]` → `C[取下一个文件]`

**Why it's weak**: 当验证步骤发现残留原始命令时，流程图将控制流返回到 `C`（取下一个文件），但没有任何退出条件。如果某个文件的原始命令无法被替换（例如命令在注释中、在条件块中、或替换规则表未覆盖的情况），这个循环将无限执行。与主流程图中明确的"重试次数 ≤ 3"退出条件相比，命令替换流程图的错误处理明显不完整。

**What must improve**: 在命令替换流程图的修复循环中增加退出条件：最大重试次数（如 3 次）或人工介入分支；明确"记录残留位置"后的下一步是人工修复还是自动修复，以及自动修复失败时的终止节点。

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1: e2e-setup 表格缺少失败条件和输出格式 | ✅ | Section 5.1 新增了"失败条件"行（`package.json` 不存在 → exit 1）和"输出格式"行（成功/失败输出），两个目标表格结构现已一致 |
| Attack 2: 命令替换流程无可视化 | ✅ | 新增了完整的命令替换流程图，包含逐文件扫描、替换、验证步骤 |
| Attack 3: Scope 与 User Stories 脱节（9 个文件无故事覆盖） | ✅ | 新增 Story 5 覆盖 fix-bug、run-tasks、task-executor、error-fixer、execute-task、record-task、improve-harness 七个文件的构建/测试场景 |

---

## Verdict

- **Score**: 83/100
- **Target**: 80/100
- **Gap**: +3 points (target reached)
- **Action**: Target reached — iteration complete. Remaining issues (breakdown-tasks spec gap, slug definition, loop exit condition) are recommended improvements but do not block acceptance.
