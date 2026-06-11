---
created: 2026-05-24
author: faner
status: Draft
---

# Proposal: Test Recipe Unification

## Problem

Forge 的测试基础设施已完成去 Profile 化重构（Convention 驱动 + config.yaml），但 justfile recipe 模型未同步适配，导致三层错位：submit 门禁跑全量测试浪费时间、模板仍硬编码 Playwright、配置项命名过时。

### Evidence

- Breaking 任务 submit 运行 `just test`（即 `go test -race ./...`），单次 ~90s，agent 全程等待
- init-justfile 模板对所有语言（Go/Python/Rust）都生成 Playwright 的 `e2e-test` recipe，与实际项目不符
- `auto.e2eTest` 配置项名称仅描述 e2e 场景，但实际控制的是高级测试 pipeline 的生成，不仅限于 e2e
- `journey_isolation.go` 硬编码 `just e2e-test`，不走 config 驱动
- `testrunner.go` 探测链首选项 `just test`，但 `test` 在新模型下不再代表单元测试

### Urgency

当前 breaking 任务每次 submit 等待 ~90s（全量测试），10 个 breaking 任务累计 15 分钟纯等待。同时 init-justfile 生成的模板与实际项目语言不匹配，导致用户需要手动修改 every time。v3.0.0 分支正在进行测试能力重构，现在对齐可避免后续返工。

## Proposed Solution

引入两层测试 recipe 模型，解耦语言级单元测试与 surface 级高级测试：

- **`just unit-test`**：语言级单元测试（Go: `go test ./...`，Python: `pytest`），快速反馈，用于 per-task submit 门禁
- **`just test`**：Surface 级高级测试（Web UI → e2e，API → 集成测试），用于 all-completed 门禁
- **`just test-setup`**（可选）：测试环境准备（Playwright 安装、DB seed）
- **`just probe`**（可选）：服务健康检查

淘汰 `e2e-test`、`e2e-setup`、`e2e-verify` 命名，统一归入 `test` + 辅助 recipe。配置项 `auto.e2eTest` 更名为 `auto.test`。任务 key `run-e2e-tests` 更名为 `run-test`。

### Innovation Highlights

非创新性方案。将业界标准的分层测试（unit → integration → e2e）映射到 Forge 的 justfile 约定，核心决策是让 Forge 保持 surface-agnostic——不区分 e2e 还是集成测试，只调用 `just test`，由 recipe 内部按 surface 分发。这避免了 Forge 与具体测试框架的耦合。

## Requirements Analysis

### Key Scenarios

- **Breaking 任务 submit**：agent 完成编码后 submit，门禁跑 compile → fmt → lint → unit-test，快速反馈（<30s），通过则提交
- **非 Breaking 任务 submit**：跑 compile → fmt → lint（无变化）
- **All-completed hook**：跑 compile → fmt → lint → unit-test → test → probe（如需），完整验证
- **Journey isolation**：`just test <journeyName>` 运行单个 journey 的高级测试
- **init-justfile 生成**：按语言生成 `unit-test`，按 surface 生成 `test`，不再硬编码 Playwright
- **无 Fallback**：v3.0.0 直接要求 `unit-test` recipe，不回落到 `test`

### Non-Functional Requirements

- **性能**：Breaking 任务 submit 门禁耗时从 ~90s 降至 ~20s（Go 项目无 `-race` 的单元测试）
- **无向后兼容负担**：v3.0.0 大版本重构，直接重命名，不保留旧键名/旧 recipe 的兼容逻辑
- **可发现性**：`auto.test` 命名直观，反映控制的是高级测试 pipeline

### Constraints & Dependencies

- init-justfile 模板需按语言模板生成 `unit-test`（Go/Node/Python/Rust 各不同）
- `journey_isolation.go` 需从 `just e2e-test` 迁移到 `just test`
- `forgeconfig.Config` 的 YAML 直接使用 `test` 键名，不保留 `e2eTest` 兼容
- 项目根 justfile 需同步更新（增加 `unit-test` recipe）

## Alternatives & Industry Benchmarking

### Industry Solutions

分层测试（Test Pyramid）是业界标准：单元测试（快速、大量）→ 集成测试（中速、适量）→ E2E 测试（慢、少量）。CI/CD 通常在 commit 时跑单元测试，merge/merge request 时跑集成+E2E。本方案将相同原则映射到 Forge 的 submit vs all-completed 时机。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 每次提交等待 90 秒；模板与语言不匹配；配置项命名过时 | Rejected: 累计等待时间不可接受 |
| 仅加 `test-quick` recipe | — | 最小改动 | 新增 recipe 但不解决模板/配置命名/Playwright 硬编码问题 | Rejected: 半吊子方案 |
| 统一 + 缓存 + 并行 | — | 最大性能提升 | 复杂度过高，缓存机制引入新 bug 面 | Rejected: 过度设计 |
| **两层 recipe 模型** | Test Pyramid | 结构清晰；Forge surface-agnostic；向后兼容 | 需更新多个组件 | **Selected: 直击痛点，复杂度可控** |

