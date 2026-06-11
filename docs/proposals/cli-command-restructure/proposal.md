---
created: 2026-05-22
author: "faner"
status: Draft
---

# Proposal: CLI Command Restructure — Grouped Subcommands & Dead Code Removal

## Problem

Forge CLI 的 17 个顶层命令全部平铺在 `internal/cmd/` 单目录下（45 个 Go 文件），缺乏组织层次。其中 `forge e2e` 组的 6 个子命令（validate-specs, run, setup, verify, compile, discover）已被 justfile 委托模式取代，`forge probe` 的 HTTP 探测功能仅在测试中使用，两者均为死代码。

### Evidence

- `forge e2e` 的实际执行全部委托给 just recipes（`just e2e-run`, `just e2e-setup` 等），CLI 仅做透传
- `forge probe` 仅被 justfile-integration 和 scope-resolution 测试文件引用，无生产调用方
- 当前 45 个 Go 文件平铺在 `internal/cmd/` 下，随着命令增长难以维护

### Urgency

v3.0.0 分支正在进行重大架构变更，是清理死代码和重组目录结构的最佳时机。延迟会导致新功能开发在混乱基础上叠加。

## Proposed Solution

1. **删除** `forge e2e` 组（6 子命令）和 `forge probe`
2. **重组目录** 从平铺文件改为子目录分组
3. `forge test` 保持不变

目标结构：

```
forge-cli/
├── pkg/                        # 共享业务包 (已有)
│   ├── testrunner/             # ← 接收迁入的辅助代码
│   │   ├── testrunner.go       # (已有) RunProjectTests, PrintHookJSON
│   │   ├── test_results.go     # ← 从 internal/cmd 迁入
│   │   └── journey_isolation.go # ← 从 internal/cmd 迁入
│   ├── e2e/                    # ← 删除：随 forge e2e 命令删除变为死代码
│   ├── (contract/, feature/, forgeconfig/, git/, index/, just/,
│   │  lesson/, project/, prompt/, proposal/, task/, template/,
│   │  version/, e2eprobe/)  # 保留不变
│   └── ...
├── internal/cmd/
│   ├── cmd.go                  # rootCmd 注册
│   ├── errors.go               # AIError (留在 cmd，深度耦合命令执行)
│   ├── output.go               # PrintBlock/PrintField (留在 cmd，CLI 输出)
│   ├── task/                   # forge task *
│   ├── test/                   # forge test * (保持不变)
│   ├── feature/                # forge feature *
│   ├── worktree/               # forge worktree *
│   ├── forensic/               # forge forensic *
│   ├── prompt/                 # forge prompt *
│   ├── cleanup.go              # forge cleanup
│   ├── quality_gate.go         # forge quality-gate
│   ├── verify_task_done.go     # forge verify-task-done
│   ├── config.go               # forge config *
│   ├── init.go                 # forge init
│   ├── version.go              # forge version (hidden)
│   ├── claude.go               # forge claude
│   ├── proposal.go             # forge proposal
│   ├── lesson.go               # forge lesson
│   └── docs/                   # 嵌入文档
```

### Innovation Highlights

标准的 cobra 子命令分组模式，无创新。选择此时执行是因为 v3.0.0 的破坏性变更窗口。

## Requirements Analysis

### Key Scenarios

- 用户运行 `forge test promote <journey>` — 行为不变
- 用户运行 `forge test run-journey <name>` — 行为不变
- 用户运行 `forge test verify` — 行为不变
- `forge task`、`forge feature`、`forge worktree`、`forge forensic`、`forge prompt` 行为不变
- 被删除的命令（`forge e2e run/setup/compile/discover/validate-specs`、`forge probe`）不再可用

### Non-Functional Requirements

- **编译速度**: 子包拆分减少单包体积，可能微幅提升增量编译速度

### Constraints & Dependencies

