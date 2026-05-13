# Hard Rules Compliance: 任务文件 Hard Rules 强制遵从机制

## Background

`task-executor` subagent 在执行 T-test-3（run-e2e-tests）时，读取了任务文件中的 `## Implementation Notes`（明确指定 `just test-e2e --feature <slug>`），但选择自行构造 `npx playwright test` 命令运行测试。这导致缺少 server 启动、环境变量设置、testIgnore 绕过等 `just test-e2e` 提供的保障。

**根本原因**：task-executor 对 `## Implementation Notes` 的遵从优先级低于其自身行为模式。当任务文件的建议命令与 agent 直觉冲突时，agent 按直觉走。

**教训**（`docs/lessons/gotcha-task-executor-ignores-implementation-notes.md`）建议将关键命令从 `## Implementation Notes` 提升为 `## Hard Rules`。但目前 0 个任务文件使用了 `## Hard Rules`，且模板、agent 定义、prompt 模板均不支持该机制。

## Proposal

在 task-executor 的三层架构中引入 `## Hard Rules` 强制遵从机制，确保带有 MUST/MUST NOT 指令的任务约束不会被 agent 忽略或替代。

### 改动范围

| 层级 | 文件 | 改动 |
|------|------|------|
| Agent 定义 | `plugins/forge/agents/task-executor.md` | 新增 Hard Constraints 规则 7 |
| 任务模板（业务） | `breakdown-tasks/templates/task.md` | 新增 `## Hard Rules` section |
| 任务模板（业务） | `quick-tasks/templates/task.md` | 新增 `## Hard Rules` section |
| 任务模板（gate） | `breakdown-tasks/templates/gate-task.md` | "no new code" 从 Impl Notes 迁移到 Hard Rules |
| 技能指令 | `breakdown-tasks/SKILL.md` | 指导何时填写 Hard Rules |
| 技能指令 | `quick-tasks/SKILL.md` | 指导何时填写 Hard Rules |
| 测试任务 Prompt 模板 ×6 | `forge-cli/pkg/prompt/data/test-pipeline-*.md` | 内置 Hard Rules |
| Prompt 模板 ×4 | `forge-cli/pkg/prompt/data/*.md` | 定制化 Hard Rules clause |
| 索引规范化（Go） | `task-cli/pkg/task/` (新文件或 build.go) | normalizeTaskMD 自动清理空 section |

---

## 1. Agent 定义

**文件**: `plugins/forge/agents/task-executor.md`

在 `<EXTREMELY-IMPORTANT>` 块中新增规则 7：

```
7. HARD RULES OVERRIDE
   - Task files may contain ## Hard Rules with MUST/MUST NOT directives
   - These directives override agent judgment, ## Implementation Notes, and strategy defaults
   - Never substitute, modify, or skip a Hard Rules directive
```

**理由**：`<EXTREMELY-IMPORTANT>` 是 agent 遵从度最高的约束层。将 Hard Rules 认知放在这里，确保无论哪个 prompt 模板或任务类型，agent 都能识别并遵守 Hard Rules。

---

## 2. 任务模板（业务任务）

### 2.1 breakdown-tasks/templates/task.md

在 `## Implementation Notes` 之前新增：

```markdown
## Hard Rules
{{HARD_RULES}}

## Implementation Notes
{{NOTES}}
```

### 2.2 quick-tasks/templates/task.md

同上结构。

### 2.3 breakdown-tasks/templates/gate-task.md

将现有 Implementation Notes 中的 "This is a verification-only task. No new feature code should be written." 迁移为：

```markdown
## Hard Rules
- MUST NOT write new feature code — this is verification only
```

---

## 3. 技能指令

### 3.1 breakdown-tasks/SKILL.md

在生成任务的指导中添加 Hard Rules 填写规则：

**何时填写 Hard Rules**：
- 任务必须通过 justfile recipe 执行（而非直接调用底层工具）
- 命令带有环境变量、server lifecycle 等隐含依赖
- agent 自行构造命令会导致测试环境不完整
- 任务有明确的文件修改范围限制（MUST NOT touch 某些文件）

