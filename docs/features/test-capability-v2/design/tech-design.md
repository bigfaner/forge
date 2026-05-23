---
created: 2026-05-23
prd: prd/prd-spec.md
status: Draft
---

# Technical Design: test-capability-v2

## Overview

Test Capability 2.0 将 Forge 测试管线从双路径架构统一为 Journey-Contract 单一路径，增加深度测试能力（风险驱动密度、surface 差异化、Run-to-Learn），并扩展通用性（内置 Convention 扩充、test-guide 自动检测）。

**核心约束**：主要改动在 Forge Plugin 层完成（skills/commands/rubrics/hooks），变更为 markdown 文件，通过 Plugin 分发机制到达用户环境。唯一例外：Fact Table 需要 Go CLI 子命令（`forge fact`）支持可靠的读写操作。

## Architecture

### Layer Placement

```
Forge Plugin Layer (plugins/forge/)
├── skills/          — Skill 定义 + 模板 + 规则 + rubrics
│   ├── gen-journeys/
│   │   └── rules/
│   │       └── surface-*.md — [NEW] 各 surface 类型的检测规则与策略指导
│   ├── gen-contracts/
│   ├── gen-test-scripts/
│   ├── run-tests/
│   └── ...
├── commands/        — 斜杠命令入口
├── hooks/           — 生命周期钩子
└── agents/          — subagent 定义（无变更）

Go CLI Layer (forge binary)
└── forge fact       — [NEW] Fact Table 读写子命令
    ├── list         — 列表展示 fact 摘要
    ├── get <id>     — 查看单条 fact 完整内容
    └── summary      — 按 source/confidence/kind 统计

Project Data
└── .forge/
    └── fact-table.json — [NEW] 跨 feature 的系统事实表
```

本功能不涉及 API 层、应用层、或数据库层。变更分布在 Plugin 层和 Go CLI 层。

### Component Diagram

```
                ┌─────────────────────────────────────────────────────┐
                │              Forge Plugin v3.0.0                    │
                │                                                      │
 ┌──────────┐   │  ┌──────────────┐  ┌───────────────┐               │
 │gen-test- │   │  │ gen-journeys │  │ gen-contracts │               │
 │cases/    │DEL │  │ (enhanced)   │  │ (enhanced)    │               │
 └──────────┘───→│  │ +surface det│  │ +required_out │               │
                │  │ +risk assign│  │ +risk density │               │
 ┌──────────┐   │  └──────┬───────┘  └───────┬───────┘               │
 │eval-test │   │         │                  │                        │
 │-cases    │DEL│         ▼                  ▼                        │
 └──────────┘───→│  ┌──────────────┐  ┌───────────────┐             │
                │  │ eval/        │  │ eval/         │             │
 6 per-type     │  │ +journey.md  │  │ +contract.md  │             │
 rubrics DEL    │  │ rubric       │  │ rubric        │             │
 ───────────→   │  └──────────────┘  └───────────────┘             │
                │                                                      │
                │  ┌──────────────┐  ┌───────────────┐               │
                │  │gen-test-     │  │ run-tests     │               │
                │  │scripts       │  │ (enhanced)    │               │
                │  │(enhanced)    │  │ +env-check    │               │
                │  │ +quality-val │  │ +confidence   │               │
                │  └──────────────┘  └───────────────┘               │
                │                                                      │
                │  ┌──────────────┐  ┌───────────────┐               │
                │  │ test-guide   │  │ gen-journeys/ │               │
                │  │ (enhanced)   │  │ rules/        │               │
                │  │ +auto-detect │  │ surface-*.md  │               │
                │  │ +draft-gen   │  │ [NEW]         │               │
                │  └──────────────┘  └───────────────┘               │
                │                                                      │
                │  ┌──────────────┐  ┌───────────────┐               │
                │  │ conventions/ │  │ commands/     │               │
                │  │ (migrated)   │  │ +eval-journey │               │
                │  │ +pytest.md   │  │ +eval-contract│               │
                │  │ +junit.md    │  │ -eval-test-cs │               │
                │  │ +rust.md     │  └───────────────┘               │
                │  └──────────────┘                                   │
                └─────────────────────────────────────────────────────┘

                ┌─────────────────────────────────────────────────────┐
                │              Go CLI Layer                           │
                │  ┌──────────────┐  ┌───────────────┐               │
                │  │ forge fact   │  │ .forge/       │               │
                │  │ [NEW]        │──│ fact-table.json│              │
                │  │ list/get/    │  │ [NEW]         │               │
                │  │ summary      │  └───────────────┘               │
                │  └──────────────┘                                   │
                └─────────────────────────────────────────────────────┘
```

