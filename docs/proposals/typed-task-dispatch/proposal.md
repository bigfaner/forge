---
created: 2026-05-11
author: fanhuifeng
status: Draft
---

# Proposal: Forge 任务执行策略化改造 — CLI 驱动的 Prompt 合成

## Problem

forge 的任务执行系统存在两个相互关联的结构性问题：

**问题一：task-executor 硬编码 TDD 流程，无法适配多种任务类型。**

当前系统中存在多种任务类型，但 task-executor 只有一套流程：

- **phase-summary 任务**：靠 `noTest: true` 强行跳过 TDD 步骤，流程错位
- **T-test-* 任务**（gen-test-cases、gen-test-scripts、run-e2e-tests 等）：通过 `mainSession: true` + 任务文件内嵌指令绕过 agent，本质上是在规避 task-executor 而非使用它
- `noTest` 字段是补丁式设计，需要在每个文档类任务上手动标注，而非由任务类型自动推导

**问题二：策略逻辑耦合在 markdown prompt 中，不可测试、不可版本化。**

task-executor 的执行策略写在 agent markdown 文件里，无法单元测试，每次修改需要人工验证，且随着任务类型增加，agent prompt 会持续膨胀。

### Evidence

- task-executor.md 已有 259 行，且仍只覆盖一种任务类型
- T-test-* 任务（7 个）全部通过 `mainSession: true` 绕过 task-executor，说明当前 agent 对这类任务完全不适用
- `noTest` 字段在 phase-summary、gate、T-test-* 等多个模板中重复出现，是补丁的直接证据

**运营事故（2026-04）**：调整 `test-pipeline.eval-cases` 的执行逻辑时，开发者需同时修改任务模板（调整内嵌指令）、run-tasks 路由（确认 mainSession 分支正确）、task-executor.md（确认不会误拦截），三处改动无法通过单元测试验证，只能人工运行完整任务链路——单次验证耗时约 2 小时，且改动后仍出现一次 record 文件未写入的静默失败，排查额外耗时 1 小时。根本原因：策略逻辑分散在 markdown 中，字段组合（`mainSession: true` + `noTest: true`）的实际行为无法在不运行 agent 的情况下推断。

**新增任务类型的固定成本**：每新增一种任务类型，需触及至少 3 处文件（task-executor.md 添加分支、index.schema.json 添加枚举值、任务模板添加字段），且无自动化回归手段，每次修改后需人工运行完整任务链路验证，平均耗时 ~2 小时/次。

### Urgency

随着 forge 功能扩展，任务类型只会增加。当前的补丁式绕过已经是技术债：
- 新增任务类型时需要同时修改 agent、schema、模板三处
- 调试困难：agent 行为取决于多个隐式字段的组合，而非显式类型
- 策略逻辑在 markdown 中，无法通过 Go 单元测试验证正确性

## Proposed Solution

**将策略逻辑从 agent markdown 下沉到 CLI，通过 `task prompt <id>` 命令合成类型专属 agent prompt，task-executor 变为薄执行器。**

### 架构变化

**当前架构**：
```
run-tasks → task claim → dispatch to task-executor (fat agent, all strategies in markdown)
                      ↘ error-fixer (record 缺失时的异常恢复)
```

**目标架构**：
```
run-tasks → task claim
  ├─ mainSession=true → Step 1.5（主会话执行，不经过 task prompt）
  ├─ 正常任务 → task prompt <id> → Agent(task-executor, prompt=<synthesized>)
  └─ record 缺失 → task prompt <id> --fix-record-missed → Agent(task-executor, prompt=<synthesized>)
```

### 两层职责分离

| 层 | 位置 | 内容 |
|----|------|------|
| **执行约束层** | task-executor.md（agent 定义，始终生效） | 硬规则：ONE TASK PER INVOCATION、record-task 强制、无后台任务、最多 3 次 subagent 调用、STOP 规则 |
| **任务策略层** | CLI 合成 prompt（作为 Agent() 的 prompt 参数传入） | 步骤定义：TDD 流程 / 文档生成流程 / fix 诊断流程等，随 type 变化 |

两层不合并：约束层内嵌在 agent 定义里，策略层通过 prompt 参数注入。约束层始终生效，不依赖策略层的正确性。

### 核心变更

**1. 新增 `task prompt` 子命令**

```bash
task prompt <id>                        # 根据任务 type 合成类型专属 agent prompt（stdout 纯文本）
task prompt <id> --fix-record-missed    # 合成 record 缺失恢复专用 prompt（内联 dispatch，不入 index）
```

