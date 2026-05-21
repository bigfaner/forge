---
created: 2026-05-21
author: faner
status: Draft
---

# Proposal: Decouple run-tests from e2e/just Hardcoding

## Problem

`run-e2e-tests` skill 硬编码了 `just e2e-setup` / `just e2e-verify` / `just e2e-test` 三个命令，将 Forge 绑定到 `just` 任务运行器和 e2e 测试语义上。非 e2e 项目（Go unit tests、Vitest integration tests 等）无法使用此 skill。

### Evidence

- SKILL.md 中 3 处硬编码 `just e2e-*` 命令（第 87、99、109 行）
- 前置检查 `grep -q "^e2e-setup:" Justfile` 强制要求 justfile
- skill 名称 `run-e2e-tests` 在整个 pipeline 中传播 e2e 语义（gen-test-scripts 的 Related Skills 表、task chain 中的 T-test-3）
- `gen-test-scripts/templates/` 目录包含 27MB 死代码（validate-specs.mjs + node_modules），SKILL.md 零引用

### Urgency

v3.0.0 正在重构测试体系（test profile 系统、Convention-driven 模型）。保留 e2e/just 硬编码会与框架无关化方向矛盾，越晚改迁移成本越高。

## Proposed Solution

将 `run-e2e-tests` 重新定位为纯执行器 `run-tests`，执行命令从 `.forge/config.yaml` 的 `test.execution` 节点读取。skill 只负责三件事：执行配置中的命令 → 解析结果 → 生成报告。

### Innovation Highlights

这不是创新，是职责归位。Convention 管"怎么写测试"（framework、assertion、result format），Config 管"怎么跑测试"（执行命令模板）。两者已在各自的文件中，只是缺少 `test.execution` 这个桥接点。

## Convention vs Config 职责边界

当前 Convention Result Format 段落包含三个字段：`Output flags`、`Format type`、`Execution command`。其中 `Execution command` 本质上是"怎么跑"的信息，与 Convention "怎么写"的定位矛盾。

**修订方案**：Convention Result Format 废弃 `Execution command` 字段，只保留解析元数据（`Output flags`、`Format type`）。`test.execution.run` 成为唯一的执行命令来源。

**对 Convention 文件的影响**：
- `docs/conventions/testing-conventions.md`：Result Format 段落移除 `Execution command` 行
- `docs/conventions/testing-go.md`：移除 `- **Execution command**: ...` 行
- `docs/conventions/testing-vitest.md`：同上
- `docs/conventions/testing-ginkgo.md`：同上

**对 init-justfile 的影响**：
- Step 3a 生成 e2e-test recipe 时，改为从 `.forge/config.yaml` 的 `test.execution.run` 读取执行命令（而非 Convention）
- Convention 仍提供 `Output flags` 和 `Format type`，用于结果解析

**数据流变更**：
```
Before: Convention Result Format → init-justfile 生成 recipe → run-e2e-tests 调用 recipe
After:  Config test.execution   → run-tests 直接执行
        Config test.execution   → init-justfile 生成 recipe（供手动使用 / CI 调用）
        Convention Result Format → run-tests 解析结果（仅 format-type + output-flags）
```

## Config Schema

### `test.execution` 完整定义

```yaml
# .forge/config.yaml
test:
  execution:
    # 必填：执行测试的命令模板
    run: "just e2e-test --feature {slug}"
    # 可选：前置设置命令
    setup: "just e2e-setup"
    # 可选：前置校验命令
    pre-check: "just e2e-verify --feature {slug}"
    # 可选：后置清理命令（在任何退出情况下执行）
    teardown: "just e2e-teardown"
    # 可选：结果目录（默认从 Convention 或 journey 路径推导）
    results-dir: "tests/{journey}/results"
    # 可选：整体执行超时（秒），默认 600
    timeout: 300
```

### 模板变量

| 变量 | 来源 | 示例值 | 缺失时行为 |
|------|------|--------|-----------|
| `{slug}` | `forge feature` | `user-auth` | 报错（无法确定测试范围） |
| `{journey}` | Convention 或目录扫描 | `e2e` | 使用默认值 `e2e` |
| `{test-dir}` | Convention Framework | `tests/e2e` | 使用默认值 `tests` |
| `{results-dir}` | `test.execution.results-dir` | `tests/e2e/results` | 使用默认值 `tests/{journey}/results` |

**转义规则**：`{{slug}}` 转义为字面量 `{slug}`。

### Output Flags 一致性验证

`run-tests` 在执行前增加验证步骤：

