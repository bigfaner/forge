---
created: "2026-06-02"
author: "fanhuifeng"
status: Draft
intent: "refactor"
---

# Proposal: Surface-First 测试约定重构

## Problem

Forge 的测试约定（Convention）体系以框架（Go/Vitest/pytest）为首要组织轴，导致 agent 在编写测试时无法根据 Surface 类型（cli/api/web/tui/mobile）获取正确的测试策略和文件位置规范。

### Evidence

1. **文件位置缺失**：test-guide 生成的 Convention 文件不包含 per-surface 的目录规则。agent 执行任务后写测试时，不知道 CLI 测试应放 `tests/cli/`、Web E2E 放 `tests/e2e/`。
2. **测试策略缺失**：Convention 文件只覆盖框架语法（断言库、runner 命令），不包含 Surface 特有的隔离模型、超时策略、断言重点等策略信息。
3. **职责分裂**：gen-test-scripts 的 `types/cli.md` 等文件包含完整测试策略，但 Convention 文件（agent 在非 skill 调用场景下的参考）中完全没有。agent 在闲聊或手动写测试时只能靠猜测。
4. **test-type-model.md 未分发**：`docs/reference/test-type-model.md` 定义了 Surface → Test Type 的权威映射，但该文件位于 Forge 项目自身 docs/ 下，不分发给用户。用户项目的 agent 无法访问。

### Urgency

Surface 是 Forge 测试管线的核心概念——orchestration 生命周期、test type 分类、justfile recipe 命名全部围绕 Surface 组织。测试约定不与 Surface 对齐，意味着每次 agent 写测试都是"策略盲区"。随着 surface 类型扩展，问题会持续恶化。

## Proposed Solution

将测试约定从 **framework-first** 重构为 **surface-first**，同时将核心测试知识融入 Forge plugin 确保用户可达。

### 核心变更

**1. Convention 文件结构重组**（用户项目层）

```
docs/conventions/testing/
  index.md              # 速查表：Surface → 类型 → 位置 → 断言重点
  cli/
    index.md            # 文档索引
    core.md             # CLI 测试策略（语言无关）
  api/
    index.md
    core.md             # API 测试策略（语言无关）
  web/
    index.md
    core.md             # Web E2E 策略（语言无关）
  tui/
    index.md
    core.md
  mobile/
    index.md
    core.md
```

每个 `core.md` 包含：文件位置、隔离模型、断言重点、超时策略、生命周期、Contract/Journey 比例、fixture 规则、反模式、断言偏好表（per-framework 一行）。不包含 framework 文件——LLM 自带语言知识，无需额外约定。

**2. guide.md hook 补充测试速查表**（Plugin 层，始终可达）

在 guide.md 新增 Testing section，包含：
- Surface → Test Type 映射表
- e2e 术语约束
- 测试文件位置规则
- 引导用户运行 `/test-guide` 获取完整策略

**3. test-type-model.md 融入 Plugin**（Plugin 层）

精华版（映射表 + e2e 约束）融入 guide.md。完整版（分类标准 + 语义定义）放入 `plugins/forge/skills/test-guide/references/test-type-model.md`。

**4. test-guide skill 重写**

- 新增 Step：读取 `.forge/config.yaml` 获取 Surface 类型
- 从内置模板（`templates/surfaces/*.md`）生成 per-surface 的 index.md + core.md
- 生成顶层 `docs/conventions/testing/index.md` 速查表
- 移除旧框架检测 → 单文件生成的流程

**5. 下游 skill 适配**

- gen-test-scripts：Convention 加载改为 surface 目录遍历（`testing/{surface}/core.md`）
- run-tests：Convention 读取路径改为 `testing/{surface}/`
- init-justfile：Test recipe 根据新目录结构生成

### Innovation Highlights

这不是行业常见做法的移植，而是从第一性原理推导的设计决策：