### Dependencies

| 依赖 | 类型 | 用途 |
|------|------|------|
| `docs/conventions/testing-*.md` | 内部 | Convention 文件被 gen-test-scripts 和 run-tests 读取 |
| `plugins/forge/skills/eval/` | 内部 | eval-journey/eval-contract 复用现有 eval 框架 |
| `docs/conventions/testing-journey-contract.md` | 内部 | Journey-Contract 模型定义 |
| `docs/conventions/testing-conventions.md` | 内部 | Convention 文件结构定义（需迁移） |
| `forge fact` CLI | 内部 | gen-test-scripts/run-tests 通过 CLI 读取 Fact Table |
| `.forge/fact-table.json` | 内部 | 跨 feature 的系统事实数据 |
| Maestro CLI | 外部（可选） | Mobile surface 就绪检测 |
| pytest / JUnit / cargo test | 外部（用户侧） | 新增内置 Convention 对应的测试框架 |

## Interfaces

### Interface 1: Surface Rules (Markdown)

各 surface 类型的检测与策略指导，以 markdown rule 文件形式存在于 `skills/gen-journeys/rules/surface-*.md`。

每个 surface rule 文件包含以下指导性内容：

**检测信号** — 帮助 LLM 识别项目属于哪种 surface：
- 哪些文件模式、依赖、目录结构暗示该 surface 类型
- 哪些信号应排除（如 package.json 含 React 依赖排除 CLI 判断）
- 检测置信度判断原则

**通用测试指导原则** — 该 surface 类型应遵循的测试原则：
- CLI：测试 exit code、stdout/stderr、参数校验、并发安全；**必须隔离被测可执行文件**（独立编译产物、环境变量隔离、临时目录隔离）
- API：测试 status code、response schema、认证/授权、幂等性
- WebUI：测试用户交互流程、状态转换、可访问性；浏览器自动化框架由 Convention 定义（Playwright/Cypress/Selenium 等）
- TUI：测试键盘输入、终端渲染输出、异步 Cmd 超时处理
- Mobile：测试 app lifecycle、navigation、deep link；复杂场景标记 manual-only

**测试策略指导** — 以自然语言描述策略 reasoning：
- 该 surface 适合的测试层级侧重（contract vs journey 比例的指导原则及原因）
- 适用的执行模型（subprocess、browser automation、HTTP client 等）
- 环境就绪检测的关注点

**必须 Outcome 参考** — 该 surface 常见的边界/异常 Outcome 示例：
- 作为参考锚点而非硬性清单，LLM 结合项目实际情况判断
- 示例：CLI surface 常见边界包括"资源不存在"、"资源已存在"

**扩展方式**：新增 surface 类型只需添加一个 `surface-<type>.md` 文件，无需修改管线代码。

### Interface 2: Fact Table Entry (JSON)

```typescript
FactEntry = {
  fact_id: string         // 格式: "{subject}-{kind}-{nonce}"
  source: "static" | "runtime" | "manual"
  subject: string         // 如 "cli.forge", "api.GET /tasks"
  kind: "signature" | "output_format" | "error_code" | "side_effect" | "precondition" | "compilation_error" | "runtime_crash"
  value: object           // 自由格式 JSON，结构随 kind 变化
  confidence: "confirmed" | "inferred" | "assumed"
  updated_at: string      // ISO 8601 时间戳
}
```

**存储位置**：`.forge/fact-table.json`（项目级，跨 feature 共享。Monorepo 中每个子项目有独立 `.forge/` 目录）

**CLI 访问**：`forge fact` 子命令
- `forge fact list [--source static|runtime] [--confidence confirmed|inferred]` — 列表展示 fact 摘要
- `forge fact get <fact_id>` — 查看单条 fact 完整内容（value 原样 pretty-print）
- `forge fact summary` — 按 source/confidence/kind 统计，含覆盖率指标

**更新策略**（由 CLI 保证正确性）:
- Runtime fact 替换同一 `subject`+`kind` 的 static fact
- 若 runtime `confidence` 非 `confirmed`，static fact 保留为 fallback（不删除）
- 同一 `subject` 的不同 `kind` 共存

### Interface 3: Eval Rubric Dimension (journey.md / contract.md)

共享评分维度框架，总分 1000：

```yaml
# rubrics/journey.md frontmatter
scale: 1000
target: 850
iterations: 3
type: journey
context:
  conventions: []
  business-rules: auto
```

