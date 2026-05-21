---
created: 2026-05-21
author: fanhuifeng
status: Approved
---

# Proposal: Decouple Test Tasks from Languages

## Problem

Test pipeline task generation is coupled to language detection, causing silent task skipping when detection fails for monorepo/subdirectory module projects.

### Evidence

Forge 自身就是一个典型案例：`go.mod` 位于 `forge-cli/` 子目录而非项目根目录。`DetectLanguages` 只检查根目录，检测失败返回空 slice → `interfaces` 为空 → `GetBreakdownTestTasks` 守卫 `len(interfaces) == 0` 直接返回 nil。结果：`forge task index` 生成的任务数从 32 降到 24，缺少 11 个测试管道任务，且无任何错误或警告。

详细分析记录在 `docs/lessons/gotcha-test-pipeline-no-languages.md`。

### Urgency

测试管道是 Forge 质量保证的核心。当前状态意味着任何非标准项目结构（monorepo、子目录模块）的测试管道任务都会静默跳过，且用户无法感知。每次 `forge task index` 都会受到影响。

## Proposed Solution

将测试任务从语言绑定解耦为纯接口类型（场景）绑定。`interfaces` 配置项成为唯一的接口类型来源，完全由用户在 `config.yaml` 中声明，移除所有语言检测逻辑。

**改动前后对比**：

Task key 变化（breakdown 模式，2 个 interface）：

| 改动前（language×type） | 改动后（type-only） |
|---|---|
| `gen-test-scripts-go-api` | `gen-test-scripts-api` |
| `gen-test-scripts-go-cli` | `gen-test-scripts-cli` |
| `run-e2e-tests-go` | `run-e2e-tests` |
| `graduate-tests-go` | `graduate-tests` |

Task key 变化（quick 模式，2 个 interface）：

| 改动前 | 改动后 |
|---|---|
| `quick-test-cases-go` | `quick-test-cases` |
| `quick-gen-and-run-go-api` | `quick-gen-and-run-api` |
| `quick-gen-and-run-go-cli` | `quick-gen-and-run-cli` |
| `quick-graduate-go` | `quick-graduate` |

### Innovation Highlights

这不是创新而是简化。核心洞察：测试任务的粒度应该由接口类型（CLI、API、WebUI、TUI）决定，而非由实现语言决定。一个 Go CLI 项目和一个 Rust CLI 项目的 CLI 测试任务是同构的——只是执行方式不同，而执行方式由 `just` 统一管理。

## Requirements Analysis

### Key Scenarios

1. **用户未配置 interfaces**：`forge task index` 正常生成业务任务和 stage-gates，跳过测试管道任务，输出明确警告提示用户配置
2. **用户配置了 interfaces**：正常生成所有测试管道任务，按接口类型组织
3. **多接口类型项目**：每个接口类型生成独立的 gen-scripts 任务，共享 run/graduate/verify 任务
4. **现有项目迁移**：已有 `config.yaml` 中的 `languages` 字段被忽略（向后兼容），用户需改为配置 `interfaces`

### Non-Functional Requirements

- 向后兼容：`languages` 字段保留在 Config struct 中但不影响测试管道，避免解析旧 config.yaml 报错
- 性能：移除文件检测逻辑后，`forge task index` 速度略微提升

### Constraints & Dependencies

- `plugins/forge/` 目录零引用这些语言相关函数，改动范围限于 `forge-cli/`
- 现有 `docs/features/*/tasks/index.json` 中的旧格式 task key 不受影响（已生成的任务不会自动重命名）

## Alternatives & Industry Benchmarking

### Industry Solutions

大多数测试框架（Jest、pytest、Go testing）不关心"语言"——它们通过测试发现机制（文件匹配、构建标签）自动定位测试。测试管道按接口类型组织而非按语言组织，与行业实践一致。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 测试管道在 monorepo 下静默失效，用户体验差 | Rejected: 核心功能不可用 |
| 增强语言检测 | — | 自动化 | 检测深度边界难以确定（1层？2层？递归？），误检风险 | Rejected: 增加复杂度，不解决根本问题 |
| 检测 + 告警兜底 | — | 兼顾自动化和安全性 | 保留了两套路径，代码复杂 | Rejected: 用户选择了纯配置驱动 |
| **纯接口绑定 + 配置驱动** | 行业实践 | 简单、确定性高、无检测边界问题 | 用户需手动配置 interfaces | **Selected: 用户选择，且符合简化理念** |