- **Surface 作为首要组织轴**：CLI 测试在 Go 和 Python 之间的共同点（子进程隔离、exit code 断言、stdout 解析），远大于 CLI 测试和 Web E2E 测试在 Go 中的共同点。测试策略的本质差异在 Surface，不在语言。
- **零 Framework 文件**：LLM 已有语言/framework 知识，Convention 文件只补充 LLM 推断不出的项目特有规则。消除了传统"代码模板"思路带来的维护负担。
- **guide.md 作为冷启动答案**：用户初始项目不存在任何 Convention 文件，但 guide.md hook 始终在 context 中，确保 agent 从第一个会话就能正确回答测试相关问题。

## Requirements Analysis

### Key Scenarios

1. **用户初始项目闲聊**：用户问"CLI 测试应该放哪？"，agent 从 guide.md 的速查表直接回答，无需任何文件加载。
2. **test-guide 生成**：用户运行 `/test-guide`，skill 检测 surface 类型，从内置模板生成 per-surface convention 文件。
3. **gen-test-scripts 消费**：skill 加载 `docs/conventions/testing/{surface}/core.md`，获取隔离模型和断言重点，生成测试代码。
4. **run-tests 执行**：skill 从 `docs/conventions/testing/{surface}/core.md` 读取超时策略和生命周期规则。
5. **agent 手动写测试**：agent 读 guide.md 知道文件位置，读 core.md 知道策略细节，用 LLM 知识写具体代码。
6. **多 Surface 项目**：项目同时有 CLI + API surface，test-guide 为每个 surface 生成独立的 core.md，agent 按需加载。

### Non-Functional Requirements

- **guide.md 增量 <= 20 行**：guide.md 是每次会话的固定 token 成本，必须精简。
- **Convention 加载路径向后兼容**：如果用户项目中存在旧结构 convention 文件，下游 skill 应能识别并给出迁移提示，而非静默失败。
- **零 framework 文件维护**：新增语言/framework 支持时，只需 LLM 知识即可，不需新增 Convention 文件。

### Constraints & Dependencies

- test-guide 的 surface 检测依赖 `.forge/config.yaml` 中的 surfaces 配置
- gen-test-scripts 和 run-tests 的 convention 加载逻辑需要同步修改
- guide.md 是 hook 文件，修改后影响所有用户的每次会话

## Alternatives & Industry Benchmarking

### Industry Solutions

测试约定的组织方式通常有两种：
1. **按测试层级**：unit/integration/e2e 目录（大多数框架的默认结构）
2. **按功能模块**：每个模块下放该模块的所有测试

两者都不以"系统入口类型"为首要轴。Forge 的 Surface 概念是测试管线特有的——它决定了测试的执行模型（子进程 vs HTTP vs 浏览器），这比测试层级或功能模块更根本。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零改动 | agent 每次猜测试位置和策略；test-type-model 不分发 | Rejected: 问题持续恶化 |
| 框架优先 + Surface 段落 | 现有结构 | 改动最小 | 单文件膨胀；新增 surface 仍需改所有框架文件 | Rejected: 根本问题未解决 |
| 渐进叠加 | — | 安全，可逐步迁移 | 长期维护两套结构；加载逻辑复杂 | Rejected: 技术债 |
| **Surface-first 全链路重构** | 第一性原理推导 | 信息就近；新增 surface 只加目录；guide.md 冷启动 | 全量改动影响面大 | **Selected: 根本性解决问题** |

## Feasibility Assessment

### Technical Feasibility

完全可行。所有改动都在 plugin 层面（skills、hooks、templates），不涉及 forge CLI 二进制或配置格式变更。Surface 检测已有 `.forge/config.yaml` 支持。

### Resource & Timeline

改动涉及 4 个 skill（test-guide、gen-test-scripts、run-tests、init-justfile）+ 1 个 hook（guide.md）+ Convention 模板。scope 清晰，无外部依赖。

### Dependency Readiness

