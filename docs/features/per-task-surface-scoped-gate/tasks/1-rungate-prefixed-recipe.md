---
id: "1"
title: "RunGate() 增加 prefixed recipe 解析"
priority: "P0"
estimated_time: "1.5h"
complexity: "high"
dependencies: []
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 1: RunGate() 增加 prefixed recipe 解析

## Description
当前 per-task quality gate 的 `RunGate()` 通过 `ResolveScope()` 探测参数模式（`compile <scope>`），但 `mixed.just` 的 `compile` recipe 不接受参数 → scope 永远不生效 → 始终执行全量 compile/fmt/lint/unit-test。需要增加 prefixed recipe 解析逻辑：当任务有 surface-key 时优先使用 `<key>-<recipe>`（如 `backend-compile`），找不到则回退到 generic recipe（如 `compile`）。

## Reference Files
- `docs/proposals/per-task-surface-scoped-gate/proposal.md` — Proposed Solution, Key Scenarios, Constraints & Dependencies, Key Risks, Success Criteria
- `forge-cli/pkg/just/just.go`: 增加 resolvePrefixedRecipe() 函数，修改 RunGate() 使用 prefixed resolution 替代 ResolveScope() (ref: Proposed Solution)
- `forge-cli/internal/cmd/qualitygate/quality_gate_lifecycle.go`: resolveRecipe() fallback 模式参考 (ref: Innovation Highlights)

## Acceptance Criteria
- [ ] mixed.just 项目中，`surface-key: backend` 的任务提交时执行 `just backend-compile`/`just backend-lint`（而非 `just compile`/`just lint`），prefixed recipe 不存在时回退到 generic
- [ ] 单 surface 项目中，`RunGate()` 的 scope 为空 → 跳过 prefixed 分支 → recipe 名称仍为 `compile`/`lint`（非 `backend-compile`），退出码和 stdout 与改动前一致
- [ ] `RunGate()` 的 `scope=""` 调用路径（feature-level gate）行为不变
- [ ] Prefixed recipe 执行失败时，`onFail` 回调的 step name 为 `<key>-<recipe>`（如 `backend-lint`），错误信息包含 surface 上下文；回退到 generic recipe 失败时 step name 为原始 recipe 名（如 `lint`），两者可区分
- [ ] 通过 `go test ./internal/cmd/task/...` 和 `go test ./pkg/just/...` 全部测试

## Implementation Notes
- `resolvePrefixedRecipe()` 替代 `ResolveScope()` 在 `RunGate()` 中的作用，调用方无需改动
- 关键守卫条件：`scope != ""` 确保 feature-level gate（传空 scope）不触发 prefixed 解析，保持全量验证行为
- NormalizeSurfaceKey 字符集约束：输出限定为 `[a-z][a-z0-9-]*`，保证 prefixed recipe 名称合法
- 参考 `quality_gate_lifecycle.go` 中 `resolveRecipe()` 的 fallback 模式（surface-specific 优先 → generic 兜底），但注意语义差异：lifecycle 是扇出（fan-out），gate 是选择（select）
- HasRecipe() 探测性能：4 步 gate × 最多 2 次探测 = 8 次 `just --dry-run` fork，Windows 上约 400ms 额外开销；无 surface-key 时跳过整个 prefixed 分支（0 probe）