**何时不填写**：
- 普通代码任务（agent 默认 TDD + quality gate 行为已足够）
- agent 可以自行决定执行策略的灵活任务

### 3.2 quick-tasks/SKILL.md

同上。

---

## 4. 测试任务 Prompt 模板（内置 Hard Rules）

将 Hard Rules 直接内嵌到各 test-pipeline prompt 模板中，取代 Go lookup table 方案。理由：prompt template 是 agent 逐步遵循的执行策略，Hard Rules 在此层的遵从度高于 task file。

### 4.1 test-pipeline-run.md

当前问题模板——原始事故发生地。在 Workflow 之前新增 Hard Rules：

```markdown
## Hard Rules

<HARD-RULE>
- MUST invoke `Skill(skill="forge:run-e2e-tests")` to execute tests
- MUST NOT run `npx playwright test` or any direct test runner command
- The skill handles profile resolution, server lifecycle, result parsing, and reporting
</HARD-RULE>
```

### 4.2 test-pipeline-graduate.md

```markdown
## Hard Rules

<HARD-RULE>
- MUST invoke `Skill(skill="forge:graduate-tests")` for test graduation
- MUST NOT manually move, copy, or rewrite test files outside the skill
</HARD-RULE>
```

### 4.3 test-pipeline-verify-regression.md

已有内联"Do NOT start dev server"，形式化为 Hard Rules：

```markdown
## Hard Rules

<HARD-RULE>
- MUST use `just test-e2e` for regression verification
- MUST NOT start dev server manually — `just test-e2e` handles server lifecycle
- MUST NOT expand fixes beyond minimal scope (source code or test selectors only)
</HARD-RULE>
```

### 4.4 test-pipeline-gen-scripts.md

```markdown
## Hard Rules

<HARD-RULE>
- MUST invoke `Skill(skill="forge:gen-test-scripts")` to generate scripts
- MUST NOT write test scripts manually — the skill generates them from test cases
</HARD-RULE>
```

### 4.5 test-pipeline-gen-cases.md

```markdown
## Hard Rules

<HARD-RULE>
- MUST invoke `Skill(skill="forge:gen-test-cases")` to generate test cases
- MUST NOT write test cases manually — the skill generates them from PRD acceptance criteria
</HARD-RULE>
```

### 4.6 test-pipeline-eval-cases.md

已有 `<EXTREMELY-IMPORTANT>` 约束 main session。追加 Hard Rules 与其互补：

```markdown
## Hard Rules

<HARD-RULE>
- MUST invoke `Skill(skill="forge:eval-test-cases")` to evaluate test cases
- MUST NOT skip evaluation or rubber-stamp test cases without running the skill
</HARD-RULE>
```

---

## 5. Prompt 模板（高风险类型）

仅修改直接运行命令的 4 个高风险模板，添加定制化 Hard Rules clause：

### 5.1 implementation.md

在 Step 1 (Read Task Definition) 之后、Step 2 (TDD Implementation) 之前添加：

```markdown
<IMPORTANT>
If the task file contains ## Hard Rules with MUST/MUST NOT directives:
- Follow them exactly during the entire TDD cycle
- Hard Rules override your default approach for any step they address
- Do not rationalize bypassing a Hard Rule based on "I know a better way"
</IMPORTANT>
```

### 5.2 gate.md

在 Step 1 之后添加：

```markdown
<IMPORTANT>
If the task file contains ## Hard Rules with MUST/MUST NOT directives:
- Treat every MUST as a pass/fail criterion — no partial credit
- Treat every MUST NOT as a red line — violation means the gate fails
- Hard Rules override your judgment about what constitutes "good enough"
</IMPORTANT>
```

### 5.3 fix.md

在 Step 1 之后添加：

```markdown
<IMPORTANT>
If the task file contains ## Hard Rules with MUST/MUST NOT directives:
- Respect file scope restrictions (MUST NOT touch X) even if touching X seems like a cleaner fix
- Respect command restrictions (MUST use X) even if you think Y is equivalent
- Hard Rules define the fix boundary — do not expand beyond it
</IMPORTANT>
```

