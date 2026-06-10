---
id: "1"
title: "Scaffold 核心框架 + cli/tui 模板"
priority: "P0"
estimated_time: "2h"
complexity: "high"
dependencies: []
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.feature"
mainSession: false
---

# 1: Scaffold 核心框架 + cli/tui 模板

## Description

实现 `forge justfile scaffold` Cobra 子命令的核心框架：命令注册、参数解析（--type / --key）、参数校验、占位符注入机制、boundary marker 包装。并基于此框架实现 cli 和 tui 两种简单 surface type 的 recipe 模板（仅 test + teardown + quality recipes，无 dev/probe 生命周期）。

cli/tui 是简单路径：不需要服务启动/健康检查/PID 管理，只生成 test + teardown + compile + fmt + lint + unit-test。

## Reference Files
- `docs/proposals/init-justfile-slim/proposal.md` — 新增：`forge justfile scaffold` CLI 命令 (输入参数、占位符清单、Recipe 命名统一模型)
- `plugins/forge/skills/init-justfile/rules/server-lifecycle.md` — 当前 bash 模板参考（teardown 部分的 PID 清理逻辑）
- `plugins/forge/skills/init-justfile/rules/surfaces/cli.md` — cli surface rule 模板参考
- `forge-cli/internal/cmd/` — 现有命令目录结构参考

## Acceptance Criteria
- [ ] `forge justfile scaffold --type cli` 输出包含 test + teardown + compile + fmt + lint + unit-test 的 valid just recipe，所有占位符使用 `<<...>>` 语法
- [ ] 参数校验：unknown surface type 报错退出；scalar surface 传入 `--key` 报错；named surface 未传 `--key` 报错
- [ ] cli 和 tui surface type 生成正确的 recipe 集（test + teardown + quality only，无 dev/probe），tui 与 cli 输出结构一致
- [ ] 所有 lifecycle 和 quality recipes 标记 `# user-customized`；scalar surface（无 --key）生成的 recipe 无前缀
- [ ] 所有 recipe 包含 `[unix]`（Linux + macOS）和 `[windows]` 双平台变体

## Hard Rules
- 占位符语法必须使用 `<<PLACEHOLDER>>` 而非 `{{...}}`，避免与 Go template 和 justfile 变量冲突
- 新建文件位于 `forge-cli/internal/cmd/scaffold/` 目录下

## Implementation Notes
- 建议使用 table-driven pattern：定义 SurfaceSpec 结构体，按 type 索引 recipe 模板集
- cli/tui 共享同一模板结构（test + teardown + quality），可合并为一个模板组
- 先读 `docs/conventions/forge-distribution.md` 了解 Forge 分发模型
- 占位符清单（本任务涉及）：COMPILE_CMD、UNIT_TEST_CMD、LINT_CMD、FMT_CMD、BUILD_CMD、CLEAN_CMD、INSTALL_CMD、TEST_CMD、URL_KEY
