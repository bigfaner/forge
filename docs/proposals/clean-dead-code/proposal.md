---
created: "2026-05-27"
author: "faner"
status: Draft
---

# Proposal: Forge CLI 死代码清理

## Problem

forge-cli 中存在死包、僵尸推断规则和空壳文件。这些代码增加了阅读者的认知负担，且可能误导新人认为这些代码仍在使用。

### Evidence

通过静态引用分析（grep import + LSP findReferences）逐一确认，并于 2026-05-29 全面复核：

| 项 | 文件 | 引用数 | 证据 | 状态 |
|----|------|--------|------|------|
| `internal/docsync/` 残留测试 | `internal/docsync/*.go` | 0 | 生产代码已删除，但 2 个测试文件（~1100 行）残留，无生产代码可供测试 | 待清理 |
| `test.verify-regression` 残留 | `pkg/task/infer.go`, `pkg/task/autoconfig_test.go`, `internal/cmd/quality_gate_test.go` | 0 | 常量已清除，但 `infer.go` 推断规则仍返回 `"test.verify-regression"`；测试文件有陈旧字符串 | 待清理 |
| 空壳文件 | `internal/cmd/errors.go`, `internal/cmd/worktree/worktree.go` | 0 | 仅包声明和注释，零功能代码 | 待清理 |
| ~~`mitchellh/hashstructure/v2`, `dustin/go-humanize`~~ | `go.mod` | — | `go mod graph` 确认为 `charmbracelet/huh` 和 `bubbles` 的合法传递依赖，`go mod tidy` 不会移除 | **假设错误，不移除** |
| ~~`pkg/e2eprobe` 包~~ | — | — | — | **已清理** |

### 排查结论（2026-05-29）

- **无新死包**：全部 27 个包均有外部引用
- **go.mod 无残留依赖**：`hashstructure/v2` 和 `go-humanize` 为 `charmbracelet/huh` 和 `bubbles` 的合法传递依赖，`go mod tidy` 不会移除
- **101 个未使用导出符号**：跨 10 个包，不纳入本次清理范围（降级为未导出属于 `forge-cli-clean-code` 提案）

### Urgency

低。纯代码卫生，无功能影响。但与 `forge-cli-clean-code` 提案互补，建议在执行该提案前先完成死代码清理，避免在重复逻辑合并和文件拆分时误操作这些死代码。

## Proposed Solution

两类操作：(1) 删除死包和空壳文件；(2) 清理 verify-regression 残留推断规则和陈旧测试引用。每步操作后 `go build ./...` 和 `go test ./...` 验证。

### Innovation Highlights

无。标准的死代码清理实践。

## Requirements Analysis

### Key Scenarios

- 新人阅读 `pkg/task/infer.go` 时不再看到已不存在的 `verify-regression` 推断规则
- `internal/docsync/` 不再作为孤立测试目录存在
- 空壳文件不再干扰包浏览和代码搜索

### Non-Functional Requirements

- **零行为变更**：所有 CLI 命令输入输出不变
- **构建稳定**：`go build ./...` 和 `go test ./...` 全部通过

### Constraints & Dependencies

- Go 1.25 工具链
- 纯删除：不引入新依赖，不修改任何导出符号

## Alternatives & Industry Benchmarking

### Industry Solutions

Go 社区标准做法：`go vet` + 静态分析工具（如 `deadcode`）识别 + 手动清理。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 认知噪音持续累积 | Rejected: 与 forge-cli-clean-code 提案精神冲突 |
| 只用 golangci-lint deadcode | Go 工具 | 自动化 | 不覆盖空壳文件 | Rejected: 覆盖面不足 |
| **手动清理** | Go 社区实践 | 覆盖全部 3 项，低风险 | 需人工审查 ~8 个文件 | **Selected: 覆盖面完整且风险低** |

## Feasibility Assessment

### Technical Feasibility

所有修改都是文件/目录删除或代码行删除，无逻辑变更。Go 工具链会立即捕获任何遗漏的引用。