维度定义见 PRD Scope "Eval Rubric 评分维度框架"。每维度最低阈值：
- Completeness（完整性）: 120/200
- Semantic Purity（语义纯度）: 120/200
- Precondition Exclusivity（前置条件互斥性）: 90/150
- Fact Alignment（事实依据）: 90/150（区分"事实依据"和"合理推理"：LLM 衍生的 required_outcomes 归为推理声明，需标注规则依据）
- Surface Fitness（surface 适配）: 90/150
- Internal Consistency（一致性）: 90/150

### Interface 4: Convention File Schema (Migrated)

**迁移路径**: {Framework, Assertion, Tags, Result Format} → {framework, discovery, structure, assertions}

```markdown
## framework
- **name**: pytest                                    # 字符串，必需
- **version**: ">=7.0"                               # 字符串，可选
- **language**: python                                # 字符串，必需
- **runner_command**: "pytest {test_dir} {flags}"     # 字符串，必需

## discovery
- **test_dir**: tests/                                # 字符串，必需
- **file_pattern**: "test_*.py"                       # glob，必需
- **exclude_pattern**: ""                             # glob，可选

## structure
- **suite_pattern**: "class Test{Feature}"            # 字符串，必需
- **case_pattern**: "def test_{tc_id}_{description}"  # 字符串，必需
- **hook_pattern**: "@pytest.fixture"                 # 字符串，必需

## assertions
- **style**: assert                                   # 枚举: assert|expect|should，必需
- **custom_matchers**: ""                             # 字符串，可选
```

**迁移策略**：直接迁移，不保留旧 schema 兼容。
1. Phase 3: 一次性更新现有 3 个 Convention 文件到新 section 名称
2. gen-test-scripts 只读取新 section 名称（不回退旧名称）
3. 完全移除旧 section 名称处理逻辑

## Data Models

### Model 1: SurfaceDetectionResult

```
SurfaceDetectionResult = {
  detected_surface: string    // "cli" | "tui" | "webui" | "mobile" | "api" | "unknown"
  matched_signals: string[]   // 触发检测的信号列表
  confidence: "high" | "medium" | "low"
  all_signals: string[]       // 所有检测到的信号（诊断用）
}
```

### Model 2: ConfidenceRating

```
ConfidenceRating = {
  level: "HIGH" | "MEDIUM" | "LOW"
  confirmed_fact_ratio: float  // 0.0 - 1.0
  total_outcomes: int
  confirmed_outcomes: int
  eval_skipped: boolean
  eval_bypassed: boolean
}
```

### Model 3: RunToLearnConfig

```
RunToLearnConfig = {
  enabled: boolean                          // 通过 .forge/config.yaml 或 CLI 标志启用
  max_iterations: 3                         // 硬上限
  coverage_threshold: 0.80                  // 绝对退出阈值
  timeout_per_test: "60s"                   // 每个骨架测试超时，可通过 .forge/config.yaml 配置
  skip_on_env_failure: boolean              // 环境未就绪时跳过 R2L
}
```

## Error Handling

### Error Types & Exit Codes

Exit codes follow BIZ-error-reporting-001:

| Exit Code | Condition | Recovery |
|-----------|-----------|----------|
| 1 (retryable) | PRD 不存在或质量前置检查未通过 | 补全 PRD，重新运行 |
| 1 (retryable) | Surface 类型未知/混合 | 用户确认类型，重新运行 |
| 1 (retryable) | eval 迭代耗尽 (PAUSE) | 用户选择：跳过/放弃/修改 |
| 1 (retryable) | 环境检测失败 | 用户修复环境，重新运行 |
| 1 (retryable) | Convention 草稿被拒绝 2 次 | 用户手动编辑草稿 |
| 1 (retryable) | gen-contracts 合约 schema 验证失败 | 自动重试 1 次，然后手动 |
| 1 (retryable) | gen-test-scripts 质量验证失败 | 自动重试 1 次，失败则跳过 |
| 1 (retryable) | `forge fact` JSON 损坏或格式无效 | 用户手动修复或删除重建 |
| 2 (blocking) | eval 评分解析失败（重试后仍失败） | 检查 rubric 配置 |
| 0 (success) | FIX_DECIDE 重试耗尽 | 报告包含失败详情 |

### Propagation Strategy