所有前置条件已就绪：Surface 类型定义（5 种）、surface 检测机制（`.forge/config.yaml`）、gen-test-scripts 的 type rules 已有 per-surface 策略可参考。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 测试约定应按框架组织，因为不同语言语法差异大 | Assumption Flip：如果按 surface 组织呢？ | Overturned：CLI 测试在 Go 和 Python 间的策略共性远大于 CLI 和 Web E2E 在 Go 中的共性 |
| Convention 文件需要包含 framework 代码模板 | Occam's Razor：LLM 已有 framework 知识，还需要吗？ | Overturned：只需补充 LLM 不知道的 Forge 项目特有规则 |
| test-type-model.md 放在 docs/reference/ 即可 | XY Detection：用户的 agent 能访问吗？ | Overturned：docs/ 不分发，必须融入 plugin |
| guide.md 不应放测试内容 | Stress Test：用户初始项目没有任何 Convention 文件时，agent 怎么回答？ | Refined：guide.md 必须包含最小测试规则集，但控制在 20 行内 |

## Scope

### In Scope

**Plugin 层（分发）**

1. 重写 `test-guide` skill：新增 surface 检测，生成 surface-first convention 文件
2. 新增 surface 策略模板：`plugins/forge/skills/test-guide/templates/surfaces/`（cli/api/web/tui/mobile 各一个 core.md 模板）
3. 新增完整 test-type-model：`plugins/forge/skills/test-guide/references/test-type-model.md`
4. 更新 `guide.md` hook：新增 Testing section（速查表 + 映射表 + e2e 约束）
5. 更新 `gen-test-scripts`：Convention 加载逻辑改为 surface 目录遍历
6. 更新 `run-tests`：Convention 读取路径改为 `testing/{surface}/`
7. 更新 `init-justfile`：Test recipe 根据新目录结构生成

**Forge 项目层（不分发）**

8. 删除旧 convention 文件：`docs/conventions/testing/` 下 6 个框架文件 + index.md
9. 删除旧 test-type-model：`docs/reference/test-type-model.md`（内容已移入 plugin）
10. 用新 test-guide 重新生成 forge 项目自身的 `docs/conventions/testing/cli/`

### Out of Scope

- gen-test-scripts type rules（`types/cli.md` 等）— 已有完整策略，不改动
- run-tests surface orchestration rules（`rules/surfaces/`）— 生命周期不变
- `.forge/config.yaml` 的 surface 配置机制 — 不变
- 用户项目中旧 convention 文件的自动迁移脚本
- test-guide 的 framework 检测逻辑移除（保留，用于 core.md 中的断言偏好表填充）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 全量改动影响面大，下游 skill 加载逻辑遗漏 | M | H | 逐 skill 验证：每个 skill 的 Convention 加载点逐一排查，确保无遗漏 |
| guide.md 新增内容超出 20 行限制 | L | M | 写完后量化行数，超标则压缩 |
| 用户已有旧结构 convention 文件，新版 skill 无法识别 | M | M | gen-test-scripts/run-tests 检测到旧结构时输出迁移提示，而非静默失败 |
| core.md 模板中的断言偏好表不覆盖用户实际使用的 framework | L | L | 表格包含 Forge 已支持的 6 个框架，未覆盖的由 LLM 知识兜底 |

## Success Criteria

- [ ] guide.md 新增 Testing section <= 20 行，包含完整的 Surface 速查表和 e2e 约束
- [ ] test-guide 能根据 `.forge/config.yaml` 检测 surface 类型并生成 per-surface convention 文件（index.md + core.md）
- [ ] 每个 core.md 包含 7 个必要段落：文件位置、隔离模型、断言重点、超时策略、生命周期、Contract/Journey 比例、反模式
- [ ] gen-test-scripts 能从 `docs/conventions/testing/{surface}/core.md` 加载策略信息
- [ ] run-tests 能从 `docs/conventions/testing/{surface}/core.md` 读取超时和生命周期规则
- [ ] 用户初始项目中不存在任何 convention 文件时，agent 能从 guide.md 正确回答测试位置和策略问题
- [ ] `docs/reference/test-type-model.md` 的内容 100% 融入 plugin（guide.md 精华版 + test-guide/references/ 完整版）

## Next Steps

- Proceed to `/write-prd` to formalize requirements