CLI 自动从 `.forge/state.json` 和 `index.json` 读取 feature、scope 等上下文。CLI 自身实现 phase boundary detection（与当前 run-tasks 逻辑等价，提取为硬性约束）：扫描 index.json 中所有已完成任务，取最大 phase 编号；若当前任务 phase > 该值，则为新 phase 第一个任务，自动定位并注入前一 phase 的 summary 文件路径（`docs/features/<slug>/tasks/records/<prev-phase>-summary.md`）。模板合成采用标准字符串替换（`{{TASK_ID}}`、`{{SCOPE}}` 等占位符），无额外依赖。失败时输出到 stderr，退出码非零，run-tasks 检测到后标记任务 blocked。

**2. 新增 `task migrate` 命令**

根据现有字段自动推断并填充 type 字段，保留所有现有任务状态。**迁移前提：所有 in_progress 任务必须先完成或手动标记为 completed/blocked，task migrate 不处理在途任务。**

| 推断规则 | 推断结果 |
|---------|---------|
| ID 以 `.summary` 结尾 | `doc-generation.summary` |
| ID 以 `.gate` 结尾 | `gate` |
| ID 为 `T-test-1` | `test-pipeline.gen-cases` |
| ID 为 `T-test-1b` | `test-pipeline.eval-cases` |
| ID 为 `T-test-2` | `test-pipeline.gen-scripts` |
| ID 为 `T-test-3` | `test-pipeline.run` |
| ID 为 `T-test-4` | `test-pipeline.graduate` |
| ID 为 `T-test-4.5` | `test-pipeline.verify-regression` |
| ID 为 `T-test-5` | `doc-generation.consolidate` |
| ID 以 `fix-` 或 `disc-` 开头 | `fix` |
| 其余 | `implementation` |

**3. 新增 `type` 字段，废弃 `noTest`**

完整类型体系（11 种业务类型 + 1 种内部操作模式）：

| type | 子类型 | 描述 | prompt 模板 | 路由 |
|------|--------|------|------------|------|
| `implementation` | — | 编码任务，TDD 流程 + quality gate | Go embed markdown | task prompt |
| `doc-generation` | `summary` | Phase summary 生成 | 独立模板 | task prompt |
| `doc-generation` | `consolidate` | 提取业务规则和技术规范 | 独立模板 | task prompt |
| `test-pipeline` | `gen-cases` | 生成测试用例文档 | 独立模板 | task prompt |
| `test-pipeline` | `eval-cases` | 对抗评估测试用例（需 spawn doc-scorer/doc-reviser subagent） | 独立模板 | task prompt（**永久例外**：run-tasks 在主会话中直接按 prompt 执行，不 dispatch subagent。原因：平台硬性限制——subagent 无法再开启 subagent） |
| `test-pipeline` | `gen-scripts` | 生成 e2e 测试脚本 | 独立模板 | task prompt |
| `test-pipeline` | `run` | 运行 e2e 测试 | 独立模板 | task prompt |
| `test-pipeline` | `graduate` | 迁移测试脚本到回归套件 | 独立模板 | task prompt |
| `test-pipeline` | `verify-regression` | 验证回归套件 | 独立模板 | task prompt |
| `fix` | — | 修复任务（编译错误、测试失败、lint 问题等），诊断 → 定位 → 修复 → 验证 → 提交 | 独立模板（承接 error-fixer 全部能力） | task prompt |
| `gate` | — | 阶段出口验证，直接进入 quality gate。自修复还是追加 fix 任务的判断标准由 prompt 模板自描述 | 独立模板 | task prompt |
| `fix-record-missed` | — | record 缺失恢复（内部操作模式，不入 index，不在 type 枚举中） | 独立模板 | `--fix-record-missed` flag |

`noTest` 字段废弃，由 type 自动推导。`mainSession` 字段保留，但仅用于 `test-pipeline.eval-cases`（需在主会话中执行以 spawn subagent）。其他任务不应设置 mainSession。

**4. 路由优先级**

```
run-tasks claim 任务后：
  if type == test-pipeline.eval-cases:
    → task prompt <id> 得到 prompt
    → 在主会话中直接按 prompt 执行（不 dispatch subagent）
  else:
    → task prompt <id> → Agent(forge:task-executor, prompt=<synthesized>)

record 缺失时：
  → task prompt <id> --fix-record-missed → Agent(forge:task-executor, prompt=<synthesized>)（内联，不入 index）
```

