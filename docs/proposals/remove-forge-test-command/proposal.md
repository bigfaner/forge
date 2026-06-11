---
created: "2026-05-26"
author: "fanhuifeng"
status: Draft
---

# Proposal: Remove `forge test` Command Group

## Problem

`forge test` 命令组的 3 个子命令（`promote`、`run-journey`、`verify`）在实际工作流中已不再使用。这些 CLI 命令只是对 `just test <journey>` 和文件操作的薄封装，skill 层已直接调用 just 完成相同工作，CLI 层完全多余。

### Evidence

- `run-tests` skill 直接调用 `just test`，不经过 `forge test run-journey`
- 标签晋升（@feature → @regression）由 skill 在流程中自行处理，不依赖 `forge test promote`
- contract 验证已在 skill 层通过其他方式完成，`forge test verify` 无调用场景
- `forge quality-gate` 和 `forge feature complete` 直接使用 `testrunner` 包函数，不走 `forge test` 命令

### Urgency

死代码增加维护负担：每次修改测试相关逻辑都需同步考虑 CLI 命令层，而该层无实际用户。清理后减少约 20+ 个文件的维护面。

## Proposed Solution

删除 `forge test` 命令组及其所有引用，保留 `testrunner` 包中被 `quality-gate` 和 `feature complete` 使用的共享函数，删除仅服务于 `forge test` 的 journey 隔离函数。删除整个 `contract` 包。

### Innovation Highlights

纯清理操作，无创新点。标准的死代码删除。

## Requirements Analysis

### Key Scenarios

- 用户执行 `forge test` → 收到 "unknown command" 错误（符合预期）
- skill 层执行测试流程 → 不受影响，已直接调用 just
- `forge quality-gate` 执行 → 不受影响，直接使用 testrunner 包

### Non-Functional Requirements

- 编译通过，无断引用
- 现有非 test 相关的 e2e 测试不受影响

### Constraints & Dependencies

- `testrunner` 包的 `RunProjectTests`、`WriteUnitTestRawOutput`、`WriteRegressionRawOutput`、`Capitalize`、`PrintHookJSON` 被 `quality-gate` 和 `feature complete` 使用，**必须保留**
- `testrunner` 包的 journey 隔离函数仅被 `forge test` 使用，可删除
- `contract` 包仅被 `forge test verify` 使用，可整删

## Alternatives & Industry Benchmarking

### Comparison Table

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| Do nothing | 零风险 | 持续维护死代码，混淆用户 | Rejected: 死代码无保留价值 |
| 废弃标记（deprecated） | 渐进过渡 | 增加废弃逻辑，最终仍需删除 | Rejected: 无用户需要迁移窗口 |
| **直接删除** | 干净利落，减少维护面 | 一次性清理工作量 | **Selected: 无实际用户，无需过渡** |

## Feasibility Assessment

### Technical Feasibility

纯删除操作，无技术风险。`testrunner` 共享函数的保留边界已明确。

### Resource & Timeline

单次清理，约 1 小时工作量。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| forge test 命令仍有用户调用 | Codebase Evidence | Confirmed 废弃：skill 层已绕过 CLI 直接调用 just |
| testrunner 包可整删 | Dependency Analysis | Overturned：quality-gate 和 feature complete 依赖其中 5 个函数 |
| contract 包可整删 | Dependency Analysis | Confirmed：仅 forge test verify 使用 |

## Scope

### In Scope

- 删除 `forge-cli/internal/cmd/test/` 目录（命令实现 + 测试）
- 删除 `forge-cli/pkg/contract/` 目录（仅 verify 使用）
- 删除 `forge-cli/internal/cmd/test_test.go` 和 `forge-cli/internal/cmd/test_verify_test.go`（cmd 根目录下的集成测试；注意 `cmd/test/` 子目录下的同名文件随目录一起删除）
- 从 `testrunner` 包中删除 journey 隔离模块：`journey_isolation.go`（`ResolveJourneyExecutionConfig`、`CreateJourneyWorkDir`、`ExecuteJourneyInIsolation`、`CopyFileToWorkDir` 及相关类型 `JourneyResult`、`ContractFailure`、`JourneyExecutionConfig`）及其测试文件 `journey_isolation_test.go`
- 从 `root.go` 移除 test 命令注册
- 从 `root_test.go` 删除 `testpkg` import（第 10 行）和 `TestRootCmd_TestGroupHasSubcommands` 函数（引用 `testpkg.Cmd`，删除 test 包后编译失败）
- 更新 `quality_gate.go` 中对 `forge test promote` 的提示文本（如有）
- 清理 README、OVERVIEW.md、OVERVIEW.zh.md、ARCHITECTURE.md 中 `forge test` 相关段落
- 清理 conventions 文档中的 `forge test` 引用：`docs/conventions/forge-cli-reference.md`（`forge test` 命令表及 `forge e2e → forge test` 映射）、`docs/conventions/forge-distribution.md`（测试管线中的 `forge test promote` 引用）
- 清理 `docs/profile-authoring.md` 中过时的 `forge testing` 引用（旧命令名，双重过时）
- 清理 skill 文档中的 `forge test` 引用：`run-tests`（SKILL.md）、`consolidate-specs`（SKILL.md）、`gen-contracts`（`rules/journey-contract-model.md`）、`gen-journeys`（`rules/journey-contract-model.md`）
- 清理 command 文档中的 `forge test` 引用：`plugins/forge/commands/run-tasks.md`（第 119 行建议执行 `forge test promote`）
- 删除 `tests/command-regression/` 整个目录（`removed_commands_test.go`、`main_test.go`、`contracts/step-1-removed-test-commands.md`）——该目录专门测试已移除的 `forge test` 子命令（detect/get/interfaces/framework），与本次删除的命令组同属死代码

### Out of Scope

- journey/contract/promote 测试模型本身不变
- `testrunner` 包的共享函数（`RunProjectTests` 等）保留不动
- feature 文档（`docs/features/`）中的历史记录不清理——它们是设计决策的历史证据

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 遗漏某个引用导致编译失败 | L | L | `go build ./...` 和 `go test ./...` 验证 |
| skill 中仍有对 forge test 的隐式调用 | L | M | 全文搜索 `forge test` 确认零残留 |
| 未来需要 CLI 层的测试编排 | L | L | 届时重新添加，YAGNI |

## Success Criteria

- [ ] `forge test` 执行返回 "unknown command" 错误
- [ ] `go build ./...` 编译通过，零错误
- [ ] `go test ./...` 全部通过
- [ ] 全文搜索 `forge test promote`、`forge test run-journey`、`forge test verify` 返回零结果（排除 feature 历史文档）
- [ ] `forge quality-gate` 和 `forge feature complete` 命令执行返回 exit code 0

## Next Steps

- 此为纯清理任务，无需 PRD 或 tech design，可直接进入任务拆分和执行
