---
created: "2026-05-10"
status: draft
supersedes: docs/forensics/doc-reviser-stuck/report.md
related: docs/forensics/reviser-stuck/report.md
---

# Proposal: test-cases.md 迁移 YAML + 三阶段修订管道

## 1. 问题陈述

### 1.1 现象

`/eval-test-cases` 技能在修订 `test-cases.md` 时，doc-reviser subagent 每次卡死 30+ 分钟，用户被迫中断。三个连续会话（a2023113、98648a17、6860d3e2）均复现相同模式：

```
读文件（7-18s）→ 推理生成（30-43min）→ 用户中断 → 0 个 Edit 完成
```

### 1.2 根因

doc-reviser 需要生成 **83 个 Edit tool call** 来修改一个 32KB / 802 行 / 48 个 TC 的文件。

| 指标 | 成功的 task-executor agent | doc-reviser 任务 |
|------|--------------------------|-----------------|
| 最大 Edit 数 | 26 | 83（3.2x 超限）|
| 推理时间 | ~9 min | ~30 min（未完成）|
| 输出 token 速率 | ~30 tok/s | ~30 tok/s（相同）|
| 隐性推理开销 | 1x | 3.3x（exact-string 回忆）|

**瓶颈不是 agent 行为，而是 Edit tool 的精确字符串匹配要求**：每个 Edit 调用必须从 32KB 文件中逐字符回忆 exact old_string，83 次精确回忆的推理成本是普通生成的 3.3 倍。

### 1.3 之前的修复无效

前一份 forensic report 提出的 3 项修复已全部实施：

| 修复 | 状态 | 效果 |
|------|------|------|
| 禁止 scope-creep（读 PRD 文件） | 已实施 | agent 只读不写 PRD |
| 禁止 TaskCreate/TaskUpdate | 已实施 | agent 不再创建任务 |
| 要求用 Edit 不用 Write | 已实施 | 但 Edit 的精确匹配才是瓶颈 |

Agent 行为已正确，但 83-Edit 的工作量超出模型单次响应能力。

---

## 2. 提案概述

**将 test-cases.md 从 markdown table 格式迁移到 YAML，并重构 eval-test-cases 的修订管道为三阶段分离。**

### 为什么是 YAML

当前 markdown table 格式的问题：

```markdown
### TC-UI-004: Slow tool call node highlighted

| Field | Value |
|-------|-------|
| **Test ID** | ui/tree/slow-node-highlight |
| **Target** | internal/model/calltree |
| **Route** | ??? |
| **Element** | ??? |
| **Steps** | 1. Press `j` to navigate...\| 2. ... |
| **Expected** | The node is rendered with yellow color/highlight |
```

- 添加一个字段 = 插入一行 table row，需要 regex 解析
- Model 修改 = Edit tool 精确匹配 old_string，容易出错
- 多行字段（Steps、Expected）用 `|` 拼接，不直观

YAML 格式：

```yaml
- id: TC-UI-004
  title: "Slow tool call node highlighted"
  test_id: ui/tree/slow-node-highlight
  target: internal/model/calltree
  route: main-tui/call-tree          # ← 加字段就是加一行
  element: "[role='treeitem'][aria-label*='Bash']"
  type: UI
  source: "Story 2 AC: ..."
  priority: P0
  preconditions:
    - "A session file with a tool call duration >= 30 seconds"
  steps:
    - "Press `j` to navigate to the slow tool call node"
  expected: "Node text contains ANSI yellow escape (\\033[33m)"
```

| 操作 | Markdown table | YAML |
|------|---------------|------|
| 添加 Route 字段 | regex 解析 table → 插入行 | 加一行 `route: value` |
| Model 修改文件 | Edit tool（精确匹配） | 输出 value-map，script merge |
| Script 解析 | fragile regex | `yaml.load()` 一行 |
| 多行字段 | `\|` 管道拼接 | 原生 list 和多行字符串 |
| Diff 可读性 | 改一个字段 diff 整行 | 改一行 diff 一行 |

---

## 3. 三阶段修订管道