execute-task 命令与 run-tasks 保持相同路由逻辑：eval-cases 在主会话执行，其余 dispatch 给 task-executor。

**5. task-executor 变为薄执行器**

task-executor.md 只保留执行约束层（硬规则），移除所有策略块（TDD 步骤、quality gate 步骤等）。CLI 合成的 prompt 作为 Agent() 的 prompt 参数传入，提供任务策略。

### Observable Behavior Example

**`task prompt T-impl-1` 的实际输出（implementation 类型）**

```
You are executing task T-impl-1 in feature auth-refresh.

Scope: src/auth/refresh.go, src/auth/refresh_test.go

Steps:
1. Read task file at docs/features/auth-refresh/tasks/T-impl-1.md
2. Write failing tests first (TDD)
3. Implement until tests pass: go test ./src/auth/...
4. Run quality gate: golangci-lint run ./src/auth/...
5. Commit with message: feat(auth): implement token refresh logic
6. Call forge:record-task to write execution record

Constraints (from task-executor):
- ONE TASK ONLY. Do not claim or start any other task.
- record-task is mandatory before stopping.
- No background processes.
```

对比当前行为：run-tasks 直接 dispatch `Agent(forge:task-executor)`，task-executor.md 内嵌 TDD 步骤和 quality gate 步骤，所有类型共用同一套流程。

**`task prompt T-fix-1` 的实际输出（fix 类型）**

```
You are executing task T-fix-1 in feature auth-refresh.

Scope: src/auth/refresh.go, src/auth/refresh_test.go

Steps:
1. Read task file at docs/features/auth-refresh/tasks/T-fix-1.md
2. Diagnose: run `go build ./src/auth/...` and `go test ./src/auth/...` to reproduce the failure
3. Locate: identify the root cause from compiler output or test failure message
4. Fix: apply the minimal change that resolves the failure
5. Verify: re-run the failing command and confirm exit code is 0
6. Commit with message: fix(auth): <describe what was broken>
7. Call forge:record-task to write execution record

Constraints (from task-executor):
- ONE TASK ONLY. Do not claim or start any other task.
- record-task is mandatory before stopping.
- No background processes.
```

两层分离在 fix 类型上的体现：策略层（task prompt）注入"诊断 → 定位 → 修复 → 验证"的五步流程，执行约束层（task-executor.md）仍强制 ONE TASK 和 record-task 规则，两层互不感知。当前 error-fixer 的修复逻辑内嵌在 agent markdown 中，无法在不运行 agent 的情况下检查；迁移后可通过 `task prompt T-fix-1` 直接查看合成结果，调试时无需启动 agent。

**task-executor.md 改造前后对比**

```diff
- ## Execution Flow
-
- ### TDD Steps (implementation tasks)
- 1. Read task file
- 2. Write failing tests
- 3. Implement until tests pass
- 4. Run golangci-lint
- 5. Commit
-
- ### Doc Generation Steps (noTest: true)
- 1. Read task file
- 2. Generate document
- 3. Commit
-
- ### Quality Gate Steps
- 1. Run full test suite
- 2. Run lint
- 3. Report result
-
  ## Hard Constraints
  - ONE TASK PER INVOCATION
  - record-task is mandatory before stopping
  - No background tasks
  - Maximum 3 subagent calls
  - STOP after record-task
```

改造后 task-executor.md 仅保留 Hard Constraints 块，预计从 259 行缩减至 ~40 行。执行步骤由 `task prompt <id>` 按类型注入，不再内嵌于 agent 定义。

**6. error-fixer 废弃**

error-fixer 的全部能力（编译错误修复、测试失败修复、lint 修复、record 缺失恢复）由以下两个路径承接：
- 编译/测试/lint 修复 → `type: fix` 的 prompt 模板
- record 缺失恢复 → `task prompt <id> --fix-record-missed`

**7. task validate 扩展**

现有 `task validate` 命令扩展支持 type 字段验证：检查所有任务的 type 字段是否为合法枚举值，以及 type 与其他字段（如 mainSession）的一致性。

**8. execute-task 同步更新**

execute-task 命令与 run-tasks 保持一致：调用 `task prompt <id>` 合成 prompt，mainSession 路由逻辑相同。

### 全链路改造范围

实施顺序（后续步骤依赖前序步骤完成）。假设团队规模：1 名开发者。

