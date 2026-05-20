---
title: "Journey-Contract Test Model"
domains: [testing, journey, contract, e2e, migration]
---

# Journey-Contract Test Model

Forge 的 e2e 测试管道以用户工作流（Journey）为组织单元，通过 Contract 六维声明定义预期行为，由 Tag-Based Promotion 管理测试生命周期。

## Core Concepts

### Journey

Journey 描述用户完成某个目标的真实工作流，是测试的主要组织单元——一个 Journey 对应一个连贯的用户工作流。

| Property | Description |
|----------|-------------|
| Name | kebab-case 标识符（如 `task-lifecycle`） |
| Risk | `High`（有状态变更/数据丢失风险）、`Medium`（多步交互无不可逆副作用）、`Low`（只读操作） |
| Steps | 有序的用户动作序列，每个带预期结果 |
| Invariants | 跨步骤不变量，必须在整个 Journey 中成立 |

每个 Journey 在独立的临时工作目录中执行，防止并行执行时的交叉干扰。

### Step

Step 是 Journey 中的单个用户动作。每个 Step 映射到一个 Contract，包含一个或多个 Outcome。

| Property | Description |
|----------|-------------|
| Sequence number | Journey 内的 1-based 序号 |
| User action | 用户执行的操作（运行命令、点击按钮、发送请求等） |
| Expected outcomes | 一个或多个 Outcome 声明，各有独立的 Preconditions |

### Contract

Contract 是 Step 的验证机制，通过六维声明定义系统预期行为。所有维度在 Outcome 层声明，Invariants 额外在 Journey 层声明。

### Outcome

Outcome 是特定场景（成功、错误变体、边界情况）下完整的 Contract 维度声明集合。同一 Step 内的 Outcome 通过 Preconditions 区分，且必须互斥——同一系统状态下最多只有一个 Outcome 的 Preconditions 可被满足。

| Property | Required | Description |
|----------|----------|-------------|
| Name | Yes | 描述性标签（如 `success`、`not-in-progress`） |
| Preconditions | Yes | 使此 Outcome 成为活跃状态所需的系统状态 |
| Input | Yes | 用户提供给系统的输入 |
| Output | Yes | 系统产生的输出（语义描述符） |
| State | Yes | 系统状态变更 |
| Side-effect | No | 外部副作用（默认 `none`） |
| Invariants | No | 步骤级不变量（默认无约束） |

## Six-Dimension Declaration Rules

每个 Outcome 声明六个维度。四个必选（非空），两个可选。

### Mandatory Dimensions

| Dimension | Content | Source |
|-----------|---------|--------|
| Preconditions | 执行前必须成立的系统状态 | Journey 前置条件 + Fact Table 状态信息 |
| Input | 进入系统的输入 | Journey 用户动作 + Fact Table 命令/参数/端点信息 |
| Output | 系统产生的输出 | Journey 预期结果 + Fact Table 输出模式 |
| State | 系统状态变更 | Fact Table 状态存储信息 + 输出推断 |

### Optional Dimensions

| Dimension | Content | When to Include |
|-----------|---------|-----------------|
| Side-effect | 外部副作用（hooks、网络调用、async Cmds） | Fact Table 发现副作用时 |
| Invariants (step-level) | 步骤内的属性约束 | 步骤有内部一致性要求时 |

Side-effect 省略时默认 `none`。步骤级 Invariants 省略时表示无约束。

### Preconditions Mutual Exclusivity

同一 Step 内不同 Outcome 的 Preconditions 必须互斥：对于任意系统状态，最多只有一个 Outcome 的 Preconditions 可被满足。

重叠时处理方式：
1. 区分 Preconditions（添加区分条件）
2. 无法区分时合并为单一 Outcome（析取 Preconditions）
3. 不得写出 Preconditions 可同时满足的 Outcome

超过 5 个 Outcome 的 Step 触发审阅检查点。

## Semantic Descriptors

所有维度值使用语义描述符——自然语言描述预期行为，表达业务意图而非精确匹配模式。

**规则**：
- 禁止包含 regex 语法（`\d`、`.*`、`[^...]`、`(?:...)`、`\s`、`\w`、`\b`、`$`、`^` 等）
- 禁止包含框架特定的断言模式
- 必须是表达业务意图的自然语言

正确示例：
- `"success confirmation containing feature-slug"`
- `"task status changed from pending to in_progress"`
- `"stderr contains error message about missing feature"`