### 3.1 架构

```
当前: scorer → [reviser: 83 Edit calls, 30min stall] → scorer → ...

提案: scorer → script(preprocess) → model(generate values) → script(apply) → scorer → ...
```

### 3.2 Phase 1: Script 自动提取（无需 model，<1s）

Script 解析 YAML 文件，自动填充可推导的值：

```bash
node preprocess.js --input test-cases.yaml --out test-cases-preprocessed.yaml
```

| 自动提取项 | 来源 | 数量 |
|-----------|------|------|
| CLI Route 值 | 从 Steps 字段提取命令 | 3 |
| API Route 值 | 从 Steps 字段提取函数调用 | ~11 |
| Traceability table | 遍历所有 TC 生成 flat 表 | 48 行 |
| Route Validation table | 遍历所有 TC 生成骨架 | 1 |

```yaml
# 自动提取示例
- id: TC-CLI-002
  steps:
    - "Run `agent-forensic --lang en`"
  # script 自动推导:
  route: "agent-forensic --lang en"
```

### 3.3 Phase 2: Model 生成结构化数据（纯输出，~3-4 min）

Model 读取 preprocessed YAML + attack points，**只输出需要填写的值**。

**不再使用 Edit tool。** 这跟 scorer 的工作模式一样——读取文件，生成结构化文本，不修改文件。

Model 输出格式（value-map）：

```yaml
# model-output.yaml — 由 model 生成

tc_values:
  TC-UI-001:
    route: "main-tui/sessions-panel"
    element: "[data-testid='session-item']:nth-child(1)"
  TC-UI-002:
    route: "main-tui/call-tree"
    element: "[role='treeitem']:first-child"
  # ... 30 UI TCs

  TC-API-005:
    route: "detector.IsSlowToolCall(call model.ToolCall) bool"
  # ... 4 个无法自动推导的 API Routes

expected_rewrites:
  TC-CLI-002: "Status bar displays 'j/k:nav Enter:expand Tab:detail /:search' (English labels)"
  TC-UI-004: "Node text contains ANSI yellow escape sequence (\\033[33m)"
  TC-UI-005: "Node text contains ANSI red escape sequence (\\033[31m)"
  # ... 6 个 vague Expected 重写

new_test_cases:
  - id: TC-UI-031
    title: "Replay slow step yellow highlight in timeline"
    test_id: "ui/replay/slow-step-highlight"
    target: "internal/ui/timeline"
    route: "main-tui/timeline"
    element: "[data-testid='timeline-node'][data-duration='slow']"
    type: UI
    source: "Story 5 AC-3: 耗时 >=30 秒的步骤在时间轴上以黄色标记高亮显示"
    priority: P1
    preconditions:
      - "A session file with a tool call duration >= 30 seconds"
      - "Application is in replay mode"
    steps:
      - "Press `n` to enter replay mode"
      - "Navigate to a step with duration >= 30s"
    expected: "Timeline node displays with yellow (\\033[33m) highlight color"
  # ... 4 个新 TCs

route_validation:
  - route: "main-tui/sessions-panel"
    tc_ids: [TC-UI-001, TC-UI-003]
    verified: false
  # ...
```

### 3.4 Phase 3: Script 应用所有值（<1s）

```bash
node apply-values.js \
  --base test-cases-preprocessed.yaml \
  --values model-output.yaml \
  --out test-cases.yaml
```

Script 执行：
1. 读取 base YAML + model value-map
2. 按 TC ID merge route/element 值
3. 替换 Expected 字段
4. 插入新 TCs
5. 替换 traceability table 和 route validation table
6. 写回 test-cases.yaml

### 3.5 性能对比

| 阶段 | 当前 | 提案 |
|------|------|------|
| Phase 1 (script) | — | <1s |
| Phase 2 (model) | 30+ min（stall） | ~3-4 min（纯生成） |
| Phase 3 (script) | — | <1s |
| **总计** | **30+ min（未完成）** | **~4 min** |
| 加速比 | — | **~8x** |

---