**阶段一：CLI 基础能力（L）**
1. **task-cli 新增 `task prompt` 命令**：`task prompt <id>` 和 `task prompt <id> --fix-record-missed`，Go embed markdown 模板，11 种类型各有独立模板文件，CLI 自身实现 phase boundary detection，标准字符串替换（`{{TASK_ID}}`、`{{SCOPE}}` 等占位符）
2. **task-cli 新增 `task migrate` 命令**：根据现有字段自动推断 type，保留任务状态（迁移前提：所有 in_progress 任务已完成）
3. **task-cli 扩展 `task validate` 命令**：支持 type 字段验证

阶段一完成标准：`task prompt <id>` 对 11 种 type 各输出正确 prompt，Go 单元测试全部通过；`task migrate` 和 `task validate` 命令可正常运行，`task validate` 对迁移后的 index.json 无报错。

**阶段二：Schema 与模板（M）**
4. **`index.schema.json`**：新增 `type` 枚举（必填，不兼容旧 index.json），废弃 `noTest`
5. **所有任务模板**（`task.md`、`phase-summary-task.md`、`gate-task.md`、`gen-test-cases.md` 等）：frontmatter 加 `type` 字段，移除 `noTest`
6. **`breakdown-tasks` skill**：生成任务时自动设置 type
7. **`quick-tasks` skill**：同步更新

阶段二完成标准：所有任务模板 frontmatter 包含 type 字段，`noTest` 字段已移除；breakdown-tasks 和 quick-tasks 生成的任务包含正确 type，`task validate` 对生成结果无报错。

**阶段三：Agent 与命令更新（S）**
8. **`task-executor.md`**：精简为薄执行器，只保留执行约束层，移除策略块
9. **`execute-task.md`**：同步更新，调用 task prompt 合成 prompt
10. **`run-tasks` 命令**：调用 `task prompt <id>` 合成 prompt，record 缺失时调用 `task prompt <id> --fix-record-missed`，不再 dispatch error-fixer

阶段三完成标准：所有 11 种类型的任务实际运行行为等价（Success Criteria 前 8 条全部通过）。

**阶段四：清理（S）**
11. **`error-fixer` agent**：废弃（阶段三验证通过后）
12. **不兼容旧任务**：无 type 字段的旧 index.json 通过 `task migrate` 迁移（type 字段必填，不提供兼容期）

阶段四完成标准：error-fixer 在 run-tasks、execute-task、guide.md 等文件中无孤立引用；`noTest` 字段从 schema 和所有模板中完全移除。

## Alternatives Considered

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| Do nothing | 零改动成本 | 技术债持续累积，新任务类型无法干净接入 | Rejected：不可持续 |
| Agent 内部策略分发（type 字段 + agent 内部 if/else） | 只改 markdown，快速落地 | 策略逻辑不可测试；agent prompt 持续膨胀 | Rejected：短期方案 |
| 多 agent（每种类型一个 agent） | 每个 agent 职责单一 | run-tasks 需要维护 dispatch 表；共享逻辑重复 | Rejected：维护成本高 |
| **CLI 合成 prompt + 薄执行器** | 策略逻辑可测试；agent prompt 精简；两层职责分离 | 实现成本高（需改 CLI + run-tasks + 所有模板）；调试面扩大——任务失败时需逐层排查 CLI 合成、模板内容、执行器约束三处，而非当前的单一 agent；11 个 prompt 模板从零编写，初始质量无保障（与 Risk 表中"高/高"评级一致） | **Selected**：调试面扩大可通过 `task prompt <id>` 独立输出 prompt 缓解——可在不运行 agent 的情况下检查合成结果，将调试范围缩小到具体层；模板冷启动风险通过逐类型对比验证控制（参照现有 agent markdown 逐条提取），不影响架构选择的正确性 |

## Scope

### In Scope
- task-cli 新增 `task prompt` 命令（11 种类型，Go embed 模板，含 phase boundary detection）
- task-cli 新增 `task migrate` 命令（自动推断 type，保留任务状态）
- task-cli 扩展 `task validate` 命令（支持 type 字段验证）
- `index.schema.json` 新增 `type` 字段（必填枚举，不兼容旧 index.json），废弃 `noTest`
- task-executor agent 精简为薄执行器（保留执行约束层）
- execute-task 命令同步更新（调用 task prompt）
- 所有任务模板 frontmatter 加 `type` 字段，移除 `noTest`
- breakdown-tasks 和 quick-tasks 生成任务时自动设置 type
- run-tasks 调用 `task prompt <id>` 合成 prompt，统一 dispatch
- error-fixer agent 废弃，能力由 fix 模板和 --fix-record-missed 承接
- 旧 index.json 通过 `task migrate` 迁移（不提供兼容期）