1. 读取 Convention Result Format 的 `format-type` 和 `Output flags`
2. 检查 `test.execution.run` 命令中是否包含所需的 output flags
3. 如果不匹配，报错：
   > Convention declares format-type `json-stream` which requires output flags `-json`, but `test.execution.run` does not include these flags. Either add the flags to your run command, or change Convention's format-type to `text-verbose`.

此验证在 Step 0（加载 Convention）和 Step 1（读取 config）之后、执行之前进行。

## Requirements Analysis

### Key Scenarios

1. **Happy path**: 用户在 config.yaml 配置 `test.execution.run`，调用 `/run-tests`，skill 执行命令、解析结果、生成报告
2. **完整生命周期**: 用户配置 setup → pre-check → run → teardown，skill 按顺序执行
3. **无配置**: config.yaml 中没有 test.execution → 报错并提示配置方法
4. **pre-check 失败**: 用户配置了 pre-check 步骤，返回非零退出码 → 报错并提示回到上游 skill（如 gen-test-scripts）
5. **结果解析失败**: Convention Result Format 与实际输出不匹配 → 使用 text-verbose fallback
6. **output-flags 不匹配**: Convention 声明 `json-stream` 但 run 命令缺少 `-json` → 执行前报错
7. **无 feature slug**: `{slug}` 变量无法解析 → 报错提示运行 `forge feature <slug>`
8. **teardown 异常中断**: 用户会话中断，teardown 未执行 → 下次启动时检测遗留状态文件并执行清理
9. **超时**: 执行超过 `timeout` 秒 → 终止进程，标记所有测试为 FAIL(timeout)

### Non-Functional Requirements

- 执行命令支持模板变量替换（`{slug}`、`{journey}`、`{test-dir}`、`{results-dir}`）
- Convention Result Format 解析逻辑不受影响（已框架无关）
- 报告模板适用于任何测试类型（非 e2e 专属）
- teardown 保证执行：使用状态文件（`.forge/test-state.json`）记录待清理命令，skill 启动时检测并清理遗留状态

### Constraints & Dependencies

- 依赖 `.forge/config.yaml` 存在且包含 `test.execution.run` 配置
- 依赖 Convention 文件提供 Result Format 的 `format-type` 和 `Output flags`（用于结果解析和一致性验证）
- 不依赖 justfile（但用户可以在 config 中使用 just 命令）
- `test.execution.run` 是唯一的执行命令来源（Convention 的 `Execution command` 字段废弃）

## Alternatives & Industry Benchmarking

### Industry Solutions

测试运行器的配置通常在项目配置文件中（package.json scripts、Justfile、Makefile）。本方案将这种模式提升为 Forge config 的标准字段，使 skill 不关心具体运行器。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零改动成本 | 绑定 just+e2e，非 TS/Playwright 项目无法使用 | Rejected: 与 v3.0 框架无关化方向矛盾 |
| Convention Execution 段落 | 自研 | Convention 已加载，无新文件 | 混淆"怎么写"和"怎么跑"的职责边界 | Rejected: 职责不清 |
| **Config test.execution** | 自研 | 职责清晰（Convention=写/解析，Config=跑），支持任意执行器，集中管理，单一真实源 | 新增 config schema，需废弃 Convention Execution command 字段 | **Selected: 最干净的职责分离，解决 two-sources-of-truth** |
| CLI 参数传入 | CLI pattern | 零配置 | 无法复用，每次调用都要输入完整命令 | Rejected: 操作成本高 |

## Feasibility Assessment

### Technical Feasibility

完全可行。核心变更：(1) 读取 config.yaml 的 test.execution 字段、(2) 替换模板变量、(3) bash 执行。现有 Result Format 解析（仅 `format-type` + `Output flags`）和报告生成逻辑不变。

### Resource & Timeline

改动量中等：1 个 skill 重写 + 1 个 config schema + Convention 文件修订（移除 Execution command 字段） + init-justfile Step 3a 数据源切换 + 若干引用更新。2-3 个 coding task 可完成。

### Dependency Readiness