## 4. 文件格式迁移

### 4.1 当前格式（markdown table）

```markdown
---
feature: "agent-forensic"
generated: "2026-05-10"
source: prd/prd-spec.md, prd/prd-user-stories.md, prd/prd-ui-functions.md
---

# Test Cases: Agent Forensic

## CLI Tests

### TC-CLI-001: Missing ~/.claude/ directory shows error

| Field | Value |
|-------|-------|
| **Test ID** | cli/launch/missing-claude-dir |
| **Target** | cmd/root |
| **Type** | CLI |
| **Source** | Story 8 AC: ... |
| **Priority** | P0 |
| **Pre-conditions** | ... |
| **Steps** | 1. ... |
| **Expected** | ... |

---

## Traceability Matrix
| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| ... | ... | ... | ... | ... |
```

### 4.2 目标格式（YAML）

```yaml
---
feature: "agent-forensic"
generated: "2026-05-10"
sources:
  - prd/prd-spec.md
  - prd/prd-user-stories.md
  - prd/prd-ui-functions.md
---

# Test Cases: Agent Forensic

## CLI Tests

- id: TC-CLI-001
  title: "Missing ~/.claude/ directory shows error and exits"
  test_id: cli/launch/missing-claude-dir
  target: cmd/root
  type: CLI
  source: "Story 8 AC: Given ~/.claude/ 目录不存在, When 启动 agent-forensic, Then 显示错误提示并退出"
  priority: P0
  preconditions:
    - "~/.claude/ directory does not exist (or is set to a non-existent path via env/config)"
  steps:
    - "Ensure no ~/.claude/ directory exists (or set a custom non-existent path)"
    - "Run `agent-forensic`"
  expected: "Application prints error message '未找到 ~/.claude/ 目录' (or its English equivalent based on --lang) and exits with non-zero code"

## API Tests

- id: TC-API-001
  title: "Parse valid JSONL session file"
  ...

## UI Tests

- id: TC-UI-001
  title: "Session list displays available JSONL files"
  route: main-tui/sessions-panel
  element: "[data-testid='session-item']"
  ...

## Traceability

- tc_id: TC-CLI-001
  source: "Story 8 AC"
  type: CLI
  target: cmd/root
  priority: P0
- tc_id: TC-CLI-002
  source: "prd-spec.md i18n Requirements"
  type: CLI
  target: cmd/root
  priority: P1
...

## Route Validation

- route: "main-tui/sessions-panel"
  tc_ids: [TC-UI-001, TC-UI-003]
  verified: false
...
```

### 4.3 格式迁移工具

```bash
# 一次性迁移
node migrate-format.js \
  --input docs/features/<slug>/testing/test-cases.md \
  --output docs/features/<slug>/testing/test-cases.yaml
```

---

## 5. 下游影响评估

### 5.1 影响矩阵

| 技能 | 角色 | 影响 | 改动范围 |
|------|------|------|---------|
| gen-test-cases | 创建 test-cases 文件 | **中** | 模板从 markdown table 改为 YAML；SKILL.md 输出指令调整 |
| eval-test-cases | 评估质量 | **低** | rubric 按字段名评分，不依赖格式；scorer 读内容不解析格式；reviser 改为 value-map 输出模式 |
| gen-test-scripts | 消费 test-cases | **低** | LLM 解析 YAML 比解析 markdown table 更容易；字段名不变 |
| breakdown-tasks | 编排管道 | **无** | 只引用文件路径 |
| quick-tasks | 快速管道 | **无** | 只引用文件路径 |
| run-e2e-tests | 执行脚本 | **无** | 不读 test-cases |

### 5.2 关键发现：无解析代码

所有对 test-cases.md 的"解析"都是 **LLM agent 通过 prompt 指令完成的**，没有任何脚本或 parser 代码需要重写。改动范围限于：

1. **模板文件** — markdown table → YAML 示例
2. **SKILL.md 指令** — 输出/读取格式描述
3. **新增脚本** — preprocess.js、apply-values.js、migrate-format.js

### 5.3 rubric 兼容性

