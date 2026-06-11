---
status: "completed"
started: "2026-06-08 22:07"
completed: "2026-06-08 22:13"
time_spent: "~6m"
---

# Task Record: 1 Create server-lifecycle.md rule with PID tracking, idempotent start, and health check patterns

## Summary
Created rules/server-lifecycle.md with PID tracking, idempotent start, health check, and multi-service orchestration patterns extracted from 6 language templates

## Changes

### Files Created
- plugins/forge/skills/init-justfile/rules/server-lifecycle.md

### Files Modified
无

### Key Decisions
无

## Document Metrics
5 sections, ~320 lines, 4 complete recipe snippets with [linux]/[windows] variants, 13 slot placeholders

## Referenced Documents
- docs/proposals/agent-driven-justfile-generation/proposal.md
- plugins/forge/skills/init-justfile/SKILL.md
- plugins/forge/skills/init-justfile/templates/mixed.just
- plugins/forge/skills/init-justfile/templates/go.just
- plugins/forge/skills/init-justfile/templates/node.just
- plugins/forge/skills/init-justfile/templates/python.just
- plugins/forge/skills/init-justfile/templates/rust.just
- plugins/forge/skills/init-justfile/templates/generic.just

## Review Status
final

## Acceptance Criteria
- [x] server-lifecycle.md 文件已创建
- [x] PID 追踪模式: 路径约定 .forge/<surfaceKey>.pid, 原子写入, stale PID 检测
- [x] 幂等启动模式: 三层检测(进程/端口/启动), 端口占用+备选端口, 启动/重启逻辑
- [x] 健康检查模式: HTTP/TCP probe, 3次重试/5秒间隔, 超时处理
- [x] multi-service 场景: per-service PID 隔离, 端口感知启动顺序, 依赖声明
- [x] 可直接使用的 bash 代码片段(带插槽占位符), [linux]+[windows] 双平台

## Notes
PID path migrated from templates' tests/results/.pid-* to spec's .forge/<surfaceKey>.pid. Health check retry adjusted from templates' 10x/3s to spec's 3x/5s. Defensive measures for \r contamination, stale PID, PID recycling, and early crash detection are documented in a dedicated section.
