---
title: "Forge CLI 测试 Journey-Contract 组织模式对齐"
status: Draft
date: 2026-05-20
---

# Forge CLI 测试 Journey-Contract 组织模式对齐

## Problem

Forge 的测试管道已采纳 Journey-Driven Test Model（见 `docs/proposals/contract-journey-test-model/`），规范要求测试按用户真实工作流（Journey）组织，每个 Journey 下有 `contracts/` 目录和测试文件。但当前 forge CLI 的 E2E 测试仍按 **feature/command 扁平组织**，未对齐该模型。

### 当前结构（按 feature/command 组织）

```
tests/e2e/                                    (Go module: e2e-tests, package: e2e)
  main_test.go                                (TestMain: binary alias)
  forge_binary.go                             (init: build binary, shared via import)
  helpers_test.go                             (shared helpers: parseBlock, hasField, etc.)
  feature_set_command_cli_test.go             ← 按 command 名
  task_lifecycle_hardening_cli_test.go        ← 按 feature 名
  task_record_immutability_cli_test.go        ← 按 feature 名
  quality_gate_fix_task_loop_breaker_cli_test.go ← 按 feature 名
  quick_test_slim_cli_test.go                 ← 按 feature 名
  test_scripts_per_type_cli_test.go           ← 按 feature 名
  e2e_test_quality_cleanup_cli_test.go        ← 元测试
  gen_journeys_skill_cli_test.go              ← 按 skill 名
  simplify_e2e_tests_cli_test.go              ← 按 feature 名
  features/
    cli-list-reverse-chronological/           (package: e2eclilistr)
    fix-task-claim-priority/                  (package: e2efixclaim)
    proposal-status-lifecycle/                (package: e2epsl)
    task-type-refinement/                     (package: e2etasktype)
    test-knowledge-convention-driven/         (package: e2etestconv, 5 test files)
  justfile-canonical-e2e/                     (separate sub-package)
  .graduated/                                 (graduated test markers)
```

### 目标结构（按 Journey 组织）

```
tests/e2e/
  testkit/                                    (共享基础设施包，供所有 Journey import)
    forge_binary.go                           (init: build binary → ForgeBinary)
    helpers.go                                (parseBlock, hasField, withRetry 等)
  task-lifecycle/                             (Journey: 任务生命周期)
    contracts/
      step-1-task-claim.md
      step-2-task-submit.md
    task_lifecycle_test.go                    (合并 task_lifecycle_hardening + fix-task-claim-priority)
    task_record_test.go                       (来自 task_record_immutability)
  quality-gate/                               (Journey: 质量门禁)
    contracts/
      step-1-quality-gate.md
    quality_gate_test.go                      (来自 quality_gate_fix_task_loop_breaker)
  feature-management/                         (Journey: Feature 管理)
    contracts/
      step-1-feature-set.md
      step-2-feature-list.md
      step-3-feature-status.md
    feature_set_test.go                       (来自 feature_set_command)
    cli_list_reverse_chronological_test.go    (来自 features/cli-list-reverse-chronological)
    proposal_status_lifecycle_test.go         (来自 features/proposal-status-lifecycle)
  test-generation/                            (Journey: 测试生成管道)
    contracts/
      step-1-task-index.md
      step-2-gen-test-scripts.md
      step-3-run-tests.md
    quick_test_slim_test.go                   (来自 quick_test_slim)
    test_scripts_per_type_test.go             (来自 test_scripts_per_type)
    gen_test_scripts_test.go                  (来自 test-knowledge-convention-driven/gen_test_scripts)
    integration_test.go                       (来自 test-knowledge-convention-driven/integration)
    forge_commands_test.go                    (来自 test-knowledge-convention-driven/forge_commands)
    test_guide_test.go                        (来自 test-knowledge-convention-driven/test_guide)
  command-regression/                         (Journey: 命令回归守卫)
    contracts/
      step-1-removed-commands.md
    removed_commands_test.go                  (来自 test-knowledge-convention-driven/removed_commands)
  task-type-system/                           (Journey: 任务类型系统)
    contracts/
      step-1-list-types.md
      step-2-validate-index.md
      step-3-migrate.md
    task_type_refinement_test.go              (来自 features/task-type-refinement)
  e2e-pipeline/                               (Journey: E2E 测试管道)
    contracts/
      step-1-e2e-setup.md
      step-2-e2e-run.md
    justfile_canonical_test.go                (来自 justfile-canonical-e2e)
  test-suite-health/                          (Journey: 测试套件健康度守卫)
    contracts/
      step-1-test-quality.md
    quality_cleanup_test.go                   (来自 e2e_test_quality_cleanup)
    gen_journeys_skill_test.go                (来自 gen_journeys_skill)
    simplify_test.go                          (来自 simplify_e2e_tests)
```

### Evidence

- 23 个测试文件分散在根目录和 6 个 feature 子目录中，无 `contracts/` 目录
- 测试按技术分类（command 名）而非用户工作流组织
- 无 Contract 规范文件，无法使用 `forge test verify` 契约断裂检测
- 标签使用 `//go:build e2e` 统一标记，未区分 `@feature`/`@regression` 生命周期

## Solution

将现有 E2E 测试从扁平 feature 组织重组为 Journey-Contract 模型，具体工作：

### 1. 提取共享基础设施到 testkit 包