所有 skill 错误通过 LLM 输出面向用户会话。无 HTTP/API 错误传播（CLI-only 功能）。错误信息包含：
1. 具体失败原因
2. 恢复提示（遵循 BIZ-error-reporting-002）
3. 当前管线状态上下文

## Testing Strategy

### Per-Layer Test Plan

| 层级 | 测试类型 | 工具 | 测试内容 | 覆盖目标 |
|------|----------|------|----------|----------|
| 管线统一 | 集成测试 | forge e2e testkit | 现有 Go/Vitest/Ginkgo 项目跑完整管线无报错 | 回归无损失 |
| Eval rubrics | 集成测试 | forge e2e testkit | eval-journey/eval-contract 评分、阈值门禁 | 100% |
| Per-surface 策略 | 按 surface 验证 | forge e2e testkit | 不同 surface 类型生成符合各自策略的测试（CLI→subprocess+可执行文件隔离、API→HTTP、WebUI→Convention 定义的浏览器框架、TUI→terminal I/O、Mobile→Maestro skeleton） | 每个 surface 至少 1 个 fixture 项目验证 |
| Convention 迁移 | Schema 验证 | Schema check | 迁移后的文件通过新 schema 验证 | 全部 3 个已有 + 3 个新增 |
| Run-to-Learn | 集成测试 | forge e2e testkit | 骨架测试生成、fact table 更新、超时保护 | 1 个端到端流程 |
| Fact Table CLI | 单元测试 | Go testing | `forge fact` list/get/summary、runtime→static 合并、fallback 逻辑 | 90% |

### Key Test Scenarios

1. **管线回归**：现有 Go/Vitest/Ginkgo 项目在迁移后跑完整管线无报错
2. **Per-surface 策略差异化**：CLI fixture → subprocess 断言 + 可执行文件隔离；API fixture → HTTP client 测试；WebUI fixture → Convention 定义的浏览器框架测试；TUI fixture → terminal I/O；Mobile fixture → Maestro YAML 骨架
3. **风险驱动密度**：高风险 Journey 生成 ≥ 13 个 Outcome，低风险 ≤ 7
4. **Eval 门禁**：Journey 评分低于阈值触发修订；某维度低于最低阈值即使总分 ≥ 850 也失败
5. **Convention 迁移**：迁移后的 Convention 文件通过新 schema 验证
6. **Run-to-Learn**：骨架测试输出丰富 Fact Table；超时保护生效
7. **Mobile best-effort**：未安装 Maestro CLI 仍可生成 Maestro YAML（无硬依赖）

### Overall Coverage Target

每个 surface 类型至少 1 个 fixture 项目验证策略差异化。关键集成路径（管线回归、eval 门禁、R2L）覆盖集成测试。配置文件（Convention markdown）通过 schema 验证。非关键增强（test-guide auto-detect、draft generation）以人工验证为主。

## Security Considerations

### Threat Model

| 风险 | 可能性 | 影响 |
|------|--------|------|
| Run-to-Learn 骨架测试执行任意代码 | M | H |
| Fact Table 数据泄露 | L | L |

### Mitigations

- **骨架测试沙箱**：在临时目录中运行，带超时保护（默认 60s，可通过 `.forge/config.yaml` 配置）。写操作使用 t.TempDir()。API 骨架测试仅发送 GET 请求；写操作生成回滚语句。
- **Fact Table**：本地存储于 `.forge/fact-table.json`。无外部传输。

## PRD Coverage Map