eval-test-cases 的 rubric 定义了 5 个维度，按字段名评分：

| 维度 | 检查的字段 | YAML 兼容？ |
|------|-----------|------------|
| PRD Traceability (25pt) | Source, traceability table | 是（字段名不变） |
| Step Actionability (25pt) | Steps, Expected, Pre-conditions | 是 |
| Route & Element Accuracy (20pt) | Route, Element | 是 |
| Completeness (20pt) | Type, boundary cases | 是 |
| Structure & ID Integrity (10pt) | TC IDs, classification, summary | 需微调：summary 格式描述 |

唯一需要调整的是 rubric 中 "required sections" 的格式描述（从 markdown table 改为 YAML list）。

### 5.4 gen-test-scripts 兼容性

gen-test-scripts SKILL.md 的提取指令：

> "Parse each test case — extract TC ID, title, type, route, feature, pre-conditions, steps, expected result, priority. Group by type."

字段名完全相同。LLM 从 YAML 提取结构化数据比从 markdown table 提取更准确。

Route 和 Element 字段的使用方式不变：
- Route → `page.goto()` 或 `curl()` 调用
- Element → sitemap.json 元素匹配
- Type → 模板选择

---

## 6. 实施计划

### Phase A: 基础设施（不改现有管道）

| 步骤 | 文件 | 内容 |
|------|------|------|
| A1 | `skills/eval-test-cases/bin/preprocess.js` | YAML 解析 + 自动提取 CLI/API Route + 生成 traceability table |
| A2 | `skills/eval-test-cases/bin/apply-values.js` | 读取 base YAML + model value-map → merge → 写回 |
| A3 | `skills/eval-test-cases/bin/migrate-format.js` | markdown table → YAML 一次性迁移工具 |

### Phase B: 格式迁移

| 步骤 | 文件 | 内容 |
|------|------|------|
| B1 | `skills/gen-test-cases/templates/test-cases.yaml` | 新模板（替代 test-cases.md） |
| B2 | `skills/gen-test-cases/SKILL.md` | 输出指令改为 YAML |
| B3 | 迁移现有文件 | `migrate-format.js` 转换所有现有 test-cases.md |

### Phase C: 管道重构

| 步骤 | 文件 | 内容 |
|------|------|------|
| C1 | `skills/eval-test-cases/SKILL.md` Step 4 | 三阶段修订：preprocess → model generate → apply |
| C2 | `agents/doc-reviser.md` | 新增 value-map 输出模式；或新建 `doc-value-generator` agent |
| C3 | `skills/eval-test-cases/templates/rubric.md` | required sections 格式描述微调 |

### 依赖关系

```
A1, A2, A3 (并行) → B3 (依赖 A3) → C1 (依赖 A1, A2, B3)
B1, B2 (并行) → B3
C2, C3 (并行) → C1
```

### 风险

| 风险 | 缓解 |
|------|------|
| 现有 test-cases.md 文件需要迁移 | migrate-format.js 一次性转换；保留 .md 备份 |
| YAML 缩进错误导致解析失败 | apply-values.js 包含 schema 验证；失败时回退到原文件 |
| gen-test-scripts agent 不适配 YAML | 字段名不变，LLM 天然适配；如需可加 format hint |
| 多个 feature 同时进行管道迁移 | Phase A 的脚本是增量的，不影响现有 markdown 管道 |

---

## 7. 预期效果

| 指标 | 当前 | 提案后 |
|------|------|--------|
| eval-test-cases 单次迭代时间 | 30+ min（stall） | ~5 min（scorer 2min + model 3min + script <1s） |
| Model 输出方式 | 83 Edit tool calls（精确匹配） | 1 份 value-map YAML（纯生成） |
| 文件修改方式 | Edit tool（fragile） | Script merge（reliable） |
| 下游解析可靠性 | LLM 解析 markdown table | LLM 解析 YAML（更准确） |
| 可扩展性 | 文件越大越慢（O(n) Edit calls） | Model 输出与 TC 数量线性，Script O(n) 即时 |