## Feasibility Assessment

### Technical Feasibility

完全可行。改动范围明确：
- `forge-cli/pkg/task/testgen.go` — 重构任务生成逻辑
- `forge-cli/pkg/task/build.go` — 移除 languages 参数
- `forge-cli/pkg/task/infer.go` — 简化 ID 模式匹配
- `forge-cli/pkg/forgeconfig/detect.go` — 移除检测函数
- 相关测试文件更新

### Resource & Timeline

约 6-8 个 coding tasks，包含测试更新。可在单次迭代内完成。

### Dependency Readiness

无外部依赖。所有代码在 `forge-cli/` 内部。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 测试任务需要按语言区分 | XY Detection | 真正的区分维度是接口类型（CLI/API/WebUI），不是语言。语言只影响执行方式，而执行由 `just` 统一管理 |
| 自动检测比手动配置更好 | Assumption Flip | 自动检测在边界场景下静默失效，比明确的"未配置"警告更危险 |
| 多语言项目需要独立测试管道 | Occam's Razor | 多语言项目的测试执行由 `just` 统一调度，不需要 Forge 为每种语言生成独立的 run/graduate 任务 |

## Scope

### In Scope

- 移除 `testgen.go` 中的 language 循环，改为 interface-only 任务生成
- 移除 `profileSuffix`、`suffixLetter` 函数
- 简化 `infer.go` 的 ID 模式匹配（移除 `profileSuffixedID`）
- 简化 `ReadInterfaces` 为纯配置读取
- 移除 `DetectLanguages`、`ReadLanguages`、`UnionLanguageInterfaces`、`defaultInterfaces`
- 移除 `languageCapabilities`、`KnownLanguages`、`IsKnownLanguage`
- 移除 `detectPytest`、`fileExists`、`dirExists` 辅助函数
- `BuildIndex` 添加 `interfaces` 未配置时的明确警告
- 更新所有相关测试
- 在 `docs/lessons/` 中更新 gotcha 文档记录新模型

### Out of Scope

- `forge init` 添加 `interfaces` 交互式配置（可后续迭代）
- `forge config get interfaces` 命令支持（可后续迭代）
- 已有 feature 的 index.json 迁移（自然淘汰）
- `Languages` 字段从 Config struct 中移除（保留以向后兼容旧 config.yaml）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 用户不知道需要配置 interfaces | M | M | BuildIndex 中输出醒目警告，包含配置示例 |
| 现有测试任务 key 变化导致已生成任务的 index.json 不匹配 | L | L | 新旧 key 格式不同，InferType 更新后旧任务 ID 仍可识别 |
| 遗漏某些引用 languages 的代码路径 | L | M | 全局搜索确认 `forge-cli/` 内仅 `build.go` 一个生产调用者 |

## Success Criteria

- [ ] `forge task index` 在 config.yaml 配置 `interfaces: [api, cli]` 时正确生成所有测试管道任务
- [ ] `forge task index` 在未配置 interfaces 时跳过测试管道任务并输出明确警告
- [ ] `DetectLanguages`、`ReadLanguages`、`UnionLanguageInterfaces` 函数从 `detect.go` 中移除
- [ ] 测试任务 key 不包含语言名（如 `gen-test-scripts-api` 而非 `gen-test-scripts-go-api`）
- [ ] `profileSuffix` 和 `suffixLetter` 函数从 `testgen.go` 中移除
- [ ] `profileSuffixedID` 函数从 `infer.go` 中移除
- [ ] 所有现有测试通过（更新后的测试）
- [ ] `forge-cli/` 中无 `Languages` 相关的生产代码引用（Config struct 字段定义除外）

## Next Steps

- Proceed to `/quick-tasks` for task breakdown
