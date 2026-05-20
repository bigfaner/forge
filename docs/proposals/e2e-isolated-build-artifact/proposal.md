---
created: 2026-05-20
author: "fanhuifeng"
status: Draft
---

# Proposal: E2E 测试使用隔离构建产物

## Problem

CLI 集成测试在 4 个不同位置使用了 3 种不同的二进制解析策略，其中 3 种依赖共享产物（系统 PATH 或 `forge-cli/bin/forge`），导致开发者在分支间切换时测试二进制与当前源码不一致，产生误导性测试结果。

### Evidence

- `forge-cli/tests/e2e/helpers_test.go` 和 `testkit/helpers.go` 使用 `exec.Command("forge", ...)`，依赖系统 PATH 上的 `~/.forge/bin/forge`
- `tests/e2e/justfile-canonical-e2e/helpers_test.go` 使用共享路径 `forge-cli/bin/forge`
- 部分特性测试使用 walk-up 模式回退到 `forge-cli/bin/forge`
- 仅 `tests/e2e/main_test.go` 的 TestMain 模式做到了真正的隔离构建
- `docs/conventions/testing-isolation.md` 的 TEST-isolation-004 已定义约定但未被全面执行

### Urgency

每次分支切换后运行测试都可能得到不可靠的结果，开发者需要手动运行 `e2e-setup` 来确保二进制一致。这降低了测试的可信度，且违反了已建立的隔离约定。

## Proposed Solution

将所有 E2E 测试模块统一为 TestMain 自动构建模式：每个测试模块在 TestMain 中从当前源码编译独立二进制到临时目录，测试结束后自动清理。消除对系统 PATH 和共享构建产物的依赖。

### Innovation Highlights

这是标准实践——Go 生态中 `TestMain` 构建测试二进制是常见模式（如 `kubectl`、`helm` 的 e2e 测试）。本方案的创新点不在于技术，而在于将已有的 `tests/e2e/main_test.go` 验证模式统一推广到所有测试位置，并简化 justfile 流程使 `e2e-setup` 从必需步骤变为可选缓存优化。

## Requirements Analysis

### Key Scenarios

- **开发者切换分支后直接运行测试**：TestMain 自动从当前源码构建，无需手动 e2e-setup
- **CI 中运行测试**：同样自动构建，无需预构建步骤
- **并行运行多组测试**：每个 Go test package 的 TestMain 独立构建，互不干扰
- **构建失败时**：TestMain 输出构建错误并立即退出，不会运行任何测试

### Non-Functional Requirements

- **构建时间**：Go 编译 forge 二进制约 2-5 秒，在 E2E 测试（通常数分钟）中可接受
- **磁盘空间**：临时目录在测试结束后自动清理，不占空间

### Constraints & Dependencies

- `tests/e2e/` 和 `forge-cli/tests/e2e/` 是不同的 Go 模块，无法共享 testkit 代码
- `tests/e2e/justfile-canonical-e2e/` 也是独立 Go 模块
- 各模块需要各自实现 TestMain 构建逻辑（可接受少量代码重复）

## Alternatives & Industry Benchmarking

### Industry Solutions

Go CLI 项目的 e2e 测试通常使用 TestMain 构建（`kubectl`、`helm`、`terraform`），或使用 Makefile 预构建到固定位置（较不推荐）。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零改动 | 分支切换测试不可靠，违反 TEST-isolation-004 | Rejected: 问题持续存在 |
| 共享 testkit 模块 | 本方案备选 | DRY | 跨 Go 模块依赖管理复杂，引入新模块 | Rejected: 复杂度不值得 |
| **Per-Module TestMain** | `tests/e2e/main_test.go` 已验证 | 零配置、自包含、可靠 | 少量代码重复 | **Selected: 简单有效** |

## Feasibility Assessment

### Technical Feasibility

完全可行。`tests/e2e/main_test.go` 已验证此模式。各模块只需在 TestMain 中添加 `go build` 调用，并将 `exec.Command("forge", ...)` 改为使用构建产物路径。

### Resource & Timeline

改动量小，预计 3-5 个任务可完成。每个模块的改动模式一致。

### Dependency Readiness

无外部依赖。所有需要修改的代码均在项目内。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 共享 testkit 能减少维护成本 | Occam's Razor | 跨模块 testkit 的维护成本 > 少量重复的 TestMain 代码。3-4 个模块的重复远比管理一个跨模块依赖简单 |
| e2e-setup 是测试前置必需步骤 | Assumption Flip | 如果测试自动构建，e2e-setup 变为可选的缓存优化步骤。大多数场景下开发者不再需要手动运行它 |
| 代码重复是不可接受的 | Stress Test | 3-4 个 TestMain 函数，每个约 15 行，重复是可接受的。统一约定比统一代码更重要 |

## Scope

### In Scope

- `forge-cli/tests/e2e/` — helpers_test.go 和 testkit/helpers.go 从系统 PATH 改为 TestMain 自动构建
- `tests/e2e/justfile-canonical-e2e/` — 从共享路径改为 TestMain 自动构建
- `tests/e2e/` 下的特性测试 — 审查所有测试文件，确保统一使用 TestMain 构建的二进制，消除 walk-up 回退路径
- `justfile` — 简化 e2e-setup（构建步骤变为可选/缓存优化），不再作为测试前置必需步骤
- `docs/conventions/testing-isolation.md` — 更新 TEST-isolation-004 scope 字段，覆盖所有测试位置

### Out of Scope

- 单元测试（不涉及 CLI 二进制调用）
- CI pipeline 修改（自动构建在 CI 同样适用，无需特殊处理）
- 构建性能优化（2-5 秒构建时间在 E2E 测试中可接受）
- 共享 testkit 模块提取（当前不做，未来可按需引入）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 构建时间拖慢 CI | L | M | 构建仅 2-5 秒，E2E 测试通常数分钟；CI 可缓存构建产物 |
| 遗漏某个使用系统 PATH 的测试文件 | M | M | 用 grep 搜索 `exec.Command("forge"` 确保全覆盖 |
| go.build 在某些环境下失败 | L | H | TestMain 输出详细构建错误信息，便于诊断 |

## Success Criteria

- [ ] `grep -r 'exec.Command("forge"' forge-cli/tests/ tests/e2e/` 返回零结果（无测试依赖系统 PATH）
- [ ] 所有 E2E 测试模块的 TestMain 从源码构建独立二进制到临时目录
- [ ] 分支切换后无需运行 `e2e-setup` 即可直接运行 E2E 测试
- [ ] `justfile` 中 `e2e-setup` 的构建步骤标记为可选
- [ ] TEST-isolation-004 的 scope 覆盖所有测试位置

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