错误示例（这些属于 gen-test-scripts 阶段）：
- `"Feature\s+([\w-]+)\s+created"`（regex）
- `"assert.Equal(t, 0, exitCode)"`（框架断言）

## Contract File Format

每个 Contract 存储为结构化 Markdown 文件，人可读且机器可解析。

### Template

```markdown
# Contract: <journey-name> / Step <N>: <step-description>

## Outcome "<outcome-name>"
- Preconditions: "<semantic description>"
- Input: <semantic description>
- Output: <semantic description>
- State: <semantic description>
- Side-effect: <semantic description or "none">
- Invariants: <step-level invariants or omit>

## Outcome "<outcome-name-2>"
- Preconditions: "<semantic description>"
- Input: <semantic description>
- Output: <semantic description>
- State: <semantic description>

## Journey Invariants
- <invariant description 1>
- <invariant description 2>
```

### Parseable Structure Rules

1. **Journey name**: 从文件路径提取：`tests/<journey-name>/_contracts/step-N-*.md`
2. **Step sequence**: 从文件名提取：`step-<N>-<slug>.md`
3. **Outcome sections**: `## Outcome "<name>"` 标题声明新 Outcome 块
4. **Dimension format**: 每行 `- <DimensionName>: <value>`
5. **Journey Invariants**: `## Journey Invariants` 节在每个 Contract 文件中必须出现且仅出现一次

## State Verification Levels

项目未暴露状态查询接口时，State 维度逐级降级：

| Level | Declaration | When |
|-------|-------------|------|
| `full` | 默认，无需标注 | 项目暴露状态查询接口（CLI flag、API 端点、文件读取） |
| `partial` | `<!-- state-verification: partial -->` | 状态字段仅能从 Output 推断 |
| `deferred` | `<!-- state-verification: deferred -->` + `limitations` 节 | 部分状态字段无法推断 |

## Fact Table Construction

gen-contracts 阶段必须从源代码中提取真实值，构建 Fact Table 作为 Contract 维度声明的事实基础。

### Generic Reconnaissance

| Source | What to Extract |
|--------|-----------------|
| CLI entry points | 命令名、flag 名、flag 类型、输出模式 |
| API handlers | 请求/响应 schema、状态码、中间件 |
| TUI model files | Model 结构体字段、Cmd 定义、Msg 类型、View 渲染 |
| Config files | 端口号、基础路径、超时值、认证机制 |
| State storage | 文件路径、JSON schema、数据库表 |
| Hook definitions | Hook 名称、触发条件、参数 schema |

### TUI-Specific Reconnaissance

| Source | What to Extract |
|--------|-----------------|
| Cmd definitions | Cmd 函数名、异步行为 |
| Batch usage | `tea.Batch()` 调用及其 Cmd 参数 |
| Timeout configs | 超时常量、默认等待时长 |
| Model transitions | Init → Idle → Processing → Result 状态 |

Fact Table 每条目必须标注源文件和行号。未知来源标记为 `UNKNOWN`，不得编造。

## Validation Rules

生成所有 Contract 后，逐一验证：

| Check | Rule |
|-------|------|
| Mandatory dimensions | 每个 Outcome 必须有非空的 Preconditions、Input、Output、State |
| Semantic descriptor purity | 维度值不得包含 regex 语法 |
| Outcome name uniqueness | 同一 Step 内 Outcome 名称必须唯一 |
| Preconditions exclusivity | 不同 Outcome 的 Preconditions 必须可区分 |
| Journey Invariants | 每个 Contract 文件必须有 `## Journey Invariants` 节且至少 1 条 |
| Side-effect default | 省略时默认 `none` |
| Outcome count | > 5 个 Outcome 触发审阅警告 |
| Unclassified points | 无法映射到维度的验证点归入 Invariants 并标注 `dimension: unclassified` |

## TUI Async Await Semantics

TUI 步骤涉及异步操作时，在 Contract 中声明 `await` 语义：

- `await <N>ms` = 等待所有待处理 Cmd 完成，最多 N 毫秒
- 默认超时：`.forge/config.yaml` 中的 `tui-await-timeout`，未配置时 3000ms
- 超时行为：fail-fast，报告超时的 Cmd 名称
- `tea.Batch(cmd1, cmd2)`：所有并发 Cmd 必须在继续前完成

异步 TUI 步骤必须包含一个超时 Outcome：

