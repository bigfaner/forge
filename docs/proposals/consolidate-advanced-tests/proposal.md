---
created: 2026-06-06
author: faner
status: Draft
intent: refactor
---

# Proposal: Consolidate Advanced Tests

## Problem

CLI 功能测试（高级测试）分散在两个目录：顶层 `tests/`（规范正确位置）和 `forge-cli/tests/`（agent 开发时未遵循 forge 规范的产物）。`just test` 只运行 `tests/`，`forge-cli/tests/` 中的 27 个测试文件从未被测试管道执行。

### Evidence

- `forge-cli/tests/` 含 7 个 journey、27 个测试文件、独立 testkit、独立 config.yaml
- 与 `tests/` 存在 2 个重叠 journey（`task-lifecycle`、`task-type-system`），文件名不同但主题相同
- `just test` 执行 `cd tests && go test -tags=cli_functional`，完全不含 `forge-cli/tests/`
- 历史上曾创建 `forge-cli-test-spec-alignment` 功能（status: completed）但 10 个任务全部仍为 pending，从未执行

### Urgency

v3.0.0 分支正在重构测试基础设施（test-recipe-unification 提案已进入实现阶段）。现在整合可避免后续返工，且确保所有高级测试都能被标准管道覆盖。

## Proposed Solution

将 `forge-cli/tests/` 全部内容迁移至顶层 `tests/`：

1. **合并 testkit**：将 `forge-cli/tests/testkit` 中顶层缺失的 3 个函数（`RunCLIExitCode`、`ProjectRoot`、`ReadProjectFile`）加入 `tests/testkit/`
2. **迁移 5 个独有 journey**：`error-handling`、`forge-commands`、`justfile-integration`、`scope-resolution`、`skill-ops` 整体移入 `tests/`
3. **合并 2 个重叠 journey**：将 `forge-cli/tests/` 中的测试文件加入 `tests/` 对应目录，共存
4. **统一初始化模式**：所有迁移的 `main_test.go` 改用顶层 testkit 的 `ForgeBinary` init 模式
5. **更新活跃引用**：修正 conventions 和 skill rules 中指向 `forge-cli/tests/` 的路径
6. **删除 `forge-cli/tests/`**

### Innovation Highlights

非创新性方案。标准的目录整合迁移——将散落的测试文件归位到规范位置。核心价值是消除"死测试"（写了但永远不跑的测试），确保 `just test` 一条命令覆盖所有高级测试。

## Requirements Analysis

### Key Scenarios

- **迁移后全量测试通过**：`just test` 执行时，所有 17 个 journey（原有 10 + 迁移 5 + 合并 2）的测试全部通过
- **重叠 journey 合并无冲突**：`task-lifecycle` 和 `task-type-system` 各自的测试文件在同一目录下共存，无命名冲突（文件名均不同）
- **testkit API 向后兼容**：原有 `tests/` 的测试代码无需修改，仅新增函数

### Non-Functional Requirements

- **性能**：迁移后 `just test` 耗时增幅 <20%（从 10 journey 增至 17 journey，符合线性预期）
- **无功能变更**：迁移不修改测试逻辑、断言内容或 contract 文件

### Constraints & Dependencies

- 目标目录 `tests/` 是独立 Go module（`module forge-tests`），迁移的测试文件必须使用该模块的 import 路径
- 所有测试文件使用 `//go:build cli_functional` build tag
- `tests/testkit/` 的 `ForgeBinary` 通过 `init()` 自建二进制，迁移后的 `main_test.go` 需适配此模式

## Alternatives & Industry Benchmarking

### Industry Solutions

测试文件位置整合是 monorepo 管理的基本实践。Go 社区惯例是将同一类型的测试放在同一模块下，通过 build tag 区分类型（如 `//go:build cli_functional`），而非分散在不同模块中。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 27 个测试文件永远不跑；违反 forge 规范 | Rejected: 死测试等于没有测试 |
| 复用 forge-cli-test-spec-alignment 旧功能 | 历史功能 | 已有任务定义 | 范围过时（含已不存在的 tests/e2e/）；10 个任务全 pending 需重新设计 | Rejected: 过时范围，不如新建 |
| **新建精简提案** | 标准 monorepo 整合 | 范围精确；风险低；聚焦单一目标 | 需写提案 + 执行 | **Selected: 最小可行方案** |

## Feasibility Assessment

### Technical Feasibility

完全可行。改动类型：
- 文件移动（7 个 journey 目录）
- import 路径替换（`forge-cli/tests/testkit` → `forge-tests/testkit`）
- `main_test.go` 重写（统一初始化模式）
- testkit 新增 3 个函数
- 删除 `forge-cli/tests/` 目录