## Feasibility Assessment

### Technical Feasibility

完全可行。Go 侧改动集中在 `pkg/just/just.go`（新增 gate sequence）、`submit.go`（使用新 sequence）、`quality_gate.go`（调整 step 2/3）、`forgeconfig/config.go`（字段重命名）。模板改动为文本替换。无外部依赖。

### Resource & Timeline

涉及 ~42 个文件，影响范围评估如下。预估 10-12 个 coding task。

### Dependency Readiness

无外部依赖。Forge 自身 justfile 作为第一个适配目标。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| Breaking 任务 submit 需要跑完整测试 | 5 Whys | Refined: 快速反馈只需 unit-test，完整测试留给 all-completed |
| Forge 需要按 surface 提供专用 recipe（e2e-test、integration-test） | Occam's Razor | Overturned: `just test` 抽象 surface 差异，Forge 无需知道内部实现 |
| `test-verify` 是必要的标准 recipe | Need Gate | Overturned: 可内联为 skill 中的 grep 命令，不占标准配方 |
| `auto.e2eTest` 命名准确反映其职责 | Assumption Flip | Refined: 它控制的是所有高级测试 pipeline，不仅 e2e，更名 `auto.test` |

## Impact Analysis

影响范围评估覆盖三组重命名，按优先级分层。

### Tier 1: Go Source Code（必须改动）

| File | Change |
|------|--------|
| `pkg/just/just.go` | `DefaultGateSequence()` 的 `test` → `unit-test`；新增 `UnitGateSequence()` |
| `pkg/testrunner/testrunner.go` | `RunProjectTests()` 探测和调用 `unit-test` 替代 `test` |
| `pkg/testrunner/journey_isolation.go` | `just e2e-test` → `just test` |
| `internal/cmd/task/submit.go` | breaking 任务使用新 gate sequence |
| `internal/cmd/quality_gate.go` | Step 2 用 `unit-test`，Step 3 用 `test`；`addFixTask` 映射更新 |
| `pkg/forgeconfig/config.go` | `E2eTest` → `Test`，YAML tag `e2eTest` → `test`，双键兼容读取 |
| `pkg/task/autogen.go` | `auto.E2eTest` → `auto.Test`；`Key: "run-e2e-tests"` → `"run-test"` |
| `internal/cmd/init.go` | init wizard 提示文案更新 |

### Tier 2: Prompt Templates + Skill/Command Markdown（必须改动）

| File | Change |
|------|--------|
| `pkg/prompt/data/gate.md` | `just test` → `just unit-test` |
| `pkg/prompt/data/fix-record-missed.md` | 同上 |
| `pkg/prompt/data/validation-code.md` | 同上 |
| `commands/fix-bug.md` | 所有 `just test` → `just unit-test`；`just e2e-test` → `just test` |
| `skills/clean-code/SKILL.md` | `just test` → `just unit-test` |
| `skills/gen-test-scripts/rules/run-to-learn.md` | `just test` → `just unit-test` |
| `skills/init-justfile/SKILL.md` | Standard Target Contract 全面更新 |
| `skills/run-tests/SKILL.md` | config-schema 示例中 recipe 名更新 |
| `skills/run-tests/references/config-schema.md` | 同上 |

### Tier 3: Justfile Templates（必须改动）

| File | Change |
|------|--------|
| `templates/generic.just` | `test` → `unit-test`；淘汰 `e2e-test`/`e2e-setup`/`e2e-verify`；增加 `test`/`test-setup` |
| `templates/go.just` | 同上 |
| `templates/node.just` | 同上 |
| `templates/python.just` | 同上 |
| `templates/rust.just` | 同上 |
| `templates/mixed.just` | 同上 |
| Root `justfile` | `test` → `unit-test`；新增 `test` recipe；`ci` recipe 更新 |

### Tier 4: Go Test Files（必须改动）

| File | Change |
|------|--------|
| `pkg/just/just_test.go` | 断言 `test` → `unit-test` |
| `pkg/testrunner/testrunner_test.go` | justfile fixture `test:` → `unit-test:` |
| `internal/cmd/quality_gate_test.go` | `HasRecipe(dir, "test")` → `HasRecipe(dir, "unit-test")` |
| `forgeconfig/config_test.go` | 所有 `E2eTest` 断言 → `Test` |
| `task/autoconfig_test.go` | `auto.E2eTest` → `auto.Test` |
| `task/autogen_test.go` | `E2eTest` fixture + `run-e2e-tests` key → `run-test` |
| `task/submit_test.go` | `run-e2e-tests` fixtures → `run-test` |
| `task/status_test.go` | `run-e2e-tests` fixtures → `run-test` |
| `tests/justfile-integration/mixed_cli_test.go` | `just test` 断言 → `just unit-test`；修复已失败的 TC_005/TC_015/TC_016 |
| `tests/justfile-integration/forge_detection_test.go` | recipe 列表 `test` → `unit-test` |
| `tests/task-type-system/task_types_dispatch_test.go` | `just test` 断言 → `just unit-test` |