### 5.4 test-pipeline-verify-regression.md

在 Step 1 之后添加：

```markdown
<IMPORTANT>
If the task file contains ## Hard Rules with MUST/MUST NOT directives:
- Use exactly the specified commands for running regression tests
- Do not substitute alternative commands or tools
- Respect file modification restrictions when applying minimal fixes
</IMPORTANT>
```

---

## 6. 索引规范化

**文件**: `task-cli/pkg/task/` (新建 `normalize.go` 或追加到 `build.go`)

### 6.1 normalizeTaskMD 函数

```go
func normalizeTaskMD(content []byte) []byte {
    // 移除空的 ## Hard Rules section:
    // - "## Hard Rules" 后跟空行 + 下一个 ## heading
    // - "## Hard Rules" 后跟空行 + EOF
    // 使用正则或字符串操作匹配并移除
}
```

### 6.2 调用时机

在 `BuildIndex` 中，扫描完所有 .md 文件后，对每个文件调用 `normalizeTaskMD`，如果内容变化则写回：

```go
// After scanning all .md files
for _, entry := range entries {
    // ... existing read logic ...
    normalized := normalizeTaskMD(content)
    if !bytes.Equal(normalized, content) {
        os.WriteFile(filePath, normalized, 0644)
    }
}
```

---

## 不改动的部分

| 文件 | 原因 |
|------|------|
| Skill-delegating prompt 模板 ×6 | 薄包装层，实际工作由 skill 完成，agent 定义层约束已覆盖 |
| fix-record-missed.md | 低风险，已有 EXTREMELY-IMPORTANT 约束 |
| doc-generation-summary.md | 纯生成任务，无命令替代风险 |
| doc-generation-consolidate.md | 委托给 skill |
| test-pipeline-gen-cases/eval-cases/gen-scripts.md | 委托给 skill |

---

## 实现顺序

1. Agent 定义（最高影响，最小改动）
2. 测试任务 Prompt 模板 ×6（内置 Hard Rules）
3. 任务模板 + gate 模板迁移
4. 技能指令（SKILL.md）
5. Prompt 模板 ×4（高风险定制化 clause）
6. 索引规范化（normalizeTaskMD）
7. 测试覆盖

## Alternatives Analysis

### 证据对照

同一架构中，两种 prompt 模板风格已有实际运行数据：

| 模板 | Step 2 写法 | 事故记录 |
|------|-----------|---------|
| `test-pipeline-run.md` | `Skill(skill="forge:run-e2e-tests")` | **agent 跳过 skill，自行 `npx playwright test`** |
| `test-pipeline-verify-regression.md` | `just test-e2e`（直接命令） | 无已知事故记录 |

> **注意**：verify-regression 的"无已知事故"不等于"已验证有效"。该任务可能从未被实际执行，或运行时机不同（pipeline 末尾，风险场景不可比）。直接命令比 skill 调用少一层间接，**推测**更可靠，但缺少对照数据。

### 三个方案

#### P0: 当前提案（Hard Rules 加固 skill 调用）

改动 `test-pipeline-run.md` 仍然保持 `Skill(skill="forge:run-e2e-tests")`，通过 agent 定义 + task file Hard Rules + prompt template clause 三层加固定点执行。

```
Agent → prompt template → [Hard Rules 加固] → invoke skill → just test-e2e
```

**防御模型**：agent 仍然可以跳过 skill，但多层 Hard Rules 增加了跳过的"心理代价"。

#### P1: 直接命令（消除 skill 中间层）

将 `test-pipeline-run.md` 的 Step 2 从 "invoke skill" 改为直接命令，与 `verify-regression.md` 风格一致：

```markdown
### Step 2: Run E2E Tests

Run e2e tests via justfile:

```bash
just e2e-setup
just test-e2e --feature {{FEATURE_SLUG}}
```

If tests fail, diagnose per failure ratio (see skill's Failure Diagnosis section).
```

`run-e2e-tests` skill 本身不变——用户仍可手动 `/run-e2e-tests`。