```
Outcome "timeout":
  Preconditions: "<包含异步 Cmd 超过 await 时长的条件>"
  Input: <同主 Outcome>
  Output: "error message containing timed-out Cmd name, fail-fast"
  State: "unchanged from pre-Cmd state"
```

## Directory Convention

```
tests/
  <journey-name>/                     # Journey 目录（kebab-case，永久位置）
    _contracts/                       # Contract 规范目录
      step-1-<action-slug>.md         # Step 1 的 Contract
      step-2-<action-slug>.md         # Step 2 的 Contract
      step-N-<action-slug>.md
    <test-file-1>                     # 生成的测试文件（直接在最终位置）
    <test-file-2>
```

### Rules

1. **Journey 目录**：以用户工作流命名（kebab-case），是所有相关测试文件的永久位置
2. **Contract 目录**：`_contracts/` 包含每个 Step 的 Contract 规范文件，命名 `step-<N>-<action-slug>.md`
3. **测试文件**：由 gen-test-scripts 直接生成到 Journey 目录，遵循项目测试框架命名约定
4. **无 staging 区域**：测试直接生成到最终位置，不经过中间目录

## Tag-Based Promotion

测试通过标签而非文件移动管理生命周期：

| Stage | Tag | Action |
|-------|-----|--------|
| New | `@feature` | 新生成的测试自动注入 |
| Promoted | `@regression` | `forge test promote <journey>` 升级 |
| CI selection | — | `forge test run --tags regression` 或 `--tags feature` |

标签使用测试框架的原生机制（Go build tag、pytest marker、describe 分组等）。

---

# Migration Guide: 旧测试模型 → Journey-Contract 模型

旧模型按接口类型分类测试、使用单步 TC 格式、依赖 staging+graduation 生命周期。本指南提供迁移到 Journey-Contract 模型的完整映射和步骤。

## Concept Mapping

| 旧模型概念 | 新模型概念 | 变化说明 |
|-----------|-----------|---------|
| 按接口类型分类（CLI/API/TUI/UI/Mobile） | 按 Journey（用户工作流）组织 | 组织维度从"技术接口"转向"用户场景" |
| 单步 TC（一条命令/一个端点） | Step + Contract（工作流中的一步） | TC 成为 Journey 中的 Step，保留完整工作流上下文 |
| TC Steps（自由格式操作列表） | Contract 六维声明（结构化验证） | 从非结构化描述升级为六维规范 |
| TC Expected（自由文本预期） | Outcome（Preconditions 互斥的多结果） | 支持每步多结果（成功、失败、边界情况） |
| Type 字段（CLI/API/TUI/UI/Mobile） | Journey 维度内的 interface-specific 描述 | 类型信息融入 Contract 维度而非顶层分类 |
| Setup 声明（API 数据准备） | Contract Preconditions + Fact Table | 数据准备作为 Preconditions 约束 |
| `{stepN.field}` 引用 | Contract 内的语义描述符 + gen-test-scripts 的精确匹配 | 声明阶段用自然语言，代码生成阶段用 Fact Table 精确匹配 |
| staging 目录 + graduation 流程 | Tag-Based Promotion（`@feature` → `@regression`） | 标签驱动而非文件移动 |
| 6 个硬编码 language profile | Convention 驱动（`docs/conventions/testing-<scope>.md`） | 可扩展、可编辑的 Convention 文件 |

## Directory Restructuring

### Old Structure

```
tests/e2e/
  features/
    <feature-name>/
      *_test.go            # 按 feature 分组的测试
  *_test.go                # 顶层散落的测试
```

### New Structure

```
tests/
  <journey-name>/           # 按用户工作流分组
    _contracts/
      step-1-<action>.md
      step-N-<action>.md
    <test-files>            # 直接在最终位置
```

### Restructuring Steps

1. **识别工作流**：从现有 TC 的 Steps 描述中提取连贯的用户工作流（如 `feature create → task claim → task submit`）
2. **创建 Journey 目录**：`tests/<journey-name>/`，如 `tests/task-lifecycle/`
3. **创建 Contract 目录**：`tests/<journey-name>/_contracts/`
4. **迁移测试文件**：将对应工作流的测试文件移入 Journey 目录
5. **移除旧目录**：确认迁移完成后删除 `tests/e2e/features/<feature-name>/` 空目录

## TC Format → Contract Conversion

### Single-Step TC → Single-Outcome Contract

