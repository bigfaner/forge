---
id: "2"
title: "api/web/mobile service surface 模板"
priority: "P0"
estimated_time: "2h"
complexity: "high"
dependencies: [1]
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.feature"
mainSession: false
---

# 2: api/web/mobile service surface 模板

## Description

在 Task 1 建立的 scaffold 框架上，实现 api、web、mobile 三种 service surface type 的 recipe 模板。与 cli/tui 的关键区别：这些 surface 需要完整的运行时生命周期管理——dev server 启动、健康检查（probe）、测试执行、teardown，以及编排 recipe（`<key>` recipe 将 dev→probe→test→teardown 串联）。

## Reference Files
- `docs/proposals/init-justfile-slim/proposal.md` — 新增：`forge justfile scaffold` CLI 命令 (每个 surface type 生成的 recipes 表、多服务编排模式)
- `plugins/forge/skills/init-justfile/rules/server-lifecycle.md` — PID 管理、idempotent start、健康检查重试的 bash 模板（需迁移到 Go）
- `plugins/forge/skills/init-justfile/rules/surfaces/api.md` — api surface rule 参考（dev→probe→test→teardown 序列）

## Acceptance Criteria
- [ ] `forge justfile scaffold --type api --key backend` 输出的 recipe 集包含：backend-dev + backend-probe + backend-test + backend-teardown + backend（dev→probe→test→teardown 编排）+ quality recipes（compile/fmt/lint/unit-test），所有 recipe 名以 `backend-` 为前缀
- [ ] web surface 与 api 生成相同的 recipe 结构；mobile surface 额外生成 test-setup recipe，编排 recipe 为 test-setup→dev→probe→test→teardown
- [ ] 所有生命周期 recipe 包含 PID 文件管理（启动时写 PID、teardown 时清理）和健康检查重试逻辑（Go 模板形式，带 <<PLACEHOLDER>>）
- [ ] 所有 lifecycle 和 quality recipes 标记 `# user-customized`，包含 `[unix]` + `[windows]` 双平台变体

## Implementation Notes
- api/web 共享 dev→probe→test→teardown 序列，可复用同一模板组；mobile 在此基础上增加 test-setup
- `<key>` 编排 recipe 调用同 surface 的其他 recipe（如 `just backend-dev && just backend-probe && just backend-test && just backend-teardown`）
- 参考 `server-lifecycle.md` 中的 bash 模板结构进行 Go 移植：
  - PID 文件路径：`/tmp/forge-<<URL_KEY>>.pid`（unix）/ `%TEMP%\forge-<<URL_KEY>>.pid`（windows）
  - Idempotent start：检测 PID 文件是否存在且进程存活，避免重复启动
  - 健康检查：循环重试 <<HEALTH_URL>>，间隔 1s，最大重试次数可配置
- 占位符清单（本任务新增涉及）：START_CMD、PORT、HEALTH_URL、SERVICE_LIST
