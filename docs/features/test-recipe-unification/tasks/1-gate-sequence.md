---
id: "1"
title: "Restructure gate sequences for two-layer test model"
priority: "P0"
estimated_time: "1-2h"
dependencies: []
scope: "backend"
breaking: true
type: "coding.enhancement"
mainSession: false
---

# 1: Restructure gate sequences for two-layer test model

## Description

引入两层测试 recipe 模型所需的核心 gate sequence 重构。将 `DefaultGateSequence()` 重命名为 `FullGateSequence()` 并修改其 step，同时新增 `UnitGateSequence()` 供 breaking 任务 submit 使用。更新 `submit.go` 让 breaking 任务使用新的 `UnitGateSequence`。

## Reference Files
- `proposal.md#Proposed-Solution` — defines two-layer test recipe model and gate sequence structure (UnitGateSequence, NonBreakingGateSequence, FullGateSequence)
- `proposal.md#Requirements-Analysis` — Key Scenarios for breaking/non-breaking submit and all-completed gate timing
- `proposal.md#Feasibility-Assessment` — Gate sequence migration table with exact step composition and callers
- `proposal.md#Key-Risks` — risk of `just test` semantic inversion and mitigation via template comments

## Acceptance Criteria
- `DefaultGateSequence()` renamed to `FullGateSequence()` with steps: `compile → fmt → lint → unit-test → test → probe`
- `UnitGateSequence()` added with steps: `compile → fmt → lint → unit-test`
- `NonBreakingGateSequence()` remains unchanged: `compile → fmt → lint`
- `submit.go` uses `UnitGateSequence` for breaking tasks, `NonBreakingGateSequence` for non-breaking
- No `DefaultGateSequence` symbol remains in the codebase

## Hard Rules
- Gate sequence 中无 fallback——如果 recipe 不存在，gate 报错提示运行 `init-justfile`
- `UnitGateSequence` 的 `unit-test` step 不回落到 `test`（两者独立）

## Implementation Notes
- `FullGateSequence` 命名消除 "Default" 歧义——迁移后每个调用方显式选择对应 sequence
- 确保 `submit.go` 中 breaking 判断逻辑正确选择 `UnitGateSequence` vs `NonBreakingGateSequence`