无外部依赖，无框架变更。

### Resource & Timeline

~35 个文件迁移 + 3 个 testkit 函数新增 + ~5 个引用更新。预估 3-4 个 coding task，1 人执行约 1-2 天。

### Dependency Readiness

无外部依赖。`tests/` 模块已稳定运行。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| forge-cli/tests/ 的测试与 tests/ 有大量重复需要去重 | Occam's Razor | Overturned: 仅 2 个 journey 名称重叠，且内部测试文件完全不同（不同文件名、不同测试函数），无需去重，直接合并共存即可 |
| 迁移需要修改测试逻辑以适配顶层 testkit | 5 Whys | Refined: 测试逻辑不变，仅需改 import 路径和 main_test.go 初始化方式。底层 testkit 需补充 3 个缺失函数确保 API 兼容 |
| forge-cli-test-spec-alignment 功能已完成迁移 | Evidence check | Overturned: 功能标记 completed 但 10 个任务全部 pending，迁移从未执行 |

## Scope

### In Scope

**testkit 合并（1 file 新增 + 2 files 修改）**
- `tests/testkit/helpers.go`：新增 `RunCLIExitCode`、`ProjectRoot`、`ReadProjectFile` 三个函数
- `tests/testkit/forge_binary.go`：确保 `ForgeBinary` 路径在迁移后仍正确解析 forge-cli 目录

**5 个独有 journey 迁移（~14 test files）**
- `error-handling/`（1 test file）→ `tests/error-handling/`
- `forge-commands/`（4 test files）→ `tests/forge-commands/`
- `justfile-integration/`（4 test files）→ `tests/justfile-integration/`
- `scope-resolution/`（1 test file）→ `tests/scope-resolution/`
- `skill-ops/`（4 test files）→ `tests/skill-ops/`

**2 个重叠 journey 合并（~5 test files）**
- `forge-cli/tests/task-lifecycle/` 3 个测试文件 + contracts → `tests/task-lifecycle/`
- `forge-cli/tests/task-type-system/` 2 个测试文件 + contracts → `tests/task-type-system/`

**初始化模式统一（7 个 main_test.go）**
- 所有迁移的 `main_test.go` 统一为顶层 testkit 的 `ForgeBinary` init 模式

**活跃引用更新**
- `docs/conventions/forge-distribution.md`：移除 `forge-cli/tests/` 路径引用
- `plugins/forge/skills/run-tests/rules/test-isolation.md`：更新路径示例
- `tests/test-suite-health/contracts/step-1-test-suite-health.md`：更新路径

**清理**
- 删除 `forge-cli/tests/` 整个目录

### Out of Scope

- 测试逻辑或断言内容的修改
- 新增测试用例
- 历史文档更新（`docs/features/`、`docs/lessons/`、`docs/proposals/` 中的历史引用）
- test-recipe-unification 提案的执行
- `go.mod` 依赖变更
- CI/CD 流水线变更（`just test` 已正确指向 `tests/`）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 迁移后测试编译失败（import 路径或 init 模式不匹配） | M | M | 逐 journey 迁移，每迁移一个即运行 `just test <journey>` 验证 |
| testkit 函数签名差异导致迁移的测试调用失败 | M | L | 逐一对比 `forge-cli/tests/testkit` 与 `tests/testkit` 的函数签名，确保行为一致 |
| contracts 目录内容与顶层 tests/ 的 contracts 有冲突 | L | L | contracts 是 markdown 文档，可直接共存；如有内容冲突再合并 |
| forge-cli/tests/ 中部分测试依赖 forge-cli 内部包路径 | M | M | 迁移前检查每个测试的 import 列表，识别非 testkit 的 forge-cli 内部引用 |

## Success Criteria

- [ ] `forge-cli/tests/` 目录完全删除，不存在于工作树中
- [ ] `tests/` 下包含所有 17 个 journey（原 10 + 迁移 5 + 合并 2 新增测试文件）
- [ ] `just test` 执行所有 17 个 journey 的测试，无编译错误
- [ ] `tests/testkit/` 包含 `RunCLIExitCode`、`ProjectRoot`、`ReadProjectFile` 三个新函数
- [ ] 所有迁移的测试文件 import `testkit "forge-tests/testkit"`，不含 `"forge-cli/tests/testkit"` 引用
- [ ] `grep -r 'forge-cli/tests' tests/` 返回空（源码中无残留路径引用）
- [ ] `docs/conventions/forge-distribution.md` 和 `plugins/forge/skills/run-tests/rules/test-isolation.md` 中无 `forge-cli/tests/` 引用

## Next Steps

- Proceed to `/quick` or `/tech-design` to formalize implementation details
