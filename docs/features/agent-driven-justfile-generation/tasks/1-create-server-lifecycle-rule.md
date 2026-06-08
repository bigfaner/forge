---
id: "1"
title: "Create server-lifecycle.md rule with PID tracking, idempotent start, and health check patterns"
priority: "P0"
estimated_time: "1.5h"
dependencies: []
type: "doc"
mainSession: false
# Note: surface-key and surface-type fields are intentionally absent from doc tasks.
# Doc tasks produce non-compilable output (markdown, specs, templates) and do not
# interact with the quality gate or test pipeline, so surface routing is unnecessary.
---

# 1: Create server-lifecycle.md rule with PID tracking, idempotent start, and health check patterns

## Description

Create a new rule file `rules/server-lifecycle.md` that extracts server lifecycle bash patterns from the 6 language templates into a single, reusable reference. This rule provides the agent with complete, ready-to-use bash code snippets (with slot placeholders) for server dev recipe generation, replacing the scattered lifecycle code across templates.

The rule must cover: PID file management (path conventions, atomic write, stale PID detection), idempotent start (process-alive check before launch, port-occupancy detection, graceful restart), health check (HTTP/TCP probe with retry and timeout), and multi-service orchestration (per-service PID isolation, port-aware startup order, dependency declaration).

## Reference Files
- `docs/proposals/agent-driven-justfile-generation/proposal.md` — Proposed Solution, Scope > In Scope, Key Risks (ref: ## Proposed Solution, ## Scope > ### In Scope, ## Key Risks)
- `plugins/forge/skills/init-justfile/SKILL.md` — current Step 3 surface recipe generation patterns (ref: ### Step 3: Generate Recipes and Assemble Justfile)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/init-justfile/rules/server-lifecycle.md` | Server lifecycle bash patterns for dev/probe/teardown recipes |

### Modify
| File | Changes |
|------|---------|

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] `plugins/forge/skills/init-justfile/rules/server-lifecycle.md` 文件已创建
- [ ] 包含 PID 追踪模式：PID 文件路径约定（`.forge/<surfaceKey>.pid`）、原子写入、stale PID 检测（进程不存在时的清理逻辑）
- [ ] 包含幂等启动模式：检测已有进程（避免重复启动）、端口占用检查（检测后选择备选端口或报错）、启动/重启逻辑
- [ ] 包含健康检查模式：HTTP/TCP probe 实现、重试策略（最多 3 次、5 秒间隔）、超时处理
- [ ] 包含 multi-service 场景指导：per-service PID 文件隔离、端口感知启动顺序、启动顺序依赖声明
- [ ] 提供可直接使用的 bash 代码片段（带插槽占位符如 `<PORT>`、`<START_CMD>`），agent 优先复用而非从头生成

## Implementation Notes

- 从现有 6 个模板中提取通用 server lifecycle 模式作为参考基线（特别是 `mixed.just` 中最复杂的场景）
- 代码片段必须同时支持 `[linux]` 和 `[windows]` 两种平台变体
- PID 文件路径约定应与 Forge 的 `.forge/` 工作目录对齐
- 已知限制（来自 proposal Key Risks）：verification step 无法覆盖 server lifecycle 边界条件（PID 文件残留、进程被外部 kill 后 PID 被回收、Windows 上 `\r` 污染），代码片段应包含防御性处理
