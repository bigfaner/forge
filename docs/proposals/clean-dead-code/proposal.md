---
created: "2026-05-27"
author: "faner"
status: Draft
---

# Proposal: Forge CLI 死代码清理

## Problem

forge-cli 中存在 6 处明确的死代码：无引用的包、僵尸类型常量、未使用的导出符号和残留的 go.mod 依赖。这些代码增加了阅读者的认知负担，且可能误导新人认为这些代码仍在使用。

### Evidence

通过静态引用分析（grep import + LSP findReferences）逐一确认：

| 项 | 文件 | 引用数 | 证据 |
|----|------|--------|------|
| `internal/docsync/` 包 | 整个目录 | 0 | 只有 `_test.go` 文件，无生产代码，无外部导入 |
| `test.verify-regression` 类型 | `pkg/task/types.go`, `pkg/task/infer.go` | 0 | 已从 autogen 移除（commit 96fc1587），但常量和推断规则残留 |
| `pkg/project` 未用导出（10 个符号） | `pkg/project/types.go`, `markers.go`, `root.go` | 0 | `FindRootInfo`, `FindVCSRoot`, `GetProjectRootFromEnv` 等从未被包外调用 |
| `pkg/version.Name` / `GetName()` | `pkg/version/version.go` | 0 | 仅 `Version` 和 `GetVersion()` 被使用 |
| `pkg/e2eprobe` 内部导出 | `pkg/e2eprobe/e2eprobe.go` | 0 | `ProbeEndpoint`, `ExtractYAMLStringField` 仅包内使用 |
| `mitchellh/hashstructure/v2`, `dustin/go-humanize` | `go.mod` | 0 | 零 import 引用，`go mod tidy` 会自动移除 |

### Urgency

低。纯代码卫生，无功能影响。但与 `forge-cli-clean-code` 提案互补，建议在执行该提案前先完成死代码清理，避免在重复逻辑合并和文件拆分时误操作这些死代码。

## Proposed Solution

分三类操作：(1) 删除死包/死常量；(2) 降级未使用导出为未导出；(3) 清理 go.mod 残留依赖。每步操作后 `go build ./...` 和 `go test ./...` 验证。

### Innovation Highlights

无。标准的死代码清理实践。

## Requirements Analysis

### Key Scenarios

- 开发者浏览 `pkg/project/` 时不再被 `FindVCSRoot` 等未用函数误导
- 新人阅读 `pkg/task/types.go` 时不再看到已不存在的 `verify-regression` 类型
- `go.mod` 不再包含无引用的间接依赖

### Non-Functional Requirements

- **零行为变更**：所有 CLI 命令输入输出不变
- **构建稳定**：`go build ./...` 和 `go test ./...` 全部通过

### Constraints & Dependencies

- Go 1.25 工具链
- 纯重构：不引入新依赖

## Alternatives & Industry Benchmarking

### Industry Solutions

Go 社区标准做法：`go vet` + 静态分析工具（如 `deadcode`）识别 + 手动清理。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 认知噪音持续累积 | Rejected: 与 forge-cli-clean-code 提案精神冲突 |
| 只用 golangci-lint deadcode | Go 工具 | 自动化 | 不覆盖僵尸类型常量和 go.mod 残留 | Rejected: 覆盖面不足 |
| **手动清理 + go mod tidy** | Go 社区实践 | 覆盖全部 6 项，低风险 | 需人工审查 ~15 个文件 | **Selected: 覆盖面完整且风险低** |

## Feasibility Assessment

### Technical Feasibility

所有修改都是删除或重命名（大写→小写），无逻辑变更。Go 工具链会立即捕获任何遗漏的引用。

### Resource & Timeline

约 5 个独立任务，每个任务可独立验证，预计 1-2 小时。

### Dependency Readiness

无外部依赖。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "docsync 是正在使用的包" | 5 Whys | Overturned: 目录中只有 `_test.go`，无生产代码，测试文件引用的包 (`pkg/feature`, `pkg/task`) 也不导入 docsync |
| "verify-regression 类型仍被引用" | 代码搜索 | Overturned: commit 96fc1587 已从 autogen 移除，但常量/注册/推断规则残留。autogen_test.go 明确注释 "verify-regression task no longer exists" |
| "未使用的导出可能被外部调用者使用" | Assumption Flip | Confirmed: forge-cli 是 CLI 工具，不是库，无外部调用者。所有未用导出可安全降级 |

## Scope

### In Scope

**A. 删除死包和僵尸类型**
- 删除 `internal/docsync/` 整个目录（2 个仅测试文件，无生产代码）
- 从 `pkg/task/types.go` 移除 `TypeTestVerifyRegression` 常量及其在 `TaskTypeRegistry`、`ValidTypes`、`SystemTypes` 中的注册
- 从 `pkg/task/infer.go` 移除 `T-test-verify-regression` 和 `T-quick-verify-regression` 推断规则

**B. 降级未使用导出为未导出**
- `pkg/project/`：将 `RootInfo`, `RootType`, `RootTypeUnknown/VCS/Workspace/Project`, `Marker`, `FindRootInfo`, `FindRootInfoFrom`, `FindVCSRoot`, `FindVCSRootFrom`, `FindProjectRootFrom`, `GetProjectRootFromEnv` 降级为未导出
- `pkg/version/`：将 `Name` 和 `GetName()` 降级为未导出
- `pkg/e2eprobe/`：将 `ProbeEndpoint` 和 `ExtractYAMLStringField` 降级为未导出

**C. 清理 go.mod**
- 运行 `go mod tidy` 移除 `mitchellh/hashstructure/v2` 和 `dustin/go-humanize` 残留间接依赖

### Out of Scope

- 遗留代码路径（`e2eprobe` 包、`runTestRegressionLegacy`）—— 被 legacy 质量门使用，属于 `forge-cli-clean-code` 提案范围
- 迁移代码（`Scope` 字段、`e2eTest` 配置迁移）—— 已计划 v3.1.0 移除
- 未使用函数参数（`_ string` in claim.go, state.go）—— API 兼容性保留
- 重复逻辑、超大文件、反模式 —— 属于 `forge-cli-clean-code` 提案范围
- `pkg/project` 内部类型的降级可能影响跨包测试——需在执行时逐个验证

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| `internal/docsync/` 测试文件被其他测试间接引用 | L | L | grep 确认无引用后删除 |
| 降级导出后测试中使用了大写名称 | M | L | 降级时同步修改包内测试，`go test` 立即捕获 |
| `go mod tidy` 移除了实际需要的间接依赖 | L | M | tidy 后运行完整测试套件验证 |
| verify-regression 常量在测试 fixture 中被硬编码引用 | M | L | grep 搜索所有测试文件确认 |

## Success Criteria

- [ ] `internal/docsync/` 目录已删除
- [ ] `TypeTestVerifyRegression` 常量及其注册从 `types.go` 移除
- [ ] 2 条 verify-regression 推断规则从 `infer.go` 移除
- [ ] `pkg/project/` 中 10 个未使用导出符号降级为未导出
- [ ] `pkg/version.Name` 和 `GetName()` 降级为未导出
- [ ] `pkg/e2eprobe.ProbeEndpoint` 和 `ExtractYAMLStringField` 降级为未导出
- [ ] `go.mod` 中 `mitchellh/hashstructure/v2` 和 `dustin/go-humanize` 不再出现
- [ ] `go build ./...` 零错误
- [ ] `go test ./...` 全部通过

## Next Steps

- Proceed to `/quick-tasks` to generate task breakdown