### Out of Scope
- CLI 二进制重命名（task→forge）——独立提案，后续实施
- 技术设计时选定 e2e 测试框架——独立提案
- 新增任务类型（本次只改造现有类型）
- 任务执行的并发/并行调度
- 跨 feature 的任务依赖

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| task prompt 命令失败（type 未知、模板错误）导致 dispatch 链路中断 | 中（迁移期）/ 低（稳定后） | 高 | 迁移期：11 个模板从零编写，模板错误是主要失败机制，与下方模板质量风险同级（中）；稳定后：模板经过逐类型验证，失败概率降低（低）。缓解：stderr + 非零退出码；run-tasks 检测到后标记任务 blocked，不静默失败 |
| 11 个 prompt 模板从零编写，初始质量无保障 | 高 | 高 | 参照现有 agent markdown（task-executor.md、error-fixer.md、T-test-* 任务文件）逐条提取逻辑；每种类型迁移后立即运行对应任务验证行为等价；模板质量稳定后，上方 task prompt 失败概率随之降低 |
| phase boundary detection 迁移到 CLI 后引入 regression | 中 | 高 | 迁移后与原 run-tasks 逻辑对比测试；phase-summary 注入有无的行为差异通过集成测试覆盖 |
| test-pipeline 子类型 prompt 模板质量不足，导致执行结果与当前不一致 | 中 | 高 | 每个子类型独立模板文件，迁移时逐个子类型与当前行为对比验证 |
| fix 模板未能完整承接 error-fixer 的全部修复能力 | 中 | 高 | 将 error-fixer.md 的完整逻辑迁移到 fix 模板；通过对比测试验证等价性 |
| task migrate 推断规则不完整，导致部分任务 type 错误 | 中 | 中 | 推断规则覆盖所有已知 ID 模式；migrate 输出推断结果，运行 task validate 验证 |
| 旧 index.json 无 type 字段，task migrate 前所有 task 命令失败 | 高 | 高 | 明确不兼容，文档说明需先运行 task migrate；task migrate 是迁移的第一步；在途任务需先完成或手动标记后再迁移 |

## Success Criteria

- [ ] `task prompt <id>` 对 11 种 type 各输出正确的类型专属 prompt，Go 单元测试覆盖所有 type
- [ ] `task prompt <id> --fix-record-missed` 合成的 prompt 能正确补写缺失 record
- [ ] task-executor 执行 `doc-generation.summary` 任务时，不出现任何 TDD 相关步骤
- [ ] task-executor 执行 `test-pipeline.gen-cases` 任务时，在 `docs/features/<slug>/test-cases.md` 生成测试用例文件，文件包含 ≥1 个测试用例条目，且 `task validate` 对该任务通过（与当前 T-test-1 执行结果等价的可验证定义）
- [ ] task-executor 执行 `fix` 任务时，针对编译错误场景：`go build ./...` 在修复后退出码为 0；针对测试失败场景：`go test ./...` 在修复后退出码为 0；针对 lint 场景：`golangci-lint run ./...` 在修复后退出码为 0（三个场景各需一次实际运行验证）
- [ ] 每种 type 的 prompt 模板迁移后，实际运行对应任务，等价定义为：(a) 任务执行后 record 文件写入 `docs/features/<slug>/tasks/records/<task-id>.md`；(b) `task validate` 对该任务通过；(c) 任务产物（代码文件、文档文件、测试文件）与迁移前同类任务产物的文件路径模式一致
- [ ] `task migrate` 对所有已知 ID 模式推断出正确 type，`task validate` 通过
- [ ] breakdown-tasks 生成的所有任务均包含正确的 type 字段，`task validate` 通过
- [ ] `test-pipeline.eval-cases` 任务：`task prompt <id>` 输出正确 prompt，run-tasks 在主会话中直接执行（不 dispatch subagent）
- [ ] error-fixer agent 废弃后无孤立引用（run-tasks、execute-task、guide.md 等均不再引用 error-fixer）
- [ ] `noTest` 字段从 schema 中移除，现有使用 noTest 的模板全部迁移到 type 字段
- [ ] task-cli 单元测试覆盖率 ≥ 80%（含 task prompt、task migrate、task validate 扩展）

## Next Steps

- Proceed to `/write-prd` to formalize requirements