### Resource & Timeline

3 组操作（A: 删除死包和空壳文件、B: 清理 verify-regression 残留、C: go mod tidy），每组可独立验证，预计 30 分钟。

### Dependency Readiness

无外部依赖。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "docsync 是正在使用的包" | 5 Whys | Overturned: 目录中只有 `_test.go`，无生产代码，测试文件引用的包也不导入 docsync |
| "verify-regression 类型仍被引用" | 代码搜索 | Overturned: 常量已从 types.go 清除，但 infer.go 推断规则和测试陈旧字符串残留 |
| "空壳文件可能有隐藏用途" | 代码搜索 | Confirmed: `errors.go` 和 `worktree.go` 零功能代码，同目录其他文件独立完成所有工作 |
| "hashstructure/go-humanize 是传递依赖所需" | `go mod graph` | **Confirmed**: 两者均为 `charmbracelet/huh` 和 `bubbles` 的合法传递依赖，`go mod tidy` 不会移除 |

## Scope

### In Scope

**A. 删除死包和空壳文件**
- 删除 `internal/docsync/` 整个目录（2 个仅测试文件，~1100 行，无生产代码可供测试）
- 删除 `internal/cmd/errors.go`（5 行，仅包声明 + 注释重定向到 base 子包）
- 删除 `internal/cmd/worktree/worktree.go`（16 行，仅包文档，零功能代码）

**B. 清理 verify-regression 残留**
- ~~从 `pkg/task/types.go` 移除 `TypeTestVerifyRegression` 常量及其注册~~ — 已完成
- 从 `pkg/task/infer.go` 移除返回 `"test.verify-regression"` 的推断规则（`T-test-verify-regression` 和 `T-quick-verify-regression` 分支）
- 从 `pkg/task/autoconfig_test.go` 清理 `T-test-verify-regression` 相关陈旧错误消息（行 264, 280）
- 从 `internal/cmd/quality_gate_test.go` 清理 `T-quick-verify-regression` 测试数据（行 1330）

### Out of Scope

- 未使用导出符号降级（101 个）—— 属于 `forge-cli-clean-code` 提案范围
- 遗留代码路径（`runTestRegressionLegacy`）—— 被 legacy 质量门使用，属于 `forge-cli-clean-code` 提案范围
- 迁移代码（`Scope` 字段、`e2eTest` 配置迁移）—— 已计划 v3.1.0 移除
- 未使用函数参数（`_ string` in claim.go, state.go）—— API 兼容性保留
- 重复逻辑、超大文件、反模式 —— 属于 `forge-cli-clean-code` 提案范围

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| `internal/docsync/` 测试文件被其他测试间接引用 | L | L | grep 确认无引用后删除 |
| 空壳文件被其他 `.go` 文件通过 `//go:embed` 或 build tag 引用 | L | L | 删除前 grep 确认无引用 |
| verify-regression 字符串在测试 fixture 中被硬编码引用 | M | L | grep 搜索所有测试文件确认 |

## Success Criteria

- [ ] `internal/docsync/` 目录已删除
- [ ] `internal/cmd/errors.go` 已删除
- [ ] `internal/cmd/worktree/worktree.go` 已删除
- [x] ~~`TypeTestVerifyRegression` 常量及其注册从 `types.go` 移除~~ — 已完成
- [ ] 2 条 verify-regression 推断规则从 `infer.go` 移除
- [ ] `autoconfig_test.go` 和 `quality_gate_test.go` 中 verify-regression 陈旧引用已清理
- [x] ~~`pkg/e2eprobe` 包已移除~~ — 已完成
- ~~[ ] `go.mod` 中 `mitchellh/hashstructure/v2` 和 `dustin/go-humanize` 不再出现~~ — 假设错误，两者为合法传递依赖
- [ ] `go build ./...` 零错误
- [ ] `go test ./...` 全部通过

## Next Steps

- Proceed to `/quick-tasks` to generate task breakdown