- Go 语言规定每个子目录是独立 package，`test_results.go` 和 `journey_isolation.go` 迁入 `pkg/testrunner/` 后需更新包名和引用
- cobra 的 `init()` 注册模式需要调整，子包通过显式 `Register()` 函数注入父命令

## Alternatives & Industry Benchmarking

### Industry Solutions

CLI 工具的命令组织是成熟领域，cobra/urfavecli 等框架原生支持子命令分组。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零风险 | 45 文件平铺持续恶化，死代码累积 | Rejected: v3.0.0 窗口关闭后更难清理 |
| 仅删除死代码，不重组目录 | — | 风险最低 | 未解决组织问题 | Rejected: 一次性解决比分步更高效 |
| **删除死代码 + 重组目录** | cobra 标准模式 | 彻底清理，结构清晰 | 改动量中等 | **Selected: v3.0.0 窗口值得投入** |

## Feasibility Assessment

### Technical Feasibility

完全可行。Go 子包拆分是标准操作，cobra 命令注册机制天然支持。

### Resource & Timeline

单一开发者 + AI 辅助，预计 1-2 天完成全部改动。

### Dependency Readiness

无外部依赖。所有涉及的代码均在本仓库内。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| `forge e2e` 组已无用 | 代码搜索 | Confirmed: 所有执行委托给 just，CLI 仅透传 |
| `forge probe` 可安全删除 | 引用分析 | Confirmed: 仅测试文件引用，无生产调用方 |
| 子目录拆分不会引起循环导入 | 架构分析 | Confirmed: 各命令组无互相依赖，共享代码留在 cmd 包 |

## Scope

### In Scope

- 删除 `forge e2e` 组（e2e_parent.go, e2e_run.go, e2e_setup.go, e2e_validate_specs.go, e2e_verify.go, e2e_compile.go, e2e_discover.go）
- 删除 `forge probe`（probe.go）
- 将 `test_results.go` 和 `journey_isolation.go` 从 `internal/cmd/` 迁入 `pkg/testrunner/`（与 testrunner 同域）
- `errors.go` 和 `output.go` 保留在 `internal/cmd/`（深度耦合命令执行，无处可归）
- 删除 `pkg/e2e/` 整个包（随 forge e2e 命令删除变为死代码）
- 将 `forge task` 组迁入 `task/` 子目录
- 将 `forge test` 组迁入 `test/` 子目录
- 将 `forge feature` 组迁入 `feature/` 子目录
- 将 `forge worktree` 组迁入 `worktree/` 子目录
- 将 `forge forensic` 组迁入 `forensic/` 子目录
- 将 `forge prompt` 组迁入 `prompt/` 子目录
- 更新 `root.go` 命令注册
- 删除/更新受影响的测试文件和合约（仅 e2e 旧组和 probe 相关）

### Out of Scope

- 历史 proposal/design/task 文档中的引用更新（非功能性）
- `forge config`、`forge init`、`forge claude` 等顶层命令的重组
- 命令行为的任何变更
- 新增命令或功能
- `forge test` 的任何重命名

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Go 子包拆分导致循环导入 | L | H | 各命令子包仅依赖 `cmd` 顶层和 `pkg/*` 不互相依赖 |
| 遗漏引用导致运行时错误 | M | M | 全局搜索 `forge e2e`、`forge probe` 验证无遗漏 |
| 测试覆盖下降（删除 probe/e2e 相关测试） | M | L | 被删命令本身是死代码，测试失去意义 |

## Success Criteria

- [ ] `forge --help` 不显示已删除的命令（e2e 旧组、probe）
- [ ] `forge test promote/run-journey/verify` 行为完全不变
- [ ] 所有保留的命令行为不变（task, feature, worktree, forensic, prompt, config, etc.）
- [ ] `go build ./...` 和 `go test ./...` 通过
- [ ] `internal/cmd/` 下无平铺的命令组文件（仅顶层命令和子目录）

## Next Steps

- Proceed to `/write-prd` to formalize requirements
