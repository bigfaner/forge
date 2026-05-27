---
created: "2026-05-26"
author: faner
status: Draft
---

# Proposal: Justfile Test Recipe Fix

## Problem

`just test` recipe 因依赖未安装的 `go-junit-report`，导致 e2e 测试每次重复执行两遍，且 XML 报告从未成功生成。

### Evidence

- `go-junit-report` 不在系统 PATH 中，`test` recipe 的 pipe 命令因 `pipefail` 失败
- `||` 触发 fallback 重跑一遍不带 `-json` 的 `go test`，每次 `just test` 实际执行两次完整测试
- `results/report.xml` 从未生成（`results/` 中仅有 `unit-raw-output.txt`）
- `go-junit-report` 和 `report.xml` 在项目其他文件中无任何引用

### Urgency

每次 e2e 测试浪费时间翻倍。作为开发者高频使用的命令，应立即修复。

## Proposed Solution

移除 `go-junit-report` 依赖，将 `test` recipe 简化为直接运行 `go test`，保留 `-json` 输出到文件以备调试。`unit-test` recipe 无需修改。

### Innovation Highlights

纯简化操作，无创新点。移除无用依赖层。

## Requirements Analysis

### Key Scenarios

- `just test` — 运行全部 e2e 测试，执行一次，成功/失败结果正确退出
- `just test <journey>` — 运行指定 journey 的 e2e 测试
- `just unit-test` — 运行 forge-cli 单元测试（已有 CGO/race 逻辑，无需修改）

### Constraints & Dependencies

- Windows 环境（Git Bash），GNU sed 可用
- e2e 测试使用 `//go:build e2e` 标签，需要 `-tags=e2e`
- `tests/` 是独立 Go module（`forge-tests`），不在 `forge-cli/` 内

## Alternatives & Industry Benchmarking

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 不改动 | 测试双倍执行时间 | Rejected: 浪费时间 |
| 安装 go-junit-report | CI 常见做法 | 生成 JUnit XML | 增加外部依赖，项目无 CI 需求 | Rejected: 无实际使用场景 |
| **简化 test recipe** | 最小化原则 | 移除依赖，单次执行 | 无 XML 报告 | **Selected: 项目不需要 JUnit 报告** |

## Feasibility Assessment

### Technical Feasibility

完全可行，仅需修改 justfile 中 `test` recipe 的 4 行代码。

### Resource & Timeline

单人 5 分钟完成。

## Scope

### In Scope

- 简化 `test` recipe：移除 `go-junit-report` pipe，直接运行 `go test`
- 验证 `just unit-test` 正确运行
- 验证修复后 `just test` 正确运行

### Out of Scope

- `test-setup`、`ci` 等其他 recipe 的修改
- init-justfile 模板的修改（已有 `test-recipe-unification` 提案覆盖）
- `forge test` 命令组的清理（已有 `remove-forge-test-command` 提案覆盖）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 无 | — | — | — |

## Success Criteria

- [ ] `just unit-test` 在 forge-cli 中执行通过，退出码 0
- [ ] `just test` 执行 e2e 测试仅运行一次，退出码 0
- [ ] `just test` 不依赖 `go-junit-report`