### Tier 5: Documentation（应当更新）

| Area | Files | Priority |
|------|-------|----------|
| CLI docs (OVERVIEW.md, WORKFLOW.md, zh versions) | 4 | High |
| ARCHITECTURE.md | 1 | High |
| business-rules/quality-gate.md | 1 | High |
| docs/conventions/testing/go.md | 1 | Medium |
| docs/conventions/forge-distribution.md | 1 | Medium |
| lessons/ (~9 files) | ~9 | Low（历史记录，可选更新） |
| proposals/ (历史提案中的 e2eTest 引用) | ~5 | Low（历史文档，不强制更新） |

### Config & Testdata

| File | Change |
|------|--------|
| `.forge/config.yaml` | `e2eTest` → `test` |
| `internal/cmd/testdata/forge-config.example.yaml` | `e2eTest` → `test` |
| `internal/cmd/testdata/forge-config.schema.json` | `e2eTest` → `test` |

## Scope

### In Scope

**Go 代码（8 files）**
- `pkg/just/just.go`：`DefaultGateSequence()` 的 `test` → `unit-test`；新增 `UnitGateSequence()`
- `pkg/testrunner/testrunner.go`：探测和调用 `unit-test`
- `pkg/testrunner/journey_isolation.go`：`e2e-test` → `test`
- `internal/cmd/task/submit.go`：breaking 任务使用 `UnitGateSequence`
- `internal/cmd/quality_gate.go`：Step 2 `unit-test`，Step 3 `test`，`addFixTask` 映射
- `pkg/forgeconfig/config.go`：`E2eTest` → `Test`，直接替换无兼容层
- `pkg/task/autogen.go`：`E2eTest` → `Test`，`run-e2e-tests` → `run-test`
- `internal/cmd/init.go`：init wizard 文案

**Prompt 模板（3 files）**
- `pkg/prompt/data/gate.md`、`fix-record-missed.md`、`validation-code.md`

**Skill/Command Markdown（5+ files）**
- `commands/fix-bug.md`、`skills/clean-code/SKILL.md`、`skills/gen-test-scripts/rules/run-to-learn.md`
- `skills/init-justfile/SKILL.md`、`skills/run-tests/SKILL.md` + references

**Justfile 模板（7 files）**
- 6 个 `templates/*.just`：`test` → `unit-test`，淘汰 `e2e-*`，增加 `test`/`test-setup`
- 项目根 `justfile`

**测试文件（11+ files）**
- Go unit tests + integration tests 中的 recipe 名和 config 字段断言

**文档（6+ files）**
- CLI docs、ARCHITECTURE.md、quality-gate.md、conventions

**Config（3 files）**
- `.forge/config.yaml`、testdata example/schema

### Out of Scope

- Gate 步骤并行执行（独立优化）
- Gate 结果缓存（独立优化）
- HasRecipe/ResolveScope 探测缓存（独立优化）
- gen-test-scripts 核心逻辑改动（仅 run-to-learn.md 引用更新）
- test Convention 文件改动（`docs/conventions/testing/` 内容不变）
- `auto.test` 的 surface 感知逻辑（当前 `test` recipe 按 surface 生成已足够）
- 历史 lessons/proposals 中的 `e2eTest` 引用（不影响功能）
- `tests/test-generation/` 和 `tests/e2e-pipeline/` 中的集成测试 contract 文件（随 task 执行时更新）
- Go `//go:build e2e` build tags（不相关，保持不变）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 现有项目 justfile 缺少 `unit-test` recipe | M | M | v3.0.0 要求重新运行 `init-justfile` 生成新 justfile |
| `auto.e2eTest` → `auto.test` 需更新现有 config.yaml | H | L | 直接重命名，用户运行 `forge init` 或手动更新即可 |
| `journey_isolation.go` 迁移影响现有 journey 测试 | L | H | `just test` 需支持 `--feature` 参数，模板统一生成 |
| 集成测试 TC_005/TC_015/TC_016 本身已在失败 | H | M | 本次一并修复，测试对齐实际文件内容 |
| 改动面大（~42 files）导致遗漏 | M | M | 按 Tier 分批执行，每批跑测试验证 |

## Success Criteria

- [ ] Breaking 任务 submit 门禁运行 `just unit-test`（非 `just test`），耗时 <30s（Go 项目）
- [ ] All-completed 门禁运行 `just unit-test` + `just test`（完整覆盖，无 `e2e-test` 调用）
- [ ] 无 `unit-test` recipe 时 quality gate 报错提示运行 `init-justfile`（不 fallback）
- [ ] `auto.e2eTest` 完全移除，仅支持 `auto.test`
- [ ] `run-e2e-tests` 任务 key 全部迁移为 `run-test`
- [ ] init-justfile 为 Go/Node/Python/Rust 各生成语言对应的 `unit-test` recipe
- [ ] 所有 Go 测试通过（`go test -race ./...`）
- [ ] 已失败的集成测试 TC_005/TC_015/TC_016 修复

## Next Steps

- Proceed to `/tech-design` to formalize implementation details
