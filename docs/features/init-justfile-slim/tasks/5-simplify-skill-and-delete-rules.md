---
id: "5"
title: "精简 SKILL.md 并删除废弃 rule 文件"
priority: "P0"
estimated_time: "1h"
complexity: "medium"
dependencies: [1]
surface-key: ""
surface-type: "cli"
breaking: false
type: "doc"
mainSession: false
---

# 5: 精简 SKILL.md 并删除废弃 rule 文件

## Description

将 init-justfile skill 的 SKILL.md 从 548 行精简到 ~250 行（加上 self-correction.md 34 行，prompt 层总计 ≤ 280 行），删除 6 个已被 CLI scaffold 取代的 rule 文件，更新 agent 工作流步骤。

SKILL.md 保留的职责：流程编排、语言检测、Convention 加载、占位符填值、验证、输出确认。删除的内容：Step 1d Load Server Lifecycle Patterns、Step 3b Surface recipe 生成细节、Surface rule 加载逻辑、Phase 1 Consistency Verification、Surface Gate Targets 段、EXTREMELY-IMPORTANT 重复项、Notes 重复。

## Reference Files
- `docs/proposals/init-justfile-slim/proposal.md` — 精简：SKILL.md, 删除：surface rule 文件 + server-lifecycle.md, Agent 新流程, 保留, Consumer Impact
- `plugins/forge/skills/init-justfile/SKILL.md` — 当前 SKILL.md（548 行），需重写
- `plugins/forge/skills/init-justfile/rules/server-lifecycle.md` — 删除（745 行）
- `plugins/forge/skills/init-justfile/rules/surfaces/` — 删除目录（api.md/cli.md/mobile.md/tui.md/web.md）

## Acceptance Criteria
- [ ] SKILL.md + 保留的 rules（self-correction.md）总行数 ≤ 280 行
- [ ] SKILL.md 引用 `forge justfile scaffold` CLI 命令替代手动模板生成，Agent 工作流更新为提案的 Step 0-5 新流程
- [ ] 删除 6 个文件：`rules/server-lifecycle.md`、`rules/surfaces/api.md`、`rules/surfaces/cli.md`、`rules/surfaces/mobile.md`、`rules/surfaces/tui.md`、`rules/surfaces/web.md`
- [ ] `rules/self-correction.md`（34 行）保留不动
- [ ] Convention Cold Start Fallback 策略以 5-10 行摘要保留在 SKILL.md 中

## Hard Rules
- 修改 `plugins/forge/` 下的文件前，必须先读 `docs/conventions/forge-distribution.md` 了解分发模型

## Implementation Notes
- 删除文件清单（精确路径）：
  - `plugins/forge/skills/init-justfile/rules/server-lifecycle.md`
  - `plugins/forge/skills/init-justfile/rules/surfaces/api.md`
  - `plugins/forge/skills/init-justfile/rules/surfaces/cli.md`
  - `plugins/forge/skills/init-justfile/rules/surfaces/mobile.md`
  - `plugins/forge/skills/init-justfile/rules/surfaces/tui.md`
  - `plugins/forge/skills/init-justfile/rules/surfaces/web.md`
- Phase 1 Consistency Verification 简化为轻量级 `just --list` 验证
- `rules/surfaces/` 目录删除后若为空则一并删除