`.forge/config.yaml` 已存在，forge CLI 已有 `forge config get` 命令。Convention 文件结构修改是纯删除操作（移除 `Execution command` 行），向后兼容。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "skill 需要管理测试生命周期（setup/verify/teardown）" | Occam's Razor | Refined: 生命周期由用户在 config 中定义，skill 只是按顺序执行。不内置任何生命周期逻辑 |
| "validate-specs.mjs 是 gen-test-scripts 的必要组件" | Evidence check (grep) | Overturned: SKILL.md 零引用，470 行 Playwright 专属代码 + 27MB node_modules 是死代码 |
| "justfile 是 Forge 的标准执行层" | Assumption Flip | Refined: justfile 是 Forge 推荐的执行层，但不是唯一选择。Config 中可以写 make、npm run、cargo test 等任何命令 |
| "Convention Result Format 包含 Execution command 是合理的" | Assumption Flip | Overturned: `Execution command` 是"怎么跑"的信息，不应出现在"怎么写"的 Convention 中。废弃此字段，Config 成为唯一执行命令来源 |
| "verify 步骤是测试执行的必要环节" | Need Gate | Refined: verify 语义绑定到 gen-test-scripts 的 `// VERIFY:` 标记，不是通用概念。泛化为 `pre-check`，语义由用户定义 |

## Scope

### In Scope

- 重命名 `run-e2e-tests` → `run-tests`，重构 SKILL.md 为纯执行器
- 定义 `.forge/config.yaml` 的 `test.execution` schema（run 必填，setup/pre-check/teardown/results-dir/timeout 可选）
- 定义模板变量体系（`{slug}`、`{journey}`、`{test-dir}`、`{results-dir}`）和转义规则
- 增加 output-flags 一致性验证（执行前检查 Convention format-type 与 run 命令的匹配）
- 废弃 Convention Result Format 的 `Execution command` 字段（从 4 个 Convention 文件和 testing-conventions.md 中移除）
- 修改 `init-justfile` Step 3a 数据源：从 Config `test.execution.run` 读取执行命令（而非 Convention Execution command）
- 删除 `gen-test-scripts/templates/` 整个目录（validate-specs.mjs + test + fixtures + 27MB node_modules）。已确认 init-justfile、gen-test-scripts、run-e2e-tests 三个 skill 均不引用此目录
- 泛化报告模板（`e2e-report.md` → `test-report.md`）：去除 e2e 专属文案、Screenshots 段落条件渲染、测试类型分类动态生成
- 更新所有 `run-e2e-tests` 引用：
  - `plugins/forge/commands/run-tasks.md` 第 97 行
  - `plugins/forge/skills/gen-contracts/SKILL.md` 第 176 行
  - `plugins/forge/skills/gen-test-cases/SKILL.md` 第 141 行
  - `plugins/forge/skills/gen-test-scripts/SKILL.md` 第 372 行

### Out of Scope

- `gen-test-scripts` skill 的核心逻辑修改（只更新 Related Skills 引用）
- `graduate-tests` skill 的重命名（如果它存在 `/run-e2e-tests` 引用则一并更新）
- 其他 skill 中 e2e 命名引用的全面清理（如 hooks/guide.md 中的命令列表）
- 其他 Convention 文件的非 Result Format 部分修改

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 已有项目缺少 test.execution 配置导致 skill 无法使用 | M | M | 报错信息明确：提示用户在 config.yaml 中添加 test.execution 配置，并给出示例 |
| 用户配置的 run 命令缺少 output flags 导致解析失败 | M | H | 执行前验证 Convention format-type 与 run 命令的 output flags 一致性，不匹配时报错 |
| Convention 废弃 Execution command 后 init-justfile 生成 recipe 的数据源断裂 | L | H | init-justfile Step 3a 切换到读取 Config test.execution.run 作为 recipe 源 |
| teardown 在会话中断时未执行 | M | M | 使用状态文件 `.forge/test-state.json` 记录待执行 teardown 命令，下次启动时检测并清理 |
| 重命名后用户缓存的旧命令名失效 | L | L | 不保留 alias，在 upgrade 文档中明确列出 breaking change |

## Success Criteria

- [ ] `run-tests` skill 不包含任何硬编码的 `just` 或 `e2e` 命令
- [ ] `test.execution.run` 配置缺失时，skill 报错并给出配置示例
- [ ] 所有模板变量（`{slug}`、`{journey}`、`{test-dir}`、`{results-dir}`）在命令模板中正确替换
- [ ] Convention format-type 与 run 命令 output flags 不匹配时，执行前报错
- [ ] Convention 文件中不再包含 `Execution command` 字段
- [ ] init-justfile 从 Config 读取执行命令生成 recipe（而非 Convention）
- [ ] `gen-test-scripts/templates/` 目录已被删除（减少 27MB）
- [ ] 报告模板中不含 "E2E" 字样，Screenshots 段落条件渲染
- [ ] 所有 5 处 `run-e2e-tests` 引用已更新为 `run-tests`