```
Agent → prompt template → just test-e2e（无中间层）
```

**防御模型**：命令直接在 strategy 中，agent 需要主动违抗明确指令才能跑偏。

#### P2: 混合方案（直接命令 + skill 做报告）

Step 2 用直接命令执行，Step 3 可选调用 skill 做结果解析和报告生成：

```markdown
### Step 2: Run E2E Tests

```bash
just e2e-setup
just test-e2e --feature {{FEATURE_SLUG}}
```

### Step 3: Generate Report (optional)

Invoke the skill for result parsing and structured reporting:

```
Skill(skill="forge:run-e2e-tests")
```
```

```
Agent → prompt template → just test-e2e（必须）→ skill 报告（可选增强）
```

### 对比

| 维度 | P0: Hard Rules 加固 | P1: 直接命令 | P2: 混合 |
|------|-------------------|------------|---------|
| **防御原始问题** | 中 — 取决于模型对 Hard Rules 的遵从度，与原始问题同质（都是 prompt engineering 防线） | **高** — 命令直接在 strategy 中，无中间层可跳过。`verify-regression` 已验证有效 | 高 — 执行命令同 P1 |
| **skill 价值保留** | 完整 — skill 全流程执行 | 执行路径被替代，skill 退化为用户手动工具 | skill 保留为报告生成器（执行后增值） |
| **结果报告** | skill 自动生成结构化报告 | **丢失** — 除非在 prompt template 中内联报告步骤 | skill 提供报告（可选） |
| **profile awareness** | skill 运行时解析 profile | `just test-e2e` 已是 profile-aware；结果解析需 profile → **丢失** | 同 P1 执行 + skill 报告 |
| **改动文件数** | 10+ 文件，3 层 | 1 文件（prompt template） | 1-2 文件 |
| **后续维护** | 高 — 需维护 Hard Rules 内容在多个文件中一致 | 低 — prompt template 是唯一维护点 | 中 |
| **泛化能力** | **高** — Hard Rules 机制对所有任务类型生效 | 无 — 只修复了 test-pipeline.run | 无泛化 |
| **失败模式** | Hard Rules 被 agent 忽略 → 回到原始问题 | 结果报告功能丢失 | 报告生成可能被跳过 |

### 核心判断

**P0 的根本矛盾**：Hard Rules 对 agent 的约束力，与 prompt template 中 "invoke skill" 的约束力，来自同一个源头——模型对 prompt 的遵从度。如果 agent 能忽略 "invoke skill"，它同样能忽略 "MUST use just test-e2e"。多层加固增加的是冗余，不是结构性保障。

**P1 的核心优势**：消除了中间层。`verify-regression.md` 已用相同模式（直接命令）零事故运行——这是已证实的防御，不是推测。

**P1 的核心代价**：丢失 skill 的结果解析和报告生成（约 150 行领域逻辑）。但这部分可以：
- 内联到 prompt template（template 变长但自包含）
- 或按 P2 保留为可选 skill 调用

### 建议：P1 + P0 通用部分

1. **test-pipeline-run.md**：采用 P1（直接命令），消除 skill 中间层
2. **通用 Hard Rules 机制**：保留 P0 中除 test-pipeline-run 以外的所有改动
   - Agent 定义规则 7
   - 业务任务模板 `## Hard Rules` section
   - Gate 模板迁移
   - 技能指令
   - 4 个高风险 prompt template 的定制化 clause
   - 索引规范化

这样：
- test-pipeline.run 的修复是**结构性的**（直接命令，已验证有效）
- 通用 Hard Rules 机制为其他任务类型提供**增量保障**
- 两者不冲突，分别解决不同层面的问题

---

## 验证方式

- `go test ./...` 通过
- `task index --feature <slug>` 生成的测试任务 .md 正常工作
- 空的 Hard Rules section 被自动清理
- 现有业务任务 .md 不受影响（除非手动添加 Hard Rules）
- `task prompt T-test-3` 输出的策略包含 Hard Rules（不再依赖 task file）
- `test-pipeline-run.md` prompt template 内置 Hard Rules（原始事故场景）
