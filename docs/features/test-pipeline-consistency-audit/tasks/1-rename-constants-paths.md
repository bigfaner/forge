---
id: "1"
title: "重命名 E2E 常量为 surface-neutral 名称并扁平化 tests/ 目录"
priority: "P0"
estimated_time: "2h"
dependencies: []
surface-key: "."
surface-type: "cli"
breaking: true
type: "coding.refactor"
mainSession: false
---

# 1: 重命名 E2E 常量为 surface-neutral 名称并扁平化 tests/ 目录

## Description
将 `pkg/feature/constants.go` 和 `pkg/feature/paths.go` 中所有 `E2E*` 前缀常量和函数重命名为 surface-neutral 名称，删除 staging/graduation 相关路径函数，并将物理目录从 `tests/e2e/` 扁平化为 `tests/`。同时确认 `testrunner` 路径一致性（MR #173 已修复 `WriteRegressionRawOutput`，需验证 `WriteUnitTestRawOutput` 路径）。

## Reference Files
- `proposal.md#Layer-1-Go-代码层术语路径统一` — 定义了常量重命名映射和路径扁平化规则
- `proposal.md#新目录结构` — 物理目录扁平化的具体前后对比
- `proposal.md#Success-Criteria` — grep 验证条件（tests/e2e、E2EStaging 等）

## Acceptance Criteria
- [ ] `E2ETestsBaseDir` 重命名为 `TestBaseDir`，值为 `"tests"`
- [ ] `E2EStagingDir` / `E2EGraduatedDir` 常量已删除
- [ ] `GetE2EStagingDir()` / `GetE2EGraduatedMarker()` / `GetE2ETargetDir()` 函数已删除
- [ ] 新增 `GetTestResultsDir()` 和 `GetTestConfigPath()` 函数
- [ ] `tests/e2e/config.yaml` 移动到 `tests/config.yaml`
- [ ] `tests/e2e/results/` 移动到 `tests/results/`
- [ ] `tests/e2e/features/` 和 `tests/e2e/.graduated/` 目录已删除
- [ ] `tests/e2e/` 目录不再存在
- [ ] `grep -rn "tests/e2e" forge-cli/pkg/ forge-cli/internal/`（排除 _test.go）返回 0 结果
- [ ] `WriteUnitTestRawOutput` 路径与 `WriteRegressionRawOutput` 一致，均使用新路径
- [ ] `go build ./...` 通过

## Implementation Notes
- MR #173 已修复 `init.go` 和 `testrunner` 的 `WriteRegressionRawOutput` 路径，不要重复修改
- 重命名常量后需同步更新所有引用点（`quality_gate.go`、`autogen.go`、`testrunner` 等）
- 物理目录移动后确认 `tests/e2e/` 下无其他遗漏文件

### Integration Test Impact
- Affected test suite(s): `forge-cli/pkg/feature/`, `forge-cli/pkg/testrunner/`
- Expected fixture changes: 测试中引用旧常量名的断言需同步更新
- Risk level: medium
