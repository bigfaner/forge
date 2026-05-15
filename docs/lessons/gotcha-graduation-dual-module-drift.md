---
created: "2026-05-15"
tags: [testing, architecture, local-dev-deployment]
---

# test graduation 流向了错误的 Go module：monorepo 双 e2e 模块分裂

## Problem

spec-drift-detection 的 test graduation 执行后，测试文件被放到了 `forge-cli/tests/e2e/spec_drift_detection_cli_test.go`，但项目中**已有的 e2e regression suite** 位于 `tests/e2e/`（项目根目录下的独立 Go module `e2e-tests`）。

具体表现：
- 之前 5 个 feature 的 graduation 全部写入 `tests/e2e/`（repo root，module `e2e-tests`）
- T-quick-4 的 graduation 却写入了 `forge-cli/tests/e2e/`（forge-cli 子目录，module `forge-cli`）
- 同一个 feature `forge-info-commands` 的 staging 测试同时存在于两个 `features/` 目录
- 结果：**两套独立的 e2e 测试基础设施，互不感知，测试运行和发现完全隔离**

## Root Cause

**症状**：graduation 把文件放到了 `forge-cli/tests/e2e/` 而不是 `tests/e2e/`

**直接原因**：gen-test-scripts 和 run-e2e-tests 阶段的测试文件已经位于 `forge-cli/tests/e2e/features/spec-drift-detection/`（因为 T-quick-2/gen 和 T-quick-3/run 阶段就用了 forge-cli 路径），graduate-tests 自然从同一路径取文件并迁移到相邻的 `forge-cli/tests/e2e/`

**根本原因**：整个 test pipeline（gen -> run -> graduate）没有统一的 module root 解析逻辑。三个 skill 各自独立解析路径：

1. **gen-test-scripts** 在 `forge-cli/tests/e2e/features/<slug>/` 生成脚本（可能因为 `go.mod` 搜索找到了 forge-cli 目录）
2. **run-e2e-tests** 在同一 forge-cli 路径运行测试（找到了就用了）
3. **graduate-tests** 从同一路径迁移到 `forge-cli/tests/e2e/`（就近迁移）

而**之前成功的 graduation**（tui-ui-design, justfile-canonical-e2e 等）都是基于 `tests/e2e/`（repo root），因为那些 feature 的 gen 阶段使用了 repo root 的 `e2e-tests` module。

**触发条件**：这个 feature 的测试生成恰好在 forge-cli 子目录下进行，pipeline 从第一步就走偏了，后续步骤没有纠偏机制。

## Solution

**短期修复**：将 `forge-cli/tests/e2e/spec_drift_detection_cli_test.go` 迁移到 `tests/e2e/spec_drift_detection_cli_test.go`，与其他 regression tests 保持一致。

**长期方案**：整个 test pipeline 需要一个统一的路径基准点。建议：

```
1. 在 .forge/config.yaml 或 profile manifest 中声明 e2e module root
2. gen/run/graduate 三个 skill 都读取同一配置
3. graduation 不应 "就近迁移"，而应读取配置后统一迁移到声明的目标
```

## Reusable Pattern

**Pipeline 的一致性依赖于路径基准点的单一声明，而非每个 step 各自推断。**

| 反模式 | 正确做法 |
|--------|----------|
| 每个 skill 独立搜索 `go.mod` 定位 module root | 在 config 中声明一次，所有 skill 读取 |
| gen 在 A 路径生成，run 在 A 路径运行，graduate 从 A 迁移到 A 的邻居 | graduation 目标由配置决定，不受 source 位置影响 |
| 两个独立的 `tests/e2e/` 目录各自演化 | 只维护一套 regression suite |

**判断标准**：当 monorepo 中存在多个 Go module 且都有 `tests/e2e/` 子目录时，e2e regression suite 只能有一个。选择哪个 module 拥有 e2e suite 应在项目初始化时显式声明。

## Example

```
# 当前状态（分裂）
tests/e2e/                                    ← module: e2e-tests (5 个 feature 的 graduation)
  .graduated/tui-ui-design
  .graduated/justfile-canonical-e2e
  tui_ui_design_cli_test.go
  ...
forge-cli/tests/e2e/                          ← module: forge-cli (1 个 feature 的 graduation)
  .graduated/spec-drift-detection             ← 错误位置！
  spec_drift_detection_cli_test.go            ← 应该在上面那个目录

# 应该是
tests/e2e/
  .graduated/tui-ui-design
  .graduated/justfile-canonical-e2e
  .graduated/spec-drift-detection             ← marker 在这里
  tui_ui_design_cli_test.go
  spec_drift_detection_cli_test.go            ← 测试文件在这里
```

## Related Files

- `tests/e2e/go.mod` (repo root e2e module, `module e2e-tests`)
- `forge-cli/go.mod` (forge-cli module, 包含自己的 `tests/e2e/`)
- `forge-cli/tests/e2e/.graduated/spec-drift-detection` (错误位置的 graduation marker)
- `tests/e2e/.graduated/tui-ui-design` (正确位置的 graduation marker 参照)
- `docs/lessons/gotcha-e2e-skill-monorepo-path-mismatch.md` (上一条相关 lesson，覆盖 run-e2e-tests 阶段)

## References

- `docs/lessons/gotcha-e2e-skill-monorepo-path-mismatch.md` -- 同一问题的 run-e2e-tests 阶段表现
- `plugins/forge/skills/graduate-tests/SKILL.md` -- skill 定义，路径假设所在
- `plugins/forge/skills/gen-test-scripts/SKILL.md` -- gen 阶段的路径逻辑
