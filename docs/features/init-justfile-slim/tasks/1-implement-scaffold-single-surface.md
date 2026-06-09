---
id: "1"
title: "实现 scaffold 单 surface recipe 生成"
priority: "P0"
estimated_time: "3h"
complexity: "high"
dependencies: []
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.feature"
mainSession: false
---

# 1: 实现 scaffold 单 surface recipe 生成

## Description

实现 `forge justfile scaffold --type <type> [--key <key>]` CLI 子命令的核心逻辑：为单个 surface type 生成完整的 just recipe 集（lifecycle + quality），输出带 `<<PLACEHOLDER>>` 占位符的 just 代码到 stdout。

这是提案的核心机制，将 recipe 代码生成从 prompt 层下沉到 Forge CLI。CLI 内置 5 种 surface type（cli/tui/api/web/mobile）的差异化模板逻辑，取代当前分散在 `server-lifecycle.md` 和 5 个 surface rule 文件中的 bash 模板。

## Reference Files
- `docs/proposals/init-justfile-slim/proposal.md` — 新增：`forge justfile scaffold` CLI 命令, 核心思路, Recipe 命名统一模型, 风险与缓解
- `plugins/forge/skills/init-justfile/rules/server-lifecycle.md` — 当前 bash 模板（745 行），需迁移到 Go 代码
- `plugins/forge/skills/init-justfile/rules/surfaces/cli.md` — cli surface rule，模板参考
- `forge-cli/pkg/just/just.go` — just 包，scaffold 命令可能需要复用 HasRecipe 等

## Acceptance Criteria
- [ ] `forge justfile scaffold --type cli` 输出包含 test + teardown + compile + fmt + lint + unit-test 的 valid just recipe，所有占位符使用 `<<...>>` 语法
- [ ] `forge justfile scaffold --type api --key backend` 输出的所有 recipe 名以 `backend-` 为前缀（如 `backend-dev`、`backend-test`），且包含 dev + probe + test + teardown + `backend`（dev→probe→test→teardown 编排）+ quality recipes
- [ ] 5 种 surface type（cli/tui/api/web/mobile）均按提案 Recipe 表生成正确的 recipe 集：cli/tui 无 dev/probe，api/web 有 dev/probe + `<key>` 编排 recipe，mobile 额外有 test-setup
- [ ] 所有 lifecycle 和 quality recipes 标记 `# user-customized`；scalar surface（无 --key）生成的 recipe 无前缀
- [ ] 参数校验：unknown surface type 报错；scalar surface 传入 `--key` 报错；named surface 未传 `--key` 报错
- [ ] 所有 recipe 包含 `[unix]`（Linux + macOS）和 `[windows]` 双平台变体

## Hard Rules
- 占位符语法必须使用 `<<PLACEHOLDER>>` 而非 `{{...}}`，避免与 Go template 和 justfile 变量冲突
- 新建文件位于 `forge-cli/internal/cmd/scaffold/` 目录下

## Implementation Notes
- 参考 `server-lifecycle.md` 中的 bash 模板结构（PID 文件管理、idempotent start、健康检查重试）进行 Go 移植
- 提案占位符清单共 13 个：START_CMD、PORT、HEALTH_URL、COMPILE_CMD、UNIT_TEST_CMD、LINT_CMD、FMT_CMD、BUILD_CMD、CLEAN_CMD、INSTALL_CMD、TEST_CMD、URL_KEY、SERVICE_LIST
- 先读 `docs/conventions/forge-distribution.md` 了解 Forge 分发模型，确保新命令在用户项目目录下正确解析路径
- 建议用 table-driven pattern 按 surface type 映射 recipe 模板，减少重复