| PRD 需求 | 设计组件 | 文件/接口 |
|----------|----------|-----------|
| **Phase 1: 管线统一** | | |
| 删除 gen-test-cases skill | 移除目录 | `skills/gen-test-cases/` |
| 删除 test.graduate 任务类型 | 更新 hooks + ARCHITECTURE.md | `hooks/hooks.json` |
| 删除 eval-test-cases 命令 | 移除文件 | `commands/eval-test-cases.md` |
| 删除 test-cases rubrics | 移除 6 个文件 | `rubrics/test-cases.md` + 5 个 per-type |
| **Phase 2: 深度增强** | | |
| eval-journey skill | 新 rubric + 命令 | `rubrics/journey.md`, `commands/eval-journey.md` |
| eval-contract skill | 新 rubric + 命令 | `rubrics/contract.md`, `commands/eval-contract.md` |
| Surface 检测 | Surface rule 文件 | `skills/gen-journeys/rules/surface-*.md` |
| 风险驱动密度 | gen-contracts 规则增强 | `skills/gen-contracts/rules/risk-density.md` |
| 必须 Outcome | Surface rules + gen-contracts 增强 | `skills/gen-journeys/rules/surface-*.md` |
| gen-test-scripts 增强 | 按 surface 策略规则 | `skills/gen-test-scripts/types/*.md` |
| Run-to-Learn | gen-test-scripts 新规则 | `skills/gen-test-scripts/rules/run-to-learn.md` |
| 环境就绪检测 | run-tests 增强 | `skills/run-tests/rules/env-check.md` |
| 置信度评级 | run-tests 增强 | `skills/run-tests/rules/confidence.md` |
| Fact Table CLI | Go CLI 子命令 | `forge fact` (list/get/summary) |
| **Phase 3: 通用性扩展** | | |
| Convention schema 迁移 | 更新 schema + 迁移 | `docs/conventions/testing-conventions.md` |
| pytest Convention | 新文件 | `docs/conventions/testing-pytest.md` |
| JUnit Convention | 新文件 | `docs/conventions/testing-junit.md` |
| Rust Convention | 新文件 | `docs/conventions/testing-rust.md` |
| test-guide 自动检测 | 增强 | `skills/test-guide/rules/signal-detection.md` |
| test-guide 草稿生成 | 增强 | `skills/test-guide/rules/draft-generation.md` |
| **Stories** | | |
| Story 1: 管线统一 | Phase 1 全部组件 | 见上 |
| Story 2: 深度测试 | risk-density + required_outcomes + strategy | `gen-contracts/rules/`, `gen-journeys/rules/surface-*.md` |
| Story 3: Surface 差异化 | 按 surface 策略规则 | `gen-journeys/rules/surface-*.md`, `gen-test-scripts/types/` |
| Story 4: Eval 门禁 | eval rubrics + 命令 | `rubrics/journey.md`, `rubrics/contract.md` |
| Story 5: 通用性 | Convention 文件 + test-guide | `docs/conventions/`, `test-guide/` |
| Story 6: Run-to-Learn | R2L 机制 + Fact Table CLI | `gen-test-scripts/rules/run-to-learn.md`, `.forge/fact-table.json`, `forge fact` CLI |
| Story 7: 可扩展 surface | surface rules + 扩展指南 | `gen-journeys/rules/surface-*.md` |

## Open Questions

- [x] ~~Convention schema 迁移：是否需要 `forge migrate-convention` 命令？~~ **已解决**：直接迁移，不保留兼容。Phase 3 中原位更新文件。
- [x] ~~Surface 配置格式：YAML 配置文件 vs markdown rules？~~ **已解决**：Markdown rules 存放在 `skills/gen-journeys/rules/surface-*.md`。LLM-based skills 消费自然语言指导，不需要结构化配置文件。
- [x] ~~Fact Table 存储位置：`docs/features/<slug>/testing/` vs `.forge/`？~~ **已解决**：`.forge/fact-table.json`（项目级，跨 feature 共享）。通过 `forge fact` CLI 子命令读写，保证更新语义正确性。

## Appendix

### Alternatives Considered

| 方案 | 优点 | 缺点 | 未选择原因 |
|------|------|------|-----------|
| 为 eval-journey/eval-contract 新建独立 skill | 职责清晰 | 重复 eval 框架逻辑，更多文件维护 | 现有 eval 框架已支持参数化 rubric |
| 在 Go 二进制中硬编码 surface 检测 | 更快、类型安全 | 需要 Go 代码改动，不可扩展 | Markdown rules 允许社区扩展而无需重新构建 |
| YAML 配置文件定义 surface 类型 | 结构化、可机器解析 | 对 LLM 消费过于死板，增加不必要的抽象层 | Markdown rules 的原则 + 示例给 LLM 兼具指导性和灵活性 |
| 保留旧 Convention schema + 新增可选 | 零迁移成本 | 两种 schema 永久共存，造成混乱 | 直接迁移更简单；仅 3 个现有文件需更新 |
| Fact Table 用 SQLite | 更好的查询能力 | 数据量级不匹配，增加依赖 | JSON 文件足够且透明 |
| LLM 直接读写 JSON Fact Table | 无需 Go 代码 | LLM 操作 JSON 易出错，更新语义（merge/fallback）无法保证 | CLI 子命令保证结构化读写正确性 |

### References
- PRD Spec: `docs/features/test-capability-v2/prd/prd-spec.md`
- Journey-Contract Model: `docs/conventions/testing-journey-contract.md`
- Convention File Structure: `docs/conventions/testing-conventions.md`
- Forge Distribution Model: `docs/conventions/forge-distribution.md`