将 `forge_binary.go` 和 `helpers_test.go` 提取为独立的 `testkit` 包，供所有 Journey 包 import。这替换当前根目录的共享基础设施。

- `tests/e2e/testkit/forge_binary.go` — `ForgeBinary` 变量 + `ForgeCmd()` 函数（从 forge_binary.go 移出）
- `tests/e2e/testkit/helpers.go` — `ParseBlock`, `HasField`, `FieldValue`, `WithRetry` 等公开函数（从 helpers_test.go 移出）

### 2. 创建 Journey 目录并迁移测试文件

按用户工作流将现有测试文件分组到 8 个 Journey 目录。每个 Journey 是独立的 Go 包（有自己的 `TestMain` 和 helpers），import testkit 获取共享二进制和工具函数。

> **目录层级说明**：规范定义路径为 `tests/<journey-name>/`，但本提案保持在 `tests/e2e/<journey-name>/` 下。原因：`tests/e2e/` 是一个独立 Go module（`e2e-tests`），与 forge-cli 模块解耦；直接放在 `tests/` 下需要创建新的 Go module 或修改 go.work，增加不必要的复杂度。

映射关系：

| Journey | 现有文件 | 合并策略 |
|---------|---------|---------|
| `task-lifecycle/` | task_lifecycle_hardening, fix-task-claim-priority, task_record_immutability | 合并为 2 个文件：task_lifecycle_test.go + task_record_test.go |
| `quality-gate/` | quality_gate_fix_task_loop_breaker | 保持 1 个文件 |
| `feature-management/` | feature_set_command, cli-list-reverse-chronological, proposal-status-lifecycle | 保持 3 个文件 |
| `test-generation/` | quick_test_slim, test_scripts_per_type, test-knowledge-convention-driven/(gen_test_scripts, integration, forge_commands, test_guide) | 保持 7 个文件（移入并重命名） |
| `task-type-system/` | task-type-refinement | 保持 1 个文件 |
| `e2e-pipeline/` | justfile-canonical-e2e | 保持 1 个文件 |
| `command-regression/` | test-knowledge-convention-driven/removed_commands | 保持 1 个文件 |
| `test-suite-health/` | e2e_test_quality_cleanup, gen_journeys_skill, simplify_e2e_tests | 保持 3 个文件 |

### 3. 为每个 Journey 创建 Contract 规范

在每个 Journey 目录下创建 `contracts/` 目录，从现有测试用例描述中提取 Contract 规范（六维度），遵循 model-and-directory-spec.md 定义的格式。

### 4. 更新构建配置

- 删除根目录的 `main_test.go`、`forge_binary.go`、`helpers_test.go`（已迁移到 testkit）
- 每个 Journey 包添加自己的 `main_test.go`（调用 testkit 设置）
- 更新 justfile 中的 `test-e2e` 命令以适应新目录结构
- 更新 go.work 或模块配置（如果适用）

## Alternatives

| 方案 | 优势 | 劣势 |
|------|------|------|
| **重组到 Journey 目录（当前方案）** | 完全对齐规范，Contract 可用，清晰的工作流组织 | 迁移工作量大，需更新 CI/justfile |
| **仅添加 contracts/ 但不迁移** | 最小改动，不破坏现有结构 | 目录结构仍不符合规范 |
| **渐进迁移** | 分批进行，风险低 | 长期维护两套结构 |

## Scope

### In Scope

- 提取共享基础设施到 testkit 包
- 迁移现有测试文件到 8 个 Journey 目录
- 为每个 Journey 创建 `contracts/` 目录和 Contract 规范文件
- 每个 Journey 包的 `main_test.go` 和 `helpers_test.go`
- 更新 justfile 中的 E2E 测试命令
- 清理旧的根目录文件和 features/ 目录
- 更新 `.graduated/` 标记文件位置（如果适用）
- **去除重复或无效的测试用例**：迁移过程中识别并删除重复测试（相同命令+相同场景的多次覆盖）、已失效测试（测试已移除命令或已变更行为的测试）

### Out of Scope

- 新增测试用例（不增加新覆盖）
- 修改 forge-cli 内部的集成测试（integration_test.go）
- 修改 forge CLI 源代码
- 实现新的 forge test 子命令
- 修改 `forge test promote` 或 `forge test verify` 的实现

## Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 迁移后测试失败（package 名/import 路径错误） | Medium | High | 每个 Journey 迁移后立即运行验证 |
| Go 模块结构不兼容 | Low | High | 所有 Journey 包在同一个 go module 内，仅 package 名不同 |
| contracts 内容不准确 | Low | Low | Contract 从现有测试描述提取，非凭空编写 |
| CI/justfile 配置更新遗漏 | Medium | Medium | 迁移完成后运行完整 CI 验证 |

## Success Criteria

- [ ] 8 个 Journey 目录创建完成，每个包含 `contracts/` 目录
- [ ] 所有 23 个现有测试文件成功迁移到对应 Journey 目录
- [ ] 共享基础设施提取到 testkit 包，所有 Journey 包正确 import
- [ ] `just test-e2e` 通过所有迁移后的测试（去除重复/无效用例后数量可能减少）
- [ ] 旧的根目录测试文件和 features/ 目录已清理
- [ ] 每个 Contract 规范文件包含完整的六维度声明