```markdown
<!-- 旧格式 -->
## TC-001: Create feature
- **Type**: CLI
- **Steps**: Run `forge feature my-feature`
- **Expected**: Feature created successfully

<!-- 新格式：tests/create-feature/_contracts/step-1-feature-create.md -->
# Contract: create-feature / Step 1: feature create

## Outcome "success"
- Preconditions: "no feature with slug matching arg exists"
- Input: "forge feature my-feature"
- Output: "success confirmation containing feature-slug", exit code 0
- State: "feature directory created with manifest.md and tasks/index.json"

## Journey Invariants
- feature_slug consistent across all steps
```

### Multi-Step TC → Multi-Step Contract

```markdown
<!-- 旧格式 -->
## TC-067: Agent Flow — claim to submit lifecycle
- **Steps**:
  1. Run `forge feature wf-test-1`
  2. Run `forge task claim`
  3. Run `forge task submit T-wf-1 --result success`

<!-- 新格式：每个 Step 独立 Contract 文件 -->
# tests/task-lifecycle/_contracts/step-1-feature-create.md
## Outcome "success"
- Preconditions: "no feature with slug exists"
- Input: "forge feature wf-test-1"
- Output: "success confirmation containing wf-test-1"
- State: "feature directory created"

# tests/task-lifecycle/_contracts/step-2-task-claim.md
## Outcome "success"
- Preconditions: "feature exists, tasks available"
- Input: "forge task claim"
- Output: "task claimed, showing task_id"
- State: "task status changed from pending to in_progress"

# tests/task-lifecycle/_contracts/step-3-task-submit.md
## Outcome "success"
- Preconditions: "task in_progress status"
- Input: "forge task submit T-wf-1 --result success"
- Output: "submit confirmation"
- State: "task status changed to completed, record created"

## Journey Invariants (每个文件)
- feature_slug consistent across all steps
- task_id stable once assigned
```

### Conversion Rules

1. **TC Type → Journey interface**：TC 的 `Type` 字段（CLI/API/TUI）不再作为顶层分类，而是融入 Contract 维度描述
2. **TC Steps → Contract files**：多步 TC 的每个步骤变为独立的 Contract 文件（`step-N-<slug>.md`）
3. **TC Expected → Outcome**：预期结果转为 Outcome 的 Output + State 维度
4. **隐含前置条件 → Preconditions**：从步骤顺序和描述中提取前置条件（如"submit"隐含"task is in_progress"）
5. **边缘情况 → 额外 Outcome**：同一 Step 内的多个场景用互斥 Preconditions 的不同 Outcome 表达
6. **Integration TC → CLI/API Journey**：旧 `Type: Integration` 映射为 CLI 或 API Journey

## Lifecycle Migration: Staging/Graduation → Tag-Based Promotion

| 旧流程 | 新流程 | Migration Step |
|--------|--------|---------------|
| 测试生成到 `tests/e2e/` | 测试直接生成到 `tests/<journey>/` | 无 staging 目录 |
| `forge test stage <journey>` | 自动 `@feature` tag | 无需显式 staging |
| `forge test graduate <journey>` | `forge test promote <journey>` | 标签从 `@feature` 变为 `@regression` |
| 文件从 staging 移到 final location | 文件始终在最终位置 | 无文件移动 |
| CI 运行 final location 的测试 | `forge test run --tags regression` 或 `--tags feature` | 按标签选择 |

### Migration Steps

1. 确认现有测试已有 build tags 或 markers
2. 统一标签为 `@feature`（新建）或 `@regression`（已毕业）
3. 将测试文件移入对应的 Journey 目录
4. 运行 `forge test run --tags feature` 验证通过
5. 删除旧 staging 目录和 graduation 状态文件

## Migration Checklist

- [ ] 识别现有 TC 中的用户工作流，分组为 Journey
- [ ] 为每个 Journey 创建 `tests/<journey-name>/` 目录
- [ ] 为每个 Journey 创建 `_contracts/` 子目录
- [ ] 将 TC 转换为 Contract 文件（每个 Step 一个文件）
- [ ] 提取 Preconditions 确保互斥性
- [ ] 声明 Journey Invariants（至少 1 条）
- [ ] 迁移测试文件到 Journey 目录
- [ ] 统一标签为 `@feature` 或 `@regression`
- [ ] 删除旧 staging 目录和空 feature 子目录
- [ ] 运行 `forge test run --tags feature` 验证所有测试通过
- [ ] 运行 `/consolidate-specs` 将新 Contract 文件纳入漂移检测
