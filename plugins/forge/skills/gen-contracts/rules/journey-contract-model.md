---
name: journey-contract-model
description: Core concepts of the Journey-Contract test model, directory conventions, tag-based promotion, and migration guide from old test model
---

# Journey-Contract Test Model

Forge 测试管道以用户工作流（Journey）为组织单元，通过 Contract 六维声明定义预期行为，由 Tag-Based Promotion 管理测试生命周期。

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

1. **Journey name**: 从文件路径提取：`docs/features/<slug>/testing/<journey>/contracts/step-N-*.md`
2. **Step sequence**: 从文件名提取：`step-<N>-<slug>.md`
3. **Outcome sections**: `## Outcome "<name>"` 标题声明新 Outcome 块
4. **Dimension format**: 每行 `- <DimensionName>: <value>`
5. **Journey Invariants**: `## Journey Invariants` 节在每个 Contract 文件中必须出现且仅出现一次

## Directory Convention

```
docs/features/<slug>/testing/
  <journey-name>/                     # Journey 目录（kebab-case）
    journey.md                        # Journey 叙述文档
    contracts/                        # Contract 规范目录
      step-1-<action-slug>.md         # Step 1 的 Contract
      step-2-<action-slug>.md         # Step 2 的 Contract
      step-N-<action-slug>.md

tests/
  <journey-name>/                     # 生成的测试文件
    <test-file-1>
    <test-file-2>
```

### Rules

1. **Journey 目录**：`docs/features/<slug>/testing/<journey>/` 以用户工作流命名（kebab-case），包含 journey.md 和 contracts/
2. **Contract 目录**：`testing/<journey>/contracts/` 包含每个 Step 的 Contract 规范文件，命名 `step-<N>-<action-slug>.md`
3. **测试文件**：由 gen-test-scripts 直接生成到 `tests/<journey>/`，遵循项目测试框架命名约定
4. **无 staging 区域**：测试直接生成到最终位置，不经过中间目录

## Tag-Based Promotion

测试通过标签而非文件移动管理生命周期：

| Stage | Tag | Action |
|-------|-----|--------|
| New | `@feature` | 新生成的测试自动注入 |
| Promoted | `@regression` | `/run-tests` 自动升级 |
| CI selection | -- | Use test framework's native tag filter (e.g., `go test -tags=regression`, `pytest -m regression`) |

标签使用测试框架的原生机制（Go build tag、pytest marker、describe 分组等）。

## Migration Guide: 旧测试模型 -> Journey-Contract 模型

旧模型按接口类型分类测试、使用单步 TC 格式、依赖 staging+graduation 生命周期。本指南提供迁移到 Journey-Contract 模型的完整映射和步骤。

### Concept Mapping

| 旧模型概念 | 新模型概念 | 变化说明 |
|-----------|-----------|---------|
| 按接口类型分类（CLI/API/TUI/Web/Mobile） | 按 Journey（用户工作流）组织 | 组织维度从"技术接口"转向"用户场景" |
| 单步 TC（一条命令/一个端点） | Step + Contract（工作流中的一步） | TC 成为 Journey 中的 Step，保留完整工作流上下文 |
| TC Steps（自由格式操作列表） | Contract 六维声明（结构化验证） | 从非结构化描述升级为六维规范 |
| TC Expected（自由文本预期） | Outcome（Preconditions 互斥的多结果） | 支持每步多结果（成功、失败、边界情况） |
| Type 字段（CLI/API/TUI/Web/Mobile） | Journey 维度内的 interface-specific 描述 | 类型信息融入 Contract 维度而非顶层分类 |
| Setup 声明（API 数据准备） | Contract Preconditions + Fact Table | 数据准备作为 Preconditions 约束 |
| `{stepN.field}` 引用 | Contract 内的语义描述符 + gen-test-scripts 的精确匹配 | 声明阶段用自然语言，代码生成阶段用 Fact Table 精确匹配 |
| staging 目录 + graduation 流程 | Tag-Based Promotion（`@feature` -> `@regression`） | 标签驱动而非文件移动 |
| 6 个硬编码 language profile | Convention 驱动（`docs/conventions/testing-<scope>.md`） | 可扩展、可编辑的 Convention 文件 |

### Directory Restructuring

#### Old Structure

```
tests/e2e/
  features/
    <feature-name>/
      *_test.go            # 按 feature 分组的测试
  *_test.go                # 顶层散落的测试
```

#### New Structure

```
docs/features/<slug>/testing/
  <journey-name>/               # 按用户工作流分组
    journey.md
    contracts/
      step-1-<action>.md
      step-N-<action>.md

tests/
  <journey-name>/               # 生成的测试文件
    <test-files>
```

### Tag-Based Promotion

| 流程 | 说明 |
|------|------|
| 测试直接生成到 `tests/<journey>/` | 无 staging 目录 |
| 自动 `@feature` tag | 无需显式 staging |
| `/run-tests` (tag promotion) | 标签从 `@feature` 变为 `@regression` |
| 文件始终在最终位置 | 无文件移动 |
| 使用测试框架原生的标签过滤机制（如 `go test -tags=regression`） | 按标签选择 |
