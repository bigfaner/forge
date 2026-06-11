---
created: "2026-05-15"
tags: [testing, architecture, skill-design]
---

# E2E Skill 假定扁平项目布局，在 monorepo 下测试路径全面失配

## Problem

执行 `run-e2e-tests` skill 运行 spec-drift-detection 的 e2e 测试时，skill 的三个核心前提检查全部失败，耗费大量时间在"文件不存在"的死胡同中排查：

1. **Justfile 缺少 `e2e-setup` recipe** -- 根 Justfile 无任何 e2e 相关 target
2. **测试目录 `tests/e2e/features/<slug>/` 不存在** -- 实际在 `forge-cli/tests/e2e/features/<slug>/`
3. **Glob 找到文件但 ls 找不到** -- glob 匹配了 git index，实际文件在 monorepo 子目录下

最终发现测试脚本位于 `forge-cli/tests/e2e/features/spec-drift-detection/`（Go module 子目录），而非 skill 预期的项目根 `tests/e2e/features/`。直接用 `go test` 命令绕过了整个 Justfile 流程才完成测试。

## Root Cause

**症状**：三个前提检查全部报"不存在"，但文件确实存在

**直接原因**：skill 的路径假设是 `{project_root}/tests/e2e/features/<slug>/`，但 forge 项目是 monorepo 结构，Go module 位于 `forge-cli/` 子目录

**根本原因**：skill 设计时的隐含假设 -- "项目根 = Go module 根" -- 从未被显式声明或验证。当 forge 自身作为被测项目时，这个假设不成立：

1. **Justfile 假设**：skill 要求 `Justfile` 含 `e2e-setup` recipe，但 forge 的根 Justfile 只有 `compile/build/test/lint` 等标准 recipe，e2e 测试由 `forge-cli` 内部的 Go test 框架管理
2. **路径假设**：`forge feature` 返回 `spec-drift-detection`，skill 在 `{root}/tests/e2e/features/spec-drift-detection/` 查找，但文件在 `{root}/forge-cli/tests/e2e/features/spec-drift-detection/`
3. **profile manifest 假设**：go-test profile 的 `test-directory: tests/e2e/` 是相对 Go module 根的路径，但 skill 解析为相对项目根

## Solution

**短期**：识别 Go module 根（`go.mod` 所在目录），以此为基础路径发现测试文件

**长期**：profile manifest 应增加 `module-root` 字段，skill 在路径解析时先定位 module root 再拼接 `test-directory`

```yaml
# profile manifest 建议增加
module-root: forge-cli/   # 相对项目根，空表示项目根即 module 根
test-directory: tests/e2e/ # 相对 module-root
```

## Reusable Pattern

**当 skill 依赖文件系统路径时，路径基准点必须显式定义，不能隐含假设。**

| 路径基准 | 适用场景 | 风险 |
|----------|----------|------|
| 项目根（git root） | 通用文件、文档 | monorepo 子模块路径失配 |
| Module 根（go.mod/Cargo.toml 所在） | 测试、构建产物 | 文档或跨模块引用路径失配 |
| Skill 工作目录 | 临时产物 | 不同调用场景下 cwd 不同 |

**判断标准**：如果 skill 涉及编译、测试、构建等语言/工具链操作，路径基准应该是 **module 根**（语言工具链的视角），而非项目根（版本控制的视角）。两者在扁平项目中相同，在 monorepo 中不同。

**反模式**：skill 写 `tests/e2e/features/<slug>/` 而不说明这个路径相对于哪里 -- 隐含 "相对于项目根" 但从未声明。

## Example

```
# skill 当前假设
tests/e2e/features/spec-drift-detection/  ← 从项目根查找，找不到

# 实际布局
forge-cli/                          ← Go module root (go.mod)
  tests/e2e/features/spec-drift-detection/  ← 测试实际位置

# 修正后的查找逻辑
1. 定位 module root: find . -name go.mod -maxdepth 2
2. 拼接 test-directory: {module_root}/tests/e2e/features/<slug>/
3. Justfile 也应在 module root 查找（或 skill 应接受"无 Justfile，直接用工具链"的 fallback）
```

## Related Files

- `plugins/forge/skills/run-e2e-tests/SKILL.md` (skill 定义，路径假设所在)
- `plugins/forge/skills/run-e2e-tests/templates/e2e-report.md` (报告模板)
- `forge-cli/tests/e2e/features/spec-drift-detection/spec_drift_detection_cli_test.go` (实际测试文件)
- `docs/features/spec-drift-detection/tasks/quick-run-tests-go-test.md` (T-quick-3 任务定义)
